package wallet

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestDeriveFromMnemonicCanonicalVector asserts the default Ethereum
// derivation path produces the well-known address for the canonical
// "test test test ..." mnemonic used across the Ethereum ecosystem.
func TestDeriveFromMnemonicCanonicalVector(t *testing.T) {
	const mnemonic = "test test test test test test test test test test test junk"
	const wantAddr = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	priv, path, addr, err := DeriveFromMnemonic(mnemonic, "", "", 0)
	if err != nil {
		t.Fatalf("DeriveFromMnemonic: %v", err)
	}
	if priv == nil {
		t.Fatal("priv is nil")
	}
	if path != "m/44'/60'/0'/0/0" {
		t.Fatalf("path = %q, want m/44'/60'/0'/0/0", path)
	}
	if !strings.EqualFold(addr.Hex(), wantAddr) {
		t.Fatalf("addr = %s, want %s", addr.Hex(), wantAddr)
	}
}

// TestDeriveFromMnemonicIndexBumpsFinalSegment asserts that an empty
// explicit path + a non-zero index produces the path with the final
// segment replaced, matching cast's --mnemonic-index.
func TestDeriveFromMnemonicIndexBumpsFinalSegment(t *testing.T) {
	const mnemonic = "test test test test test test test test test test test junk"
	_, path, _, err := DeriveFromMnemonic(mnemonic, "", "", 3)
	if err != nil {
		t.Fatalf("DeriveFromMnemonic: %v", err)
	}
	if path != "m/44'/60'/0'/0/3" {
		t.Fatalf("path = %q, want m/44'/60'/0'/0/3", path)
	}
}

// TestParseDerivationPathRejections covers the malformed-path branches.
func TestParseDerivationPathRejections(t *testing.T) {
	for _, p := range []string{
		"",
		"44'/60'/0'/0/0",   // missing leading m
		"m/",               // empty segment
		"m/44'/zz'/0'/0/0", // non-numeric
	} {
		p := p
		t.Run(p, func(t *testing.T) {
			if _, err := ParseDerivationPath(p); err == nil {
				t.Fatalf("expected error for %q", p)
			}
		})
	}
}

// TestResolveKeystoreDirFlagWins asserts that the explicit override
// beats every other source and is made absolute.
func TestResolveKeystoreDirFlagWins(t *testing.T) {
	t.Setenv("ETH_KEYSTORE", "/should/be/ignored")
	dir, err := ResolveKeystoreDir("/tmp/custom-ks", false)
	if err != nil {
		t.Fatalf("ResolveKeystoreDir: %v", err)
	}
	if dir != "/tmp/custom-ks" {
		t.Fatalf("dir = %q, want /tmp/custom-ks", dir)
	}
}

// TestResolveKeystoreDirEnv asserts that ETH_KEYSTORE is honoured when
// the override is empty.
func TestResolveKeystoreDirEnv(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("ETH_KEYSTORE", tmp)
	t.Setenv("HOME", filepath.Join(tmp, "no-such-home"))
	got, err := ResolveKeystoreDir("", false)
	if err != nil {
		t.Fatalf("ResolveKeystoreDir: %v", err)
	}
	if got != tmp {
		t.Fatalf("dir = %q, want %q", got, tmp)
	}
}

// TestResolveKeystoreDirFoundryExists asserts that ~/.foundry/keystores
// is preferred when it already exists.
func TestResolveKeystoreDirFoundryExists(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("ETH_KEYSTORE", "")
	foundry := filepath.Join(home, ".foundry", "keystores")
	if err := os.MkdirAll(foundry, 0o700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	got, err := ResolveKeystoreDir("", false)
	if err != nil {
		t.Fatalf("ResolveKeystoreDir: %v", err)
	}
	if got != foundry {
		t.Fatalf("dir = %q, want %q", got, foundry)
	}
}

// TestResolveKeystoreDirDefaultCreate asserts that the polycli fallback
// is created when createDefault is true.
func TestResolveKeystoreDirDefaultCreate(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("ETH_KEYSTORE", "")
	got, err := ResolveKeystoreDir("", true)
	if err != nil {
		t.Fatalf("ResolveKeystoreDir: %v", err)
	}
	want := filepath.Join(home, ".polycli", "keystores")
	if got != want {
		t.Fatalf("dir = %q, want %q", got, want)
	}
	if st, err := os.Stat(want); err != nil || !st.IsDir() {
		t.Fatalf("expected default dir created; err=%v", err)
	}
}

// TestResolveKeystoreDirDefaultNoCreate asserts that the polycli
// fallback path is returned but NOT created when createDefault is
// false. This is the signing path; it should not silently materialise
// empty keystores.
func TestResolveKeystoreDirDefaultNoCreate(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("ETH_KEYSTORE", "")
	got, err := ResolveKeystoreDir("", false)
	if err != nil {
		t.Fatalf("ResolveKeystoreDir: %v", err)
	}
	want := filepath.Join(home, ".polycli", "keystores")
	if got != want {
		t.Fatalf("dir = %q, want %q", got, want)
	}
	if _, err := os.Stat(want); !os.IsNotExist(err) {
		t.Fatalf("default dir should not be created when createDefault=false, err=%v", err)
	}
}
