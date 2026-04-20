package wallet

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// personalSignPrefix is the EIP-191 "\x19Ethereum Signed Message:\n"
// prefix used by personal_sign / eth_sign clients. Followed by the
// ASCII-decimal length of the message and the message itself.
const personalSignPrefix = "\x19Ethereum Signed Message:\n"

// personalSignHash returns keccak256 of the EIP-191 prefixed message.
// Matches the hash that MetaMask, ethers.js, viem, cast, and all other
// mainstream Ethereum tooling use for message signing.
func personalSignHash(message []byte) []byte {
	prefix := []byte(fmt.Sprintf("%s%d", personalSignPrefix, len(message)))
	payload := append(prefix, message...)
	return crypto.Keccak256(payload)
}

// signPersonal signs message with priv using the EIP-191 personal_sign
// scheme. The returned 65-byte signature uses the canonical v = 27/28
// convention (not the pre-EIP-155 v = 0/1 emitted by
// crypto.Sign).
func signPersonal(priv *ecdsa.PrivateKey, message []byte) ([]byte, error) {
	sig, err := crypto.Sign(personalSignHash(message), priv)
	if err != nil {
		return nil, fmt.Errorf("signing message: %w", err)
	}
	// Raise v from {0,1} to {27,28} so the signature matches what
	// eth_sign / cast / ethers would emit.
	sig[64] += 27
	return sig, nil
}

// signRawHash signs a 32-byte hash directly without EIP-191
// framing. Emits v in {27,28} for consistency with signPersonal.
func signRawHash(priv *ecdsa.PrivateKey, hash []byte) ([]byte, error) {
	if len(hash) != 32 {
		return nil, fmt.Errorf("hash must be 32 bytes, got %d", len(hash))
	}
	sig, err := crypto.Sign(hash, priv)
	if err != nil {
		return nil, fmt.Errorf("signing hash: %w", err)
	}
	sig[64] += 27
	return sig, nil
}

// verifyPersonal verifies that sig (0x-hex) was produced over message
// with the private key of expected under EIP-191 personal_sign
// framing. Accepts signatures with v in {0,1} or {27,28}.
func verifyPersonal(expected common.Address, message []byte, sig []byte) (bool, error) {
	return verifySignature(expected, personalSignHash(message), sig)
}

// verifyRaw verifies sig against hash without EIP-191 framing.
func verifyRaw(expected common.Address, hash, sig []byte) (bool, error) {
	if len(hash) != 32 {
		return false, fmt.Errorf("hash must be 32 bytes, got %d", len(hash))
	}
	return verifySignature(expected, hash, sig)
}

func verifySignature(expected common.Address, digest, sig []byte) (bool, error) {
	if len(sig) != 65 {
		return false, fmt.Errorf("signature must be 65 bytes, got %d", len(sig))
	}
	// Normalise v: crypto.SigToPub wants {0,1}. Accept {27,28} and
	// subtract; reject anything else so we do not silently accept
	// malformed signatures.
	normalised := make([]byte, 65)
	copy(normalised, sig)
	switch normalised[64] {
	case 0, 1:
		// already normalised
	case 27, 28:
		normalised[64] -= 27
	default:
		return false, fmt.Errorf("signature v byte must be 0/1 or 27/28, got %d", sig[64])
	}
	pub, err := crypto.SigToPub(digest, normalised)
	if err != nil {
		return false, fmt.Errorf("recovering public key: %w", err)
	}
	got := crypto.PubkeyToAddress(*pub)
	return got == expected, nil
}

// parseSignatureHex decodes a 0x-prefixed or bare hex signature into
// the 65-byte binary form.
func parseSignatureHex(input string) ([]byte, error) {
	s := strings.TrimSpace(input)
	s = strings.TrimPrefix(strings.TrimPrefix(s, "0x"), "0X")
	if len(s) != 130 {
		return nil, fmt.Errorf("signature must be 65 bytes (130 hex chars), got %d", len(s))
	}
	raw, err := hex.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("decoding signature: %w", err)
	}
	return raw, nil
}
