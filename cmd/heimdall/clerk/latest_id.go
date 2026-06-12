package clerk

import (
	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/cmdutil"
)

// newLatestIDCmd builds `state-sync latest-id` → GET
// /clerk/event-records/latest-id. Requires L1 RPC on the node; on gRPC
// code 13 (or a transport-level `connection refused`) the command
// prints a human-friendly hint pointing at missing `eth_rpc_url`
// configuration before propagating the error. Some Heimdalls surface
// gRPC code 13 on 2xx with the envelope body; L1Hint checks that too.
func newLatestIDCmd() *cobra.Command {
	return pkg.NewGetCmd(cmdutil.Get{
		Use:    "latest-id",
		Short:  "Latest L1 state-sync counter (needs eth_rpc_url).",
		Path:   "/clerk/event-records/latest-id",
		Label:  "clerk latest-id",
		L1Hint: true,
	})
}
