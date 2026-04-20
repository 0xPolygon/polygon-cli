// Package config resolves the runtime configuration for the
// `polycli heimdall` command tree.
//
// Precedence (highest wins): explicit command-line flag > environment
// variable > config file (~/.polycli/heimdall.toml) > network preset >
// built-in default. A missing config file is not an error.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

// Exit codes mirror cast conventions (see requirements §2.1).
const (
	ExitOK       = 0
	ExitNodeErr  = 1
	ExitNetErr   = 2
	ExitUsageErr = 3
	ExitSignErr  = 4
)

// Default file location for the optional user config.
const defaultConfigRelPath = ".polycli/heimdall.toml"

// Environment-variable names (requirements §2.2).
const (
	EnvNetwork    = "HEIMDALL_NETWORK"
	EnvRESTURL    = "HEIMDALL_REST_URL"
	EnvRPCURL     = "HEIMDALL_RPC_URL"
	EnvChainID    = "HEIMDALL_CHAIN_ID"
	EnvDenom      = "HEIMDALL_FEE_DENOM"
	EnvTimeout    = "HEIMDALL_TIMEOUT"
	EnvRPCHeaders = "HEIMDALL_RPC_HEADERS"
	EnvNoColor    = "NO_COLOR"
)

// Config holds the fully resolved runtime configuration.
type Config struct {
	Network    string
	RESTURL    string
	RPCURL     string
	ChainID    string
	Denom      string
	Timeout    time.Duration
	RPCHeaders map[string]string
	Insecure   bool
	JSON       bool
	Curl       bool
	Color      string // auto|always|never
	Raw        bool
}

// Preset is a named set of defaults for a known Heimdall network.
type Preset struct {
	Name    string
	RESTURL string
	RPCURL  string
	ChainID string
}

var presets = map[string]Preset{
	"amoy": {
		Name:    "amoy",
		RESTURL: "https://heimdall-api-amoy.polygon.technology",
		RPCURL:  "https://tendermint-api-amoy.polygon.technology",
		ChainID: "heimdallv2-80002",
	},
	"mainnet": {
		Name:    "mainnet",
		RESTURL: "https://heimdall-api.polygon.technology",
		RPCURL:  "https://tendermint-api.polygon.technology",
		ChainID: "heimdallv2-137",
	},
}

// Preset names in stable order (used for error messages + help text).
func PresetNames() []string {
	return []string{"amoy", "mainnet"}
}

// GetPreset returns a copy of the named preset, or (Preset{}, false).
func GetPreset(name string) (Preset, bool) {
	p, ok := presets[name]
	return p, ok
}

// Flags holds the raw command-line flag state before resolution.
// Persistent flags bind to the fields of this struct.
type Flags struct {
	Mainnet    bool
	Amoy       bool
	Network    string
	RESTURL    string
	RPCURL     string
	ChainID    string
	Denom      string
	TimeoutSec int
	RPCHeaders string
	Insecure   bool
	JSON       bool
	Curl       bool
	Color      string
	NoColor    bool
	Raw        bool
	ConfigPath string
}

// Register binds the persistent heimdall flags to the given command's
// PersistentFlags set and wires them into f.
func (f *Flags) Register(cmd *cobra.Command) {
	p := cmd.PersistentFlags()
	p.BoolVar(&f.Mainnet, "mainnet", false, "shortcut for --network mainnet")
	p.BoolVar(&f.Amoy, "amoy", false, "shortcut for --network amoy (default)")
	p.StringVarP(&f.Network, "network", "N", "", "named network preset (amoy|mainnet)")
	p.StringVarP(&f.RESTURL, "rest-url", "r", "", "heimdall REST gateway URL")
	p.StringVarP(&f.RPCURL, "rpc-url", "R", "", "cometBFT RPC URL")
	p.StringVar(&f.ChainID, "chain-id", "", "chain id used for signing")
	p.StringVar(&f.Denom, "denom", "", "fee denom")
	p.IntVar(&f.TimeoutSec, "timeout", 0, "HTTP timeout in seconds")
	p.StringVar(&f.RPCHeaders, "rpc-headers", "", "extra request headers, comma-separated key=value pairs")
	p.BoolVarP(&f.Insecure, "insecure", "k", false, "accept invalid TLS certs")
	p.BoolVar(&f.JSON, "json", false, "emit JSON instead of key/value")
	p.BoolVar(&f.Curl, "curl", false, "print the equivalent curl command instead of executing")
	p.StringVar(&f.Color, "color", "auto", "color mode (auto|always|never)")
	p.BoolVar(&f.NoColor, "no-color", false, "disable color output")
	p.BoolVar(&f.Raw, "raw", false, "preserve raw bytes (no 0x-hex normalization)")
	p.StringVar(&f.ConfigPath, "heimdall-config", "", "path to heimdall config TOML (default ~/.polycli/heimdall.toml)")
}

// Resolve assembles the final Config by layering, in order: built-in
// defaults -> preset -> optional config file -> env vars -> flags.
// Returns an error if the selected network is unknown or a provided
// config file is malformed.
func Resolve(f *Flags) (*Config, error) {
	network, err := resolveNetworkName(f)
	if err != nil {
		return nil, err
	}

	preset, ok := presets[network]
	if !ok {
		file, ferr := loadConfigFile(f.ConfigPath)
		if ferr != nil {
			return nil, ferr
		}
		if file != nil {
			if p, pok := file.Networks[network]; pok {
				preset = Preset{Name: network, RESTURL: p.RESTURL, RPCURL: p.RPCURL, ChainID: p.ChainID}
				ok = true
			}
		}
		if !ok {
			return nil, fmt.Errorf("unknown network %q (known: %s)", network, strings.Join(PresetNames(), ", "))
		}
	}

	// Merge config file on top of preset defaults.
	file, err := loadConfigFile(f.ConfigPath)
	if err != nil {
		return nil, err
	}
	cfg := &Config{
		Network: network,
		RESTURL: preset.RESTURL,
		RPCURL:  preset.RPCURL,
		ChainID: preset.ChainID,
		Denom:   "pol",
		Timeout: 30 * time.Second,
		Color:   "auto",
	}
	if file != nil {
		if n, ok := file.Networks[network]; ok {
			overlayNetwork(cfg, n)
		}
		overlayGlobal(cfg, file)
	}

	// Environment variables.
	overlayEnv(cfg)

	// Command-line flags (highest precedence).
	overlayFlags(cfg, f)

	// Validate.
	if cfg.RESTURL == "" {
		return nil, errors.New("no REST URL resolved (set --rest-url or HEIMDALL_REST_URL)")
	}
	if cfg.RPCURL == "" {
		return nil, errors.New("no RPC URL resolved (set --rpc-url or HEIMDALL_RPC_URL)")
	}
	if cfg.ChainID == "" {
		return nil, errors.New("no chain id resolved (set --chain-id or HEIMDALL_CHAIN_ID)")
	}
	if cfg.Timeout <= 0 {
		return nil, fmt.Errorf("timeout must be positive, got %s", cfg.Timeout)
	}

	return cfg, nil
}

// resolveNetworkName returns the network name honouring the
// flag > env > default chain. --mainnet/--amoy are sugar and must
// not conflict.
func resolveNetworkName(f *Flags) (string, error) {
	if f.Mainnet && f.Amoy {
		return "", errors.New("--mainnet and --amoy are mutually exclusive")
	}
	switch {
	case f.Mainnet:
		return "mainnet", nil
	case f.Amoy:
		return "amoy", nil
	}
	if f.Network != "" {
		return f.Network, nil
	}
	if v := os.Getenv(EnvNetwork); v != "" {
		return v, nil
	}
	return "amoy", nil
}

type fileConfig struct {
	DefaultNetwork string                    `toml:"default_network"`
	RESTURL        string                    `toml:"rest_url"`
	RPCURL         string                    `toml:"rpc_url"`
	ChainID        string                    `toml:"chain_id"`
	Denom          string                    `toml:"denom"`
	TimeoutSec     int                       `toml:"timeout"`
	Networks       map[string]fileConfigNet  `toml:"networks"`
}

type fileConfigNet struct {
	RESTURL string `toml:"rest_url"`
	RPCURL  string `toml:"rpc_url"`
	ChainID string `toml:"chain_id"`
}

func loadConfigFile(explicit string) (*fileConfig, error) {
	path := explicit
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, nil
		}
		path = filepath.Join(home, defaultConfigRelPath)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) && explicit == "" {
			return nil, nil
		}
		return nil, fmt.Errorf("reading heimdall config %s: %w", path, err)
	}
	out := &fileConfig{}
	if err := toml.Unmarshal(raw, out); err != nil {
		return nil, fmt.Errorf("parsing heimdall config %s: %w", path, err)
	}
	return out, nil
}

func overlayNetwork(cfg *Config, n fileConfigNet) {
	if n.RESTURL != "" {
		cfg.RESTURL = n.RESTURL
	}
	if n.RPCURL != "" {
		cfg.RPCURL = n.RPCURL
	}
	if n.ChainID != "" {
		cfg.ChainID = n.ChainID
	}
}

func overlayGlobal(cfg *Config, f *fileConfig) {
	if f.RESTURL != "" {
		cfg.RESTURL = f.RESTURL
	}
	if f.RPCURL != "" {
		cfg.RPCURL = f.RPCURL
	}
	if f.ChainID != "" {
		cfg.ChainID = f.ChainID
	}
	if f.Denom != "" {
		cfg.Denom = f.Denom
	}
	if f.TimeoutSec > 0 {
		cfg.Timeout = time.Duration(f.TimeoutSec) * time.Second
	}
}

func overlayEnv(cfg *Config) {
	if v := os.Getenv(EnvRESTURL); v != "" {
		cfg.RESTURL = v
	}
	if v := os.Getenv(EnvRPCURL); v != "" {
		cfg.RPCURL = v
	}
	if v := os.Getenv(EnvChainID); v != "" {
		cfg.ChainID = v
	}
	if v := os.Getenv(EnvDenom); v != "" {
		cfg.Denom = v
	}
	if v := os.Getenv(EnvTimeout); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.Timeout = time.Duration(n) * time.Second
		}
	}
	if v := os.Getenv(EnvRPCHeaders); v != "" {
		cfg.RPCHeaders = parseHeaders(v)
	}
	if os.Getenv(EnvNoColor) != "" {
		cfg.Color = "never"
	}
}

func overlayFlags(cfg *Config, f *Flags) {
	if f.RESTURL != "" {
		cfg.RESTURL = f.RESTURL
	}
	if f.RPCURL != "" {
		cfg.RPCURL = f.RPCURL
	}
	if f.ChainID != "" {
		cfg.ChainID = f.ChainID
	}
	if f.Denom != "" {
		cfg.Denom = f.Denom
	}
	if f.TimeoutSec > 0 {
		cfg.Timeout = time.Duration(f.TimeoutSec) * time.Second
	}
	if f.RPCHeaders != "" {
		headers := parseHeaders(f.RPCHeaders)
		if cfg.RPCHeaders == nil {
			cfg.RPCHeaders = headers
		} else {
			for k, v := range headers {
				cfg.RPCHeaders[k] = v
			}
		}
	}
	cfg.Insecure = f.Insecure
	cfg.JSON = f.JSON
	cfg.Curl = f.Curl
	cfg.Raw = f.Raw
	if f.NoColor {
		cfg.Color = "never"
	} else if f.Color != "" && f.Color != "auto" {
		cfg.Color = f.Color
	}
}

// parseHeaders splits "K1=V1,K2=V2" into a map. Malformed pairs are
// silently ignored to avoid blocking a query over bad operator input —
// logged upstream if needed.
func parseHeaders(raw string) map[string]string {
	out := map[string]string{}
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		eq := strings.IndexByte(part, '=')
		if eq <= 0 {
			continue
		}
		k := strings.TrimSpace(part[:eq])
		v := strings.TrimSpace(part[eq+1:])
		if k == "" {
			continue
		}
		out[k] = v
	}
	return out
}
