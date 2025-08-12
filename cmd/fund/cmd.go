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
	RpcUrl     *string
	PrivateKey *string

	WalletsNumber      *uint64
	UseHDDerivation    *bool
	WalletAddresses    *[]string
	FundingAmountInWei *big.Int
	OutputFile         *string

	KeyFile *string

	FunderAddress *string
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
		params.RpcUrl = flag_loader.GetRpcUrlFlagValue(cmd)
		params.PrivateKey = flag_loader.GetPrivateKeyFlagValue(cmd)
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
	flagSet := FundCmd.Flags()

	p.RpcUrl = flagSet.StringP("rpc-url", "r", "http://localhost:8545", "The RPC endpoint url")
	p.PrivateKey = flagSet.String("private-key", defaultPrivateKey, "The hex encoded private key that we'll use to send transactions")

	// Wallet parameters.
	p.WalletsNumber = flagSet.Uint64P("number", "n", 10, "The number of wallets to fund")
	p.UseHDDerivation = flagSet.Bool("hd-derivation", true, "Derive wallets to fund from the private key in a deterministic way")
	p.WalletAddresses = flagSet.StringSlice("addresses", nil, "Comma-separated list of wallet addresses to fund")
	p.FundingAmountInWei = defaultFundingInWei
	flagSet.Var(&flag_loader.BigIntValue{Val: p.FundingAmountInWei}, "eth-amount", "The amount of wei to send to each wallet")
	p.KeyFile = flagSet.String("key-file", "", "The file containing the accounts private keys, one per line.")

	p.OutputFile = flagSet.StringP("file", "f", "wallets.json", "The output JSON file path for storing the addresses and private keys of funded wallets")

	// Marking flags as mutually exclusive
	FundCmd.MarkFlagsMutuallyExclusive("addresses", "number")
	FundCmd.MarkFlagsMutuallyExclusive("addresses", "hd-derivation")
	FundCmd.MarkFlagsMutuallyExclusive("key-file", "addresses")
	FundCmd.MarkFlagsMutuallyExclusive("key-file", "number")
	FundCmd.MarkFlagsMutuallyExclusive("key-file", "hd-derivation")

	// Funder contract parameters.
	p.FunderAddress = flagSet.String("contract-address", "", "The address of a pre-deployed Funder contract")

	params = *p
}

func checkFlags() error {
	// Check rpc url flag.
	if params.RpcUrl == nil {
		panic("RPC URL is empty")
	}
	if err := util.ValidateUrl(*params.RpcUrl); err != nil {
		return err
	}

	// Check private key flag.
	if params.PrivateKey != nil && *params.PrivateKey == "" {
		return errors.New("the private key is empty")
	}

	// Check that exactly one method is used to specify target accounts
	hasAddresses := params.WalletAddresses != nil && len(*params.WalletAddresses) > 0
	hasKeyFile := params.KeyFile != nil && *params.KeyFile != ""
	hasNumberFlag := params.WalletsNumber != nil && *params.WalletsNumber > 0
	if !hasAddresses && !hasKeyFile && !hasNumberFlag {
		return errors.New("must specify target accounts via --addresses, --key-file, or --number")
	}

	minValue := big.NewInt(1000000000)
	if params.FundingAmountInWei != nil && params.FundingAmountInWei.Cmp(minValue) <= 0 {
		return errors.New("the funding amount must be greater than 1000000000")
	}
	if params.OutputFile != nil && *params.OutputFile == "" {
		return errors.New("the output file is not specified")
	}

	return nil
}
