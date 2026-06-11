package validator

import (
	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/cmdutil"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newSignerCmd builds `validator signer <ADDR>` → GET
// /stake/signer/{addr}. Tolerates a missing `0x` prefix and lower/upper
// case input; the REST endpoint accepts the prefixed form uniformly.
//
// Kept as a custom RunE rather than a cmdutil.Get spec: this command
// historically exposes no --field flag, which the builder would add.
func newSignerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "signer <ADDR>",
		Short: "Fetch a validator by hex signer address.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			addr, err := cmdutil.NormalizeHex(args[0], 20, "signer")
			if err != nil {
				return err
			}
			rest, cfg, err := pkg.RESTClient(cmd)
			if err != nil {
				return err
			}
			body, status, err := rest.Get(cmd.Context(), "/stake/signer/"+addr, nil)
			if err != nil {
				return err
			}
			if status == 0 && body == nil {
				return nil
			}
			opts := cmdutil.RenderOpts(cmd, cfg, nil)
			m, err := cmdutil.DecodeJSONMap(body, "signer")
			if err != nil {
				return err
			}
			if opts.JSON {
				return render.RenderJSON(cmd.OutOrStdout(), m, opts)
			}
			return renderValidatorKV(cmd, m, opts)
		},
	}
	return cmd
}
