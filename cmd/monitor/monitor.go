package monitor

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/0xPolygon/polygon-cli/util"
	lru "github.com/hashicorp/golang-lru"

	_ "embed"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	ethrpc "github.com/ethereum/go-ethereum/rpc"

	"github.com/0xPolygon/polygon-cli/cmd/monitor/ui"
	"github.com/0xPolygon/polygon-cli/metrics"
	"github.com/0xPolygon/polygon-cli/rpctypes"
	"github.com/cenkalti/backoff/v4"
	termui "github.com/gizak/termui/v3"
	"github.com/rs/zerolog/log"
)

var errBatchRequestsNotSupported = errors.New("batch requests are not supported")

var (
	// windowSize determines the number of blocks to display in the monitor UI at one time.
	windowSize int

	// batchSize holds the number of blocks to fetch in one batch.
	// It can be adjusted dynamically based on network conditions.
	batchSize SafeBatchSize

	// interval specifies the time duration to wait between each update cycle.
	interval time.Duration

	// one and zero are big.Int representations of 1 and 0, used for convenience in calculations.
	one  = big.NewInt(1)
	zero = big.NewInt(0)

	// observedPendingTxs holds a historical record of the number of pending transactions.
	observedPendingTxs historicalRange

	// maxDataPoints defines the maximum number of data points to keep in historical records.
	maxDataPoints = 1000

	// maxConcurrency defines the maximum number of goroutines that can fetch block data concurrently.
	maxConcurrency = 10

	// semaphore is a channel used to control the concurrency of block data fetch operations.
	semaphore = make(chan struct{}, maxConcurrency)
)

type (
	monitorStatus struct {
		TopDisplayedBlock    *big.Int
		UpperBlock           *big.Int
		LowerBlock           *big.Int
		ChainID              *big.Int
		ForkID               uint64
		HeadBlock            *big.Int
		SafeBlock            *big.Int
		FinalizedBlock       *big.Int
		PeerCount            uint64
		GasPrice             *big.Int
		TxPoolStatus         txPoolStatus
		ZkEVMBatches         zkEVMBatches
		SelectedBlock        rpctypes.PolyBlock
		SelectedTransaction  rpctypes.PolyTransaction
		BlockCache           *lru.Cache   `json:"-"`
		BlocksLock           sync.RWMutex `json:"-"`
		RollupAddress        string
		RollupManagerAddress string
	}
	chainState struct {
		HeadBlock            uint64
		SafeBlock            uint64
		FinalizedBlock       uint64
		ChainID              *big.Int
		PeerCount            uint64
		GasPrice             *big.Int
		TxPoolStatus         txPoolStatus
		ZkEVMBatches         zkEVMBatches
		ForkID               uint64
		RollupAddress        string
		RollupManagerAddress string
	}
	txPoolStatus struct {
		pending uint64
		queued  uint64
	}
	zkEVMBatches struct {
		trusted  uint64
		virtual  uint64
		verified uint64
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
	monitorModeSelectBlock
	monitorModeBlock
	monitorModeTransaction
)

func monitor(ctx context.Context) error {
	rpc, err := ethrpc.DialContext(ctx, rpcUrl)
	if err != nil {
		log.Error().Err(err).Msg("Unable to dial rpc")
		return err
	}
	ec := ethclient.NewClient(rpc)
	latestBlockNumber, err := ec.BlockNumber(ctx)
	if err != nil {
		return err
	}

	// Check if batch requests are supported.
	if err = checkBatchRequestsSupport(ctx, ec.Client()); err != nil {
		return errBatchRequestsNotSupported
	}

	// Check if tx pool status is supported.
	txPoolStatusSupported := false
	if _, _, err = util.GetTxPoolStatus(rpc); err != nil {
		log.Debug().Err(err).Msg("Unable to get tx pool status")
	} else {
		txPoolStatusSupported = true
	}

	// Check if zkevm batches are supported.
	zkEVMBatchesSupported := false
	if _, _, _, err = util.GetZkEVMBatches(rpc); err != nil {
		log.Debug().Err(err).Msg("Unable to get zkevm batches")
	} else {
		zkEVMBatchesSupported = true
	}

	// check if net peerCount is supported
	peerCountSupported := false
	if _, err = ec.PeerCount(ctx); err != nil {
		log.Debug().Err(err).Msg("Unable to fake peer count")
	} else {
		peerCountSupported = true
	}

	// check if EIP-1559 is supported
	eip1559Supported := false
	latestBlock, err := ec.BlockByNumber(ctx, big.NewInt(0).SetUint64(latestBlockNumber))
	if err != nil {
		log.Debug().Err(err).Msg("Unable to get latest block")
	} else {
		if latestBlock.BaseFee() == nil {
			log.Debug().Msg("EIP-1559 not supported")
			eip1559Supported = false
		} else {
			log.Debug().Msg("EIP-1559 supported")
			eip1559Supported = true
		}
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
	ms.TxPoolStatus = txPoolStatus{}
	ms.ZkEVMBatches = zkEVMBatches{}

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
		case <-ctx.Done():
			return
		default:
			for {
				err = fetchCurrentBlockData(ctx, ec, ms, isUiRendered, txPoolStatusSupported, zkEVMBatchesSupported, peerCountSupported)
				if err != nil {
					log.Error().Msg(fmt.Sprintf("Error: unable to fetch current block data: %v", err))
					// Send the error to the errChan channel to return.
					errChan <- err
					return
				}
				if ms.TopDisplayedBlock == nil || ms.SelectedBlock == nil {
					ms.TopDisplayedBlock = ms.HeadBlock
				}
				if !isUiRendered {
					go func() {
						errChan <- renderMonitorUI(ctx, ec, ms, rpc, txPoolStatusSupported, zkEVMBatchesSupported, eip1559Supported)
					}()
					isUiRendered = true
				}

				time.Sleep(interval)
			}
		}
	}()

	if !isUiRendered {
		// If UI cannot be rendered and returns, close the goroutine.
		ctx.Done()
	}
	// Receive the errors from errChan and return to exit.
	err = <-errChan
	return err
}

func getChainState(ctx context.Context, ec *ethclient.Client, txPoolStatusSupported, zkEVMBatchesSupported, peerCountSupported bool) (*chainState, error) {
	var err error
	cs := new(chainState)
	cs.HeadBlock, err = ec.BlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("couldn't fetch block number: %s", err.Error())
	}

	safeBlock, err := ec.HeaderByNumber(ctx, big.NewInt(int64(rpc.SafeBlockNumber)))
	if err != nil {
		log.Debug().Err(err).Msg("Unable to fetch safe block number")
		cs.SafeBlock = 0
	} else if safeBlock != nil {
		cs.SafeBlock = safeBlock.Number.Uint64()
	}

	finalizedBlock, err := ec.HeaderByNumber(ctx, big.NewInt(int64(rpc.FinalizedBlockNumber)))
	if err != nil {
		log.Debug().Err(err).Msg("Unable to fetch finalized block number")
		cs.FinalizedBlock = 0
	} else if finalizedBlock != nil {
		cs.FinalizedBlock = finalizedBlock.Number.Uint64()
	}

	cs.ChainID, err = ec.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("couldn't fetch chain id: %s", err.Error())
	}

	if peerCountSupported {
		cs.PeerCount, err = ec.PeerCount(ctx)
		if err != nil {
			log.Debug().Err(err).Msg("Using fake peer count")
			cs.PeerCount = 0
		}
	}

	cs.GasPrice, err = ec.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("couldn't estimate gas: %s", err.Error())
	}

	if txPoolStatusSupported {
		cs.TxPoolStatus.pending, cs.TxPoolStatus.queued, err = util.GetTxPoolStatus(ec.Client())
		if err != nil {
			log.Debug().Err(err).Msg("Unable to get tx pool status")
		}
	}

	if zkEVMBatchesSupported {
		cs.ZkEVMBatches.trusted, cs.ZkEVMBatches.virtual, cs.ZkEVMBatches.verified, err = util.GetZkEVMBatches(ec.Client())
		if err != nil {
			log.Debug().Err(err).Msg("Unable to get zkevm batches")
		}

		cs.ForkID, err = util.GetForkID(ec.Client())
		if err != nil {
			log.Debug().Err(err).Msg("Unable to get fork id")
		}

		cs.RollupAddress, err = util.GetRollupAddress(ec.Client())
		if err != nil {
			log.Debug().Err(err).Msg("Unable to get rollup address")
		}

		cs.RollupManagerAddress, err = util.GetRollupManagerAddress(ec.Client())
		if err != nil {
			log.Debug().Err(err).Msg("Unable to get rollup manager address")
		}
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

func fetchCurrentBlockData(ctx context.Context, ec *ethclient.Client, ms *monitorStatus, isUiRendered, txPoolStatusSupported, zkEVMBatchesSupported, peerCountSupported bool) (err error) {
	var cs *chainState
	cs, err = getChainState(ctx, ec, txPoolStatusSupported, zkEVMBatchesSupported, peerCountSupported)
	if err != nil {
		log.Error().Err(err).Msg("Encountered issue fetching network information")
		time.Sleep(interval)
		return err
	}
	observedPendingTxs = append(observedPendingTxs, historicalDataPoint{SampleTime: time.Now(), SampleValue: float64(cs.TxPoolStatus.pending)})
	if len(observedPendingTxs) > maxDataPoints {
		observedPendingTxs = observedPendingTxs[len(observedPendingTxs)-maxDataPoints:]
	}

	log.Debug().Uint64("PeerCount", cs.PeerCount).Uint64("ChainID", cs.ChainID.Uint64()).Uint64("HeadBlock", cs.HeadBlock).Uint64("GasPrice", cs.GasPrice.Uint64()).Msg("Fetching blocks")

	if batchSize.Get() == 100 && batchSize.Auto() {
		newBatchSize := blockCacheLimit
		batchSize.Set(newBatchSize, true)
		log.Debug().Msgf("Auto-adjusted batchSize to %d based on cache limit", newBatchSize)
	}

	ms.HeadBlock = new(big.Int).SetUint64(cs.HeadBlock)
	ms.SafeBlock = new(big.Int).SetUint64(cs.SafeBlock)
	ms.FinalizedBlock = new(big.Int).SetUint64(cs.FinalizedBlock)

	ms.ChainID = cs.ChainID
	ms.PeerCount = cs.PeerCount
	ms.GasPrice = cs.GasPrice
	ms.TxPoolStatus = cs.TxPoolStatus
	ms.ZkEVMBatches = cs.ZkEVMBatches
	ms.ForkID = cs.ForkID
	ms.RollupAddress = cs.RollupAddress
	ms.RollupManagerAddress = cs.RollupManagerAddress

	return
}

func (ms *monitorStatus) getBlockRange(ctx context.Context, to *big.Int, rpc *ethrpc.Client) error {
	desiredBatchSize := new(big.Int).SetInt64(int64(batchSize.Get()))

	halfBatchSize := new(big.Int).Div(desiredBatchSize, big.NewInt(2))

	provisionalStartBlock := new(big.Int).Sub(to, halfBatchSize)
	provisionalEndBlock := new(big.Int).Add(to, halfBatchSize)

	log.Debug().Int64("desiredBatchSize", int64(batchSize.Get()))

	startBlock := big.NewInt(0).Set(provisionalStartBlock)
	if startBlock.Cmp(zero) < 0 {
		startBlock.SetInt64(0)
	}

	endBlock := big.NewInt(0).Set(provisionalEndBlock)
	if endBlock.Cmp(ms.HeadBlock) > 0 {
		endBlock.Set(ms.HeadBlock)
	}

	if new(big.Int).Sub(endBlock, startBlock).Cmp(desiredBatchSize) < 0 {
		if startBlock.Cmp(zero) == 0 {
			possibleEndBlock := new(big.Int).Add(startBlock, desiredBatchSize)
			if possibleEndBlock.Cmp(ms.HeadBlock) <= 0 {
				endBlock.Set(possibleEndBlock)
			} else {
				endBlock.Set(ms.HeadBlock)
			}
		} else if endBlock.Cmp(ms.HeadBlock) == 0 {
			possibleStartBlock := new(big.Int).Sub(endBlock, desiredBatchSize)
			if possibleStartBlock.Cmp(zero) >= 0 {
				startBlock.Set(possibleStartBlock)
			} else {
				startBlock.SetInt64(0)
			}
		}
	}

	ms.LowerBlock = startBlock
	ms.UpperBlock = endBlock

	blms := make([]ethrpc.BatchElem, 0)
	for i := new(big.Int).Set(startBlock); i.Cmp(endBlock) <= 0; i.Add(i, one) {
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

	err := ms.processBatchesConcurrently(ctx, rpc, blms)
	if err != nil {
		log.Error().Err(err).Msg("Error processing batches concurrently")
		return err
	}

	return nil
}

func (ms *monitorStatus) processBatchesConcurrently(ctx context.Context, rpc *ethrpc.Client, blms []ethrpc.BatchElem) error {
	var wg sync.WaitGroup
	var errs []error = make([]error, 0)
	var errorsMutex sync.Mutex

	for i := 0; i < len(blms); i += subBatchSize {
		semaphore <- struct{}{}
		wg.Add(1)
		go func(i int) {
			defer func() {
				<-semaphore
				wg.Done()
			}()
			end := i + subBatchSize
			if end > len(blms) {
				end = len(blms)
			}
			subBatch := blms[i:end]

			b := backoff.NewExponentialBackOff()
			b.MaxElapsedTime = 3 * time.Minute
			retryable := func() error {
				err := rpc.BatchCallContext(ctx, subBatch)
				if err != nil {
					log.Error().Err(err).Msg("BatchCallContext error - retry loop")
					if strings.Contains(err.Error(), "limit") {
						return backoff.Permanent(err)
					}
				}
				return nil
			}
			if err := backoff.Retry(retryable, b); err != nil {
				log.Error().Err(err).Msg("unable to retry")
				errorsMutex.Lock()
				errs = append(errs, err)
				errorsMutex.Unlock()
				return
			}

			for _, elem := range subBatch {
				if elem.Error != nil {
					log.Error().Str("Method", elem.Method).Interface("Args", elem.Args).Err(elem.Error).Msg("Failed batch element")
				} else {
					pb := rpctypes.NewPolyBlock(elem.Result.(*rpctypes.RawBlockResponse))
					ms.BlocksLock.Lock()
					ms.BlockCache.Add(pb.Number().String(), pb)
					ms.BlocksLock.Unlock()
				}
			}

		}(i)
	}

	wg.Wait()

	return errors.Join(errs...)
}

func renderMonitorUI(ctx context.Context, ec *ethclient.Client, ms *monitorStatus, rpc *ethrpc.Client, txPoolStatusSupported, zkEVMBatchesSupported, eip1559Supported bool) error {
	if err := termui.Init(); err != nil {
		log.Error().Err(err).Msg("Failed to initialize UI")
		return err
	}
	defer termui.Close()

	currentMode := monitorModeExplorer

	blockTable, blockInfo, transactionList, transactionInformationList, transactionInfo, grid, selectGrid, blockGrid, transactionGrid, skeleton := ui.SetUISkeleton(txPoolStatusSupported, zkEVMBatchesSupported, eip1559Supported)

	termWidth, termHeight := termui.TerminalDimensions()
	windowSize = termHeight/2 - 4
	grid.SetRect(0, 0, termWidth, termHeight)
	selectGrid.SetRect(0, 0, termWidth, termHeight)
	blockGrid.SetRect(0, 0, termWidth, termHeight)
	transactionGrid.SetRect(0, 0, termWidth, termHeight)
	// Initial render needed I assume to avoid the first bad redraw
	termui.Render(grid)

	var setBlock = false
	var renderedBlocks rpctypes.SortableBlocks

	redraw := func(ms *monitorStatus, force ...bool) {
		if currentMode == monitorModeHelp {
			// TODO add some help context?
		} else if currentMode == monitorModeSelectBlock {
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
			rows, title := ui.GetSelectedBlocksList(renderedBlocks)
			blockTable.Rows = rows
			blockTable.Title = title

			// in monitorSelectModeTransaction, blocks will always be selected
			transactionColumnRatio := []int{30, 5, 20, 20, 5, 10}
			// Sanity check for renderedBlocks table to avoid panic
			if len(renderedBlocks)-blockTable.SelectedRow < 0 {
				return
			} else if len(renderedBlocks) == 0 {
				return
			} else {
				ms.SelectedBlock = renderedBlocks[len(renderedBlocks)-blockTable.SelectedRow]
			}
			blockInfo.Rows = ui.GetSimpleBlockFields(ms.SelectedBlock)
			transactionInfo.ColumnWidths = getColumnWidths(transactionColumnRatio, transactionInfo.Dx())
			transactionInfo.Rows = ui.GetBlockTxTable(ms.SelectedBlock, ms.ChainID)
			transactionInfo.Title = fmt.Sprintf("Latest Transactions for Block #%s", ms.SelectedBlock.Number().String())

			termui.Clear()
			termui.Render(selectGrid)
			return
		} else if currentMode == monitorModeBlock {
			if ms.SelectedBlock == nil {
				currentMode = monitorModeExplorer
				blockTable.SelectedRow = 0
				termui.Clear()
				termui.Render(grid)
				return
			}

			// render a block
			skeleton.BlockInfo.Rows = ui.GetSimpleBlockFields(ms.SelectedBlock)
			rows, title := ui.GetTransactionsList(ms.SelectedBlock, ms.ChainID)
			transactionList.Rows = rows
			transactionList.Title = title

			baseFee := ms.SelectedBlock.BaseFee()
			if transactionList.SelectedRow != 0 {
				transactions := ms.SelectedBlock.Transactions()
				if len(transactions) > 0 {
					index := transactionList.SelectedRow - 1
					if index >= 0 && index < len(transactions) {
						ms.SelectedTransaction = transactions[index]
					} else {
						log.Error().
							Int("row", transactionList.SelectedRow).
							Msg("Selected row is out of range for transactions")
					}
				} else {
					log.Debug().
						Int("block", int(ms.SelectedBlock.Number().Uint64())).
						Msg("No transactions available in the selected block")
				}
				transactionInformationList.Rows = ui.GetSimpleTxFields(ms.SelectedTransaction, ms.ChainID, baseFee)
			}
			termui.Clear()
			termui.Render(blockGrid)

			log.Debug().
				Int("skeleton.TransactionList.SelectedRow", transactionList.SelectedRow).
				Msg("Redrawing block mode")

			return
		} else if currentMode == monitorModeTransaction {
			baseFee := ms.SelectedBlock.BaseFee()
			transactions := ms.SelectedBlock.Transactions()
			if len(transactions) > 0 {
				index := transactionList.SelectedRow - 1
				if index >= 0 && index < len(transactions) {
					tx := transactions[index]
					skeleton.TxInfo.Rows = ui.GetSimpleTxFields(tx, ms.ChainID, baseFee)
				} else {
					log.Error().
						Int("row", transactionList.SelectedRow).
						Msg("Selected row is out of range for transactions")
				}
			} else {
				log.Debug().
					Int("block", int(ms.SelectedBlock.Number().Uint64())).
					Msg("No transactions available in the selected block")
			}
			skeleton.Receipts.Rows = ui.GetSimpleReceipt(ctx, rpc, ms.SelectedTransaction)

			termui.Clear()
			termui.Render(transactionGrid)

			log.Debug().
				Int("skeleton.TransactionList.SelectedRow", transactionList.SelectedRow).
				Msg("Redrawing transaction mode")

			return
		}

		log.Debug().
			Str("TopDisplayedBlock", ms.TopDisplayedBlock.String()).
			Int("BatchSize", batchSize.Get()).
			Str("UpperBlock", ms.UpperBlock.String()).
			Str("LowerBlock", ms.LowerBlock.String()).
			Str("ChainID", ms.ChainID.String()).
			Str("HeadBlock", ms.HeadBlock.String()).
			Uint64("PeerCount", ms.PeerCount).
			Str("GasPrice", ms.GasPrice.String()).
			Interface("TxPoolStatus", ms.TxPoolStatus).
			Interface("ZkEVMBatches", ms.ZkEVMBatches).
			Msg("Redrawing")

		if blockTable.SelectedRow == 0 {
			bottomBlockNumber := new(big.Int).Sub(ms.HeadBlock, big.NewInt(int64(windowSize-1)))
			if bottomBlockNumber.Cmp(zero) < 0 {
				bottomBlockNumber.SetInt64(0)
			}

			err := ms.getBlockRange(ctx, ms.TopDisplayedBlock, rpc)
			if err != nil {
				log.Error().Err(err).Msg("There was an issue fetching the block range")
			}
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

		var renderedBlocksMeanGasPrice []float64
		if eip1559Supported {
			renderedBlocksMeanGasPrice = metrics.GetMeanBaseFeePerBlock(renderedBlocks)
		} else {
			renderedBlocksMeanGasPrice = metrics.GetMeanGasPricePerBlock(renderedBlocks)
		}

		// First initialization will render no gas price because the GasPriceChart will have no data.
		if renderedBlocksMeanGasPrice == nil {
			skeleton.Current.Text = ui.GetCurrentText(skeleton.Current, ms.HeadBlock, ms.SafeBlock, ms.FinalizedBlock, "--", ms.PeerCount, ms.ChainID, rpcUrl)
		} else {
			if len(renderedBlocksMeanGasPrice) >= 1 {
				// Under normal cases, the gas price will be derived from the last element of the GasPriceChart with 2 decimal places precision.
				gasPriceStr := strconv.FormatFloat(renderedBlocksMeanGasPrice[len(renderedBlocksMeanGasPrice)-1]/1000000000, 'f', 2, 64)
				skeleton.Current.Text = ui.GetCurrentText(skeleton.Current, ms.HeadBlock, ms.SafeBlock, ms.FinalizedBlock, gasPriceStr, ms.PeerCount, ms.ChainID, rpcUrl)
			}
		}

		if txPoolStatusSupported {
			skeleton.TxPool.Text = ui.GetTxPoolText(skeleton.TxPool, ms.TxPoolStatus.pending, ms.TxPoolStatus.queued)
		}

		// if zkEVMBatchesSupported == true, this means the network will also support ForkIDs.
		if zkEVMBatchesSupported {
			skeleton.ZkEVM.Text = ui.GetZkEVMText(skeleton.ZkEVM, ms.ZkEVMBatches.trusted, ms.ZkEVMBatches.virtual, ms.ZkEVMBatches.verified)

			skeleton.Rollup.Text = ui.GetRollupText(skeleton.Rollup, ms.ForkID, ms.RollupAddress, ms.RollupManagerAddress)
		}

		skeleton.TxPerBlockChart.Data = metrics.GetTxsPerBlock(renderedBlocks)
		skeleton.GasPriceChart.Data = renderedBlocksMeanGasPrice // equivalent to metrics.GetMeanGasPricePerBlock(renderedBlocks)
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
		transactionColumnRatio := []int{30, 5, 20, 20, 5, 10}
		if blockTable.SelectedRow > 0 && blockTable.SelectedRow <= len(blockTable.Rows) && (len(renderedBlocks)-blockTable.SelectedRow) >= 0 {
			// Only changed the selected block when the user presses the up down keys.
			// Otherwise this will adjust when the table is updated automatically.
			if setBlock && ms.SelectedBlock != nil {
				log.Debug().
					Int("blockTable.SelectedRow", blockTable.SelectedRow).
					Int("renderedBlocks", len(renderedBlocks)).
					Msg("setBlock")

				blockInfo.Rows = ui.GetSimpleBlockFields(ms.SelectedBlock)
				transactionInfo.ColumnWidths = getColumnWidths(transactionColumnRatio, transactionInfo.Dx())
				transactionInfo.Rows = ui.GetBlockTxTable(ms.SelectedBlock, ms.ChainID)
				transactionInfo.Title = fmt.Sprintf("Latest Transactions for Block #%s", ms.SelectedBlock.Number().String())
				setBlock = false
				log.Debug().Uint64("blockNumber", ms.SelectedBlock.Number().Uint64()).Msg("Selected block changed")
			}
		} else {
			ms.SelectedBlock = nil
			blockInfo.Rows = []string{}
			transactionInfo.ColumnWidths = getColumnWidths(transactionColumnRatio, transactionInfo.Dx())
			if len(renderedBlocks) > 0 {
				i := len(renderedBlocks) - 1
				transactionInfo.Rows = ui.GetBlockTxTable(renderedBlocks[i], ms.ChainID)
				transactionInfo.Title = fmt.Sprintf("Latest Transactions for Block #%s", renderedBlocks[i].Number().String())
			}
		}

		termui.Render(grid)
	}

	currentBn := ms.HeadBlock
	uiEvents := termui.PollEvents()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	redraw(ms)

	for {
		forceRedraw := false
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return nil
			case "<Escape>":
				if currentMode == monitorModeExplorer {
					ms.TopDisplayedBlock = ms.HeadBlock
					blockTable.SelectedRow = 0

					toBlockNumber := new(big.Int).Sub(ms.TopDisplayedBlock, big.NewInt(int64(windowSize-1)))
					if toBlockNumber.Cmp(zero) < 0 {
						toBlockNumber.SetInt64(0)
					}

					err := ms.getBlockRange(ctx, ms.TopDisplayedBlock, rpc)
					if err != nil {
						log.Error().Err(err).Msg("There was an issue fetching the block range")
						break
					}
				} else if currentMode == monitorModeBlock {
					currentMode = monitorModeExplorer
					blockTable.SelectedRow = 0
				} else if currentMode == monitorModeSelectBlock {
					currentMode = monitorModeExplorer
					blockTable.SelectedRow = 0
				} else if currentMode == monitorModeTransaction {
					currentMode = monitorModeBlock
					blockTable.SelectedRow = 0
				}
			case "<Enter>":
				if (currentMode == monitorModeExplorer || currentMode == monitorModeSelectBlock) && blockTable.SelectedRow > 0 {
					currentMode = monitorModeBlock
				} else if transactionList.SelectedRow > 0 {
					currentMode = monitorModeTransaction
				}
			case "<Resize>":
				payload := e.Payload.(termui.Resize)
				grid.SetRect(0, 0, payload.Width, payload.Height)
				selectGrid.SetRect(0, 0, payload.Width, payload.Height)
				blockGrid.SetRect(0, 0, payload.Width, payload.Height)
				transactionGrid.SetRect(0, 0, payload.Width, payload.Height)
				_, termHeight = termui.TerminalDimensions()
				windowSize = termHeight/2 - 4
				termui.Clear()
			case "<Up>", "<Down>", "<MouseWheelUp>", "<MouseWheelDown>":
				up := e.ID == "<Up>" || e.ID == "<MouseWheelUp>"
				down := e.ID == "<Down>" || e.ID == "<MouseWheelDown>"

				if currentMode == monitorModeBlock {
					if len(transactionList.Rows) != 0 && down {
						transactionList.ScrollDown()
					} else if len(transactionList.Rows) != 0 && up {
						transactionList.ScrollUp()
					}
					break
				}

				if blockTable.SelectedRow == 0 {
					blockTable.SelectedRow = 1
					setBlock = true
					currentMode = monitorModeSelectBlock
					break
				}

				if down {
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

						if !ms.isBlockInCache(toBlockNumber) {
							err := ms.getBlockRange(ctx, toBlockNumber, rpc)
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
				} else if up {
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
						if !ms.isBlockInCache(nextTopBlockNumber) {
							err := ms.getBlockRange(ctx, new(big.Int).Add(nextTopBlockNumber, big.NewInt(int64(windowSize))), rpc)
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
				blockTable.SelectedRow = 1
				setBlock = true
			case "G", "<End>":
				if len(renderedBlocks) < windowSize {
					blockTable.SelectedRow = len(renderedBlocks)
				} else {
					blockTable.SelectedRow = max(windowSize, len(renderedBlocks))
				}
				setBlock = true
			case "<C-f>", "<PageDown>":
				// When pressing PageDown beyond the genesis block, redraw the monitor screen to avoid freezing at the previous rendered blocks.
				if len(renderedBlocks) == 0 {
					currentMode = monitorModeExplorer
					blockTable.SelectedRow = 0
					forceRedraw = true
					redraw(ms, true)
					break
				}

				if renderedBlocks[0].Number().String() == "0" || renderedBlocks[0].Number().String() == "1" {
					blockTable.SelectedRow = len(renderedBlocks)
					forceRedraw = true
					redraw(ms, true)
					break
				}

				// If the length of the renderedBlocks (the visible blocks in the monitor) is less than the windowSize, select the last element of the renderedBlocks
				if len(renderedBlocks) < windowSize {
					blockTable.SelectedRow = len(renderedBlocks)
				}

				// PageDown key also enforces monitorModeSelectBlock in addition to the arrow Down key
				if blockTable.SelectedRow == 0 {
					blockTable.SelectedRow = 1
					setBlock = true
					currentMode = monitorModeSelectBlock
					break
				}

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
					err := ms.getBlockRange(ctx, nextTopBlockNumber, rpc)
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
				// PageUp key also enforces monitorModeSelectBlock in addition to the arrow Up key
				if blockTable.SelectedRow == 0 {
					blockTable.SelectedRow = 1
					setBlock = true
					currentMode = monitorModeSelectBlock
					break
				}

				nextTopBlockNumber := new(big.Int).Add(ms.TopDisplayedBlock, big.NewInt(int64(windowSize)))
				if nextTopBlockNumber.Cmp(ms.HeadBlock) > 0 {
					nextTopBlockNumber.SetInt64(ms.HeadBlock.Int64())
				}

				toBlockNumber := new(big.Int).Sub(nextTopBlockNumber, big.NewInt(int64(windowSize-1)))
				if toBlockNumber.Cmp(zero) < 0 {
					toBlockNumber.SetInt64(0)
				}

				err := ms.getBlockRange(ctx, nextTopBlockNumber, rpc)
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
					Msg("PageUp")

				forceRedraw = true
				redraw(ms, true)
			default:
				log.Trace().Str("id", e.ID).Msg("Unknown ui event")
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

func (ms *monitorStatus) isBlockInCache(blockNumber *big.Int) bool {
	ms.BlocksLock.RLock()
	_, exists := ms.BlockCache.Get(blockNumber.String())
	ms.BlocksLock.RUnlock()
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
