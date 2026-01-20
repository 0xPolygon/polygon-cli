package metrics

import (
	"sync"
	"time"

	"github.com/0xPolygon/polygon-cli/rpctypes"
)

// BlockTimeMetric calculates average block time over a rolling window
type BlockTimeMetric struct {
	mu          sync.RWMutex
	blocks      []blockTimeInfo
	windowSize  int           // number of blocks to track
	averageTime time.Duration // current average block time
}

type blockTimeInfo struct {
	timestamp uint64
	blockTime time.Duration // time between this block and previous (chronologically)
}

// NewBlockTimeMetric creates a new block time calculator
func NewBlockTimeMetric() *BlockTimeMetric {
	return &BlockTimeMetric{
		blocks:     make([]blockTimeInfo, 0),
		windowSize: 50, // Track last 50 blocks for average
	}
}

// Name returns the metric identifier
func (b *BlockTimeMetric) Name() string {
	return "blockTime"
}

// ProcessBlock adds a new block to calculate block time
func (b *BlockTimeMetric) ProcessBlock(block rpctypes.PolyBlock) {
	b.mu.Lock()
	defer b.mu.Unlock()

	currentTimestamp := block.Time()

	// Add block to the list (they come in newest-first order)
	info := blockTimeInfo{
		timestamp: currentTimestamp,
		blockTime: 0, // Will be calculated after sorting
	}
	b.blocks = append(b.blocks, info)

	// Maintain window size
	if len(b.blocks) > b.windowSize {
		b.blocks = b.blocks[1:]
	}

	// Sort blocks by timestamp (oldest first) to calculate block times correctly
	b.sortAndCalculateBlockTimes()

	// Recalculate average
	b.calculateAverage()
}

// sortAndCalculateBlockTimes sorts blocks chronologically and calculates block times
func (b *BlockTimeMetric) sortAndCalculateBlockTimes() {
	if len(b.blocks) < 2 {
		return
	}

	// Sort by timestamp (oldest first)
	// Using a simple bubble sort since we typically only add one block at a time
	for i := 0; i < len(b.blocks)-1; i++ {
		for j := 0; j < len(b.blocks)-i-1; j++ {
			if b.blocks[j].timestamp > b.blocks[j+1].timestamp {
				b.blocks[j], b.blocks[j+1] = b.blocks[j+1], b.blocks[j]
			}
		}
	}

	// Calculate block times (time between consecutive blocks)
	for i := 1; i < len(b.blocks); i++ {
		timeDiff := b.blocks[i].timestamp - b.blocks[i-1].timestamp
		b.blocks[i].blockTime = time.Duration(timeDiff) * time.Second
	}

	// First block has no previous block, so no block time
	if len(b.blocks) > 0 {
		b.blocks[0].blockTime = 0
	}
}

// calculateAverage computes the average block time from the current window
func (b *BlockTimeMetric) calculateAverage() {
	if len(b.blocks) == 0 {
		b.averageTime = 0
		return
	}

	totalTime := time.Duration(0)
	validBlocks := 0

	// Only include blocks with non-zero block times (skip first block)
	for _, block := range b.blocks {
		if block.blockTime > 0 {
			totalTime += block.blockTime
			validBlocks++
		}
	}

	if validBlocks > 0 {
		b.averageTime = totalTime / time.Duration(validBlocks)
	} else {
		b.averageTime = 0
	}
}

// GetMetric returns the current block time statistics
func (b *BlockTimeMetric) GetMetric() any {
	b.mu.RLock()
	defer b.mu.RUnlock()

	// Calculate additional statistics (exclude zero block times)
	var minTime, maxTime time.Duration
	first := true
	validBlocks := 0

	for _, block := range b.blocks {
		if block.blockTime > 0 {
			validBlocks++
			if first {
				minTime = block.blockTime
				maxTime = block.blockTime
				first = false
			} else {
				if block.blockTime < minTime {
					minTime = block.blockTime
				}
				if block.blockTime > maxTime {
					maxTime = block.blockTime
				}
			}
		}
	}

	return BlockTimeStats{
		AverageBlockTime: b.averageTime,
		MinBlockTime:     minTime,
		MaxBlockTime:     maxTime,
		WindowSize:       validBlocks,
		MaxWindowSize:    b.windowSize,
	}
}

// GetUpdateInterval returns how often this metric should be updated
func (b *BlockTimeMetric) GetUpdateInterval() time.Duration {
	return 1 * time.Second
}

// BlockTimeStats provides detailed block time statistics
type BlockTimeStats struct {
	AverageBlockTime time.Duration // Average time between blocks
	MinBlockTime     time.Duration // Shortest block time in window
	MaxBlockTime     time.Duration // Longest block time in window
	WindowSize       int           // Current number of blocks tracked
	MaxWindowSize    int           // Maximum window size
}
