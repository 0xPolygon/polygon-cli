package loadtest

import (
	"strconv"
	"testing"

	"github.com/0xPolygon/polygon-cli/loadtest/config"
	"github.com/ethereum/go-ethereum/common"
)

func TestComputeReverseTxsPerAccount(t *testing.T) {
	newPool := func(n int) *AccountPool {
		accounts := make([]*Account, n)
		for i := range accounts {
			accounts[i] = &Account{address: common.BytesToAddress([]byte{byte(i + 1), byte(i >> 8)}), ready: true}
		}
		return newTestAccountPool(accounts...)
	}

	tests := []struct {
		name              string
		concurrency       int64
		requests          int64
		accounts          int
		wantErr           bool
		wantTxsPerAccount uint64
	}{
		{
			name:              "even split",
			concurrency:       250,
			requests:          10,
			accounts:          500,
			wantTxsPerAccount: 5,
		},
		{
			name:              "single account",
			concurrency:       2,
			requests:          10,
			accounts:          1,
			wantTxsPerAccount: 20,
		},
		{
			name:        "uneven split",
			concurrency: 3,
			requests:    7,
			accounts:    2,
			wantErr:     true,
		},
		{
			name:        "more accounts than requests",
			concurrency: 1,
			requests:    250,
			accounts:    500,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Runner{
				cfg: &config.Config{
					Concurrency: tt.concurrency,
					Requests:    tt.requests,
				},
				accountPool: newPool(tt.accounts),
			}

			err := r.computeReverseTxsPerAccount()
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr = %v", err, tt.wantErr)
			}
			if err == nil && r.reverseTxsPerAccount != tt.wantTxsPerAccount {
				t.Errorf("reverseTxsPerAccount = %d, want %d", r.reverseTxsPerAccount, tt.wantTxsPerAccount)
			}
		})
	}
}

func TestNonceTooLowRegexp(t *testing.T) {
	tests := []struct {
		name          string
		errMsg        string
		wantMatch     bool
		wantNextNonce uint64
	}{
		{
			name:          "geth format",
			errMsg:        "nonce too low: next nonce 78, tx nonce 44",
			wantMatch:     true,
			wantNextNonce: 78,
		},
		{
			name:          "geth format wrapped with context",
			errMsg:        "failed to send transaction: nonce too low: next nonce 78, tx nonce 44",
			wantMatch:     true,
			wantNextNonce: 78,
		},
		{
			name:      "nonce too low without details",
			errMsg:    "nonce too low",
			wantMatch: false,
		},
		{
			name:      "unrelated error",
			errMsg:    "insufficient funds for gas * price + value",
			wantMatch: false,
		},
		{
			name:      "different client format",
			errMsg:    "nonce too low: address 0xabc, tx: 44 state: 78",
			wantMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match := nonceTooLowRegexp.FindStringSubmatch(tt.errMsg)
			if (match != nil) != tt.wantMatch {
				t.Fatalf("match = %v, wantMatch = %v", match, tt.wantMatch)
			}
			if match == nil {
				return
			}
			nextNonce, err := strconv.ParseUint(match[1], 10, 64)
			if err != nil {
				t.Fatalf("failed to parse next nonce: %v", err)
			}
			if nextNonce != tt.wantNextNonce {
				t.Errorf("nextNonce = %d, want %d", nextNonce, tt.wantNextNonce)
			}
		})
	}
}
