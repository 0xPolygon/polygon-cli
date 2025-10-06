package loadtest

import (
	"crypto/ecdsa"
	_ "embed"
	"fmt"
	"math/big"
	"math/rand"
	"sync"
	"time"

	"github.com/0xPolygon/polygon-cli/flag"
	"github.com/0xPolygon/polygon-cli/rpctypes"
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
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		inputLoadTestParams.RPCUrl, err = flag.GetRPCURL(cmd)
		if err != nil {
			return err
		}
		inputLoadTestParams.PrivateKey, err = flag.GetPrivateKey(cmd)
		if err != nil {
			return err
		}
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
	pf.StringVarP(&ltp.RPCUrl, flag.RPCURL, "r", flag.DefaultRPCURL, "the RPC endpoint URL")
	pf.Int64VarP(&ltp.Requests, "requests", "n", 1, "number of requests to perform for the benchmarking session (default of 1 leads to non-representative results)")
	pf.Int64VarP(&ltp.Concurrency, "concurrency", "c", 1, "number of requests to perform concurrently (default: one at a time)")
	pf.Int64VarP(&ltp.TimeLimit, "time-limit", "t", -1, "maximum seconds to spend benchmarking (default: no limit)")
	pf.StringVar(&ltp.PrivateKey, flag.PrivateKey, codeQualityPrivateKey, "hex encoded private key to use for sending transactions")
	pf.Uint64Var(&ltp.ChainID, "chain-id", 0, "chain ID for the transactions")
	pf.StringVar(&ltp.ToAddress, "to-address", "0xDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEF", "recipient address for transactions")
	pf.BoolVar(&ltp.RandomRecipients, "random-recipients", false, "send to random addresses instead of fixed address in transfer tests")
	pf.BoolVar(&ltp.EthCallOnly, "eth-call-only", false, "call contracts without sending transactions (incompatible with adaptive rate limiting and summarization)")
	pf.BoolVar(&ltp.EthCallOnlyLatestBlock, "eth-call-only-latest", false, "execute on latest block instead of original block in call-only mode with recall")
	pf.BoolVar(&ltp.OutputRawTxOnly, "output-raw-tx-only", false, "output raw signed transaction hex without sending (works with most modes except RPC and UniswapV3)")
	pf.Uint64Var(&ltp.EthAmountInWei, "eth-amount-in-wei", 0, "amount of ether in wei to send per transaction")
	pf.Float64Var(&ltp.RateLimit, "rate-limit", 4, "requests per second limit (use negative value to remove limit)")
	pf.BoolVar(&ltp.AdaptiveRateLimit, "adaptive-rate-limit", false, "enable AIMD-style congestion control to automatically adjust request rate")
	pf.Uint64Var(&ltp.AdaptiveTargetSize, "adaptive-target-size", 1000, "target queue size for adaptive rate limiting (speed up if smaller, back off if larger)")
	pf.Uint64Var(&ltp.AdaptiveRateLimitIncrement, "adaptive-rate-limit-increment", 50, "size of additive increases for adaptive rate limiting")
	pf.Uint64Var(&ltp.AdaptiveCycleDuration, "adaptive-cycle-duration-seconds", 10, "interval in seconds to check queue size and adjust rates for adaptive rate limiting")
	pf.Float64Var(&ltp.AdaptiveBackoffFactor, "adaptive-backoff-factor", 2, "multiplicative decrease factor for adaptive rate limiting")
	pf.Float64Var(&ltp.GasPriceMultiplier, "gas-price-multiplier", 1, "a multiplier to increase or decrease the gas price")
	pf.Int64Var(&ltp.Seed, "seed", 123456, "a seed for generating random values and addresses")
	pf.Uint64Var(&ltp.ForceGasLimit, "gas-limit", 0, "manually specify gas limit (useful to avoid eth_estimateGas or when auto-computation fails)")
	pf.Uint64Var(&ltp.ForceGasPrice, "gas-price", 0, "manually specify gas price (useful when auto-detection fails)")
	pf.Uint64Var(&ltp.StartNonce, "nonce", 0, "use this flag to manually set the starting nonce")
	pf.Uint64Var(&ltp.ForcePriorityGasPrice, "priority-gas-price", 0, "gas tip price for EIP-1559 transactions")
	pf.BoolVar(&ltp.ShouldProduceSummary, "summarize", false, "produce execution summary after load test (can take a long time for large tests)")
	pf.Uint64Var(&ltp.BatchSize, "batch-size", 999, "batch size for receipt fetching (default: 999)")
	pf.StringVar(&ltp.SummaryOutputMode, "output-mode", "text", "format mode for summary output (json | text)")
	pf.BoolVar(&ltp.LegacyTransactionMode, "legacy", false, "send a legacy transaction instead of an EIP1559 transaction")
	pf.BoolVar(&ltp.FireAndForget, "fire-and-forget", false, "send transactions and load without waiting for it to be mined")
	pf.BoolVar(&ltp.FireAndForget, "send-only", false, "alias for --fire-and-forget")

	// Local flags.
	f := LoadtestCmd.Flags()
	f.Uint64Var(&ltp.BlobFeeCap, "blob-fee-cap", 100000, "blob fee cap, or maximum blob fee per chunk, in Gwei")
	f.Uint64Var(&ltp.SendingAccountsCount, "sending-accounts-count", 0, "number of sending accounts to use (avoids pool account queue)")
	ltp.AccountFundingAmount = defaultAccountFundingAmount
	f.Var(&flag.BigIntValue{Val: ltp.AccountFundingAmount}, "account-funding-amount", "amount in wei to fund sending accounts (set to 0 to disable)")
	f.BoolVar(&ltp.PreFundSendingAccounts, "pre-fund-sending-accounts", false, "fund all sending accounts at start instead of on first use")
	f.BoolVar(&ltp.RefundRemainingFunds, "refund-remaining-funds", false, "refund remaining balance to funding account after completion")
	f.StringVar(&ltp.SendingAccountsFile, "sending-accounts-file", "", "file with sending account private keys, one per line (avoids pool queue and preserves accounts across runs)")
	f.Uint64Var(&ltp.MaxBaseFeeWei, "max-base-fee-wei", 0, "maximum base fee in wei (pause sending new transactions when exceeded, useful during network congestion)")
	f.StringSliceVarP(&ltp.Modes, "mode", "m", []string{"t"}, `testing mode (can specify multiple like "d,t"):
2, erc20 - send ERC20 tokens
7, erc721 - mint ERC721 tokens
b, blob - send blob transactions
cc, contract-call - make contract calls
d, deploy - deploy contracts
inc, increment - increment a counter
r, random - random modes (excludes: blob, call, inscription, recall, rpc, uniswapv3)
R, recall - replay or simulate transactions
rpc - call random rpc methods
s, store - store bytes in a dynamic byte array
t, transaction - send transactions
v3, uniswapv3 - perform UniswapV3 swaps`)
	f.Uint64Var(&ltp.StoreDataSize, "store-data-size", 1024, "number of bytes to store in contract for store mode")
	f.StringVar(&ltp.LoadtestContractAddress, "loadtest-contract-address", "", "address of pre-deployed load test contract")
	f.StringVar(&ltp.ERC20Address, "erc20-address", "", "address of pre-deployed ERC20 contract")
	f.StringVar(&ltp.ERC721Address, "erc721-address", "", "address of pre-deployed ERC721 contract")
	f.Uint64Var(&ltp.RecallLength, "recall-blocks", 50, "number of blocks that we'll attempt to fetch for recall")
	f.StringVar(&ltp.ContractAddress, "contract-address", "", "contract address for --mode contract-call (requires --calldata)")
	f.StringVar(&ltp.ContractCallData, "calldata", "", "hex encoded calldata: function signature + encoded arguments (requires --mode contract-call and --contract-address)")
	f.BoolVar(&ltp.ContractCallPayable, "contract-call-payable", false, "mark function as payable using value from --eth-amount-in-wei (requires --mode contract-call and --contract-address)")
	f.StringVar(&ltp.Proxy, "proxy", "", "use the proxy specified")
	f.BoolVar(&ltp.WaitForReceipt, "wait-for-receipt", false, "wait for transaction receipt to be mined instead of just sending")
	f.UintVar(&ltp.ReceiptRetryMax, "receipt-retry-max", 30, "maximum polling attempts for transaction receipt with --wait-for-receipt")
	f.UintVar(&ltp.ReceiptRetryInitialDelayMs, "receipt-retry-initial-delay-ms", 100, "initial delay in milliseconds for receipt polling (uses exponential backoff with jitter)")

	// TODO Compression
}

func initSubCommands() {
	LoadtestCmd.AddCommand(uniswapV3LoadTestCmd)
}
