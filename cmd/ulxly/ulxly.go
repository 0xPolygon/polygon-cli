package ulxly

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/maticnetwork/polygon-cli/bindings/ulxly"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type uLxLyArgs struct {
	FromBlock     *int64
	ToBlock       *int64
	RPCURL        *string
	BridgeAddress *string
}

var ulxlyInputArgs uLxLyArgs

var ULxLyCmd = &cobra.Command{
	Use:   "ulxly",
	Short: "Utilities for interacting with the lxly bridge",
	Long:  "TODO",
	Args:  cobra.NoArgs,
}

var DepositsCmd = &cobra.Command{
	Use:     "deposits",
	Short:   "get a range of deposits",
	Long:    "TODO",
	Args:    cobra.NoArgs,
	PreRunE: checkDepositArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		// Dial the Ethereum RPC server.
		rpc, err := ethrpc.DialContext(ctx, *ulxlyInputArgs.RPCURL)
		if err != nil {
			log.Error().Err(err).Msg("Unable to dial rpc")
			return err
		}
		defer rpc.Close()
		ec := ethclient.NewClient(rpc)
		bridge, err := ulxly.NewUlxly(common.HexToAddress(*ulxlyInputArgs.BridgeAddress), ec)
		if err != nil {
			return err
		}
		fromBlock := uint64(*ulxlyInputArgs.FromBlock)
		toBlock := uint64(*ulxlyInputArgs.ToBlock)
		opts := bind.FilterOpts{
			Start:   fromBlock,
			End:     &toBlock,
			Context: ctx,
		}
		evtIterator, err := bridge.FilterBridgeEvent(&opts)
		if err != nil {
			return err
		}
		defer evtIterator.Close()
		for evtIterator.Next() {
			evt := evtIterator.Event
			log.Info().Uint32("deposit", evt.DepositCount).Msg("Found Deposit")
			jBytes, err := json.Marshal(evt)
			if err != nil {
				return err
			}
			fmt.Println(string(jBytes))
		}
		return nil
	},
}

func checkDepositArgs(cmd *cobra.Command, args []string) error {
	if *ulxlyInputArgs.BridgeAddress == "" {
		return fmt.Errorf("please provide the bridge address")
	}
	return nil
}

func init() {
	ULxLyCmd.AddCommand(DepositsCmd)
	//   - When blockNr is -1 the chain pending header is returned.
	//   - When blockNr is -2 the chain latest header is returned.
	//   - When blockNr is -3 the chain finalized header is returned.
	//   - When blockNr is -4 the chain safe header is returned.
	ulxlyInputArgs.FromBlock = DepositsCmd.PersistentFlags().Int64("from-block", 0, "The block height to start query at.")
	ulxlyInputArgs.ToBlock = DepositsCmd.PersistentFlags().Int64("to-block", -2, "The block height to start query at.")
	ulxlyInputArgs.RPCURL = DepositsCmd.PersistentFlags().String("rpc-url", "http://127.0.0.1:8545", "The RPC to query for events")

	ulxlyInputArgs.BridgeAddress = DepositsCmd.Flags().String("bridge-address", "", "The address of the lxly bridge")

}
