package p2p

import (
	"context"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/clique"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

const validatorSetFixture = "../internal/heimdall/client/testdata/rest/stake_validators_set.json"

// knownSigner is the first signer in the validator-set fixture.
var knownSigner = common.HexToAddress("0x02f615e95563ef16f10354dba9e584e58d2d4314")

// readFixture reads the validator-set fixture, failing the test if it cannot.
func readFixture(t *testing.T) []byte {
	t.Helper()
	body, err := os.ReadFile(validatorSetFixture)
	if err != nil {
		t.Fatalf("reading fixture: %v", err)
	}
	return body
}

// newFixtureServer serves the given body for /stake/validators-set.
func newFixtureServer(t *testing.T, status int, body []byte) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/stake/validators-set" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(status)
		_, _ = w.Write(body)
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestValidatorSetStartAndHasSigner(t *testing.T) {
	srv := newFixtureServer(t, http.StatusOK, readFixture(t))

	s := NewValidatorSet(srv.URL, time.Hour)
	if err := s.Start(context.Background()); err != nil {
		t.Fatalf("Start: %v", err)
	}

	if !s.HasSigner(knownSigner) {
		t.Errorf("expected %s to be a known signer", knownSigner)
	}
	if s.HasSigner(common.HexToAddress("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")) {
		t.Errorf("expected an unknown address to be rejected")
	}
}

func TestValidatorSetStartEmptyFails(t *testing.T) {
	srv := newFixtureServer(t, http.StatusOK, []byte(`{"validator_set":{"validators":[]}}`))
	if err := NewValidatorSet(srv.URL, time.Hour).Start(context.Background()); err == nil {
		t.Fatal("expected error when the validator set has no signers")
	}
}

func TestValidatorSetStartBadStatusFails(t *testing.T) {
	srv := newFixtureServer(t, http.StatusInternalServerError, []byte(`{}`))
	if err := NewValidatorSet(srv.URL, time.Hour).Start(context.Background()); err == nil {
		t.Fatal("expected error on a non-2xx status")
	}
}

func TestValidatorSetRefreshKeepsPreviousSetOnError(t *testing.T) {
	good := readFixture(t)

	// First request succeeds, later requests fail.
	var calls int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls == 1 {
			_, _ = w.Write(good)
			return
		}
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	t.Cleanup(srv.Close)

	s := NewValidatorSet(srv.URL, time.Hour)
	if err := s.Start(context.Background()); err != nil {
		t.Fatalf("Start: %v", err)
	}
	// A failed refresh must not wipe the previously loaded set.
	if _, ferr := s.fetch(context.Background()); ferr == nil {
		t.Fatal("expected the refresh fetch to fail")
	}
	if !s.HasSigner(knownSigner) {
		t.Error("expected the previous signer set to be retained after a failed refresh")
	}
}

func TestValidatorSetRefreshLoopExitsOnCancel(t *testing.T) {
	srv := newFixtureServer(t, http.StatusOK, readFixture(t))

	s := NewValidatorSet(srv.URL, time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	if err := s.Start(ctx); err != nil {
		t.Fatalf("Start: %v", err)
	}

	done := make(chan struct{})
	go func() {
		s.refreshLoop(ctx)
		close(done)
	}()

	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("refreshLoop did not exit on context cancellation")
	}
}

func TestConnsRecoverSigner(t *testing.T) {
	// nil validator set => validation disabled: everything is "known".
	if _, known, err := (&Conns{}).RecoverSigner(&types.Header{}); !known || err != nil {
		t.Errorf("nil validator set: want known=true err=nil, got known=%v err=%v", known, err)
	}

	// Build a header signed by a generated key.
	priv, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("generating key: %v", err)
	}
	signer := crypto.PubkeyToAddress(priv.PublicKey)
	header := &types.Header{
		Number:     big.NewInt(1),
		Difficulty: big.NewInt(1),
		Extra:      make([]byte, crypto.SignatureLength),
	}
	sig, err := crypto.Sign(clique.SealHash(header).Bytes(), priv)
	if err != nil {
		t.Fatalf("signing header: %v", err)
	}
	copy(header.Extra[len(header.Extra)-crypto.SignatureLength:], sig)

	// Signer in the set => known, and the recovered address is returned.
	inSet := &Conns{validators: &ValidatorSet{}}
	inSet.validators.signers.Set(map[common.Address]struct{}{signer: {}})
	if addr, known, err := inSet.RecoverSigner(header); !known || err != nil || addr != signer {
		t.Errorf("known signer: want known=true err=nil addr=%s, got known=%v err=%v addr=%s", signer, known, err, addr)
	}

	// Signer not in the set => not known, but the address is still recovered.
	notInSet := &Conns{validators: &ValidatorSet{}}
	notInSet.validators.signers.Set(map[common.Address]struct{}{
		common.HexToAddress("0x1111111111111111111111111111111111111111"): {},
	})
	if addr, known, err := notInSet.RecoverSigner(header); known || err != nil || addr != signer {
		t.Errorf("unknown signer: want known=false err=nil addr=%s, got known=%v err=%v addr=%s", signer, known, err, addr)
	}

	// Header with no recoverable signature => not known, with an error.
	if _, known, err := notInSet.RecoverSigner(&types.Header{}); known || err == nil {
		t.Errorf("unrecoverable header: want known=false err!=nil, got known=%v err=%v", known, err)
	}
}
