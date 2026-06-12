//go:build heimdall_integration

package heimdallutil

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

// Integration tests reach the live Heimdall v2 node on 172.19.0.2
// unless overridden via HEIMDALL_TEST_REST_URL / HEIMDALL_TEST_RPC_URL.

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

func buildLiveRoot() *cobra.Command {
	root := &cobra.Command{Use: "h", SilenceUsage: true}
	f := &config.Flags{}
	f.Register(root)
	flags = f
	util := &cobra.Command{Use: "util", Args: cobra.NoArgs}
	util.AddCommand(newVersionCmd())
	root.AddCommand(util)
	return root
}

func TestIntegrationVersionNode(t *testing.T) {
	root := buildLiveRoot()
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)
	root.SetArgs([]string{
		"--rest-url", liveREST(),
		"--rpc-url", liveRPC(),
		"--timeout", "10",
		"--json",
		"util", "version", "--node",
	})
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := root.ExecuteContext(ctx); err != nil {
		t.Fatalf("version --node: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, buf.String())
	}
	if v, ok := m["cometbft_version"].(string); !ok || v == "" {
		t.Errorf("missing or empty cometbft_version: %v", m)
	}
	if v, ok := m["network_id"].(string); !ok || !strings.HasPrefix(v, "heimdall") {
		t.Errorf("unexpected network_id: %v", m["network_id"])
	}
	if v, ok := m["polycli_version"].(string); !ok || v == "" {
		t.Errorf("missing polycli_version: %v", m)
	}
}
