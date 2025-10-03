package fund

import (
	"errors"
	"math/big"

	_ "embed"

	"github.com/0xPolygon/polygon-cli/util"
	"github.com/spf13/cobra"
)

// The default private key used to send transactions.
const defaultPrivateKey = "0x42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa"

// The default mnemonic  used to derive wallets.
const defaultMnemonic = "code code code code code code code code code code code quality"

// The default password used to create a wallet for HD derivation.
const defaultPassword = "password"

// cmdFundParams holds the command-line parameters for the fund command.
type cmdFundParams struct {
	RpcUrl     string
	PrivateKey string

	WalletsNumber      uint64
	UseHDDerivation    bool
	WalletAddresses    []string
	FundingAmountInWei *big.Int
	OutputFile         string

	KeyFile string

	FunderAddress string
}

var (
	//go:embed usage.md
	usage               string
	params              cmdFundParams
	defaultFundingInWei = big.NewInt(50000000000000000) // 0.05 ETH
)

// FundCmd represents the fund command.
var FundCmd = &cobra.Command{
	Use:   "fund",
	Short: "Bulk fund crypto wallets automatically.",
	Long:  usage,
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		params.RpcUrl, err = util.GetRPCURL(cmd)
		if err != nil {
			return err
		}
		params.PrivateKey, err = util.GetPrivateKey(cmd)
		if err != nil {
			return err
		}
		return checkFlags()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return runFunding(cmd.Context())
	},
}

func init() {
	f := FundCmd.Flags()

	f.StringVarP(&params.RpcUrl, "rpc-url", "r", "http://localhost:8545", "RPC endpoint URL")
	f.StringVar(&params.PrivateKey, "private-key", defaultPrivateKey, "hex encoded private key to use for sending transactions")

	// Wallet parameters.
	f.Uint64VarP(&params.WalletsNumber, "number", "n", 10, "number of wallets to fund")
	f.BoolVar(&params.UseHDDerivation, "hd-derivation", true, "derive wallets to fund from private key in deterministic way")
	f.StringSliceVar(&params.WalletAddresses, "addresses", nil, "comma-separated list of wallet addresses to fund")
	params.FundingAmountInWei = defaultFundingInWei
	f.Var(&util.BigIntValue{Val: params.FundingAmountInWei}, "eth-amount", "amount of wei to send to each wallet")
	f.StringVar(&params.KeyFile, "key-file", "", "file containing accounts private keys, one per line")

	f.StringVarP(&params.OutputFile, "file", "f", "wallets.json", "output JSON file path for storing addresses and private keys of funded wallets")

	// Marking flags as mutually exclusive
	FundCmd.MarkFlagsMutuallyExclusive("addresses", "number")
	FundCmd.MarkFlagsMutuallyExclusive("addresses", "hd-derivation")
	FundCmd.MarkFlagsMutuallyExclusive("key-file", "addresses")
	FundCmd.MarkFlagsMutuallyExclusive("key-file", "number")
	FundCmd.MarkFlagsMutuallyExclusive("key-file", "hd-derivation")

	// Require at least one method to specify target accounts
	FundCmd.MarkFlagsOneRequired("addresses", "key-file", "number")

	// Funder contract parameters.
	f.StringVar(&params.FunderAddress, "contract-address", "", "address of pre-deployed Funder contract")
}

func checkFlags() error {
	// Validate funding amount
	minValue := big.NewInt(1000000000)
	if params.FundingAmountInWei != nil && params.FundingAmountInWei.Cmp(minValue) <= 0 {
		return errors.New("the funding amount must be greater than 1000000000")
	}
	return nil
}
