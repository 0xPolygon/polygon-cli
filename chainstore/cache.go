package chainstore

import (
	"math/big"
	"sync"
	"time"
)

// CachedValue represents a cached value with TTL
type CachedValue[T any] struct {
	value     T
	timestamp time.Time
	ttl       time.Duration
}

// NewCachedValue creates a new cached value
func NewCachedValue[T any](value T, ttl time.Duration) *CachedValue[T] {
	return &CachedValue[T]{
		value:     value,
		timestamp: time.Now(),
		ttl:       ttl,
	}
}

// IsValid returns true if the cached value is still valid
func (cv *CachedValue[T]) IsValid() bool {
	if cv.ttl == 0 {
		return true // Never expires
	}
	return time.Since(cv.timestamp) < cv.ttl
}

// Get returns the cached value and whether it's valid
func (cv *CachedValue[T]) Get() (T, bool) {
	if cv == nil {
		var zero T
		return zero, false
	}
	return cv.value, cv.IsValid()
}

// ChainCache manages cached chain information with different TTLs
type ChainCache struct {
	mu sync.RWMutex
	
	// Static data (never expires)
	chainID *CachedValue[*big.Int]
	
	// Semi-static data (5-15 minute TTL)
	safeBlock      *CachedValue[*big.Int]
	finalizedBlock *CachedValue[*big.Int]
	
	// Block-aligned data (expires when new block)
	baseFee      *CachedValue[*big.Int]
	baseFeeBlock *big.Int
	
	// Frequent data (30-60 second TTL)
	gasPrice   *CachedValue[*big.Int]
	feeHistory *CachedValue[*FeeHistoryResult]
	
	// Very frequent data (5-10 second TTL)
	pendingTxCount *CachedValue[*big.Int]
	queuedTxCount  *CachedValue[*big.Int]
}

// NewChainCache creates a new chain cache
func NewChainCache() *ChainCache {
	return &ChainCache{}
}

// GetChainID gets cached chain ID
func (cc *ChainCache) GetChainID() (*big.Int, bool) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	if cc.chainID == nil {
		return nil, false
	}
	return cc.chainID.Get()
}

// SetChainID caches chain ID (never expires)
func (cc *ChainCache) SetChainID(chainID *big.Int) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc.chainID = NewCachedValue(chainID, 0) // Never expires
}

// GetSafeBlock gets cached safe block
func (cc *ChainCache) GetSafeBlock(ttl time.Duration) (*big.Int, bool) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	if cc.safeBlock == nil {
		return nil, false
	}
	return cc.safeBlock.Get()
}

// SetSafeBlock caches safe block
func (cc *ChainCache) SetSafeBlock(block *big.Int, ttl time.Duration) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc.safeBlock = NewCachedValue(block, ttl)
}

// GetFinalizedBlock gets cached finalized block
func (cc *ChainCache) GetFinalizedBlock(ttl time.Duration) (*big.Int, bool) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	if cc.finalizedBlock == nil {
		return nil, false
	}
	return cc.finalizedBlock.Get()
}

// SetFinalizedBlock caches finalized block
func (cc *ChainCache) SetFinalizedBlock(block *big.Int, ttl time.Duration) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc.finalizedBlock = NewCachedValue(block, ttl)
}

// GetBaseFee gets cached base fee
func (cc *ChainCache) GetBaseFee(currentBlock *big.Int) (*big.Int, bool) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	if cc.baseFee == nil || cc.baseFeeBlock == nil {
		return nil, false
	}
	// Base fee is valid only for the same block
	if currentBlock != nil && cc.baseFeeBlock.Cmp(currentBlock) != 0 {
		return nil, false
	}
	return cc.baseFee.Get()
}

// SetBaseFee caches base fee for a specific block
func (cc *ChainCache) SetBaseFee(baseFee *big.Int, blockNumber *big.Int) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc.baseFee = NewCachedValue(baseFee, time.Hour) // Long TTL since it's block-aligned
	cc.baseFeeBlock = new(big.Int).Set(blockNumber)
}

// GetGasPrice gets cached gas price
func (cc *ChainCache) GetGasPrice(ttl time.Duration) (*big.Int, bool) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	if cc.gasPrice == nil {
		return nil, false
	}
	return cc.gasPrice.Get()
}

// SetGasPrice caches gas price
func (cc *ChainCache) SetGasPrice(gasPrice *big.Int, ttl time.Duration) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc.gasPrice = NewCachedValue(gasPrice, ttl)
}

// GetFeeHistory gets cached fee history
func (cc *ChainCache) GetFeeHistory(ttl time.Duration) (*FeeHistoryResult, bool) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	if cc.feeHistory == nil {
		return nil, false
	}
	return cc.feeHistory.Get()
}

// SetFeeHistory caches fee history
func (cc *ChainCache) SetFeeHistory(feeHistory *FeeHistoryResult, ttl time.Duration) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc.feeHistory = NewCachedValue(feeHistory, ttl)
}

// GetPendingTxCount gets cached pending transaction count
func (cc *ChainCache) GetPendingTxCount(ttl time.Duration) (*big.Int, bool) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	if cc.pendingTxCount == nil {
		return nil, false
	}
	return cc.pendingTxCount.Get()
}

// SetPendingTxCount caches pending transaction count
func (cc *ChainCache) SetPendingTxCount(count *big.Int, ttl time.Duration) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc.pendingTxCount = NewCachedValue(count, ttl)
}

// GetQueuedTxCount gets cached queued transaction count
func (cc *ChainCache) GetQueuedTxCount(ttl time.Duration) (*big.Int, bool) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	if cc.queuedTxCount == nil {
		return nil, false
	}
	return cc.queuedTxCount.Get()
}

// SetQueuedTxCount caches queued transaction count
func (cc *ChainCache) SetQueuedTxCount(count *big.Int, ttl time.Duration) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc.queuedTxCount = NewCachedValue(count, ttl)
}