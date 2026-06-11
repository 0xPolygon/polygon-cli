package clerk

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/cmdutil"
)

// buildTxEventQuery validates the <TX_HASH> <LOG_INDEX> argument pair
// shared by `sequence` and `is-old` and turns it into the
// tx_hash/log_index query both endpoints expect.
func buildTxEventQuery(args []string) (url.Values, error) {
	hash, err := cmdutil.NormalizeTxHash(args[0])
	if err != nil {
		return nil, err
	}
	logIndex, err := strconv.ParseUint(args[1], 10, 64)
	if err != nil {
		return nil, &client.UsageError{Msg: fmt.Sprintf("log_index must be a non-negative integer, got %q", args[1])}
	}
	q := url.Values{}
	q.Set("tx_hash", hash)
	q.Set("log_index", strconv.FormatUint(logIndex, 10))
	return q, nil
}

// newSequenceCmd builds `state-sync sequence <TX_HASH> <LOG_INDEX>` →
// GET /clerk/sequence?tx_hash=…&log_index=…. Like `is-old`, this
// endpoint fans out to L1 on the server side; a node without
// `eth_rpc_url` will return gRPC code 13, which we surface as an
// L1-not-configured hint.
func newSequenceCmd() *cobra.Command {
	return pkg.NewGetCmd(cmdutil.Get{
		Use:    "sequence <TX_HASH> <LOG_INDEX>",
		Short:  "Dedup sequence key for an L1 state-sync event.",
		Args:   cobra.ExactArgs(2),
		Label:  "clerk sequence",
		L1Hint: true,
		Build: func(cmd *cobra.Command, args []string) (string, url.Values, error) {
			q, err := buildTxEventQuery(args)
			if err != nil {
				return "", nil, err
			}
			return "/clerk/sequence", q, nil
		},
	})
}
