package chainparams

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
)

// testdataPath resolves a path under cmd/heimdall/chainparams/testdata/
// from this test file's location.
func testdataPath(t *testing.T, name string) string {
	t.Helper()
	_, thisFile, _, _ := runtime.Caller(0)
	base := filepath.Join(filepath.Dir(thisFile), "testdata")
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

// restRoute is a single route entry for newRESTFixtureServer.
type restRoute struct {
	status int
	body   []byte
	// match allows a route to inspect query parameters.
	match func(url.Values) bool
}

// newRESTFixtureServer routes request paths to canned bodies. A route
// with status==0 is treated as a 200.
func newRESTFixtureServer(t *testing.T, routes map[string][]restRoute) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		candidates, ok := routes[r.URL.Path]
		if !ok {
			http.Error(w, "no route for "+r.URL.Path, 404)
			return
		}
		for _, route := range candidates {
			if route.match != nil && !route.match(r.URL.Query()) {
				continue
			}
			w.Header().Set("Content-Type", "application/json")
			status := route.status
			if status == 0 {
				status = 200
			}
			w.WriteHeader(status)
			_, _ = w.Write(route.body)
			return
		}
		http.Error(w, "no matching route for "+r.URL.String(), 404)
	}))
	t.Cleanup(srv.Close)
	return srv
}

// runCmd assembles a fresh root command with the chainparams umbrella
// wired in, using the given REST URL, and executes argv. Each call
// creates new cobra command instances so tests don't share state.
func runCmd(t *testing.T, restURL string, args ...string) (string, string, error) {
	t.Helper()

	local := &cobra.Command{
		Use:     "chainmanager",
		Aliases: []string{"cm"},
		Short:   "chainmanager",
		Args:    cobra.NoArgs,
	}
	local.AddCommand(
		newParamsCmd(),
		newAddressesCmd(),
	)

	root := &cobra.Command{Use: "h", SilenceUsage: true}
	f := &config.Flags{}
	f.Register(root)
	flags = f
	root.AddCommand(local)

	var stdout, stderr bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stderr)
	all := []string{}
	if restURL != "" {
		all = append(all, "--rest-url", restURL, "--rpc-url", restURL)
	}
	all = append(all, "chainmanager")
	all = append(all, args...)
	root.SetArgs(all)
	err := root.ExecuteContext(context.Background())
	return stdout.String(), stderr.String(), err
}

func mustContain(t *testing.T, s, substr string) {
	t.Helper()
	if !strings.Contains(s, substr) {
		t.Fatalf("expected %q in output:\n%s", substr, s)
	}
}

// newLocalCmdWithAlias mirrors the alias wiring from Register so alias
// routing (`cm` → `chainmanager`) can be exercised in tests.
func newLocalCmdWithAlias() *cobra.Command {
	local := &cobra.Command{
		Use:     "chainmanager",
		Aliases: []string{"cm"},
		Short:   "chainmanager",
		Args:    cobra.NoArgs,
	}
	local.AddCommand(
		newParamsCmd(),
		newAddressesCmd(),
	)
	return local
}

// runRootWithAlias assembles a fresh root, attaches local, and executes
// argv WITHOUT prepending the literal `chainmanager` token so callers
// can route through the `cm` alias.
func runRootWithAlias(t *testing.T, restURL string, local *cobra.Command, args ...string) (string, error) {
	t.Helper()
	root := &cobra.Command{Use: "h", SilenceUsage: true}
	f := &config.Flags{}
	f.Register(root)
	flags = f
	root.AddCommand(local)

	var stdout, stderr bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stderr)
	all := []string{}
	if restURL != "" {
		all = append(all, "--rest-url", restURL, "--rpc-url", restURL)
	}
	all = append(all, args...)
	root.SetArgs(all)
	err := root.ExecuteContext(context.Background())
	_ = stderr
	return stdout.String(), err
}
