package monitorv2

import (
	"context"
	"fmt"
	"net/http"

	"github.com/0xPolygon/polygon-cli/chainstore"
	"github.com/0xPolygon/polygon-cli/cmd/monitorv2/renderer"
	"github.com/0xPolygon/polygon-cli/indexer"
	"github.com/0xPolygon/polygon-cli/util"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	_ "embed"
	_ "net/http/pprof" // Import pprof HTTP handlers
)

//go:embed monitorv2Usage.md
var usage string

var (
	rpcURL       string
	rendererType string
	pprofAddr    string
)

var MonitorV2Cmd = &cobra.Command{
	Use:   "monitorv2",
	Short: "Monitor v2 command stub",
	Long:  usage,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Set default verbosity to Error level (300) if not explicitly set by user
		verbosityFlag := cmd.Flag("verbosity")
		if verbosityFlag != nil && !verbosityFlag.Changed {
			util.SetLogLevel(300) // Error level
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if rpcURL == "" {
			return fmt.Errorf("--rpc-url is required")
		}

		// Start pprof server if requested
		if pprofAddr != "" {
			go func() {
				log.Info().Str("addr", pprofAddr).Msg("Starting pprof server")
				if err := http.ListenAndServe(pprofAddr, nil); err != nil {
					log.Error().Err(err).Msg("pprof server failed")
				}
			}()
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
	MonitorV2Cmd.Flags().StringVar(&rendererType, "renderer", "tui", "Renderer type (json, tview, tui)")
	MonitorV2Cmd.Flags().StringVar(&pprofAddr, "pprof", "", "Enable pprof server on specified address (e.g. 127.0.0.1:6060)")
}
