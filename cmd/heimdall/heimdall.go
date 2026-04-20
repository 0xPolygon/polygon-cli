// Package heimdall implements the `polycli heimdall` command group, a
// cast-like CLI for querying Heimdall v2 REST + CometBFT endpoints and
// broadcasting signed Heimdall transactions.
package heimdall

import (
	_ "embed"

	"github.com/spf13/cobra"
)

//go:embed usage.md
var usage string

// HeimdallCmd is the root command for the heimdall subcommand tree.
var HeimdallCmd = &cobra.Command{
	Use:     "heimdall",
	Aliases: []string{"h"},
	Short:   "Query and interact with a Heimdall v2 node.",
	Long:    usage,
	Args:    cobra.NoArgs,
}
