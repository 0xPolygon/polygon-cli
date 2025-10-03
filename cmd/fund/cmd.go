package fund

import (
	"errors"
	"math/big"

	_ "embed"

	"github.com/0xPolygon/polygon-cli/cmd/flag_loader"
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
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		rpcUrl := flag_loader.GetRpcUrlFlagValue(cmd)
		params.RpcUrl = *rpcUrl
		privateKey := flag_loader.GetPrivateKeyFlagValue(cmd)
		params.PrivateKey = *privateKey
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return checkFlags()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return runFunding(cmd.Context())
	},
}

func init() {
	p := new(cmdFundParams)
	f := FundCmd.Flags()

	f.StringVarP(&p.RpcUrl, "rpc-url", "r", "http://localhost:8545", "RPC endpoint URL")
	f.StringVar(&p.PrivateKey, "private-key", defaultPrivateKey, "hex encoded private key to use for sending transactions")

	// Wallet parameters.
	f.Uint64VarP(&p.WalletsNumber, "number", "n", 10, "number of wallets to fund")
	f.BoolVar(&p.UseHDDerivation, "hd-derivation", true, "derive wallets to fund from private key in deterministic way")
	f.StringSliceVar(&p.WalletAddresses, "addresses", nil, "comma-separated list of wallet addresses to fund")
	p.FundingAmountInWei = defaultFundingInWei
	f.Var(&flag_loader.BigIntValue{Val: p.FundingAmountInWei}, "eth-amount", "amount of wei to send to each wallet")
	f.StringVar(&p.KeyFile, "key-file", "", "file containing accounts private keys, one per line")

	f.StringVarP(&p.OutputFile, "file", "f", "wallets.json", "output JSON file path for storing addresses and private keys of funded wallets")

	// Marking flags as mutually exclusive
	FundCmd.MarkFlagsMutuallyExclusive("addresses", "number")
	FundCmd.MarkFlagsMutuallyExclusive("addresses", "hd-derivation")
	FundCmd.MarkFlagsMutuallyExclusive("key-file", "addresses")
	FundCmd.MarkFlagsMutuallyExclusive("key-file", "number")
	FundCmd.MarkFlagsMutuallyExclusive("key-file", "hd-derivation")

	// Funder contract parameters.
	f.StringVar(&p.FunderAddress, "contract-address", "", "address of pre-deployed Funder contract")

	params = *p
}

func checkFlags() error {
	// Check rpc url flag.
	if params.RpcUrl == "" {
		panic("RPC URL is empty")
	}
	if err := util.ValidateUrl(params.RpcUrl); err != nil {
		return err
	}

	// Check private key flag.
	if params.PrivateKey == "" {
		return errors.New("the private key is empty")
	}

	// Check that exactly one method is used to specify target accounts
	hasAddresses := len(params.WalletAddresses) > 0
	hasKeyFile := params.KeyFile != ""
	hasNumberFlag := params.WalletsNumber > 0
	if !hasAddresses && !hasKeyFile && !hasNumberFlag {
		return errors.New("must specify target accounts via --addresses, --key-file, or --number")
	}

	minValue := big.NewInt(1000000000)
	if params.FundingAmountInWei != nil && params.FundingAmountInWei.Cmp(minValue) <= 0 {
		return errors.New("the funding amount must be greater than 1000000000")
	}
	if params.OutputFile == "" {
		return errors.New("the output file is not specified")
	}

	return nil
}
