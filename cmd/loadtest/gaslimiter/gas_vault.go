package gaslimiter

import (
	"math"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type GasVault struct {
	mu                 *sync.Mutex
	gasBudgetAvailable uint64
}

func NewGasVault() *GasVault {
	return &GasVault{
		mu:                 &sync.Mutex{},
		gasBudgetAvailable: 0,
	}
}

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

func (o *GasVault) SpendOrWaitAvailableBudget(gas uint64) {
	o.mu.Lock()
	log.Trace().Uint64("gas", gas).Uint64("available_budget", o.gasBudgetAvailable).Msg("requesting gas from vault")
	o.mu.Unlock()
	for {
		o.mu.Lock()
		if gas <= o.gasBudgetAvailable {
			o.gasBudgetAvailable -= gas
			log.Trace().Uint64("gas", gas).Uint64("available_budget", o.gasBudgetAvailable).Msg("gas spent from vault")
			o.mu.Unlock()
			break
		}
		o.mu.Unlock()
		time.Sleep(100 * time.Millisecond)
	}
}

func (o *GasVault) GetAvailableBudget() uint64 {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.gasBudgetAvailable
}
