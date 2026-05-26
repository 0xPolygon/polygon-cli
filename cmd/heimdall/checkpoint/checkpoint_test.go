package checkpoint

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
)

// --- normalizeCheckpointHash ---

func TestNormalizeCheckpointHash(t *testing.T) {
	const raw = "94297F18F736A0C018E4871A5257384450673AC8441F8F7956523231D74D2A29"
	cases := []struct {
		in      string
		want    string
		wantErr bool
	}{
		{raw, strings.ToLower(raw), false},
		{"0x" + raw, strings.ToLower(raw), false},
		{"0X" + raw, strings.ToLower(raw), false},
		{strings.ToLower(raw), strings.ToLower(raw), false},
		{"", "", true},
		{"0x12", "", true},
		{"zz" + raw[2:], "", true},
	}
	for _, c := range cases {
		got, err := normalizeCheckpointHash(c.in)
		if (err != nil) != c.wantErr {
			t.Errorf("in=%q err=%v wantErr=%v", c.in, err, c.wantErr)
			continue
		}
		if !c.wantErr && got != c.want {
			t.Errorf("in=%q got=%q want=%q", c.in, got, c.want)
		}
	}
}

// --- params ---

func TestParamsHumanOutput(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/checkpoints/params": {body: loadFixture(t, "rest", "checkpoints_params.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "params")
	if err != nil {
		t.Fatalf("params: %v", err)
	}
	for _, want := range []string{"avg_checkpoint_length", "max_checkpoint_length", "checkpoint_buffer_time"} {
		mustContain(t, stdout, want)
	}
	mustContain(t, stdout, "256")
	mustContain(t, stdout, "1500s")
}

func TestParamsJSONPassthrough(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/checkpoints/params": {body: loadFixture(t, "rest", "checkpoints_params.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "params", "--json")
	if err != nil {
		t.Fatalf("params --json: %v", err)
	}
	var v any
	if jerr := json.Unmarshal([]byte(stdout), &v); jerr != nil {
		t.Fatalf("output not valid JSON: %v\n%s", jerr, stdout)
	}
}

// --- count ---

func TestCountBareInteger(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/checkpoints/count": {body: loadFixture(t, "rest", "checkpoints_count.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "count")
	if err != nil {
		t.Fatalf("count: %v", err)
	}
	if strings.TrimSpace(stdout) != "38871" {
		t.Errorf("count stdout = %q, want 38871", stdout)
	}
}

// --- latest / get ---

func TestLatestUnwrapsEnvelope(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/checkpoints/latest": {body: loadFixture(t, "rest", "checkpoints_latest.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "latest")
	if err != nil {
		t.Fatalf("latest: %v", err)
	}
	// Fields from the inner checkpoint object should appear without
	// the "checkpoint." prefix, and root_hash should be hex.
	mustContain(t, stdout, "id")
	mustContain(t, stdout, "root_hash")
	mustContain(t, stdout, "0x")
	// Base64 root_hash from fixture should not leak through.
	mustNotContain(t, stdout, "NRfvvV9YAjjav+cR70om6WDob+IIZjPVyIYMAcrzxy4=")
}

func TestLatestRawPreservesBase64(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/checkpoints/latest": {body: loadFixture(t, "rest", "checkpoints_latest.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "latest", "--raw")
	if err != nil {
		t.Fatalf("latest --raw: %v", err)
	}
	mustContain(t, stdout, "NRfvvV9YAjjav+cR70om6WDob+IIZjPVyIYMAcrzxy4=")
}

func TestGetByExplicitSubcommand(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/checkpoints/38871": {body: loadFixture(t, "rest", "checkpoints_by_id.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "get", "38871")
	if err != nil {
		t.Fatalf("get 38871: %v", err)
	}
	mustContain(t, stdout, "38871")
	mustContain(t, stdout, "0x")
}

func TestGetBareIntegerShortcut(t *testing.T) {
	// `checkpoint 38871` (no explicit `get`) should route to runGet.
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/checkpoints/38871": {body: loadFixture(t, "rest", "checkpoints_by_id.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "38871")
	if err != nil {
		t.Fatalf("checkpoint 38871: %v", err)
	}
	mustContain(t, stdout, "38871")
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

// --- buffer ---

func TestBufferPopulated(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/checkpoints/buffer": {body: loadFixture(t, "rest", "checkpoints_buffer.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "buffer")
	if err != nil {
		t.Fatalf("buffer: %v", err)
	}
	mustContain(t, stdout, "proposer")
	mustContain(t, stdout, "0x4ad84f7014b7b44f723f284a85b1662337971439")
	mustNotContain(t, stdout, "empty")
}

func TestBufferEmptyPrintsEmpty(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/checkpoints/buffer": {body: loadFixture(t, "rest", "checkpoints_buffer_empty.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "buffer")
	if err != nil {
		t.Fatalf("buffer (empty): %v", err)
	}
	mustContain(t, stdout, "empty")
	// The buffer-empty hint should accompany it.
	mustContain(t, stdout, "no checkpoint in flight")
}

func TestBufferJSONPassesZerosThrough(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/checkpoints/buffer": {body: loadFixture(t, "rest", "checkpoints_buffer_empty.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "buffer", "--json")
	if err != nil {
		t.Fatalf("buffer --json: %v", err)
	}
	var v any
	if jerr := json.Unmarshal([]byte(stdout), &v); jerr != nil {
		t.Fatalf("output not valid JSON: %v\n%s", jerr, stdout)
	}
	// --json should not produce the `empty` human-readable line.
	mustNotContain(t, stdout, "\nempty\n")
}

// --- last-no-ack ---

func TestLastNoAckAnnotatesAge(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/checkpoints/last-no-ack": {body: loadFixture(t, "rest", "checkpoints_last_no_ack.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "last-no-ack")
	if err != nil {
		t.Fatalf("last-no-ack: %v", err)
	}
	// Unix seconds (1776695056 == 2026-04-20 UTC) and the annotated form.
	mustContain(t, stdout, "1776695056")
	mustContain(t, stdout, "UTC")
	// The annotator always suffixes with "ago" or "from now".
	if !(strings.Contains(stdout, "ago") || strings.Contains(stdout, "from now")) {
		t.Errorf("expected age suffix in output: %q", stdout)
	}
}

// --- next ---

func TestNextSuccess(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/checkpoints/prepare-next": {body: loadFixture(t, "rest", "checkpoints_prepare_next.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "next")
	if err != nil {
		t.Fatalf("next: %v", err)
	}
	mustContain(t, stdout, "proposer")
	// account_root_hash and root_hash normalized to hex.
	mustContain(t, stdout, "0x")
}

func TestNextL1NotConfiguredEmitsHint(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/checkpoints/prepare-next": {
			status: 500,
			body:   loadFixture(t, "rest", "checkpoints_prepare_next_l1_unconfigured.json"),
		},
	})
	_, stderr, err := runCmd(t, srv.URL, "next")
	if err == nil {
		t.Fatal("expected error from prepare-next HTTP 500")
	}
	// The L1-not-configured hint should show up on stderr.
	mustContain(t, stderr, "eth_rpc_url")
}

// --- list ---

func TestListDefaultsAndTable(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/checkpoints/list": {body: loadFixture(t, "rest", "checkpoints_list.json")},
	})
	stdout, stderr, err := runCmd(t, srv.URL, "list")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	// Column headers from the union of fields.
	mustContain(t, stdout, "id")
	mustContain(t, stdout, "root_hash")
	mustContain(t, stdout, "proposer")
	// At least one row value.
	mustContain(t, stdout, "38871")
	// next_key surfaces on stderr.
	mustContain(t, stderr, "next_key=")
}

func TestListJSONPassthrough(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/checkpoints/list": {body: loadFixture(t, "rest", "checkpoints_list.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "list", "--json")
	if err != nil {
		t.Fatalf("list --json: %v", err)
	}
	var v any
	if jerr := json.Unmarshal([]byte(stdout), &v); jerr != nil {
		t.Fatalf("output not valid JSON: %v\n%s", jerr, stdout)
	}
}

// --- signatures ---

func TestSignaturesPopulated(t *testing.T) {
	const txHash = "94297f18f736a0c018e4871a5257384450673ac8441f8f7956523231d74d2a29"
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/checkpoints/signatures/" + txHash: {body: loadFixture(t, "rest", "checkpoints_signatures.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "signatures", "0x"+txHash)
	if err != nil {
		t.Fatalf("signatures: %v", err)
	}
	var v any
	if jerr := json.Unmarshal([]byte(stdout), &v); jerr != nil {
		t.Fatalf("signatures output not JSON: %v\n%s", jerr, stdout)
	}
	mustContain(t, stdout, "signer")
	// root_hash-like fields don't appear; signatures are normalized to hex.
	mustContain(t, stdout, "0x")
}

func TestSignaturesToleratesMissingPrefix(t *testing.T) {
	const txHash = "94297f18f736a0c018e4871a5257384450673ac8441f8f7956523231d74d2a29"
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/checkpoints/signatures/" + txHash: {body: loadFixture(t, "rest", "checkpoints_signatures.json")},
	})
	_, _, err := runCmd(t, srv.URL, "signatures", txHash)
	if err != nil {
		t.Fatalf("signatures (no 0x): %v", err)
	}
}

func TestSignaturesRejectsBadHash(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{})
	_, _, err := runCmd(t, srv.URL, "signatures", "0x1234")
	if err == nil {
		t.Fatal("expected error for short hash")
	}
	var uErr *client.UsageError
	if !errors.As(err, &uErr) {
		t.Fatalf("got %T, want *UsageError", err)
	}
}

// --- overview ---

func TestOverviewJSON(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/checkpoints/overview": {body: loadFixture(t, "rest", "checkpoints_overview.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "overview")
	if err != nil {
		t.Fatalf("overview: %v", err)
	}
	var v any
	if jerr := json.Unmarshal([]byte(stdout), &v); jerr != nil {
		t.Fatalf("overview output not JSON: %v\n%s", jerr, stdout)
	}
	mustContain(t, stdout, "ack_count")
	mustContain(t, stdout, "validator_set")
}

// --- isBufferEmpty unit test ---

func TestIsBufferEmpty(t *testing.T) {
	cases := []struct {
		name string
		in   map[string]any
		want bool
	}{
		{
			name: "zero-address proposer is empty",
			in: map[string]any{
				"checkpoint": map[string]any{
					"proposer": "0x0000000000000000000000000000000000000000",
				},
			},
			want: true,
		},
		{
			name: "empty-string proposer is empty",
			in: map[string]any{
				"checkpoint": map[string]any{
					"proposer": "",
				},
			},
			want: true,
		},
		{
			name: "non-zero proposer is not empty",
			in: map[string]any{
				"checkpoint": map[string]any{
					"proposer": "0x4ad84f7014b7b44f723f284a85b1662337971439",
				},
			},
			want: false,
		},
		{
			name: "missing checkpoint wrapper is not empty",
			in:   map[string]any{"foo": "bar"},
			want: false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := isBufferEmpty(c.in); got != c.want {
				t.Errorf("got=%v want=%v", got, c.want)
			}
		})
	}
}

// --- isL1Unreachable unit test ---

func TestIsL1Unreachable(t *testing.T) {
	code13 := []byte(`{"code":13,"message":"dial tcp"}`)
	other := []byte(`{"code":5,"message":"not found"}`)
	notJSON := []byte(`<html>502 Bad Gateway</html>`)

	if !isL1Unreachable(code13, nil) {
		t.Error("expected true for code 13")
	}
	if isL1Unreachable(other, nil) {
		t.Error("expected false for non-13 code")
	}
	if isL1Unreachable(notJSON, nil) {
		t.Error("expected false for non-JSON body")
	}
	// Wrapped HTTPError with the body on the error itself still works.
	hErr := &client.HTTPError{StatusCode: 500, Body: code13}
	if !isL1Unreachable(nil, hErr) {
		t.Error("expected true when body comes from HTTPError")
	}
}
