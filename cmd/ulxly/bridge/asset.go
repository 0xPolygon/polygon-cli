package bridge

import (
	_ "embed"
	"math/big"
	"strings"

	"github.com/0xPolygon/polygon-cli/bindings/tokens"
	ulxlycommon "github.com/0xPolygon/polygon-cli/cmd/ulxly/common"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

//go:embed bridgeAssetUsage.md
var bridgeAssetUsage string

// AssetCmd represents the bridge asset command.
var AssetCmd = &cobra.Command{
	Use:          "asset",
	Short:        "Move ETH or an ERC20 between to chains.",
	Long:         bridgeAssetUsage,
	PreRunE:      ulxlycommon.PrepInputs,
	RunE:         bridgeAsset,
	SilenceUsage: true,
}

func bridgeAsset(cmd *cobra.Command, _ []string) error {
	bridgeAddr := ulxlycommon.InputArgs.BridgeAddress
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

	client, err := ulxlycommon.CreateEthClient(cmd.Context(), rpcURL)
	if err != nil {
		log.Error().Err(err).Msg("Unable to Dial RPC")
		return err
	}
	defer client.Close()

	// Initialize and assign variables required to send transaction payload
	bridgeV2, toAddress, auth, err := ulxlycommon.GenerateTransactionPayload(cmd.Context(), client, bridgeAddr, privateKey, gasLimit, destinationAddress, chainID)
	if err != nil {
		log.Error().Err(err).Msg("error generating transaction payload")
		return err
	}

	bridgeAddress := common.HexToAddress(bridgeAddr)
	value, _ := big.NewInt(0).SetString(amount, 0)
	tokenAddress := common.HexToAddress(tokenAddr)
	callData := common.Hex2Bytes(strings.TrimPrefix(callDataString, "0x"))

	if tokenAddress == common.HexToAddress("0x0000000000000000000000000000000000000000") {
		auth.Value = value
	} else {
		// in case it's a token transfer, we need to ensure that the bridge contract
		// has enough allowance to transfer the tokens on behalf of the user
		tokenContract, iErr := tokens.NewERC20(tokenAddress, client)
		if iErr != nil {
			log.Error().Err(iErr).Msg("error getting token contract")
			return iErr
		}

		allowance, iErr := tokenContract.Allowance(&bind.CallOpts{Pending: false}, auth.From, bridgeAddress)
		if iErr != nil {
			log.Error().Err(iErr).Msg("error getting token allowance")
			return iErr
		}

		if allowance.Cmp(value) < 0 {
			log.Info().
				Str("amount", value.String()).
				Str("tokenAddress", tokenAddress.String()).
				Str("bridgeAddress", bridgeAddress.String()).
				Str("userAddress", auth.From.String()).
				Msg("approving bridge contract to spend tokens on behalf of user")

			// Approve the bridge contract to spend the tokens on behalf of the user
			approveTxn, iErr := tokenContract.Approve(auth, bridgeAddress, value)
			if iErr = ulxlycommon.LogAndReturnJSONError(cmd.Context(), client, approveTxn, auth, iErr); iErr != nil {
				return iErr
			}
			log.Info().Msg("approveTxn: " + approveTxn.Hash().String())
			if iErr = ulxlycommon.WaitMineTransaction(cmd.Context(), client, approveTxn, timeoutTxnReceipt); iErr != nil {
				return iErr
			}
		}
	}

	bridgeTxn, err := bridgeV2.BridgeAsset(auth, destinationNetwork, toAddress, value, tokenAddress, isForced, callData)
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
