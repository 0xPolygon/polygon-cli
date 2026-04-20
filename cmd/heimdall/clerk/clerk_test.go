package clerk

import (
	"encoding/json"
	"errors"
	"net/url"
	"strings"
	"testing"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
)

// --- count ---

func TestCountBareInteger(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/clerk/event-records/count": {{body: loadFixture(t, "rest", "clerk_event_records_count.json")}},
	})
	stdout, _, err := runCmd(t, srv.URL, "count")
	if err != nil {
		t.Fatalf("count: %v", err)
	}
	// Fixture: {"count":"36610"}. Default output is just the number.
	if strings.TrimSpace(stdout) != "36610" {
		t.Errorf("expected bare count 36610, got %q", stdout)
	}
}

func TestCountJSON(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/clerk/event-records/count": {{body: loadFixture(t, "rest", "clerk_event_records_count.json")}},
	})
	stdout, _, err := runCmd(t, srv.URL, "count", "--json")
	if err != nil {
		t.Fatalf("count --json: %v", err)
	}
	var m map[string]any
	if jerr := json.Unmarshal([]byte(stdout), &m); jerr != nil {
		t.Fatalf("count --json not valid JSON: %v\n%s", jerr, stdout)
	}
	if got := m["count"]; got != "36610" {
		t.Errorf("expected count=\"36610\", got %v", got)
	}
}

// --- latest-id ---

// TestLatestIDL1UnconfiguredEmitsHint asserts that a node without
// `eth_rpc_url` (gRPC code 13 on /clerk/event-records/latest-id)
// surfaces the L1-not-configured hint. The hint must travel on stderr
// so --json / -f output on stdout stays clean for scripts.
func TestLatestIDL1UnconfiguredEmitsHint(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/clerk/event-records/latest-id": {{
			status: 500,
			body:   loadFixture(t, "rest", "clerk_latest_id_l1_unconfigured.json"),
		}},
	})
	_, stderr, err := runCmd(t, srv.URL, "latest-id")
	if err == nil {
		t.Fatal("expected error on L1-unreachable")
	}
	mustContain(t, stderr, "eth_rpc_url")
}

// --- get / bare integer ---

func TestGetBareInteger(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/clerk/event-records/36610": {{body: loadFixture(t, "rest", "clerk_event_record_by_id.json")}},
	})
	stdout, _, err := runCmd(t, srv.URL, "36610")
	if err != nil {
		t.Fatalf("state-sync 36610: %v", err)
	}
	// Envelope unwrapped: body fields visible.
	mustContain(t, stdout, "contract")
	mustContain(t, stdout, "tx_hash")
	mustContain(t, stdout, "log_index")
	// The fixture's data base64 starts with "h6eB..." which decodes to
	// 0x87a7811f4bfedea3... . Assert the hex prefix is surfaced so we
	// prove `data` normalization kicks in.
	mustContain(t, stdout, "0x87a7811f4bfedea3")
	// The raw base64 must NOT leak into default output.
	mustNotContain(t, stdout, "h6eBH0v+3qPTQa0W")
}

func TestGetExplicitSubcommand(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/clerk/event-records/36610": {{body: loadFixture(t, "rest", "clerk_event_record_by_id.json")}},
	})
	stdout, _, err := runCmd(t, srv.URL, "get", "36610")
	if err != nil {
		t.Fatalf("get 36610: %v", err)
	}
	mustContain(t, stdout, "tx_hash")
}

func TestGetBase64PreservesRaw(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/clerk/event-records/36610": {{body: loadFixture(t, "rest", "clerk_event_record_by_id.json")}},
	})
	stdout, _, err := runCmd(t, srv.URL, "get", "36610", "--base64")
	if err != nil {
		t.Fatalf("get --base64: %v", err)
	}
	mustContain(t, stdout, "h6eBH0v+3qPTQa0W")
	mustNotContain(t, stdout, "0x87a7811f4bfedea3")
}

func TestGetRawFlagAlsoPreservesBase64(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/clerk/event-records/36610": {{body: loadFixture(t, "rest", "clerk_event_record_by_id.json")}},
	})
	stdout, _, err := runCmd(t, srv.URL, "get", "36610", "--raw")
	if err != nil {
		t.Fatalf("get --raw: %v", err)
	}
	mustContain(t, stdout, "h6eBH0v+3qPTQa0W")
}

func TestGetInvalidArgIsUsageError(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{})
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
	srv := newRESTFixtureServer(t, map[string][]restRoute{})
	_, _, err := runCmd(t, srv.URL, "banana")
	if err == nil {
		t.Fatal("expected usage error for unknown umbrella arg")
	}
	var uErr *client.UsageError
	if !errors.As(err, &uErr) {
		t.Fatalf("got %T, want *UsageError", err)
	}
}

// --- list (page-based, NOT cosmos pagination) ---

// TestListPageBasedQueryParams verifies that `list` emits bare `page`
// and `limit` query params — not Cosmos `pagination.*` — mirroring the
// upstream /clerk/event-records/list shape.
func TestListPageBasedQueryParams(t *testing.T) {
	var gotQuery url.Values
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/clerk/event-records/list": {{
			match: func(q url.Values) bool { gotQuery = q; return true },
			body:  loadFixture(t, "rest", "clerk_event_records_list.json"),
		}},
	})
	_, _, err := runCmd(t, srv.URL, "list", "--page", "2", "--limit", "5")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if got := gotQuery.Get("page"); got != "2" {
		t.Errorf("page=%q, want 2", got)
	}
	if got := gotQuery.Get("limit"); got != "5" {
		t.Errorf("limit=%q, want 5", got)
	}
	// Must NOT emit cosmos-style params.
	if got := gotQuery.Get("pagination.limit"); got != "" {
		t.Errorf("unexpected pagination.limit=%q", got)
	}
	if got := gotQuery.Get("pagination.key"); got != "" {
		t.Errorf("unexpected pagination.key=%q", got)
	}
}

// TestListWithoutLimitEmitsHint asserts that the pagination-limit
// hint is surfaced when --limit is omitted. The hint travels on stderr
// so it does not leak into scripting pipelines.
func TestListWithoutLimitEmitsHint(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/clerk/event-records/list": {{body: loadFixture(t, "rest", "clerk_event_records_list.json")}},
	})
	_, stderr, err := runCmd(t, srv.URL, "list")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	mustContain(t, stderr, "pagination.limit")
}

func TestListWithLimitNoHint(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/clerk/event-records/list": {{body: loadFixture(t, "rest", "clerk_event_records_list.json")}},
	})
	_, stderr, err := runCmd(t, srv.URL, "list", "--limit", "5")
	if err != nil {
		t.Fatalf("list --limit: %v", err)
	}
	mustNotContain(t, stderr, "pagination.limit")
}

// TestListTableOutput renders the fixture and asserts the summary
// columns appear. The `data` blob is intentionally dropped by
// renderRecordTable; confirm it did not bleed into stdout.
func TestListTableOutput(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/clerk/event-records/list": {{body: loadFixture(t, "rest", "clerk_event_records_list.json")}},
	})
	stdout, _, err := runCmd(t, srv.URL, "list", "--limit", "3")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	mustContain(t, stdout, "tx_hash")
	mustContain(t, stdout, "contract")
	mustContain(t, stdout, "log_index")
	// data column removed from summary table.
	for _, line := range strings.Split(stdout, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "data ") {
			t.Errorf("table should not include data column: %q", line)
		}
	}
}

// TestListPage0RejectedByServerPropagates: when the server returns
// HTTP 400 (the actual upstream behaviour for page=0), the command
// must surface the error and NOT swallow it.
func TestListPage0RejectedByServerPropagates(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/clerk/event-records/list": {{
			status: 400,
			body:   loadFixture(t, "rest", "clerk_event_records_list_page_0.json"),
		}},
	})
	_, _, err := runCmd(t, srv.URL, "list", "--page", "0", "--limit", "5")
	if err == nil {
		t.Fatal("expected error for page=0")
	}
}

// --- range (/clerk/time) ---

func TestRangeRequiresFromID(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{})
	_, _, err := runCmd(t, srv.URL, "range")
	if err == nil {
		t.Fatal("expected usage error when --from-id missing")
	}
	var uErr *client.UsageError
	if !errors.As(err, &uErr) {
		t.Fatalf("got %T, want *UsageError", err)
	}
}

// TestRangeQueryParams asserts that /clerk/time is called with
// `from_id` / `to_time` / `pagination.limit` — note this is the ONE
// clerk endpoint that uses cosmos pagination on the server, unlike
// /clerk/event-records/list which is page-based.
func TestRangeQueryParams(t *testing.T) {
	var gotQuery url.Values
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/clerk/time": {{
			match: func(q url.Values) bool { gotQuery = q; return true },
			body:  loadFixture(t, "rest", "clerk_time_range.json"),
		}},
	})
	_, _, err := runCmd(t, srv.URL, "range",
		"--from-id", "36600",
		"--to-time", "2026-04-20T23:00:00Z",
		"--limit", "3")
	if err != nil {
		t.Fatalf("range: %v", err)
	}
	if got := gotQuery.Get("from_id"); got != "36600" {
		t.Errorf("from_id=%q, want 36600", got)
	}
	if got := gotQuery.Get("to_time"); got != "2026-04-20T23:00:00Z" {
		t.Errorf("to_time=%q, want RFC3339 upper bound", got)
	}
	if got := gotQuery.Get("pagination.limit"); got != "3" {
		t.Errorf("pagination.limit=%q, want 3", got)
	}
	// Must NOT emit bare page/limit (that's the /list endpoint).
	if got := gotQuery.Get("page"); got != "" {
		t.Errorf("unexpected page=%q", got)
	}
	if got := gotQuery.Get("limit"); got != "" {
		t.Errorf("unexpected bare limit=%q", got)
	}
}

// --- sequence / is-old ---

func TestSequenceL1UnconfiguredEmitsHint(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/clerk/sequence": {{
			status: 500,
			body:   loadFixture(t, "rest", "clerk_sequence_l1_unconfigured.json"),
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

func TestIsOldL1UnconfiguredEmitsHint(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/clerk/is-old-tx": {{
			status: 500,
			body:   loadFixture(t, "rest", "clerk_is_old_tx_l1_unconfigured.json"),
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

func TestIsOldBadTxHashIsUsageError(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{})
	_, _, err := runCmd(t, srv.URL, "is-old", "0xdeadbeef", "0")
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

// TestGetNotFoundPropagates: the upstream returns HTTP 500 + gRPC code
// 13 for a missing id (not 404). We don't transform that — the shape
// is a plain HTTPError surfaced to the caller for exit-code mapping.
func TestGetNotFoundPropagates(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string][]restRoute{
		"/clerk/event-records/99999999": {{
			status: 500,
			body:   loadFixture(t, "rest", "clerk_event_record_not_found.json"),
		}},
	})
	_, _, err := runCmd(t, srv.URL, "99999999")
	if err == nil {
		t.Fatal("expected error for missing id")
	}
	var hErr *client.HTTPError
	if !errors.As(err, &hErr) {
		t.Fatalf("got %T, want *HTTPError", err)
	}
}
