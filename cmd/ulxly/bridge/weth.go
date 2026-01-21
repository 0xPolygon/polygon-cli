package bridge

import (
	_ "embed"
	"fmt"
	"math/big"
	"strings"

	ulxlycommon "github.com/0xPolygon/polygon-cli/cmd/ulxly/common"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

//go:embed bridgeWethUsage.md
var bridgeWethUsage string

// WethCmd represents the bridge WETH command.
var WethCmd = &cobra.Command{
	Use:          "weth",
	Short:        "For L2's that use a gas token, use this to transfer WETH to another chain.",
	Long:         bridgeWethUsage,
	PreRunE:      ulxlycommon.PrepInputs,
	RunE:         bridgeWETHMessage,
	SilenceUsage: true,
}

func bridgeWETHMessage(cmd *cobra.Command, _ []string) error {
	bridgeAddress := ulxlycommon.InputArgs.BridgeAddress
	privateKey := ulxlycommon.InputArgs.PrivateKey
	gasLimit := ulxlycommon.InputArgs.GasLimit
	destinationAddress := ulxlycommon.InputArgs.DestAddress
	chainID := ulxlycommon.InputArgs.ChainID
	amount := ulxlycommon.InputArgs.Value
	callDataString := ulxlycommon.InputArgs.CallData
	destinationNetwork := ulxlycommon.InputArgs.DestNetwork
	isForced := ulxlycommon.InputArgs.ForceUpdate
	timeoutTxnReceipt := ulxlycommon.InputArgs.Timeout
	rpcURL := ulxlycommon.InputArgs.RPCURL

	// Dial the Ethereum RPC server.
	client, err := ulxlycommon.CreateEthClient(cmd.Context(), rpcURL)
	if err != nil {
		log.Error().Err(err).Msg("Unable to Dial RPC")
		return err
	}
	defer client.Close()

	// Initialize and assign variables required to send transaction payload
	bridgeV2, toAddress, auth, err := ulxlycommon.GenerateTransactionPayload(cmd.Context(), client, bridgeAddress, privateKey, gasLimit, destinationAddress, chainID)
	if err != nil {
		log.Error().Err(err).Msg("error generating transaction payload")
		return err
	}

	// Check if WETH is allowed
	wethAddress, err := bridgeV2.WETHToken(&bind.CallOpts{Pending: false})
	if err != nil {
		log.Error().Err(err).Msg("error getting WETH address from the bridge smc")
		return err
	}
	if wethAddress == (common.Address{}) {
		return fmt.Errorf("bridge WETH not allowed. Native ETH token configured in this network. This tx will fail")
	}

	value, _ := big.NewInt(0).SetString(amount, 0)
	callData := common.Hex2Bytes(strings.TrimPrefix(callDataString, "0x"))

	bridgeTxn, err := bridgeV2.BridgeMessageWETH(auth, destinationNetwork, toAddress, value, isForced, callData)
	if err = ulxlycommon.LogAndReturnJSONError(cmd.Context(), client, bridgeTxn, auth, err); err != nil {
		log.Info().Err(err).Str("calldata", callDataString).Msg("Bridge transaction failed")
		return err
	}
	log.Info().Msg("bridgeTxn: " + bridgeTxn.Hash().String())
	if err = ulxlycommon.WaitMineTransaction(cmd.Context(), client, bridgeTxn, timeoutTxnReceipt); err != nil {
		return err
	}
	depositCount, err := parseDepositCountFromTransaction(cmd.Context(), client, bridgeTxn.Hash(), bridgeV2)
	if err != nil {
		return err
	}

	log.Info().Uint32("depositCount", depositCount).Msg("Bridge deposit count parsed from logs")
	return nil
}
