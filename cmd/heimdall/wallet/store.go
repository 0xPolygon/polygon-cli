package wallet

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"os"
	"strings"

	accounts "github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"

	"github.com/0xPolygon/polygon-cli/gethkeystore"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
)

// newKeyStore returns a KeyStore rooted at dir using light scrypt
// parameters. LightScryptN/P gives the cast-compatible "fast enough
// on a laptop" encryption — matches the Foundry default.
func newKeyStore(dir string) *keystore.KeyStore {
	return keystore.NewKeyStore(dir, keystore.LightScryptN, keystore.LightScryptP)
}

// findAccount resolves a CLI identifier to a keystore account. The
// identifier can be an `0x`-prefixed address or a path to a keystore
// file. A path that names a file under the keystore directory is
// converted to the address stored inside that file.
func findAccount(ks *keystore.KeyStore, identifier string) (accounts.Account, error) {
	identifier = strings.TrimSpace(identifier)
	if identifier == "" {
		return accounts.Account{}, &client.UsageError{Msg: "empty address or file path"}
	}
	// File path — honour both exact paths and bare file names inside
	// the keystore directory. We delegate address extraction to
	// go-ethereum by reading the JSON body.
	if strings.ContainsAny(identifier, "/\\") || strings.HasSuffix(identifier, ".json") || strings.HasPrefix(identifier, "UTC--") {
		addr, err := addressFromKeystoreFile(identifier)
		if err != nil {
			return accounts.Account{}, err
		}
		identifier = addr.Hex()
	}
	if !common.IsHexAddress(identifier) {
		return accounts.Account{}, &client.UsageError{Msg: fmt.Sprintf("%q is neither an address nor a keystore file path", identifier)}
	}
	addr := common.HexToAddress(identifier)
	target := accounts.Account{Address: addr}
	got, err := ks.Find(target)
	if err != nil {
		return accounts.Account{}, fmt.Errorf("account %s not found in keystore: %w", addr.Hex(), err)
	}
	return got, nil
}

// addressFromKeystoreFile reads a v3 JSON keystore from path and
// returns the address it encodes. Works for both keystores that
// include the `address` field at the top level (go-ethereum + foundry
// do) and for ones that only have `crypto`.
func addressFromKeystoreFile(path string) (common.Address, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return common.Address{}, fmt.Errorf("reading keystore file %s: %w", path, err)
	}
	// RawKeystoreData has an explicit Address field; re-use it rather
	// than hand-unmarshal here.
	var raw gethkeystore.RawKeystoreData
	if err := unmarshalJSON(data, &raw); err != nil {
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

// decryptKeystoreAccount loads the raw JSON for acc and decrypts it
// with password, returning the raw ECDSA private key. It is the lower
// level of ks.Unlock; we need the key material directly for signing
// utilities that are not part of the keystore's own signing surface.
func decryptKeystoreAccount(acc accounts.Account, password string) (*ecdsa.PrivateKey, error) {
	data, err := os.ReadFile(acc.URL.Path)
	if err != nil {
		return nil, fmt.Errorf("reading keystore file %s: %w", acc.URL.Path, err)
	}
	return gethkeystore.DecryptKeystoreFile(data, password)
}

// ErrAccountExists is returned when an import would overwrite an
// existing key for the same address.
var ErrAccountExists = errors.New("account already exists in keystore")
