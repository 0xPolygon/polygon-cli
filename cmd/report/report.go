package report

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"time"

	_ "embed"

	"github.com/0xPolygon/polygon-cli/util"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type (
	reportParams struct {
		RpcUrl     string
		StartBlock uint64
		EndBlock   uint64
		OutputFile string
		Format     string
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

		log.Info().
			Str("rpc-url", inputReport.RpcUrl).
			Uint64("start-block", inputReport.StartBlock).
			Uint64("end-block", inputReport.EndBlock).
			Msg("Starting block analysis")

		// Fetch chain ID
		var chainIDHex string
		err = ec.CallContext(ctx, &chainIDHex, "eth_chainId")
		if err != nil {
			return fmt.Errorf("failed to fetch chain ID: %w", err)
		}
		chainID := hexToUint64(chainIDHex)

		// Initialize the report
		report := &BlockReport{
			ChainID:     chainID,
			RpcUrl:      inputReport.RpcUrl,
			StartBlock:  inputReport.StartBlock,
			EndBlock:    inputReport.EndBlock,
			GeneratedAt: time.Now(),
			Blocks:      []BlockInfo{},
		}

		// Generate the report
		err = generateReport(ctx, ec, report)
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
	f.Uint64Var(&inputReport.StartBlock, "start-block", 0, "starting block number for analysis")
	f.Uint64Var(&inputReport.EndBlock, "end-block", 0, "ending block number for analysis")
	f.StringVarP(&inputReport.OutputFile, "output", "o", "", "output file path (default: stdout for JSON, report.html for HTML)")
	f.StringVarP(&inputReport.Format, "format", "f", "json", "output format [json, html]")
}

func checkFlags() error {
	// Validate RPC URL
	if err := util.ValidateUrl(inputReport.RpcUrl); err != nil {
		return err
	}

	// Validate block range
	if inputReport.EndBlock < inputReport.StartBlock {
		return fmt.Errorf("end-block must be greater than or equal to start-block")
	}

	// Validate format
	if inputReport.Format != "json" && inputReport.Format != "html" {
		return fmt.Errorf("format must be either 'json' or 'html'")
	}

	// Set default output file for HTML if not specified
	if inputReport.Format == "html" && inputReport.OutputFile == "" {
		inputReport.OutputFile = "report.html"
	}

	return nil
}

// generateReport analyzes the block range and generates a report
func generateReport(ctx context.Context, ec *ethrpc.Client, report *BlockReport) error {
	log.Info().Msg("Fetching and analyzing blocks")

	totalTxCount := uint64(0)
	totalGasUsed := uint64(0)
	totalBaseFee := big.NewInt(0)
	blockCount := uint64(0)
	uniqueSenders := make(map[string]bool)
	uniqueRecipients := make(map[string]bool)

	// Fetch blocks in the range
	for blockNum := report.StartBlock; blockNum <= report.EndBlock; blockNum++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		blockInfo, err := fetchBlockInfo(ctx, ec, blockNum)
		if err != nil {
			log.Warn().Err(err).Uint64("block", blockNum).Msg("Failed to fetch block, skipping")
			continue
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

		if blockNum%100 == 0 || blockNum == report.EndBlock {
			log.Info().Uint64("block", blockNum).Uint64("progress", blockNum-report.StartBlock+1).Uint64("total", report.EndBlock-report.StartBlock+1).Msg("Progress")
		}
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
			report.Summary.AvgBaseFeePerGas = avgBaseFee.Uint64()
		}
	}

	// Calculate top 10 statistics
	report.Top10 = calculateTop10Stats(report.Blocks)

	return nil
}

// fetchBlockInfo retrieves information about a specific block and its transactions
func fetchBlockInfo(ctx context.Context, ec *ethrpc.Client, blockNum uint64) (*BlockInfo, error) {
	var result map[string]any
	err := ec.CallContext(ctx, &result, "eth_getBlockByNumber", fmt.Sprintf("0x%x", blockNum), true)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, fmt.Errorf("block not found")
	}

	blockInfo := &BlockInfo{
		Number:    blockNum,
		Timestamp: hexToUint64(result["timestamp"]),
		GasUsed:   hexToUint64(result["gasUsed"]),
		GasLimit:  hexToUint64(result["gasLimit"]),
		Transactions: []TransactionInfo{},
	}

	// Parse base fee if present (EIP-1559)
	if baseFee, ok := result["baseFeePerGas"].(string); ok {
		bf := new(big.Int)
		bf.SetString(baseFee[2:], 16) // Remove "0x" prefix
		blockInfo.BaseFeePerGas = bf
	}

	// Process transactions
	if txs, ok := result["transactions"].([]any); ok {
		blockInfo.TxCount = uint64(len(txs))

		// Fetch transaction receipts to get actual gas used
		for _, txData := range txs {
			txMap, ok := txData.(map[string]any)
			if !ok {
				continue
			}

			txHash, _ := txMap["hash"].(string)
			from, _ := txMap["from"].(string)
			to, _ := txMap["to"].(string)
			gasPrice := hexToUint64(txMap["gasPrice"])
			gasLimit := hexToUint64(txMap["gas"])

			// Fetch transaction receipt for gas used
			var receipt map[string]any
			err := ec.CallContext(ctx, &receipt, "eth_getTransactionReceipt", txHash)
			if err != nil || receipt == nil {
				continue
			}

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
	for i := 0; i < len(blocksByTxCount)-1; i++ {
		for j := i + 1; j < len(blocksByTxCount); j++ {
			if blocksByTxCount[j].TxCount > blocksByTxCount[i].TxCount {
				blocksByTxCount[i], blocksByTxCount[j] = blocksByTxCount[j], blocksByTxCount[i]
			}
		}
	}
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
	for i := 0; i < len(blocksByGasUsed)-1; i++ {
		for j := i + 1; j < len(blocksByGasUsed); j++ {
			if blocksByGasUsed[j].GasUsed > blocksByGasUsed[i].GasUsed {
				blocksByGasUsed[i], blocksByGasUsed[j] = blocksByGasUsed[j], blocksByGasUsed[i]
			}
		}
	}
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
	for i := 0; i < len(allTxsByGasUsed)-1; i++ {
		for j := i + 1; j < len(allTxsByGasUsed); j++ {
			if allTxsByGasUsed[j].GasUsed > allTxsByGasUsed[i].GasUsed {
				allTxsByGasUsed[i], allTxsByGasUsed[j] = allTxsByGasUsed[j], allTxsByGasUsed[i]
			}
		}
	}
	if len(allTxsByGasUsed) > 10 {
		top10.TransactionsByGas = allTxsByGasUsed[:10]
	} else {
		top10.TransactionsByGas = allTxsByGasUsed
	}

	// Top 10 transactions by gas limit
	// Sort transactions by gas limit descending
	for i := 0; i < len(allTxsByGasLimit)-1; i++ {
		for j := i + 1; j < len(allTxsByGasLimit); j++ {
			if allTxsByGasLimit[j].GasLimit > allTxsByGasLimit[i].GasLimit {
				allTxsByGasLimit[i], allTxsByGasLimit[j] = allTxsByGasLimit[j], allTxsByGasLimit[i]
			}
		}
	}
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
	for i := 0; i < len(gasPriceFreqs)-1; i++ {
		for j := i + 1; j < len(gasPriceFreqs); j++ {
			if gasPriceFreqs[j].Count > gasPriceFreqs[i].Count {
				gasPriceFreqs[i], gasPriceFreqs[j] = gasPriceFreqs[j], gasPriceFreqs[i]
			}
		}
	}
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
	for i := 0; i < len(gasLimitFreqs)-1; i++ {
		for j := i + 1; j < len(gasLimitFreqs); j++ {
			if gasLimitFreqs[j].Count > gasLimitFreqs[i].Count {
				gasLimitFreqs[i], gasLimitFreqs[j] = gasLimitFreqs[j], gasLimitFreqs[i]
			}
		}
	}
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
