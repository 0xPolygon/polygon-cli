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

func (pt *PreconfTracker) Stats() {
	log.Info().Uint64("total_tasks", pt.totalTasks.Load()).
		Uint64("preconf_success", pt.preconfSuccess.Load()).
		Uint64("preconf_fail", pt.preconfFail.Load()).
		Uint64("both_failed", pt.bothFailedCount.Load()).
		Uint64("ineffective_preconf", pt.ineffectivePreconf.Load()).
		Uint64("false_positives", pt.falsePositiveCount.Load()).
		Uint64("confidence", pt.confidence.Load()).
		Uint64("receipt_success", pt.receiptSuccess.Load()).
		Uint64("receipt_fail", pt.receiptFail.Load()).
		Uint64("total_gas_used", pt.totalGasUsed.Load()).
		Msg("Preconf Tracker Stats")

	if pt.statsFilePath == "" {
		return
	}

	// Copy txResults under lock
	pt.mu.Lock()
	txResults := make([]PreconfTxResult, len(pt.txResults))
	copy(txResults, pt.txResults)
	pt.mu.Unlock()

	// Build JSON output
	output := PreconfStats{
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
		},
		Transactions: txResults,
	}

	// Write JSON file
	timestamp := time.Now().Format(time.RFC3339)
	path := pt.statsFilePath + "-" + timestamp + ".json"

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		log.Error().Err(err).Msg("Error marshaling preconf stats")
		return
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		log.Error().Err(err).Msg("Error writing preconf stats file")
		return
	}

	log.Info().Str("path", path).Msg("Dumped preconf stats into file")
}
