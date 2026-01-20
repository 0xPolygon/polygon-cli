package bridge

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/0xPolygon/polygon-cli/bindings/ulxly"
	ulxlycommon "github.com/0xPolygon/polygon-cli/cmd/ulxly/common"
	smcerror "github.com/0xPolygon/polygon-cli/errors"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
)

// logAndReturnJSONError logs and returns a JSON-RPC error with additional context.
func logAndReturnJSONError(ctx context.Context, client *ethclient.Client, tx *types.Transaction, opts *bind.TransactOpts, err error) error {
	var callErr error
	if tx != nil {
		// in case the error came down to gas estimation, we can sometimes get more information by doing a call
		_, callErr = client.CallContract(ctx, ethereum.CallMsg{
			From:          opts.From,
			To:            tx.To(),
			Gas:           tx.Gas(),
			GasPrice:      tx.GasPrice(),
			GasFeeCap:     tx.GasFeeCap(),
			GasTipCap:     tx.GasTipCap(),
			Value:         tx.Value(),
			Data:          tx.Data(),
			AccessList:    tx.AccessList(),
			BlobGasFeeCap: tx.BlobGasFeeCap(),
			BlobHashes:    tx.BlobHashes(),
		}, nil)

		if ulxlycommon.InputArgs.DryRun {
			castCmd := "cast call"
			castCmd += fmt.Sprintf(" --rpc-url %s", ulxlycommon.InputArgs.RPCURL)
			castCmd += fmt.Sprintf(" --from %s", opts.From.String())
			castCmd += fmt.Sprintf(" --gas-limit %d", tx.Gas())
			if tx.Type() == types.LegacyTxType {
				castCmd += fmt.Sprintf(" --gas-price %s", tx.GasPrice().String())
			} else {
				castCmd += fmt.Sprintf(" --gas-price %s", tx.GasFeeCap().String())
				castCmd += fmt.Sprintf(" --priority-gas-price %s", tx.GasTipCap().String())
			}
			castCmd += fmt.Sprintf(" --value %s", tx.Value().String())
			castCmd += fmt.Sprintf(" %s", tx.To().String())
			castCmd += fmt.Sprintf(" %s", common.Bytes2Hex(tx.Data()))
			log.Info().Str("cmd", castCmd).Msg("use this command to replicate the call")
		}
	}

	if err == nil {
		return nil
	}

	var jsonError ulxlycommon.JSONError
	jsonErrorBytes, jsErr := json.Marshal(err)
	if jsErr != nil {
		log.Error().Err(err).Msg("Unable to interact with the bridge contract")
		return err
	}

	jsErr = json.Unmarshal(jsonErrorBytes, &jsonError)
	if jsErr != nil {
		log.Error().Err(err).Msg("Unable to interact with the bridge contract")
		return err
	}

	reason, decodeErr := smcerror.DecodeSmcErrorCode(jsonError.Data)
	if decodeErr != nil {
		log.Error().Err(err).Msg("unable to decode smart contract error")
		return err
	}
	errLog := log.Error().
		Err(err).
		Str("message", jsonError.Message).
		Int("code", jsonError.Code).
		Interface("data", jsonError.Data).
		Str("reason", reason)

	if callErr != nil {
		errLog = errLog.Err(callErr)
	}

	customErr := errors.New(err.Error() + ": " + reason)
	if errCode, isValid := jsonError.Data.(string); isValid && errCode == "0x646cf558" {
		// I don't want to bother with the additional error logging for previously claimed deposits
		return customErr
	}

	errLog.Msg("Unable to interact with bridge contract")
	return customErr
}

// parseDepositCountFromTransaction extracts the deposit count from a bridge transaction receipt.
func parseDepositCountFromTransaction(ctx context.Context, client *ethclient.Client, txHash common.Hash, bridgeContract *ulxly.Ulxly) (uint32, error) {
	receipt, err := client.TransactionReceipt(ctx, txHash)
	if err != nil {
		return 0, err
	}

	// Check if the transaction was successful before trying to parse logs
	if receipt.Status == 0 {
		log.Error().Str("txHash", receipt.TxHash.String()).Msg("Bridge transaction failed")
		return 0, fmt.Errorf("bridge transaction failed with hash: %s", receipt.TxHash.String())
	}

	// Convert []*types.Log to []types.Log
	logs := make([]types.Log, len(receipt.Logs))
	for i, l := range receipt.Logs {
		logs[i] = *l
	}

	depositCount, err := ParseBridgeDepositCount(logs, bridgeContract)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse deposit count from logs")
		return 0, err
	}

	return depositCount, nil
}

// ParseBridgeDepositCount parses the deposit count from bridge transaction logs.
func ParseBridgeDepositCount(logs []types.Log, bridgeContract *ulxly.Ulxly) (uint32, error) {
	for _, l := range logs {
		// Try to parse the log as a BridgeEvent using the contract's filterer
		bridgeEvent, err := bridgeContract.ParseBridgeEvent(l)
		if err != nil {
			// This log is not a bridge event, continue to next log
			continue
		}

		// Successfully parsed a bridge event, return the deposit count
		return bridgeEvent.DepositCount, nil
	}

	return 0, fmt.Errorf("bridge event not found in logs")
}
