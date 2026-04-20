package heimdallutil

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/btcutil/bech32"
	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
)

// DefaultHRP is the bech32 human-readable part used by Heimdall v2.
// Confirmed from heimdall-v2/API_REFERENCE.md ("Cosmos bech32 cosmos1…")
// and the absence of any sdk.Config.SetBech32PrefixForAccount override
// in heimdall-v2/cmd/heimdalld/cmd/commands.go — Heimdall v2 uses the
// default cosmos-sdk account prefix.
const DefaultHRP = "cosmos"

// newAddrCmd builds `util addr <VALUE>` which converts between 0x-hex
// and bech32 representations of a Heimdall address. The direction is
// auto-detected: a `0x`-prefixed value is encoded to bech32, anything
// else is decoded from bech32 to 0x-hex. --all prints both forms.
func newAddrCmd() *cobra.Command {
	var showAll bool
	var hrp string
	cmd := &cobra.Command{
		Use:   "addr <value>",
		Short: "Convert an address between 0x-hex and bech32.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			hexAddr, bechAddr, err := ConvertAddress(args[0], hrp)
			if err != nil {
				return err
			}
			w := cmd.OutOrStdout()
			if showAll {
				fmt.Fprintf(w, "hex=%s\n", hexAddr)
				fmt.Fprintf(w, "bech32=%s\n", bechAddr)
				return nil
			}
			// Default: print the *other* form from what the user supplied.
			if isHexAddressInput(args[0]) {
				fmt.Fprintln(w, bechAddr)
			} else {
				fmt.Fprintln(w, hexAddr)
			}
			return nil
		},
	}
	f := cmd.Flags()
	f.BoolVarP(&showAll, "all", "a", false, "print both hex and bech32 forms")
	f.StringVar(&hrp, "hrp", DefaultHRP, "bech32 human-readable part")
	return cmd
}

// ConvertAddress returns both the canonical 0x-hex and bech32 encodings
// of value, given the target bech32 prefix hrp. Direction is inferred
// from the input: 0x-prefixed hex inputs are encoded; anything else is
// decoded as bech32.
func ConvertAddress(value, hrp string) (hexAddr, bechAddr string, err error) {
	if hrp == "" {
		hrp = DefaultHRP
	}
	value = strings.TrimSpace(value)
	if value == "" {
		return "", "", &client.UsageError{Msg: "address value is empty"}
	}
	if isHexAddressInput(value) {
		raw, err := decodeHexAddress(value)
		if err != nil {
			return "", "", err
		}
		bech, err := encodeBech32(hrp, raw)
		if err != nil {
			return "", "", err
		}
		return "0x" + hex.EncodeToString(raw), bech, nil
	}
	gotHRP, raw, err := decodeBech32(value)
	if err != nil {
		return "", "", err
	}
	if gotHRP != hrp {
		return "", "", &client.UsageError{Msg: fmt.Sprintf(
			"bech32 prefix mismatch: got %q, want %q (override with --hrp)", gotHRP, hrp)}
	}
	// Re-encode with the canonical hrp so the returned bech32 matches
	// the intended output even if the user supplied mixed case.
	bech, err := encodeBech32(hrp, raw)
	if err != nil {
		return "", "", err
	}
	return "0x" + hex.EncodeToString(raw), bech, nil
}

// isHexAddressInput reports whether v looks like a 0x-prefixed hex
// address (case-insensitive 0x plus hex).
func isHexAddressInput(v string) bool {
	s := strings.TrimSpace(v)
	if len(s) < 3 {
		return false
	}
	if s[0] != '0' || (s[1] != 'x' && s[1] != 'X') {
		return false
	}
	for _, c := range s[2:] {
		switch {
		case c >= '0' && c <= '9',
			c >= 'a' && c <= 'f',
			c >= 'A' && c <= 'F':
		default:
			return false
		}
	}
	return true
}

func decodeHexAddress(v string) ([]byte, error) {
	raw, err := hex.DecodeString(strings.TrimPrefix(strings.TrimPrefix(v, "0x"), "0X"))
	if err != nil {
		return nil, &client.UsageError{Msg: fmt.Sprintf("decoding hex address %q: %v", v, err)}
	}
	if len(raw) != 20 {
		return nil, &client.UsageError{Msg: fmt.Sprintf(
			"invalid address length: got %d bytes, want 20", len(raw))}
	}
	return raw, nil
}

// encodeBech32 wraps the 5-bit regrouping + bech32 encoding step.
func encodeBech32(hrp string, raw []byte) (string, error) {
	conv, err := bech32.ConvertBits(raw, 8, 5, true)
	if err != nil {
		return "", fmt.Errorf("regrouping bits: %w", err)
	}
	out, err := bech32.Encode(hrp, conv)
	if err != nil {
		return "", fmt.Errorf("encoding bech32: %w", err)
	}
	return out, nil
}

// decodeBech32 returns the hrp + 20 raw bytes for a bech32 address. The
// bech32 library only enforces its own checksum; we additionally
// validate the decoded byte length because a Heimdall account address
// is always 20 bytes.
func decodeBech32(v string) (string, []byte, error) {
	hrp, conv, err := bech32.Decode(v)
	if err != nil {
		return "", nil, &client.UsageError{Msg: fmt.Sprintf("decoding bech32 %q: %v", v, err)}
	}
	raw, err := bech32.ConvertBits(conv, 5, 8, false)
	if err != nil {
		return "", nil, &client.UsageError{Msg: fmt.Sprintf("regrouping bech32 bits: %v", err)}
	}
	if len(raw) != 20 {
		return "", nil, &client.UsageError{Msg: fmt.Sprintf(
			"invalid bech32 payload length: got %d bytes, want 20", len(raw))}
	}
	return hrp, raw, nil
}
