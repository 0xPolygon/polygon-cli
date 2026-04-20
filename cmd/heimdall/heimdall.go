// Package heimdall implements the `polycli heimdall` command group, a
// cast-like CLI for querying Heimdall v2 REST + CometBFT endpoints and
// broadcasting signed Heimdall transactions.
package heimdall

import (
	_ "embed"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
)

//go:embed usage.md
var usage string

// PersistentFlags holds the raw flag state shared across every
// heimdall subcommand. Subcommand RunE functions call
// config.Resolve(&PersistentFlags) to obtain a fully resolved
// *config.Config.
var PersistentFlags = &config.Flags{}

// HeimdallCmd is the root command for the heimdall subcommand tree.
var HeimdallCmd = &cobra.Command{
	Use:     "heimdall",
	Aliases: []string{"h"},
	Short:   "Query and interact with a Heimdall v2 node.",
	Long:    usage,
	Args:    cobra.NoArgs,
}

func init() {
	PersistentFlags.Register(HeimdallCmd)
}
