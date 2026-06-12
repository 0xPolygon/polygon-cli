package validator

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/cmdutil"
)

// newIsOldStakeTxCmd builds `validator is-old-stake-tx <TX_HASH>
// <LOG_INDEX>` → GET /stake/is-old-tx?tx_hash=…&log_index=…. On gRPC
// code 13 (or a transport-level `connection refused` / `dial tcp`)
// the command prints a human-friendly hint pointing at missing
// `eth_rpc_url` configuration before propagating the error (L1Hint).
func newIsOldStakeTxCmd() *cobra.Command {
	return pkg.NewGetCmd(cmdutil.Get{
		Use:    "is-old-stake-tx <TX_HASH> <LOG_INDEX>",
		Short:  "Check whether an L1 stake event was already replayed.",
		Args:   cobra.ExactArgs(2),
		Label:  "is-old-tx",
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
			return "/stake/is-old-tx", q, nil
		},
	})
}
