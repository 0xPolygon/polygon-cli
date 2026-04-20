package validator

import (
	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newProposerCmd builds `validator proposer` → GET
// /stake/proposers/current. Single validator envelope, same shape as
// `/stake/validator/{id}`.
func newProposerCmd() *cobra.Command {
	var fields []string
	cmd := &cobra.Command{
		Use:   "proposer",
		Short: "Show the current proposer.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			rest, cfg, err := newRESTClient(cmd)
			if err != nil {
				return err
			}
			body, status, err := rest.Get(cmd.Context(), "/stake/proposers/current", nil)
			if err != nil {
				return err
			}
			if status == 0 && body == nil {
				return nil
			}
			opts := renderOpts(cmd, cfg, fields)
			m, err := decodeJSONMap(body, "proposer")
			if err != nil {
				return err
			}
			if opts.JSON {
				return render.RenderJSON(cmd.OutOrStdout(), m, opts)
			}
			return renderValidatorKV(cmd, m, opts)
		},
	}
	cmd.Flags().StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable)")
	return cmd
}
