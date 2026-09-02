package mode

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	ethrpc "github.com/ethereum/go-ethereum/rpc"

	"github.com/0xPolygon/polygon-cli/loadtest/config"
)

// syncRPC is a fake JSON-RPC endpoint standing in for a node implementing
// eth_sendRawTransactionSync. It records the params it was called with so the
// wire format can be asserted, and replies with a canned result or error.
type syncRPC struct {
	t *testing.T

	// result is marshalled as the JSON-RPC result when err is nil.
	result any
	// errCode/errMsg/errData produce a JSON-RPC error response instead.
	errCode int
	errMsg  string
	errData any

	mu     sync.Mutex
	params []any
	calls  int
}

func newSyncRPC(t *testing.T, s *syncRPC) *ethrpc.Client {
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

func (s *syncRPC) handle(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID     json.RawMessage `json:"id"`
		Method string          `json:"method"`
		Params []any           `json:"params"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Method != "eth_sendRawTransactionSync" {
		http.Error(w, "unexpected method "+req.Method, http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	s.params = req.Params
	s.calls++
	s.mu.Unlock()

	var payload []byte
	var err error
	if s.errCode != 0 {
		body := map[string]any{"code": s.errCode, "message": s.errMsg}
		if s.errData != nil {
			body["data"] = s.errData
		}
		payload, err = json.Marshal(map[string]any{
			"jsonrpc": "2.0", "id": json.RawMessage(req.ID), "error": body,
		})
	} else {
		payload, err = json.Marshal(map[string]any{
			"jsonrpc": "2.0", "id": json.RawMessage(req.ID), "result": s.result,
		})
	}
	if err != nil {
		s.t.Errorf("failed to marshal rpc response: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, werr := w.Write(payload); werr != nil {
		s.t.Errorf("failed to write rpc response: %v", werr)
	}
}

func (s *syncRPC) lastParams() []any {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.params
}

// testTx returns a signed transaction to submit.
func testTx(t *testing.T) *types.Transaction {
	t.Helper()
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}
	signer := types.LatestSignerForChainID(big.NewInt(1337))
	tx, err := types.SignNewTx(key, signer, &types.DynamicFeeTx{
		ChainID:   big.NewInt(1337),
		Nonce:     7,
		GasTipCap: big.NewInt(1e9),
		GasFeeCap: big.NewInt(2e9),
		Gas:       21000,
		To:        &common.Address{},
	})
	if err != nil {
		t.Fatalf("failed to sign tx: %v", err)
	}
	return tx
}

// canonicalReceipt is a receipt with the block fields filled in.
func canonicalReceipt(hash common.Hash) map[string]any {
	return map[string]any{
		"transactionHash": hash.Hex(),
		"status":          "0x1",
		"gasUsed":         "0x5208",
		"blockNumber":     "0x2a",
		"blockHash":       common.HexToHash("0xbeef").Hex(),
	}
}

// speculativeReceipt is a receipt a preconfirming node can return before the
// transaction lands in a canonical block.
func speculativeReceipt(hash common.Hash) map[string]any {
	return map[string]any{
		"transactionHash": hash.Hex(),
		"status":          "0x1",
		"gasUsed":         "0x5208",
		"blockNumber":     nil,
		"blockHash":       nil,
	}
}

func send(t *testing.T, rpc *syncRPC, cfg *config.Config, tx *types.Transaction) (*SyncTracker, error) {
	t.Helper()
	client := newSyncRPC(t, rpc)
	tracker := NewSyncTracker()
	deps := &Dependencies{SendRPCClient: client, SyncTracker: tracker}
	return tracker, SendRawTransactionSync(t.Context(), deps, cfg, tx)
}

func TestSendRawTransactionSyncOmitsUnsetTimeout(t *testing.T) {
	tx := testTx(t)
	rpc := &syncRPC{result: canonicalReceipt(tx.Hash())}

	if _, err := send(t, rpc, &config.Config{}, tx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	params := rpc.lastParams()
	if len(params) != 1 {
		t.Fatalf("params = %v, want just the raw transaction", params)
	}
	raw, ok := params[0].(string)
	if !ok || len(raw) < 4 || raw[:2] != "0x" {
		t.Errorf("params[0] = %v, want a 0x-prefixed raw transaction", params[0])
	}
}

func TestSendRawTransactionSyncSendsTimeoutAsHexByDefault(t *testing.T) {
	// bor takes this parameter as *hexutil.Uint64, which unmarshals only from a
	// quoted hex quantity and rejects a bare JSON number, so pin the encoding.
	tx := testTx(t)
	rpc := &syncRPC{result: canonicalReceipt(tx.Hash())}

	if _, err := send(t, rpc, &config.Config{SyncTxTimeout: 5 * time.Second}, tx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	params := rpc.lastParams()
	if len(params) != 2 {
		t.Fatalf("params = %v, want raw transaction and timeout", params)
	}
	got, ok := params[1].(string)
	if !ok {
		t.Fatalf("params[1] = %#v (%T), want a hex quantity string", params[1], params[1])
	}
	if got != "0x1388" {
		t.Errorf("timeout = %s, want 0x1388 (5000ms)", got)
	}
}

func TestSendRawTransactionSyncSendsTimeoutAsIntegerWhenAsked(t *testing.T) {
	// EIP-7966 describes the parameter as an integer, which spec-literal
	// servers require and which hexutil.Uint64 rejects. --sync-tx-timeout-int
	// selects it.
	tx := testTx(t)
	rpc := &syncRPC{result: canonicalReceipt(tx.Hash())}

	cfg := &config.Config{SyncTxTimeout: 5 * time.Second, SyncTxTimeoutInt: true}
	if _, err := send(t, rpc, cfg, tx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	params := rpc.lastParams()
	if len(params) != 2 {
		t.Fatalf("params = %v, want raw transaction and timeout", params)
	}
	// encoding/json decodes JSON numbers into float64; a hex string would not
	// satisfy this assertion.
	ms, ok := params[1].(float64)
	if !ok {
		t.Fatalf("params[1] = %#v (%T), want a JSON number of milliseconds", params[1], params[1])
	}
	if ms != 5000 {
		t.Errorf("timeout = %v ms, want 5000", ms)
	}
}

func TestSendRawTransactionSyncBorPreconfirmationReceipt(t *testing.T) {
	// Bor's preconfirmation receipt: no block hash, plus an explicit marker.
	tx := testTx(t)
	receipt := map[string]any{
		"transactionHash": tx.Hash().Hex(),
		"status":          "0x1",
		"gasUsed":         "0x5208",
		"blockNumber":     "0x2a",
		"blockHash":       nil,
		"preconfirmation": true,
	}
	rpc := &syncRPC{result: receipt}

	tracker, err := send(t, rpc, &config.Config{}, tx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := tracker.speculative.Load(); got != 1 {
		t.Errorf("speculative count = %d, want 1", got)
	}
	if got := tracker.canonical.Load(); got != 0 {
		t.Errorf("canonical count = %d, want 0", got)
	}
}

func TestSendRawTransactionSyncMarkerBeatsBlockFields(t *testing.T) {
	// A node that fills in block fields but says preconfirmation=false is
	// canonical; the marker is authoritative when present.
	tx := testTx(t)
	receipt := canonicalReceipt(tx.Hash())
	receipt["preconfirmation"] = false
	rpc := &syncRPC{result: receipt}

	tracker, err := send(t, rpc, &config.Config{}, tx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := tracker.canonical.Load(); got != 1 {
		t.Errorf("canonical count = %d, want 1", got)
	}
}

func TestSendRawTransactionSyncBorErrorCodes(t *testing.T) {
	// Bor implements only code 4 from EIP-7966; anything it refuses to accept
	// fails in SubmitTransaction and returns bor's own codes.
	tests := []struct {
		name    string
		code    int
		counter func(*SyncTracker) uint64
	}{
		{"nonce too high is a gap", SyncErrBorNonceTooHigh, func(t *SyncTracker) uint64 { return t.nonceGaps.Load() }},
		{"nonce too low is a rejection", SyncErrBorNonceTooLow, func(t *SyncTracker) uint64 { return t.rejected.Load() }},
		{"intrinsic gas", SyncErrBorIntrinsicGas, func(t *SyncTracker) uint64 { return t.rejected.Load() }},
		{"insufficient funds", SyncErrBorInsufficientFunds, func(t *SyncTracker) uint64 { return t.rejected.Load() }},
		{"client limit", SyncErrBorClientLimit, func(t *SyncTracker) uint64 { return t.rejected.Load() }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := testTx(t)
			rpc := &syncRPC{errCode: tt.code, errMsg: tt.name}

			tracker, err := send(t, rpc, &config.Config{}, tx)
			if err == nil {
				t.Fatal("expected the rpc error to be returned")
			}
			if got := tt.counter(tracker); got != 1 {
				t.Errorf("counter = %d, want 1", got)
			}
			if got := tracker.otherErrs.Load(); got != 0 {
				t.Errorf("other error count = %d, want 0 (should be classified)", got)
			}
		})
	}
}

func TestSendRawTransactionSyncCanonicalReceipt(t *testing.T) {
	tx := testTx(t)
	rpc := &syncRPC{result: canonicalReceipt(tx.Hash())}

	tracker, err := send(t, rpc, &config.Config{}, tx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	results := tracker.Results()
	if len(results) != 1 {
		t.Fatalf("results = %d, want 1", len(results))
	}
	r := results[0]
	if r.Speculative {
		t.Error("receipt with a block should not be classified speculative")
	}
	if r.BlockNumber != 42 {
		t.Errorf("block number = %d, want 42", r.BlockNumber)
	}
	if r.GasUsed != 21000 {
		t.Errorf("gas used = %d, want 21000", r.GasUsed)
	}
	if r.Status != types.ReceiptStatusSuccessful {
		t.Errorf("status = %d, want 1", r.Status)
	}
	if got := tracker.canonical.Load(); got != 1 {
		t.Errorf("canonical count = %d, want 1", got)
	}
	if got := tracker.speculative.Load(); got != 0 {
		t.Errorf("speculative count = %d, want 0", got)
	}
}

func TestSendRawTransactionSyncSpeculativeReceipt(t *testing.T) {
	tx := testTx(t)
	rpc := &syncRPC{result: speculativeReceipt(tx.Hash())}

	tracker, err := send(t, rpc, &config.Config{}, tx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	results := tracker.Results()
	if len(results) != 1 || !results[0].Speculative {
		t.Fatalf("results = %+v, want one speculative receipt", results)
	}
	if got := tracker.speculative.Load(); got != 1 {
		t.Errorf("speculative count = %d, want 1", got)
	}
	if got := tracker.receipts.Load(); got != 1 {
		t.Errorf("receipt count = %d, want 1", got)
	}
}

func TestSendRawTransactionSyncErrorCodes(t *testing.T) {
	tests := []struct {
		name    string
		code    int
		data    any
		counter func(*SyncTracker) uint64
	}{
		{"timeout", SyncErrTimeout, "0xabc", func(t *SyncTracker) uint64 { return t.timeouts.Load() }},
		{"queued", SyncErrQueued, "0xabc", func(t *SyncTracker) uint64 { return t.queued.Load() }},
		{"nonce gap", SyncErrNonceGap, "0x8", func(t *SyncTracker) uint64 { return t.nonceGaps.Load() }},
		{"unspecified code", -32000, nil, func(t *SyncTracker) uint64 { return t.otherErrs.Load() }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := testTx(t)
			rpc := &syncRPC{errCode: tt.code, errMsg: tt.name, errData: tt.data}

			tracker, err := send(t, rpc, &config.Config{}, tx)
			if err == nil {
				t.Fatal("expected the rpc error to be returned")
			}
			if code, ok := SyncErrorCode(err); !ok || code != tt.code {
				t.Errorf("error code = %d (present %v), want %d", code, ok, tt.code)
			}
			if got := tt.counter(tracker); got != 1 {
				t.Errorf("counter = %d, want 1", got)
			}
			if got := tracker.receipts.Load(); got != 0 {
				t.Errorf("receipt count = %d, want 0 on error", got)
			}
			results := tracker.Results()
			if len(results) != 1 || results[0].ErrorCode != tt.code {
				t.Errorf("results = %+v, want one entry with code %d", results, tt.code)
			}
		})
	}
}

func TestSendRawTransactionSyncEmptyReceipt(t *testing.T) {
	// A node that answers with a null result and no error leaves us nothing to
	// record; it must not count as a receipt.
	tx := testTx(t)
	rpc := &syncRPC{result: nil}

	tracker, err := send(t, rpc, &config.Config{}, tx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := tracker.receipts.Load(); got != 0 {
		t.Errorf("receipt count = %d, want 0", got)
	}
	if got := tracker.otherErrs.Load(); got != 1 {
		t.Errorf("other error count = %d, want 1", got)
	}
}

func TestSendRawTransactionSyncRevertedReceipt(t *testing.T) {
	tx := testTx(t)
	receipt := canonicalReceipt(tx.Hash())
	receipt["status"] = "0x0"
	rpc := &syncRPC{result: receipt}

	tracker, err := send(t, rpc, &config.Config{}, tx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := tracker.reverted.Load(); got != 1 {
		t.Errorf("reverted count = %d, want 1", got)
	}
}

func TestSendRawTransactionSyncReceiptWithoutStatus(t *testing.T) {
	tx := testTx(t)
	receipt := canonicalReceipt(tx.Hash())
	delete(receipt, "status")
	rpc := &syncRPC{result: receipt}

	tracker, err := send(t, rpc, &config.Config{}, tx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := tracker.noStatus.Load(); got != 1 {
		t.Errorf("no-status count = %d, want 1", got)
	}
	if got := tracker.reverted.Load(); got != 0 {
		t.Errorf("reverted count = %d, want 0 (absent status is not a revert)", got)
	}
}

func TestSyncReceiptClassification(t *testing.T) {
	zero := common.Hash{}
	block := common.HexToHash("0xbeef")
	num := (*hexutil.Big)(big.NewInt(9))
	yes, no := true, false
	zeroNum := (*hexutil.Big)(big.NewInt(0))

	tests := []struct {
		name            string
		receipt         *SyncReceipt
		wantSpeculative bool
	}{
		{"nil receipt", nil, false},
		{"no block fields", &SyncReceipt{}, true},
		{"zero block number", &SyncReceipt{BlockNumber: zeroNum, BlockHash: &block}, true},
		{"zero block hash", &SyncReceipt{BlockNumber: num, BlockHash: &zero}, true},
		{"nil block hash", &SyncReceipt{BlockNumber: num}, true},
		{"full block fields", &SyncReceipt{BlockNumber: num, BlockHash: &block}, false},
		{"marker true beats block fields", &SyncReceipt{BlockNumber: num, BlockHash: &block, Preconfirmation: &yes}, true},
		{"marker false beats absent block", &SyncReceipt{Preconfirmation: &no}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.receipt.Speculative(); got != tt.wantSpeculative {
				t.Errorf("Speculative() = %v, want %v", got, tt.wantSpeculative)
			}
		})
	}
}

func TestSyncTrackerNilIsSafe(t *testing.T) {
	var tracker *SyncTracker
	tracker.Record(common.Hash{}, &SyncReceipt{}, time.Second, nil)
	tracker.Stats()
	if got := tracker.Results(); got != nil {
		t.Errorf("Results() = %v, want nil", got)
	}
}

func TestSyncTrackerConcurrentRecord(t *testing.T) {
	tracker := NewSyncTracker()
	const goroutines = 32
	const each = 20

	var wg sync.WaitGroup
	for i := range goroutines {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := range each {
				hash := common.BigToHash(big.NewInt(int64(i*each + j)))
				if j%2 == 0 {
					tracker.Record(hash, &SyncReceipt{}, time.Millisecond, nil)
				} else {
					tracker.Record(hash, nil, time.Millisecond, fmt.Errorf("boom"))
				}
			}
		}(i)
	}
	wg.Wait()

	if got := tracker.submitted.Load(); got != goroutines*each {
		t.Errorf("submitted = %d, want %d", got, goroutines*each)
	}
	if got := len(tracker.Results()); got != goroutines*each {
		t.Errorf("results = %d, want %d", got, goroutines*each)
	}
	tracker.Stats()
}

func TestSyncTrackerCanonicalReceiptWithoutBlockNumber(t *testing.T) {
	// preconfirmation=false is authoritative, so this receipt is canonical even
	// though it carries no block fields to read a block number from.
	no := false
	tracker := NewSyncTracker()
	tracker.Record(common.Hash{}, &SyncReceipt{Preconfirmation: &no}, time.Millisecond, nil)

	if got := tracker.canonical.Load(); got != 1 {
		t.Errorf("canonical count = %d, want 1", got)
	}
	results := tracker.Results()
	if len(results) != 1 {
		t.Fatalf("results = %d, want 1", len(results))
	}
	if results[0].BlockNumber != 0 {
		t.Errorf("block number = %d, want 0", results[0].BlockNumber)
	}
}
