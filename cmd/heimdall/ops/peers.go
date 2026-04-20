package ops

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// cometNetInfo is the subset of /net_info used for the default table.
type cometNetInfo struct {
	NPeers string `json:"n_peers"`
	Peers  []struct {
		NodeInfo struct {
			ID         string `json:"id"`
			Moniker    string `json:"moniker"`
			ListenAddr string `json:"listen_addr"`
			Network    string `json:"network"`
			Version    string `json:"version"`
		} `json:"node_info"`
		IsOutbound bool   `json:"is_outbound"`
		RemoteIP   string `json:"remote_ip"`
	} `json:"peers"`
}

// newPeersCmd builds `ops peers`. Default output is a table of
// node_id/remote_ip/moniker; --verbose defers to --json-style full
// peer structure; --json always wins and passes the raw /net_info
// response through.
func newPeersCmd() *cobra.Command {
	var verbose bool
	var fields []string
	cmd := &cobra.Command{
		Use:   "peers",
		Short: "List peers from CometBFT /net_info.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			rpc, cfg, err := newRPCClient(cmd)
			if err != nil {
				return err
			}
			raw, err := callEmpty(cmd.Context(), rpc, "net_info")
			if err != nil {
				return err
			}
			if raw == nil {
				return nil // --curl
			}
			opts := renderOpts(cmd, cfg, fields)
			if opts.JSON || verbose {
				generic, derr := decodeGeneric(raw)
				if derr != nil {
					return derr
				}
				// --verbose without --json renders the full decoded
				// struct as pretty JSON too; it's the only sane format
				// for a peer list with per-peer connection metrics.
				return render.RenderJSON(cmd.OutOrStdout(), generic, render.Options{
					JSON: true, Raw: cfg.Raw, Fields: fields,
					Color: cfg.Color, IsTTY: opts.IsTTY,
				})
			}
			var ni cometNetInfo
			if err := json.Unmarshal(raw, &ni); err != nil {
				return fmt.Errorf("decoding net_info: %w", err)
			}
			records := make([]map[string]any, 0, len(ni.Peers))
			for _, p := range ni.Peers {
				direction := "inbound"
				if p.IsOutbound {
					direction = "outbound"
				}
				records = append(records, map[string]any{
					"node_id":   p.NodeInfo.ID,
					"remote_ip": p.RemoteIP,
					"moniker":   p.NodeInfo.Moniker,
					"direction": direction,
					"version":   p.NodeInfo.Version,
				})
			}
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "n_peers  %s\n", ni.NPeers); err != nil {
				return err
			}
			return render.RenderTable(cmd.OutOrStdout(), records, opts)
		},
	}
	f := cmd.Flags()
	f.BoolVar(&verbose, "verbose", false, "emit full per-peer JSON (connection metrics, channels, etc)")
	f.StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable)")
	return cmd
}
