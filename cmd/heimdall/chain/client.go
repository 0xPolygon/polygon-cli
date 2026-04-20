package chain

import (
	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newClientCmd builds `client`. Surfaces the Heimdall app version
// (from /abci_info.response.version + response.data) and the
// CometBFT binary version (from /status.node_info.version).
func newClientCmd() *cobra.Command {
	var fields []string
	cmd := &cobra.Command{
		Use:   "client",
		Short: "Show Heimdall app + CometBFT versions.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			rpc, cfg, err := newRPCClient(cmd)
			if err != nil {
				return err
			}
			ctx := cmd.Context()
			abci, err := fetchABCIInfo(ctx, rpc)
			if err != nil {
				return err
			}
			if abci == nil {
				return nil
			}
			st, err := fetchStatus(ctx, rpc)
			if err != nil {
				return err
			}
			if st == nil {
				return nil
			}

			opts := renderOpts(cmd, cfg, fields)
			out := map[string]any{
				"heimdall_app":     abci.Response.Data,
				"heimdall_version": abci.Response.Version,
				"cometbft_version": st.NodeInfo.Version,
				"moniker":          st.NodeInfo.Moniker,
				"network":          st.NodeInfo.Network,
			}
			if opts.JSON {
				return render.RenderJSON(cmd.OutOrStdout(), out, opts)
			}
			return render.RenderKV(cmd.OutOrStdout(), out, opts)
		},
	}
	cmd.Flags().StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable)")
	return cmd
}
