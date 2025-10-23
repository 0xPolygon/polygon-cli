package gasmanager

type PriceStrategy interface {
	GetGasPrice() *uint64
}

type GasPricer struct {
	strategy PriceStrategy
}

func NewGasPricer(strategy PriceStrategy) *GasPricer {
	return &GasPricer{
		strategy: strategy,
	}
}

func (gp *GasPricer) GetGasPrice() *uint64 {
	return gp.strategy.GetGasPrice()
}
