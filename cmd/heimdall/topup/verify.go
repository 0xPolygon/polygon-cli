package topup

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newVerifyCmd builds `topup verify <ADDR> <PROOF>` → GET
// /topup/account-proof/{address}/verify?proof=…. Confirmed from
// heimdall-v2 query.proto: the upstream uses GET (not POST) and the
// proof travels as a query parameter.
func newVerifyCmd() *cobra.Command {
	var fields []string
	cmd := &cobra.Command{
		Use:   "verify <ADDR> <PROOF>",
		Short: "Verify a submitted Merkle proof for a dividend account.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			addr, err := normalizeAddress(args[0])
			if err != nil {
				return err
			}
			proof, err := normalizeHexBytes(args[1], "proof")
			if err != nil {
				return err
			}
			rest, cfg, err := newRESTClient(cmd)
			if err != nil {
				return err
			}
			q := url.Values{}
			q.Set("proof", proof)
			body, status, err := rest.Get(cmd.Context(), fmt.Sprintf("/topup/account-proof/%s/verify", addr), q)
			if err != nil {
				return err
			}
			if status == 0 && body == nil {
				return nil
			}
			opts := renderOpts(cmd, cfg, fields)
			m, err := decodeJSONMap(body, "topup verify")
			if err != nil {
				return err
			}
			if opts.JSON {
				return render.RenderJSON(cmd.OutOrStdout(), m, opts)
			}
			// Default human text: print the bool directly.
			if v, ok := m["is_verified"].(bool); ok {
				if v {
					_, err = fmt.Fprintln(cmd.OutOrStdout(), "true")
				} else {
					_, err = fmt.Fprintln(cmd.OutOrStdout(), "false")
				}
				return err
			}
			return render.RenderKV(cmd.OutOrStdout(), m, opts)
		},
	}
	cmd.Flags().StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable)")
	return cmd
}
