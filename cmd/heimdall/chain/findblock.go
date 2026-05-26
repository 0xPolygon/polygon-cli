package chain

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
)

// newFindBlockCmd builds `find-block <TIMESTAMP>`. Binary-searches
// CometBFT /block to find the height whose block time is closest to
// TIMESTAMP. Accepts either unix seconds or an RFC3339 string.
func newFindBlockCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "find-block <TIMESTAMP>",
		Short: "Find the block height closest to a timestamp.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target, err := parseTimestamp(args[0])
			if err != nil {
				return &client.UsageError{Msg: err.Error()}
			}

			rpc, _, err := newRPCClient(cmd)
			if err != nil {
				return err
			}
			ctx := cmd.Context()

			st, err := fetchStatus(ctx, rpc)
			if err != nil {
				return err
			}
			if st == nil {
				return nil // --curl
			}

			lo, err := strconv.ParseInt(st.SyncInfo.EarliestBlockHeight, 10, 64)
			if err != nil {
				return fmt.Errorf("parsing earliest height %q: %w", st.SyncInfo.EarliestBlockHeight, err)
			}
			hi, err := strconv.ParseInt(st.SyncInfo.LatestBlockHeight, 10, 64)
			if err != nil {
				return fmt.Errorf("parsing latest height %q: %w", st.SyncInfo.LatestBlockHeight, err)
			}
			if lo > hi {
				return fmt.Errorf("inconsistent sync info: earliest %d > latest %d", lo, hi)
			}

			h, err := findBlockAt(ctx, rpc, lo, hi, target)
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), h)
			return err
		},
	}
}

// parseTimestamp accepts either a bare unix-second integer or an
// RFC3339 / RFC3339Nano string. Returns a time.Time in UTC.
func parseTimestamp(s string) (time.Time, error) {
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

// findBlockAt binary-searches the closed range [lo, hi] for the
// height whose block time is closest to target. Cancellable via ctx.
//
// Heimdall block times are monotonically non-decreasing so bsearch
// applies. The fetch-per-step cost is the dominant factor so we cap
// at log2(hi-lo+1) + 1 probes.
func findBlockAt(ctx context.Context, rpc *client.RPCClient, lo, hi int64, target time.Time) (int64, error) {
	if lo == hi {
		return lo, nil
	}

	// Collect the bracketing heights so we can pick the closer one
	// once the search collapses.
	var best int64 = lo
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
	loTime, err := blockTime(ctx, rpc, lo)
	if err != nil {
		return 0, err
	}
	consider(lo, loTime)
	if target.Before(loTime) {
		return lo, nil
	}

	hiTime, err := blockTime(ctx, rpc, hi)
	if err != nil {
		return 0, err
	}
	consider(hi, hiTime)
	if target.After(hiTime) {
		return hi, nil
	}

	left, right := lo, hi
	for right-left > 1 {
		if err := ctx.Err(); err != nil {
			return 0, err
		}
		mid := left + (right-left)/2
		midTime, err := blockTime(ctx, rpc, mid)
		if err != nil {
			return 0, err
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
	leftTime, err := blockTime(ctx, rpc, left)
	if err != nil {
		return 0, err
	}
	consider(left, leftTime)
	rightTime, err := blockTime(ctx, rpc, right)
	if err != nil {
		return 0, err
	}
	consider(right, rightTime)
	return best, nil
}

// blockTime fetches /block at h and returns the header's time parsed
// as a Go time.Time. Separate from fetchBlock for tightness.
func blockTime(ctx context.Context, rpc *client.RPCClient, h int64) (time.Time, error) {
	blk, raw, err := fetchBlock(ctx, rpc, strconv.FormatInt(h, 10))
	if err != nil {
		return time.Time{}, err
	}
	if raw == nil {
		return time.Time{}, fmt.Errorf("find-block does not support --curl")
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
