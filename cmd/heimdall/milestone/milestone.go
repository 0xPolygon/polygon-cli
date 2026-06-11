// Package milestone implements the `polycli heimdall milestone`
// umbrella command (alias `ms`) and its subcommands targeting Heimdall
// v2's `x/milestone` module: params, count, latest, get.
//
// Per HEIMDALLCAST_REQUIREMENTS.md §3.2.3 these endpoints live under a
// single umbrella rather than at the top level. The umbrella also
// accepts a bare integer (`milestone 11602043`) as a shorthand for
// `milestone get 11602043`.
package milestone

import (
	_ "embed"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/cmdutil"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

//go:embed usage.md
var usage string

// pkg carries the package name and the flag struct injected via
// Register; cmdutil derives clients and render options from it.
var pkg = &cmdutil.Pkg{Name: "milestone"}

// MilestoneCmd is the umbrella `milestone` command. Subcommands are
// attached by Register.
var MilestoneCmd = &cobra.Command{
	Use:     "milestone [NUMBER]",
	Aliases: []string{"ms"},
	Short:   "Query milestone module endpoints.",
	Long:    usage,
	Args:    cobra.MaximumNArgs(1),
	// Bare-number shorthand: `milestone 11602043` forwards to
	// `milestone get 11602043`.
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		if _, err := strconv.ParseUint(args[0], 10, 64); err != nil {
			return &client.UsageError{Msg: fmt.Sprintf("unknown milestone subcommand or number %q", args[0])}
		}
		return runGet(cmd, args[0])
	},
}

// Register attaches the milestone umbrella command and all of its
// subcommands to parent, wiring in the shared flag struct.
func Register(parent *cobra.Command, f *config.Flags) {
	pkg.Flags = f
	MilestoneCmd.AddCommand(
		newParamsCmd(),
		newCountCmd(),
		newLatestCmd(),
		newGetCmd(),
		newVotesCmd(),
	)
	render.EnableWatchTree(MilestoneCmd)
	parent.AddCommand(MilestoneCmd)
}
