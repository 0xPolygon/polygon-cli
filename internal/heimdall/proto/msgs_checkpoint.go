package proto

import "fmt"

// MsgCheckpointTypeURL is the Any type URL for MsgCheckpoint
// (heimdallv2/proto/heimdallv2/checkpoint/tx.proto).
const MsgCheckpointTypeURL = "/heimdallv2.checkpoint.MsgCheckpoint"

// MsgCheckpoint mirrors heimdallv2.checkpoint.MsgCheckpoint. Fields are
// ordered to match the .proto source; the Marshal emits them in field
// number order regardless.
type MsgCheckpoint struct {
	Proposer        string
	StartBlock      uint64
	EndBlock        uint64
	RootHash        []byte
	AccountRootHash []byte
	BorChainID      string
}

// Marshal encodes MsgCheckpoint.
func (m *MsgCheckpoint) Marshal() []byte {
	if m == nil {
		return nil
	}
	var out []byte
	out = appendString(out, 1, m.Proposer)
	out = appendUint64(out, 2, m.StartBlock)
	out = appendUint64(out, 3, m.EndBlock)
	out = appendBytes(out, 4, m.RootHash)
	out = appendBytes(out, 5, m.AccountRootHash)
	out = appendString(out, 6, m.BorChainID)
	return out
}

// UnmarshalMsgCheckpoint parses a MsgCheckpoint from its proto bytes.
func UnmarshalMsgCheckpoint(b []byte) (*MsgCheckpoint, error) {
	out := &MsgCheckpoint{}
	for len(b) > 0 {
		num, _, val, n, err := consumeField(b)
		if err != nil {
			return nil, fmt.Errorf("MsgCheckpoint: %w", err)
		}
		switch num {
		case 1:
			out.Proposer = string(val)
		case 2:
			v, err := varint(val)
			if err != nil {
				return nil, err
			}
			out.StartBlock = v
		case 3:
			v, err := varint(val)
			if err != nil {
				return nil, err
			}
			out.EndBlock = v
		case 4:
			out.RootHash = append([]byte(nil), val...)
		case 5:
			out.AccountRootHash = append([]byte(nil), val...)
		case 6:
			out.BorChainID = string(val)
		}
		b = b[n:]
	}
	return out, nil
}

// AsAny wraps the message as a google.protobuf.Any.
func (m *MsgCheckpoint) AsAny() *Any {
	return &Any{TypeURL: MsgCheckpointTypeURL, Value: m.Marshal()}
}

// MsgCpAckTypeURL is the Any type URL for MsgCpAck.
const MsgCpAckTypeURL = "/heimdallv2.checkpoint.MsgCpAck"

// MsgCpAck mirrors heimdallv2.checkpoint.MsgCpAck.
type MsgCpAck struct {
	From       string
	Number     uint64
	Proposer   string
	StartBlock uint64
	EndBlock   uint64
	RootHash   []byte
}

// Marshal encodes MsgCpAck.
func (m *MsgCpAck) Marshal() []byte {
	if m == nil {
		return nil
	}
	var out []byte
	out = appendString(out, 1, m.From)
	out = appendUint64(out, 2, m.Number)
	out = appendString(out, 3, m.Proposer)
	out = appendUint64(out, 4, m.StartBlock)
	out = appendUint64(out, 5, m.EndBlock)
	out = appendBytes(out, 6, m.RootHash)
	return out
}

// UnmarshalMsgCpAck parses MsgCpAck bytes.
func UnmarshalMsgCpAck(b []byte) (*MsgCpAck, error) {
	out := &MsgCpAck{}
	for len(b) > 0 {
		num, _, val, n, err := consumeField(b)
		if err != nil {
			return nil, fmt.Errorf("MsgCpAck: %w", err)
		}
		switch num {
		case 1:
			out.From = string(val)
		case 2:
			v, err := varint(val)
			if err != nil {
				return nil, err
			}
			out.Number = v
		case 3:
			out.Proposer = string(val)
		case 4:
			v, err := varint(val)
			if err != nil {
				return nil, err
			}
			out.StartBlock = v
		case 5:
			v, err := varint(val)
			if err != nil {
				return nil, err
			}
			out.EndBlock = v
		case 6:
			out.RootHash = append([]byte(nil), val...)
		}
		b = b[n:]
	}
	return out, nil
}

// AsAny wraps the message.
func (m *MsgCpAck) AsAny() *Any {
	return &Any{TypeURL: MsgCpAckTypeURL, Value: m.Marshal()}
}

// MsgCpNoAckTypeURL is the Any type URL for MsgCpNoAck.
const MsgCpNoAckTypeURL = "/heimdallv2.checkpoint.MsgCpNoAck"

// MsgCpNoAck mirrors heimdallv2.checkpoint.MsgCpNoAck (single field).
type MsgCpNoAck struct {
	From string
}

// Marshal encodes MsgCpNoAck.
func (m *MsgCpNoAck) Marshal() []byte {
	if m == nil {
		return nil
	}
	return appendString(nil, 1, m.From)
}

// UnmarshalMsgCpNoAck parses MsgCpNoAck bytes.
func UnmarshalMsgCpNoAck(b []byte) (*MsgCpNoAck, error) {
	out := &MsgCpNoAck{}
	for len(b) > 0 {
		num, _, val, n, err := consumeField(b)
		if err != nil {
			return nil, fmt.Errorf("MsgCpNoAck: %w", err)
		}
		if num == 1 {
			out.From = string(val)
		}
		b = b[n:]
	}
	return out, nil
}

// AsAny wraps the message.
func (m *MsgCpNoAck) AsAny() *Any {
	return &Any{TypeURL: MsgCpNoAckTypeURL, Value: m.Marshal()}
}
