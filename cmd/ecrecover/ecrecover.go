package ecrecover

import (
	_ "embed"
	"encoding/json"
	"io"
	"math/big"
	"os"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/maticnetwork/polygon-cli/util"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	//go:embed usage.md
	usage string

	rpcUrl      string
	blockNumber uint64
	filePath    string
)

var EcRecoverCmd = &cobra.Command{
	Use:   "ecrecover",
	Short: "Recovers and returns the public key of the signature",
	Long:  usage,
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return checkFlags()
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		var block *types.Block
		var err error

		if rpcUrl == "" {
			var blockJSON []byte
			if filePath != "" {
				blockJSON, err = os.ReadFile(filePath)
				if err != nil {
					log.Error().Err(err).Msg("Unable to read file")
					return
				}
			} else {
				blockJSON, err = io.ReadAll(os.Stdin)
				if err != nil {
					log.Error().Err(err).Msg("Unable to read stdin")
					return
				}
			}

			var header types.Header
			if err = json.Unmarshal(blockJSON, &header); err != nil {
				log.Error().Err(err).Msg("Unable to unmarshal JSON")
				return
			}

			block = types.NewBlockWithHeader(&header)
			blockNumber = header.Number.Uint64()
		} else {
			var rpc *ethrpc.Client
			rpc, err = ethrpc.DialContext(ctx, rpcUrl)
			if err != nil {
				log.Error().Err(err).Msg("Unable to dial rpc")
				return
			}

			ec := ethclient.NewClient(rpc)

			if blockNumber == 0 {
				blockNumber, err = ec.BlockNumber(ctx)
				if err != nil {
					log.Error().Err(err).Msg("Unable to retrieve latest block number")
					return
				}
			}

			block, err = ec.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
			if err != nil {
				log.Error().Err(err).Msg("Unable to retrieve block")
				return
			}
		}

		signerBytes, err := util.Ecrecover(block)
		if err != nil {
			log.Error().Err(err).Msg("Unable to recover signature")
			return
		}
		cmd.Printf("Recovering signer from block #%d\n", blockNumber)
		cmd.Println(ethcommon.BytesToAddress(signerBytes))
	},
}

func init() {
	EcRecoverCmd.PersistentFlags().StringVarP(&rpcUrl, "rpc-url", "r", "", "The RPC endpoint url")
	EcRecoverCmd.PersistentFlags().Uint64VarP(&blockNumber, "block-number", "b", 0, "Block number to check the extra data for (default: latest)")
	EcRecoverCmd.PersistentFlags().StringVarP(&filePath, "file", "f", "", "Path to a file containing block information in JSON format")
}

func checkFlags() (err error) {
	if rpcUrl != "" && util.ValidateUrl(rpcUrl) != nil {
		return
	}

	return nil
}
