package wallet

import (
	"crypto/ecdsa"

	accounts "github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"

	sharedwallet "github.com/0xPolygon/polygon-cli/internal/heimdall/wallet"
)

// ErrAccountExists is re-exported from the shared wallet package so
// existing callers compile unchanged.
var ErrAccountExists = sharedwallet.ErrAccountExists

// newKeyStore is a thin wrapper around
// internal/heimdall/wallet.NewKeyStore.
func newKeyStore(dir string) *keystore.KeyStore {
	return sharedwallet.NewKeyStore(dir)
}

// findAccount is a thin wrapper around
// internal/heimdall/wallet.FindAccount.
func findAccount(ks *keystore.KeyStore, identifier string) (accounts.Account, error) {
	return sharedwallet.FindAccount(ks, identifier)
}

// addressFromKeystoreFile is a thin wrapper around
// internal/heimdall/wallet.AddressFromKeystoreFile.
func addressFromKeystoreFile(path string) (common.Address, error) {
	return sharedwallet.AddressFromKeystoreFile(path)
}

// decryptKeystoreAccount is a thin wrapper around
// internal/heimdall/wallet.DecryptKeystoreAccount.
func decryptKeystoreAccount(acc accounts.Account, password string) (*ecdsa.PrivateKey, error) {
	return sharedwallet.DecryptKeystoreAccount(acc, password)
}
