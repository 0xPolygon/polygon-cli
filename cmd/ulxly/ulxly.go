// Package ulxly provides utilities for interacting with the uLxLy bridge.
package ulxly

import (
	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/cmd/ulxly/bridge"
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/claim"
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/events"
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/proof"
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/tree"
)

// ULxLyCmd is the root command for ulxly utilities.
var ULxLyCmd = &cobra.Command{
	Use:   "ulxly",
	Short: "Utilities for interacting with the uLxLy bridge.",
	Long:  "Basic utility commands for interacting with the bridge contracts, bridge services, and generating proofs.",
	Args:  cobra.NoArgs,
}

func init() {
	// Bridge and claim commands
	ULxLyCmd.AddCommand(bridge.BridgeCmd)
	ULxLyCmd.AddCommand(claim.ClaimCmd)
	ULxLyCmd.AddCommand(claim.ClaimEverythingCmd)

	// Proof commands
	ULxLyCmd.AddCommand(proof.ProofCmd)
	ULxLyCmd.AddCommand(proof.RollupsProofCmd)
	ULxLyCmd.AddCommand(proof.EmptyProofCmd)
	ULxLyCmd.AddCommand(proof.ZeroProofCmd)

	// Event commands
	ULxLyCmd.AddCommand(events.GetDepositCmd)
	ULxLyCmd.AddCommand(events.GetClaimCmd)
	ULxLyCmd.AddCommand(events.GetVerifyBatchesCmd)

	// Tree commands
	ULxLyCmd.AddCommand(tree.BalanceTreeCmd)
	ULxLyCmd.AddCommand(tree.NullifierTreeCmd)
	ULxLyCmd.AddCommand(tree.NullifierAndBalanceTreeCmd)
}
