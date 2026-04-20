package validator

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
)

// --- normalizeSignerAddress / normalizeTxHash ---

func TestNormalizeSignerAddress(t *testing.T) {
	const raw = "4AD84F7014B7B44F723F284A85B1662337971439"
	cases := []struct {
		in      string
		want    string
		wantErr bool
	}{
		{raw, "0x" + strings.ToLower(raw), false},
		{"0x" + raw, "0x" + strings.ToLower(raw), false},
		{"0X" + raw, "0x" + strings.ToLower(raw), false},
		{strings.ToLower(raw), "0x" + strings.ToLower(raw), false},
		{"", "", true},
		{"0x12", "", true},
		{"zz" + raw[2:], "", true},
	}
	for _, c := range cases {
		got, err := normalizeSignerAddress(c.in)
		if (err != nil) != c.wantErr {
			t.Errorf("in=%q err=%v wantErr=%v", c.in, err, c.wantErr)
			continue
		}
		if !c.wantErr && got != c.want {
			t.Errorf("in=%q got=%q want=%q", c.in, got, c.want)
		}
	}
}

func TestNormalizeTxHash(t *testing.T) {
	const raw = "94297F18F736A0C018E4871A5257384450673AC8441F8F7956523231D74D2A29"
	cases := []struct {
		in      string
		want    string
		wantErr bool
	}{
		{raw, "0x" + strings.ToLower(raw), false},
		{"0x" + raw, "0x" + strings.ToLower(raw), false},
		{"", "", true},
		{"0x12", "", true},
		{"zz" + raw[2:], "", true},
	}
	for _, c := range cases {
		got, err := normalizeTxHash(c.in)
		if (err != nil) != c.wantErr {
			t.Errorf("in=%q err=%v wantErr=%v", c.in, err, c.wantErr)
			continue
		}
		if !c.wantErr && got != c.want {
			t.Errorf("in=%q got=%q want=%q", c.in, got, c.want)
		}
	}
}

// --- total-power ---

func TestTotalPowerBareInteger(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/stake/total-power": {body: loadFixture(t, "rest", "stake_total_power.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "total-power")
	if err != nil {
		t.Fatalf("total-power: %v", err)
	}
	if strings.TrimSpace(stdout) != "632197800" {
		t.Errorf("total-power stdout = %q, want 632197800", stdout)
	}
}

func TestTotalPowerJSON(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/stake/total-power": {body: loadFixture(t, "rest", "stake_total_power.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "total-power", "--json")
	if err != nil {
		t.Fatalf("total-power --json: %v", err)
	}
	var m map[string]any
	if jerr := json.Unmarshal([]byte(stdout), &m); jerr != nil {
		t.Fatalf("total-power --json not JSON: %v\n%s", jerr, stdout)
	}
	if got := m["total_power"]; got != "632197800" {
		t.Errorf("total_power=%v, want \"632197800\"", got)
	}
}

// --- get / bare integer ---

func TestGetByExplicitSubcommand(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/stake/validator/16": {body: loadFixture(t, "rest", "stake_validator_by_id.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "get", "16")
	if err != nil {
		t.Fatalf("get 16: %v", err)
	}
	mustContain(t, stdout, "val_id")
	mustContain(t, stdout, "16")
	mustContain(t, stdout, "0x02f615e95563ef16f10354dba9e584e58d2d4314")
}

func TestGetBareIntegerShortcut(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/stake/validator/16": {body: loadFixture(t, "rest", "stake_validator_by_id.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "16")
	if err != nil {
		t.Fatalf("validator 16: %v", err)
	}
	mustContain(t, stdout, "16")
}

func TestGetInvalidIDIsUsageError(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{})
	_, _, err := runCmd(t, srv.URL, "get", "notanumber")
	if err == nil {
		t.Fatal("expected error for non-integer id")
	}
	var uErr *client.UsageError
	if !errors.As(err, &uErr) {
		t.Fatalf("got %T, want *UsageError", err)
	}
}

func TestUmbrellaUnknownArgIsUsageError(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{})
	_, _, err := runCmd(t, srv.URL, "banana")
	if err == nil {
		t.Fatal("expected usage error for unknown umbrella arg")
	}
	var uErr *client.UsageError
	if !errors.As(err, &uErr) {
		t.Fatalf("got %T, want *UsageError", err)
	}
}

// --- signer ---

func TestSignerToleratesMissingPrefix(t *testing.T) {
	const addr = "0x02f615e95563ef16f10354dba9e584e58d2d4314"
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/stake/signer/" + addr: {body: loadFixture(t, "rest", "stake_signer.json")},
	})
	// Without 0x:
	stdout, _, err := runCmd(t, srv.URL, "signer", "02f615e95563ef16f10354dba9e584e58d2d4314")
	if err != nil {
		t.Fatalf("signer (no 0x): %v", err)
	}
	mustContain(t, stdout, "val_id")
	mustContain(t, stdout, "16")
}

func TestSignerRejectsBadHex(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{})
	_, _, err := runCmd(t, srv.URL, "signer", "0xshort")
	if err == nil {
		t.Fatal("expected error for short signer")
	}
	var uErr *client.UsageError
	if !errors.As(err, &uErr) {
		t.Fatalf("got %T, want *UsageError", err)
	}
}

// --- status / is_old -> is_current rename ---

// TestStatusRenamesIsOldToIsCurrent is the flagship test for this
// subcommand: the upstream name `is_old` must never appear in default
// KV output; the renamed `is_current` must.
func TestStatusRenamesIsOldToIsCurrent(t *testing.T) {
	const addr = "0x4ad84f7014b7b44f723f284a85b1662337971439"
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/stake/validator-status/" + addr: {body: loadFixture(t, "rest", "stake_validator_status.json")},
	})
	stdout, stderr, err := runCmd(t, srv.URL, "status", addr)
	if err != nil {
		t.Fatalf("status: %v", err)
	}
	mustContain(t, stdout, "is_current")
	mustNotContain(t, stdout, "is_old")
	// The rename hint goes on stderr.
	mustContain(t, stderr, "is_current")
	mustContain(t, stderr, "is_old")
	mustContain(t, stderr, "renamed")
}

// TestStatusJSONRenamesIsOldToIsCurrent asserts the same rename in the
// --json output.
func TestStatusJSONRenamesIsOldToIsCurrent(t *testing.T) {
	const addr = "0x4ad84f7014b7b44f723f284a85b1662337971439"
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/stake/validator-status/" + addr: {body: loadFixture(t, "rest", "stake_validator_status.json")},
	})
	stdout, stderr, err := runCmd(t, srv.URL, "status", addr, "--json")
	if err != nil {
		t.Fatalf("status --json: %v", err)
	}
	var m map[string]any
	if jerr := json.Unmarshal([]byte(stdout), &m); jerr != nil {
		t.Fatalf("status --json not JSON: %v\n%s", jerr, stdout)
	}
	if _, ok := m["is_old"]; ok {
		t.Errorf("is_old must not appear in JSON output: %v", m)
	}
	v, ok := m["is_current"]
	if !ok {
		t.Fatalf("is_current missing from JSON output: %v", m)
	}
	if b, ok := v.(bool); !ok || !b {
		t.Errorf("is_current should be true, got %v", v)
	}
	// Hint still emitted for scripters who pipe JSON through jq.
	mustContain(t, stderr, "renamed")
}

// TestRenameIsOldToIsCurrentUnit is a direct unit test on the helper so
// we catch changes to the rename semantics even when the command plumbing
// is bypassed.
func TestRenameIsOldToIsCurrentUnit(t *testing.T) {
	m := map[string]any{"is_old": true}
	if !renameIsOldToIsCurrent(m) {
		t.Fatal("expected rename to return true")
	}
	if _, stillThere := m["is_old"]; stillThere {
		t.Error("is_old should be removed")
	}
	if v, ok := m["is_current"].(bool); !ok || !v {
		t.Errorf("expected is_current=true, got %v", m["is_current"])
	}

	// Idempotent on already-renamed or empty input.
	empty := map[string]any{}
	if renameIsOldToIsCurrent(empty) {
		t.Error("empty map should not trigger a rename")
	}
}

// --- proposer / proposers ---

func TestProposerKV(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/stake/proposers/current": {body: loadFixture(t, "rest", "stake_proposers_current.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "proposer")
	if err != nil {
		t.Fatalf("proposer: %v", err)
	}
	mustContain(t, stdout, "val_id")
	mustContain(t, stdout, "0x4ad84f7014b7b44f723f284a85b1662337971439")
}

func TestProposersDefaultsToOne(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/stake/proposers/1": {body: []byte(`{"proposers":[]}`)},
	})
	// `proposers` with no N must call /stake/proposers/1
	stdout, _, err := runCmd(t, srv.URL, "proposers")
	if err != nil {
		t.Fatalf("proposers (default): %v", err)
	}
	mustContain(t, stdout, "(no proposers)")
}

func TestProposersExplicitN(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/stake/proposers/5": {body: loadFixture(t, "rest", "stake_proposers_n.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "proposers", "5")
	if err != nil {
		t.Fatalf("proposers 5: %v", err)
	}
	// Table header columns from the union of validator fields.
	mustContain(t, stdout, "val_id")
	mustContain(t, stdout, "signer")
}

func TestProposersRejectsZero(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{})
	_, _, err := runCmd(t, srv.URL, "proposers", "0")
	if err == nil {
		t.Fatal("expected error for N=0")
	}
	var uErr *client.UsageError
	if !errors.As(err, &uErr) {
		t.Fatalf("got %T, want *UsageError", err)
	}
}

// --- set / validators alias ---

// TestSetSortedByPowerDefault asserts the default power-desc sort. The
// fixture's top validator by voting_power is val_id 5
// (voting_power 80000015) even though it is not first in the JSON
// array.
func TestSetSortedByPowerDefault(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/stake/validators-set": {body: loadFixture(t, "rest", "stake_validators_set.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "set", "--limit", "1")
	if err != nil {
		t.Fatalf("set: %v", err)
	}
	// First row's val_id must be the highest-power one (val_id 5).
	rows := strings.Split(strings.TrimRight(stdout, "\n"), "\n")
	if len(rows) < 2 {
		t.Fatalf("expected at least header+1 row, got %d rows:\n%s", len(rows), stdout)
	}
	first := rows[1]
	if !strings.Contains(first, "80000015") {
		t.Errorf("expected top-power validator (voting_power 80000015) in first row, got: %q", first)
	}
	if !strings.Contains(first, "0x6dc2dd54f24979ec26212794c71afefed722280c") {
		t.Errorf("expected signer 0x6dc2dd… in first row, got: %q", first)
	}
}

// TestSetSortedByID asserts --sort id uses ascending val_id order.
func TestSetSortedByID(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/stake/validators-set": {body: loadFixture(t, "rest", "stake_validators_set.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "set", "--sort", "id", "--limit", "3")
	if err != nil {
		t.Fatalf("set --sort id: %v", err)
	}
	// Minimum val_ids in fixture sorted ascending: 1, 4, 5.
	rows := strings.Split(strings.TrimRight(stdout, "\n"), "\n")
	if len(rows) < 4 {
		t.Fatalf("expected header+3 rows, got %d:\n%s", len(rows), stdout)
	}
	// val_id column is the second-to-last; asserting substrings rather
	// than exact column parsing to stay robust against column re-order.
	row1 := rows[1]
	row2 := rows[2]
	row3 := rows[3]
	if !strings.Contains(row1, "6ab3d36c46ecfb9b9c0bd51cb1c3da5a2c81cea6") {
		// val_id 1's signer.
		t.Errorf("expected val_id 1 first, got: %q", row1)
	}
	if !strings.Contains(row2, "4ad84f7014b7b44f723f284a85b1662337971439") {
		t.Errorf("expected val_id 4 second, got: %q", row2)
	}
	if !strings.Contains(row3, "6dc2dd54f24979ec26212794c71afefed722280c") {
		t.Errorf("expected val_id 5 third, got: %q", row3)
	}
}

// TestSetSortedBySigner asserts --sort signer uses ascending address
// order.
func TestSetSortedBySigner(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/stake/validators-set": {body: loadFixture(t, "rest", "stake_validators_set.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "set", "--sort", "signer", "--limit", "1")
	if err != nil {
		t.Fatalf("set --sort signer: %v", err)
	}
	rows := strings.Split(strings.TrimRight(stdout, "\n"), "\n")
	if len(rows) < 2 {
		t.Fatalf("expected header+row, got %d:\n%s", len(rows), stdout)
	}
	// Lowest signer alphabetically in the fixture is 0x02f615e95… (val_id 16).
	if !strings.Contains(rows[1], "0x02f615e95563ef16f10354dba9e584e58d2d4314") {
		t.Errorf("expected lowest-signer first, got: %q", rows[1])
	}
}

// TestSetLimitTruncates asserts --limit truncates to the first N rows
// (header excluded).
func TestSetLimitTruncates(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/stake/validators-set": {body: loadFixture(t, "rest", "stake_validators_set.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "set", "--limit", "2")
	if err != nil {
		t.Fatalf("set --limit: %v", err)
	}
	rows := strings.Split(strings.TrimRight(stdout, "\n"), "\n")
	// header + 2 data rows
	if len(rows) != 3 {
		t.Errorf("expected 3 rows (header+2), got %d:\n%s", len(rows), stdout)
	}
}

func TestSetUnknownSortIsUsageError(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{})
	_, _, err := runCmd(t, srv.URL, "set", "--sort", "banana")
	if err == nil {
		t.Fatal("expected usage error")
	}
	var uErr *client.UsageError
	if !errors.As(err, &uErr) {
		t.Fatalf("got %T, want *UsageError", err)
	}
}

// TestValidatorsTopLevelAlias asserts the top-level `validators`
// command is identical in behaviour to `validator set`.
func TestValidatorsTopLevelAlias(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/stake/validators-set": {body: loadFixture(t, "rest", "stake_validators_set.json")},
	})
	stdout, _, err := runCmdNamed(t, srv.URL, "validators", "--limit", "1", "--sort", "id")
	if err != nil {
		t.Fatalf("validators alias: %v", err)
	}
	// val_id 1's signer should be the only data row.
	mustContain(t, stdout, "6ab3d36c46ecfb9b9c0bd51cb1c3da5a2c81cea6")
}

// --- is-old-stake-tx ---

// TestIsOldStakeTxL1UnconfiguredEmitsHint: the HTTP 500 gRPC-code-13
// envelope must produce the L1-not-configured hint on stderr while
// still propagating the error.
func TestIsOldStakeTxL1UnconfiguredEmitsHint(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/stake/is-old-tx": {
			status: 500,
			body:   loadFixture(t, "rest", "stake_is_old_tx_l1_unconfigured.json"),
			wantQuery: map[string]string{
				"tx_hash":   "0x94297f18f736a0c018e4871a5257384450673ac8441f8f7956523231d74d2a29",
				"log_index": "0",
			},
		},
	})
	_, stderr, err := runCmd(t, srv.URL,
		"is-old-stake-tx",
		"0x94297f18f736a0c018e4871a5257384450673ac8441f8f7956523231d74d2a29",
		"0",
	)
	if err == nil {
		t.Fatal("expected error from HTTP 500")
	}
	mustContain(t, stderr, "eth_rpc_url")
	// The error surface should be a node error (exit code 1), since the
	// HTTP 500 still came back from Heimdall itself.
	if code := client.ExitCode(err); code == 0 || code == 3 {
		t.Errorf("unexpected exit code %d for node error", code)
	}
}

func TestIsOldStakeTxBadLogIndexIsUsageError(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{})
	_, _, err := runCmd(t, srv.URL,
		"is-old-stake-tx",
		"0x94297f18f736a0c018e4871a5257384450673ac8441f8f7956523231d74d2a29",
		"notanumber",
	)
	if err == nil {
		t.Fatal("expected usage error")
	}
	var uErr *client.UsageError
	if !errors.As(err, &uErr) {
		t.Fatalf("got %T, want *UsageError", err)
	}
}

// TestIsL1UnreachableUnit mirrors checkpoint's unit test but exercises
// the extra fallback on err-string inspection (for cases where the
// transport returns an error with no HTTP body at all).
func TestIsL1UnreachableUnit(t *testing.T) {
	code13 := []byte(`{"code":13,"message":"dial tcp"}`)
	other := []byte(`{"code":5,"message":"not found"}`)
	notJSON := []byte(`<html>502 Bad Gateway</html>`)

	if !isL1Unreachable(code13, nil) {
		t.Error("expected true for code 13 body")
	}
	if isL1Unreachable(other, nil) {
		t.Error("expected false for non-13 code")
	}
	if isL1Unreachable(notJSON, nil) {
		t.Error("expected false for non-JSON body")
	}
	hErr := &client.HTTPError{StatusCode: 500, Body: code13}
	if !isL1Unreachable(nil, hErr) {
		t.Error("expected true when body comes from HTTPError")
	}
	// Transport-level error string without body.
	connErr := errors.New("dial tcp 172.19.0.2:1317: connect: connection refused")
	if !isL1Unreachable(nil, connErr) {
		t.Error("expected true for connection-refused transport error")
	}
}
