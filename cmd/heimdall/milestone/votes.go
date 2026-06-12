package milestone

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/cmdutil"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/comet"
	hproto "github.com/0xPolygon/polygon-cli/internal/heimdall/proto"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

const (
	// defaultVoteRange is the number of vote heights scanned when no
	// range flags are given.
	defaultVoteRange = 1000
	// votesFetchRetries is the per-height attempt budget; any height
	// that exhausts it aborts the whole scan (all-or-nothing so the
	// output is deterministic).
	votesFetchRetries = 3
	// votesProgressEvery controls stderr progress logging cadence.
	votesProgressEvery = 500
)

// newVotesCmd builds `milestone votes`: scan a range of heimdall vote
// heights, decode each carrier block's injected ExtendedCommitInfo
// (txs[0]), and emit one record per (vote height, validator) showing
// that validator's milestone proposition — the raw material for
// spotting validators that are late or missing.
//
// Range semantics: --from/--to are VOTE heights H. The votes for
// height H travel in block H+1 (the proposer of H+1 injects the
// previous height's extended commit as the special tx at index 0), so
// the command fetches /block and /block_results at H+1.
func newVotesCmd() *cobra.Command {
	var (
		fromArg     string
		toArg       string
		fromTime    string
		toTime      string
		concurrency int
		summary     bool
		missingOnly bool
		fields      []string
	)
	cmd := &cobra.Command{
		Use:   "votes",
		Short: "Dump per-validator milestone votes over a height range.",
		Long: strings.TrimSpace(`
Scan a range of heimdall heights and report, for every validator at every
height, whether it signed the commit and which milestone proposition (bor
block range) its vote extension carried. Heights where a milestone was
finalized are correlated so late or missing validators stand out.

--from/--to are vote heights H; the data is read from block H+1, which
carries the previous height's vote extensions as its first transaction.
`),
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if fromArg != "" && fromTime != "" {
				return &client.UsageError{Msg: "--from and --from-time are mutually exclusive"}
			}
			if toArg != "" && toTime != "" {
				return &client.UsageError{Msg: "--to and --to-time are mutually exclusive"}
			}
			if concurrency < 1 {
				return &client.UsageError{Msg: fmt.Sprintf("--concurrency must be >= 1, got %d", concurrency)}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			rpc, cfg, err := pkg.RPCClient(cmd)
			if err != nil {
				return err
			}
			ctx := cmd.Context()

			if cfg.Curl {
				// Print one representative request pair and stop; the
				// scan itself issues the same calls per height.
				_, _, _ = comet.FetchBlock(ctx, rpc, "")
				_, _ = comet.FetchBlockResults(ctx, rpc, 0)
				return nil
			}

			from, to, err := resolveVoteRange(ctx, rpc, fromArg, toArg, fromTime, toTime)
			if err != nil {
				return err
			}

			rest, _, err := pkg.RESTClient(cmd)
			if err != nil {
				return err
			}
			valIDs, err := fetchValIDMap(ctx, rest)
			if err != nil {
				return err
			}

			results, err := scanVoteRange(ctx, rpc, from, to, concurrency)
			if err != nil {
				return err
			}

			records, milestones := buildVoteRecords(results, valIDs)
			if missingOnly {
				records = filterMissing(records)
			}

			opts := cmdutil.RenderOpts(cmd, cfg, fields)
			out := cmd.OutOrStdout()
			if cfg.JSON {
				env := map[string]any{
					"from":       from,
					"to":         to,
					"votes":      voteRecordMaps(records, true),
					"milestones": milestoneMaps(milestones),
				}
				if summary {
					env["summary"] = summarizeVotes(records, true)
				}
				return render.RenderJSON(out, env, opts)
			}
			if summary {
				return render.RenderTable(out, summarizeVotes(records, false), opts)
			}
			if err := render.RenderTable(out, voteRecordMaps(records, false), opts); err != nil {
				return err
			}
			if len(milestones) > 0 {
				if _, err := fmt.Fprintf(out, "\nmilestones finalized in range:\n"); err != nil {
					return err
				}
				return render.RenderTable(out, milestoneMaps(milestones), opts)
			}
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVar(&fromArg, "from", "", "first vote height to scan")
	f.StringVar(&toArg, "to", "", "last vote height to scan")
	f.StringVar(&fromTime, "from-time", "", "start of scan as unix seconds or RFC3339 (resolved to a height)")
	f.StringVar(&toTime, "to-time", "", "end of scan as unix seconds or RFC3339 (resolved to a height)")
	f.IntVar(&concurrency, "concurrency", 8, "number of heights fetched in parallel")
	f.BoolVar(&summary, "summary", false, "aggregate to one row per validator")
	f.BoolVar(&missingOnly, "missing-only", false, "only show votes that did not commit or carried no proposition")
	f.StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable)")
	return cmd
}

// resolveVoteRange turns the from/to flag values (heights or
// timestamps) into a concrete closed vote-height range, clamped to
// what the node can serve: block H+1 must exist, so vote heights run
// from max(1, earliest-1) to latest-1.
func resolveVoteRange(ctx context.Context, rpc *client.RPCClient, fromArg, toArg, fromTime, toTime string) (int64, int64, error) {
	st, err := comet.FetchStatus(ctx, rpc)
	if err != nil {
		return 0, 0, err
	}
	earliest, err := strconv.ParseInt(st.SyncInfo.EarliestBlockHeight, 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("parsing earliest height %q: %w", st.SyncInfo.EarliestBlockHeight, err)
	}
	latest, err := strconv.ParseInt(st.SyncInfo.LatestBlockHeight, 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("parsing latest height %q: %w", st.SyncInfo.LatestBlockHeight, err)
	}
	minVote := earliest - 1
	if minVote < 1 {
		minVote = 1
	}
	maxVote := latest - 1
	if maxVote < minVote {
		return 0, 0, fmt.Errorf("node has no scannable heights (earliest %d, latest %d)", earliest, latest)
	}

	resolve := func(label, heightArg, timeArg string) (int64, bool, error) {
		if heightArg != "" {
			h, perr := strconv.ParseInt(heightArg, 10, 64)
			if perr != nil || h < 1 {
				return 0, false, &client.UsageError{Msg: fmt.Sprintf("invalid %s height %q (want positive integer)", label, heightArg)}
			}
			return h, true, nil
		}
		if timeArg != "" {
			target, perr := comet.ParseTimestamp(timeArg)
			if perr != nil {
				return 0, false, &client.UsageError{Msg: perr.Error()}
			}
			h, ferr := comet.FindBlockAt(ctx, rpc, earliest, latest, target)
			if ferr != nil {
				return 0, false, ferr
			}
			// FindBlockAt returns the block closest to the timestamp;
			// the votes assembled at that moment are for the previous
			// height.
			return h - 1, true, nil
		}
		return 0, false, nil
	}

	to, toSet, err := resolve("--to", toArg, toTime)
	if err != nil {
		return 0, 0, err
	}
	if !toSet {
		to = maxVote
	}
	from, fromSet, err := resolve("--from", fromArg, fromTime)
	if err != nil {
		return 0, 0, err
	}
	if !fromSet {
		from = to - (defaultVoteRange - 1)
	}

	if from < minVote {
		from = minVote
	}
	if to > maxVote {
		log.Warn().Int64("to", to).Int64("max", maxVote).Msg("clamping --to to the latest available vote height")
		to = maxVote
	}
	if from > to {
		return 0, 0, &client.UsageError{Msg: fmt.Sprintf("empty range: from %d > to %d", from, to)}
	}
	return from, to, nil
}

// voteHeightResult is the decoded payload for one vote height.
type voteHeightResult struct {
	voteHeight int64
	time       string
	round      int32
	votes      []hproto.ExtendedVoteInfo
	milestones []milestoneEvent
}

// milestoneEvent is one `milestone` event from /block_results.
type milestoneEvent struct {
	finalizedAt int64
	number      string
	proposer    string
	startBlock  uint64
	endBlock    uint64
	hash        string
	milestoneID string
}

// scanVoteRange fetches and decodes every vote height in [from, to].
// Heights are fetched concurrently but the result slice is ordered;
// any height that fails after retries aborts the scan.
func scanVoteRange(ctx context.Context, rpc *client.RPCClient, from, to int64, concurrency int) ([]voteHeightResult, error) {
	total := to - from + 1
	results := make([]voteHeightResult, total)
	var done atomic.Int64

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(concurrency)
	for h := from; h <= to; h++ {
		h := h
		g.Go(func() error {
			res, err := fetchVoteHeight(ctx, rpc, h)
			if err != nil {
				return fmt.Errorf("vote height %d: %w", h, err)
			}
			results[h-from] = *res
			if n := done.Add(1); total > votesProgressEvery && n%votesProgressEvery == 0 {
				log.Info().Int64("done", n).Int64("total", total).Msg("scanning vote heights")
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return results, nil
}

// fetchVoteHeight retrieves block and block_results at voteHeight+1
// and decodes the extended commit. Each RPC is retried before giving
// up.
func fetchVoteHeight(ctx context.Context, rpc *client.RPCClient, voteHeight int64) (*voteHeightResult, error) {
	carrier := voteHeight + 1

	var blk *comet.Block
	err := withRetry(ctx, func() error {
		b, raw, ferr := comet.FetchBlock(ctx, rpc, strconv.FormatInt(carrier, 10))
		if ferr != nil {
			return ferr
		}
		if raw == nil {
			return fmt.Errorf("empty block response")
		}
		blk = b
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(blk.Block.Data.Txs) == 0 {
		return nil, fmt.Errorf("block %d has no transactions; vote extensions may not be enabled at this height", carrier)
	}
	rawTx, err := base64.StdEncoding.DecodeString(blk.Block.Data.Txs[0])
	if err != nil {
		return nil, fmt.Errorf("decoding txs[0] of block %d: %w", carrier, err)
	}
	extCommit, err := hproto.UnmarshalExtendedCommitInfo(rawTx)
	if err != nil {
		return nil, fmt.Errorf("decoding extended commit in block %d (vote extensions may not be enabled at this height): %w", carrier, err)
	}

	var br *comet.BlockResults
	err = withRetry(ctx, func() error {
		r, ferr := comet.FetchBlockResults(ctx, rpc, carrier)
		if ferr != nil {
			return ferr
		}
		if r == nil {
			return fmt.Errorf("empty block_results response")
		}
		br = r
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &voteHeightResult{
		voteHeight: voteHeight,
		time:       blk.Block.Header.Time,
		round:      extCommit.Round,
		votes:      extCommit.Votes,
		milestones: parseMilestoneEvents(br, carrier),
	}, nil
}

// withRetry runs fn up to votesFetchRetries times with linear backoff,
// honouring context cancellation between attempts.
func withRetry(ctx context.Context, fn func() error) error {
	var err error
	for attempt := 1; attempt <= votesFetchRetries; attempt++ {
		if err = fn(); err == nil {
			return nil
		}
		if attempt == votesFetchRetries {
			break
		}
		timer := time.NewTimer(time.Duration(attempt) * 250 * time.Millisecond)
		select {
		case <-timer.C:
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		}
		timer.Stop()
	}
	return fmt.Errorf("failed after %d attempts: %w", votesFetchRetries, err)
}

// parseMilestoneEvents extracts `milestone` events from block_results.
func parseMilestoneEvents(br *comet.BlockResults, carrier int64) []milestoneEvent {
	var out []milestoneEvent
	for _, ev := range br.FinalizeBlockEvents {
		if ev.Type != "milestone" {
			continue
		}
		attrs := map[string]string{}
		for _, a := range ev.Attributes {
			attrs[a.Key] = a.Value
		}
		start, _ := strconv.ParseUint(attrs["start_block"], 10, 64)
		end, _ := strconv.ParseUint(attrs["end_block"], 10, 64)
		out = append(out, milestoneEvent{
			finalizedAt: carrier,
			number:      attrs["number"],
			proposer:    attrs["proposer"],
			startBlock:  start,
			endBlock:    end,
			hash:        attrs["hash"],
			milestoneID: attrs["milestone_id"],
		})
	}
	return out
}

// voteRecord is one (vote height, validator) output row.
type voteRecord struct {
	voteHeight int64
	time       string
	round      int32
	address    string // 0x-prefixed lowercase hex
	valID      string
	power      int64
	flag       hproto.BlockIDFlag
	prop       *propInfo
	// majorityEnd is the highest bor block covered by >=2/3 of the
	// voting power at this height; nil when no majority exists.
	majorityEnd *uint64
	// lag = majorityEnd - prop end; nil without majority or proposition.
	lag *int64
	// milestone correlation, set on heights where a milestone finalized.
	msNumber   string
	msCovered  bool
	msRelevant bool // validator committed with a proposition at a finalize height
}

type propInfo struct {
	startBlock uint64
	endBlock   uint64
	endHash    string
	numHashes  int
	parentHash string
}

// validatorInfo is the slice of /stake/validators-set we keep.
type validatorInfo struct {
	ValID string
}

// fetchValIDMap fetches the current validator set once and returns a
// map from 0x-prefixed lowercase signer address to validator info.
// The mapping reflects the CURRENT set; validators rotated out since
// the scanned heights render with val_id "-".
func fetchValIDMap(ctx context.Context, rest *client.RESTClient) (map[string]validatorInfo, error) {
	body, status, err := rest.Get(ctx, "/stake/validators-set", nil)
	if err != nil {
		return nil, fmt.Errorf("fetching validator set: %w", err)
	}
	if status == 0 && body == nil {
		return map[string]validatorInfo{}, nil // --curl
	}
	m, err := cmdutil.DecodeJSONMap(body, "validator set")
	if err != nil {
		return nil, err
	}
	out := map[string]validatorInfo{}
	set, _ := m["validator_set"].(map[string]any)
	vals, _ := set["validators"].([]any)
	for _, v := range vals {
		vm, ok := v.(map[string]any)
		if !ok {
			continue
		}
		signer, _ := vm["signer"].(string)
		id, _ := vm["val_id"].(string)
		if signer == "" {
			continue
		}
		out[strings.ToLower(signer)] = validatorInfo{ValID: id}
	}
	return out, nil
}

// buildVoteRecords flattens scan results into per-validator records
// and the list of milestones finalized in the range.
func buildVoteRecords(results []voteHeightResult, valIDs map[string]validatorInfo) ([]voteRecord, []milestoneEvent) {
	var records []voteRecord
	var milestones []milestoneEvent
	for _, res := range results {
		milestones = append(milestones, res.milestones...)

		decoded := make([]*propInfo, len(res.votes))
		var totalPower int64
		for i, v := range res.votes {
			totalPower += v.Validator.Power
			if v.BlockIDFlag != hproto.BlockIDFlagCommit || len(v.VoteExtension) == 0 {
				continue
			}
			ve, err := hproto.UnmarshalVoteExtension(v.VoteExtension)
			if err != nil || ve.MilestoneProposition == nil {
				continue
			}
			mp := ve.MilestoneProposition
			if len(mp.BlockHashes) == 0 {
				continue
			}
			decoded[i] = &propInfo{
				startBlock: mp.StartBlockNumber,
				endBlock:   mp.StartBlockNumber + uint64(len(mp.BlockHashes)) - 1,
				endHash:    "0x" + hex.EncodeToString(mp.BlockHashes[len(mp.BlockHashes)-1]),
				numHashes:  len(mp.BlockHashes),
				parentHash: "0x" + hex.EncodeToString(mp.ParentHash),
			}
		}

		majorityEnd, hasMajority := majorityEndBlock(res.votes, decoded, totalPower)

		// At most one milestone normally finalizes per height; if
		// several did, correlate against the last (highest) one.
		var ms *milestoneEvent
		if len(res.milestones) > 0 {
			ms = &res.milestones[len(res.milestones)-1]
		}

		for i, v := range res.votes {
			rec := voteRecord{
				voteHeight: res.voteHeight,
				time:       res.time,
				round:      res.round,
				address:    "0x" + strings.ToLower(hex.EncodeToString(v.Validator.Address)),
				power:      v.Validator.Power,
				flag:       v.BlockIDFlag,
				prop:       decoded[i],
			}
			if info, ok := valIDs[rec.address]; ok {
				rec.valID = info.ValID
			}
			if hasMajority {
				me := majorityEnd
				rec.majorityEnd = &me
				if rec.prop != nil {
					lag := int64(majorityEnd) - int64(rec.prop.endBlock)
					rec.lag = &lag
				}
			}
			if ms != nil && v.BlockIDFlag == hproto.BlockIDFlagCommit && rec.prop != nil {
				rec.msRelevant = true
				rec.msNumber = ms.number
				rec.msCovered = rec.prop.startBlock <= ms.endBlock && ms.endBlock <= rec.prop.endBlock
			}
			records = append(records, rec)
		}
	}
	sortVoteRecords(records)
	return records, milestones
}

// sortVoteRecords orders by vote height ascending, then power
// descending, then address for determinism.
func sortVoteRecords(records []voteRecord) {
	sort.Slice(records, func(i, j int) bool {
		a, b := records[i], records[j]
		if a.voteHeight != b.voteHeight {
			return a.voteHeight < b.voteHeight
		}
		if a.power != b.power {
			return a.power > b.power
		}
		return a.address < b.address
	})
}

// majorityEndBlock returns the highest bor block number covered by
// >= 2/3+1 of the total voting power, mirroring heimdall's
// majorityVP := totalVotingPower*2/3 + 1 with only COMMIT votes
// counted. The maximum such block is always some proposition's end
// block, so only those are candidates.
func majorityEndBlock(votes []hproto.ExtendedVoteInfo, decoded []*propInfo, totalPower int64) (uint64, bool) {
	if totalPower <= 0 {
		return 0, false
	}
	threshold := totalPower*2/3 + 1

	candidates := map[uint64]bool{}
	for _, p := range decoded {
		if p != nil {
			candidates[p.endBlock] = true
		}
	}
	ends := make([]uint64, 0, len(candidates))
	for e := range candidates {
		ends = append(ends, e)
	}
	sort.Slice(ends, func(i, j int) bool { return ends[i] > ends[j] })

	for _, end := range ends {
		var covered int64
		for i, v := range votes {
			p := decoded[i]
			if p == nil || v.BlockIDFlag != hproto.BlockIDFlagCommit {
				continue
			}
			if p.startBlock <= end && end <= p.endBlock {
				covered += v.Validator.Power
			}
		}
		if covered >= threshold {
			return end, true
		}
	}
	return 0, false
}

// filterMissing keeps only records that signal a problem: the
// validator did not commit, or committed without a proposition.
func filterMissing(records []voteRecord) []voteRecord {
	out := make([]voteRecord, 0, len(records))
	for _, r := range records {
		if r.flag != hproto.BlockIDFlagCommit || r.prop == nil {
			out = append(out, r)
		}
	}
	return out
}

// voteRecordMaps converts records to the render-friendly form shared
// by the table and JSON outputs. Missing values render as "-" in
// tables and null in JSON so JSON consumers see typed columns.
func voteRecordMaps(records []voteRecord, jsonOut bool) []map[string]any {
	blank := func() any {
		if jsonOut {
			return nil
		}
		return "-"
	}
	out := make([]map[string]any, 0, len(records))
	for _, r := range records {
		m := map[string]any{
			"height": r.voteHeight,
			"time":   r.time,
			"val_id": dash(r.valID),
			"signer": r.address,
			"power":  r.power,
			"flag":   r.flag.String(),
		}
		if r.prop != nil {
			m["prop_start"] = r.prop.startBlock
			m["prop_end"] = r.prop.endBlock
			m["end_hash"] = r.prop.endHash
			m["n_hashes"] = r.prop.numHashes
		} else {
			m["prop_start"], m["prop_end"], m["end_hash"], m["n_hashes"] = blank(), blank(), blank(), blank()
		}
		if r.lag != nil {
			m["lag"] = *r.lag
		} else {
			m["lag"] = blank()
		}
		switch {
		case r.msRelevant && r.msCovered:
			m["milestone"] = r.msNumber
		case r.msRelevant:
			m["milestone"] = "miss"
		default:
			m["milestone"] = blank()
		}
		out = append(out, m)
	}
	return out
}

// milestoneMaps converts milestone events to the render-friendly form.
func milestoneMaps(milestones []milestoneEvent) []map[string]any {
	out := make([]map[string]any, 0, len(milestones))
	for _, ms := range milestones {
		out = append(out, map[string]any{
			"vote_height":  ms.finalizedAt - 1,
			"finalized_at": ms.finalizedAt,
			"number":       ms.number,
			"proposer":     ms.proposer,
			"start_block":  ms.startBlock,
			"end_block":    ms.endBlock,
			"hash":         "0x" + strings.TrimPrefix(ms.hash, "0x"),
			"milestone_id": ms.milestoneID,
		})
	}
	return out
}

// summarizeVotes aggregates records to one row per validator:
// commit/miss counts, proposition coverage, and lag statistics.
func summarizeVotes(records []voteRecord, jsonOut bool) []map[string]any {
	type agg struct {
		valID    string
		power    int64
		expected int
		signed   int
		noProp   int
		lagSum   int64
		lagCount int64
		maxLag   int64
		covered  int
		msTotal  int
	}
	byAddr := map[string]*agg{}
	var order []string
	for _, r := range records {
		a, ok := byAddr[r.address]
		if !ok {
			a = &agg{valID: r.valID}
			byAddr[r.address] = a
			order = append(order, r.address)
		}
		a.power = r.power
		a.expected++
		if r.flag == hproto.BlockIDFlagCommit {
			a.signed++
			if r.prop == nil {
				a.noProp++
			}
		}
		if r.lag != nil {
			a.lagSum += *r.lag
			a.lagCount++
			if *r.lag > a.maxLag {
				a.maxLag = *r.lag
			}
		}
		if r.msRelevant {
			a.msTotal++
			if r.msCovered {
				a.covered++
			}
		}
	}
	sort.Slice(order, func(i, j int) bool {
		if byAddr[order[i]].power != byAddr[order[j]].power {
			return byAddr[order[i]].power > byAddr[order[j]].power
		}
		return order[i] < order[j]
	})
	out := make([]map[string]any, 0, len(order))
	for _, addr := range order {
		a := byAddr[addr]
		m := map[string]any{
			"val_id":     dash(a.valID),
			"signer":     addr,
			"power":      a.power,
			"expected":   a.expected,
			"signed":     a.signed,
			"missed":     a.expected - a.signed,
			"no_prop":    a.noProp,
			"ms_covered": a.covered,
			"ms_total":   a.msTotal,
		}
		if a.lagCount > 0 {
			avg := fmt.Sprintf("%.2f", float64(a.lagSum)/float64(a.lagCount))
			if avg == "-0.00" {
				avg = "0.00"
			}
			m["avg_lag"] = avg
			m["max_lag"] = a.maxLag
		} else if jsonOut {
			m["avg_lag"], m["max_lag"] = nil, nil
		} else {
			m["avg_lag"], m["max_lag"] = "-", "-"
		}
		out = append(out, m)
	}
	return out
}

func dash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}
