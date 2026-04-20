package decode

import (
	"fmt"

	"google.golang.org/protobuf/encoding/protowire"
)

// consumeField consumes one proto wire-format record from the start of
// b. Returns (number, type, value bytes, total bytes consumed, error).
//
// Duplicates internal/heimdall/proto.consumeField because the decode
// package is in a different module tree from the parser helpers. Keeping
// the copy tiny is safer than exporting proto internals.
func consumeField(b []byte) (protowire.Number, protowire.Type, []byte, int, error) {
	num, typ, tagLen := protowire.ConsumeTag(b)
	if tagLen < 0 {
		return 0, 0, nil, 0, fmt.Errorf("proto: invalid tag: %w", protowire.ParseError(tagLen))
	}
	switch typ {
	case protowire.VarintType:
		v, n := protowire.ConsumeVarint(b[tagLen:])
		if n < 0 {
			return 0, 0, nil, 0, fmt.Errorf("proto: invalid varint: %w", protowire.ParseError(n))
		}
		buf := protowire.AppendVarint(nil, v)
		return num, typ, buf, tagLen + n, nil
	case protowire.BytesType:
		v, n := protowire.ConsumeBytes(b[tagLen:])
		if n < 0 {
			return 0, 0, nil, 0, fmt.Errorf("proto: invalid bytes: %w", protowire.ParseError(n))
		}
		return num, typ, v, tagLen + n, nil
	case protowire.Fixed32Type:
		v, n := protowire.ConsumeFixed32(b[tagLen:])
		if n < 0 {
			return 0, 0, nil, 0, fmt.Errorf("proto: invalid fixed32: %w", protowire.ParseError(n))
		}
		buf := protowire.AppendFixed32(nil, v)
		return num, typ, buf, tagLen + n, nil
	case protowire.Fixed64Type:
		v, n := protowire.ConsumeFixed64(b[tagLen:])
		if n < 0 {
			return 0, 0, nil, 0, fmt.Errorf("proto: invalid fixed64: %w", protowire.ParseError(n))
		}
		buf := protowire.AppendFixed64(nil, v)
		return num, typ, buf, tagLen + n, nil
	default:
		return 0, 0, nil, 0, fmt.Errorf("proto: unsupported wire type %d", typ)
	}
}

// rawVarint reads a raw varint from a slice returned by consumeField.
func rawVarint(b []byte) (uint64, error) {
	v, n := protowire.ConsumeVarint(b)
	if n < 0 {
		return 0, fmt.Errorf("proto: invalid varint: %w", protowire.ParseError(n))
	}
	return v, nil
}
