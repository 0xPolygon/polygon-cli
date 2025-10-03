package publish

import (
	"context"
	_ "embed"
	"time"

	"github.com/0xPolygon/polygon-cli/flag"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/time/rate"
)

const (
	ArgForkID = "fork-id"
)

//go:embed publish.md
var cmdUsage string

var Cmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish transactions to the network with high-throughput",
	Long:  cmdUsage,
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		publishInputArgs.rpcURL, err = flag.GetRPCURL(cmd)
		if err != nil {
			return err
		}
		return nil
	},
	RunE: publish,
}

type inputArgs struct {
	rpcURL        string
	concurrency   uint64
	jobQueueSize  uint64
	inputFileName string
	rateLimit     uint64
}

var publishInputArgs = inputArgs{}

func publish(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	inputData, dataSource, err := getInputData(&publishInputArgs.inputFileName, args)
	if err != nil {
		cmd.PrintErrf("There was an error reading input data: %s", err.Error())
		return err
	}

	client, err := ethclient.Dial(publishInputArgs.rpcURL)
	if err != nil {
		return err
	}
	defer client.Close()
	log.Info().
		Str("rpcURL", publishInputArgs.rpcURL).
		Msg("Connected to the network")

	wp := NewWorkerPool(publishInputArgs.concurrency, publishInputArgs.jobQueueSize)
	wp.Start(ctx)

	log.Info().
		Msg("Starting to publish transactions")

	limit := time.Duration(0)
	if publishInputArgs.rateLimit > 0 {
		limit = time.Second / time.Duration(publishInputArgs.rateLimit)
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
		wp.SubmitJob(func(ctx context.Context, workerID uint64) {
			summary.InputDataCount.Add(1)

			var iErr error
			tx, iErr := inputDataItemToTx(inputDataItem)
			if iErr != nil {
				log.Error().
					Err(iErr).
					Str("inputDataItem", inputDataItem).
					Uint64("workerID", workerID).
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
					Uint64("workerID", workerID).
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

	return nil
}

func init() {
	f := Cmd.Flags()
	f.StringVar(&publishInputArgs.rpcURL, flag.RPCURL, flag.DefaultRPCURL, "RPC URL of network")
	f.Uint64VarP(&publishInputArgs.concurrency, "concurrency", "c", 1, "number of txs to send concurrently (default: one at a time)")
	f.Uint64Var(&publishInputArgs.jobQueueSize, "job-queue-size", 100, "number of jobs we can put in the job queue for workers to process")
	f.StringVar(&publishInputArgs.inputFileName, "file", "", "provide a filename with transactions to publish")
	f.Uint64Var(&publishInputArgs.rateLimit, "rate-limit", 0, "rate limit in txs per second (default: no limit)")
}
