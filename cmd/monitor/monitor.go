package monitor

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/maticnetwork/polygon-cli/util"

	_ "embed"

	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"

	"github.com/cenkalti/backoff/v4"
	termui "github.com/gizak/termui/v3"
	"github.com/maticnetwork/polygon-cli/cmd/monitor/ui"
	"github.com/maticnetwork/polygon-cli/metrics"
	"github.com/maticnetwork/polygon-cli/rpctypes"
	"github.com/rs/zerolog/log"
)

var errBatchRequestsNotSupported = errors.New("batch requests are not supported")

var (
	windowSize         int
	batchSize          int
	interval           time.Duration
	one                = big.NewInt(1)
	zero               = big.NewInt(0)
	observedPendingTxs historicalRange
	maxDataPoints      = 1000
)

type (
	monitorStatus struct {
		TopDisplayedBlock *big.Int
		UpperBlock        *big.Int
		LowerBlock        *big.Int
		ChainID           *big.Int
		HeadBlock         *big.Int
		PeerCount         uint64
		GasPrice          *big.Int
		PendingCount      uint64
		SelectedBlock     rpctypes.PolyBlock
		BlockCache        *lru.Cache   `json:"-"`
		BlocksLock        sync.RWMutex `json:"-"`
	}
	chainState struct {
		HeadBlock    uint64
		ChainID      *big.Int
		PeerCount    uint64
		GasPrice     *big.Int
		PendingCount uint64
	}
	historicalDataPoint struct {
		SampleTime  time.Time
		SampleValue float64
	}
	historicalRange []historicalDataPoint
	monitorMode     int
)

const (
	monitorModeHelp monitorMode = iota
	monitorModeExplorer
	monitorModeBlock
	monitorModeTransaction
)

func monitor(ctx context.Context) error {
	// Dial rpc.
	rpc, err := ethrpc.DialContext(ctx, rpcUrl)
	if err != nil {
		log.Error().Err(err).Msg("Unable to dial rpc")
		return err
	}
	ec := ethclient.NewClient(rpc)
	if _, err = ec.BlockNumber(ctx); err != nil {
		return err
	}

	// Check if batch requests are supported.
	if err = checkBatchRequestsSupport(ctx, ec.Client()); err != nil {
		return errBatchRequestsNotSupported
	}

	ms := new(monitorStatus)
	ms.BlocksLock.Lock()
	ms.BlockCache, err = lru.New(blockCacheLimit)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create new LRU cache")
		return err
	}
	ms.BlocksLock.Unlock()

	ms.ChainID = big.NewInt(0)
	ms.PendingCount = 0

	observedPendingTxs = make(historicalRange, 0)

	isUiRendered := false
	errChan := make(chan error)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error().Msg(fmt.Sprintf("Recovered in f: %v", r))
			}
		}()
		select {
		case <-ctx.Done(): // listens for a cancellation signal
			return // exit the goroutine when the context is done
		default:
			for {
				err = fetchCurrentBlockData(ctx, ec, ms, rpc, isUiRendered)
				if err != nil {
					continue
				}
				if ms.TopDisplayedBlock == nil || ms.SelectedBlock == nil {
					ms.TopDisplayedBlock = ms.HeadBlock
					// from := new(big.Int).Sub(ms.HeadBlock, big.NewInt(int64(batchSize-1)))

					// if from.Cmp(zero) < 0 {
					// 	from.SetInt64(0)
					// }

					// err = ms.getBlockRange(ctx, from, ms.HeadBlock, rpc)
					// if err != nil {
					// 	continue
					// }
				}
				if !isUiRendered {
					go func() {
						// if ms.TopDisplayedBlock == nil {
						// 	ms.TopDisplayedBlock = ms.HeadBlock
						// }
						errChan <- renderMonitorUI(ctx, ec, ms, rpc)
					}()
					isUiRendered = true
				}

				time.Sleep(interval)
			}
		}
	}()

	err = <-errChan
	return err
}

func getChainState(ctx context.Context, ec *ethclient.Client) (*chainState, error) {
	var err error
	cs := new(chainState)
	cs.HeadBlock, err = ec.BlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("couldn't fetch block number: %s", err.Error())
	}

	cs.ChainID, err = ec.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("couldn't fetch chain id: %s", err.Error())
	}

	cs.PeerCount, err = ec.PeerCount(ctx)
	if err != nil {
		log.Debug().Err(err).Msg("Using fake peer count")
		cs.PeerCount = 0
	}

	cs.GasPrice, err = ec.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("couldn't estimate gas: %s", err.Error())
	}

	cs.PendingCount, err = util.GetTxPoolSize(ec.Client())
	if err != nil {
		log.Debug().Err(err).Msg("Unable to get pending transaction count")
		cs.PendingCount = 0
	}

	return cs, nil

}

func (h historicalRange) getValues(limit int) []float64 {
	values := make([]float64, len(h))
	for idx, v := range h {
		values[idx] = v.SampleValue
	}
	if limit < len(values) {
		values = values[len(values)-limit:]
	}
	return values
}

func fetchCurrentBlockData(ctx context.Context, ec *ethclient.Client, ms *monitorStatus, rpc *ethrpc.Client, isUiRendered bool) (err error) {
	var cs *chainState
	cs, err = getChainState(ctx, ec)
	if err != nil {
		log.Error().Err(err).Msg("Encountered issue fetching network information")
		time.Sleep(interval)
		return err
	}
	observedPendingTxs = append(observedPendingTxs, historicalDataPoint{SampleTime: time.Now(), SampleValue: float64(cs.PendingCount)})
	if len(observedPendingTxs) > maxDataPoints {
		observedPendingTxs = observedPendingTxs[len(observedPendingTxs)-maxDataPoints:]
	}

	log.Debug().Uint64("PeerCount", cs.PeerCount).Uint64("ChainID", cs.ChainID.Uint64()).Uint64("HeadBlock", cs.HeadBlock).Uint64("GasPrice", cs.GasPrice.Uint64()).Msg("Fetching blocks")

	if isUiRendered && batchSize < 0 {
		_, termHeight := termui.TerminalDimensions()
		batchSize = termHeight/2 - 4
	}

	ms.HeadBlock = new(big.Int).SetUint64(cs.HeadBlock)
	ms.ChainID = cs.ChainID
	ms.PeerCount = cs.PeerCount
	ms.GasPrice = cs.GasPrice
	ms.PendingCount = cs.PendingCount

	return
}

func (ms *monitorStatus) getBlockRange(ctx context.Context, from, to *big.Int, rpc *ethrpc.Client) error {
	blms := make([]ethrpc.BatchElem, 0)

	ms.UpperBlock = to
	ms.LowerBlock = from
	log.Debug().Msgf("FETCH BLOCKS %d %d", ms.UpperBlock, ms.LowerBlock)
	for i := new(big.Int).Set(from); i.Cmp(to) <= 0; i.Add(i, one) {
		ms.BlocksLock.RLock()
		_, found := ms.BlockCache.Get(i.String())
		ms.BlocksLock.RUnlock()
		if found {
			continue
		}
		r := new(rpctypes.RawBlockResponse)
		blms = append(blms, ethrpc.BatchElem{
			Method: "eth_getBlockByNumber",
			Args:   []interface{}{"0x" + i.Text(16), true},
			Result: r,
			Error:  nil,
		})
	}

	if len(blms) == 0 {
		return nil
	}

	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 3 * time.Minute
	retryable := func() error {
		return rpc.BatchCallContext(ctx, blms)
	}
	if err := backoff.Retry(retryable, b); err != nil {
		return err
	}

	ms.BlocksLock.Lock()
	defer ms.BlocksLock.Unlock()
	for _, b := range blms {
		if b.Error != nil {
			continue
		}
		pb := rpctypes.NewPolyBlock(b.Result.(*rpctypes.RawBlockResponse))
		ms.BlockCache.Add(pb.Number().String(), pb)
	}

	return nil
}

func renderMonitorUI(ctx context.Context, ec *ethclient.Client, ms *monitorStatus, rpc *ethrpc.Client) error {
	if err := termui.Init(); err != nil {
		log.Error().Err(err).Msg("Failed to initialize UI")
		return err
	}
	defer termui.Close()

	currentMode := monitorModeExplorer

	blockTable, blockInfo, transactionInfo, grid, blockGrid, skeleton := ui.SetUISkeleton()

	termWidth, termHeight := termui.TerminalDimensions()
	windowSize = termHeight/2 - 4
	grid.SetRect(0, 0, termWidth, termHeight)
	blockGrid.SetRect(0, 0, termWidth, termHeight)

	var setBlock = false
	var renderedBlocks rpctypes.SortableBlocks

	redraw := func(ms *monitorStatus, force ...bool) {
		log.Debug().
			Str("TopDisplayedBlock", ms.TopDisplayedBlock.String()).
			Str("UpperBlock", ms.UpperBlock.String()).
			Str("LowerBlock", ms.LowerBlock.String()).
			Str("ChainID", ms.ChainID.String()).
			Str("HeadBlock", ms.HeadBlock.String()).
			Uint64("PeerCount", ms.PeerCount).
			Str("GasPrice", ms.GasPrice.String()).
			Uint64("PendingCount", ms.PendingCount).
			Msg("Redrawing")

		if currentMode == monitorModeHelp {
			// TODO add some help context?
		} else if currentMode == monitorModeBlock {
			// render a block
			skeleton.BlockInfo.Rows = ui.GetSimpleBlockFields(ms.SelectedBlock)
			rows, title := ui.GetTransactionsList(ms.SelectedBlock, ms.ChainID)
			skeleton.TransactionList.Rows = rows
			skeleton.TransactionList.Title = title

			termui.Clear()
			termui.Render(blockGrid)
			return
		}

		if blockTable.SelectedRow == 0 {
			// ms.TopDisplayedBlock = ms.HeadBlock
			// log.Debug().
			// 	Str("LowerBlock", ms.LowerBlock.String()).
			// 	Str("UpperBlock", ms.UpperBlock.String()).
			// 	Str("ms.HeadBlock", ms.HeadBlock.String()).
			// 	Msg("TEST TEST")

			bottomBlockNumber := new(big.Int).Sub(ms.HeadBlock, big.NewInt(int64(windowSize-1)))
			if bottomBlockNumber.Cmp(zero) < 0 {
				bottomBlockNumber.SetInt64(0)
			}

			// if ms.LowerBlock == nil || ms.LowerBlock.Cmp(bottomBlockNumber) > 0 {
			err := ms.getBlockRange(ctx, bottomBlockNumber, ms.TopDisplayedBlock, rpc)
			if err != nil {
				log.Error().Err(err).Msg("There was an issue fetching the block range")
			}
			// }
		}
		toBlockNumber := ms.TopDisplayedBlock
		fromBlockNumber := new(big.Int).Sub(toBlockNumber, big.NewInt(int64(windowSize-1)))
		if fromBlockNumber.Cmp(zero) < 0 {
			fromBlockNumber.SetInt64(0) // We cannot have block numbers less than 0.
		}
		renderedBlocksTemp := make([]rpctypes.PolyBlock, 0, windowSize)
		ms.BlocksLock.RLock()
		for i := new(big.Int).Set(fromBlockNumber); i.Cmp(toBlockNumber) <= 0; i.Add(i, big.NewInt(1)) {
			if block, ok := ms.BlockCache.Get(i.String()); ok {
				renderedBlocksTemp = append(renderedBlocksTemp, block.(rpctypes.PolyBlock))
			} else {
				// If for some reason the block is not in the cache after fetching, handle this case.
				log.Warn().Str("blockNumber", i.String()).Msg("Block should be in cache but is not")
			}
		}
		ms.BlocksLock.RUnlock()
		renderedBlocks = renderedBlocksTemp

		skeleton.Current.Text = ui.GetCurrentBlockInfo(ms.HeadBlock, ms.GasPrice, ms.PeerCount, ms.PendingCount, ms.ChainID, renderedBlocks)
		skeleton.TxPerBlockChart.Data = metrics.GetTxsPerBlock(renderedBlocks)
		skeleton.GasPriceChart.Data = metrics.GetMeanGasPricePerBlock(renderedBlocks)
		skeleton.BlockSizeChart.Data = metrics.GetSizePerBlock(renderedBlocks)
		// skeleton.pendingTxChart.Data = metrics.GetUnclesPerBlock(renderedBlocks)
		skeleton.PendingTxChart.Data = observedPendingTxs.getValues(25)
		skeleton.GasChart.Data = metrics.GetGasPerBlock(renderedBlocks)

		// If a row has not been selected, continue to update the list with new blocks.
		rows, title := ui.GetBlocksList(renderedBlocks)
		blockTable.Rows = rows
		blockTable.Title = title

		blockTable.TextStyle = termui.NewStyle(termui.ColorWhite)
		blockTable.SelectedRowStyle = termui.NewStyle(termui.ColorWhite, termui.ColorRed, termui.ModifierBold)
		if blockTable.SelectedRow > 0 && blockTable.SelectedRow <= len(blockTable.Rows) {
			// Only changed the selected block when the user presses the up down keys.
			// Otherwise this will adjust when the table is updated automatically.
			if setBlock {
				log.Debug().
					Int("blockTable.SelectedRow", blockTable.SelectedRow).
					Int("renderedBlocks", len(renderedBlocks)).
					Msg("setBlock")

				ms.SelectedBlock = renderedBlocks[len(renderedBlocks)-blockTable.SelectedRow]
				blockInfo.Rows = ui.GetSimpleBlockFields(ms.SelectedBlock)
				columnRatio := []int{30, 5, 5, 20, 20, 5, 10}
				transactionInfo.ColumnWidths = getColumnWidths(columnRatio, transactionInfo.Dx())
				transactionInfo.Rows = ui.GetBlockTxTable(ms.SelectedBlock, ms.ChainID)
				transactionInfo.Title = fmt.Sprintf("Latest Transactions for Block #%s", ms.SelectedBlock.Number().String())

				setBlock = false
				log.Debug().Uint64("blockNumber", ms.SelectedBlock.Number().Uint64()).Msg("Selected block changed")
			}
		} else {
			ms.SelectedBlock = nil
			transactionInfo.Title = "Latest Transactions"
			blockInfo.Rows = []string{}
			transactionInfo.Rows = [][]string{{""}, {""}}
		}

		termui.Render(grid)
	}

	currentBn := ms.HeadBlock
	uiEvents := termui.PollEvents()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	redraw(ms)

	previousKey := ""
	for {
		forceRedraw := false
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return nil
			case "<Escape>":
				ms.TopDisplayedBlock = ms.HeadBlock
				blockTable.SelectedRow = 0
				currentMode = monitorModeExplorer

				toBlockNumber := new(big.Int).Sub(ms.TopDisplayedBlock, big.NewInt(int64(windowSize-1)))
				if toBlockNumber.Cmp(zero) < 0 {
					toBlockNumber.SetInt64(0)
				}

				err := ms.getBlockRange(ctx, toBlockNumber, ms.TopDisplayedBlock, rpc)
				if err != nil {
					log.Error().Err(err).Msg("There was an issue fetching the block range")
					break
				}
			case "<Enter>":
				if blockTable.SelectedRow > 0 {
					currentMode = monitorModeBlock
				}
			case "<Resize>":
				payload := e.Payload.(termui.Resize)
				grid.SetRect(0, 0, payload.Width, payload.Height)
				blockGrid.SetRect(0, 0, payload.Width, payload.Height)
				_, termHeight = termui.TerminalDimensions()
				windowSize = termHeight/2 - 4
				termui.Clear()
			case "<Up>", "<Down>":
				if currentMode == monitorModeBlock {
					if len(skeleton.TransactionList.Rows) != 0 && e.ID == "<Down>" {
						skeleton.TransactionList.ScrollDown()
					} else if len(skeleton.TransactionList.Rows) != 0 && e.ID == "<Up>" {
						skeleton.TransactionList.ScrollUp()
					}
					break
				}

				if blockTable.SelectedRow == 0 {
					blockTable.SelectedRow = 1
					setBlock = true
					break
				}

				if e.ID == "<Down>" {
					log.Debug().
						Int("blockTable.SelectedRow", blockTable.SelectedRow).
						Int("windowSize", windowSize).
						Int("renderedBlocks", len(renderedBlocks)).
						Int("dy", blockTable.Dy()).
						Msg("Down")

					if blockTable.SelectedRow > windowSize-1 {
						nextTopBlockNumber := new(big.Int).Sub(ms.TopDisplayedBlock, one)
						if nextTopBlockNumber.Cmp(zero) < 0 {
							nextTopBlockNumber.SetInt64(0)
						}

						toBlockNumber := new(big.Int).Sub(nextTopBlockNumber, big.NewInt(int64(windowSize-1)))
						if toBlockNumber.Cmp(zero) < 0 {
							toBlockNumber.SetInt64(0)
						}

						if !isBlockInCache(ms.BlockCache, toBlockNumber) {
							err := ms.getBlockRange(ctx, new(big.Int).Sub(nextTopBlockNumber, big.NewInt(int64(windowSize))), toBlockNumber, rpc)
							if err != nil {
								log.Warn().Err(err).Msg("Failed to fetch blocks on page down")
								break
							}
						}

						ms.TopDisplayedBlock = nextTopBlockNumber

						blockTable.SelectedRow = len(renderedBlocks)
						setBlock = true

						forceRedraw = true
						redraw(ms, true)
						break
					}
					// blockTable.SelectedRow += 1
					blockTable.ScrollDown()

					setBlock = true
				} else if e.ID == "<Up>" {
					log.Debug().Int("blockTable.SelectedRow", blockTable.SelectedRow).Int("windowSize", windowSize).Msg("Up")

					// the last row of current window size
					if blockTable.SelectedRow == 1 {
						// Calculate the range of block numbers we are trying to page down to
						nextTopBlockNumber := new(big.Int).Add(ms.TopDisplayedBlock, one)
						if nextTopBlockNumber.Cmp(ms.HeadBlock) > 0 {
							nextTopBlockNumber.SetInt64(ms.HeadBlock.Int64())
						}

						// Calculate the 'to' block number based on the next top block number
						toBlockNumber := new(big.Int).Sub(nextTopBlockNumber, big.NewInt(int64(windowSize-1)))
						if toBlockNumber.Cmp(zero) < 0 {
							toBlockNumber.SetInt64(0)
						}

						// Fetch the blocks in the new range if they are missing
						if !isBlockInCache(ms.BlockCache, nextTopBlockNumber) {
							err := ms.getBlockRange(ctx, toBlockNumber, new(big.Int).Add(nextTopBlockNumber, big.NewInt(int64(windowSize))), rpc)
							if err != nil {
								log.Warn().Err(err).Msg("Failed to fetch blocks on page up")
								break
							}
						}

						// Update the top displayed block number
						ms.TopDisplayedBlock = nextTopBlockNumber

						blockTable.SelectedRow = 1
						setBlock = true

						// Force redraw to update the UI with the new page of blocks
						forceRedraw = true
						redraw(ms, true)
						break
					}
					// blockTable.SelectedRow -= 1
					blockTable.ScrollUp()
					setBlock = true
				}
			case "<Home>":
				ms.TopDisplayedBlock = ms.HeadBlock
				blockTable.SelectedRow = 1
				setBlock = true
			case "g":
				if previousKey == "g" {
					ms.TopDisplayedBlock = ms.HeadBlock
					blockTable.SelectedRow = 1
					setBlock = true
				}
			case "G", "<End>":
				if len(renderedBlocks) < windowSize {
					ms.TopDisplayedBlock = ms.HeadBlock
					blockTable.SelectedRow = len(renderedBlocks)
				} else {
					// windowOffset = len(allBlocks) - windowSize
					blockTable.SelectedRow = max(windowSize, len(renderedBlocks))
				}
				setBlock = true
			case "<C-f>", "<PageDown>":
				nextTopBlockNumber := new(big.Int).Sub(ms.TopDisplayedBlock, big.NewInt(int64(windowSize)))
				if nextTopBlockNumber.Cmp(zero) < 0 {
					nextTopBlockNumber.SetInt64(0)
				}

				bottomBlockNumber := new(big.Int).Sub(nextTopBlockNumber, big.NewInt(int64(windowSize-1)))
				if bottomBlockNumber.Cmp(zero) < 0 {
					bottomBlockNumber.SetInt64(0)
				}

				if ms.LowerBlock.Cmp(bottomBlockNumber) > 0 {
					log.Debug().Msgf("TEST NOT HERE %d %d", ms.LowerBlock, bottomBlockNumber)
					err := ms.getBlockRange(ctx, bottomBlockNumber, nextTopBlockNumber, rpc)
					if err != nil {
						log.Warn().Err(err).Msg("Failed to fetch blocks on page down")
						break
					}
				}

				ms.TopDisplayedBlock = nextTopBlockNumber

				blockTable.SelectedRow = 1
				setBlock = true

				log.Debug().
					Int("TopDisplayedBlock", int(ms.TopDisplayedBlock.Int64())).
					Int("bottomBlockNumber", int(bottomBlockNumber.Int64())).
					Msg("PageDown")

				forceRedraw = true
				redraw(ms, true)
			case "<C-b>", "<PageUp>":
				nextTopBlockNumber := new(big.Int).Add(ms.TopDisplayedBlock, big.NewInt(int64(windowSize)))
				if nextTopBlockNumber.Cmp(ms.HeadBlock) > 0 {
					nextTopBlockNumber.SetInt64(ms.HeadBlock.Int64())
				}

				toBlockNumber := new(big.Int).Sub(nextTopBlockNumber, big.NewInt(int64(windowSize-1)))
				if toBlockNumber.Cmp(zero) < 0 {
					toBlockNumber.SetInt64(0)
				}

				err := ms.getBlockRange(ctx, toBlockNumber, nextTopBlockNumber, rpc)
				if err != nil {
					log.Warn().Err(err).Msg("Failed to fetch blocks on page up")
					break
				}

				ms.TopDisplayedBlock = nextTopBlockNumber

				blockTable.SelectedRow = 1
				setBlock = true

				log.Debug().
					Int("TopDisplayedBlock", int(ms.TopDisplayedBlock.Int64())).
					Int("toBlockNumber", int(toBlockNumber.Int64())).
					Msg("PageDown")

				forceRedraw = true
				redraw(ms, true)
			default:
				log.Trace().Str("id", e.ID).Msg("Unknown ui event")
			}

			if previousKey == "g" {
				previousKey = ""
			} else {
				previousKey = e.ID
			}

			if !forceRedraw {
				redraw(ms)
			}
		case <-ticker.C:
			if currentBn != ms.HeadBlock {
				currentBn = ms.HeadBlock
				redraw(ms)
			}
		}
	}
}

func isBlockInCache(cache *lru.Cache, blockNumber *big.Int) bool {
	_, exists := cache.Get(blockNumber.String())
	return exists
}

func max(nums ...int) int {
	m := nums[0]
	for _, n := range nums {
		if m < n {
			m = n
		}
	}
	return m
}

// checkBatchRequestsSupport checks if batch requests are supported by making a sample batch request.
// https://ethereum.org/en/developers/docs/apis/json-rpc/#eth_blocknumber
func checkBatchRequestsSupport(ctx context.Context, ec *ethrpc.Client) error {
	batchRequest := []ethrpc.BatchElem{
		{Method: "eth_blockNumber"},
		{Method: "eth_blockNumber"},
	}
	return ec.BatchCallContext(ctx, batchRequest)
}

func getColumnWidths(columnRatio []int, width int) (columnWidths []int) {
	totalRatio := 0
	for _, ratio := range columnRatio {
		totalRatio += ratio
	}

	columnWidths = make([]int, len(columnRatio))
	for i, ratio := range columnRatio {
		columnWidths[i] = width * ratio / totalRatio
	}

	return
}
