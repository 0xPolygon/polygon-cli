package proto

import (
	"fmt"
	"sort"
	"sync"
)

// Decoder takes an Any.value byte slice and returns a decoded Go value
// (typically a typed *Msg or a map[string]any) along with an error.
//
// The registry is populated by package init() functions in this file.
// Additional Msg types added in the future should register their
// decoder alongside their Marshal/Unmarshal helpers.
type Decoder func([]byte) (any, error)

var (
	registryMu     sync.RWMutex
	typeURLDecoder = map[string]Decoder{}
)

// Register associates typeURL with a decoder. Panics on a duplicate
// registration so the error is caught at init time. Empty typeURL or
// nil decoder also panic.
func Register(typeURL string, decoder Decoder) {
	if typeURL == "" {
		panic("proto.Register: typeURL is empty")
	}
	if decoder == nil {
		panic("proto.Register: decoder is nil")
	}
	registryMu.Lock()
	defer registryMu.Unlock()
	if _, ok := typeURLDecoder[typeURL]; ok {
		panic("proto.Register: duplicate type URL " + typeURL)
	}
	typeURLDecoder[typeURL] = decoder
}

// Lookup returns the decoder registered for typeURL, or (nil, false)
// if none is registered.
func Lookup(typeURL string) (Decoder, bool) {
	registryMu.RLock()
	defer registryMu.RUnlock()
	d, ok := typeURLDecoder[typeURL]
	return d, ok
}

// Decode resolves typeURL in the registry and invokes its decoder. If
// typeURL is unknown, Decode returns an error including the type URL
// so callers can forward it to the user verbatim.
func Decode(typeURL string, value []byte) (any, error) {
	d, ok := Lookup(typeURL)
	if !ok {
		return nil, fmt.Errorf("unknown type URL %q", typeURL)
	}
	return d(value)
}

// KnownTypeURLs returns a sorted list of every registered type URL.
// Useful for diagnostics and for `decode msg` help text.
func KnownTypeURLs() []string {
	registryMu.RLock()
	defer registryMu.RUnlock()
	out := make([]string, 0, len(typeURLDecoder))
	for k := range typeURLDecoder {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func init() {
	// Topup.
	Register(MsgWithdrawFeeTxTypeURL, func(b []byte) (any, error) { return UnmarshalMsgWithdrawFeeTx(b) })
	Register(MsgTopupTxTypeURL, func(b []byte) (any, error) { return UnmarshalMsgTopupTx(b) })

	// Checkpoint.
	Register(MsgCheckpointTypeURL, func(b []byte) (any, error) { return UnmarshalMsgCheckpoint(b) })
	Register(MsgCpAckTypeURL, func(b []byte) (any, error) { return UnmarshalMsgCpAck(b) })
	Register(MsgCpNoAckTypeURL, func(b []byte) (any, error) { return UnmarshalMsgCpNoAck(b) })

	// Bor.
	Register(MsgProposeSpanTypeURL, func(b []byte) (any, error) { return UnmarshalMsgProposeSpan(b) })
	Register(MsgBackfillSpansTypeURL, func(b []byte) (any, error) { return UnmarshalMsgBackfillSpans(b) })
	Register(MsgVoteProducersTypeURL, func(b []byte) (any, error) { return UnmarshalMsgVoteProducers(b) })
	Register(MsgSetProducerDowntimeTypeURL, func(b []byte) (any, error) { return UnmarshalMsgSetProducerDowntime(b) })

	// Stake.
	Register(MsgValidatorJoinTypeURL, func(b []byte) (any, error) { return UnmarshalMsgValidatorJoin(b) })
	Register(MsgStakeUpdateTypeURL, func(b []byte) (any, error) { return UnmarshalMsgStakeUpdate(b) })
	Register(MsgSignerUpdateTypeURL, func(b []byte) (any, error) { return UnmarshalMsgSignerUpdate(b) })
	Register(MsgValidatorExitTypeURL, func(b []byte) (any, error) { return UnmarshalMsgValidatorExit(b) })

	// Clerk.
	Register(MsgEventRecordTypeURL, func(b []byte) (any, error) { return UnmarshalMsgEventRecord(b) })
}
