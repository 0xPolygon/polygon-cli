package checkpoint

import (
	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/cmdutil"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newLatestCmd builds `checkpoint latest` → GET /checkpoints/latest.
// The single-checkpoint envelope is unwrapped for KV output; timestamp
// (if present) is annotated with the human-readable age.
func newLatestCmd() *cobra.Command {
	return pkg.NewGetCmd(cmdutil.Get{
		Use:    "latest",
		Short:  "Show the latest acked checkpoint.",
		Path:   "/checkpoints/latest",
		Label:  "checkpoint latest",
		Render: renderCheckpointKV,
	})
}

// renderCheckpointKV unwraps the { "checkpoint": {...} } envelope,
// annotates the timestamp with human-readable age, and renders with
// the shared KV formatter.
func renderCheckpointKV(cmd *cobra.Command, m map[string]any, opts render.Options) error {
	inner, ok := m["checkpoint"].(map[string]any)
	if !ok {
		return render.RenderKV(cmd.OutOrStdout(), m, opts)
	}
	if ts, ok := inner["timestamp"].(string); ok && ts != "" {
		inner["timestamp"] = render.AnnotateUnixSeconds(ts)
	}
	return render.RenderKV(cmd.OutOrStdout(), inner, opts)
}
