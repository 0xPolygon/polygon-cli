package loadtest

import (
	"context"
	"math"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

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
	if err := ap.AddReusableNonce(ctx, addr, 14); err != nil {
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
