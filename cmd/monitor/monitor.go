package monitor

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/maticnetwork/polygon-cli/util"

	_ "embed"

	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"

	"github.com/cenkalti/backoff/v4"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/maticnetwork/polygon-cli/metrics"
	"github.com/maticnetwork/polygon-cli/rpctypes"
	"github.com/rs/zerolog/log"
)

var (
	windowSize                   int
	batchSize                    int
	interval                     time.Duration
	one                          = big.NewInt(1)
	zero                         = big.NewInt(0)
	selectedBlock                rpctypes.PolyBlock
	currentlyFetchingHistoryLock sync.RWMutex
	observedPendingTxs           historicalRange
)

type (
	monitorStatus struct {
		ChainID      *big.Int
		HeadBlock    *big.Int
		PeerCount    uint64
		GasPrice     *big.Int
		PendingCount uint64
		BlockCache   *lru.Cache

		MaxBlockRetrieved *big.Int
		MinBlockRetrieved *big.Int
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
	uiSkeleton      struct {
		h0  *widgets.Paragraph
		h1  *widgets.Paragraph
		h2  *widgets.Paragraph
		h3  *widgets.Paragraph
		h4  *widgets.Paragraph
		sl0 *widgets.Sparkline
		sl1 *widgets.Sparkline
		sl2 *widgets.Sparkline
		sl3 *widgets.Sparkline
		sl4 *widgets.Sparkline
		b1  *widgets.List
		b2  *widgets.List
	}
	monitorMode int
)

const (
	monitorModeHelp monitorMode = iota
	monitorModeExplorer
	monitorModeBlock
)

func monitor(ctx context.Context) error {
	rpc, err := ethrpc.DialContext(ctx, rpcUrl)
	if err != nil {
		log.Error().Err(err).Msg("Unable to dial rpc")
		return err
	}
	ec := ethclient.NewClient(rpc)

	ms := new(monitorStatus)
	ms.BlockCache, _ = lru.New(1000)
	ms.MaxBlockRetrieved = big.NewInt(0)

	ms.ChainID = big.NewInt(0)
	ms.PendingCount = 0
	observedPendingTxs = make(historicalRange, 0)

	isUiRendered := false
	errChan := make(chan error)
	go func() {
		select {
		case <-ctx.Done(): // listens for a cancellation signal
			return // exit the goroutine when the context is done
		default:
			for {
				err = fetchBlocks(ctx, ec, ms, rpc, isUiRendered)
				if err != nil {
					continue
				}

				if !isUiRendered {
					go func() {
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
func prependLatestBlocks(ctx context.Context, ms *monitorStatus, rpc *ethrpc.Client) {
	from := new(big.Int).Sub(ms.HeadBlock, big.NewInt(int64(batchSize-1)))
	// Prevent getBlockRange from fetching duplicate blocks.
	if ms.MaxBlockRetrieved.Cmp(from) == 1 {
		from.Add(ms.MaxBlockRetrieved, big.NewInt(1))
	}

	if from.Cmp(zero) < 0 {
		from.SetInt64(0)
	}

	log.Debug().
		Int64("from", from.Int64()).
		Int64("to", ms.HeadBlock.Int64()).
		Int64("max", ms.MaxBlockRetrieved.Int64()).
		Msg("Fetching latest blocks")

	err := ms.getBlockRange(ctx, from, ms.HeadBlock, rpc)
	if err != nil {
		log.Error().Err(err).Msg("There was an issue fetching the block range")
	}
}

func appendOlderBlocks(ctx context.Context, ms *monitorStatus, rpc *ethrpc.Client) error {
	if ms.MinBlockRetrieved == nil {
		log.Warn().Msg("Nil min block")
		return fmt.Errorf("the min block is nil")
	}
	if !currentlyFetchingHistoryLock.TryLock() {
		return fmt.Errorf("the function is currently locked")
	}
	defer currentlyFetchingHistoryLock.Unlock()

	to := new(big.Int).Sub(ms.MinBlockRetrieved, one)
	from := new(big.Int).Sub(to, big.NewInt(int64(batchSize-1)))
	if from.Cmp(zero) < 0 {
		from.SetInt64(0)
	}

	log.Debug().
		Int64("from", from.Int64()).
		Int64("to", to.Int64()).
		Int64("min", ms.MinBlockRetrieved.Int64()).
		Msg("Fetching older blocks")

	err := ms.getBlockRange(ctx, from, to, rpc)
	if err != nil {
		log.Error().Err(err).Msg("There was an issue fetching the block range")
		return err
	}
	return nil
}

const maxHistoricalPoints = 1000 // set a limit to the number of historical points

func fetchBlocks(ctx context.Context, ec *ethclient.Client, ms *monitorStatus, rpc *ethrpc.Client, isUiRendered bool) (err error) {
	var cs *chainState
	cs, err = getChainState(ctx, ec)
	if err != nil {
		log.Error().Err(err).Msg("Encountered issue fetching network information")
		time.Sleep(interval)
		return err
	}
	if len(observedPendingTxs) >= maxHistoricalPoints {
		// remove the oldest data point
		observedPendingTxs = observedPendingTxs[1:]
	}
	observedPendingTxs = append(observedPendingTxs, historicalDataPoint{SampleTime: time.Now(), SampleValue: float64(cs.PendingCount)})

	log.Debug().Uint64("PeerCount", cs.PeerCount).Uint64("ChainID", cs.ChainID.Uint64()).Uint64("HeadBlock", cs.HeadBlock).Uint64("GasPrice", cs.GasPrice.Uint64()).Msg("Fetching blocks")

	if isUiRendered && batchSize < 0 {
		_, termHeight := ui.TerminalDimensions()
		batchSize = termHeight/2 - 4
	} else {
		batchSize = 50
	}

	ms.HeadBlock = new(big.Int).SetUint64(cs.HeadBlock)
	ms.ChainID = cs.ChainID
	ms.PeerCount = cs.PeerCount
	ms.GasPrice = cs.GasPrice
	ms.PendingCount = cs.PendingCount

	prependLatestBlocks(ctx, ms, rpc)
	if shouldLoadMoreHistory(ctx, ms) {
		err = appendOlderBlocks(ctx, ms, rpc)
		if err != nil {
			log.Warn().Err(err).Msg("Unable to append more history")
		}
	}

	return
}

// shouldLoadMoreHistory is meant to decide if we should keep fetching more block history. The idea is that if the user
// hasn't scrolled within a batch size of the minimum of the page, we won't  keep loading more history
func shouldLoadMoreHistory(ctx context.Context, ms *monitorStatus) bool {
	if ms.MinBlockRetrieved == nil {
		return false
	}
	if selectedBlock == nil {
		return false
	}
	minBlockNumber := ms.MinBlockRetrieved.Int64()
	selectedBlockNumber := selectedBlock.Number().Int64()
	if minBlockNumber == 0 {
		return false
	}
	if minBlockNumber < selectedBlockNumber-(5*int64(batchSize)) {
		return false
	}
	return true
}

func (ms *monitorStatus) getBlockRange(ctx context.Context, from, to *big.Int, rpc *ethrpc.Client) error {
	blms := make([]ethrpc.BatchElem, 0)
	for i := from; i.Cmp(to) != 1; i.Add(i, one) {
		r := new(rpctypes.RawBlockResponse)
		var err error
		blms = append(blms, ethrpc.BatchElem{
			Method: "eth_getBlockByNumber",
			Args:   []interface{}{"0x" + i.Text(16), true},
			Result: r,
			Error:  err,
		})
	}
	if len(blms) == 0 {
		return nil
	}
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 3 * time.Minute
	retryable := func() error {
		err := rpc.BatchCallContext(ctx, blms)
		return err
	}
	err := backoff.Retry(retryable, b)
	if err != nil {
		return err
	}
	for _, b := range blms {
		if b.Error != nil {
			return b.Error
		}
		pb := rpctypes.NewPolyBlock(b.Result.(*rpctypes.RawBlockResponse))

		ms.BlockCache.Add(pb.Number().String(), pb)

		if ms.MaxBlockRetrieved.Cmp(pb.Number()) == -1 {
			ms.MaxBlockRetrieved = pb.Number()
		}
		if ms.MinBlockRetrieved == nil || (ms.MinBlockRetrieved.Cmp(pb.Number()) == 1 && pb.Number().Cmp(zero) == 1) {
			ms.MinBlockRetrieved = pb.Number()
		}
	}

	return nil
}

func setUISkeleton() (blockTable *widgets.List, grid *ui.Grid, blockGrid *ui.Grid, termUi uiSkeleton) {
	blockTable = widgets.NewList()
	blockTable.TextStyle = ui.NewStyle(ui.ColorWhite)
	termUi = uiSkeleton{}

	termUi.h0 = widgets.NewParagraph()
	termUi.h0.Title = "Current"

	termUi.h1 = widgets.NewParagraph()
	termUi.h1.Title = "Gas Price"

	termUi.h2 = widgets.NewParagraph()
	termUi.h2.Title = "Current"

	termUi.h3 = widgets.NewParagraph()
	termUi.h3.Title = "Chain ID"

	termUi.h4 = widgets.NewParagraph()
	termUi.h4.Title = "Avg Block Time"

	termUi.sl0 = widgets.NewSparkline()
	termUi.sl0.LineColor = ui.ColorRed
	slg0 := widgets.NewSparklineGroup(termUi.sl0)
	slg0.Title = "TXs / Block"

	termUi.sl1 = widgets.NewSparkline()
	termUi.sl1.LineColor = ui.ColorGreen
	slg1 := widgets.NewSparklineGroup(termUi.sl1)
	slg1.Title = "Gas Price"

	termUi.sl2 = widgets.NewSparkline()
	termUi.sl2.LineColor = ui.ColorYellow
	slg2 := widgets.NewSparklineGroup(termUi.sl2)
	slg2.Title = "Block Size"

	termUi.sl3 = widgets.NewSparkline()
	termUi.sl3.LineColor = ui.ColorBlue
	slg3 := widgets.NewSparklineGroup(termUi.sl3)
	slg3.Title = "Pending Tx"

	termUi.sl4 = widgets.NewSparkline()
	termUi.sl4.LineColor = ui.ColorMagenta
	slg4 := widgets.NewSparklineGroup(termUi.sl4)
	slg4.Title = "Gas Used"

	grid = ui.NewGrid()
	blockGrid = ui.NewGrid()

	b0 := widgets.NewParagraph()
	b0.Title = "Block Headers"
	b0.Text = "Use the arrow keys to scroll through the transactions. Press <Esc> to go back to the explorer view"

	termUi.b1 = widgets.NewList()
	termUi.b1.Title = "Block Info"
	termUi.b1.TextStyle = ui.NewStyle(ui.ColorYellow)
	termUi.b1.WrapText = false

	termUi.b2 = widgets.NewList()
	termUi.b2.Title = "Transactions"
	termUi.b2.TextStyle = ui.NewStyle(ui.ColorGreen)
	termUi.b2.WrapText = true

	blockGrid.Set(
		ui.NewRow(1.0/10, b0),

		ui.NewRow(9.0/10,
			ui.NewCol(1.0/2, termUi.b1),
			ui.NewCol(1.0/2, termUi.b2),
		),
	)

	grid.Set(
		ui.NewRow(1.0/10,
			ui.NewCol(1.0/5, termUi.h0),
			ui.NewCol(1.0/5, termUi.h1),
			ui.NewCol(1.0/5, termUi.h2),
			ui.NewCol(1.0/5, termUi.h3),
			ui.NewCol(1.0/5, termUi.h4),
		),

		ui.NewRow(4.0/10,
			ui.NewCol(1.0/5, slg0),
			ui.NewCol(1.0/5, slg1),
			ui.NewCol(1.0/5, slg2),
			ui.NewCol(1.0/5, slg3),
			ui.NewCol(1.0/5, slg4),
		),
		ui.NewRow(5.0/10, blockTable),
	)

	return
}

func updateAllBlocks(ms *monitorStatus) []rpctypes.PolyBlock {
	var blocks []rpctypes.PolyBlock

	// Retrieve all current items from the LRU cache.
	// Since the cache has no inherent order, we will need to sort them if necessary.
	for _, key := range ms.BlockCache.Keys() {
		if value, ok := ms.BlockCache.Peek(key); ok {
			block, ok := value.(rpctypes.PolyBlock)
			if ok {
				blocks = append(blocks, block)
			}
		}
	}

	// Assuming blocks need to be sorted, you'd sort them here.
	// This assumes that metrics.SortableBlocks is a type that can be sorted.
	sort.Sort(metrics.SortableBlocks(blocks))

	return blocks
}

func renderMonitorUI(ctx context.Context, ec *ethclient.Client, ms *monitorStatus, rpc *ethrpc.Client) error {
	if err := ui.Init(); err != nil {
		return err
	}
	defer ui.Close()

	currentMode := monitorModeExplorer

	blockTable, grid, blockGrid, termUi := setUISkeleton()

	termWidth, termHeight := ui.TerminalDimensions()
	windowSize = termHeight/2 - 4
	grid.SetRect(0, 0, termWidth, termHeight)
	blockGrid.SetRect(0, 0, termWidth, termHeight)

	var setBlock = false
	var allBlocks metrics.SortableBlocks
	var renderedBlocks metrics.SortableBlocks
	windowOffset := 0

	redraw := func(ms *monitorStatus, force ...bool) {
		log.Debug().Interface("ms", ms).Msg("Redrawing")

		if currentMode == monitorModeHelp {
			// TODO add some help context?
		} else if currentMode == monitorModeBlock {
			// render a block
			termUi.b1.Rows = metrics.GetSimpleBlockFields(selectedBlock)
			termUi.b2.Rows = metrics.GetSimpleBlockTxFields(selectedBlock, ms.ChainID)

			ui.Clear()
			ui.Render(blockGrid)
			return
		}

		if blockTable.SelectedRow == 0 || len(force) > 0 && force[0] {
			allBlocks = updateAllBlocks(ms)
			sort.Sort(allBlocks)
		}
		start := len(allBlocks) - windowSize - windowOffset
		if start < 0 {
			start = 0
		}
		end := len(allBlocks) - windowOffset
		renderedBlocks = allBlocks[start:end]

		termUi.h0.Text = fmt.Sprintf("Height: %s\nTime: %s", ms.HeadBlock.String(), time.Now().Format("02 Jan 06 15:04:05 MST"))
		gasGwei := new(big.Int).Div(ms.GasPrice, metrics.UnitShannon)
		termUi.h1.Text = fmt.Sprintf("%s gwei", gasGwei.String())
		termUi.h2.Text = fmt.Sprintf("%d Peers\n%d Pending Tx", ms.PeerCount, ms.PendingCount)
		termUi.h3.Text = ms.ChainID.String()
		termUi.h4.Text = fmt.Sprintf("%0.2f", metrics.GetMeanBlockTime(renderedBlocks))

		termUi.sl0.Data = metrics.GetTxsPerBlock(renderedBlocks)
		termUi.sl1.Data = metrics.GetMeanGasPricePerBlock(renderedBlocks)
		termUi.sl2.Data = metrics.GetSizePerBlock(renderedBlocks)
		// termUi.sl3.Data = metrics.GetUnclesPerBlock(renderedBlocks)
		termUi.sl3.Data = observedPendingTxs.getValues(25)
		termUi.sl4.Data = metrics.GetGasPerBlock(renderedBlocks)

		// If a row has not been selected, continue to update the list with new blocks.
		rows, title := metrics.GetSimpleBlockRecords(renderedBlocks)
		blockTable.Rows = rows
		blockTable.Title = title

		blockTable.TextStyle = ui.NewStyle(ui.ColorWhite)
		blockTable.SelectedRowStyle = ui.NewStyle(ui.ColorWhite, ui.ColorRed, ui.ModifierBold)
		if blockTable.SelectedRow > 0 && blockTable.SelectedRow <= len(blockTable.Rows) {
			// Only changed the selected block when the user presses the up down keys.
			// Otherwise this will adjust when the table is updated automatically.
			if setBlock {
				selectedBlock = renderedBlocks[len(renderedBlocks)-blockTable.SelectedRow]
				setBlock = false
				log.Debug().Uint64("blockNumber", selectedBlock.Number().Uint64()).Msg("Selected block changed")
			}
		}

		ui.Render(grid)
	}

	currentBn := ms.HeadBlock
	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	redraw(ms)

	currIdx := 0
	previousKey := ""
	for {
		forceRedraw := false
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return nil
			case "<Escape>":
				blockTable.SelectedRow = 0
				currentMode = monitorModeExplorer
				windowOffset = 0
			case "<Enter>":
				if blockTable.SelectedRow > 0 {
					currentMode = monitorModeBlock
				}
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				grid.SetRect(0, 0, payload.Width, payload.Height)
				blockGrid.SetRect(0, 0, payload.Width, payload.Height)
				_, termHeight = ui.TerminalDimensions()
				windowSize = termHeight/2 - 4
				ui.Clear()
			case "<Up>", "<Down>":
				if currentMode == monitorModeBlock {
					if len(termUi.b2.Rows) != 0 && e.ID == "<Down>" {
						termUi.b2.ScrollDown()
					} else if len(termUi.b2.Rows) != 0 && e.ID == "<Up>" {
						termUi.b2.ScrollUp()
					}
					break
				}

				if blockTable.SelectedRow == 0 {
					currIdx = 1
					blockTable.SelectedRow = currIdx
					setBlock = true
					break
				}
				currIdx = blockTable.SelectedRow

				if e.ID == "<Down>" {
					log.Debug().
						Int("currIdx", currIdx).
						Int("windowSize", windowSize).
						Int("renderedBlocks", len(renderedBlocks)).
						Int("dy", blockTable.Dy()).
						Int("windowOffset", windowOffset).
						Int("allBlocks", len(allBlocks)).
						Msg("Down")

					// the last row of current window size
					if currIdx > windowSize-1 {
						if windowOffset+windowSize < len(allBlocks) {
							windowOffset += 1
						} else {
							err := appendOlderBlocks(ctx, ms, rpc)
							if err != nil {
								log.Warn().Err(err).Msg("Unable to append more history")
							}
							forceRedraw = true
							redraw(ms, true)
							break
						}
					}
					currIdx += 1
					setBlock = true
				} else if e.ID == "<Up>" {
					log.Debug().Int("currIdx", currIdx).Int("windowSize", windowSize).Msg("Up")
					if currIdx <= 1 && windowOffset > 0 {
						windowOffset -= 1
						break
					}
					currIdx -= 1
					setBlock = true
				}
				// need a better way to understand how many rows are visible
				if currIdx > 0 && currIdx <= windowSize && currIdx <= len(renderedBlocks) {
					blockTable.SelectedRow = currIdx
				}
			case "<Home>":
				windowOffset = 0
				blockTable.SelectedRow = 1
				setBlock = true
			case "g":
				if previousKey == "g" {
					windowOffset = 0
					blockTable.SelectedRow = 1
					setBlock = true
				}
			case "G", "<End>":
				if len(renderedBlocks) < windowSize {
					windowOffset = 0
					blockTable.SelectedRow = len(renderedBlocks)
				} else {
					windowOffset = len(allBlocks) - windowSize
					blockTable.SelectedRow = max(windowSize, len(renderedBlocks))
				}
				setBlock = true
			case "<C-f>", "<PageDown>":
				if len(renderedBlocks) < windowSize {
					windowOffset = 0
					blockTable.SelectedRow = len(renderedBlocks)
					break
				}
				windowOffset += windowSize
				// good to go to next page but not enough blocks to fill page
				if windowOffset > len(allBlocks)-windowSize {
					err := appendOlderBlocks(ctx, ms, rpc)
					if err != nil {
						log.Warn().Err(err).Msg("Unable to append more history")
					}
					forceRedraw = true
					redraw(ms, true)
				}
				blockTable.SelectedRow = len(renderedBlocks)
				setBlock = true
			case "<C-b>", "<PageUp>":
				windowOffset -= windowSize
				if windowOffset < 0 {
					windowOffset = 0
					blockTable.SelectedRow = 1
				}
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

func max(nums ...int) int {
	m := nums[0]
	for _, n := range nums {
		if m < n {
			m = n
		}
	}
	return m
}
