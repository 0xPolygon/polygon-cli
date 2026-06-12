package milestone

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	hproto "github.com/0xPolygon/polygon-cli/internal/heimdall/proto"
)

// Test fleet: A and B match signers in the stake_validators_set.json
// fixture (val_id 16 and 13); C is not in the set.
var (
	valA = mustHexAddr("02f615e95563ef16f10354dba9e584e58d2d4314") // power 100
	valB = mustHexAddr("04ba3ef4c023c1006019a0f9baf6e70455e41fcf") // power 10
	valC = mustHexAddr("1111111111111111111111111111111111111111") // power 20
)

func mustHexAddr(s string) []byte {
	out := make([]byte, len(s)/2)
	if _, err := fmt.Sscanf(s, "%40x", &out); err != nil {
		panic(err)
	}
	return out
}

// makeProp builds a MilestoneProposition covering [start, start+n-1].
func makeProp(start uint64, n int) *hproto.MilestoneProposition {
	hashes := make([][]byte, n)
	tds := make([]uint64, n)
	for i := range hashes {
		h := make([]byte, 32)
		h[0] = byte(start>>8) + byte(i)
		h[31] = byte(start) + byte(i)
		hashes[i] = h
		tds[i] = start + uint64(i)
	}
	return &hproto.MilestoneProposition{
		BlockHashes:      hashes,
		StartBlockNumber: start,
		ParentHash:       make([]byte, 32),
		BlockTDs:         tds,
	}
}

// makeExtCommit builds the txs[0] payload for one carrier block.
// A commits with propA (nil for no proposition), B with propB, C is
// always ABSENT.
func makeExtCommit(voteHeight int64, propA, propB *hproto.MilestoneProposition) string {
	mkVE := func(prop *hproto.MilestoneProposition) []byte {
		ve := &hproto.VoteExtension{
			BlockHash:            make([]byte, 32),
			Height:               voteHeight,
			MilestoneProposition: prop,
		}
		return ve.Marshal()
	}
	ec := &hproto.ExtendedCommitInfo{
		Round: 0,
		Votes: []hproto.ExtendedVoteInfo{
			{
				Validator:          hproto.ExtValidator{Address: valA, Power: 100},
				VoteExtension:      mkVE(propA),
				ExtensionSignature: []byte("sig-a"),
				BlockIDFlag:        hproto.BlockIDFlagCommit,
			},
			{
				Validator:          hproto.ExtValidator{Address: valB, Power: 10},
				VoteExtension:      mkVE(propB),
				ExtensionSignature: []byte("sig-b"),
				BlockIDFlag:        hproto.BlockIDFlagCommit,
			},
			{
				Validator:   hproto.ExtValidator{Address: valC, Power: 20},
				BlockIDFlag: hproto.BlockIDFlagAbsent,
			},
		},
	}
	return base64.StdEncoding.EncodeToString(ec.Marshal())
}

func milestoneEventJSON(number string, start, end uint64) map[string]any {
	attr := func(k, v string) map[string]any { return map[string]any{"key": k, "value": v} }
	return map[string]any{
		"type": "milestone",
		"attributes": []any{
			attr("proposer", "0x02f615e95563ef16f10354dba9e584e58d2d4314"),
			attr("hash", "00aa"),
			attr("start_block", fmt.Sprintf("%d", start)),
			attr("end_block", fmt.Sprintf("%d", end)),
			attr("bor_chain_id", "80002"),
			attr("milestone_id", "ms-id-"+number),
			attr("timestamp", "1770000000"),
			attr("number", number),
		},
	}
}

// votesScenario describes the canned chain served by the test server:
// vote heights 29999..30002 (carrier blocks 30000..30003). Heights
// 30001/30002 carry propositions and finalize milestones 77/78.
//
//	30001: A proposes [100..102], B commits without a proposition,
//	       C absent. Milestone 77 (end 102) finalizes -> A covers.
//	30002: A proposes [103..104], B proposes a stale [100..101],
//	       C absent. Milestone 78 (end 104) finalizes -> B misses,
//	       lag 3.
func votesScenario() (map[int64]string, map[int64][]any) {
	txs := map[int64]string{
		30000: makeExtCommit(29999, makeProp(95, 2), nil),
		30001: makeExtCommit(30000, makeProp(97, 3), nil),
		30002: makeExtCommit(30001, makeProp(100, 3), nil),
		30003: makeExtCommit(30002, makeProp(103, 2), makeProp(100, 2)),
	}
	events := map[int64][]any{
		30002: {milestoneEventJSON("77", 100, 102)},
		30003: {milestoneEventJSON("78", 103, 104)},
	}
	return txs, events
}

// newVotesTestServer serves both protocols from one URL: GET requests
// hit the REST fixture routes, POST requests are treated as CometBFT
// JSON-RPC (status/block/block_results) over the canned scenario.
// blockFailures, when non-nil, maps a carrier height to the number of
// times /block should fail before succeeding (retry-path testing).
func newVotesTestServer(t *testing.T, blockFailures map[int64]*atomic.Int64) *httptest.Server {
	t.Helper()
	txs, events := votesScenario()
	valSet := loadFixture(t, "rest", "stake_validators_set.json")

	rpcResult := func(method string, height int64) (any, error) {
		switch method {
		case "status":
			return map[string]any{
				"node_info": map[string]any{"network": "heimdallv2-80002"},
				"sync_info": map[string]any{
					"earliest_block_height": "30000",
					"latest_block_height":   "30003",
					"catching_up":           false,
				},
			}, nil
		case "block":
			tx, ok := txs[height]
			if !ok {
				return nil, fmt.Errorf("height %d not in scenario", height)
			}
			return map[string]any{
				"block_id": map[string]any{"hash": "AA"},
				"block": map[string]any{
					"header": map[string]any{
						"chain_id": "heimdallv2-80002",
						"height":   fmt.Sprintf("%d", height),
						"time":     fmt.Sprintf("2026-06-11T07:00:%02dZ", height-30000),
					},
					"data": map[string]any{"txs": []any{tx}},
				},
			}, nil
		case "block_results":
			evs := events[height]
			if evs == nil {
				evs = []any{}
			}
			return map[string]any{
				"height":                fmt.Sprintf("%d", height),
				"finalize_block_events": evs,
			}, nil
		}
		return nil, fmt.Errorf("unexpected method %q", method)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if r.URL.Path != "/stake/validators-set" {
				http.Error(w, "no route for "+r.URL.Path, 404)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(valSet)
			return
		}
		var req struct {
			ID     uint64         `json:"id"`
			Method string         `json:"method"`
			Params map[string]any `json:"params"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		var height int64
		if h, ok := req.Params["height"].(string); ok {
			_, _ = fmt.Sscanf(h, "%d", &height)
		}
		if req.Method == "block" && blockFailures != nil {
			if rem, ok := blockFailures[height]; ok && rem.Add(-1) >= 0 {
				http.Error(w, "transient failure", 500)
				return
			}
		}
		result, err := rpcResult(req.Method, height)
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			_ = json.NewEncoder(w).Encode(map[string]any{
				"jsonrpc": "2.0", "id": req.ID,
				"error": map[string]any{"code": -32603, "message": err.Error()},
			})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"jsonrpc": "2.0", "id": req.ID, "result": result})
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestVotesTable(t *testing.T) {
	srv := newVotesTestServer(t, nil)
	out, _, err := runCmd(t, srv.URL, "votes", "--from", "30001", "--to", "30002")
	if err != nil {
		t.Fatalf("votes: %v", err)
	}
	// Validator A: mapped to val_id 16, covers milestone 77 at 30001.
	mustContain(t, out, "16")
	mustContain(t, out, "0x02f615e95563ef16f10354dba9e584e58d2d4314")
	mustContain(t, out, "77")
	// Validator B misses milestone 78 with lag 3 at 30002.
	mustContain(t, out, "miss")
	// Validator C is absent and unmapped.
	mustContain(t, out, "ABSENT")
	mustContain(t, out, "0x1111111111111111111111111111111111111111")
	// Milestones section appended after the votes table.
	mustContain(t, out, "milestones finalized in range:")
	mustContain(t, out, "ms-id-78")
}

func TestVotesJSON(t *testing.T) {
	srv := newVotesTestServer(t, nil)
	out, _, err := runCmd(t, srv.URL, "votes", "--from", "30001", "--to", "30002", "--json")
	if err != nil {
		t.Fatalf("votes --json: %v", err)
	}
	var env struct {
		From       int64            `json:"from"`
		To         int64            `json:"to"`
		Votes      []map[string]any `json:"votes"`
		Milestones []map[string]any `json:"milestones"`
	}
	if err := json.Unmarshal([]byte(out), &env); err != nil {
		t.Fatalf("decoding JSON output: %v\n%s", err, out)
	}
	if env.From != 30001 || env.To != 30002 {
		t.Errorf("range: got [%d, %d], want [30001, 30002]", env.From, env.To)
	}
	if len(env.Votes) != 6 {
		t.Fatalf("votes: got %d records, want 6", len(env.Votes))
	}
	if len(env.Milestones) != 2 {
		t.Fatalf("milestones: got %d, want 2", len(env.Milestones))
	}

	find := func(height float64, signer string) map[string]any {
		for _, v := range env.Votes {
			if v["height"] == height && v["signer"] == signer {
				return v
			}
		}
		t.Fatalf("no vote record for height %v signer %s", height, signer)
		return nil
	}

	a := find(30002, "0x02f615e95563ef16f10354dba9e584e58d2d4314")
	if a["lag"] != float64(0) || a["milestone"] != "78" || a["prop_end"] != float64(104) {
		t.Errorf("validator A record diverged: %v", a)
	}
	b := find(30002, "0x04ba3ef4c023c1006019a0f9baf6e70455e41fcf")
	if b["lag"] != float64(3) || b["milestone"] != "miss" || b["val_id"] != "13" {
		t.Errorf("validator B record diverged: %v", b)
	}
	c := find(30002, "0x1111111111111111111111111111111111111111")
	if c["flag"] != "ABSENT" || c["lag"] != nil || c["prop_start"] != nil || c["milestone"] != nil || c["val_id"] != "-" {
		t.Errorf("validator C record diverged: %v", c)
	}
	// B committed without a proposition at 30001: typed nulls, no
	// milestone relevance.
	b1 := find(30001, "0x04ba3ef4c023c1006019a0f9baf6e70455e41fcf")
	if b1["flag"] != "COMMIT" || b1["prop_start"] != nil || b1["milestone"] != nil {
		t.Errorf("validator B@30001 record diverged: %v", b1)
	}

	ms := env.Milestones[1]
	if ms["number"] != "78" || ms["vote_height"] != float64(30002) || ms["finalized_at"] != float64(30003) ||
		ms["end_block"] != float64(104) || ms["hash"] != "0x00aa" {
		t.Errorf("milestone record diverged: %v", ms)
	}
}

func TestVotesSummary(t *testing.T) {
	srv := newVotesTestServer(t, nil)
	out, _, err := runCmd(t, srv.URL, "votes", "--from", "30001", "--to", "30002", "--summary", "--json")
	if err != nil {
		t.Fatalf("votes --summary: %v", err)
	}
	var env struct {
		Summary []map[string]any `json:"summary"`
	}
	if err := json.Unmarshal([]byte(out), &env); err != nil {
		t.Fatalf("decoding JSON output: %v\n%s", err, out)
	}
	if len(env.Summary) != 3 {
		t.Fatalf("summary: got %d rows, want 3", len(env.Summary))
	}
	// Rows are power-descending: A (100), C (20), B (10).
	a, c, b := env.Summary[0], env.Summary[1], env.Summary[2]
	if a["val_id"] != "16" || a["signed"] != float64(2) || a["missed"] != float64(0) ||
		a["ms_covered"] != float64(2) || a["ms_total"] != float64(2) || a["max_lag"] != float64(0) {
		t.Errorf("summary A diverged: %v", a)
	}
	if c["missed"] != float64(2) || c["signed"] != float64(0) || c["val_id"] != "-" {
		t.Errorf("summary C diverged: %v", c)
	}
	if b["no_prop"] != float64(1) || b["ms_covered"] != float64(0) || b["ms_total"] != float64(1) ||
		b["max_lag"] != float64(3) || b["avg_lag"] != "3.00" {
		t.Errorf("summary B diverged: %v", b)
	}
}

func TestVotesMissingOnly(t *testing.T) {
	srv := newVotesTestServer(t, nil)
	out, _, err := runCmd(t, srv.URL, "votes", "--from", "30001", "--to", "30002", "--missing-only", "--json")
	if err != nil {
		t.Fatalf("votes --missing-only: %v", err)
	}
	var env struct {
		Votes []map[string]any `json:"votes"`
	}
	if err := json.Unmarshal([]byte(out), &env); err != nil {
		t.Fatalf("decoding JSON output: %v\n%s", err, out)
	}
	// C absent at both heights + B's no-proposition commit at 30001.
	if len(env.Votes) != 3 {
		t.Fatalf("votes: got %d records, want 3: %v", len(env.Votes), env.Votes)
	}
	for _, v := range env.Votes {
		if v["flag"] == "COMMIT" && v["prop_start"] != nil {
			t.Errorf("unexpected healthy record in --missing-only output: %v", v)
		}
	}
}

func TestVotesDefaultRangeClampsToNode(t *testing.T) {
	// No --from/--to: to = latest-1 = 30002, from = to-999 clamped to
	// earliest-1 = 29999.
	srv := newVotesTestServer(t, nil)
	out, _, err := runCmd(t, srv.URL, "votes", "--json")
	if err != nil {
		t.Fatalf("votes: %v", err)
	}
	var env struct {
		From  int64            `json:"from"`
		To    int64            `json:"to"`
		Votes []map[string]any `json:"votes"`
	}
	if err := json.Unmarshal([]byte(out), &env); err != nil {
		t.Fatalf("decoding JSON output: %v\n%s", err, out)
	}
	if env.From != 29999 || env.To != 30002 {
		t.Errorf("range: got [%d, %d], want [29999, 30002]", env.From, env.To)
	}
	if len(env.Votes) != 12 {
		t.Errorf("votes: got %d records, want 12", len(env.Votes))
	}
}

func TestVotesRetriesTransientFailures(t *testing.T) {
	// /block at carrier 30002 fails twice; the third attempt succeeds
	// within the votesFetchRetries budget.
	var failures atomic.Int64
	failures.Store(2)
	srv := newVotesTestServer(t, map[int64]*atomic.Int64{30002: &failures})
	_, _, err := runCmd(t, srv.URL, "votes", "--from", "30001", "--to", "30002")
	if err != nil {
		t.Fatalf("votes with transient failures: %v", err)
	}
}

func TestVotesFailsAfterRetryBudget(t *testing.T) {
	var failures atomic.Int64
	failures.Store(99)
	srv := newVotesTestServer(t, map[int64]*atomic.Int64{30002: &failures})
	_, _, err := runCmd(t, srv.URL, "votes", "--from", "30001", "--to", "30002")
	if err == nil {
		t.Fatal("expected error after exhausting retries")
	}
	mustContain(t, err.Error(), "vote height 30001")
	mustContain(t, err.Error(), "failed after 3 attempts")
}

func TestVotesFlagValidation(t *testing.T) {
	srv := newVotesTestServer(t, nil)
	cases := [][]string{
		{"votes", "--from", "30001", "--from-time", "2026-06-11T00:00:00Z"},
		{"votes", "--to", "30002", "--to-time", "2026-06-11T00:00:00Z"},
		{"votes", "--concurrency", "0"},
		{"votes", "--from", "abc"},
		{"votes", "--from", "30002", "--to", "30001"},
	}
	for _, args := range cases {
		if _, _, err := runCmd(t, srv.URL, args...); err == nil {
			t.Errorf("expected error for %v", args)
		}
	}
}

func TestVotesEmptyTxsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ID     uint64         `json:"id"`
			Method string         `json:"method"`
			Params map[string]any `json:"params"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)
		var result any
		switch req.Method {
		case "status":
			result = map[string]any{"sync_info": map[string]any{
				"earliest_block_height": "1", "latest_block_height": "10",
			}}
		case "block":
			result = map[string]any{"block": map[string]any{
				"header": map[string]any{"height": "6", "time": "2026-06-11T07:00:00Z"},
				"data":   map[string]any{"txs": []any{}},
			}}
		default:
			result = map[string]any{}
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"jsonrpc": "2.0", "id": req.ID, "result": result})
	}))
	t.Cleanup(srv.Close)
	_, _, err := runCmd(t, srv.URL, "votes", "--from", "5", "--to", "5")
	if err == nil {
		t.Fatal("expected error for block without transactions")
	}
	mustContain(t, err.Error(), "vote extensions may not be enabled")
}

func TestWithRetryHonoursCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	calls := 0
	err := withRetry(ctx, func() error {
		calls++
		return fmt.Errorf("always fails")
	})
	if err != context.Canceled {
		t.Errorf("got %v, want context.Canceled", err)
	}
	if calls != 1 {
		t.Errorf("fn called %d times, want 1 (no retries after cancellation)", calls)
	}
}

func TestMajorityEndBlock(t *testing.T) {
	votes := []hproto.ExtendedVoteInfo{
		{Validator: hproto.ExtValidator{Power: 100}, BlockIDFlag: hproto.BlockIDFlagCommit},
		{Validator: hproto.ExtValidator{Power: 10}, BlockIDFlag: hproto.BlockIDFlagCommit},
		{Validator: hproto.ExtValidator{Power: 20}, BlockIDFlag: hproto.BlockIDFlagAbsent},
	}
	// threshold = 130*2/3+1 = 87.
	cases := []struct {
		name    string
		decoded []*propInfo
		want    uint64
		wantOK  bool
	}{
		{
			name:    "majority from single large validator",
			decoded: []*propInfo{{startBlock: 100, endBlock: 104}, {startBlock: 100, endBlock: 101}, nil},
			want:    104, wantOK: true,
		},
		{
			name:    "no propositions",
			decoded: []*propInfo{nil, nil, nil},
			wantOK:  false,
		},
		{
			name: "absent validator's proposition does not count",
			decoded: []*propInfo{
				nil,
				{startBlock: 100, endBlock: 104},
				{startBlock: 100, endBlock: 104}, // ABSENT
			},
			wantOK: false,
		},
	}
	for _, c := range cases {
		got, ok := majorityEndBlock(votes, c.decoded, 130)
		if ok != c.wantOK || (ok && got != c.want) {
			t.Errorf("%s: got (%d, %v), want (%d, %v)", c.name, got, ok, c.want, c.wantOK)
		}
	}
}
