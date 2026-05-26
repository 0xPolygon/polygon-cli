//go:build heimdall_integration

package milestone

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
// runs `milestone …`.
func execLive(t *testing.T, args ...string) (string, string, error) {
	t.Helper()

	local := &cobra.Command{
		Use:     "milestone [NUMBER]",
		Aliases: []string{"ms"},
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
		"milestone",
	}
	all = append(all, args...)
	root.SetArgs(all)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	err := root.ExecuteContext(ctx)
	return stdout.String(), stderr.String(), err
}

func TestIntegrationMilestoneParams(t *testing.T) {
	stdout, _, err := execLive(t, "params")
	if err != nil {
		t.Fatalf("params: %v", err)
	}
	if !strings.Contains(stdout, "ff_milestone_threshold") {
		t.Errorf("params missing ff_milestone_threshold: %q", stdout)
	}
}

// TestIntegrationMilestoneCountPositive asserts that the network has
// produced at least one milestone. Stronger than >=0 to catch an
// obvious misparse.
func TestIntegrationMilestoneCountPositive(t *testing.T) {
	stdout, _, err := execLive(t, "count")
	if err != nil {
		t.Fatalf("count: %v", err)
	}
	n, err := strconv.ParseUint(strings.TrimSpace(stdout), 10, 64)
	if err != nil {
		t.Fatalf("count output not an integer: %q (%v)", stdout, err)
	}
	if n == 0 {
		t.Errorf("expected milestone count > 0")
	}
}

// TestIntegrationMilestoneLatestHasHexHash pulls `latest` and verifies
// that the `hash` field has been re-encoded to 0x-prefixed hex (i.e.
// the renderer wiring works against a live response).
func TestIntegrationMilestoneLatestHasHexHash(t *testing.T) {
	stdout, _, err := execLive(t, "latest", "--json")
	if err != nil {
		t.Fatalf("latest: %v", err)
	}
	var env struct {
		Milestone struct {
			Hash        string `json:"hash"`
			MilestoneID string `json:"milestone_id"`
		} `json:"milestone"`
	}
	if jerr := json.Unmarshal([]byte(stdout), &env); jerr != nil {
		t.Fatalf("latest not valid JSON: %v\n%s", jerr, stdout)
	}
	if !strings.HasPrefix(env.Milestone.Hash, "0x") {
		t.Errorf("expected latest.hash to start with 0x, got %q", env.Milestone.Hash)
	}
	if env.Milestone.MilestoneID == "" {
		t.Errorf("expected latest.milestone_id to be non-empty")
	}
}

// TestIntegrationMilestoneByCount exercises `milestone <count>` →
// /milestones/{number}. The count-valued number is guaranteed to exist.
func TestIntegrationMilestoneByCount(t *testing.T) {
	countOut, _, err := execLive(t, "count")
	if err != nil {
		t.Fatalf("count: %v", err)
	}
	count := strings.TrimSpace(countOut)
	if count == "" {
		t.Fatalf("count output empty")
	}
	stdout, _, err := execLive(t, count)
	if err != nil {
		t.Fatalf("milestone %s: %v", count, err)
	}
	// Both number and milestone_id must be rendered.
	if !strings.Contains(stdout, "milestone_id") {
		t.Errorf("expected milestone_id in output: %q", stdout)
	}
	if !strings.Contains(stdout, "number") {
		t.Errorf("expected number label in output: %q", stdout)
	}
	if !strings.Contains(stdout, count) {
		t.Errorf("expected URL-path number %s in output: %q", count, stdout)
	}
}

// TestIntegrationMilestoneNumberNotEqualToMilestoneID proves the
// footgun on live data: the URL-path sequence number is *not* the
// value returned in milestone_id.
func TestIntegrationMilestoneNumberNotEqualToMilestoneID(t *testing.T) {
	countOut, _, err := execLive(t, "count")
	if err != nil {
		t.Fatalf("count: %v", err)
	}
	count := strings.TrimSpace(countOut)
	stdout, _, err := execLive(t, count, "--json")
	if err != nil {
		t.Fatalf("milestone %s --json: %v", count, err)
	}
	var env struct {
		Milestone struct {
			Number      string `json:"number"`
			MilestoneID string `json:"milestone_id"`
		} `json:"milestone"`
	}
	if jerr := json.Unmarshal([]byte(stdout), &env); jerr != nil {
		t.Fatalf("output not JSON: %v\n%s", jerr, stdout)
	}
	if env.Milestone.Number != count {
		t.Errorf("expected number=%s, got %q", count, env.Milestone.Number)
	}
	if env.Milestone.MilestoneID == env.Milestone.Number {
		t.Errorf("number and milestone_id must differ; both were %q", env.Milestone.Number)
	}
}

// TestIntegrationMilestoneZeroIsOutOfRange asserts that `milestone 0`
// triggers the hint.
func TestIntegrationMilestoneZeroIsOutOfRange(t *testing.T) {
	_, stderr, err := execLive(t, "0")
	if err == nil {
		t.Fatal("expected error on milestone 0")
	}
	if !strings.Contains(stderr, "valid range is") {
		t.Errorf("expected out-of-range hint, got stderr=%q", stderr)
	}
}

// TestIntegrationMilestoneCountPlus10kIsOutOfRange asserts that a
// number comfortably above `count` triggers the hint.
func TestIntegrationMilestoneCountPlus10kIsOutOfRange(t *testing.T) {
	countOut, _, err := execLive(t, "count")
	if err != nil {
		t.Fatalf("count: %v", err)
	}
	count, perr := strconv.ParseUint(strings.TrimSpace(countOut), 10, 64)
	if perr != nil {
		t.Fatalf("count output not an integer: %v", perr)
	}
	far := strconv.FormatUint(count+10000, 10)
	_, stderr, err := execLive(t, far)
	if err == nil {
		t.Fatalf("expected error on milestone %s", far)
	}
	if !strings.Contains(stderr, "valid range is") {
		t.Errorf("expected out-of-range hint for %s, got stderr=%q", far, stderr)
	}
}
