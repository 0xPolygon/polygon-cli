//go:build heimdall_integration

package ops

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
)

// Integration tests run against the live Amoy-backed node at
// 172.19.0.2:26657 unless overridden by HEIMDALL_TEST_RPC_URL.

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

func execLive(t *testing.T, args ...string) (string, string, error) {
	t.Helper()
	root := &cobra.Command{Use: "h", SilenceUsage: true}
	f := &config.Flags{}
	f.Register(root)
	Register(root, f)
	var stdout, stderr bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stderr)
	all := append([]string{
		"--rest-url", liveREST(),
		"--rpc-url", liveRPC(),
		"--timeout", "10",
	}, args...)
	root.SetArgs(all)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := root.ExecuteContext(ctx)
	return stdout.String(), stderr.String(), err
}

func TestIntegrationStatus(t *testing.T) {
	stdout, _, err := execLive(t, "ops", "status")
	if err != nil {
		t.Fatalf("status: %v", err)
	}
	for _, want := range []string{"node_id", "moniker", "latest_block_height"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("missing %q:\n%s", want, stdout)
		}
	}
}

func TestIntegrationHealth(t *testing.T) {
	stdout, _, err := execLive(t, "ops", "health")
	if err != nil {
		t.Fatalf("health: %v", err)
	}
	if strings.TrimSpace(stdout) != "OK" {
		t.Fatalf("health = %q, want OK", stdout)
	}
}

func TestIntegrationPeers(t *testing.T) {
	stdout, _, err := execLive(t, "ops", "peers")
	if err != nil {
		t.Fatalf("peers: %v", err)
	}
	if !strings.Contains(stdout, "n_peers") {
		t.Errorf("missing n_peers:\n%s", stdout)
	}
}

func TestIntegrationABCIInfo(t *testing.T) {
	stdout, _, err := execLive(t, "ops", "abci-info")
	if err != nil {
		t.Fatalf("abci-info: %v", err)
	}
	for _, want := range []string{"app", "last_block_height", "last_block_app_hash"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("missing %q:\n%s", want, stdout)
		}
	}
}

func TestIntegrationCommitLatest(t *testing.T) {
	stdout, _, err := execLive(t, "ops", "commit")
	if err != nil {
		t.Fatalf("commit: %v", err)
	}
	for _, want := range []string{"chain_id", "height", "block_hash"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("missing %q:\n%s", want, stdout)
		}
	}
}

func TestIntegrationValidatorsCometBFT(t *testing.T) {
	stdout, stderr, err := execLive(t, "ops", "validators-cometbft")
	if err != nil {
		t.Fatalf("validators-cometbft: %v", err)
	}
	if !strings.Contains(stdout, "block_height") {
		t.Errorf("missing block_height:\n%s", stdout)
	}
	if !strings.Contains(stderr, "heimdall validator") {
		t.Errorf("missing stderr hint:\n%s", stderr)
	}
}

func TestIntegrationTxPool(t *testing.T) {
	stdout, _, err := execLive(t, "ops", "tx-pool")
	if err != nil {
		t.Fatalf("tx-pool: %v", err)
	}
	if !strings.Contains(stdout, "n_txs") {
		t.Errorf("missing n_txs:\n%s", stdout)
	}
}

func TestIntegrationTxPoolList(t *testing.T) {
	// Pool is almost always empty on Amoy; we just verify it doesn't
	// error out and prints the expected headers.
	stdout, _, err := execLive(t, "ops", "tx-pool", "--list", "--limit", "5")
	if err != nil {
		t.Fatalf("tx-pool --list: %v", err)
	}
	if !strings.Contains(stdout, "n_txs") {
		t.Errorf("missing n_txs:\n%s", stdout)
	}
}

func TestIntegrationStatusJSON(t *testing.T) {
	stdout, _, err := execLive(t, "--json", "ops", "status")
	if err != nil {
		t.Fatalf("status --json: %v", err)
	}
	var v any
	if err := json.Unmarshal([]byte(stdout), &v); err != nil {
		t.Fatalf("output not JSON: %v\n%s", err, stdout)
	}
}

func TestIntegrationCommitHeight(t *testing.T) {
	// First fetch status to pick a known-valid recent height.
	stdout, _, err := execLive(t, "ops", "status", "--field", "latest_block_height")
	if err != nil {
		t.Fatalf("status: %v", err)
	}
	latest, err := strconv.ParseInt(strings.TrimSpace(stdout), 10, 64)
	if err != nil {
		t.Fatalf("bad height %q: %v", stdout, err)
	}
	// Pick a height five blocks behind tip to dodge the brief window
	// before the commit is indexed.
	target := strconv.FormatInt(latest-5, 10)
	out, _, err := execLive(t, "ops", "commit", target)
	if err != nil {
		t.Fatalf("commit %s: %v", target, err)
	}
	if !strings.Contains(out, target) {
		t.Errorf("expected height %s in output:\n%s", target, out)
	}
}

func TestIntegrationConsensusDisabledReportsCleanly(t *testing.T) {
	// Amoy public endpoints typically disable consensus dumps. We
	// treat a clean JSON-RPC error as a successful test of the error
	// path: the command should surface the node-side disable message
	// rather than crashing.
	_, _, err := execLive(t, "ops", "consensus")
	if err == nil {
		t.Skip("consensus endpoints are enabled on this node; nothing to assert")
	}
	var rpcErr *client.RPCError
	if !errors.As(err, &rpcErr) {
		t.Fatalf("consensus error %T = %v, want *client.RPCError", err, err)
	}
}
