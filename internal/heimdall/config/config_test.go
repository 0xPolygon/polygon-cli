package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// clearEnv unsets every heimdall-relevant env var so that tests don't
// leak through the process environment.
func clearEnv(t *testing.T) {
	t.Helper()
	for _, k := range []string{
		EnvNetwork, EnvRESTURL, EnvRPCURL, EnvChainID,
		EnvDenom, EnvTimeout, EnvRPCHeaders, EnvNoColor,
	} {
		t.Setenv(k, "")
	}
}

func TestResolveDefaultsToAmoy(t *testing.T) {
	clearEnv(t)
	cfg, err := Resolve(&Flags{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Network != "amoy" {
		t.Errorf("Network = %q, want amoy", cfg.Network)
	}
	if cfg.ChainID != "heimdallv2-80002" {
		t.Errorf("ChainID = %q, want heimdallv2-80002", cfg.ChainID)
	}
	if cfg.Denom != "pol" {
		t.Errorf("Denom = %q, want pol", cfg.Denom)
	}
	if cfg.Timeout != 30*time.Second {
		t.Errorf("Timeout = %s, want 30s", cfg.Timeout)
	}
}

func TestResolveMainnetShortcut(t *testing.T) {
	clearEnv(t)
	cfg, err := Resolve(&Flags{Mainnet: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Network != "mainnet" {
		t.Errorf("Network = %q, want mainnet", cfg.Network)
	}
	if cfg.ChainID != "heimdallv2-137" {
		t.Errorf("ChainID = %q, want heimdallv2-137", cfg.ChainID)
	}
}

func TestResolveMutuallyExclusiveShortcuts(t *testing.T) {
	clearEnv(t)
	if _, err := Resolve(&Flags{Mainnet: true, Amoy: true}); err == nil {
		t.Fatal("expected error for --mainnet + --amoy, got nil")
	}
}

func TestResolveUnknownNetwork(t *testing.T) {
	clearEnv(t)
	if _, err := Resolve(&Flags{Network: "does-not-exist"}); err == nil {
		t.Fatal("expected error for unknown network, got nil")
	}
}

func TestPrecedenceFlagBeatsEnvBeatsDefault(t *testing.T) {
	clearEnv(t)

	// Env alone: env wins over default.
	t.Setenv(EnvChainID, "heimdallv2-env")
	cfg, err := Resolve(&Flags{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ChainID != "heimdallv2-env" {
		t.Errorf("ChainID via env = %q, want heimdallv2-env", cfg.ChainID)
	}

	// Flag + env: flag wins.
	cfg, err = Resolve(&Flags{ChainID: "heimdallv2-flag"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ChainID != "heimdallv2-flag" {
		t.Errorf("ChainID via flag = %q, want heimdallv2-flag", cfg.ChainID)
	}
}

func TestPrecedenceConfigFileBeatsPreset(t *testing.T) {
	clearEnv(t)

	// Config file with a named network override.
	dir := t.TempDir()
	path := filepath.Join(dir, "heimdall.toml")
	contents := `
default_network = "amoy"

[networks.amoy]
rest_url = "http://config-rest:1317"
rpc_url  = "http://config-rpc:26657"
chain_id = "heimdallv2-from-config"
`
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatalf("writing config: %v", err)
	}

	cfg, err := Resolve(&Flags{ConfigPath: path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ChainID != "heimdallv2-from-config" {
		t.Errorf("ChainID = %q, want heimdallv2-from-config", cfg.ChainID)
	}
	if cfg.RESTURL != "http://config-rest:1317" {
		t.Errorf("RESTURL = %q, want http://config-rest:1317", cfg.RESTURL)
	}

	// Flag still beats config file.
	cfg, err = Resolve(&Flags{ConfigPath: path, ChainID: "heimdallv2-flag"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ChainID != "heimdallv2-flag" {
		t.Errorf("ChainID = %q, want heimdallv2-flag", cfg.ChainID)
	}
}

func TestConfigFileDefinesCustomNetwork(t *testing.T) {
	clearEnv(t)
	dir := t.TempDir()
	path := filepath.Join(dir, "heimdall.toml")
	contents := `
[networks.devnet]
rest_url = "http://10.0.0.5:1317"
rpc_url  = "http://10.0.0.5:26657"
chain_id = "heimdallv2-dev"
`
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatalf("writing config: %v", err)
	}

	cfg, err := Resolve(&Flags{Network: "devnet", ConfigPath: path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ChainID != "heimdallv2-dev" {
		t.Errorf("ChainID = %q, want heimdallv2-dev", cfg.ChainID)
	}
}

func TestMissingConfigFileIsOK(t *testing.T) {
	clearEnv(t)
	// Rely on the default path - it may or may not exist. Pass no
	// explicit path and ensure we still resolve.
	cfg, err := Resolve(&Flags{})
	if err != nil {
		t.Fatalf("expected success without config file, got %v", err)
	}
	if cfg.Network != "amoy" {
		t.Errorf("Network = %q, want amoy", cfg.Network)
	}
}

func TestMalformedConfigFileErrors(t *testing.T) {
	clearEnv(t)
	dir := t.TempDir()
	path := filepath.Join(dir, "heimdall.toml")
	if err := os.WriteFile(path, []byte("not = valid = toml ] [["), 0o600); err != nil {
		t.Fatalf("writing config: %v", err)
	}
	if _, err := Resolve(&Flags{ConfigPath: path}); err == nil {
		t.Fatal("expected parse error, got nil")
	}
}

func TestParseHeaders(t *testing.T) {
	got := parseHeaders("X-Auth=secret, X-Trace = t1 , , bad")
	want := map[string]string{"X-Auth": "secret", "X-Trace": "t1"}
	if len(got) != len(want) {
		t.Fatalf("parseHeaders len = %d, want %d: %v", len(got), len(want), got)
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("%s = %q, want %q", k, got[k], v)
		}
	}
}

func TestNoColorEnvForcesNever(t *testing.T) {
	clearEnv(t)
	t.Setenv(EnvNoColor, "1")
	cfg, err := Resolve(&Flags{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Color != "never" {
		t.Errorf("Color = %q, want never", cfg.Color)
	}
}

func TestNoColorFlagForcesNever(t *testing.T) {
	clearEnv(t)
	cfg, err := Resolve(&Flags{NoColor: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Color != "never" {
		t.Errorf("Color = %q, want never", cfg.Color)
	}
}
