package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeTempCalldata writes content to a temp file and returns its path.
func writeTempCalldata(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "calldata.txt")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("writing temp calldata file: %v", err)
	}
	return path
}

func TestResolveContractCallData_NoFile(t *testing.T) {
	// Without --calldata-file, ContractCallData is left untouched.
	c := &Config{ContractCallData: "0xdeadbeef"}
	if err := c.resolveContractCallData(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.ContractCallData != "0xdeadbeef" {
		t.Fatalf("ContractCallData changed: got %q", c.ContractCallData)
	}
}

func TestResolveContractCallData_ReadsAndTrims(t *testing.T) {
	// A trailing newline (as written by `... > calldata-file.txt`) must be
	// trimmed so it doesn't corrupt the hex.
	hexNoPrefix := strings.Repeat("ab", 1000)
	c := &Config{ContractCallDataFile: writeTempCalldata(t, hexNoPrefix+"\n")}
	if err := c.resolveContractCallData(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.ContractCallData != hexNoPrefix {
		t.Fatalf("expected trimmed hex, got %q", c.ContractCallData)
	}
}

func TestResolveContractCallData_StripsInteriorWhitespace(t *testing.T) {
	// Hex dumps from `xxd`/`od` are line-wrapped; interior newlines (and any
	// stray spaces) must be stripped, not just the surrounding bytes, so the
	// result decodes as clean hex.
	wrapped := "deadbeef\n" + "cafe babe\n" + "0011\t2233\n"
	want := "deadbeefcafebabe00112233"
	c := &Config{ContractCallDataFile: writeTempCalldata(t, wrapped)}
	if err := c.resolveContractCallData(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.ContractCallData != want {
		t.Fatalf("expected %q, got %q", want, c.ContractCallData)
	}
}

func TestResolveContractCallData_AcceptsHexPrefix(t *testing.T) {
	c := &Config{ContractCallDataFile: writeTempCalldata(t, "0xc0ffee")}
	if err := c.resolveContractCallData(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.ContractCallData != "0xc0ffee" {
		t.Fatalf("expected 0xc0ffee, got %q", c.ContractCallData)
	}
}

func TestResolveContractCallData_LargePayloadExceedsArgLimit(t *testing.T) {
	// The whole point of --calldata-file: carry a payload larger than the OS
	// single-argument limit (~128 KiB) that --calldata cannot.
	hexBig := strings.Repeat("00", 130000) // 260000 hex chars
	c := &Config{ContractCallDataFile: writeTempCalldata(t, hexBig)}
	if err := c.resolveContractCallData(); err != nil {
		t.Fatalf("unexpected error for large payload: %v", err)
	}
	if c.ContractCallData != hexBig {
		t.Fatalf("large payload not loaded correctly (len got %d, want %d)", len(c.ContractCallData), len(hexBig))
	}
}

func TestResolveContractCallData_MutuallyExclusive(t *testing.T) {
	c := &Config{
		ContractCallData:     "0xdeadbeef",
		ContractCallDataFile: writeTempCalldata(t, "0xc0ffee"),
	}
	err := c.resolveContractCallData()
	if err == nil {
		t.Fatal("expected error when both --calldata and --calldata-file are set")
	}
	if !strings.Contains(err.Error(), "mutually exclusive") {
		t.Fatalf("expected mutual-exclusivity error, got: %v", err)
	}
}

func TestResolveContractCallData_EmptyFile(t *testing.T) {
	c := &Config{ContractCallDataFile: writeTempCalldata(t, "  \n\t ")}
	err := c.resolveContractCallData()
	if err == nil || !strings.Contains(err.Error(), "empty") {
		t.Fatalf("expected empty-file error, got: %v", err)
	}
}

func TestResolveContractCallData_PrefixOnlyIsEmpty(t *testing.T) {
	// A file holding only "0x" decodes to zero bytes without a hex error; it
	// must still be rejected as empty (the empty-check runs on the hex body
	// after the 0x prefix is removed).
	c := &Config{ContractCallDataFile: writeTempCalldata(t, "0x\n")}
	err := c.resolveContractCallData()
	if err == nil || !strings.Contains(err.Error(), "empty") {
		t.Fatalf("expected empty-file error for \"0x\"-only file, got: %v", err)
	}
}

func TestResolveContractCallData_CRLF(t *testing.T) {
	// Windows-style CRLF line endings must be stripped like LF.
	c := &Config{ContractCallDataFile: writeTempCalldata(t, "dead\r\nbeef\r\n")}
	if err := c.resolveContractCallData(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.ContractCallData != "deadbeef" {
		t.Fatalf("expected deadbeef, got %q", c.ContractCallData)
	}
}

func TestResolveContractCallData_MissingFile(t *testing.T) {
	c := &Config{ContractCallDataFile: filepath.Join(t.TempDir(), "does-not-exist.txt")}
	err := c.resolveContractCallData()
	if err == nil || !strings.Contains(err.Error(), "unable to read") {
		t.Fatalf("expected read error, got: %v", err)
	}
}

func TestResolveContractCallData_InvalidHex(t *testing.T) {
	c := &Config{ContractCallDataFile: writeTempCalldata(t, "nothexatall")}
	err := c.resolveContractCallData()
	if err == nil || !strings.Contains(err.Error(), "valid hex") {
		t.Fatalf("expected invalid-hex error, got: %v", err)
	}
}

func TestValidate_ResolvesCalldataFile(t *testing.T) {
	// End-to-end through Validate(): a valid --calldata-file is loaded into
	// ContractCallData. The other fields are set to values that pass the
	// unrelated Validate() checks.
	c := &Config{
		AdaptiveBackoffFactor: 2,
		GasPriceMultiplier:    1,
		ContractCallDataFile:  writeTempCalldata(t, "0xabcdef\n"),
	}
	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected Validate error: %v", err)
	}
	if c.ContractCallData != "0xabcdef" {
		t.Fatalf("Validate did not resolve calldata-file, got %q", c.ContractCallData)
	}
}
