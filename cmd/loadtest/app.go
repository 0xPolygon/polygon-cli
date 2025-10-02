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
		RPCUrl                     string
		Requests                   int64
		Concurrency                int64
		BatchSize                  uint64
		TimeLimit                  int64
		RandomRecipients           bool
		EthCallOnly                bool
		EthCallOnlyLatestBlock     bool
		ChainID                    uint64
		PrivateKey                 string
		ToAddress                  string
		EthAmountInWei             uint64
		RateLimit                  float64
		AdaptiveRateLimit          bool
		AdaptiveTargetSize         uint64
		AdaptiveRateLimitIncrement uint64
		AdaptiveCycleDuration      uint64
		AdaptiveBackoffFactor      float64
		Modes                      []string
		StoreDataSize              uint64
		Seed                       int64
		LoadtestContractAddress    string
		ERC20Address               string
		ERC721Address              string
		DelAddress                 string
		ForceGasLimit              uint64
		ForceGasPrice              uint64
		ForcePriorityGasPrice      uint64
		ShouldProduceSummary       bool
		SummaryOutputMode          string
		LegacyTransactionMode      bool
		FireAndForget              bool
		RecallLength               uint64
		ContractAddress            string
		ContractCallData           string
		ContractCallPayable        bool
		BlobFeeCap                 uint64
		StartNonce                 uint64
		GasPriceMultiplier         float64
		SendingAccountsCount       uint64
		AccountFundingAmount       *big.Int
		PreFundSendingAccounts     bool
		RefundRemainingFunds       bool
		SendingAccountsFile        string
		Proxy                      string
		WaitForReceipt             bool
		ReceiptRetryMax            uint
		ReceiptRetryInitialDelayMs uint
		MaxBaseFeeWei              uint64
		OutputRawTxOnly            bool

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
		rpcUrl := flag_loader.GetRpcUrlFlagValue(cmd)
		privateKey := flag_loader.GetPrivateKeyFlagValue(cmd)
		if rpcUrl != nil {
			inputLoadTestParams.RPCUrl = *rpcUrl
		}
		if privateKey != nil {
			inputLoadTestParams.PrivateKey = *privateKey
		}
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
	if ltp.RPCUrl == "" {
		return fmt.Errorf("RPC URL is empty")
	}
	if err := util.ValidateUrl(ltp.RPCUrl); err != nil {
		return err
	}

	if ltp.AdaptiveBackoffFactor <= 0.0 {
		return fmt.Errorf("the backoff factor needs to be non-zero positive. Given: %f", ltp.AdaptiveBackoffFactor)
	}

	if ltp.WaitForReceipt && ltp.ReceiptRetryMax <= 1 {
		return fmt.Errorf("when waiting for a receipt, use a max retry greater than 1")
	}

	if ltp.PreFundSendingAccounts && ltp.AccountFundingAmount != nil && ltp.AccountFundingAmount.Uint64() == 0 {
		return fmt.Errorf("a non-zero funding amount is required when pre-funding sending accounts")
	}
	if ltp.EthCallOnly {
		if ltp.PreFundSendingAccounts || ltp.SendingAccountsFile != "" || ltp.SendingAccountsCount > 0 {
			return fmt.Errorf("pre-funding accounts with call only mode doesn't make sense")
		}
		if ltp.WaitForReceipt {
			return fmt.Errorf("waiting for receipts doesn't make sense with call only mode")
		}
	}
	if ltp.GasPriceMultiplier == 0 {
		return fmt.Errorf("gas price multiplier should be non-zero")
	}

	return nil
}

func init() {
	initFlags()
	initSubCommands()
}

func initFlags() {
	ltp := &inputLoadTestParams

	// Persistent flags.
	pf := LoadtestCmd.PersistentFlags()
	pf.StringVarP(&ltp.RPCUrl, "rpc-url", "r", "http://localhost:8545", "the RPC endpoint URL")
	pf.Int64VarP(&ltp.Requests, "requests", "n", 1, "number of requests to perform for the benchmarking session. The default is to just perform a single request which usually leads to non-representative benchmarking results")
	pf.Int64VarP(&ltp.Concurrency, "concurrency", "c", 1, "number of requests to perform concurrently. Default is one request at a time")
	pf.Int64VarP(&ltp.TimeLimit, "time-limit", "t", -1, "maximum number of seconds to spend for benchmarking. Use this to benchmark within a fixed total amount of time. Per default there is no time limit")
	pf.StringVar(&ltp.PrivateKey, "private-key", codeQualityPrivateKey, "the hex encoded private key that we'll use to send transactions")
	pf.Uint64Var(&ltp.ChainID, "chain-id", 0, "the chain ID for the transactions")
	pf.StringVar(&ltp.ToAddress, "to-address", "0xDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEF", "the address that we're going to send to")
	pf.BoolVar(&ltp.RandomRecipients, "random-recipients", false, "when doing a transfer test, should we send to random addresses rather than DEADBEEFx5")
	pf.BoolVar(&ltp.EthCallOnly, "eth-call-only", false, "when using this mode, rather than sending a transaction, we'll just call. This mode is incompatible with adaptive rate limiting, summarization, and a few other features")
	pf.BoolVar(&ltp.EthCallOnlyLatestBlock, "eth-call-only-latest", false, "when using call only mode with recall, should we execute on the latest block or on the original block")
	pf.BoolVar(&ltp.OutputRawTxOnly, "output-raw-tx-only", false, "when using this mode, rather than sending a transaction, we'll just output the raw signed transaction hex. Works with most load test modes except RPC and UniswapV3")
	pf.Uint64Var(&ltp.EthAmountInWei, "eth-amount-in-wei", 0, "the amount of ether in wei to send on every transaction")
	pf.Float64Var(&ltp.RateLimit, "rate-limit", 4, "an overall limit to the number of requests per second. Give a number less than zero to remove this limit all together")
	pf.BoolVar(&ltp.AdaptiveRateLimit, "adaptive-rate-limit", false, "enable AIMD-style congestion control to automatically adjust request rate")
	pf.Uint64Var(&ltp.AdaptiveTargetSize, "adaptive-target-size", 1000, "when using adaptive rate limiting, this value sets the target queue size. If the queue is smaller than this value, we'll speed up. If the queue is smaller than this value, we'll back off")
	pf.Uint64Var(&ltp.AdaptiveRateLimitIncrement, "adaptive-rate-limit-increment", 50, "when using adaptive rate limiting, this flag controls the size of the additive increases")
	pf.Uint64Var(&ltp.AdaptiveCycleDuration, "adaptive-cycle-duration-seconds", 10, "when using adaptive rate limiting, this flag controls how often we check the queue size and adjust the rates")
	pf.Float64Var(&ltp.AdaptiveBackoffFactor, "adaptive-backoff-factor", 2, "when using adaptive rate limiting, this flag controls our multiplicative decrease value")
	pf.Float64Var(&ltp.GasPriceMultiplier, "gas-price-multiplier", 1, "a multiplier to increase or decrease the gas price")
	pf.Int64Var(&ltp.Seed, "seed", 123456, "a seed for generating random values and addresses")
	pf.Uint64Var(&ltp.ForceGasLimit, "gas-limit", 0, "in environments where the gas limit can't be computed on the fly, we can specify it manually. This can also be used to avoid eth_estimateGas")
	pf.Uint64Var(&ltp.ForceGasPrice, "gas-price", 0, "in environments where the gas price can't be determined automatically, we can specify it manually")
	pf.Uint64Var(&ltp.StartNonce, "nonce", 0, "use this flag to manually set the starting nonce")
	pf.Uint64Var(&ltp.ForcePriorityGasPrice, "priority-gas-price", 0, "specify gas tip price in the case of EIP-1559")
	pf.BoolVar(&ltp.ShouldProduceSummary, "summarize", false, "should we produce an execution summary after the load test has finished. If you're running a large load test, this can take a long time")
	pf.Uint64Var(&ltp.BatchSize, "batch-size", 999, "number of batches to perform at a time for receipt fetching. Default is 999 requests at a time")
	pf.StringVar(&ltp.SummaryOutputMode, "output-mode", "text", "format mode for summary output (json | text)")
	pf.BoolVar(&ltp.LegacyTransactionMode, "legacy", false, "send a legacy transaction instead of an EIP1559 transaction")
	pf.BoolVar(&ltp.FireAndForget, "fire-and-forget", false, "send transactions and load without waiting for it to be mined")
	pf.BoolVar(&ltp.FireAndForget, "send-only", false, "alias for --fire-and-forget")

	// Local flags.
	f := LoadtestCmd.Flags()
	f.Uint64Var(&ltp.BlobFeeCap, "blob-fee-cap", 100000, "blob fee cap, or maximum blob fee per chunk, in Gwei")
	f.Uint64Var(&ltp.SendingAccountsCount, "sending-accounts-count", 0, "number of sending accounts to use. This is useful for avoiding pool account queue")
	ltp.AccountFundingAmount = defaultAccountFundingAmount
	f.Var(&flag_loader.BigIntValue{Val: ltp.AccountFundingAmount}, "account-funding-amount", "the amount in wei to fund the sending accounts with. Set to 0 to disable account funding (useful for eth-call-only mode or pre-funded accounts)")
	f.BoolVar(&ltp.PreFundSendingAccounts, "pre-fund-sending-accounts", false, "if set to true, the sending accounts will be funded at the start of the execution, otherwise all accounts will be funded when used for the first time")
	f.BoolVar(&ltp.RefundRemainingFunds, "refund-remaining-funds", false, "if set to true, the funded amount will be refunded to the funding account. Otherwise, the funded amount will remain in the sending accounts")
	f.StringVar(&ltp.SendingAccountsFile, "sending-accounts-file", "", "file containing sending accounts private keys, one per line. This is useful for avoiding pool account queue but also to keep same sending accounts for different execution cycles")
	f.Uint64Var(&ltp.MaxBaseFeeWei, "max-base-fee-wei", 0, "maximum base fee in wei. If base fee exceeds this value, sending tx will be paused and while paused, existing in-flight transactions continue to confirmation, but no additional SendTransaction calls occur. This is useful to avoid sending transactions when network is congested")
	f.StringSliceVarP(&ltp.Modes, "mode", "m", []string{"t"}, `the testing mode to use. It can be multiple like: "d,t"
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
	f.Uint64Var(&ltp.StoreDataSize, "store-data-size", 1024, "if we're in store mode, this controls how many bytes we'll try to store in our contract")
	f.StringVar(&ltp.LoadtestContractAddress, "loadtest-contract-address", "", "address of pre-deployed load test contract")
	f.StringVar(&ltp.ERC20Address, "erc20-address", "", "address of pre-deployed ERC20 contract")
	f.StringVar(&ltp.ERC721Address, "erc721-address", "", "address of pre-deployed ERC721 contract")
	f.Uint64Var(&ltp.RecallLength, "recall-blocks", 50, "number of blocks that we'll attempt to fetch for recall")
	f.StringVar(&ltp.ContractAddress, "contract-address", "", "address of contract that will be used in --mode contract-call. This must be paired up with --mode contract-call and --calldata")
	f.StringVar(&ltp.ContractCallData, "calldata", "", "hex encoded calldata passed in. Format is function signature + arguments encoded together. This must be paired up with --mode contract-call and --contract-address")
	f.BoolVar(&ltp.ContractCallPayable, "contract-call-payable", false, "use this flag if the function is payable, the value amount passed will be from --eth-amount-in-wei. This must be paired up with --mode contract-call and --contract-address")
	f.StringVar(&ltp.Proxy, "proxy", "", "use the proxy specified")
	f.BoolVar(&ltp.WaitForReceipt, "wait-for-receipt", false, "if set to true, the load test will wait for the transaction receipt to be mined. If set to false, the load test will not wait for the transaction receipt and will just send the transaction")
	f.UintVar(&ltp.ReceiptRetryMax, "receipt-retry-max", 30, "maximum number of attempts to poll for transaction receipt when --wait-for-receipt is enabled")
	f.UintVar(&ltp.ReceiptRetryInitialDelayMs, "receipt-retry-initial-delay-ms", 100, "initial delay in milliseconds for receipt polling retry. Uses exponential backoff with jitter")

	// TODO Compression
}

func initSubCommands() {
	LoadtestCmd.AddCommand(uniswapV3LoadTestCmd)
}
