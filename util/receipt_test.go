package util

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
)

// receiptRPC is a fake JSON-RPC endpoint for eth_getTransactionReceipt that
// answers null for the first pendingPolls calls and the canned receipt after.
type receiptRPC struct {
	t *testing.T

	pendingPolls int
	receipt      map[string]any

	mu    sync.Mutex
	calls int
}

func (s *receiptRPC) callCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.calls
}

func (s *receiptRPC) handle(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID     json.RawMessage `json:"id"`
		Method string          `json:"method"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Method != "eth_getTransactionReceipt" {
		http.Error(w, "unexpected method "+req.Method, http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	s.calls++
	var result any
	if s.calls > s.pendingPolls {
		result = s.receipt
	}
	s.mu.Unlock()

	payload, err := json.Marshal(map[string]any{
		"jsonrpc": "2.0", "id": json.RawMessage(req.ID), "result": result,
	})
	if err != nil {
		s.t.Errorf("failed to marshal rpc response: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, werr := w.Write(payload); werr != nil {
		s.t.Errorf("failed to write rpc response: %v", werr)
	}
}

func newReceiptRPC(t *testing.T, s *receiptRPC) *ethrpc.Client {
	t.Helper()
	s.t = t
	server := httptest.NewServer(http.HandlerFunc(s.handle))
	t.Cleanup(server.Close)

	client, err := ethrpc.DialContext(t.Context(), server.URL)
	if err != nil {
		t.Fatalf("failed to dial fake rpc: %v", err)
	}
	t.Cleanup(client.Close)
	return client
}

func TestWaitReceiptRawFixedInterval(t *testing.T) {
	txHash := common.HexToHash("0xf00d")
	rpc := &receiptRPC{
		pendingPolls: 3,
		receipt: map[string]any{
			"transactionHash": txHash.Hex(),
			"status":          "0x1",
			// A field types.Receipt doesn't require, proving raw passthrough.
			"effectiveGasPrice": "0x77359400",
		},
	}
	client := newReceiptRPC(t, rpc)

	raw, err := WaitReceiptRaw(t.Context(), client, txHash, ReceiptWaitOpts{
		Interval: 10 * time.Millisecond,
		// MaxAttempts must be ignored in interval mode: fewer than the polls
		// needed to see the receipt.
		MaxAttempts: 2,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := rpc.callCount(); got != 4 {
		t.Errorf("poll count = %d, want 4 (3 pending + 1 found)", got)
	}
	if !strings.Contains(string(raw), `"effectiveGasPrice":"0x77359400"`) {
		t.Errorf("raw receipt not passed through verbatim: %s", raw)
	}
}

func TestWaitReceiptRawFixedIntervalTimesOut(t *testing.T) {
	txHash := common.HexToHash("0xf00d")
	rpc := &receiptRPC{pendingPolls: 1 << 30} // never found
	client := newReceiptRPC(t, rpc)

	start := time.Now()
	_, err := WaitReceiptRaw(t.Context(), client, txHash, ReceiptWaitOpts{
		Interval: 20 * time.Millisecond,
		Timeout:  150 * time.Millisecond,
	})
	if err == nil {
		t.Fatal("expected a timeout error")
	}
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Errorf("wait took %s, expected to stop at the ~150ms timeout", elapsed)
	}
	if got := rpc.callCount(); got < 3 {
		t.Errorf("poll count = %d, want steady polling before the timeout", got)
	}
}

func TestWaitReceiptRawBackoffExhaustsAttempts(t *testing.T) {
	txHash := common.HexToHash("0xf00d")
	rpc := &receiptRPC{pendingPolls: 1 << 30} // never found
	client := newReceiptRPC(t, rpc)

	_, err := WaitReceiptRaw(t.Context(), client, txHash, ReceiptWaitOpts{
		MaxAttempts:  3,
		InitialDelay: 10 * time.Millisecond,
	})
	if err == nil || !strings.Contains(err.Error(), "failed to get receipt after 3 attempts") {
		t.Fatalf("err = %v, want attempts-exhausted error", err)
	}
	if got := rpc.callCount(); got != 3 {
		t.Errorf("poll count = %d, want exactly 3", got)
	}
}

func TestWaitReceiptRawBackoffFindsReceipt(t *testing.T) {
	txHash := common.HexToHash("0xf00d")
	rpc := &receiptRPC{
		pendingPolls: 2,
		receipt:      map[string]any{"transactionHash": txHash.Hex(), "status": "0x1"},
	}
	client := newReceiptRPC(t, rpc)

	raw, err := WaitReceiptRaw(t.Context(), client, txHash, ReceiptWaitOpts{
		MaxAttempts:  10,
		InitialDelay: 10 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(raw), txHash.Hex()) {
		t.Errorf("raw receipt missing tx hash: %s", raw)
	}
}

func TestWaitReceiptRawContextCancellation(t *testing.T) {
	for _, tt := range []struct {
		name string
		opts ReceiptWaitOpts
	}{
		{"fixed interval", ReceiptWaitOpts{Interval: 50 * time.Millisecond}},
		{"backoff", ReceiptWaitOpts{InitialDelay: 50 * time.Millisecond}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			txHash := common.HexToHash("0xf00d")
			rpc := &receiptRPC{pendingPolls: 1 << 30} // never found
			client := newReceiptRPC(t, rpc)

			ctx, cancel := context.WithCancel(t.Context())
			done := make(chan error, 1)
			go func() {
				_, err := WaitReceiptRaw(ctx, client, txHash, tt.opts)
				done <- err
			}()

			time.Sleep(20 * time.Millisecond)
			cancel()

			select {
			case err := <-done:
				if err == nil {
					t.Error("expected a cancellation error")
				}
			case <-time.After(2 * time.Second):
				t.Fatal("WaitReceiptRaw did not stop promptly on cancellation")
			}
		})
	}
}
