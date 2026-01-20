package events

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/0xPolygon/polygon-cli/bindings/ulxly/polygonrollupmanager"
	ulxlycommon "github.com/0xPolygon/polygon-cli/cmd/ulxly/common"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	ethclient "github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	getVerifyBatchesEvent   = &ulxlycommon.GetEvent{}
	getVerifyBatchesOptions = &ulxlycommon.GetVerifyBatchesOptions{}
)

//go:embed getVerifyBatchesUsage.md
var getVerifyBatchesUsage string

var GetVerifyBatchesCmd = &cobra.Command{
	Use:          "get-verify-batches",
	Short:        "Generate ndjson for each verify batch over a particular range of blocks.",
	Long:         getVerifyBatchesUsage,
	RunE:         readVerifyBatches,
	SilenceUsage: true,
}

func init() {
	getVerifyBatchesEvent.AddFlags(GetVerifyBatchesCmd)
	getVerifyBatchesOptions.AddFlags(GetVerifyBatchesCmd)
}

func readVerifyBatches(cmd *cobra.Command, _ []string) error {
	rollupManagerAddress := getVerifyBatchesOptions.RollupManagerAddress
	rpcURL := getVerifyBatchesEvent.URL
	toBlock := getVerifyBatchesEvent.ToBlock
	fromBlock := getVerifyBatchesEvent.FromBlock
	filter := getVerifyBatchesEvent.FilterSize

	var rpc *ethrpc.Client
	var err error

	if getVerifyBatchesEvent.Insecure {
		client, clientErr := ulxlycommon.CreateInsecureEthClient(rpcURL)
		if clientErr != nil {
			log.Error().Err(clientErr).Msg("Unable to create insecure client")
			return clientErr
		}
		defer client.Close()
		rpc = client.Client()
	} else {
		rpc, err = ethrpc.DialContext(cmd.Context(), rpcURL)
		if err != nil {
			log.Error().Err(err).Msg("Unable to Dial RPC")
			return err
		}
		defer rpc.Close()
	}

	client := ethclient.NewClient(rpc)
	rm := common.HexToAddress(rollupManagerAddress)
	rollupManager, err := polygonrollupmanager.NewPolygonrollupmanager(rm, client)
	if err != nil {
		return err
	}
	verifyBatchesTrustedAggregatorSignatureHash := crypto.Keccak256Hash([]byte("VerifyBatchesTrustedAggregator(uint32,uint64,bytes32,bytes32,address)"))

	currentBlock := fromBlock
	for currentBlock < toBlock {
		endBlock := min(currentBlock+filter, toBlock)
		// Filter 0xd1ec3a1216f08b6eff72e169ceb548b782db18a6614852618d86bb19f3f9b0d3
		query := ethereum.FilterQuery{
			FromBlock: new(big.Int).SetUint64(currentBlock),
			ToBlock:   new(big.Int).SetUint64(endBlock),
			Addresses: []common.Address{rm},
			Topics:    [][]common.Hash{{verifyBatchesTrustedAggregatorSignatureHash}},
		}
		logs, err := client.FilterLogs(cmd.Context(), query)
		if err != nil {
			return err
		}

		for _, vLog := range logs {
			vb, err := rollupManager.ParseVerifyBatchesTrustedAggregator(vLog)
			if err != nil {
				return err
			}
			log.Info().Uint32("RollupID", vb.RollupID).Uint64("block-number", vb.Raw.BlockNumber).Msg("Found rollupmanager VerifyBatchesTrustedAggregator event")
			var jBytes []byte
			jBytes, err = json.Marshal(vb)
			if err != nil {
				return err
			}
			fmt.Println(string(jBytes))
		}
		currentBlock = endBlock + 1
	}

	return nil
}
