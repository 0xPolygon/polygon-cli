package ecrecover

import (
	_ "embed"
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
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

		rpc, err := ethrpc.DialContext(ctx, rpcUrl)
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

		block, err := ec.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
		if err != nil {
			log.Error().Err(err).Msg("Unable to retrieve block")
			return
		}

		if len(block.Transactions()) == 0 {
			log.Error().Msg("no transaction to derive public key fromk")
			return
		}

		signerBytes, err := util.Ecrecover(block)
		if err != nil {
			log.Error().Err(err).Msg("Unable to recover signature")
			return
		}
		cmd.Println(ethcommon.BytesToAddress(signerBytes))
	},
}

func init() {
	EcRecoverCmd.PersistentFlags().StringVarP(&rpcUrl, "rpc-url", "r", "http://localhost:8545", "The RPC endpoint url")
	EcRecoverCmd.PersistentFlags().Uint64VarP(&blockNumber, "block-number", "b", 0, "Block number to check the extra data for")
}

func checkFlags() (err error) {
	if err = util.ValidateUrl(rpcUrl); err != nil {
		return
	}

	return nil
}
