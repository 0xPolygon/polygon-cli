package ops

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
)

// fixturePath returns the absolute path to a testdata JSON file.
func fixturePath(t *testing.T, name string) string {
	t.Helper()
	_, thisFile, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(thisFile), "testdata", name)
}

// loadFixture reads testdata/<name>.
func loadFixture(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(fixturePath(t, name))
	if err != nil {
		t.Fatalf("reading fixture %s: %v", name, err)
	}
	return b
}

// newTestServer routes CometBFT JSON-RPC methods to canned fixture
// bytes. If the map value is a fully-formed envelope (has a "jsonrpc"
// key), the server rewrites its "id" to match the request.
func newTestServer(t *testing.T, routes map[string][]byte) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req struct {
			Method string         `json:"method"`
			Params map[string]any `json:"params"`
			ID     uint64         `json:"id"`
		}
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		data, ok := routes[req.Method]
		if !ok {
			http.Error(w, "no route for "+req.Method, 404)
			return
		}
		var envelope map[string]any
		_ = json.Unmarshal(data, &envelope)
		envelope["id"] = req.ID
		out, _ := json.Marshal(envelope)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(out)
	}))
	t.Cleanup(srv.Close)
	return srv
}

// runCmd dispatches args against a fresh heimdall root with ops
// registered, returning stdout, stderr, and any error.
func runCmd(t *testing.T, srv *httptest.Server, args ...string) (string, string, error) {
	t.Helper()
	root := &cobra.Command{Use: "h", SilenceUsage: true}
	f := &config.Flags{}
	f.Register(root)
	Register(root, f)
	var stdout, stderr bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stderr)
	full := append([]string{
		"--rest-url", srv.URL,
		"--rpc-url", srv.URL,
	}, args...)
	root.SetArgs(full)
	err := root.ExecuteContext(context.Background())
	return stdout.String(), stderr.String(), err
}

// ---- status ----

func TestStatusSummary(t *testing.T) {
	srv := newTestServer(t, map[string][]byte{
		"status": loadFixture(t, "status.json"),
	})
	stdout, _, err := runCmd(t, srv, "ops", "status")
	if err != nil {
		t.Fatalf("status: %v", err)
	}
	for _, want := range []string{"node_id", "moniker", "network", "latest_block_height", "catching_up", "validator_address"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("missing %q in output:\n%s", want, stdout)
		}
	}
	if !strings.Contains(stdout, "heimdallv2-80002") {
		t.Errorf("expected network in output:\n%s", stdout)
	}
}

func TestStatusJSON(t *testing.T) {
	srv := newTestServer(t, map[string][]byte{
		"status": loadFixture(t, "status.json"),
	})
	stdout, _, err := runCmd(t, srv, "--json", "ops", "status")
	if err != nil {
		t.Fatalf("status --json: %v", err)
	}
	var v any
	if err := json.Unmarshal([]byte(stdout), &v); err != nil {
		t.Fatalf("output not json: %v\n%s", err, stdout)
	}
}

func TestStatusField(t *testing.T) {
	srv := newTestServer(t, map[string][]byte{
		"status": loadFixture(t, "status.json"),
	})
	stdout, _, err := runCmd(t, srv, "ops", "status", "--field", "network")
	if err != nil {
		t.Fatalf("status --field: %v", err)
	}
	if strings.TrimSpace(stdout) != "heimdallv2-80002" {
		t.Fatalf("field output = %q, want heimdallv2-80002", stdout)
	}
}

// ---- health ----

func TestHealthOK(t *testing.T) {
	srv := newTestServer(t, map[string][]byte{
		"health": loadFixture(t, "health.json"),
	})
	stdout, _, err := runCmd(t, srv, "ops", "health")
	if err != nil {
		t.Fatalf("health: %v", err)
	}
	if strings.TrimSpace(stdout) != "OK" {
		t.Fatalf("health output = %q, want OK", stdout)
	}
}

func TestHealthRPCError(t *testing.T) {
	// Route returns a JSON-RPC error envelope.
	errorEnv, _ := json.Marshal(map[string]any{
		"jsonrpc": "2.0", "id": 1,
		"error": map[string]any{"code": -32603, "message": "Internal error"},
	})
	srv := newTestServer(t, map[string][]byte{
		"health": errorEnv,
	})
	_, _, err := runCmd(t, srv, "ops", "health")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if client.ExitCode(err) == 0 {
		t.Fatalf("expected non-zero exit code, got 0 for err %v", err)
	}
}

// ---- peers ----

func TestPeersDefault(t *testing.T) {
	srv := newTestServer(t, map[string][]byte{
		"net_info": loadFixture(t, "net_info.json"),
	})
	stdout, _, err := runCmd(t, srv, "ops", "peers")
	if err != nil {
		t.Fatalf("peers: %v", err)
	}
	// Fixture from Amoy has 10 peers.
	if !strings.Contains(stdout, "n_peers") {
		t.Errorf("missing n_peers:\n%s", stdout)
	}
	for _, col := range []string{"node_id", "remote_ip", "moniker"} {
		if !strings.Contains(stdout, col) {
			t.Errorf("missing column %q:\n%s", col, stdout)
		}
	}
}

func TestPeersVerbose(t *testing.T) {
	srv := newTestServer(t, map[string][]byte{
		"net_info": loadFixture(t, "net_info.json"),
	})
	stdout, _, err := runCmd(t, srv, "ops", "peers", "--verbose")
	if err != nil {
		t.Fatalf("peers --verbose: %v", err)
	}
	// Verbose emits JSON; must be parseable.
	var v any
	if err := json.Unmarshal([]byte(stdout), &v); err != nil {
		t.Fatalf("verbose output not JSON: %v\n%s", err, stdout)
	}
}

// ---- consensus ----

func TestConsensusSummary(t *testing.T) {
	srv := newTestServer(t, map[string][]byte{
		"dump_consensus_state": loadFixture(t, "dump_consensus_state.json"),
	})
	stdout, stderr, err := runCmd(t, srv, "ops", "consensus")
	if err != nil {
		t.Fatalf("consensus: %v", err)
	}
	for _, want := range []string{"height", "round", "step", "proposer_address"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("missing %q:\n%s", want, stdout)
		}
	}
	if !strings.Contains(stderr, "expensive") {
		t.Errorf("expected cost warning on stderr:\n%s", stderr)
	}
}

func TestConsensusEndpointDisabled(t *testing.T) {
	// The live Amoy node returns a JSON-RPC error when consensus
	// endpoints are disabled; we must surface it as an error (non-zero
	// exit) rather than silently swallowing it.
	srv := newTestServer(t, map[string][]byte{
		"dump_consensus_state": loadFixture(t, "dump_consensus_state_error.json"),
	})
	_, _, err := runCmd(t, srv, "ops", "consensus")
	if err == nil {
		t.Fatal("expected error for disabled consensus endpoint")
	}
	var rpcErr *client.RPCError
	if !errors.As(err, &rpcErr) {
		t.Fatalf("err type %T, want *client.RPCError", err)
	}
}

// ---- tx-pool ----

func TestTxPoolSummary(t *testing.T) {
	srv := newTestServer(t, map[string][]byte{
		"num_unconfirmed_txs": loadFixture(t, "num_unconfirmed_txs.json"),
	})
	stdout, _, err := runCmd(t, srv, "ops", "tx-pool")
	if err != nil {
		t.Fatalf("tx-pool: %v", err)
	}
	for _, want := range []string{"n_txs", "total_bytes"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("missing %q:\n%s", want, stdout)
		}
	}
}

func TestTxPoolListEmpty(t *testing.T) {
	srv := newTestServer(t, map[string][]byte{
		"unconfirmed_txs": loadFixture(t, "unconfirmed_txs.json"),
	})
	stdout, _, err := runCmd(t, srv, "ops", "tx-pool", "--list", "--limit", "5")
	if err != nil {
		t.Fatalf("tx-pool --list: %v", err)
	}
	// Empty fixture: only headers, no hash lines.
	if strings.Count(stdout, "0x") != 0 {
		t.Errorf("expected no tx hashes for empty pool:\n%s", stdout)
	}
}

func TestTxPoolListWithTxs(t *testing.T) {
	// Build a fixture with two canned base64 txs.
	env := map[string]any{
		"jsonrpc": "2.0", "id": 1,
		"result": map[string]any{
			"n_txs": "2", "total": "2", "total_bytes": "10",
			"txs": []string{"aGVsbG8=", "d29ybGQ="}, // "hello", "world"
		},
	}
	body, _ := json.Marshal(env)
	srv := newTestServer(t, map[string][]byte{"unconfirmed_txs": body})
	stdout, _, err := runCmd(t, srv, "ops", "tx-pool", "--list", "--limit", "2")
	if err != nil {
		t.Fatalf("tx-pool --list: %v", err)
	}
	// Expect two "0x" hash lines.
	lines := strings.Split(strings.TrimSpace(stdout), "\n")
	var hashes []string
	for _, l := range lines {
		if strings.HasPrefix(l, "0x") && len(l) > 10 {
			hashes = append(hashes, l)
		}
	}
	if len(hashes) != 2 {
		t.Fatalf("expected 2 hashes, got %d\n%s", len(hashes), stdout)
	}
	// sha256("hello") = 2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824
	if !strings.EqualFold(hashes[0], "0x2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824") {
		t.Errorf("hash0 = %q", hashes[0])
	}
}

func TestTxPoolBadLimit(t *testing.T) {
	srv := newTestServer(t, map[string][]byte{
		"unconfirmed_txs": loadFixture(t, "unconfirmed_txs.json"),
	})
	_, _, err := runCmd(t, srv, "ops", "tx-pool", "--list", "--limit", "0")
	if err == nil {
		t.Fatal("expected error for --limit 0")
	}
	var uErr *client.UsageError
	if !errors.As(err, &uErr) {
		t.Fatalf("err type = %T, want *UsageError", err)
	}
}

// ---- abci-info ----

func TestABCIInfoSummary(t *testing.T) {
	srv := newTestServer(t, map[string][]byte{
		"abci_info": loadFixture(t, "abci_info.json"),
	})
	stdout, _, err := runCmd(t, srv, "ops", "abci-info")
	if err != nil {
		t.Fatalf("abci-info: %v", err)
	}
	for _, want := range []string{"app", "version", "last_block_height", "last_block_app_hash"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("missing %q:\n%s", want, stdout)
		}
	}
}

// ---- commit ----

func TestCommitLatest(t *testing.T) {
	srv := newTestServer(t, map[string][]byte{
		"commit": loadFixture(t, "commit.json"),
	})
	stdout, _, err := runCmd(t, srv, "ops", "commit")
	if err != nil {
		t.Fatalf("commit: %v", err)
	}
	for _, want := range []string{"chain_id", "height", "proposer_address", "app_hash", "block_hash"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("missing %q:\n%s", want, stdout)
		}
	}
}

func TestCommitExplicitHeight(t *testing.T) {
	// Verify the server sees a concrete height param.
	var seenHeight any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req struct {
			Method string         `json:"method"`
			Params map[string]any `json:"params"`
			ID     uint64         `json:"id"`
		}
		_ = json.Unmarshal(body, &req)
		seenHeight = req.Params["height"]
		// Reuse the latest fixture as the response.
		b, _ := os.ReadFile(fixturePath(t, "commit.json"))
		var env map[string]any
		_ = json.Unmarshal(b, &env)
		env["id"] = req.ID
		out, _ := json.Marshal(env)
		_, _ = w.Write(out)
	}))
	t.Cleanup(srv.Close)
	_, _, err := runCmd(t, srv, "ops", "commit", "12345")
	if err != nil {
		t.Fatalf("commit HEIGHT: %v", err)
	}
	if seenHeight != "12345" {
		t.Fatalf("height param = %v, want \"12345\"", seenHeight)
	}
}

func TestCommitBadHeight(t *testing.T) {
	srv := newTestServer(t, map[string][]byte{
		"commit": loadFixture(t, "commit.json"),
	})
	_, _, err := runCmd(t, srv, "ops", "commit", "notanumber")
	if err == nil {
		t.Fatal("expected error for bogus height")
	}
	var uErr *client.UsageError
	if !errors.As(err, &uErr) {
		t.Fatalf("err type = %T, want *UsageError", err)
	}
}

// ---- validators-cometbft ----

func TestValidatorsCometBFT(t *testing.T) {
	srv := newTestServer(t, map[string][]byte{
		"validators": loadFixture(t, "validators.json"),
	})
	stdout, stderr, err := runCmd(t, srv, "ops", "validators-cometbft")
	if err != nil {
		t.Fatalf("validators-cometbft: %v", err)
	}
	// Header lines.
	for _, want := range []string{"block_height", "count", "address", "voting_power"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("missing %q:\n%s", want, stdout)
		}
	}
	// Must route operators to the stake command for the real set.
	if !strings.Contains(stderr, "heimdall validator") {
		t.Errorf("expected stderr hint mentioning 'heimdall validator':\n%s", stderr)
	}
}

func TestValidatorsCometBFTJSON(t *testing.T) {
	srv := newTestServer(t, map[string][]byte{
		"validators": loadFixture(t, "validators.json"),
	})
	stdout, _, err := runCmd(t, srv, "--json", "ops", "validators-cometbft")
	if err != nil {
		t.Fatalf("validators-cometbft --json: %v", err)
	}
	var v any
	if err := json.Unmarshal([]byte(stdout), &v); err != nil {
		t.Fatalf("output not json: %v\n%s", err, stdout)
	}
}

// ---- helper function unit tests ----

func TestTxHashFromBase64Known(t *testing.T) {
	// sha256("hello") known-good digest.
	got, err := txHashFromBase64("aGVsbG8=")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	want := "0x2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	if !strings.EqualFold(got, want) {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestTxHashFromBase64Bad(t *testing.T) {
	if _, err := txHashFromBase64("not base64 !!!"); err == nil {
		t.Fatal("expected error on bad base64")
	}
}
