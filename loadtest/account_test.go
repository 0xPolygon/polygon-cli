package loadtest

import (
	"context"
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
