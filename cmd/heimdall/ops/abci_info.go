package ops

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// cometABCIInfo is the subset of /abci_info used for the summary.
type cometABCIInfo struct {
	Response struct {
		Data             string `json:"data"`
		Version          string `json:"version"`
		LastBlockHeight  string `json:"last_block_height"`
		LastBlockAppHash string `json:"last_block_app_hash"`
	} `json:"response"`
}

// newABCIInfoCmd builds `ops abci-info`. Default output is a KV
// summary of the ABCI-reported app identity and latest block hash.
func newABCIInfoCmd() *cobra.Command {
	var fields []string
	cmd := &cobra.Command{
		Use:   "abci-info",
		Short: "Show CometBFT /abci_info app identity.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			rpc, cfg, err := newRPCClient(cmd)
			if err != nil {
				return err
			}
			raw, err := callEmpty(cmd.Context(), rpc, "abci_info")
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
			var info cometABCIInfo
			if err := json.Unmarshal(raw, &info); err != nil {
				return fmt.Errorf("decoding abci_info: %w", err)
			}
			out := map[string]any{
				"app":                 info.Response.Data,
				"version":             info.Response.Version,
				"last_block_height":   info.Response.LastBlockHeight,
				"last_block_app_hash": info.Response.LastBlockAppHash,
			}
			return render.RenderKV(cmd.OutOrStdout(), out, opts)
		},
	}
	f := cmd.Flags()
	f.StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable)")
	return cmd
}
