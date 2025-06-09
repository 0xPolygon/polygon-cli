package metrics

import (
	"math/big"
	"sync"
	"time"

	"github.com/0xPolygon/polygon-cli/rpctypes"
)

// BaseFeeMetric calculates average base fees over block-count windows
type BaseFeeMetric struct {
	mu        sync.RWMutex
	blocks    []baseFeeBlockInfo
	maxBlocks int // 30 to cover both 10 and 30 block windows
}

type baseFeeBlockInfo struct {
	timestamp uint64
	baseFee   *big.Int
}

// NewBaseFeeMetric creates a new base fee calculator
func NewBaseFeeMetric() *BaseFeeMetric {
	return &BaseFeeMetric{
		blocks:    make([]baseFeeBlockInfo, 0),
		maxBlocks: 30, // Track last 30 blocks to cover both windows
	}
}

// Name returns the metric identifier
func (b *BaseFeeMetric) Name() string {
	return "basefee"
}

// ProcessBlock adds a new block to calculate base fee metrics
func (b *BaseFeeMetric) ProcessBlock(block rpctypes.PolyBlock) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Get base fee from block (may be nil for non-EIP-1559 chains)
	baseFee := block.BaseFee()
	if baseFee == nil {
		baseFee = big.NewInt(0)
	}

	// Add new block info to the front of the slice (newest first)
	info := baseFeeBlockInfo{
		timestamp: block.Time(),
		baseFee:   new(big.Int).Set(baseFee), // Copy to avoid mutation
	}

	// Prepend new block (newest first)
	b.blocks = append([]baseFeeBlockInfo{info}, b.blocks...)

	// Maintain window size (keep only last 30 blocks)
	if len(b.blocks) > b.maxBlocks {
		b.blocks = b.blocks[:b.maxBlocks]
	}
}

// GetMetric returns the current base fee statistics
func (b *BaseFeeMetric) GetMetric() interface{} {
	b.mu.RLock()
	defer b.mu.RUnlock()

	stats := BaseFeeStats{
		BaseFee10:       b.calculateAverageBaseFee(10),
		BaseFee30:       b.calculateAverageBaseFee(30),
		BlocksAvailable: len(b.blocks),
	}

	return stats
}

// calculateAverageBaseFee calculates average base fee over the specified number of blocks
func (b *BaseFeeMetric) calculateAverageBaseFee(blockCount int) *big.Int {
	if len(b.blocks) < blockCount || blockCount == 0 {
		return big.NewInt(0)
	}

	// Take the first blockCount blocks (newest to oldest)
	windowBlocks := b.blocks[:blockCount]

	// Calculate total base fee
	totalBaseFee := big.NewInt(0)
	for _, block := range windowBlocks {
		totalBaseFee.Add(totalBaseFee, block.baseFee)
	}

	// Calculate average
	avgBaseFee := new(big.Int).Div(totalBaseFee, big.NewInt(int64(blockCount)))

	return avgBaseFee
}

// GetUpdateInterval returns how often this metric should be updated
func (b *BaseFeeMetric) GetUpdateInterval() time.Duration {
	return 1 * time.Second
}

// BaseFeeStats provides detailed base fee statistics
type BaseFeeStats struct {
	BaseFee10       *big.Int // Average base fee (10 blocks)
	BaseFee30       *big.Int // Average base fee (30 blocks)
	BlocksAvailable int      // Number of blocks available for calculation
}
