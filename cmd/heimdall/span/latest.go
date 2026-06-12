package span

import (
	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/cmdutil"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newLatestCmd builds `span latest` → GET /bor/spans/latest. The
// single-span envelope is unwrapped for KV output; the deeply-nested
// validator_set is emitted as a JSON blob on its own line.
func newLatestCmd() *cobra.Command {
	return pkg.NewGetCmd(cmdutil.Get{
		Use:    "latest",
		Short:  "Show the current (latest) span.",
		Path:   "/bor/spans/latest",
		Label:  "span latest",
		Render: renderSpanKV,
	})
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
