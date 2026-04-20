package span

import (
	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newParamsCmd builds `span params` → GET /bor/params. Prints sprint
// duration, span duration, and producer count.
func newParamsCmd() *cobra.Command {
	var fields []string
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Show bor module parameters.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			rest, cfg, err := newRESTClient(cmd)
			if err != nil {
				return err
			}
			body, status, err := rest.Get(cmd.Context(), "/bor/params", nil)
			if err != nil {
				return err
			}
			if status == 0 && body == nil {
				return nil
			}
			opts := renderOpts(cmd, cfg, fields)
			m, err := decodeJSONMap(body, "bor params")
			if err != nil {
				return err
			}
			if opts.JSON {
				return render.RenderJSON(cmd.OutOrStdout(), m, opts)
			}
			// Unwrap the { "params": { ... } } envelope for KV output.
			if inner, ok := m["params"].(map[string]any); ok {
				return render.RenderKV(cmd.OutOrStdout(), inner, opts)
			}
			return render.RenderKV(cmd.OutOrStdout(), m, opts)
		},
	}
	cmd.Flags().StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable)")
	return cmd
}
