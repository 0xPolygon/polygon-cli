package heimdallutil

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
)

// supportedShells enumerates the shells cobra can emit completions for.
// Kept in stable order for help text.
var supportedShells = []string{"bash", "zsh", "fish", "powershell"}

// newCompletionsCmd builds `util completions <shell>`. At invocation
// time it walks up to the *root* cobra command and uses cobra's
// GenXCompletion helpers so the completions cover the whole polycli
// tree (matching the stock `polycli completion <shell>` behaviour).
//
// The parent arg is accepted for symmetry with the other subcommand
// builders in this package; the root is located at run-time from the
// live command's parent chain rather than captured at construction
// time, because tests build their own parents per-case.
func newCompletionsCmd(_ *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:       "completions <shell>",
		Short:     "Generate shell completion script.",
		Args:      cobra.ExactArgs(1),
		ValidArgs: supportedShells,
		RunE: func(cmd *cobra.Command, args []string) error {
			shell := strings.ToLower(args[0])
			root := rootOf(cmd)
			w := cmd.OutOrStdout()
			switch shell {
			case "bash":
				return root.GenBashCompletionV2(w, true)
			case "zsh":
				return root.GenZshCompletion(w)
			case "fish":
				return root.GenFishCompletion(w, true)
			case "powershell":
				return root.GenPowerShellCompletionWithDesc(w)
			default:
				return &client.UsageError{Msg: fmt.Sprintf(
					"unsupported shell %q (want one of: %s)",
					args[0], strings.Join(supportedShells, ", "))}
			}
		},
	}
}

// rootOf walks up the cobra parent chain to find the top-level command.
// Returns cmd itself if it has no parent (only in tests that drive a
// standalone subcommand).
func rootOf(cmd *cobra.Command) *cobra.Command {
	for c := cmd; c != nil; c = c.Parent() {
		if c.Parent() == nil {
			return c
		}
	}
	return cmd
}
