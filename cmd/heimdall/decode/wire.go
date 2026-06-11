package decode

import (
	"google.golang.org/protobuf/encoding/protowire"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/proto"
)

// consumeField consumes one proto wire-format record from the start of
// b. Returns (number, type, value bytes, total bytes consumed, error).
// It delegates to internal/heimdall/proto so the decode command does
// not carry its own copy of the parser.
func consumeField(b []byte) (protowire.Number, protowire.Type, []byte, int, error) {
	return proto.ConsumeField(b)
}

// rawVarint reads a raw varint from a slice returned by consumeField.
func rawVarint(b []byte) (uint64, error) {
	return proto.Varint(b)
}
