package sensor

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"

	"github.com/0xPolygon/polygon-cli/p2p"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
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

// handleRPC sets up the JSON-RPC server for receiving and broadcasting transactions.
// It handles eth_sendRawTransaction requests, validates transaction signatures,
// and broadcasts valid transactions to all connected peers.
// Supports both single requests and batch requests per JSON-RPC 2.0 specification.
func handleRPC(conns *p2p.Conns, networkID uint64) {
	// Use network ID as chain ID for signature validation
	chainID := new(big.Int).SetUint64(networkID)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Read request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, -32700, "Parse error", nil)
			return
		}
		defer r.Body.Close()

		// Check if this is a batch request (starts with '[') or single request
		trimmed := strings.TrimSpace(string(body))
		if len(trimmed) > 0 && trimmed[0] == '[' {
			// Handle batch request
			handleBatchRequest(w, body, conns, chainID)
			return
		}

		// Parse single JSON-RPC request
		var req rpcRequest
		if err := json.Unmarshal(body, &req); err != nil {
			writeError(w, -32700, "Parse error", nil)
			return
		}

		// Handle eth_sendRawTransaction
		if req.Method == "eth_sendRawTransaction" {
			handleSendRawTransaction(w, req, conns, chainID)
			return
		}

		// Method not found
		writeError(w, -32601, "Method not found", req.ID)
	})

	addr := fmt.Sprintf(":%d", inputSensorParams.RPCPort)
	log.Info().Str("addr", addr).Msg("Starting JSON-RPC server")
	if err := http.ListenAndServe(addr, nil); err != nil {
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
	json.NewEncoder(w).Encode(response)
}

// writeResult writes a JSON-RPC 2.0 success response with the specified result and request ID.
func writeResult(w http.ResponseWriter, result any, id any) {
	w.Header().Set("Content-Type", "application/json")
	response := rpcResponse{
		JSONRPC: "2.0",
		Result:  result,
		ID:      id,
	}
	json.NewEncoder(w).Encode(response)
}

// handleBatchRequest processes JSON-RPC 2.0 batch requests, validates all transactions,
// and broadcasts valid transactions to connected peers. Returns a batch response with
// results or errors for each request in the batch.
func handleBatchRequest(w http.ResponseWriter, body []byte, conns *p2p.Conns, chainID *big.Int) {
	// Parse batch of requests
	var requests []rpcRequest
	if err := json.Unmarshal(body, &requests); err != nil {
		writeError(w, -32700, "Parse error", nil)
		return
	}

	// Validate batch is not empty
	if len(requests) == 0 {
		writeError(w, -32600, "Invalid request: empty batch", nil)
		return
	}

	// Process all requests and collect valid transactions for batch broadcasting
	responses := make([]rpcResponse, 0, len(requests))
	txs := make(types.Transactions, 0)

	for _, req := range requests {
		if req.Method != "eth_sendRawTransaction" {
			response := processSingleRequest(req, conns, chainID)
			responses = append(responses, response)
			continue
		}

		tx, response := validateTransaction(req, chainID)
		if tx == nil {
			responses = append(responses, response)
			continue
		}

		txs = append(txs, tx)
		responses = append(responses, rpcResponse{
			JSONRPC: "2.0",
			Result:  tx.Hash().Hex(),
			ID:      req.ID,
		})
	}

	// Broadcast all valid transactions in a single batch if there are any
	if len(txs) > 0 {
		log.Info().
			Int("txs", len(txs)).
			Int("requests", len(requests)).
			Msg("Broadcasting batch of transactions")

		count := conns.BroadcastTxs(txs)

		log.Info().
			Int("txCount", len(txs)).
			Int("peers", count).
			Msg("Batch broadcast complete")
	}

	// Write batch response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responses); err != nil {
		log.Error().Err(err).Msg("Failed to encode batch response")
	}
}

// validateTransaction validates a transaction from a JSON-RPC request by decoding the raw
// transaction hex, unmarshaling it, and verifying the signature. Returns the transaction
// if valid, or an error response if validation fails.
func validateTransaction(req rpcRequest, chainID *big.Int) (*types.Transaction, rpcResponse) {
	// Check params
	if len(req.Params) == 0 {
		return nil, rpcResponse{
			JSONRPC: "2.0",
			Error: &rpcError{
				Code:    -32602,
				Message: "Invalid params: missing raw transaction",
			},
			ID: req.ID,
		}
	}

	// Extract raw transaction hex string
	hex, ok := req.Params[0].(string)
	if !ok {
		return nil, rpcResponse{
			JSONRPC: "2.0",
			Error: &rpcError{
				Code:    -32602,
				Message: "Invalid params: raw transaction must be a hex string",
			},
			ID: req.ID,
		}
	}

	// Decode hex string to bytes
	bytes, err := hexutil.Decode(hex)
	if err != nil {
		return nil, rpcResponse{
			JSONRPC: "2.0",
			Error: &rpcError{
				Code:    -32602,
				Message: fmt.Sprintf("Invalid transaction hex: %v", err),
			},
			ID: req.ID,
		}
	}

	// Unmarshal transaction
	tx := new(types.Transaction)
	if err := tx.UnmarshalBinary(bytes); err != nil {
		return nil, rpcResponse{
			JSONRPC: "2.0",
			Error: &rpcError{
				Code:    -32602,
				Message: fmt.Sprintf("Invalid transaction encoding: %v", err),
			},
			ID: req.ID,
		}
	}

	// Validate transaction signature
	signer := types.LatestSignerForChainID(chainID)
	sender, err := types.Sender(signer, tx)
	if err != nil {
		return nil, rpcResponse{
			JSONRPC: "2.0",
			Error: &rpcError{
				Code:    -32602,
				Message: fmt.Sprintf("Invalid transaction signature: %v", err),
			},
			ID: req.ID,
		}
	}

	// Log the transaction
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

// processSingleRequest processes a single JSON-RPC request and returns a response.
// Currently supports eth_sendRawTransaction method.
func processSingleRequest(req rpcRequest, conns *p2p.Conns, chainID *big.Int) rpcResponse {
	// Handle eth_sendRawTransaction
	if req.Method == "eth_sendRawTransaction" {
		return processSendRawTransaction(req, conns, chainID)
	}

	// Method not found
	return rpcResponse{
		JSONRPC: "2.0",
		Error: &rpcError{
			Code:    -32601,
			Message: "Method not found",
		},
		ID: req.ID,
	}
}

// processSendRawTransaction processes a single eth_sendRawTransaction request, validates
// the transaction, broadcasts it to all connected peers, and returns the transaction hash.
func processSendRawTransaction(req rpcRequest, conns *p2p.Conns, chainID *big.Int) rpcResponse {
	// Check params
	if len(req.Params) == 0 {
		return rpcResponse{
			JSONRPC: "2.0",
			Error: &rpcError{
				Code:    -32602,
				Message: "Invalid params: missing raw transaction",
			},
			ID: req.ID,
		}
	}

	// Extract raw transaction hex string
	rawTxHex, ok := req.Params[0].(string)
	if !ok {
		return rpcResponse{
			JSONRPC: "2.0",
			Error: &rpcError{
				Code:    -32602,
				Message: "Invalid params: raw transaction must be a hex string",
			},
			ID: req.ID,
		}
	}

	// Decode hex string to bytes
	txBytes, err := hexutil.Decode(rawTxHex)
	if err != nil {
		return rpcResponse{
			JSONRPC: "2.0",
			Error: &rpcError{
				Code:    -32602,
				Message: fmt.Sprintf("Invalid transaction hex: %v", err),
			},
			ID: req.ID,
		}
	}

	// Unmarshal transaction
	tx := new(types.Transaction)
	if err := tx.UnmarshalBinary(txBytes); err != nil {
		return rpcResponse{
			JSONRPC: "2.0",
			Error: &rpcError{
				Code:    -32602,
				Message: fmt.Sprintf("Invalid transaction encoding: %v", err),
			},
			ID: req.ID,
		}
	}

	// Validate transaction signature
	signer := types.LatestSignerForChainID(chainID)
	sender, err := types.Sender(signer, tx)
	if err != nil {
		return rpcResponse{
			JSONRPC: "2.0",
			Error: &rpcError{
				Code:    -32602,
				Message: fmt.Sprintf("Invalid transaction signature: %v", err),
			},
			ID: req.ID,
		}
	}

	// Log the transaction
	toAddr := "nil"
	if tx.To() != nil {
		toAddr = tx.To().Hex()
	}

	log.Info().
		Str("hash", tx.Hash().Hex()).
		Str("from", sender.Hex()).
		Str("to", toAddr).
		Str("value", tx.Value().String()).
		Uint64("gas", tx.Gas()).
		Msg("Broadcasting transaction")

	// Broadcast to all peers
	broadcastCount := conns.BroadcastTx(tx)

	log.Info().
		Str("hash", tx.Hash().Hex()).
		Int("peers", broadcastCount).
		Msg("Transaction broadcast complete")

	// Return transaction hash
	return rpcResponse{
		JSONRPC: "2.0",
		Result:  tx.Hash().Hex(),
		ID:      req.ID,
	}
}

// handleSendRawTransaction processes eth_sendRawTransaction requests, validates the
// transaction, broadcasts it to all connected peers, and writes the transaction hash
// as a JSON-RPC response.
func handleSendRawTransaction(w http.ResponseWriter, req rpcRequest, conns *p2p.Conns, chainID *big.Int) {
	// Check params
	if len(req.Params) == 0 {
		writeError(w, -32602, "Invalid params: missing raw transaction", req.ID)
		return
	}

	// Extract raw transaction hex string
	rawTxHex, ok := req.Params[0].(string)
	if !ok {
		writeError(w, -32602, "Invalid params: raw transaction must be a hex string", req.ID)
		return
	}

	// Decode hex string to bytes
	txBytes, err := hexutil.Decode(rawTxHex)
	if err != nil {
		writeError(w, -32602, fmt.Sprintf("Invalid transaction hex: %v", err), req.ID)
		return
	}

	// Unmarshal transaction
	tx := new(types.Transaction)
	if err := tx.UnmarshalBinary(txBytes); err != nil {
		writeError(w, -32602, fmt.Sprintf("Invalid transaction encoding: %v", err), req.ID)
		return
	}

	// Validate transaction signature
	signer := types.LatestSignerForChainID(chainID)
	sender, err := types.Sender(signer, tx)
	if err != nil {
		writeError(w, -32602, fmt.Sprintf("Invalid transaction signature: %v", err), req.ID)
		return
	}

	to := "nil"
	if tx.To() != nil {
		to = tx.To().Hex()
	}

	log.Info().
		Str("hash", tx.Hash().Hex()).
		Str("from", sender.Hex()).
		Str("to", to).
		Str("value", tx.Value().String()).
		Uint64("gas", tx.Gas()).
		Msg("Broadcasting transaction")

	// Broadcast to all peers
	count := conns.BroadcastTx(tx)

	log.Info().
		Str("hash", tx.Hash().Hex()).
		Int("peers", count).
		Msg("Transaction broadcast complete")

	// Return transaction hash
	writeResult(w, tx.Hash().Hex(), req.ID)
}
