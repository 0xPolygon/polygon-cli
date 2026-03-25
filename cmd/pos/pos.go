// Package pos provides utilities for interacting with Polygon PoS.
package pos

import (
	"github.com/0xPolygon/polygon-cli/cmd/pos/exitproof"
	"github.com/spf13/cobra"
)

// POSCmd is the root command for Polygon PoS utilities.
var POSCmd = &cobra.Command{
	Use:   "pos",
	Short: "Utilities for Polygon PoS.",
	Long:  "Commands for generating exit proofs and other Polygon PoS-specific operations.",
	Args:  cobra.NoArgs,
}

func init() {
	POSCmd.AddCommand(exitproof.Cmd)
}
