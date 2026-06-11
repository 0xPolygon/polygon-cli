package chainparams

import (
	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/cmdutil"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newParamsCmd builds `chainmanager params` → GET /chainmanager/params.
// The response shape is:
//
//	{
//	  "params": {
//	    "chain_params": { ... addresses ... },
//	    "main_chain_tx_confirmations": "64",
//	    "bor_chain_tx_confirmations": "512"
//	  }
//	}
//
// Default human output unwraps the `params` envelope for KV rendering.
// --json preserves the raw server shape.
func newParamsCmd() *cobra.Command {
	return pkg.NewGetCmd(cmdutil.Get{
		Use:   "params",
		Short: "Fetch the chainmanager module parameters.",
		Path:  "/chainmanager/params",
		Label: "chainmanager params",
		// --field addresses the raw server shape so the envelope
		// stays visible in the path. Only unwrap the `params`
		// envelope for the default (no --field) KV render — which is
		// why this is a Render hook rather than UnwrapKey.
		Render: func(cmd *cobra.Command, m map[string]any, opts render.Options) error {
			if len(opts.Fields) > 0 {
				return render.RenderKV(cmd.OutOrStdout(), m, opts)
			}
			if inner, ok := m["params"].(map[string]any); ok {
				return render.RenderKV(cmd.OutOrStdout(), inner, opts)
			}
			return render.RenderKV(cmd.OutOrStdout(), m, opts)
		},
	})
}
