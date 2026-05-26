package validator

import (
	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newSignerCmd builds `validator signer <ADDR>` → GET
// /stake/signer/{addr}. Tolerates a missing `0x` prefix and lower/upper
// case input; the REST endpoint accepts the prefixed form uniformly.
func newSignerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "signer <ADDR>",
		Short: "Fetch a validator by hex signer address.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			addr, err := normalizeSignerAddress(args[0])
			if err != nil {
				return err
			}
			rest, cfg, err := newRESTClient(cmd)
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
			opts := renderOpts(cmd, cfg, nil)
			m, err := decodeJSONMap(body, "signer")
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
