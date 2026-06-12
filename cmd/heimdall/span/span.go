// Package span implements the `polycli heimdall span` umbrella command
// (alias `sp`) and its subcommands targeting Heimdall v2's `x/bor`
// module: params, latest, get, list, producers, seed, votes, downtime,
// scores, find.
//
// Per HEIMDALLCAST_REQUIREMENTS.md §3.2.2 these endpoints live under a
// single umbrella rather than at the top level. The umbrella also
// accepts a bare integer (`span 5982`) as a shorthand for `span get
// 5982`.
package span

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
var pkg = &cmdutil.Pkg{Name: "span"}

// SpanCmd is the umbrella `span` command. Subcommands are attached by
// Register.
var SpanCmd = &cobra.Command{
	Use:     "span [ID]",
	Aliases: []string{"sp"},
	Short:   "Query bor/span module endpoints.",
	Long:    usage,
	Args:    cobra.MaximumNArgs(1),
	// Bare-id shorthand: `span 5982` forwards to `span get 5982`.
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		if _, err := strconv.ParseUint(args[0], 10, 64); err != nil {
			return &client.UsageError{Msg: fmt.Sprintf("unknown span subcommand or id %q", args[0])}
		}
		return runGet(cmd, args[0])
	},
}

// Register attaches the span umbrella command and all of its
// subcommands to parent, wiring in the shared flag struct.
//
// Every span subcommand is read-only, so we apply render.EnableWatchTree
// once here and every descendant gets a `--watch DURATION` flag.
func Register(parent *cobra.Command, f *config.Flags) {
	pkg.Flags = f
	SpanCmd.AddCommand(
		newParamsCmd(),
		newLatestCmd(),
		newGetCmd(),
		newListCmd(),
		newProducersCmd(),
		newSeedCmd(),
		newVotesCmd(),
		newDowntimeCmd(),
		newScoresCmd(),
		newFindCmd(),
	)
	render.EnableWatchTree(SpanCmd)
	parent.AddCommand(SpanCmd)
}

// parseSpanID validates a CLI-provided span/validator/producer id and
// returns it as uint64. We accept only unsigned base-10 integers.
func parseSpanID(label, raw string) (uint64, error) {
	v, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, &client.UsageError{Msg: fmt.Sprintf("%s must be a positive integer, got %q", label, raw)}
	}
	return v, nil
}
