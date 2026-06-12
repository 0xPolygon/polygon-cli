package topup

import (
	"encoding/json"
	"errors"
	"net/url"
	"strings"
	"testing"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
)

// --- root ---

// TestRootDefault0xHex asserts that the default (human) output of
// `topup root` converts the base64 `account_root_hash` to 0x-hex so
// the root is easy to eyeball.
func TestRootDefault0xHex(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/topup/dividend-account-root": {{body: loadFixture(t, "topup_dividend_account_root.json")}},
	})
	stdout, _, err := runCmd(t, srv.URL, "root")
	if err != nil {
		t.Fatalf("root: %v", err)
	}
	got := strings.TrimSpace(stdout)
	// Fixture decodes to this 0x-hex digest.
	want := "0x4b6b994b99d24e35e8626af419a087ec78498d976265e5879d7c22b9241c3b98"
	if got != want {
		t.Errorf("root: got %q, want %q", got, want)
	}
}

// TestRootRawPreservesBase64 asserts that --raw leaves the upstream
// base64 string intact.
func TestRootRawPreservesBase64(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/topup/dividend-account-root": {{body: loadFixture(t, "topup_dividend_account_root.json")}},
	})
	stdout, _, err := runCmd(t, srv.URL, "--raw", "root")
	if err != nil {
		t.Fatalf("root --raw: %v", err)
	}
	got := strings.TrimSpace(stdout)
	if got != "S2uZS5nSTjXoYmr0GaCH7HhJjZdiZeWHnXwiuSQcO5g=" {
		t.Errorf("root --raw: got %q, want base64", got)
	}
	// Must not contain 0x-hex.
	mustNotContain(t, stdout, "0x4b6b")
}

// TestRootJSON asserts that --json emits the full wrapper object.
// The byte-field normalization in render.RenderJSON still converts
// the hash to 0x-hex (because "root" is in the byte-field suffix
// list), unless --raw is also set.
func TestRootJSON(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/topup/dividend-account-root": {{body: loadFixture(t, "topup_dividend_account_root.json")}},
	})
	stdout, _, err := runCmd(t, srv.URL, "root", "--json")
	if err != nil {
		t.Fatalf("root --json: %v", err)
	}
	var m map[string]any
	if jerr := json.Unmarshal([]byte(stdout), &m); jerr != nil {
		t.Fatalf("root --json not valid JSON: %v\n%s", jerr, stdout)
	}
	if _, ok := m["account_root_hash"]; !ok {
		t.Errorf("expected account_root_hash key, got %v", m)
	}
}

// --- account ---

func TestAccountHappyPath(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/topup/dividend-account/0x02f615e95563ef16f10354dba9e584e58d2d4314": {{
			body: loadFixture(t, "topup_dividend_account.json"),
		}},
	})
	stdout, _, err := runCmd(t, srv.URL, "account", "0x02f615e95563ef16f10354dba9e584e58d2d4314")
	if err != nil {
		t.Fatalf("account: %v", err)
	}
	// Envelope unwrapped → inner fields visible.
	mustContain(t, stdout, "user")
	mustContain(t, stdout, "fee_amount")
	mustContain(t, stdout, "1000000000000000000")
}

func TestAccountNormalizesAddress(t *testing.T) {
	// Upper-case with 0x gets lower-cased before the URL is built.
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/topup/dividend-account/0x02f615e95563ef16f10354dba9e584e58d2d4314": {{
			body: loadFixture(t, "topup_dividend_account.json"),
		}},
	})
	_, _, err := runCmd(t, srv.URL, "account", "0x02F615E95563EF16F10354DBA9E584E58D2D4314")
	if err != nil {
		t.Fatalf("account upper-case: %v", err)
	}
}

func TestAccountAcceptsBareHex(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/topup/dividend-account/0x02f615e95563ef16f10354dba9e584e58d2d4314": {{
			body: loadFixture(t, "topup_dividend_account.json"),
		}},
	})
	_, _, err := runCmd(t, srv.URL, "account", "02f615e95563ef16f10354dba9e584e58d2d4314")
	if err != nil {
		t.Fatalf("account bare hex: %v", err)
	}
}

func TestAccountBadAddressIsUsageError(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{})
	_, _, err := runCmd(t, srv.URL, "account", "0xdeadbeef")
	if err == nil {
		t.Fatal("expected usage error for short address")
	}
	var uErr *client.UsageError
	if !errors.As(err, &uErr) {
		t.Fatalf("got %T, want *UsageError", err)
	}
}

func TestAccountNotFoundPropagates(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/topup/dividend-account/0x0000000000000000000000000000000000000000": {{
			status: 500,
			body:   loadFixture(t, "topup_dividend_account_not_found.json"),
		}},
	})
	_, _, err := runCmd(t, srv.URL, "account", "0x0000000000000000000000000000000000000000")
	if err == nil {
		t.Fatal("expected error for missing account")
	}
	var hErr *client.HTTPError
	if !errors.As(err, &hErr) {
		t.Fatalf("got %T, want *HTTPError", err)
	}
}

// --- proof ---

// TestProofL1UnconfiguredEmitsHint asserts that a node without
// `eth_rpc_url` (gRPC code 13 on /topup/account-proof/…) surfaces the
// L1-not-configured hint on stderr. Hint must NOT leak into stdout.
func TestProofL1UnconfiguredEmitsHint(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/topup/account-proof/0x02f615e95563ef16f10354dba9e584e58d2d4314": {{
			status: 500,
			body:   loadFixture(t, "topup_account_proof_l1_unconfigured.json"),
		}},
	})
	stdout, stderr, err := runCmd(t, srv.URL,
		"proof", "0x02f615e95563ef16f10354dba9e584e58d2d4314")
	if err == nil {
		t.Fatal("expected error on L1-unreachable")
	}
	mustContain(t, stderr, "eth_rpc_url")
	mustNotContain(t, stdout, "eth_rpc_url")
}

func TestProofHappyPath(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/topup/account-proof/0x02f615e95563ef16f10354dba9e584e58d2d4314": {{
			body: loadFixture(t, "topup_account_proof.json"),
		}},
	})
	stdout, _, err := runCmd(t, srv.URL,
		"proof", "0x02f615e95563ef16f10354dba9e584e58d2d4314")
	if err != nil {
		t.Fatalf("proof: %v", err)
	}
	// Envelope unwrapped; address + index visible.
	mustContain(t, stdout, "address")
	mustContain(t, stdout, "index")
	mustContain(t, stdout, "3")
}

// --- verify ---

func TestVerifyHappyPathFalse(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/topup/account-proof/0x02f615e95563ef16f10354dba9e584e58d2d4314/verify": {{
			body: loadFixture(t, "topup_verify_false.json"),
		}},
	})
	stdout, _, err := runCmd(t, srv.URL,
		"verify",
		"0x02f615e95563ef16f10354dba9e584e58d2d4314",
		"0x0000000000000000000000000000000000000000000000000000000000000000")
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if strings.TrimSpace(stdout) != "false" {
		t.Errorf("verify false: got %q, want \"false\"", stdout)
	}
}

func TestVerifyHappyPathTrue(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/topup/account-proof/0x02f615e95563ef16f10354dba9e584e58d2d4314/verify": {{
			body: loadFixture(t, "topup_verify_true.json"),
		}},
	})
	stdout, _, err := runCmd(t, srv.URL,
		"verify",
		"0x02f615e95563ef16f10354dba9e584e58d2d4314",
		"deadbeef")
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if strings.TrimSpace(stdout) != "true" {
		t.Errorf("verify true: got %q", stdout)
	}
}

// TestVerifyProofInQueryParam asserts that the proof travels on the
// `proof` query string (not in a POST body), matching the upstream
// GET /topup/account-proof/{address}/verify?proof=… route.
func TestVerifyProofInQueryParam(t *testing.T) {
	var gotQuery url.Values
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/topup/account-proof/0x02f615e95563ef16f10354dba9e584e58d2d4314/verify": {{
			match: func(q url.Values) bool { gotQuery = q; return true },
			body:  loadFixture(t, "topup_verify_false.json"),
		}},
	})
	_, _, err := runCmd(t, srv.URL,
		"verify",
		"0x02f615e95563ef16f10354dba9e584e58d2d4314",
		"0xDEADBEEF")
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if got := gotQuery.Get("proof"); got != "deadbeef" {
		t.Errorf("proof query param: got %q, want deadbeef", got)
	}
}

func TestVerifyEmptyProofIsUsageError(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{})
	_, _, err := runCmd(t, srv.URL, "verify",
		"0x02f615e95563ef16f10354dba9e584e58d2d4314",
		"0x")
	if err == nil {
		t.Fatal("expected usage error for empty proof")
	}
	var uErr *client.UsageError
	if !errors.As(err, &uErr) {
		t.Fatalf("got %T, want *UsageError", err)
	}
}

func TestVerifyBadProofServerPropagates(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/topup/account-proof/0x02f615e95563ef16f10354dba9e584e58d2d4314/verify": {{
			status: 400,
			body:   loadFixture(t, "topup_verify_bad_proof.json"),
		}},
	})
	_, _, err := runCmd(t, srv.URL,
		"verify",
		"0x02f615e95563ef16f10354dba9e584e58d2d4314",
		// 2 bytes — server rejects proofs not multiples of 32.
		"deadbeef")
	if err == nil {
		t.Fatal("expected error for invalid proof length")
	}
	var hErr *client.HTTPError
	if !errors.As(err, &hErr) {
		t.Fatalf("got %T, want *HTTPError", err)
	}
}

// --- sequence ---

func TestSequenceL1UnconfiguredEmitsHint(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/topup/sequence": {{
			status: 500,
			body:   loadFixture(t, "topup_sequence_l1_unconfigured.json"),
		}},
	})
	_, stderr, err := runCmd(t, srv.URL,
		"sequence",
		"0x48bd44a37ff84c7cc584e5df4bf43bbda6116d5708d41080e7b9b030195c6bf8",
		"423")
	if err == nil {
		t.Fatal("expected error on L1-unreachable")
	}
	mustContain(t, stderr, "eth_rpc_url")
}

func TestSequenceHappyPath(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/topup/sequence": {{body: loadFixture(t, "topup_sequence.json")}},
	})
	stdout, _, err := runCmd(t, srv.URL,
		"sequence",
		"0x48bd44a37ff84c7cc584e5df4bf43bbda6116d5708d41080e7b9b030195c6bf8",
		"423")
	if err != nil {
		t.Fatalf("sequence: %v", err)
	}
	// Default text mode prints just the sequence.
	if strings.TrimSpace(stdout) != "11602043000000" {
		t.Errorf("sequence: got %q, want 11602043000000", stdout)
	}
}

// TestSequenceQueryParams asserts that `sequence` sends tx_hash +
// log_index as query parameters (confirmed from heimdall-v2 query.proto).
func TestSequenceQueryParams(t *testing.T) {
	var gotQuery url.Values
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/topup/sequence": {{
			match: func(q url.Values) bool { gotQuery = q; return true },
			body:  loadFixture(t, "topup_sequence.json"),
		}},
	})
	_, _, err := runCmd(t, srv.URL,
		"sequence",
		"0x48bd44a37ff84c7cc584e5df4bf43bbda6116d5708d41080e7b9b030195c6bf8",
		"423")
	if err != nil {
		t.Fatalf("sequence: %v", err)
	}
	if got := gotQuery.Get("tx_hash"); got != "0x48bd44a37ff84c7cc584e5df4bf43bbda6116d5708d41080e7b9b030195c6bf8" {
		t.Errorf("tx_hash=%q", got)
	}
	if got := gotQuery.Get("log_index"); got != "423" {
		t.Errorf("log_index=%q, want 423", got)
	}
}

func TestSequenceBadTxHashIsUsageError(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{})
	_, _, err := runCmd(t, srv.URL, "sequence", "0xdeadbeef", "0")
	if err == nil {
		t.Fatal("expected usage error for short tx hash")
	}
	var uErr *client.UsageError
	if !errors.As(err, &uErr) {
		t.Fatalf("got %T, want *UsageError", err)
	}
}

func TestSequenceBadLogIndexIsUsageError(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{})
	_, _, err := runCmd(t, srv.URL, "sequence",
		"0x48bd44a37ff84c7cc584e5df4bf43bbda6116d5708d41080e7b9b030195c6bf8",
		"notanumber")
	if err == nil {
		t.Fatal("expected usage error for non-integer log_index")
	}
	var uErr *client.UsageError
	if !errors.As(err, &uErr) {
		t.Fatalf("got %T, want *UsageError", err)
	}
}

// --- is-old ---

func TestIsOldL1UnconfiguredEmitsHint(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/topup/is-old-tx": {{
			status: 500,
			body:   loadFixture(t, "topup_is_old_tx_l1_unconfigured.json"),
		}},
	})
	_, stderr, err := runCmd(t, srv.URL,
		"is-old",
		"0x48bd44a37ff84c7cc584e5df4bf43bbda6116d5708d41080e7b9b030195c6bf8",
		"423")
	if err == nil {
		t.Fatal("expected error on L1-unreachable")
	}
	mustContain(t, stderr, "eth_rpc_url")
}

func TestIsOldHappyPathTrue(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/topup/is-old-tx": {{body: loadFixture(t, "topup_is_old_tx_true.json")}},
	})
	stdout, _, err := runCmd(t, srv.URL,
		"is-old",
		"0x48bd44a37ff84c7cc584e5df4bf43bbda6116d5708d41080e7b9b030195c6bf8",
		"423")
	if err != nil {
		t.Fatalf("is-old: %v", err)
	}
	if strings.TrimSpace(stdout) != "true" {
		t.Errorf("is-old true: got %q", stdout)
	}
}

func TestIsOldHappyPathFalse(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/topup/is-old-tx": {{body: loadFixture(t, "topup_is_old_tx_false.json")}},
	})
	stdout, _, err := runCmd(t, srv.URL,
		"is-old",
		"0x48bd44a37ff84c7cc584e5df4bf43bbda6116d5708d41080e7b9b030195c6bf8",
		"0")
	if err != nil {
		t.Fatalf("is-old: %v", err)
	}
	if strings.TrimSpace(stdout) != "false" {
		t.Errorf("is-old false: got %q", stdout)
	}
}

// TestIsOldJSONPassthrough asserts that --json emits the raw server
// payload (no bare-bool shortcut).
func TestIsOldJSONPassthrough(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/topup/is-old-tx": {{body: loadFixture(t, "topup_is_old_tx_true.json")}},
	})
	stdout, _, err := runCmd(t, srv.URL,
		"is-old", "--json",
		"0x48bd44a37ff84c7cc584e5df4bf43bbda6116d5708d41080e7b9b030195c6bf8",
		"423")
	if err != nil {
		t.Fatalf("is-old --json: %v", err)
	}
	var m map[string]any
	if jerr := json.Unmarshal([]byte(stdout), &m); jerr != nil {
		t.Fatalf("not valid JSON: %v\n%s", jerr, stdout)
	}
	if got, ok := m["is_old"].(bool); !ok || !got {
		t.Errorf("is_old: got %v, want true", m["is_old"])
	}
}
