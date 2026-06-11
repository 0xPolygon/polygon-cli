// Package checkpoint implements the `polycli heimdall checkpoint`
// umbrella command (alias `cp`) and its subcommands: params, count,
// latest, get, buffer, last-no-ack, next, list, signatures, overview.
//
// Per HEIMDALLCAST_REQUIREMENTS.md §3.2.1 these endpoints live under
// a single umbrella rather than at the top level of the heimdall
// tree. The umbrella also accepts a bare integer (`checkpoint 38871`)
// as a shorthand for `checkpoint get 38871`.
package checkpoint

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"

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
var pkg = &cmdutil.Pkg{Name: "checkpoint"}

// CheckpointCmd is the umbrella `checkpoint` command. Subcommands are
// attached by Register.
var CheckpointCmd = &cobra.Command{
	Use:     "checkpoint [ID]",
	Aliases: []string{"cp"},
	Short:   "Query checkpoint module endpoints.",
	Long:    usage,
	Args:    cobra.MaximumNArgs(1),
	// Bare-id shorthand: `checkpoint 38871` forwards to `checkpoint get 38871`.
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		if _, err := strconv.ParseUint(args[0], 10, 64); err != nil {
			return &client.UsageError{Msg: fmt.Sprintf("unknown checkpoint subcommand or id %q", args[0])}
		}
		return runGet(cmd, args[0])
	},
}

// Register attaches the checkpoint umbrella command and all of its
// subcommands to parent, wiring in the shared flag struct.
//
// Every checkpoint subcommand is read-only, so we apply
// render.EnableWatchTree once here.
func Register(parent *cobra.Command, f *config.Flags) {
	pkg.Flags = f
	CheckpointCmd.AddCommand(
		newParamsCmd(),
		newCountCmd(),
		newLatestCmd(),
		newGetCmd(),
		newBufferCmd(),
		newLastNoAckCmd(),
		newNextCmd(),
		newListCmd(),
		newSignaturesCmd(),
		newOverviewCmd(),
	)
	render.EnableWatchTree(CheckpointCmd)
	parent.AddCommand(CheckpointCmd)
}

// normalizeCheckpointHash accepts a checkpoint tx hash with or without
// the `0x` prefix and returns the lower-case, unprefixed hex form
// expected by /checkpoints/signatures/{hash} on Heimdall. Returns a
// UsageError for non-hex or non-32-byte inputs. The validation lives
// in cmdutil.NormalizeTxHash; this endpoint just wants the prefix
// stripped back off.
func normalizeCheckpointHash(raw string) (string, error) {
	h, err := cmdutil.NormalizeTxHash(raw)
	if err != nil {
		return "", err
	}
	return strings.TrimPrefix(h, "0x"), nil
}
