package blockstore

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/0xPolygon/polygon-cli/rpctypes"
)

// PassthroughStore is a blockstore implementation that doesn't store anything
// and always passes through requests directly to the RPC endpoint
type PassthroughStore struct {
	client *rpc.Client
}

// NewPassthroughStore creates a new passthrough store with the given RPC client
func NewPassthroughStore(rpcURL string) (*PassthroughStore, error) {
	client, err := rpc.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC: %w", err)
	}
	
	return &PassthroughStore{
		client: client,
	}, nil
}

// GetBlock retrieves a block by hash or number
func (s *PassthroughStore) GetBlock(ctx context.Context, blockHashOrNumber interface{}) (rpctypes.PolyBlock, error) {
	var raw rpctypes.RawBlockResponse
	
	switch v := blockHashOrNumber.(type) {
	case common.Hash:
		err := s.client.CallContext(ctx, &raw, "eth_getBlockByHash", v, true)
		if err != nil {
			return nil, fmt.Errorf("failed to get block by hash: %w", err)
		}
	case *big.Int:
		err := s.client.CallContext(ctx, &raw, "eth_getBlockByNumber", fmt.Sprintf("0x%x", v), true)
		if err != nil {
			return nil, fmt.Errorf("failed to get block by number: %w", err)
		}
	case int64:
		err := s.client.CallContext(ctx, &raw, "eth_getBlockByNumber", fmt.Sprintf("0x%x", v), true)
		if err != nil {
			return nil, fmt.Errorf("failed to get block by number: %w", err)
		}
	case string:
		// Could be "latest", "pending", "earliest" or a hex number
		err := s.client.CallContext(ctx, &raw, "eth_getBlockByNumber", v, true)
		if err != nil {
			return nil, fmt.Errorf("failed to get block by tag: %w", err)
		}
	default:
		return nil, fmt.Errorf("invalid block identifier type: %T", blockHashOrNumber)
	}
	
	return rpctypes.NewPolyBlock(&raw), nil
}

// GetTransaction retrieves a transaction by hash
func (s *PassthroughStore) GetTransaction(ctx context.Context, txHash common.Hash) (rpctypes.PolyTransaction, error) {
	var raw rpctypes.RawTransactionResponse
	err := s.client.CallContext(ctx, &raw, "eth_getTransactionByHash", txHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}
	
	return rpctypes.NewPolyTransaction(&raw), nil
}

// GetReceipt retrieves a transaction receipt by transaction hash
func (s *PassthroughStore) GetReceipt(ctx context.Context, txHash common.Hash) (rpctypes.PolyReceipt, error) {
	var raw rpctypes.RawTxReceipt
	err := s.client.CallContext(ctx, &raw, "eth_getTransactionReceipt", txHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get receipt: %w", err)
	}
	
	return rpctypes.NewPolyReceipt(&raw), nil
}

// GetLatestBlock retrieves the most recent block
func (s *PassthroughStore) GetLatestBlock(ctx context.Context) (rpctypes.PolyBlock, error) {
	return s.GetBlock(ctx, "latest")
}

// GetBlockByNumber retrieves a block by its number
func (s *PassthroughStore) GetBlockByNumber(ctx context.Context, number *big.Int) (rpctypes.PolyBlock, error) {
	return s.GetBlock(ctx, number)
}

// GetBlockByHash retrieves a block by its hash
func (s *PassthroughStore) GetBlockByHash(ctx context.Context, hash common.Hash) (rpctypes.PolyBlock, error) {
	return s.GetBlock(ctx, hash)
}

// Close closes the store and releases any resources
func (s *PassthroughStore) Close() error {
	if s.client != nil {
		s.client.Close()
	}
	return nil
}