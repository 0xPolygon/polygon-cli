// Package claim provides commands for claiming deposits on a particular chain.
package claim

import (
	"time"

	"github.com/0xPolygon/polygon-cli/cmd/ulxly/claim/asset"
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/claim/message"
	ulxlycommon "github.com/0xPolygon/polygon-cli/cmd/ulxly/common"
	"github.com/0xPolygon/polygon-cli/flag"
	"github.com/spf13/cobra"
)

var ClaimCmd = &cobra.Command{
	Use:   "claim",
	Short: "Commands for claiming deposits on a particular chain.",
	Args:  cobra.NoArgs,
}

func init() {
	ClaimCmd.AddCommand(asset.Cmd)
	ClaimCmd.AddCommand(message.Cmd)

	// Add shared transaction flags (rpc-url, bridge-address, private-key, etc.)
	ulxlycommon.AddTransactionFlags(ClaimCmd)

	// Claim-specific persistent flags
	f := ClaimCmd.PersistentFlags()
	f.Uint32Var(&ulxlycommon.InputArgs.DepositCount, ulxlycommon.ArgDepositCount, 0, "deposit count of the bridge transaction")
	f.Uint32Var(&ulxlycommon.InputArgs.DepositNetwork, ulxlycommon.ArgDepositNetwork, 0, "rollup ID of the network where the deposit was made")
	f.StringVar(&ulxlycommon.InputArgs.BridgeServiceURL, ulxlycommon.ArgBridgeServiceURL, "", "URL of the bridge service")
	f.StringVar(&ulxlycommon.InputArgs.GlobalIndex, ulxlycommon.ArgGlobalIndex, "", "an override of the global index value")
	f.DurationVar(&ulxlycommon.InputArgs.Wait, ulxlycommon.ArgWait, time.Duration(0), "retry claiming until deposit is ready, up to specified duration (available for claim asset and claim message)")
	f.StringVar(&ulxlycommon.InputArgs.ProofGER, ulxlycommon.ArgProofGER, "", "if specified and using legacy mode, the proof will be generated against this GER")
	f.Uint32Var(&ulxlycommon.InputArgs.ProofL1InfoTreeIndex, ulxlycommon.ArgProofL1InfoTreeIndex, 0, "if specified and using aggkit mode, the proof will be generated against this L1 Info Tree Index")
	flag.MarkPersistentFlagsRequired(ClaimCmd, ulxlycommon.ArgDepositCount, ulxlycommon.ArgDepositNetwork, ulxlycommon.ArgBridgeServiceURL)
	ClaimCmd.MarkFlagsMutuallyExclusive(ulxlycommon.ArgProofGER, ulxlycommon.ArgProofL1InfoTreeIndex)
}
