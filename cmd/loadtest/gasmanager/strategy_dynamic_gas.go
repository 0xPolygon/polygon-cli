package gasmanager

import (
	"math/rand/v2"
	"sync/atomic"
)

// DynamicGasPriceConfig holds the configuration for the DynamicGasPriceStrategy.
type DynamicGasPriceConfig struct {
	GasPrices []uint64
	Variation float64
}

// DynamicGasPriceStrategy provides gas prices from a predefined list with random variation.
type DynamicGasPriceStrategy struct {
	config DynamicGasPriceConfig
	i      atomic.Int32
}

// NewDynamicGasPriceStrategy creates a new DynamicGasPriceStrategy with the given configuration.
func NewDynamicGasPriceStrategy(config DynamicGasPriceConfig) *DynamicGasPriceStrategy {
	s := &DynamicGasPriceStrategy{
		config: config,
	}
	return s
}

// GetGasPrice retrieves the next gas price from the list, applying random variation.
func (s *DynamicGasPriceStrategy) GetGasPrice() *uint64 {
	idx := s.i.Load()
	s.i.Store((idx + 1) % int32(len(s.config.GasPrices)))
	gp := s.config.GasPrices[idx]
	if gp == 0 {
		return nil
	}

	variationMin := float64(1) - s.config.Variation
	variationMax := float64(1) + s.config.Variation
	factor := variationMin + rand.Float64()*(variationMax-variationMin)
	varied := uint64(float64(gp) * factor)
	return &varied
}
