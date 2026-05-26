package chain

import (
	"bytes"
	"context"
	"encoding/json"
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

// testdataPath returns the absolute path to a captured RPC fixture.
// The internal testdata directory is shared across heimdall tests.
func testdataPath(t *testing.T, name string) string {
	t.Helper()
	_, thisFile, _, _ := runtime.Caller(0)
	// cmd/heimdall/chain/<thisFile> -> ../../../internal/heimdall/client/testdata/rpc
	base := filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "internal", "heimdall", "client", "testdata", "rpc")
	return filepath.Join(base, name)
}

func loadFixture(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(testdataPath(t, name))
	if err != nil {
		t.Fatalf("reading fixture %s: %v", name, err)
	}
	return b
}

// newTestServer returns an httptest.Server that routes CometBFT RPC
// methods to canned fixture bytes. The mapping is method -> raw JSON
// response body (the full envelope including the "jsonrpc" key).
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
		// Rewrite the id so the RPCClient's response decoder matches.
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

// newTestCmd sets up the chain subcommand group against the fixture
// server. Returns a cobra.Command with the target subcommand already
// positioned as root so tests can invoke with cmd.SetArgs.
func newTestCmd(t *testing.T, srv *httptest.Server) *cobra.Command {
	t.Helper()
	root := &cobra.Command{Use: "h", SilenceUsage: true}
	f := &config.Flags{}
	f.Register(root)
	Register(root, f)
	// Force the fixture server as both REST and RPC so config.Resolve
	// is happy. `--chain-id` is auto-resolved via the amoy preset.
	root.SetArgs([]string{
		"--rest-url", srv.URL,
		"--rpc-url", srv.URL,
	})
	return root
}

// runCmd executes the given root command with args prepended by the
// persistent URL overrides, returns stdout/stderr.
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

// --- Tests ---

func TestBlockNumberFromFixture(t *testing.T) {
	srv := newTestServer(t, map[string][]byte{
		"status": loadFixture(t, "status.json"),
	})
	stdout, _, err := runCmd(t, srv, "block-number")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	got := strings.TrimSpace(stdout)
	if got != "32620626" {
		t.Fatalf("block-number = %q, want %q", got, "32620626")
	}
}

func TestChainIDFromFixture(t *testing.T) {
	srv := newTestServer(t, map[string][]byte{
		"status": loadFixture(t, "status.json"),
	})
	stdout, _, err := runCmd(t, srv, "chain-id")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if strings.TrimSpace(stdout) != "heimdallv2-80002" {
		t.Fatalf("chain-id = %q", stdout)
	}
}

func TestChainNameKnown(t *testing.T) {
	srv := newTestServer(t, map[string][]byte{
		"status": loadFixture(t, "status.json"),
	})
	stdout, _, err := runCmd(t, srv, "chain")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if strings.TrimSpace(stdout) != "Polygon Amoy Testnet" {
		t.Fatalf("chain = %q", stdout)
	}
}

func TestChainNameUnknown(t *testing.T) {
	var status map[string]any
	_ = json.Unmarshal(loadFixture(t, "status.json"), &status)
	result := status["result"].(map[string]any)
	result["node_info"].(map[string]any)["network"] = "heimdallv2-999999"
	altered, _ := json.Marshal(status)
	srv := newTestServer(t, map[string][]byte{"status": altered})
	stdout, _, err := runCmd(t, srv, "chain")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !strings.HasPrefix(strings.TrimSpace(stdout), "unknown chain heimdallv2-999999") {
		t.Fatalf("chain = %q, want unknown-chain prefix", stdout)
	}
}

func TestClientVersions(t *testing.T) {
	srv := newTestServer(t, map[string][]byte{
		"abci_info": loadFixture(t, "abci_info.json"),
		"status":    loadFixture(t, "status.json"),
	})
	stdout, _, err := runCmd(t, srv, "client")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	out := stdout
	for _, want := range []string{"cometbft_version", "heimdall_app", "heimdall_version"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q: %q", want, out)
		}
	}
	if !strings.Contains(out, "0.38.19") {
		t.Errorf("output missing cometbft 0.38.19: %q", out)
	}
}

func TestBlockDefaultSummary(t *testing.T) {
	srv := newTestServer(t, map[string][]byte{
		"block": loadFixture(t, "block_latest.json"),
	})
	stdout, _, err := runCmd(t, srv, "block")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	for _, want := range []string{"chain_id", "height", "time", "proposer", "num_txs"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("output missing key %q: %q", want, stdout)
		}
	}
	if !strings.Contains(stdout, "32620627") {
		t.Errorf("expected block height 32620627 in output: %q", stdout)
	}
	if !strings.Contains(stdout, "0xB4D5335E0D89F4666B824BA098F920D83264A69A") {
		t.Errorf("expected 0x-prefixed proposer in output: %q", stdout)
	}
	// Default summary omits the full tx array (which is rendered as
	// a key literally named "txs"). "num_txs" must not trigger a false
	// positive, so anchor the check to a line starting with "txs ".
	for _, line := range strings.Split(stdout, "\n") {
		if strings.HasPrefix(line, "txs ") || strings.HasPrefix(line, "txs\t") {
			t.Errorf("default output should not include txs list: %q", stdout)
		}
	}
}

func TestBlockFullIncludesTxs(t *testing.T) {
	srv := newTestServer(t, map[string][]byte{
		"block": loadFixture(t, "block_latest.json"),
	})
	stdout, _, err := runCmd(t, srv, "block", "--full")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !strings.Contains(stdout, "txs") {
		t.Errorf("expected txs key in --full output: %q", stdout)
	}
}

func TestBlockFieldPluck(t *testing.T) {
	srv := newTestServer(t, map[string][]byte{
		"block": loadFixture(t, "block_latest.json"),
	})
	stdout, _, err := runCmd(t, srv, "block", "--field", "height")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	got := strings.TrimSpace(stdout)
	if got != "32620627" {
		t.Fatalf("field=height output = %q, want 32620627", got)
	}
}

func TestBlockEarliestViaStatus(t *testing.T) {
	// earliest tag probes /status first; we stub both.
	srv := newTestServer(t, map[string][]byte{
		"status": loadFixture(t, "status.json"),
		"block":  loadFixture(t, "block_earliest.json"),
	})
	stdout, _, err := runCmd(t, srv, "block", "earliest")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !strings.Contains(stdout, "29992725") {
		t.Errorf("expected earliest height 29992725: %q", stdout)
	}
}

func TestBlockRejectsEthereumTags(t *testing.T) {
	srv := newTestServer(t, map[string][]byte{
		"block":  loadFixture(t, "block_latest.json"),
		"status": loadFixture(t, "status.json"),
	})
	for _, tag := range []string{"finalized", "safe", "pending"} {
		_, _, err := runCmd(t, srv, "block", tag)
		if err == nil {
			t.Fatalf("tag %q should have errored", tag)
		}
		var uErr *client.UsageError
		if !errorsAs(err, &uErr) {
			t.Fatalf("tag %q returned %T, want *UsageError", tag, err)
		}
	}
}

func TestBlockInvalidHeight(t *testing.T) {
	srv := newTestServer(t, map[string][]byte{
		"block": loadFixture(t, "block_latest.json"),
	})
	_, _, err := runCmd(t, srv, "block", "notanumber")
	if err == nil {
		t.Fatal("expected error for bogus height")
	}
	var uErr *client.UsageError
	if !errorsAs(err, &uErr) {
		t.Fatalf("error type = %T, want *UsageError", err)
	}
}

func TestAgeRendersTimestamp(t *testing.T) {
	srv := newTestServer(t, map[string][]byte{
		"block": loadFixture(t, "block_latest.json"),
	})
	stdout, _, err := runCmd(t, srv, "age")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	// block_latest fixture is 2026-04-20T15:10:48.60553605Z, which is
	// unix second 1776697848.
	if !strings.Contains(stdout, "1776697848") {
		t.Errorf("expected unix 1776697848 in output: %q", stdout)
	}
	if !strings.Contains(stdout, "2026-04-20") {
		t.Errorf("expected annotated date in output: %q", stdout)
	}
}

func TestParseTimestampFormats(t *testing.T) {
	cases := []struct {
		name   string
		in     string
		wantOK bool
	}{
		{"unix", "1776697848", true},
		{"rfc3339", "2026-04-20T15:10:48Z", true},
		{"rfc3339nano", "2026-04-20T15:10:48.605Z", true},
		{"bogus", "yesterday", false},
		{"empty", "", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := parseTimestamp(c.in)
			gotOK := err == nil
			if gotOK != c.wantOK {
				t.Fatalf("parseTimestamp(%q) ok=%v, want %v (err=%v)", c.in, gotOK, c.wantOK, err)
			}
		})
	}
}

func TestResolveHeightCases(t *testing.T) {
	// No RPC call should be made for latest/bad-tag/bad-int.
	srv := newTestServer(t, map[string][]byte{
		"status": loadFixture(t, "status.json"),
	})
	rpc := client.NewRPCClient(srv.URL, 0, nil, false)

	cases := []struct {
		name    string
		in      string
		want    string
		wantErr bool
	}{
		{"empty is latest", "", "", false},
		{"latest", "latest", "", false},
		{"earliest", "earliest", "29992725", false},
		{"finalized rejected", "finalized", "", true},
		{"pending rejected", "pending", "", true},
		{"positive integer", "123", "123", false},
		{"zero rejected", "0", "", true},
		{"negative rejected", "-1", "", true},
		{"garbage rejected", "abc", "", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := resolveHeight(context.Background(), rpc, c.in)
			if (err != nil) != c.wantErr {
				t.Fatalf("err=%v wantErr=%v", err, c.wantErr)
			}
			if got != c.want {
				t.Errorf("got=%q want=%q", got, c.want)
			}
		})
	}
}

func TestUnixFromRFC3339Nano(t *testing.T) {
	got, err := unixFromRFC3339Nano("2026-04-20T15:10:48.60553605Z")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got != "1776697848" {
		t.Fatalf("got %s, want 1776697848", got)
	}
	if _, err := unixFromRFC3339Nano("not a time"); err == nil {
		t.Fatal("expected error on bad input")
	}
}

// --- fake-server find-block test ---

func TestFindBlockNarrowsToTarget(t *testing.T) {
	// 20 synthetic blocks at 1s intervals starting at t=1000.
	base := int64(1000)
	type blk struct {
		height int64
		time   string
	}
	blocks := make(map[string]blk)
	for h := int64(1); h <= 20; h++ {
		blocks[formatInt(h)] = blk{height: h, time: formatRFC(base + h)}
	}

	status := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"result": map[string]any{
			"node_info": map[string]any{"network": "heimdallv2-80002", "version": "test"},
			"sync_info": map[string]any{
				"earliest_block_height": "1",
				"latest_block_height":   "20",
				"latest_block_time":     formatRFC(base + 20),
				"catching_up":           false,
			},
		},
	}
	statusBody, _ := json.Marshal(status)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req struct {
			Method string         `json:"method"`
			Params map[string]any `json:"params"`
			ID     uint64         `json:"id"`
		}
		_ = json.Unmarshal(body, &req)
		switch req.Method {
		case "status":
			var env map[string]any
			_ = json.Unmarshal(statusBody, &env)
			env["id"] = req.ID
			out, _ := json.Marshal(env)
			_, _ = w.Write(out)
		case "block":
			h, _ := req.Params["height"].(string)
			b, ok := blocks[h]
			if !ok {
				http.Error(w, "no block", 500)
				return
			}
			env := map[string]any{
				"jsonrpc": "2.0", "id": req.ID,
				"result": map[string]any{
					"block_id": map[string]any{"hash": "DEAD"},
					"block": map[string]any{
						"header": map[string]any{
							"chain_id":         "heimdallv2-80002",
							"height":           formatInt(b.height),
							"time":             b.time,
							"proposer_address": "ABCD",
						},
						"data": map[string]any{"txs": []string{}},
					},
				},
			}
			out, _ := json.Marshal(env)
			_, _ = w.Write(out)
		default:
			http.Error(w, "no route", 404)
		}
	}))
	t.Cleanup(srv.Close)

	target := base + 7 // height 7 exactly.
	stdout, _, err := runCmd(t, srv, "find-block", formatInt(target))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if strings.TrimSpace(stdout) != "7" {
		t.Fatalf("got height %q, want 7", stdout)
	}
}

// formatInt / formatRFC are tiny helpers local to tests to avoid
// repeated strconv/time juggling.
func formatInt(n int64) string { return jsonNum(n) }
func formatRFC(unixSec int64) string {
	// UTC, nanoseconds 0, so deterministic.
	return isoTime(unixSec)
}

// errorsAs wraps errors.As without importing "errors" at test top.
func errorsAs(err error, target any) bool {
	return errorsAsImpl(err, target)
}
