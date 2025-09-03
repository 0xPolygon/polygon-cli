package contract

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/0xPolygon/polygon-cli/cmd/flag_loader"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	ArgRpcURL     = "rpc-url"
	ArgAddress    = "address"
	defaultRPCURL = "http://localhost:8545"
)

var (
	//go:embed usage.md
	usage string
)

type contractInputArgs struct {
	rpcURL  *string
	address *string
}

var inputArgs = contractInputArgs{}

type ContractInfo struct {
	Address         string             `json:"address"`
	Balance         uint64             `json:"balance,omitempty"`
	CreationTx      *types.Transaction `json:"creation_tx,omitempty"`
	CreationReceipt *types.Receipt     `json:"creation_receipt,omitempty"`
}

var Cmd = &cobra.Command{
	Use:   "contract",
	Short: "Interact with smart contracts and fetch contract information from the blockchain",
	Long:  usage,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		inputArgs.rpcURL = flag_loader.GetRpcUrlFlagValue(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return contract(cmd)
	},
}

func contract(cmd *cobra.Command) error {
	ctx := cmd.Context()

	rpcURL := *inputArgs.rpcURL
	address := *inputArgs.address

	// connect to the blockchain node and fetch contract information
	c, err := ethclient.Dial(rpcURL)
	if err != nil {
		return fmt.Errorf("failed to connect to rpc url %s: %w", rpcURL, err)
	}
	defer c.Close()

	tx, receipt, err := fetchContractCreationTx(ctx, c, address)
	if err != nil {
		return fmt.Errorf("failed to fetch contract creation tx: %w", err)
	}

	balance, err := c.BalanceAt(ctx, common.HexToAddress(address), nil)
	if err != nil {
		return fmt.Errorf("failed to fetch contract balance: %w", err)
	}

	info := ContractInfo{
		Address:         address,
		CreationReceipt: receipt,
		CreationTx:      tx,
		Balance:         balance.Uint64(),
	}

	output, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal contract info: %w", err)
	}
	fmt.Println(string(output))
	return nil

}

func fetchContractCreationTx(ctx context.Context, client *ethclient.Client, contractAddress string) (*types.Transaction, *types.Receipt, error) {
	log.Info().Msg("Checking if contract exists")
	// check if contract exists
	code, err := client.CodeAt(ctx, common.HexToAddress(contractAddress), nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch contract code: %w", err)
	}
	if len(code) == 0 {
		return nil, nil, fmt.Errorf("no contract found at address %s", contractAddress)
	}
	log.Info().Msg("Contract found!")

	latestBlock, err := client.BlockNumber(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch latest block: %w", err)
	}
	start := uint64(0)
	end := latestBlock

	var creationTx *types.Transaction
	var creationReceipt *types.Receipt

	log.Info().Msg("Finding contract creation block")
	for start <= end {
		mid := (start + end) / 2

		code, iErr := client.CodeAt(ctx, common.HexToAddress(contractAddress), new(big.Int).SetUint64(mid))
		if iErr != nil {
			return nil, nil, fmt.Errorf("failed to fetch contract code: %w", iErr)
		}

		if len(code) == 0 {
			start = mid + 1
		} else {
			end = mid - 1
		}
	}
	log.Info().Uint64("blockNumber", start).Msg("Block found!")

	// Fetch the block at which the contract was created
	block, err := client.BlockByNumber(ctx, new(big.Int).SetUint64(start))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch block %d: %w", start, err)
	}

	log.Info().Msg("Searching for contract creation transaction in the block")
	// Iterate through the transactions in the block to find the contract creation transaction
	for _, tx := range block.Transactions() {
		if tx.To() != nil {
			continue // Not a contract creation transaction
		}
		if len(tx.Data()) == 0 {
			continue // No data, not a contract creation transaction
		}

		receipt, err := client.TransactionReceipt(ctx, tx.Hash())
		if err != nil {
			return nil, nil, fmt.Errorf("failed to fetch transaction receipt: %w", err)
		}

		if receipt != nil && strings.EqualFold(receipt.ContractAddress.Hex(), contractAddress) {
			creationTx = tx
			creationReceipt = receipt
			break
		}
	}

	if creationTx == nil || creationReceipt == nil {
		return nil, nil, fmt.Errorf("contract creation transaction not found")
	}
	log.Info().Str("hash", creationTx.Hash().String()).Msg("Contract creation transaction found!")

	return creationTx, creationReceipt, nil
}

func init() {
	inputArgs.rpcURL = Cmd.PersistentFlags().String(ArgRpcURL, defaultRPCURL, "The RPC URL of the network containing the contract")
	inputArgs.address = Cmd.PersistentFlags().String(ArgAddress, "", "The contract address")

	_ = Cmd.MarkPersistentFlagRequired(ArgAddress)
}
