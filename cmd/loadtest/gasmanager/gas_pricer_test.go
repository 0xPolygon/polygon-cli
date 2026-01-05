package gasmanager

import (
	"math"
	"testing"
)

func TestNewGasPricer(t *testing.T) {
	strategy := NewEstimatedGasPriceStrategy()
	pricer := NewGasPricer(strategy)

	if pricer == nil {
		t.Fatal("NewGasPricer returned nil")
	}
	if pricer.strategy == nil {
		t.Fatal("GasPricer strategy is nil")
	}
}

func TestGasPricer_GetGasPrice(t *testing.T) {
	fixedPrice := uint64(1000000000) // 1 Gwei
	strategy := NewFixedGasPriceStrategy(FixedGasPriceConfig{
		GasPriceWei: fixedPrice,
	})
	pricer := NewGasPricer(strategy)

	price := pricer.GetGasPrice()
	if price == nil {
		t.Fatal("GetGasPrice returned nil")
	}
	if *price != fixedPrice {
		t.Errorf("Expected gas price %d, got %d", fixedPrice, *price)
	}
}

// Estimated Gas Price Strategy Tests

func TestEstimatedGasPriceStrategy(t *testing.T) {
	strategy := NewEstimatedGasPriceStrategy()
	if strategy == nil {
		t.Fatal("NewEstimatedGasPriceStrategy returned nil")
	}
}

func TestEstimatedGasPriceStrategy_GetGasPrice(t *testing.T) {
	strategy := NewEstimatedGasPriceStrategy()
	price := strategy.GetGasPrice()

	// Estimated strategy should return nil to indicate network price should be used
	if price != nil {
		t.Errorf("Expected nil (network-estimated price), got %v", price)
	}
}

// Fixed Gas Price Strategy Tests

func TestFixedGasPriceStrategy(t *testing.T) {
	tests := []struct {
		name     string
		gasPrice uint64
	}{
		{"1 Gwei", 1000000000},
		{"50 Gwei", 50000000000},
		{"100 Gwei", 100000000000},
		{"zero", 0},
		{"max uint64", math.MaxUint64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strategy := NewFixedGasPriceStrategy(FixedGasPriceConfig{
				GasPriceWei: tt.gasPrice,
			})

			if strategy == nil {
				t.Fatal("NewFixedGasPriceStrategy returned nil")
			}

			price := strategy.GetGasPrice()
			if price == nil {
				t.Fatal("GetGasPrice returned nil")
			}
			if *price != tt.gasPrice {
				t.Errorf("Expected gas price %d, got %d", tt.gasPrice, *price)
			}
		})
	}
}

func TestFixedGasPriceStrategy_ConsistentPrice(t *testing.T) {
	fixedPrice := uint64(25000000000) // 25 Gwei
	strategy := NewFixedGasPriceStrategy(FixedGasPriceConfig{
		GasPriceWei: fixedPrice,
	})

	// Call multiple times and verify price is always the same
	for i := 0; i < 100; i++ {
		price := strategy.GetGasPrice()
		if price == nil {
			t.Fatal("GetGasPrice returned nil")
		}
		if *price != fixedPrice {
			t.Errorf("Call %d: expected gas price %d, got %d", i, fixedPrice, *price)
		}
	}
}

// Dynamic Gas Price Strategy Tests

func TestDynamicGasPriceStrategy_InvalidConfig(t *testing.T) {
	tests := []struct {
		name      string
		gasPrices []uint64
	}{
		{"empty slice", []uint64{}},
		{"nil slice", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strategy, err := NewDynamicGasPriceStrategy(DynamicGasPriceConfig{
				GasPrices: tt.gasPrices,
				Variation: 0.3,
			})

			if err == nil {
				t.Error("Expected error for empty gas prices, got nil")
			}
			if strategy != nil {
				t.Error("Expected nil strategy for invalid config, got non-nil")
			}
		})
	}
}

func TestDynamicGasPriceStrategy_SinglePrice(t *testing.T) {
	basePrice := uint64(5000000000) // 5 Gwei
	strategy, err := NewDynamicGasPriceStrategy(DynamicGasPriceConfig{
		GasPrices: []uint64{basePrice},
		Variation: 0.2, // ±20%
	})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Call multiple times, should always use the same base price (with variation)
	for i := 0; i < 10; i++ {
		price := strategy.GetGasPrice()
		if price == nil {
			t.Fatal("GetGasPrice returned nil")
		}

		// Price should be within ±20% of base price
		minPrice := uint64(float64(basePrice) * 0.8)
		maxPrice := uint64(float64(basePrice) * 1.2)
		if *price < minPrice || *price > maxPrice {
			t.Errorf("Price %d is outside expected range [%d, %d]", *price, minPrice, maxPrice)
		}
	}
}

func TestDynamicGasPriceStrategy_Cycling(t *testing.T) {
	gasPrices := []uint64{1000000000, 2000000000, 3000000000} // 1, 2, 3 Gwei
	strategy, err := NewDynamicGasPriceStrategy(DynamicGasPriceConfig{
		GasPrices: gasPrices,
		Variation: 0, // No variation for predictable testing
	})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should cycle through prices in order
	for round := 0; round < 3; round++ {
		for i, expectedBase := range gasPrices {
			price := strategy.GetGasPrice()
			if price == nil {
				t.Fatalf("Round %d, index %d: GetGasPrice returned nil", round, i)
			}
			if *price != expectedBase {
				t.Errorf("Round %d, index %d: expected %d, got %d", round, i, expectedBase, *price)
			}
		}
	}
}

func TestDynamicGasPriceStrategy_ZeroPriceReturnsNil(t *testing.T) {
	gasPrices := []uint64{1000000000, 0, 2000000000} // Middle one is 0
	strategy, err := NewDynamicGasPriceStrategy(DynamicGasPriceConfig{
		GasPrices: gasPrices,
		Variation: 0,
	})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// First call: non-zero price
	price := strategy.GetGasPrice()
	if price == nil || *price != 1000000000 {
		t.Errorf("Expected 1000000000, got %v", price)
	}

	// Second call: zero means use network price (nil)
	price = strategy.GetGasPrice()
	if price != nil {
		t.Errorf("Expected nil for zero gas price, got %v", price)
	}

	// Third call: back to non-zero
	price = strategy.GetGasPrice()
	if price == nil || *price != 2000000000 {
		t.Errorf("Expected 2000000000, got %v", price)
	}
}

func TestDynamicGasPriceStrategy_VariationRange(t *testing.T) {
	basePrice := uint64(10000000000) // 10 Gwei
	variation := 0.3                 // ±30%
	strategy, err := NewDynamicGasPriceStrategy(DynamicGasPriceConfig{
		GasPrices: []uint64{basePrice},
		Variation: variation,
	})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	minExpected := uint64(float64(basePrice) * (1 - variation))
	maxExpected := uint64(float64(basePrice) * (1 + variation))

	// Sample many times to check range
	samples := 1000
	foundMin := false
	foundMax := false

	for i := 0; i < samples; i++ {
		price := strategy.GetGasPrice()
		if price == nil {
			t.Fatal("GetGasPrice returned nil")
		}

		// Check within valid range
		if *price < minExpected || *price > maxExpected {
			t.Errorf("Price %d is outside expected range [%d, %d]", *price, minExpected, maxExpected)
		}

		// Track if we've seen prices near the extremes
		if *price <= minExpected+uint64(float64(basePrice)*0.05) {
			foundMin = true
		}
		if *price >= maxExpected-uint64(float64(basePrice)*0.05) {
			foundMax = true
		}
	}

	// With 1000 samples, we should see prices across the range
	if !foundMin || !foundMax {
		t.Logf("Warning: In %d samples, didn't see full range. foundMin=%v, foundMax=%v", samples, foundMin, foundMax)
	}
}

func TestDynamicGasPriceStrategy_NoVariation(t *testing.T) {
	gasPrices := []uint64{5000000000, 10000000000}
	strategy, err := NewDynamicGasPriceStrategy(DynamicGasPriceConfig{
		GasPrices: gasPrices,
		Variation: 0, // No variation
	})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// With no variation, should return exact prices
	price := strategy.GetGasPrice()
	if price == nil || *price != gasPrices[0] {
		t.Errorf("Expected exact price %d, got %v", gasPrices[0], price)
	}

	price = strategy.GetGasPrice()
	if price == nil || *price != gasPrices[1] {
		t.Errorf("Expected exact price %d, got %v", gasPrices[1], price)
	}
}

func TestDynamicGasPriceStrategy_HighVariation(t *testing.T) {
	basePrice := uint64(1000000000) // 1 Gwei
	variation := 0.9                // ±90%
	strategy, err := NewDynamicGasPriceStrategy(DynamicGasPriceConfig{
		GasPrices: []uint64{basePrice},
		Variation: variation,
	})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	minExpected := uint64(float64(basePrice) * (1 - variation))
	maxExpected := uint64(float64(basePrice) * (1 + variation))

	// Check that high variation still stays within bounds
	for i := 0; i < 100; i++ {
		price := strategy.GetGasPrice()
		if price == nil {
			t.Fatal("GetGasPrice returned nil")
		}

		if *price < minExpected || *price > maxExpected {
			t.Errorf("Price %d is outside expected range [%d, %d]", *price, minExpected, maxExpected)
		}
	}
}

func TestDynamicGasPriceStrategy_ConcurrentAccess(t *testing.T) {
	gasPrices := []uint64{1000000000, 2000000000, 3000000000}
	strategy, err := NewDynamicGasPriceStrategy(DynamicGasPriceConfig{
		GasPrices: gasPrices,
		Variation: 0.1,
	})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Launch multiple goroutines calling GetGasPrice concurrently
	const numGoroutines = 100
	results := make(chan *uint64, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			price := strategy.GetGasPrice()
			results <- price
		}()
	}

	// Collect all results
	prices := make([]*uint64, 0, numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		price := <-results
		if price == nil {
			t.Error("GetGasPrice returned nil in concurrent access")
		}
		prices = append(prices, price)
	}

	// Verify all prices are valid (no panics, no invalid values)
	for i, price := range prices {
		if price == nil {
			continue // Already reported above
		}

		// Should be within reasonable bounds (base prices with variation)
		minValid := uint64(float64(gasPrices[0]) * 0.9)
		maxValid := uint64(float64(gasPrices[len(gasPrices)-1]) * 1.1)
		if *price < minValid || *price > maxValid {
			t.Errorf("Price %d at index %d is outside reasonable bounds [%d, %d]", *price, i, minValid, maxValid)
		}
	}
}
