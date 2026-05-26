//go:build heimdall_integration

package validator

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
// runs a validator subcommand. Each call re-constructs the umbrella to
// avoid subcommand double-registration.
func execLive(t *testing.T, args ...string) (string, string, error) {
	t.Helper()

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
	all := []string{
		"--rest-url", liveREST(),
		"--rpc-url", liveRPC(),
		"--timeout", "10",
	}
	all = append(all, args...)
	root.SetArgs(all)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := root.ExecuteContext(ctx)
	return stdout.String(), stderr.String(), err
}

func TestIntegrationTotalPower(t *testing.T) {
	stdout, _, err := execLive(t, "validator", "total-power")
	if err != nil {
		t.Fatalf("total-power: %v", err)
	}
	n, perr := strconv.ParseUint(strings.TrimSpace(stdout), 10, 64)
	if perr != nil {
		t.Fatalf("total-power not an integer: %q (%v)", stdout, perr)
	}
	if n == 0 {
		t.Errorf("expected total-power > 0, got %d", n)
	}
}

func TestIntegrationSetLimit(t *testing.T) {
	stdout, _, err := execLive(t, "validator", "set", "--limit", "5")
	if err != nil {
		t.Fatalf("set --limit 5: %v", err)
	}
	lines := strings.Split(strings.TrimRight(stdout, "\n"), "\n")
	// header + at most 5 data rows
	if len(lines) == 0 {
		t.Fatalf("no output for set --limit 5")
	}
	// Subtract 1 for the header row.
	dataRows := len(lines) - 1
	if dataRows > 5 {
		t.Errorf("set --limit 5 returned %d data rows, expected <= 5", dataRows)
	}
	if dataRows == 0 {
		t.Errorf("set --limit 5 returned 0 data rows")
	}
}

// TestIntegrationSignerRoundTrip asserts that the signer of the first
// validator from `set` round-trips to the same id via
// `signer <addr>`.
func TestIntegrationSignerRoundTrip(t *testing.T) {
	// Fetch as JSON so we can parse deterministically.
	stdout, _, err := execLive(t, "validator", "set", "--limit", "1", "--json")
	if err != nil {
		t.Fatalf("set --limit 1 --json: %v", err)
	}
	var resp struct {
		ValidatorSet struct {
			Validators []struct {
				ValID  string `json:"val_id"`
				Signer string `json:"signer"`
			} `json:"validators"`
		} `json:"validator_set"`
	}
	if jerr := json.Unmarshal([]byte(stdout), &resp); jerr != nil {
		t.Fatalf("unmarshal: %v\n%s", jerr, stdout)
	}
	if len(resp.ValidatorSet.Validators) == 0 {
		t.Fatal("no validators returned")
	}
	first := resp.ValidatorSet.Validators[0]
	if first.Signer == "" || first.ValID == "" {
		t.Fatalf("first validator missing signer/id: %+v", first)
	}
	// Round-trip.
	sigStdout, _, err := execLive(t, "validator", "signer", first.Signer, "--json")
	if err != nil {
		t.Fatalf("signer %s: %v", first.Signer, err)
	}
	var sigResp struct {
		Validator struct {
			ValID  string `json:"val_id"`
			Signer string `json:"signer"`
		} `json:"validator"`
	}
	if jerr := json.Unmarshal([]byte(sigStdout), &sigResp); jerr != nil {
		t.Fatalf("signer unmarshal: %v\n%s", jerr, sigStdout)
	}
	if sigResp.Validator.ValID != first.ValID {
		t.Errorf("round-trip val_id mismatch: set=%q signer=%q", first.ValID, sigResp.Validator.ValID)
	}
	if !strings.EqualFold(sigResp.Validator.Signer, first.Signer) {
		t.Errorf("round-trip signer mismatch: set=%q signer=%q", first.Signer, sigResp.Validator.Signer)
	}
}

func TestIntegrationProposer(t *testing.T) {
	stdout, _, err := execLive(t, "validator", "proposer")
	if err != nil {
		t.Fatalf("proposer: %v", err)
	}
	if !strings.Contains(stdout, "signer") {
		t.Errorf("proposer output missing signer: %q", stdout)
	}
}

func TestIntegrationGetRoundTrip(t *testing.T) {
	// Grab the current proposer (which exposes a val_id) and round-trip
	// through `get`.
	stdout, _, err := execLive(t, "validator", "proposer", "--json")
	if err != nil {
		t.Fatalf("proposer --json: %v", err)
	}
	var resp struct {
		Validator struct {
			ValID string `json:"val_id"`
		} `json:"validator"`
	}
	if jerr := json.Unmarshal([]byte(stdout), &resp); jerr != nil {
		t.Fatalf("unmarshal: %v\n%s", jerr, stdout)
	}
	if resp.Validator.ValID == "" {
		t.Fatal("proposer missing val_id")
	}
	getStdout, _, err := execLive(t, "validator", "get", resp.Validator.ValID)
	if err != nil {
		t.Fatalf("get %s: %v", resp.Validator.ValID, err)
	}
	if !strings.Contains(getStdout, resp.Validator.ValID) {
		t.Errorf("get output missing val_id=%q: %s", resp.Validator.ValID, getStdout)
	}
}
