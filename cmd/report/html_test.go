package report

import (
	"math/big"
	"strings"
	"testing"
	"time"
)

// TestFormatNumber tests the formatNumber function for comma separation
func TestFormatNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    uint64
		expected string
	}{
		{
			name:     "zero",
			input:    0,
			expected: "0",
		},
		{
			name:     "single digit",
			input:    5,
			expected: "5",
		},
		{
			name:     "three digits",
			input:    999,
			expected: "999",
		},
		{
			name:     "four digits",
			input:    1000,
			expected: "1,000",
		},
		{
			name:     "five digits",
			input:    12345,
			expected: "12,345",
		},
		{
			name:     "six digits",
			input:    123456,
			expected: "123,456",
		},
		{
			name:     "seven digits",
			input:    1234567,
			expected: "1,234,567",
		},
		{
			name:     "million",
			input:    1000000,
			expected: "1,000,000",
		},
		{
			name:     "large number",
			input:    1234567890,
			expected: "1,234,567,890",
		},
		{
			name:     "max uint64",
			input:    18446744073709551615,
			expected: "18,446,744,073,709,551,615",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatNumber(tt.input)
			if result != tt.expected {
				t.Errorf("formatNumber(%d) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

// TestFormatNumberWithUnits tests the formatNumberWithUnits function
func TestFormatNumberWithUnits(t *testing.T) {
	tests := []struct {
		name     string
		input    uint64
		expected string
	}{
		{
			name:     "zero",
			input:    0,
			expected: "0",
		},
		{
			name:     "hundred",
			input:    100,
			expected: "100",
		},
		{
			name:     "thousand - uses 2 decimals for values < 10",
			input:    1500,
			expected: "1.50K",
		},
		{
			name:     "million - uses 2 decimals for values < 10",
			input:    2500000,
			expected: "2.50M",
		},
		{
			name:     "billion - uses 2 decimals for values < 10",
			input:    3500000000,
			expected: "3.50B",
		},
		{
			name:     "trillion - uses 2 decimals for values < 10",
			input:    4500000000000,
			expected: "4.50T",
		},
		{
			name:     "large value - uses 1 decimal for values >= 10",
			input:    15000000,
			expected: "15.0M",
		},
		{
			name:     "very large value - uses 0 decimals for values >= 100",
			input:    150000000,
			expected: "150M",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatNumberWithUnits(tt.input)
			if result != tt.expected {
				t.Errorf("formatNumberWithUnits(%d) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

// TestGenerateStatCards tests the stat cards generation
func TestGenerateStatCards(t *testing.T) {
	report := &BlockReport{
		Summary: SummaryStats{
			TotalBlocks:       100,
			TotalTransactions: 5000,
			UniqueSenders:     250,
			UniqueRecipients:  300,
			AvgTxPerBlock:     50.0,
			TotalGasUsed:      10000000,
			AvgGasPerBlock:    100000.0,
			AvgBaseFeePerGas:  "20000000000", // 20 Gwei
		},
	}

	result := generateStatCards(report)

	// Check that result contains expected elements
	if !strings.Contains(result, "Total Blocks") {
		t.Error("expected stat cards to contain 'Total Blocks'")
	}
	if !strings.Contains(result, "100") {
		t.Error("expected stat cards to contain '100' for total blocks")
	}
	if !strings.Contains(result, "5,000") {
		t.Error("expected stat cards to contain '5,000' for total transactions")
	}
	if !strings.Contains(result, "Avg Base Fee (Gwei)") {
		t.Error("expected stat cards to contain 'Avg Base Fee (Gwei)'")
	}
	if !strings.Contains(result, "20.00") {
		t.Error("expected stat cards to contain '20.00' for avg base fee in Gwei")
	}
}

// TestGenerateStatCardsWithoutBaseFee tests stat cards when base fee is not available
func TestGenerateStatCardsWithoutBaseFee(t *testing.T) {
	report := &BlockReport{
		Summary: SummaryStats{
			TotalBlocks:       100,
			TotalTransactions: 5000,
			AvgBaseFeePerGas:  "", // No base fee
		},
	}

	result := generateStatCards(report)

	// Check that base fee card is not included
	if strings.Contains(result, "Avg Base Fee") {
		t.Error("expected stat cards not to contain 'Avg Base Fee' when not available")
	}
}

// TestGenerateTxCountChart tests transaction count chart generation
func TestGenerateTxCountChart(t *testing.T) {
	t.Run("empty blocks", func(t *testing.T) {
		report := &BlockReport{
			Blocks: []BlockInfo{},
		}

		result := generateTxCountChart(report)

		if result != "" {
			t.Error("expected empty string for empty blocks")
		}
	})

	t.Run("single block", func(t *testing.T) {
		report := &BlockReport{
			Blocks: []BlockInfo{
				{Number: 100, TxCount: 50},
			},
		}

		result := generateTxCountChart(report)

		// Should handle single block without division by zero
		if !strings.Contains(result, "Transaction Count by Block") {
			t.Error("expected chart to contain title")
		}
		if !strings.Contains(result, "<svg") {
			t.Error("expected chart to contain SVG")
		}
	})

	t.Run("multiple blocks", func(t *testing.T) {
		report := &BlockReport{
			Blocks: []BlockInfo{
				{Number: 100, TxCount: 10},
				{Number: 101, TxCount: 20},
				{Number: 102, TxCount: 30},
			},
		}

		result := generateTxCountChart(report)

		if !strings.Contains(result, "Transaction Count by Block") {
			t.Error("expected chart to contain title")
		}
		if !strings.Contains(result, "<svg") {
			t.Error("expected chart to contain SVG")
		}
	})
}

// TestGenerateGasUsageChart tests gas usage chart generation
func TestGenerateGasUsageChart(t *testing.T) {
	t.Run("empty blocks", func(t *testing.T) {
		report := &BlockReport{
			Blocks: []BlockInfo{},
		}

		result := generateGasUsageChart(report)

		if result != "" {
			t.Error("expected empty string for empty blocks")
		}
	})

	t.Run("single block", func(t *testing.T) {
		report := &BlockReport{
			Blocks: []BlockInfo{
				{Number: 100, GasUsed: 1000000},
			},
		}

		result := generateGasUsageChart(report)

		// Should handle single block without division by zero
		if !strings.Contains(result, "Gas Usage by Block") {
			t.Error("expected chart to contain title")
		}
		if !strings.Contains(result, "<svg") {
			t.Error("expected chart to contain SVG")
		}
	})

	t.Run("multiple blocks", func(t *testing.T) {
		report := &BlockReport{
			Blocks: []BlockInfo{
				{Number: 100, GasUsed: 1000000},
				{Number: 101, GasUsed: 2000000},
				{Number: 102, GasUsed: 3000000},
			},
		}

		result := generateGasUsageChart(report)

		if !strings.Contains(result, "Gas Usage by Block") {
			t.Error("expected chart to contain title")
		}
		if !strings.Contains(result, "<svg") {
			t.Error("expected chart to contain SVG")
		}
	})
}

// TestGenerateTop10Sections tests top 10 sections generation
func TestGenerateTop10Sections(t *testing.T) {
	report := &BlockReport{
		Top10: Top10Stats{
			BlocksByTxCount: []TopBlock{
				{Number: 100, TxCount: 50},
				{Number: 101, TxCount: 45},
			},
			BlocksByGasUsed: []TopBlock{
				{Number: 100, GasUsed: 1000000, GasLimit: 2000000, GasUsedPercent: 50.0},
			},
			TransactionsByGas: []TopTransaction{
				{Hash: "0xabc123", BlockNumber: 100, GasUsed: 50000, BlockGasLimit: 1000000, GasUsedPercent: 5.0},
			},
			TransactionsByGasLimit: []TopTransaction{
				{Hash: "0xdef456", BlockNumber: 101, GasLimit: 100000, GasUsed: 80000},
			},
			MostUsedGasPrices: []GasPriceFreq{
				{GasPrice: 20000000000, Count: 100},
			},
			MostUsedGasLimits: []GasLimitFreq{
				{GasLimit: 21000, Count: 50},
			},
		},
	}

	result := generateTop10Sections(report)

	// Check for expected content
	if !strings.Contains(result, "Top 10 Blocks by Transaction Count") {
		t.Error("expected top 10 sections to contain 'Top 10 Blocks by Transaction Count'")
	}
	if !strings.Contains(result, "Top 10 Blocks by Gas Used") {
		t.Error("expected top 10 sections to contain 'Top 10 Blocks by Gas Used'")
	}
	if !strings.Contains(result, "Top 10 Transactions by Gas Used") {
		t.Error("expected top 10 sections to contain 'Top 10 Transactions by Gas Used'")
	}
}

// TestGenerateHTMLEscaping tests that HTML generation properly escapes user-controlled data
func TestGenerateHTMLEscaping(t *testing.T) {
	// Create a report with potentially malicious data
	report := &BlockReport{
		ChainID:     1,
		RPCURL:      "<script>alert('xss')</script>",
		StartBlock:  0,
		EndBlock:    100,
		GeneratedAt: time.Now(),
		Summary: SummaryStats{
			TotalBlocks: 100,
		},
		Blocks: []BlockInfo{
			{
				Number:  100,
				TxCount: 1,
				Transactions: []TransactionInfo{
					{
						Hash:        "<script>alert('xss')</script>",
						BlockNumber: 100,
						GasUsed:     21000,
					},
				},
			},
		},
		Top10: Top10Stats{
			TransactionsByGas: []TopTransaction{
				{
					Hash:        "<script>alert('xss')</script>",
					BlockNumber: 100,
					GasUsed:     21000,
				},
			},
			TransactionsByGasLimit: []TopTransaction{
				{
					Hash:        "<script>alert('xss')</script>",
					BlockNumber: 100,
					GasLimit:    21000,
				},
			},
		},
	}

	result := generateHTML(report)

	// Check that malicious script tags are escaped
	if strings.Contains(result, "<script>alert('xss')</script>") {
		t.Error("expected HTML to escape script tags in RPC URL or transaction hashes")
	}

	// Check that escaped version is present
	if !strings.Contains(result, "&lt;script&gt;") {
		t.Error("expected HTML to contain escaped script tags")
	}
}

// TestGenerateHTMLWithBaseFee tests HTML generation with base fee data
func TestGenerateHTMLWithBaseFee(t *testing.T) {
	baseFee := new(big.Int)
	baseFee.SetString("20000000000", 10) // 20 Gwei

	report := &BlockReport{
		ChainID:     1,
		RPCURL:      "http://localhost:8545",
		StartBlock:  0,
		EndBlock:    100,
		GeneratedAt: time.Now(),
		Summary: SummaryStats{
			TotalBlocks:      100,
			AvgBaseFeePerGas: "20000000000",
		},
		Blocks: []BlockInfo{
			{
				Number:        100,
				TxCount:       1,
				GasUsed:       21000,
				GasLimit:      30000000,
				BaseFeePerGas: baseFee,
			},
		},
	}

	result := generateHTML(report)

	// Check that base fee is included
	if !strings.Contains(result, "20.00") {
		t.Error("expected HTML to contain base fee in Gwei (20.00)")
	}
}

// TestGenerateHTMLWithoutBaseFee tests HTML generation without base fee data (pre-EIP-1559)
func TestGenerateHTMLWithoutBaseFee(t *testing.T) {
	report := &BlockReport{
		ChainID:     1,
		RPCURL:      "http://localhost:8545",
		StartBlock:  0,
		EndBlock:    100,
		GeneratedAt: time.Now(),
		Summary: SummaryStats{
			TotalBlocks:      100,
			AvgBaseFeePerGas: "", // No base fee
		},
		Blocks: []BlockInfo{
			{
				Number:        100,
				TxCount:       1,
				GasUsed:       21000,
				GasLimit:      30000000,
				BaseFeePerGas: nil,
			},
		},
	}

	result := generateHTML(report)

	// Check that base fee section is not included
	if strings.Contains(result, "Avg Base Fee (Gwei)") {
		t.Error("expected HTML not to contain base fee when not available")
	}
}
