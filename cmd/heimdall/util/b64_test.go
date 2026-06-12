package heimdallutil

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
)

func TestConvertBase64Auto(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"hex_prefixed_to_base64", "0x01020304", "AQIDBA=="},
		{"upper_hex_to_base64", "0X01020304", "AQIDBA=="},
		{"base64_to_hex", "AQIDBA==", "0x01020304"},
		{"base64_urlsafe_to_hex", "_-4=", "0xffee"},
		{"empty_bytes_hex", "0x", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ConvertBase64(tc.in, directionAuto)
			if err != nil {
				t.Fatalf("convert: %v", err)
			}
			if got != tc.want {
				t.Errorf("got %q want %q", got, tc.want)
			}
		})
	}
}

func TestConvertBase64ExplicitDirection(t *testing.T) {
	// Same input "AQIDBA==" decoded via --to hex becomes 0x01020304.
	got, err := ConvertBase64("AQIDBA==", directionHex)
	if err != nil {
		t.Fatalf("to hex: %v", err)
	}
	if got != "0x01020304" {
		t.Errorf("to hex: got %q want 0x01020304", got)
	}
	// Input "0x01020304" encoded via --to base64 becomes "AQIDBA==".
	got, err = ConvertBase64("0x01020304", directionBase64)
	if err != nil {
		t.Fatalf("to base64: %v", err)
	}
	if got != "AQIDBA==" {
		t.Errorf("to base64: got %q want AQIDBA==", got)
	}
}

func TestConvertBase64Invalid(t *testing.T) {
	cases := []struct {
		name      string
		in        string
		direction string
	}{
		{"empty", "", directionAuto},
		{"bad_hex_chars", "0xZZZZ", directionAuto},
		{"odd_hex_nibbles", "0x123", directionAuto},
		{"bad_base64", "not!base64!", directionAuto},
		{"unknown_direction", "AQIDBA==", "gibberish"},
		{"explicit_base64_but_value_not_hex", "AQIDBA==", directionBase64},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ConvertBase64(tc.in, tc.direction)
			if err == nil {
				t.Fatalf("expected error for %q (direction %q)", tc.in, tc.direction)
			}
			// empty + unknown_direction must surface as UsageError.
			if tc.name == "empty" || tc.name == "unknown_direction" {
				var uErr *client.UsageError
				if !errors.As(err, &uErr) {
					t.Errorf("got %T, want *UsageError", err)
				}
			}
		})
	}
}

func TestConvertBase64RoundTrip(t *testing.T) {
	originals := []string{
		"0x00",
		"0xdeadbeef",
		"0x" + strings.Repeat("ff", 32),
	}
	for _, in := range originals {
		b64, err := ConvertBase64(in, directionAuto)
		if err != nil {
			t.Fatalf("to b64 %q: %v", in, err)
		}
		back, err := ConvertBase64(b64, directionAuto)
		if err != nil {
			t.Fatalf("back from b64 %q: %v", b64, err)
		}
		if back != in {
			t.Errorf("round-trip mismatch: got %q want %q", back, in)
		}
	}
}

func TestB64CmdAuto(t *testing.T) {
	stdout := runB64(t, "AQIDBA==")
	if strings.TrimSpace(stdout) != "0x01020304" {
		t.Errorf("got %q", stdout)
	}
}

func TestB64CmdExplicit(t *testing.T) {
	stdout := runB64(t, "--to", "base64", "0x01020304")
	if strings.TrimSpace(stdout) != "AQIDBA==" {
		t.Errorf("got %q", stdout)
	}
}

func runB64(t *testing.T, args ...string) string {
	t.Helper()
	cmd := newB64Cmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs(args)
	if err := cmd.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("b64: %v", err)
	}
	return buf.String()
}
