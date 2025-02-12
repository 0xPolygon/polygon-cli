package nodekey

import (
	"bytes"
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"math/big"

	secp256k1 "github.com/btcsuite/btcd/btcec/v2"
	"github.com/cometbft/cometbft/crypto"
	"github.com/cometbft/cometbft/privval"

	cmtjson "github.com/cometbft/cometbft/libs/json"
	gethcrypto "github.com/ethereum/go-ethereum/crypto"
)

const (
	PrivKeyName = "comet/PrivKeySecp256k1Uncompressed"
	PubKeyName  = "comet/PubKeySecp256k1Uncompressed"

	KeyType     = "secp256k1"
	PrivKeySize = 32
	// PubKeySize (uncompressed) is composed of 65 bytes for two field elements (x and y)
	// and a prefix byte (0x04) to indicate that it is uncompressed.
	PubKeySize = 65
	// SigSize is the size of the ECDSA signature.
	SigSize = 65
)

var _ crypto.PrivKey = PrivKey{}
var _ crypto.PubKey = PubKey{}

// -------------------------------------
// PrivKey type
// -------------------------------------

type PrivKey []byte

// Bytes marshals the private key using amino encoding.
func (privKey PrivKey) Bytes() []byte {
	return []byte(privKey)
}

// PubKey performs the point-scalar multiplication from the privKey on the
// generator point to get the pubkey.
func (privKey PrivKey) PubKey() crypto.PubKey {
	privateObject, err := gethcrypto.ToECDSA(privKey)
	if err != nil {
		panic(err)
	}

	pk := gethcrypto.FromECDSAPub(&privateObject.PublicKey)

	return PubKey(pk)
}

// Equals - you probably don't need to use this.
// Runs in constant time based on length of the keys.
func (privKey PrivKey) Equals(other crypto.PrivKey) bool {
	if otherSecp, ok := other.(PrivKey); ok {
		return subtle.ConstantTimeCompare(privKey[:], otherSecp[:]) == 1
	}
	return false
}

func (privKey PrivKey) Type() string {
	return KeyType
}

// Sign creates an ECDSA signature on curve Secp256k1, using SHA256 on the msg.
// The returned signature will be of the form R || S || V (in lower-S form).
func (privKey PrivKey) Sign(msg []byte) ([]byte, error) {
	privateObject, err := gethcrypto.ToECDSA(privKey)
	if err != nil {
		return nil, err
	}

	return gethcrypto.Sign(gethcrypto.Keccak256(msg), privateObject)
}

// -------------------------------------
// PubKey type
// -------------------------------------

type PubKey []byte

// Bytes returns the pubkey marshaled with amino encoding.
func (pubKey PubKey) Bytes() []byte {
	return []byte(pubKey)
}

// Address returns a Ethereym style addresses: Last_20_Bytes(KECCAK256(pubkey))
func (pubKey PubKey) Address() crypto.Address {
	if len(pubKey) != PubKeySize {
		panic(fmt.Sprintf("length of pubkey is incorrect %d != %d", len(pubKey), PubKeySize))
	}
	return gethcrypto.Keccak256(pubKey[1:])[12:]
}

func (pubKey PubKey) Equals(other crypto.PubKey) bool {
	if otherSecp, ok := other.(PubKey); ok {
		return bytes.Equal(pubKey[:], otherSecp[:])
	}
	return false
}

func (pubKey PubKey) Type() string {
	return KeyType
}

// VerifySignature verifies a signature of the form R || S || V.
// It rejects signatures which are not in lower-S form.
func (pubKey PubKey) VerifySignature(msg []byte, sigStr []byte) bool {
	if len(sigStr) != SigSize {

		return false
	}

	hash := gethcrypto.Keccak256(msg)
	return gethcrypto.VerifySignature(pubKey, hash, sigStr[:64])
}

func init() {
	cmtjson.RegisterType(PubKey{}, PubKeyName)
	cmtjson.RegisterType(PrivKey{}, PrivKeyName)
}

// Generate an secp256k1 private key from a secret.
// Most of the logic has been copy/pasted from 0xPolygon/cometbft's fork.
// https://github.com/0xPolygon/cometbft/blob/v0.1.2-beta-polygon/crypto/secp256k1/secp256k1.go
// Notes:
// - It is not possible to import the package yet because go.mod declares its path as github.com/cometbft/cometbft instead of github.com/0xpolygon/cometbft.
// - This logic will need to be updated to support newer versions.
func generateSecp256k1PrivateKey(secret []byte) PrivKey {
	// To guarantee that we have a valid field element, we use the approach of: "Suite B Implementer’s Guide to FIPS 186-3", A.2.1
	// https://apps.nsa.gov/iaarchive/library/ia-guidance/ia-solutions-for-classified/algorithm-guidance/suite-b-implementers-guide-to-fips-186-3-ecdsa.cfm
	// See also https://github.com/golang/go/blob/0380c9ad38843d523d9c9804fe300cb7edd7cd3c/src/crypto/ecdsa/ecdsa.go#L89-L101
	secretHash := sha256.Sum256(secret)
	fe := new(big.Int).SetBytes(secretHash[:])

	one := new(big.Int).SetInt64(1)
	n := new(big.Int).Sub(secp256k1.S256().N, one)
	fe.Mod(fe, n)
	fe.Add(fe, one)

	feB := fe.Bytes()
	privKey32 := make([]byte, PrivKeySize)
	// Copy feB over to fixed 32 byte privKey32 and pad (if necessary).
	copy(privKey32[32-len(feB):32], feB)

	return PrivKey(privKey32)
}

func displayHeimdallV2PrivValidatorKey(privKey crypto.PrivKey) error {
	nodeKey := privval.FilePVKey{
		Address: privKey.PubKey().Address(),
		PubKey:  privKey.PubKey(),
		PrivKey: privKey,
	}
	jsonBytes, err := cmtjson.MarshalIndent(nodeKey, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(jsonBytes))
	return nil
}
