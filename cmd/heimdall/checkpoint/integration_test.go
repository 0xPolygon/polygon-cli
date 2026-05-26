//go:build heimdall_integration

package checkpoint

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
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
// runs `checkpoint …`. Each call re-constructs the umbrella to avoid
// subcommand double-registration.
func execLive(t *testing.T, args ...string) (string, string, error) {
	t.Helper()

	local := &cobra.Command{
		Use:     "checkpoint [ID]",
		Aliases: []string{"cp"},
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}
			return runGet(cmd, args[0])
		},
	}
	local.AddCommand(
		newParamsCmd(),
		newCountCmd(),
		newLatestCmd(),
		newGetCmd(),
		newBufferCmd(),
		newLastNoAckCmd(),
		newNextCmd(),
		newListCmd(),
		newSignaturesCmd(),
		newOverviewCmd(),
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
		"checkpoint",
	}
	all = append(all, args...)
	root.SetArgs(all)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := root.ExecuteContext(ctx)
	return stdout.String(), stderr.String(), err
}

func TestIntegrationCheckpointCount(t *testing.T) {
	stdout, _, err := execLive(t, "count")
	if err != nil {
		t.Fatalf("count: %v", err)
	}
	n, perr := strconv.ParseUint(strings.TrimSpace(stdout), 10, 64)
	if perr != nil {
		t.Fatalf("count not an integer: %q (%v)", stdout, perr)
	}
	if n == 0 {
		t.Errorf("expected count > 0, got %d", n)
	}
}

func TestIntegrationCheckpointRoundTrip(t *testing.T) {
	// Fetch count to pick a valid id.
	stdout, _, err := execLive(t, "count")
	if err != nil {
		t.Fatalf("count: %v", err)
	}
	count, perr := strconv.ParseUint(strings.TrimSpace(stdout), 10, 64)
	if perr != nil || count == 0 {
		t.Fatalf("count parse failed: %q (%v)", stdout, perr)
	}
	id := strconv.FormatUint(count, 10)
	getStdout, _, err := execLive(t, "get", id)
	if err != nil {
		t.Fatalf("get %s: %v", id, err)
	}
	if !strings.Contains(getStdout, id) {
		t.Errorf("get %s output missing id: %q", id, getStdout)
	}
	if !strings.Contains(getStdout, "root_hash") {
		t.Errorf("get %s missing root_hash field: %q", id, getStdout)
	}
}

func TestIntegrationCheckpointLatest(t *testing.T) {
	stdout, _, err := execLive(t, "latest")
	if err != nil {
		t.Fatalf("latest: %v", err)
	}
	if !strings.Contains(stdout, "root_hash") {
		t.Errorf("latest missing root_hash: %q", stdout)
	}
	if !strings.Contains(stdout, "0x") {
		t.Errorf("latest root_hash not hex-normalized: %q", stdout)
	}
}

func TestIntegrationCheckpointBuffer(t *testing.T) {
	stdout, _, err := execLive(t, "buffer")
	if err != nil {
		t.Fatalf("buffer: %v", err)
	}
	// Either populated (has `proposer`) or the empty form.
	if !strings.Contains(stdout, "proposer") && !strings.Contains(stdout, "empty") {
		t.Errorf("buffer output neither populated nor empty: %q", stdout)
	}
}

func TestIntegrationCheckpointList(t *testing.T) {
	stdout, _, err := execLive(t, "list", "--limit", "3")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	// Rough shape check — at least the headers.
	if !strings.Contains(stdout, "id") || !strings.Contains(stdout, "root_hash") {
		t.Errorf("list missing expected columns: %q", stdout)
	}
}

func TestIntegrationCheckpointSignatures(t *testing.T) {
	// The live Amoy node's tx index does not contain checkpoint/ack
	// txs (verified at fixture-capture time). Use tx_search for a
	// known-indexed action (topup) to obtain a valid 32-byte hash and
	// assert that signatures responds coherently. The REST endpoint
	// either returns a payload or a gRPC error envelope — both count.
	hash := pickRecentIndexedTxHash(t)
	if hash == "" {
		t.Skip("no indexed tx available on live node")
	}
	stdout, _, err := execLive(t, "signatures", hash)
	// Invalid hash (from topup) is expected to produce an HTTPError or
	// a JSON envelope with `code: 5` (not set); we just require the
	// command not to panic and to emit something recognisable.
	if err != nil && !strings.Contains(err.Error(), "signatures") && !strings.Contains(err.Error(), "HTTP") {
		t.Fatalf("signatures %s: unexpected error: %v", hash, err)
	}
	_ = stdout
}

// pickRecentIndexedTxHash queries CometBFT /tx_search for a recent
// topup tx (a reliably indexed action on Amoy). Returns "" on miss.
func pickRecentIndexedTxHash(t *testing.T) string {
	t.Helper()
	httpC := &http.Client{Timeout: 20 * time.Second}
	searchBody := `{"jsonrpc":"2.0","id":1,"method":"tx_search","params":{"query":"message.action='/heimdallv2.topup.MsgTopupTx'","per_page":"1","page":"1","order_by":"desc"}}`
	req, _ := http.NewRequest(http.MethodPost, liveRPC(), strings.NewReader(searchBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpC.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	buf, _ := io.ReadAll(resp.Body)
	var out struct {
		Result struct {
			Txs []struct {
				Hash string `json:"hash"`
			} `json:"txs"`
		} `json:"result"`
	}
	if json.Unmarshal(buf, &out) != nil || len(out.Result.Txs) == 0 {
		return ""
	}
	return out.Result.Txs[0].Hash
}

func TestIntegrationCheckpointOverview(t *testing.T) {
	stdout, _, err := execLive(t, "overview")
	if err != nil {
		t.Fatalf("overview: %v", err)
	}
	var v any
	if jerr := json.Unmarshal([]byte(stdout), &v); jerr != nil {
		t.Fatalf("overview output not JSON: %v\n%s", jerr, stdout)
	}
	if !strings.Contains(stdout, "ack_count") {
		t.Errorf("overview missing ack_count: %q", stdout)
	}
}
