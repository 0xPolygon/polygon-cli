package report

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"os"
	"slices"
	"sort"
	"strconv"
	"sync"
	"time"

	_ "embed"

	"github.com/0xPolygon/polygon-cli/util"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/time/rate"
)

const (
	// DefaultBlockRange is the default number of blocks to analyze when start/end blocks are not specified
	DefaultBlockRange = 500
	// BlockNotSet is a sentinel value to indicate a block number flag was not set by the user
	BlockNotSet = math.MaxUint64
)

type (
	reportParams struct {
		RpcUrl      string
		StartBlock  uint64
		EndBlock    uint64
		OutputFile  string
		Format      string
		Concurrency int
		RateLimit   float64
	}
)

var (
	//go:embed usage.md
	usage       string
	inputReport reportParams = reportParams{}
)

// ReportCmd represents the report command
var ReportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate a report analyzing a range of blocks from an Ethereum-compatible blockchain.",
	Long:  usage,
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return checkFlags()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		// Connect to RPC
		ec, err := ethrpc.DialContext(ctx, inputReport.RpcUrl)
		if err != nil {
			return fmt.Errorf("failed to connect to RPC: %w", err)
		}
		defer ec.Close()

		// Fetch chain ID
		var chainIDHex string
		err = ec.CallContext(ctx, &chainIDHex, "eth_chainId")
		if err != nil {
			return fmt.Errorf("failed to fetch chain ID: %w", err)
		}
		chainID := hexToUint64(chainIDHex)

		// Determine block range with smart defaults
		startBlock := inputReport.StartBlock
		endBlock := inputReport.EndBlock

		// Fetch latest block if needed for auto-detection
		var latestBlock uint64
		needsLatest := startBlock == BlockNotSet || endBlock == BlockNotSet
		if needsLatest {
			var latestBlockHex string
			err = ec.CallContext(ctx, &latestBlockHex, "eth_blockNumber")
			if err != nil {
				return fmt.Errorf("failed to fetch latest block number: %w", err)
			}
			latestBlock = hexToUint64(latestBlockHex)
			log.Info().Uint64("latest-block", latestBlock).Msg("Auto-detected latest block")
		}

		// Apply smart defaults based on which flags were set
		if startBlock == BlockNotSet && endBlock == BlockNotSet {
			// Both unspecified: analyze latest DefaultBlockRange blocks
			endBlock = latestBlock
			if latestBlock >= DefaultBlockRange-1 {
				startBlock = latestBlock - (DefaultBlockRange - 1)
			} else {
				startBlock = 0
			}
		} else if startBlock == BlockNotSet {
			// Only start-block unspecified: analyze previous DefaultBlockRange blocks from end-block
			if endBlock >= DefaultBlockRange-1 {
				startBlock = endBlock - (DefaultBlockRange - 1)
			} else {
				startBlock = 0
			}
		} else if endBlock == BlockNotSet {
			// Only end-block unspecified: analyze next DefaultBlockRange blocks from start-block
			// But don't exceed the latest block
			endBlock = startBlock + (DefaultBlockRange - 1)
			if endBlock > latestBlock {
				endBlock = latestBlock
			}
		}
		// If both are set by user (including 0,0), use them as-is

		log.Info().
			Str("rpc-url", inputReport.RpcUrl).
			Uint64("start-block", startBlock).
			Uint64("end-block", endBlock).
			Msg("Starting block analysis")

		// Initialize the report
		report := &BlockReport{
			ChainID:     chainID,
			RpcUrl:      inputReport.RpcUrl,
			StartBlock:  startBlock,
			EndBlock:    endBlock,
			GeneratedAt: time.Now(),
			Blocks:      []BlockInfo{},
		}

		// Generate the report
		err = generateReport(ctx, ec, report, inputReport.Concurrency, inputReport.RateLimit)
		if err != nil {
			return fmt.Errorf("failed to generate report: %w", err)
		}

		// Output the report
		if err := outputReport(report, inputReport.Format, inputReport.OutputFile); err != nil {
			return fmt.Errorf("failed to output report: %w", err)
		}

		log.Info().Msg("Report generation completed")
		return nil
	},
}

func init() {
	f := ReportCmd.Flags()
	f.StringVar(&inputReport.RpcUrl, "rpc-url", "http://localhost:8545", "RPC endpoint URL")
	f.Uint64Var(&inputReport.StartBlock, "start-block", BlockNotSet, "starting block number (default: auto-detect based on end-block or latest)")
	f.Uint64Var(&inputReport.EndBlock, "end-block", BlockNotSet, "ending block number (default: auto-detect based on start-block or latest)")
	f.StringVarP(&inputReport.OutputFile, "output", "o", "", "output file path (default: stdout for JSON, report.html for HTML, report.pdf for PDF)")
	f.StringVarP(&inputReport.Format, "format", "f", "json", "output format [json, html, pdf]")
	f.IntVar(&inputReport.Concurrency, "concurrency", 10, "number of concurrent RPC requests")
	f.Float64Var(&inputReport.RateLimit, "rate-limit", 4, "requests per second limit")
}

func checkFlags() error {
	// Validate RPC URL
	if err := util.ValidateUrl(inputReport.RpcUrl); err != nil {
		return err
	}

	// Validate block range only if both are explicitly specified by the user
	if inputReport.StartBlock != BlockNotSet && inputReport.EndBlock != BlockNotSet {
		if inputReport.EndBlock < inputReport.StartBlock {
			return fmt.Errorf("end-block must be greater than or equal to start-block")
		}
	}

	// Validate format
	if inputReport.Format != "json" && inputReport.Format != "html" && inputReport.Format != "pdf" {
		return fmt.Errorf("format must be either 'json', 'html', or 'pdf'")
	}

	// Set default output file for HTML if not specified
	if inputReport.Format == "html" && inputReport.OutputFile == "" {
		inputReport.OutputFile = "report.html"
	}

	// Set default output file for PDF if not specified
	if inputReport.Format == "pdf" && inputReport.OutputFile == "" {
		inputReport.OutputFile = "report.pdf"
	}

	return nil
}

// generateReport analyzes the block range and generates a report
func generateReport(ctx context.Context, ec *ethrpc.Client, report *BlockReport, concurrency int, rateLimit float64) error {
	log.Info().Msg("Fetching and analyzing blocks")

	// Create a cancellable context for workers
	workerCtx, cancelWorkers := context.WithCancel(ctx)
	defer cancelWorkers() // Ensure workers are stopped when function returns

	// Create rate limiter
	rateLimiter := rate.NewLimiter(rate.Limit(rateLimit), 1)

	totalBlocks := report.EndBlock - report.StartBlock + 1
	blockChan := make(chan uint64, totalBlocks)
	resultChan := make(chan *BlockInfo, concurrency)
	errorChan := make(chan error, 1) // Only need space for one error

	// Fill the block channel with block numbers to fetch
	for blockNum := report.StartBlock; blockNum <= report.EndBlock; blockNum++ {
		blockChan <- blockNum
	}
	close(blockChan)

	// Start worker goroutines
	var wg sync.WaitGroup
	for range concurrency {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for blockNum := range blockChan {
				// Check if worker context is canceled
				if workerCtx.Err() != nil {
					return
				}

				blockInfo, err := fetchBlockInfo(workerCtx, ec, blockNum, rateLimiter)
				if err != nil {
					// Check for context cancellation errors (user interrupt or internal cancellation)
					if workerCtx.Err() != nil {
						return
					}
					log.Warn().Err(err).Uint64("block", blockNum).Msg("Failed to fetch block, skipping")
					continue
				}

				// Send result with context check to avoid blocking
				select {
				case resultChan <- blockInfo:
				case <-workerCtx.Done():
					return
				}
			}
		}()
	}

	// Close result channel when all workers are done
	go func() {
		wg.Wait()
		close(resultChan)
		close(errorChan)
	}()

	// Collect results
	totalTxCount := uint64(0)
	totalGasUsed := uint64(0)
	totalBaseFee := big.NewInt(0)
	blockCount := uint64(0)
	uniqueSenders := make(map[string]bool)
	uniqueRecipients := make(map[string]bool)
	processedBlocks := uint64(0)

	// Process results and check for context cancellation
	for {
		select {
		case blockInfo, ok := <-resultChan:
			if !ok {
				// Channel closed, all results processed
				goto done
			}
			report.Blocks = append(report.Blocks, *blockInfo)
			totalTxCount += blockInfo.TxCount
			totalGasUsed += blockInfo.GasUsed
			if blockInfo.BaseFeePerGas != nil {
				totalBaseFee.Add(totalBaseFee, blockInfo.BaseFeePerGas)
			}
			blockCount++

			// Track unique addresses
			for _, tx := range blockInfo.Transactions {
				if tx.From != "" {
					uniqueSenders[tx.From] = true
				}
				if tx.To != "" {
					uniqueRecipients[tx.To] = true
				}
			}

			processedBlocks++
			if processedBlocks%100 == 0 || processedBlocks == totalBlocks {
				log.Info().Uint64("progress", processedBlocks).Uint64("total", totalBlocks).Msg("Progress")
			}
		case <-ctx.Done():
			// Parent context canceled (e.g., user pressed Ctrl+C)
			// cancelWorkers() will be called by defer to stop all workers
			return ctx.Err()
		}
	}
done:
	// Sort blocks by block number to ensure correct ordering for charts and analysis
	slices.SortFunc(report.Blocks, func(a, b BlockInfo) int {
		if a.Number < b.Number {
			return -1
		} else if a.Number > b.Number {
			return 1
		}
		return 0
	})

	// Warn if no blocks were successfully fetched
	if len(report.Blocks) == 0 {
		log.Warn().Msg("No blocks were successfully fetched. Report will be empty.")
	}

	// Calculate summary statistics
	report.Summary = SummaryStats{
		TotalBlocks:       blockCount,
		TotalTransactions: totalTxCount,
		TotalGasUsed:      totalGasUsed,
		UniqueSenders:     uint64(len(uniqueSenders)),
		UniqueRecipients:  uint64(len(uniqueRecipients)),
	}

	if blockCount > 0 {
		report.Summary.AvgTxPerBlock = float64(totalTxCount) / float64(blockCount)
		report.Summary.AvgGasPerBlock = float64(totalGasUsed) / float64(blockCount)
		if totalBaseFee.Cmp(big.NewInt(0)) > 0 {
			avgBaseFee := new(big.Int).Div(totalBaseFee, big.NewInt(int64(blockCount)))
			report.Summary.AvgBaseFeePerGas = avgBaseFee.String()
		}
	}

	// Calculate top 10 statistics
	report.Top10 = calculateTop10Stats(report.Blocks)

	return nil
}

// fetchBlockInfo retrieves information about a specific block and its transactions
func fetchBlockInfo(ctx context.Context, ec *ethrpc.Client, blockNum uint64, rateLimiter *rate.Limiter) (*BlockInfo, error) {
	// Wait for rate limiter before making RPC call
	if err := rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiter error: %w", err)
	}

	var result map[string]any
	err := ec.CallContext(ctx, &result, "eth_getBlockByNumber", fmt.Sprintf("0x%x", blockNum), true)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, fmt.Errorf("block not found")
	}

	blockInfo := &BlockInfo{
		Number:       blockNum,
		Timestamp:    hexToUint64(result["timestamp"]),
		GasUsed:      hexToUint64(result["gasUsed"]),
		GasLimit:     hexToUint64(result["gasLimit"]),
		Transactions: []TransactionInfo{},
	}

	// Parse base fee if present (EIP-1559)
	if baseFee, ok := result["baseFeePerGas"].(string); ok && baseFee != "" {
		bf := new(big.Int)
		// Remove "0x" prefix if present
		if len(baseFee) > 2 && baseFee[:2] == "0x" {
			baseFee = baseFee[2:]
		}
		if _, success := bf.SetString(baseFee, 16); success {
			blockInfo.BaseFeePerGas = bf
		}
	}

	// Process transactions
	if txs, ok := result["transactions"].([]any); ok {
		blockInfo.TxCount = uint64(len(txs))

		// Fetch all receipts for this block in a single call
		if len(txs) > 0 {
			// Wait for rate limiter before making RPC call
			if err := rateLimiter.Wait(ctx); err != nil {
				return nil, fmt.Errorf("rate limiter error: %w", err)
			}

			var receipts []map[string]any
			err := ec.CallContext(ctx, &receipts, "eth_getBlockReceipts", fmt.Sprintf("0x%x", blockNum))
			if err != nil {
				return nil, fmt.Errorf("failed to fetch block receipts: %w", err)
			}

			if len(receipts) != len(txs) {
				return nil, fmt.Errorf("mismatch between transactions (%d) and receipts (%d)", len(txs), len(receipts))
			}

			// Process each transaction with its corresponding receipt
			for i, txData := range txs {
				txMap, ok := txData.(map[string]any)
				if !ok {
					continue
				}

				txHash, _ := txMap["hash"].(string)
				from, _ := txMap["from"].(string)
				to, _ := txMap["to"].(string)
				gasPrice := hexToUint64(txMap["gasPrice"])
				gasLimit := hexToUint64(txMap["gas"])

				receipt := receipts[i]
				gasUsed := hexToUint64(receipt["gasUsed"])
				gasUsedPercent := 0.0
				if blockInfo.GasLimit > 0 {
					gasUsedPercent = (float64(gasUsed) / float64(blockInfo.GasLimit)) * 100
				}

				txInfo := TransactionInfo{
					Hash:           txHash,
					From:           from,
					To:             to,
					BlockNumber:    blockNum,
					GasUsed:        gasUsed,
					GasLimit:       gasLimit,
					GasPrice:       gasPrice,
					BlockGasLimit:  blockInfo.GasLimit,
					GasUsedPercent: gasUsedPercent,
				}

				blockInfo.Transactions = append(blockInfo.Transactions, txInfo)
			}
		}
	}

	return blockInfo, nil
}

// hexToUint64 converts a hex string to uint64
func hexToUint64(v any) uint64 {
	if v == nil {
		return 0
	}
	str, ok := v.(string)
	if !ok {
		return 0
	}
	if len(str) > 2 && str[:2] == "0x" {
		str = str[2:]
	}
	val, _ := strconv.ParseUint(str, 16, 64)
	return val
}

// calculateTop10Stats calculates various top 10 lists from the collected blocks
func calculateTop10Stats(blocks []BlockInfo) Top10Stats {
	top10 := Top10Stats{}

	// Top 10 blocks by transaction count
	blocksByTxCount := make([]TopBlock, len(blocks))
	for i, block := range blocks {
		blocksByTxCount[i] = TopBlock{
			Number:  block.Number,
			TxCount: block.TxCount,
		}
	}
	// Sort by tx count descending
	sort.Slice(blocksByTxCount, func(i, j int) bool {
		return blocksByTxCount[i].TxCount > blocksByTxCount[j].TxCount
	})
	if len(blocksByTxCount) > 10 {
		top10.BlocksByTxCount = blocksByTxCount[:10]
	} else {
		top10.BlocksByTxCount = blocksByTxCount
	}

	// Top 10 blocks by gas used
	blocksByGasUsed := make([]TopBlock, len(blocks))
	for i, block := range blocks {
		gasUsedPercent := 0.0
		if block.GasLimit > 0 {
			gasUsedPercent = (float64(block.GasUsed) / float64(block.GasLimit)) * 100
		}
		blocksByGasUsed[i] = TopBlock{
			Number:         block.Number,
			GasUsed:        block.GasUsed,
			GasLimit:       block.GasLimit,
			GasUsedPercent: gasUsedPercent,
		}
	}
	// Sort by gas used descending
	sort.Slice(blocksByGasUsed, func(i, j int) bool {
		return blocksByGasUsed[i].GasUsed > blocksByGasUsed[j].GasUsed
	})
	if len(blocksByGasUsed) > 10 {
		top10.BlocksByGasUsed = blocksByGasUsed[:10]
	} else {
		top10.BlocksByGasUsed = blocksByGasUsed
	}

	// Collect all transactions and track gas prices and gas limits
	var allTxsByGasUsed []TopTransaction
	var allTxsByGasLimit []TopTransaction
	gasPriceMap := make(map[uint64]uint64)
	gasLimitMap := make(map[uint64]uint64)

	for _, block := range blocks {
		for _, tx := range block.Transactions {
			allTxsByGasUsed = append(allTxsByGasUsed, TopTransaction{
				Hash:           tx.Hash,
				BlockNumber:    tx.BlockNumber,
				GasLimit:       tx.GasLimit,
				GasUsed:        tx.GasUsed,
				BlockGasLimit:  tx.BlockGasLimit,
				GasUsedPercent: tx.GasUsedPercent,
			})
			allTxsByGasLimit = append(allTxsByGasLimit, TopTransaction{
				Hash:        tx.Hash,
				BlockNumber: tx.BlockNumber,
				GasLimit:    tx.GasLimit,
				GasUsed:     tx.GasUsed,
			})
			gasPriceMap[tx.GasPrice]++
			gasLimitMap[tx.GasLimit]++
		}
	}

	// Top 10 transactions by gas used
	// Sort transactions by gas used descending
	sort.Slice(allTxsByGasUsed, func(i, j int) bool {
		return allTxsByGasUsed[i].GasUsed > allTxsByGasUsed[j].GasUsed
	})
	if len(allTxsByGasUsed) > 10 {
		top10.TransactionsByGas = allTxsByGasUsed[:10]
	} else {
		top10.TransactionsByGas = allTxsByGasUsed
	}

	// Top 10 transactions by gas limit
	// Sort transactions by gas limit descending
	sort.Slice(allTxsByGasLimit, func(i, j int) bool {
		return allTxsByGasLimit[i].GasLimit > allTxsByGasLimit[j].GasLimit
	})
	if len(allTxsByGasLimit) > 10 {
		top10.TransactionsByGasLimit = allTxsByGasLimit[:10]
	} else {
		top10.TransactionsByGasLimit = allTxsByGasLimit
	}

	// Top 10 most used gas prices
	gasPriceFreqs := make([]GasPriceFreq, 0, len(gasPriceMap))
	for price, count := range gasPriceMap {
		gasPriceFreqs = append(gasPriceFreqs, GasPriceFreq{
			GasPrice: price,
			Count:    count,
		})
	}
	// Sort by count descending
	sort.Slice(gasPriceFreqs, func(i, j int) bool {
		return gasPriceFreqs[i].Count > gasPriceFreqs[j].Count
	})
	if len(gasPriceFreqs) > 10 {
		top10.MostUsedGasPrices = gasPriceFreqs[:10]
	} else {
		top10.MostUsedGasPrices = gasPriceFreqs
	}

	// Top 10 most used gas limits
	gasLimitFreqs := make([]GasLimitFreq, 0, len(gasLimitMap))
	for limit, count := range gasLimitMap {
		gasLimitFreqs = append(gasLimitFreqs, GasLimitFreq{
			GasLimit: limit,
			Count:    count,
		})
	}
	// Sort by count descending
	sort.Slice(gasLimitFreqs, func(i, j int) bool {
		return gasLimitFreqs[i].Count > gasLimitFreqs[j].Count
	})
	if len(gasLimitFreqs) > 10 {
		top10.MostUsedGasLimits = gasLimitFreqs[:10]
	} else {
		top10.MostUsedGasLimits = gasLimitFreqs
	}

	return top10
}

// outputReport writes the report to the specified output
func outputReport(report *BlockReport, format, outputFile string) error {
	switch format {
	case "json":
		return outputJSON(report, outputFile)
	case "html":
		return outputHTML(report, outputFile)
	case "pdf":
		return outputPDF(report, outputFile)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

// outputJSON writes the report as JSON
func outputJSON(report *BlockReport, outputFile string) error {
	jsonData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report to JSON: %w", err)
	}

	if outputFile == "" {
		fmt.Println(string(jsonData))
		return nil
	}

	if err := os.WriteFile(outputFile, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	log.Info().Str("file", outputFile).Msg("JSON report written")
	return nil
}

// outputHTML generates an HTML report from the JSON data
func outputHTML(report *BlockReport, outputFile string) error {
	html := generateHTML(report)

	if err := os.WriteFile(outputFile, []byte(html), 0644); err != nil {
		return fmt.Errorf("failed to write HTML file: %w", err)
	}

	log.Info().Str("file", outputFile).Msg("HTML report written")
	return nil
}
