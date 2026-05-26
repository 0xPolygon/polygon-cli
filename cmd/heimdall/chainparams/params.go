package chainparams

import (
	"github.com/spf13/cobra"

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
	var fields []string
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Fetch the chainmanager module parameters.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			rest, cfg, err := newRESTClient(cmd)
			if err != nil {
				return err
			}
			body, status, err := rest.Get(cmd.Context(), "/chainmanager/params", nil)
			if err != nil {
				return err
			}
			if status == 0 && body == nil {
				// --curl short-circuit.
				return nil
			}
			opts := renderOpts(cmd, cfg, fields)
			m, err := decodeJSONMap(body, "chainmanager params")
			if err != nil {
				return err
			}
			if opts.JSON {
				return render.RenderJSON(cmd.OutOrStdout(), m, opts)
			}
			// --field addresses the raw server shape so the envelope
			// stays visible in the path. Only unwrap the `params`
			// envelope for the default (no --field) KV render.
			if len(opts.Fields) > 0 {
				return render.RenderKV(cmd.OutOrStdout(), m, opts)
			}
			if inner, ok := m["params"].(map[string]any); ok {
				return render.RenderKV(cmd.OutOrStdout(), inner, opts)
			}
			return render.RenderKV(cmd.OutOrStdout(), m, opts)
		},
	}
	cmd.Flags().StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable)")
	return cmd
}
