package topup

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/cmdutil"
)

// newAccountCmd builds `topup account <ADDR>` → GET
// /topup/dividend-account/{address}. Prints the `user` and
// `fee_amount` fields of the dividend account. The
// { "dividend_account": {...} } envelope is unwrapped for KV output.
func newAccountCmd() *cobra.Command {
	return pkg.NewGetCmd(cmdutil.Get{
		Use:   "account <ADDR>",
		Short: "Fetch the dividend account for an address.",
		Args:  cobra.ExactArgs(1),
		Label: "topup dividend-account",
		Build: func(_ *cobra.Command, args []string) (string, url.Values, error) {
			addr, err := cmdutil.NormalizeAddress(args[0])
			if err != nil {
				return "", nil, err
			}
			return fmt.Sprintf("/topup/dividend-account/%s", addr), nil, nil
		},
		UnwrapKey: "dividend_account",
	})
}
