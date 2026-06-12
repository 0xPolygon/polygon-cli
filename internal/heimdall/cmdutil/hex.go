package cmdutil

import (
	"fmt"
	"strings"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
)

// NormalizeHex accepts a fixed-length hex value with or without the
// `0x` prefix and returns the lower-case, `0x`-prefixed form expected
// by the Heimdall REST endpoints. label names the value in errors
// (e.g. "address", "signer", "tx hash").
func NormalizeHex(raw string, byteLen int, label string) (string, error) {
	s := strings.TrimSpace(raw)
	s = strings.TrimPrefix(strings.TrimPrefix(s, "0x"), "0X")
	if len(s) != byteLen*2 {
		return "", &client.UsageError{Msg: fmt.Sprintf("%s must be %d bytes (%d hex chars), got %d", label, byteLen, byteLen*2, len(s))}
	}
	if r, ok := firstNonHex(s); ok {
		return "", &client.UsageError{Msg: fmt.Sprintf("invalid %s %q (non-hex character %q)", label, raw, r)}
	}
	return "0x" + strings.ToLower(s), nil
}

// NormalizeAddress normalizes a 20-byte Ethereum address.
func NormalizeAddress(raw string) (string, error) {
	return NormalizeHex(raw, 20, "address")
}

// NormalizeTxHash normalizes a 32-byte transaction hash. The REST
// endpoints expect the `0x` prefix and will 500 without it, so we
// re-add it unconditionally.
func NormalizeTxHash(raw string) (string, error) {
	return NormalizeHex(raw, 32, "tx hash")
}

// NormalizeHexBytes accepts a variable-length hex string with or
// without the `0x` prefix and returns the lower-case form WITHOUT the
// prefix (for use as a bare query param). Empty input is an error.
func NormalizeHexBytes(raw, label string) (string, error) {
	s := strings.TrimSpace(raw)
	s = strings.TrimPrefix(strings.TrimPrefix(s, "0x"), "0X")
	if s == "" {
		return "", &client.UsageError{Msg: fmt.Sprintf("%s must not be empty", label)}
	}
	if len(s)%2 != 0 {
		return "", &client.UsageError{Msg: fmt.Sprintf("%s must have an even number of hex chars, got %d", label, len(s))}
	}
	if r, ok := firstNonHex(s); ok {
		return "", &client.UsageError{Msg: fmt.Sprintf("invalid %s %q (non-hex character %q)", label, raw, r)}
	}
	return strings.ToLower(s), nil
}

// firstNonHex returns the first rune of s that is not a hex digit.
func firstNonHex(s string) (rune, bool) {
	for _, r := range s {
		switch {
		case r >= '0' && r <= '9':
		case r >= 'a' && r <= 'f':
		case r >= 'A' && r <= 'F':
		default:
			return r, true
		}
	}
	return 0, false
}
