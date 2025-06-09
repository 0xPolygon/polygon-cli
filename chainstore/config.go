package chainstore

import "time"

// ChainStoreConfig holds configuration for the ChainStore
type ChainStoreConfig struct {
	// Cache TTLs for different data types
	StaticTTL       time.Duration // Never expire (0)
	SemiStaticTTL   time.Duration // 5 minutes
	FrequentTTL     time.Duration // 30 seconds
	VeryFrequentTTL time.Duration // 5 seconds

	// Capability detection
	CapabilityTTL time.Duration // 1 hour

	// Feature toggles
	EnableTxPoolMonitoring   bool
	EnableFinalityTracking   bool
	EnableFeeHistoryTracking bool
}

// DefaultChainStoreConfig returns default configuration
func DefaultChainStoreConfig() *ChainStoreConfig {
	return &ChainStoreConfig{
		StaticTTL:       0,                // Never expire
		SemiStaticTTL:   5 * time.Minute,  // 5 minutes
		FrequentTTL:     30 * time.Second, // 30 seconds
		VeryFrequentTTL: 5 * time.Second,  // 5 seconds
		CapabilityTTL:   1 * time.Hour,    // 1 hour

		EnableTxPoolMonitoring:   true,
		EnableFinalityTracking:   true,
		EnableFeeHistoryTracking: true,
	}
}
