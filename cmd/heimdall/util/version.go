package heimdallutil

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/cmd/version"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newVersionCmd builds `util version`. Default output prints the
// polycli build metadata; --node additionally reaches the configured
// CometBFT RPC for /status and reports the remote node version.
func newVersionCmd() *cobra.Command {
	var contactNode bool
	var fields []string
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print polycli and (optionally) node version.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			out := map[string]any{
				"polycli_version": version.Version,
				"polycli_commit":  version.Commit,
				"polycli_built":   version.Date,
			}
			cfg, err := resolveIfPossible()
			if err == nil && cfg != nil {
				out["chain_id"] = cfg.ChainID
				out["network"] = cfg.Network
			}
			if contactNode {
				if err != nil {
					return err
				}
				info, nerr := fetchNodeStatus(cmd, cfg)
				if nerr != nil {
					return nerr
				}
				if info != nil {
					out["cometbft_version"] = info.NodeInfo.Version
					out["moniker"] = info.NodeInfo.Moniker
					out["network_id"] = info.NodeInfo.Network
					out["catching_up"] = info.SyncInfo.CatchingUp
					out["latest_block_height"] = info.SyncInfo.LatestBlockHeight
				}
			}
			opts := render.Options{
				JSON:   cfg != nil && cfg.JSON,
				Fields: fields,
				Color:  colorMode(cfg),
			}
			if opts.JSON {
				return render.RenderJSON(cmd.OutOrStdout(), out, opts)
			}
			return render.RenderKV(cmd.OutOrStdout(), out, opts)
		},
	}
	f := cmd.Flags()
	f.BoolVar(&contactNode, "node", false, "also fetch the connected node version via CometBFT /status")
	f.StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable)")
	return cmd
}

// resolveIfPossible returns a resolved *config.Config when the flag set
// has been wired in; otherwise it returns (nil, nil). The plain
// (no --node) version subcommand must work even without a fully
// configured flag set, so missing flags is not an error here.
func resolveIfPossible() (*config.Config, error) {
	if flags == nil {
		return nil, nil
	}
	return config.Resolve(flags)
}

func colorMode(cfg *config.Config) string {
	if cfg == nil {
		return "auto"
	}
	return cfg.Color
}

// nodeStatus is a minimal subset of the CometBFT /status response.
type nodeStatus struct {
	NodeInfo struct {
		Version string `json:"version"`
		Moniker string `json:"moniker"`
		Network string `json:"network"`
	} `json:"node_info"`
	SyncInfo struct {
		LatestBlockHeight string `json:"latest_block_height"`
		CatchingUp        bool   `json:"catching_up"`
	} `json:"sync_info"`
}

// fetchNodeStatus dials the configured CometBFT RPC and decodes its
// /status response. Returns (nil, nil) when running under --curl.
func fetchNodeStatus(cmd *cobra.Command, cfg *config.Config) (*nodeStatus, error) {
	if cfg == nil {
		return nil, &client.UsageError{Msg: "cannot contact node: config not resolved"}
	}
	rpc := client.NewRPCClient(cfg.RPCURL, cfg.Timeout, cfg.RPCHeaders, cfg.Insecure)
	if cfg.Curl {
		rpc.Transport = &client.CurlTransport{Out: cmd.OutOrStdout(), Headers: cfg.RPCHeaders}
	}
	raw, err := rpc.Call(cmd.Context(), "status", nil)
	if err != nil {
		return nil, fmt.Errorf("fetching status: %w", err)
	}
	if raw == nil {
		return nil, nil
	}
	var st nodeStatus
	if err := json.Unmarshal(raw, &st); err != nil {
		return nil, fmt.Errorf("decoding status: %w", err)
	}
	return &st, nil
}
