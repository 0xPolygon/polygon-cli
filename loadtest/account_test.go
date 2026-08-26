package loadtest

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// rpcRequest is the subset of the JSON-RPC request the fake server needs.
type rpcRequest struct {
	ID     json.RawMessage `json:"id"`
	Method string          `json:"method"`
}

// fakeRPC is a minimal JSON-RPC server that serves the calls the account pool
// makes during construction and nonce fetching. It records the peak number of
// concurrent eth_getTransactionCount calls and can fail the first N of them.
type fakeRPC struct {
	server *httptest.Server

	mu           sync.Mutex
	inFlight     int
	maxInFlight  int
	nonceCalls   int
	failFirstN   int
	nonceLatency time.Duration

	totalNonceCalls atomic.Int64
}

func newFakeRPC(t *testing.T, failFirstN int, nonceLatency time.Duration) *fakeRPC {
	t.Helper()
	f := &fakeRPC{failFirstN: failFirstN, nonceLatency: nonceLatency}
	f.server = httptest.NewServer(http.HandlerFunc(f.handle))
	t.Cleanup(f.server.Close)
	return f
}

func (f *fakeRPC) handle(w http.ResponseWriter, r *http.Request) {
	var req rpcRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeResult := func(result string) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%s,"result":%q}`, req.ID, result)
	}

	switch req.Method {
	case "eth_chainId":
		writeResult("0x89")
	case "eth_blockNumber":
		writeResult("0x1")
	case "eth_getTransactionCount":
		f.totalNonceCalls.Add(1)

		f.mu.Lock()
		f.nonceCalls++
		shouldFail := f.nonceCalls <= f.failFirstN
		f.inFlight++
		if f.inFlight > f.maxInFlight {
			f.maxInFlight = f.inFlight
		}
		f.mu.Unlock()

		defer func() {
			f.mu.Lock()
			f.inFlight--
			f.mu.Unlock()
		}()

		if f.nonceLatency > 0 {
			time.Sleep(f.nonceLatency)
		}

		if shouldFail {
			// Mimic a load-balanced endpoint shedding load.
			http.Error(w, "the server encountered an error", http.StatusInternalServerError)
			return
		}
		writeResult("0x2a")
	default:
		http.Error(w, "unexpected method "+req.Method, http.StatusBadRequest)
	}
}

func (f *fakeRPC) peakConcurrency() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.maxInFlight
}

// newTestPool builds an account pool pointed at the fake RPC with n accounts
// that have no nonce yet, so FetchNoncesInParallel has work to do.
func newTestPool(t *testing.T, f *fakeRPC, concurrency int64, n int) *AccountPool {
	t.Helper()
	ctx := context.Background()

	client, err := ethclient.DialContext(ctx, f.server.URL)
	if err != nil {
		t.Fatalf("failed to dial fake rpc: %v", err)
	}
	t.Cleanup(client.Close)

	fundingKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed to generate funding key: %v", err)
	}

	ap, err := NewAccountPool(ctx, client, &AccountPoolConfig{
		FundingPrivateKey: fundingKey,
		FundingAmount:     big.NewInt(0),
		Concurrency:       concurrency,
	})
	if err != nil {
		t.Fatalf("failed to create account pool: %v", err)
	}

	for range n {
		key, err := crypto.GenerateKey()
		if err != nil {
			t.Fatalf("failed to generate account key: %v", err)
		}
		if err := ap.Add(ctx, key, nil); err != nil {
			t.Fatalf("failed to add account: %v", err)
		}
	}

	return ap
}

// newTestAccountPool builds a bare account pool around pre-constructed
// accounts, with no RPC client behind it.
func newTestAccountPool(accounts ...*Account) *AccountPool {
	ap := &AccountPool{
		accounts:          accounts,
		accountsPositions: make(map[common.Address]int),
		cfg:               &AccountPoolConfig{},
	}
	for i, acc := range accounts {
		ap.accountsPositions[acc.address] = i
	}
	return ap
}

func TestFetchNoncesInParallelRespectsConcurrency(t *testing.T) {
	const accounts = 50
	const concurrency = 4

	f := newFakeRPC(t, 0, 5*time.Millisecond)
	ap := newTestPool(t, f, concurrency, accounts)

	if err := ap.FetchNoncesInParallel(context.Background()); err != nil {
		t.Fatalf("FetchNoncesInParallel failed: %v", err)
	}

	if peak := f.peakConcurrency(); peak > concurrency {
		t.Errorf("peak in-flight nonce requests = %d, want <= %d", peak, concurrency)
	}

	ready, rdyCount, total := ap.AllAccountsReady()
	if !ready {
		t.Errorf("accounts not all ready: %d/%d", rdyCount, total)
	}
	for _, acc := range ap.accounts {
		if acc.nonce != 0x2a || acc.startNonce != 0x2a {
			t.Errorf("account %s nonce = %d/%d, want 42/42", acc.address, acc.nonce, acc.startNonce)
		}
	}
}

func TestFetchNoncesInParallelClampsConcurrency(t *testing.T) {
	tests := []struct {
		name        string
		concurrency int64
	}{
		{"zero falls back to one", 0},
		{"negative falls back to one", -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newFakeRPC(t, 0, time.Millisecond)
			ap := newTestPool(t, f, tt.concurrency, 5)

			if err := ap.FetchNoncesInParallel(context.Background()); err != nil {
				t.Fatalf("FetchNoncesInParallel failed: %v", err)
			}
			if peak := f.peakConcurrency(); peak != 1 {
				t.Errorf("peak in-flight nonce requests = %d, want 1", peak)
			}
		})
	}
}

func TestFetchNoncesInParallelRetriesTransientFailures(t *testing.T) {
	// Fail the first two calls so the single account needs retries to succeed.
	f := newFakeRPC(t, 2, 0)
	ap := newTestPool(t, f, 1, 1)

	if err := ap.FetchNoncesInParallel(context.Background()); err != nil {
		t.Fatalf("FetchNoncesInParallel failed despite retries: %v", err)
	}
	if calls := f.totalNonceCalls.Load(); calls != 3 {
		t.Errorf("nonce calls = %d, want 3 (two failures plus one success)", calls)
	}
	if ready, rdyCount, total := ap.AllAccountsReady(); !ready {
		t.Errorf("accounts not all ready: %d/%d", rdyCount, total)
	}
}

func TestFetchNoncesInParallelFailsAfterMaxAttempts(t *testing.T) {
	f := newFakeRPC(t, nonceFetchMaxAttempts, 0)
	ap := newTestPool(t, f, 1, 1)

	err := ap.FetchNoncesInParallel(context.Background())
	if err == nil {
		t.Fatal("expected error after exhausting retries, got nil")
	}
	if calls := f.totalNonceCalls.Load(); calls != int64(nonceFetchMaxAttempts) {
		t.Errorf("nonce calls = %d, want %d", calls, nonceFetchMaxAttempts)
	}
}

func TestFetchNoncesInParallelHonorsCancellation(t *testing.T) {
	// Every request fails, so the fetch is stuck in backoff when we cancel.
	f := newFakeRPC(t, 1000, 0)
	ap := newTestPool(t, f, 1, 20)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	done := make(chan error, 1)
	go func() { done <- ap.FetchNoncesInParallel(ctx) }()

	select {
	case err := <-done:
		if err == nil {
			t.Fatal("expected error on cancellation, got nil")
		}
	case <-time.After(10 * time.Second):
		t.Fatal("FetchNoncesInParallel did not return after context cancellation")
	}
}

func TestFetchNoncesInParallelNoop(t *testing.T) {
	f := newFakeRPC(t, 0, 0)
	// A pool whose accounts all have a forced start nonce needs no fetching.
	ap := newTestPool(t, f, 4, 0)

	ctx := context.Background()
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}
	startNonce := uint64(7)
	if err := ap.Add(ctx, key, &startNonce); err != nil {
		t.Fatalf("failed to add account: %v", err)
	}

	if err := ap.FetchNoncesInParallel(ctx); err != nil {
		t.Fatalf("FetchNoncesInParallel failed: %v", err)
	}
	if calls := f.totalNonceCalls.Load(); calls != 0 {
		t.Errorf("nonce calls = %d, want 0", calls)
	}
}

func TestFastForwardNonce(t *testing.T) {
	addr := common.HexToAddress("0x1")

	tests := []struct {
		name               string
		nonce              uint64
		reusableNonces     []uint64
		nextNonce          uint64
		wantUpdated        bool
		wantNonce          uint64
		wantReusableNonces []uint64
	}{
		{
			name:        "fast forward when next nonce is higher",
			nonce:       44,
			nextNonce:   78,
			wantUpdated: true,
			wantNonce:   78,
		},
		{
			name:        "no rewind when next nonce is lower",
			nonce:       100,
			nextNonce:   78,
			wantUpdated: false,
			wantNonce:   100,
		},
		{
			name:        "no update when next nonce is equal",
			nonce:       78,
			nextNonce:   78,
			wantUpdated: false,
			wantNonce:   78,
		},
		{
			name:               "drops reusable nonces below next nonce",
			nonce:              44,
			reusableNonces:     []uint64{40, 42, 77, 78, 80},
			nextNonce:          78,
			wantUpdated:        true,
			wantNonce:          78,
			wantReusableNonces: []uint64{78, 80},
		},
		{
			name:               "keeps reusable nonces when not rewinding",
			nonce:              100,
			reusableNonces:     []uint64{90, 95},
			nextNonce:          78,
			wantUpdated:        false,
			wantNonce:          100,
			wantReusableNonces: []uint64{90, 95},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reusable := make([]uint64, len(tt.reusableNonces))
			copy(reusable, tt.reusableNonces)
			acc := &Account{
				address:        addr,
				nonce:          tt.nonce,
				reusableNonces: reusable,
			}
			ap := newTestAccountPool(acc)

			updated, err := ap.FastForwardNonce(context.Background(), addr, tt.nextNonce)
			if err != nil {
				t.Fatalf("FastForwardNonce returned error: %v", err)
			}
			if updated != tt.wantUpdated {
				t.Errorf("updated = %v, want %v", updated, tt.wantUpdated)
			}
			if acc.nonce != tt.wantNonce {
				t.Errorf("nonce = %d, want %d", acc.nonce, tt.wantNonce)
			}
			if len(acc.reusableNonces) != len(tt.wantReusableNonces) {
				t.Fatalf("reusableNonces = %v, want %v", acc.reusableNonces, tt.wantReusableNonces)
			}
			for i, n := range tt.wantReusableNonces {
				if acc.reusableNonces[i] != n {
					t.Fatalf("reusableNonces = %v, want %v", acc.reusableNonces, tt.wantReusableNonces)
				}
			}
		})
	}
}

func TestPrepareReverseNonces(t *testing.T) {
	addr1 := common.HexToAddress("0x1")
	addr2 := common.HexToAddress("0x2")

	acc1 := &Account{address: addr1, nonce: 0, ready: true}
	acc2 := &Account{address: addr2, nonce: 100, ready: true}
	ap := newTestAccountPool(acc1, acc2)

	if err := ap.PrepareReverseNonces(5); err != nil {
		t.Fatalf("PrepareReverseNonces returned error: %v", err)
	}

	if acc1.startNonce != 0 || acc1.nonce != 4 {
		t.Errorf("acc1 = {startNonce: %d, nonce: %d}, want {0, 4}", acc1.startNonce, acc1.nonce)
	}
	if acc2.startNonce != 100 || acc2.nonce != 104 {
		t.Errorf("acc2 = {startNonce: %d, nonce: %d}, want {100, 104}", acc2.startNonce, acc2.nonce)
	}
}

func TestPrepareReverseNoncesErrors(t *testing.T) {
	addr := common.HexToAddress("0x1")

	t.Run("zero txs per account", func(t *testing.T) {
		ap := newTestAccountPool(&Account{address: addr, ready: true})
		if err := ap.PrepareReverseNonces(0); err == nil {
			t.Error("expected error for zero txsPerAccount, got nil")
		}
	})

	t.Run("account not ready", func(t *testing.T) {
		ap := newTestAccountPool(&Account{address: addr, ready: false})
		if err := ap.PrepareReverseNonces(5); err == nil {
			t.Error("expected error for not-ready account, got nil")
		}
	})

	t.Run("nonce overflow", func(t *testing.T) {
		ap := newTestAccountPool(&Account{address: addr, nonce: math.MaxUint64 - 1, ready: true})
		if err := ap.PrepareReverseNonces(5); err == nil {
			t.Error("expected error for nonce overflow, got nil")
		}
	})
}

func TestNextReverseNonceOrder(t *testing.T) {
	addr1 := common.HexToAddress("0x1")
	addr2 := common.HexToAddress("0x2")

	acc1 := &Account{address: addr1, nonce: 0, ready: true}
	acc2 := &Account{address: addr2, nonce: 100, ready: true}
	ap := newTestAccountPool(acc1, acc2)
	ap.cfg.ReverseNonceOrder = true

	const txsPerAccount = 3
	if err := ap.PrepareReverseNonces(txsPerAccount); err != nil {
		t.Fatalf("PrepareReverseNonces returned error: %v", err)
	}

	// Round-robin across 2 accounts, each walking its own range downward.
	want := []struct {
		addr  common.Address
		nonce uint64
	}{
		{addr1, 2}, {addr2, 102},
		{addr1, 1}, {addr2, 101},
		{addr1, 0}, {addr2, 100},
	}

	ctx := context.Background()
	for i, w := range want {
		acc, err := ap.Next(ctx)
		if err != nil {
			t.Fatalf("Next() call %d returned error: %v", i, err)
		}
		if acc.Address() != w.addr || acc.Nonce() != w.nonce {
			t.Errorf("Next() call %d = {%s, %d}, want {%s, %d}", i, acc.Address(), acc.Nonce(), w.addr, w.nonce)
		}
	}

	// All ranges are exhausted; further calls must fail rather than underflow
	// below the floor nonce.
	if _, err := ap.Next(ctx); err == nil {
		t.Error("expected error after all account ranges exhausted, got nil")
	}
}

func TestNextReverseNonceOrderReusableNonce(t *testing.T) {
	addr := common.HexToAddress("0x1")

	acc := &Account{address: addr, nonce: 10, ready: true}
	ap := newTestAccountPool(acc)
	ap.cfg.ReverseNonceOrder = true

	if err := ap.PrepareReverseNonces(5); err != nil {
		t.Fatalf("PrepareReverseNonces returned error: %v", err)
	}

	ctx := context.Background()

	// First call hands the top of the range (14) and moves the counter down.
	first, err := ap.Next(ctx)
	if err != nil {
		t.Fatalf("Next() returned error: %v", err)
	}
	if first.Nonce() != 14 {
		t.Fatalf("first nonce = %d, want 14", first.Nonce())
	}

	// Simulate a failed send of nonce 14: it goes back as reusable and must be
	// handed out next without moving the descending counter.
	if err = ap.AddReusableNonce(ctx, addr, 14); err != nil {
		t.Fatalf("AddReusableNonce returned error: %v", err)
	}

	retry, err := ap.Next(ctx)
	if err != nil {
		t.Fatalf("Next() returned error: %v", err)
	}
	if retry.Nonce() != 14 {
		t.Errorf("retry nonce = %d, want reused 14", retry.Nonce())
	}

	next, err := ap.Next(ctx)
	if err != nil {
		t.Fatalf("Next() returned error: %v", err)
	}
	if next.Nonce() != 13 {
		t.Errorf("nonce after retry = %d, want 13", next.Nonce())
	}
}

func TestFastForwardNonceUnknownAccount(t *testing.T) {
	ap := newTestAccountPool()

	updated, err := ap.FastForwardNonce(context.Background(), common.HexToAddress("0x2"), 10)
	if err == nil {
		t.Fatal("expected error for unknown account, got nil")
	}
	if updated {
		t.Error("updated = true, want false for unknown account")
	}
}
