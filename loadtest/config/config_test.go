package config

import (
	"strings"
	"testing"
	"time"
)

// validConfig returns a minimal config that passes Validate.
func validConfig() *Config {
	return &Config{
		AdaptiveBackoffFactor: 2,
		GasPriceMultiplier:    1,
		RateLimit:             4,
		Modes:                 []string{"t"},
	}
}

func TestValidateRateLimitRampDuration(t *testing.T) {
	tests := []struct {
		name              string
		rampDuration      time.Duration
		rateLimit         float64
		adaptiveRateLimit bool
		wantErr           string
	}{
		{
			name:      "no ramp",
			rateLimit: 4,
		},
		{
			name:         "valid ramp",
			rampDuration: 3 * time.Minute,
			rateLimit:    100,
		},
		{
			name:         "negative duration",
			rampDuration: -time.Minute,
			rateLimit:    100,
			wantErr:      "--rate-limit-ramp-duration must be positive",
		},
		{
			name:              "mutually exclusive with adaptive",
			rampDuration:      3 * time.Minute,
			rateLimit:         100,
			adaptiveRateLimit: true,
			wantErr:           "--rate-limit-ramp-duration and --adaptive-rate-limit are mutually exclusive",
		},
		{
			name:         "requires positive rate limit",
			rampDuration: 3 * time.Minute,
			rateLimit:    -1,
			wantErr:      "--rate-limit-ramp-duration requires a positive --rate-limit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := validConfig()
			cfg.RateLimitRampDuration = tt.rampDuration
			cfg.RateLimit = tt.rateLimit
			cfg.AdaptiveRateLimit = tt.adaptiveRateLimit

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
