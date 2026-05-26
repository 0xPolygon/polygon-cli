package ops

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// cometStatus is the subset of /status we surface as the summary.
type cometStatus struct {
	NodeInfo struct {
		ID      string `json:"id"`
		Network string `json:"network"`
		Moniker string `json:"moniker"`
		Version string `json:"version"`
	} `json:"node_info"`
	SyncInfo struct {
		LatestBlockHeight string `json:"latest_block_height"`
		LatestBlockTime   string `json:"latest_block_time"`
		CatchingUp        bool   `json:"catching_up"`
	} `json:"sync_info"`
	ValidatorInfo struct {
		Address     string `json:"address"`
		VotingPower string `json:"voting_power"`
	} `json:"validator_info"`
}

// newStatusCmd builds `ops status`. Default output is a KV summary;
// --json passes the full /status result through.
func newStatusCmd() *cobra.Command {
	var fields []string
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show CometBFT /status: height, sync, moniker, own validator.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			rpc, cfg, err := newRPCClient(cmd)
			if err != nil {
				return err
			}
			raw, err := callEmpty(cmd.Context(), rpc, "status")
			if err != nil {
				return err
			}
			if raw == nil {
				return nil // --curl
			}
			opts := renderOpts(cmd, cfg, fields)
			if opts.JSON {
				generic, derr := decodeGeneric(raw)
				if derr != nil {
					return derr
				}
				return render.RenderJSON(cmd.OutOrStdout(), generic, opts)
			}
			var st cometStatus
			if err := json.Unmarshal(raw, &st); err != nil {
				return fmt.Errorf("decoding status: %w", err)
			}
			out := map[string]any{
				"node_id":             st.NodeInfo.ID,
				"moniker":             st.NodeInfo.Moniker,
				"network":             st.NodeInfo.Network,
				"cometbft_version":    st.NodeInfo.Version,
				"latest_block_height": st.SyncInfo.LatestBlockHeight,
				"latest_block_time":   st.SyncInfo.LatestBlockTime,
				"catching_up":         st.SyncInfo.CatchingUp,
				"validator_address":   "0x" + st.ValidatorInfo.Address,
				"voting_power":        st.ValidatorInfo.VotingPower,
			}
			return render.RenderKV(cmd.OutOrStdout(), out, opts)
		},
	}
	f := cmd.Flags()
	f.StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable)")
	return cmd
}
