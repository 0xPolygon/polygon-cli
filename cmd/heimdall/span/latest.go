package span

import (
	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newLatestCmd builds `span latest` → GET /bor/spans/latest. The
// single-span envelope is unwrapped for KV output; the deeply-nested
// validator_set is emitted as a JSON blob on its own line.
func newLatestCmd() *cobra.Command {
	var fields []string
	cmd := &cobra.Command{
		Use:   "latest",
		Short: "Show the current (latest) span.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			rest, cfg, err := newRESTClient(cmd)
			if err != nil {
				return err
			}
			body, status, err := rest.Get(cmd.Context(), "/bor/spans/latest", nil)
			if err != nil {
				return err
			}
			if status == 0 && body == nil {
				return nil
			}
			opts := renderOpts(cmd, cfg, fields)
			m, err := decodeJSONMap(body, "span latest")
			if err != nil {
				return err
			}
			if opts.JSON {
				return render.RenderJSON(cmd.OutOrStdout(), m, opts)
			}
			return renderSpanKV(cmd, m, opts)
		},
	}
	cmd.Flags().StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable)")
	return cmd
}

// renderSpanKV unwraps the { "span": {...} } envelope and renders with
// the shared KV formatter. Nested objects (validator_set,
// selected_producers) are emitted inline as JSON by the KV renderer.
func renderSpanKV(cmd *cobra.Command, m map[string]any, opts render.Options) error {
	inner, ok := m["span"].(map[string]any)
	if !ok {
		return render.RenderKV(cmd.OutOrStdout(), m, opts)
	}
	return render.RenderKV(cmd.OutOrStdout(), inner, opts)
}
