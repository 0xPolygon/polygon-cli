package gasmanager

import (
	"math/rand/v2"
	"sync/atomic"
)

type DynamicGasPriceConfig struct {
	GasPrices []uint64
}

type DynamicGasPriceStrategy struct {
	config DynamicGasPriceConfig
	i      atomic.Int32
}

func NewDynamicGasPriceStrategy(config DynamicGasPriceConfig) *DynamicGasPriceStrategy {
	s := &DynamicGasPriceStrategy{
		config: config,
	}
	return s
}

func (s *DynamicGasPriceStrategy) GetGasPrice() *uint64 {
	idx := s.i.Load()
	s.i.Store((idx + 1) % int32(len(s.config.GasPrices)))
	gp := s.config.GasPrices[idx]
	if gp == 0 {
		return nil
	}

	// introduce random  variation
	factor := 0.7 + rand.Float64()*0.6 // random value in [0.7, 1.3]
	varied := uint64(float64(gp) * factor)
	return &varied
}
