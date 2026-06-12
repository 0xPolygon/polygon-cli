// Package clerk implements the `polycli heimdall state-sync` umbrella
// command (aliases `clerk` and `ss`) and its subcommands targeting
// Heimdall v2's `x/clerk` module: count, latest-id, get, list, range,
// sequence, is-old.
//
// Per HEIMDALLCAST_REQUIREMENTS.md §3.2.5 these endpoints live under a
// single umbrella rather than at the top level of the heimdall tree.
// The umbrella also accepts a bare integer (`state-sync 36610`) as a
// shorthand for `state-sync get 36610`.
//
// Pagination note: `/clerk/event-records/list` is page-based (page +
// limit query params), NOT Cosmos pagination, and rejects `page=0`
// with HTTP 400. `/clerk/time` is Cosmos-paginated (pagination.limit).
package clerk

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
var pkg = &cmdutil.Pkg{Name: "clerk"}

// ClerkCmd is the umbrella `state-sync` command (aliases `clerk`,
// `ss`). Subcommands are attached by Register.
var ClerkCmd = &cobra.Command{
	Use:     "state-sync [ID]",
	Aliases: []string{"clerk", "ss"},
	Short:   "Query state-sync (clerk) module endpoints.",
	Long:    usage,
	Args:    cobra.MaximumNArgs(1),
	// Bare-id shorthand: `state-sync 36610` forwards to `state-sync get 36610`.
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		if _, err := strconv.ParseUint(args[0], 10, 64); err != nil {
			return &client.UsageError{Msg: fmt.Sprintf("unknown state-sync subcommand or id %q", args[0])}
		}
		return runGet(cmd, args[0], false)
	},
}

// Register attaches the state-sync umbrella command and all of its
// subcommands to parent, wiring in the shared flag struct.
func Register(parent *cobra.Command, f *config.Flags) {
	pkg.Flags = f
	ClerkCmd.AddCommand(
		newCountCmd(),
		newLatestIDCmd(),
		newGetCmd(),
		newListCmd(),
		newRangeCmd(),
		newSequenceCmd(),
		newIsOldCmd(),
	)
	render.EnableWatchTree(ClerkCmd)
	parent.AddCommand(ClerkCmd)
}
