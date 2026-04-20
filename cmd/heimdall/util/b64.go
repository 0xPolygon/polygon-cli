package heimdallutil

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
)

// Conversion directions for the b64 subcommand.
const (
	directionAuto   = "auto"
	directionHex    = "hex"
	directionBase64 = "base64"
)

// newB64Cmd builds `util b64 <VALUE>` — convert between base64 and
// 0x-hex. Auto-detects direction: a 0x-prefixed input is treated as
// hex and encoded to base64; anything else is treated as base64 and
// decoded to 0x-hex. --to overrides.
func newB64Cmd() *cobra.Command {
	var direction string
	cmd := &cobra.Command{
		Use:   "b64 <value>",
		Short: "Convert between base64 and 0x-hex.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			out, err := ConvertBase64(args[0], direction)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), out)
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVar(&direction, "to", directionAuto, "target format (auto|hex|base64)")
	return cmd
}

// ConvertBase64 returns the converted representation of value according
// to direction. When direction is "auto" a 0x-prefixed input is
// converted to base64; anything else is converted from base64 to
// 0x-hex. When direction is explicit, the input must be in the
// opposite format.
func ConvertBase64(value, direction string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", &client.UsageError{Msg: "value is empty"}
	}
	target := direction
	if target == "" {
		target = directionAuto
	}
	switch target {
	case directionAuto:
		if strings.HasPrefix(value, "0x") || strings.HasPrefix(value, "0X") {
			return hexToBase64(value)
		}
		return base64ToHex(value)
	case directionBase64:
		return hexToBase64(value)
	case directionHex:
		return base64ToHex(value)
	default:
		return "", &client.UsageError{Msg: fmt.Sprintf(
			"invalid --to value %q (want auto|hex|base64)", direction)}
	}
}

func hexToBase64(v string) (string, error) {
	trimmed := strings.TrimPrefix(strings.TrimPrefix(v, "0x"), "0X")
	raw, err := hex.DecodeString(trimmed)
	if err != nil {
		return "", &client.UsageError{Msg: fmt.Sprintf("decoding hex %q: %v", v, err)}
	}
	return base64.StdEncoding.EncodeToString(raw), nil
}

func base64ToHex(v string) (string, error) {
	// Accept either std or URL-safe base64 to save users a footgun.
	raw, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		// Try URL-safe before giving up.
		if urlRaw, urlErr := base64.URLEncoding.DecodeString(v); urlErr == nil {
			return "0x" + hex.EncodeToString(urlRaw), nil
		}
		return "", &client.UsageError{Msg: fmt.Sprintf("decoding base64 %q: %v", v, err)}
	}
	return "0x" + hex.EncodeToString(raw), nil
}
