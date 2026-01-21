// Package deposits provides the get-deposits command.
package deposits

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/0xPolygon/polygon-cli/bindings/ulxly"
	ulxlycommon "github.com/0xPolygon/polygon-cli/cmd/ulxly/common"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethclient "github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	getDepositEvent      = &ulxlycommon.GetEvent{}
	getDepositSmcOptions = &ulxlycommon.GetSmcOptions{}
)

//go:embed usage.md
var usage string

var Cmd = &cobra.Command{
	Use:          "get-deposits",
	Short:        "Generate ndjson for each bridge deposit over a particular range of blocks.",
	Long:         usage,
	RunE:         readDeposit,
	SilenceUsage: true,
}

func init() {
	getDepositEvent.AddFlags(Cmd)
	getDepositSmcOptions.AddFlags(Cmd)
}

func readDeposit(cmd *cobra.Command, _ []string) error {
	bridgeAddress := getDepositSmcOptions.BridgeAddress
	rpcURL := getDepositEvent.URL
	toBlock := getDepositEvent.ToBlock
	fromBlock := getDepositEvent.FromBlock
	filter := getDepositEvent.FilterSize

	var rpc *ethrpc.Client
	var err error

	if getDepositEvent.Insecure {
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

	ec := ethclient.NewClient(rpc)

	bridgeV2, err := ulxly.NewUlxly(common.HexToAddress(bridgeAddress), ec)
	if err != nil {
		return err
	}
	currentBlock := fromBlock
	for currentBlock < toBlock {
		endBlock := min(currentBlock+filter, toBlock)

		opts := bind.FilterOpts{
			Start:   currentBlock,
			End:     &endBlock,
			Context: cmd.Context(),
		}
		evtV2Iterator, err := bridgeV2.FilterBridgeEvent(&opts)
		if err != nil {
			return err
		}

		for evtV2Iterator.Next() {
			evt := evtV2Iterator.Event
			log.Info().Uint32("deposit", evt.DepositCount).Uint64("block-number", evt.Raw.BlockNumber).Msg("Found ulxly Deposit")
			var jBytes []byte
			jBytes, err = json.Marshal(evt)
			if err != nil {
				return err
			}
			fmt.Println(string(jBytes))
		}
		err = evtV2Iterator.Close()
		if err != nil {
			log.Error().Err(err).Msg("error closing event iterator")
		}
		currentBlock = endBlock + 1
	}

	return nil
}
