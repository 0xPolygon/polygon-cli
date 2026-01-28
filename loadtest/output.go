package loadtest

import (
	"context"
	"encoding/json"
	"math"
	"math/big"
	"slices"
	"sync"
	"time"

	"github.com/montanaflynn/stats"

	"golang.org/x/time/rate"

	"github.com/0xPolygon/polygon-cli/loadtest/config"
	"github.com/0xPolygon/polygon-cli/rpctypes"
	"github.com/0xPolygon/polygon-cli/util"
	"golang.org/x/exp/constraints"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// SummarizeResults handles the post-load-test summarization.
func SummarizeResults(ctx context.Context, c *ethclient.Client, rpc *ethrpc.Client, cfg *config.Config, ap *AccountPool, results []Sample, startBlockNumber, lastBlockNumber uint64) error {
	var err error

	log.Trace().Msg("Starting block range capture")
	// confirm start block number is ok
	_, err = c.BlockByNumber(ctx, new(big.Int).SetUint64(startBlockNumber))
	if err != nil {
		return err
	}
	rawBlocks, err := util.GetBlockRange(ctx, startBlockNumber, lastBlockNumber, rpc, false)
	if err != nil {
		return err
	}

	batchSize := cfg.BatchSize
	goRoutineLimit := cfg.Concurrency
	var txGroup sync.WaitGroup
	threadPool := make(chan bool, goRoutineLimit)
	log.Trace().Msg("Starting tx receipt capture")
	rawTxReceipts := make([]*json.RawMessage, 0)
	var rawTxReceiptsLock sync.Mutex
	var txGroupErr error

	startReceipt := time.Now()
	for k := range rawBlocks {
		threadPool <- true
		txGroup.Add(1)
		go func(b *json.RawMessage) {
			var receipt []*json.RawMessage
			receipt, err = util.GetReceipts(ctx, []*json.RawMessage{b}, rpc, batchSize)
			if err != nil {
				txGroupErr = err
				return
			}
			rawTxReceiptsLock.Lock()
			rawTxReceipts = append(rawTxReceipts, receipt...)
			rawTxReceiptsLock.Unlock()
			<-threadPool
			txGroup.Done()
		}(rawBlocks[k])
	}

	endReceipt := time.Now()
	txGroup.Wait()
	if txGroupErr != nil {
		log.Error().Err(err).Msg("One of the threads fetching tx receipts failed")
		return err
	}

	blocks := make([]rpctypes.RawBlockResponse, 0)
	for _, b := range rawBlocks {
		var block rpctypes.RawBlockResponse
		err = json.Unmarshal(*b, &block)
		if err != nil {
			log.Error().Err(err).Msg("Error decoding block response")
			return err
		}
		blocks = append(blocks, block)
	}
	log.Info().Int("len", len(blocks)).Msg("Block summary")

	txReceipts := make([]rpctypes.RawTxReceipt, 0)
	log.Trace().Int("len", len(rawTxReceipts)).Msg("Raw receipts")
	for _, r := range rawTxReceipts {
		if isEmptyJSONResponse(r) {
			continue
		}
		var receipt rpctypes.RawTxReceipt
		err = json.Unmarshal(*r, &receipt)
		if err != nil {
			log.Error().Err(err).Msg("Error decoding tx receipt response")
			return err
		}
		txReceipts = append(txReceipts, receipt)
	}
	log.Info().Int("len", len(txReceipts)).Msg("Receipt summary")

	blockData := make(map[uint64]BlockSummary)
	for k, b := range blocks {
		bs := BlockSummary{}
		bs.Block = &blocks[k]
		bs.Receipts = make(map[ethcommon.Hash]rpctypes.RawTxReceipt)
		bs.Latencies = make(map[uint64]time.Duration)
		blockData[b.Number.ToUint64()] = bs
	}

	for _, r := range txReceipts {
		bn := r.BlockNumber.ToUint64()
		bs := blockData[bn]
		if bs.Receipts == nil {
			log.Error().Uint64("blocknumber", bn).Msg("Block number from receipts does not exist in block data")
		}
		bs.Receipts[r.TransactionHash.ToHash()] = r
		blockData[bn] = bs
	}

	nonceTimes := make(map[uint64]time.Time)
	for _, ltr := range results {
		nonceTimes[ltr.Nonce] = ltr.RequestTime
	}

	minLatency := time.Millisecond * 100
	for _, bs := range blockData {
		for _, tx := range bs.Block.Transactions {
			mineTime := time.Unix(bs.Block.Timestamp.ToInt64(), 0)
			requestTime := nonceTimes[tx.Nonce.ToUint64()]
			txLatency := mineTime.Sub(requestTime)
			if txLatency.Hours() > 2 {
				log.Debug().
					Float64("txHours", txLatency.Hours()).
					Uint64("nonce", tx.Nonce.ToUint64()).
					Uint64("blockNumber", bs.Block.Number.ToUint64()).
					Time("mineTime", mineTime).
					Time("requestTime", requestTime).
					Msg("Encountered transaction with more than 2 hours latency")
			}
			bs.Latencies[tx.Nonce.ToUint64()] = txLatency

			if txLatency < minLatency {
				minLatency = txLatency
			}
		}
	}
	// Adjust latencies for time discrepancies
	if minLatency < time.Millisecond*100 {
		log.Trace().Str("minLatency", minLatency.String()).Msg("Minimum latency is below expected threshold")
		shiftSize := ((time.Millisecond * 100) - minLatency) + time.Millisecond + 100
		for _, bs := range blockData {
			for _, tx := range bs.Block.Transactions {
				bs.Latencies[tx.Nonce.ToUint64()] += shiftSize
			}
		}
	}

	printBlockSummary(c, cfg, ap, blockData)

	log.Trace().Str("summaryTime", (endReceipt.Sub(startReceipt)).String()).Msg("Total Summary Time")

	return nil
}

func printBlockSummary(c *ethclient.Client, cfg *config.Config, ap *AccountPool, bs map[uint64]BlockSummary) {
	filterBlockSummary(ap, bs)
	mapKeys := getSortedMapKeys(bs)
	if len(mapKeys) == 0 {
		return
	}

	var totalTransactions uint64 = 0
	var totalGasUsed uint64 = 0
	p := message.NewPrinter(language.English)

	allLatencies := make([]time.Duration, 0)
	summaryOutputMode := cfg.SummaryOutputMode
	jsonSummaryList := []Summary{}
	for _, v := range mapKeys {
		summary := bs[v]
		gasUsed := getTotalGasUsed(summary.Receipts)
		blockLatencies := getMapValues(summary.Latencies)
		minLatency, medianLatency, maxLatency := getMinMedianMax(blockLatencies)
		allLatencies = append(allLatencies, blockLatencies...)
		blockUtilization := float64(gasUsed) / summary.Block.GasLimit.ToFloat64()
		if gasUsed == 0 {
			blockUtilization = 0
		}
		// if we're at trace, debug, or info level we'll output the block level metrics
		if zerolog.GlobalLevel() <= zerolog.InfoLevel {
			switch summaryOutputMode {
			case "text":
				_, _ = p.Printf("Block number: %v\tTime: %s\tGas Limit: %v\tGas Used: %v\tNum Tx: %v\tUtilization %v\tLatencies: %v\t%v\t%v\n",
					number.Decimal(summary.Block.Number.ToUint64()),
					time.Unix(summary.Block.Timestamp.ToInt64(), 0),
					number.Decimal(summary.Block.GasLimit.ToUint64()),
					number.Decimal(gasUsed),
					number.Decimal(len(summary.Block.Transactions)),
					number.Percent(blockUtilization),
					number.Decimal(minLatency.Seconds()),
					number.Decimal(medianLatency.Seconds()),
					number.Decimal(maxLatency.Seconds()))
			case "json":
				jsonSummary := Summary{}
				jsonSummary.BlockNumber = summary.Block.Number.ToUint64()
				jsonSummary.Time = time.Unix(summary.Block.Timestamp.ToInt64(), 0)
				jsonSummary.GasLimit = summary.Block.GasLimit.ToUint64()
				jsonSummary.GasUsed = gasUsed
				jsonSummary.NumTx = len(summary.Block.Transactions)
				jsonSummary.Utilization = blockUtilization
				latencies := Latency{}
				latencies.Min = minLatency.Seconds()
				latencies.Median = medianLatency.Seconds()
				latencies.Max = maxLatency.Seconds()
				jsonSummary.Latencies = latencies
				jsonSummaryList = append(jsonSummaryList, jsonSummary)
			default:
				log.Error().Str("mode", summaryOutputMode).Msg("Invalid mode for summary output")
			}
		}
		totalTransactions += uint64(len(summary.Block.Transactions))
		totalGasUsed += gasUsed
	}
	parentOfFirstBlock, _ := c.BlockByNumber(context.Background(), big.NewInt(bs[mapKeys[0]].Block.Number.ToInt64()-1))
	lastBlock := bs[mapKeys[len(mapKeys)-1]].Block
	totalMiningTime := time.Duration(lastBlock.Timestamp.ToUint64()-parentOfFirstBlock.Time()) * time.Second
	tps := float64(totalTransactions) / totalMiningTime.Seconds()
	gaspersec := float64(totalGasUsed) / totalMiningTime.Seconds()
	minLatency, medianLatency, maxLatency := getMinMedianMax(allLatencies)
	successfulTx, totalTx := getSuccessfulTransactionCount(bs)
	meanBlocktime, medianBlocktime, minBlocktime, maxBlocktime, stddevBlocktime, varianceBlocktime := getTimestampBlockSummary(bs)

	switch summaryOutputMode {
	case "text":
		// In the case where no transaction receipts could be retrieved, return.
		if successfulTx == 0 {
			log.Error().Msg("No transaction could be retrieved from the receipts")
			return
		}
		p.Printf("Successful Tx: %v\tTotal Tx: %v\n", number.Decimal(successfulTx), number.Decimal(totalTx))
		p.Printf("Total Mining Time: %s\n", totalMiningTime)
		p.Printf("Total Transactions: %v\n", number.Decimal(totalTransactions))
		p.Printf("Total Gas Used: %v\n", number.Decimal(totalGasUsed))
		p.Printf("Transactions per sec: %v\n", number.Decimal(tps))
		p.Printf("Gas Per Second: %v\n", number.Decimal(gaspersec))
		p.Printf("Latencies - Min: %v\tMedian: %v\tMax: %v\n", number.Decimal(minLatency.Seconds()), number.Decimal(medianLatency.Seconds()), number.Decimal(maxLatency.Seconds()))
		// Blocktime related metrics can only be calculated when there are at least two blocks
		if len(bs) > 1 {
			p.Printf("Mean Blocktime: %vs\n", number.Decimal(meanBlocktime))
			p.Printf("Median Blocktime: %vs\n", number.Decimal(medianBlocktime))
			p.Printf("Minimum Blocktime: %vs\n", number.Decimal(minBlocktime))
			p.Printf("Maximum Blocktime: %vs\n", number.Decimal(maxBlocktime))
			p.Printf("Blocktime Standard Deviation: %vs\n", number.Decimal(stddevBlocktime))
			p.Printf("Blocktime Variance: %vs\n", number.Decimal(varianceBlocktime))
		} else {
			log.Debug().Int("Length of blockSummary", len(bs)).Msg("blockSummary is empty")
		}
	case "json":
		summaryOutput := SummaryOutput{}
		summaryOutput.Summaries = jsonSummaryList
		summaryOutput.SuccessfulTx = successfulTx
		summaryOutput.TotalTx = totalTx
		summaryOutput.TotalMiningTime = totalMiningTime
		summaryOutput.TotalGasUsed = totalGasUsed
		summaryOutput.TransactionsPerSec = tps
		summaryOutput.GasPerSecond = gaspersec

		latencies := Latency{}
		latencies.Min = minLatency.Seconds()
		latencies.Median = medianLatency.Seconds()
		latencies.Max = maxLatency.Seconds()
		summaryOutput.Latencies = latencies

		val, _ := json.MarshalIndent(summaryOutput, "", "    ")
		p.Println(string(val))
	default:
		log.Error().Str("mode", summaryOutputMode).Msg("Invalid mode for summary output")
	}
}

func filterBlockSummary(ap *AccountPool, blockSummaries map[uint64]BlockSummary) {
	validTx := make(map[ethcommon.Hash]struct{})
	var minBlock uint64 = math.MaxUint64
	var maxBlock uint64 = 0
	for _, bs := range blockSummaries {
		for _, tx := range bs.Block.Transactions {
			startNonce, endNonce := ap.NoncesOf(tx.From.ToAddress())
			if tx.Nonce.ToUint64() >= startNonce && tx.Nonce.ToUint64() <= endNonce {
				validTx[tx.Hash.ToHash()] = struct{}{}
				if tx.BlockNumber.ToUint64() < minBlock {
					minBlock = tx.BlockNumber.ToUint64()
				}
				if tx.BlockNumber.ToUint64() > maxBlock {
					maxBlock = tx.BlockNumber.ToUint64()
				}
			}
		}
	}
	keys := getSortedMapKeys(blockSummaries)
	for _, k := range keys {
		if k < minBlock {
			delete(blockSummaries, k)
		}
		if k > maxBlock {
			delete(blockSummaries, k)
		}
	}

	for _, bs := range blockSummaries {
		filteredTransactions := make([]rpctypes.RawTransactionResponse, 0)
		for txKey, tx := range bs.Block.Transactions {
			if _, hasKey := validTx[tx.Hash.ToHash()]; hasKey {
				filteredTransactions = append(filteredTransactions, bs.Block.Transactions[txKey])
			}
		}
		bs.Block.Transactions = filteredTransactions
		filteredReceipts := make(map[ethcommon.Hash]rpctypes.RawTxReceipt)
		for receiptKey, receipt := range bs.Receipts {
			if _, hasKey := validTx[receipt.TransactionHash.ToHash()]; hasKey {
				filteredReceipts[receipt.TransactionHash.ToHash()] = bs.Receipts[receiptKey]
			}
		}
		bs.Receipts = filteredReceipts
	}
}

func getMapValues[K constraints.Ordered, V any](m map[K]V) []V {
	newSlice := make([]V, 0)
	for _, val := range m {
		newSlice = append(newSlice, val)
	}
	return newSlice
}

func getMinMedianMax[V constraints.Float | constraints.Integer](values []V) (V, V, V) {
	if len(values) == 0 {
		return 0, 0, 0
	}
	slices.Sort(values)
	half := len(values) / 2
	median := values[half]
	if len(values)%2 == 0 {
		median = (median + values[half-1]) / V(2)
	}
	var min V
	var max V
	for k, v := range values {
		if k == 0 {
			min = v
			max = v
			continue
		}
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return min, median, max
}

func getSortedMapKeys[V any, K constraints.Ordered](m map[K]V) []K {
	keys := make([]K, 0)
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

func getSuccessfulTransactionCount(bs map[uint64]BlockSummary) (successful, total int64) {
	for _, block := range bs {
		total += int64(len(block.Receipts))
		for _, receipt := range block.Receipts {
			successful += receipt.Status.ToInt64()
		}
	}
	return
}

func getTimestampBlockSummary(bs map[uint64]BlockSummary) (float64, float64, float64, float64, float64, float64) {
	blockTimestamps := make([]float64, 0)
	var prevBlockTimestamp float64 = 0
	mapKeys := getSortedMapKeys(bs)
	for _, v := range mapKeys {
		currBlockTimestamp := bs[v].Block.Timestamp.ToFloat64()
		if prevBlockTimestamp == 0 {
			prevBlockTimestamp = currBlockTimestamp
			continue
		} else if currBlockTimestamp > prevBlockTimestamp {
			blockTimeDiff := currBlockTimestamp - prevBlockTimestamp
			blockTimestamps = append(blockTimestamps, float64(blockTimeDiff))
			prevBlockTimestamp = currBlockTimestamp
		}
	}
	meanBlocktime, _ := stats.Mean(blockTimestamps)
	medianBlocktime, _ := stats.Median(blockTimestamps)
	minBlocktime, _ := stats.Min(blockTimestamps)
	maxBlocktime, _ := stats.Max(blockTimestamps)
	stddevBlocktime, _ := stats.StandardDeviation(blockTimestamps)
	varianceBlocktime, _ := stats.Variance(blockTimestamps)
	return meanBlocktime, medianBlocktime, minBlocktime, maxBlocktime, stddevBlocktime, varianceBlocktime
}

func getTotalGasUsed(receipts map[ethcommon.Hash]rpctypes.RawTxReceipt) uint64 {
	var totalGasUsed uint64 = 0
	for _, receipt := range receipts {
		totalGasUsed += receipt.GasUsed.ToUint64()
	}
	return totalGasUsed
}

func isEmptyJSONResponse(r *json.RawMessage) bool {
	rawJson := []byte(*r)
	return len(rawJson) == 0
}

// LightSummary prints a quick summary of load test results.
func LightSummary(results []Sample, startTime, endTime time.Time, rl *rate.Limiter) {
	if len(results) == 0 {
		log.Error().Msg("No results recorded")
		return
	}

	log.Info().Msg("* Results")
	log.Info().Int("samples", len(results)).Msg("Samples")

	var numErrors uint64 = 0

	latencies := make([]float64, 0)
	for _, s := range results {
		if s.IsError {
			numErrors++
		}
		latencies = append(latencies, s.WaitTime.Seconds())
	}

	testDuration := endTime.Sub(startTime)
	rps := float64(len(results)) / testDuration.Seconds()
	tps := float64(len(results)-int(numErrors)) / testDuration.Seconds()

	var rlLimit float64
	if rl != nil {
		rlLimit = float64(rl.Limit())
	}
	meanLat, _ := stats.Mean(latencies)
	medianLat, _ := stats.Median(latencies)
	minLat, _ := stats.Min(latencies)
	maxLat, _ := stats.Max(latencies)
	stddevLat, _ := stats.StandardDeviation(latencies)
	lastLTSample := lastSample(results)

	log.Info().Time("startTime", startTime).Msg("Start time of loadtest (first transaction sent)")
	log.Info().Time("loadStopTime", lastLTSample.RequestTime).Msg("End of load generation (last transaction sent)")
	log.Info().Time("endTime", endTime).Msg("End time of loadtest (final transaction mined)")
	log.Info().Float64("tps", tps).Msg("Successful Requests Per Second")
	if tps != rps {
		log.Error().Float64("rps", rps).Msg("Total Requests Per Second (both successful and unsuccessful transactions)")
	}
	log.Info().
		Float64("mean", meanLat).
		Float64("median", medianLat).
		Float64("min", minLat).
		Float64("max", maxLat).
		Float64("stddev", stddevLat).
		Msg("Request Latency of Transactions Stats")
	log.Info().
		Float64("testDuration", testDuration.Seconds()).
		Float64("finalRateLimit", rlLimit).
		Msg("Rough test summary")
	log.Info().Uint64("numErrors", numErrors).Msg("Num errors")
}

func lastSample(results []Sample) Sample {
	var maxTime time.Time
	var maxIdx int
	for idx, lt := range results {
		if maxTime.Before(lt.RequestTime) {
			maxTime = lt.RequestTime
			maxIdx = idx
		}
	}
	return results[maxIdx]
}
