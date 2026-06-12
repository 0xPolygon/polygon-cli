//go:build heimdall_integration

package chain

import (
	"bytes"
	"context"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
)

// integration tests run against the live Amoy-backed node at
// 172.19.0.2:26657 unless overridden by HEIMDALL_TEST_RPC_URL.

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

func execLive(t *testing.T, args ...string) (string, string, error) {
	t.Helper()
	root := &cobra.Command{Use: "h", SilenceUsage: true}
	f := &config.Flags{}
	f.Register(root)
	Register(root, f)
	var stdout, stderr bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stderr)
	all := append([]string{
		"--rest-url", liveREST(),
		"--rpc-url", liveRPC(),
		"--timeout", "10",
	}, args...)
	root.SetArgs(all)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := root.ExecuteContext(ctx)
	return stdout.String(), stderr.String(), err
}

func TestIntegrationBlockNumber(t *testing.T) {
	stdout, _, err := execLive(t, "block-number")
	if err != nil {
		t.Fatalf("block-number: %v", err)
	}
	n, err := strconv.ParseInt(strings.TrimSpace(stdout), 10, 64)
	if err != nil {
		t.Fatalf("block-number output %q is not an integer: %v", stdout, err)
	}
	if n <= 0 {
		t.Fatalf("block-number %d should be > 0", n)
	}
}

func TestIntegrationChainID(t *testing.T) {
	stdout, _, err := execLive(t, "chain-id")
	if err != nil {
		t.Fatalf("chain-id: %v", err)
	}
	got := strings.TrimSpace(stdout)
	if got != "heimdallv2-80002" {
		t.Fatalf("chain-id = %q, want heimdallv2-80002", got)
	}
}

func TestIntegrationClient(t *testing.T) {
	stdout, _, err := execLive(t, "client")
	if err != nil {
		t.Fatalf("client: %v", err)
	}
	if !strings.Contains(stdout, "cometbft_version") {
		t.Errorf("expected cometbft_version key: %q", stdout)
	}
	// Expect any non-empty version string on the cometbft_version line.
	for _, line := range strings.Split(stdout, "\n") {
		if strings.HasPrefix(line, "cometbft_version") {
			parts := strings.Fields(line)
			if len(parts) < 2 {
				t.Errorf("cometbft_version line has no value: %q", line)
			}
			if parts[1] == "" {
				t.Errorf("cometbft_version is empty")
			}
		}
	}
}

func TestIntegrationFindBlockRecent(t *testing.T) {
	// Grab tip, then ask for a block 300s behind tip's wall clock
	// time and confirm the returned height lies within the
	// [tip-600, tip] window — generous window keeps the test stable
	// against chain rate jitter.
	tipOut, _, err := execLive(t, "block-number")
	if err != nil {
		t.Fatalf("block-number: %v", err)
	}
	tip, err := strconv.ParseInt(strings.TrimSpace(tipOut), 10, 64)
	if err != nil {
		t.Fatalf("bad tip %q: %v", tipOut, err)
	}
	target := strconv.FormatInt(time.Now().Unix()-300, 10)
	found, _, err := execLive(t, "find-block", target)
	if err != nil {
		t.Fatalf("find-block %s: %v", target, err)
	}
	h, err := strconv.ParseInt(strings.TrimSpace(found), 10, 64)
	if err != nil {
		t.Fatalf("find-block returned non-int %q: %v", found, err)
	}
	// Heimdall makes ~1 block/s so 300s back should be about
	// tip-300. Allow ±400 to tolerate rate changes + test drift.
	lo, hi := tip-700, tip
	if h < lo || h > hi {
		t.Fatalf("find-block %s returned %d, expected in [%d, %d]", target, h, lo, hi)
	}
}
