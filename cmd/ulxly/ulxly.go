// Package ulxly provides utilities for interacting with the uLxLy bridge.
package ulxly

import (
	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/cmd/ulxly/bridge"
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/claim"
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/claim/everything"
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/events/claims"
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/events/deposits"
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/events/verifybatches"
	proofdeposits "github.com/0xPolygon/polygon-cli/cmd/ulxly/proof/deposits"
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/proof/empty"
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/proof/rollups"
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/proof/zero"
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/tree/balance"
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/tree/combined"
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/tree/nullifier"
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
	ULxLyCmd.AddCommand(everything.Cmd)

	// Proof commands
	ULxLyCmd.AddCommand(proofdeposits.Cmd)
	ULxLyCmd.AddCommand(rollups.Cmd)
	ULxLyCmd.AddCommand(empty.Cmd)
	ULxLyCmd.AddCommand(zero.Cmd)

	// Event commands
	ULxLyCmd.AddCommand(deposits.Cmd)
	ULxLyCmd.AddCommand(claims.Cmd)
	ULxLyCmd.AddCommand(verifybatches.Cmd)

	// Tree commands
	ULxLyCmd.AddCommand(balance.Cmd)
	ULxLyCmd.AddCommand(nullifier.Cmd)
	ULxLyCmd.AddCommand(combined.Cmd)
}
