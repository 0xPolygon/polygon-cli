package config

import (
	"strings"
	"testing"
)

// validConfig returns a minimal config that passes Validate.
func validConfig() *Config {
	return &Config{
		AdaptiveBackoffFactor: 2,
		GasPriceMultiplier:    1,
		Modes:                 []string{"t"},
	}
}

func TestValidateSendRPCURL(t *testing.T) {
	tests := []struct {
		name       string
		sendRPCURL string
		privateTxs bool
		modes      []string
		wantErr    string
	}{
		{
			name:  "unset send-rpc-url with any mode",
			modes: []string{"erc20"},
		},
		{
			name:       "supported modes",
			sendRPCURL: "http://localhost:8546",
			modes:      []string{"t", "transaction", "b", "blob", "cc", "contract-call", "R", "recall"},
		},
		{
			name:       "unsupported mode",
			sendRPCURL: "http://localhost:8546",
			modes:      []string{"erc20"},
			wantErr:    `--send-rpc-url is not supported for mode "erc20"`,
		},
		{
			name:       "invalid url scheme",
			sendRPCURL: "localhost:8546",
			modes:      []string{"t"},
			wantErr:    "invalid --send-rpc-url",
		},
		{
			name:       "combined with private-txs",
			sendRPCURL: "https://private.example.com",
			privateTxs: true,
			modes:      []string{"t"},
		},
		{
			name:       "private-txs unsupported mode",
			privateTxs: true,
			modes:      []string{"uniswapv3"},
			wantErr:    `--private-txs is not supported for mode "uniswapv3"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := validConfig()
			cfg.SendRPCURL = tt.sendRPCURL
			cfg.PrivateTxs = tt.privateTxs
			cfg.Modes = tt.modes

			err := cfg.Validate()
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("Validate() unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("Validate() expected error containing %q, got nil", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("Validate() error %q does not contain %q", err.Error(), tt.wantErr)
			}
		})
	}
}
