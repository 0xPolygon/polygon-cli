package ecrecover

import (
	_ "embed"
	"fmt"

	"github.com/maticnetwork/polygon-cli/util"
	"github.com/spf13/cobra"
)

var (
	//go:embed usage.md
	usage string

	rpcUrl      string
	blockNumber int
	extraData   string
)

var EcRecoverCmd = &cobra.Command{
	Use:   "ecrecover",
	Short: "Recovers and returns the public key of the signature",
	Long:  usage,
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return checkFlags()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return ecrecover(cmd.Context())
	},
}

func init() {
	EcRecoverCmd.PersistentFlags().StringVarP(&rpcUrl, "rpc-url", "r", "http://localhost:8545", "The RPC endpoint url")
	EcRecoverCmd.PersistentFlags().IntVarP(&blockNumber, "block-number", "b", 0, "Block number to check the extra data for")
	EcRecoverCmd.PersistentFlags().StringVarP(&extraData, "extra-data", "e", "", "Raw extra data")
}

func checkFlags() (err error) {
	if err = util.ValidateUrl(rpcUrl); err != nil {
		return
	}

	if blockNumber <= 0 {
		return fmt.Errorf("block-number should be greater than 0")
	}

	// if extraData == "" {
	// 	return fmt.Errorf("block-number should be greater than 0")
	// }

	return nil
}
