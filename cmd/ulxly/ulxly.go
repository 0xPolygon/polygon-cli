package ulxly

import (
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/bridge"
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/claim"
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/common"
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/events"
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/proof"
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/tree"
	"github.com/0xPolygon/polygon-cli/flag"
	"github.com/spf13/cobra"
)

var ULxLyCmd = &cobra.Command{
	Use:   "ulxly",
	Short: "Utilities for interacting with the uLxLy bridge.",
	Long:  "Basic utility commands for interacting with the bridge contracts, bridge services, and generating proofs.",
	Args:  cobra.NoArgs,
}

// Hidden parent command for bridge and claim to share flags
var ulxlyBridgeAndClaimCmd = &cobra.Command{
	Args:   cobra.NoArgs,
	Hidden: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		common.InputUlxlyArgs.RPCURL, err = flag.GetRequiredRPCURL(cmd)
		if err != nil {
			return err
		}
		common.InputUlxlyArgs.PrivateKey, err = flag.GetRequiredPrivateKey(cmd)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	// Arguments for both bridge and claim
	fBridgeAndClaim := ulxlyBridgeAndClaimCmd.PersistentFlags()
	fBridgeAndClaim.StringVar(&common.InputUlxlyArgs.RPCURL, common.ArgRPCURL, "", "RPC URL to send the transaction")
	fBridgeAndClaim.StringVar(&common.InputUlxlyArgs.BridgeAddress, common.ArgBridgeAddress, "", "address of the lxly bridge")
	fBridgeAndClaim.Uint64Var(&common.InputUlxlyArgs.GasLimit, common.ArgGasLimit, 0, "force specific gas limit for transaction")
	fBridgeAndClaim.StringVar(&common.InputUlxlyArgs.ChainID, common.ArgChainID, "", "chain ID to use in the transaction")
	fBridgeAndClaim.StringVar(&common.InputUlxlyArgs.PrivateKey, common.ArgPrivateKey, "", "hex encoded private key for sending transaction")
	fBridgeAndClaim.StringVar(&common.InputUlxlyArgs.DestAddress, common.ArgDestAddress, "", "destination address for the bridge")
	fBridgeAndClaim.Uint64Var(&common.InputUlxlyArgs.Timeout, common.ArgTimeout, 60, "timeout in seconds to wait for transaction receipt confirmation")
	fBridgeAndClaim.StringVar(&common.InputUlxlyArgs.GasPrice, common.ArgGasPrice, "", "gas price to use")
	fBridgeAndClaim.BoolVar(&common.InputUlxlyArgs.DryRun, common.ArgDryRun, false, "do all of the transaction steps but do not send the transaction")
	fBridgeAndClaim.BoolVar(&common.InputUlxlyArgs.Insecure, common.ArgInsecure, false, "skip TLS certificate verification")
	fBridgeAndClaim.BoolVar(&common.InputUlxlyArgs.Legacy, common.ArgLegacy, true, "force usage of legacy bridge service")
	flag.MarkPersistentFlagsRequired(ulxlyBridgeAndClaimCmd, common.ArgBridgeAddress)

	// Bridge and Claim subcommands under hidden parent
	ulxlyBridgeAndClaimCmd.AddCommand(bridge.BridgeCmd)
	ulxlyBridgeAndClaimCmd.AddCommand(claim.ClaimCmd)
	ulxlyBridgeAndClaimCmd.AddCommand(claim.ClaimEverythingCmd)

	// Add hidden parent to root
	ULxLyCmd.AddCommand(ulxlyBridgeAndClaimCmd)

	// Add bridge and claim directly to root (so they're visible in help)
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
