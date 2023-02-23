package fork

import (
	"context"
	"fmt"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	rpcURL     string
	blockHash  ethcommon.Hash
	retryLimit = 30

	errRetryLimitExceeded = fmt.Errorf("Unable to process request after hitting retry limit")
)
var ForkCmd = &cobra.Command{
	Use:   "fork blockhash http://polygon-rpc.com",
	Short: "Take a forked block and walk up the chain to do analysis",
	Long: `
TODO
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info().Msg("Hi there")
		log.Info().Str("rpc", rpcURL).Str("blockHash", blockHash.String()).Msg("Starting Analysis")
		c, err := ethclient.Dial(rpcURL)
		if err != nil {
			log.Error().Err(err).Str("rpc", rpcURL).Msg("Could not rpc dial connection")
			return err
		}
		walkTheBlocks(blockHash, c)
		return nil
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("two arguments required a block hash and an RPC URL")
		}
		blockHash = ethcommon.HexToHash(args[0])
		rpcURL = args[1]
		return nil
	},
}

func walkTheBlocks(inputBlockHash ethcommon.Hash, client *ethclient.Client) error {
	log.Info().Msg("Starting block analysis")
	ctx := context.Background()
	bn, err := client.BlockNumber(ctx)
	if err != nil {
		log.Error().Err(err).Msg("unable to get current block number from chain")
		return err
	}
	log.Info().Uint64("headBlock", bn).Msg("retrieved current head of the chain")

	for {
		potentialForkedBlock, err := getBlockByHash(ctx, inputBlockHash, client)
		if err != nil {
			log.Error().Err(err).Str("blockhash", inputBlockHash.String()).Msg("unable to fetch block")
			return err
		}
		log.Info().Uint64("number", potentialForkedBlock.NumberU64()).Msg("successfully retrieved starting block hash")

		canonicalBlock, err := client.BlockByNumber(ctx, potentialForkedBlock.Number())
		if err != nil {
			log.Error().Err(err).Uint64("number", potentialForkedBlock.NumberU64()).Msg("unable to retrieve block by number")
			return err
		}
		if potentialForkedBlock.Hash().String() == canonicalBlock.Hash().String() {
			log.Info().Uint64("number", canonicalBlock.NumberU64()).Str("blockHash", canonicalBlock.Hash().String()).Msg("the current block seems to be canonical in the chain. Stopping analysis")
			break
		} else {
			log.Info().
				Uint64("number", potentialForkedBlock.NumberU64()).
				Str("forkedBlockHash", potentialForkedBlock.Hash().String()).
				Str("canonicalBlockHash", canonicalBlock.Hash().String()).
				Msg("Identified forked block. Continuing traversal")
		}
		// Ever higher
		inputBlockHash = potentialForkedBlock.ParentHash()
	}
	return nil
}

// getBlockByHash will try to get a block by hash in a loop. Unless we have a dedicated node that we know has the forked blocks it's going to be tricky to consistently get results from a fork. So it requires some brute force
func getBlockByHash(ctx context.Context, bh ethcommon.Hash, client *ethclient.Client) (*types.Block, error) {
	for i := 0; i < retryLimit; i = i + 1 {
		block, err := client.BlockByHash(ctx, bh)
		if err != nil {
			log.Warn().Err(err).Int("attempt", i).Str("blockhash", bh.String()).Msg("unable to fetch block")
		} else {
			return block, nil
		}
		time.Sleep(2 * time.Second)
	}
	log.Error().Err(errRetryLimitExceeded).Str("blockhash", bh.String()).Int("retryLimit", retryLimit).Msg("unable to fetch block after retrying")
	return nil, errRetryLimitExceeded
}

func init() {
	// flagSet := ForkCmd.PersistentFlags()
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

}
