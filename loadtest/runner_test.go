package loadtest

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/0xPolygon/polygon-cli/loadtest/config"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func TestComputeReverseTxsPerAccount(t *testing.T) {
	newPool := func(n int) *AccountPool {
		accounts := make([]*Account, n)
		for i := range accounts {
			accounts[i] = &Account{address: common.BytesToAddress([]byte{byte(i + 1), byte(i >> 8)}), ready: true}
		}
		return newTestAccountPool(accounts...)
	}

	tests := []struct {
		name              string
		concurrency       int64
		requests          int64
		accounts          int
		wantErr           bool
		wantTxsPerAccount uint64
	}{
		{
			name:              "even split",
			concurrency:       250,
			requests:          10,
			accounts:          500,
			wantTxsPerAccount: 5,
		},
		{
			name:              "single account",
			concurrency:       2,
			requests:          10,
			accounts:          1,
			wantTxsPerAccount: 20,
		},
		{
			name:        "uneven split",
			concurrency: 3,
			requests:    7,
			accounts:    2,
			wantErr:     true,
		},
		{
			name:        "more accounts than requests",
			concurrency: 1,
			requests:    250,
			accounts:    500,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Runner{
				cfg: &config.Config{
					Concurrency: tt.concurrency,
					Requests:    tt.requests,
				},
				accountPool: newPool(tt.accounts),
			}

			err := r.computeReverseTxsPerAccount()
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr = %v", err, tt.wantErr)
			}
			if err == nil && r.reverseTxsPerAccount != tt.wantTxsPerAccount {
				t.Errorf("reverseTxsPerAccount = %d, want %d", r.reverseTxsPerAccount, tt.wantTxsPerAccount)
			}
		})
	}
}

func TestNonceTooLowRegexp(t *testing.T) {
	tests := []struct {
		name          string
		errMsg        string
		wantMatch     bool
		wantNextNonce uint64
	}{
		{
			name:          "geth format",
			errMsg:        "nonce too low: next nonce 78, tx nonce 44",
			wantMatch:     true,
			wantNextNonce: 78,
		},
		{
			name:          "geth format wrapped with context",
			errMsg:        "failed to send transaction: nonce too low: next nonce 78, tx nonce 44",
			wantMatch:     true,
			wantNextNonce: 78,
		},
		{
			name:      "nonce too low without details",
			errMsg:    "nonce too low",
			wantMatch: false,
		},
		{
			name:      "unrelated error",
			errMsg:    "insufficient funds for gas * price + value",
			wantMatch: false,
		},
		{
			name:      "different client format",
			errMsg:    "nonce too low: address 0xabc, tx: 44 state: 78",
			wantMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match := nonceTooLowRegexp.FindStringSubmatch(tt.errMsg)
			if (match != nil) != tt.wantMatch {
				t.Fatalf("match = %v, wantMatch = %v", match, tt.wantMatch)
			}
			if match == nil {
				return
			}
			nextNonce, err := strconv.ParseUint(match[1], 10, 64)
			if err != nil {
				t.Fatalf("failed to parse next nonce: %v", err)
			}
			if nextNonce != tt.wantNextNonce {
				t.Errorf("nextNonce = %d, want %d", nextNonce, tt.wantNextNonce)
			}
		})
	}
}

func TestCachedPricesNeverReturnsNil(t *testing.T) {
	// Every fallback in getSuggestedGasPrices can be reached before the cache
	// is written, and configureTransactOpts dereferences what it gets, so a nil
	// here used to panic the whole run.
	tests := []struct {
		name       string
		gasPrice   *big.Int
		gasTipCap  *big.Int
		wantPrice  int64
		wantTipCap int64
	}{
		{"empty cache", nil, nil, 0, 0},
		{"price only", big.NewInt(7), nil, 7, 0},
		{"tip cap only", nil, big.NewInt(3), 0, 3},
		{"both cached", big.NewInt(7), big.NewInt(3), 7, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Runner{cachedGasPrice: tt.gasPrice, cachedGasTipCap: tt.gasTipCap}

			gasPrice, gasTipCap := r.cachedPrices()
			if gasPrice == nil || gasTipCap == nil {
				t.Fatalf("cachedPrices() = (%v, %v), want non-nil values", gasPrice, gasTipCap)
			}
			if gasPrice.Int64() != tt.wantPrice {
				t.Errorf("gas price = %s, want %d", gasPrice, tt.wantPrice)
			}
			if gasTipCap.Int64() != tt.wantTipCap {
				t.Errorf("gas tip cap = %s, want %d", gasTipCap, tt.wantTipCap)
			}
		})
	}
}

func TestCachedPricesDoesNotAliasTheCache(t *testing.T) {
	// A caller that clamps the tip cap must not mutate the cached value.
	r := &Runner{cachedGasPrice: big.NewInt(10), cachedGasTipCap: big.NewInt(4)}

	_, gasTipCap := r.cachedPrices()
	gasTipCap.SetInt64(99)

	if got := r.cachedGasTipCap.Int64(); got != 4 {
		t.Errorf("cached tip cap = %d, want 4 (unchanged)", got)
	}
}

func TestWaitForFinalBlockRespectsCancellation(t *testing.T) {
	// The retry loop used a bare time.Sleep, so 30 retries held the process for
	// 150s past a Ctrl+C. It must return as soon as the context is done.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var body struct {
			ID     json.RawMessage `json:"id"`
			Method string          `json:"method"`
		}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// A node whose nonce never advances: the wait can never succeed.
		result := "0x1"
		if body.Method == "eth_getTransactionCount" {
			result = "0x0"
		}
		w.Header().Set("Content-Type", "application/json")
		if _, err := fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%s,"result":%q}`, body.ID, result); err != nil {
			t.Errorf("failed to write rpc response: %v", err)
		}
	}))
	defer server.Close()

	client, err := ethclient.DialContext(t.Context(), server.URL)
	if err != nil {
		t.Fatalf("failed to dial fake rpc: %v", err)
	}
	defer client.Close()

	pool := newTestAccountPool(&Account{
		address:    common.BytesToAddress([]byte{1}),
		ready:      true,
		startNonce: 0,
		nonce:      5, // outstanding transactions that will never be mined
	})

	// Negative rate limit means "no limit"; building a limiter from it used to
	// block every nonce check after the first forever.
	r := &Runner{
		cfg:         &config.Config{RateLimit: -1},
		client:      client,
		accountPool: pool,
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		_, waitErr := r.waitForFinalBlock(ctx)
		done <- waitErr
	}()

	// Give it long enough to enter the retry sleep, then cancel.
	time.Sleep(500 * time.Millisecond)
	cancel()

	select {
	case waitErr := <-done:
		if waitErr == nil {
			t.Error("expected an error after cancellation")
		}
	case <-time.After(10 * time.Second):
		t.Fatal("waitForFinalBlock did not return after context cancellation")
	}
}

func TestWaitForFinalBlockCompletesWithoutRateLimit(t *testing.T) {
	// With --rate-limit -1 and all nonces mined, the wait must finish promptly
	// rather than blocking on a limiter whose tokens never come back.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var body struct {
			ID     json.RawMessage `json:"id"`
			Method string          `json:"method"`
		}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		result := "0x63" // block number, and a nonce well past what we expect
		w.Header().Set("Content-Type", "application/json")
		if _, err := fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%s,"result":%q}`, body.ID, result); err != nil {
			t.Errorf("failed to write rpc response: %v", err)
		}
	}))
	defer server.Close()

	client, err := ethclient.DialContext(t.Context(), server.URL)
	if err != nil {
		t.Fatalf("failed to dial fake rpc: %v", err)
	}
	defer client.Close()

	// Several accounts: with a negative-rate limiter, the first would pass and
	// the rest would hang.
	accounts := make([]*Account, 5)
	for i := range accounts {
		accounts[i] = &Account{
			address:    common.BytesToAddress([]byte{byte(i + 1)}),
			ready:      true,
			startNonce: 0,
			nonce:      3,
		}
	}

	r := &Runner{
		cfg:         &config.Config{RateLimit: -1},
		client:      client,
		accountPool: newTestAccountPool(accounts...),
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		blockNumber, waitErr := r.waitForFinalBlock(t.Context())
		if waitErr != nil {
			t.Errorf("waitForFinalBlock failed: %v", waitErr)
		}
		if blockNumber != 0x63 {
			t.Errorf("block number = %d, want 99", blockNumber)
		}
	}()

	select {
	case <-done:
	case <-time.After(15 * time.Second):
		t.Fatal("waitForFinalBlock blocked with rate limiting disabled")
	}
}
