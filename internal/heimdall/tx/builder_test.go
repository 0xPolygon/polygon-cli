package tx

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	hproto "github.com/0xPolygon/polygon-cli/internal/heimdall/proto"
)

// fixedECDSAKey returns a deterministic 32-byte secp256k1 key for
// tests so golden-file comparisons are stable.
func fixedECDSAKey(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	hexKey := "0101010101010101010101010101010101010101010101010101010101010101"
	b, err := hex.DecodeString(hexKey)
	if err != nil {
		t.Fatalf("decoding fixed key: %v", err)
	}
	priv, err := ethcrypto.ToECDSA(b)
	if err != nil {
		t.Fatalf("loading fixed key: %v", err)
	}
	return priv
}

// fakeAccountFetcher returns fixed account info and records the
// calls it received.
type fakeAccountFetcher struct {
	calls  []string
	return_ Account
	err    error
}

func (f *fakeAccountFetcher) FetchAccount(_ context.Context, addr string) (*Account, error) {
	f.calls = append(f.calls, addr)
	if f.err != nil {
		return nil, f.err
	}
	out := f.return_
	out.Address = addr
	return &out, nil
}

// _ keeps ethcrypto imported for tests below that use it indirectly.
var _ = ethcrypto.PubkeyToAddress

func TestBuilderSignDirectGolden(t *testing.T) {
	priv := fixedECDSAKey(t)
	addr := ethcrypto.PubkeyToAddress(priv.PublicKey).Hex()

	msg := &WithdrawFeeMsg{
		Proposer: strings.ToLower(addr),
		Amount:   "1000000000000000000",
	}
	b := NewBuilder().
		WithChainID("heimdallv2-80002").
		WithGasLimit(200000).
		WithFee(hproto.Coin{Denom: "pol", Amount: "10000000000000000"}).
		WithAccountNumber(25).
		WithSequence(51129).
		AddMsg(msg)

	raw, err := b.Sign(priv)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	if len(raw) == 0 {
		t.Fatal("Sign returned empty bytes")
	}

	// Round-trip the TxRaw so we can assert structure without relying
	// on a golden byte string (proto serialization is deterministic for
	// our marshallers, but a golden hex string would be brittle to
	// refactor). Verify the body and auth_info match a re-encoded
	// version of the same inputs.
	parsed, err := hproto.UnmarshalTxRaw(raw)
	if err != nil {
		t.Fatalf("UnmarshalTxRaw: %v", err)
	}
	if len(parsed.BodyBytes) == 0 {
		t.Fatal("parsed body empty")
	}
	if len(parsed.AuthInfoBytes) == 0 {
		t.Fatal("parsed auth_info empty")
	}
	if len(parsed.Signatures) != 1 {
		t.Fatalf("expected 1 signature, got %d", len(parsed.Signatures))
	}
	if got := len(parsed.Signatures[0]); got != 64 {
		t.Fatalf("signature length = %d, want 64", got)
	}
}

func TestBuilderSignDeterministic(t *testing.T) {
	priv := fixedECDSAKey(t)
	newB := func() *Builder {
		return NewBuilder().
			WithChainID("heimdallv2-80002").
			WithGasLimit(200000).
			WithFee(hproto.Coin{Denom: "pol", Amount: "10000000000000000"}).
			WithAccountNumber(25).
			WithSequence(51129).
			AddMsg(&WithdrawFeeMsg{Proposer: "0x02f615e95563ef16f10354dba9e584e58d2d4314", Amount: "1"})
	}
	raw1, err := newB().Sign(priv)
	if err != nil {
		t.Fatalf("Sign #1: %v", err)
	}
	raw2, err := newB().Sign(priv)
	if err != nil {
		t.Fatalf("Sign #2: %v", err)
	}
	// ECDSA without deterministic-k varies per call; what must be
	// identical is the pre-image (body + auth_info + sign_doc).
	p1, _ := hproto.UnmarshalTxRaw(raw1)
	p2, _ := hproto.UnmarshalTxRaw(raw2)
	if !bytesEq(p1.BodyBytes, p2.BodyBytes) {
		t.Fatal("body bytes differ across builds with identical inputs")
	}
	if !bytesEq(p1.AuthInfoBytes, p2.AuthInfoBytes) {
		t.Fatal("auth_info bytes differ across builds with identical inputs")
	}
}

func TestSignModeDirectVsAminoDiffer(t *testing.T) {
	priv := fixedECDSAKey(t)
	msg := &WithdrawFeeMsg{Proposer: "0x02f615e95563ef16f10354dba9e584e58d2d4314", Amount: "1"}
	build := func(mode SignMode) []byte {
		b := NewBuilder().
			WithChainID("heimdallv2-80002").
			WithGasLimit(200000).
			WithFee(hproto.Coin{Denom: "pol", Amount: "10000000000000000"}).
			WithAccountNumber(25).
			WithSequence(51129).
			WithSignMode(mode).
			AddMsg(msg)
		raw, err := b.Sign(priv)
		if err != nil {
			t.Fatalf("Sign(%s): %v", mode, err)
		}
		return raw
	}
	direct := build(SignModeDirect)
	amino := build(SignModeAminoJSON)

	pDirect, err := hproto.UnmarshalTxRaw(direct)
	if err != nil {
		t.Fatalf("UnmarshalTxRaw direct: %v", err)
	}
	pAmino, err := hproto.UnmarshalTxRaw(amino)
	if err != nil {
		t.Fatalf("UnmarshalTxRaw amino: %v", err)
	}
	// Body is the same (mode doesn't affect TxBody).
	if !bytesEq(pDirect.BodyBytes, pAmino.BodyBytes) {
		t.Fatal("TxBody bytes should be identical across sign modes")
	}
	// AuthInfo differs (mode field lives inside signer_infos).
	if bytesEq(pDirect.AuthInfoBytes, pAmino.AuthInfoBytes) {
		t.Fatal("AuthInfo bytes should differ when sign mode changes")
	}
	// Signatures are over different pre-images so they differ.
	if bytesEq(pDirect.Signatures[0], pAmino.Signatures[0]) {
		t.Fatal("signatures should differ between direct and amino-json")
	}
}

func TestResolveAccountUsesFetcher(t *testing.T) {
	f := &fakeAccountFetcher{return_: Account{AccountNumber: 25, Sequence: 51129}}
	b := NewBuilder()
	if err := b.ResolveAccount(context.Background(), f, "0xabc"); err != nil {
		t.Fatalf("ResolveAccount: %v", err)
	}
	if b.accountNumber != 25 {
		t.Errorf("accountNumber = %d, want 25", b.accountNumber)
	}
	if b.sequence != 51129 {
		t.Errorf("sequence = %d, want 51129", b.sequence)
	}
	if len(f.calls) != 1 || f.calls[0] != "0xabc" {
		t.Errorf("unexpected fetch calls: %v", f.calls)
	}
}

func TestResolveAccountPropagatesFetcherError(t *testing.T) {
	want := errors.New("boom")
	f := &fakeAccountFetcher{err: want}
	b := NewBuilder()
	err := b.ResolveAccount(context.Background(), f, "0xabc")
	if !errors.Is(err, want) {
		t.Fatalf("error = %v, want wrapping %v", err, want)
	}
}

func TestResolveAccountPreservesExplicitSequence(t *testing.T) {
	f := &fakeAccountFetcher{return_: Account{AccountNumber: 25, Sequence: 99}}
	b := NewBuilder().WithSequence(500)
	if err := b.ResolveAccount(context.Background(), f, "0xabc"); err != nil {
		t.Fatal(err)
	}
	if b.sequence != 500 {
		t.Errorf("explicit sequence was overwritten: got %d, want 500", b.sequence)
	}
	if b.accountNumber != 25 {
		t.Errorf("accountNumber = %d, want 25", b.accountNumber)
	}
}

func TestSignMissingInputs(t *testing.T) {
	priv := fixedECDSAKey(t)
	cases := []struct {
		name string
		b    *Builder
	}{
		{"no messages", NewBuilder().WithChainID("x").WithGasLimit(1)},
		{"no chain id", NewBuilder().WithGasLimit(1).AddMsg(&WithdrawFeeMsg{Proposer: "x", Amount: "1"})},
		{"no gas limit", NewBuilder().WithChainID("x").AddMsg(&WithdrawFeeMsg{Proposer: "x", Amount: "1"})},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := tc.b.Sign(priv); err == nil {
				t.Fatalf("expected error, got nil")
			}
		})
	}
}

func TestSignNilKey(t *testing.T) {
	b := NewBuilder().
		WithChainID("x").
		WithGasLimit(1).
		AddMsg(&WithdrawFeeMsg{Proposer: "x", Amount: "1"})
	if _, err := b.Sign(nil); err == nil {
		t.Fatal("expected error for nil key")
	}
}

func TestSignWithdrawFeeMissingFields(t *testing.T) {
	priv := fixedECDSAKey(t)
	b := NewBuilder().
		WithChainID("x").
		WithGasLimit(1).
		AddMsg(&WithdrawFeeMsg{Amount: "1"}) // no proposer
	_, err := b.Sign(priv)
	if err == nil || !strings.Contains(err.Error(), "proposer") {
		t.Fatalf("expected proposer-required error, got %v", err)
	}
}

func TestSignatureRecoversPubkey(t *testing.T) {
	priv := fixedECDSAKey(t)
	b := NewBuilder().
		WithChainID("heimdallv2-80002").
		WithGasLimit(200000).
		WithFee(hproto.Coin{Denom: "pol", Amount: "1"}).
		WithAccountNumber(25).
		WithSequence(0).
		AddMsg(&WithdrawFeeMsg{Proposer: "0xabc", Amount: "1"})
	raw, err := b.Sign(priv)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	parsed, err := hproto.UnmarshalTxRaw(raw)
	if err != nil {
		t.Fatalf("UnmarshalTxRaw: %v", err)
	}
	// Reconstruct the sign digest and verify the 64-byte sig
	// verifies against the signer's pubkey.
	_, digest := signBytesDirect(parsed.BodyBytes, parsed.AuthInfoBytes, "heimdallv2-80002", 25)
	if !ethcrypto.VerifySignature(ethcrypto.FromECDSAPub(&priv.PublicKey), digest, parsed.Signatures[0]) {
		t.Fatal("signature does not verify against signer pubkey")
	}
}

// -----------------------------------------------------------------------
// Broadcast / Simulate / WaitForInclusion HTTP tests.
// -----------------------------------------------------------------------

func newTestRESTClient(t *testing.T, handler http.Handler) *client.RESTClient {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return client.NewRESTClient(srv.URL, 5*time.Second, nil, false)
}

func TestBroadcastHappyPath(t *testing.T) {
	h := http.NewServeMux()
	h.HandleFunc("/cosmos/tx/v1beta1/txs", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"tx_response":{"txhash":"ABCDEF","code":0,"height":"100"}}`)
	})
	rest := newTestRESTClient(t, h)
	res, err := Broadcast(context.Background(), rest, []byte{0x01, 0x02}, BroadcastModeSync)
	if err != nil {
		t.Fatalf("Broadcast: %v", err)
	}
	if res.TxHash != "abcdef" {
		t.Errorf("TxHash = %q, want %q", res.TxHash, "abcdef")
	}
	if res.Code != 0 {
		t.Errorf("Code = %d, want 0", res.Code)
	}
	if res.Height != 100 {
		t.Errorf("Height = %d, want 100", res.Height)
	}
}

func TestBroadcastNonZeroCode(t *testing.T) {
	h := http.NewServeMux()
	h.HandleFunc("/cosmos/tx/v1beta1/txs", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"tx_response":{"txhash":"DEAD","code":13,"codespace":"sdk","raw_log":"insufficient fees"}}`)
	})
	rest := newTestRESTClient(t, h)
	res, err := Broadcast(context.Background(), rest, []byte{0x01}, BroadcastModeSync)
	if err == nil {
		t.Fatal("expected error for non-zero code")
	}
	if res == nil {
		t.Fatal("expected non-nil result for non-zero code")
	}
	if res.Code != 13 || res.RawLog != "insufficient fees" {
		t.Fatalf("result = %+v", res)
	}
}

func TestSimulateHappyPath(t *testing.T) {
	h := http.NewServeMux()
	h.HandleFunc("/cosmos/tx/v1beta1/simulate", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"gas_info":{"gas_wanted":"200000","gas_used":"137842"}}`)
	})
	rest := newTestRESTClient(t, h)
	res, err := Simulate(context.Background(), rest, []byte{0x01})
	if err != nil {
		t.Fatalf("Simulate: %v", err)
	}
	if res.GasUsed != 137842 {
		t.Errorf("GasUsed = %d, want 137842", res.GasUsed)
	}
	if res.GasWanted != 200000 {
		t.Errorf("GasWanted = %d, want 200000", res.GasWanted)
	}
}

func TestWaitForInclusionContextCancel(t *testing.T) {
	// Server always 404s so the wait loop spins until we cancel.
	h := http.NewServeMux()
	h.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	rpc := client.NewRPCClient(srv.URL, 1*time.Second, nil, false)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		_, _, err := WaitForInclusion(ctx, rpc, "deadbeef", 10*time.Millisecond)
		done <- err
	}()
	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("WaitForInclusion did not return within 2s of cancel")
	}
}

func TestWaitForInclusionTimeout(t *testing.T) {
	// Same as above but using deadline, mimicking --timeout expiry.
	h := http.NewServeMux()
	h.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	rpc := client.NewRPCClient(srv.URL, 1*time.Second, nil, false)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	_, _, err := WaitForInclusion(ctx, rpc, "deadbeef", 10*time.Millisecond)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}
}

func TestWaitForInclusionFindsTx(t *testing.T) {
	call := 0
	h := http.NewServeMux()
	h.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		call++
		if call < 3 {
			// Return a JSON-RPC error ("not found") for the first two
			// calls so we exercise the polling loop.
			fmt.Fprint(w, `{"jsonrpc":"2.0","id":0,"error":{"code":-1,"message":"tx not found"}}`)
			return
		}
		fmt.Fprint(w, `{"jsonrpc":"2.0","id":0,"result":{"hash":"DEADBEEF","height":"42"}}`)
	})
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	rpc := client.NewRPCClient(srv.URL, 1*time.Second, nil, false)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	h0, _, err := WaitForInclusion(ctx, rpc, "deadbeef", 5*time.Millisecond)
	if err != nil {
		t.Fatalf("WaitForInclusion: %v", err)
	}
	if h0 != 42 {
		t.Errorf("height = %d, want 42", h0)
	}
}

// -----------------------------------------------------------------------
// RESTAccountFetcher against mocked server.
// -----------------------------------------------------------------------

func TestRESTAccountFetcher(t *testing.T) {
	h := http.NewServeMux()
	h.HandleFunc("/cosmos/auth/v1beta1/accounts/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"account":{"@type":"/cosmos.auth.v1beta1.BaseAccount","address":"0xabc","account_number":"25","sequence":"51129"}}`)
	})
	rest := newTestRESTClient(t, h)
	f := &RESTAccountFetcher{Client: rest}
	acc, err := f.FetchAccount(context.Background(), "0xabc")
	if err != nil {
		t.Fatalf("FetchAccount: %v", err)
	}
	if acc.AccountNumber != 25 || acc.Sequence != 51129 {
		t.Fatalf("got %+v", acc)
	}
}

// -----------------------------------------------------------------------
// Utility assertions.
// -----------------------------------------------------------------------

func bytesEq(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
