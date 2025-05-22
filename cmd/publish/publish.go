package publish

import (
	"context"
	_ "embed"

	"github.com/0xPolygon/polygon-cli/cmd/flag_loader"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	ArgRpcURL = "rpc-url"
	ArgForkID = "fork-id"

	defaultRPCURL = "http://localhost:8545"
)

//go:embed publish.md
var cmdUsage string

var Cmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish transactions to the network with high-throughput",
	Long:  cmdUsage,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		publishInputArgs.rpcURL = flag_loader.GetRpcUrlFlagValue(cmd)
	},
	RunE: publish,
}

type inputArgs struct {
	rpcURL        *string
	concurrency   *uint64
	inputFileName *string
}

var publishInputArgs = inputArgs{}

func publish(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	inputData, dataSource, err := getInputData(publishInputArgs.inputFileName, args)
	if err != nil {
		cmd.PrintErrf("There was an error reading input data: %s", err.Error())
		return err
	}

	client, err := ethclient.Dial(*publishInputArgs.rpcURL)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("rpcURL", *publishInputArgs.rpcURL).
			Msg("Failed to connect to the network")
	}

	wp := NewWorkerPool(int(*publishInputArgs.concurrency))
	wp.Start(ctx)

	summary := output{InputDataSource: dataSource}
	summary.Start()
	for inputDataItem := range inputData {
		wp.SubmitJob(func(ctx context.Context, workerID int) {
			var iErr error

			summary.InputDataCount.Add(1)

			tx, iErr := inputDataItemToTx(inputDataItem)
			if iErr != nil {
				log.Error().
					Err(iErr).
					Str("inputDataItem", inputDataItem).
					Int("workerID", workerID).
					Msg("failed to decode input data item to transaction")
				summary.InvalidInputs.Add(1)
				return
			}
			summary.ValidInputs.Add(1)

			iErr = client.SendTransaction(ctx, tx)
			if iErr != nil {
				log.Error().
					Err(iErr).
					Str("inputDataItem", inputDataItem).
					Int("workerID", workerID).
					Msg("failed to send transaction")
				summary.TxsSentUnsuccessfully.Add(1)
				return
			}
			log.Info().
				Str("inputDataItem", inputDataItem).
				Int("workerID", workerID).
				Msg("transaction sent successfully")

			summary.TxsSentSuccessfully.Add(1)
		})
	}

	wp.Stop()

	summary.Stop()
	summary.Print()

	return err
}

func init() {
	publishInputArgs.rpcURL = Cmd.PersistentFlags().String(ArgRpcURL, defaultRPCURL, "The RPC URL of the network")
	publishInputArgs.concurrency = Cmd.PersistentFlags().Uint64P("concurrency", "c", 1, "Number of txs to send concurrently. Default is one request at a time.")
	publishInputArgs.inputFileName = Cmd.PersistentFlags().String("file", "", "Provide a filename with transactions to publish")
}
