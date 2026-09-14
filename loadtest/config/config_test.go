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

func TestValidateReverseNonceOrder(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*Config)
		wantErr string
	}{
		{
			name: "valid with fire-and-forget",
			mutate: func(c *Config) {
				c.FireAndForget = true
			},
		},
		{
			name:    "requires fire-and-forget",
			mutate:  func(c *Config) {},
			wantErr: "--reverse-nonce-order requires --fire-and-forget",
		},
		{
			name: "incompatible with wait-for-receipt",
			mutate: func(c *Config) {
				c.FireAndForget = true
				c.WaitForReceipt = true
				c.ReceiptRetryMax = 5
			},
			wantErr: "--reverse-nonce-order is incompatible with --wait-for-receipt",
		},
		{
			name: "incompatible with eth-call-only",
			mutate: func(c *Config) {
				c.FireAndForget = true
				c.EthCallOnly = true
			},
			wantErr: "--reverse-nonce-order doesn't make sense with --eth-call-only",
		},
		{
			name: "incompatible with duplicate-nonce-rate",
			mutate: func(c *Config) {
				c.FireAndForget = true
				c.DuplicateNonceRate = 1
			},
			wantErr: "--reverse-nonce-order is incompatible with --duplicate-nonce-rate",
		},
		{
			name: "incompatible with adaptive-rate-limit",
			mutate: func(c *Config) {
				c.FireAndForget = true
				c.AdaptiveRateLimit = true
			},
			wantErr: "--reverse-nonce-order is incompatible with --adaptive-rate-limit",
		},
		{
			name: "requires positive requests",
			mutate: func(c *Config) {
				c.FireAndForget = true
				c.Requests = 0
			},
			wantErr: "--reverse-nonce-order requires positive --requests and --concurrency",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := validConfig()
			cfg.ReverseNonceOrder = true
			cfg.Requests = 10
			cfg.Concurrency = 2
			tt.mutate(cfg)

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

func TestValidateSyncTxs(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*Config)
		wantErr string
	}{
		{
			name:   "sync txs alone",
			mutate: func(c *Config) { c.SyncTxs = true },
		},
		{
			name: "sync txs with timeout",
			mutate: func(c *Config) {
				c.SyncTxs = true
				c.SyncTxTimeout = 2 * time.Second
			},
		},
		{
			name: "sync and private are exclusive",
			mutate: func(c *Config) {
				c.SyncTxs = true
				c.PrivateTxs = true
			},
			wantErr: "--sync-txs and --private-txs are mutually exclusive",
		},
		{
			name: "sync with call only",
			mutate: func(c *Config) {
				c.SyncTxs = true
				c.EthCallOnly = true
			},
			wantErr: "--sync-txs doesn't make sense with --eth-call-only",
		},
		{
			name: "sync with raw output",
			mutate: func(c *Config) {
				c.SyncTxs = true
				c.OutputRawTxOnly = true
			},
			wantErr: "--sync-txs doesn't make sense with --output-raw-tx-only",
		},
		{
			name: "sync with an unsupported mode",
			mutate: func(c *Config) {
				c.SyncTxs = true
				c.Modes = []string{"d"}
			},
			wantErr: "--sync-txs is not supported for mode \"d\"",
		},
		{
			name:    "negative timeout",
			mutate:  func(c *Config) { c.SyncTxTimeout = -time.Second },
			wantErr: "--sync-tx-timeout must not be negative",
		},
		{
			name:    "sub-millisecond timeout rounds to zero",
			mutate:  func(c *Config) { c.SyncTxTimeout = 500 * time.Microsecond },
			wantErr: "rounds to zero",
		},
		{
			// The pairing that measures speculative then canonical inclusion.
			name: "sync with wait for receipt",
			mutate: func(c *Config) {
				c.SyncTxs = true
				c.WaitForReceipt = true
				c.ReceiptRetryMax = 30
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := validConfig()
			tt.mutate(c)

			err := c.Validate()
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("Validate() = %v, want nil", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("Validate() = nil, want error containing %q", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("Validate() = %q, want it to contain %q", err, tt.wantErr)
			}
		})
	}
}

func TestValidateReceiptPollInterval(t *testing.T) {
	tests := []struct {
		name            string
		waitForReceipt  bool
		pollInterval    time.Duration
		receiptRetryMax uint
		wantErr         string
	}{
		{
			name:            "backoff mode still requires retry max",
			waitForReceipt:  true,
			receiptRetryMax: 1,
			wantErr:         "use a max retry greater than 1",
		},
		{
			name:            "interval mode ignores retry max",
			waitForReceipt:  true,
			pollInterval:    50 * time.Millisecond,
			receiptRetryMax: 1,
		},
		{
			name:           "negative interval",
			waitForReceipt: true,
			pollInterval:   -time.Second,
			wantErr:        "--receipt-poll-interval must not be negative",
		},
		{
			name:         "interval requires wait-for-receipt",
			pollInterval: 50 * time.Millisecond,
			wantErr:      "--receipt-poll-interval requires --wait-for-receipt",
		},
		{
			name:            "valid interval",
			waitForReceipt:  true,
			pollInterval:    50 * time.Millisecond,
			receiptRetryMax: 30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := validConfig()
			cfg.WaitForReceipt = tt.waitForReceipt
			cfg.ReceiptPollInterval = tt.pollInterval
			cfg.ReceiptRetryMax = tt.receiptRetryMax

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
