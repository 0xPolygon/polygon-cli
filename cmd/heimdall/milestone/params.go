package milestone

import (
	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/cmdutil"
)

// newParamsCmd builds `milestone params` → GET /milestones/params.
// Prints thresholds + interval. The { "params": { ... } } envelope is
// unwrapped for KV output.
func newParamsCmd() *cobra.Command {
	return pkg.NewGetCmd(cmdutil.Get{
		Use:       "params",
		Short:     "Show milestone module parameters.",
		Path:      "/milestones/params",
		Label:     "milestone params",
		UnwrapKey: "params",
	})
}
