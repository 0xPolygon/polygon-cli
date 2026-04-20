package tx

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

	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
)

// --- fixture loaders ---

func testdataPath(t *testing.T, subdir, name string) string {
	t.Helper()
	_, thisFile, _, _ := runtime.Caller(0)
	// cmd/heimdall/tx/<thisFile> -> ../../../internal/heimdall/client/testdata/<subdir>
	base := filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "internal", "heimdall", "client", "testdata", subdir)
	return filepath.Join(base, name)
}

func loadFixture(t *testing.T, subdir, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(testdataPath(t, subdir, name))
	if err != nil {
		t.Fatalf("reading fixture %s/%s: %v", subdir, name, err)
	}
	return b
}

// --- fixture RPC server ---

// newRPCFixtureServer returns a test server that routes CometBFT RPC
// method names to canned fixture bodies. Each fixture is the full
// JSON-RPC envelope (jsonrpc/id/result|error); the id is rewritten to
// the request id so the client's decoder matches.
func newRPCFixtureServer(t *testing.T, routes map[string][]byte) *httptest.Server {
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

// newRESTFixtureServer returns a test server that routes URL paths
// to canned REST JSON bodies. Query strings are ignored by default so
// a single fixture can back multiple queries.
func newRESTFixtureServer(t *testing.T, routes map[string][]byte) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, ok := routes[r.URL.Path]
		if !ok {
			http.Error(w, "no route for "+r.URL.Path, 404)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}))
	t.Cleanup(srv.Close)
	return srv
}

// runCmd assembles a root heimdall cobra command wired to the given
// REST + RPC URLs (either or both may be empty) and runs the argv.
// Returns stdout, stderr, and any error from Execute.
func runCmd(t *testing.T, restURL, rpcURL string, args ...string) (string, string, error) {
	t.Helper()
	root := &cobra.Command{Use: "h", SilenceUsage: true}
	f := &config.Flags{}
	f.Register(root)
	Register(root, f)
	var stdout, stderr bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stderr)
	all := []string{}
	if restURL != "" {
		all = append(all, "--rest-url", restURL)
	}
	if rpcURL != "" {
		all = append(all, "--rpc-url", rpcURL)
	}
	all = append(all, args...)
	root.SetArgs(all)
	err := root.ExecuteContext(context.Background())
	return stdout.String(), stderr.String(), err
}

// errorsAs is a thin wrapper so tests can dodge importing errors
// twice (linter noise on repeated imports in large test suites).
func errorsAs(err error, target any) bool { return errors.As(err, target) }

// mustContain fails the test with a helpful message if substr is not
// in s.
func mustContain(t *testing.T, s, substr string) {
	t.Helper()
	if !strings.Contains(s, substr) {
		t.Fatalf("expected %q in output:\n%s", substr, s)
	}
}
