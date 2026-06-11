package topup

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/cmdutil"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newSequenceCmd builds `topup sequence <TX_HASH> <LOG_INDEX>` → GET
// /topup/sequence?tx_hash=…&log_index=…. Requires L1 RPC on the
// Heimdall node; on gRPC code 13 the command surfaces an
// L1-not-configured hint on stderr before propagating the error.
func newSequenceCmd() *cobra.Command {
	return pkg.NewGetCmd(cmdutil.Get{
		Use:    "sequence <TX_HASH> <LOG_INDEX>",
		Short:  "Dedup sequence key for an L1 topup tx.",
		Args:   cobra.ExactArgs(2),
		Label:  "topup sequence",
		L1Hint: true,
		Build: func(_ *cobra.Command, args []string) (string, url.Values, error) {
			hash, err := cmdutil.NormalizeTxHash(args[0])
			if err != nil {
				return "", nil, err
			}
			logIndex, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return "", nil, &client.UsageError{Msg: fmt.Sprintf("log_index must be a non-negative integer, got %q", args[1])}
			}
			q := url.Values{}
			q.Set("tx_hash", hash)
			q.Set("log_index", strconv.FormatUint(logIndex, 10))
			return "/topup/sequence", q, nil
		},
		Render: func(cmd *cobra.Command, m map[string]any, opts render.Options) error {
			// Default text output: print the bare sequence if present.
			if seq, ok := m["sequence"].(string); ok && len(opts.Fields) == 0 {
				_, err := fmt.Fprintln(cmd.OutOrStdout(), seq)
				return err
			}
			return render.RenderKV(cmd.OutOrStdout(), m, opts)
		},
	})
}
