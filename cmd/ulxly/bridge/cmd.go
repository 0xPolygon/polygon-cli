// Package bridge provides commands for bridging assets and messages between chains.
package bridge

import (
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/bridge/asset"
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/bridge/message"
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/bridge/weth"
	ulxlycommon "github.com/0xPolygon/polygon-cli/cmd/ulxly/common"
	"github.com/0xPolygon/polygon-cli/flag"
	"github.com/spf13/cobra"
)

// BridgeCmd is the parent command for all bridge operations.
var BridgeCmd = &cobra.Command{
	Use:   "bridge",
	Short: "Commands for moving funds and sending messages from one chain to another.",
	Args:  cobra.NoArgs,
}

func init() {
	// Add subcommands
	BridgeCmd.AddCommand(asset.Cmd)
	BridgeCmd.AddCommand(message.Cmd)
	BridgeCmd.AddCommand(weth.Cmd)

	// Add shared transaction flags (rpc-url, bridge-address, private-key, etc.)
	ulxlycommon.AddTransactionFlags(BridgeCmd)

	// Bridge-specific persistent flags
	f := BridgeCmd.PersistentFlags()
	f.BoolVar(&ulxlycommon.InputArgs.ForceUpdate, ulxlycommon.ArgForceUpdate, true, "update the new global exit root")
	f.StringVar(&ulxlycommon.InputArgs.Value, ulxlycommon.ArgValue, "0", "amount in wei to send with the transaction")
	f.Uint32Var(&ulxlycommon.InputArgs.DestNetwork, ulxlycommon.ArgDestNetwork, 0, "rollup ID of the destination network")
	f.StringVar(&ulxlycommon.InputArgs.TokenAddress, ulxlycommon.ArgTokenAddress, "0x0000000000000000000000000000000000000000", "address of ERC20 token to use")
	f.StringVar(&ulxlycommon.InputArgs.CallData, ulxlycommon.ArgCallData, "0x", "call data to be passed directly with bridge-message or as an ERC20 Permit")
	f.StringVar(&ulxlycommon.InputArgs.CallDataFile, ulxlycommon.ArgCallDataFile, "", "a file containing hex encoded call data")
	flag.MarkPersistentFlagsRequired(BridgeCmd, ulxlycommon.ArgDestNetwork)
}
