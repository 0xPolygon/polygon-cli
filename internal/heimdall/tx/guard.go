package tx

import "fmt"

// L1MirroringMsgTypes is the set of Msg types that are produced by the
// Heimdall bridge after observing an L1 event. Operators almost never
// hand-build these, and submitting one manually is either a no-op or a
// replay that competes with the real bridge path.
//
// Subcommands flagged for these types call RequireForce before
// building; it returns a user-facing refusal error unless the caller
// has set --force.
//
// Wording is taken verbatim from HEIMDALLCAST_REQUIREMENTS.md §3.3.
var L1MirroringMsgTypes = map[string]struct{}{
	"MsgTopupTx":           {},
	"MsgCheckpointAck":     {},
	"MsgCpAck":             {},
	"MsgCheckpointNoAck":   {},
	"MsgCpNoAck":           {},
	"MsgEventRecordRequest": {},
	"MsgClerkRecord":       {},
	"MsgValidatorJoin":     {},
	"MsgStakeJoin":         {},
	"MsgValidatorUpdate":   {},
	"MsgStakeUpdate":       {},
	"MsgSignerUpdate":      {},
	"MsgValidatorExit":     {},
	"MsgStakeExit":         {},
}

// RequireForce returns an error if msgType is an L1-mirroring type and
// force is false. Returns nil for safe types or when force is true.
//
// The error text matches HEIMDALLCAST_REQUIREMENTS.md §3.3:
//
//	"this message is produced by the bridge after observing an L1
//	 event; you almost certainly do not want to build one by hand.
//	 Re-run with --force to bypass."
//
// msgType is matched against the short name only (the last segment of
// the type URL, e.g. "MsgTopupTx" rather than the full URL) to keep
// call-sites short. Both short and fully-qualified names are accepted.
func RequireForce(msgType string, force bool) error {
	if force {
		return nil
	}
	short := shortMsgName(msgType)
	if _, ok := L1MirroringMsgTypes[short]; !ok {
		return nil
	}
	return fmt.Errorf(
		"%s is produced by the bridge after observing an L1 event; you almost certainly do not want to build one by hand. Re-run with --force to bypass",
		short,
	)
}

// shortMsgName returns the last segment of a cosmos-sdk Any type URL
// so callers can pass either "/heimdallv2.topup.MsgTopupTx" or
// "MsgTopupTx". Empty input is returned unchanged.
func shortMsgName(t string) string {
	for i := len(t) - 1; i >= 0; i-- {
		c := t[i]
		if c == '.' || c == '/' {
			return t[i+1:]
		}
	}
	return t
}
