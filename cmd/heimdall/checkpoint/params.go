package checkpoint

import (
	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/cmdutil"
)

// newParamsCmd builds `checkpoint params` → GET /checkpoints/params.
// Prints interval, buffer time, max/avg length, chain interval. The
// { "params": { ... } } envelope is unwrapped for KV output.
func newParamsCmd() *cobra.Command {
	return pkg.NewGetCmd(cmdutil.Get{
		Use:       "params",
		Short:     "Show checkpoint module parameters.",
		Path:      "/checkpoints/params",
		Label:     "checkpoint params",
		UnwrapKey: "params",
	})
}
