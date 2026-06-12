// Package tx implements the shared Heimdall transaction builder used
// by `polycli heimdall mktx`, `send`, and `estimate`. The builder
// assembles a cosmos.tx.v1beta1.TxRaw from one or more messages,
// fetches the signer's account number + sequence (unless overridden),
// signs with a secp256k1-eth key, and — via the broadcast helper —
// submits it through the REST gateway.
//
// The signing scheme is the Heimdall v2 `PubKeySecp256k1eth` variant:
// curve secp256k1, but the digest is keccak256 (not sha256) and the
// 65-byte (r||s||v) output from eth signing is truncated to 64 bytes
// (r||s) before landing in TxRaw.signatures. This matches heimdall-v2
// /crypto/keys/secp256k1.PrivKey.Sign.
package tx

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/proto"
)

// SignMode is the signing mode used to compute the signature pre-image.
type SignMode int32

const (
	// SignModeDirect signs the SIGN_MODE_DIRECT doc (proto-serialized
	// SignDoc). This is the default and matches cosmos-sdk's default.
	SignModeDirect SignMode = SignMode(proto.SignModeDirect)
	// SignModeAminoJSON signs the legacy amino-JSON document. Kept for
	// compatibility with older tooling; not otherwise recommended.
	SignModeAminoJSON SignMode = SignMode(proto.SignModeAminoJSON)
)

// String returns the cobra-flag-friendly name of the sign mode.
func (m SignMode) String() string {
	switch m {
	case SignModeDirect:
		return "direct"
	case SignModeAminoJSON:
		return "amino-json"
	default:
		return "unspecified"
	}
}

// ParseSignMode maps the flag-style name ("direct" / "amino-json") to a
// SignMode. Unknown names return an error so usage mistakes surface at
// the boundary rather than silently falling back.
func ParseSignMode(name string) (SignMode, error) {
	switch name {
	case "", "direct", "DIRECT":
		return SignModeDirect, nil
	case "amino-json", "amino_json", "AMINO_JSON":
		return SignModeAminoJSON, nil
	default:
		return 0, fmt.Errorf("unknown sign mode %q: expected direct or amino-json", name)
	}
}

// uncompressedPubkey derives the 65-byte uncompressed secp256k1 pubkey
// (0x04 || X || Y) from an ECDSA private key. Heimdall's
// PubKeySecp256k1eth stores the uncompressed form and ethcrypto.
// FromECDSAPub writes exactly that.
func uncompressedPubkey(priv *ecdsa.PrivateKey) []byte {
	return ethcrypto.FromECDSAPub(&priv.PublicKey)
}

// EthAddress derives the 20-byte Ethereum-style address that Heimdall
// uses for the `proposer` / `from` fields on messages. Matches
// PubKey.Address() in heimdall-v2's secp256k1 package (keccak256 of
// the uncompressed pubkey minus the 0x04 prefix, right-most 20 bytes).
func EthAddress(priv *ecdsa.PrivateKey) common.Address {
	return ethcrypto.PubkeyToAddress(priv.PublicKey)
}

// signDigest signs digest (a 32-byte hash) with priv and returns the
// 64-byte r||s payload that Cosmos SDK stores in TxRaw.signatures.
// go-ethereum's Sign returns 65 bytes (r||s||v); Cosmos drops v.
func signDigest(priv *ecdsa.PrivateKey, digest []byte) ([]byte, error) {
	if len(digest) != 32 {
		return nil, fmt.Errorf("signDigest: expected 32-byte digest, got %d", len(digest))
	}
	sig, err := ethcrypto.Sign(digest, priv)
	if err != nil {
		return nil, fmt.Errorf("signing digest: %w", err)
	}
	if len(sig) != 65 {
		return nil, fmt.Errorf("signDigest: expected 65-byte signature, got %d", len(sig))
	}
	return sig[:64], nil
}

// signBytesDirect returns the canonical SIGN_MODE_DIRECT pre-image and
// its keccak256 digest. Callers sign the digest.
func signBytesDirect(bodyBytes, authInfoBytes []byte, chainID string, accountNumber uint64) ([]byte, []byte) {
	doc := &proto.SignDoc{
		BodyBytes:     bodyBytes,
		AuthInfoBytes: authInfoBytes,
		ChainID:       chainID,
		AccountNumber: accountNumber,
	}
	raw := doc.Marshal()
	digest := ethcrypto.Keccak256(raw)
	return raw, digest
}

// signBytesAminoJSON returns the legacy amino-JSON pre-image and its
// keccak256 digest. Amino-JSON signing serializes a StdSignDoc JSON
// document with sorted keys:
//
//	{
//	  "account_number": "<n>",
//	  "chain_id": "<id>",
//	  "fee": { "amount": [...], "gas": "<n>" },
//	  "memo": "<...>",
//	  "msgs": [ { "type": "<amino.name>", "value": { ... } }, ... ],
//	  "sequence": "<n>"
//	}
//
// The `type` field on each msg is the `amino.name` option declared on
// the proto; the `value` is the msg fields in natural JSON form
// (numbers as strings for uint64/int). Only a minimal subset of Msg
// types is supported here — enough for MsgWithdrawFeeTx. Unknown msg
// types return an error so operators aren't silently signing wrong
// bytes.
//
// The document is marshalled with sorted keys at every level, which is
// the cosmos-sdk canonical form. We use encoding/json + manual sort.
func signBytesAminoJSON(b *Builder, signerAccountNumber uint64) ([]byte, []byte, error) {
	if len(b.msgs) == 0 {
		return nil, nil, fmt.Errorf("amino-json sign: no messages")
	}
	msgs := make([]map[string]any, 0, len(b.msgs))
	for _, m := range b.msgs {
		js, err := m.AminoJSON()
		if err != nil {
			return nil, nil, fmt.Errorf("amino-json sign: %w", err)
		}
		msgs = append(msgs, map[string]any{
			"type":  m.AminoName(),
			"value": js,
		})
	}
	feeAmounts := make([]map[string]any, 0, len(b.fee.Amount))
	for _, c := range b.fee.Amount {
		feeAmounts = append(feeAmounts, map[string]any{
			"amount": c.Amount,
			"denom":  c.Denom,
		})
	}
	doc := map[string]any{
		"account_number": fmt.Sprintf("%d", signerAccountNumber),
		"chain_id":       b.chainID,
		"fee": map[string]any{
			"amount": feeAmounts,
			"gas":    fmt.Sprintf("%d", b.fee.GasLimit),
		},
		"memo":     b.memo,
		"msgs":     msgs,
		"sequence": fmt.Sprintf("%d", b.sequence),
	}
	raw, err := marshalCanonicalJSON(doc)
	if err != nil {
		return nil, nil, fmt.Errorf("amino-json sign: %w", err)
	}
	digest := ethcrypto.Keccak256(raw)
	return raw, digest, nil
}

// marshalCanonicalJSON is a sorted-keys JSON encoder used for
// amino-JSON sign bytes. encoding/json already sorts map[string]…
// keys but does not handle nested slices of maps deterministically
// beyond that; we walk the tree and re-marshal.
func marshalCanonicalJSON(v any) ([]byte, error) {
	canon := canonicalize(v)
	return json.Marshal(canon)
}

// canonicalize sorts every map's keys at every level so encoding/json
// produces byte-identical output across machines. Slices are preserved
// in their input order (cosmos-sdk amino-JSON does not sort slice
// elements).
func canonicalize(v any) any {
	switch x := v.(type) {
	case map[string]any:
		keys := make([]string, 0, len(x))
		for k := range x {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		out := make(orderedMap, 0, len(keys))
		for _, k := range keys {
			out = append(out, orderedEntry{K: k, V: canonicalize(x[k])})
		}
		return out
	case []any:
		out := make([]any, len(x))
		for i, e := range x {
			out[i] = canonicalize(e)
		}
		return out
	case []map[string]any:
		out := make([]any, len(x))
		for i, e := range x {
			out[i] = canonicalize(e)
		}
		return out
	default:
		return v
	}
}

// orderedMap is a key-sorted map serialized as a JSON object. Used to
// bypass encoding/json's default map ordering (which is Go-random for
// map[string]any) while preserving the insertion (sorted) order.
type orderedMap []orderedEntry

type orderedEntry struct {
	K string
	V any
}

// MarshalJSON produces a deterministic `{"k":v,...}` in insertion order.
func (m orderedMap) MarshalJSON() ([]byte, error) {
	buf := []byte{'{'}
	for i, e := range m {
		if i > 0 {
			buf = append(buf, ',')
		}
		kb, err := json.Marshal(e.K)
		if err != nil {
			return nil, err
		}
		buf = append(buf, kb...)
		buf = append(buf, ':')
		vb, err := json.Marshal(e.V)
		if err != nil {
			return nil, err
		}
		buf = append(buf, vb...)
	}
	buf = append(buf, '}')
	return buf, nil
}
