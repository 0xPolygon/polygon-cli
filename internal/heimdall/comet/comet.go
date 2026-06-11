// Package comet provides typed fetch helpers for the CometBFT JSON-RPC
// endpoints shared by the heimdall subcommand families (chain,
// milestone). We decode only the fields actually used; the raw JSON is
// returned alongside where callers need passthrough.
package comet

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
)

// Status is the subset of the CometBFT /status response we consume.
type Status struct {
	NodeInfo struct {
		Version string `json:"version"`
		Network string `json:"network"`
		Moniker string `json:"moniker"`
	} `json:"node_info"`
	SyncInfo struct {
		LatestBlockHeight   string `json:"latest_block_height"`
		LatestBlockTime     string `json:"latest_block_time"`
		EarliestBlockHeight string `json:"earliest_block_height"`
		CatchingUp          bool   `json:"catching_up"`
	} `json:"sync_info"`
}

// Block is the subset of the CometBFT /block response we consume.
type Block struct {
	BlockID struct {
		Hash string `json:"hash"`
	} `json:"block_id"`
	Block struct {
		Header struct {
			ChainID         string `json:"chain_id"`
			Height          string `json:"height"`
			Time            string `json:"time"`
			ProposerAddress string `json:"proposer_address"`
		} `json:"header"`
		Data struct {
			Txs []string `json:"txs"`
		} `json:"data"`
	} `json:"block"`
}

// ABCIEvent is one event from a CometBFT /block_results response.
// Attribute keys/values are plain strings on CometBFT 0.38.
type ABCIEvent struct {
	Type       string `json:"type"`
	Attributes []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"attributes"`
}

// BlockResults is the subset of /block_results we consume.
type BlockResults struct {
	Height              string      `json:"height"`
	FinalizeBlockEvents []ABCIEvent `json:"finalize_block_events"`
}

// FetchStatus calls the CometBFT /status RPC and returns the decoded
// struct. Returns (nil, nil) when running under --curl.
func FetchStatus(ctx context.Context, rpc *client.RPCClient) (*Status, error) {
	raw, err := rpc.Call(ctx, "status", nil)
	if err != nil {
		return nil, fmt.Errorf("fetching status: %w", err)
	}
	if raw == nil {
		return nil, nil
	}
	var st Status
	if err := json.Unmarshal(raw, &st); err != nil {
		return nil, fmt.Errorf("decoding status: %w", err)
	}
	return &st, nil
}

// FetchBlock calls CometBFT /block at the given height (empty ==
// latest). Returns the typed struct, the raw result bytes (for --json
// passthrough), and any error. Both return values are nil when --curl
// short-circuits.
func FetchBlock(ctx context.Context, rpc *client.RPCClient, height string) (*Block, json.RawMessage, error) {
	// CometBFT's reflect-based RPC requires an explicit `height`
	// key in params; a missing or empty params object returns
	// "reflect: Call with too few input arguments". Pass nil height
	// to request the latest block.
	params := map[string]any{"height": nil}
	if height != "" {
		params["height"] = height
	}
	raw, err := rpc.Call(ctx, "block", params)
	if err != nil {
		return nil, nil, fmt.Errorf("fetching block: %w", err)
	}
	if raw == nil {
		return nil, nil, nil
	}
	var blk Block
	if err := json.Unmarshal(raw, &blk); err != nil {
		return nil, nil, fmt.Errorf("decoding block: %w", err)
	}
	return &blk, raw, nil
}

// FetchBlockResults calls CometBFT /block_results at the given height.
// Returns (nil, nil) when --curl short-circuits.
func FetchBlockResults(ctx context.Context, rpc *client.RPCClient, height int64) (*BlockResults, error) {
	params := map[string]any{"height": strconv.FormatInt(height, 10)}
	raw, err := rpc.Call(ctx, "block_results", params)
	if err != nil {
		return nil, fmt.Errorf("fetching block_results: %w", err)
	}
	if raw == nil {
		return nil, nil
	}
	var out BlockResults
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("decoding block_results: %w", err)
	}
	return &out, nil
}

// ParseTimestamp accepts either a bare unix-second integer or an
// RFC3339 / RFC3339Nano string. Returns a time.Time in UTC.
func ParseTimestamp(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, fmt.Errorf("empty timestamp")
	}
	// All-digit input (optionally with leading minus): treat as unix seconds.
	if isAllDigits(s) {
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return time.Time{}, fmt.Errorf("parsing unix seconds %q: %w", s, err)
		}
		return time.Unix(n, 0).UTC(), nil
	}
	// RFC3339 / RFC3339Nano.
	if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
		return t.UTC(), nil
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t.UTC(), nil
	}
	return time.Time{}, fmt.Errorf("timestamp %q not recognised (want unix seconds or RFC3339)", s)
}

func isAllDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// FindBlockAt binary-searches the closed range [lo, hi] for the
// height whose block time is closest to target. Cancellable via ctx.
//
// Heimdall block times are monotonically non-decreasing so bsearch
// applies. The fetch-per-step cost is the dominant factor so we cap
// at log2(hi-lo+1) + 1 probes.
func FindBlockAt(ctx context.Context, rpc *client.RPCClient, lo, hi int64, target time.Time) (int64, error) {
	if lo == hi {
		return lo, nil
	}

	// Collect the bracketing heights so we can pick the closer one
	// once the search collapses.
	best := lo
	var bestDelta time.Duration = 1<<62 - 1

	consider := func(h int64, t time.Time) {
		delta := t.Sub(target)
		if delta < 0 {
			delta = -delta
		}
		if delta < bestDelta {
			bestDelta = delta
			best = h
		}
	}

	// Prime with the edges so the initial "best" is real.
	loTime, err := BlockTime(ctx, rpc, lo)
	if err != nil {
		return 0, err
	}
	consider(lo, loTime)
	if target.Before(loTime) {
		return lo, nil
	}

	hiTime, err := BlockTime(ctx, rpc, hi)
	if err != nil {
		return 0, err
	}
	consider(hi, hiTime)
	if target.After(hiTime) {
		return hi, nil
	}

	left, right := lo, hi
	for right-left > 1 {
		if cerr := ctx.Err(); cerr != nil {
			return 0, cerr
		}
		mid := left + (right-left)/2
		midTime, merr := BlockTime(ctx, rpc, mid)
		if merr != nil {
			return 0, merr
		}
		consider(mid, midTime)
		if midTime.Before(target) {
			left = mid
		} else {
			right = mid
		}
	}
	// Evaluate the final endpoints one more time in case they were
	// never explicitly considered.
	leftTime, err := BlockTime(ctx, rpc, left)
	if err != nil {
		return 0, err
	}
	consider(left, leftTime)
	rightTime, err := BlockTime(ctx, rpc, right)
	if err != nil {
		return 0, err
	}
	consider(right, rightTime)
	return best, nil
}

// BlockTime fetches /block at h and returns the header's time parsed
// as a Go time.Time. Separate from FetchBlock for tightness.
func BlockTime(ctx context.Context, rpc *client.RPCClient, h int64) (time.Time, error) {
	blk, raw, err := FetchBlock(ctx, rpc, strconv.FormatInt(h, 10))
	if err != nil {
		return time.Time{}, err
	}
	if raw == nil {
		return time.Time{}, fmt.Errorf("block time lookup does not support --curl")
	}
	t, err := time.Parse(time.RFC3339Nano, blk.Block.Header.Time)
	if err != nil {
		t, err = time.Parse(time.RFC3339, blk.Block.Header.Time)
		if err != nil {
			return time.Time{}, fmt.Errorf("parsing block %d time %q: %w", h, blk.Block.Header.Time, err)
		}
	}
	return t.UTC(), nil
}
