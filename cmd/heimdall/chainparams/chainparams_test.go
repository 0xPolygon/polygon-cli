package chainparams

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
)

// --- params ---

// TestParamsHappyPathKV asserts the default human output unwraps the
// `params` envelope and surfaces the confirmation depths plus a nested
// chain_params blob.
func TestParamsHappyPathKV(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/chainmanager/params": {{body: loadFixture(t, "chainmanager_params.json")}},
	})
	stdout, _, err := runCmd(t, srv.URL, "params")
	if err != nil {
		t.Fatalf("params: %v", err)
	}
	mustContain(t, stdout, "main_chain_tx_confirmations")
	mustContain(t, stdout, "64")
	mustContain(t, stdout, "bor_chain_tx_confirmations")
	mustContain(t, stdout, "512")
	// chain_params should be rendered inline as JSON on a single line;
	// contents should include at least one well-known address key.
	mustContain(t, stdout, "root_chain_address")
	mustContain(t, stdout, "0xbd07d7e1e93c8d4b2a261327f3c28a8ea7167209")
}

// TestParamsJSONPassthrough asserts --json emits the raw server shape
// (with the `params` envelope preserved) and parses cleanly.
func TestParamsJSONPassthrough(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/chainmanager/params": {{body: loadFixture(t, "chainmanager_params.json")}},
	})
	stdout, _, err := runCmd(t, srv.URL, "params", "--json")
	if err != nil {
		t.Fatalf("params --json: %v", err)
	}
	var m map[string]any
	if jerr := json.Unmarshal([]byte(stdout), &m); jerr != nil {
		t.Fatalf("not valid JSON: %v\n%s", jerr, stdout)
	}
	params, ok := m["params"].(map[string]any)
	if !ok {
		t.Fatalf("expected params object, got %v", m)
	}
	if _, ok := params["chain_params"].(map[string]any); !ok {
		t.Errorf("expected chain_params object, got %v", params)
	}
	if got := params["main_chain_tx_confirmations"]; got != "64" {
		t.Errorf("main_chain_tx_confirmations=%v want \"64\"", got)
	}
}

// TestParamsFieldPluck asserts --field can drill into the nested
// params object using dot-notation.
func TestParamsFieldPluck(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/chainmanager/params": {{body: loadFixture(t, "chainmanager_params.json")}},
	})
	stdout, _, err := runCmd(t, srv.URL, "params",
		"--field", "params.chain_params.root_chain_address")
	if err != nil {
		t.Fatalf("params --field: %v", err)
	}
	got := strings.TrimSpace(stdout)
	if got != "0xbd07d7e1e93c8d4b2a261327f3c28a8ea7167209" {
		t.Errorf("root_chain_address: got %q", got)
	}
}

// TestParamsNotImplementedPropagates asserts that upstream gRPC-gateway
// errors (code 12, status 501) surface as *HTTPError so the outer
// exit-code mapper can return a non-zero status.
func TestParamsNotImplementedPropagates(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/chainmanager/params": {{
			status: 501,
			body:   loadFixture(t, "chainmanager_not_implemented.json"),
		}},
	})
	_, _, err := runCmd(t, srv.URL, "params")
	if err == nil {
		t.Fatal("expected error for 501 Not Implemented")
	}
	var hErr *client.HTTPError
	if !errors.As(err, &hErr) {
		t.Fatalf("got %T, want *HTTPError", err)
	}
}

// --- addresses ---

// TestAddressesHappyPath asserts the derived view lists every
// `*_address` entry + the two chain ids, one per line, alphabetized.
func TestAddressesHappyPath(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/chainmanager/params": {{body: loadFixture(t, "chainmanager_params.json")}},
	})
	stdout, _, err := runCmd(t, srv.URL, "addresses")
	if err != nil {
		t.Fatalf("addresses: %v", err)
	}
	// All eight *_address entries + both chain ids must appear.
	expected := []string{
		"bor_chain_id=80002",
		"heimdall_chain_id=heimdallv2-80002",
		"pol_token_address=0x3fd0a53f4bf853985a95f4eb3f9c9fde1f8e2b53",
		"staking_manager_address=0x4ae8f648b1ec892b6cc68c89cc088583964d08be",
		"slash_manager_address=0x9e699267858ce513eacf3b66420334785f9c8e4c",
		"root_chain_address=0xbd07d7e1e93c8d4b2a261327f3c28a8ea7167209",
		"staking_info_address=0x5e3111a5d928d24718c1a7897261d0b9087002ed",
		"state_sender_address=0x49e307fa5a58ff1834e0f8a60eb2a9609e6a5f50",
		"state_receiver_address=0x0000000000000000000000000000000000001001",
		"validator_set_address=0x0000000000000000000000000000000000001000",
	}
	for _, line := range expected {
		mustContain(t, stdout, line)
	}
	// Confirmation depths should NOT leak into the addresses view.
	if strings.Contains(stdout, "tx_confirmations") {
		t.Errorf("addresses leaked confirmation depths:\n%s", stdout)
	}
	// Verify alphabetical ordering for the first three keys.
	idxBor := strings.Index(stdout, "bor_chain_id=")
	idxHeimdall := strings.Index(stdout, "heimdall_chain_id=")
	idxPol := strings.Index(stdout, "pol_token_address=")
	if !(idxBor < idxHeimdall && idxHeimdall < idxPol) {
		t.Errorf("addresses not alphabetized: bor=%d heimdall=%d pol=%d", idxBor, idxHeimdall, idxPol)
	}
}

// TestAddressesJSON asserts --json emits the same derived map.
func TestAddressesJSON(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/chainmanager/params": {{body: loadFixture(t, "chainmanager_params.json")}},
	})
	stdout, _, err := runCmd(t, srv.URL, "addresses", "--json")
	if err != nil {
		t.Fatalf("addresses --json: %v", err)
	}
	var m map[string]any
	if jerr := json.Unmarshal([]byte(stdout), &m); jerr != nil {
		t.Fatalf("not valid JSON: %v\n%s", jerr, stdout)
	}
	if got := m["bor_chain_id"]; got != "80002" {
		t.Errorf("bor_chain_id=%v", got)
	}
	if got := m["root_chain_address"]; got != "0xbd07d7e1e93c8d4b2a261327f3c28a8ea7167209" {
		t.Errorf("root_chain_address=%v", got)
	}
	// Confirmation depth fields must not appear.
	if _, ok := m["main_chain_tx_confirmations"]; ok {
		t.Errorf("confirmation depth leaked into addresses --json output: %v", m)
	}
}

// TestAddressesFieldPluck asserts --field works over the derived view
// when combined with --json.
func TestAddressesFieldPluck(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/chainmanager/params": {{body: loadFixture(t, "chainmanager_params.json")}},
	})
	stdout, _, err := runCmd(t, srv.URL, "addresses", "--json",
		"--field", "root_chain_address")
	if err != nil {
		t.Fatalf("addresses --field: %v", err)
	}
	got := strings.TrimSpace(stdout)
	// Single --field produces the bare value, quoted by JSON string encoding.
	want := `"0xbd07d7e1e93c8d4b2a261327f3c28a8ea7167209"`
	if got != want && got != "0xbd07d7e1e93c8d4b2a261327f3c28a8ea7167209" {
		t.Errorf("addresses --field: got %q, want %q", got, want)
	}
}

// TestAddressesServerErrorPropagates asserts that upstream errors
// surface as *HTTPError, same as the params command.
func TestAddressesServerErrorPropagates(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/chainmanager/params": {{
			status: 500,
			body:   []byte(`{"code":2,"message":"internal","details":[]}`),
		}},
	})
	_, _, err := runCmd(t, srv.URL, "addresses")
	if err == nil {
		t.Fatal("expected error for 500")
	}
	var hErr *client.HTTPError
	if !errors.As(err, &hErr) {
		t.Fatalf("got %T, want *HTTPError", err)
	}
}

// TestAddressesMalformedBody asserts that a response lacking the
// params.chain_params envelope surfaces a clear error rather than
// silently emitting an empty map.
func TestAddressesMalformedBody(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/chainmanager/params": {{body: []byte(`{"params":{}}`)}},
	})
	_, _, err := runCmd(t, srv.URL, "addresses")
	if err == nil {
		t.Fatal("expected error for missing chain_params")
	}
	if !strings.Contains(err.Error(), "chain_params") {
		t.Errorf("error does not mention chain_params: %v", err)
	}
}

// TestAliasCM asserts the `cm` alias reaches the same subcommands.
func TestAliasCM(t *testing.T) {
	// The helper injects the literal word `chainmanager` into argv; to
	// test the alias path we drive the root cobra command directly.
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/chainmanager/params": {{body: loadFixture(t, "chainmanager_params.json")}},
	})
	// Re-implement the minimal runner inline so we can use the alias.
	local := newLocalCmdWithAlias()
	stdout, err := runRootWithAlias(t, srv.URL, local, "cm", "params")
	if err != nil {
		t.Fatalf("cm params: %v", err)
	}
	mustContain(t, stdout, "main_chain_tx_confirmations")
}
