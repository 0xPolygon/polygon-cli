package monitorv2

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/0xPolygon/polygon-cli/chainstore"
	"github.com/0xPolygon/polygon-cli/cmd/monitorv2/renderer"
	"github.com/0xPolygon/polygon-cli/indexer"

	_ "embed"
)

//go:embed monitorv2Usage.md
var usage string

var (
	rpcURL         string
	rendererType   string
)

var MonitorV2Cmd = &cobra.Command{
	Use:   "monitorv2",
	Short: "Monitor v2 command stub",
	Long:  usage,
	RunE: func(cmd *cobra.Command, args []string) error {
		if rpcURL == "" {
			return fmt.Errorf("--rpc-url is required")
		}

		// Create store
		store, err := chainstore.NewPassthroughStore(rpcURL)
		if err != nil {
			return fmt.Errorf("failed to create store: %w", err)
		}
		defer store.Close()

		// Create indexer
		idx := indexer.NewIndexer(store, indexer.DefaultConfig())

		// Start indexer first
		if err := idx.Start(); err != nil {
			return fmt.Errorf("failed to start indexer: %w", err)
		}
		defer idx.Stop()

		// Create renderer based on type
		var r renderer.Renderer
		switch rendererType {
		case "json":
			r = renderer.NewJSONRenderer(idx)
		case "tview", "tui":
			r = renderer.NewTviewRenderer(idx)
		default:
			return fmt.Errorf("unknown renderer type: %s (supported: json, tview, tui)", rendererType)
		}

		// Start rendering
		ctx := context.Background()
		return r.Start(ctx)
	},
}

func init() {
	MonitorV2Cmd.Flags().StringVar(&rpcURL, "rpc-url", "", "RPC endpoint URL (required)")
	MonitorV2Cmd.Flags().StringVar(&rendererType, "renderer", "json", "Renderer type (json, tview, tui)")
}