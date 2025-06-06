package chainstore

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/0xPolygon/polygon-cli/rpctypes"
)

// FeeHistoryResult represents the result of eth_feeHistory
type FeeHistoryResult struct {
	OldestBlock  *big.Int     `json:"oldestBlock"`
	BaseFeePerGas []*big.Int `json:"baseFeePerGas"`
	GasUsedRatio []float64   `json:"gasUsedRatio"`
	Reward       [][]*big.Int `json:"reward,omitempty"`
}

// ChainStore defines the unified interface for storing and retrieving all chain-related data
type ChainStore interface {
	// === BLOCK DATA (existing BlockStore methods) ===
	
	// GetBlock retrieves a block by hash or number
	GetBlock(ctx context.Context, blockHashOrNumber interface{}) (rpctypes.PolyBlock, error)

	// GetTransaction retrieves a transaction by hash
	GetTransaction(ctx context.Context, txHash common.Hash) (rpctypes.PolyTransaction, error)

	// GetReceipt retrieves a transaction receipt by transaction hash
	GetReceipt(ctx context.Context, txHash common.Hash) (rpctypes.PolyReceipt, error)

	// GetLatestBlock retrieves the most recent block
	GetLatestBlock(ctx context.Context) (rpctypes.PolyBlock, error)

	// GetBlockByNumber retrieves a block by its number
	GetBlockByNumber(ctx context.Context, number *big.Int) (rpctypes.PolyBlock, error)

	// GetBlockByHash retrieves a block by its hash
	GetBlockByHash(ctx context.Context, hash common.Hash) (rpctypes.PolyBlock, error)

	// === CHAIN METADATA (new functionality) ===
	
	// Static info (fetch once, cache indefinitely)
	GetChainID(ctx context.Context) (*big.Int, error)
	
	// Semi-static info (cache for minutes)
	GetSafeBlock(ctx context.Context) (*big.Int, error)
	GetFinalizedBlock(ctx context.Context) (*big.Int, error)
	
	// Block-aligned info (cache per block)
	GetBaseFee(ctx context.Context) (*big.Int, error)
	GetBaseFeeForBlock(ctx context.Context, blockNumber *big.Int) (*big.Int, error)
	
	// Frequent info (cache for seconds)
	GetGasPrice(ctx context.Context) (*big.Int, error)
	GetFeeHistory(ctx context.Context, blockCount int, newestBlock string, rewardPercentiles []float64) (*FeeHistoryResult, error)
	
	// Very frequent info (minimal cache)
	GetPendingTransactionCount(ctx context.Context) (*big.Int, error)
	GetQueuedTransactionCount(ctx context.Context) (*big.Int, error)
	
	// === CAPABILITY & MANAGEMENT ===
	IsMethodSupported(method string) bool
	RefreshCapabilities(ctx context.Context) error
	GetSupportedMethods() []string
	
	// Close closes the store and releases any resources
	Close() error
}