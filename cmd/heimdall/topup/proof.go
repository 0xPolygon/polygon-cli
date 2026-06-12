package topup

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/cmdutil"
)

// newProofCmd builds `topup proof <ADDR>` → GET
// /topup/account-proof/{address}. Requires L1 RPC on the Heimdall
// node; on gRPC code 13 the command surfaces an L1-not-configured
// hint on stderr before propagating the error. The { "proof": {...} }
// envelope is unwrapped for KV output.
func newProofCmd() *cobra.Command {
	return pkg.NewGetCmd(cmdutil.Get{
		Use:   "proof <ADDR>",
		Short: "Fetch the Merkle proof for a dividend account.",
		Args:  cobra.ExactArgs(1),
		Label: "topup account-proof",
		Build: func(_ *cobra.Command, args []string) (string, url.Values, error) {
			addr, err := cmdutil.NormalizeAddress(args[0])
			if err != nil {
				return "", nil, err
			}
			return fmt.Sprintf("/topup/account-proof/%s", addr), nil, nil
		},
		L1Hint:    true,
		UnwrapKey: "proof",
	})
}
