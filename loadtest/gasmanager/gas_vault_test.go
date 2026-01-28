package gasmanager

import (
	"context"
	"math"
	"sync"
	"testing"
	"time"
)

func TestNewGasVault(t *testing.T) {
	vault := NewGasVault()
	if vault == nil {
		t.Fatal("NewGasVault returned nil")
	}
	if vault.GetAvailableBudget() != 0 {
		t.Errorf("Expected initial budget to be 0, got %d", vault.GetAvailableBudget())
	}
}

func TestGasVault_AddGas(t *testing.T) {
	tests := []struct {
		name          string
		initialBudget uint64
		addAmounts    []uint64
		expectedFinal uint64
	}{
		{
			name:          "add single amount",
			initialBudget: 0,
			addAmounts:    []uint64{1000},
			expectedFinal: 1000,
		},
		{
			name:          "add multiple amounts",
			initialBudget: 500,
			addAmounts:    []uint64{200, 300, 100},
			expectedFinal: 1100,
		},
		{
			name:          "add zero",
			initialBudget: 1000,
			addAmounts:    []uint64{0},
			expectedFinal: 1000,
		},
		{
			name:          "overflow protection - cap at max uint64",
			initialBudget: math.MaxUint64 - 100,
			addAmounts:    []uint64{200},
			expectedFinal: math.MaxUint64,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vault := NewGasVault()
			if tt.initialBudget > 0 {
				vault.AddGas(tt.initialBudget)
			}

			for _, amount := range tt.addAmounts {
				vault.AddGas(amount)
			}

			got := vault.GetAvailableBudget()
			if got != tt.expectedFinal {
				t.Errorf("Expected final budget %d, got %d", tt.expectedFinal, got)
			}
		})
	}
}

func TestGasVault_SpendOrWaitAvailableBudget_Success(t *testing.T) {
	vault := NewGasVault()
	vault.AddGas(1000)

	ctx := context.Background()
	err := vault.SpendOrWaitAvailableBudget(ctx, 500)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	remaining := vault.GetAvailableBudget()
	if remaining != 500 {
		t.Errorf("Expected remaining budget 500, got %d", remaining)
	}
}

func TestGasVault_SpendOrWaitAvailableBudget_ExactAmount(t *testing.T) {
	vault := NewGasVault()
	vault.AddGas(1000)

	ctx := context.Background()
	err := vault.SpendOrWaitAvailableBudget(ctx, 1000)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	remaining := vault.GetAvailableBudget()
	if remaining != 0 {
		t.Errorf("Expected remaining budget 0, got %d", remaining)
	}
}

func TestGasVault_SpendOrWaitAvailableBudget_ContextCancelled(t *testing.T) {
	vault := NewGasVault()
	// Don't add any budget, so it will block

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := vault.SpendOrWaitAvailableBudget(ctx, 500)
	if err == nil {
		t.Error("Expected error due to cancelled context, got nil")
	}
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled error, got %v", err)
	}

	// Budget should remain unchanged
	remaining := vault.GetAvailableBudget()
	if remaining != 0 {
		t.Errorf("Expected budget to remain 0, got %d", remaining)
	}
}

func TestGasVault_SpendOrWaitAvailableBudget_ContextTimeout(t *testing.T) {
	vault := NewGasVault()
	// Don't add any budget

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	start := time.Now()
	err := vault.SpendOrWaitAvailableBudget(ctx, 500)
	elapsed := time.Since(start)

	if err == nil {
		t.Error("Expected error due to timeout, got nil")
	}
	if err != context.DeadlineExceeded {
		t.Errorf("Expected context.DeadlineExceeded error, got %v", err)
	}
	if elapsed < 200*time.Millisecond {
		t.Errorf("Expected to wait at least 200ms, waited %v", elapsed)
	}
}

func TestGasVault_SpendOrWaitAvailableBudget_WaitThenSuccess(t *testing.T) {
	vault := NewGasVault()
	vault.AddGas(100) // Not enough initially

	ctx := context.Background()

	// Start goroutine that will add budget after a delay
	go func() {
		time.Sleep(200 * time.Millisecond)
		vault.AddGas(1000) // Now enough budget available
	}()

	start := time.Now()
	err := vault.SpendOrWaitAvailableBudget(ctx, 500)
	elapsed := time.Since(start)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if elapsed < 200*time.Millisecond {
		t.Errorf("Expected to wait at least 200ms for budget, waited %v", elapsed)
	}

	remaining := vault.GetAvailableBudget()
	expected := uint64(600) // 100 + 1000 - 500
	if remaining != expected {
		t.Errorf("Expected remaining budget %d, got %d", expected, remaining)
	}
}

func TestGasVault_ConcurrentAccess(t *testing.T) {
	vault := NewGasVault()
	vault.AddGas(10000)

	ctx := context.Background()
	const numGoroutines = 100
	const spendAmount = 50

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines)

	// Spawn multiple goroutines trying to spend gas concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := vault.SpendOrWaitAvailableBudget(ctx, spendAmount)
			if err != nil {
				errors <- err
			}
		}()
	}

	wg.Wait()
	close(errors)

	// Check for errors
	for err := range errors {
		t.Errorf("Unexpected error in concurrent access: %v", err)
	}

	// Verify final budget
	expected := uint64(10000 - (numGoroutines * spendAmount))
	remaining := vault.GetAvailableBudget()
	if remaining != expected {
		t.Errorf("Expected remaining budget %d after concurrent access, got %d", expected, remaining)
	}
}

func TestGasVault_ConcurrentAddAndSpend(t *testing.T) {
	vault := NewGasVault()
	vault.AddGas(5000)

	ctx := context.Background()
	const numAdders = 50
	const numSpenders = 50
	const addAmount = 100
	const spendAmount = 100

	var wg sync.WaitGroup

	// Adders
	for i := 0; i < numAdders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			vault.AddGas(addAmount)
		}()
	}

	// Spenders
	errors := make(chan error, numSpenders)
	for i := 0; i < numSpenders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := vault.SpendOrWaitAvailableBudget(ctx, spendAmount)
			if err != nil {
				errors <- err
			}
		}()
	}

	wg.Wait()
	close(errors)

	// Check for errors
	for err := range errors {
		t.Errorf("Unexpected error in concurrent add/spend: %v", err)
	}

	// Verify final budget: initial + added - spent
	expected := uint64(5000 + (numAdders * addAmount) - (numSpenders * spendAmount))
	remaining := vault.GetAvailableBudget()
	if remaining != expected {
		t.Errorf("Expected remaining budget %d, got %d", expected, remaining)
	}
}

func TestGasVault_MultipleSpendersWaiting(t *testing.T) {
	vault := NewGasVault()
	// Start with no budget

	ctx := context.Background()
	const numSpenders = 5
	const spendAmount = 100

	var wg sync.WaitGroup
	successCount := make(chan int, numSpenders)

	// Start multiple spenders that will all wait
	for i := 0; i < numSpenders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := vault.SpendOrWaitAvailableBudget(ctx, spendAmount)
			if err == nil {
				successCount <- 1
			}
		}()
	}

	// Give goroutines time to start waiting
	time.Sleep(100 * time.Millisecond)

	// Add enough budget for all spenders
	vault.AddGas(uint64(numSpenders * spendAmount))

	wg.Wait()
	close(successCount)

	// Count successful spends
	count := 0
	for range successCount {
		count++
	}

	if count != numSpenders {
		t.Errorf("Expected %d successful spends, got %d", numSpenders, count)
	}

	// Budget should be depleted
	remaining := vault.GetAvailableBudget()
	if remaining != 0 {
		t.Errorf("Expected remaining budget 0, got %d", remaining)
	}
}
