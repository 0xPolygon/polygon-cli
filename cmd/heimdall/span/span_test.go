package span

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
)

// --- params ---

func TestParamsHumanOutput(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/bor/params": {body: loadFixture(t, "rest", "bor_params.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "params")
	if err != nil {
		t.Fatalf("params: %v", err)
	}
	for _, want := range []string{"sprint_duration", "span_duration", "producer_count"} {
		mustContain(t, stdout, want)
	}
}

func TestParamsJSONPassthrough(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/bor/params": {body: loadFixture(t, "rest", "bor_params.json")},
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

// --- latest / get / bare-id ---

func TestLatestUnwrapsEnvelope(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/bor/spans/latest": {body: loadFixture(t, "rest", "bor_spans_latest.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "latest")
	if err != nil {
		t.Fatalf("latest: %v", err)
	}
	mustContain(t, stdout, "start_block")
	mustContain(t, stdout, "end_block")
	mustContain(t, stdout, "bor_chain_id")
}

func TestGetByExplicitSubcommand(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/bor/spans/5982": {body: loadFixture(t, "rest", "bor_spans_by_id.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "get", "5982")
	if err != nil {
		t.Fatalf("get 5982: %v", err)
	}
	mustContain(t, stdout, "5982")
	mustContain(t, stdout, "selected_producers")
}

func TestGetBareIntegerShortcut(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/bor/spans/5982": {body: loadFixture(t, "rest", "bor_spans_by_id.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "5982")
	if err != nil {
		t.Fatalf("span 5982: %v", err)
	}
	mustContain(t, stdout, "5982")
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

// --- list ---

func TestListDefaultsAndTable(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/bor/spans/list": {body: loadFixture(t, "rest", "bor_spans_list.json")},
	})
	stdout, stderr, err := runCmd(t, srv.URL, "list")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	// Columns include the summary fields.
	mustContain(t, stdout, "id")
	mustContain(t, stdout, "start_block")
	mustContain(t, stdout, "end_block")
	mustContain(t, stdout, "producers")
	// A known value from the fixture.
	mustContain(t, stdout, "5982")
	// Pagination key surfaces on stderr.
	mustContain(t, stderr, "next_key=")
}

func TestListJSONPassthrough(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/bor/spans/list": {body: loadFixture(t, "rest", "bor_spans_list.json")},
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

// --- producers (derived) ---

func TestProducersListsOnlySelected(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/bor/spans/5982": {body: loadFixture(t, "rest", "bor_spans_by_id.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "producers", "5982")
	if err != nil {
		t.Fatalf("producers: %v", err)
	}
	// The fixture has exactly one selected producer with val_id=5.
	mustContain(t, stdout, "val_id")
	mustContain(t, stdout, "signer")
	// selected_producers[0].signer from the fixture.
	mustContain(t, stdout, "0x6dc2dd54f24979ec26212794c71afefed722280c")
}

func TestProducersJSONIsArray(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/bor/spans/5982": {body: loadFixture(t, "rest", "bor_spans_by_id.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "producers", "5982", "--json")
	if err != nil {
		t.Fatalf("producers --json: %v", err)
	}
	var arr []any
	if jerr := json.Unmarshal([]byte(stdout), &arr); jerr != nil {
		t.Fatalf("output not a JSON array: %v\n%s", jerr, stdout)
	}
	if len(arr) != 1 {
		t.Errorf("expected 1 selected producer in fixture, got %d", len(arr))
	}
}

// --- seed ---

func TestSeedRendersKV(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/bor/spans/seed/5982": {body: loadFixture(t, "rest", "bor_spans_seed.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "seed", "5982")
	if err != nil {
		t.Fatalf("seed: %v", err)
	}
	mustContain(t, stdout, "seed")
	mustContain(t, stdout, "seed_author")
	mustContain(t, stdout, "0x3d4a70bfe707923a644449b661a5a89fb84dcd787a0811212a5ab874666003f8")
}

// --- votes (both arities) ---

func TestVotesAllIsJSON(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/bor/producer-votes": {body: loadFixture(t, "rest", "bor_producer_votes.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "votes")
	if err != nil {
		t.Fatalf("votes: %v", err)
	}
	var v any
	if jerr := json.Unmarshal([]byte(stdout), &v); jerr != nil {
		t.Fatalf("votes output not JSON: %v\n%s", jerr, stdout)
	}
	mustContain(t, stdout, "all_votes")
}

func TestVotesByID(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/bor/producer-votes/4": {body: loadFixture(t, "rest", "bor_producer_votes_by_id.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "votes", "4")
	if err != nil {
		t.Fatalf("votes 4: %v", err)
	}
	mustContain(t, stdout, "votes")
	// The fixture is {"votes":["4","5"]}.
	mustContain(t, stdout, "\"4\"")
}

func TestVotesInvalidIDIsUsageError(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{})
	_, _, err := runCmd(t, srv.URL, "votes", "notanumber")
	if err == nil {
		t.Fatal("expected error for non-integer voter id")
	}
	var uErr *client.UsageError
	if !errors.As(err, &uErr) {
		t.Fatalf("got %T, want *UsageError", err)
	}
}

// --- downtime ---

func TestDowntimePopulated(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/bor/producers/planned-downtime/4": {body: loadFixture(t, "rest", "bor_producers_planned_downtime.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "downtime", "4")
	if err != nil {
		t.Fatalf("downtime 4: %v", err)
	}
	mustContain(t, stdout, "start_block")
	mustContain(t, stdout, "end_block")
	mustNotContain(t, stdout, "none")
}

func TestDowntimeNotFoundPrintsNone(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/bor/producers/planned-downtime/99": {
			status: 404,
			body:   loadFixture(t, "rest", "bor_producers_planned_downtime_not_found.json"),
		},
	})
	stdout, _, err := runCmd(t, srv.URL, "downtime", "99")
	if err != nil {
		t.Fatalf("downtime 99: %v", err)
	}
	if strings.TrimSpace(stdout) != "none" {
		t.Errorf("expected exactly 'none', got %q", stdout)
	}
}

// --- scores ---

func TestScoresSortedDesc(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{
		"/bor/validator-performance-score": {body: loadFixture(t, "rest", "bor_validator_performance_score.json")},
	})
	stdout, _, err := runCmd(t, srv.URL, "scores")
	if err != nil {
		t.Fatalf("scores: %v", err)
	}
	// val_id=14 has the highest score (10694588) in the fixture.
	// It should appear on the first data line.
	lines := strings.Split(strings.TrimSpace(stdout), "\n")
	if len(lines) < 2 {
		t.Fatalf("expected at least 2 lines (header + first data), got %d:\n%s", len(lines), stdout)
	}
	first := lines[1]
	if !strings.Contains(first, "14") || !strings.Contains(first, "10694588") {
		t.Errorf("expected first data line to be val_id=14 score=10694588, got %q", first)
	}
}

// --- sortScoresDesc unit test ---

func TestSortScoresDeterministicTieBreak(t *testing.T) {
	in := map[string]string{
		"1":  "100",
		"2":  "100",
		"10": "100",
		"5":  "200",
	}
	rows := sortScoresDesc(in)
	want := []scoreRow{
		{id: "5", score: "200"},
		// Ties on 100: break on ascending numeric id.
		{id: "1", score: "100"},
		{id: "2", score: "100"},
		{id: "10", score: "100"},
	}
	if len(rows) != len(want) {
		t.Fatalf("len mismatch: got %d want %d", len(rows), len(want))
	}
	for i := range rows {
		if rows[i] != want[i] {
			t.Errorf("row %d: got %+v want %+v", i, rows[i], want[i])
		}
	}
}

// --- find: exercises the core algorithm via a fake spanFinder ---

// fakeFinder implements spanFinder with an in-memory slice of spans.
type fakeFinder struct {
	sprintDuration string
	// spans is indexed by span id (contiguous from 0).
	spans []spanRecord
	// latestID, if set, overrides len(spans)-1. Useful for exercising
	// the "advertised-by-latest-but-missing" paradox path.
	latestID int
}

func newFakeFinder(sprint string, spans []spanRecord) *fakeFinder {
	return &fakeFinder{
		sprintDuration: sprint,
		spans:          spans,
		latestID:       -1,
	}
}

func (f *fakeFinder) Params(_ context.Context) (borParamsEnvelope, error) {
	var env borParamsEnvelope
	env.Params.SprintDuration = f.sprintDuration
	env.Params.SpanDuration = "6400"
	env.Params.ProducerCount = "11"
	return env, nil
}

func (f *fakeFinder) Latest(_ context.Context) (spanRecord, error) {
	if len(f.spans) == 0 {
		return spanRecord{}, errors.New("no spans")
	}
	if f.latestID >= 0 && f.latestID < len(f.spans) {
		return f.spans[f.latestID], nil
	}
	return f.spans[len(f.spans)-1], nil
}

func (f *fakeFinder) ByID(_ context.Context, id uint64) (spanRecord, error) {
	if int(id) >= len(f.spans) {
		return spanRecord{}, &client.HTTPError{StatusCode: 404, Body: []byte("not found")}
	}
	return f.spans[id], nil
}

// buildFakeSpans builds N contiguous spans starting at startBlock with
// the given span length, each having the given producers. Producers
// are specified as (val_id, signer) pairs.
func buildFakeSpans(n int, startBlock, spanLen uint64, producers [][2]string) []spanRecord {
	ps := make([]spanProducer, len(producers))
	for i, p := range producers {
		ps[i] = spanProducer{ValID: p[0], Signer: p[1]}
	}
	out := make([]spanRecord, n)
	cur := startBlock
	for i := 0; i < n; i++ {
		out[i] = spanRecord{
			ID:                itoa(i),
			StartBlock:        itoa64(cur),
			EndBlock:          itoa64(cur + spanLen - 1),
			BorChainID:        "80002",
			SelectedProducers: ps,
		}
		cur += spanLen
	}
	return out
}

func itoa(i int) string      { return strconv.FormatUint(uint64(i), 10) }
func itoa64(u uint64) string { return strconv.FormatUint(u, 10) }

func TestSpanFind(t *testing.T) {
	// Fixture world:
	//  - sprint_duration = 16
	//  - span length     = 64 (4 sprints per span)
	//  - 3 selected producers cycling per sprint
	//  - Spans: id=0  [0..63],  id=1 [64..127],  id=2 [128..191]
	producers := [][2]string{
		{"4", "0xaaaa"},
		{"5", "0xbbbb"},
		{"6", "0xcccc"},
	}
	spans := buildFakeSpans(3, 0, 64, producers)
	f := newFakeFinder("16", spans)

	cases := []struct {
		name          string
		block         uint64
		wantSpanID    string
		wantProducer  string
		wantSprintIdx uint64
		beforeAny     bool
		afterLatest   bool
	}{
		{
			name:          "block at span start_block",
			block:         64, // start of span 1, sprint 0 -> producer index 0
			wantSpanID:    "1",
			wantProducer:  "4",
			wantSprintIdx: 0,
		},
		{
			name:          "block at span end_block",
			block:         127, // end of span 1; sprint (127-64)/16 = 3 -> 3 % 3 = 0
			wantSpanID:    "1",
			wantProducer:  "4",
			wantSprintIdx: 3,
		},
		{
			name:          "block at sprint boundary within span",
			block:         80, // span 1, sprint (80-64)/16 = 1 -> producer 5
			wantSpanID:    "1",
			wantProducer:  "5",
			wantSprintIdx: 1,
		},
		{
			name:          "block mid-sprint within span",
			block:         100, // span 1, sprint (100-64)/16 = 2 -> producer 6
			wantSpanID:    "1",
			wantProducer:  "6",
			wantSprintIdx: 2,
		},
		{
			name:          "block in span 0",
			block:         0, // span 0, sprint 0 -> producer 4
			wantSpanID:    "0",
			wantProducer:  "4",
			wantSprintIdx: 0,
		},
		{
			name:          "block in last span",
			block:         191, // span 2, sprint 3 -> producer 4
			wantSpanID:    "2",
			wantProducer:  "4",
			wantSprintIdx: 3,
		},
		{
			name:        "block after latest",
			block:       192,
			afterLatest: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out, err := runFind(context.Background(), f, c.block)
			if err != nil {
				t.Fatalf("runFind(%d): %v", c.block, err)
			}
			if c.afterLatest {
				if !out.AfterLatest {
					t.Fatalf("expected AfterLatest=true, got %+v", out)
				}
				return
			}
			if out.AfterLatest || out.BeforeAnySpan {
				t.Fatalf("unexpected edge flags: %+v", out)
			}
			if out.Span.ID != c.wantSpanID {
				t.Errorf("span id: got %q want %q", out.Span.ID, c.wantSpanID)
			}
			if out.SprintIndex != c.wantSprintIdx {
				t.Errorf("sprint index: got %d want %d", out.SprintIndex, c.wantSprintIdx)
			}
			if out.DesignatedProducer.ValID != c.wantProducer {
				t.Errorf("producer val_id: got %q want %q", out.DesignatedProducer.ValID, c.wantProducer)
			}
		})
	}
}

func TestSpanFindBeforeAnySpan(t *testing.T) {
	// Give span 0 a non-zero start_block so we can probe the
	// before-any-span branch.
	producers := [][2]string{{"4", "0xaaaa"}}
	spans := buildFakeSpans(2, 1000, 64, producers)
	f := newFakeFinder("16", spans)

	out, err := runFind(context.Background(), f, 42)
	if err != nil {
		t.Fatalf("runFind(42): %v", err)
	}
	if !out.BeforeAnySpan {
		t.Fatalf("expected BeforeAnySpan=true, got %+v", out)
	}
}

func TestSpanFindAfterLatest(t *testing.T) {
	producers := [][2]string{{"4", "0xaaaa"}}
	spans := buildFakeSpans(1, 0, 64, producers)
	f := newFakeFinder("16", spans)

	out, err := runFind(context.Background(), f, 999999)
	if err != nil {
		t.Fatalf("runFind(999999): %v", err)
	}
	if !out.AfterLatest {
		t.Fatalf("expected AfterLatest=true, got %+v", out)
	}
	if out.LatestEndBlock != 63 {
		t.Errorf("LatestEndBlock: got %d want 63", out.LatestEndBlock)
	}
}

func TestSpanFindSpanWithNoProducers(t *testing.T) {
	// Span covers the block but has no selected_producers.
	spans := []spanRecord{
		{
			ID:                "0",
			StartBlock:        "0",
			EndBlock:          "63",
			BorChainID:        "80002",
			SelectedProducers: nil,
		},
	}
	f := newFakeFinder("16", spans)

	out, err := runFind(context.Background(), f, 10)
	if err != nil {
		t.Fatalf("runFind(10): %v", err)
	}
	if out.Span.ID != "0" {
		t.Errorf("span id: got %q want %q", out.Span.ID, "0")
	}
	if out.DesignatedProducer.ValID != "" {
		t.Errorf("expected empty producer for span with no producers, got %+v", out.DesignatedProducer)
	}
}

func TestSpanFindInvalidBlock(t *testing.T) {
	srv := newRESTFixtureServer(t, map[string]restRoute{})
	_, _, err := runCmd(t, srv.URL, "find", "notanumber")
	if err == nil {
		t.Fatal("expected error for non-integer bor block")
	}
	var uErr *client.UsageError
	if !errors.As(err, &uErr) {
		t.Fatalf("got %T, want *UsageError", err)
	}
}

func TestSpanFindCaveatOnStderr(t *testing.T) {
	// Exercise the full `find` command end-to-end with a REST server
	// that pattern-matches /bor/spans/{id} to a single fixture. The
	// fixture has id=5982 spanning [36983659..36990058], so any block
	// within that range is "covered" by whatever id the binary search
	// probes — that's fine for asserting the Veblop note reaches
	// stderr.
	spanBody := loadFixture(t, "rest", "bor_spans_by_id.json")
	paramsBody := loadFixture(t, "rest", "bor_params.json")
	latestBody := loadFixture(t, "rest", "bor_spans_latest.json")

	srv := httptestServer(t, func(path string) ([]byte, int, bool) {
		switch {
		case path == "/bor/params":
			return paramsBody, 200, true
		case path == "/bor/spans/latest":
			return latestBody, 200, true
		case strings.HasPrefix(path, "/bor/spans/"):
			// Match any /bor/spans/{integer}, serve the canned span.
			rest := strings.TrimPrefix(path, "/bor/spans/")
			if _, err := strconv.ParseUint(rest, 10, 64); err == nil {
				return spanBody, 200, true
			}
		}
		return nil, 404, false
	})

	_, stderr, err := runCmd(t, srv.URL, "find", "36983700")
	if err != nil {
		t.Fatalf("find 36983700: %v", err)
	}
	mustContain(t, stderr, "Veblop")
}
