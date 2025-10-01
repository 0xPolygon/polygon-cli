package loadtest

import (
	"crypto/ecdsa"
	_ "embed"
	"fmt"
	"math/big"
	"math/rand"
	"sync"
	"time"

	"github.com/0xPolygon/polygon-cli/cmd/flag_loader"
	"github.com/0xPolygon/polygon-cli/rpctypes"
	"github.com/0xPolygon/polygon-cli/util"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"golang.org/x/time/rate"
)

type (
	blockSummary struct {
		Block     *rpctypes.RawBlockResponse
		Receipts  map[ethcommon.Hash]rpctypes.RawTxReceipt
		Latencies map[uint64]time.Duration
	}
	hexwordReader struct {
	}
	loadTestSample struct {
		GoRoutineID int64
		RequestID   int64
		RequestTime time.Time
		WaitTime    time.Duration // Wait time for transaction to be broadcasted
		Receipt     string
		IsError     bool
		Nonce       uint64
	}
	loadTestParams struct {
		// inputs
		RPCUrl                     *string
		Requests                   *int64
		Concurrency                *int64
		BatchSize                  *uint64
		TimeLimit                  *int64
		RandomRecipients           *bool
		EthCallOnly                *bool
		EthCallOnlyLatestBlock     *bool
		ChainID                    *uint64
		PrivateKey                 *string
		ToAddress                  *string
		EthAmountInWei             *uint64
		RateLimit                  *float64
		AdaptiveRateLimit          *bool
		AdaptiveTargetSize         *uint64
		AdaptiveRateLimitIncrement *uint64
		AdaptiveCycleDuration      *uint64
		AdaptiveBackoffFactor      *float64
		Modes                      *[]string
		StoreDataSize              *uint64
		Seed                       *int64
		LoadtestContractAddress    *string
		ERC20Address               *string
		ERC721Address              *string
		DelAddress                 *string
		ForceGasLimit              *uint64
		ForceGasPrice              *uint64
		ForcePriorityGasPrice      *uint64
		ShouldProduceSummary       *bool
		SummaryOutputMode          *string
		LegacyTransactionMode      *bool
		FireAndForget              *bool
		RecallLength               *uint64
		ContractAddress            *string
		ContractCallData           *string
		ContractCallPayable        *bool
		BlobFeeCap                 *uint64
		StartNonce                 *uint64
		GasPriceMultiplier         *float64
		SendingAccountsCount       *uint64
		AccountFundingAmount       *big.Int
		PreFundSendingAccounts     *bool
		RefundRemainingFunds       *bool
		SendingAccountsFile        *string
		Proxy                      *string
		WaitForReceipt             *bool
		ReceiptRetryMax            *uint
		ReceiptRetryInitialDelayMs *uint
		MaxBaseFeeWei              *uint64
		OutputRawTxOnly            *bool

		// Computed
		CurrentGasPrice       *big.Int
		CurrentGasTipCap      *big.Int
		CurrentNonce          *uint64
		ECDSAPrivateKey       *ecdsa.PrivateKey
		FromETHAddress        *ethcommon.Address
		ToETHAddress          *ethcommon.Address
		ContractETHAddress    *ethcommon.Address
		SendAmount            *big.Int
		ChainSupportBaseFee   bool
		Mode                  loadTestMode
		ParsedModes           []loadTestMode
		MultiMode             bool
		BigGasPriceMultiplier *big.Float
	}
)

var (
	//go:embed loadtestUsage.md
	loadTestUsage        string
	inputLoadTestParams  loadTestParams
	loadTestResults      []loadTestSample
	loadTestResultsMutex sync.RWMutex
	startBlockNumber     uint64
	finalBlockNumber     uint64
	rl                   *rate.Limiter
	accountPool          *AccountPool

	hexwords = []byte{
		0x00, 0x0F, 0xF1, 0xCE,
		0x00, 0xBA, 0xB1, 0x0C,
		0x1B, 0xAD, 0xB0, 0x02,
		0x8B, 0xAD, 0xF0, 0x0D,
		0xAB, 0xAD, 0xBA, 0xBE,
		0xB1, 0x05, 0xF0, 0x0D,
		0xB1, 0x6B, 0x00, 0xB5,
		0x0B, 0x00, 0xB1, 0x35,
		0xBA, 0xAA, 0xAA, 0xAD,
		0xBA, 0xAD, 0xF0, 0x0D,
		0xBA, 0xD2, 0x22, 0x22,
		0xBA, 0xDD, 0xCA, 0xFE,
		0xCA, 0xFE, 0xB0, 0xBA,
		0xB0, 0xBA, 0xBA, 0xBE,
		0xBE, 0xEF, 0xBA, 0xBE,
		0xC0, 0x00, 0x10, 0xFF,
		0xCA, 0xFE, 0xBA, 0xBE,
		0xCA, 0xFE, 0xD0, 0x0D,
		0xCE, 0xFA, 0xED, 0xFE,
		0x0D, 0x15, 0xEA, 0x5E,
		0xDA, 0xBB, 0xAD, 0x00,
		0xDE, 0xAD, 0x2B, 0xAD,
		0xDE, 0xAD, 0xBA, 0xAD,
		0xDE, 0xAD, 0xBA, 0xBE,
		0xDE, 0xAD, 0xBE, 0xAF,
		0xDE, 0xAD, 0xBE, 0xEF,
		0xDE, 0xAD, 0xC0, 0xDE,
		0xDE, 0xAD, 0xDE, 0xAD,
		0xDE, 0xAD, 0xD0, 0x0D,
		0xDE, 0xAD, 0xFA, 0x11,
		0xDE, 0xAD, 0x10, 0xCC,
		0xDE, 0xAD, 0xFE, 0xED,
		0xDE, 0xCA, 0xFB, 0xAD,
		0xDE, 0xFE, 0xC8, 0xED,
		0xD0, 0xD0, 0xCA, 0xCA,
		0xE0, 0x11, 0xCF, 0xD0,
		0xFA, 0xCE, 0xFE, 0xED,
		0xFB, 0xAD, 0xBE, 0xEF,
		0xFE, 0xE1, 0xDE, 0xAD,
		0xFE, 0xED, 0xBA, 0xBE,
		0xFE, 0xED, 0xC0, 0xDE,
		0xFF, 0xBA, 0xDD, 0x11,
		0xF0, 0x0D, 0xBA, 0xBE,
	}

	randSrc                     *rand.Rand
	defaultAccountFundingAmount = new(big.Int).SetUint64(0) // 1 ETH
)

// LoadtestCmd represents the loadtest command
var LoadtestCmd = &cobra.Command{
	Use:   "loadtest",
	Short: "Run a generic load test against an Eth/EVM style JSON-RPC endpoint.",
	Long:  loadTestUsage,
	Args:  cobra.NoArgs,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		inputLoadTestParams.RPCUrl = flag_loader.GetRpcUrlFlagValue(cmd)
		inputLoadTestParams.PrivateKey = flag_loader.GetPrivateKeyFlagValue(cmd)
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		zerolog.DurationFieldUnit = time.Second
		zerolog.DurationFieldInteger = true

		return checkLoadtestFlags()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return runLoadTest(cmd.Context())
	},
}

func checkLoadtestFlags() error {
	ltp := inputLoadTestParams

	// Check `rpc-url` flag.
	if ltp.RPCUrl == nil {
		return fmt.Errorf("RPC URL is empty")
	}
	if err := util.ValidateUrl(*ltp.RPCUrl); err != nil {
		return err
	}

	if ltp.AdaptiveBackoffFactor != nil && *ltp.AdaptiveBackoffFactor <= 0.0 {
		return fmt.Errorf("the backoff factor needs to be non-zero positive. Given: %f", *ltp.AdaptiveBackoffFactor)
	}

	if *ltp.WaitForReceipt && *ltp.ReceiptRetryMax <= 1 {
		return fmt.Errorf("when waiting for a receipt, use a max retry greater than 1")
	}

	if *ltp.PreFundSendingAccounts && ltp.AccountFundingAmount != nil && ltp.AccountFundingAmount.Uint64() == 0 {
		return fmt.Errorf("a non-zero funding amount is required when pre-funding sending accounts")
	}
	if *ltp.EthCallOnly {
		if *ltp.PreFundSendingAccounts || *ltp.SendingAccountsFile != "" || *ltp.SendingAccountsCount > 0 {
			return fmt.Errorf("pre-funding accounts with call only mode doesn't make sense")
		}
		if *ltp.WaitForReceipt {
			return fmt.Errorf("waiting for receipts doesn't make sense with call only mode")
		}
	}
	if *ltp.GasPriceMultiplier == 0 {
		return fmt.Errorf("gas price multiplier should be non-zero")
	}

	return nil
}

func init() {
	initFlags()
	initSubCommands()
}

func initFlags() {
	ltp := new(loadTestParams)

	// Persistent flags.
	ltp.RPCUrl = LoadtestCmd.PersistentFlags().StringP("rpc-url", "r", "http://localhost:8545", "The RPC endpoint url")
	ltp.Requests = LoadtestCmd.PersistentFlags().Int64P("requests", "n", 1, "Number of requests to perform for the benchmarking session. The default is to just perform a single request which usually leads to non-representative benchmarking results.")
	ltp.Concurrency = LoadtestCmd.PersistentFlags().Int64P("concurrency", "c", 1, "Number of requests to perform concurrently. Default is one request at a time.")
	ltp.TimeLimit = LoadtestCmd.PersistentFlags().Int64P("time-limit", "t", -1, "Maximum number of seconds to spend for benchmarking. Use this to benchmark within a fixed total amount of time. Per default there is no time limit.")
	ltp.PrivateKey = LoadtestCmd.PersistentFlags().String("private-key", codeQualityPrivateKey, "The hex encoded private key that we'll use to send transactions")
	ltp.ChainID = LoadtestCmd.PersistentFlags().Uint64("chain-id", 0, "The chain id for the transactions.")
	ltp.ToAddress = LoadtestCmd.PersistentFlags().String("to-address", "0xDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEF", "The address that we're going to send to")
	ltp.RandomRecipients = LoadtestCmd.PersistentFlags().Bool("random-recipients", false, "When doing a transfer test, should we send to random addresses rather than DEADBEEFx5")
	ltp.EthCallOnly = LoadtestCmd.PersistentFlags().Bool("eth-call-only", false, "When using this mode, rather than sending a transaction, we'll just call. This mode is incompatible with adaptive rate limiting, summarization, and a few other features.")
	ltp.EthCallOnlyLatestBlock = LoadtestCmd.PersistentFlags().Bool("eth-call-only-latest", false, "When using call only mode with recall, should we execute on the latest block or on the original block")
	ltp.OutputRawTxOnly = LoadtestCmd.PersistentFlags().Bool("output-raw-tx-only", false, "When using this mode, rather than sending a transaction, we'll just output the raw signed transaction hex. Works with most load test modes except RPC and UniswapV3.")
	ltp.EthAmountInWei = LoadtestCmd.PersistentFlags().Uint64("eth-amount-in-wei", 0, "The amount of ether in wei to send on every transaction")
	ltp.RateLimit = LoadtestCmd.PersistentFlags().Float64("rate-limit", 4, "An overall limit to the number of requests per second. Give a number less than zero to remove this limit all together")
	ltp.AdaptiveRateLimit = LoadtestCmd.PersistentFlags().Bool("adaptive-rate-limit", false, "Enable AIMD-style congestion control to automatically adjust request rate")
	ltp.AdaptiveTargetSize = LoadtestCmd.PersistentFlags().Uint64("adaptive-target-size", 1000, "When using adaptive rate limiting, this value sets the target queue size. If the queue is smaller than this value, we'll speed up. If the queue is smaller than this value, we'll back off.")
	ltp.AdaptiveRateLimitIncrement = LoadtestCmd.PersistentFlags().Uint64("adaptive-rate-limit-increment", 50, "When using adaptive rate limiting, this flag controls the size of the additive increases.")
	ltp.AdaptiveCycleDuration = LoadtestCmd.PersistentFlags().Uint64("adaptive-cycle-duration-seconds", 10, "When using adaptive rate limiting, this flag controls how often we check the queue size and adjust the rates")
	ltp.AdaptiveBackoffFactor = LoadtestCmd.PersistentFlags().Float64("adaptive-backoff-factor", 2, "When using adaptive rate limiting, this flag controls our multiplicative decrease value.")
	ltp.GasPriceMultiplier = LoadtestCmd.PersistentFlags().Float64("gas-price-multiplier", 1, "A multiplier to increase or decrease the gas price")
	ltp.Seed = LoadtestCmd.PersistentFlags().Int64("seed", 123456, "A seed for generating random values and addresses")
	ltp.ForceGasLimit = LoadtestCmd.PersistentFlags().Uint64("gas-limit", 0, "In environments where the gas limit can't be computed on the fly, we can specify it manually. This can also be used to avoid eth_estimateGas")
	ltp.ForceGasPrice = LoadtestCmd.PersistentFlags().Uint64("gas-price", 0, "In environments where the gas price can't be determined automatically, we can specify it manually")
	ltp.StartNonce = LoadtestCmd.PersistentFlags().Uint64("nonce", 0, "Use this flag to manually set the starting nonce")
	ltp.ForcePriorityGasPrice = LoadtestCmd.PersistentFlags().Uint64("priority-gas-price", 0, "Specify Gas Tip Price in the case of EIP-1559")
	ltp.ShouldProduceSummary = LoadtestCmd.PersistentFlags().Bool("summarize", false, "Should we produce an execution summary after the load test has finished. If you're running a large load test, this can take a long time")
	ltp.BatchSize = LoadtestCmd.PersistentFlags().Uint64("batch-size", 999, "Number of batches to perform at a time for receipt fetching. Default is 999 requests at a time.")
	ltp.SummaryOutputMode = LoadtestCmd.PersistentFlags().String("output-mode", "text", "Format mode for summary output (json | text)")
	ltp.LegacyTransactionMode = LoadtestCmd.PersistentFlags().Bool("legacy", false, "Send a legacy transaction instead of an EIP1559 transaction.")
	ltp.FireAndForget = LoadtestCmd.PersistentFlags().Bool("fire-and-forget", false, "Send transactions and load without waiting for it to be mined.")
	LoadtestCmd.PersistentFlags().BoolVar(ltp.FireAndForget, "send-only", false, "Alias for --fire-and-forget.")
	ltp.BlobFeeCap = LoadtestCmd.Flags().Uint64("blob-fee-cap", 100000, "The blob fee cap, or the maximum blob fee per chunk, in Gwei.")
	ltp.SendingAccountsCount = LoadtestCmd.Flags().Uint64("sending-accounts-count", 0, "The number of sending accounts to use. This is useful for avoiding pool account queue.")
	ltp.AccountFundingAmount = defaultAccountFundingAmount
	LoadtestCmd.Flags().Var(&flag_loader.BigIntValue{Val: ltp.AccountFundingAmount}, "account-funding-amount", "The amount in wei to fund the sending accounts with. Set to 0 to disable account funding (useful for eth-call-only mode or pre-funded accounts).")
	ltp.PreFundSendingAccounts = LoadtestCmd.Flags().Bool("pre-fund-sending-accounts", false, "If set to true, the sending accounts will be funded at the start of the execution, otherwise all accounts will be funded when used for the first time.")
	ltp.RefundRemainingFunds = LoadtestCmd.Flags().Bool("refund-remaining-funds", false, "If set to true, the funded amount will be refunded to the funding account. Otherwise, the funded amount will remain in the sending accounts.")
	ltp.SendingAccountsFile = LoadtestCmd.Flags().String("sending-accounts-file", "", "The file containing the sending accounts private keys, one per line. This is useful for avoiding pool account queue but also to keep the same sending accounts for different execution cycles.")
	ltp.MaxBaseFeeWei = LoadtestCmd.Flags().Uint64("max-base-fee-wei", 0, "The maximum base fee in wei. If the base fee exceeds this value, sending tx will be paused and while paused, existing in-flight transactions continue to confirmation, but no additional SendTransaction calls occur. This is useful to avoid sending transactions when the network is congested.")

	// Local flags.
	ltp.Modes = LoadtestCmd.Flags().StringSliceP("mode", "m", []string{"t"}, `The testing mode to use. It can be multiple like: "d,t"
2, erc20 - Send ERC20 tokens
7, erc721 - Mint ERC721 tokens
b, blob - Send blob transactions
cc, contract-call - Make contract calls
d, deploy - Deploy contracts
inc, increment - Increment a counter
r, random - Random modes (does not include the following modes: blob, call, inscription, recall, rpc, uniswapv3)
R, recall - Replay or simulate transactions
rpc - Call random rpc methods
s, store - Store bytes in a dynamic byte array
t, transaction - Send transactions
v3, uniswapv3 - Perform UniswapV3 swaps`)
	ltp.StoreDataSize = LoadtestCmd.Flags().Uint64("store-data-size", 1024, "If we're in store mode, this controls how many bytes we'll try to store in our contract")
	ltp.LoadtestContractAddress = LoadtestCmd.Flags().String("loadtest-contract-address", "", "The address of a pre-deployed load test contract")
	ltp.ERC20Address = LoadtestCmd.Flags().String("erc20-address", "", "The address of a pre-deployed ERC20 contract")
	ltp.ERC721Address = LoadtestCmd.Flags().String("erc721-address", "", "The address of a pre-deployed ERC721 contract")
	ltp.RecallLength = LoadtestCmd.Flags().Uint64("recall-blocks", 50, "The number of blocks that we'll attempt to fetch for recall")
	ltp.ContractAddress = LoadtestCmd.Flags().String("contract-address", "", "The address of the contract that will be used in --mode contract-call. This must be paired up with --mode contract-call and --calldata")
	ltp.ContractCallData = LoadtestCmd.Flags().String("calldata", "", "The hex encoded calldata passed in. The format is function signature + arguments encoded together. This must be paired up with --mode contract-call and --contract-address")
	ltp.ContractCallPayable = LoadtestCmd.Flags().Bool("contract-call-payable", false, "Use this flag if the function is payable, the value amount passed will be from --eth-amount-in-wei. This must be paired up with --mode contract-call and --contract-address")
	ltp.Proxy = LoadtestCmd.Flags().String("proxy", "", "Use the proxy specified")
	ltp.WaitForReceipt = LoadtestCmd.Flags().Bool("wait-for-receipt", false, "If set to true, the load test will wait for the transaction receipt to be mined. If set to false, the load test will not wait for the transaction receipt and will just send the transaction.")
	ltp.ReceiptRetryMax = LoadtestCmd.Flags().Uint("receipt-retry-max", 30, "Maximum number of attempts to poll for transaction receipt when --wait-for-receipt is enabled.")
	ltp.ReceiptRetryInitialDelayMs = LoadtestCmd.Flags().Uint("receipt-retry-initial-delay-ms", 100, "Initial delay in milliseconds for receipt polling retry. Uses exponential backoff with jitter.")

	inputLoadTestParams = *ltp

	// TODO Compression
}

func initSubCommands() {
	LoadtestCmd.AddCommand(uniswapV3LoadTestCmd)
}
