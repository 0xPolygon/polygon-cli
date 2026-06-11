package checkpoint

import (
	"net/url"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/cmdutil"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newSignaturesCmd builds `checkpoint signatures <TX_HASH>` → GET
// /checkpoints/signatures/{hash}. Tolerates the `0x` prefix. Output is
// always JSON; the shape is too nested for the KV renderer.
func newSignaturesCmd() *cobra.Command {
	return pkg.NewGetCmd(cmdutil.Get{
		Use:   "signatures <TX_HASH>",
		Short: "Aggregated validator signatures for a checkpoint tx.",
		Args:  cobra.ExactArgs(1),
		Label: "signatures",
		Build: func(cmd *cobra.Command, args []string) (string, url.Values, error) {
			hash, err := normalizeCheckpointHash(args[0])
			if err != nil {
				return "", nil, err
			}
			return "/checkpoints/signatures/" + hash, nil, nil
		},
		Render: func(cmd *cobra.Command, m map[string]any, opts render.Options) error {
			return render.RenderJSON(cmd.OutOrStdout(), m, opts)
		},
	})
}
