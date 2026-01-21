// Package claims provides the get-claims command.
package claims

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
	getClaimEvent      = &ulxlycommon.GetEvent{}
	getClaimSmcOptions = &ulxlycommon.GetSmcOptions{}
)

//go:embed usage.md
var usage string

var Cmd = &cobra.Command{
	Use:          "get-claims",
	Short:        "Generate ndjson for each bridge claim over a particular range of blocks.",
	Long:         usage,
	RunE:         readClaim,
	SilenceUsage: true,
}

func init() {
	getClaimEvent.AddFlags(Cmd)
	getClaimSmcOptions.AddFlags(Cmd)
}

func readClaim(cmd *cobra.Command, _ []string) error {
	bridgeAddress := getClaimSmcOptions.BridgeAddress
	rpcURL := getClaimEvent.URL
	toBlock := getClaimEvent.ToBlock
	fromBlock := getClaimEvent.FromBlock
	filter := getClaimEvent.FilterSize

	var rpc *ethrpc.Client
	var err error

	if getClaimEvent.Insecure {
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
		evtV2Iterator, err := bridgeV2.FilterClaimEvent(&opts)
		if err != nil {
			return err
		}

		for evtV2Iterator.Next() {
			evt := evtV2Iterator.Event
			var (
				mainnetFlag                     bool
				rollupIndex, localExitRootIndex uint32
			)
			mainnetFlag, rollupIndex, localExitRootIndex, err = ulxlycommon.DecodeGlobalIndex(evt.GlobalIndex)
			if err != nil {
				log.Error().Err(err).Msg("error decoding globalIndex")
				return err
			}
			log.Info().Bool("claim-mainnetFlag", mainnetFlag).Uint32("claim-RollupIndex", rollupIndex).Uint32("claim-LocalExitRootIndex", localExitRootIndex).Uint64("block-number", evt.Raw.BlockNumber).Msg("Found Claim")
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
