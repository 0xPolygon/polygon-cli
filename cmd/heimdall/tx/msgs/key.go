package msgs

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	sharedwallet "github.com/0xPolygon/polygon-cli/internal/heimdall/wallet"
)

// ResolvedSigner carries everything a msg subcommand needs to sign
// and address a Msg: the secp256k1 private key plus the derived
// 20-byte Ethereum-style address used as the signer identifier on
// Heimdall messages.
type ResolvedSigner struct {
	Key     *ecdsa.PrivateKey
	Address common.Address
}

// ResolveSigningKey returns the signing key for the current TxOpts.
// Precedence (highest first):
//  1. --private-key (hex). Logs a warning because the key is visible
//     in shell history / `ps`.
//  2. --mnemonic plus --mnemonic-index / --derivation-path.
//  3. --keystore-file (explicit JSON path).
//  4. --account / --from against the resolved keystore directory.
//
// At least one source must be provided; otherwise a UsageError is
// returned so the command exits with rc=3 per §2.1.
//
// BIP-39/32 derivation and keystore access are delegated to
// internal/heimdall/wallet so this path stays consistent with
// `polycli heimdall wallet`. Keystore-dir resolution is called with
// createDefault=false because signing is not a keystore-management
// operation — if the default dir does not exist, the caller should see
// a clear "account not found" error rather than have polycli silently
// materialise an empty keystore dir.
func ResolveSigningKey(opts *TxOpts, stdin io.Reader) (*ResolvedSigner, error) {
	switch {
	case opts.PrivateKey != "":
		log.Warn().Msg("using --private-key exposes key material via shell history; prefer a keystore for anything beyond local dev")
		priv, err := parsePrivateKeyHex(opts.PrivateKey)
		if err != nil {
			return nil, err
		}
		return signerFromKey(priv), nil
	case opts.Mnemonic != "":
		priv, _, addr, err := sharedwallet.DeriveFromMnemonic(opts.Mnemonic, "", opts.DerivationPath, opts.MnemonicIndex)
		if err != nil {
			return nil, err
		}
		return &ResolvedSigner{Key: priv, Address: addr}, nil
	}

	// Keystore path: --keystore-file wins over --account / --from so
	// operators can point at a specific file even with an ambient
	// keystore directory. `--account` then `--from` — both accept an
	// address or (for --account) a keystore index.
	identifier := opts.KeystoreFile
	if identifier == "" {
		identifier = opts.Account
	}
	if identifier == "" {
		identifier = opts.From
	}
	if identifier == "" {
		return nil, &client.UsageError{Msg: "one of --private-key, --mnemonic, --keystore-file, --account, or --from is required"}
	}

	dir, err := sharedwallet.ResolveKeystoreDir(opts.KeystoreDir, false)
	if err != nil {
		return nil, err
	}
	ks := sharedwallet.NewKeyStore(dir)
	acc, err := sharedwallet.FindAccount(ks, identifier)
	if err != nil {
		return nil, err
	}
	password, err := readPassword(opts, stdin)
	if err != nil {
		return nil, err
	}
	priv, err := sharedwallet.DecryptKeystoreAccount(acc, password)
	if err != nil {
		return nil, fmt.Errorf("decrypting keystore entry: %w", err)
	}
	return signerFromKey(priv), nil
}

// parsePrivateKeyHex decodes a 0x-prefixed or bare 32-byte hex string
// into an ECDSA private key. This is the only key-input helper we keep
// local to msgs/: it is not in the shared wallet package because the
// wallet package has its own flag-driven variant with different error
// wording.
func parsePrivateKeyHex(input string) (*ecdsa.PrivateKey, error) {
	s := strings.TrimSpace(input)
	s = strings.TrimPrefix(strings.TrimPrefix(s, "0x"), "0X")
	if len(s) != 64 {
		return nil, fmt.Errorf("private key must be 32 bytes (64 hex chars), got %d", len(s))
	}
	raw, err := hex.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("decoding private key: %w", err)
	}
	return ethcrypto.ToECDSA(raw)
}

// signerFromKey derives the Ethereum address for priv and returns a
// populated ResolvedSigner.
func signerFromKey(priv *ecdsa.PrivateKey) *ResolvedSigner {
	return &ResolvedSigner{Key: priv, Address: ethcrypto.PubkeyToAddress(priv.PublicKey)}
}

// readPassword returns the keystore password per --password /
// --password-file. Interactive prompt falls back to stdin without a
// terminal cue because we do not want to depend on tty detection in
// the tx path; operators running `send` from scripts should always
// provide --password-file.
func readPassword(opts *TxOpts, stdin io.Reader) (string, error) {
	if opts.Password != "" && opts.PasswordFile != "" {
		return "", &client.UsageError{Msg: "--password and --password-file are mutually exclusive"}
	}
	if opts.Password != "" {
		return opts.Password, nil
	}
	if opts.PasswordFile != "" {
		raw, err := os.ReadFile(opts.PasswordFile)
		if err != nil {
			return "", fmt.Errorf("reading password file %s: %w", opts.PasswordFile, err)
		}
		return trimTrailingNewline(string(raw)), nil
	}
	if stdin == nil {
		return "", &client.UsageError{Msg: "no password source (provide --password or --password-file)"}
	}
	// Read a single line from stdin. We don't attempt tty-suppressing
	// echo: the wallet package owns the interactive UX; tx path is
	// primarily scripted. Operators who want interactive signing can
	// run `polycli heimdall wallet ...` first.
	buf := make([]byte, 0, 256)
	tmp := make([]byte, 1)
	for {
		n, err := stdin.Read(tmp)
		if n > 0 {
			if tmp[0] == '\n' {
				break
			}
			buf = append(buf, tmp[0])
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", fmt.Errorf("reading password: %w", err)
		}
	}
	return trimTrailingNewline(string(buf)), nil
}

func trimTrailingNewline(s string) string {
	if n := len(s); n > 0 && s[n-1] == '\n' {
		if n >= 2 && s[n-2] == '\r' {
			return s[:n-2]
		}
		return s[:n-1]
	}
	return strings.TrimRight(s, "\r")
}
