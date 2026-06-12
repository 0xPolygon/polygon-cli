// Package ops implements the `polycli heimdall ops` umbrella command
// and its CometBFT JSON-RPC-facing subcommands: status, health, peers,
// consensus, tx-pool, abci-info, commit, and validators-cometbft.
//
// All calls target the CometBFT RPC endpoint (`:26657`). The
// Heimdall REST gateway is unused here. The umbrella keeps node-ops
// commands grouped so operators can find them without scrolling
// through the flat cast-like tree.
package ops

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
// Register; cmdutil derives clients and render options from it.
var pkg = &cmdutil.Pkg{Name: "ops"}

// newOpsCmd builds a fresh `ops` umbrella. Constructed per Register
// call so tests that re-wire a parent do not accumulate duplicate
// subcommands on a shared command tree.
func newOpsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ops",
		Short: "Node-operator commands backed by CometBFT JSON-RPC.",
		Long:  usage,
		Args:  cobra.NoArgs,
	}
	cmd.AddCommand(
		newStatusCmd(),
		newHealthCmd(),
		newPeersCmd(),
		newConsensusCmd(),
		newTxPoolCmd(),
		newABCIInfoCmd(),
		newCommitCmd(),
		newValidatorsCometBFTCmd(),
	)
	return cmd
}

// Register attaches the ops umbrella and its subcommands to parent,
// wiring in the shared flag struct used for config resolution.
//
// Every ops subcommand is read-only, so we apply render.EnableWatchTree
// once so each descendant picks up `--watch DURATION`.
func Register(parent *cobra.Command, f *config.Flags) {
	pkg.Flags = f
	cmd := newOpsCmd()
	render.EnableWatchTree(cmd)
	parent.AddCommand(cmd)
}
