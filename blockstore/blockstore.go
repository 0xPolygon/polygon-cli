package blockstore

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/0xPolygon/polygon-cli/rpctypes"
)

// BlockStore defines the interface for storing and retrieving blockchain data
type BlockStore interface {
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

	// Close closes the store and releases any resources
	Close() error
}