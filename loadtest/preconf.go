package loadtest

import (
	"context"
	"encoding/json"
	"os"
	"strings"
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
	TotalTasks uint64 `json:"total_tasks"`

	// Preconf totals
	PreconfSuccess uint64 `json:"preconf_success"` // preconf=true
	PreconfFail    uint64 `json:"preconf_fail"`    // preconf=false

	// Receipt totals
	ReceiptSuccess uint64 `json:"receipt_success"` // receipt=yes
	ReceiptFail    uint64 `json:"receipt_fail"`    // receipt=no

	// 2x2 outcome matrix
	BothConfirmed    uint64 `json:"both_confirmed"`    // preconf=true AND receipt=yes
	PreconfOnly      uint64 `json:"preconf_only"`      // preconf=true AND receipt=no
	ReceiptOnly      uint64 `json:"receipt_only"`      // preconf=false AND receipt=yes
	NeitherConfirmed uint64 `json:"neither_confirmed"` // preconf=false AND receipt=no

	Confidence   uint64 `json:"confidence"` // both_confirmed AND block_diff < 10
	TotalGasUsed uint64 `json:"total_gas_used"`

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
	totalTasks     atomic.Uint64
	preconfSuccess atomic.Uint64
	preconfFail    atomic.Uint64
	receiptSuccess atomic.Uint64
	receiptFail    atomic.Uint64

	// 2x2 outcome matrix
	bothConfirmed    atomic.Uint64
	preconfOnly      atomic.Uint64
	receiptOnly      atomic.Uint64
	neitherConfirmed atomic.Uint64

	confidence   atomic.Uint64
	totalGasUsed atomic.Uint64

	// Cached block number to reduce RPC calls
	cachedBlock     atomic.Uint64
	cachedBlockTime atomic.Int64

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

// getBlockNumber returns the current block number, cached to reduce RPC calls.
// The cache is refreshed at most once per second.
func (pt *PreconfTracker) getBlockNumber() uint64 {
	now := time.Now().UnixNano()
	lastUpdate := pt.cachedBlockTime.Load()
	if now-lastUpdate < int64(time.Second) {
		return pt.cachedBlock.Load()
	}
	if pt.cachedBlockTime.CompareAndSwap(lastUpdate, now) {
		if block, err := pt.client.BlockNumber(context.Background()); err == nil {
			pt.cachedBlock.Store(block)
		}
	}
	return pt.cachedBlock.Load()
}

// RegisterTx adds a transaction hash to be tracked. Non-blocking.
func (pt *PreconfTracker) RegisterTx(hash common.Hash) {
	pt.pendingMu.Lock()
	pt.pending[hash] = &trackedTx{
		hash:         hash,
		registeredAt: time.Now(),
		startBlock:   pt.getBlockNumber(),
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
		log.Warn().Err(err).Int("count", len(hashes)).Msg("Batch receipt call failed")
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
			// "missing required field" errors typically mean the tx isn't mined yet
			// (some RPCs return {} instead of null for pending txs)
			// Don't treat as permanent error - continue polling
			if strings.Contains(batch[i].Error.Error(), "missing required field") {
				continue
			}
			tx.receiptError = batch[i].Error
			log.Warn().Err(batch[i].Error).Str("hash", hash.Hex()).Msg("Receipt batch element error")
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
		log.Warn().Err(err).Int("count", len(hashes)).Msg("Batch preconf call failed")
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
			log.Warn().Err(batch[i].Error).Str("hash", hash.Hex()).Msg("Preconf batch element error")
			continue
		}

		// Record preconf result (true = confirmed, false = not confirmed)
		tx.preconfResult = &results[i]
		tx.preconfTime = now.Sub(tx.registeredAt)
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

	preconfTrue := tx.preconfResult != nil && *tx.preconfResult
	receiptOK := tx.receipt != nil

	// Track preconf totals
	if preconfTrue {
		pt.preconfSuccess.Add(1)
	} else {
		pt.preconfFail.Add(1)
	}

	// Track receipt totals
	if receiptOK {
		pt.receiptSuccess.Add(1)
		pt.totalGasUsed.Add(tx.receipt.GasUsed)
	} else {
		pt.receiptFail.Add(1)
	}

	// Track 2x2 matrix
	switch {
	case preconfTrue && receiptOK:
		pt.bothConfirmed.Add(1)
		if tx.startBlock > 0 {
			blockDiff := tx.receipt.BlockNumber.Uint64() - tx.startBlock
			if blockDiff < 10 {
				pt.confidence.Add(1)
			}
		}
	case preconfTrue && !receiptOK:
		pt.preconfOnly.Add(1)
	case !preconfTrue && receiptOK:
		pt.receiptOnly.Add(1)
	case !preconfTrue && !receiptOK:
		pt.neitherConfirmed.Add(1)
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
			TotalTasks:     pt.totalTasks.Load(),
			PreconfSuccess: pt.preconfSuccess.Load(),
			PreconfFail:    pt.preconfFail.Load(),
			ReceiptSuccess: pt.receiptSuccess.Load(),
			ReceiptFail:    pt.receiptFail.Load(),

			BothConfirmed:    pt.bothConfirmed.Load(),
			PreconfOnly:      pt.preconfOnly.Load(),
			ReceiptOnly:      pt.receiptOnly.Load(),
			NeitherConfirmed: pt.neitherConfirmed.Load(),

			Confidence:   pt.confidence.Load(),
			TotalGasUsed: pt.totalGasUsed.Load(),

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
