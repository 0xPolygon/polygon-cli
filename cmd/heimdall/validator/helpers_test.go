package validator

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
	// When set, the server will assert that the incoming query string
	// contains these values.
	wantQuery map[string]string
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
		for k, want := range route.wantQuery {
			if got := r.URL.Query().Get(k); got != want {
				http.Error(w, "bad query value for "+k, 400)
				return
			}
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

// runCmd assembles a fresh root command with the validator umbrella
// wired in, using the given REST URL, and executes argv. Each call
// creates new cobra command instances so tests don't share state.
func runCmd(t *testing.T, restURL string, args ...string) (string, string, error) {
	t.Helper()
	return runCmdNamed(t, restURL, "validator", args...)
}

// runCmdNamed is like runCmd but lets the caller pick between the
// `validator` umbrella and the top-level `validators` alias.
func runCmdNamed(t *testing.T, restURL, topLevel string, args ...string) (string, string, error) {
	t.Helper()

	// Reset set flags between runs so --limit/--sort from one test do
	// not bleed into another.
	setFlags.sort = "power"
	setFlags.limit = 0
	setFlags.fields = nil

	local := &cobra.Command{
		Use:     "validator [ID]",
		Aliases: []string{"val"},
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}
			return runGet(cmd, args[0])
		},
	}
	local.AddCommand(
		newSetCmd(),
		newTotalPowerCmd(),
		newGetCmd(),
		newSignerCmd(),
		newStatusCmd(),
		newProposerCmd(),
		newProposersCmd(),
		newIsOldStakeTxCmd(),
	)

	validators := &cobra.Command{
		Use:  "validators",
		Args: cobra.NoArgs,
		RunE: runSet,
	}
	attachSetFlags(validators.Flags())

	root := &cobra.Command{Use: "h", SilenceUsage: true}
	f := &config.Flags{}
	f.Register(root)
	flags = f
	root.AddCommand(local, validators)

	var stdout, stderr bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stderr)
	all := []string{}
	if restURL != "" {
		all = append(all, "--rest-url", restURL, "--rpc-url", restURL)
	}
	all = append(all, topLevel)
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

func mustNotContain(t *testing.T, s, substr string) {
	t.Helper()
	if strings.Contains(s, substr) {
		t.Fatalf("expected %q NOT in output:\n%s", substr, s)
	}
}
