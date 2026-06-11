package proto

import "fmt"

// BlockIDFlag mirrors tendermint.types.BlockIDFlag.
type BlockIDFlag int32

// BlockIDFlag values. UNKNOWN is encoded as 0 and omitted when unset.
const (
	BlockIDFlagUnknown BlockIDFlag = 0
	BlockIDFlagAbsent  BlockIDFlag = 1
	BlockIDFlagCommit  BlockIDFlag = 2
	BlockIDFlagNil     BlockIDFlag = 3
)

// String returns the short enum name (without the BLOCK_ID_FLAG_ prefix
// used in the .proto source) for compact table rendering.
func (f BlockIDFlag) String() string {
	switch f {
	case BlockIDFlagAbsent:
		return "ABSENT"
	case BlockIDFlagCommit:
		return "COMMIT"
	case BlockIDFlagNil:
		return "NIL"
	case BlockIDFlagUnknown:
		return "UNKNOWN"
	default:
		return fmt.Sprintf("BlockIDFlag(%d)", int32(f))
	}
}

// ExtValidator mirrors tendermint.abci.Validator. Field 2 is reserved
// in the .proto source (it once held a pub_key).
type ExtValidator struct {
	Address []byte // 20-byte CometBFT address == heimdall signer address
	Power   int64
}

// Marshal encodes ExtValidator.
func (v ExtValidator) Marshal() []byte {
	var out []byte
	out = appendBytes(out, 1, v.Address)
	out = appendUint64(out, 3, uint64(v.Power))
	return out
}

func unmarshalExtValidator(b []byte) (ExtValidator, error) {
	var out ExtValidator
	for len(b) > 0 {
		num, _, val, n, err := consumeField(b)
		if err != nil {
			return out, fmt.Errorf("ExtValidator: %w", err)
		}
		switch num {
		case 1:
			out.Address = append([]byte(nil), val...)
		case 3:
			v, err := varint(val)
			if err != nil {
				return out, err
			}
			out.Power = int64(v)
		}
		b = b[n:]
	}
	return out, nil
}

// ExtendedVoteInfo mirrors tendermint.abci.ExtendedVoteInfo as defined
// by the 0xPolygon cometbft fork (v0.3.6-polygon), which adds the
// non-RP extension fields 6 and 7 on top of the upstream message.
// Field 2 is reserved in the .proto source.
type ExtendedVoteInfo struct {
	Validator               ExtValidator
	VoteExtension           []byte
	ExtensionSignature      []byte
	BlockIDFlag             BlockIDFlag
	NonRpVoteExtension      []byte
	NonRpExtensionSignature []byte
}

// Marshal encodes ExtendedVoteInfo.
func (v ExtendedVoteInfo) Marshal() []byte {
	var out []byte
	val := v.Validator
	out = appendSubmessage(out, 1, func() []byte { return val.Marshal() })
	out = appendBytes(out, 3, v.VoteExtension)
	out = appendBytes(out, 4, v.ExtensionSignature)
	out = appendInt32(out, 5, int32(v.BlockIDFlag))
	out = appendBytes(out, 6, v.NonRpVoteExtension)
	out = appendBytes(out, 7, v.NonRpExtensionSignature)
	return out
}

func unmarshalExtendedVoteInfo(b []byte) (ExtendedVoteInfo, error) {
	var out ExtendedVoteInfo
	for len(b) > 0 {
		num, _, val, n, err := consumeField(b)
		if err != nil {
			return out, fmt.Errorf("ExtendedVoteInfo: %w", err)
		}
		switch num {
		case 1:
			v, err := unmarshalExtValidator(val)
			if err != nil {
				return out, err
			}
			out.Validator = v
		case 3:
			out.VoteExtension = append([]byte(nil), val...)
		case 4:
			out.ExtensionSignature = append([]byte(nil), val...)
		case 5:
			v, err := varint(val)
			if err != nil {
				return out, err
			}
			out.BlockIDFlag = BlockIDFlag(int32(v))
		case 6:
			out.NonRpVoteExtension = append([]byte(nil), val...)
		case 7:
			out.NonRpExtensionSignature = append([]byte(nil), val...)
		}
		b = b[n:]
	}
	return out, nil
}

// ExtendedCommitInfo mirrors tendermint.abci.ExtendedCommitInfo. On
// heimdall-v2 chains the proposer injects the previous height's
// extended commit (all validators' vote extensions) as the special
// transaction at index 0 of every block; decoding block.data.txs[0]
// with this function recovers each validator's vote extension bytes.
// Unknown fields are skipped for forward compatibility with fork
// drift.
type ExtendedCommitInfo struct {
	Round int32
	Votes []ExtendedVoteInfo
}

// Marshal encodes ExtendedCommitInfo.
func (e *ExtendedCommitInfo) Marshal() []byte {
	if e == nil {
		return nil
	}
	var out []byte
	out = appendInt32(out, 1, e.Round)
	for _, v := range e.Votes {
		v := v
		out = appendSubmessage(out, 2, func() []byte { return v.Marshal() })
	}
	return out
}

// UnmarshalExtendedCommitInfo parses the marshaled ExtendedCommitInfo
// bytes found at block.data.txs[0].
func UnmarshalExtendedCommitInfo(b []byte) (*ExtendedCommitInfo, error) {
	out := &ExtendedCommitInfo{}
	for len(b) > 0 {
		num, _, val, n, err := consumeField(b)
		if err != nil {
			return nil, fmt.Errorf("ExtendedCommitInfo: %w", err)
		}
		switch num {
		case 1:
			v, err := varint(val)
			if err != nil {
				return nil, err
			}
			out.Round = int32(v)
		case 2:
			vote, err := unmarshalExtendedVoteInfo(val)
			if err != nil {
				return nil, err
			}
			out.Votes = append(out.Votes, vote)
		}
		b = b[n:]
	}
	return out, nil
}
