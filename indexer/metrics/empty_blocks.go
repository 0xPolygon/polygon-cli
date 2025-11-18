package metrics

import (
	"sync"
	"time"

	"github.com/0xPolygon/polygon-cli/rpctypes"
)

// EmptyBlockMetric tracks the rate of empty blocks
type EmptyBlockMetric struct {
	mu           sync.RWMutex
	totalBlocks  int
	emptyBlocks  int
	recentWindow []bool // true if block was empty
	windowSize   int    // number of blocks to track
}

// NewEmptyBlockMetric creates a new empty block rate calculator
func NewEmptyBlockMetric() *EmptyBlockMetric {
	return &EmptyBlockMetric{
		recentWindow: make([]bool, 0),
		windowSize:   100, // Track last 100 blocks
	}
}

// Name returns the metric identifier
func (e *EmptyBlockMetric) Name() string {
	return "emptyBlockRate"
}

// ProcessBlock processes a new block to update empty block statistics
func (e *EmptyBlockMetric) ProcessBlock(block rpctypes.PolyBlock) {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Check if block is empty
	isEmpty := len(block.Transactions()) == 0

	// Update totals
	e.totalBlocks++
	if isEmpty {
		e.emptyBlocks++
	}

	// Update recent window
	e.recentWindow = append(e.recentWindow, isEmpty)

	// Maintain window size
	if len(e.recentWindow) > e.windowSize {
		// Remove oldest entry
		// (Note: we don't decrement totalBlocks as that's cumulative)
		e.recentWindow = e.recentWindow[1:]
	}
}

// GetMetric returns the current empty block statistics
func (e *EmptyBlockMetric) GetMetric() any {
	e.mu.RLock()
	defer e.mu.RUnlock()

	// Calculate rates
	var overallRate, recentRate float64

	if e.totalBlocks > 0 {
		overallRate = float64(e.emptyBlocks) / float64(e.totalBlocks)
	}

	// Calculate recent rate from window
	recentEmpty := 0
	for _, isEmpty := range e.recentWindow {
		if isEmpty {
			recentEmpty++
		}
	}

	if len(e.recentWindow) > 0 {
		recentRate = float64(recentEmpty) / float64(len(e.recentWindow))
	}

	return EmptyBlockStats{
		TotalBlocks:      e.totalBlocks,
		EmptyBlocks:      e.emptyBlocks,
		OverallRate:      overallRate,
		RecentRate:       recentRate,
		RecentWindowSize: len(e.recentWindow),
	}
}

// GetUpdateInterval returns how often this metric should be updated
func (e *EmptyBlockMetric) GetUpdateInterval() time.Duration {
	return 1 * time.Second
}

// EmptyBlockStats provides detailed empty block statistics
type EmptyBlockStats struct {
	TotalBlocks      int     // Total blocks observed
	EmptyBlocks      int     // Total empty blocks
	OverallRate      float64 // Empty blocks / total blocks (all time)
	RecentRate       float64 // Empty blocks / total blocks (recent window)
	RecentWindowSize int     // Number of blocks in recent window
}
