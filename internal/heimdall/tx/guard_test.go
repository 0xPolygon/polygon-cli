package tx

import (
	"strings"
	"testing"
)

func TestRequireForceFlagsL1MirroringTypes(t *testing.T) {
	cases := []struct {
		name string
		typ  string
	}{
		{"topup short", "MsgTopupTx"},
		{"topup full url", "/heimdallv2.topup.MsgTopupTx"},
		{"cpack short", "MsgCpAck"},
		{"checkpoint ack", "MsgCheckpointAck"},
		{"cpnoack", "MsgCpNoAck"},
		{"clerk record", "MsgEventRecordRequest"},
		{"clerk short", "MsgClerkRecord"},
		{"validator join", "MsgValidatorJoin"},
		{"stake join", "MsgStakeJoin"},
		{"signer update", "MsgSignerUpdate"},
		{"stake update", "MsgStakeUpdate"},
		{"stake exit", "MsgStakeExit"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := RequireForce(tc.typ, false)
			if err == nil {
				t.Fatalf("expected error for %q", tc.typ)
			}
			if !strings.Contains(err.Error(), "bridge") {
				t.Errorf("expected 'bridge' in error, got %q", err.Error())
			}
			if !strings.Contains(err.Error(), "--force") {
				t.Errorf("expected '--force' in error, got %q", err.Error())
			}
			// --force allows through.
			if err := RequireForce(tc.typ, true); err != nil {
				t.Errorf("with force=true: %v", err)
			}
		})
	}
}

func TestRequireForceAllowsSafeTypes(t *testing.T) {
	cases := []string{
		"MsgWithdrawFeeTx",
		"/heimdallv2.topup.MsgWithdrawFeeTx",
		"MsgProposeSpan",
		"MsgBackfillSpan",
		"MsgVoteProducers",
		"MsgCheckpoint",
		"",
		"RandomOtherMsg",
	}
	for _, typ := range cases {
		t.Run(typ, func(t *testing.T) {
			if err := RequireForce(typ, false); err != nil {
				t.Errorf("expected nil for safe type %q, got %v", typ, err)
			}
		})
	}
}

func TestShortMsgName(t *testing.T) {
	cases := map[string]string{
		"MsgTopupTx":                       "MsgTopupTx",
		"/heimdallv2.topup.MsgTopupTx":     "MsgTopupTx",
		"heimdallv2.topup.MsgTopupTx":      "MsgTopupTx",
		"":                                  "",
		"no/path":                           "path",
	}
	for in, want := range cases {
		if got := shortMsgName(in); got != want {
			t.Errorf("shortMsgName(%q) = %q, want %q", in, got, want)
		}
	}
}
