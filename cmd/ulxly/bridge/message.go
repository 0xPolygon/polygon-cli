package bridge

import (
	_ "embed"
	"math/big"
	"strings"

	ulxlycommon "github.com/0xPolygon/polygon-cli/cmd/ulxly/common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

//go:embed bridgeMessageUsage.md
var bridgeMessageUsage string

// MessageCmd represents the bridge message command.
var MessageCmd = &cobra.Command{
	Use:          "message",
	Short:        "Send some ETH along with data from one chain to another chain.",
	Long:         bridgeMessageUsage,
	PreRunE:      ulxlycommon.PrepInputs,
	RunE:         bridgeMessage,
	SilenceUsage: true,
}

func bridgeMessage(cmd *cobra.Command, _ []string) error {
	bridgeAddress := ulxlycommon.InputArgs.BridgeAddress
	privateKey := ulxlycommon.InputArgs.PrivateKey
	gasLimit := ulxlycommon.InputArgs.GasLimit
	destinationAddress := ulxlycommon.InputArgs.DestAddress
	chainID := ulxlycommon.InputArgs.ChainID
	amount := ulxlycommon.InputArgs.Value
	tokenAddr := ulxlycommon.InputArgs.TokenAddress
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

	value, _ := big.NewInt(0).SetString(amount, 0)
	tokenAddress := common.HexToAddress(tokenAddr)
	callData := common.Hex2Bytes(strings.TrimPrefix(callDataString, "0x"))

	if tokenAddress == common.HexToAddress("0x0000000000000000000000000000000000000000") {
		auth.Value = value
	}

	bridgeTxn, err := bridgeV2.BridgeMessage(auth, destinationNetwork, toAddress, isForced, callData)
	if err = logAndReturnJSONError(cmd.Context(), client, bridgeTxn, auth, err); err != nil {
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
