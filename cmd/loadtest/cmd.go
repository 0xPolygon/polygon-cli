package loadtest

import (
	_ "embed"
	"math/big"
	"time"

	"github.com/0xPolygon/polygon-cli/flag"
	"github.com/0xPolygon/polygon-cli/loadtest"
	"github.com/0xPolygon/polygon-cli/loadtest/config"
	"github.com/0xPolygon/polygon-cli/loadtest/uniswapv3"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

const (
	codeQualityPrivateKey = "42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa"
)

var (
	//go:embed loadtestUsage.md
	loadtestUsage string

	//go:embed uniswapv3Usage.md
	uniswapv3Usage string
)

// cfg is the shared loadtest configuration instance.
// CLI flags bind directly to this instance.
var cfg = &config.Config{
	AccountFundingAmount: new(big.Int),
}

// uniswapCfg holds UniswapV3-specific configuration.
var uniswapCfg = &config.UniswapV3Config{}

// LoadtestCmd represents the loadtest command.
var LoadtestCmd = &cobra.Command{
	Use:   "loadtest",
	Short: "Run a generic load test against an Eth/EVM style JSON-RPC endpoint.",
	Long:  loadtestUsage,
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		cfg.RPCUrl, err = flag.GetRPCURL(cmd)
		if err != nil {
			return err
		}
		cfg.PrivateKey, err = flag.GetPrivateKey(cmd)
		if err != nil {
			return err
		}
		zerolog.DurationFieldUnit = time.Second
		zerolog.DurationFieldInteger = true

		return cfg.Validate()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return loadtest.Run(cmd.Context(), cfg)
	},
}

// uniswapv3Cmd represents the uniswapv3 subcommand.
var uniswapv3Cmd = &cobra.Command{
	Use:   "uniswapv3",
	Short: "Run UniswapV3-like load test against an Eth/EVM style JSON-RPC endpoint.",
	Long:  uniswapv3Usage,
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return uniswapCfg.Validate()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Override mode to uniswapv3 and attach UniswapV3 config.
		cfg.Modes = []string{"v3"}
		cfg.UniswapV3 = uniswapCfg

		return loadtest.Run(cmd.Context(), cfg)
	},
}

func init() {
	initPersistentFlags()
	initFlags()
	initSubCommands()
}

func initPersistentFlags() {
	pf := LoadtestCmd.PersistentFlags()
	pf.StringVarP(&cfg.RPCUrl, flag.RPCURL, "r", flag.DefaultRPCURL, "the RPC endpoint URL")
	pf.Int64VarP(&cfg.Requests, "requests", "n", 1, "number of requests to perform for the benchmarking session (default of 1 leads to non-representative results)")
	pf.Int64VarP(&cfg.Concurrency, "concurrency", "c", 1, "number of requests to perform concurrently (default: one at a time)")
	pf.Int64VarP(&cfg.TimeLimit, "time-limit", "t", -1, "maximum seconds to spend benchmarking (default: no limit)")
	pf.StringVar(&cfg.PrivateKey, flag.PrivateKey, codeQualityPrivateKey, "hex encoded private key to use for sending transactions")
	pf.Uint64Var(&cfg.ChainID, "chain-id", 0, "chain ID for the transactions")
	pf.StringVar(&cfg.ToAddress, "to-address", "0xDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEF", "recipient address for transactions")
	pf.BoolVar(&cfg.RandomRecipients, "random-recipients", false, "send to random addresses instead of fixed address in transfer tests")
	pf.BoolVar(&cfg.EthCallOnly, "eth-call-only", false, "call contracts without sending transactions (incompatible with adaptive rate limiting and summarization)")
	pf.BoolVar(&cfg.EthCallOnlyLatestBlock, "eth-call-only-latest", false, "execute on latest block instead of original block in call-only mode with recall")
	pf.BoolVar(&cfg.OutputRawTxOnly, "output-raw-tx-only", false, "output raw signed transaction hex without sending (works with most modes except RPC and UniswapV3)")
	pf.Uint64Var(&cfg.EthAmountInWei, "eth-amount-in-wei", 0, "amount of ether in wei to send per transaction")
	pf.Float64Var(&cfg.RateLimit, "rate-limit", 4, "requests per second limit (use negative value to remove limit)")
	pf.BoolVar(&cfg.AdaptiveRateLimit, "adaptive-rate-limit", false, "enable AIMD-style congestion control to automatically adjust request rate")
	pf.Uint64Var(&cfg.AdaptiveTargetSize, "adaptive-target-size", 1000, "target queue size for adaptive rate limiting (speed up if smaller, back off if larger)")
	pf.Uint64Var(&cfg.AdaptiveRateLimitIncrement, "adaptive-rate-limit-increment", 50, "size of additive increases for adaptive rate limiting")
	pf.Uint64Var(&cfg.AdaptiveCycleDuration, "adaptive-cycle-duration-seconds", 10, "interval in seconds to check queue size and adjust rates for adaptive rate limiting")
	pf.Float64Var(&cfg.AdaptiveBackoffFactor, "adaptive-backoff-factor", 2, "multiplicative decrease factor for adaptive rate limiting")
	pf.Float64Var(&cfg.GasPriceMultiplier, "gas-price-multiplier", 1, "a multiplier to increase or decrease the gas price")
	pf.Int64Var(&cfg.Seed, "seed", 123456, "a seed for generating random values and addresses")
	pf.Uint64Var(&cfg.ForceGasLimit, "gas-limit", 0, "manually specify gas limit (useful to avoid eth_estimateGas or when auto-computation fails)")
	pf.Uint64Var(&cfg.ForceGasPrice, "gas-price", 0, "manually specify gas price (useful when auto-detection fails)")
	pf.Uint64Var(&cfg.StartNonce, "nonce", 0, "use this flag to manually set the starting nonce")
	pf.Uint64Var(&cfg.ForcePriorityGasPrice, "priority-gas-price", 0, "gas tip price for EIP-1559 transactions")
	pf.BoolVar(&cfg.ShouldProduceSummary, "summarize", false, "produce execution summary after load test (can take a long time for large tests)")
	pf.Uint64Var(&cfg.BatchSize, "batch-size", 999, "batch size for receipt fetching (default: 999)")
	pf.StringVar(&cfg.SummaryOutputMode, "output-mode", "text", "format mode for summary output (json | text)")
	pf.BoolVar(&cfg.LegacyTxMode, "legacy", false, "send a legacy transaction instead of an EIP1559 transaction")
	pf.BoolVar(&cfg.FireAndForget, "fire-and-forget", false, "send transactions and load without waiting for it to be mined")
	pf.BoolVar(&cfg.FireAndForget, "send-only", false, "alias for --fire-and-forget")
}

func initFlags() {
	f := LoadtestCmd.Flags()
	f.Uint64Var(&cfg.BlobFeeCap, "blob-fee-cap", 100000, "blob fee cap, or maximum blob fee per chunk, in Gwei")
	f.Uint64Var(&cfg.SendingAccountsCount, "sending-accounts-count", 0, "number of sending accounts to use (avoids pool account queue)")
	f.Var(&flag.BigIntValue{Val: cfg.AccountFundingAmount}, "account-funding-amount", "amount in wei to fund sending accounts (set to 0 to disable)")
	f.BoolVar(&cfg.PreFundSendingAccounts, "pre-fund-sending-accounts", false, "fund all sending accounts at start instead of on first use")
	f.BoolVar(&cfg.RefundRemainingFunds, "refund-remaining-funds", false, "refund remaining balance to funding account after completion")
	f.StringVar(&cfg.SendingAccountsFile, "sending-accounts-file", "", "file with sending account private keys, one per line (avoids pool queue and preserves accounts across runs)")
	f.Uint64Var(&cfg.MaxBaseFeeWei, "max-base-fee-wei", 0, "maximum base fee in wei (pause sending new transactions when exceeded, useful during network congestion)")
	f.StringSliceVarP(&cfg.Modes, "mode", "m", []string{"t"}, `testing mode (can specify multiple like "d,t"):
2, erc20 - send ERC20 tokens
7, erc721 - mint ERC721 tokens
b, blob - send blob transactions
cc, contract-call - make contract calls
d, deploy - deploy contracts
inc, increment - increment a counter
r, random - random modes (excludes: blob, call, recall, rpc, uniswapv3)
R, recall - replay or simulate transactions
rpc - call random rpc methods
s, store - store bytes in a dynamic byte array
t, transaction - send transactions
v3, uniswapv3 - perform UniswapV3 swaps`)
	f.Uint64Var(&cfg.StoreDataSize, "store-data-size", 1024, "number of bytes to store in contract for store mode")
	f.StringVar(&cfg.LoadTestContractAddress, "loadtest-contract-address", "", "address of pre-deployed load test contract")
	f.StringVar(&cfg.ERC20Address, "erc20-address", "", "address of pre-deployed ERC20 contract")
	f.StringVar(&cfg.ERC721Address, "erc721-address", "", "address of pre-deployed ERC721 contract")
	f.Uint64Var(&cfg.RecallLength, "recall-blocks", 50, "number of blocks that we'll attempt to fetch for recall")
	f.Uint64Var(&cfg.BlockBatchSize, "block-batch-size", 25, "number of blocks to fetch per RPC batch request for recall and rpc modes")
	f.StringVar(&cfg.ContractAddress, "contract-address", "", "contract address for --mode contract-call (requires --calldata)")
	f.StringVar(&cfg.ContractCallData, "calldata", "", "hex encoded calldata: function signature + encoded arguments (requires --mode contract-call and --contract-address)")
	f.BoolVar(&cfg.ContractCallPayable, "contract-call-payable", false, "mark function as payable using value from --eth-amount-in-wei (requires --mode contract-call and --contract-address)")
	f.StringVar(&cfg.Proxy, "proxy", "", "use the proxy specified")
	f.BoolVar(&cfg.WaitForReceipt, "wait-for-receipt", false, "wait for transaction receipt to be mined instead of just sending")
	f.UintVar(&cfg.ReceiptRetryMax, "receipt-retry-max", 30, "maximum polling attempts for transaction receipt with --wait-for-receipt")
	f.UintVar(&cfg.ReceiptRetryDelay, "receipt-retry-initial-delay-ms", 100, "initial delay in milliseconds for receipt polling (uses exponential backoff with jitter)")
	f.BoolVar(&cfg.CheckBalanceBeforeFunding, "check-balance-before-funding", false, "check account balance before funding sending accounts (saves gas when accounts are already funded)")
}

func initSubCommands() {
	LoadtestCmd.AddCommand(uniswapv3Cmd)
	initUniswapv3Flags()
}

func initUniswapv3Flags() {
	f := uniswapv3Cmd.Flags()

	// Pre-deployed addresses.
	f.StringVar(&uniswapCfg.FactoryV3, "uniswap-factory-v3-address", "", "address of pre-deployed UniswapFactoryV3 contract")
	f.StringVar(&uniswapCfg.Multicall, "uniswap-multicall-address", "", "address of pre-deployed Multicall contract")
	f.StringVar(&uniswapCfg.ProxyAdmin, "uniswap-proxy-admin-address", "", "address of pre-deployed ProxyAdmin contract")
	f.StringVar(&uniswapCfg.TickLens, "uniswap-tick-lens-address", "", "address of pre-deployed TickLens contract")
	f.StringVar(&uniswapCfg.NFTDescriptorLib, "uniswap-nft-descriptor-lib-address", "", "address of pre-deployed NFTDescriptor library contract")
	f.StringVar(&uniswapCfg.NonfungibleTokenPositionDescriptor, "uniswap-nft-position-descriptor-address", "", "address of pre-deployed NonfungibleTokenPositionDescriptor contract")
	f.StringVar(&uniswapCfg.TransparentUpgradeableProxy, "uniswap-upgradeable-proxy-address", "", "address of pre-deployed TransparentUpgradeableProxy contract")
	f.StringVar(&uniswapCfg.NonfungiblePositionManager, "uniswap-non-fungible-position-manager-address", "", "address of pre-deployed NonfungiblePositionManager contract")
	f.StringVar(&uniswapCfg.Migrator, "uniswap-migrator-address", "", "address of pre-deployed Migrator contract")
	f.StringVar(&uniswapCfg.Staker, "uniswap-staker-address", "", "address of pre-deployed Staker contract")
	f.StringVar(&uniswapCfg.QuoterV2, "uniswap-quoter-v2-address", "", "address of pre-deployed QuoterV2 contract")
	f.StringVar(&uniswapCfg.SwapRouter, "uniswap-swap-router-address", "", "address of pre-deployed SwapRouter contract")
	f.StringVar(&uniswapCfg.WETH9, "weth9-address", "", "address of pre-deployed WETH9 contract")
	f.StringVar(&uniswapCfg.PoolToken0, "uniswap-pool-token-0-address", "", "address of pre-deployed ERC20 contract used in Uniswap pool Token0 // Token1")
	f.StringVar(&uniswapCfg.PoolToken1, "uniswap-pool-token-1-address", "", "address of pre-deployed ERC20 contract used in Uniswap pool Token0 // Token1")

	// Pool and swap parameters.
	f.Float64VarP(&uniswapCfg.PoolFees, "pool-fees", "f", float64(uniswapv3.StandardTier), "trading fees for UniswapV3 liquidity pool swaps (e.g. 0.3 means 0.3%)")
	f.Uint64VarP(&uniswapCfg.SwapAmountInput, "swap-amount", "a", uniswapv3.SwapAmountInput.Uint64(), "amount of inbound token given as swap input")
}
