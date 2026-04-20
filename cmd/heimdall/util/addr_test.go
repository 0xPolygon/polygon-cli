package heimdallutil

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
)

// Known-good address pairs. The hex values are canonical lower-case
// 20-byte Heimdall validator signers (lifted from the checked-in
// testdata fixtures); the bech32 forms are derived locally using the
// default `cosmos` prefix and verified by round-tripping.
var addrVectors = []struct {
	name string
	hex  string
	bech string
}{
	{
		name: "validator_02f615",
		hex:  "0x02f615e95563ef16f10354dba9e584e58d2d4314",
		bech: "cosmos1qtmpt624v0h3dugr2nd6nevyukxj6sc54tvenp",
	},
	{
		name: "root_chain_contract",
		hex:  "0xbd07d7e1e93c8d4b2a261327f3c28a8ea7167209",
		bech: "cosmos1h5ra0c0f8jx5k23xzvnl8s5236n3vusfjkjpwr",
	},
	{
		name: "checkpoints_proposer",
		hex:  "0x4ad84f7014b7b44f723f284a85b1662337971439",
		bech: "cosmos1ftvy7uq5k76y7u3l9p9gtvtxyvmew9pee27vjt",
	},
}

func TestConvertAddressRoundTrip(t *testing.T) {
	for _, tc := range addrVectors {
		t.Run(tc.name, func(t *testing.T) {
			gotHex, gotBech, err := ConvertAddress(tc.hex, DefaultHRP)
			if err != nil {
				t.Fatalf("hex->: %v", err)
			}
			if gotHex != tc.hex {
				t.Errorf("hex: got %q want %q", gotHex, tc.hex)
			}
			if gotBech != tc.bech {
				t.Errorf("bech32: got %q want %q", gotBech, tc.bech)
			}
			// Reverse direction.
			revHex, revBech, err := ConvertAddress(tc.bech, DefaultHRP)
			if err != nil {
				t.Fatalf("bech32->: %v", err)
			}
			if revHex != tc.hex {
				t.Errorf("reverse hex: got %q want %q", revHex, tc.hex)
			}
			if revBech != tc.bech {
				t.Errorf("reverse bech32: got %q want %q", revBech, tc.bech)
			}
		})
	}
}

// TestConvertAddressCaseInsensitive ensures mixed-case hex input is
// normalized to lower-case and that the bech32 output is unaffected.
func TestConvertAddressCaseInsensitive(t *testing.T) {
	mixed := "0x02F615E95563eF16F10354Dba9e584E58D2D4314"
	gotHex, gotBech, err := ConvertAddress(mixed, DefaultHRP)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	if gotHex != strings.ToLower(mixed) {
		t.Errorf("hex: got %q want %q", gotHex, strings.ToLower(mixed))
	}
	if gotBech != "cosmos1qtmpt624v0h3dugr2nd6nevyukxj6sc54tvenp" {
		t.Errorf("bech32 mismatch: %q", gotBech)
	}
}

func TestConvertAddressCustomHRP(t *testing.T) {
	raw := "0x02f615e95563ef16f10354dba9e584e58d2d4314"
	_, bech, err := ConvertAddress(raw, "heimdall")
	if err != nil {
		t.Fatalf("convert with hrp heimdall: %v", err)
	}
	if !strings.HasPrefix(bech, "heimdall1") {
		t.Errorf("expected heimdall1 prefix, got %q", bech)
	}
	// Reverse with the same hrp must succeed.
	_, _, err = ConvertAddress(bech, "heimdall")
	if err != nil {
		t.Fatalf("reverse: %v", err)
	}
	// Reverse with the default hrp must surface a prefix mismatch.
	_, _, err = ConvertAddress(bech, DefaultHRP)
	if err == nil {
		t.Fatal("expected prefix mismatch error")
	}
	var uErr *client.UsageError
	if !errors.As(err, &uErr) {
		t.Errorf("got %T, want *UsageError", err)
	}
}

func TestConvertAddressInvalid(t *testing.T) {
	cases := []struct {
		name string
		in   string
	}{
		{"empty", ""},
		{"spaces_only", "   "},
		{"hex_wrong_length", "0x1234"},
		{"hex_non_hex_chars", "0xZZe95563ef16f10354dba9e584e58d2d4314xx"},
		{"bech32_gibberish", "not-a-bech32-address"},
		{"bech32_bad_checksum", "cosmos1qtmpt624v0h3dugr2nd6nevyukxj6sc54tvenq"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := ConvertAddress(tc.in, DefaultHRP)
			if err == nil {
				t.Fatalf("expected error for %q", tc.in)
			}
		})
	}
}

// TestAddrCmdDefaultPrintsOpposite asserts that without --all the
// command prints only the *other* encoding from what was supplied.
func TestAddrCmdDefaultPrintsOpposite(t *testing.T) {
	// hex -> expect bech32
	stdout := runAddr(t, "0x02f615e95563ef16f10354dba9e584e58d2d4314")
	if strings.TrimSpace(stdout) != "cosmos1qtmpt624v0h3dugr2nd6nevyukxj6sc54tvenp" {
		t.Errorf("expected bech32, got %q", stdout)
	}
	// bech32 -> expect hex
	stdout = runAddr(t, "cosmos1qtmpt624v0h3dugr2nd6nevyukxj6sc54tvenp")
	if strings.TrimSpace(stdout) != "0x02f615e95563ef16f10354dba9e584e58d2d4314" {
		t.Errorf("expected hex, got %q", stdout)
	}
}

func TestAddrCmdAllPrintsBoth(t *testing.T) {
	stdout := runAddr(t, "--all", "0x02f615e95563ef16f10354dba9e584e58d2d4314")
	if !strings.Contains(stdout, "hex=0x02f615e95563ef16f10354dba9e584e58d2d4314") {
		t.Errorf("missing hex= line: %s", stdout)
	}
	if !strings.Contains(stdout, "bech32=cosmos1qtmpt624v0h3dugr2nd6nevyukxj6sc54tvenp") {
		t.Errorf("missing bech32= line: %s", stdout)
	}
}

func runAddr(t *testing.T, args ...string) string {
	t.Helper()
	cmd := newAddrCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs(args)
	if err := cmd.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("addr: %v", err)
	}
	return buf.String()
}

// ensure the raw cobra command rejects wrong arg counts (defense-in-depth
// against future refactors that might relax ExactArgs).
func TestAddrCmdArgCount(t *testing.T) {
	cmd := newAddrCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{})
	if err := cmd.ExecuteContext(context.Background()); err == nil {
		t.Fatal("expected error for zero args")
	}
}
