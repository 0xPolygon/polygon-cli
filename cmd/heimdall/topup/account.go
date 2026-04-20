package topup

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newAccountCmd builds `topup account <ADDR>` → GET
// /topup/dividend-account/{address}. Prints the `user` and
// `fee_amount` fields of the dividend account.
func newAccountCmd() *cobra.Command {
	var fields []string
	cmd := &cobra.Command{
		Use:   "account <ADDR>",
		Short: "Fetch the dividend account for an address.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			addr, err := normalizeAddress(args[0])
			if err != nil {
				return err
			}
			rest, cfg, err := newRESTClient(cmd)
			if err != nil {
				return err
			}
			body, status, err := rest.Get(cmd.Context(), fmt.Sprintf("/topup/dividend-account/%s", addr), nil)
			if err != nil {
				return err
			}
			if status == 0 && body == nil {
				return nil
			}
			opts := renderOpts(cmd, cfg, fields)
			m, err := decodeJSONMap(body, "topup dividend-account")
			if err != nil {
				return err
			}
			if opts.JSON {
				return render.RenderJSON(cmd.OutOrStdout(), m, opts)
			}
			// Unwrap the { "dividend_account": {...} } envelope for KV.
			if inner, ok := m["dividend_account"].(map[string]any); ok {
				return render.RenderKV(cmd.OutOrStdout(), inner, opts)
			}
			return render.RenderKV(cmd.OutOrStdout(), m, opts)
		},
	}
	cmd.Flags().StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable)")
	return cmd
}
