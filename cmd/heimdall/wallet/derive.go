package wallet

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"

	sharedwallet "github.com/0xPolygon/polygon-cli/internal/heimdall/wallet"
)

// DefaultDerivationPath re-exports the shared constant so existing
// callers under cmd/heimdall/wallet keep their current import surface.
const DefaultDerivationPath = sharedwallet.DefaultDerivationPath

// deriveFromMnemonic is a thin wrapper around
// internal/heimdall/wallet.DeriveFromMnemonic retained so existing
// in-package call sites do not need to be rewritten in the same commit
// that introduces the shared package.
func deriveFromMnemonic(mnemonic, passphrase, path string, index uint32) (*ecdsa.PrivateKey, string, common.Address, error) {
	return sharedwallet.DeriveFromMnemonic(mnemonic, passphrase, path, index)
}

// parseDerivationPath is a thin wrapper around
// internal/heimdall/wallet.ParseDerivationPath.
func parseDerivationPath(path string) ([]uint32, error) {
	return sharedwallet.ParseDerivationPath(path)
}
