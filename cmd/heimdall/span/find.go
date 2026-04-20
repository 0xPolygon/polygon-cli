package span

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// borParamsEnvelope is the minimal shape of GET /bor/params we need.
type borParamsEnvelope struct {
	Params struct {
		SprintDuration string `json:"sprint_duration"`
		SpanDuration   string `json:"span_duration"`
		ProducerCount  string `json:"producer_count"`
	} `json:"params"`
}

// spanProducer is the minimal shape we need from selected_producers[].
type spanProducer struct {
	ValID  string `json:"val_id"`
	Signer string `json:"signer"`
}

// spanRecord is the minimal shape we need from /bor/spans/{id} and
// /bor/spans/latest.
type spanRecord struct {
	ID                string         `json:"id"`
	StartBlock        string         `json:"start_block"`
	EndBlock          string         `json:"end_block"`
	BorChainID        string         `json:"bor_chain_id"`
	SelectedProducers []spanProducer `json:"selected_producers"`
}

// spanEnvelope wraps a spanRecord in its upstream { "span": ... } form.
type spanEnvelope struct {
	Span spanRecord `json:"span"`
}

// spanFinder abstracts span lookups so tests can substitute a fake.
// Production uses restSpanFinder (below), which hits the REST gateway.
type spanFinder interface {
	// Latest returns the current (highest-id) span.
	Latest(ctx context.Context) (spanRecord, error)
	// ByID returns the span with the given numeric id.
	ByID(ctx context.Context, id uint64) (spanRecord, error)
	// Params returns the bor module parameters.
	Params(ctx context.Context) (borParamsEnvelope, error)
}

// restSpanFinder is the production spanFinder implementation.
type restSpanFinder struct {
	rest *client.RESTClient
}

func (f *restSpanFinder) Latest(ctx context.Context) (spanRecord, error) {
	body, _, err := f.rest.Get(ctx, "/bor/spans/latest", nil)
	if err != nil {
		return spanRecord{}, err
	}
	var env spanEnvelope
	if jerr := json.Unmarshal(body, &env); jerr != nil {
		return spanRecord{}, fmt.Errorf("decoding /bor/spans/latest: %w", jerr)
	}
	return env.Span, nil
}

func (f *restSpanFinder) ByID(ctx context.Context, id uint64) (spanRecord, error) {
	body, _, err := f.rest.Get(ctx, fmt.Sprintf("/bor/spans/%d", id), nil)
	if err != nil {
		return spanRecord{}, err
	}
	var env spanEnvelope
	if jerr := json.Unmarshal(body, &env); jerr != nil {
		return spanRecord{}, fmt.Errorf("decoding /bor/spans/%d: %w", id, jerr)
	}
	return env.Span, nil
}

func (f *restSpanFinder) Params(ctx context.Context) (borParamsEnvelope, error) {
	body, _, err := f.rest.Get(ctx, "/bor/params", nil)
	if err != nil {
		return borParamsEnvelope{}, err
	}
	var env borParamsEnvelope
	if jerr := json.Unmarshal(body, &env); jerr != nil {
		return borParamsEnvelope{}, fmt.Errorf("decoding /bor/params: %w", jerr)
	}
	return env, nil
}

// newFindCmd builds `span find <BOR_BLOCK>` — the most-requested
// operator query. Purely client-side: fetch bor params + enough spans
// (binary search) to locate the covering span, then compute the
// designated sprint producer via
// `(block - span.start_block) / sprint_duration mod len(selected_producers)`.
func newFindCmd() *cobra.Command {
	var fields []string
	cmd := &cobra.Command{
		Use:   "find <BOR_BLOCK>",
		Short: "Find span covering a Bor block.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			block, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return &client.UsageError{Msg: fmt.Sprintf("bor block must be a non-negative integer, got %q", args[0])}
			}
			rest, cfg, err := newRESTClient(cmd)
			if err != nil {
				return err
			}
			finder := &restSpanFinder{rest: rest}
			out, err := runFind(cmd.Context(), finder, block)
			if err != nil {
				return err
			}
			opts := renderOpts(cmd, cfg, fields)
			if err := renderFindResult(cmd, opts, out); err != nil {
				return err
			}
			// Print the Veblop caveat on stderr so it doesn't pollute
			// structured output.
			_, _ = fmt.Fprintln(cmd.ErrOrStderr(), veblopCaveat)
			return nil
		},
	}
	cmd.Flags().StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable)")
	return cmd
}

const veblopCaveat = "note: post-Rio, the actual block producer may differ from this designated one because Veblop rotates producers based on performance scores and planned downtime; query the Bor block itself for the on-chain author."

// findResult is what `span find` reports. Designated* fields are
// populated when the span has at least one selected producer.
type findResult struct {
	Block              uint64
	Span               spanRecord
	SprintDuration     uint64
	SprintIndex        uint64 // sprint number within the span
	ProducerIndex      int    // 0-based index into selected_producers
	DesignatedProducer spanProducer
	// BeforeAnySpan is set true when the block is below the earliest
	// span (id=0). Span / ProducerIndex etc. are zero-valued.
	BeforeAnySpan bool
	// AfterLatest is set true when the block is beyond the latest
	// span's end_block.
	AfterLatest bool
	// LatestEndBlock is populated on AfterLatest so the caller can
	// explain how far off the block is.
	LatestEndBlock uint64
}

// runFind is the core algorithm, decoupled from the REST client so
// TestSpanFind can exercise it with handcrafted fixtures.
func runFind(ctx context.Context, f spanFinder, block uint64) (findResult, error) {
	params, err := f.Params(ctx)
	if err != nil {
		return findResult{}, fmt.Errorf("fetching bor params: %w", err)
	}
	sprint, err := strconv.ParseUint(params.Params.SprintDuration, 10, 64)
	if err != nil || sprint == 0 {
		return findResult{}, fmt.Errorf("invalid sprint_duration in bor params: %q", params.Params.SprintDuration)
	}

	latest, err := f.Latest(ctx)
	if err != nil {
		return findResult{}, fmt.Errorf("fetching latest span: %w", err)
	}
	latestEnd, err := strconv.ParseUint(latest.EndBlock, 10, 64)
	if err != nil {
		return findResult{}, fmt.Errorf("invalid end_block on latest span: %q", latest.EndBlock)
	}
	if block > latestEnd {
		return findResult{AfterLatest: true, LatestEndBlock: latestEnd, Block: block, SprintDuration: sprint}, nil
	}

	latestID, err := strconv.ParseUint(latest.ID, 10, 64)
	if err != nil {
		return findResult{}, fmt.Errorf("invalid id on latest span: %q", latest.ID)
	}

	// Span 0 is the genesis span. Heimdall always has id=0 as the
	// earliest span; a block below span-0's start_block (normally 0
	// or 1) is before any span.
	span0, err := f.ByID(ctx, 0)
	if err != nil {
		// If span 0 isn't present, the next-lowest span we can fetch
		// is the latest-1-steps-back bound. Fall back to binary search
		// without it; we'll still detect before-any-span via the first
		// probe.
		span0 = spanRecord{}
	}
	if span0.StartBlock != "" {
		start0, perr := strconv.ParseUint(span0.StartBlock, 10, 64)
		if perr == nil && block < start0 {
			return findResult{BeforeAnySpan: true, Block: block, SprintDuration: sprint}, nil
		}
	}

	// Binary search: span ids are contiguous [0, latestID] and sorted
	// by start_block ascending. Probe by id; at each step fetch the
	// midpoint span and decide which half to descend into.
	lo := uint64(0)
	hi := latestID
	var covering spanRecord
	found := false
	for lo <= hi {
		mid := lo + (hi-lo)/2
		s, ferr := f.ByID(ctx, mid)
		if ferr != nil {
			// Some intermediate ids may be missing on historical
			// networks. Treat missing spans as "too low" (try the
			// upper half) — if we exhaust the range, we report
			// BeforeAnySpan.
			var hErr *client.HTTPError
			if errors.As(ferr, &hErr) && hErr.NotFound() {
				if mid == latestID {
					// Paradoxical: latest advertises this id but /bor/spans/{id}
					// 404s. Give up.
					return findResult{}, fmt.Errorf("span id %d advertised by /latest but missing from /bor/spans/%d", mid, mid)
				}
				lo = mid + 1
				continue
			}
			return findResult{}, fmt.Errorf("fetching span %d: %w", mid, ferr)
		}
		sStart, perr := strconv.ParseUint(s.StartBlock, 10, 64)
		if perr != nil {
			return findResult{}, fmt.Errorf("invalid start_block on span %d: %q", mid, s.StartBlock)
		}
		sEnd, perr := strconv.ParseUint(s.EndBlock, 10, 64)
		if perr != nil {
			return findResult{}, fmt.Errorf("invalid end_block on span %d: %q", mid, s.EndBlock)
		}
		switch {
		case block < sStart:
			if mid == 0 {
				return findResult{BeforeAnySpan: true, Block: block, SprintDuration: sprint}, nil
			}
			hi = mid - 1
		case block > sEnd:
			lo = mid + 1
		default:
			covering = s
			found = true
			// Found it; break out of the loop by forcing lo > hi.
			lo = hi + 1
		}
	}
	if !found {
		// Should be unreachable — block <= latestEnd and span 0 covers
		// or is below block. But don't crash if upstream lies.
		return findResult{}, fmt.Errorf("no span covers bor block %d (searched ids 0..%d)", block, latestID)
	}

	// Compute the designated producer.
	start, _ := strconv.ParseUint(covering.StartBlock, 10, 64)
	producers := covering.SelectedProducers
	if len(producers) == 0 {
		return findResult{
			Block:          block,
			Span:           covering,
			SprintDuration: sprint,
		}, nil
	}
	sprintIdx := (block - start) / sprint
	producerIdx := int(sprintIdx % uint64(len(producers)))
	return findResult{
		Block:              block,
		Span:               covering,
		SprintDuration:     sprint,
		SprintIndex:        sprintIdx,
		ProducerIndex:      producerIdx,
		DesignatedProducer: producers[producerIdx],
	}, nil
}

func renderFindResult(cmd *cobra.Command, opts render.Options, r findResult) error {
	if r.BeforeAnySpan {
		out := map[string]any{
			"block":  r.Block,
			"result": "before any span",
		}
		if opts.JSON {
			return render.RenderJSON(cmd.OutOrStdout(), out, opts)
		}
		_, err := fmt.Fprintf(cmd.OutOrStdout(), "bor block %d is before any known span\n", r.Block)
		return err
	}
	if r.AfterLatest {
		out := map[string]any{
			"block":            r.Block,
			"latest_end_block": r.LatestEndBlock,
			"result":           "after latest span",
		}
		if opts.JSON {
			return render.RenderJSON(cmd.OutOrStdout(), out, opts)
		}
		_, err := fmt.Fprintf(cmd.OutOrStdout(), "bor block %d is past the latest known span (end_block=%d)\n", r.Block, r.LatestEndBlock)
		return err
	}

	summary := map[string]any{
		"block":           r.Block,
		"span_id":         r.Span.ID,
		"start_block":     r.Span.StartBlock,
		"end_block":       r.Span.EndBlock,
		"sprint_duration": r.SprintDuration,
		"sprint_index":    r.SprintIndex,
		"bor_chain_id":    r.Span.BorChainID,
	}
	if len(r.Span.SelectedProducers) == 0 {
		summary["designated_producer"] = "(span has no selected_producers)"
	} else {
		summary["designated_producer_val_id"] = r.DesignatedProducer.ValID
		summary["designated_producer_signer"] = r.DesignatedProducer.Signer
		summary["producer_index"] = r.ProducerIndex
		summary["producer_count"] = len(r.Span.SelectedProducers)
	}
	if opts.JSON {
		return render.RenderJSON(cmd.OutOrStdout(), summary, opts)
	}
	return render.RenderKV(cmd.OutOrStdout(), summary, opts)
}
