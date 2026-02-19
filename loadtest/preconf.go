package loadtest

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/0xPolygon/polygon-cli/util"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/montanaflynn/stats"
	"github.com/rs/zerolog/log"
)

// PreconfTxResult holds per-transaction preconf and receipt data.
type PreconfTxResult struct {
	TxHash            string `json:"tx_hash"`
	PreconfDurationMs int64  `json:"preconf_duration_ms,omitempty"`
	ReceiptDurationMs int64  `json:"receipt_duration_ms,omitempty"`
	BlockDiff         uint64 `json:"block_diff,omitempty"`
	GasUsed           uint64 `json:"gas_used,omitempty"`
	Status            uint64 `json:"status,omitempty"` // 1 = success, 0 = fail
}

// PreconfSummary holds aggregate stats from the preconf tracker.
type PreconfSummary struct {
	TotalTasks         uint64 `json:"total_tasks"`
	PreconfSuccess     uint64 `json:"preconf_success"`
	PreconfFail        uint64 `json:"preconf_fail"`
	BothFailed         uint64 `json:"both_failed"`
	IneffectivePreconf uint64 `json:"ineffective_preconf"`
	FalsePositives     uint64 `json:"false_positives"`
	Confidence         uint64 `json:"confidence"`
	ReceiptSuccess     uint64 `json:"receipt_success"`
	ReceiptFail        uint64 `json:"receipt_fail"`
	TotalGasUsed       uint64 `json:"total_gas_used"`

	// Preconf duration percentiles (milliseconds)
	PreconfP50 float64 `json:"preconf_p50,omitempty"`
	PreconfP75 float64 `json:"preconf_p75,omitempty"`
	PreconfP90 float64 `json:"preconf_p90,omitempty"`
	PreconfP95 float64 `json:"preconf_p95,omitempty"`
	PreconfP99 float64 `json:"preconf_p99,omitempty"`

	// Receipt duration percentiles (milliseconds)
	ReceiptP50 float64 `json:"receipt_p50,omitempty"`
	ReceiptP75 float64 `json:"receipt_p75,omitempty"`
	ReceiptP90 float64 `json:"receipt_p90,omitempty"`
	ReceiptP95 float64 `json:"receipt_p95,omitempty"`
	ReceiptP99 float64 `json:"receipt_p99,omitempty"`
}

// PreconfStats is the JSON output structure containing summary and per-tx data.
type PreconfStats struct {
	Summary      PreconfSummary    `json:"summary"`
	Transactions []PreconfTxResult `json:"transactions"`
}

type PreconfTracker struct {
	client        *ethclient.Client
	statsFilePath string

	// preconf metrics
	preconfSuccess     atomic.Uint64
	preconfFail        atomic.Uint64
	totalTasks         atomic.Uint64
	bothFailedCount    atomic.Uint64
	ineffectivePreconf atomic.Uint64
	falsePositiveCount atomic.Uint64
	confidence         atomic.Uint64

	// receipt metrics
	receiptSuccess atomic.Uint64
	receiptFail    atomic.Uint64
	totalGasUsed   atomic.Uint64

	mu        sync.Mutex
	txResults []PreconfTxResult
}

func NewPreconfTracker(client *ethclient.Client, statsFilePath string) *PreconfTracker {
	return &PreconfTracker{
		client:        client,
		statsFilePath: statsFilePath,
		txResults:     make([]PreconfTxResult, 0, 1024),
	}
}

func (pt *PreconfTracker) Track(txHash common.Hash) {
	currentBlock, err := pt.client.BlockNumber(context.Background())
	if err != nil {
		return
	}

	// wait for preconf
	var wg sync.WaitGroup
	var preconfStatus bool
	var preconfError error
	var preconfDuration time.Duration
	wg.Add(1)
	go func() {
		defer wg.Done()

		preconfStartTime := time.Now()
		defer func() {
			preconfDuration = time.Since(preconfStartTime)
		}()

		preconfStatus, preconfError = util.WaitPreconf(context.Background(), pt.client, txHash, time.Minute)
	}()

	// wait for receipt
	var receipt *types.Receipt
	var receiptError error
	var receiptDuration time.Duration
	wg.Add(1)
	go func() {
		defer wg.Done()

		time.Sleep(100 * time.Millisecond)

		receiptTime := time.Now()
		defer func() {
			receiptDuration = time.Since(receiptTime)
		}()

		receipt, receiptError = util.WaitReceiptWithTimeout(context.Background(), pt.client, txHash, time.Minute)
	}()

	wg.Wait()

	// Build per-transaction result
	result := PreconfTxResult{
		TxHash: txHash.Hex(),
	}

	pt.totalTasks.Add(1)
	if preconfStatus {
		pt.preconfSuccess.Add(1)
		result.PreconfDurationMs = preconfDuration.Milliseconds()
	} else {
		pt.preconfFail.Add(1)
	}

	// Track receipt metrics
	if receiptError == nil {
		pt.receiptSuccess.Add(1)
		pt.totalGasUsed.Add(receipt.GasUsed)
		result.ReceiptDurationMs = receiptDuration.Milliseconds()
		result.GasUsed = receipt.GasUsed
		result.Status = receipt.Status
		result.BlockDiff = receipt.BlockNumber.Uint64() - currentBlock
	} else {
		pt.receiptFail.Add(1)
	}

	// Append result under lock
	pt.mu.Lock()
	pt.txResults = append(pt.txResults, result)
	pt.mu.Unlock()

	switch {
	case preconfError != nil && receiptError != nil:
		// Both failed: no tx inclusion in txpool or block
		pt.bothFailedCount.Add(1)

	case preconfError == nil && receiptError != nil:
		// False positive: preconf said tx is included but never got executed
		pt.falsePositiveCount.Add(1)

	case preconfError != nil && receiptError == nil:
		// Receipt arrived but preconf failed: preconf wasn't effective
		pt.ineffectivePreconf.Add(1)

	case preconfError == nil && receiptError == nil:
		// Both succeeded
		if preconfDuration > receiptDuration {
			// Receipt arrived before preconf: preconf wasn't effective
			pt.ineffectivePreconf.Add(1)
		}
		// Track confidence (block diff < 10)
		if result.BlockDiff < 10 {
			pt.confidence.Add(1)
		}
	}
}

// Percentiles holds p50, p75, p90, p95, p99 values.
type Percentiles struct {
	P50 float64
	P75 float64
	P90 float64
	P95 float64
	P99 float64
}

// calculatePercentiles computes p50, p75, p90, p95, p99 for a slice of durations.
// Returns zero values if the input slice is empty.
func calculatePercentiles(durations []float64) Percentiles {
	if len(durations) == 0 {
		return Percentiles{}
	}
	p50, _ := stats.Percentile(durations, 50)
	p75, _ := stats.Percentile(durations, 75)
	p90, _ := stats.Percentile(durations, 90)
	p95, _ := stats.Percentile(durations, 95)
	p99, _ := stats.Percentile(durations, 99)
	return Percentiles{P50: p50, P75: p75, P90: p90, P95: p95, P99: p99}
}

// Start begins periodic stats file writing every 2 seconds until context is cancelled.
func (pt *PreconfTracker) Start(ctx context.Context) {
	if pt.statsFilePath == "" {
		return
	}
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				pt.writeStatsFile()
			case <-ctx.Done():
				return
			}
		}
	}()
}

// Stats logs the final summary and writes the stats file.
func (pt *PreconfTracker) Stats() {
	output := pt.buildStats()
	log.Info().Any("summary", output.Summary).Msg("Preconf tracker stats")

	if pt.statsFilePath == "" {
		return
	}
	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal preconf stats")
		return
	}
	if err := os.WriteFile(pt.statsFilePath, data, 0644); err != nil {
		log.Error().Err(err).Msg("Failed to write preconf stats file")
	}
}

func (pt *PreconfTracker) writeStatsFile() {
	data, err := json.MarshalIndent(pt.buildStats(), "", "  ")
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal preconf stats")
		return
	}
	if err := os.WriteFile(pt.statsFilePath, data, 0644); err != nil {
		log.Error().Err(err).Msg("Failed to write preconf stats file")
	}
}

func (pt *PreconfTracker) buildStats() PreconfStats {
	pt.mu.Lock()
	txResults := make([]PreconfTxResult, len(pt.txResults))
	copy(txResults, pt.txResults)
	pt.mu.Unlock()

	var preconfDurations, receiptDurations []float64
	for _, tx := range txResults {
		if tx.PreconfDurationMs > 0 {
			preconfDurations = append(preconfDurations, float64(tx.PreconfDurationMs))
		}
		if tx.ReceiptDurationMs > 0 {
			receiptDurations = append(receiptDurations, float64(tx.ReceiptDurationMs))
		}
	}

	preconfPct := calculatePercentiles(preconfDurations)
	receiptPct := calculatePercentiles(receiptDurations)

	return PreconfStats{
		Summary: PreconfSummary{
			TotalTasks:         pt.totalTasks.Load(),
			PreconfSuccess:     pt.preconfSuccess.Load(),
			PreconfFail:        pt.preconfFail.Load(),
			BothFailed:         pt.bothFailedCount.Load(),
			IneffectivePreconf: pt.ineffectivePreconf.Load(),
			FalsePositives:     pt.falsePositiveCount.Load(),
			Confidence:         pt.confidence.Load(),
			ReceiptSuccess:     pt.receiptSuccess.Load(),
			ReceiptFail:        pt.receiptFail.Load(),
			TotalGasUsed:       pt.totalGasUsed.Load(),

			PreconfP50: preconfPct.P50,
			PreconfP75: preconfPct.P75,
			PreconfP90: preconfPct.P90,
			PreconfP95: preconfPct.P95,
			PreconfP99: preconfPct.P99,

			ReceiptP50: receiptPct.P50,
			ReceiptP75: receiptPct.P75,
			ReceiptP90: receiptPct.P90,
			ReceiptP95: receiptPct.P95,
			ReceiptP99: receiptPct.P99,
		},
		Transactions: txResults,
	}
}
