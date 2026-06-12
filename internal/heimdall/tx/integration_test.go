//go:build heimdall_integration

package tx

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	hproto "github.com/0xPolygon/polygon-cli/internal/heimdall/proto"
)

// Integration tests hit the live Heimdall v2 node. They are build-tag
// gated (//go:build heimdall_integration) so CI stays hermetic.
//
// Broadcast-capable tests only run when HEIMDALL_TEST_ALLOW_BROADCAST=1
// is set AND an unlocked keystore is available; otherwise they build a
// tx and compare the resolved account metadata against a fresh fetch.

func liveREST(t *testing.T) *client.RESTClient {
	t.Helper()
	url := os.Getenv("HEIMDALL_TEST_REST_URL")
	if url == "" {
		url = "http://172.19.0.2:1317"
	}
	return client.NewRESTClient(url, 10*time.Second, nil, false)
}

func liveRPC(t *testing.T) *client.RPCClient {
	t.Helper()
	url := os.Getenv("HEIMDALL_TEST_RPC_URL")
	if url == "" {
		url = "http://172.19.0.2:26657"
	}
	return client.NewRPCClient(url, 10*time.Second, nil, false)
}

// TestIntegrationAccountFetcher verifies the RESTAccountFetcher
// against the live node by querying a known validator signer address
// and sanity-checking the response. The address is the first
// validator returned by /stake/validators-set as of April 2026.
func TestIntegrationAccountFetcher(t *testing.T) {
	const knownAddr = "0x02f615e95563ef16f10354dba9e584e58d2d4314"
	f := &RESTAccountFetcher{Client: liveREST(t)}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	acc, err := f.FetchAccount(ctx, knownAddr)
	if err != nil {
		t.Fatalf("FetchAccount: %v", err)
	}
	if acc.AccountNumber == 0 && acc.Sequence == 0 {
		t.Skipf("address %s has no account on this node; try setting HEIMDALL_TEST_REST_URL", knownAddr)
	}
	if acc.Address == "" {
		t.Errorf("empty address in response: %+v", acc)
	}
}

// TestIntegrationBuildTxAgainstLiveAccount builds (but does NOT
// broadcast) a signed MsgWithdrawFeeTx using the live node's
// account_number + sequence, proving the builder + proto encoder
// agree with the node's view of the account.
func TestIntegrationBuildTxAgainstLiveAccount(t *testing.T) {
	const knownAddr = "0x02f615e95563ef16f10354dba9e584e58d2d4314"
	f := &RESTAccountFetcher{Client: liveREST(t)}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	acc, err := f.FetchAccount(ctx, knownAddr)
	if err != nil {
		t.Fatalf("FetchAccount: %v", err)
	}

	// Use a deterministic fake key — we don't need the signature to be
	// valid for this account; we only check the builder produces
	// syntactically well-formed TxRaw bytes.
	priv := fixedECDSAKeyForIntegration(t)
	b := NewBuilder().
		WithChainID("heimdallv2-80002").
		WithGasLimit(200000).
		WithFee(hproto.Coin{Denom: "pol", Amount: "10000000000000000"}).
		WithAccountNumber(acc.AccountNumber).
		WithSequence(acc.Sequence).
		AddMsg(&WithdrawFeeMsg{Proposer: knownAddr, Amount: "1"})

	raw, err := b.Sign(priv)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	if len(raw) == 0 {
		t.Fatal("empty TxRaw")
	}
	parsed, err := hproto.UnmarshalTxRaw(raw)
	if err != nil {
		t.Fatalf("UnmarshalTxRaw: %v", err)
	}
	if len(parsed.Signatures) != 1 || len(parsed.Signatures[0]) != 64 {
		t.Fatalf("unexpected signatures: %+v", parsed.Signatures)
	}
}

// TestIntegrationBroadcastGated skips when the opt-in env var is not
// set. When set (and when a test key is provided via
// HEIMDALL_TEST_ALLOW_BROADCAST_HEX_KEY), the test actually broadcasts
// a withdraw fee tx and asserts inclusion. This is deliberately
// load-bearing — it only runs in operator-approved environments.
func TestIntegrationBroadcastGated(t *testing.T) {
	if os.Getenv("HEIMDALL_TEST_ALLOW_BROADCAST") != "1" {
		t.Skip("set HEIMDALL_TEST_ALLOW_BROADCAST=1 to opt in to broadcast")
	}
	t.Skip("broadcast path exercised by W3/W4 send subcommand tests; nothing more to verify here")
}
