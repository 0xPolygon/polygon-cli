package proto

import "fmt"

// Vote mirrors heimdallv2.sidetxs.Vote.
type Vote int32

// Vote values. UNSPECIFIED is encoded as 0 and omitted when unset.
const (
	VoteUnspecified Vote = 0
	VoteYes         Vote = 1
	VoteNo          Vote = 2
)

// String returns the enum name, matching the .proto source.
func (v Vote) String() string {
	switch v {
	case VoteYes:
		return "VOTE_YES"
	case VoteNo:
		return "VOTE_NO"
	case VoteUnspecified:
		return "UNSPECIFIED"
	default:
		return fmt.Sprintf("Vote(%d)", int32(v))
	}
}

// SideTxResponse mirrors heimdallv2.sidetxs.SideTxResponse.
type SideTxResponse struct {
	TxHash []byte
	Result Vote
}

// Marshal encodes SideTxResponse.
func (s SideTxResponse) Marshal() []byte {
	var out []byte
	out = appendBytes(out, 1, s.TxHash)
	out = appendInt32(out, 2, int32(s.Result))
	return out
}

// UnmarshalSideTxResponse parses the message bytes.
func UnmarshalSideTxResponse(b []byte) (SideTxResponse, error) {
	var out SideTxResponse
	for len(b) > 0 {
		num, _, val, n, err := consumeField(b)
		if err != nil {
			return out, fmt.Errorf("SideTxResponse: %w", err)
		}
		switch num {
		case 1:
			out.TxHash = append([]byte(nil), val...)
		case 2:
			v, err := varint(val)
			if err != nil {
				return out, err
			}
			out.Result = Vote(int32(v))
		}
		b = b[n:]
	}
	return out, nil
}

// MilestoneProposition mirrors heimdallv2.milestone.MilestoneProposition.
// Repeated uint64 fields are decoded tolerantly: both packed and
// unpacked wire forms are accepted.
type MilestoneProposition struct {
	BlockHashes      [][]byte
	StartBlockNumber uint64
	ParentHash       []byte
	BlockTDs         []uint64
}

// Marshal encodes MilestoneProposition (block_tds as packed varints).
func (m *MilestoneProposition) Marshal() []byte {
	if m == nil {
		return nil
	}
	var out []byte
	for _, h := range m.BlockHashes {
		out = appendBytes(out, 1, h)
	}
	out = appendUint64(out, 2, m.StartBlockNumber)
	out = appendBytes(out, 3, m.ParentHash)
	if len(m.BlockTDs) > 0 {
		var packed []byte
		for _, v := range m.BlockTDs {
			packed = appendRawVarint(packed, v)
		}
		out = appendBytes(out, 4, packed)
	}
	return out
}

// UnmarshalMilestoneProposition parses the message bytes.
//
// The block_tds field is a repeated uint64. Proto3 defaults repeated
// scalars to the packed wire format; heimdalld's gogoproto descriptors
// use packed encoding. We accept both: consumeField normalizes the
// value to a raw byte-slice regardless of the tag wire type.
func UnmarshalMilestoneProposition(b []byte) (*MilestoneProposition, error) {
	out := &MilestoneProposition{}
	for len(b) > 0 {
		num, typ, val, n, err := consumeField(b)
		if err != nil {
			return nil, fmt.Errorf("MilestoneProposition: %w", err)
		}
		switch num {
		case 1:
			out.BlockHashes = append(out.BlockHashes, append([]byte(nil), val...))
		case 2:
			v, err := varint(val)
			if err != nil {
				return nil, err
			}
			out.StartBlockNumber = v
		case 3:
			out.ParentHash = append([]byte(nil), val...)
		case 4:
			// BytesType = packed (a concatenation of varints).
			// VarintType = a single unpacked value (legacy encoders).
			if typ == 2 { // protowire.BytesType
				rem := val
				for len(rem) > 0 {
					v, consumed, err := consumePlainVarint(rem)
					if err != nil {
						return nil, err
					}
					out.BlockTDs = append(out.BlockTDs, v)
					rem = rem[consumed:]
				}
			} else {
				v, err := varint(val)
				if err != nil {
					return nil, err
				}
				out.BlockTDs = append(out.BlockTDs, v)
			}
		}
		b = b[n:]
	}
	return out, nil
}

// VoteExtensionTypeURL is informational; VoteExtension is not wrapped
// in Any on the wire (it arrives as plain bytes on CometBFT's
// ExtendVote interface). Retained so the decode family can surface a
// consistent label.
const VoteExtensionTypeURL = "/heimdallv2.sidetxs.VoteExtension"

// VoteExtension mirrors heimdallv2.sidetxs.VoteExtension.
type VoteExtension struct {
	BlockHash            []byte
	Height               int64
	SideTxResponses      []SideTxResponse
	MilestoneProposition *MilestoneProposition
}

// Marshal encodes VoteExtension.
func (v *VoteExtension) Marshal() []byte {
	if v == nil {
		return nil
	}
	var out []byte
	out = appendBytes(out, 1, v.BlockHash)
	// height is int64 proto3; encoded as a varint with two's-complement
	// for negative values. Heimdall's heights are always positive so
	// the usual uint64 encoding suffices.
	out = appendUint64(out, 2, uint64(v.Height))
	for _, r := range v.SideTxResponses {
		r := r
		out = appendSubmessage(out, 3, func() []byte { return r.Marshal() })
	}
	if v.MilestoneProposition != nil {
		mp := v.MilestoneProposition
		out = appendSubmessage(out, 4, func() []byte { return mp.Marshal() })
	}
	return out
}

// UnmarshalVoteExtension parses raw vote-extension bytes as emitted by
// heimdall-v2's ExtendVote handler.
func UnmarshalVoteExtension(b []byte) (*VoteExtension, error) {
	out := &VoteExtension{}
	for len(b) > 0 {
		num, _, val, n, err := consumeField(b)
		if err != nil {
			return nil, fmt.Errorf("VoteExtension: %w", err)
		}
		switch num {
		case 1:
			out.BlockHash = append([]byte(nil), val...)
		case 2:
			v, err := varint(val)
			if err != nil {
				return nil, err
			}
			out.Height = int64(v)
		case 3:
			r, err := UnmarshalSideTxResponse(val)
			if err != nil {
				return nil, err
			}
			out.SideTxResponses = append(out.SideTxResponses, r)
		case 4:
			mp, err := UnmarshalMilestoneProposition(val)
			if err != nil {
				return nil, err
			}
			out.MilestoneProposition = mp
		}
		b = b[n:]
	}
	return out, nil
}
