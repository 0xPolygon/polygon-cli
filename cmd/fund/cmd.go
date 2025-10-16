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
	Seed    string

	FunderAddress     string
	Multicall3Address string

	RateLimit float64

	// ERC20 specific parameters
	TokenAddress   string
	TokenAmount    *big.Int
	ApproveSpender string
	ApproveAmount  *big.Int
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
	f := FundCmd.Flags()

	f.StringVarP(&params.RpcUrl, "rpc-url", "r", "http://localhost:8545", "RPC endpoint URL")
	f.StringVar(&params.PrivateKey, "private-key", defaultPrivateKey, "hex encoded private key to use for sending transactions")

	// Wallet parameters.
	f.Uint64VarP(&params.WalletsNumber, "number", "n", 10, "number of wallets to fund")
	f.BoolVar(&params.UseHDDerivation, "hd-derivation", true, "derive wallets to fund from private key in deterministic way")
	f.StringSliceVar(&params.WalletAddresses, "addresses", nil, "comma-separated list of wallet addresses to fund")
	params.FundingAmountInWei = defaultFundingInWei
	f.Var(&flag_loader.BigIntValue{Val: params.FundingAmountInWei}, "eth-amount", "amount of wei to send to each wallet")
	f.StringVar(&params.KeyFile, "key-file", "", "file containing accounts private keys, one per line")
	f.StringVar(&params.Seed, "seed", "", "seed string for deterministic wallet generation (e.g., 'ephemeral_test')")

	// Output parameters.
	f.StringVarP(&params.OutputFile, "file", "f", "wallets.json", "output JSON file path for storing addresses and private keys of funded wallets")

	// ERC20 parameters
	f.StringVar(&params.TokenAddress, "token-address", "", "address of the ERC20 token contract to mint and fund (if provided, enables ERC20 mode)")
	params.TokenAmount = new(big.Int)
	params.TokenAmount.SetString("1000000000000000000", 10) // 1 token
	f.Var(&flag_loader.BigIntValue{Val: params.TokenAmount}, "token-amount", "amount of ERC20 tokens to transfer from private-key wallet to each wallet")
	f.StringVar(&params.ApproveSpender, "approve-spender", "", "address to approve for spending tokens from each funded wallet")
	params.ApproveAmount = new(big.Int)
	params.ApproveAmount.SetString("1000000000000000000000", 10) // 1000 tokens default
	f.Var(&flag_loader.BigIntValue{Val: params.ApproveAmount}, "approve-amount", "amount of ERC20 tokens to approve for the spender")

	// Marking flags as mutually exclusive
	FundCmd.MarkFlagsMutuallyExclusive("addresses", "number")
	FundCmd.MarkFlagsMutuallyExclusive("addresses", "hd-derivation")
	FundCmd.MarkFlagsMutuallyExclusive("addresses", "seed")
	FundCmd.MarkFlagsMutuallyExclusive("key-file", "addresses")
	FundCmd.MarkFlagsMutuallyExclusive("key-file", "number")
	FundCmd.MarkFlagsMutuallyExclusive("key-file", "hd-derivation")
	FundCmd.MarkFlagsMutuallyExclusive("key-file", "seed")
	FundCmd.MarkFlagsMutuallyExclusive("seed", "hd-derivation")

	// contract parameters.
	f.StringVar(&params.FunderAddress, "funder-address", "", "address of pre-deployed funder contract")
	f.StringVar(&params.Multicall3Address, "multicall3-address", "", "address of pre-deployed multicall3 contract")

	// RPC parameters.
	f.Float64Var(&params.RateLimit, "rate-limit", 4, "requests per second limit (use negative value to remove limit)")

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
	hasSeed := params.Seed != ""
	hasNumberWithoutSeed := params.WalletsNumber > 0 && !hasSeed && !hasAddresses && !hasKeyFile

	methodCount := 0
	if hasAddresses {
		methodCount++
	}
	if hasKeyFile {
		methodCount++
	}
	if hasNumberWithoutSeed {
		methodCount++
	}
	if hasSeed {
		methodCount++
	}

	if methodCount == 0 {
		return errors.New("must specify target accounts via --addresses, --key-file, --number, or --seed")
	}
	if methodCount > 1 {
		return errors.New("cannot use multiple wallet specification methods simultaneously")
	}

	// When using seed, require a number of wallets to generate
	if hasSeed && params.WalletsNumber <= 0 {
		return errors.New("when using --seed, must also specify --number > 0 to indicate how many wallets to generate")
	}

	minValue := big.NewInt(1000000000)
	if params.FundingAmountInWei != nil && params.FundingAmountInWei.Cmp(minValue) <= 0 {
		return errors.New("the funding amount must be greater than 1000000000")
	}
	if params.OutputFile == "" {
		return errors.New("the output file is not specified")
	}

	// ERC20 specific validations
	if params.TokenAddress != "" {
		// ERC20 mode - validate token parameters
		if params.TokenAmount == nil || params.TokenAmount.Cmp(big.NewInt(0)) <= 0 {
			return errors.New("token amount must be greater than 0 when using ERC20 mode")
		}
		// Validate approve parameters if provided
		if params.ApproveSpender != "" {
			if params.ApproveAmount == nil || params.ApproveAmount.Cmp(big.NewInt(0)) <= 0 {
				return errors.New("approve amount must be greater than 0 when approve spender is specified")
			}
		}
		// In ERC20 mode, ETH funding is still supported alongside token minting
	}

	return nil
}
