package chainstore

import (
	"context"
	"math/big"
	"time"

	"github.com/0xPolygon/polygon-cli/rpctypes"
	"github.com/ethereum/go-ethereum/common"
)

// FeeHistoryResult represents the result of eth_feeHistory
type FeeHistoryResult struct {
	OldestBlock   *big.Int     `json:"oldestBlock"`
	BaseFeePerGas []*big.Int   `json:"baseFeePerGas"`
	GasUsedRatio  []float64    `json:"gasUsedRatio"`
	Reward        [][]*big.Int `json:"reward,omitempty"`
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
	GetClientVersion(ctx context.Context) (string, error)

	// Semi-static info (cache for minutes)
	GetSyncStatus(ctx context.Context) (interface{}, error)
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
	GetTxPoolStatus(ctx context.Context) (map[string]interface{}, error)
	GetNetPeerCount(ctx context.Context) (*big.Int, error)

	// === CAPABILITY & MANAGEMENT ===
	IsMethodSupported(method string) bool
	RefreshCapabilities(ctx context.Context) error
	GetSupportedMethods() []string

	// === CONNECTION INFO ===
	GetRPCURL() string

	// === SIGNATURE LOOKUP ===
	// GetSignature retrieves function/event signatures from 4byte.directory
	GetSignature(ctx context.Context, hexSignature string) ([]Signature, error)

	// Close closes the store and releases any resources
	Close() error
}

// Signature represents a function or event signature from 4byte.directory
type Signature struct {
	ID             int       `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	TextSignature  string    `json:"text_signature"`
	HexSignature   string    `json:"hex_signature"`
	BytesSignature string    `json:"bytes_signature"`
}

// SignatureResponse represents the paginated response from 4byte.directory API
type SignatureResponse struct {
	Count    int         `json:"count"`
	Next     *string     `json:"next"`
	Previous *string     `json:"previous"`
	Results  []Signature `json:"results"`
}
