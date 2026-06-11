package clerk

import (
	"net/url"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/cmdutil"
)

// newIsOldCmd builds `state-sync is-old <TX_HASH> <LOG_INDEX>` → GET
// /clerk/is-old-tx?tx_hash=…&log_index=…. On gRPC code 13 (or a
// transport-level `connection refused` / `dial tcp`) the command
// surfaces an L1-not-configured hint before propagating the error,
// matching the shape of `validator is-old-stake-tx`.
func newIsOldCmd() *cobra.Command {
	return pkg.NewGetCmd(cmdutil.Get{
		Use:    "is-old <TX_HASH> <LOG_INDEX>",
		Short:  "Check whether an L1 state-sync event was already replayed.",
		Args:   cobra.ExactArgs(2),
		Label:  "clerk is-old-tx",
		L1Hint: true,
		Build: func(cmd *cobra.Command, args []string) (string, url.Values, error) {
			q, err := buildTxEventQuery(args)
			if err != nil {
				return "", nil, err
			}
			return "/clerk/is-old-tx", q, nil
		},
	})
}
