// Package decode implements the `polycli heimdall decode` umbrella
// command and its subcommands (tx, msg, hash-tx, ve). All decoders are
// offline: they read the cached proto registry in internal/heimdall/proto
// and never reach the network.
package decode

import (
	_ "embed"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
)

//go:embed usage.md
var usage string

// flags is injected by Register. Decoders do not use the network, but
// keep the shared Flags handle in case future subcommands need it (for
// instance, a --json flag inherited from the root).
var flags *config.Flags

// Register attaches the decode umbrella and its children to parent.
// The umbrella command is created fresh on every call so tests that
// build a throwaway root do not accumulate duplicate subcommands on a
// shared global.
func Register(parent *cobra.Command, f *config.Flags) {
	flags = f
	cmd := &cobra.Command{
		Use:   "decode",
		Short: "Offline proto decoders for Heimdall tx / msg / vote-extension bytes.",
		Long:  usage,
		Args:  cobra.NoArgs,
	}
	cmd.AddCommand(
		newTxCmd(),
		newMsgCmd(),
		newHashTxCmd(),
		newVECmd(),
	)
	parent.AddCommand(cmd)
}

// decodeInput accepts either base64 (standard or URL alphabet, with or
// without padding) or 0x-prefixed hex (or bare hex). Returns the raw
// bytes.
//
// Callers pass a label (e.g. "tx") for error context.
func decodeInput(label, raw string) ([]byte, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return nil, &client.UsageError{Msg: fmt.Sprintf("%s: empty input", label)}
	}
	// Hex form: explicit 0x prefix OR all-hex characters of even length
	// long enough to plausibly be hex (>= 2).
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		h := strings.TrimPrefix(strings.TrimPrefix(s, "0x"), "0X")
		b, err := hex.DecodeString(h)
		if err != nil {
			return nil, &client.UsageError{Msg: fmt.Sprintf("%s: invalid hex: %v", label, err)}
		}
		return b, nil
	}
	// Try hex first when the string looks hex (no base64 padding/special
	// chars and length is even). This makes inputs copied from CometBFT
	// logs just work.
	if looksLikeHex(s) {
		if b, err := hex.DecodeString(s); err == nil {
			return b, nil
		}
	}
	// Standard base64 with padding.
	if b, err := base64.StdEncoding.DecodeString(s); err == nil {
		return b, nil
	}
	if b, err := base64.RawStdEncoding.DecodeString(s); err == nil {
		return b, nil
	}
	if b, err := base64.URLEncoding.DecodeString(s); err == nil {
		return b, nil
	}
	if b, err := base64.RawURLEncoding.DecodeString(s); err == nil {
		return b, nil
	}
	return nil, &client.UsageError{Msg: fmt.Sprintf("%s: could not decode as base64 or hex", label)}
}

func looksLikeHex(s string) bool {
	if len(s) == 0 || len(s)%2 != 0 {
		return false
	}
	for _, c := range s {
		ok := (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
		if !ok {
			return false
		}
	}
	return true
}
