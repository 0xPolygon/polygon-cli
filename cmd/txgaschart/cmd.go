package txgaschart

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/time/rate"
)

type txGasChartConfig struct {
	rateLimiter *rate.Limiter
	concurrency uint64
	output      string

	targetAddr string
	startBlock uint64
	endBlock   uint64

	scale string
}

type blocksMetadata struct {
	blocks []block

	minTxGasLimit    uint64
	maxTxGasLimit    uint64
	minTxGasPrice    uint64
	maxTxGasPrice    uint64
	maxBlockGasLimit uint64
	avgBlockGasUsed  uint64
	txCount          uint64
	targetTxCount    uint64
}

type block struct {
	number      uint64
	avgGasPrice uint64
	gasLimit    uint64
	txsGasLimit uint64
	gasUsed     uint64
	txs         []transaction
}

type transaction struct {
	hash     common.Hash
	gasPrice uint64
	gasLimit uint64
	target   bool
}

func buildChart(cmd *cobra.Command) error {
	ctx := context.Background()
	log.Info().
		Str("rpc_url", inputArgs.rpcURL).
		Float64("rate_limit", inputArgs.rateLimit).
		Msg("RPC connection parameters")

	log.Info().
		Uint64("start_block", inputArgs.startBlock).
		Uint64("end_block", inputArgs.endBlock).
		Str("target_address", inputArgs.targetAddr).
		Msg("Chart generation parameters")

	client, err := ethclient.DialContext(ctx, inputArgs.rpcURL)
	if err != nil {
		return err
	}
	defer client.Close()

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return err
	}

	config, err := parseFlags(ctx, client)
	if err != nil {
		return err
	}

	bm := loadBlocksMetadata(ctx, config, client, chainID)

	chartMetadata := txGasChartMetadata{
		rpcURL:  inputArgs.rpcURL,
		chainID: chainID.Uint64(),

		targetAddr: config.targetAddr,
		startBlock: config.startBlock,
		endBlock:   config.endBlock,

		blocksMetadata: bm,

		scale: config.scale,

		outputPath: config.output,
	}

	logMostUsedGasPrices(bm)

	return plotChart(chartMetadata)
}

func logMostUsedGasPrices(bm blocksMetadata) {
	x := map[uint64]uint64{}
	for _, b := range bm.blocks {
		for _, t := range b.txs {
			x[t.gasPrice]++
		}
	}

	ox := []struct {
		gasPrice uint64
		count    uint64
	}{}
	for k, v := range x {
		ox = append(ox, struct {
			gasPrice uint64
			count    uint64
		}{
			gasPrice: k,
			count:    v,
		})
	}

	slices.SortFunc(ox, func(a, b struct {
		gasPrice uint64
		count    uint64
	}) int {
		if a.count < b.count {
			return 1
		} else if a.count > b.count {
			return -1
		}
		return 0
	})

	if len(ox) > 0 {
		log.Debug().Msg("most used gas prices:")
		max := 20
		for _, v := range ox {
			log.Debug().Uint64("gas_price_wei", v.gasPrice).
				Uint64("count", v.count).
				Msg("gas price usage")
			max--
			if max <= 0 {
				break
			}
		}
	}
}

func parseFlags(ctx context.Context, client *ethclient.Client) (*txGasChartConfig, error) {
	config := &txGasChartConfig{}

	config.startBlock = inputArgs.startBlock
	config.endBlock = inputArgs.endBlock

	h, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}

	if config.endBlock == math.MaxUint64 || config.endBlock > h.Number.Uint64() {
		config.endBlock = h.Number.Uint64()
		log.Warn().Uint64("end_block", config.endBlock).Msg("end block was not set or set to a value higher than the latest block in the network, defaulting to latest block")
	}

	if config.startBlock > config.endBlock {
		return nil, fmt.Errorf("start block %d cannot be greater than end block %d", config.startBlock, config.endBlock)
	}

	const defaultBlockRange = 500

	if config.startBlock == 0 {
		if config.endBlock < defaultBlockRange {
			config.startBlock = 0
		} else {
			config.startBlock = config.endBlock - defaultBlockRange
		}

		log.Warn().Uint64("start_block", config.startBlock).
			Msg("start block was 0, defaulting to last blocks")
	}

	config.rateLimiter = nil
	if inputArgs.rateLimit > 0.0 {
		config.rateLimiter = rate.NewLimiter(rate.Limit(inputArgs.rateLimit), 1)
	}

	if len(inputArgs.targetAddr) > 0 && !common.IsHexAddress(inputArgs.targetAddr) {
		return nil, fmt.Errorf("target address %s is not a valid hex address", inputArgs.targetAddr)
	}

	config.targetAddr = inputArgs.targetAddr
	config.concurrency = inputArgs.concurrency
	config.output = inputArgs.output
	config.scale = inputArgs.scale

	return config, nil
}

func loadBlocksMetadata(ctx context.Context, config *txGasChartConfig, client *ethclient.Client, chainID *big.Int) blocksMetadata {
	// prepare worker pool
	wrk := make(chan struct{}, config.concurrency)
	for i := 0; i < cap(wrk); i++ {
		wrk <- struct{}{}
	}

	blockMutex := &sync.Mutex{}
	blocks := blocksMetadata{
		minTxGasLimit: math.MaxUint64,
		maxTxGasLimit: 0,
		minTxGasPrice: math.MaxUint64,
		maxTxGasPrice: 0,
		txCount:       0,
		targetTxCount: 0,
	}

	blocks.blocks = make([]block, config.endBlock-config.startBlock+1)
	offset := config.startBlock

	log.Info().Msg("reading blocks")

	wg := sync.WaitGroup{}
	totalGasUsed := big.NewInt(0)
	for blockNumber := config.startBlock; blockNumber <= config.endBlock; blockNumber++ {
		wg.Add(1) // notify block to process
		go func(blockNumber uint64) {
			defer wg.Done()                      // notify block done
			<-wrk                                // wait for worker slot
			defer func() { wrk <- struct{}{} }() // release worker slot

			for {
				log.Trace().Uint64("block_number", blockNumber).Msg("processing block")
				if config.rateLimiter != nil {
					_ = config.rateLimiter.Wait(ctx)
				}
				blockFromNetwork, err := client.BlockByNumber(ctx, new(big.Int).SetUint64(blockNumber))
				if err != nil {
					log.Error().Err(err).Uint64("block_number", blockNumber).Msg("failed to fetch block, retrying...")
					time.Sleep(time.Second)
					continue
				}

				txs := blockFromNetwork.Transactions()

				b := block{
					number:   blockFromNetwork.NumberU64(),
					gasLimit: blockFromNetwork.GasLimit(),
					gasUsed:  blockFromNetwork.GasUsed(),
					txs:      make([]transaction, len(blockFromNetwork.Transactions())),
				}

				blockMutex.Lock()
				blocks.maxBlockGasLimit = max(blocks.maxBlockGasLimit, b.gasLimit)
				totalGasUsed = totalGasUsed.Add(totalGasUsed, big.NewInt(int64(b.gasUsed)))
				blockMutex.Unlock()

				totalGasPrice := uint64(0)
				totalGasLimit := uint64(0)
				for txi, tx := range txs {
					signer := types.LatestSignerForChainID(chainID)
					from, _ := types.Sender(signer, tx)
					target := strings.EqualFold(from.String(), config.targetAddr)
					if !target {
						target = tx.To() != nil && strings.EqualFold(tx.To().String(), config.targetAddr)
					}
					gasPrice := tx.GasPrice().Uint64()
					gasLimit := tx.Gas()

					b.txs[txi] = transaction{
						hash:     tx.Hash(),
						gasPrice: gasPrice,
						gasLimit: gasLimit,
						target:   target,
					}

					totalGasPrice += gasPrice
					totalGasLimit += gasLimit

					blockMutex.Lock()
					blocks.minTxGasLimit = min(blocks.minTxGasLimit, gasLimit)
					blocks.maxTxGasLimit = max(blocks.maxTxGasLimit, gasLimit)
					blocks.minTxGasPrice = min(blocks.minTxGasPrice, gasPrice)
					blocks.maxTxGasPrice = max(blocks.maxTxGasPrice, gasPrice)

					blocks.txCount++
					if target {
						blocks.targetTxCount++
						log.Info().
							Uint64("block", b.number).
							Stringer("txHash", tx.Hash()).
							Uint64("gas_price_wei", gasPrice).
							Uint64("gas_limit", gasLimit).
							Msg("target tx found")
					}
					blockMutex.Unlock()
				}
				if len(txs) > 0 {
					b.avgGasPrice = uint64(totalGasPrice / uint64(len(txs)))
				} else {
					b.avgGasPrice = 1
				}

				b.txsGasLimit = totalGasLimit

				blocks.blocks[blockNumber-offset] = b
				break
			}
		}(blockNumber)
	}
	wg.Wait()

	blocks.avgBlockGasUsed = big.NewInt(0).Div(totalGasUsed, big.NewInt(int64(len(blocks.blocks)))).Uint64()

	return blocks
}

// 26683746818
