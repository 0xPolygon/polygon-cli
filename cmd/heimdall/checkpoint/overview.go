package checkpoint

import (
	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/cmdutil"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newOverviewCmd builds `checkpoint overview` → GET /checkpoints/overview.
// The response is a dashboard bundle (ack count, buffer, validator set).
// We emit JSON by default because the shape is too nested for the KV
// renderer to be useful.
func newOverviewCmd() *cobra.Command {
	return pkg.NewGetCmd(cmdutil.Get{
		Use:   "overview",
		Short: "Checkpoint module dashboard bundle.",
		Path:  "/checkpoints/overview",
		Label: "overview",
		Render: func(cmd *cobra.Command, m map[string]any, opts render.Options) error {
			return render.RenderJSON(cmd.OutOrStdout(), m, opts)
		},
	})
}
