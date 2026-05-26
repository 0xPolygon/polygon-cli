package proto

import (
	"fmt"

	"google.golang.org/protobuf/encoding/protowire"
)

// Any mirrors google.protobuf.Any — a polymorphic envelope carrying a
// type_url and a proto-encoded value. We use this for both tx messages
// and pubkeys.
type Any struct {
	TypeURL string
	Value   []byte
}

// Marshal encodes the Any to proto3 wire format. Type URLs and values
// are both required for a non-empty Any; zero-Any marshals to an empty
// byte slice.
func (a *Any) Marshal() []byte {
	if a == nil {
		return nil
	}
	var out []byte
	out = appendString(out, 1, a.TypeURL)
	out = appendBytes(out, 2, a.Value)
	return out
}

// UnmarshalAny parses a length-prefixed Any from b and returns it.
// Unknown fields are skipped.
func UnmarshalAny(b []byte) (*Any, error) {
	out := &Any{}
	for len(b) > 0 {
		num, _, val, n, err := consumeField(b)
		if err != nil {
			return nil, err
		}
		switch num {
		case 1:
			out.TypeURL = string(val)
		case 2:
			out.Value = append([]byte(nil), val...)
		}
		b = b[n:]
	}
	return out, nil
}

// Coin mirrors cosmos.base.v1beta1.Coin.
type Coin struct {
	Denom  string
	Amount string
}

// Marshal encodes the Coin.
func (c Coin) Marshal() []byte {
	var out []byte
	out = appendString(out, 1, c.Denom)
	out = appendString(out, 2, c.Amount)
	return out
}

// PubKeyTypeURL is the Any type URL for a Heimdall / Cosmos SDK
// secp256k1 public key. heimdall-v2 registers this concrete type with
// the Ethereum-style 65-byte uncompressed key; the proto shape is
// identical to the upstream cosmos-sdk type.
const PubKeyTypeURL = "/cosmos.crypto.secp256k1.PubKey"

// PubKeyAny returns an Any wrapping a cosmos.crypto.secp256k1.PubKey
// whose single `bytes key = 1` field is keyBytes. keyBytes should be
// the 65-byte uncompressed secp256k1 pubkey (0x04 || X || Y) for
// PubKeySecp256k1eth.
func PubKeyAny(keyBytes []byte) *Any {
	var val []byte
	val = appendBytes(val, 1, keyBytes)
	return &Any{TypeURL: PubKeyTypeURL, Value: val}
}

// SignModeDirect and SignModeAminoJSON are the two values of
// cosmos.tx.signing.v1beta1.SignMode that polycli supports.
const (
	SignModeDirect     int32 = 1
	SignModeAminoJSON  int32 = 127
	SignModeUnspecif   int32 = 0
)

// ModeInfoSingle is the ModeInfo.Single sub-message.
type ModeInfoSingle struct {
	Mode int32
}

// Marshal encodes ModeInfo.Single.
func (m ModeInfoSingle) Marshal() []byte {
	return appendInt32(nil, 1, m.Mode)
}

// ModeInfo represents the ModeInfo oneof; only Single is supported.
type ModeInfo struct {
	Single *ModeInfoSingle
}

// Marshal encodes ModeInfo.
func (m *ModeInfo) Marshal() []byte {
	if m == nil {
		return nil
	}
	var out []byte
	if m.Single != nil {
		s := m.Single
		out = appendSubmessage(out, 1, func() []byte { return s.Marshal() })
	}
	return out
}

// SignerInfo mirrors cosmos.tx.v1beta1.SignerInfo.
type SignerInfo struct {
	PublicKey *Any
	ModeInfo  *ModeInfo
	Sequence  uint64
}

// Marshal encodes SignerInfo.
func (s *SignerInfo) Marshal() []byte {
	if s == nil {
		return nil
	}
	var out []byte
	if s.PublicKey != nil {
		pk := s.PublicKey
		out = appendSubmessage(out, 1, func() []byte { return pk.Marshal() })
	}
	if s.ModeInfo != nil {
		mi := s.ModeInfo
		out = appendSubmessage(out, 2, func() []byte { return mi.Marshal() })
	}
	out = appendUint64(out, 3, s.Sequence)
	return out
}

// Fee mirrors cosmos.tx.v1beta1.Fee.
type Fee struct {
	Amount   []Coin
	GasLimit uint64
	Payer    string
	Granter  string
}

// Marshal encodes Fee.
func (f *Fee) Marshal() []byte {
	if f == nil {
		return nil
	}
	var out []byte
	for i := range f.Amount {
		c := f.Amount[i]
		out = appendSubmessage(out, 1, func() []byte { return c.Marshal() })
	}
	out = appendUint64(out, 2, f.GasLimit)
	out = appendString(out, 3, f.Payer)
	out = appendString(out, 4, f.Granter)
	return out
}

// AuthInfo mirrors cosmos.tx.v1beta1.AuthInfo.
type AuthInfo struct {
	SignerInfos []*SignerInfo
	Fee         *Fee
}

// Marshal encodes AuthInfo.
func (a *AuthInfo) Marshal() []byte {
	if a == nil {
		return nil
	}
	var out []byte
	for _, si := range a.SignerInfos {
		si := si
		out = appendSubmessage(out, 1, func() []byte { return si.Marshal() })
	}
	if a.Fee != nil {
		fee := a.Fee
		out = appendSubmessage(out, 2, func() []byte { return fee.Marshal() })
	}
	return out
}

// TxBody mirrors cosmos.tx.v1beta1.TxBody.
type TxBody struct {
	Messages                    []*Any
	Memo                        string
	TimeoutHeight               uint64
	ExtensionOptions            []*Any
	NonCriticalExtensionOptions []*Any
}

// Marshal encodes TxBody.
func (t *TxBody) Marshal() []byte {
	if t == nil {
		return nil
	}
	var out []byte
	for _, m := range t.Messages {
		m := m
		out = appendSubmessage(out, 1, func() []byte { return m.Marshal() })
	}
	out = appendString(out, 2, t.Memo)
	out = appendUint64(out, 3, t.TimeoutHeight)
	for _, e := range t.ExtensionOptions {
		e := e
		out = appendSubmessage(out, 1023, func() []byte { return e.Marshal() })
	}
	for _, e := range t.NonCriticalExtensionOptions {
		e := e
		out = appendSubmessage(out, 2047, func() []byte { return e.Marshal() })
	}
	return out
}

// TxRaw mirrors cosmos.tx.v1beta1.TxRaw — the canonical signed tx.
type TxRaw struct {
	BodyBytes     []byte
	AuthInfoBytes []byte
	Signatures    [][]byte
}

// Marshal encodes TxRaw.
func (t *TxRaw) Marshal() []byte {
	if t == nil {
		return nil
	}
	var out []byte
	out = appendBytes(out, 1, t.BodyBytes)
	out = appendBytes(out, 2, t.AuthInfoBytes)
	for _, sig := range t.Signatures {
		out = appendBytes(out, 3, sig)
	}
	return out
}

// UnmarshalTxRaw parses a TxRaw from bytes.
func UnmarshalTxRaw(b []byte) (*TxRaw, error) {
	out := &TxRaw{}
	for len(b) > 0 {
		num, _, val, n, err := consumeField(b)
		if err != nil {
			return nil, fmt.Errorf("txraw: %w", err)
		}
		switch num {
		case 1:
			out.BodyBytes = append([]byte(nil), val...)
		case 2:
			out.AuthInfoBytes = append([]byte(nil), val...)
		case 3:
			out.Signatures = append(out.Signatures, append([]byte(nil), val...))
		}
		b = b[n:]
	}
	return out, nil
}

// SignDoc mirrors cosmos.tx.v1beta1.SignDoc — the pre-image signed in
// SIGN_MODE_DIRECT.
type SignDoc struct {
	BodyBytes     []byte
	AuthInfoBytes []byte
	ChainID       string
	AccountNumber uint64
}

// Marshal encodes SignDoc.
func (s *SignDoc) Marshal() []byte {
	if s == nil {
		return nil
	}
	var out []byte
	out = appendBytes(out, 1, s.BodyBytes)
	out = appendBytes(out, 2, s.AuthInfoBytes)
	out = appendString(out, 3, s.ChainID)
	out = appendUint64(out, 4, s.AccountNumber)
	return out
}

// MsgWithdrawFeeTx is the starter Msg we exercise end-to-end in W2.
// Matches heimdall-v2/proto/heimdallv2/topup/tx.proto. Additional
// message types are added in W3/W4 with the same pattern.
type MsgWithdrawFeeTx struct {
	Proposer string
	// Amount is carried as a string because cosmossdk.io/math.Int is
	// encoded as its decimal string representation on the wire.
	Amount string
}

// MsgWithdrawFeeTxTypeURL is the Any type URL for the withdraw fee
// message as registered in heimdall-v2.
const MsgWithdrawFeeTxTypeURL = "/heimdallv2.topup.MsgWithdrawFeeTx"

// Marshal encodes MsgWithdrawFeeTx.
func (m *MsgWithdrawFeeTx) Marshal() []byte {
	if m == nil {
		return nil
	}
	var out []byte
	out = appendString(out, 1, m.Proposer)
	out = appendString(out, 2, m.Amount)
	return out
}

// UnmarshalMsgWithdrawFeeTx parses a MsgWithdrawFeeTx from bytes.
func UnmarshalMsgWithdrawFeeTx(b []byte) (*MsgWithdrawFeeTx, error) {
	out := &MsgWithdrawFeeTx{}
	for len(b) > 0 {
		num, _, val, n, err := consumeField(b)
		if err != nil {
			return nil, fmt.Errorf("MsgWithdrawFeeTx: %w", err)
		}
		switch num {
		case 1:
			out.Proposer = string(val)
		case 2:
			out.Amount = string(val)
		}
		b = b[n:]
	}
	return out, nil
}

// AsAny wraps the withdraw-fee message as a google.protobuf.Any for
// inclusion in TxBody.messages.
func (m *MsgWithdrawFeeTx) AsAny() *Any {
	return &Any{TypeURL: MsgWithdrawFeeTxTypeURL, Value: m.Marshal()}
}

// ensure protowire is referenced so goimports doesn't drop it in
// future edits that pare the surface down.
var _ = protowire.Number(0)
