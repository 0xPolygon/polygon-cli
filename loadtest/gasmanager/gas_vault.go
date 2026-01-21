package gasmanager

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// GasVault manages a budget of gas that can be added to and spent from.
type GasVault struct {
	mu                 sync.Mutex
	gasBudgetAvailable uint64
}

// NewGasVault creates a new GasVault instance.
func NewGasVault() *GasVault {
	return &GasVault{}
}

// AddGas adds the specified amount of gas to the vault's available budget.
func (o *GasVault) AddGas(gas uint64) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.gasBudgetAvailable+gas < o.gasBudgetAvailable {
		o.gasBudgetAvailable = math.MaxUint64
		log.Trace().Uint64("available_budget", o.gasBudgetAvailable).Msg("Gas budget capped to max uint64")
	} else {
		o.gasBudgetAvailable += gas
		log.Trace().Uint64("available_budget", o.gasBudgetAvailable).Msg("Gas added to vault")
	}
}

// SpendOrWaitAvailableBudget attempts to spend the specified amount of gas from the vault's available budget.
// It blocks until sufficient budget is available or the context is cancelled.
func (o *GasVault) SpendOrWaitAvailableBudget(ctx context.Context, gas uint64) error {
	const intervalToCheckBudgetAvailability = 100 * time.Millisecond
	ticker := time.NewTicker(intervalToCheckBudgetAvailability)
	defer ticker.Stop()

	for {
		if spent := o.trySpendBudget(gas); spent {
			return nil
		}
		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// trySpendBudget tries to spend the specified amount of gas from the vault's available budget.
// It returns true if the gas was successfully spent, or false if there was insufficient budget.
func (o *GasVault) trySpendBudget(gas uint64) bool {
	o.mu.Lock()
	defer o.mu.Unlock()
	if gas <= o.gasBudgetAvailable {
		o.gasBudgetAvailable -= gas
		return true
	}
	return false
}

// GetAvailableBudget returns the current available gas budget in the vault.
func (o *GasVault) GetAvailableBudget() uint64 {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.gasBudgetAvailable
}
