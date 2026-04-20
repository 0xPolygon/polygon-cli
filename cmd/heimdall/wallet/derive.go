package wallet

import (
	"crypto/ecdsa"
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

// DefaultDerivationPath is the standard Ethereum BIP-44 path at index
// 0. Matches `cast wallet new-mnemonic` and most hardware wallet
// defaults.
const DefaultDerivationPath = "m/44'/60'/0'/0/0"

// deriveFromMnemonic returns the ECDSA private key, derivation path,
// and Ethereum address for mnemonic at the given path / index.
// If path is empty it is built from DefaultDerivationPath with the
// final component replaced by index.
func deriveFromMnemonic(mnemonic, passphrase, path string, index uint32) (*ecdsa.PrivateKey, string, common.Address, error) {
	mnemonic = strings.TrimSpace(mnemonic)
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, "", common.Address{}, fmt.Errorf("invalid BIP-39 mnemonic")
	}
	finalPath := path
	if finalPath == "" {
		// Strip the trailing index and re-append the requested one.
		base := strings.TrimSuffix(DefaultDerivationPath, "/0")
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
	priv, err := crypto.ToECDSA(current.Key)
	if err != nil {
		return nil, "", common.Address{}, fmt.Errorf("converting derived key: %w", err)
	}
	return priv, finalPath, crypto.PubkeyToAddress(priv.PublicKey), nil
}

// parseDerivationPath turns a path like "m/44'/60'/0'/0/0" into the
// list of BIP-32 child indices. Hardened components are marked with a
// trailing apostrophe (') and offset by bip32.FirstHardenedChild.
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
		// uint32 overflow guard: ParseUint already restricts to 32
		// bits, and adding FirstHardenedChild (2^31) to a value < 2^31
		// stays within uint32. A non-hardened segment >= 2^31 would
		// conflict with the hardened half of the tree and should be
		// expressed with the apostrophe instead.
		if base == 0 && n >= uint64(bip32.FirstHardenedChild) {
			return nil, fmt.Errorf("non-hardened segment %s out of range (use %s' to harden)", p, p)
		}
		out = append(out, uint32(n)+base)
	}
	return out, nil
}
