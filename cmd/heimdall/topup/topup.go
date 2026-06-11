// Package topup implements the `polycli heimdall topup` umbrella
// command and its subcommands targeting Heimdall v2's `x/topup`
// module: root, account, proof, verify, sequence, is-old.
//
// Per HEIMDALLCAST_REQUIREMENTS.md §3.2.6 these endpoints live under a
// single umbrella rather than at the top level of the heimdall tree.
//
// Endpoints confirmed from heimdall-v2
// proto/heimdallv2/topup/query.proto:
//
//   - GET /topup/dividend-account-root
//   - GET /topup/dividend-account/{address}
//   - GET /topup/account-proof/{address}
//   - GET /topup/account-proof/{address}/verify?proof=…
//   - GET /topup/sequence?tx_hash=…&log_index=…
//   - GET /topup/is-old-tx?tx_hash=…&log_index=…
//
// The `proof`, `sequence`, and `is-old` endpoints fan out to L1 on the
// server side; a Heimdall node without `eth_rpc_url` returns gRPC code
// 13, which we surface as an L1-not-configured hint on stderr before
// propagating the error.
package topup

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
var pkg = &cmdutil.Pkg{Name: "topup"}

// TopupCmd is the umbrella `topup` command. Subcommands are attached
// by Register.
var TopupCmd = &cobra.Command{
	Use:   "topup",
	Short: "Query topup (dividend account) module endpoints.",
	Long:  usage,
	Args:  cobra.NoArgs,
}

// Register attaches the topup umbrella command and all of its
// subcommands to parent, wiring in the shared flag struct.
func Register(parent *cobra.Command, f *config.Flags) {
	pkg.Flags = f
	TopupCmd.AddCommand(
		newRootCmd(),
		newAccountCmd(),
		newProofCmd(),
		newVerifyCmd(),
		newSequenceCmd(),
		newIsOldCmd(),
	)
	render.EnableWatchTree(TopupCmd)
	parent.AddCommand(TopupCmd)
}
