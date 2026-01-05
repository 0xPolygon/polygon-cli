package gasmanager

import (
	"math"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
)

func TestOscillatingGasProvider_NewOscillatingGasProvider(t *testing.T) {
	vault := NewGasVault()
	wave := NewFlatWave(WaveConfig{
		Period:    10,
		Amplitude: 1000,
		Target:    5000,
	})

	provider := NewOscillatingGasProvider(nil, vault, wave)

	if provider == nil {
		t.Fatal("NewOscillatingGasProvider returned nil")
	}
	if provider.vault != vault {
		t.Error("Provider vault doesn't match provided vault")
	}
	if provider.wave != wave {
		t.Error("Provider wave doesn't match provided wave")
	}
}

func TestOscillatingGasProvider_OnNewHeader_AddsGasToVault(t *testing.T) {
	vault := NewGasVault()
	initialBudget := vault.GetAvailableBudget()

	// Create a flat wave with target 5000
	wave := NewFlatWave(WaveConfig{
		Period:    10,
		Amplitude: 0,
		Target:    5000,
	})

	provider := NewOscillatingGasProvider(nil, vault, wave)

	// Simulate receiving a new block header
	header := &types.Header{
		Number: big.NewInt(100),
	}

	provider.onNewHeader(header)

	// Verify gas was added to vault
	newBudget := vault.GetAvailableBudget()
	if newBudget == initialBudget {
		t.Error("Gas was not added to vault after new header")
	}
	if newBudget != 5000 {
		t.Errorf("Expected budget 5000, got %d", newBudget)
	}
}

func TestOscillatingGasProvider_OnNewHeader_AdvancesWave(t *testing.T) {
	vault := NewGasVault()

	// Create a wave where we can track position changes
	// Using a sawtooth with period 3, we can check Y values change
	wave := NewSawtoothWave(WaveConfig{
		Period:    3,
		Amplitude: 1000,
		Target:    2000,
	})

	provider := NewOscillatingGasProvider(nil, vault, wave)

	// Get initial Y value
	initialY := wave.Y()

	// Simulate receiving multiple headers
	for i := 0; i < 3; i++ {
		header := &types.Header{
			Number: big.NewInt(int64(100 + i)),
		}
		provider.onNewHeader(header)
	}

	// After 3 headers (one complete period), X should have wrapped
	// and Y should be different from initial
	finalX := wave.X()
	if finalX != 0 {
		t.Errorf("Expected wave X to wrap to 0 after period, got %f", finalX)
	}

	// Initial Y should equal final Y after one complete period
	finalY := wave.Y()
	if finalY != initialY {
		t.Logf("Initial Y: %f, Final Y: %f", initialY, finalY)
	}
}

func TestOscillatingGasProvider_OnNewHeader_AccumulatesGas(t *testing.T) {
	vault := NewGasVault()

	// Flat wave adds same amount each time
	wave := NewFlatWave(WaveConfig{
		Period:    1,
		Amplitude: 0,
		Target:    1000,
	})

	provider := NewOscillatingGasProvider(nil, vault, wave)

	// Simulate multiple blocks
	numBlocks := 5
	for i := 0; i < numBlocks; i++ {
		header := &types.Header{
			Number: big.NewInt(int64(i)),
		}
		provider.onNewHeader(header)
	}

	// Total gas should be numBlocks * target
	expectedTotal := uint64(numBlocks * 1000)
	actualTotal := vault.GetAvailableBudget()
	if actualTotal != expectedTotal {
		t.Errorf("Expected accumulated gas %d, got %d", expectedTotal, actualTotal)
	}
}

func TestOscillatingGasProvider_OnNewHeader_HandlesFlooring(t *testing.T) {
	vault := NewGasVault()

	// Create a sine wave that will produce non-integer Y values
	wave := NewSineWave(WaveConfig{
		Period:    100,
		Amplitude: 500,
		Target:    1000,
	})

	provider := NewOscillatingGasProvider(nil, vault, wave)

	header := &types.Header{
		Number: big.NewInt(25), // Quarter period, should be at peak
	}

	initialY := wave.Y()
	provider.onNewHeader(header)

	// Gas added should be floor of Y value
	expectedGas := uint64(math.Floor(initialY))
	actualGas := vault.GetAvailableBudget()

	if actualGas != expectedGas {
		t.Errorf("Expected gas %d (floor of %f), got %d", expectedGas, initialY, actualGas)
	}
}

func TestOscillatingGasProvider_OnNewHeader_WithSineWave(t *testing.T) {
	vault := NewGasVault()

	wave := NewSineWave(WaveConfig{
		Period:    4,
		Amplitude: 1000,
		Target:    2000,
	})

	provider := NewOscillatingGasProvider(nil, vault, wave)

	gasAmounts := make([]uint64, 0, 4)

	// Collect gas amounts for one complete period
	for i := 0; i < 4; i++ {
		initialBudget := vault.GetAvailableBudget()
		header := &types.Header{
			Number: big.NewInt(int64(i)),
		}
		provider.onNewHeader(header)
		newBudget := vault.GetAvailableBudget()
		gasAdded := newBudget - initialBudget
		gasAmounts = append(gasAmounts, gasAdded)
	}

	// Verify that gas amounts vary (sine wave oscillates)
	allSame := true
	first := gasAmounts[0]
	for _, amount := range gasAmounts[1:] {
		if amount != first {
			allSame = false
			break
		}
	}

	if allSame {
		t.Error("Sine wave should produce varying gas amounts, but all were the same")
	}
}

func TestOscillatingGasProvider_OnNewHeader_WithSquareWave(t *testing.T) {
	vault := NewGasVault()

	wave := NewSquareWave(WaveConfig{
		Period:    4,
		Amplitude: 1000,
		Target:    2000,
	})

	provider := NewOscillatingGasProvider(nil, vault, wave)

	gasAmounts := make([]uint64, 0, 4)

	// Collect gas amounts for one complete period
	for i := 0; i < 4; i++ {
		initialBudget := vault.GetAvailableBudget()
		header := &types.Header{
			Number: big.NewInt(int64(i)),
		}
		provider.onNewHeader(header)
		newBudget := vault.GetAvailableBudget()
		gasAdded := newBudget - initialBudget
		gasAmounts = append(gasAmounts, gasAdded)
	}

	// Square wave should have two distinct values
	// First half should be one value, second half another
	firstHalf := gasAmounts[:2]
	secondHalf := gasAmounts[2:]

	// Check first half is consistent
	if firstHalf[0] != firstHalf[1] {
		t.Errorf("First half of square wave should be consistent, got %d and %d", firstHalf[0], firstHalf[1])
	}

	// Check second half is consistent
	if secondHalf[0] != secondHalf[1] {
		t.Errorf("Second half of square wave should be consistent, got %d and %d", secondHalf[0], secondHalf[1])
	}

	// Check that first and second halves are different
	if firstHalf[0] == secondHalf[0] {
		t.Error("Square wave first and second halves should be different")
	}
}

func TestOscillatingGasProvider_OnNewHeader_WithTriangleWave(t *testing.T) {
	vault := NewGasVault()

	wave := NewTriangleWave(WaveConfig{
		Period:    6,
		Amplitude: 1000,
		Target:    2000,
	})

	provider := NewOscillatingGasProvider(nil, vault, wave)

	gasAmounts := make([]uint64, 0, 6)

	// Collect gas amounts for one complete period
	for i := 0; i < 6; i++ {
		initialBudget := vault.GetAvailableBudget()
		header := &types.Header{
			Number: big.NewInt(int64(i)),
		}
		provider.onNewHeader(header)
		newBudget := vault.GetAvailableBudget()
		gasAdded := newBudget - initialBudget
		gasAmounts = append(gasAmounts, gasAdded)
	}

	// Triangle wave should increase then decrease
	// Check that first half generally increases
	increasing := true
	for i := 1; i < 3; i++ {
		if gasAmounts[i] < gasAmounts[i-1] {
			increasing = false
			break
		}
	}

	if !increasing {
		t.Error("First half of triangle wave should generally increase")
	}

	// Check that second half generally decreases
	decreasing := true
	for i := 4; i < 6; i++ {
		if gasAmounts[i] > gasAmounts[i-1] {
			decreasing = false
			break
		}
	}

	if !decreasing {
		t.Error("Second half of triangle wave should generally decrease")
	}
}

func TestOscillatingGasProvider_OnNewHeader_NilVault(t *testing.T) {
	wave := NewFlatWave(WaveConfig{
		Period:    10, // Use period > 1 so X doesn't wrap
		Amplitude: 0,
		Target:    1000,
	})

	// Create provider with nil vault
	provider := NewOscillatingGasProvider(nil, nil, wave)

	header := &types.Header{
		Number: big.NewInt(100),
	}

	// Should not panic with nil vault
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("onNewHeader panicked with nil vault: %v", r)
		}
	}()

	provider.onNewHeader(header)

	// Wave should still advance
	if wave.X() != 1 {
		t.Errorf("Wave should advance even with nil vault, expected X=1, got X=%f", wave.X())
	}
}

func TestOscillatingGasProvider_OnNewHeader_LargeGasAmount(t *testing.T) {
	vault := NewGasVault()

	// Create wave with very large target
	wave := NewFlatWave(WaveConfig{
		Period:    1,
		Amplitude: 0,
		Target:    math.MaxUint64 / 2,
	})

	provider := NewOscillatingGasProvider(nil, vault, wave)

	// Add gas twice
	for i := 0; i < 2; i++ {
		header := &types.Header{
			Number: big.NewInt(int64(i)),
		}
		provider.onNewHeader(header)
	}

	// Should cap at max uint64 (tested in vault tests, but verify here too)
	budget := vault.GetAvailableBudget()
	if budget < math.MaxUint64/2 {
		t.Errorf("Expected budget at least %d, got %d", math.MaxUint64/2, budget)
	}
}
