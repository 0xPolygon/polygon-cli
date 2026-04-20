package wallet

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	accounts "github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"

	"github.com/0xPolygon/polygon-cli/gethkeystore"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
)

// ErrAccountExists is returned when an import would overwrite an
// existing key for the same address.
var ErrAccountExists = errors.New("account already exists in keystore")

// NewKeyStore returns a KeyStore rooted at dir using light scrypt
// parameters. LightScryptN/P gives the cast-compatible "fast enough on
// a laptop" encryption — matches the Foundry default.
func NewKeyStore(dir string) *keystore.KeyStore {
	return keystore.NewKeyStore(dir, keystore.LightScryptN, keystore.LightScryptP)
}

// FindAccount resolves a CLI identifier to a keystore account.
// Supported identifier forms:
//   - `0x`-prefixed address
//   - path to a keystore JSON file (UTC-- / .json)
//   - decimal index into ks.Accounts() (for `--account 0` style use)
//
// The index form is what the tx/mktx/send paths call "account by
// position"; the wallet package does not currently expose that surface,
// but supporting it here costs nothing.
func FindAccount(ks *keystore.KeyStore, identifier string) (accounts.Account, error) {
	identifier = strings.TrimSpace(identifier)
	if identifier == "" {
		return accounts.Account{}, &client.UsageError{Msg: "empty address or file path"}
	}
	// Integer index — operators often prefer `--account 0` to referring
	// to an address by heart. Only accept this when the string is a
	// pure unsigned integer that does not also look like an address.
	if n, err := strconv.ParseUint(identifier, 10, 32); err == nil {
		list := ks.Accounts()
		if int(n) >= len(list) {
			return accounts.Account{}, &client.UsageError{
				Msg: fmt.Sprintf("keystore has %d accounts; index %d out of range", len(list), n),
			}
		}
		return list[int(n)], nil
	}
	// File path — honour both exact paths and bare file names inside
	// the keystore directory.
	if strings.ContainsAny(identifier, "/\\") || strings.HasSuffix(identifier, ".json") || strings.HasPrefix(identifier, "UTC--") {
		addr, err := AddressFromKeystoreFile(identifier)
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

// AddressFromKeystoreFile reads a v3 JSON keystore from path and
// returns the address it encodes. Works for keystores that include the
// `address` field at the top level (go-ethereum + foundry do).
func AddressFromKeystoreFile(path string) (common.Address, error) {
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

// DecryptKeystoreAccount loads the raw JSON for acc and decrypts it
// with password, returning the raw ECDSA private key. It is the lower
// level of ks.Unlock; we need the key material directly for signing
// utilities that are not part of the keystore's own signing surface.
func DecryptKeystoreAccount(acc accounts.Account, password string) (*ecdsa.PrivateKey, error) {
	data, err := os.ReadFile(acc.URL.Path)
	if err != nil {
		return nil, fmt.Errorf("reading keystore file %s: %w", acc.URL.Path, err)
	}
	return gethkeystore.DecryptKeystoreFile(data, password)
}
