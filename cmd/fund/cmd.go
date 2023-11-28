package fund

import (
	"errors"
	"math"

	_ "embed"

	"github.com/maticnetwork/polygon-cli/util"
	"github.com/spf13/cobra"
)

// The default private key used to send transactions.
const defaultPrivateKey = "0x42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa"

// The default mnemonic  used to derive wallets.
const defaultMnemonic = "code code code code code code code code code code code quality"

// The default password used to create a wallet for HD derivation.
const defaultPassword = "password"

// cmdParams holds the command-line parameters for the fund command.
type cmdFundParams struct {
	RpcUrl     *string
	PrivateKey *string

	WalletsNumber      *uint64
	UseHDDerivation    *bool
	WalletAddresses    *[]string
	FundingAmountInEth *float64
	OutputFile         *string

	FunderAddress *string
}

var (
	//go:embed usage.md
	usage  string
	params cmdFundParams
)

// FundCmd represents the fund command.
var FundCmd = &cobra.Command{
	Use:   "fund",
	Short: "Bulk fund crypto wallets automatically.",
	Long:  usage,
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
	p.FundingAmountInEth = flagSet.Float64P("eth-amount", "a", 0.05, "The amount of ether to send to each wallet")
	p.OutputFile = flagSet.StringP("file", "f", "wallets.json", "The output JSON file path for storing the addresses and private keys of funded wallets")

	// Marking flags as mutually exclusive
	FundCmd.MarkFlagsMutuallyExclusive("addresses", "number")
	FundCmd.MarkFlagsMutuallyExclusive("addresses", "hd-derivation")

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

	// Check wallet flags.
	if params.WalletsNumber != nil && *params.WalletsNumber == 0 {
		return errors.New("the number of wallets to fund is set to zero")
	}
	if params.FundingAmountInEth != nil && math.Abs(*params.FundingAmountInEth) <= 1e-9 {
		return errors.New("the amount of eth to send to each wallet is set to zero")
	}
	if params.OutputFile != nil && *params.OutputFile == "" {
		return errors.New("the output file is not specified")
	}

	return nil
}
