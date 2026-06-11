package span

import (
	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/cmdutil"
)

// newParamsCmd builds `span params` → GET /bor/params. Prints sprint
// duration, span duration, and producer count. The { "params": { ... } }
// envelope is unwrapped for KV output.
func newParamsCmd() *cobra.Command {
	return pkg.NewGetCmd(cmdutil.Get{
		Use:       "params",
		Short:     "Show bor module parameters.",
		Path:      "/bor/params",
		Label:     "bor params",
		UnwrapKey: "params",
	})
}
