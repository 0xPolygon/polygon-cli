package indexer

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/0xPolygon/polygon-cli/blockstore"
	"github.com/0xPolygon/polygon-cli/rpctypes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
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
	latestHeight int64        // Latest block height we've indexed
	mu           sync.RWMutex // Protects state fields

	// Control channels
	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct{}

	// Block channel for publishing new blocks to renderers
	blockChan chan rpctypes.PolyBlock

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
		LookbackDepth:   128,
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
		blockChan:       make(chan rpctypes.PolyBlock, 100), // Buffered channel
		workerSem:       make(chan struct{}, cfg.MaxConcurrency),
	}
}

// BlockChannel returns the channel where new blocks are published
func (i *Indexer) BlockChannel() <-chan rpctypes.PolyBlock {
	return i.blockChan
}

// GetBlock retrieves a block by hash or number through the store
func (i *Indexer) GetBlock(ctx context.Context, blockHashOrNumber interface{}) (rpctypes.PolyBlock, error) {
	return i.store.GetBlock(ctx, blockHashOrNumber)
}

// GetTransaction retrieves a transaction by hash through the store
func (i *Indexer) GetTransaction(ctx context.Context, txHash common.Hash) (rpctypes.PolyTransaction, error) {
	return i.store.GetTransaction(ctx, txHash)
}

// GetReceipt retrieves a transaction receipt by hash through the store
func (i *Indexer) GetReceipt(ctx context.Context, txHash common.Hash) (rpctypes.PolyReceipt, error) {
	return i.store.GetReceipt(ctx, txHash)
}

// LatestHeight returns the latest block height that has been indexed
func (i *Indexer) LatestHeight() int64 {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.latestHeight
}

// Start begins the indexing process
func (i *Indexer) Start() error {
	log.Info().Msg("Starting indexer")

	go i.indexingLoop()
	return nil
}

// Stop gracefully stops the indexer
func (i *Indexer) Stop() error {
	log.Info().Msg("Stopping indexer")
	i.cancel()
	close(i.blockChan)
	<-i.done
	return nil
}

// indexingLoop is the main loop that polls for new blocks and publishes them
func (i *Indexer) indexingLoop() {
	defer close(i.done)

	// First, do initial catchup to get recent blocks for context
	if err := i.initialCatchup(); err != nil {
		log.Error().Err(err).Msg("Error during initial catchup")
		return
	}

	ticker := time.NewTicker(i.pollingInterval)
	defer ticker.Stop()

	log.Info().Dur("interval", i.pollingInterval).Msg("Starting indexing loop")

	for {
		select {
		case <-i.ctx.Done():
			log.Info().Msg("Indexing loop stopped")
			return
		case <-ticker.C:
			if err := i.checkForNewBlocks(); err != nil {
				log.Error().Err(err).Msg("Error checking for new blocks")
			}
		}
	}
}

// checkForNewBlocks fetches the latest block and publishes it if it's new
func (i *Indexer) checkForNewBlocks() error {
	latestBlock, err := i.store.GetLatestBlock(i.ctx)
	if err != nil {
		return err
	}

	currentTip := latestBlock.Number().Int64()

	i.mu.Lock()
	lastProcessed := i.latestHeight
	i.mu.Unlock()

	// If we've missed blocks, fetch them all to maintain order
	if currentTip > lastProcessed {
		log.Debug().
			Int64("currentTip", currentTip).
			Int64("lastProcessed", lastProcessed).
			Int64("gap", currentTip-lastProcessed).
			Msg("Catching up missed blocks")

		// Fetch missed blocks in parallel and publish them in order
		startHeight := lastProcessed + 1
		blocks, err := i.fetchBlocksInParallel(startHeight, currentTip)
		if err != nil {
			return err
		}

		// Publish blocks to channel in order
		for height, block := range blocks {
			if block == nil {
				// Skip blocks that failed to fetch
				continue
			}

			select {
			case i.blockChan <- block:
				blockHeight := startHeight + int64(height)
				log.Debug().Int64("height", blockHeight).Str("hash", block.Hash().Hex()).Msg("Published block")
				// Update latest height after successful publish
				i.mu.Lock()
				i.latestHeight = blockHeight
				i.mu.Unlock()
			case <-i.ctx.Done():
				return i.ctx.Err()
			}
		}
	}

	return nil
}

// initialCatchup fetches recent blocks to provide context when starting
func (i *Indexer) initialCatchup() error {
	// Get the current tip of the chain
	latestBlock, err := i.store.GetLatestBlock(i.ctx)
	if err != nil {
		return err
	}
	
	currentTip := latestBlock.Number().Int64()
	
	// Calculate starting height (tip - lookbackDepth)
	startHeight := currentTip - i.lookbackDepth
	if startHeight < 0 {
		startHeight = 0
	}
	
	log.Info().
		Int64("currentTip", currentTip).
		Int64("startHeight", startHeight).
		Int64("lookbackDepth", i.lookbackDepth).
		Msg("Starting initial catchup")
	
	// Fetch blocks in parallel and publish them in order
	blocks, err := i.fetchBlocksInParallel(startHeight, currentTip)
	if err != nil {
		return err
	}
	
	// Publish blocks to channel in order
	for height, block := range blocks {
		if block == nil {
			// Skip blocks that failed to fetch
			continue
		}
		
		select {
		case i.blockChan <- block:
			log.Debug().Int64("height", startHeight+int64(height)).Str("hash", block.Hash().Hex()).Msg("Published catchup block")
		case <-i.ctx.Done():
			return i.ctx.Err()
		}
	}
	
	// Update our latest height to the current tip
	i.mu.Lock()
	i.latestHeight = currentTip
	i.mu.Unlock()
	
	log.Info().
		Int64("blocksProcessed", int64(len(blocks))).
		Int64("latestHeight", currentTip).
		Msg("Initial catchup completed")
	
	return nil
}

// fetchBlocksInParallel fetches a range of blocks concurrently while maintaining order
func (i *Indexer) fetchBlocksInParallel(startHeight, endHeight int64) ([]rpctypes.PolyBlock, error) {
	if startHeight > endHeight {
		return nil, nil
	}
	
	blockCount := endHeight - startHeight + 1
	blocks := make([]rpctypes.PolyBlock, blockCount)
	var wg sync.WaitGroup
	
	log.Debug().
		Int64("startHeight", startHeight).
		Int64("endHeight", endHeight).
		Int64("blockCount", blockCount).
		Int("maxConcurrency", i.maxConcurrency).
		Msg("Starting parallel block fetch")
	
	// Fetch blocks in parallel with concurrency control
	for idx := int64(0); idx < blockCount; idx++ {
		select {
		case <-i.ctx.Done():
			return nil, i.ctx.Err()
		case i.workerSem <- struct{}{}: // Acquire semaphore
		}
		
		wg.Add(1)
		go func(index int64, height int64) {
			defer func() {
				<-i.workerSem // Release semaphore
				wg.Done()
			}()
			
			block, err := i.store.GetBlockByNumber(i.ctx, big.NewInt(height))
			if err != nil {
				log.Error().
					Err(err).
					Int64("height", height).
					Msg("Error fetching block in parallel")
				// Leave blocks[index] as nil to indicate failure
				return
			}
			
			blocks[index] = block
			log.Trace().
				Int64("height", height).
				Str("hash", block.Hash().Hex()).
				Msg("Fetched block in parallel")
		}(idx, startHeight+idx)
	}
	
	// Wait for all goroutines to complete
	wg.Wait()
	
	// Count successful fetches
	successCount := int64(0)
	for _, block := range blocks {
		if block != nil {
			successCount++
		}
	}
	
	log.Debug().
		Int64("requested", blockCount).
		Int64("successful", successCount).
		Int64("failed", blockCount-successCount).
		Msg("Parallel block fetch completed")
	
	return blocks, nil
}
