package checkpoint

import (
	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newLatestCmd builds `checkpoint latest` → GET /checkpoints/latest.
// The single-checkpoint envelope is unwrapped for KV output; timestamp
// (if present) is annotated with the human-readable age.
func newLatestCmd() *cobra.Command {
	var fields []string
	cmd := &cobra.Command{
		Use:   "latest",
		Short: "Show the latest acked checkpoint.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			rest, cfg, err := newRESTClient(cmd)
			if err != nil {
				return err
			}
			body, status, err := rest.Get(cmd.Context(), "/checkpoints/latest", nil)
			if err != nil {
				return err
			}
			if status == 0 && body == nil {
				return nil
			}
			opts := renderOpts(cmd, cfg, fields)
			m, err := decodeJSONMap(body, "checkpoint latest")
			if err != nil {
				return err
			}
			if opts.JSON {
				return render.RenderJSON(cmd.OutOrStdout(), m, opts)
			}
			return renderCheckpointKV(cmd, m, opts)
		},
	}
	cmd.Flags().StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable)")
	return cmd
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
