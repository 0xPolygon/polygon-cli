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

// newIsOldCmd builds `topup is-old <TX_HASH> <LOG_INDEX>` → GET
// /topup/is-old-tx?tx_hash=…&log_index=…. Requires L1 RPC on the
// Heimdall node; on gRPC code 13 the command surfaces an
// L1-not-configured hint on stderr before propagating the error.
//
// Default text output is a bare `true`/`false` so shell scripts can
// pipe it without parsing JSON.
func newIsOldCmd() *cobra.Command {
	return pkg.NewGetCmd(cmdutil.Get{
		Use:    "is-old <TX_HASH> <LOG_INDEX>",
		Short:  "Check whether an L1 topup tx was already processed.",
		Args:   cobra.ExactArgs(2),
		Label:  "topup is-old-tx",
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
			return "/topup/is-old-tx", q, nil
		},
		Render: func(cmd *cobra.Command, m map[string]any, opts render.Options) error {
			// Default text output: bare bool when no --field filter.
			if v, ok := m["is_old"].(bool); ok && len(opts.Fields) == 0 {
				var err error
				if v {
					_, err = fmt.Fprintln(cmd.OutOrStdout(), "true")
				} else {
					_, err = fmt.Fprintln(cmd.OutOrStdout(), "false")
				}
				return err
			}
			return render.RenderKV(cmd.OutOrStdout(), m, opts)
		},
	})
}
