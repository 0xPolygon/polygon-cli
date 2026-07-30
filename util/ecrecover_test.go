package util

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/clique"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// TestEcrecoverDoesNotPanicOnUntrustedHeader covers every trailing optional field
// clique's encodeSigHeader panics on. Headers arrive from peers, and the sensor's
// write path calls Ecrecover on the peer's own goroutine, which geth does not
// recover around -- so a panic here is a remote process kill, not a bad row.
//
// If a geth bump adds another panicking field, this test keeps passing (the
// recover is generic) but the new field belongs in the table so the coverage
// stays visible.
func TestEcrecoverDoesNotPanicOnUntrustedHeader(t *testing.T) {
	hash := common.HexToHash("0xdead")
	gas := uint64(1)

	for _, tc := range []struct {
		name  string
		apply func(*types.Header)
	}{
		{"withdrawals_hash", func(h *types.Header) { h.WithdrawalsHash = &hash }},
		{"excess_blob_gas", func(h *types.Header) { h.ExcessBlobGas = &gas }},
		{"blob_gas_used", func(h *types.Header) { h.BlobGasUsed = &gas }},
		{"parent_beacon_root", func(h *types.Header) { h.ParentBeaconRoot = &hash }},
		{"slot_number", func(h *types.Header) { h.SlotNumber = &gas }}, // EIP-7843, geth v1.17.4
	} {
		t.Run(tc.name, func(t *testing.T) {
			h := &types.Header{
				Number:     big.NewInt(1),
				Difficulty: big.NewInt(1),
				Extra:      make([]byte, crypto.SignatureLength),
			}
			tc.apply(h)

			// Must return an error, and must not panic.
			signer, err := Ecrecover(h)
			if err == nil {
				t.Fatalf("expected an error for a header clique rejects, got signer %x", signer)
			}
			if signer != nil {
				t.Fatalf("expected nil signer alongside the error, got %x", signer)
			}
		})
	}
}

// A well-formed clique-sealed header must still recover, so the recover() above
// cannot be masking a real regression.
func TestEcrecoverStillRecoversValidHeader(t *testing.T) {
	priv, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	h := &types.Header{
		Number:     big.NewInt(1),
		Difficulty: big.NewInt(1),
		Extra:      make([]byte, crypto.SignatureLength),
	}
	sig, err := crypto.Sign(clique.SealHash(h).Bytes(), priv)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	copy(h.Extra[len(h.Extra)-crypto.SignatureLength:], sig)

	signer, err := Ecrecover(h)
	if err != nil {
		t.Fatalf("expected recovery to succeed: %v", err)
	}
	if got, want := common.BytesToAddress(signer), crypto.PubkeyToAddress(priv.PublicKey); got != want {
		t.Fatalf("signer mismatch: got %s want %s", got, want)
	}
}
