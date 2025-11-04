package gasmanager

// PriceStrategy defines the interface for different gas price strategies.
type PriceStrategy interface {
	GetGasPrice() *uint64
}

// GasPricer uses a PriceStrategy to determine the gas price.
type GasPricer struct {
	strategy PriceStrategy
}

// NewGasPricer creates a new GasPricer with the given PriceStrategy.
func NewGasPricer(strategy PriceStrategy) *GasPricer {
	return &GasPricer{
		strategy: strategy,
	}
}

// GetGasPrice retrieves the gas price using the configured PriceStrategy.
func (gp *GasPricer) GetGasPrice() *uint64 {
	return gp.strategy.GetGasPrice()
}
