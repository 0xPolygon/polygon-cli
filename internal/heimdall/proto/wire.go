// Package proto implements the narrow subset of the Cosmos SDK and
// Heimdall v2 protobuf wire format that polycli's heimdall tx builder
// needs. See README.md for the rationale — we encode/decode by hand
// rather than pull in cosmos-sdk's fork-pinned go module.
package proto

import (
	"fmt"

	"google.golang.org/protobuf/encoding/protowire"
)

// appendString writes (tag, length-prefixed bytes) for a proto3 string
// field. Zero-length strings are omitted to match the default proto3
// encoding.
func appendString(b []byte, fieldNum protowire.Number, v string) []byte {
	if v == "" {
		return b
	}
	b = protowire.AppendTag(b, fieldNum, protowire.BytesType)
	b = protowire.AppendString(b, v)
	return b
}

// appendBytes writes (tag, length-prefixed bytes) for a proto3 bytes
// field. Zero-length slices are omitted.
func appendBytes(b []byte, fieldNum protowire.Number, v []byte) []byte {
	if len(v) == 0 {
		return b
	}
	b = protowire.AppendTag(b, fieldNum, protowire.BytesType)
	b = protowire.AppendBytes(b, v)
	return b
}

// appendUint64 writes a proto3 varint field. Zero values are omitted.
func appendUint64(b []byte, fieldNum protowire.Number, v uint64) []byte {
	if v == 0 {
		return b
	}
	b = protowire.AppendTag(b, fieldNum, protowire.VarintType)
	b = protowire.AppendVarint(b, v)
	return b
}

// appendInt32 writes a proto3 varint field treating the value as a
// signed int32 encoded as a varint (per the proto3 spec for enums and
// non-zigzag int32). Zero values are omitted.
func appendInt32(b []byte, fieldNum protowire.Number, v int32) []byte {
	if v == 0 {
		return b
	}
	b = protowire.AppendTag(b, fieldNum, protowire.VarintType)
	b = protowire.AppendVarint(b, uint64(v))
	return b
}

// appendSubmessage encodes a nested message m as (tag, length, inner).
// Passing a nil encoder omits the field entirely (nullable submessages
// in proto3 are absent when unset).
func appendSubmessage(b []byte, fieldNum protowire.Number, encode func() []byte) []byte {
	if encode == nil {
		return b
	}
	inner := encode()
	b = protowire.AppendTag(b, fieldNum, protowire.BytesType)
	b = protowire.AppendBytes(b, inner)
	return b
}

// consumeField reads one (number, type, value, n) record from the
// start of b. Returns an error on malformed input.
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
		// Re-encode the varint as raw bytes so callers can re-decode
		// uniformly; callers interpret based on fieldNum + expected
		// type.
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

// varint reads a varint from a slice returned by consumeField.
func varint(b []byte) (uint64, error) {
	v, n := protowire.ConsumeVarint(b)
	if n < 0 {
		return 0, fmt.Errorf("proto: invalid varint: %w", protowire.ParseError(n))
	}
	return v, nil
}

// appendRawVarint appends a raw varint (no tag) to b. Used by packed
// repeated scalar encodings.
func appendRawVarint(b []byte, v uint64) []byte {
	return protowire.AppendVarint(b, v)
}

// consumePlainVarint reads one raw varint (no tag) from b and returns
// its value plus the number of bytes consumed.
func consumePlainVarint(b []byte) (uint64, int, error) {
	v, n := protowire.ConsumeVarint(b)
	if n < 0 {
		return 0, 0, fmt.Errorf("proto: invalid varint: %w", protowire.ParseError(n))
	}
	return v, n, nil
}
