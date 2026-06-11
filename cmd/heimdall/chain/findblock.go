package chain

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/comet"
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

			rpc, _, err := pkg.RPCClient(cmd)
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

// parseTimestamp delegates to comet.ParseTimestamp.
func parseTimestamp(s string) (time.Time, error) {
	return comet.ParseTimestamp(s)
}

// findBlockAt delegates to comet.FindBlockAt.
func findBlockAt(ctx context.Context, rpc *client.RPCClient, lo, hi int64, target time.Time) (int64, error) {
	return comet.FindBlockAt(ctx, rpc, lo, hi, target)
}
