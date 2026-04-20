//go:build heimdall_integration

package span

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
// runs `span …`. Each call re-constructs the umbrella to avoid
// subcommand double-registration.
func execLive(t *testing.T, args ...string) (string, string, error) {
	t.Helper()

	local := &cobra.Command{
		Use:     "span [ID]",
		Aliases: []string{"sp"},
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
		newLatestCmd(),
		newGetCmd(),
		newListCmd(),
		newProducersCmd(),
		newSeedCmd(),
		newVotesCmd(),
		newDowntimeCmd(),
		newScoresCmd(),
		newFindCmd(),
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
		"span",
	}
	all = append(all, args...)
	root.SetArgs(all)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	err := root.ExecuteContext(ctx)
	return stdout.String(), stderr.String(), err
}

func TestIntegrationSpanParams(t *testing.T) {
	stdout, _, err := execLive(t, "params")
	if err != nil {
		t.Fatalf("params: %v", err)
	}
	if !strings.Contains(stdout, "sprint_duration") {
		t.Errorf("params missing sprint_duration: %q", stdout)
	}
}

// TestIntegrationSpanLatest verifies that latest returns a span whose
// end_block is strictly greater than its start_block.
func TestIntegrationSpanLatest(t *testing.T) {
	stdout, _, err := execLive(t, "latest", "--json")
	if err != nil {
		t.Fatalf("latest: %v", err)
	}
	var env struct {
		Span struct {
			ID         string `json:"id"`
			StartBlock string `json:"start_block"`
			EndBlock   string `json:"end_block"`
		} `json:"span"`
	}
	if jerr := json.Unmarshal([]byte(stdout), &env); jerr != nil {
		t.Fatalf("latest not valid JSON: %v\n%s", jerr, stdout)
	}
	start, err := strconv.ParseUint(env.Span.StartBlock, 10, 64)
	if err != nil {
		t.Fatalf("invalid start_block %q: %v", env.Span.StartBlock, err)
	}
	end, err := strconv.ParseUint(env.Span.EndBlock, 10, 64)
	if err != nil {
		t.Fatalf("invalid end_block %q: %v", env.Span.EndBlock, err)
	}
	if end <= start {
		t.Errorf("expected end_block > start_block, got start=%d end=%d", start, end)
	}
}

// TestIntegrationSpanFindAtLatestStart checks that span find correctly
// identifies the current span when given a block just past
// latest.start_block.
func TestIntegrationSpanFindAtLatestStart(t *testing.T) {
	stdout, _, err := execLive(t, "latest", "--json")
	if err != nil {
		t.Fatalf("latest: %v", err)
	}
	var env struct {
		Span struct {
			ID         string `json:"id"`
			StartBlock string `json:"start_block"`
		} `json:"span"`
	}
	if jerr := json.Unmarshal([]byte(stdout), &env); jerr != nil {
		t.Fatalf("latest not JSON: %v", jerr)
	}
	start, err := strconv.ParseUint(env.Span.StartBlock, 10, 64)
	if err != nil {
		t.Fatalf("invalid start_block: %v", err)
	}
	target := strconv.FormatUint(start+1, 10)
	findOut, stderr, err := execLive(t, "find", target)
	if err != nil {
		t.Fatalf("find %s: %v\nstderr=%s", target, err, stderr)
	}
	if !strings.Contains(findOut, env.Span.ID) {
		t.Errorf("find %s output missing span id %s: %q", target, env.Span.ID, findOut)
	}
	if !strings.Contains(stderr, "Veblop") {
		t.Errorf("find stderr missing Veblop caveat: %q", stderr)
	}
}

// TestIntegrationSpanFindBeforeAnySpan verifies that find with block 0
// returns a helpful message rather than panicking. Heimdall's span 0
// may start at block 0 (in which case block 0 is covered) or at a
// later block; either outcome is acceptable, but the command must
// succeed.
func TestIntegrationSpanFindBeforeAnySpan(t *testing.T) {
	stdout, stderr, err := execLive(t, "find", "0")
	if err != nil {
		t.Fatalf("find 0: %v\nstdout=%s\nstderr=%s", err, stdout, stderr)
	}
	if stdout == "" {
		t.Errorf("find 0 produced no stdout")
	}
}

func TestIntegrationSpanScores(t *testing.T) {
	stdout, _, err := execLive(t, "scores")
	if err != nil {
		t.Fatalf("scores: %v", err)
	}
	if !strings.Contains(stdout, "val_id") || !strings.Contains(stdout, "score") {
		t.Errorf("scores missing expected columns: %q", stdout)
	}
}

func TestIntegrationSpanVotesAll(t *testing.T) {
	stdout, _, err := execLive(t, "votes")
	if err != nil {
		t.Fatalf("votes: %v", err)
	}
	var v any
	if jerr := json.Unmarshal([]byte(stdout), &v); jerr != nil {
		t.Fatalf("votes output not JSON: %v", jerr)
	}
}

func TestIntegrationSpanList(t *testing.T) {
	stdout, _, err := execLive(t, "list", "--limit", "3")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if !strings.Contains(stdout, "id") || !strings.Contains(stdout, "start_block") {
		t.Errorf("list missing columns: %q", stdout)
	}
}

func TestIntegrationSpanDowntimeNone(t *testing.T) {
	// A very large producer id is unlikely to have planned downtime.
	stdout, _, err := execLive(t, "downtime", "999999")
	if err != nil {
		t.Fatalf("downtime: %v", err)
	}
	if !strings.Contains(stdout, "none") && !strings.Contains(stdout, "start_block") {
		t.Errorf("downtime 999999 should print 'none' or be populated, got %q", stdout)
	}
}
