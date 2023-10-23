package fork

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/clique"
	"github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	rpcURL                string
	blockHash             ethcommon.Hash
	retryLimit            = 30
	errRetryLimitExceeded = fmt.Errorf("unable to process request after hitting retry limit")
)

var ForkCmd = &cobra.Command{
	Use:   "fork blockhash url",
	Short: "Take a forked block and walk up the chain to do analysis.",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info().Str("rpc", rpcURL).Str("blockHash", blockHash.String()).Msg("Starting Analysis")
		c, err := ethclient.Dial(rpcURL)
		if err != nil {
			log.Error().Err(err).Str("rpc", rpcURL).Msg("Could not rpc dial connection")
			return err
		}
		return walkTheBlocks(blockHash, c)
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
		log.Error().Err(err).Msg("Unable to get current block number from chain")
		return err
	}
	log.Info().Uint64("headBlock", bn).Msg("Retrieved current head of the chain")

	folderName := fmt.Sprintf("fork-analysis-%d", time.Now().Unix())
	if err := os.Mkdir(folderName, os.ModePerm); err != nil {
		log.Error().Err(err).Msg("Unable to create output folder")
		return err
	}

	for {
		potentialForkedBlock, err := getBlockByHash(ctx, inputBlockHash, client)
		if err != nil {
			log.Error().Err(err).Str("blockhash", inputBlockHash.String()).Msg("Unable to fetch block")
			return err
		}
		log.Info().Uint64("number", potentialForkedBlock.NumberU64()).Msg("Successfully retrieved starting block hash")

		canonicalBlock, err := client.BlockByNumber(ctx, potentialForkedBlock.Number())
		if err != nil {
			log.Error().Err(err).Uint64("number", potentialForkedBlock.NumberU64()).Msg("Unable to retrieve block by number")

			return err
		}
		if potentialForkedBlock.Hash().String() == canonicalBlock.Hash().String() {
			err = writeBlock(folderName, canonicalBlock, true)
			if err != nil {
				log.Error().Err(err).Msg("Failed to save final canonical block")
			}
			log.Info().Uint64("number", canonicalBlock.NumberU64()).Str("blockHash", canonicalBlock.Hash().String()).Msg("The current block seems to be canonical in the chain. Stopping analysis")
			break
		}
		log.Info().
			Uint64("number", potentialForkedBlock.NumberU64()).
			Str("forkedBlockHash", potentialForkedBlock.Hash().String()).
			Str("canonicalBlockHash", canonicalBlock.Hash().String()).
			Msg("Identified forked block. Continuing traversal")

		err = writeBlock(folderName, potentialForkedBlock, false)
		if err != nil {
			log.Error().Err(err).Msg("Unable to save forked block")
			return err
		}
		err = writeBlock(folderName, canonicalBlock, true)
		if err != nil {
			log.Error().Err(err).Msg("Unable to save canonical block")
			return err
		}
		// Ever higher
		inputBlockHash = potentialForkedBlock.ParentHash()
	}
	return nil
}

func writeBlock(folderName string, block *types.Block, isCanonical bool) error {
	rawHeader, err := block.Header().MarshalJSON()
	if err != nil {
		log.Error().Err(err).Msg("Unable to json marshal the header")
		return err
	}
	fields := make(map[string]interface{}, 0)
	err = json.Unmarshal(rawHeader, &fields)
	if err != nil {
		log.Error().Err(err).Msg("Unable to convert header to map type")
		return err
	}
	fields["transactions"] = block.Transactions()
	// TODO in the future if this is used in other chains or with different types of consensus this would need to be revised
	signer, err := ecrecover(block)
	if err != nil {
		log.Error().Err(err).Msg("Unable to recover signature")
		return err
	}
	fields["_signer"] = ethcommon.BytesToAddress(signer)

	jsonData, err := json.Marshal(fields)
	if err != nil {
		log.Error().Err(err).Msg("Unable to marshal block to json")
		return err
	}
	blockType := "c"
	if !isCanonical {
		blockType = "f"
	}
	fileName := fmt.Sprintf("%s/%d-%s-%s.json", folderName, block.NumberU64(), blockType, block.Hash().String())
	err = os.WriteFile(fileName, jsonData, 0744)
	return err
}

// getBlockByHash will try to get a block by hash in a loop. Unless we have a dedicated node that we know has the forked blocks it's going to be tricky to consistently get results from a fork. So it requires some brute force
func getBlockByHash(ctx context.Context, bh ethcommon.Hash, client *ethclient.Client) (*types.Block, error) {
	for i := 0; i < retryLimit; i = i + 1 {
		block, err := client.BlockByHash(ctx, bh)
		if err != nil {
			log.Warn().Err(err).Int("attempt", i).Str("blockhash", bh.String()).Msg("Unable to fetch block")
		} else {
			return block, nil
		}
		time.Sleep(2 * time.Second)
	}
	log.Error().Err(errRetryLimitExceeded).Str("blockhash", bh.String()).Int("retryLimit", retryLimit).Msg("Unable to fetch block after retrying")
	return nil, errRetryLimitExceeded
}

func ecrecover(block *types.Block) ([]byte, error) {
	header := block.Header()
	sigStart := len(header.Extra) - ethcrypto.SignatureLength
	if sigStart < 0 || sigStart > len(header.Extra) {
		return nil, fmt.Errorf("unable to recover signature")
	}
	signature := header.Extra[sigStart:]
	pubkey, err := ethcrypto.Ecrecover(clique.SealHash(header).Bytes(), signature)
	if err != nil {
		return nil, err
	}
	signer := ethcrypto.Keccak256(pubkey[1:])[12:]

	return signer, nil
}
