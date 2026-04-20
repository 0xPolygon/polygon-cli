package wallet

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
)

// fixturePath returns the absolute path to a testdata file so the
// tests do not depend on the caller's working directory.
func fixturePath(t *testing.T, name string) string {
	t.Helper()
	_, thisFile, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(thisFile), "testdata", name)
}

// runWallet executes the wallet umbrella command with args against a
// temporary HOME/keystore directory so the test never touches the
// operator's ~/.foundry or ~/.polycli. Returns stdout/stderr/err.
func runWallet(t *testing.T, keystoreDir string, stdin io.Reader, args ...string) (string, string, error) {
	t.Helper()
	// Sanitise environment so tests are hermetic.
	t.Setenv("ETH_KEYSTORE", "")
	// Force HOME to a temp dir so the default fallback can't pick
	// up a real ~/.foundry or ~/.polycli.
	t.Setenv("HOME", t.TempDir())

	root := &cobra.Command{Use: "h", SilenceUsage: true}
	f := &config.Flags{}
	f.Register(root)
	Register(root, f)
	var stdout, stderr bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stderr)
	if stdin != nil {
		root.SetIn(stdin)
	}
	full := append([]string{"wallet"}, args...)
	if keystoreDir != "" && subcommandTakesKeystoreDir(args) {
		full = append(full, "--keystore-dir", keystoreDir)
	}
	root.SetArgs(full)
	err := root.ExecuteContext(context.Background())
	return stdout.String(), stderr.String(), err
}

// TestNewCreatesKey exercises `wallet new`: we pass a password via
// --password so the command can run non-interactively, then verify
// the keystore directory contains a single account with a 0x address.
func TestNewCreatesKey(t *testing.T) {
	ksDir := t.TempDir()
	stdout, _, err := runWallet(t, ksDir, nil, "new", "--password", "test")
	if err != nil {
		t.Fatalf("wallet new: %v", err)
	}
	if !strings.Contains(stdout, "address") || !strings.Contains(stdout, "keyfile") {
		t.Fatalf("missing address/keyfile in output:\n%s", stdout)
	}
	// Directory should contain exactly one UTC-- file.
	entries, err := os.ReadDir(ksDir)
	if err != nil {
		t.Fatalf("read keystore dir: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 keyfile, got %d", len(entries))
	}
	if !strings.HasPrefix(entries[0].Name(), "UTC--") {
		t.Errorf("keyfile name %q not UTC-- prefixed", entries[0].Name())
	}
}

// TestListEmpty exercises `wallet list` against a fresh directory.
func TestListEmpty(t *testing.T) {
	ksDir := t.TempDir()
	stdout, _, err := runWallet(t, ksDir, nil, "list")
	if err != nil {
		t.Fatalf("wallet list: %v", err)
	}
	if !strings.Contains(stdout, "no keys") {
		t.Errorf("expected 'no keys' message, got:\n%s", stdout)
	}
}

// TestRoundTripCreateDecryptReEncrypt verifies create -> decrypt ->
// re-encrypt keeps the address stable.
func TestRoundTripCreateDecryptReEncrypt(t *testing.T) {
	ksDir := t.TempDir()
	stdout, _, err := runWallet(t, ksDir, nil, "new", "--password", "original")
	if err != nil {
		t.Fatalf("wallet new: %v", err)
	}
	addr := extractAddress(t, stdout)

	// List and confirm.
	stdout, _, err = runWallet(t, ksDir, nil, "list")
	if err != nil {
		t.Fatalf("wallet list: %v", err)
	}
	if !strings.Contains(stdout, addr) {
		t.Fatalf("list missing address %s:\n%s", addr, stdout)
	}

	// Change password.
	_, _, err = runWallet(t, ksDir, nil, "change-password", addr,
		"--password", "original", "--new-password", "updated")
	if err != nil {
		t.Fatalf("change-password: %v", err)
	}
	// Address should remain the same after list.
	stdout, _, err = runWallet(t, ksDir, nil, "list")
	if err != nil {
		t.Fatalf("wallet list (post-change): %v", err)
	}
	if !strings.Contains(stdout, addr) {
		t.Fatalf("address lost after password change:\n%s", stdout)
	}

	// Decrypt with the new password works; old password fails.
	_, _, err = runWallet(t, ksDir, nil, "private-key", addr,
		"--password", "original", "--i-understand-the-risks")
	if err == nil {
		t.Fatal("expected old password to fail after change-password")
	}
	stdout, _, err = runWallet(t, ksDir, nil, "private-key", addr,
		"--password", "updated", "--i-understand-the-risks")
	if err != nil {
		t.Fatalf("private-key with new password: %v", err)
	}
	if !strings.HasPrefix(strings.TrimSpace(stdout), "0x") || len(strings.TrimSpace(stdout)) != 66 {
		t.Errorf("private-key output %q malformed", strings.TrimSpace(stdout))
	}
}

// TestFoundryKeystoreCompat decrypts a foundry-format keystore
// fixture and checks the derived address matches the expected value.
func TestFoundryKeystoreCompat(t *testing.T) {
	const expectedAddr = "0x7E5F4552091A69125d5DfCb7b8C2659029395Bdf"
	const expectedPriv = "0x0000000000000000000000000000000000000000000000000000000000000001"
	fixtureDir := fixturePath(t, "foundry")
	ksDir := t.TempDir()
	// Copy the fixture into our test keystore dir so ks.Accounts() sees it.
	entries, err := os.ReadDir(fixtureDir)
	if err != nil {
		t.Fatalf("read fixture dir: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("no foundry fixture files")
	}
	for _, e := range entries {
		data, err := os.ReadFile(filepath.Join(fixtureDir, e.Name()))
		if err != nil {
			t.Fatalf("read fixture file: %v", err)
		}
		if err := os.WriteFile(filepath.Join(ksDir, e.Name()), data, 0o600); err != nil {
			t.Fatalf("copy fixture file: %v", err)
		}
	}

	stdout, _, err := runWallet(t, ksDir, nil, "list", "--addresses-only")
	if err != nil {
		t.Fatalf("wallet list: %v", err)
	}
	if !strings.Contains(stdout, expectedAddr) {
		t.Errorf("expected %s in list:\n%s", expectedAddr, stdout)
	}

	// Export via private-key.
	stdout, _, err = runWallet(t, ksDir, nil, "private-key", expectedAddr,
		"--password", "test", "--i-understand-the-risks")
	if err != nil {
		t.Fatalf("private-key: %v", err)
	}
	got := strings.TrimSpace(stdout)
	if got != expectedPriv {
		t.Fatalf("private-key = %q, want %q", got, expectedPriv)
	}

	// Decrypt-keystore against the file path also works.
	entries, _ = os.ReadDir(ksDir)
	filePath := filepath.Join(ksDir, entries[0].Name())
	stdout, _, err = runWallet(t, ksDir, nil, "decrypt-keystore", filePath,
		"--password", "test", "--i-understand-the-risks")
	if err != nil {
		t.Fatalf("decrypt-keystore: %v", err)
	}
	if strings.TrimSpace(stdout) != expectedPriv {
		t.Fatalf("decrypt-keystore = %q, want %q", strings.TrimSpace(stdout), expectedPriv)
	}
}

// TestMnemonicKnownVector verifies that the canonical BIP-39 zero
// mnemonic derives the well-known address at m/44'/60'/0'/0/0.
// This is the same vector documented by BIP-39 reference wallets and
// cast / ethers / viem all produce the same result.
func TestMnemonicKnownVector(t *testing.T) {
	const mnemonic = "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	const expectedAddr = "0x9858EfFD232B4033E47d90003D41EC34EcaEda94"

	priv, path, addr, err := deriveFromMnemonic(mnemonic, "", "", 0)
	if err != nil {
		t.Fatalf("deriveFromMnemonic: %v", err)
	}
	if path != "m/44'/60'/0'/0/0" {
		t.Errorf("path = %q, want m/44'/60'/0'/0/0", path)
	}
	if addr.Hex() != expectedAddr {
		t.Errorf("addr = %s, want %s", addr.Hex(), expectedAddr)
	}
	if priv == nil {
		t.Fatal("nil private key")
	}

	// Via the command, with --print-only.
	stdout, _, err := runWallet(t, t.TempDir(), nil, "derive",
		"--mnemonic", mnemonic, "--count", "1")
	if err != nil {
		t.Fatalf("wallet derive: %v", err)
	}
	if !strings.Contains(stdout, expectedAddr) {
		t.Errorf("expected %s in output:\n%s", expectedAddr, stdout)
	}
}

// TestSignVerifyRoundTrip signs a message with a known private key
// and verifies it, both via the internal helpers and end-to-end.
func TestSignVerifyRoundTrip(t *testing.T) {
	// Well-known private key 0x...01 -> 0x7E5F4552...
	const privHex = "0000000000000000000000000000000000000000000000000000000000000001"
	const expectedAddr = "0x7E5F4552091A69125d5DfCb7b8C2659029395Bdf"
	const message = "hello world"

	priv, err := parsePrivateKeyHex("0x" + privHex)
	if err != nil {
		t.Fatalf("parsePrivateKeyHex: %v", err)
	}
	sig, err := signPersonal(priv, []byte(message))
	if err != nil {
		t.Fatalf("signPersonal: %v", err)
	}
	if len(sig) != 65 {
		t.Fatalf("sig length = %d, want 65", len(sig))
	}
	// v must be 27 or 28 per our normalisation.
	if sig[64] != 27 && sig[64] != 28 {
		t.Errorf("v byte = %d, want 27 or 28", sig[64])
	}

	addr := crypto.PubkeyToAddress(priv.PublicKey).Hex()
	if addr != expectedAddr {
		t.Fatalf("address derivation broken: %s vs %s", addr, expectedAddr)
	}

	// Round-trip via helpers.
	ok, err := verifyPersonal(common.HexToAddress(expectedAddr), []byte(message), sig)
	if err != nil {
		t.Fatalf("verifyPersonal: %v", err)
	}
	if !ok {
		t.Fatalf("verifyPersonal failed for self-signed message")
	}

	// End-to-end via command.
	stdout, _, err := runWallet(t, t.TempDir(), nil, "sign", message,
		"--private-key", privHex)
	if err != nil {
		t.Fatalf("wallet sign: %v", err)
	}
	sigHex := strings.TrimSpace(stdout)
	if !strings.HasPrefix(sigHex, "0x") {
		t.Fatalf("sig hex missing 0x: %q", sigHex)
	}
	// Verify via the command.
	_, _, err = runWallet(t, t.TempDir(), nil, "verify", expectedAddr, message, sigHex)
	if err != nil {
		t.Fatalf("wallet verify: %v", err)
	}

	// Tamper with the message and ensure verify fails.
	_, _, err = runWallet(t, t.TempDir(), nil, "verify", expectedAddr, "tampered", sigHex)
	if err == nil {
		t.Fatal("verify should have failed for tampered message")
	}
	var uErr *client.UsageError
	if !errors.As(err, &uErr) {
		t.Fatalf("expected *UsageError, got %T", err)
	}
}

// TestSignRawHash exercises the --raw path.
func TestSignRawHash(t *testing.T) {
	const privHex = "0000000000000000000000000000000000000000000000000000000000000001"
	const expectedAddr = "0x7E5F4552091A69125d5DfCb7b8C2659029395Bdf"
	// 32 bytes of 0xff.
	hashHex := "0x" + strings.Repeat("ff", 32)
	stdout, _, err := runWallet(t, t.TempDir(), nil, "sign", hashHex,
		"--private-key", privHex, "--raw")
	if err != nil {
		t.Fatalf("wallet sign --raw: %v", err)
	}
	sigHex := strings.TrimSpace(stdout)
	_, _, err = runWallet(t, t.TempDir(), nil, "verify", expectedAddr, hashHex, sigHex, "--raw")
	if err != nil {
		t.Fatalf("wallet verify --raw: %v", err)
	}
}

// TestImportAndAddress uses `wallet import --private-key` then
// exercises `wallet address <addr>` / `wallet list`.
func TestImportAndAddress(t *testing.T) {
	ksDir := t.TempDir()
	const priv = "0x0000000000000000000000000000000000000000000000000000000000000001"
	const expectedAddr = "0x7E5F4552091A69125d5DfCb7b8C2659029395Bdf"
	stdout, _, err := runWallet(t, ksDir, nil, "import",
		"--private-key", priv, "--password", "pw")
	if err != nil {
		t.Fatalf("import: %v", err)
	}
	if !strings.Contains(stdout, expectedAddr) {
		t.Errorf("expected %s in output:\n%s", expectedAddr, stdout)
	}

	stdout, _, err = runWallet(t, ksDir, nil, "address")
	if err != nil {
		t.Fatalf("address: %v", err)
	}
	if !strings.Contains(stdout, expectedAddr) {
		t.Errorf("expected %s in address output:\n%s", expectedAddr, stdout)
	}

	// Check private-key via --private-key flag (no keystore).
	stdout, _, err = runWallet(t, t.TempDir(), nil, "address",
		"--private-key", priv)
	if err != nil {
		t.Fatalf("address --private-key: %v", err)
	}
	if strings.TrimSpace(stdout) != expectedAddr {
		t.Errorf("expected address %s, got %q", expectedAddr, stdout)
	}
}

// TestRemove exercises the deletion path, including the --yes flag
// bypassing the confirm prompt.
func TestRemove(t *testing.T) {
	ksDir := t.TempDir()
	const priv = "0x0000000000000000000000000000000000000000000000000000000000000001"
	const expectedAddr = "0x7E5F4552091A69125d5DfCb7b8C2659029395Bdf"
	_, _, err := runWallet(t, ksDir, nil, "import",
		"--private-key", priv, "--password", "pw")
	if err != nil {
		t.Fatalf("import: %v", err)
	}

	// Without --yes and piping an "n" should abort.
	_, _, err = runWallet(t, ksDir, strings.NewReader("n\n"), "remove", expectedAddr,
		"--password", "pw")
	if err == nil {
		t.Fatal("expected abort error when answering n")
	}
	var uErr *client.UsageError
	if !errors.As(err, &uErr) {
		t.Errorf("err type = %T, want *UsageError", err)
	}

	// Now with --yes.
	_, _, err = runWallet(t, ksDir, nil, "remove", expectedAddr,
		"--password", "pw", "--yes")
	if err != nil {
		t.Fatalf("remove --yes: %v", err)
	}
	// Directory is now empty.
	entries, _ := os.ReadDir(ksDir)
	if len(entries) != 0 {
		t.Errorf("expected empty keystore dir, got %d entries", len(entries))
	}
}

// TestPublicKey exercises the public-key command against an in-flight
// --private-key, avoiding any keystore round-trip.
func TestPublicKey(t *testing.T) {
	const priv = "0x0000000000000000000000000000000000000000000000000000000000000001"
	stdout, _, err := runWallet(t, t.TempDir(), nil, "public-key",
		"--private-key", priv)
	if err != nil {
		t.Fatalf("public-key: %v", err)
	}
	if !strings.Contains(stdout, "uncompressed") || !strings.Contains(stdout, "compressed") {
		t.Errorf("expected both pub key forms:\n%s", stdout)
	}
	// Uncompressed for 0x01 is 0x04 || Gx || Gy (65 bytes = 130 hex).
	// Compressed starts with 0x02 or 0x03.
	for _, line := range strings.Split(stdout, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "uncompressed") {
			parts := strings.Fields(line)
			if len(parts) != 2 {
				t.Fatalf("uncompressed line malformed: %q", line)
			}
			raw, err := hex.DecodeString(strings.TrimPrefix(parts[1], "0x"))
			if err != nil || len(raw) != 65 || raw[0] != 0x04 {
				t.Errorf("uncompressed = %q, want 65-byte 0x04-prefixed", parts[1])
			}
		}
		if strings.HasPrefix(line, "compressed") {
			parts := strings.Fields(line)
			if len(parts) != 2 {
				t.Fatalf("compressed line malformed: %q", line)
			}
			raw, err := hex.DecodeString(strings.TrimPrefix(parts[1], "0x"))
			if err != nil || len(raw) != 33 || (raw[0] != 0x02 && raw[0] != 0x03) {
				t.Errorf("compressed = %q, want 33-byte 0x02/0x03-prefixed", parts[1])
			}
		}
	}
}

// TestPrivateKeyRequiresAck verifies the friction flag is enforced.
func TestPrivateKeyRequiresAck(t *testing.T) {
	ksDir := t.TempDir()
	const priv = "0x0000000000000000000000000000000000000000000000000000000000000001"
	const expectedAddr = "0x7E5F4552091A69125d5DfCb7b8C2659029395Bdf"
	_, _, err := runWallet(t, ksDir, nil, "import",
		"--private-key", priv, "--password", "pw")
	if err != nil {
		t.Fatalf("import: %v", err)
	}
	_, _, err = runWallet(t, ksDir, nil, "private-key", expectedAddr, "--password", "pw")
	if err == nil {
		t.Fatal("expected friction-flag error")
	}
	var uErr *client.UsageError
	if !errors.As(err, &uErr) {
		t.Fatalf("err type = %T, want *UsageError", err)
	}
	if !strings.Contains(uErr.Msg, "i-understand-the-risks") {
		t.Errorf("error message missing friction flag name: %q", uErr.Msg)
	}
}

// TestRejectedFlagsAndSubcommands ensures hardware flags and the
// dropped subcommands surface the intended error.
func TestRejectedFlagsAndSubcommands(t *testing.T) {
	// --ledger on a subcommand.
	_, _, err := runWallet(t, t.TempDir(), nil, "new", "--ledger")
	if err == nil {
		t.Fatal("expected error for --ledger")
	}
	if !strings.Contains(err.Error(), "hardware wallets") {
		t.Errorf("expected hardware-wallets message, got %q", err.Error())
	}
	// vanity subcommand.
	_, _, err = runWallet(t, t.TempDir(), nil, "vanity", "--starts-with", "0xabc")
	if err == nil {
		t.Fatal("expected error for vanity")
	}
	if !strings.Contains(err.Error(), "cast wallet vanity") {
		t.Errorf("expected vanity pointer, got %q", err.Error())
	}
	// sign-auth subcommand.
	_, _, err = runWallet(t, t.TempDir(), nil, "sign-auth", "0xabc")
	if err == nil {
		t.Fatal("expected error for sign-auth")
	}
	if !strings.Contains(err.Error(), "cast wallet sign-auth") {
		t.Errorf("expected sign-auth pointer, got %q", err.Error())
	}
}

// TestResolveKeystoreDirPrecedence exercises the fallback chain.
// Each case uses t.Setenv + a fabricated HOME and confirms which
// directory is chosen.
func TestResolveKeystoreDirPrecedence(t *testing.T) {
	t.Run("flag wins", func(t *testing.T) {
		dir, err := resolveKeystoreDir("/tmp/custom-ks")
		if err != nil {
			t.Fatalf("resolve: %v", err)
		}
		if dir != "/tmp/custom-ks" {
			t.Errorf("flag path = %q, want /tmp/custom-ks", dir)
		}
	})
	t.Run("env wins over foundry/polycli", func(t *testing.T) {
		home := t.TempDir()
		t.Setenv("HOME", home)
		// Create a foundry dir that should lose to env.
		_ = os.MkdirAll(filepath.Join(home, ".foundry", "keystores"), 0o700)
		envDir := t.TempDir()
		t.Setenv("ETH_KEYSTORE", envDir)
		got, err := resolveKeystoreDir("")
		if err != nil {
			t.Fatalf("resolve: %v", err)
		}
		if got != envDir {
			t.Errorf("env dir = %q, want %q", got, envDir)
		}
	})
	t.Run("foundry wins over polycli", func(t *testing.T) {
		home := t.TempDir()
		t.Setenv("HOME", home)
		t.Setenv("ETH_KEYSTORE", "")
		foundryDir := filepath.Join(home, ".foundry", "keystores")
		_ = os.MkdirAll(foundryDir, 0o700)
		got, err := resolveKeystoreDir("")
		if err != nil {
			t.Fatalf("resolve: %v", err)
		}
		if got != foundryDir {
			t.Errorf("got %q, want %q", got, foundryDir)
		}
	})
	t.Run("polycli default when nothing else", func(t *testing.T) {
		home := t.TempDir()
		t.Setenv("HOME", home)
		t.Setenv("ETH_KEYSTORE", "")
		got, err := resolveKeystoreDir("")
		if err != nil {
			t.Fatalf("resolve: %v", err)
		}
		want := filepath.Join(home, ".polycli", "keystores")
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
		// And the dir should exist — resolve creates it.
		if st, err := os.Stat(got); err != nil || !st.IsDir() {
			t.Errorf("default dir not created: %v", err)
		}
	})
}

// TestParseDerivationPathBad covers a few malformed paths.
func TestParseDerivationPathBad(t *testing.T) {
	cases := []string{
		"",
		"not/m",
		"m/44'/60'//0",
		"m/-1",
	}
	for _, p := range cases {
		if _, err := parseDerivationPath(p); err == nil {
			t.Errorf("expected error for %q", p)
		}
	}
}

// --- helpers ---

// subcommandTakesKeystoreDir returns true if the first positional
// arg in args is a subcommand that accepts --keystore-dir. Most do;
// `derive` and `verify` do not.
func subcommandTakesKeystoreDir(args []string) bool {
	if len(args) == 0 {
		return false
	}
	switch args[0] {
	case "derive", "verify", "vanity", "sign-auth":
		return false
	}
	return true
}

func extractAddress(t *testing.T, out string) string {
	t.Helper()
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "address") {
			fields := strings.Fields(line)
			if len(fields) == 2 {
				return fields[1]
			}
		}
	}
	t.Fatalf("no address line in:\n%s", out)
	return ""
}
