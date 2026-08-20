package loadtest

import (
	"strconv"
	"testing"
)

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
