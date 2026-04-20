package msgs

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
)

// parseHexBytes returns the decoded bytes of a 0x-prefixed or bare hex
// string. Empty input returns a nil slice and no error so msg fields
// that accept optional bytes can pass --flag="" through unchanged.
// expectedLen == 0 disables the length check.
func parseHexBytes(flagName, raw string, expectedLen int) ([]byte, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return nil, nil
	}
	s = strings.TrimPrefix(strings.TrimPrefix(s, "0x"), "0X")
	if len(s)%2 != 0 {
		return nil, &client.UsageError{Msg: fmt.Sprintf("--%s: odd hex length %d", flagName, len(s))}
	}
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, &client.UsageError{Msg: fmt.Sprintf("--%s: invalid hex: %v", flagName, err)}
	}
	if expectedLen > 0 && len(b) != expectedLen {
		return nil, &client.UsageError{Msg: fmt.Sprintf("--%s must be %d bytes (got %d)", flagName, expectedLen, len(b))}
	}
	return b, nil
}

// requireNonEmptyString returns a UsageError when s is blank.
func requireNonEmptyString(flagName, s string) error {
	if strings.TrimSpace(s) == "" {
		return &client.UsageError{Msg: fmt.Sprintf("--%s is required", flagName)}
	}
	return nil
}

// lowerEthAddress normalises s to lowercase 0x-prefixed hex. Also
// validates the 20-byte length. Returns a UsageError otherwise.
func lowerEthAddress(flagName, s string) (string, error) {
	s = strings.TrimSpace(s)
	trimmed := strings.TrimPrefix(strings.TrimPrefix(s, "0x"), "0X")
	if len(trimmed) != 40 {
		return "", &client.UsageError{Msg: fmt.Sprintf("--%s must be a 20-byte hex address", flagName)}
	}
	for _, c := range trimmed {
		ok := (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
		if !ok {
			return "", &client.UsageError{Msg: fmt.Sprintf("--%s must be hex", flagName)}
		}
	}
	return "0x" + strings.ToLower(trimmed), nil
}
