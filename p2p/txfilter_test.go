package p2p

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/prometheus/client_golang/prometheus/testutil"

	ds "github.com/0xPolygon/polygon-cli/p2p/datastructures"
)

const testChainID = 137

var testSigner = types.LatestSignerForChainID(big.NewInt(testChainID))

// newTestKey returns a key and its address.
func newTestKey(t *testing.T) (*ecdsa.PrivateKey, common.Address) {
	t.Helper()
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}
	return key, crypto.PubkeyToAddress(key.PublicKey)
}

// newTestTx signs a dynamic-fee transaction with the given nonce and tip.
func newTestTx(t *testing.T, key *ecdsa.PrivateKey, nonce uint64, tipGwei int64) *types.Transaction {
	t.Helper()
	tip := new(big.Int).Mul(big.NewInt(tipGwei), big.NewInt(1e9))
	tx, err := types.SignNewTx(key, testSigner, &types.DynamicFeeTx{
		ChainID:   big.NewInt(testChainID),
		Nonce:     nonce,
		GasTipCap: tip,
		GasFeeCap: new(big.Int).Add(tip, big.NewInt(1e9)),
		Gas:       21000,
		To:        &common.Address{},
	})
	if err != nil {
		t.Fatalf("failed to sign tx: %v", err)
	}
	return tx
}

// newTestFilter builds a filter, failing the test if the options produce none.
func newTestFilter(t *testing.T, opts TxFilterOptions) *TxFilter {
	t.Helper()
	if opts.ChainID == 0 {
		opts.ChainID = testChainID
	}
	if opts.NonceCache.MaxSize == 0 {
		opts.NonceCache.MaxSize = 1024
	}
	f, err := NewTxFilter(t.Context(), opts)
	if err != nil {
		t.Fatalf("failed to create filter: %v", err)
	}
	if f == nil {
		t.Fatal("expected a filter, got nil")
	}
	t.Cleanup(f.Close)
	return f
}

// hashesOf is a readable form for comparing allow results.
func hashesOf(txs []*types.Transaction) []common.Hash {
	out := make([]common.Hash, len(txs))
	for i, tx := range txs {
		out[i] = tx.Hash()
	}
	return out
}

func TestTxFilterDisabledByDefault(t *testing.T) {
	f, err := NewTxFilter(t.Context(), TxFilterOptions{ChainID: testChainID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f != nil {
		t.Fatal("expected nil filter when no gate is enabled")
	}

	// A nil filter must behave as allow-everything for both entry points.
	key, _ := newTestKey(t)
	txs := []*types.Transaction{newTestTx(t, key, 0, 30)}
	if got := f.Allow(txs); len(got) != 1 {
		t.Errorf("nil filter allowed %d txs, want 1", len(got))
	}
	f.ObserveMined(txs) // must not panic
}

func TestTxFilterGasFloor(t *testing.T) {
	f := newTestFilter(t, TxFilterOptions{MinTip: 25e9})

	key, _ := newTestKey(t)
	lowTip := newTestTx(t, key, 0, 1)
	atFloor := newTestTx(t, key, 1, 25)
	aboveFloor := newTestTx(t, key, 2, 30)

	before := testutil.ToFloat64(txFilterDropped.WithLabelValues(dropReasonLowTip))
	allowed := f.Allow([]*types.Transaction{lowTip, atFloor, aboveFloor})
	after := testutil.ToFloat64(txFilterDropped.WithLabelValues(dropReasonLowTip))

	want := []common.Hash{atFloor.Hash(), aboveFloor.Hash()}
	if got := hashesOf(allowed); len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("allowed = %v, want %v", got, want)
	}
	if delta := after - before; delta != 1 {
		t.Errorf("low_tip drops = %v, want 1", delta)
	}
}

func TestTxFilterStaleNonceFromObservedBlock(t *testing.T) {
	f := newTestFilter(t, TxFilterOptions{GateStaleTxs: true})

	key, addr := newTestKey(t)

	// Before observing anything, an unknown sender is not judged stale.
	old := newTestTx(t, key, 3, 30)
	if got := f.Allow([]*types.Transaction{old}); len(got) != 1 {
		t.Fatalf("unknown sender: allowed %d txs, want 1 (fail open)", len(got))
	}

	// A mined transaction at nonce 7 proves the next nonce is 8.
	f.ObserveMined([]*types.Transaction{newTestTx(t, key, 7, 30)})

	if next, ok := f.nonces.Peek(addr); !ok || next != 8 {
		t.Fatalf("tracked nonce = %d (present %v), want 8", next, ok)
	}

	stale := newTestTx(t, key, 7, 30)
	alsoStale := newTestTx(t, key, 0, 30)
	current := newTestTx(t, key, 8, 30)
	future := newTestTx(t, key, 12, 30)

	before := testutil.ToFloat64(txFilterDropped.WithLabelValues(dropReasonStaleNonce))
	allowed := f.Allow([]*types.Transaction{stale, alsoStale, current, future})
	after := testutil.ToFloat64(txFilterDropped.WithLabelValues(dropReasonStaleNonce))

	want := []common.Hash{current.Hash(), future.Hash()}
	if got := hashesOf(allowed); len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("allowed = %v, want %v", got, want)
	}
	if delta := after - before; delta != 2 {
		t.Errorf("stale_nonce drops = %v, want 2", delta)
	}
}

func TestTxFilterObservedNonceNeverGoesBackwards(t *testing.T) {
	f := newTestFilter(t, TxFilterOptions{GateStaleTxs: true})
	key, addr := newTestKey(t)

	f.ObserveMined([]*types.Transaction{newTestTx(t, key, 9, 30)})
	// An out-of-order block carrying an older transaction must not lower the
	// tracked nonce, or stale transactions would start passing again.
	f.ObserveMined([]*types.Transaction{newTestTx(t, key, 2, 30)})

	if next, _ := f.nonces.Peek(addr); next != 10 {
		t.Errorf("tracked nonce = %d, want 10", next)
	}
}

func TestTxFilterRateLimit(t *testing.T) {
	// One token per second with a burst of 2: the third transaction is capped.
	f := newTestFilter(t, TxFilterOptions{RateLimit: 1, RateLimitBurst: 2})

	key, _ := newTestKey(t)
	txs := []*types.Transaction{
		newTestTx(t, key, 0, 30),
		newTestTx(t, key, 1, 30),
		newTestTx(t, key, 2, 30),
		newTestTx(t, key, 3, 30),
	}

	before := testutil.ToFloat64(txFilterDropped.WithLabelValues(dropReasonRateLimited))
	allowed := f.Allow(txs)
	after := testutil.ToFloat64(txFilterDropped.WithLabelValues(dropReasonRateLimited))

	if len(allowed) != 2 {
		t.Errorf("allowed %d txs, want 2 (the burst)", len(allowed))
	}
	if delta := after - before; delta != 2 {
		t.Errorf("rate_limited drops = %v, want 2", delta)
	}
}

func TestTxFilterRejectedTxsDoNotSpendRateBudget(t *testing.T) {
	// The gates run cheapest-first and the limiter runs last, so a flood of
	// unminable transactions must not consume the budget for good ones.
	f := newTestFilter(t, TxFilterOptions{MinTip: 25e9, RateLimit: 1, RateLimitBurst: 1})

	key, _ := newTestKey(t)
	txs := []*types.Transaction{
		newTestTx(t, key, 0, 1),
		newTestTx(t, key, 1, 1),
		newTestTx(t, key, 2, 1),
		newTestTx(t, key, 3, 1),
		newTestTx(t, key, 4, 30),
	}

	allowed := f.Allow(txs)
	if len(allowed) != 1 {
		t.Fatalf("allowed %d txs, want 1", len(allowed))
	}
	if allowed[0].Nonce() != 4 {
		t.Errorf("allowed tx nonce = %d, want 4 (the well-priced one)", allowed[0].Nonce())
	}
}

func TestTxFilterLogOnly(t *testing.T) {
	f := newTestFilter(t, TxFilterOptions{MinTip: 25e9, GateStaleTxs: true, LogOnly: true})

	key, _ := newTestKey(t)
	f.ObserveMined([]*types.Transaction{newTestTx(t, key, 5, 30)})

	txs := []*types.Transaction{
		newTestTx(t, key, 0, 1),  // low tip and stale
		newTestTx(t, key, 6, 30), // fine
	}

	before := testutil.ToFloat64(txFilterDropped.WithLabelValues(dropReasonLowTip))
	allowed := f.Allow(txs)
	after := testutil.ToFloat64(txFilterDropped.WithLabelValues(dropReasonLowTip))

	if len(allowed) != 2 {
		t.Errorf("log-only allowed %d txs, want 2 (nothing withheld)", len(allowed))
	}
	if delta := after - before; delta != 1 {
		t.Errorf("low_tip drops = %v, want 1 (counted but not enforced)", delta)
	}
}

func TestTxFilterUnsignedTxIsNotJudgedStale(t *testing.T) {
	f := newTestFilter(t, TxFilterOptions{GateStaleTxs: true})

	// No signature, so the sender cannot be recovered. The decode path already
	// accepted it; the staleness gate has no opinion.
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   big.NewInt(testChainID),
		Nonce:     0,
		GasTipCap: big.NewInt(1e9),
		GasFeeCap: big.NewInt(2e9),
		Gas:       21000,
		To:        &common.Address{},
	})

	if got := f.Allow([]*types.Transaction{tx}); len(got) != 1 {
		t.Errorf("allowed %d txs, want 1", len(got))
	}
}

// nonceRPC is a fake JSON-RPC endpoint serving eth_getTransactionCount.
type nonceRPC struct {
	t      *testing.T
	nonces map[common.Address]uint64
	calls  atomic.Int64
}

func newNonceRPC(t *testing.T, nonces map[common.Address]uint64) (*nonceRPC, string) {
	t.Helper()
	n := &nonceRPC{t: t, nonces: nonces}
	server := httptest.NewServer(http.HandlerFunc(n.handle))
	t.Cleanup(server.Close)
	return n, server.URL
}

func (n *nonceRPC) handle(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID     json.RawMessage `json:"id"`
		Method string          `json:"method"`
		Params []any           `json:"params"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var result string
	switch req.Method {
	case "eth_chainId":
		result = "0x89"
	case "eth_getTransactionCount":
		n.calls.Add(1)
		addr := common.HexToAddress(req.Params[0].(string))
		result = fmt.Sprintf("0x%x", n.nonces[addr])
	default:
		http.Error(w, "unexpected method "+req.Method, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%s,"result":%q}`, req.ID, result); err != nil {
		n.t.Errorf("failed to write rpc response: %v", err)
	}
}

func TestTxFilterNonceRPCFallback(t *testing.T) {
	key, addr := newTestKey(t)
	rpc, url := newNonceRPC(t, map[common.Address]uint64{addr: 42})

	f := newTestFilter(t, TxFilterOptions{GateStaleTxs: true, NonceRPCURL: url})

	// First sight of an unknown sender: allowed, and a lookup is queued.
	if got := f.Allow([]*types.Transaction{newTestTx(t, key, 5, 30)}); len(got) != 1 {
		t.Fatalf("allowed %d txs, want 1 (fail open on unknown sender)", len(got))
	}

	// Once the lookup lands, the sender is known and stale txs are withheld.
	deadline := time.Now().Add(5 * time.Second)
	for {
		if next, ok := f.nonces.Peek(addr); ok && next == 42 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("nonce lookup did not populate the cache")
		}
		time.Sleep(10 * time.Millisecond)
	}

	if got := f.Allow([]*types.Transaction{newTestTx(t, key, 5, 30)}); len(got) != 0 {
		t.Errorf("allowed %d stale txs after lookup, want 0", len(got))
	}
	if got := f.Allow([]*types.Transaction{newTestTx(t, key, 42, 30)}); len(got) != 1 {
		t.Errorf("allowed %d current txs, want 1", len(got))
	}
	if calls := rpc.calls.Load(); calls != 1 {
		t.Errorf("rpc lookups = %d, want 1", calls)
	}
}

func TestTxFilterNonceRPCFallbackDedupesInFlight(t *testing.T) {
	key, addr := newTestKey(t)
	rpc, url := newNonceRPC(t, map[common.Address]uint64{addr: 3})

	f := newTestFilter(t, TxFilterOptions{GateStaleTxs: true, NonceRPCURL: url})

	// Many transactions from the same unknown sender must not fan out into many
	// lookups.
	txs := make([]*types.Transaction, 0, 50)
	for i := range uint64(50) {
		txs = append(txs, newTestTx(t, key, i, 30))
	}
	f.Allow(txs)

	deadline := time.Now().Add(5 * time.Second)
	for {
		if _, ok := f.nonces.Peek(addr); ok {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("nonce lookup did not populate the cache")
		}
		time.Sleep(10 * time.Millisecond)
	}

	if calls := rpc.calls.Load(); calls != 1 {
		t.Errorf("rpc lookups = %d, want 1 for a single sender", calls)
	}
}

func TestTxFilterNoRPCFallbackWhenUnset(t *testing.T) {
	f := newTestFilter(t, TxFilterOptions{GateStaleTxs: true})
	if f.lookups != nil || f.rpc != nil {
		t.Error("expected no RPC fallback machinery when NonceRPCURL is empty")
	}

	key, _ := newTestKey(t)
	if got := f.Allow([]*types.Transaction{newTestTx(t, key, 0, 30)}); len(got) != 1 {
		t.Errorf("allowed %d txs, want 1", len(got))
	}
}

func TestTxFilterConcurrentAccess(t *testing.T) {
	f := newTestFilter(t, TxFilterOptions{
		GateStaleTxs: true,
		MinTip:       25e9,
		RateLimit:    1e6,
		NonceCache:   ds.LRUOptions{MaxSize: 128},
	})

	const goroutines = 32
	var wg sync.WaitGroup
	for i := range goroutines {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key, _ := newTestKey(t)
			for n := range uint64(20) {
				tx := newTestTx(t, key, n, 30)
				f.ObserveMined([]*types.Transaction{tx})
				f.Allow([]*types.Transaction{tx})
			}
		}(i)
	}
	wg.Wait()
}

func TestTxFilterAllowEmpty(t *testing.T) {
	f := newTestFilter(t, TxFilterOptions{MinTip: 25e9})
	if got := f.Allow(nil); got != nil {
		t.Errorf("Allow(nil) = %v, want nil", got)
	}
}

func TestTxFilterContextCancellationStopsWorkers(t *testing.T) {
	key, addr := newTestKey(t)
	_, url := newNonceRPC(t, map[common.Address]uint64{addr: 1})

	ctx, cancel := context.WithCancel(context.Background())
	f, err := NewTxFilter(ctx, TxFilterOptions{
		ChainID:      testChainID,
		GateStaleTxs: true,
		NonceCache:   ds.LRUOptions{MaxSize: 16},
		NonceRPCURL:  url,
	})
	if err != nil {
		t.Fatalf("failed to create filter: %v", err)
	}
	defer f.Close()

	cancel()
	// Queued work after cancellation must not block the caller.
	done := make(chan struct{})
	go func() {
		defer close(done)
		f.Allow([]*types.Transaction{newTestTx(t, key, 0, 30)})
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("Allow blocked after context cancellation")
	}
}
