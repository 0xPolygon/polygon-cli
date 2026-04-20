// Package heimdallutil implements the `polycli heimdall util` umbrella
// command and its local utility subcommands (addr, b64, version,
// completions). These helpers are deliberately offline — only `version
// --node` reaches the network.
//
// Package directory is named `util` under cmd/heimdall/ and the Go
// package is `heimdallutil` to avoid colliding with the repo's
// top-level util/ package.
package heimdallutil

import (
	_ "embed"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
)

//go:embed usage.md
var usage string

// flags is injected by Register. The version subcommand reads it via
// config.Resolve to dial the CometBFT RPC when --node is supplied.
var flags *config.Flags

// Cmd is the `util` umbrella command. Subcommands are attached by
// Register so tests can wire their own parent for isolation.
var Cmd = &cobra.Command{
	Use:   "util",
	Short: "Local helpers for addresses, base64, versions, and completions.",
	Long:  usage,
	Args:  cobra.NoArgs,
}

// Register attaches the util umbrella command and its subcommands to
// parent, wiring in the shared flag struct.
func Register(parent *cobra.Command, f *config.Flags) {
	flags = f
	Cmd.AddCommand(
		newAddrCmd(),
		newB64Cmd(),
		newVersionCmd(),
		newCompletionsCmd(parent),
	)
	parent.AddCommand(Cmd)
}
