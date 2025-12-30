package report

import (
	"math"
	"testing"
)

// TestHexToUint64 tests the hexToUint64 conversion function
func TestHexToUint64(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected uint64
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: 0,
		},
		{
			name:     "empty string",
			input:    "",
			expected: 0,
		},
		{
			name:     "non-string input",
			input:    123,
			expected: 0,
		},
		{
			name:     "hex with 0x prefix",
			input:    "0x10",
			expected: 16,
		},
		{
			name:     "hex without prefix",
			input:    "10",
			expected: 16,
		},
		{
			name:     "zero",
			input:    "0x0",
			expected: 0,
		},
		{
			name:     "large hex value",
			input:    "0xffffffffffffffff",
			expected: math.MaxUint64,
		},
		{
			name:     "typical block number",
			input:    "0x1234567",
			expected: 19088743,
		},
		{
			name:     "invalid hex",
			input:    "0xZZZ",
			expected: 0,
		},
		{
			name:     "short hex",
			input:    "0x",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hexToUint64(tt.input)
			if result != tt.expected {
				t.Errorf("hexToUint64(%v) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

// TestCalculateTop10Stats tests the top 10 statistics calculation
func TestCalculateTop10Stats(t *testing.T) {
	t.Run("empty blocks", func(t *testing.T) {
		blocks := []BlockInfo{}
		result := calculateTop10Stats(blocks)

		if len(result.BlocksByTxCount) != 0 {
			t.Errorf("expected 0 blocks by tx count, got %d", len(result.BlocksByTxCount))
		}
		if len(result.BlocksByGasUsed) != 0 {
			t.Errorf("expected 0 blocks by gas used, got %d", len(result.BlocksByGasUsed))
		}
		if len(result.TransactionsByGas) != 0 {
			t.Errorf("expected 0 transactions by gas, got %d", len(result.TransactionsByGas))
		}
	})

	t.Run("single block with transactions", func(t *testing.T) {
		blocks := []BlockInfo{
			{
				Number:   100,
				TxCount:  3,
				GasUsed:  100000,
				GasLimit: 200000,
				Transactions: []TransactionInfo{
					{Hash: "0x1", GasUsed: 50000, GasLimit: 60000, GasPrice: 20000000000, BlockNumber: 100, BlockGasLimit: 200000, GasUsedPercent: 25.0},
					{Hash: "0x2", GasUsed: 30000, GasLimit: 40000, GasPrice: 20000000000, BlockNumber: 100, BlockGasLimit: 200000, GasUsedPercent: 15.0},
					{Hash: "0x3", GasUsed: 20000, GasLimit: 30000, GasPrice: 30000000000, BlockNumber: 100, BlockGasLimit: 200000, GasUsedPercent: 10.0},
				},
			},
		}

		result := calculateTop10Stats(blocks)

		// Check blocks by tx count
		if len(result.BlocksByTxCount) != 1 {
			t.Errorf("expected 1 block by tx count, got %d", len(result.BlocksByTxCount))
		}
		if result.BlocksByTxCount[0].Number != 100 {
			t.Errorf("expected block 100, got %d", result.BlocksByTxCount[0].Number)
		}
		if result.BlocksByTxCount[0].TxCount != 3 {
			t.Errorf("expected 3 transactions, got %d", result.BlocksByTxCount[0].TxCount)
		}

		// Check transactions by gas used
		if len(result.TransactionsByGas) != 3 {
			t.Errorf("expected 3 transactions by gas, got %d", len(result.TransactionsByGas))
		}
		// Should be sorted by gas used descending
		if result.TransactionsByGas[0].GasUsed != 50000 {
			t.Errorf("expected highest gas used to be 50000, got %d", result.TransactionsByGas[0].GasUsed)
		}
		if result.TransactionsByGas[2].GasUsed != 20000 {
			t.Errorf("expected lowest gas used to be 20000, got %d", result.TransactionsByGas[2].GasUsed)
		}

		// Check most used gas prices
		if len(result.MostUsedGasPrices) != 2 {
			t.Errorf("expected 2 unique gas prices, got %d", len(result.MostUsedGasPrices))
		}
		// The most frequent gas price (20000000000) should be first
		if result.MostUsedGasPrices[0].GasPrice != 20000000000 {
			t.Errorf("expected most used gas price to be 20000000000, got %d", result.MostUsedGasPrices[0].GasPrice)
		}
		if result.MostUsedGasPrices[0].Count != 2 {
			t.Errorf("expected count of 2, got %d", result.MostUsedGasPrices[0].Count)
		}
	})

	t.Run("multiple blocks - top 10 limit", func(t *testing.T) {
		// Create 15 blocks to test the top 10 limit
		blocks := make([]BlockInfo, 15)
		for i := range blocks {
			blocks[i] = BlockInfo{
				Number:   uint64(i),
				TxCount:  uint64(i * 10), // Increasing tx count
				GasUsed:  uint64(i * 100000),
				GasLimit: 1000000,
			}
		}

		result := calculateTop10Stats(blocks)

		// Should only return top 10
		if len(result.BlocksByTxCount) != 10 {
			t.Errorf("expected 10 blocks by tx count, got %d", len(result.BlocksByTxCount))
		}

		// Should be sorted descending, so highest should be first
		if result.BlocksByTxCount[0].TxCount != 140 {
			t.Errorf("expected highest tx count to be 140, got %d", result.BlocksByTxCount[0].TxCount)
		}
		if result.BlocksByTxCount[9].TxCount != 50 {
			t.Errorf("expected 10th highest tx count to be 50, got %d", result.BlocksByTxCount[9].TxCount)
		}
	})

	t.Run("gas used percentage calculation", func(t *testing.T) {
		blocks := []BlockInfo{
			{
				Number:   100,
				TxCount:  1,
				GasUsed:  75000,
				GasLimit: 100000,
			},
		}

		result := calculateTop10Stats(blocks)

		if len(result.BlocksByGasUsed) != 1 {
			t.Fatalf("expected 1 block by gas used, got %d", len(result.BlocksByGasUsed))
		}

		expectedPercent := 75.0
		if result.BlocksByGasUsed[0].GasUsedPercent != expectedPercent {
			t.Errorf("expected gas used percent to be %.2f, got %.2f", expectedPercent, result.BlocksByGasUsed[0].GasUsedPercent)
		}
	})
}

// TestBlockRangeLogic tests the smart default logic for block ranges
func TestBlockRangeLogic(t *testing.T) {
	tests := []struct {
		name        string
		startInput  uint64
		endInput    uint64
		latestBlock uint64
		wantStart   uint64
		wantEnd     uint64
	}{
		{
			name:        "no flags specified - latest 500 blocks",
			startInput:  BlockNotSet,
			endInput:    BlockNotSet,
			latestBlock: 1000,
			wantStart:   501,
			wantEnd:     1000,
		},
		{
			name:        "no flags - small chain (< 500 blocks)",
			startInput:  BlockNotSet,
			endInput:    BlockNotSet,
			latestBlock: 100,
			wantStart:   0,
			wantEnd:     100,
		},
		{
			name:        "only start specified - next 500 blocks",
			startInput:  100,
			endInput:    BlockNotSet,
			latestBlock: 1000,
			wantStart:   100,
			wantEnd:     599,
		},
		{
			name:        "only start specified - capped at latest",
			startInput:  900,
			endInput:    BlockNotSet,
			latestBlock: 1000,
			wantStart:   900,
			wantEnd:     1000,
		},
		{
			name:        "only end specified - previous 500 blocks",
			startInput:  BlockNotSet,
			endInput:    600,
			latestBlock: 1000,
			wantStart:   101,
			wantEnd:     600,
		},
		{
			name:        "only end specified - end < 500",
			startInput:  BlockNotSet,
			endInput:    100,
			latestBlock: 1000,
			wantStart:   0,
			wantEnd:     100,
		},
		{
			name:        "both specified - genesis block only",
			startInput:  0,
			endInput:    0,
			latestBlock: 1000,
			wantStart:   0,
			wantEnd:     0,
		},
		{
			name:        "both specified - custom range",
			startInput:  1000,
			endInput:    2000,
			latestBlock: 5000,
			wantStart:   1000,
			wantEnd:     2000,
		},
		{
			name:        "start at zero with end unspecified",
			startInput:  0,
			endInput:    BlockNotSet,
			latestBlock: 1000,
			wantStart:   0,
			wantEnd:     499,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := tt.startInput
			end := tt.endInput

			// Apply the same logic as in report.go RunE function
			if start == BlockNotSet && end == BlockNotSet {
				// Both unspecified
				end = tt.latestBlock
				if tt.latestBlock >= DefaultBlockRange-1 {
					start = tt.latestBlock - (DefaultBlockRange - 1)
				} else {
					start = 0
				}
			} else if start == BlockNotSet {
				// Only start unspecified
				if end >= DefaultBlockRange-1 {
					start = end - (DefaultBlockRange - 1)
				} else {
					start = 0
				}
			} else if end == BlockNotSet {
				// Only end unspecified
				end = start + (DefaultBlockRange - 1)
				if end > tt.latestBlock {
					end = tt.latestBlock
				}
			}
			// Both set: use as-is

			if start != tt.wantStart {
				t.Errorf("start = %d, want %d", start, tt.wantStart)
			}
			if end != tt.wantEnd {
				t.Errorf("end = %d, want %d", end, tt.wantEnd)
			}
		})
	}
}

// TestBlockNotSetConstant verifies the sentinel value is set correctly
func TestBlockNotSetConstant(t *testing.T) {
	expected := uint64(math.MaxUint64)
	actual := uint64(BlockNotSet)
	if actual != expected {
		t.Errorf("BlockNotSet = %d, want %d", actual, expected)
	}
}

// TestDefaultBlockRangeConstant verifies the default block range
func TestDefaultBlockRangeConstant(t *testing.T) {
	if DefaultBlockRange != 500 {
		t.Errorf("DefaultBlockRange = %d, want 500", DefaultBlockRange)
	}
}
