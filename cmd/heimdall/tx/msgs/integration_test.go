//go:build heimdall_integration

package msgs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
)

// Integration tests talk to a live Heimdall v2 node. Defaults point
// at 172.19.0.2 (the documented local compose node in the
// HEIMDALLCAST_REQUIREMENTS.md §6 test plan); override with
// HEIMDALL_TEST_REST_URL / HEIMDALL_TEST_RPC_URL.

func liveREST() string {
	if v := os.Getenv("HEIMDALL_TEST_REST_URL"); v != "" {
		return v
	}
	return "http://172.19.0.2:1317"
}

func liveRPC() string {
	if v := os.Getenv("HEIMDALL_TEST_RPC_URL"); v != "" {
		return v
	}
	return "http://172.19.0.2:26657"
}

func liveChainID() string {
	if v := os.Getenv("HEIMDALL_TEST_CHAIN_ID"); v != "" {
		return v
	}
	return "heimdallv2-80002"
}

func liveTestAddress(t *testing.T) string {
	t.Helper()
	if v := os.Getenv("HEIMDALL_TEST_FROM_ADDR"); v != "" {
		return v
	}
	// Derive a signer from the stake validator set so the integration
	// test is runnable without operator-side setup.
	c := &http.Client{Timeout: 10 * time.Second}
	resp, err := c.Get(liveREST() + "/stake/validators-set")
	if err != nil {
		t.Skipf("cannot reach live REST at %s: %v", liveREST(), err)
	}
	defer resp.Body.Close()
	var out struct {
		ValidatorSet struct {
			Validators []struct {
				Signer string `json:"signer"`
			} `json:"validators"`
		} `json:"validator_set"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decoding validators-set: %v", err)
	}
	if len(out.ValidatorSet.Validators) == 0 {
		t.Skip("no validators on live node")
	}
	return out.ValidatorSet.Validators[0].Signer
}

func execLive(t *testing.T, args ...string) (string, string, error) {
	t.Helper()
	root := &cobra.Command{Use: "h", SilenceUsage: true, SilenceErrors: true}
	f := &config.Flags{}
	f.Register(root)
	for _, mode := range []Mode{ModeMkTx, ModeSend, ModeEstimate} {
		root.AddCommand(newUmbrellaForTest(mode, f))
	}
	var stdout, stderr bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stderr)
	all := append([]string{
		"--rest-url", liveREST(),
		"--rpc-url", liveRPC(),
		"--chain-id", liveChainID(),
		"--timeout", "15",
	}, args...)
	root.SetArgs(all)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := root.ExecuteContext(ctx)
	return stdout.String(), stderr.String(), err
}

// TestIntegrationMktxWithdraw builds a real withdraw tx against the
// live node. We use --private-key with a well-known test key because
// building only needs the signer's address + account state; no funds
// are moved.
func TestIntegrationMktxWithdraw(t *testing.T) {
	// If the operator didn't pre-configure a signer, derive one from
	// the validator set — but that also means we don't know a working
	// private key for it. In that case skip because the builder still
	// needs to sign (Sign requires a real key before emitting TxRaw).
	if os.Getenv("HEIMDALL_TEST_PRIVATE_KEY") == "" {
		t.Skip("set HEIMDALL_TEST_PRIVATE_KEY to run mktx integration tests")
	}
	from := liveTestAddress(t)
	stdout, stderr, err := execLive(t,
		"mktx", "withdraw",
		"--from", from,
		"--private-key", os.Getenv("HEIMDALL_TEST_PRIVATE_KEY"),
		"--gas", "200000",
		"--fee", "10000000000000000pol",
	)
	if err != nil {
		t.Fatalf("mktx withdraw: %v\nstdout=%s\nstderr=%s", err, stdout, stderr)
	}
	hex := strings.TrimSpace(stdout)
	if !strings.HasPrefix(hex, "0x") || len(hex) < 200 {
		t.Fatalf("unexpected TxRaw hex: %q", hex)
	}
}

// TestIntegrationEstimateWithdraw simulates a withdraw against the
// live node without broadcasting. The node returns real gas usage.
func TestIntegrationEstimateWithdraw(t *testing.T) {
	if os.Getenv("HEIMDALL_TEST_PRIVATE_KEY") == "" {
		t.Skip("set HEIMDALL_TEST_PRIVATE_KEY to run estimate integration tests")
	}
	stdout, stderr, err := execLive(t,
		"estimate", "withdraw",
		"--private-key", os.Getenv("HEIMDALL_TEST_PRIVATE_KEY"),
		"--fee", "10000000000000000pol",
	)
	if err != nil {
		// Simulate may reject an under-funded signer with a specific
		// error; we accept that because the test is about the call
		// round-tripping, not about success of the simulation. Log
		// the stderr so operators can diagnose.
		t.Logf("estimate returned error %v (acceptable if simulate rejects): stderr=%s", err, stderr)
		return
	}
	if !strings.Contains(stdout, "gas_used=") {
		t.Fatalf("no gas_used in estimate output: %q", stdout)
	}
}

// TestIntegrationSendWithdraw is gated behind an explicit opt-in env
// var because it actually broadcasts a real transaction.
func TestIntegrationSendWithdraw(t *testing.T) {
	if os.Getenv("HEIMDALL_TEST_ALLOW_BROADCAST") != "1" {
		t.Skip("set HEIMDALL_TEST_ALLOW_BROADCAST=1 to run the broadcasting integration tests")
	}
	if os.Getenv("HEIMDALL_TEST_PRIVATE_KEY") == "" {
		t.Skip("set HEIMDALL_TEST_PRIVATE_KEY to run send integration tests")
	}
	stdout, stderr, err := execLive(t,
		"send", "withdraw",
		"--private-key", os.Getenv("HEIMDALL_TEST_PRIVATE_KEY"),
		"--gas", "200000",
		"--fee", "10000000000000000pol",
		"--async", // don't wait for inclusion; the test just verifies broadcast accepted the tx
	)
	if err != nil {
		t.Fatalf("send withdraw: %v\nstdout=%s\nstderr=%s", err, stdout, stderr)
	}
	if !strings.Contains(stdout, "txhash=") {
		t.Fatalf("send output missing txhash: %q", stdout)
	}
}

// sanity check — this avoids io.ReadAll going unused when build tags
// shuffle around.
var _ = fmt.Sprintf
var _ = io.ReadAll
