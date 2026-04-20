package proto

import "fmt"

// MsgProposeSpanTypeURL is the Any type URL for MsgProposeSpan
// (heimdallv2/proto/heimdallv2/bor/tx.proto).
const MsgProposeSpanTypeURL = "/heimdallv2.bor.MsgProposeSpan"

// MsgProposeSpan mirrors heimdallv2.bor.MsgProposeSpan.
type MsgProposeSpan struct {
	SpanID     uint64
	Proposer   string
	StartBlock uint64
	EndBlock   uint64
	ChainID    string
	Seed       []byte
	SeedAuthor string
}

// Marshal encodes MsgProposeSpan.
func (m *MsgProposeSpan) Marshal() []byte {
	if m == nil {
		return nil
	}
	var out []byte
	out = appendUint64(out, 1, m.SpanID)
	out = appendString(out, 2, m.Proposer)
	out = appendUint64(out, 3, m.StartBlock)
	out = appendUint64(out, 4, m.EndBlock)
	out = appendString(out, 5, m.ChainID)
	out = appendBytes(out, 6, m.Seed)
	out = appendString(out, 7, m.SeedAuthor)
	return out
}

// UnmarshalMsgProposeSpan parses the message bytes.
func UnmarshalMsgProposeSpan(b []byte) (*MsgProposeSpan, error) {
	out := &MsgProposeSpan{}
	for len(b) > 0 {
		num, _, val, n, err := consumeField(b)
		if err != nil {
			return nil, fmt.Errorf("MsgProposeSpan: %w", err)
		}
		switch num {
		case 1:
			v, err := varint(val)
			if err != nil {
				return nil, err
			}
			out.SpanID = v
		case 2:
			out.Proposer = string(val)
		case 3:
			v, err := varint(val)
			if err != nil {
				return nil, err
			}
			out.StartBlock = v
		case 4:
			v, err := varint(val)
			if err != nil {
				return nil, err
			}
			out.EndBlock = v
		case 5:
			out.ChainID = string(val)
		case 6:
			out.Seed = append([]byte(nil), val...)
		case 7:
			out.SeedAuthor = string(val)
		}
		b = b[n:]
	}
	return out, nil
}

// AsAny wraps the message.
func (m *MsgProposeSpan) AsAny() *Any {
	return &Any{TypeURL: MsgProposeSpanTypeURL, Value: m.Marshal()}
}

// MsgBackfillSpansTypeURL is the Any type URL for MsgBackfillSpans.
const MsgBackfillSpansTypeURL = "/heimdallv2.bor.MsgBackfillSpans"

// MsgBackfillSpans mirrors heimdallv2.bor.MsgBackfillSpans.
type MsgBackfillSpans struct {
	Proposer        string
	ChainID         string
	LatestSpanID    uint64
	LatestBorSpanID uint64
}

// Marshal encodes MsgBackfillSpans.
func (m *MsgBackfillSpans) Marshal() []byte {
	if m == nil {
		return nil
	}
	var out []byte
	out = appendString(out, 1, m.Proposer)
	out = appendString(out, 2, m.ChainID)
	out = appendUint64(out, 3, m.LatestSpanID)
	out = appendUint64(out, 4, m.LatestBorSpanID)
	return out
}

// UnmarshalMsgBackfillSpans parses the message bytes.
func UnmarshalMsgBackfillSpans(b []byte) (*MsgBackfillSpans, error) {
	out := &MsgBackfillSpans{}
	for len(b) > 0 {
		num, _, val, n, err := consumeField(b)
		if err != nil {
			return nil, fmt.Errorf("MsgBackfillSpans: %w", err)
		}
		switch num {
		case 1:
			out.Proposer = string(val)
		case 2:
			out.ChainID = string(val)
		case 3:
			v, err := varint(val)
			if err != nil {
				return nil, err
			}
			out.LatestSpanID = v
		case 4:
			v, err := varint(val)
			if err != nil {
				return nil, err
			}
			out.LatestBorSpanID = v
		}
		b = b[n:]
	}
	return out, nil
}

// AsAny wraps the message.
func (m *MsgBackfillSpans) AsAny() *Any {
	return &Any{TypeURL: MsgBackfillSpansTypeURL, Value: m.Marshal()}
}

// MsgVoteProducersTypeURL is the Any type URL for MsgVoteProducers.
const MsgVoteProducersTypeURL = "/heimdallv2.bor.MsgVoteProducers"

// ProducerVotes mirrors heimdallv2.bor.ProducerVotes (repeated uint64
// field `votes`, encoded packed or repeated per proto3).
type ProducerVotes struct {
	Votes []uint64
}

// Marshal encodes ProducerVotes with packed varints (the proto3 default
// for repeated scalar fields). Accepts either packed or unpacked on the
// wire when unmarshalling.
func (p ProducerVotes) Marshal() []byte {
	if len(p.Votes) == 0 {
		return nil
	}
	// Packed encoding: field 1, wire type bytes, then concatenated
	// varints.
	var inner []byte
	for _, v := range p.Votes {
		inner = appendRawVarint(inner, v)
	}
	return appendBytes(nil, 1, inner)
}

// MsgVoteProducers mirrors heimdallv2.bor.MsgVoteProducers.
type MsgVoteProducers struct {
	Voter   string
	VoterID uint64
	Votes   ProducerVotes
}

// Marshal encodes MsgVoteProducers.
func (m *MsgVoteProducers) Marshal() []byte {
	if m == nil {
		return nil
	}
	var out []byte
	out = appendString(out, 1, m.Voter)
	out = appendUint64(out, 2, m.VoterID)
	if inner := m.Votes.Marshal(); len(inner) > 0 {
		out = appendBytes(out, 3, inner)
	}
	return out
}

// UnmarshalMsgVoteProducers parses the message bytes. Votes may arrive
// either packed or unpacked per proto3 rules.
func UnmarshalMsgVoteProducers(b []byte) (*MsgVoteProducers, error) {
	out := &MsgVoteProducers{}
	for len(b) > 0 {
		num, _, val, n, err := consumeField(b)
		if err != nil {
			return nil, fmt.Errorf("MsgVoteProducers: %w", err)
		}
		switch num {
		case 1:
			out.Voter = string(val)
		case 2:
			v, err := varint(val)
			if err != nil {
				return nil, err
			}
			out.VoterID = v
		case 3:
			// Nested ProducerVotes. Walk its inner bytes.
			inner := val
			for len(inner) > 0 {
				inNum, _, inVal, inN, err := consumeField(inner)
				if err != nil {
					return nil, fmt.Errorf("MsgVoteProducers.votes: %w", err)
				}
				if inNum == 1 {
					// Packed: inVal is a byte-string of concatenated
					// varints.
					rem := inVal
					for len(rem) > 0 {
						v, consumed, err := consumePlainVarint(rem)
						if err != nil {
							return nil, fmt.Errorf("MsgVoteProducers.votes packed: %w", err)
						}
						out.Votes.Votes = append(out.Votes.Votes, v)
						rem = rem[consumed:]
					}
				}
				inner = inner[inN:]
			}
		}
		b = b[n:]
	}
	return out, nil
}

// AsAny wraps the message.
func (m *MsgVoteProducers) AsAny() *Any {
	return &Any{TypeURL: MsgVoteProducersTypeURL, Value: m.Marshal()}
}

// MsgSetProducerDowntimeTypeURL is the Any type URL for
// MsgSetProducerDowntime.
const MsgSetProducerDowntimeTypeURL = "/heimdallv2.bor.MsgSetProducerDowntime"

// BlockRange mirrors heimdallv2.bor.BlockRange.
type BlockRange struct {
	StartBlock uint64
	EndBlock   uint64
}

// Marshal encodes BlockRange.
func (r BlockRange) Marshal() []byte {
	var out []byte
	out = appendUint64(out, 1, r.StartBlock)
	out = appendUint64(out, 2, r.EndBlock)
	return out
}

// MsgSetProducerDowntime mirrors heimdallv2.bor.MsgSetProducerDowntime.
type MsgSetProducerDowntime struct {
	Producer      string
	DowntimeRange BlockRange
}

// Marshal encodes MsgSetProducerDowntime.
func (m *MsgSetProducerDowntime) Marshal() []byte {
	if m == nil {
		return nil
	}
	var out []byte
	out = appendString(out, 1, m.Producer)
	inner := m.DowntimeRange.Marshal()
	if len(inner) > 0 {
		out = appendBytes(out, 2, inner)
	} else {
		// A zero-valued BlockRange still needs to be present for the
		// signer since the proto is non-nullable. Emit an explicit
		// empty submessage.
		out = appendBytes(out, 2, []byte{})
	}
	return out
}

// UnmarshalMsgSetProducerDowntime parses the message bytes.
func UnmarshalMsgSetProducerDowntime(b []byte) (*MsgSetProducerDowntime, error) {
	out := &MsgSetProducerDowntime{}
	for len(b) > 0 {
		num, _, val, n, err := consumeField(b)
		if err != nil {
			return nil, fmt.Errorf("MsgSetProducerDowntime: %w", err)
		}
		switch num {
		case 1:
			out.Producer = string(val)
		case 2:
			// Nested BlockRange.
			inner := val
			for len(inner) > 0 {
				inNum, _, inVal, inN, err := consumeField(inner)
				if err != nil {
					return nil, fmt.Errorf("MsgSetProducerDowntime.range: %w", err)
				}
				switch inNum {
				case 1:
					v, err := varint(inVal)
					if err != nil {
						return nil, err
					}
					out.DowntimeRange.StartBlock = v
				case 2:
					v, err := varint(inVal)
					if err != nil {
						return nil, err
					}
					out.DowntimeRange.EndBlock = v
				}
				inner = inner[inN:]
			}
		}
		b = b[n:]
	}
	return out, nil
}

// AsAny wraps the message.
func (m *MsgSetProducerDowntime) AsAny() *Any {
	return &Any{TypeURL: MsgSetProducerDowntimeTypeURL, Value: m.Marshal()}
}
