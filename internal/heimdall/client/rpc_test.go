package client

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestRPCClientCallSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		body, _ := io.ReadAll(r.Body)
		var req RPCRequest
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("unmarshal req: %v", err)
		}
		if req.Method != "status" {
			t.Errorf("method = %q, want status", req.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"jsonrpc":"2.0","id":1,"result":{"latest_block_height":"42"}}`)
	}))
	defer srv.Close()

	c := NewRPCClient(srv.URL, 5*time.Second, nil, false)
	res, err := c.Call(context.Background(), "status", nil)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !strings.Contains(string(res), `latest_block_height`) {
		t.Errorf("result = %q", res)
	}
}

func TestRPCClientErrorEnvelope(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"jsonrpc":"2.0","id":1,"error":{"code":-32602,"message":"invalid params","data":"height too large"}}`)
	}))
	defer srv.Close()

	c := NewRPCClient(srv.URL, 5*time.Second, nil, false)
	_, err := c.Call(context.Background(), "block", map[string]any{"height": "999999999999"})
	if err == nil {
		t.Fatal("expected rpc error, got nil")
	}
	var rpcErr *RPCError
	if !errors.As(err, &rpcErr) {
		t.Fatalf("error type = %T, want *RPCError", err)
	}
	if rpcErr.Code != -32602 {
		t.Errorf("Code = %d, want -32602", rpcErr.Code)
	}
}

func TestRPCClient5xxIsHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = io.WriteString(w, "upstream lost")
	}))
	defer srv.Close()

	c := NewRPCClient(srv.URL, 5*time.Second, nil, false)
	_, err := c.Call(context.Background(), "health", nil)
	if err == nil {
		t.Fatal("expected 5xx error, got nil")
	}
	var hErr *HTTPError
	if !errors.As(err, &hErr) {
		t.Fatalf("error type = %T, want *HTTPError", err)
	}
	if hErr.StatusCode != 502 {
		t.Errorf("StatusCode = %d, want 502", hErr.StatusCode)
	}
}

func TestRPCClientNonJSONBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "not json at all")
	}))
	defer srv.Close()

	c := NewRPCClient(srv.URL, 5*time.Second, nil, false)
	_, err := c.Call(context.Background(), "status", nil)
	if err == nil {
		t.Fatal("expected decode error, got nil")
	}
	if !strings.Contains(err.Error(), "decoding") {
		t.Errorf("error = %v, want decode hint", err)
	}
}
