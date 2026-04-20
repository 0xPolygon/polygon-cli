package proto

import "fmt"

// MsgTopupTxTypeURL is the Any type URL for MsgTopupTx
// (heimdallv2/proto/heimdallv2/topup/tx.proto).
const MsgTopupTxTypeURL = "/heimdallv2.topup.MsgTopupTx"

// MsgTopupTx mirrors heimdallv2.topup.MsgTopupTx. Fee is carried as a
// decimal string (math.Int) identical to MsgWithdrawFeeTx.
type MsgTopupTx struct {
	Proposer    string
	User        string
	Fee         string
	TxHash      []byte
	LogIndex    uint64
	BlockNumber uint64
}

// Marshal encodes MsgTopupTx.
func (m *MsgTopupTx) Marshal() []byte {
	if m == nil {
		return nil
	}
	var out []byte
	out = appendString(out, 1, m.Proposer)
	out = appendString(out, 2, m.User)
	out = appendString(out, 3, m.Fee)
	out = appendBytes(out, 4, m.TxHash)
	out = appendUint64(out, 5, m.LogIndex)
	out = appendUint64(out, 6, m.BlockNumber)
	return out
}

// UnmarshalMsgTopupTx parses the message bytes.
func UnmarshalMsgTopupTx(b []byte) (*MsgTopupTx, error) {
	out := &MsgTopupTx{}
	for len(b) > 0 {
		num, _, val, n, err := consumeField(b)
		if err != nil {
			return nil, fmt.Errorf("MsgTopupTx: %w", err)
		}
		switch num {
		case 1:
			out.Proposer = string(val)
		case 2:
			out.User = string(val)
		case 3:
			out.Fee = string(val)
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
		}
		b = b[n:]
	}
	return out, nil
}

// AsAny wraps the message.
func (m *MsgTopupTx) AsAny() *Any {
	return &Any{TypeURL: MsgTopupTxTypeURL, Value: m.Marshal()}
}
