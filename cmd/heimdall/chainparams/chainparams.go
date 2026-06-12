// Package chainparams implements the `polycli heimdall chainmanager`
// umbrella command (alias `cm`) and its subcommands targeting Heimdall
// v2's `x/chainmanager` module.
//
// Per HEIMDALLCAST_REQUIREMENTS.md §3.2.7 the chainmanager module holds
// the L1/L2 chain ids, tx confirmation depths, and L1 contract
// addresses. Upstream exposes a single HTTP route
// (`/chainmanager/params`, confirmed in
// heimdall-v2/proto/heimdallv2/chainmanager/query.proto); the
// `addresses` subcommand is a derived view over the same response.
//
// Package directory is named `chainparams` (not `chain`) because the
// top-level `chain` command is already claimed by the CometBFT-facing
// cast-like commands in cmd/heimdall/chain.
package chainparams

import (
	_ "embed"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/cmdutil"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

//go:embed usage.md
var usage string

// pkg carries the package name and the flag struct injected via
// Register; cmdutil derives clients and render options from it. The
// name is "chainmanager" (the command name) so the "not registered"
// error reads the same as before the cmdutil extraction.
var pkg = &cmdutil.Pkg{Name: "chainmanager"}

// Cmd is the umbrella `chainmanager` command (alias `cm`).
// Subcommands are attached by Register.
var Cmd = &cobra.Command{
	Use:     "chainmanager",
	Aliases: []string{"cm"},
	Short:   "Query chainmanager module endpoints.",
	Long:    usage,
	Args:    cobra.NoArgs,
}

// Register attaches the chainmanager umbrella command and its
// subcommands to parent, wiring in the shared flag struct.
func Register(parent *cobra.Command, f *config.Flags) {
	pkg.Flags = f
	Cmd.AddCommand(
		newParamsCmd(),
		newAddressesCmd(),
	)
	render.EnableWatchTree(Cmd)
	parent.AddCommand(Cmd)
}
