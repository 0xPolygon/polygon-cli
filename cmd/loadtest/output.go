package loadtest

import (
	"context"
	"encoding/json"
	"math"
	"math/big"
	"sort"
	"sync"
	"time"

	"github.com/montanaflynn/stats"

	"golang.org/x/time/rate"

	"github.com/maticnetwork/polygon-cli/rpctypes"
	"github.com/maticnetwork/polygon-cli/util"
	"golang.org/x/exp/constraints"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"

	_ "embed"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func printBlockSummary(c *ethclient.Client, bs map[uint64]blockSummary, startNonce, endNonce uint64) {
	filterBlockSummary(bs, startNonce, endNonce)
	mapKeys := getSortedMapKeys(bs)
	if len(mapKeys) == 0 {
		return
	}

	var totalTransactions uint64 = 0
	var totalGasUsed uint64 = 0
	p := message.NewPrinter(language.English)

	allLatencies := make([]time.Duration, 0)
	summaryOutputMode := *inputLoadTestParams.SummaryOutputMode
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
			if summaryOutputMode == "text" {
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
			} else if summaryOutputMode == "json" {
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
			} else {
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

	if summaryOutputMode == "text" {
		p.Printf("Successful Tx: %v\tTotal Tx: %v\n", number.Decimal(successfulTx), number.Decimal(totalTx))
		p.Printf("Total Mining Time: %s\n", totalMiningTime)
		p.Printf("Total Transactions: %v\n", number.Decimal(totalTransactions))
		p.Printf("Total Gas Used: %v\n", number.Decimal(totalGasUsed))
		p.Printf("Transactions per sec: %v\n", number.Decimal(tps))
		p.Printf("Gas Per Second: %v\n", number.Decimal(gaspersec))
		p.Printf("Latencies - Min: %v\tMedian: %v\tMax: %v\n", number.Decimal(minLatency.Seconds()), number.Decimal(medianLatency.Seconds()), number.Decimal(maxLatency.Seconds()))
		p.Printf("Mean Blocktime: %vs\n", number.Decimal(meanBlocktime))
		p.Printf("Median Blocktime: %vs\n", number.Decimal(medianBlocktime))
		p.Printf("Minimum Blocktime: %vs\n", number.Decimal(minBlocktime))
		p.Printf("Maximum Blocktime: %vs\n", number.Decimal(maxBlocktime))
		p.Printf("Blocktime Standard Deviation: %vs\n", number.Decimal(stddevBlocktime))
		p.Printf("Blocktime Variance: %vs\n", number.Decimal(varianceBlocktime))
	} else if summaryOutputMode == "json" {
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
	} else {
		log.Error().Str("mode", summaryOutputMode).Msg("Invalid mode for summary output")
	}
}
func filterBlockSummary(blockSummaries map[uint64]blockSummary, startNonce, endNonce uint64) {
	validTx := make(map[ethcommon.Hash]struct{}, 0)
	var minBlock uint64 = math.MaxUint64
	var maxBlock uint64 = 0
	for _, bs := range blockSummaries {
		for _, tx := range bs.Block.Transactions {
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
		filteredReceipts := make(map[ethcommon.Hash]rpctypes.RawTxReceipt, 0)
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
	sort.Slice(values, func(i, j int) bool {
		return values[i] < values[j]
	})
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
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return keys
}

func getSuccessfulTransactionCount(bs map[uint64]blockSummary) (successful, total int64) {
	for _, block := range bs {
		total += int64(len(block.Receipts))
		for _, receipt := range block.Receipts {
			successful += receipt.Status.ToInt64()
		}
	}
	return
}

func getTimestampBlockSummary(bs map[uint64]blockSummary) (float64, float64, float64, float64, float64, float64) {
	blockTimestamps := make([]float64, 0)
	var prevBlockTimestamp float64 = 0
	// Keys for BlockSummary must be sorted, because it is not guaranteed that the keys of BlockSummary map is in order.
	mapKeys := getSortedMapKeys(bs)
	// Iterate through the BlockSummary elements and calculate the blocktime by comparing current and previous blocks' timestamps.
	for _, v := range mapKeys {
		currBlockTimestamp := bs[v].Block.Timestamp.ToFloat64()
		// Since the first block will not have a previous block to compare, continue.
		if prevBlockTimestamp == 0 {
			prevBlockTimestamp = currBlockTimestamp
			continue
			// Sanity check to make sure that the current blocktime is always greater than the previous blocktime.
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

type Latency struct {
	Min    float64
	Median float64
	Max    float64
}

type Summary struct {
	BlockNumber uint64
	Time        time.Time
	GasLimit    uint64
	GasUsed     uint64
	NumTx       int
	Utilization float64
	Latencies   Latency
}

type SummaryOutput struct {
	Summaries          []Summary
	SuccessfulTx       int64
	TotalTx            int64
	TotalMiningTime    time.Duration
	TotalGasUsed       uint64
	TransactionsPerSec float64
	GasPerSecond       float64
	Latencies          Latency
}

func summarizeTransactions(ctx context.Context, c *ethclient.Client, rpc *ethrpc.Client, startBlockNumber, startNonce, lastBlockNumber, endNonce uint64) error {
	ltp := inputLoadTestParams
	var err error

	log.Trace().Msg("Starting block range capture")
	// confirm start block number is ok
	_, err = c.BlockByNumber(ctx, new(big.Int).SetUint64(startBlockNumber))
	if err != nil {
		return err
	}
	rawBlocks, err := util.GetBlockRange(ctx, startBlockNumber, lastBlockNumber, rpc)
	if err != nil {
		return err
	}
	// TODO: Add some kind of decimation to avoid summarizing for 10 minutes?
	batchSize := *ltp.BatchSize
	goRoutineLimit := *ltp.Concurrency
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

	blockData := make(map[uint64]blockSummary, 0)
	for k, b := range blocks {
		bs := blockSummary{}
		bs.Block = &blocks[k]
		bs.Receipts = make(map[ethcommon.Hash]rpctypes.RawTxReceipt, 0)
		bs.Latencies = make(map[uint64]time.Duration, 0)
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

	nonceTimes := make(map[uint64]time.Time, 0)
	for _, ltr := range loadTestResults {
		nonceTimes[ltr.Nonce] = ltr.RequestTime
	}

	minLatency := time.Millisecond * 100
	for _, bs := range blockData {
		for _, tx := range bs.Block.Transactions {
			// TODO: What happens when the system clock of the load tester isn't in sync with the system clock of the miner?
			// TODO: the timestamp in the chain only has granularity down to the second. How to deal with this
			mineTime := time.Unix(bs.Block.Timestamp.ToInt64(), 0)
			requestTime := nonceTimes[tx.Nonce.ToUint64()]
			txLatency := mineTime.Sub(requestTime)
			if txLatency.Hours() > 2 {
				log.Debug().Float64("txHours", txLatency.Hours()).Uint64("nonce", tx.Nonce.ToUint64()).Uint64("blockNumber", bs.Block.Number.ToUint64()).Time("mineTime", mineTime).Time("requestTime", requestTime).Msg("Encountered transaction with more than 2 hours latency")
			}
			bs.Latencies[tx.Nonce.ToUint64()] = txLatency

			if txLatency < minLatency {
				minLatency = txLatency
			}
		}
	}
	// TODO this might be a hack, but not sure what's a better way to deal with time discrepancies
	if minLatency < time.Millisecond*100 {
		log.Trace().Str("minLatency", minLatency.String()).Msg("Minimum latency is below expected threshold")
		shiftSize := ((time.Millisecond * 100) - minLatency) + time.Millisecond + 100
		for _, bs := range blockData {
			for _, tx := range bs.Block.Transactions {
				bs.Latencies[tx.Nonce.ToUint64()] += shiftSize
			}
		}
	}

	printBlockSummary(c, blockData, startNonce, endNonce)

	log.Trace().Str("summaryTime", (endReceipt.Sub(startReceipt)).String()).Msg("Total Summary Time")

	return nil
}

func isEmptyJSONResponse(r *json.RawMessage) bool {
	rawJson := []byte(*r)
	return len(rawJson) == 0
}

func lightSummary(lts []loadTestSample, startTime, endTime time.Time, rl *rate.Limiter) {
	if len(lts) == 0 {
		log.Error().Msg("No results recorded")
		return
	}

	log.Info().Msg("* Results")
	log.Info().Int("samples", len(lts)).Msg("Samples")

	var numErrors uint64 = 0

	// latencies refers to the delay for the transactions to be relayed.
	latencies := make([]float64, 0)
	for _, s := range lts {
		if s.IsError {
			numErrors++
		}
		latencies = append(latencies, s.WaitTime.Seconds())
	}

	testDuration := endTime.Sub(startTime)
	tps := float64(len(loadTestResults)) / testDuration.Seconds()

	var rlLimit float64
	if rl != nil {
		rlLimit = float64(rl.Limit())
	}
	meanLat, _ := stats.Mean(latencies)
	medianLat, _ := stats.Median(latencies)
	minLat, _ := stats.Min(latencies)
	maxLat, _ := stats.Max(latencies)
	stddevLat, _ := stats.StandardDeviation(latencies)

	log.Info().Time("startTime", startTime).Msg("Start time of loadtest (first transaction sent)")
	log.Info().Time("endTime", endTime).Msg("End time of loadtest (final transaction mined)")
	log.Info().Float64("tps", tps).Msg("Overall Requests Per Second")
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
