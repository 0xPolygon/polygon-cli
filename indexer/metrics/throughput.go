package metrics

import (
	"sync"
	"time"

	"github.com/0xPolygon/polygon-cli/rpctypes"
)

// ThroughputMetric calculates throughput metrics over block-count windows
type ThroughputMetric struct {
	mu          sync.RWMutex
	blocks      []throughputBlockInfo
	maxBlocks   int // 30 to cover both 10 and 30 block windows
}

type throughputBlockInfo struct {
	timestamp uint64
	txCount   int
	gasUsed   uint64
}

// NewThroughputMetric creates a new throughput calculator
func NewThroughputMetric() *ThroughputMetric {
	return &ThroughputMetric{
		blocks:    make([]throughputBlockInfo, 0),
		maxBlocks: 30, // Track last 30 blocks to cover both windows
	}
}

// Name returns the metric identifier
func (t *ThroughputMetric) Name() string {
	return "throughput"
}

// ProcessBlock adds a new block to calculate throughput metrics
func (t *ThroughputMetric) ProcessBlock(block rpctypes.PolyBlock) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Add new block info to the front of the slice (newest first)
	info := throughputBlockInfo{
		timestamp: block.Time(),
		txCount:   len(block.Transactions()),
		gasUsed:   block.GasUsed(),
	}
	
	// Prepend new block (newest first)
	t.blocks = append([]throughputBlockInfo{info}, t.blocks...)

	// Maintain window size (keep only last 30 blocks)
	if len(t.blocks) > t.maxBlocks {
		t.blocks = t.blocks[:t.maxBlocks]
	}
}

// GetMetric returns the current throughput statistics
func (t *ThroughputMetric) GetMetric() interface{} {
	t.mu.RLock()
	defer t.mu.RUnlock()

	stats := ThroughputStats{
		TPS10: t.calculateTPS(10),
		TPS30: t.calculateTPS(30),
		GPS10: t.calculateGPS(10),
		GPS30: t.calculateGPS(30),
		BlocksAvailable: len(t.blocks),
	}

	return stats
}

// calculateTPS calculates transactions per second over the specified number of blocks
func (t *ThroughputMetric) calculateTPS(blockCount int) float64 {
	if len(t.blocks) < blockCount || len(t.blocks) < 2 {
		return 0.0
	}

	// Take the first blockCount blocks (newest to oldest)
	windowBlocks := t.blocks[:blockCount]

	// Calculate total transactions
	totalTxs := 0
	for _, block := range windowBlocks {
		totalTxs += block.txCount
	}

	// Calculate time span (newest - oldest)
	newestTime := windowBlocks[0].timestamp
	oldestTime := windowBlocks[blockCount-1].timestamp
	
	if newestTime <= oldestTime {
		return 0.0
	}

	timeSpan := newestTime - oldestTime
	if timeSpan == 0 {
		return 0.0
	}

	return float64(totalTxs) / float64(timeSpan)
}

// calculateGPS calculates gas per second over the specified number of blocks
func (t *ThroughputMetric) calculateGPS(blockCount int) float64 {
	if len(t.blocks) < blockCount || len(t.blocks) < 2 {
		return 0.0
	}

	// Take the first blockCount blocks (newest to oldest)
	windowBlocks := t.blocks[:blockCount]

	// Calculate total gas used
	totalGas := uint64(0)
	for _, block := range windowBlocks {
		totalGas += block.gasUsed
	}

	// Calculate time span (newest - oldest)
	newestTime := windowBlocks[0].timestamp
	oldestTime := windowBlocks[blockCount-1].timestamp
	
	if newestTime <= oldestTime {
		return 0.0
	}

	timeSpan := newestTime - oldestTime
	if timeSpan == 0 {
		return 0.0
	}

	return float64(totalGas) / float64(timeSpan)
}

// GetUpdateInterval returns how often this metric should be updated
func (t *ThroughputMetric) GetUpdateInterval() time.Duration {
	return 1 * time.Second
}

// ThroughputStats provides detailed throughput statistics
type ThroughputStats struct {
	TPS10           float64 // Transactions per second (10 blocks)
	TPS30           float64 // Transactions per second (30 blocks)
	GPS10           float64 // Gas per second (10 blocks)
	GPS30           float64 // Gas per second (30 blocks)
	BlocksAvailable int     // Number of blocks available for calculation
}