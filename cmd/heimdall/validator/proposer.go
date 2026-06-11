package validator

import (
	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/cmdutil"
)

// newProposerCmd builds `validator proposer` → GET
// /stake/proposers/current. Single validator envelope, same shape as
// `/stake/validator/{id}`; the { "validator": { … } } envelope is
// unwrapped for KV output.
func newProposerCmd() *cobra.Command {
	return pkg.NewGetCmd(cmdutil.Get{
		Use:       "proposer",
		Short:     "Show the current proposer.",
		Path:      "/stake/proposers/current",
		Label:     "proposer",
		UnwrapKey: "validator",
	})
}
