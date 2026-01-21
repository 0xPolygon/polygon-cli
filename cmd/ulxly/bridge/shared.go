package bridge

import (
	"context"
	"fmt"

	"github.com/0xPolygon/polygon-cli/bindings/ulxly"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
)

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
