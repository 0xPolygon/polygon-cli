package loadtest

import (
	"crypto/ecdsa"
	_ "embed"
	"fmt"
	"math/big"
	"math/rand"
	"sync"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/maticnetwork/polygon-cli/rpctypes"
	"github.com/maticnetwork/polygon-cli/util"
	"github.com/spf13/cobra"
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
		WaitTime    time.Duration
		Receipt     string
		IsError     bool
		Nonce       uint64
	}
	loadTestParams struct {
		// inputs
		RPCUrl                              *string
		Requests                            *int64
		Concurrency                         *int64
		BatchSize                           *uint64
		TimeLimit                           *int64
		ToRandom                            *bool
		CallOnly                            *bool
		CallOnlyLatestBlock                 *bool
		ChainID                             *uint64
		PrivateKey                          *string
		ToAddress                           *string
		HexSendAmount                       *string
		RateLimit                           *float64
		AdaptiveRateLimit                   *bool
		SteadyStateTxPoolSize               *uint64
		AdaptiveRateLimitIncrement          *uint64
		AdaptiveCycleDuration               *uint64
		AdaptiveBackoffFactor               *float64
		Modes                               *[]string
		Function                            *uint64
		Iterations                          *uint64
		ByteCount                           *uint64
		Seed                                *int64
		LtAddress                           *string
		ERC20Address                        *string
		ERC721Address                       *string
		DelAddress                          *string
		ContractCallNumberOfBlocksToWaitFor *uint64
		ContractCallBlockInterval           *uint64
		ForceContractDeploy                 *bool
		ForceGasLimit                       *uint64
		ForceGasPrice                       *uint64
		ForcePriorityGasPrice               *uint64
		ShouldProduceSummary                *bool
		SummaryOutputMode                   *string
		LegacyTransactionMode               *bool
		RecallLength                        *uint64

		// Computed
		CurrentGasPrice     *big.Int
		CurrentGasTipCap    *big.Int
		CurrentNonce        *uint64
		ECDSAPrivateKey     *ecdsa.PrivateKey
		FromETHAddress      *ethcommon.Address
		ToETHAddress        *ethcommon.Address
		SendAmount          *big.Int
		CurrentBaseFee      *big.Int
		ChainSupportBaseFee bool
		Mode                loadTestMode
		ParsedModes         []loadTestMode
		MultiMode           bool
	}
)

var (
	//go:embed usage.md
	usage               string
	inputLoadTestParams loadTestParams
	loadTestResults     []loadTestSample
	loadTestResutsMutex sync.RWMutex

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

	randSrc *rand.Rand
)

// LoadtestCmd represents the loadtest command
var LoadtestCmd = &cobra.Command{
	Use:   "loadtest",
	Short: "Run a generic load test against an Eth/EVM style JSON-RPC endpoint.",
	Long:  usage,
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return checkFlags()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return runLoadTest(cmd.Context())
	},
}

func init() {
	ltp := new(loadTestParams)

	ltp.RPCUrl = LoadtestCmd.PersistentFlags().StringP("rpc-url", "r", "http://localhost:8545", "The RPC endpoint url")
	ltp.Requests = LoadtestCmd.PersistentFlags().Int64P("requests", "n", 1, "Number of requests to perform for the benchmarking session. The default is to just perform a single request which usually leads to non-representative benchmarking results.")
	ltp.Concurrency = LoadtestCmd.PersistentFlags().Int64P("concurrency", "c", 1, "Number of requests to perform concurrently. Default is one request at a time.")
	ltp.TimeLimit = LoadtestCmd.PersistentFlags().Int64P("time-limit", "t", -1, "Maximum number of seconds to spend for benchmarking. Use this to benchmark within a fixed total amount of time. Per default there is no time limit.")
	ltp.PrivateKey = LoadtestCmd.PersistentFlags().String("private-key", codeQualityPrivateKey, "The hex encoded private key that we'll use to send transactions")
	ltp.ChainID = LoadtestCmd.PersistentFlags().Uint64("chain-id", 0, "The chain id for the transactions.")
	ltp.ToAddress = LoadtestCmd.PersistentFlags().String("to-address", "0xDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEF", "The address that we're going to send to")
	ltp.ToRandom = LoadtestCmd.PersistentFlags().Bool("to-random", false, "When doing a transfer test, should we send to random addresses rather than DEADBEEFx5")
	ltp.CallOnly = LoadtestCmd.PersistentFlags().Bool("call-only", false, "When using this mode, rather than sending a transaction, we'll just call. This mode is incompatible with adaptive rate limiting, summarization, and a few other features.")
	ltp.CallOnlyLatestBlock = LoadtestCmd.PersistentFlags().Bool("call-only-latest", false, "When using call only mode with recall, should we execute on the latest block or on the original block")
	ltp.HexSendAmount = LoadtestCmd.PersistentFlags().String("send-amount", "0x38D7EA4C68000", "The amount of wei that we'll send every transaction")
	ltp.RateLimit = LoadtestCmd.PersistentFlags().Float64("rate-limit", 4, "An overall limit to the number of requests per second. Give a number less than zero to remove this limit all together")
	ltp.AdaptiveRateLimit = LoadtestCmd.PersistentFlags().Bool("adaptive-rate-limit", false, "Enable AIMD-style congestion control to automatically adjust request rate")
	ltp.SteadyStateTxPoolSize = LoadtestCmd.PersistentFlags().Uint64("steady-state-tx-pool-size", 1000, "When using adaptive rate limiting, this value sets the target queue size. If the queue is smaller than this value, we'll speed up. If the queue is smaller than this value, we'll back off.")
	ltp.AdaptiveRateLimitIncrement = LoadtestCmd.PersistentFlags().Uint64("adaptive-rate-limit-increment", 50, "When using adaptive rate limiting, this flag controls the size of the additive increases.")
	ltp.AdaptiveCycleDuration = LoadtestCmd.PersistentFlags().Uint64("adaptive-cycle-duration-seconds", 10, "When using adaptive rate limiting, this flag controls how often we check the queue size and adjust the rates")
	ltp.AdaptiveBackoffFactor = LoadtestCmd.PersistentFlags().Float64("adaptive-backoff-factor", 2, "When using adaptive rate limiting, this flag controls our multiplicative decrease value.")
	ltp.Modes = LoadtestCmd.PersistentFlags().StringSliceP("mode", "m", []string{"t"}, `The testing mode to use. It can be multiple like: "t,c,d,f"
t - sending transactions
d - deploy contract
c - call random contract functions
f - call specific contract function
p - call random precompiled contracts
a - call a specific precompiled contract address
s - store mode
r - random modes
2 - ERC20 Transfers
7 - ERC721 Mints
R - total recall
rpc - call random rpc methods`)
	ltp.Function = LoadtestCmd.PersistentFlags().Uint64P("function", "f", 1, "A specific function to be called if running with `--mode f` or a specific precompiled contract when running with `--mode a`")
	ltp.Iterations = LoadtestCmd.PersistentFlags().Uint64P("iterations", "i", 1, "If we're making contract calls, this controls how many times the contract will execute the instruction in a loop. If we are making ERC721 Mints, this indicates the minting batch size")
	ltp.ByteCount = LoadtestCmd.PersistentFlags().Uint64P("byte-count", "b", 1024, "If we're in store mode, this controls how many bytes we'll try to store in our contract")
	ltp.Seed = LoadtestCmd.PersistentFlags().Int64("seed", 123456, "A seed for generating random values and addresses")
	ltp.LtAddress = LoadtestCmd.PersistentFlags().String("lt-address", "", "The address of a pre-deployed load test contract")
	ltp.ERC20Address = LoadtestCmd.PersistentFlags().String("erc20-address", "", "The address of a pre-deployed erc 20 contract")
	ltp.ERC721Address = LoadtestCmd.PersistentFlags().String("erc721-address", "", "The address of a pre-deployed erc 721 contract")
	ltp.ContractCallNumberOfBlocksToWaitFor = LoadtestCmd.PersistentFlags().Uint64("contract-call-nb-blocks-to-wait-for", 30, "The number of blocks to wait for before giving up on a contract deployment")
	ltp.ContractCallBlockInterval = LoadtestCmd.PersistentFlags().Uint64("contract-call-block-interval", 1, "During deployment, this flag controls if we should check every block, every other block, or every nth block to determine that the contract has been deployed")
	ltp.ForceContractDeploy = LoadtestCmd.PersistentFlags().Bool("force-contract-deploy", false, "Some load test modes don't require a contract deployment. Set this flag to true to force contract deployments. This will still respect the --lt-address flags.")
	ltp.ForceGasLimit = LoadtestCmd.PersistentFlags().Uint64("gas-limit", 0, "In environments where the gas limit can't be computed on the fly, we can specify it manually. This can also be used to avoid eth_estimateGas")
	ltp.ForceGasPrice = LoadtestCmd.PersistentFlags().Uint64("gas-price", 0, "In environments where the gas price can't be determined automatically, we can specify it manually")
	ltp.ForcePriorityGasPrice = LoadtestCmd.PersistentFlags().Uint64("priority-gas-price", 0, "Specify Gas Tip Price in the case of EIP-1559")
	ltp.ShouldProduceSummary = LoadtestCmd.PersistentFlags().Bool("summarize", false, "Should we produce an execution summary after the load test has finished. If you're running a large load test, this can take a long time")
	ltp.BatchSize = LoadtestCmd.PersistentFlags().Uint64("batch-size", 999, "Number of batches to perform at a time for receipt fetching. Default is 999 requests at a time.")
	ltp.SummaryOutputMode = LoadtestCmd.PersistentFlags().String("output-mode", "text", "Format mode for summary output (json | text)")
	ltp.LegacyTransactionMode = LoadtestCmd.PersistentFlags().Bool("legacy", false, "Send a legacy transaction instead of an EIP1559 transaction.")
	ltp.RecallLength = LoadtestCmd.PersistentFlags().Uint64("recall-blocks", 50, "The number of blocks that we'll attempt to fetch for recall")
	inputLoadTestParams = *ltp

	// TODO Compression
}

func checkFlags() error {
	ltp := inputLoadTestParams

	// Check `rpc-url` flag.
	if ltp.RPCUrl == nil {
		panic("RPC URL is empty")
	}
	if err := util.ValidateUrl(*ltp.RPCUrl); err != nil {
		return err
	}

	if ltp.AdaptiveBackoffFactor != nil && *ltp.AdaptiveBackoffFactor <= 0.0 {
		return fmt.Errorf("the backoff factor needs to be non-zero positive")
	}

	if ltp.ContractCallBlockInterval != nil && *ltp.ContractCallBlockInterval == 0 {
		return fmt.Errorf("the contract call block interval must be strictly positive")
	}

	return nil
}
