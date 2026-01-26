package gasmanager

import (
	"context"
	"math"
	"sync"

	"github.com/rs/zerolog/log"
)

// GasVault manages a budget of gas that can be added to and spent from.
type GasVault struct {
	mu                 sync.Mutex
	cond               *sync.Cond
	gasBudgetAvailable uint64
}

// NewGasVault creates a new GasVault instance.
func NewGasVault() *GasVault {
	v := &GasVault{}
	v.cond = sync.NewCond(&v.mu)
	return v
}

// AddGas adds the specified amount of gas to the vault's available budget.
func (v *GasVault) AddGas(gas uint64) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if v.gasBudgetAvailable+gas < v.gasBudgetAvailable {
		v.gasBudgetAvailable = math.MaxUint64
		log.Warn().Uint64("available_budget", v.gasBudgetAvailable).Msg("Gas budget overflow, capped to max uint64")
	} else {
		v.gasBudgetAvailable += gas
		log.Trace().Uint64("available_budget", v.gasBudgetAvailable).Msg("Gas added to vault")
	}

	// Signal waiters that gas was added
	v.cond.Broadcast()
}

// SpendOrWaitAvailableBudget attempts to spend the specified amount of gas from the vault's available budget.
// It blocks until sufficient budget is available or the context is cancelled.
func (v *GasVault) SpendOrWaitAvailableBudget(ctx context.Context, gas uint64) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	// Start a goroutine to handle context cancellation
	done := make(chan struct{})
	defer close(done)

	go func() {
		select {
		case <-ctx.Done():
			v.cond.Broadcast() // Wake up the waiter so it can check context
		case <-done:
		}
	}()

	for gas > v.gasBudgetAvailable {
		// Check if context was cancelled before waiting
		if ctx.Err() != nil {
			return ctx.Err()
		}
		v.cond.Wait()
	}

	// Final context check after waking up
	if ctx.Err() != nil {
		return ctx.Err()
	}

	v.gasBudgetAvailable -= gas
	return nil
}

// GetAvailableBudget returns the current available gas budget in the vault.
func (v *GasVault) GetAvailableBudget() uint64 {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.gasBudgetAvailable
}
