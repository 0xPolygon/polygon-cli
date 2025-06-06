package indexer

import (
	"context"
	"sync"
	"time"

	"github.com/0xPolygon/polygon-cli/blockstore"
)

// Indexer is responsible for fetching blockchain data and populating the store
type Indexer struct {
	// Store is the blockstore to populate with data
	store blockstore.BlockStore

	// Configuration
	pollingInterval time.Duration // How often to poll for new blocks
	lookbackDepth   int64         // How many blocks to keep in the store
	reorgDepth      int64         // How many blocks back to check for reorgs
	maxConcurrency  int           // Maximum concurrent requests to the store

	// State tracking
	latestHeight    int64      // Latest block height we've indexed
	mu              sync.RWMutex // Protects state fields
	
	// Control channels
	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct{}

	// Worker pool
	workerSem chan struct{} // Semaphore for controlling concurrency
}

// Config holds the configuration for the indexer
type Config struct {
	// PollingInterval is how often to check for new blocks
	PollingInterval time.Duration
	
	// LookbackDepth is how many blocks to keep in the store
	// 0 means keep all blocks
	LookbackDepth int64
	
	// ReorgDepth is how many blocks back to check for reorgs
	// This should be based on the chain's finality assumptions
	ReorgDepth int64
	
	// MaxConcurrency limits the number of concurrent requests
	MaxConcurrency int
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		PollingInterval: 2 * time.Second,
		LookbackDepth:   1000,
		ReorgDepth:      128,
		MaxConcurrency:  10,
	}
}

// NewIndexer creates a new indexer with the given store and configuration
func NewIndexer(store blockstore.BlockStore, cfg *Config) *Indexer {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	ctx, cancel := context.WithCancel(context.Background())
	
	return &Indexer{
		store:           store,
		pollingInterval: cfg.PollingInterval,
		lookbackDepth:   cfg.LookbackDepth,
		reorgDepth:      cfg.ReorgDepth,
		maxConcurrency:  cfg.MaxConcurrency,
		latestHeight:    -1,
		ctx:             ctx,
		cancel:          cancel,
		done:            make(chan struct{}),
		workerSem:       make(chan struct{}, cfg.MaxConcurrency),
	}
}

// LatestHeight returns the latest block height that has been indexed
func (i *Indexer) LatestHeight() int64 {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.latestHeight
}

// Start begins the indexing process
func (i *Indexer) Start() error {
	// TODO: Implement indexing logic
	return nil
}

// Stop gracefully stops the indexer
func (i *Indexer) Stop() error {
	i.cancel()
	<-i.done
	return nil
}