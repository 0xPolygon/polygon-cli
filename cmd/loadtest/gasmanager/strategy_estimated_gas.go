package gasmanager

type EstimatedGasPriceStrategy struct {
}

func NewEstimatedGasPriceStrategy() *EstimatedGasPriceStrategy {
	return &EstimatedGasPriceStrategy{}
}

func (s *EstimatedGasPriceStrategy) GetGasPrice() *uint64 {
	return nil
}
