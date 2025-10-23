package gasmanager

type FixedGasPriceConfig struct {
	GasPriceWei uint64
}

type FixedGasPriceStrategy struct {
	config FixedGasPriceConfig
}

func NewFixedGasPriceStrategy(config FixedGasPriceConfig) *FixedGasPriceStrategy {
	return &FixedGasPriceStrategy{
		config: config,
	}
}

func (s *FixedGasPriceStrategy) GetGasPrice() *uint64 {
	return &s.config.GasPriceWei
}
