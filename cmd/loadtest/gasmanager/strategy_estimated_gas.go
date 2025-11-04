package gasmanager

// EstimatedGasPriceStrategy provides gas prices estimated from the network.
type EstimatedGasPriceStrategy struct {
}

// NewEstimatedGasPriceStrategy creates a new EstimatedGasPriceStrategy.
func NewEstimatedGasPriceStrategy() *EstimatedGasPriceStrategy {
	return &EstimatedGasPriceStrategy{}
}

// GetGasPrice retrieves the estimated gas price from the network.
// For this strategy, we return nil to indicate that the default network gas price should be used.
func (s *EstimatedGasPriceStrategy) GetGasPrice() *uint64 {
	return nil
}
