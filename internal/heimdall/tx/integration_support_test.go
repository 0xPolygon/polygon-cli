//go:build heimdall_integration

package tx

import (
	"crypto/ecdsa"
	"encoding/hex"
	"testing"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

// fixedECDSAKeyForIntegration duplicates fixedECDSAKey from
// builder_test.go (which is not compiled under the integration build
// tag) so integration tests can get a deterministic signing key
// without pulling the unit test helpers into the integration build.
func fixedECDSAKeyForIntegration(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	b, err := hex.DecodeString("0101010101010101010101010101010101010101010101010101010101010101")
	if err != nil {
		t.Fatalf("decoding fixed key: %v", err)
	}
	priv, err := ethcrypto.ToECDSA(b)
	if err != nil {
		t.Fatalf("loading fixed key: %v", err)
	}
	return priv
}
