package tail

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "embed"

	"github.com/0xPolygon/polygon-cli/flag"
	"github.com/0xPolygon/polygon-cli/util"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type tailParams struct {
	RPCURL       string
	BlocksBack   uint64
	Follow       bool
	BatchSize    uint64
	PollInterval time.Duration
}

var (
	//go:embed usage.md
	usage     string
	inputTail = tailParams{}
)

var TailCmd = &cobra.Command{
	Use:   "tail",
	Short: "Tail full blocks from a JSON-RPC endpoint as NDJSON.",
	Long:  usage,
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, _ []string) (err error) {
		inputTail.RPCURL, err = flag.GetRPCURL(cmd)
		if err != nil {
			return err
		}
		if inputTail.BatchSize == 0 {
			return fmt.Errorf("batch-size must be greater than 0")
		}
		if inputTail.BlocksBack == 0 {
			return fmt.Errorf("blocks-back must be greater than 0")
		}
		if inputTail.PollInterval <= 0 {
			return fmt.Errorf("poll-interval must be greater than 0")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()
		ec, err := ethrpc.DialContext(ctx, inputTail.RPCURL)
		if err != nil {
			return err
		}

		latestBlock, err := getLatestBlockNumber(ctx, ec)
		if err != nil {
			return err
		}

		startBlock := uint64(0)
		if latestBlock+1 > inputTail.BlocksBack {
			startBlock = latestBlock - inputTail.BlocksBack + 1
		}

		nextBlock := startBlock
		log.Info().
			Uint64("latest", latestBlock).
			Uint64("start", startBlock).
			Bool("follow", inputTail.Follow).
			Msg("Starting tail")

		for {
			latestBlock, err = getLatestBlockNumber(ctx, ec)
			if err != nil {
				if !inputTail.Follow {
					return err
				}
				log.Warn().Err(err).Msg("Unable to fetch latest block number; retrying")
			} else if nextBlock <= latestBlock {
				if err := writeBlockRange(ctx, ec, nextBlock, latestBlock); err != nil {
					if !inputTail.Follow {
						return err
					}
					log.Warn().Err(err).Msg("Unable to fetch block range; retrying")
				} else {
					nextBlock = latestBlock + 1
				}
			}

			if !inputTail.Follow {
				return nil
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(inputTail.PollInterval):
			}
		}
	},
}

func init() {
	f := TailCmd.Flags()
	f.StringVarP(&inputTail.RPCURL, flag.RPCURL, "r", flag.DefaultRPCURL, "the RPC endpoint URL")
	f.Uint64VarP(&inputTail.BlocksBack, "blocks-back", "n", 10, "number of latest blocks to output before following")
	f.BoolVar(&inputTail.Follow, "follow", false, "poll for and stream newly produced blocks")
	f.Uint64VarP(&inputTail.BatchSize, "batch-size", "b", 150, "batch size for block requests")
	f.DurationVar(&inputTail.PollInterval, "poll-interval", 2*time.Second, "poll interval when --follow is enabled")
}

func getLatestBlockNumber(ctx context.Context, ec *ethrpc.Client) (uint64, error) {
	var result string
	if err := ec.CallContext(ctx, &result, "eth_blockNumber"); err != nil {
		return 0, err
	}
	blockNumber, err := strconv.ParseUint(result, 0, 64)
	if err != nil {
		return 0, fmt.Errorf("unable to parse latest block number %q: %w", result, err)
	}
	return blockNumber, nil
}

func writeBlockRange(ctx context.Context, ec *ethrpc.Client, start, end uint64) error {
	blocks, err := util.GetBlockRangeInPages(
		ctx,
		start,
		end,
		inputTail.BatchSize,
		ec,
		false,
	)
	if err != nil {
		return err
	}

	for _, block := range blocks {
		if _, err := fmt.Fprintln(os.Stdout, string(*block)); err != nil {
			return err
		}
	}

	return nil
}
