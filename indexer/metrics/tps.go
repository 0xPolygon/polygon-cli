package metrics

import (
	"sync"
	"time"

	"github.com/0xPolygon/polygon-cli/rpctypes"
)

// TPSMetric calculates transactions per second over a rolling window
type TPSMetric struct {
	mu         sync.RWMutex
	window     []blockInfo
	windowSize time.Duration
	tps        float64
}

type blockInfo struct {
	timestamp uint64
	txCount   int
}

// NewTPSMetric creates a new TPS calculator with a 30-second window
func NewTPSMetric() *TPSMetric {
	return &TPSMetric{
		window:     make([]blockInfo, 0),
		windowSize: 30 * time.Second,
		tps:        0,
	}
}

// Name returns the metric identifier
func (t *TPSMetric) Name() string {
	return "tps"
}

// ProcessBlock adds a new block to the calculation window
func (t *TPSMetric) ProcessBlock(block rpctypes.PolyBlock) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Add new block info
	info := blockInfo{
		timestamp: block.Time(),
		txCount:   len(block.Transactions()),
	}
	t.window = append(t.window, info)

	// Prune old blocks outside the window
	t.pruneOldBlocks()

	// Recalculate TPS
	t.calculateTPS()
}

// pruneOldBlocks removes blocks outside the time window
func (t *TPSMetric) pruneOldBlocks() {
	if len(t.window) == 0 {
		return
	}

	cutoff := uint64(time.Now().Unix()) - uint64(t.windowSize.Seconds())

	// Find the first block within the window
	startIdx := 0
	for i, block := range t.window {
		if block.timestamp >= cutoff {
			startIdx = i
			break
		}
	}

	// Keep only blocks within the window
	if startIdx > 0 {
		t.window = t.window[startIdx:]
	}
}

// calculateTPS calculates the current TPS based on blocks in the window
func (t *TPSMetric) calculateTPS() {
	if len(t.window) < 2 {
		t.tps = 0
		return
	}

	// Calculate total transactions and time span
	totalTxs := 0
	for _, block := range t.window {
		totalTxs += block.txCount
	}

	// Get time span from oldest to newest block
	oldestTime := t.window[0].timestamp
	newestTime := t.window[len(t.window)-1].timestamp
	timeSpan := newestTime - oldestTime

	if timeSpan > 0 {
		t.tps = float64(totalTxs) / float64(timeSpan)
	} else {
		t.tps = 0
	}
}

// GetMetric returns the current TPS value
func (t *TPSMetric) GetMetric() interface{} {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.tps
}

// GetUpdateInterval returns how often this metric should be updated
func (t *TPSMetric) GetUpdateInterval() time.Duration {
	return 1 * time.Second
}

// TPSStats provides detailed TPS statistics
type TPSStats struct {
	CurrentTPS   float64
	WindowSize   time.Duration
	BlockCount   int
	TotalTxCount int
}
