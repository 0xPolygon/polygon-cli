package publish

import (
	"context"
	_ "embed"
	"time"

	"github.com/0xPolygon/polygon-cli/cmd/flag_loader"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/time/rate"
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
	rateLimit     *uint64
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
	log.Info().
		Str("rpcURL", *publishInputArgs.rpcURL).
		Msg("Connected to the network")

	wp := NewWorkerPool(int(*publishInputArgs.concurrency))
	wp.Start(ctx)

	log.Info().
		Msg("Starting to publish transactions")

	limit := time.Duration(0)
	if *publishInputArgs.rateLimit > 0 {
		limit = time.Second / time.Duration(*publishInputArgs.rateLimit)
	}
	rl := rate.NewLimiter(rate.Every(limit), 1)

	summary := output{InputDataSource: dataSource}
	summary.Start()

	for inputDataItem := range inputData {
		err = rl.Wait(ctx)
		if err != nil {
			log.Error().
				Err(err).
				Str("inputDataItem", inputDataItem).
				Msg("failed to wait for rate limiter")
			continue
		}
		wp.SubmitJob(func(ctx context.Context, workerID int) {
			summary.InputDataCount.Add(1)

			var iErr error
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
			// log.Info().
			// 	Str("inputDataItem", inputDataItem).
			// 	Int("workerID", workerID).
			// 	Msg("decoded input data item to transaction successfully")

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
			// log.Info().
			// 	Str("inputDataItem", inputDataItem).
			// 	Int("workerID", workerID).
			// 	Msg("transaction sent successfully")

			summary.TxsSentSuccessfully.Add(1)
		})
	}

	wp.Stop()

	summary.Stop()

	log.Info().
		Msg("Finished publishing transactions")

	summary.Print()

	return err
}

func init() {
	publishInputArgs.rpcURL = Cmd.PersistentFlags().String(ArgRpcURL, defaultRPCURL, "The RPC URL of the network")
	publishInputArgs.concurrency = Cmd.PersistentFlags().Uint64P("concurrency", "c", 1, "Number of txs to send concurrently. Default is one request at a time.")
	publishInputArgs.inputFileName = Cmd.PersistentFlags().String("file", "", "Provide a filename with transactions to publish")
	publishInputArgs.rateLimit = Cmd.PersistentFlags().Uint64("rate-limit", 0, "Rate limit in txs per second. Default is no rate limit.")
}
