package checkpoint

import (
	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newSignaturesCmd builds `checkpoint signatures <TX_HASH>` → GET
// /checkpoints/signatures/{hash}. Tolerates the `0x` prefix.
func newSignaturesCmd() *cobra.Command {
	var fields []string
	cmd := &cobra.Command{
		Use:   "signatures <TX_HASH>",
		Short: "Aggregated validator signatures for a checkpoint tx.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			hash, err := normalizeCheckpointHash(args[0])
			if err != nil {
				return err
			}
			rest, cfg, err := newRESTClient(cmd)
			if err != nil {
				return err
			}
			body, status, err := rest.Get(cmd.Context(), "/checkpoints/signatures/"+hash, nil)
			if err != nil {
				return err
			}
			if status == 0 && body == nil {
				return nil
			}
			opts := renderOpts(cmd, cfg, fields)
			m, err := decodeJSONMap(body, "signatures")
			if err != nil {
				return err
			}
			return render.RenderJSON(cmd.OutOrStdout(), m, opts)
		},
	}
	cmd.Flags().StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable)")
	return cmd
}
