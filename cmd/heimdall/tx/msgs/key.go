package msgs

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	accounts "github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"

	"github.com/0xPolygon/polygon-cli/gethkeystore"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
)

// defaultDerivationPath is the standard Ethereum BIP-44 path at
// index 0. Matches `cast wallet new-mnemonic` and the wallet package.
const defaultDerivationPath = "m/44'/60'/0'/0/0"

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
		priv, _, addr, err := deriveFromMnemonic(opts.Mnemonic, "", opts.DerivationPath, opts.MnemonicIndex)
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

	dir, err := resolveKeystoreDir(opts.KeystoreDir)
	if err != nil {
		return nil, err
	}
	ks := newLightKeyStore(dir)
	acc, err := findKeystoreAccount(ks, identifier)
	if err != nil {
		return nil, err
	}
	password, err := readPassword(opts, stdin)
	if err != nil {
		return nil, err
	}
	priv, err := decryptKeystoreAccount(acc, password)
	if err != nil {
		return nil, fmt.Errorf("decrypting keystore entry: %w", err)
	}
	return signerFromKey(priv), nil
}

// parsePrivateKeyHex decodes a 0x-prefixed or bare 32-byte hex string
// into an ECDSA private key. Duplicated with the wallet package on
// purpose: we don't import cmd/heimdall/wallet to avoid an import
// cycle with cmd/heimdall.
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

// --- BIP-39 / BIP-32 derivation (subset of cmd/heimdall/wallet/derive.go). ---

func deriveFromMnemonic(mnemonic, passphrase, path string, index uint32) (*ecdsa.PrivateKey, string, common.Address, error) {
	mnemonic = strings.TrimSpace(mnemonic)
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, "", common.Address{}, fmt.Errorf("invalid BIP-39 mnemonic")
	}
	finalPath := path
	if finalPath == "" {
		base := strings.TrimSuffix(defaultDerivationPath, "/0")
		finalPath = fmt.Sprintf("%s/%d", base, index)
	}
	seed := bip39.NewSeed(mnemonic, passphrase)
	parts, err := parseDerivationPath(finalPath)
	if err != nil {
		return nil, "", common.Address{}, err
	}
	master, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, "", common.Address{}, fmt.Errorf("deriving master key: %w", err)
	}
	current := master
	for i, idx := range parts {
		current, err = current.NewChildKey(idx)
		if err != nil {
			return nil, "", common.Address{}, fmt.Errorf("deriving child at position %d (%s): %w", i+1, finalPath, err)
		}
	}
	priv, err := ethcrypto.ToECDSA(current.Key)
	if err != nil {
		return nil, "", common.Address{}, fmt.Errorf("converting derived key: %w", err)
	}
	return priv, finalPath, ethcrypto.PubkeyToAddress(priv.PublicKey), nil
}

func parseDerivationPath(path string) ([]uint32, error) {
	if path == "" {
		return nil, fmt.Errorf("empty derivation path")
	}
	pieces := strings.Split(path, "/")
	if pieces[0] != "m" {
		return nil, fmt.Errorf("derivation path must start with \"m\", got %q", pieces[0])
	}
	out := make([]uint32, 0, len(pieces)-1)
	for _, p := range pieces[1:] {
		if p == "" {
			return nil, fmt.Errorf("empty segment in derivation path %q", path)
		}
		var base uint32
		if strings.HasSuffix(p, "'") {
			base = bip32.FirstHardenedChild
			p = strings.TrimSuffix(p, "'")
		}
		n, err := strconv.ParseUint(p, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid derivation path segment %q: %w", p, err)
		}
		if base == 0 && n >= uint64(bip32.FirstHardenedChild) {
			return nil, fmt.Errorf("non-hardened segment %s out of range (use %s' to harden)", p, p)
		}
		out = append(out, uint32(n)+base)
	}
	return out, nil
}

// --- Keystore helpers. ---

// resolveKeystoreDir returns the keystore directory per the same
// precedence rule as cmd/heimdall/wallet: flag > ETH_KEYSTORE >
// ~/.foundry/keystores (if exists) > ~/.polycli/keystores.
//
// Unlike the wallet package this implementation does NOT create the
// default directory on demand: `mktx`/`send`/`estimate` are signing
// operations, not keystore-management commands. If the default dir
// doesn't exist and no other source is configured, the caller should
// see a clear "keystore not found" error from findKeystoreAccount.
func resolveKeystoreDir(override string) (string, error) {
	switch {
	case override != "":
		abs, err := filepath.Abs(override)
		if err != nil {
			return "", fmt.Errorf("resolving --keystore-dir %q: %w", override, err)
		}
		return abs, nil
	case os.Getenv("ETH_KEYSTORE") != "":
		abs, err := filepath.Abs(os.Getenv("ETH_KEYSTORE"))
		if err != nil {
			return "", fmt.Errorf("resolving ETH_KEYSTORE %q: %w", os.Getenv("ETH_KEYSTORE"), err)
		}
		return abs, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolving home directory: %w", err)
	}
	foundry := filepath.Join(home, ".foundry", "keystores")
	if st, err := os.Stat(foundry); err == nil && st.IsDir() {
		return foundry, nil
	}
	polycli := filepath.Join(home, ".polycli", "keystores")
	return polycli, nil
}

// newLightKeyStore returns a KeyStore rooted at dir using light
// scrypt parameters — matches Foundry / cast defaults.
func newLightKeyStore(dir string) *keystore.KeyStore {
	return keystore.NewKeyStore(dir, keystore.LightScryptN, keystore.LightScryptP)
}

// findKeystoreAccount resolves a CLI identifier to a keystore
// account. Supports:
//   - 0x-prefixed address
//   - keystore file path (UTC-- / .json)
//   - integer index into ks.Accounts() for --account 0 style use
func findKeystoreAccount(ks *keystore.KeyStore, identifier string) (accounts.Account, error) {
	identifier = strings.TrimSpace(identifier)
	if identifier == "" {
		return accounts.Account{}, &client.UsageError{Msg: "empty address or file path"}
	}
	// Integer index — operators often prefer `--account 0` to referring
	// to an address by heart. Only accept this when the string is a
	// pure unsigned integer.
	if n, err := strconv.ParseUint(identifier, 10, 32); err == nil {
		list := ks.Accounts()
		if int(n) >= len(list) {
			return accounts.Account{}, &client.UsageError{
				Msg: fmt.Sprintf("keystore has %d accounts; index %d out of range", len(list), n),
			}
		}
		return list[int(n)], nil
	}
	if strings.ContainsAny(identifier, "/\\") || strings.HasSuffix(identifier, ".json") || strings.HasPrefix(identifier, "UTC--") {
		addr, err := addressFromKeystoreFile(identifier)
		if err != nil {
			return accounts.Account{}, err
		}
		identifier = addr.Hex()
	}
	if !common.IsHexAddress(identifier) {
		return accounts.Account{}, &client.UsageError{Msg: fmt.Sprintf("%q is neither an address, keystore index, nor file path", identifier)}
	}
	addr := common.HexToAddress(identifier)
	target := accounts.Account{Address: addr}
	got, err := ks.Find(target)
	if err != nil {
		return accounts.Account{}, fmt.Errorf("account %s not found in keystore: %w", addr.Hex(), err)
	}
	return got, nil
}

func addressFromKeystoreFile(path string) (common.Address, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return common.Address{}, fmt.Errorf("reading keystore file %s: %w", path, err)
	}
	var raw gethkeystore.RawKeystoreData
	if err := json.Unmarshal(data, &raw); err != nil {
		return common.Address{}, fmt.Errorf("parsing keystore %s: %w", path, err)
	}
	if raw.Address == "" {
		return common.Address{}, fmt.Errorf("keystore %s missing address field", path)
	}
	if !common.IsHexAddress("0x" + strings.TrimPrefix(raw.Address, "0x")) {
		return common.Address{}, fmt.Errorf("keystore %s has invalid address %q", path, raw.Address)
	}
	return common.HexToAddress(raw.Address), nil
}

func decryptKeystoreAccount(acc accounts.Account, password string) (*ecdsa.PrivateKey, error) {
	data, err := os.ReadFile(acc.URL.Path)
	if err != nil {
		return nil, fmt.Errorf("reading keystore file %s: %w", acc.URL.Path, err)
	}
	return gethkeystore.DecryptKeystoreFile(data, password)
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
