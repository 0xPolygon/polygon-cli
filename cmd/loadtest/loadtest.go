/*
Copyright © 2022 Polygon <engineering@polygon.technology>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package loadtest

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"math/rand"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/maticnetwork/polygon-cli/metrics"
	"github.com/maticnetwork/polygon-cli/rpctypes"
	"github.com/maticnetwork/polygon-cli/util"
	"golang.org/x/exp/constraints"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	gssignature "github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	gstypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"

	"github.com/maticnetwork/polygon-cli/contracts"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/time/rate"
)

const (
	loadTestModeTransaction          = "t"
	loadTestModeDeploy               = "d"
	loadTestModeCall                 = "c"
	loadTestModeFunction             = "f"
	loadTestModeInc                  = "i"
	loadTestModeRandom               = "r"
	loadTestModeStore                = "s"
	loadTestModeLong                 = "l"
	loadTestModeERC20                = "2"
	loadTestModeERC721               = "7"
	loadTestModePrecompiledContracts = "p"
	loadTestModePrecompiledContract  = "a"

	codeQualitySeed       = "code code code code code code code code code code code quality"
	codeQualityPrivateKey = "42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa"
)

var (
	inputLoadTestParams loadTestParams
	loadTestResults     []loadTestSample
	loadTestResutsMutex sync.RWMutex
	validLoadTestModes  = []string{
		loadTestModeTransaction,
		loadTestModeDeploy,
		loadTestModeCall,
		loadTestModeFunction,
		loadTestModeInc,
		loadTestModeStore,
		loadTestModeLong,
		loadTestModeERC20,
		loadTestModeERC721,
		loadTestModePrecompiledContracts,
		loadTestModePrecompiledContract,
		// r should be last to exclude it from random mode selection
		loadTestModeRandom,
	}

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
)

// LoadtestCmd represents the loadtest command
var LoadtestCmd = &cobra.Command{
	Use:   "loadtest rpc-endpoint",
	Short: "A simple script for quickly running a load test",
	Long:  `Loadtest gives us a simple way to run a generic load test against an eth/EVM style json RPC endpoint`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Debug().Msg("Starting Loadtest")

		err := runLoadTest(cmd.Context())
		if err != nil {
			return err
		}
		return nil
	},
	Args: func(cmd *cobra.Command, args []string) error {
		zerolog.DurationFieldUnit = time.Second
		zerolog.DurationFieldInteger = true

		if len(args) != 1 {
			return fmt.Errorf("expected exactly one argument")
		}
		url, err := url.Parse(args[0])
		if err != nil {
			log.Error().Err(err).Msg("Unable to parse url input error")
			return err
		}
		if url.Scheme != "http" && url.Scheme != "https" && url.Scheme != "ws" && url.Scheme != "wss" {
			return fmt.Errorf("the scheme %s is not supported", url.Scheme)
		}
		inputLoadTestParams.URL = url
		r := regexp.MustCompile(fmt.Sprintf("^[%s]+$", strings.Join(validLoadTestModes, "")))
		if !r.MatchString(*inputLoadTestParams.Mode) {
			return fmt.Errorf("the mode %s is not recognized", *inputLoadTestParams.Mode)
		}
		if *inputLoadTestParams.AdaptiveBackoffFactor <= 0.0 {
			return fmt.Errorf("the backoff factor needs to be non-zero positive")
		}
		return nil
	},
}

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
		Requests                            *int64
		Concurrency                         *int64
		BatchSize                           *uint64
		TimeLimit                           *int64
		ToRandom                            *bool
		URL                                 *url.URL
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
		Mode                                *string
		Function                            *uint64
		Iterations                          *uint64
		ByteCount                           *uint64
		Seed                                *int64
		IsAvail                             *bool
		LtAddress                           *string
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

		// Computed
		CurrentGas       *big.Int
		CurrentGasTipCap *big.Int
		CurrentNonce     *uint64
		ECDSAPrivateKey  *ecdsa.PrivateKey
		FromETHAddress   *ethcommon.Address
		ToETHAddress     *ethcommon.Address
		SendAmount       *big.Int
		BaseFee          *big.Int

		ToAvailAddress   *gstypes.MultiAddress
		FromAvailAddress *gssignature.KeyringPair
		AvailRuntime     *gstypes.RuntimeVersion
	}

	txpoolStatus struct {
		Pending any `json:"pending"`
		Queued  any `json:"queued"`
	}
)

func init() {
	ltp := new(loadTestParams)
	// Apache Bench Parameters
	ltp.Requests = LoadtestCmd.PersistentFlags().Int64P("requests", "n", 1, "Number of requests to perform for the benchmarking session. The default is to just perform a single request which usually leads to non-representative benchmarking results.")
	ltp.Concurrency = LoadtestCmd.PersistentFlags().Int64P("concurrency", "c", 1, "Number of multiple requests to perform at a time. Default is one request at a time.")
	ltp.TimeLimit = LoadtestCmd.PersistentFlags().Int64P("time-limit", "t", -1, "Maximum number of seconds to spend for benchmarking. Use this to benchmark within a fixed total amount of time. Per default there is no timelimit.")
	// https://logging.apache.org/log4j/2.x/manual/customloglevels.html

	// extended parameters
	ltp.PrivateKey = LoadtestCmd.PersistentFlags().String("private-key", codeQualityPrivateKey, "The hex encoded private key that we'll use to sending transactions")
	ltp.ChainID = LoadtestCmd.PersistentFlags().Uint64P("chain-id", "", 0, "The chain id for the transactions that we're going to send")
	ltp.ToAddress = LoadtestCmd.PersistentFlags().String("to-address", "0xDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEF", "The address that we're going to send to")
	ltp.ToRandom = LoadtestCmd.PersistentFlags().Bool("to-random", false, "When doing a transfer test, should we send to random addresses rather than DEADBEEFx5")
	ltp.HexSendAmount = LoadtestCmd.PersistentFlags().String("send-amount", "0x38D7EA4C68000", "The amount of wei that we'll send every transaction")
	ltp.RateLimit = LoadtestCmd.PersistentFlags().Float64("rate-limit", 4, "An overall limit to the number of requests per second. Give a number less than zero to remove this limit all together")
	ltp.AdaptiveRateLimit = LoadtestCmd.PersistentFlags().Bool("adaptive-rate-limit", false, "Loadtest automatically adjusts request rate to maximize utilization but prevent congestion")
	ltp.SteadyStateTxPoolSize = LoadtestCmd.PersistentFlags().Uint64("steady-state-tx-pool-size", 1000, "Transaction Pool queue size which we use to either increase/decrease requests per second")
	ltp.AdaptiveRateLimitIncrement = LoadtestCmd.PersistentFlags().Uint64("adaptive-rate-limit-increment", 50, "Additive increment to rate of requests if txpool below steady state size")
	ltp.AdaptiveCycleDuration = LoadtestCmd.PersistentFlags().Uint64("adaptive-cycle-duration-seconds", 10, "Duration in seconds that adaptive load test will review txpool and determine whether to increase/decrease rate limit")
	ltp.AdaptiveBackoffFactor = LoadtestCmd.PersistentFlags().Float64("adaptive-backoff-factor", 2, "When we detect congestion we will use this factor to determine how much we slow down")
	ltp.Mode = LoadtestCmd.PersistentFlags().StringP("mode", "m", "t", `The testing mode to use. It can be multiple like: "tcdf"
t - sending transactions
d - deploy contract
c - call random contract functions
f - call specific contract function
p - call random precompiled contracts
a - call a specific precompiled contract address
s - store mode
l - long running mode
r - random modes
2 - ERC20 Transfers
7 - ERC721 Mints`)
	ltp.Function = LoadtestCmd.PersistentFlags().Uint64P("function", "f", 1, "A specific function to be called if running with `--mode f` or a specific precompiled contract when running with `--mode a`")
	ltp.Iterations = LoadtestCmd.PersistentFlags().Uint64P("iterations", "i", 100, "If we're making contract calls, this controls how many times the contract will execute the instruction in a loop. If we are making ERC721 Mints, this indicated the minting batch size")
	ltp.ByteCount = LoadtestCmd.PersistentFlags().Uint64P("byte-count", "b", 1024, "If we're in store mode, this controls how many bytes we'll try to store in our contract")
	ltp.Seed = LoadtestCmd.PersistentFlags().Int64("seed", 123456, "A seed for generating random values and addresses")
	ltp.IsAvail = LoadtestCmd.PersistentFlags().Bool("data-avail", false, "Is this a test of avail rather than an EVM / Geth Chain")
	ltp.LtAddress = LoadtestCmd.PersistentFlags().String("lt-address", "", "A pre-deployed load test contract address")
	ltp.DelAddress = LoadtestCmd.PersistentFlags().String("del-address", "", "A pre-deployed delegator contract address")
	ltp.ContractCallNumberOfBlocksToWaitFor = LoadtestCmd.PersistentFlags().Uint64("contract-call-nb-blocks-to-wait-for", 30, "The number of blocks to wait for before giving up on a contract call")
	ltp.ContractCallBlockInterval = LoadtestCmd.PersistentFlags().Uint64("contract-call-block-interval", 1, "The number of blocks to wait between contract calls")
	ltp.ForceContractDeploy = LoadtestCmd.PersistentFlags().Bool("force-contract-deploy", false, "Some loadtest modes don't require a contract deployment. Set this flag to true to force contract deployments. This will still respect the --del-address and --il-address flags.")
	ltp.ForceGasLimit = LoadtestCmd.PersistentFlags().Uint64("gas-limit", 0, "In environments where the gas limit can't be computed on the fly, we can specify it manually")
	ltp.ForceGasPrice = LoadtestCmd.PersistentFlags().Uint64("gas-price", 0, "In environments where the gas price can't be estimated, we can specify it manually")
	ltp.ForcePriorityGasPrice = LoadtestCmd.PersistentFlags().Uint64("priority-gas-price", 0, "Specify Gas Tip Price in the case of EIP-1559")
	ltp.ShouldProduceSummary = LoadtestCmd.PersistentFlags().Bool("summarize", false, "Should we produce an execution summary after the load test has finished. If you're running a large loadtest, this can take a long time")
	ltp.BatchSize = LoadtestCmd.PersistentFlags().Uint64("batch-size", 999, "Number of batches to perform at a time for receipt fetching. Default is 999 requests at a time.")
	ltp.SummaryOutputMode = LoadtestCmd.PersistentFlags().String("output-mode", "text", "Format mode for summary output (json | text)")
	ltp.LegacyTransactionMode = LoadtestCmd.PersistentFlags().Bool("legacy", false, "Send a legacy transaction instead of an EIP1559 transaction.")
	inputLoadTestParams = *ltp

	// TODO batch size
	// TODO Compression
	// TODO array of RPC endpoints to round robin?
}

func initializeLoadTestParams(ctx context.Context, c *ethclient.Client) error {
	log.Info().Msg("Connecting with RPC endpoint to initialize load test parameters")
	gas, err := c.SuggestGasPrice(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Unable to retrieve gas price")
		return err
	}
	log.Trace().Interface("gasprice", gas).Msg("Retreived current gas price")

	if !*inputLoadTestParams.LegacyTransactionMode {
		gasTipCap, _err := c.SuggestGasTipCap(ctx)
		if _err != nil {
			log.Error().Err(_err).Msg("Unable to retrieve gas tip cap")
			return _err
		}
		log.Trace().Interface("gastipcap", gasTipCap).Msg("Retreived current gas tip cap")
		inputLoadTestParams.CurrentGasTipCap = gasTipCap
	}

	privateKey, err := ethcrypto.HexToECDSA(*inputLoadTestParams.PrivateKey)
	if err != nil {
		log.Error().Err(err).Msg("Couldn't process the hex private key")
		return err
	}

	blockNumber, err := c.BlockNumber(ctx)
	bigBlockNumber := big.NewInt(int64(blockNumber))
	if err != nil {
		log.Error().Err(err).Msg("Couldn't get the current block number")
		return err
	}
	log.Trace().Uint64("blocknumber", blockNumber).Msg("Current Block Number")

	ethAddress := ethcrypto.PubkeyToAddress(privateKey.PublicKey)

	nonce, err := c.NonceAt(ctx, ethAddress, bigBlockNumber)
	if err != nil {
		log.Error().Err(err).Msg("Unable to get account nonce")
		return err
	}
	accountBal, err := c.BalanceAt(ctx, ethAddress, bigBlockNumber)
	if err != nil {
		log.Error().Err(err).Msg("Unable to get the balance for the account")
		return err
	}
	log.Trace().Interface("balance", accountBal).Msg("Current account balance")

	toAddr := ethcommon.HexToAddress(*inputLoadTestParams.ToAddress)

	amt, err := hexToBigInt(*inputLoadTestParams.HexSendAmount)
	if err != nil {
		log.Error().Err(err).Msg("Couldn't parse send amount")
		return err
	}

	header, err := c.HeaderByNumber(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg("Unable to get header")
		return err
	}

	chainID, err := c.ChainID(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Unable to fetch chain ID")
		return err
	}

	if *inputLoadTestParams.LegacyTransactionMode && *inputLoadTestParams.ForcePriorityGasPrice > 0 {
		log.Error().Msg("Cannot set priority gas price in legacy mode")
		return errors.New("cannot set priority gas price in legacy mode")
	}

	inputLoadTestParams.ToETHAddress = &toAddr
	inputLoadTestParams.SendAmount = amt
	inputLoadTestParams.CurrentGas = gas
	inputLoadTestParams.CurrentNonce = &nonce
	inputLoadTestParams.ECDSAPrivateKey = privateKey
	inputLoadTestParams.FromETHAddress = &ethAddress
	if *inputLoadTestParams.ChainID == 0 {
		*inputLoadTestParams.ChainID = chainID.Uint64()
	}
	inputLoadTestParams.BaseFee = header.BaseFee

	rand.Seed(*inputLoadTestParams.Seed)

	return nil
}

func hexToBigInt(raw any) (bi *big.Int, err error) {
	bi = big.NewInt(0)
	hexString, ok := raw.(string)
	if !ok {
		err = fmt.Errorf("could not assert value %v as a string", raw)
		return
	}
	hexString = strings.Replace(hexString, "0x", "", -1)
	if len(hexString)%2 != 0 {
		log.Trace().Str("original", hexString).Msg("Hex of odd length")
		hexString = "0" + hexString
	}

	rawGas, err := hex.DecodeString(hexString)
	if err != nil {
		log.Error().Err(err).Str("hex", hexString).Msg("Unable to decode hex string")
		return
	}
	bi.SetBytes(rawGas)
	return
}

func runLoadTest(ctx context.Context) error {
	log.Info().Msg("Starting Load Test")

	timeLimit := *inputLoadTestParams.TimeLimit
	var overallTimer *time.Timer
	if timeLimit > 0 {
		overallTimer = time.NewTimer(time.Duration(timeLimit) * time.Second)
	} else {
		overallTimer = new(time.Timer)
	}

	rpc, err := ethrpc.DialContext(ctx, inputLoadTestParams.URL.String())
	if err != nil {
		log.Error().Err(err).Msg("Unable to dial rpc")
		return err
	}
	rpc.SetHeader("Accept-Encoding", "identity")
	ec := ethclient.NewClient(rpc)

	var loopFunc func() error
	if *inputLoadTestParams.IsAvail {
		log.Info().Msg("Running in Avail mode")
		loopFunc = func() error {
			var api *gsrpc.SubstrateAPI
			api, err = gsrpc.NewSubstrateAPI(inputLoadTestParams.URL.String())
			if err != nil {
				return err
			}
			err = initAvailTestParams(ctx, api)
			if err != nil {
				return err
			}
			return availLoop(ctx, api)
		}

	} else {
		log.Info().Msg("Starting Load Test")
		loopFunc = func() error {
			err = initializeLoadTestParams(ctx, ec)
			if err != nil {
				return err
			}

			return mainLoop(ctx, ec, rpc)
		}
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	loadTestResults = make([]loadTestSample, 0)
	errCh := make(chan error)
	go func() {
		errCh <- loopFunc()
	}()

	select {
	case <-overallTimer.C:
		log.Info().Msg("Time's up")
	case <-sigCh:
		log.Info().Msg("Interrupted.. Stopping load test")
	case err = <-errCh:
		if err != nil {
			log.Fatal().Err(err).Msg("Received critical error while running load test")
		}
	}

	printResults(loadTestResults)
	if *inputLoadTestParams.IsAvail {
		log.Trace().Msg("Finished testing avail")
		return nil
	}

	// TODO this doesn't make sense for avail
	ptc, err := ec.PendingTransactionCount(ctx)
	if err != nil {
		log.Debug().Err(err).Msg("Unable to get the number of pending transactions before closing")
	} else if ptc > 0 {
		log.Info().Uint("pending", ptc).Msg("There are still oustanding transactions. There might be issues restarting with the same sending key until those transactions clear")
	}
	log.Info().Msg("Finished")
	return nil
}

func printResults(lts []loadTestSample) {
	if len(lts) == 0 {
		log.Error().Msg("No results recorded")
		return
	}

	log.Info().Msg("* Results")
	log.Info().Int("samples", len(lts)).Msg("Samples")

	var startTime = lts[0].RequestTime
	var endTime = lts[len(lts)-1].RequestTime
	var meanWait float64
	var totalWait float64 = 0
	var numErrors uint64 = 0

	for _, s := range lts {
		if s.IsError {
			numErrors += 1
		}
		totalWait = float64(s.WaitTime.Seconds()) + totalWait
	}
	meanWait = totalWait / float64(len(lts))

	log.Info().Time("startTime", startTime).Msg("Start")
	log.Info().Time("endTime", endTime).Msg("End")
	log.Info().Float64("meanWait", meanWait).Msg("Mean Wait")
	log.Info().Uint64("numErrors", numErrors).Msg("Num errors")
}

func convHexToUint64(hexString string) (uint64, error) {
	hexString = strings.TrimPrefix(hexString, "0x")
	if len(hexString)%2 != 0 {
		hexString = "0" + hexString
	}

	result, err := strconv.ParseUint(hexString, 16, 64)
	if err != nil {
		return 0, err
	}
	return uint64(result), nil
}

func tryCastToUint64(val any) (uint64, error) {
	switch t := val.(type) {
	case float64:
		return uint64(t), nil
	case string:
		return convHexToUint64(t)
	default:
		return 0, fmt.Errorf("the value %v couldn't be marshalled to uint64", t)

	}
}

func getTxPoolSize(rpc *ethrpc.Client) (uint64, error) {
	var status = new(txpoolStatus)
	err := rpc.Call(status, "txpool_status")
	if err != nil {
		return 0, err
	}
	pendingCount, err := tryCastToUint64(status.Pending)
	if err != nil {
		return 0, err
	}
	queuedCount, err := tryCastToUint64(status.Queued)
	if err != nil {
		return 0, err
	}

	return pendingCount + queuedCount, nil
}

func updateRateLimit(ctx context.Context, rl *rate.Limiter, rpc *ethrpc.Client, steadyStateQueueSize uint64, rateLimitIncrement uint64, cycleDuration time.Duration, backoff float64) {
	ticker := time.NewTicker(cycleDuration)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			txPoolSize, err := getTxPoolSize(rpc)
			if err != nil {
				log.Error().Err(err).Msg("Error getting txpool size")
				return
			}

			if txPoolSize < steadyStateQueueSize {
				// additively increment requests per second if txpool less than queue steady state
				newRateLimit := rate.Limit(float64(rl.Limit()) + float64(rateLimitIncrement))
				rl.SetLimit(newRateLimit)
				log.Info().Float64("New Rate Limit (RPS)", float64(rl.Limit())).Uint64("Current Tx Pool Size", txPoolSize).Uint64("Steady State Tx Pool Size", steadyStateQueueSize).Msg("Increased rate limit")
			} else if txPoolSize > steadyStateQueueSize {
				// halve rate limit requests per second if txpool greater than queue steady state
				rl.SetLimit(rl.Limit() / rate.Limit(backoff))
				log.Info().Float64("New Rate Limit (RPS)", float64(rl.Limit())).Uint64("Current Tx Pool Size", txPoolSize).Uint64("Steady State Tx Pool Size", steadyStateQueueSize).Msg("Backed off rate limit")
			}
		case <-ctx.Done():
			return
		}
	}
}

func mainLoop(ctx context.Context, c *ethclient.Client, rpc *ethrpc.Client) error {
	ltp := inputLoadTestParams
	log.Trace().Interface("Input Params", ltp).Msg("Params")

	routines := *ltp.Concurrency
	requests := *ltp.Requests
	currentNonce := *ltp.CurrentNonce
	chainID := new(big.Int).SetUint64(*ltp.ChainID)
	privateKey := ltp.ECDSAPrivateKey
	mode := *ltp.Mode
	steadyStateTxPoolSize := *ltp.SteadyStateTxPoolSize
	adaptiveRateLimitIncrement := *ltp.AdaptiveRateLimitIncrement
	var rl *rate.Limiter
	rl = rate.NewLimiter(rate.Limit(*ltp.RateLimit), 1)
	if *ltp.RateLimit <= 0.0 {
		rl = nil
	}
	rateLimitCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	if *ltp.AdaptiveRateLimit && rl != nil {
		go updateRateLimit(rateLimitCtx, rl, rpc, steadyStateTxPoolSize, adaptiveRateLimitIncrement, time.Duration(*ltp.AdaptiveCycleDuration)*time.Second, *ltp.AdaptiveBackoffFactor)
	}

	tops, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	tops = configureTransactOpts(tops)
	tops.GasLimit = 10000000

	if err != nil {
		log.Error().Err(err).Msg("Unable create transaction signer")
		return err
	}
	cops := new(bind.CallOpts)

	// deploy and instantiate the load tester contract
	var ltAddr ethcommon.Address
	var ltContract *contracts.LoadTester
	numberOfBlocksToWaitFor := *inputLoadTestParams.ContractCallNumberOfBlocksToWaitFor
	blockInterval := *inputLoadTestParams.ContractCallBlockInterval
	if strings.ContainsAny(mode, "rcfislpas") || *inputLoadTestParams.ForceContractDeploy {
		if *inputLoadTestParams.LtAddress == "" {
			ltAddr, _, _, err = contracts.DeployLoadTester(tops, c)
			if err != nil {
				log.Error().Err(err).Msg("Failed to create the load testing contract. Do you have the right chain id? Do you have enough funds?")
				return err
			}
		} else {
			ltAddr = ethcommon.HexToAddress(*inputLoadTestParams.LtAddress)
		}
		log.Trace().Interface("contractaddress", ltAddr).Msg("Load test contract address")
		// bump the nonce since deploying a contract should cause it to increase
		currentNonce = currentNonce + 1

		ltContract, err = contracts.NewLoadTester(ltAddr, c)
		if err != nil {
			log.Error().Err(err).Msg("Unable to instantiate new contract")
			return err
		}
		err = blockUntilSuccessful(ctx, c, func() error {
			_, err = ltContract.GetCallCounter(cops)
			return err
		}, numberOfBlocksToWaitFor, blockInterval)

		if err != nil {
			return err
		}
	}

	var erc20Addr ethcommon.Address
	var erc20Contract *contracts.ERC20
	if mode == loadTestModeERC20 || mode == loadTestModeRandom {
		erc20Addr, _, _, err = contracts.DeployERC20(tops, c)
		if err != nil {
			log.Error().Err(err).Msg("Unable to deploy ERC20 contract")
			return err
		}
		log.Trace().Interface("contractaddress", erc20Addr).Msg("ERC20 contract address")

		erc20Contract, err = contracts.NewERC20(erc20Addr, c)
		if err != nil {
			log.Error().Err(err).Msg("Unable to instantiate new erc20 contract")
			return err
		}
		currentNonce = currentNonce + 1
		err = blockUntilSuccessful(ctx, c, func() error {
			_, err = erc20Contract.BalanceOf(cops, *ltp.FromETHAddress)
			return err
		}, numberOfBlocksToWaitFor, blockInterval)
		if err != nil {
			return err
		}

		tops.Nonce = new(big.Int).SetUint64(currentNonce)

		_, err = erc20Contract.Mint(tops, metrics.UnitMegaether)
		if err != nil {
			log.Error().Err(err).Msg("There was an error minting ERC20")
			return err
		}

		currentNonce = currentNonce + 1
		err = blockUntilSuccessful(ctx, c, func() error {
			var balance *big.Int
			balance, err = erc20Contract.BalanceOf(cops, *ltp.FromETHAddress)
			if err != nil {
				return err
			}
			if balance.Uint64() == 0 {
				err = fmt.Errorf("ERC20 Balance is Zero")
				return err
			}
			return nil
		}, numberOfBlocksToWaitFor, blockInterval)
		if err != nil {
			return err
		}
	}

	var erc721Addr ethcommon.Address
	var erc721Contract *contracts.ERC721
	if mode == loadTestModeERC721 || mode == loadTestModeRandom {
		erc721Addr, _, _, err = contracts.DeployERC721(tops, c)
		if err != nil {
			log.Error().Err(err).Msg("Unable to deploy ERC721 contract")
			return err
		}
		log.Trace().Interface("contractaddress", erc721Addr).Msg("ERC721 contract address")

		erc721Contract, err = contracts.NewERC721(erc721Addr, c)
		if err != nil {
			log.Error().Err(err).Msg("Unable to instantiate new erc20 contract")
			return err
		}
		currentNonce = currentNonce + 1

		err = blockUntilSuccessful(ctx, c, func() error {
			_, err = erc721Contract.BalanceOf(cops, *ltp.FromETHAddress)
			return err
		}, numberOfBlocksToWaitFor, blockInterval)
		if err != nil {
			return err
		}

		tops.Nonce = new(big.Int).SetUint64(currentNonce)

		err = blockUntilSuccessful(ctx, c, func() error {
			_, err = erc721Contract.MintBatch(tops, *ltp.FromETHAddress, new(big.Int).SetUint64(1))
			return err
		}, numberOfBlocksToWaitFor, blockInterval)
		if err != nil {
			return err
		}
		currentNonce = currentNonce + 1
	}

	// deploy and instantiate the delegator contract
	var delegatorContract *contracts.Delegator
	if strings.ContainsAny(mode, "rl") || *inputLoadTestParams.ForceContractDeploy {
		var delegatorAddr ethcommon.Address
		if *inputLoadTestParams.DelAddress == "" {
			delegatorAddr, _, _, err = contracts.DeployDelegator(tops, c)
			if err != nil {
				log.Error().Err(err).Msg("Failed to create the load testing contract. Do you have the right chain id? Do you have enough funds?")
				return err
			}
		} else {
			delegatorAddr = ethcommon.HexToAddress(*inputLoadTestParams.DelAddress)
		}
		log.Trace().Interface("contractaddress", delegatorAddr).Msg("Delegator contract address")
		currentNonce = currentNonce + 1

		delegatorContract, err = contracts.NewDelegator(delegatorAddr, c)
		if err != nil {
			log.Error().Err(err).Msg("Unable to instantiate new contract")
			return err
		}

		err = blockUntilSuccessful(ctx, c, func() error {
			_, err = delegatorContract.Call(tops, ltAddr, []byte{0x12, 0x87, 0xa6, 0x8c})
			return err
		}, numberOfBlocksToWaitFor, blockInterval)
		if err != nil {
			return err
		}
		currentNonce = currentNonce + 1
	}

	var currentNonceMutex sync.Mutex
	var i int64
	startBlockNumber, err := c.BlockNumber(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get current block number")
		return err
	}
	startNonce := currentNonce
	log.Debug().Uint64("currentNonce", currentNonce).Msg("Starting main loadtest loop")
	var wg sync.WaitGroup
	for i = 0; i < routines; i = i + 1 {
		log.Trace().Int64("routine", i).Msg("Starting Thread")
		wg.Add(1)
		go func(i int64) {
			var j int64
			var startReq time.Time
			var endReq time.Time
			var retryForNonce bool = false
			var myNonceValue uint64

			for j = 0; j < requests; j = j + 1 {
				if rl != nil {
					err = rl.Wait(ctx)
					if err != nil {
						log.Error().Err(err).Msg("Encountered a rate limiting error")
					}
				}

				if retryForNonce {
					retryForNonce = false
				} else {
					currentNonceMutex.Lock()
					myNonceValue = currentNonce
					currentNonce = currentNonce + 1
					currentNonceMutex.Unlock()
				}

				localMode := mode
				// if there are multiple modes, iterate through them, 'r' mode is supported here
				if len(mode) > 1 {
					localMode = string(mode[int(i+j)%(len(mode))])
				}
				// if we're doing random, we'll just pick one based on the current index
				if localMode == loadTestModeRandom {
					localMode = validLoadTestModes[int(i+j)%(len(validLoadTestModes)-1)]
				}
				switch localMode {
				case loadTestModeTransaction:
					startReq, endReq, err = loadtestTransaction(ctx, c, myNonceValue)
				case loadTestModeDeploy:
					startReq, endReq, err = loadtestDeploy(ctx, c, myNonceValue)
				case loadTestModeCall:
					startReq, endReq, err = loadtestCall(ctx, c, myNonceValue, ltContract)
				case loadTestModeFunction:
					startReq, endReq, err = loadtestFunction(ctx, c, myNonceValue, ltContract)
				case loadTestModeInc:
					startReq, endReq, err = loadtestInc(ctx, c, myNonceValue, ltContract)
				case loadTestModeStore:
					startReq, endReq, err = loadtestStore(ctx, c, myNonceValue, ltContract)
				case loadTestModeLong:
					startReq, endReq, err = loadtestLong(ctx, c, myNonceValue, delegatorContract, ltAddr)
				case loadTestModeERC20:
					startReq, endReq, err = loadtestERC20(ctx, c, myNonceValue, erc20Contract, ltAddr)
				case loadTestModeERC721:
					startReq, endReq, err = loadtestERC721(ctx, c, myNonceValue, erc721Contract, ltAddr)
				case loadTestModePrecompiledContract:
					startReq, endReq, err = loadtestCallPrecompiledContracts(ctx, c, myNonceValue, ltContract, true)
				case loadTestModePrecompiledContracts:
					startReq, endReq, err = loadtestCallPrecompiledContracts(ctx, c, myNonceValue, ltContract, false)
				default:
					log.Error().Str("mode", mode).Msg("We've arrived at a load test mode that we don't recognize")
				}
				recordSample(i, j, err, startReq, endReq, myNonceValue)
				if err != nil {
					log.Error().Err(err).Uint64("nonce", myNonceValue).Msg("Recorded an error while sending transactions")
					retryForNonce = true
				}

				log.Trace().Uint64("nonce", myNonceValue).Int64("routine", i).Str("mode", localMode).Int64("request", j).Msg("Request")
			}
			wg.Done()
		}(i)
	}
	log.Trace().Msg("Finished starting go routines. Waiting..")
	wg.Wait()
	cancel()
	log.Debug().Uint64("currentNonce", currentNonce).Msg("Finished main loadtest loop")
	log.Debug().Msg("Waiting for transactions to actually be mined")
	finalBlockNumber, err := waitForFinalBlock(ctx, c, rpc, startBlockNumber, startNonce, currentNonce)
	if err != nil {
		log.Error().Err(err).Msg("there was an issue waiting for all transactions to be mined")
	}

	lightSummary(ctx, c, rpc, startBlockNumber, startNonce, finalBlockNumber, currentNonce, rl)
	if *ltp.ShouldProduceSummary {
		err = summarizeTransactions(ctx, c, rpc, startBlockNumber, startNonce, finalBlockNumber, currentNonce)
		if err != nil {
			log.Error().Err(err).Msg("There was an issue creating the load test summary")
		}
	}
	return nil
}

func lightSummary(ctx context.Context, c *ethclient.Client, rpc *ethrpc.Client, startBlockNumber, startNonce, endBlockNumber, endNonce uint64, rl *rate.Limiter) {
	startBlock, err := c.BlockByNumber(ctx, new(big.Int).SetUint64(startBlockNumber))
	if err != nil {
		log.Error().Err(err).Msg("unable to get start block for light summary")
		return
	}
	endBlock, err := c.BlockByNumber(ctx, new(big.Int).SetUint64(endBlockNumber))
	if err != nil {
		log.Error().Err(err).Msg("unable to get end block for light summary")
		return
	}
	endTime := time.Unix(int64(endBlock.Time()), 0)
	startTime := time.Unix(int64(startBlock.Time()), 0)

	testDuration := endTime.Sub(startTime)
	tps := float64(len(loadTestResults)) / testDuration.Seconds()

	log.Info().
		Time("firstBlockTime", startTime).
		Time("lastBlockTime", endTime).
		Int("transactionCount", len(loadTestResults)).
		Float64("testDuration", testDuration.Seconds()).
		Float64("tps", tps).
		Float64("final rate limit", float64(rl.Limit())).
		Msg("rough test summary (ignores errors)")
}

func blockUntilSuccessful(ctx context.Context, c *ethclient.Client, f func() error, numberOfBlocksToWaitFor, blockInterval uint64) error {
	start := time.Now()
	startBlockNumber, err := c.BlockNumber(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error getting block number")
		return err
	}
	log.Trace().
		Uint64("startBlockNumber", startBlockNumber).
		Uint64("numberOfBlocksToWaitFor", numberOfBlocksToWaitFor).
		Uint64("blockInterval", blockInterval).
		Msg("Starting blocking loop")
	var lastBlockNumber, currentBlockNumber uint64
	var lock bool
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			elapsed := time.Since(start)
			blockDiff := currentBlockNumber % startBlockNumber
			if blockDiff > numberOfBlocksToWaitFor {
				log.Error().Err(err).Dur("elapsedTimeSeconds", elapsed).Msg("Exhausted waiting period")
				return err
			}

			currentBlockNumber, err = c.BlockNumber(ctx)
			if err != nil {
				log.Error().Err(err).Msg("Error getting block number")
				return err
			} else {
				log.Trace().Uint64("newBlock", currentBlockNumber).Msg("New block")
			}

			if currentBlockNumber != lastBlockNumber {
				lock = false
			}
			if (currentBlockNumber%startBlockNumber)%blockInterval == 0 {
				if !lock {
					lock = true
					err := f()
					if err == nil {
						log.Trace().Err(err).Dur("elapsedTimeSeconds", elapsed).Msg("Function executed successfuly")
						return nil
					}
					log.Trace().Err(err).Dur("elapsedTimeSeconds", elapsed).Msg("Unable to execute function")
				}
			}
			lastBlockNumber = currentBlockNumber
			time.Sleep(time.Second)
		}
	}
}

func loadtestTransaction(ctx context.Context, c *ethclient.Client, nonce uint64) (t1 time.Time, t2 time.Time, err error) {
	ltp := inputLoadTestParams

	gasPrice := ltp.CurrentGas

	to := ltp.ToETHAddress
	if *ltp.ToRandom {
		to = getRandomAddress()
	}

	amount := ltp.SendAmount
	chainID := new(big.Int).SetUint64(*ltp.ChainID)
	privateKey := ltp.ECDSAPrivateKey

	tops, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Error().Err(err).Msg("Unable create transaction signer")
		return
	}
	tops.GasLimit = uint64(21000)
	tops = configureTransactOpts(tops)

	var tx *ethtypes.Transaction
	if *ltp.LegacyTransactionMode {
		tx = ethtypes.NewTransaction(nonce, *to, amount, tops.GasLimit, gasPrice, nil)
	} else {
		gasTipCap := tops.GasTipCap
		gasFeeCap := new(big.Int).Add(gasTipCap, ltp.BaseFee)
		dynamicFeeTx := &ethtypes.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     nonce,
			To:        to,
			Gas:       tops.GasLimit,
			GasFeeCap: gasFeeCap,
			GasTipCap: gasTipCap,
			Data:      nil,
			Value:     amount,
		}
		tx = ethtypes.NewTx(dynamicFeeTx)
	}

	stx, err := tops.Signer(*ltp.FromETHAddress, tx)
	if err != nil {
		log.Error().Err(err).Msg("Unable to sign transaction")
		return
	}

	t1 = time.Now()
	err = c.SendTransaction(ctx, stx)
	t2 = time.Now()
	return
}

func loadtestDeploy(ctx context.Context, c *ethclient.Client, nonce uint64) (t1 time.Time, t2 time.Time, err error) {
	ltp := inputLoadTestParams

	chainID := new(big.Int).SetUint64(*ltp.ChainID)
	privateKey := ltp.ECDSAPrivateKey

	tops, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Error().Err(err).Msg("Unable create transaction signer")
		return
	}
	tops.Nonce = new(big.Int).SetUint64(nonce)
	tops = configureTransactOpts(tops)

	t1 = time.Now()
	_, _, _, err = contracts.DeployLoadTester(tops, c)
	t2 = time.Now()
	return
}

func loadtestFunction(ctx context.Context, c *ethclient.Client, nonce uint64, ltContract *contracts.LoadTester) (t1 time.Time, t2 time.Time, err error) {
	ltp := inputLoadTestParams

	chainID := new(big.Int).SetUint64(*ltp.ChainID)
	privateKey := ltp.ECDSAPrivateKey
	iterations := ltp.Iterations
	f := ltp.Function

	tops, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Error().Err(err).Msg("Unable create transaction signer")
		return
	}
	tops.Nonce = new(big.Int).SetUint64(nonce)
	tops = configureTransactOpts(tops)

	t1 = time.Now()
	_, err = contracts.CallLoadTestFunctionByOpCode(*f, ltContract, tops, *iterations)
	t2 = time.Now()
	return
}

func loadtestCall(ctx context.Context, c *ethclient.Client, nonce uint64, ltContract *contracts.LoadTester) (t1 time.Time, t2 time.Time, err error) {
	ltp := inputLoadTestParams

	chainID := new(big.Int).SetUint64(*ltp.ChainID)
	privateKey := ltp.ECDSAPrivateKey
	iterations := ltp.Iterations
	f := contracts.GetRandomOPCode()

	tops, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Error().Err(err).Msg("Unable create transaction signer")
		return
	}
	tops.Nonce = new(big.Int).SetUint64(nonce)
	tops = configureTransactOpts(tops)

	t1 = time.Now()
	_, err = contracts.CallLoadTestFunctionByOpCode(f, ltContract, tops, *iterations)
	t2 = time.Now()
	return
}

func loadtestCallPrecompiledContracts(ctx context.Context, c *ethclient.Client, nonce uint64, ltContract *contracts.LoadTester, useSelectedAddress bool) (t1 time.Time, t2 time.Time, err error) {
	var f int
	ltp := inputLoadTestParams

	chainID := new(big.Int).SetUint64(*ltp.ChainID)
	privateKey := ltp.ECDSAPrivateKey
	iterations := ltp.Iterations
	if useSelectedAddress {
		f = int(*ltp.Function)
	} else {
		f = contracts.GetRandomPrecompiledContractAddress()
	}

	tops, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Error().Err(err).Msg("Unable create transaction signer")
		return
	}
	tops.Nonce = new(big.Int).SetUint64(nonce)
	tops = configureTransactOpts(tops)

	t1 = time.Now()
	_, err = contracts.CallPrecompiledContracts(f, ltContract, tops, *iterations, privateKey)
	t2 = time.Now()
	return
}

func loadtestInc(ctx context.Context, c *ethclient.Client, nonce uint64, ltContract *contracts.LoadTester) (t1 time.Time, t2 time.Time, err error) {
	ltp := inputLoadTestParams

	chainID := new(big.Int).SetUint64(*ltp.ChainID)
	privateKey := ltp.ECDSAPrivateKey

	tops, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Error().Err(err).Msg("Unable create transaction signer")
		return
	}
	tops.Nonce = new(big.Int).SetUint64(nonce)
	tops = configureTransactOpts(tops)

	t1 = time.Now()
	_, err = ltContract.Inc(tops)
	t2 = time.Now()
	return
}

func loadtestStore(ctx context.Context, c *ethclient.Client, nonce uint64, ltContract *contracts.LoadTester) (t1 time.Time, t2 time.Time, err error) {
	ltp := inputLoadTestParams

	chainID := new(big.Int).SetUint64(*ltp.ChainID)
	privateKey := ltp.ECDSAPrivateKey

	tops, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Error().Err(err).Msg("Unable create transaction signer")
		return
	}
	tops.Nonce = new(big.Int).SetUint64(nonce)
	tops = configureTransactOpts(tops)

	inputData := make([]byte, *ltp.ByteCount)
	_, _ = hexwordRead(inputData)
	t1 = time.Now()
	_, err = ltContract.Store(tops, inputData)
	t2 = time.Now()
	return
}

func loadtestLong(ctx context.Context, c *ethclient.Client, nonce uint64, delegatorContract *contracts.Delegator, ltAddress ethcommon.Address) (t1 time.Time, t2 time.Time, err error) {
	ltp := inputLoadTestParams

	chainID := new(big.Int).SetUint64(*ltp.ChainID)
	privateKey := ltp.ECDSAPrivateKey

	tops, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Error().Err(err).Msg("Unable create transaction signer")
		return
	}
	tops.Nonce = new(big.Int).SetUint64(nonce)
	tops = configureTransactOpts(tops)

	// TODO the delegated call should be a parameter
	t1 = time.Now()
	// loopBlockHashUntilLimit (verify here https://abi.hashex.org/)
	_, err = delegatorContract.LoopDelegateCall(tops, ltAddress, []byte{0xa2, 0x71, 0xb7, 0x21})
	// loopUntilLimit
	// _, err = delegatorContract.LoopDelegateCall(tops, ltAddress, []byte{0x65, 0x9b, 0xbb, 0x4f})
	t2 = time.Now()
	return
}

func loadtestERC20(ctx context.Context, c *ethclient.Client, nonce uint64, erc20Contract *contracts.ERC20, ltAddress ethcommon.Address) (t1 time.Time, t2 time.Time, err error) {
	ltp := inputLoadTestParams

	to := ltp.ToETHAddress
	if *ltp.ToRandom {
		to = getRandomAddress()
	}
	amount := ltp.SendAmount

	chainID := new(big.Int).SetUint64(*ltp.ChainID)
	privateKey := ltp.ECDSAPrivateKey

	tops, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Error().Err(err).Msg("Unable create transaction signer")
		return
	}
	tops.Nonce = new(big.Int).SetUint64(nonce)
	tops = configureTransactOpts(tops)

	t1 = time.Now()
	_, err = erc20Contract.Transfer(tops, *to, amount)
	t2 = time.Now()
	return
}

func loadtestERC721(ctx context.Context, c *ethclient.Client, nonce uint64, erc721Contract *contracts.ERC721, ltAddress ethcommon.Address) (t1 time.Time, t2 time.Time, err error) {
	ltp := inputLoadTestParams
	iterations := ltp.Iterations

	to := ltp.ToETHAddress
	if *ltp.ToRandom {
		to = getRandomAddress()
	}

	chainID := new(big.Int).SetUint64(*ltp.ChainID)
	privateKey := ltp.ECDSAPrivateKey

	tops, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Error().Err(err).Msg("Unable create transaction signer")
		return
	}
	tops.Nonce = new(big.Int).SetUint64(nonce)
	tops = configureTransactOpts(tops)

	t1 = time.Now()
	_, err = erc721Contract.MintBatch(tops, *to, new(big.Int).SetUint64(*iterations))
	t2 = time.Now()
	return
}

func recordSample(goRoutineID, requestID int64, err error, start, end time.Time, nonce uint64) {
	s := loadTestSample{}
	s.GoRoutineID = goRoutineID
	s.RequestID = requestID
	s.RequestTime = start
	s.WaitTime = end.Sub(start)
	s.Nonce = nonce
	if err != nil {
		s.IsError = true
	}
	loadTestResutsMutex.Lock()
	loadTestResults = append(loadTestResults, s)
	loadTestResutsMutex.Unlock()
}

func hexwordRead(b []byte) (int, error) {
	hw := hexwordReader{}
	return io.ReadFull(&hw, b)
}

func (hw *hexwordReader) Read(p []byte) (n int, err error) {
	hwLen := len(hexwords)
	for k := range p {
		p[k] = hexwords[k%hwLen]
	}
	n = len(p)
	return
}

func getRandomAddress() *ethcommon.Address {
	addr := make([]byte, 20)
	n, err := rand.Read(addr)
	if err != nil {
		log.Error().Err(err).Msg("There was an issue getting random bytes for the address")
	}
	if n != 20 {
		log.Error().Int("n", n).Msg("Somehow we didn't read 20 random bytes")
	}
	realAddr := ethcommon.BytesToAddress(addr)
	return &realAddr
}

func availLoop(ctx context.Context, c *gsrpc.SubstrateAPI) error {
	var err error

	ltp := inputLoadTestParams
	log.Trace().Interface("Input Params", ltp).Msg("Params")

	routines := *ltp.Concurrency
	requests := *ltp.Requests
	currentNonce := uint64(0) // *ltp.CurrentNonce
	chainID := new(big.Int).SetUint64(*ltp.ChainID)
	privateKey := ltp.ECDSAPrivateKey
	mode := *ltp.Mode

	_ = chainID
	_ = privateKey

	meta, err := c.RPC.State.GetMetadataLatest()
	if err != nil {
		return err
	}

	genesisHash, err := c.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return err
	}

	key, err := gstypes.CreateStorageKey(meta, "System", "Account", ltp.FromAvailAddress.PublicKey, nil)
	if err != nil {
		log.Error().Err(err).Msg("Could not create storage key")
		return err
	}

	var accountInfo gstypes.AccountInfo
	ok, err := c.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil {
		log.Error().Err(err).Msg("Could not load storage")
		return err
	}
	if !ok {
		err = fmt.Errorf("loaded storage is not okay")
		log.Error().Err(err).Msg("Loaded storage is not okay")
		return err
	}

	currentNonce = uint64(accountInfo.Nonce)

	rl := rate.NewLimiter(rate.Limit(*ltp.RateLimit), 1)
	if *ltp.RateLimit <= 0.0 {
		rl = nil
	}

	var currentNonceMutex sync.Mutex

	var i int64

	var wg sync.WaitGroup
	for i = 0; i < routines; i = i + 1 {
		log.Trace().Int64("routine", i).Msg("Starting Thread")
		wg.Add(1)
		go func(i int64) {
			var j int64
			var startReq time.Time
			var endReq time.Time

			for j = 0; j < requests; j = j + 1 {

				if rl != nil {
					err = rl.Wait(ctx)
					if err != nil {
						log.Error().Err(err).Msg("Encountered a rate limiting error")
					}
				}
				currentNonceMutex.Lock()
				myNonceValue := currentNonce
				currentNonce = currentNonce + 1
				currentNonceMutex.Unlock()

				localMode := mode
				// if there are multiple modes, iterate through them, 'r' mode is supported here
				if len(mode) > 1 {
					localMode = string(mode[int(i+j)%(len(mode))])
				}
				// if we're doing random, we'll just pick one based on the current index
				if localMode == loadTestModeRandom {
					localMode = validLoadTestModes[int(i+j)%(len(validLoadTestModes)-1)]
				}
				// this function should probably be abstracted
				switch localMode {
				case loadTestModeTransaction:
					startReq, endReq, err = loadtestAvailTransfer(ctx, c, myNonceValue, meta, genesisHash)
				case loadTestModeDeploy:
					startReq, endReq, err = loadtestNotImplemented(ctx, c, myNonceValue)
				case loadTestModeCall:
					startReq, endReq, err = loadtestNotImplemented(ctx, c, myNonceValue)
				case loadTestModeFunction:
					startReq, endReq, err = loadtestNotImplemented(ctx, c, myNonceValue)
				case loadTestModeInc:
					startReq, endReq, err = loadtestNotImplemented(ctx, c, myNonceValue)
				case loadTestModeStore:
					startReq, endReq, err = loadtestAvailStore(ctx, c, myNonceValue, meta, genesisHash)
				case loadTestModeLong:
					startReq, endReq, err = loadtestNotImplemented(ctx, c, myNonceValue)
				default:
					log.Error().Str("mode", mode).Msg("We've arrived at a load test mode that we don't recognize")
				}
				recordSample(i, j, err, startReq, endReq, myNonceValue)
				if err != nil {
					log.Trace().Err(err).Msg("Recorded an error while sending transactions")
				}

				log.Trace().Int64("routine", i).Str("mode", localMode).Int64("request", j).Msg("Request")
			}
			wg.Done()
		}(i)

	}
	log.Trace().Msg("Finished starting go routines. Waiting..")
	wg.Wait()
	return nil

}

func loadtestNotImplemented(ctx context.Context, c *gsrpc.SubstrateAPI, nonce uint64) (t1 time.Time, t2 time.Time, err error) {
	t1 = time.Now()
	t2 = time.Now()
	err = fmt.Errorf("this method is not implemented")
	return
}

func initAvailTestParams(ctx context.Context, c *gsrpc.SubstrateAPI) error {
	toAddr, err := gstypes.NewMultiAddressFromHexAccountID(*inputLoadTestParams.ToAddress)
	if err != nil {
		log.Error().Err(err).Msg("Unable to create new multi address")
		return err
	}

	if *inputLoadTestParams.PrivateKey == codeQualityPrivateKey {
		// Avail keys can use the same seed but the way the key is derived is different
		*inputLoadTestParams.PrivateKey = codeQualitySeed
	}

	kp, err := gssignature.KeyringPairFromSecret(*inputLoadTestParams.PrivateKey, uint8(*inputLoadTestParams.ChainID))
	if err != nil {
		log.Error().Err(err).Msg("Could not create key pair")
		return err
	}

	amt, err := hexToBigInt(*inputLoadTestParams.HexSendAmount)
	if err != nil {
		log.Error().Err(err).Msg("Couldn't parse send amount")
		return err
	}

	rv, err := c.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		log.Error().Err(err).Msg("Couldn't get runtime version")
		return err
	}

	inputLoadTestParams.AvailRuntime = rv
	inputLoadTestParams.SendAmount = amt
	inputLoadTestParams.FromAvailAddress = &kp
	inputLoadTestParams.ToAvailAddress = &toAddr
	return nil
}

func loadtestAvailTransfer(ctx context.Context, c *gsrpc.SubstrateAPI, nonce uint64, meta *gstypes.Metadata, genesisHash gstypes.Hash) (t1 time.Time, t2 time.Time, err error) {
	ltp := inputLoadTestParams

	toAddr := *ltp.ToAvailAddress
	if *ltp.ToRandom {
		pk := make([]byte, 32)
		_, err = rand.Read(pk)
		if err != nil {
			// For some reason weren't able to read the random data
			log.Error().Msg("Sending to random is not implemented for substrate yet")
		} else {
			toAddr = gstypes.NewMultiAddressFromAccountID(pk)
		}

	}

	gsCall, err := gstypes.NewCall(meta, "Balances.transfer", toAddr, gstypes.NewUCompact(ltp.SendAmount))
	if err != nil {
		return
	}

	ext := gstypes.NewExtrinsic(gsCall)
	rv := ltp.AvailRuntime
	kp := *inputLoadTestParams.FromAvailAddress

	o := gstypes.SignatureOptions{
		BlockHash:          genesisHash,
		Era:                gstypes.ExtrinsicEra{IsMortalEra: false, IsImmortalEra: true},
		GenesisHash:        genesisHash,
		Nonce:              gstypes.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                gstypes.NewUCompactFromUInt(0),
		TransactionVersion: rv.TransactionVersion,
	}

	err = ext.Sign(kp, o)
	if err != nil {
		return
	}

	t1 = time.Now()
	_, err = c.RPC.Author.SubmitExtrinsic(ext)
	t2 = time.Now()
	if err != nil {
		return
	}
	return
}

func loadtestAvailStore(ctx context.Context, c *gsrpc.SubstrateAPI, nonce uint64, meta *gstypes.Metadata, genesisHash gstypes.Hash) (t1 time.Time, t2 time.Time, err error) {
	ltp := inputLoadTestParams

	inputData := make([]byte, *ltp.ByteCount)
	_, _ = hexwordRead(inputData)

	gsCall, err := gstypes.NewCall(meta, "DataAvailability.submit_data", gstypes.NewBytes([]byte(inputData)))
	if err != nil {
		return
	}

	// Create the extrinsic
	ext := gstypes.NewExtrinsic(gsCall)

	rv := ltp.AvailRuntime

	kp := *inputLoadTestParams.FromAvailAddress

	o := gstypes.SignatureOptions{
		BlockHash:          genesisHash,
		Era:                gstypes.ExtrinsicEra{IsMortalEra: false, IsImmortalEra: true},
		GenesisHash:        genesisHash,
		Nonce:              gstypes.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                gstypes.NewUCompactFromUInt(100),
		TransactionVersion: rv.TransactionVersion,
	}
	// Sign the transaction using Alice's default account
	err = ext.Sign(kp, o)
	if err != nil {
		return
	}

	// Send the extrinsic
	t1 = time.Now()
	_, err = c.RPC.Author.SubmitExtrinsic(ext)
	t2 = time.Now()
	if err != nil {
		return
	}
	return
}

func configureTransactOpts(tops *bind.TransactOpts) *bind.TransactOpts {
	ltp := inputLoadTestParams
	if *ltp.LegacyTransactionMode {
		if ltp.ForceGasPrice != nil && *ltp.ForceGasPrice != 0 {
			tops.GasPrice = big.NewInt(0).SetUint64(*ltp.ForceGasPrice)
		} else {
			tops.GasPrice = ltp.CurrentGas
		}
	} else {
		if ltp.ForcePriorityGasPrice != nil && *ltp.ForcePriorityGasPrice != 0 {
			tops.GasTipCap = big.NewInt(0).SetUint64(*ltp.ForcePriorityGasPrice)
		} else {
			tops.GasTipCap = ltp.CurrentGasTipCap
		}
	}
	if ltp.ForceGasLimit != nil && *ltp.ForceGasLimit != 0 {
		tops.GasLimit = *ltp.ForceGasLimit
	}
	return tops
}

func waitForFinalBlock(ctx context.Context, c *ethclient.Client, rpc *ethrpc.Client, startBlockNumber, startNonce, endNonce uint64) (uint64, error) {
	ltp := inputLoadTestParams
	var err error
	var lastBlockNumber uint64
	var currentNonce uint64
	var initialWaitCount = 50
	var maxWaitCount = initialWaitCount
	for {
		lastBlockNumber, err = c.BlockNumber(ctx)
		if err != nil {
			return 0, err
		}
		currentNonce, err = c.NonceAt(ctx, *ltp.FromETHAddress, new(big.Int).SetUint64(lastBlockNumber))
		if err != nil {
			return 0, err
		}
		if currentNonce < endNonce && maxWaitCount > 0 {
			log.Trace().Uint64("endNonce", endNonce).Uint64("currentNonce", currentNonce).Msg("Not all transactions have been mined. Waiting")
			time.Sleep(5 * time.Second)
			maxWaitCount = maxWaitCount - 1
			continue
		}
		if maxWaitCount <= 0 {
			return 0, fmt.Errorf("waited for %d attempts for the transactions to be mined", initialWaitCount)
		}
		break
	}

	log.Trace().Uint64("currentNonce", currentNonce).Uint64("startblock", startBlockNumber).Uint64("endblock", lastBlockNumber).Msg("It looks like all transactions have been mined")
	return lastBlockNumber, nil
}

func summarizeTransactions(ctx context.Context, c *ethclient.Client, rpc *ethrpc.Client, startBlockNumber, startNonce, lastBlockNumber, endNonce uint64) error {
	ltp := inputLoadTestParams
	var err error

	log.Trace().Msg("Starting block range capture")
	// confirm start block number is ok
	_, err = c.BlockByNumber(ctx, new(big.Int).SetUint64(startBlockNumber))
	if err != nil {
		return err
	}
	rawBlocks, err := util.GetBlockRange(ctx, startBlockNumber, lastBlockNumber, rpc)
	if err != nil {
		return err
	}
	// TODO: Add some kind of decimation to avoid summarizing for 10 minutes?
	batchSize := *ltp.BatchSize
	goRoutineLimit := *ltp.Concurrency
	var txGroup sync.WaitGroup
	threadPool := make(chan bool, goRoutineLimit)
	log.Trace().Msg("Starting tx receipt capture")
	rawTxReceipts := make([]*json.RawMessage, 0)
	var rawTxReceiptsLock sync.Mutex
	var txGroupErr error

	startReceipt := time.Now()
	for k := range rawBlocks {
		threadPool <- true
		txGroup.Add(1)
		go func(b *json.RawMessage) {
			var receipt []*json.RawMessage
			receipt, err = util.GetReceipts(ctx, []*json.RawMessage{b}, rpc, batchSize)
			if err != nil {
				txGroupErr = err
				return
			}
			rawTxReceiptsLock.Lock()
			rawTxReceipts = append(rawTxReceipts, receipt...)
			rawTxReceiptsLock.Unlock()
			<-threadPool
			txGroup.Done()
		}(rawBlocks[k])
	}

	endReceipt := time.Now()
	txGroup.Wait()
	if txGroupErr != nil {
		log.Error().Err(err).Msg("One of the threads fetching tx receipts failed")
		return err
	}

	blocks := make([]rpctypes.RawBlockResponse, 0)
	for _, b := range rawBlocks {
		var block rpctypes.RawBlockResponse
		err = json.Unmarshal(*b, &block)
		if err != nil {
			log.Error().Err(err).Msg("Error decoding block response")
			return err
		}
		blocks = append(blocks, block)
	}
	log.Info().Int("len", len(blocks)).Msg("Block summary")

	txReceipts := make([]rpctypes.RawTxReceipt, 0)
	log.Trace().Int("len", len(rawTxReceipts)).Msg("Raw receipts")
	for _, r := range rawTxReceipts {
		if isEmptyJSONResponse(r) {
			continue
		}
		var receipt rpctypes.RawTxReceipt
		err = json.Unmarshal(*r, &receipt)
		if err != nil {
			log.Error().Err(err).Msg("Error decoding tx receipt response")
			return err
		}
		txReceipts = append(txReceipts, receipt)
	}
	log.Info().Int("len", len(txReceipts)).Msg("Receipt summary")

	blockData := make(map[uint64]blockSummary, 0)
	for k, b := range blocks {
		bs := blockSummary{}
		bs.Block = &blocks[k]
		bs.Receipts = make(map[ethcommon.Hash]rpctypes.RawTxReceipt, 0)
		bs.Latencies = make(map[uint64]time.Duration, 0)
		blockData[b.Number.ToUint64()] = bs
	}

	for _, r := range txReceipts {
		bn := r.BlockNumber.ToUint64()
		bs := blockData[bn]
		if bs.Receipts == nil {
			log.Error().Uint64("blocknumber", bn).Msg("Block number from receipts does not exist in block data")
		}
		bs.Receipts[r.TransactionHash.ToHash()] = r
		blockData[bn] = bs
	}

	nonceTimes := make(map[uint64]time.Time, 0)
	for _, ltr := range loadTestResults {
		nonceTimes[ltr.Nonce] = ltr.RequestTime
	}

	minLatency := time.Millisecond * 100
	for _, bs := range blockData {
		for _, tx := range bs.Block.Transactions {
			// TODO: What happens when the system clock of the load tester isn't in sync with the system clock of the miner?
			// TODO: the timestamp in the chain only has granularity down to the second. How to deal with this
			mineTime := time.Unix(bs.Block.Timestamp.ToInt64(), 0)
			requestTime := nonceTimes[tx.Nonce.ToUint64()]
			txLatency := mineTime.Sub(requestTime)
			if txLatency.Hours() > 2 {
				log.Debug().Float64("txHours", txLatency.Hours()).Uint64("nonce", tx.Nonce.ToUint64()).Uint64("blockNumber", bs.Block.Number.ToUint64()).Time("mineTime", mineTime).Time("requestTime", requestTime).Msg("Encountered transaction with more than 2 hours latency")
			}
			bs.Latencies[tx.Nonce.ToUint64()] = txLatency
			if txLatency < minLatency {
				minLatency = txLatency
			}
		}
	}
	// TODO this might be a hack, but not sure what's a better way to deal with time discrepancies
	if minLatency < time.Millisecond*100 {
		log.Trace().Str("minLatency", minLatency.String()).Msg("Minimum latency is below expected threshold")
		shiftSize := ((time.Millisecond * 100) - minLatency) + time.Millisecond + 100
		for _, bs := range blockData {
			for _, tx := range bs.Block.Transactions {
				bs.Latencies[tx.Nonce.ToUint64()] += shiftSize
			}
		}
	}

	printBlockSummary(c, blockData, startNonce, endNonce)

	log.Trace().Str("summaryTime", (endReceipt.Sub(startReceipt)).String()).Msg("Total Summary Time")
	return nil

}

func isEmptyJSONResponse(r *json.RawMessage) bool {
	rawJson := []byte(*r)
	return len(rawJson) == 0
}

type Latency struct {
	Min    float64
	Median float64
	Max    float64
}

type Summary struct {
	BlockNumber uint64
	Time        time.Time
	GasLimit    uint64
	GasUsed     uint64
	NumTx       int
	Utilization float64
	Latencies   Latency
}

type SummaryOutput struct {
	Summaries          []Summary
	SuccessfulTx       int64
	TotalTx            int64
	TotalMiningTime    time.Duration
	TotalGasUsed       uint64
	TransactionsPerSec float64
	GasPerSecond       float64
	Latencies          Latency
}

func printBlockSummary(c *ethclient.Client, bs map[uint64]blockSummary, startNonce, endNonce uint64) {
	filterBlockSummary(bs, startNonce, endNonce)
	mapKeys := getSortedMapKeys(bs)
	if len(mapKeys) == 0 {
		return
	}

	var totalTransactions uint64 = 0
	var totalGasUsed uint64 = 0
	p := message.NewPrinter(language.English)

	allLatencies := make([]time.Duration, 0)
	summaryOutputMode := *inputLoadTestParams.SummaryOutputMode
	jsonSummaryList := []Summary{}
	for _, v := range mapKeys {
		summary := bs[v]
		gasUsed := getTotalGasUsed(summary.Receipts)
		blockLatencies := getMapValues(summary.Latencies)
		minLatency, medianLatency, maxLatency := getMinMedianMax(blockLatencies)
		allLatencies = append(allLatencies, blockLatencies...)
		blockUtilization := float64(gasUsed) / summary.Block.GasLimit.ToFloat64()
		if gasUsed == 0 {
			blockUtilization = 0
		}
		// if we're at trace, debug, or info level we'll output the block level metrics
		if zerolog.GlobalLevel() <= zerolog.InfoLevel {
			if summaryOutputMode == "text" {
				_, _ = p.Printf("Block number: %v\tTime: %s\tGas Limit: %v\tGas Used: %v\tNum Tx: %v\tUtilization %v\tLatencies: %v\t%v\t%v\n",
					number.Decimal(summary.Block.Number.ToUint64()),
					time.Unix(summary.Block.Timestamp.ToInt64(), 0),
					number.Decimal(summary.Block.GasLimit.ToUint64()),
					number.Decimal(gasUsed),
					number.Decimal(len(summary.Block.Transactions)),
					number.Percent(blockUtilization),
					number.Decimal(minLatency.Seconds()),
					number.Decimal(medianLatency.Seconds()),
					number.Decimal(maxLatency.Seconds()))
			} else if summaryOutputMode == "json" {
				jsonSummary := Summary{}
				jsonSummary.BlockNumber = summary.Block.Number.ToUint64()
				jsonSummary.Time = time.Unix(summary.Block.Timestamp.ToInt64(), 0)
				jsonSummary.GasLimit = summary.Block.GasLimit.ToUint64()
				jsonSummary.GasUsed = gasUsed
				jsonSummary.NumTx = len(summary.Block.Transactions)
				jsonSummary.Utilization = blockUtilization
				latencies := Latency{}
				latencies.Min = minLatency.Seconds()
				latencies.Median = medianLatency.Seconds()
				latencies.Max = maxLatency.Seconds()
				jsonSummary.Latencies = latencies
				jsonSummaryList = append(jsonSummaryList, jsonSummary)
			} else {
				log.Error().Str("mode", summaryOutputMode).Msg("Invalid mode for summary output")
			}
		}
		totalTransactions += uint64(len(summary.Block.Transactions))
		totalGasUsed += gasUsed
	}
	parentOfFirstBlock, _ := c.BlockByNumber(context.Background(), big.NewInt(bs[mapKeys[0]].Block.Number.ToInt64()-1))
	lastBlock := bs[mapKeys[len(mapKeys)-1]].Block
	totalMiningTime := time.Duration(lastBlock.Timestamp.ToUint64()-parentOfFirstBlock.Time()) * time.Second
	tps := float64(totalTransactions) / totalMiningTime.Seconds()
	gaspersec := float64(totalGasUsed) / totalMiningTime.Seconds()
	minLatency, medianLatency, maxLatency := getMinMedianMax(allLatencies)
	successfulTx, totalTx := getSuccessfulTransactionCount(bs)

	if summaryOutputMode == "text" {
		p.Printf("Successful Tx: %v\tTotal Tx: %v\n", number.Decimal(successfulTx), number.Decimal(totalTx))
		p.Printf("Total Mining Time: %s\n", totalMiningTime)
		p.Printf("Total Transactions: %v\n", number.Decimal(totalTransactions))
		p.Printf("Total Gas Used: %v\n", number.Decimal(totalGasUsed))
		p.Printf("Transactions per sec: %v\n", number.Decimal(tps))
		p.Printf("Gas Per Second: %v\n", number.Decimal(gaspersec))
		p.Printf("Latencies - Min: %v\tMedian: %v\tMax: %v\n", number.Decimal(minLatency.Seconds()), number.Decimal(medianLatency.Seconds()), number.Decimal(maxLatency.Seconds()))
		// TODO: Add some kind of indication of block time variance
	} else if summaryOutputMode == "json" {
		summaryOutput := SummaryOutput{}
		summaryOutput.Summaries = jsonSummaryList
		summaryOutput.SuccessfulTx = successfulTx
		summaryOutput.TotalTx = totalTx
		summaryOutput.TotalMiningTime = totalMiningTime
		summaryOutput.TotalGasUsed = totalGasUsed
		summaryOutput.TransactionsPerSec = tps
		summaryOutput.GasPerSecond = gaspersec

		latencies := Latency{}
		latencies.Min = minLatency.Seconds()
		latencies.Median = medianLatency.Seconds()
		latencies.Max = maxLatency.Seconds()
		summaryOutput.Latencies = latencies

		val, _ := json.MarshalIndent(summaryOutput, "", "    ")
		p.Println(string(val))
	} else {
		log.Error().Str("mode", summaryOutputMode).Msg("Invalid mode for summary output")
	}
}

func getSuccessfulTransactionCount(bs map[uint64]blockSummary) (successful, total int64) {
	for _, block := range bs {
		total += int64(len(block.Receipts))
		for _, receipt := range block.Receipts {
			successful += receipt.Status.ToInt64()
		}
	}
	return
}

func getTotalGasUsed(receipts map[ethcommon.Hash]rpctypes.RawTxReceipt) uint64 {
	var totalGasUsed uint64 = 0
	for _, receipt := range receipts {
		totalGasUsed += receipt.GasUsed.ToUint64()
	}
	return totalGasUsed
}

func getMapValues[K constraints.Ordered, V any](m map[K]V) []V {
	newSlice := make([]V, 0)
	for _, val := range m {
		newSlice = append(newSlice, val)
	}
	return newSlice
}

func getMinMedianMax[V constraints.Float | constraints.Integer](values []V) (V, V, V) {
	if len(values) == 0 {
		return 0, 0, 0
	}
	sort.Slice(values, func(i, j int) bool {
		return values[i] < values[j]
	})
	half := len(values) / 2
	median := values[half]
	if len(values)%2 == 0 {
		median = (median + values[half-1]) / V(2)
	}
	var min V
	var max V
	for k, v := range values {
		if k == 0 {
			min = v
			max = v
			continue
		}
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return min, median, max
}

func filterBlockSummary(blockSummaries map[uint64]blockSummary, startNonce, endNonce uint64) {
	validTx := make(map[ethcommon.Hash]struct{}, 0)
	var minBlock uint64 = math.MaxUint64
	var maxBlock uint64 = 0
	for _, bs := range blockSummaries {
		for _, tx := range bs.Block.Transactions {
			if tx.Nonce.ToUint64() >= startNonce && tx.Nonce.ToUint64() <= endNonce {
				validTx[tx.Hash.ToHash()] = struct{}{}
				if tx.BlockNumber.ToUint64() < minBlock {
					minBlock = tx.BlockNumber.ToUint64()
				}
				if tx.BlockNumber.ToUint64() > maxBlock {
					maxBlock = tx.BlockNumber.ToUint64()
				}
			}
		}
	}
	keys := getSortedMapKeys(blockSummaries)
	for _, k := range keys {
		if k < minBlock {
			delete(blockSummaries, k)
		}
		if k > maxBlock {
			delete(blockSummaries, k)
		}
	}

	for _, bs := range blockSummaries {
		filteredTransactions := make([]rpctypes.RawTransactionResponse, 0)
		for txKey, tx := range bs.Block.Transactions {
			if _, hasKey := validTx[tx.Hash.ToHash()]; hasKey {
				filteredTransactions = append(filteredTransactions, bs.Block.Transactions[txKey])
			}
		}
		bs.Block.Transactions = filteredTransactions
		filteredReceipts := make(map[ethcommon.Hash]rpctypes.RawTxReceipt, 0)
		for receiptKey, receipt := range bs.Receipts {
			if _, hasKey := validTx[receipt.TransactionHash.ToHash()]; hasKey {
				filteredReceipts[receipt.TransactionHash.ToHash()] = bs.Receipts[receiptKey]
			}
		}
		bs.Receipts = filteredReceipts

	}
}

func getSortedMapKeys[V any, K constraints.Ordered](m map[K]V) []K {
	keys := make([]K, 0)
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return keys
}
