// Package validator implements the `polycli heimdall validator`
// umbrella command (alias `val`) and its subcommands targeting Heimdall
// v2's `x/stake` module: set/validators, total-power, get, signer,
// status, proposer, proposers, is-old-stake-tx.
//
// Per HEIMDALLCAST_REQUIREMENTS.md §3.2.4 these endpoints live under a
// single umbrella, and the umbrella also accepts a bare integer
// (`validator 4`) as a shorthand for `validator get 4`. The top-level
// `validators` command is registered separately as an alias for
// `validator set`.
package validator

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
var pkg = &cmdutil.Pkg{Name: "validator"}

// ValidatorCmd is the umbrella `validator` command. Subcommands are
// attached by Register.
var ValidatorCmd = &cobra.Command{
	Use:     "validator [ID]",
	Aliases: []string{"val"},
	Short:   "Query stake module endpoints.",
	Long:    usage,
	Args:    cobra.MaximumNArgs(1),
	// Bare-id shorthand: `validator 4` forwards to `validator get 4`.
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		if _, err := strconv.ParseUint(args[0], 10, 64); err != nil {
			return &client.UsageError{Msg: fmt.Sprintf("unknown validator subcommand or id %q", args[0])}
		}
		return runGet(cmd, args[0])
	},
}

// ValidatorsCmd is the top-level `validators` alias for `validator set`.
// It is attached to the root heimdall command alongside ValidatorCmd so
// operators can type either form.
var ValidatorsCmd = &cobra.Command{
	Use:   "validators",
	Short: "Alias for `validator set`.",
	Args:  cobra.NoArgs,
	RunE:  runSet,
}

// setFlags keeps the shared flag state for `validator set` /
// `validators`. Both commands share RunE (runSet) and must therefore
// read from the same variables.
var setFlags = struct {
	sort   string
	limit  int
	fields []string
}{}

// Register attaches the validator umbrella command (and the top-level
// `validators` alias) to parent, wiring in the shared flag struct.
func Register(parent *cobra.Command, f *config.Flags) {
	pkg.Flags = f
	ValidatorCmd.AddCommand(
		newSetCmd(),
		newTotalPowerCmd(),
		newGetCmd(),
		newSignerCmd(),
		newStatusCmd(),
		newProposerCmd(),
		newProposersCmd(),
		newIsOldStakeTxCmd(),
	)
	// Attach shared flags to the top-level `validators` alias as well.
	attachSetFlags(ValidatorsCmd.Flags())
	// Read-only umbrella: wire `--watch` into every descendant plus
	// the top-level alias.
	render.EnableWatchTree(ValidatorCmd)
	render.EnableWatch(ValidatorsCmd)
	parent.AddCommand(ValidatorCmd)
	parent.AddCommand(ValidatorsCmd)
}
