package milestone

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
)

// --- params ---

func TestParamsHumanOutput(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/milestones/params": {body: loadFixture(t, "rest", "milestones_params.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "params")
	if err != nil {
		t.Fatalf("params: %v", err)
	}
	for _, want := range []string{"max_milestone_proposition_length", "ff_milestone_threshold", "ff_milestone_block_interval"} {
		mustContain(t, stdout, want)
	}
}

func TestParamsJSONPassthrough(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/milestones/params": {body: loadFixture(t, "rest", "milestones_params.json")},
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
		"/milestones/count": {body: loadFixture(t, "rest", "milestones_count.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "count")
	if err != nil {
		t.Fatalf("count: %v", err)
	}
	// Fixture: {"count":"11597445"}. Default output is just the number.
	if strings.TrimSpace(stdout) != "11597445" {
		t.Errorf("expected bare count 11597445, got %q", stdout)
	}
}

func TestCountJSON(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/milestones/count": {body: loadFixture(t, "rest", "milestones_count.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "count", "--json")
	if err != nil {
		t.Fatalf("count --json: %v", err)
	}
	var m map[string]any
	if jerr := json.Unmarshal([]byte(stdout), &m); jerr != nil {
		t.Fatalf("count --json not valid JSON: %v\n%s", jerr, stdout)
	}
	if got := m["count"]; got != "11597445" {
		t.Errorf("expected count=\"11597445\", got %v", got)
	}
}

// --- latest ---

func TestLatestUnwrapsEnvelopeAndHashAsHex(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/milestones/latest": {body: loadFixture(t, "rest", "milestones_latest.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "latest")
	if err != nil {
		t.Fatalf("latest: %v", err)
	}
	// Envelope unwrapped: body fields visible.
	mustContain(t, stdout, "proposer")
	mustContain(t, stdout, "start_block")
	mustContain(t, stdout, "end_block")
	mustContain(t, stdout, "milestone_id")
	// Hash must be re-encoded from base64 to 0x-hex.
	mustContain(t, stdout, "0xc866a7811d701a548def61e9347455aff6e5eef37d07a155eae281b1199bbc3b")
	// Raw base64 hash must NOT leak into default KV output.
	mustNotContain(t, stdout, "yGangR1wGlSN72HpNHRVr/bl7vN9B6FV6uKBsRmbvDs=")
}

func TestLatestRawPreservesBase64(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/milestones/latest": {body: loadFixture(t, "rest", "milestones_latest.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "latest", "--raw")
	if err != nil {
		t.Fatalf("latest --raw: %v", err)
	}
	mustContain(t, stdout, "yGangR1wGlSN72HpNHRVr/bl7vN9B6FV6uKBsRmbvDs=")
}

// TestLatestPrintsMilestoneIDButNoNumber asserts our design decision
// for `latest`: the server doesn't return `number` and we don't splice
// one in (since the value would just be `milestone count`, which the
// user can ask for separately). `milestone_id` must still be visible.
func TestLatestPrintsMilestoneIDButNoNumber(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/milestones/latest": {body: loadFixture(t, "rest", "milestones_latest.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "latest")
	if err != nil {
		t.Fatalf("latest: %v", err)
	}
	mustContain(t, stdout, "milestone_id")
	mustContain(t, stdout, "5c40dfd345378cd83580475ea0c62a7584c482c24a3851ed8a3a3d76ca8066ac")
	// `number` label must NOT appear at the start of a line (the
	// timestamp-annotation helper may print the word elsewhere, but
	// no "number  " left-aligned key should be emitted). The KV
	// renderer produces `key<spaces>  value`.
	for _, line := range strings.Split(strings.TrimRight(stdout, "\n"), "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "number ") {
			t.Errorf("latest should not render a number label: %q", line)
		}
	}
}

// --- get / bare integer ---

// TestGetBareIntegerPrintsBothNumberAndMilestoneID is the single most
// important test in this package: the URL path carries the sequence
// number (`number`), and the response body carries `milestone_id`; the
// two are DIFFERENT VALUES on real Heimdall data. We must render both.
func TestGetBareIntegerPrintsBothNumberAndMilestoneID(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/milestones/1": {body: loadFixture(t, "rest", "milestones_by_number_one.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "1")
	if err != nil {
		t.Fatalf("milestone 1: %v", err)
	}
	// The number label from the URL path.
	mustContain(t, stdout, "number")
	// The milestone_id label from the body.
	mustContain(t, stdout, "milestone_id")
	// The canonical `milestone_id` for milestone #1 on Amoy is a
	// UUID + creator-address string (a genesis artefact); its mere
	// presence confirms that `milestone_id` != `number`.
	mustContain(t, stdout, "cd8b33d3-5c87-49c4-b391-7ee296a058f9")
	// The URL-path number must render as a standalone "1".
	var sawNumber bool
	for _, line := range strings.Split(strings.TrimRight(stdout, "\n"), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "number ") && strings.HasSuffix(trimmed, " 1") {
			sawNumber = true
			break
		}
	}
	if !sawNumber {
		t.Errorf("expected a 'number  1' row in output:\n%s", stdout)
	}
}

func TestGetExplicitSubcommand(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/milestones/1": {body: loadFixture(t, "rest", "milestones_by_number_one.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "get", "1")
	if err != nil {
		t.Fatalf("get 1: %v", err)
	}
	mustContain(t, stdout, "milestone_id")
}

func TestGetJSONIncludesSplicedNumber(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/milestones/1": {body: loadFixture(t, "rest", "milestones_by_number_one.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "1", "--json")
	if err != nil {
		t.Fatalf("milestone 1 --json: %v", err)
	}
	var env struct {
		Milestone struct {
			Number      string `json:"number"`
			MilestoneID string `json:"milestone_id"`
		} `json:"milestone"`
	}
	if jerr := json.Unmarshal([]byte(stdout), &env); jerr != nil {
		t.Fatalf("output not valid JSON: %v\n%s", jerr, stdout)
	}
	if env.Milestone.Number != "1" {
		t.Errorf("expected number=\"1\", got %q", env.Milestone.Number)
	}
	if env.Milestone.MilestoneID == "" {
		t.Errorf("expected milestone_id to be non-empty")
	}
	if env.Milestone.Number == env.Milestone.MilestoneID {
		t.Errorf("number and milestone_id must differ; both were %q", env.Milestone.Number)
	}
}

func TestGetHashRenderedAsHex(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/milestones/11597445": {body: loadFixture(t, "rest", "milestones_by_number.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "11597445")
	if err != nil {
		t.Fatalf("milestone 11597445: %v", err)
	}
	mustContain(t, stdout, "0xc866a7811d701a548def61e9347455aff6e5eef37d07a155eae281b1199bbc3b")
}

func TestGetInvalidArgIsUsageError(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{})
	_, _, err := runCmd(t, srv.URL, "get", "notanumber")
	if err == nil {
		t.Fatal("expected error for non-integer number")
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

// --- out-of-range hint ---

// TestGetOutOfRangeNumberEmitsHint covers the happy path of the hint:
// the server returns 404, /milestones/count is reachable, and the
// requested number exceeds the count. Hint must travel on stderr.
func TestGetOutOfRangeNumberEmitsHint(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/milestones/99999999": {status: 404, body: loadFixture(t, "rest", "milestones_out_of_range.json")},
		"/milestones/count":    {body: loadFixture(t, "rest", "milestones_count.json")},
	})
	_, stderr, err := runCmd(t, srv.URL, "99999999")
	if err == nil {
		t.Fatal("expected error on out-of-range milestone")
	}
	mustContain(t, stderr, "valid range is 1..11597445")
}

// TestGetZeroNumberEmitsHint asserts that number=0 also triggers the
// hint even though the server's error message is the same ("milestone
// number out of range"), because 0 is strictly below the valid range.
func TestGetZeroNumberEmitsHint(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/milestones/0":     {status: 404, body: loadFixture(t, "rest", "milestones_out_of_range.json")},
		"/milestones/count": {body: loadFixture(t, "rest", "milestones_count.json")},
	})
	_, stderr, err := runCmd(t, srv.URL, "0")
	if err == nil {
		t.Fatal("expected error on milestone 0")
	}
	mustContain(t, stderr, "valid range is 1..11597445")
}

// TestGetInRange404DoesNotEmitHint covers the edge case where the
// server returns 404 for a number within the valid range. Rare but
// possible if a milestone row is pruned; in that case the hint would
// be incorrect and must not be emitted.
func TestGetInRange404DoesNotEmitHint(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/milestones/5":     {status: 404, body: loadFixture(t, "rest", "milestones_out_of_range.json")},
		"/milestones/count": {body: loadFixture(t, "rest", "milestones_count.json")},
	})
	_, stderr, err := runCmd(t, srv.URL, "5")
	if err == nil {
		t.Fatal("expected error")
	}
	// 5 <= count, so no range hint.
	mustNotContain(t, stderr, "valid range is")
}
