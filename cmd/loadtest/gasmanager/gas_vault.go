package gasmanager

import (
	"math"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// GasVault manages a budget of gas that can be added to and spent from.
type GasVault struct {
	mu                 *sync.Mutex
	gasBudgetAvailable uint64
}

// NewGasVault creates a new GasVault instance.
func NewGasVault() *GasVault {
	return &GasVault{
		mu:                 &sync.Mutex{},
		gasBudgetAvailable: 0,
	}
}

// AddGas adds the specified amount of gas to the vault's available budget.
func (o *GasVault) AddGas(gas uint64) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.gasBudgetAvailable+gas < o.gasBudgetAvailable {
		o.gasBudgetAvailable = math.MaxUint64
		log.Trace().Uint64("available_budget_after", o.gasBudgetAvailable).Msg("gas budget in vault capped to max uint64")
	} else {
		o.gasBudgetAvailable += gas
		log.Trace().Uint64("available_budget_after", o.gasBudgetAvailable).Msg("new gas budget available in vault")
	}
}

// SpendOrWaitAvailableBudget attempts to spend the specified amount of gas from the vault's available budget.
func (o *GasVault) SpendOrWaitAvailableBudget(gas uint64) {
	for {
		o.mu.Lock()
		if gas <= o.gasBudgetAvailable {
			o.gasBudgetAvailable -= gas
			o.mu.Unlock()
			break
		}
		o.mu.Unlock()
		time.Sleep(100 * time.Millisecond)
	}
}

// GetAvailableBudget returns the current available gas budget in the vault.
func (o *GasVault) GetAvailableBudget() uint64 {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.gasBudgetAvailable
}
