package gaslimiter

import (
	"math"
	"sync"
	"time"
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
	} else {
		o.gasBudgetAvailable += gas
	}
}

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
