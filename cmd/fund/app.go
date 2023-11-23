package fund

// defaultPrivateKey is the default private key used to fund wallets.
const defaultPrivateKey = "42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa"

// cmdParams holds the command-line parameters for the fund command.
type cmdFundParams struct {
	RpcUrl     *string
	PrivateKey *string

	WalletCount         *uint64
	WalletFundingAmount *float64
	WalletFundingGas    *uint64
	ConcurrencyLevel    *uint64
	OutputFile          *string
}

var params cmdFundParams

func init() {
	p := new(cmdFundParams)
	flagSet := FundCmd.Flags()

	p.RpcUrl = flagSet.StringP("rpc-url", "r", "http://localhost:8545", "The RPC endpoint url")
	p.PrivateKey = flagSet.String("private-key", defaultPrivateKey, "The hex encoded private key that we'll use to send transactions")

	// Wallet parameters.
	p.WalletCount = flagSet.Uint64P("wallets", "w", 2, "The number of wallets to fund")
	p.WalletFundingAmount = flagSet.Float64P("amount", "a", 0.05, "The amount of eth to send to each wallet")
	p.WalletFundingGas = flagSet.Uint64P("gas", "g", 21000, "The cost of funding a wallet")
	p.ConcurrencyLevel = flagSet.Uint64P("concurrency", "c", 2, "The concurrency level for speeding up funding wallets")
	p.OutputFile = flagSet.StringP("file", "f", "wallets.json", "The output JSON file path for storing the addresses and private keys of funded wallets")

	params = *p
}
