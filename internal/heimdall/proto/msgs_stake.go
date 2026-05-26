package proto

import "fmt"

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
	for len(b) > 0 {
		num, _, val, n, err := consumeField(b)
		if err != nil {
			return nil, fmt.Errorf("MsgValidatorJoin: %w", err)
		}
		switch num {
		case 1:
			out.From = string(val)
		case 2:
			v, err := varint(val)
			if err != nil {
				return nil, err
			}
			out.ValID = v
		case 3:
			v, err := varint(val)
			if err != nil {
				return nil, err
			}
			out.ActivationEpoch = v
		case 4:
			out.Amount = string(val)
		case 5:
			out.SignerPubKey = append([]byte(nil), val...)
		case 6:
			out.TxHash = append([]byte(nil), val...)
		case 7:
			v, err := varint(val)
			if err != nil {
				return nil, err
			}
			out.LogIndex = v
		case 8:
			v, err := varint(val)
			if err != nil {
				return nil, err
			}
			out.BlockNumber = v
		case 9:
			v, err := varint(val)
			if err != nil {
				return nil, err
			}
			out.Nonce = v
		}
		b = b[n:]
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
	for len(b) > 0 {
		num, _, val, n, err := consumeField(b)
		if err != nil {
			return nil, fmt.Errorf("MsgStakeUpdate: %w", err)
		}
		switch num {
		case 1:
			out.From = string(val)
		case 2:
			v, err := varint(val)
			if err != nil {
				return nil, err
			}
			out.ValID = v
		case 3:
			out.NewAmount = string(val)
		case 4:
			out.TxHash = append([]byte(nil), val...)
		case 5:
			v, err := varint(val)
			if err != nil {
				return nil, err
			}
			out.LogIndex = v
		case 6:
			v, err := varint(val)
			if err != nil {
				return nil, err
			}
			out.BlockNumber = v
		case 7:
			v, err := varint(val)
			if err != nil {
				return nil, err
			}
			out.Nonce = v
		}
		b = b[n:]
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
	for len(b) > 0 {
		num, _, val, n, err := consumeField(b)
		if err != nil {
			return nil, fmt.Errorf("MsgSignerUpdate: %w", err)
		}
		switch num {
		case 1:
			out.From = string(val)
		case 2:
			v, err := varint(val)
			if err != nil {
				return nil, err
			}
			out.ValID = v
		case 3:
			out.NewSignerPubKey = append([]byte(nil), val...)
		case 4:
			out.TxHash = append([]byte(nil), val...)
		case 5:
			v, err := varint(val)
			if err != nil {
				return nil, err
			}
			out.LogIndex = v
		case 6:
			v, err := varint(val)
			if err != nil {
				return nil, err
			}
			out.BlockNumber = v
		case 7:
			v, err := varint(val)
			if err != nil {
				return nil, err
			}
			out.Nonce = v
		}
		b = b[n:]
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
	for len(b) > 0 {
		num, _, val, n, err := consumeField(b)
		if err != nil {
			return nil, fmt.Errorf("MsgValidatorExit: %w", err)
		}
		switch num {
		case 1:
			out.From = string(val)
		case 2:
			v, err := varint(val)
			if err != nil {
				return nil, err
			}
			out.ValID = v
		case 3:
			v, err := varint(val)
			if err != nil {
				return nil, err
			}
			out.DeactivationEpoch = v
		case 4:
			out.TxHash = append([]byte(nil), val...)
		case 5:
			v, err := varint(val)
			if err != nil {
				return nil, err
			}
			out.LogIndex = v
		case 6:
			v, err := varint(val)
			if err != nil {
				return nil, err
			}
			out.BlockNumber = v
		case 7:
			v, err := varint(val)
			if err != nil {
				return nil, err
			}
			out.Nonce = v
		}
		b = b[n:]
	}
	return out, nil
}

// AsAny wraps the message.
func (m *MsgValidatorExit) AsAny() *Any {
	return &Any{TypeURL: MsgValidatorExitTypeURL, Value: m.Marshal()}
}
