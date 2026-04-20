package msgs

import (
	"errors"
	"strings"
	"testing"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	htx "github.com/0xPolygon/polygon-cli/internal/heimdall/tx"
)

// TestW4SubcommandsRegistered verifies every msg subcommand added in W4
// appears in the package registry so the mktx/send/estimate umbrellas
// pick them up.
func TestW4SubcommandsRegistered(t *testing.T) {
	want := map[string]bool{
		"checkpoint":           false,
		"checkpoint-ack":       false,
		"checkpoint-noack":     false,
		"span-propose":         false,
		"span-backfill":        false,
		"span-vote-producers":  false,
		"span-set-downtime":    false,
		"topup":                false,
		"stake-join":           false,
		"stake-update":         false,
		"signer-update":        false,
		"stake-exit":           false,
		"clerk-record":         false,
	}
	for _, n := range Names() {
		if _, ok := want[n]; ok {
			want[n] = true
		}
	}
	for name, found := range want {
		if !found {
			t.Errorf("registry missing %q", name)
		}
	}
}

// TestL1MirroringGuardsOnW4 exercises the RequireForce guard for every
// L1-mirroring Msg we registered.
func TestL1MirroringGuardsOnW4(t *testing.T) {
	shorts := []string{
		checkpointAckMsgShort,
		checkpointNoAckMsgShort,
		topupMsgShort,
		validatorJoinMsgShort,
		stakeUpdateMsgShort,
		signerUpdateMsgShort,
		validatorExitMsgShort,
		clerkRecordMsgShort,
	}
	for _, s := range shorts {
		t.Run(s, func(t *testing.T) {
			if err := htx.RequireForce(s, false); err == nil {
				t.Errorf("%s should require --force", s)
			}
			if err := htx.RequireForce(s, true); err != nil {
				t.Errorf("%s with --force should pass: %v", s, err)
			}
		})
	}
}

// TestSafeMsgsDoNotRequireForce checks the W4 msgs that are validator-
// only but not L1-mirroring.
func TestSafeMsgsDoNotRequireForce(t *testing.T) {
	shorts := []string{
		checkpointMsgShort,
		proposeSpanMsgShort,
		backfillSpansMsgShort,
		voteProducersMsgShort,
		setProducerDowntimeMsgShort,
	}
	for _, s := range shorts {
		t.Run(s, func(t *testing.T) {
			if err := htx.RequireForce(s, false); err != nil {
				t.Errorf("%s should not require --force: %v", s, err)
			}
		})
	}
}

// TestCheckpointRequiresValidatorFlag verifies the --i-am-a-validator
// friction flag is enforced.
func TestCheckpointRequiresValidatorFlag(t *testing.T) {
	srv := newTestServer(t, 25, 51129, nil)
	root, _ := newRoot(t)
	_, err := runCmd(t, root, []string{
		"--rest-url", srv.URL, "--rpc-url", srv.URL,
		"--chain-id", "heimdallv2-80002",
		"mktx", "checkpoint",
		"--private-key", fixedPrivateKeyHex,
		"--start-block", "1", "--end-block", "10",
		"--bor-chain-id", "137",
		"--root-hash", "0x" + strings.Repeat("aa", 32),
	})
	if err == nil {
		t.Fatal("expected error without --i-am-a-validator")
	}
	var uErr *client.UsageError
	if !errors.As(err, &uErr) {
		t.Fatalf("expected UsageError, got %T: %v", err, err)
	}
	if !strings.Contains(err.Error(), "validator-only") {
		t.Errorf("expected 'validator-only' in error, got %q", err.Error())
	}
}

// TestCheckpointBuildsWithValidatorFlag verifies the happy path once
// --i-am-a-validator is provided.
func TestCheckpointBuildsWithValidatorFlag(t *testing.T) {
	srv := newTestServer(t, 25, 51129, nil)
	root, _ := newRoot(t)
	stdout, err := runCmd(t, root, []string{
		"--rest-url", srv.URL, "--rpc-url", srv.URL,
		"--chain-id", "heimdallv2-80002",
		"mktx", "checkpoint",
		"--i-am-a-validator",
		"--private-key", fixedPrivateKeyHex,
		"--start-block", "1", "--end-block", "10",
		"--bor-chain-id", "137",
		"--root-hash", "0x" + strings.Repeat("aa", 32),
		"--gas", "200000", "--fee", "1000pol",
	})
	if err != nil {
		t.Fatalf("mktx checkpoint: %v\n%s", err, stdout)
	}
	if !strings.HasPrefix(strings.TrimSpace(stdout), "0x") {
		t.Fatalf("expected 0x-hex output, got %q", stdout)
	}
}

// TestSpanProposeBuildsHappyPath exercises a non-L1-mirroring msg
// end-to-end in mktx.
func TestSpanProposeBuildsHappyPath(t *testing.T) {
	srv := newTestServer(t, 25, 51129, nil)
	root, _ := newRoot(t)
	stdout, err := runCmd(t, root, []string{
		"--rest-url", srv.URL, "--rpc-url", srv.URL,
		"--chain-id", "heimdallv2-80002",
		"mktx", "span-propose",
		"--private-key", fixedPrivateKeyHex,
		"--span-id", "3",
		"--start-block", "1", "--end-block", "100",
		"--bor-chain-id", "137",
		"--seed", "0x" + strings.Repeat("bb", 32),
		"--gas", "200000", "--fee", "1000pol",
	})
	if err != nil {
		t.Fatalf("mktx span-propose: %v\n%s", err, stdout)
	}
	if !strings.HasPrefix(strings.TrimSpace(stdout), "0x") {
		t.Fatalf("expected 0x-hex output, got %q", stdout)
	}
}

// TestTopupRefusesWithoutForce asserts the L1-mirroring guard fires
// before any network call is attempted.
func TestTopupRefusesWithoutForce(t *testing.T) {
	srv := newTestServer(t, 25, 51129, nil)
	root, _ := newRoot(t)
	_, err := runCmd(t, root, []string{
		"--rest-url", srv.URL, "--rpc-url", srv.URL,
		"--chain-id", "heimdallv2-80002",
		"mktx", "topup",
		"--private-key", fixedPrivateKeyHex,
		"--user", "0x" + strings.Repeat("cc", 20),
		"--fee-amount", "1000",
		"--tx-hash", "0x" + strings.Repeat("dd", 32),
		"--log-index", "1", "--block-number", "100",
		"--gas", "200000", "--fee", "1000pol",
	})
	if err == nil {
		t.Fatal("expected error without --force")
	}
	var uErr *client.UsageError
	if !errors.As(err, &uErr) {
		t.Fatalf("expected UsageError, got %T: %v", err, err)
	}
	if !strings.Contains(err.Error(), "bridge") {
		t.Errorf("expected 'bridge' in error, got %q", err.Error())
	}
	if srv.broadcastHits.Load() != 0 {
		t.Errorf("broadcast unexpectedly called before guard: %d", srv.broadcastHits.Load())
	}
}

// TestTopupBuildsWithForce checks that --force lets the Msg through.
func TestTopupBuildsWithForce(t *testing.T) {
	srv := newTestServer(t, 25, 51129, nil)
	root, _ := newRoot(t)
	stdout, err := runCmd(t, root, []string{
		"--rest-url", srv.URL, "--rpc-url", srv.URL,
		"--chain-id", "heimdallv2-80002",
		"mktx", "topup",
		"--private-key", fixedPrivateKeyHex,
		"--user", "0x" + strings.Repeat("cc", 20),
		"--fee-amount", "1000",
		"--tx-hash", "0x" + strings.Repeat("dd", 32),
		"--log-index", "1", "--block-number", "100",
		"--gas", "200000", "--fee", "1000pol",
		"--force",
	})
	if err != nil {
		t.Fatalf("mktx topup --force: %v\n%s", err, stdout)
	}
	if !strings.HasPrefix(strings.TrimSpace(stdout), "0x") {
		t.Fatalf("expected 0x-hex output, got %q", stdout)
	}
}

// TestStakeJoinRequiresForce covers another L1-mirroring msg (one more
// check beyond the generic guard unit test to exercise wiring end-to-end).
func TestStakeJoinRequiresForce(t *testing.T) {
	srv := newTestServer(t, 25, 51129, nil)
	root, _ := newRoot(t)
	_, err := runCmd(t, root, []string{
		"--rest-url", srv.URL, "--rpc-url", srv.URL,
		"--chain-id", "heimdallv2-80002",
		"mktx", "stake-join",
		"--private-key", fixedPrivateKeyHex,
		"--val-id", "42", "--amount", "1000",
		"--signer-pub-key", "0x04" + strings.Repeat("ee", 64),
		"--tx-hash", "0x" + strings.Repeat("ff", 32),
		"--gas", "200000", "--fee", "1000pol",
	})
	if err == nil {
		t.Fatal("expected error without --force")
	}
	if !strings.Contains(err.Error(), "bridge") {
		t.Errorf("expected bridge refusal, got %v", err)
	}
}

// TestCheckpointAckRequiresL1Tx enforces the --l1-tx argument even in
// the face of --force.
func TestCheckpointAckRequiresL1Tx(t *testing.T) {
	srv := newTestServer(t, 25, 51129, nil)
	root, _ := newRoot(t)
	_, err := runCmd(t, root, []string{
		"--rest-url", srv.URL, "--rpc-url", srv.URL,
		"--chain-id", "heimdallv2-80002",
		"mktx", "checkpoint-ack",
		"--private-key", fixedPrivateKeyHex,
		"--number", "5", "--start-block", "1", "--end-block", "10",
		"--root-hash", "0x" + strings.Repeat("aa", 32),
		"--gas", "200000", "--fee", "1000pol",
		"--force",
	})
	if err == nil {
		t.Fatal("expected error when --l1-tx is missing")
	}
	if !strings.Contains(err.Error(), "l1-tx") {
		t.Errorf("expected --l1-tx in error, got %v", err)
	}
}

// TestSpanVoteProducersParsesVotes confirms the CSV votes flag parses
// and surfaces through to the Msg.
func TestSpanVoteProducersParsesVotes(t *testing.T) {
	srv := newTestServer(t, 25, 51129, nil)
	root, _ := newRoot(t)
	stdout, err := runCmd(t, root, []string{
		"--rest-url", srv.URL, "--rpc-url", srv.URL,
		"--chain-id", "heimdallv2-80002",
		"mktx", "span-vote-producers",
		"--private-key", fixedPrivateKeyHex,
		"--voter-id", "1",
		"--votes", "1,2,3",
		"--gas", "200000", "--fee", "1000pol",
	})
	if err != nil {
		t.Fatalf("mktx span-vote-producers: %v\n%s", err, stdout)
	}
}

func TestParseUint64CSV(t *testing.T) {
	got, err := parseUint64CSV("1,2,3")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(got) != 3 || got[0] != 1 || got[1] != 2 || got[2] != 3 {
		t.Errorf("got %v", got)
	}
	if _, err := parseUint64CSV("1,,3"); err == nil {
		t.Error("expected error on empty entry")
	}
	if _, err := parseUint64CSV("x"); err == nil {
		t.Error("expected error on non-numeric")
	}
}

func TestLowerEthAddress(t *testing.T) {
	cases := []struct {
		in, want string
		wantErr  bool
	}{
		{"0xABCDEF0123456789abcdef0123456789abcdef01", "0xabcdef0123456789abcdef0123456789abcdef01", false},
		{"ABCDEF0123456789abcdef0123456789abcdef01", "0xabcdef0123456789abcdef0123456789abcdef01", false},
		{"0xabc", "", true},
		{"0xzzzz0123456789abcdef0123456789abcdef0101", "", true},
	}
	for _, c := range cases {
		got, err := lowerEthAddress("addr", c.in)
		if (err != nil) != c.wantErr {
			t.Errorf("%q: err=%v wantErr=%v", c.in, err, c.wantErr)
			continue
		}
		if !c.wantErr && got != c.want {
			t.Errorf("%q: got %q want %q", c.in, got, c.want)
		}
	}
}

func TestParseHexBytes(t *testing.T) {
	b, err := parseHexBytes("f", "0xdeadbeef", 4)
	if err != nil || len(b) != 4 {
		t.Errorf("happy path: %v %v", b, err)
	}
	b, err = parseHexBytes("f", "", 4)
	if err != nil || b != nil {
		t.Errorf("empty should be nil: %v %v", b, err)
	}
	if _, err := parseHexBytes("f", "0xde", 4); err == nil {
		t.Error("length check should fail")
	}
	if _, err := parseHexBytes("f", "0xxx", 0); err == nil {
		t.Error("invalid hex should fail")
	}
}
