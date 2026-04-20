package proto

import "fmt"

// MsgEventRecordTypeURL is the Any type URL for MsgEventRecord
// (heimdallv2/proto/heimdallv2/clerk/tx.proto).
const MsgEventRecordTypeURL = "/heimdallv2.clerk.MsgEventRecord"

// MsgEventRecord mirrors heimdallv2.clerk.MsgEventRecord.
type MsgEventRecord struct {
	From            string
	TxHash          string
	LogIndex        uint64
	BlockNumber     uint64
	ContractAddress string
	Data            []byte
	ID              uint64
	ChainID         string
}

// Marshal encodes MsgEventRecord.
func (m *MsgEventRecord) Marshal() []byte {
	if m == nil {
		return nil
	}
	var out []byte
	out = appendString(out, 1, m.From)
	out = appendString(out, 2, m.TxHash)
	out = appendUint64(out, 3, m.LogIndex)
	out = appendUint64(out, 4, m.BlockNumber)
	out = appendString(out, 5, m.ContractAddress)
	out = appendBytes(out, 6, m.Data)
	out = appendUint64(out, 7, m.ID)
	out = appendString(out, 8, m.ChainID)
	return out
}

// UnmarshalMsgEventRecord parses the message bytes.
func UnmarshalMsgEventRecord(b []byte) (*MsgEventRecord, error) {
	out := &MsgEventRecord{}
	for len(b) > 0 {
		num, _, val, n, err := consumeField(b)
		if err != nil {
			return nil, fmt.Errorf("MsgEventRecord: %w", err)
		}
		switch num {
		case 1:
			out.From = string(val)
		case 2:
			out.TxHash = string(val)
		case 3:
			v, err := varint(val)
			if err != nil {
				return nil, err
			}
			out.LogIndex = v
		case 4:
			v, err := varint(val)
			if err != nil {
				return nil, err
			}
			out.BlockNumber = v
		case 5:
			out.ContractAddress = string(val)
		case 6:
			out.Data = append([]byte(nil), val...)
		case 7:
			v, err := varint(val)
			if err != nil {
				return nil, err
			}
			out.ID = v
		case 8:
			out.ChainID = string(val)
		}
		b = b[n:]
	}
	return out, nil
}

// AsAny wraps the message.
func (m *MsgEventRecord) AsAny() *Any {
	return &Any{TypeURL: MsgEventRecordTypeURL, Value: m.Marshal()}
}
