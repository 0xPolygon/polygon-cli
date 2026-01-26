package plot

import (
	"bufio"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/0xPolygon/polygon-cli/rpctypes"
	"github.com/0xPolygon/polygon-cli/util"
	"github.com/ethereum/go-ethereum/common"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/time/rate"
)

type args struct {
	rpcURL      string
	rateLimit   float64
	concurrency uint64

	renderer string

	startBlock uint64
	endBlock   uint64

	targetAddr string

	output string
	cache  string
}

var inputArgs = args{}

//go:embed usage.md
var usage string
var Cmd = &cobra.Command{
	Use:   "plot",
	Short: "Plot a chart of transaction gas prices and limits.",
	Long:  usage,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		return buildChart(cmd)
	},
}

func init() {
	f := Cmd.PersistentFlags()
	f.StringVar(&inputArgs.rpcURL, "rpc-url", "http://localhost:8545", "RPC URL of network")
	f.Float64Var(&inputArgs.rateLimit, "rate-limit", 4, "requests per second limit (use negative value to remove limit)")
	f.Uint64VarP(&inputArgs.concurrency, "concurrency", "c", 1, "number of tasks to perform concurrently (default: one at a time)")

	f.StringVar(&inputArgs.renderer, "renderer", "svg", "chart renderer (options: svg, canvas)")

	f.Uint64Var(&inputArgs.startBlock, "start-block", math.MaxUint64, "starting block number (inclusive)")
	f.Uint64Var(&inputArgs.endBlock, "end-block", math.MaxUint64, "ending block number (inclusive)")
	f.StringVar(&inputArgs.targetAddr, "target-address", "", "address that will have tx sent from or to highlighted in the chart")
	f.StringVarP(&inputArgs.output, "output", "o", "plot.html", "output file path")
	f.StringVar(&inputArgs.cache, "cache", "", "cache file path for block data (.ndjson); if set, reads from cache if exists, otherwise fetches and writes to cache")
}

// txGasChartConfig holds the configuration for generating the transaction gas chart.
type txGasChartConfig struct {
	rateLimiter *rate.Limiter
	concurrency uint64
	output      string
	cache       string

	targetAddr string
	startBlock uint64
	endBlock   uint64

	renderer string
}

// blocksMetadata holds metadata about a range of blocks.
type blocksMetadata struct {
	blocks []block

	maxBlockGasLimit uint64
	avgBlockGasUsed  uint64
	txCount          uint64
	targetTxCount    uint64
}

// block holds metadata about a single block.
type block struct {
	Hash        common.Hash   `json:"hash"`
	Number      uint64        `json:"number"`
	GasLimit    uint64        `json:"gasLimit"`
	TxsGasLimit uint64        `json:"txsGasLimit"`
	GasUsed     uint64        `json:"gasUsed"`
	Txs         []transaction `json:"txs"`
}

// transaction holds metadata about a single transaction.
type transaction struct {
	Hash     common.Hash    `json:"hash"`
	From     common.Address `json:"from"`
	To       common.Address `json:"to"`
	GasPrice uint64         `json:"gasPrice"`
	GasLimit uint64         `json:"gasLimit"`
	Target   bool           `json:"-"` // computed field, not cached
}

// isTargetTx returns true if the transaction involves the target address.
func isTargetTx(from, to common.Address, targetAddr string) bool {
	if targetAddr == "" {
		return false
	}
	return strings.EqualFold(from.Hex(), targetAddr) || strings.EqualFold(to.Hex(), targetAddr)
}

// writeCache writes block data to an NDJSON cache file.
func writeCache(path string, blocks []block) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create cache file: %w", err)
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	encoder := json.NewEncoder(writer)

	for _, b := range blocks {
		if err := encoder.Encode(b); err != nil {
			return fmt.Errorf("failed to encode block %d: %w", b.Number, err)
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush cache file: %w", err)
	}

	log.Info().Str("file", path).Int("blocks", len(blocks)).Msg("Cache written")
	return nil
}

// getBlockRange returns the min and max block numbers from a slice of blocks.
func getBlockRange(blocks []block) (minBlock, maxBlock uint64) {
	if len(blocks) == 0 {
		return 0, 0
	}
	minBlock = blocks[0].Number
	maxBlock = blocks[0].Number
	for _, b := range blocks[1:] {
		if b.Number < minBlock {
			minBlock = b.Number
		}
		if b.Number > maxBlock {
			maxBlock = b.Number
		}
	}
	return minBlock, maxBlock
}

// readCache reads block data from an NDJSON cache file.
// Returns the blocks and true if cache was read successfully, or nil and false if cache doesn't exist.
func readCache(path string, targetAddr string) ([]block, bool) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false
		}
		log.Warn().Err(err).Str("file", path).Msg("Failed to open cache file")
		return nil, false
	}
	defer f.Close()

	var blocks []block
	scanner := bufio.NewScanner(f)

	// Increase scanner buffer for large lines
	const maxLineSize = 10 * 1024 * 1024 // 10MB
	scanner.Buffer(make([]byte, 0, 64*1024), maxLineSize)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		var b block
		if err := json.Unmarshal(scanner.Bytes(), &b); err != nil {
			log.Warn().Err(err).Int("line", lineNum).Msg("Failed to parse cache line, ignoring cache")
			return nil, false
		}

		// Compute target field for each transaction
		for i := range b.Txs {
			b.Txs[i].Target = isTargetTx(b.Txs[i].From, b.Txs[i].To, targetAddr)
		}

		blocks = append(blocks, b)
	}

	if err := scanner.Err(); err != nil {
		log.Warn().Err(err).Msg("Error reading cache file")
		return nil, false
	}

	log.Info().Str("file", path).Int("blocks", len(blocks)).Msg("Cache loaded")
	return blocks, true
}

// computeBlocksMetadata computes aggregate metadata from a slice of blocks.
func computeBlocksMetadata(blocks []block) blocksMetadata {
	bm := blocksMetadata{
		blocks: blocks,
	}

	totalGasUsed := big.NewInt(0)
	for _, b := range blocks {
		bm.maxBlockGasLimit = max(bm.maxBlockGasLimit, b.GasLimit)
		totalGasUsed.Add(totalGasUsed, new(big.Int).SetUint64(b.GasUsed))

		for _, t := range b.Txs {
			bm.txCount++
			if t.Target {
				bm.targetTxCount++
			}
		}
	}

	numBlocks := len(blocks)
	if numBlocks == 0 {
		bm.avgBlockGasUsed = 0
	} else {
		avgGasUsed := new(big.Int).Div(totalGasUsed, big.NewInt(int64(numBlocks)))
		if avgGasUsed.IsUint64() {
			bm.avgBlockGasUsed = avgGasUsed.Uint64()
		} else {
			bm.avgBlockGasUsed = math.MaxUint64
		}
	}

	return bm
}

// buildChart builds the transaction gas chart based on the provided command and input arguments.
func buildChart(cmd *cobra.Command) error {
	ctx := cmd.Context()
	log.Info().
		Str("rpc_url", inputArgs.rpcURL).
		Float64("rate_limit", inputArgs.rateLimit).
		Msg("RPC connection parameters")

	log.Info().
		Uint64("start_block", inputArgs.startBlock).
		Uint64("end_block", inputArgs.endBlock).
		Str("target_address", inputArgs.targetAddr).
		Msg("Chart generation parameters")

	client, err := ethrpc.DialContext(ctx, inputArgs.rpcURL)
	if err != nil {
		return err
	}
	defer client.Close()

	chainID, err := util.GetChainID(ctx, client)
	if err != nil {
		return err
	}

	config, err := parseFlags(ctx, client)
	if err != nil {
		return err
	}

	var blocks []block

	// Try to load from cache if cache path is specified
	if config.cache != "" {
		if cached, ok := readCache(config.cache, config.targetAddr); ok {
			blocks = cached
			// Update block range to match cache contents
			if len(blocks) > 0 {
				config.startBlock, config.endBlock = getBlockRange(blocks)
				log.Info().
					Uint64("start", config.startBlock).
					Uint64("end", config.endBlock).
					Msg("Using block range from cache")
			}
		} else {
			// Cache miss or error - fetch from RPC and write cache
			blocks = fetchBlocks(ctx, config, client)
			if err := writeCache(config.cache, blocks); err != nil {
				log.Warn().Err(err).Msg("Failed to write cache")
			}
		}
	} else {
		blocks = fetchBlocks(ctx, config, client)
	}

	bm := computeBlocksMetadata(blocks)

	chartMetadata := txGasChartMetadata{
		chainID: chainID.Uint64(),

		targetAddr: config.targetAddr,
		startBlock: config.startBlock,
		endBlock:   config.endBlock,

		blocksMetadata: bm,

		renderer: config.renderer,

		outputPath: config.output,
	}

	return plotChart(chartMetadata)
}

// parseFlags parses the command-line flags and returns the corresponding txGasChartConfig.
func parseFlags(ctx context.Context, client *ethrpc.Client) (*txGasChartConfig, error) {
	config := &txGasChartConfig{}

	config.startBlock = inputArgs.startBlock
	config.endBlock = inputArgs.endBlock
	config.cache = inputArgs.cache

	// Skip default block range calculation if cache is specified - we'll use cache's range
	if config.cache == "" || inputArgs.startBlock != math.MaxUint64 || inputArgs.endBlock != math.MaxUint64 {
		h, err := util.HeaderByBlockNumber(ctx, client, nil)
		if err != nil {
			return nil, err
		}

		if config.endBlock == math.MaxUint64 || config.endBlock > h.Number.Uint64() {
			config.endBlock = h.Number.Uint64()
			log.Warn().Uint64("end_block", config.endBlock).Msg("End block not set or exceeds latest, using latest block")
		}

		const defaultBlockRange = 500

		if config.startBlock == math.MaxUint64 {
			if config.endBlock < defaultBlockRange {
				config.startBlock = 0
			} else {
				config.startBlock = config.endBlock - defaultBlockRange
			}

			log.Warn().Uint64("start_block", config.startBlock).Msg("Start block not set, using last 500 blocks")
		}

		if config.startBlock > config.endBlock {
			return nil, fmt.Errorf("start block %d cannot be greater than end block %d", config.startBlock, config.endBlock)
		}
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

	renderer := strings.ToLower(inputArgs.renderer)
	if renderer != "svg" && renderer != "canvas" {
		return nil, fmt.Errorf("invalid renderer %q: must be 'svg' or 'canvas'", inputArgs.renderer)
	}
	config.renderer = renderer

	return config, nil
}

// fetchBlocks fetches blocks in the specified range concurrently using the provided Ethereum client.
func fetchBlocks(ctx context.Context, config *txGasChartConfig, client *ethrpc.Client) []block {
	numBlocks := config.endBlock - config.startBlock + 1
	blocks := make([]block, numBlocks)
	offset := config.startBlock

	// Prepare worker pool
	workers := make(chan struct{}, config.concurrency)
	for i := uint64(0); i < config.concurrency; i++ {
		workers <- struct{}{}
	}

	log.Info().Msg("Reading blocks")

	var wg sync.WaitGroup
	for blockNumber := config.startBlock; blockNumber <= config.endBlock; blockNumber++ {
		wg.Add(1)
		go func(blockNumber uint64) {
			defer wg.Done()
			<-workers
			defer func() { workers <- struct{}{} }()

			b := fetchBlock(ctx, config, client, blockNumber)
			blocks[blockNumber-offset] = b
		}(blockNumber)
	}
	wg.Wait()

	return blocks
}

// fetchBlock fetches a single block with retry logic.
func fetchBlock(ctx context.Context, config *txGasChartConfig, client *ethrpc.Client, blockNumber uint64) block {
	for {
		log.Trace().Uint64("block_number", blockNumber).Msg("Processing block")
		if config.rateLimiter != nil {
			_ = config.rateLimiter.Wait(ctx)
		}

		blocksFromNetwork, err := util.GetBlockRange(ctx, blockNumber, blockNumber, client, false)
		if err != nil {
			log.Warn().Err(err).Uint64("block_number", blockNumber).Msg("Failed to fetch block, retrying")
			time.Sleep(time.Second)
			continue
		}

		var rawBlock rpctypes.RawBlockResponse
		if err := json.Unmarshal(*blocksFromNetwork[0], &rawBlock); err != nil {
			log.Error().Err(err).Uint64("block_number", blockNumber).Msg("Unable to unmarshal block")
			time.Sleep(time.Second)
			continue
		}

		parsedBlock := rpctypes.NewPolyBlock(&rawBlock)
		return parseBlock(ctx, parsedBlock, config.targetAddr)
	}
}

// parseBlock converts a parsed block to our block struct.
func parseBlock(ctx context.Context, parsedBlock rpctypes.PolyBlock, targetAddr string) block {
	txs := parsedBlock.Transactions()

	b := block{
		Hash:     parsedBlock.Hash(),
		Number:   parsedBlock.Number().Uint64(),
		GasLimit: parsedBlock.GasLimit(),
		GasUsed:  parsedBlock.GasUsed(),
		Txs:      make([]transaction, len(txs)),
	}

	var totalGasLimit uint64
	for i, tx := range txs {
		from, err := util.GetSenderFromTx(ctx, tx)
		if err != nil {
			log.Error().Err(err).Uint64("block", b.Number).Stringer("txHash", tx.Hash()).Msg("Unable to get sender from tx, skipping")
			continue
		}

		to := tx.To()
		target := isTargetTx(from, to, targetAddr)
		gasPrice := tx.GasPrice().Uint64()
		gasLimit := tx.Gas()

		b.Txs[i] = transaction{
			Hash:     tx.Hash(),
			From:     from,
			To:       to,
			GasPrice: gasPrice,
			GasLimit: gasLimit,
			Target:   target,
		}

		totalGasLimit += gasLimit

		if target {
			log.Info().
				Uint64("block", b.Number).
				Stringer("txHash", tx.Hash()).
				Uint64("gas_price_wei", gasPrice).
				Uint64("gas_limit", gasLimit).
				Msg("target tx found")
		}
	}

	b.TxsGasLimit = totalGasLimit

	return b
}
