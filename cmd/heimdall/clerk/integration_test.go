//go:build heimdall_integration

package clerk

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"strconv"
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
// runs `state-sync …`.
func execLive(t *testing.T, args ...string) (string, string, error) {
	t.Helper()

	local := &cobra.Command{
		Use:     "state-sync [ID]",
		Aliases: []string{"clerk", "ss"},
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}
			return runGet(cmd, args[0], false)
		},
	}
	local.AddCommand(
		newCountCmd(),
		newLatestIDCmd(),
		newGetCmd(),
		newListCmd(),
		newRangeCmd(),
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
		"state-sync",
	}
	all = append(all, args...)
	root.SetArgs(all)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	err := root.ExecuteContext(ctx)
	return stdout.String(), stderr.String(), err
}

// TestIntegrationClerkCountPositive asserts that the network has
// synced at least one state-sync event. Stronger than >=0 to catch
// obvious misparses.
func TestIntegrationClerkCountPositive(t *testing.T) {
	stdout, _, err := execLive(t, "count")
	if err != nil {
		t.Fatalf("count: %v", err)
	}
	n, err := strconv.ParseUint(strings.TrimSpace(stdout), 10, 64)
	if err != nil {
		t.Fatalf("count output not an integer: %q (%v)", stdout, err)
	}
	if n == 0 {
		t.Errorf("expected state-sync count > 0")
	}
}

// TestIntegrationClerkByCount pulls `count` and then fetches the
// matching record, asserting the response shape (non-empty data,
// tx_hash, contract). The count-valued id is guaranteed to exist.
func TestIntegrationClerkByCount(t *testing.T) {
	countOut, _, err := execLive(t, "count")
	if err != nil {
		t.Fatalf("count: %v", err)
	}
	count := strings.TrimSpace(countOut)
	if count == "" {
		t.Fatalf("count output empty")
	}
	stdout, _, err := execLive(t, count, "--json")
	if err != nil {
		t.Fatalf("state-sync %s: %v", count, err)
	}
	var env struct {
		Record struct {
			ID         string `json:"id"`
			Contract   string `json:"contract"`
			Data       string `json:"data"`
			TxHash     string `json:"tx_hash"`
			LogIndex   string `json:"log_index"`
			BorChainID string `json:"bor_chain_id"`
			RecordTime string `json:"record_time"`
		} `json:"record"`
	}
	if jerr := json.Unmarshal([]byte(stdout), &env); jerr != nil {
		t.Fatalf("output not JSON: %v\n%s", jerr, stdout)
	}
	if env.Record.ID != count {
		t.Errorf("expected id=%s, got %q", count, env.Record.ID)
	}
	if env.Record.Data == "" {
		t.Errorf("expected non-empty data, got %q", env.Record.Data)
	}
	// Default JSON output normalizes byte fields to 0x-hex; `data` is
	// in the byte-field allowlist so the rendered value starts with 0x.
	if !strings.HasPrefix(env.Record.Data, "0x") {
		t.Errorf("expected data to start with 0x, got %q", env.Record.Data[:min(16, len(env.Record.Data))])
	}
	if !strings.HasPrefix(env.Record.TxHash, "0x") {
		t.Errorf("expected tx_hash to start with 0x, got %q", env.Record.TxHash)
	}
}

// TestIntegrationClerkListLimit5 asserts that page=1&limit=5 against
// the live /clerk/event-records/list returns at most 5 rows.
func TestIntegrationClerkListLimit5(t *testing.T) {
	stdout, _, err := execLive(t, "list", "--limit", "5", "--json")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	var resp struct {
		EventRecords []map[string]any `json:"event_records"`
	}
	if jerr := json.Unmarshal([]byte(stdout), &resp); jerr != nil {
		t.Fatalf("list output not JSON: %v\n%s", jerr, stdout)
	}
	if len(resp.EventRecords) == 0 {
		t.Errorf("expected at least one event record")
	}
	if len(resp.EventRecords) > 5 {
		t.Errorf("expected at most 5 records, got %d", len(resp.EventRecords))
	}
}

// TestIntegrationClerkLatestIDL1Unconfigured asserts that the live
// test node (which lacks L1 RPC) surfaces the L1-not-configured hint
// when asked for latest-id. If the live node ever gains L1 connectivity
// this test will become a no-op false positive — noted in the package
// docs.
func TestIntegrationClerkLatestIDL1Unconfigured(t *testing.T) {
	_, stderr, err := execLive(t, "latest-id")
	if err == nil {
		// L1 is configured on this node; skip rather than fail.
		t.Skip("live node has L1 RPC configured; skipping L1-unconfigured hint assertion")
	}
	if !strings.Contains(stderr, "eth_rpc_url") {
		t.Errorf("expected L1-not-configured hint, got stderr=%q", stderr)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
