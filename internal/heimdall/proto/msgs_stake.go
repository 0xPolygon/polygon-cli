package proto

import "google.golang.org/protobuf/encoding/protowire"

// Stake module Msg type URLs.
const (
	MsgValidatorJoinTypeURL = "/heimdallv2.stake.MsgValidatorJoin"
	MsgStakeUpdateTypeURL   = "/heimdallv2.stake.MsgStakeUpdate"
	MsgSignerUpdateTypeURL  = "/heimdallv2.stake.MsgSignerUpdate"
	MsgValidatorExitTypeURL = "/heimdallv2.stake.MsgValidatorExit"
)

// MsgValidatorJoin mirrors heimdallv2.stake.MsgValidatorJoin.
type MsgValidatorJoin struct {
	From            string
	ValID           uint64
	ActivationEpoch uint64
	Amount          string
	SignerPubKey    []byte
	TxHash          []byte
	LogIndex        uint64
	BlockNumber     uint64
	Nonce           uint64
}

// Marshal encodes MsgValidatorJoin.
func (m *MsgValidatorJoin) Marshal() []byte {
	if m == nil {
		return nil
	}
	var out []byte
	out = appendString(out, 1, m.From)
	out = appendUint64(out, 2, m.ValID)
	out = appendUint64(out, 3, m.ActivationEpoch)
	out = appendString(out, 4, m.Amount)
	out = appendBytes(out, 5, m.SignerPubKey)
	out = appendBytes(out, 6, m.TxHash)
	out = appendUint64(out, 7, m.LogIndex)
	out = appendUint64(out, 8, m.BlockNumber)
	out = appendUint64(out, 9, m.Nonce)
	return out
}

// UnmarshalMsgValidatorJoin parses the message bytes.
func UnmarshalMsgValidatorJoin(b []byte) (*MsgValidatorJoin, error) {
	out := &MsgValidatorJoin{}
	if err := unmarshalFields(b, "MsgValidatorJoin", map[protowire.Number]fieldHandler{
		1: setString(&out.From),
		2: setUint64(&out.ValID),
		3: setUint64(&out.ActivationEpoch),
		4: setString(&out.Amount),
		5: setBytes(&out.SignerPubKey),
		6: setBytes(&out.TxHash),
		7: setUint64(&out.LogIndex),
		8: setUint64(&out.BlockNumber),
		9: setUint64(&out.Nonce),
	}); err != nil {
		return nil, err
	}
	return out, nil
}

// AsAny wraps the message.
func (m *MsgValidatorJoin) AsAny() *Any {
	return &Any{TypeURL: MsgValidatorJoinTypeURL, Value: m.Marshal()}
}

// MsgStakeUpdate mirrors heimdallv2.stake.MsgStakeUpdate.
type MsgStakeUpdate struct {
	From        string
	ValID       uint64
	NewAmount   string
	TxHash      []byte
	LogIndex    uint64
	BlockNumber uint64
	Nonce       uint64
}

// Marshal encodes MsgStakeUpdate.
func (m *MsgStakeUpdate) Marshal() []byte {
	if m == nil {
		return nil
	}
	var out []byte
	out = appendString(out, 1, m.From)
	out = appendUint64(out, 2, m.ValID)
	out = appendString(out, 3, m.NewAmount)
	out = appendBytes(out, 4, m.TxHash)
	out = appendUint64(out, 5, m.LogIndex)
	out = appendUint64(out, 6, m.BlockNumber)
	out = appendUint64(out, 7, m.Nonce)
	return out
}

// UnmarshalMsgStakeUpdate parses the message bytes.
func UnmarshalMsgStakeUpdate(b []byte) (*MsgStakeUpdate, error) {
	out := &MsgStakeUpdate{}
	if err := unmarshalFields(b, "MsgStakeUpdate", map[protowire.Number]fieldHandler{
		1: setString(&out.From),
		2: setUint64(&out.ValID),
		3: setString(&out.NewAmount),
		4: setBytes(&out.TxHash),
		5: setUint64(&out.LogIndex),
		6: setUint64(&out.BlockNumber),
		7: setUint64(&out.Nonce),
	}); err != nil {
		return nil, err
	}
	return out, nil
}

// AsAny wraps the message.
func (m *MsgStakeUpdate) AsAny() *Any {
	return &Any{TypeURL: MsgStakeUpdateTypeURL, Value: m.Marshal()}
}

// MsgSignerUpdate mirrors heimdallv2.stake.MsgSignerUpdate.
type MsgSignerUpdate struct {
	From            string
	ValID           uint64
	NewSignerPubKey []byte
	TxHash          []byte
	LogIndex        uint64
	BlockNumber     uint64
	Nonce           uint64
}

// Marshal encodes MsgSignerUpdate.
func (m *MsgSignerUpdate) Marshal() []byte {
	if m == nil {
		return nil
	}
	var out []byte
	out = appendString(out, 1, m.From)
	out = appendUint64(out, 2, m.ValID)
	out = appendBytes(out, 3, m.NewSignerPubKey)
	out = appendBytes(out, 4, m.TxHash)
	out = appendUint64(out, 5, m.LogIndex)
	out = appendUint64(out, 6, m.BlockNumber)
	out = appendUint64(out, 7, m.Nonce)
	return out
}

// UnmarshalMsgSignerUpdate parses the message bytes.
func UnmarshalMsgSignerUpdate(b []byte) (*MsgSignerUpdate, error) {
	out := &MsgSignerUpdate{}
	if err := unmarshalFields(b, "MsgSignerUpdate", map[protowire.Number]fieldHandler{
		1: setString(&out.From),
		2: setUint64(&out.ValID),
		3: setBytes(&out.NewSignerPubKey),
		4: setBytes(&out.TxHash),
		5: setUint64(&out.LogIndex),
		6: setUint64(&out.BlockNumber),
		7: setUint64(&out.Nonce),
	}); err != nil {
		return nil, err
	}
	return out, nil
}

// AsAny wraps the message.
func (m *MsgSignerUpdate) AsAny() *Any {
	return &Any{TypeURL: MsgSignerUpdateTypeURL, Value: m.Marshal()}
}

// MsgValidatorExit mirrors heimdallv2.stake.MsgValidatorExit.
type MsgValidatorExit struct {
	From              string
	ValID             uint64
	DeactivationEpoch uint64
	TxHash            []byte
	LogIndex          uint64
	BlockNumber       uint64
	Nonce             uint64
}

// Marshal encodes MsgValidatorExit.
func (m *MsgValidatorExit) Marshal() []byte {
	if m == nil {
		return nil
	}
	var out []byte
	out = appendString(out, 1, m.From)
	out = appendUint64(out, 2, m.ValID)
	out = appendUint64(out, 3, m.DeactivationEpoch)
	out = appendBytes(out, 4, m.TxHash)
	out = appendUint64(out, 5, m.LogIndex)
	out = appendUint64(out, 6, m.BlockNumber)
	out = appendUint64(out, 7, m.Nonce)
	return out
}

// UnmarshalMsgValidatorExit parses the message bytes.
func UnmarshalMsgValidatorExit(b []byte) (*MsgValidatorExit, error) {
	out := &MsgValidatorExit{}
	if err := unmarshalFields(b, "MsgValidatorExit", map[protowire.Number]fieldHandler{
		1: setString(&out.From),
		2: setUint64(&out.ValID),
		3: setUint64(&out.DeactivationEpoch),
		4: setBytes(&out.TxHash),
		5: setUint64(&out.LogIndex),
		6: setUint64(&out.BlockNumber),
		7: setUint64(&out.Nonce),
	}); err != nil {
		return nil, err
	}
	return out, nil
}

// AsAny wraps the message.
func (m *MsgValidatorExit) AsAny() *Any {
	return &Any{TypeURL: MsgValidatorExitTypeURL, Value: m.Marshal()}
}
