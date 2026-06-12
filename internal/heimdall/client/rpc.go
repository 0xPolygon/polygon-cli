package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
)

// RPCRequest is the JSON-RPC 2.0 envelope used by CometBFT.
type RPCRequest struct {
	JSONRPC string         `json:"jsonrpc"`
	ID      uint64         `json:"id"`
	Method  string         `json:"method"`
	Params  map[string]any `json:"params,omitempty"`
}

// RPCError represents a JSON-RPC error payload.
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

// Error implements error.
func (e *RPCError) Error() string {
	if e.Data != "" {
		return fmt.Sprintf("rpc error %d: %s (%s)", e.Code, e.Message, e.Data)
	}
	return fmt.Sprintf("rpc error %d: %s", e.Code, e.Message)
}

// RPCResponse is the JSON-RPC 2.0 response envelope.
type RPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      uint64          `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *RPCError       `json:"error,omitempty"`
}

// RPCClient is a CometBFT JSON-RPC client.
type RPCClient struct {
	BaseURL   string
	Headers   map[string]string
	Transport Transport

	// Monotonic request ID counter; incremented per call. Atomic so
	// the same client can be shared across goroutines.
	nextID atomic.Uint64
}

// NewRPCClient returns an RPCClient configured from the resolved
// config.
func NewRPCClient(base string, timeout time.Duration, headers map[string]string, insecure bool) *RPCClient {
	tlsCfg := &tls.Config{}
	if insecure {
		tlsCfg.InsecureSkipVerify = true
	}
	httpClient := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: tlsCfg,
		},
	}
	return &RPCClient{
		BaseURL:   strings.TrimRight(base, "/"),
		Headers:   cloneHeaders(headers),
		Transport: &HTTPTransport{Client: httpClient},
	}
}

// Call issues a JSON-RPC call with the given method and params and
// returns the raw `result` field. Returns *RPCError on a JSON-RPC
// error, *HTTPError on HTTP 4xx/5xx, or *NetworkError on transport
// failures.
func (c *RPCClient) Call(ctx context.Context, method string, params map[string]any) (json.RawMessage, error) {
	id := c.nextID.Add(1)
	envelope := &RPCRequest{
		JSONRPC: "2.0",
		ID:      id,
		Method:  method,
		Params:  params,
	}
	body, err := json.Marshal(envelope)
	if err != nil {
		return nil, fmt.Errorf("marshal %s request: %w", method, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("building %s request: %w", method, err)
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}

	respBody, status, err := c.Transport.Do(req)
	if err != nil {
		return nil, err
	}
	if status == 0 && respBody == nil {
		// CurlTransport short-circuit.
		return nil, nil
	}
	if status >= 400 {
		return respBody, &HTTPError{Method: http.MethodPost, URL: c.BaseURL, StatusCode: status, Body: respBody}
	}

	var resp RPCResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("decoding %s response: %w (body=%q)", method, err, truncate(respBody, 256))
	}
	if resp.Error != nil {
		return resp.Result, resp.Error
	}
	return resp.Result, nil
}

func truncate(b []byte, n int) string {
	if len(b) <= n {
		return string(b)
	}
	return string(b[:n]) + "..."
}
