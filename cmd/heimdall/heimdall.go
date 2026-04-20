// Package heimdall implements the `polycli heimdall` command group, a
// cast-like CLI for querying Heimdall v2 REST + CometBFT endpoints and
// broadcasting signed Heimdall transactions.
package heimdall

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/cmd/heimdall/chain"
	"github.com/0xPolygon/polygon-cli/cmd/heimdall/chainparams"
	"github.com/0xPolygon/polygon-cli/cmd/heimdall/checkpoint"
	"github.com/0xPolygon/polygon-cli/cmd/heimdall/clerk"
	"github.com/0xPolygon/polygon-cli/cmd/heimdall/decode"
	"github.com/0xPolygon/polygon-cli/cmd/heimdall/milestone"
	"github.com/0xPolygon/polygon-cli/cmd/heimdall/ops"
	"github.com/0xPolygon/polygon-cli/cmd/heimdall/span"
	"github.com/0xPolygon/polygon-cli/cmd/heimdall/topup"
	"github.com/0xPolygon/polygon-cli/cmd/heimdall/tx"
	heimdallutil "github.com/0xPolygon/polygon-cli/cmd/heimdall/util"
	"github.com/0xPolygon/polygon-cli/cmd/heimdall/validator"
	"github.com/0xPolygon/polygon-cli/cmd/heimdall/wallet"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
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
	chain.Register(HeimdallCmd, PersistentFlags)
	tx.Register(HeimdallCmd, PersistentFlags)
	checkpoint.Register(HeimdallCmd, PersistentFlags)
	span.Register(HeimdallCmd, PersistentFlags)
	milestone.Register(HeimdallCmd, PersistentFlags)
	validator.Register(HeimdallCmd, PersistentFlags)
	clerk.Register(HeimdallCmd, PersistentFlags)
	topup.Register(HeimdallCmd, PersistentFlags)
	chainparams.Register(HeimdallCmd, PersistentFlags)
	heimdallutil.Register(HeimdallCmd, PersistentFlags)
	ops.Register(HeimdallCmd, PersistentFlags)
	wallet.Register(HeimdallCmd, PersistentFlags)
	decode.Register(HeimdallCmd, PersistentFlags)
	wireExitCodes(HeimdallCmd)
}

// wireExitCodes walks the heimdall subcommand tree and wraps every
// RunE so that on failure the process exits with the cast-style code
// produced by client.ExitCode. Cobra's default machinery returns the
// error to `rootCmd.Execute` which then calls `os.Exit(1)`, collapsing
// every failure mode into the same exit code. Operators scripting
// against polycli want to distinguish node errors (1), network errors
// (2), usage errors (3), and signing errors (4); wireExitCodes is the
// single place that guarantees that contract.
//
// We additionally set SilenceUsage + SilenceErrors on the heimdall
// subtree so cobra does not print the usage blob and duplicate error
// line on a failing RunE. We print the error ourselves before exiting.
func wireExitCodes(root *cobra.Command) {
	root.SilenceUsage = true
	root.SilenceErrors = true
	var walk func(*cobra.Command)
	walk = func(c *cobra.Command) {
		c.SilenceUsage = true
		c.SilenceErrors = true
		if c.RunE != nil {
			orig := c.RunE
			c.RunE = func(cmd *cobra.Command, args []string) error {
				err := orig(cmd, args)
				if err == nil {
					return nil
				}
				// Print the error ourselves (cobra won't, because
				// SilenceErrors is set) and map it to the correct
				// cast-style exit code.
				fmt.Fprintf(cmd.ErrOrStderr(), "Error: %s\n", err.Error())
				code := client.ExitCode(err)
				if code == 0 {
					code = config.ExitNodeErr
				}
				os.Exit(code)
				// unreachable
				return err
			}
		}
		for _, sub := range c.Commands() {
			walk(sub)
		}
	}
	walk(root)
}
