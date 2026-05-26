//go:build heimdall_integration

package topup

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
)

// Integration tests talk directly to the live Heimdall v2 node on
// 172.19.0.2 unless overridden via HEIMDALL_TEST_REST_URL /
// HEIMDALL_TEST_RPC_URL.

func liveRPC() string {
	if v := os.Getenv("HEIMDALL_TEST_RPC_URL"); v != "" {
		return v
	}
	return "http://172.19.0.2:26657"
}

func liveREST() string {
	if v := os.Getenv("HEIMDALL_TEST_REST_URL"); v != "" {
		return v
	}
	return "http://172.19.0.2:1317"
}

// execLive spins up a fresh cobra root wired to the live Heimdall and
// runs `topup …`.
func execLive(t *testing.T, args ...string) (string, string, error) {
	t.Helper()

	local := &cobra.Command{
		Use:   "topup",
		Short: "topup",
		Args:  cobra.NoArgs,
	}
	local.AddCommand(
		newRootCmd(),
		newAccountCmd(),
		newProofCmd(),
		newVerifyCmd(),
		newSequenceCmd(),
		newIsOldCmd(),
	)

	root := &cobra.Command{Use: "h", SilenceUsage: true}
	f := &config.Flags{}
	f.Register(root)
	flags = f
	root.AddCommand(local)

	var stdout, stderr bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stderr)
	all := []string{
		"--rest-url", liveREST(),
		"--rpc-url", liveRPC(),
		"--timeout", "10",
		"topup",
	}
	all = append(all, args...)
	root.SetArgs(all)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	err := root.ExecuteContext(ctx)
	return stdout.String(), stderr.String(), err
}

// TestIntegrationRootHash asserts that `topup root` returns a sane
// 32-byte (64-hex-char) 0x-hex digest.
func TestIntegrationRootHash(t *testing.T) {
	stdout, _, err := execLive(t, "root")
	if err != nil {
		t.Fatalf("root: %v", err)
	}
	got := strings.TrimSpace(stdout)
	if !strings.HasPrefix(got, "0x") {
		t.Errorf("root: expected 0x-hex, got %q", got)
	}
	if len(got) != 2+64 {
		t.Errorf("root: expected 32-byte hex (66 chars), got %d chars: %q", len(got), got)
	}
}

// TestIntegrationRootJSON asserts that --json wrapping works end to end.
func TestIntegrationRootJSON(t *testing.T) {
	stdout, _, err := execLive(t, "root", "--json")
	if err != nil {
		t.Fatalf("root --json: %v", err)
	}
	var m map[string]any
	if jerr := json.Unmarshal([]byte(stdout), &m); jerr != nil {
		t.Fatalf("root --json not valid JSON: %v\n%s", jerr, stdout)
	}
	if _, ok := m["account_root_hash"]; !ok {
		t.Errorf("missing account_root_hash: %v", m)
	}
}

// TestIntegrationProofL1Unconfigured asserts that the live test node
// (which lacks L1 RPC) surfaces the L1-not-configured hint when asked
// for an account proof. If the live node ever gains L1 connectivity
// the test is skipped rather than failed — we still want to exercise
// the request shape.
func TestIntegrationProofL1Unconfigured(t *testing.T) {
	_, stderr, err := execLive(t, "proof", "0x02f615e95563ef16f10354dba9e584e58d2d4314")
	if err == nil {
		t.Skip("live node has L1 RPC configured; skipping L1-unconfigured hint assertion")
	}
	if !strings.Contains(stderr, "eth_rpc_url") {
		t.Errorf("expected L1-not-configured hint, got stderr=%q", stderr)
	}
}

// TestIntegrationSequenceL1Unconfigured asserts the same hint surfaces
// for the sequence endpoint.
func TestIntegrationSequenceL1Unconfigured(t *testing.T) {
	_, stderr, err := execLive(t,
		"sequence",
		"0x48bd44a37ff84c7cc584e5df4bf43bbda6116d5708d41080e7b9b030195c6bf8",
		"423")
	if err == nil {
		t.Skip("live node has L1 RPC configured; skipping L1-unconfigured hint assertion")
	}
	if !strings.Contains(stderr, "eth_rpc_url") {
		t.Errorf("expected L1-not-configured hint, got stderr=%q", stderr)
	}
}

// TestIntegrationIsOldL1Unconfigured asserts the same hint surfaces
// for is-old-tx.
func TestIntegrationIsOldL1Unconfigured(t *testing.T) {
	_, stderr, err := execLive(t,
		"is-old",
		"0x48bd44a37ff84c7cc584e5df4bf43bbda6116d5708d41080e7b9b030195c6bf8",
		"423")
	if err == nil {
		t.Skip("live node has L1 RPC configured; skipping L1-unconfigured hint assertion")
	}
	if !strings.Contains(stderr, "eth_rpc_url") {
		t.Errorf("expected L1-not-configured hint, got stderr=%q", stderr)
	}
}

// TestIntegrationVerifyBadProof asserts verify returns a sane error
// for a too-short proof (server-side validation), exercising the
// non-L1 path through the command.
func TestIntegrationVerifyBadProof(t *testing.T) {
	// 2 bytes — server requires multiples of 32 bytes.
	_, _, err := execLive(t,
		"verify",
		"0x02f615e95563ef16f10354dba9e584e58d2d4314",
		"deadbeef")
	if err == nil {
		t.Fatal("expected server to reject short proof")
	}
}

// TestIntegrationVerifyZeroProof asserts verify returns a bool for a
// well-formed but never-actually-valid proof of 32 zero bytes.
func TestIntegrationVerifyZeroProof(t *testing.T) {
	stdout, _, err := execLive(t,
		"verify",
		"0x02f615e95563ef16f10354dba9e584e58d2d4314",
		"0x0000000000000000000000000000000000000000000000000000000000000000")
	if err != nil {
		t.Fatalf("verify zero: %v", err)
	}
	got := strings.TrimSpace(stdout)
	if got != "true" && got != "false" {
		t.Errorf("verify: expected bool, got %q", got)
	}
}

// TestIntegrationAccountNotFoundExits1 exercises the 404/5xx path for
// a well-formed address that has no dividend account on chain.
func TestIntegrationAccountNotFoundExits1(t *testing.T) {
	_, _, err := execLive(t, "account", "0x0000000000000000000000000000000000000000")
	if err == nil {
		t.Skip("zero address unexpectedly has a dividend account; skipping not-found assertion")
	}
}
