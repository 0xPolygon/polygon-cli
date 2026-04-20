package heimdallutil

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
)

// buildVersionRoot mirrors the runtime wiring: heimdall root owns the
// persistent flag set; util is attached; version is a subcommand. The
// flag wiring mirrors Register but without pulling in the whole
// subcommand tree.
func buildVersionRoot() (*cobra.Command, *config.Flags) {
	root := &cobra.Command{Use: "h", SilenceUsage: true}
	f := &config.Flags{}
	f.Register(root)
	flags = f // package global used by resolveIfPossible
	util := &cobra.Command{Use: "util", Args: cobra.NoArgs}
	util.AddCommand(newVersionCmd())
	root.AddCommand(util)
	return root, f
}

// TestVersionDefault asserts the plain `version` command prints the
// polycli version and resolved chain-id via config, without touching
// the network.
func TestVersionDefault(t *testing.T) {
	root, _ := buildVersionRoot()
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)
	root.SetArgs([]string{"util", "version"})
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("version: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "polycli_version") {
		t.Errorf("missing polycli_version: %s", out)
	}
	if !strings.Contains(out, "chain_id") {
		t.Errorf("missing chain_id: %s", out)
	}
	// Without --node the command must NOT emit cometbft_version.
	if strings.Contains(out, "cometbft_version") {
		t.Errorf("unexpected cometbft_version in default output: %s", out)
	}
}

// TestVersionJSON asserts --json serializes to a parseable object.
func TestVersionJSON(t *testing.T) {
	root, _ := buildVersionRoot()
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)
	root.SetArgs([]string{"--json", "util", "version"})
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("version --json: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, buf.String())
	}
	if _, ok := m["polycli_version"]; !ok {
		t.Errorf("missing polycli_version in JSON: %v", m)
	}
}

// TestVersionField asserts --field plucks a single top-level key.
func TestVersionField(t *testing.T) {
	root, _ := buildVersionRoot()
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)
	root.SetArgs([]string{"util", "version", "--field", "polycli_version"})
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("version --field: %v", err)
	}
	out := strings.TrimSpace(buf.String())
	if out == "" {
		t.Errorf("empty --field output")
	}
}
