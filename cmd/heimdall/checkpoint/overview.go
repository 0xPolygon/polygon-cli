package checkpoint

import (
	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newOverviewCmd builds `checkpoint overview` → GET /checkpoints/overview.
// The response is a dashboard bundle (ack count, buffer, validator set).
// We emit JSON by default because the shape is too nested for the KV
// renderer to be useful.
func newOverviewCmd() *cobra.Command {
	var fields []string
	cmd := &cobra.Command{
		Use:   "overview",
		Short: "Checkpoint module dashboard bundle.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			rest, cfg, err := newRESTClient(cmd)
			if err != nil {
				return err
			}
			body, status, err := rest.Get(cmd.Context(), "/checkpoints/overview", nil)
			if err != nil {
				return err
			}
			if status == 0 && body == nil {
				return nil
			}
			opts := renderOpts(cmd, cfg, fields)
			m, err := decodeJSONMap(body, "overview")
			if err != nil {
				return err
			}
			return render.RenderJSON(cmd.OutOrStdout(), m, opts)
		},
	}
	cmd.Flags().StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable)")
	return cmd
}
