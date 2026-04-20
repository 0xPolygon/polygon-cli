//go:build heimdall_integration

package chainparams

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
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
// runs `chainmanager …`.
func execLive(t *testing.T, args ...string) (string, string, error) {
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
	all := []string{
		"--rest-url", liveREST(),
		"--rpc-url", liveRPC(),
		"--timeout", "10",
		"chainmanager",
	}
	all = append(all, args...)
	root.SetArgs(all)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	err := root.ExecuteContext(ctx)
	return stdout.String(), stderr.String(), err
}

// TestIntegrationParamsJSON asserts the live chainmanager response
// carries the expected envelope and the well-known address fields.
// Values themselves can drift across upgrades; structure is what we
// pin.
func TestIntegrationParamsJSON(t *testing.T) {
	stdout, _, err := execLive(t, "params", "--json")
	if err != nil {
		t.Fatalf("params --json: %v", err)
	}
	var m map[string]any
	if jerr := json.Unmarshal([]byte(stdout), &m); jerr != nil {
		t.Fatalf("not valid JSON: %v\n%s", jerr, stdout)
	}
	params, ok := m["params"].(map[string]any)
	if !ok {
		t.Fatalf("missing params object: %v", m)
	}
	chainParams, ok := params["chain_params"].(map[string]any)
	if !ok {
		t.Fatalf("missing chain_params object: %v", params)
	}
	// Required fields per proto/heimdallv2/chainmanager/chainmanager.proto.
	for _, key := range []string{
		"bor_chain_id",
		"heimdall_chain_id",
		"pol_token_address",
		"staking_manager_address",
		"slash_manager_address",
		"root_chain_address",
		"staking_info_address",
		"state_sender_address",
		"state_receiver_address",
		"validator_set_address",
	} {
		v, ok := chainParams[key].(string)
		if !ok || v == "" {
			t.Errorf("chain_params.%s missing or empty: %v", key, chainParams[key])
		}
	}
	// The tx confirmation depths are required and non-empty strings.
	for _, key := range []string{"main_chain_tx_confirmations", "bor_chain_tx_confirmations"} {
		v, ok := params[key].(string)
		if !ok || v == "" {
			t.Errorf("params.%s missing or empty", key)
		}
	}
}

// TestIntegrationParamsField exercises the --field plucker against the
// live node.
func TestIntegrationParamsField(t *testing.T) {
	stdout, _, err := execLive(t,
		"params",
		"--field", "params.chain_params.bor_chain_id")
	if err != nil {
		t.Fatalf("params --field: %v", err)
	}
	got := strings.TrimSpace(stdout)
	if got == "" {
		t.Errorf("empty bor_chain_id")
	}
}

// TestIntegrationAddresses asserts the derived view surfaces both
// chain ids and every required address field from the live response.
func TestIntegrationAddresses(t *testing.T) {
	stdout, _, err := execLive(t, "addresses")
	if err != nil {
		t.Fatalf("addresses: %v", err)
	}
	// The derived view must include both chain ids + every *_address
	// field from the proto.
	required := []string{
		"bor_chain_id=",
		"heimdall_chain_id=",
		"pol_token_address=0x",
		"staking_manager_address=0x",
		"slash_manager_address=0x",
		"root_chain_address=0x",
		"staking_info_address=0x",
		"state_sender_address=0x",
		"state_receiver_address=0x",
		"validator_set_address=0x",
	}
	for _, line := range required {
		if !strings.Contains(stdout, line) {
			t.Errorf("addresses output missing %q\n%s", line, stdout)
		}
	}
	// Confirmation depth fields must not appear.
	if strings.Contains(stdout, "tx_confirmations") {
		t.Errorf("addresses leaked confirmation depths:\n%s", stdout)
	}
}

// TestIntegrationAliasCM exercises the `cm` alias end-to-end against the
// live node by running a direct argv without the literal `chainmanager`
// token.
func TestIntegrationAliasCM(t *testing.T) {
	local := &cobra.Command{
		Use:     "chainmanager",
		Aliases: []string{"cm"},
		Args:    cobra.NoArgs,
	}
	local.AddCommand(newParamsCmd(), newAddressesCmd())

	root := &cobra.Command{Use: "h", SilenceUsage: true}
	f := &config.Flags{}
	f.Register(root)
	flags = f
	root.AddCommand(local)

	var stdout, stderr bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stderr)
	root.SetArgs([]string{
		"--rest-url", liveREST(),
		"--rpc-url", liveRPC(),
		"--timeout", "10",
		"cm", "params", "--field", "params.chain_params.bor_chain_id",
	})
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := root.ExecuteContext(ctx); err != nil {
		t.Fatalf("cm params: %v stderr=%s", err, stderr.String())
	}
	if strings.TrimSpace(stdout.String()) == "" {
		t.Errorf("cm params produced empty output")
	}
}
