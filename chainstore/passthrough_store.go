package chainstore

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/0xPolygon/polygon-cli/rpctypes"
	"github.com/rs/zerolog/log"
)

// PassthroughStore is a chainstore implementation that doesn't store anything
// and passes through requests directly to the RPC endpoint with caching
type PassthroughStore struct {
	client       *rpc.Client
	cache        *ChainCache
	capabilities *CapabilityManager
	config       *ChainStoreConfig
	rpcURL       string
}

// NewPassthroughStore creates a new passthrough store with the given RPC client
func NewPassthroughStore(rpcURL string) (*PassthroughStore, error) {
	return NewPassthroughStoreWithConfig(rpcURL, DefaultChainStoreConfig())
}

// NewPassthroughStoreWithConfig creates a new passthrough store with custom configuration
func NewPassthroughStoreWithConfig(rpcURL string, config *ChainStoreConfig) (*PassthroughStore, error) {
	client, err := rpc.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC: %w", err)
	}
	
	store := &PassthroughStore{
		client:       client,
		cache:        NewChainCache(),
		capabilities: NewCapabilityManager(client, config.CapabilityTTL),
		config:       config,
		rpcURL:       rpcURL,
	}
	
	// Initialize capabilities in background
	go func() {
		ctx := context.Background()
		if err := store.capabilities.RefreshCapabilities(ctx); err != nil {
			log.Error().Err(err).Msg("Failed to refresh capabilities")
		}
	}()
	
	return store, nil
}

// === BLOCK DATA (existing BlockStore methods) ===

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

// === CHAIN METADATA (new functionality) ===

// GetChainID retrieves the chain ID (cached indefinitely)
func (s *PassthroughStore) GetChainID(ctx context.Context) (*big.Int, error) {
	// Check cache first
	if chainID, valid := s.cache.GetChainID(); valid {
		return chainID, nil
	}
	
	// Not supported, return error
	if !s.capabilities.IsMethodSupported("eth_chainId") {
		return nil, fmt.Errorf("eth_chainId method not supported")
	}
	
	var result string
	err := s.client.CallContext(ctx, &result, "eth_chainId")
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}
	
	chainID, err := hexToBigInt(result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse chain ID: %w", err)
	}
	
	// Cache the result
	s.cache.SetChainID(chainID)
	
	log.Debug().Str("chainID", chainID.String()).Msg("Retrieved and cached chain ID")
	return chainID, nil
}

// GetSafeBlock retrieves the safe block number (cached semi-statically)
func (s *PassthroughStore) GetSafeBlock(ctx context.Context) (*big.Int, error) {
	// Check cache first
	if block, valid := s.cache.GetSafeBlock(s.config.SemiStaticTTL); valid {
		return block, nil
	}
	
	// Try engine_forkchoiceUpdatedV3 first (if supported)
	if s.capabilities.IsMethodSupported("engine_forkchoiceUpdatedV3") {
		var result map[string]interface{}
		err := s.client.CallContext(ctx, &result, "eth_getBlockByNumber", "safe", false)
		if err == nil && result != nil {
			if numberHex, ok := result["number"].(string); ok {
				if block, err := hexToBigInt(numberHex); err == nil {
					s.cache.SetSafeBlock(block, s.config.SemiStaticTTL)
					return block, nil
				}
			}
		}
	}
	
	// Fallback: return latest block number minus some depth
	latestBlock, err := s.GetLatestBlock(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block for safe block estimation: %w", err)
	}
	
	// Estimate safe block as latest - 64 blocks
	latestNum := latestBlock.Number()
	safeNum := new(big.Int).Sub(latestNum, big.NewInt(64))
	if safeNum.Sign() < 0 {
		safeNum = big.NewInt(0)
	}
	
	s.cache.SetSafeBlock(safeNum, s.config.SemiStaticTTL)
	return safeNum, nil
}

// GetFinalizedBlock retrieves the finalized block number (cached semi-statically)
func (s *PassthroughStore) GetFinalizedBlock(ctx context.Context) (*big.Int, error) {
	// Check cache first
	if block, valid := s.cache.GetFinalizedBlock(s.config.SemiStaticTTL); valid {
		return block, nil
	}
	
	// Try to get finalized block via eth_getBlockByNumber
	var result map[string]interface{}
	err := s.client.CallContext(ctx, &result, "eth_getBlockByNumber", "finalized", false)
	if err == nil && result != nil {
		if numberHex, ok := result["number"].(string); ok {
			if block, parseErr := hexToBigInt(numberHex); parseErr == nil {
				s.cache.SetFinalizedBlock(block, s.config.SemiStaticTTL)
				return block, nil
			}
		}
	}
	
	// Fallback: return latest block number minus some depth
	latestBlock, err := s.GetLatestBlock(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block for finalized block estimation: %w", err)
	}
	
	// Estimate finalized block as latest - 128 blocks
	latestNum := latestBlock.Number()
	finalizedNum := new(big.Int).Sub(latestNum, big.NewInt(128))
	if finalizedNum.Sign() < 0 {
		finalizedNum = big.NewInt(0)
	}
	
	s.cache.SetFinalizedBlock(finalizedNum, s.config.SemiStaticTTL)
	return finalizedNum, nil
}

// GetBaseFee retrieves the current base fee (cached per block)
func (s *PassthroughStore) GetBaseFee(ctx context.Context) (*big.Int, error) {
	// Get latest block to check if we have base fee cached for current block
	latestBlock, err := s.GetLatestBlock(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block: %w", err)
	}
	
	latestNum := latestBlock.Number()
	
	// Check cache first
	if baseFee, valid := s.cache.GetBaseFee(latestNum); valid {
		return baseFee, nil
	}
	
	// Get base fee from the latest block
	if baseFee := latestBlock.BaseFee(); baseFee != nil {
		s.cache.SetBaseFee(baseFee, latestNum)
		return baseFee, nil
	}
	
	return nil, fmt.Errorf("base fee not available in latest block")
}

// GetBaseFeeForBlock retrieves the base fee for a specific block
func (s *PassthroughStore) GetBaseFeeForBlock(ctx context.Context, blockNumber *big.Int) (*big.Int, error) {
	block, err := s.GetBlockByNumber(ctx, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get block %s: %w", blockNumber.String(), err)
	}
	
	if baseFee := block.BaseFee(); baseFee != nil {
		return baseFee, nil
	}
	
	return nil, fmt.Errorf("base fee not available in block %s", blockNumber.String())
}

// GetGasPrice retrieves the current gas price (cached frequently)
func (s *PassthroughStore) GetGasPrice(ctx context.Context) (*big.Int, error) {
	// Check cache first
	if gasPrice, valid := s.cache.GetGasPrice(s.config.FrequentTTL); valid {
		return gasPrice, nil
	}
	
	if !s.capabilities.IsMethodSupported("eth_gasPrice") {
		return nil, fmt.Errorf("eth_gasPrice method not supported")
	}
	
	var result string
	err := s.client.CallContext(ctx, &result, "eth_gasPrice")
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}
	
	gasPrice, err := hexToBigInt(result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse gas price: %w", err)
	}
	
	// Cache the result
	s.cache.SetGasPrice(gasPrice, s.config.FrequentTTL)
	
	return gasPrice, nil
}

// GetFeeHistory retrieves fee history (cached frequently)
func (s *PassthroughStore) GetFeeHistory(ctx context.Context, blockCount int, newestBlock string, rewardPercentiles []float64) (*FeeHistoryResult, error) {
	// Check cache first (simple cache key for now)
	if feeHistory, valid := s.cache.GetFeeHistory(s.config.FrequentTTL); valid {
		return feeHistory, nil
	}
	
	if !s.capabilities.IsMethodSupported("eth_feeHistory") {
		return nil, fmt.Errorf("eth_feeHistory method not supported")
	}
	
	var result FeeHistoryResult
	err := s.client.CallContext(ctx, &result, "eth_feeHistory", fmt.Sprintf("0x%x", blockCount), newestBlock, rewardPercentiles)
	if err != nil {
		return nil, fmt.Errorf("failed to get fee history: %w", err)
	}
	
	// Cache the result
	s.cache.SetFeeHistory(&result, s.config.FrequentTTL)
	
	return &result, nil
}

// GetPendingTransactionCount retrieves pending transaction count (cached very frequently)
func (s *PassthroughStore) GetPendingTransactionCount(ctx context.Context) (*big.Int, error) {
	// Check cache first
	if count, valid := s.cache.GetPendingTxCount(s.config.VeryFrequentTTL); valid {
		return count, nil
	}
	
	if !s.capabilities.IsMethodSupported("txpool_status") {
		return nil, fmt.Errorf("txpool_status method not supported")
	}
	
	var result map[string]interface{}
	err := s.client.CallContext(ctx, &result, "txpool_status")
	if err != nil {
		return nil, fmt.Errorf("failed to get txpool status: %w", err)
	}
	
	// Parse pending count
	pendingCount := big.NewInt(0)
	if pending, ok := result["pending"]; ok {
		switch v := pending.(type) {
		case string:
			if count, err := hexToBigInt(v); err == nil {
				pendingCount = count
			}
		case float64:
			pendingCount = big.NewInt(int64(v))
		case json.Number:
			if i, err := v.Int64(); err == nil {
				pendingCount = big.NewInt(i)
			}
		}
	}
	
	// Cache the result
	s.cache.SetPendingTxCount(pendingCount, s.config.VeryFrequentTTL)
	
	return pendingCount, nil
}

// GetQueuedTransactionCount retrieves queued transaction count (cached very frequently)
func (s *PassthroughStore) GetQueuedTransactionCount(ctx context.Context) (*big.Int, error) {
	// Check cache first
	if count, valid := s.cache.GetQueuedTxCount(s.config.VeryFrequentTTL); valid {
		return count, nil
	}
	
	if !s.capabilities.IsMethodSupported("txpool_status") {
		return nil, fmt.Errorf("txpool_status method not supported")
	}
	
	var result map[string]interface{}
	err := s.client.CallContext(ctx, &result, "txpool_status")
	if err != nil {
		return nil, fmt.Errorf("failed to get txpool status: %w", err)
	}
	
	// Parse queued count
	queuedCount := big.NewInt(0)
	if queued, ok := result["queued"]; ok {
		switch v := queued.(type) {
		case string:
			if count, err := hexToBigInt(v); err == nil {
				queuedCount = count
			}
		case float64:
			queuedCount = big.NewInt(int64(v))
		case json.Number:
			if i, err := v.Int64(); err == nil {
				queuedCount = big.NewInt(i)
			}
		}
	}
	
	// Cache the result
	s.cache.SetQueuedTxCount(queuedCount, s.config.VeryFrequentTTL)
	
	return queuedCount, nil
}

// GetTxPoolStatus retrieves the full txpool status (cached very frequently)
func (s *PassthroughStore) GetTxPoolStatus(ctx context.Context) (map[string]interface{}, error) {
	if !s.capabilities.IsMethodSupported("txpool_status") {
		return nil, fmt.Errorf("txpool_status method not supported")
	}
	
	var result map[string]interface{}
	err := s.client.CallContext(ctx, &result, "txpool_status")
	if err != nil {
		return nil, fmt.Errorf("failed to get txpool status: %w", err)
	}
	
	return result, nil
}

// GetNetPeerCount retrieves the number of connected peers (cached very frequently)
func (s *PassthroughStore) GetNetPeerCount(ctx context.Context) (*big.Int, error) {
	if !s.capabilities.IsMethodSupported("net_peerCount") {
		return nil, fmt.Errorf("net_peerCount method not supported")
	}
	
	var result string
	err := s.client.CallContext(ctx, &result, "net_peerCount")
	if err != nil {
		return nil, fmt.Errorf("failed to get peer count: %w", err)
	}
	
	peerCount, err := hexToBigInt(result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse peer count: %w", err)
	}
	
	return peerCount, nil
}

// === CAPABILITY & MANAGEMENT ===

// IsMethodSupported checks if a method is supported
func (s *PassthroughStore) IsMethodSupported(method string) bool {
	return s.capabilities.IsMethodSupported(method)
}

// RefreshCapabilities refreshes the capability cache
func (s *PassthroughStore) RefreshCapabilities(ctx context.Context) error {
	return s.capabilities.RefreshCapabilities(ctx)
}

// GetSupportedMethods returns all supported methods
func (s *PassthroughStore) GetSupportedMethods() []string {
	return s.capabilities.GetSupportedMethods()
}

// GetRPCURL returns the RPC endpoint URL
func (s *PassthroughStore) GetRPCURL() string {
	return s.rpcURL
}

// Close closes the store and releases any resources
func (s *PassthroughStore) Close() error {
	if s.client != nil {
		s.client.Close()
	}
	return nil
}

// === UTILITY FUNCTIONS ===

// hexToBigInt converts a hex string to big.Int
func hexToBigInt(hex string) (*big.Int, error) {
	if len(hex) >= 2 && hex[:2] == "0x" {
		hex = hex[2:]
	}
	
	result := big.NewInt(0)
	result, ok := result.SetString(hex, 16)
	if !ok {
		return nil, fmt.Errorf("invalid hex string: %s", hex)
	}
	
	return result, nil
}