package ecrecover

import (
	_ "embed"
	"encoding/json"
	"io"
	"math/big"
	"os"

	"github.com/0xPolygon/polygon-cli/cmd/flag_loader"
	"github.com/0xPolygon/polygon-cli/util"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	//go:embed usage.md
	usage string

	rpcUrl      string
	blockNumber uint64
	filePath    string
	txData      string
)

var EcRecoverCmd = &cobra.Command{
	Use:   "ecrecover",
	Short: "Recovers and returns the public key of the signature",
	Long:  usage,
	Args:  cobra.NoArgs,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		rpcUrlFlagValue := flag_loader.GetRpcUrlFlagValue(cmd)
		rpcUrl = *rpcUrlFlagValue
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return checkFlags()
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		var (
			signerBytes []byte
			err         error
		)

		if filePath != "" { // block signer from file
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

			block := types.NewBlockWithHeader(&header)
			blockNumber = header.Number.Uint64()
			signerBytes, err = util.Ecrecover(block)

		} else if txData != "" { // transaction signer from data

			txBytes := ethcommon.FromHex(txData)
			var tx types.Transaction
			err = tx.UnmarshalBinary(txBytes)
			if err != nil {
				log.Error().Err(err).Msg("Unable to decode transaction")
				return
			}
			signerBytes, err = util.EcrecoverTx(&tx)
			if err != nil {
				log.Error().Err(err).Msg("Unable to retrieve block")
				return
			}

		} else { // block signer block-number, requires rcp-url
			if rpcUrl == "" {
				log.Error().Msg("No RPC URL provided")
				return
			}
			var rpc *ethrpc.Client
			rpc, err = ethrpc.DialContext(ctx, rpcUrl)
			if err != nil {
				log.Error().Err(err).Msg("Unable to dial rpc")
				return
			}
			ec := ethclient.NewClient(rpc)
			defer ec.Close()

			var block *types.Block
			if blockNumber == 0 {
				blockNumber, err = ec.BlockNumber(ctx)
				if err != nil {
					log.Error().Err(err).Msg("Unable to retrieve latest block number")
					return
				}
				cmd.Println("Using latest block number:", blockNumber)
			}
			block, err = ec.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
			if err != nil {
				log.Error().Err(err).Msg("Unable to retrieve block")
				return
			}
			signerBytes, err = util.Ecrecover(block)
		}

		if err != nil {
			log.Error().Err(err).Msg("Unable to recover signature")
			return
		}
		cmd.Println(ethcommon.BytesToAddress(signerBytes))
	},
}

func init() {
	EcRecoverCmd.PersistentFlags().StringVarP(&rpcUrl, "rpc-url", "r", "", "The RPC endpoint url")
	EcRecoverCmd.PersistentFlags().Uint64VarP(&blockNumber, "block-number", "b", 0, "Block number to check the extra data for (default: latest)")
	EcRecoverCmd.PersistentFlags().StringVarP(&filePath, "file", "f", "", "Path to a file containing block information in JSON format")
	EcRecoverCmd.PersistentFlags().StringVarP(&txData, "tx", "t", "", "Transaction data in hex format")

	// The sources of decoding are mutually exclusive
	EcRecoverCmd.MarkFlagsMutuallyExclusive("file", "block-number", "tx")
}

func checkFlags() error {
	var err error
	if rpcUrl != "" {
		if err = util.ValidateUrl(rpcUrl); err != nil {
			return err
		}
	}
	return err
}
