package checkpoint

import (
	"bytes"
	"context"
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

// testdataPath resolves a path under internal/heimdall/client/testdata/
// from this test file's location.
func testdataPath(t *testing.T, subdir, name string) string {
	t.Helper()
	_, thisFile, _, _ := runtime.Caller(0)
	// cmd/heimdall/checkpoint/<thisFile> -> ../../../internal/heimdall/client/testdata/<subdir>
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

// restRoute is a single route entry for newRESTFixtureServer.
type restRoute struct {
	status int
	body   []byte
}

// newRESTFixtureServer routes request paths to canned bodies. A route
// with status==0 is treated as a 200.
func newRESTFixtureServer(t *testing.T, routes map[string]restRoute) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route, ok := routes[r.URL.Path]
		if !ok {
			http.Error(w, "no route for "+r.URL.Path, 404)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		status := route.status
		if status == 0 {
			status = 200
		}
		w.WriteHeader(status)
		_, _ = w.Write(route.body)
	}))
	t.Cleanup(srv.Close)
	return srv
}

// runCmd assembles a fresh root command with the checkpoint umbrella
// wired in, using the given REST URL, and executes argv. Each call
// creates new cobra command instances so tests don't share state.
func runCmd(t *testing.T, restURL string, args ...string) (string, string, error) {
	t.Helper()

	// Build a fresh umbrella each run. We can't reuse CheckpointCmd
	// directly because subsequent calls to Register would double-add
	// its children.
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
	// Re-bind the package-level flags so subcommand RunE resolves config.
	flags = f
	root.AddCommand(local)

	var stdout, stderr bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stderr)
	all := []string{}
	if restURL != "" {
		all = append(all, "--rest-url", restURL, "--rpc-url", restURL)
	}
	all = append(all, "checkpoint")
	all = append(all, args...)
	root.SetArgs(all)
	err := root.ExecuteContext(context.Background())
	return stdout.String(), stderr.String(), err
}

// mustContain fails the test if substr isn't in s.
func mustContain(t *testing.T, s, substr string) {
	t.Helper()
	if !strings.Contains(s, substr) {
		t.Fatalf("expected %q in output:\n%s", substr, s)
	}
}

// mustNotContain fails the test if substr is in s.
func mustNotContain(t *testing.T, s, substr string) {
	t.Helper()
	if strings.Contains(s, substr) {
		t.Fatalf("expected %q NOT in output:\n%s", substr, s)
	}
}
