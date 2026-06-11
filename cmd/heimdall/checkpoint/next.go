package checkpoint

import (
	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/cmdutil"
)

// newNextCmd builds `checkpoint next` → GET /checkpoints/prepare-next.
// Requires the node to have L1 RPC configured; the special gRPC-code
// 13 case is surfaced (via L1Hint) with a hint about missing
// `eth_rpc_url`, both on HTTP errors carrying the gRPC envelope and on
// 2xx responses that still smuggle the envelope on older Heimdalls.
func newNextCmd() *cobra.Command {
	return pkg.NewGetCmd(cmdutil.Get{
		Use:    "next",
		Short:  "Compute the next checkpoint to propose.",
		Path:   "/checkpoints/prepare-next",
		Label:  "prepare-next",
		L1Hint: true,
		Render: renderCheckpointKV,
	})
}
