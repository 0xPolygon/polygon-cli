package sensor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"

	"github.com/0xPolygon/polygon-cli/p2p"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
)

// rpcRequest represents a JSON-RPC 2.0 request message.
type rpcRequest struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
	ID      any    `json:"id"`
}

// rpcResponse represents a JSON-RPC 2.0 response message.
type rpcResponse struct {
	JSONRPC string    `json:"jsonrpc"`
	Result  any       `json:"result,omitempty"`
	Error   *rpcError `json:"error,omitempty"`
	ID      any       `json:"id"`
}

// rpcError represents a JSON-RPC 2.0 error object.
type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// rpcProxy holds configuration for proxying RPC requests to an upstream server.
type rpcProxy struct {
	rpcURL     string
	httpClient *http.Client
}

// rpcParams holds shared parameters for processing JSON-RPC requests.
type rpcParams struct {
	conns    *p2p.Conns
	chainID  *big.Int
	gpo      *p2p.GasPriceOracle
	proxy    *rpcProxy
	requests *prometheus.CounterVec
}

// handleRPC sets up the JSON-RPC server for receiving and broadcasting transactions.
// It handles eth_sendRawTransaction requests, validates transaction signatures,
// and broadcasts valid transactions to all connected peers.
// Supports both single requests and batch requests per JSON-RPC 2.0 specification.
// If proxyRPC is enabled, unsupported methods are forwarded to the upstream rpcURL.
func handleRPC(conns *p2p.Conns, networkID uint64) {
	// Use network ID as chain ID for signature validation
	chainID := new(big.Int).SetUint64(networkID)
	gpo := p2p.NewGasPriceOracle(conns)

	params := &rpcParams{
		conns:    conns,
		chainID:  chainID,
		gpo:      gpo,
		requests: p2p.NewRPCRequestsCounter(),
	}

	if inputSensorParams.ProxyRPC {
		params.proxy = &rpcProxy{
			rpcURL: inputSensorParams.RPC,
			httpClient: &http.Client{
				Timeout: inputSensorParams.ProxyRPCTimeout,
			},
		}
		log.Info().
			Str("rpc", inputSensorParams.RPC).
			Dur("timeout", inputSensorParams.ProxyRPCTimeout).
			Msg("RPC proxy enabled")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, 5*1024*1024) // 5MB limit
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, -32700, "Parse error", nil)
			return
		}

		// Check if this is a batch request (starts with '[') or single request
		trimmed := bytes.TrimSpace(body)
		if len(trimmed) > 0 && trimmed[0] == '[' {
			// Handle batch request
			handleBatchRequest(w, r, body, params)
			return
		}

		// Parse single JSON-RPC request
		var req rpcRequest
		if err := json.Unmarshal(body, &req); err != nil {
			writeError(w, -32700, "Parse error", nil)
			return
		}

		// Process request
		var txs types.Transactions
		resp := processRequest(req, params, &txs)

		// If method not found and proxy is enabled, forward to upstream
		if isMethodNotFound(resp) && params.proxy != nil {
			params.requests.WithLabelValues(req.Method, "true").Inc()
			proxyRPCRequest(w, r, body, params.proxy)
			return
		}

		params.requests.WithLabelValues(req.Method, "false").Inc()

		// Enqueue transactions for async broadcast
		if len(txs) > 0 {
			params.conns.EnqueueTxBroadcast(txs)
		}

		// Write response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(&resp); err != nil {
			log.Error().Err(err).Msg("Failed to encode response")
		}
	})

	addr := fmt.Sprintf(":%d", inputSensorParams.RPCPort)
	log.Info().Str("addr", addr).Msg("Starting JSON-RPC server")
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Error().Err(err).Msg("Failed to start RPC server")
	}
}

// writeError writes a JSON-RPC 2.0 error response with the specified code, message, and request ID.
func writeError(w http.ResponseWriter, code int, message string, id any) {
	w.Header().Set("Content-Type", "application/json")
	response := rpcResponse{
		JSONRPC: "2.0",
		Error: &rpcError{
			Code:    code,
			Message: message,
		},
		ID: id,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error().Err(err).Msg("Failed to encode error response")
	}
}

// proxyRPCRequest forwards a JSON-RPC request to the upstream RPC server and streams
// the response back to the client. Used for methods not handled locally.
func proxyRPCRequest(w http.ResponseWriter, r *http.Request, body []byte, config *rpcProxy) {
	req, err := http.NewRequestWithContext(r.Context(), http.MethodPost, config.rpcURL, bytes.NewReader(body))
	if err != nil {
		log.Error().Err(err).Msg("Failed to create proxy request")
		writeError(w, -32603, "Internal error: failed to create proxy request", nil)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := config.httpClient.Do(req)
	if err != nil {
		log.Error().Err(err).Str("rpc", config.rpcURL).Msg("Proxy request failed")
		if r.Context().Err() != nil {
			writeError(w, -32603, "Request cancelled or timed out", nil)
			return
		}
		writeError(w, -32603, fmt.Sprintf("Upstream RPC error: %v", err), nil)
		return
	}
	defer func() { _ = resp.Body.Close() }()

	// Copy response headers and status
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)

	// Stream response body
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Error().Err(err).Msg("Failed to copy proxy response")
	}
}

// proxyBatchRequests sends only the specified requests to the proxy and updates
// responses in place.
func proxyBatchRequests(r *http.Request, requests []rpcRequest, responses []rpcResponse, indices []int, proxy *rpcProxy) {
	proxyRequests := make([]rpcRequest, len(indices))
	for i, idx := range indices {
		proxyRequests[i] = requests[idx]
	}

	body, err := json.Marshal(proxyRequests)
	if err != nil {
		log.Error().Err(err).Msg("Failed to encode proxy batch")
		return
	}

	req, err := http.NewRequestWithContext(r.Context(), http.MethodPost, proxy.rpcURL, bytes.NewReader(body))
	if err != nil {
		log.Error().Err(err).Msg("Failed to create proxy request")
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := proxy.httpClient.Do(req)
	if err != nil {
		log.Error().Err(err).Str("rpc", proxy.rpcURL).Msg("Proxy batch request failed")
		return
	}
	defer func() { _ = resp.Body.Close() }()

	var proxyResponses []rpcResponse
	if err := json.NewDecoder(resp.Body).Decode(&proxyResponses); err != nil {
		log.Error().Err(err).Msg("Failed to decode proxy batch response")
		return
	}

	for i, idx := range indices {
		if i < len(proxyResponses) {
			responses[idx] = proxyResponses[i]
		}
	}
}

// handleBatchRequest processes JSON-RPC 2.0 batch requests.
// For eth_sendRawTransaction requests, it collects valid transactions for batch
// broadcasting. Supported methods are handled locally; unsupported methods are
// proxied if configured. This ensures transactions are always broadcast locally
// and never lost.
func handleBatchRequest(w http.ResponseWriter, r *http.Request, body []byte, params *rpcParams) {
	var requests []rpcRequest
	if err := json.Unmarshal(body, &requests); err != nil {
		writeError(w, -32700, "Parse error", nil)
		return
	}

	if len(requests) == 0 {
		writeError(w, -32600, "Invalid request: empty batch", nil)
		return
	}

	responses := make([]rpcResponse, len(requests))
	var txs types.Transactions
	var indices []int

	for i, req := range requests {
		resp := processRequest(req, params, &txs)
		responses[i] = resp

		if isMethodNotFound(resp) {
			if params.proxy != nil {
				indices = append(indices, i)
			}
			params.requests.WithLabelValues(req.Method, "true").Inc()
		} else {
			params.requests.WithLabelValues(req.Method, "false").Inc()
		}
	}

	// Enqueue transactions for async broadcast
	if len(txs) > 0 {
		params.conns.EnqueueTxBroadcast(txs)
	}

	if len(indices) > 0 {
		proxyBatchRequests(r, requests, responses, indices, params.proxy)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responses); err != nil {
		log.Error().Err(err).Msg("Failed to encode batch response")
	}
}

// newResultResponse creates a success response.
func newResultResponse(result, id any) rpcResponse {
	return rpcResponse{JSONRPC: "2.0", Result: result, ID: id}
}

// newErrorResponse creates an error response.
func newErrorResponse(err *rpcError, id any) rpcResponse {
	return rpcResponse{JSONRPC: "2.0", Error: err, ID: id}
}

const rpcMethodNotFoundCode = -32601

// newMethodNotFoundResponse creates a method not found error response.
func newMethodNotFoundResponse(id any) rpcResponse {
	return rpcResponse{
		JSONRPC: "2.0",
		Error:   &rpcError{Code: rpcMethodNotFoundCode, Message: "Method not found"},
		ID:      id,
	}
}

// isMethodNotFound returns true if the response is a method not found error.
func isMethodNotFound(resp rpcResponse) bool {
	return resp.Error != nil && resp.Error.Code == rpcMethodNotFoundCode
}

// processRequest handles a single RPC request and returns a response.
// For eth_sendRawTransaction, valid transactions are appended to txs for batch broadcasting.
// Returns a method not found response if the method is not handled locally.
func processRequest(req rpcRequest, params *rpcParams, txs *types.Transactions) rpcResponse {
	switch req.Method {
	case "eth_sendRawTransaction":
		tx, errResp := validateTx(req, params.chainID)
		if tx == nil {
			return errResp
		}
		if txs != nil {
			*txs = append(*txs, tx)
		}
		return newResultResponse(tx.Hash().Hex(), req.ID)

	case "eth_chainId":
		return newResultResponse(hexutil.EncodeBig(params.chainID), req.ID)

	case "eth_blockNumber":
		head := params.conns.HeadBlock()
		if head.Block == nil {
			return newResultResponse(nil, req.ID)
		}
		return newResultResponse(hexutil.EncodeUint64(head.Block.NumberU64()), req.ID)

	case "eth_gasPrice":
		return newResultResponse(hexutil.EncodeBig(params.gpo.SuggestGasPrice()), req.ID)

	case "eth_maxPriorityFeePerGas":
		tip := params.gpo.SuggestGasTipCap()
		if tip == nil {
			tip = big.NewInt(1e9) // Default to 1 gwei
		}
		return newResultResponse(hexutil.EncodeBig(tip), req.ID)

	case "eth_getBlockByHash":
		result, err := getBlockByHash(req, params.conns)
		return handleMethodResult(result, err, req.ID)

	case "eth_getBlockByNumber":
		result, err := getBlockByNumber(req, params.conns)
		return handleMethodResult(result, err, req.ID)

	case "eth_getTransactionByHash":
		result, err := getTransactionByHash(req, params.conns)
		return handleMethodResult(result, err, req.ID)

	case "eth_getTransactionByBlockHashAndIndex":
		result, err := getTransactionByBlockHashAndIndex(req, params.conns)
		return handleMethodResult(result, err, req.ID)

	case "eth_getBlockTransactionCountByHash":
		result, err := getBlockTransactionCountByHash(req, params.conns)
		return handleMethodResult(result, err, req.ID)

	case "eth_getUncleCountByBlockHash":
		result, err := getUncleCountByBlockHash(req, params.conns)
		return handleMethodResult(result, err, req.ID)

	default:
		return newMethodNotFoundResponse(req.ID)
	}
}

// handleMethodResult converts a method's result and error into an rpcResponse.
func handleMethodResult(result any, err *rpcError, id any) rpcResponse {
	if err != nil {
		return newErrorResponse(err, id)
	}
	return newResultResponse(result, id)
}

// validateTx validates a transaction from a JSON-RPC request by decoding the raw
// transaction hex, unmarshaling it, and verifying the signature. Returns the transaction if valid
// (with an empty response), or nil transaction with an error response if validation fails.
func validateTx(req rpcRequest, chainID *big.Int) (*types.Transaction, rpcResponse) {
	invalidParams := func(msg string) rpcResponse {
		return newErrorResponse(&rpcError{Code: -32602, Message: msg}, req.ID)
	}

	if len(req.Params) == 0 {
		return nil, invalidParams("Invalid params: missing raw transaction")
	}

	hex, ok := req.Params[0].(string)
	if !ok {
		return nil, invalidParams("Invalid params: raw transaction must be a hex string")
	}

	bytes, err := hexutil.Decode(hex)
	if err != nil {
		return nil, invalidParams(fmt.Sprintf("Invalid transaction hex: %v", err))
	}

	tx := new(types.Transaction)
	if err = tx.UnmarshalBinary(bytes); err != nil {
		return nil, invalidParams(fmt.Sprintf("Invalid transaction encoding: %v", err))
	}

	signer := types.LatestSignerForChainID(chainID)
	sender, err := types.Sender(signer, tx)
	if err != nil {
		return nil, invalidParams(fmt.Sprintf("Invalid transaction signature: %v", err))
	}

	to := "nil"
	if tx.To() != nil {
		to = tx.To().Hex()
	}

	log.Debug().
		Str("hash", tx.Hash().Hex()).
		Str("from", sender.Hex()).
		Str("to", to).
		Str("value", tx.Value().String()).
		Uint64("gas", tx.Gas()).
		Msg("Validated transaction")

	return tx, rpcResponse{}
}

// parseFullTxParam extracts the fullTx boolean from params[1], defaulting to false.
func parseFullTxParam(params []any) bool {
	if len(params) >= 2 {
		if fullTx, ok := params[1].(bool); ok {
			return fullTx
		}
	}
	return false
}

// getBlockByHash retrieves a block by its hash from the cache.
func getBlockByHash(req rpcRequest, conns *p2p.Conns) (any, *rpcError) {
	if len(req.Params) < 1 {
		return nil, &rpcError{Code: -32602, Message: "missing block hash parameter"}
	}

	hashStr, ok := req.Params[0].(string)
	if !ok {
		return nil, &rpcError{Code: -32602, Message: "invalid block hash parameter"}
	}

	hash := common.HexToHash(hashStr)
	cache, ok := conns.Blocks().Get(hash)
	if !ok {
		return nil, nil // Return null for not found (per spec)
	}

	return formatBlockResponse(hash, cache, parseFullTxParam(req.Params)), nil
}

// getBlockByNumber retrieves a block by its number from the cache.
func getBlockByNumber(req rpcRequest, conns *p2p.Conns) (any, *rpcError) {
	if len(req.Params) < 1 {
		return nil, &rpcError{Code: -32602, Message: "missing block number parameter"}
	}

	blockNumParam, ok := req.Params[0].(string)
	if !ok {
		return nil, &rpcError{Code: -32602, Message: "invalid block number parameter"}
	}

	var hash common.Hash
	var cache p2p.BlockCache
	var found bool

	switch blockNumParam {
	case "latest", "pending":
		head := conns.HeadBlock()
		if head.Block == nil {
			return nil, nil
		}
		hash = head.Block.Hash()
		cache, found = conns.Blocks().Get(hash)
		if !found {
			// Construct cache from head block
			txList, _ := rlp.EncodeToRawList([]*types.Transaction(head.Block.Transactions()))
			uncleList, _ := rlp.EncodeToRawList(head.Block.Uncles())
			cache = p2p.BlockCache{
				Header: head.Block.Header(),
				Body: &eth.BlockBody{
					Transactions: txList,
					Uncles:       uncleList,
				},
				TD: head.TD,
			}
			found = true
		}
	case "earliest":
		hash, cache, found = conns.GetBlockByNumber(0)
	default:
		num, err := hexutil.DecodeUint64(blockNumParam)
		if err != nil {
			return nil, &rpcError{Code: -32602, Message: "invalid block number: " + err.Error()}
		}
		hash, cache, found = conns.GetBlockByNumber(num)
	}

	if !found {
		return nil, nil
	}

	return formatBlockResponse(hash, cache, parseFullTxParam(req.Params)), nil
}

// getTransactionByHash retrieves a transaction by its hash from the cache.
func getTransactionByHash(req rpcRequest, conns *p2p.Conns) (any, *rpcError) {
	if len(req.Params) < 1 {
		return nil, &rpcError{Code: -32602, Message: "missing transaction hash parameter"}
	}

	hashStr, ok := req.Params[0].(string)
	if !ok {
		return nil, &rpcError{Code: -32602, Message: "invalid transaction hash parameter"}
	}

	hash := common.HexToHash(hashStr)

	// First check the transactions cache
	tx, ok := conns.GetTx(hash)
	if ok {
		return formatTransactionResponse(tx, common.Hash{}, nil, 0), nil
	}

	// Search in blocks for the transaction
	for _, blockHash := range conns.Blocks().Keys() {
		cache, ok := conns.Blocks().Peek(blockHash)
		if !ok || cache.Body == nil {
			continue
		}
		txs, err := cache.Body.Transactions.Items()
		if err != nil {
			continue
		}
		for i, tx := range txs {
			if tx.Hash() == hash {
				return formatTransactionResponse(tx, blockHash, cache.Header, uint64(i)), nil
			}
		}
	}

	return nil, nil
}

// getTransactionByBlockHashAndIndex retrieves a transaction by block hash and index.
func getTransactionByBlockHashAndIndex(req rpcRequest, conns *p2p.Conns) (any, *rpcError) {
	if len(req.Params) < 2 {
		return nil, &rpcError{Code: -32602, Message: "missing block hash or index parameter"}
	}

	hashStr, ok := req.Params[0].(string)
	if !ok {
		return nil, &rpcError{Code: -32602, Message: "invalid block hash parameter"}
	}

	indexStr, ok := req.Params[1].(string)
	if !ok {
		return nil, &rpcError{Code: -32602, Message: "invalid index parameter"}
	}

	index, err := hexutil.DecodeUint64(indexStr)
	if err != nil {
		return nil, &rpcError{Code: -32602, Message: "invalid index: " + err.Error()}
	}

	blockHash := common.HexToHash(hashStr)
	cache, ok := conns.Blocks().Get(blockHash)
	if !ok || cache.Body == nil {
		return nil, nil
	}

	txs, err := cache.Body.Transactions.Items()
	if err != nil || int(index) >= len(txs) {
		return nil, nil
	}

	tx := txs[index]
	return formatTransactionResponse(tx, blockHash, cache.Header, index), nil
}

// getBlockCacheByHashParam parses a block hash from params[0] and returns the block cache.
// Returns the cache and nil error on success, or nil cache and error on parse failure.
// If the block is not found, returns nil cache with nil error (per JSON-RPC spec).
func getBlockCacheByHashParam(req rpcRequest, conns *p2p.Conns) (p2p.BlockCache, *rpcError) {
	if len(req.Params) < 1 {
		return p2p.BlockCache{}, &rpcError{Code: -32602, Message: "missing block hash parameter"}
	}

	hashStr, ok := req.Params[0].(string)
	if !ok {
		return p2p.BlockCache{}, &rpcError{Code: -32602, Message: "invalid block hash parameter"}
	}

	hash := common.HexToHash(hashStr)
	cache, ok := conns.Blocks().Get(hash)
	if !ok || cache.Body == nil {
		return p2p.BlockCache{}, nil
	}

	return cache, nil
}

// getBlockTransactionCountByHash returns the transaction count in a block.
func getBlockTransactionCountByHash(req rpcRequest, conns *p2p.Conns) (any, *rpcError) {
	cache, err := getBlockCacheByHashParam(req, conns)
	if err != nil || cache.Body == nil {
		return nil, err
	}
	return hexutil.EncodeUint64(uint64(cache.Body.Transactions.Len())), nil
}

// getUncleCountByBlockHash returns the uncle count in a block.
func getUncleCountByBlockHash(req rpcRequest, conns *p2p.Conns) (any, *rpcError) {
	cache, err := getBlockCacheByHashParam(req, conns)
	if err != nil || cache.Body == nil {
		return nil, err
	}
	return hexutil.EncodeUint64(uint64(cache.Body.Uncles.Len())), nil
}

// formatBlockResponse formats a block cache into the Ethereum JSON-RPC block format.
func formatBlockResponse(hash common.Hash, cache p2p.BlockCache, fullTx bool) map[string]any {
	header := cache.Header
	if header == nil {
		return nil
	}

	result := map[string]any{
		"hash":             hash.Hex(),
		"number":           hexutil.EncodeUint64(header.Number.Uint64()),
		"parentHash":       header.ParentHash.Hex(),
		"nonce":            hexutil.Encode(header.Nonce[:]),
		"sha3Uncles":       header.UncleHash.Hex(),
		"logsBloom":        hexutil.Encode(header.Bloom.Bytes()),
		"transactionsRoot": header.TxHash.Hex(),
		"stateRoot":        header.Root.Hex(),
		"receiptsRoot":     header.ReceiptHash.Hex(),
		"miner":            header.Coinbase.Hex(),
		"difficulty":       hexutil.EncodeBig(header.Difficulty),
		"extraData":        hexutil.Encode(header.Extra),
		"gasLimit":         hexutil.EncodeUint64(header.GasLimit),
		"gasUsed":          hexutil.EncodeUint64(header.GasUsed),
		"timestamp":        hexutil.EncodeUint64(header.Time),
		"mixHash":          header.MixDigest.Hex(),
	}

	if header.BaseFee != nil {
		result["baseFeePerGas"] = hexutil.EncodeBig(header.BaseFee)
	}

	if header.WithdrawalsHash != nil {
		result["withdrawalsRoot"] = header.WithdrawalsHash.Hex()
	}

	if header.BlobGasUsed != nil {
		result["blobGasUsed"] = hexutil.EncodeUint64(*header.BlobGasUsed)
	}

	if header.ExcessBlobGas != nil {
		result["excessBlobGas"] = hexutil.EncodeUint64(*header.ExcessBlobGas)
	}

	if header.ParentBeaconRoot != nil {
		result["parentBeaconBlockRoot"] = header.ParentBeaconRoot.Hex()
	}

	// Add total difficulty (default to 0 if not available)
	if cache.TD != nil {
		result["totalDifficulty"] = hexutil.EncodeBig(cache.TD)
	} else {
		result["totalDifficulty"] = "0x0"
	}

	// Add transactions
	if cache.Body != nil && cache.Body.Transactions.Len() > 0 {
		txs, _ := cache.Body.Transactions.Items()
		if fullTx {
			txResults := make([]map[string]any, len(txs))
			for i, tx := range txs {
				txResults[i] = formatTransactionResponse(tx, hash, header, uint64(i))
			}
			result["transactions"] = txResults
		} else {
			txHashes := make([]string, len(txs))
			for i, tx := range txs {
				txHashes[i] = tx.Hash().Hex()
			}
			result["transactions"] = txHashes
		}
	} else {
		result["transactions"] = []string{}
	}

	// Add uncles
	if cache.Body != nil && cache.Body.Uncles.Len() > 0 {
		uncles, _ := cache.Body.Uncles.Items()
		uncleHashes := make([]string, len(uncles))
		for i, uncle := range uncles {
			uncleHashes[i] = uncle.Hash().Hex()
		}
		result["uncles"] = uncleHashes
	} else {
		result["uncles"] = []string{}
	}

	// Add size (approximate based on header + body)
	result["size"] = hexutil.EncodeUint64(0) // We don't have exact size; use 0

	return result
}

// formatTransactionResponse formats a transaction into the Ethereum JSON-RPC format.
// If blockHash is empty, the transaction is considered pending.
func formatTransactionResponse(tx *types.Transaction, blockHash common.Hash, header *types.Header, index uint64) map[string]any {
	v, r, s := tx.RawSignatureValues()

	result := map[string]any{
		"hash":  tx.Hash().Hex(),
		"nonce": hexutil.EncodeUint64(tx.Nonce()),
		"gas":   hexutil.EncodeUint64(tx.Gas()),
		"value": hexutil.EncodeBig(tx.Value()),
		"input": hexutil.Encode(tx.Data()),
		"v":     hexutil.EncodeBig(v),
		"r":     hexutil.EncodeBig(r),
		"s":     hexutil.EncodeBig(s),
		"type":  hexutil.EncodeUint64(uint64(tx.Type())),
	}

	if tx.To() != nil {
		result["to"] = tx.To().Hex()
	} else {
		result["to"] = nil
	}

	// Add from address if we can derive it
	signer := types.LatestSignerForChainID(tx.ChainId())
	if from, err := types.Sender(signer, tx); err == nil {
		result["from"] = from.Hex()
	}

	// Set gas price fields based on transaction type
	switch tx.Type() {
	case types.LegacyTxType, types.AccessListTxType:
		result["gasPrice"] = hexutil.EncodeBig(tx.GasPrice())
	case types.DynamicFeeTxType, types.BlobTxType:
		result["maxFeePerGas"] = hexutil.EncodeBig(tx.GasFeeCap())
		result["maxPriorityFeePerGas"] = hexutil.EncodeBig(tx.GasTipCap())
		// For EIP-1559 txs, also set gasPrice to effective gas price if in a block
		if header != nil && header.BaseFee != nil {
			effectiveGasPrice := new(big.Int).Add(header.BaseFee, tx.GasTipCap())
			if effectiveGasPrice.Cmp(tx.GasFeeCap()) > 0 {
				effectiveGasPrice = tx.GasFeeCap()
			}
			result["gasPrice"] = hexutil.EncodeBig(effectiveGasPrice)
		} else {
			result["gasPrice"] = hexutil.EncodeBig(tx.GasFeeCap())
		}
	}

	// Add chain ID if present
	if tx.ChainId() != nil {
		result["chainId"] = hexutil.EncodeBig(tx.ChainId())
	}

	// Add yParity for typed transactions (EIP-2930+)
	if tx.Type() != types.LegacyTxType {
		result["yParity"] = hexutil.EncodeBig(v)
	}

	// Add access list if present
	if tx.AccessList() != nil {
		result["accessList"] = tx.AccessList()
	}

	// Add blob-specific fields
	if tx.Type() == types.BlobTxType {
		result["maxFeePerBlobGas"] = hexutil.EncodeBig(tx.BlobGasFeeCap())
		result["blobVersionedHashes"] = tx.BlobHashes()
	}

	// Add block info if transaction is in a block
	if blockHash != (common.Hash{}) && header != nil {
		result["blockHash"] = blockHash.Hex()
		result["blockNumber"] = hexutil.EncodeUint64(header.Number.Uint64())
		result["transactionIndex"] = hexutil.EncodeUint64(index)
	} else {
		result["blockHash"] = nil
		result["blockNumber"] = nil
		result["transactionIndex"] = nil
	}

	return result
}
