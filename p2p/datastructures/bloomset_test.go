package datastructures

import (
	"encoding/binary"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestBloomSet(t *testing.T) {
	t.Run("AddAndContains", func(t *testing.T) {
		b := NewBloomSet(DefaultBloomSetOptions())

		hash1 := common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
		hash2 := common.HexToHash("0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890")

		// Initially should not contain anything
		if b.Contains(hash1) {
			t.Error("expected hash1 not to be contained initially")
		}

		// Add hash1
		b.Add(hash1)

		// Should contain hash1
		if !b.Contains(hash1) {
			t.Error("expected hash1 to be contained after add")
		}

		// Should not contain hash2
		if b.Contains(hash2) {
			t.Error("expected hash2 not to be contained")
		}
	})

	t.Run("AddMany", func(t *testing.T) {
		b := NewBloomSet(DefaultBloomSetOptions())

		hashes := make([]common.Hash, 100)
		for i := range hashes {
			hashes[i] = common.BytesToHash([]byte{byte(i), byte(i + 1), byte(i + 2)})
		}

		b.AddMany(hashes)

		// All added hashes should be contained
		for i, hash := range hashes {
			if !b.Contains(hash) {
				t.Errorf("expected hash %d to be contained", i)
			}
		}
	})

	t.Run("FilterNotContained", func(t *testing.T) {
		b := NewBloomSet(DefaultBloomSetOptions())

		// Add some hashes
		known := []common.Hash{
			common.HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111"),
			common.HexToHash("0x2222222222222222222222222222222222222222222222222222222222222222"),
		}
		b.AddMany(known)

		// Create a mixed list
		unknown := []common.Hash{
			common.HexToHash("0x3333333333333333333333333333333333333333333333333333333333333333"),
			common.HexToHash("0x4444444444444444444444444444444444444444444444444444444444444444"),
		}
		mixed := append(known, unknown...)

		// Filter should return only unknown hashes
		result := b.FilterNotContained(mixed)

		if len(result) != 2 {
			t.Errorf("expected 2 unknown hashes, got %d", len(result))
		}

		for _, h := range result {
			if h == known[0] || h == known[1] {
				t.Errorf("known hash %s should not be in result", h.Hex())
			}
		}
	})

	t.Run("Rotate", func(t *testing.T) {
		b := NewBloomSet(DefaultBloomSetOptions())

		hash1 := common.HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111")
		hash2 := common.HexToHash("0x2222222222222222222222222222222222222222222222222222222222222222")

		// Add hash1 to current
		b.Add(hash1)

		// Rotate - hash1 moves to previous
		b.Rotate()

		// hash1 should still be found (in previous)
		if !b.Contains(hash1) {
			t.Error("expected hash1 to be found after first rotation")
		}

		// Add hash2 to current
		b.Add(hash2)

		// Rotate again - hash1's filter is now cleared, hash2 moves to previous
		b.Rotate()

		// hash1 should no longer be found
		if b.Contains(hash1) {
			t.Error("expected hash1 not to be found after second rotation")
		}

		// hash2 should still be found
		if !b.Contains(hash2) {
			t.Error("expected hash2 to be found after second rotation")
		}
	})

	t.Run("MemoryUsage", func(t *testing.T) {
		opts := DefaultBloomSetOptions()
		b := NewBloomSet(opts)

		usage := b.MemoryUsage()

		// With default 327680 bits = 5120 words per filter * 8 bytes * 2 filters = 81920 bytes
		expectedBytes := uint((327680+63)/64) * 8 * 2
		if usage != expectedBytes {
			t.Errorf("expected memory usage %d bytes, got %d", expectedBytes, usage)
		}
	})

	t.Run("FalsePositiveRate", func(t *testing.T) {
		// Test that false positive rate is approximately as expected
		b := NewBloomSet(DefaultBloomSetOptions())

		// Add 32768 unique hashes (the design capacity)
		// Use keccak256 to generate properly distributed hashes
		for i := uint64(0); i < 32768; i++ {
			b.Add(generateTestHash(i))
		}

		// Test 10000 hashes that were NOT added (different seed range)
		falsePositives := 0
		for i := uint64(100000); i < 110000; i++ {
			if b.Contains(generateTestHash(i)) {
				falsePositives++
			}
		}

		// Expected ~1% false positive rate, allow up to 3% for statistical variance
		rate := float64(falsePositives) / 10000.0
		if rate > 0.03 {
			t.Errorf("false positive rate too high: %.2f%% (expected ~1%%)", rate*100)
		}
		t.Logf("False positive rate: %.2f%%", rate*100)
	})
}

// generateTestHash creates a deterministic hash from a seed using keccak256.
func generateTestHash(seed uint64) common.Hash {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], seed)
	return crypto.Keccak256Hash(buf[:])
}

func BenchmarkBloomSetAdd(b *testing.B) {
	bloom := NewBloomSet(DefaultBloomSetOptions())
	hash := common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")

	for b.Loop() {
		bloom.Add(hash)
	}
}

func BenchmarkBloomSetContains(b *testing.B) {
	bloom := NewBloomSet(DefaultBloomSetOptions())
	hash := common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	bloom.Add(hash)

	for b.Loop() {
		bloom.Contains(hash)
	}
}

func BenchmarkBloomSetFilterNotContained(b *testing.B) {
	bloom := NewBloomSet(DefaultBloomSetOptions())

	// Add 1000 hashes
	for i := range 1000 {
		hash := common.BytesToHash([]byte{byte(i >> 8), byte(i)})
		bloom.Add(hash)
	}

	// Create a batch of 100 hashes (mix of known and unknown)
	batch := make([]common.Hash, 100)
	for i := range batch {
		batch[i] = common.BytesToHash([]byte{byte(i >> 8), byte(i)})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bloom.FilterNotContained(batch)
	}
}

func BenchmarkLRUFilterNotContained(b *testing.B) {
	cache := NewLRU[common.Hash, struct{}](LRUOptions{MaxSize: 32768})

	// Add 1000 hashes
	for i := range 1000 {
		hash := common.BytesToHash([]byte{byte(i >> 8), byte(i)})
		cache.Add(hash, struct{}{})
	}

	// Create a batch of 100 hashes (mix of known and unknown)
	batch := make([]common.Hash, 100)
	for i := range batch {
		batch[i] = common.BytesToHash([]byte{byte(i >> 8), byte(i)})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.FilterNotContained(batch)
	}
}
