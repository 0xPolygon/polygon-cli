package gasmanager

// FixedGasPriceConfig holds the configuration for the FixedGasPriceStrategy.
type FixedGasPriceConfig struct {
	GasPriceWei uint64
}

// FixedGasPriceStrategy provides a fixed gas price.
type FixedGasPriceStrategy struct {
	config FixedGasPriceConfig
}

// NewFixedGasPriceStrategy creates a new FixedGasPriceStrategy with the given configuration.
func NewFixedGasPriceStrategy(config FixedGasPriceConfig) *FixedGasPriceStrategy {
	return &FixedGasPriceStrategy{
		config: config,
	}
}

// GetGasPrice retrieves the fixed gas price.
func (s *FixedGasPriceStrategy) GetGasPrice() *uint64 {
	return &s.config.GasPriceWei
}
