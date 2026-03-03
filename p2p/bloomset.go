package p2p

import (
	"encoding/binary"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	bloomfilter "github.com/holiman/bloomfilter/v2"
)

// BloomSetOptions contains configuration for creating a BloomSet.
type BloomSetOptions struct {
	// Size is the number of bits in the bloom filter.
	// Larger size = lower false positive rate but more memory.
	// Recommended: 10 * expected_elements for ~1% false positive rate.
	Size uint

	// HashCount is the number of hash functions to use.
	// Recommended: 7 for ~1% false positive rate.
	HashCount uint
}

// DefaultBloomSetOptions returns sensible defaults for tracking ~32K elements
// with approximately 1% false positive rate.
// Memory usage: ~80KB per BloomSet (2 filters of ~40KB each).
func DefaultBloomSetOptions() BloomSetOptions {
	return BloomSetOptions{
		Size:      327680, // 32768 * 10 bits ≈ 40KB per filter
		HashCount: 7,
	}
}

// BloomSet is a memory-efficient probabilistic set for tracking seen hashes.
// It uses a rotating dual-bloom-filter design:
//   - "current" filter receives all new additions
//   - "previous" filter is checked during lookups for recency
//   - Rotate() moves current to previous and creates a fresh current
//
// Trade-offs vs LRU cache:
//   - Pro: ~10x less memory, minimal GC pressure (fixed-size arrays)
//   - Pro: O(1) add/lookup with very low constant factor
//   - Con: False positives possible (~1% with default settings)
//   - Con: No exact eviction control (use Rotate for approximate TTL)
//
// For knownTxs, false positives mean occasionally not broadcasting a tx
// to a peer that doesn't have it - acceptable since they'll get it elsewhere.
//
// This implementation wraps holiman/bloomfilter/v2, the same battle-tested
// bloom filter library used by geth for state pruning.
type BloomSet struct {
	mu       sync.RWMutex
	current  *bloomfilter.Filter
	previous *bloomfilter.Filter
	m        uint64 // bits per filter
	k        uint64 // hash functions
}

// NewBloomSet creates a new BloomSet with the given options.
// If options are zero-valued, defaults are applied.
func NewBloomSet(opts BloomSetOptions) *BloomSet {
	defaults := DefaultBloomSetOptions()
	if opts.Size == 0 {
		opts.Size = defaults.Size
	}
	if opts.HashCount == 0 {
		opts.HashCount = defaults.HashCount
	}

	m := uint64(opts.Size)
	k := uint64(opts.HashCount)

	current, _ := bloomfilter.New(m, k)
	previous, _ := bloomfilter.New(m, k)

	return &BloomSet{
		current:  current,
		previous: previous,
		m:        m,
		k:        k,
	}
}

// bloomHash converts common.Hash to uint64 for the bloom filter.
// Uses first 8 bytes - sufficient since keccak256 hashes are already
// cryptographically distributed (same approach as geth).
func bloomHash(hash common.Hash) uint64 {
	return binary.BigEndian.Uint64(hash[:8])
}

// Add adds a hash to the set.
func (b *BloomSet) Add(hash common.Hash) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.current.AddHash(bloomHash(hash))
}

// AddMany adds multiple hashes to the set efficiently.
func (b *BloomSet) AddMany(hashes []common.Hash) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, hash := range hashes {
		b.current.AddHash(bloomHash(hash))
	}
}

// Contains checks if a hash might be in the set.
// Returns true if the hash is probably in the set (may have false positives).
// Returns false if the hash is definitely not in the set.
func (b *BloomSet) Contains(hash common.Hash) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	h := bloomHash(hash)
	return b.current.ContainsHash(h) || b.previous.ContainsHash(h)
}

// FilterNotContained returns hashes that are definitely not in the set.
// Hashes that might be in the set (including false positives) are excluded.
func (b *BloomSet) FilterNotContained(hashes []common.Hash) []common.Hash {
	b.mu.RLock()
	defer b.mu.RUnlock()

	result := make([]common.Hash, 0, len(hashes))
	for _, hash := range hashes {
		h := bloomHash(hash)
		if !b.current.ContainsHash(h) && !b.previous.ContainsHash(h) {
			result = append(result, hash)
		}
	}
	return result
}

// Rotate moves the current filter to previous and creates a fresh current.
// Call this periodically to maintain approximate recency (e.g., every N minutes).
// After rotation, lookups still check the previous filter, so recently-added
// items remain "known" for one more rotation period.
func (b *BloomSet) Rotate() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.previous = b.current
	b.current, _ = bloomfilter.New(b.m, b.k)
}

// Count returns the approximate number of elements added since last rotation.
// This uses the bloom filter's internal count of added elements.
func (b *BloomSet) Count() uint {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return uint(b.current.N())
}

// Reset clears both filters.
func (b *BloomSet) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.current, _ = bloomfilter.New(b.m, b.k)
	b.previous, _ = bloomfilter.New(b.m, b.k)
}

// MemoryUsage returns the approximate memory usage in bytes.
func (b *BloomSet) MemoryUsage() uint {
	b.mu.RLock()
	defer b.mu.RUnlock()

	// Two filters, each with m bits = m/8 bytes
	// Round up to account for uint64 alignment
	bytesPerFilter := (b.m + 63) / 64 * 8
	return uint(bytesPerFilter * 2)
}
