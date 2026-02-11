package loadtest

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/0xPolygon/polygon-cli/loadtest/config"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
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

// trackedTx holds tracking state for a pending transaction.
// Resolution state is determined by non-nil fields:
// - receipt != nil OR receiptError != nil means receipt is resolved
// - preconfResult != nil OR preconfError != nil means preconf is resolved
type trackedTx struct {
	hash         common.Hash
	registeredAt time.Time
	startBlock   uint64
	// Receipt resolution - exactly one of (receipt, receiptError) will be set when resolved
	receipt      *types.Receipt
	receiptTime  time.Duration
	receiptError error
	// Preconf resolution - exactly one of (preconfResult, preconfError) will be set when resolved
	preconfResult *bool
	preconfTime   time.Duration
	preconfError  error
}

// receiptResolved returns true if the receipt has been resolved (success or error).
func (tx *trackedTx) receiptResolved() bool {
	return tx.receipt != nil || tx.receiptError != nil
}

// preconfResolved returns true if the preconf status has been resolved (success or error).
func (tx *trackedTx) preconfResolved() bool {
	return tx.preconfResult != nil || tx.preconfError != nil
}

// PreconfTracker tracks preconf and receipt status using centralized batch polling.
type PreconfTracker struct {
	client *ethclient.Client
	rpc    *ethrpc.Client
	cfg    *config.PreconfConfig

	// Pending transactions awaiting receipt/preconf
	pendingMu sync.RWMutex
	pending   map[common.Hash]*trackedTx

	// Completed transaction results
	completedMu sync.Mutex
	completed   []*trackedTx

	// Metrics (atomic for lock-free access)
	preconfSuccess     atomic.Uint64
	preconfFail        atomic.Uint64
	totalTasks         atomic.Uint64
	bothFailedCount    atomic.Uint64
	ineffectivePreconf atomic.Uint64
	falsePositiveCount atomic.Uint64
	confidence         atomic.Uint64
	receiptSuccess     atomic.Uint64
	receiptFail        atomic.Uint64
	totalGasUsed       atomic.Uint64

	// Shutdown coordination
	wg sync.WaitGroup
}

// NewPreconfTracker creates a new PreconfTracker.
func NewPreconfTracker(client *ethclient.Client, rpcClient *ethrpc.Client, cfg *config.PreconfConfig) *PreconfTracker {
	return &PreconfTracker{
		client:    client,
		rpc:       rpcClient,
		cfg:       cfg,
		pending:   make(map[common.Hash]*trackedTx),
		completed: make([]*trackedTx, 0, 1024),
	}
}

// Start begins the batch polling loops for receipts and preconf status.
func (pt *PreconfTracker) Start(ctx context.Context) {
	// Start receipt poller
	pt.wg.Add(1)
	go pt.receiptPoller(ctx)

	// Start preconf poller
	pt.wg.Add(1)
	go pt.preconfPoller(ctx)

	// Start stats file writer if path is configured
	if pt.cfg.StatsFile != "" {
		pt.wg.Add(1)
		go pt.statsWriter(ctx)
	}
}

// RegisterTx adds a transaction hash to be tracked. Non-blocking.
func (pt *PreconfTracker) RegisterTx(hash common.Hash) {
	currentBlock, err := pt.client.BlockNumber(context.Background())
	if err != nil {
		log.Warn().Err(err).Str("hash", hash.Hex()).Msg("Failed to get current block for preconf tracking")
		currentBlock = 0
	}

	pt.pendingMu.Lock()
	pt.pending[hash] = &trackedTx{
		hash:         hash,
		registeredAt: time.Now(),
		startBlock:   currentBlock,
	}
	pt.pendingMu.Unlock()
}

// receiptPoller polls for receipts in batches.
func (pt *PreconfTracker) receiptPoller(ctx context.Context) {
	defer pt.wg.Done()

	ticker := time.NewTicker(pt.cfg.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			pt.pollReceipts(ctx)
		}
	}
}

// preconfPoller polls for preconf status in batches.
func (pt *PreconfTracker) preconfPoller(ctx context.Context) {
	defer pt.wg.Done()

	ticker := time.NewTicker(pt.cfg.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			pt.pollPreconf(ctx)
		}
	}
}

// statsWriter periodically writes stats to file.
func (pt *PreconfTracker) statsWriter(ctx context.Context) {
	defer pt.wg.Done()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			pt.writeStats()
		}
	}
}

// pollReceipts fetches receipts for pending transactions in batches.
func (pt *PreconfTracker) pollReceipts(ctx context.Context) {
	pt.pendingMu.RLock()
	hashes := make([]common.Hash, 0, len(pt.pending))
	for hash, tx := range pt.pending {
		if !tx.receiptResolved() {
			hashes = append(hashes, hash)
		}
	}
	pt.pendingMu.RUnlock()

	if len(hashes) == 0 {
		return
	}

	// Process in batches
	for i := 0; i < len(hashes); i += pt.cfg.BatchSize {
		end := min(i+pt.cfg.BatchSize, len(hashes))
		batch := hashes[i:end]
		pt.getReceipts(ctx, batch)
	}

	// Check for completed/timed out transactions
	pt.checkTx()
}

// pollPreconf fetches preconf status for pending transactions in batches.
func (pt *PreconfTracker) pollPreconf(ctx context.Context) {
	pt.pendingMu.RLock()
	hashes := make([]common.Hash, 0, len(pt.pending))
	for hash, tx := range pt.pending {
		if !tx.preconfResolved() {
			hashes = append(hashes, hash)
		}
	}
	pt.pendingMu.RUnlock()

	if len(hashes) == 0 {
		return
	}

	// Process in batches
	for i := 0; i < len(hashes); i += pt.cfg.BatchSize {
		end := min(i+pt.cfg.BatchSize, len(hashes))
		batch := hashes[i:end]
		pt.getPreconfs(ctx, batch)
	}
}

// getReceipts fetches multiple receipts in a single batch RPC call.
func (pt *PreconfTracker) getReceipts(ctx context.Context, hashes []common.Hash) {
	if len(hashes) == 0 {
		return
	}

	// Create batch elements
	batch := make([]ethrpc.BatchElem, len(hashes))
	receipts := make([]*types.Receipt, len(hashes))

	for i, hash := range hashes {
		receipts[i] = new(types.Receipt)
		batch[i] = ethrpc.BatchElem{
			Method: "eth_getTransactionReceipt",
			Args:   []any{hash},
			Result: receipts[i],
		}
	}

	// Execute batch call
	if err := pt.rpc.BatchCallContext(ctx, batch); err != nil {
		log.Debug().Err(err).Int("count", len(hashes)).Msg("Batch receipt call failed")
		return
	}

	// Process results
	pt.pendingMu.Lock()
	now := time.Now()
	for i, hash := range hashes {
		tx, exists := pt.pending[hash]
		if !exists {
			continue
		}

		if batch[i].Error != nil {
			tx.receiptError = batch[i].Error
			continue
		}

		// Check if receipt is nil (not yet mined)
		if receipts[i] == nil || receipts[i].BlockNumber == nil {
			continue
		}

		tx.receipt = receipts[i]
		tx.receiptTime = now.Sub(tx.registeredAt)
	}
	pt.pendingMu.Unlock()
}

// getPreconfs checks preconf status for multiple transactions in a single batch RPC call.
func (pt *PreconfTracker) getPreconfs(ctx context.Context, hashes []common.Hash) {
	if len(hashes) == 0 {
		return
	}

	// Create batch elements
	batch := make([]ethrpc.BatchElem, len(hashes))
	results := make([]bool, len(hashes))

	for i, hash := range hashes {
		batch[i] = ethrpc.BatchElem{
			Method: "eth_checkPreconfStatus",
			Args:   []any{hash.Hex()},
			Result: &results[i],
		}
	}

	// Execute batch call
	if err := pt.rpc.BatchCallContext(ctx, batch); err != nil {
		log.Debug().Err(err).Int("count", len(hashes)).Msg("Batch preconf call failed")
		return
	}

	// Process results
	pt.pendingMu.Lock()
	now := time.Now()
	for i, hash := range hashes {
		tx, exists := pt.pending[hash]
		if !exists {
			continue
		}

		if batch[i].Error != nil {
			tx.preconfError = batch[i].Error
			continue
		}

		// Preconf result received (true = preconf confirmed)
		// If false, we keep polling until timeout
		if results[i] {
			tx.preconfResult = &results[i]
			tx.preconfTime = now.Sub(tx.registeredAt)
		}
	}
	pt.pendingMu.Unlock()
}

// checkTx moves completed or timed out transactions to the completed list.
func (pt *PreconfTracker) checkTx() {
	now := time.Now()

	pt.pendingMu.Lock()
	var toRemove []common.Hash

	for hash, tx := range pt.pending {
		timedOut := now.Sub(tx.registeredAt) >= pt.cfg.Timeout
		receiptDone := tx.receiptResolved() || timedOut
		preconfDone := tx.preconfResolved() || timedOut

		if receiptDone && preconfDone {
			toRemove = append(toRemove, hash)
			pt.recordMetrics(tx)
		}
	}

	for _, hash := range toRemove {
		delete(pt.pending, hash)
	}
	pt.pendingMu.Unlock()
}

// recordMetrics records the final metrics for a completed transaction.
func (pt *PreconfTracker) recordMetrics(tx *trackedTx) {
	pt.totalTasks.Add(1)

	preconfSuccess := tx.preconfResult != nil && *tx.preconfResult
	receiptSuccess := tx.receipt != nil

	// Track preconf metrics
	if preconfSuccess {
		pt.preconfSuccess.Add(1)
	} else {
		pt.preconfFail.Add(1)
	}

	// Track receipt metrics
	if receiptSuccess {
		pt.receiptSuccess.Add(1)
		pt.totalGasUsed.Add(tx.receipt.GasUsed)
	} else {
		pt.receiptFail.Add(1)
	}

	// Track combined metrics
	switch {
	case !preconfSuccess && !receiptSuccess:
		pt.bothFailedCount.Add(1)

	case preconfSuccess && !receiptSuccess:
		pt.falsePositiveCount.Add(1)

	case !preconfSuccess && receiptSuccess:
		pt.ineffectivePreconf.Add(1)

	case preconfSuccess && receiptSuccess:
		// Both succeeded - check if preconf was faster
		if tx.preconfTime > tx.receiptTime {
			pt.ineffectivePreconf.Add(1)
		}
		// Track confidence (block diff < 10)
		if tx.startBlock > 0 {
			blockDiff := tx.receipt.BlockNumber.Uint64() - tx.startBlock
			if blockDiff < 10 {
				pt.confidence.Add(1)
			}
		}
	}

	// Add to completed list for stats
	pt.completedMu.Lock()
	pt.completed = append(pt.completed, tx)
	pt.completedMu.Unlock()
}

// Percentiles holds p50, p75, p90, p95, p99 values.
type Percentiles struct {
	P50 float64
	P75 float64
	P90 float64
	P95 float64
	P99 float64
}

// percentiles computes p50, p75, p90, p95, p99 for a slice of durations.
// Returns zero values if the input slice is empty.
func percentiles(durations []float64) Percentiles {
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

// Stats logs the final summary and writes the stats file.
func (pt *PreconfTracker) Stats() {
	// Finalize any remaining pending transactions
	pt.finalizePending()

	output := pt.buildStats()
	log.Info().Any("summary", output.Summary).Msg("Preconf tracker stats")

	if pt.cfg.StatsFile == "" {
		return
	}
	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal preconf stats")
		return
	}
	if err := os.WriteFile(pt.cfg.StatsFile, data, 0644); err != nil {
		log.Error().Err(err).Msg("Failed to write preconf stats file")
	}
}

// finalizePending moves all remaining pending transactions to completed (as timed out).
func (pt *PreconfTracker) finalizePending() {
	pt.pendingMu.Lock()
	for _, tx := range pt.pending {
		pt.recordMetrics(tx)
	}
	pt.pending = make(map[common.Hash]*trackedTx)
	pt.pendingMu.Unlock()
}

func (pt *PreconfTracker) writeStats() {
	data, err := json.MarshalIndent(pt.buildStats(), "", "  ")
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal preconf stats")
		return
	}
	if err := os.WriteFile(pt.cfg.StatsFile, data, 0644); err != nil {
		log.Error().Err(err).Msg("Failed to write preconf stats file")
	}
}

func (pt *PreconfTracker) buildStats() PreconfStats {
	pt.completedMu.Lock()
	completed := make([]*trackedTx, len(pt.completed))
	copy(completed, pt.completed)
	pt.completedMu.Unlock()

	var pd, rd []float64
	txs := make([]PreconfTxResult, 0, len(completed))

	for _, tx := range completed {
		result := PreconfTxResult{
			TxHash: tx.hash.Hex(),
		}

		if tx.preconfResult != nil && *tx.preconfResult {
			result.PreconfDurationMs = tx.preconfTime.Milliseconds()
			pd = append(pd, float64(result.PreconfDurationMs))
		}

		if tx.receipt != nil {
			result.ReceiptDurationMs = tx.receiptTime.Milliseconds()
			result.GasUsed = tx.receipt.GasUsed
			result.Status = tx.receipt.Status
			if tx.startBlock > 0 {
				result.BlockDiff = tx.receipt.BlockNumber.Uint64() - tx.startBlock
			}
			rd = append(rd, float64(result.ReceiptDurationMs))
		}

		txs = append(txs, result)
	}

	pp := percentiles(pd)
	rp := percentiles(rd)

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

			PreconfP50: pp.P50,
			PreconfP75: pp.P75,
			PreconfP90: pp.P90,
			PreconfP95: pp.P95,
			PreconfP99: pp.P99,

			ReceiptP50: rp.P50,
			ReceiptP75: rp.P75,
			ReceiptP90: rp.P90,
			ReceiptP95: rp.P95,
			ReceiptP99: rp.P99,
		},
		Transactions: txs,
	}
}
