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
package cmd

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"sync"
	"time"

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
	loadTestModeTransaction = "t"
	loadTestModeDeploy      = "d"
	loadTestModeCall        = "c"
	loadTestModeFunction    = "f"
	loadTestModeInc         = "i"
	loadTestModeRandom      = "r"
	loadTestModeStore       = "s"
	loadTestModeLong        = "l"

	codeQualitySeed       = "code code code code code code code code code code code quality"
	codeQualityPrivateKey = "42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa"
)

var (
	inputLoadTestParams loadTestParams
	loadTestResults     []loadTestSample
	validLoadTestModes  = []string{
		loadTestModeTransaction,
		loadTestModeDeploy,
		loadTestModeCall,
		loadTestModeFunction,
		loadTestModeInc,
		loadTestModeStore,
		loadTestModeLong,
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

// loadtestCmd represents the loadtest command
var loadtestCmd = &cobra.Command{
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
		setLogLevel(inputLoadTestParams)
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
		return nil
	},
}

func setLogLevel(ltp loadTestParams) {
	verbosity := *ltp.Verbosity
	if verbosity < 100 {
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	} else if verbosity < 200 {
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	} else if verbosity < 300 {
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	} else if verbosity < 400 {
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	} else if verbosity < 500 {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else if verbosity < 600 {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}
	if *ltp.PrettyLogs {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		log.Debug().Msg("Starting logger in console mode")
	} else {
		log.Debug().Msg("Starting logger in JSON mode")
	}
}

type (
	hexwordReader struct {
	}
	loadTestSample struct {
		GoRoutineID int64
		RequestID   int64
		RequestTime time.Time
		WaitTime    time.Duration
		Receipt     string
		IsError     bool
	}
	loadTestParams struct {
		// inputs
		Requests            *int64
		Concurrency         *int64
		TimeLimit           *int64
		Verbosity           *int64
		PrettyLogs          *bool
		ToRandom            *bool
		URL                 *url.URL
		ChainID             *uint64
		PrivateKey          *string
		ToAddress           *string
		HexSendAmount       *string
		RateLimit           *float64
		Mode                *string
		Function            *uint64
		Iterations          *uint64
		ByteCount           *uint64
		Seed                *int64
		IsAvail             *bool
		AvailAppID          *uint32
		LtAddress           *string
		DelAddress          *string
		ForceContractDeploy *bool

		// Computed
		CurrentGas      *big.Int
		CurrentNonce    *uint64
		ECDSAPrivateKey *ecdsa.PrivateKey
		FromETHAddress  *ethcommon.Address
		ToETHAddress    *ethcommon.Address
		SendAmount      *big.Int

		ToAvailAddress   *gstypes.MultiAddress
		FromAvailAddress *gssignature.KeyringPair
		AvailRuntime     *gstypes.RuntimeVersion
	}
)

func init() {
	rootCmd.AddCommand(loadtestCmd)

	ltp := new(loadTestParams)
	// Apache Bench Parameters
	ltp.Requests = loadtestCmd.PersistentFlags().Int64P("requests", "n", 1, "Number of requests to perform for the benchmarking session. The default is to just perform a single request which usually leads to non-representative benchmarking results.")
	ltp.Concurrency = loadtestCmd.PersistentFlags().Int64P("concurrency", "c", 1, "Number of multiple requests to perform at a time. Default is one request at a time.")
	ltp.TimeLimit = loadtestCmd.PersistentFlags().Int64P("time-limit", "t", -1, "Maximum number of seconds to spend for benchmarking. Use this to benchmark within a fixed total amount of time. Per default there is no timelimit.")
	// https://logging.apache.org/log4j/2.x/manual/customloglevels.html
	ltp.Verbosity = loadtestCmd.PersistentFlags().Int64P("verbosity", "v", 200, "0 - Silent\n100 Fatals\n200 Errors\n300 Warnings\n400 INFO\n500 Debug\n600 Trace")

	// extended parameters
	ltp.PrettyLogs = loadtestCmd.PersistentFlags().Bool("pretty-logs", true, "Should we log in pretty format or JSON")
	ltp.PrivateKey = loadtestCmd.PersistentFlags().String("private-key", codeQualityPrivateKey, "The hex encoded private key that we'll use to sending transactions")
	ltp.ChainID = loadtestCmd.PersistentFlags().Uint64("chain-id", 1256, "The chain id for the transactions that we're going to send")
	ltp.ToAddress = loadtestCmd.PersistentFlags().String("to-address", "0xDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEF", "The address that we're going to send to")
	ltp.ToRandom = loadtestCmd.PersistentFlags().Bool("to-random", true, "When doing a transfer test, should we send to random addresses rather than DEADBEEFx5")
	ltp.HexSendAmount = loadtestCmd.PersistentFlags().String("send-amount", "0x38D7EA4C68000", "The amount of wei that we'll send every transaction")
	ltp.RateLimit = loadtestCmd.PersistentFlags().Float64("rate-limit", 4, "An overall limit to the number of requests per second. Give a number less than zero to remove this limit all together")
	ltp.Mode = loadtestCmd.PersistentFlags().StringP("mode", "m", "t", `The testing mode to use. It can be multiple like: "tcdf"
t - sending transactions
d - deploy contract
c - call random contract functions
f - call specific contract function
s - store mode
l - long running mode
r - random modes`)
	ltp.Function = loadtestCmd.PersistentFlags().Uint64P("function", "f", 1, "A specific function to be called if running with `--mode f` ")
	ltp.Iterations = loadtestCmd.PersistentFlags().Uint64P("iterations", "i", 100, "If we're making contract calls, this controls how many times the contract will execute the instruction in a loop")
	ltp.ByteCount = loadtestCmd.PersistentFlags().Uint64P("byte-count", "b", 1024, "If we're in store mode, this controls how many bytes we'll try to store in our contract")
	ltp.Seed = loadtestCmd.PersistentFlags().Int64("seed", 123456, "A seed for generating random values and addresses")
	ltp.IsAvail = loadtestCmd.PersistentFlags().Bool("data-avail", false, "Is this a test of avail rather than an EVM / Geth Chain")
	ltp.AvailAppID = loadtestCmd.PersistentFlags().Uint32("app-id", 0, "The AppID used for avail")
	ltp.LtAddress = loadtestCmd.PersistentFlags().String("lt-address", "", "A pre-deployed load test contract address")
	ltp.DelAddress = loadtestCmd.PersistentFlags().String("del-address", "", "A pre-deployed delegator contract address")
	ltp.ForceContractDeploy = loadtestCmd.PersistentFlags().Bool("force-contract-deploy", false, "Some loadtest modes don't require a contract deployment. Set this flag to true to force contract deployments. This will still respect the --del-address and --il-address flags.")

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
		log.Error().Err(err).Msg("couldn't parse send amount")
		return err
	}

	inputLoadTestParams.ToETHAddress = &toAddr
	inputLoadTestParams.SendAmount = amt
	inputLoadTestParams.CurrentGas = gas
	inputLoadTestParams.CurrentNonce = &nonce
	inputLoadTestParams.ECDSAPrivateKey = privateKey
	inputLoadTestParams.FromETHAddress = &ethAddress

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

			api, err := gsrpc.NewSubstrateAPI(inputLoadTestParams.URL.String())
			if err != nil {
				return err
			}
			err = initAvailTestParams(ctx, api)
			return availLoop(ctx, api)
		}

	} else {
		log.Info().Msg("Starting Load Test")
		loopFunc = func() error {
			err = initializeLoadTestParams(ctx, ec)
			if err != nil {
				return err
			}

			return mainLoop(ctx, ec)
		}
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	loadTestResults = make([]loadTestSample, 0)
	errCh := make(chan error)
	go func() {
		// errCh <- mainLoop(ctx, ec)
		errCh <- loopFunc()
	}()

	select {
	case <-overallTimer.C:
		log.Info().Msg("Time's up")
	case <-sigCh:
		log.Info().Msg("Interrupted.. Stopping load test")
	case err := <-errCh:
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
		log.Error().Err(err).Msg("Unable to get the number of pending transactions before closing")
	} else if ptc > 0 {
		log.Info().Uint("pending", ptc).Msg("there are still oustanding transactions. There might be issues restarting with the same sending key until those transactions clear")
	}
	return nil
}

func printResults(lts []loadTestSample) {
	if len(lts) == 0 {
		log.Error().Msg("No results recorded")
		return
	}

	fmt.Println("* Results")
	fmt.Printf("Samples: %d\n", len(lts))

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
	fmt.Printf("Start: %s\n", startTime)
	fmt.Printf("End: %s\n", endTime)
	fmt.Printf("Mean Wait: %0.4f\n", meanWait)
	fmt.Printf("Num errors: %d\n", numErrors)
}

func mainLoop(ctx context.Context, c *ethclient.Client) error {

	ltp := inputLoadTestParams
	log.Trace().Interface("Input Params", ltp).Msg("Params")

	routines := *ltp.Concurrency
	requests := *ltp.Requests
	currentNonce := *ltp.CurrentNonce
	chainID := new(big.Int).SetUint64(*ltp.ChainID)
	privateKey := ltp.ECDSAPrivateKey
	mode := *ltp.Mode

	rl := rate.NewLimiter(rate.Limit(*ltp.RateLimit), 1)
	if *ltp.RateLimit <= 0.0 {
		rl = nil
	}

	tops, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Error().Err(err).Msg("Unable create transaction signer")
		return err
	}
	cops := new(bind.CallOpts)

	// deploy and instantiate the load tester contract
	var ltAddr ethcommon.Address
	var ltContract *contracts.LoadTester
	if strings.ContainsAny(mode, "rcfisl") || *inputLoadTestParams.ForceContractDeploy {
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

		// block while the contract is pending
		waitCounter := 30
		for {
			ltCounter, err := ltContract.GetCallCounter(cops)

			if err != nil {
				log.Trace().Msg("Waiting for Load Test contract to deploy")
				time.Sleep(time.Second)
				if waitCounter < 1 {
					log.Error().Err(err).Msg("Exhausted waiting period for contract deployment")
					return err
				}
				waitCounter = waitCounter - 1
				continue
			}
			log.Trace().Interface("counter", ltCounter).Msg("Number of contract calls")
			break
		}
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

		// block while the contract is pending
		waitCounter := 30
		for {
			_, err = delegatorContract.Call(tops, ltAddr, []byte{0x12, 0x87, 0xa6, 0x8c})
			if err != nil {
				log.Trace().Msg("Waiting for Delegator contract to deploy")
				time.Sleep(time.Second)
				if waitCounter < 1 {
					log.Error().Err(err).Msg("Exhausted waiting period for contract deployment")
					return err
				}
				waitCounter = waitCounter - 1
				continue
			}
			break
		}

		currentNonce = currentNonce + 1
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
					err := rl.Wait(ctx)
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
				switch localMode {
				case loadTestModeTransaction:
					startReq, endReq, err = loadtestTransaction(ctx, c, myNonceValue)
					break
				case loadTestModeDeploy:
					startReq, endReq, err = loadtestDeploy(ctx, c, myNonceValue)
					break
				case loadTestModeCall:
					startReq, endReq, err = loadtestCall(ctx, c, myNonceValue, ltContract)
					break
				case loadTestModeFunction:
					startReq, endReq, err = loadtestFunction(ctx, c, myNonceValue, ltContract)
					break
				case loadTestModeInc:
					startReq, endReq, err = loadtestInc(ctx, c, myNonceValue, ltContract)
					break
				case loadTestModeStore:
					startReq, endReq, err = loadtestStore(ctx, c, myNonceValue, ltContract)
					break
				case loadTestModeLong:
					startReq, endReq, err = loadtestLong(ctx, c, myNonceValue, delegatorContract, ltAddr)
					break
				default:
					log.Error().Str("mode", mode).Msg("We've arrived at a load test mode that we don't recognize")
				}
				recordSample(i, j, err, startReq, endReq)
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

	gasLimit := uint64(21000)
	tx := ethtypes.NewTransaction(nonce, *to, amount, gasLimit, gasPrice, nil)
	stx, err := ethtypes.SignTx(tx, ethtypes.NewEIP155Signer(chainID), privateKey)
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

	t1 = time.Now()
	_, err = contracts.CallLoadTestFunctionByOpCode(f, ltContract, tops, *iterations)
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

	inputData := make([]byte, *ltp.ByteCount, *ltp.ByteCount)
	hexwordRead(inputData)
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
	tops.GasLimit = 10000000

	// TODO the deletgated call should be a parameter
	t1 = time.Now()
	// loopBlockHashUntilLimit
	_, err = delegatorContract.LoopDelegateCall(tops, ltAddress, []byte{0xa2, 0x71, 0xb7, 0x21})
	// loopUntilLimit
	// _, err = delegatorContract.LoopDelegateCall(tops, ltAddress, []byte{0x65, 0x9b, 0xbb, 0x4f})
	t2 = time.Now()
	return
}

func recordSample(goRoutineID, requestID int64, err error, start, end time.Time) {
	s := loadTestSample{}
	s.GoRoutineID = goRoutineID
	s.RequestID = requestID
	s.RequestTime = start
	s.WaitTime = end.Sub(start)
	if err != nil {
		s.IsError = true
	}
	loadTestResults = append(loadTestResults, s)
}

/*
// This function is unused
func createLoadTesterContract(ctx context.Context, c *ethclient.Client, nonce uint64, gasPrice *big.Int) (*ethtypes.Receipt, error) {
	var gasLimit uint64 = 0x192f64
	contract, err := contracts.GetLoadTesterBytes()
	if err != nil {
		return nil, err
	}

	ltp := inputLoadTestParams
	chainID := new(big.Int).SetUint64(*ltp.ChainID)
	privateKey := ltp.ECDSAPrivateKey

	tx := ethtypes.NewContractCreation(nonce, big.NewInt(0), gasLimit, gasPrice, contract)
	stx, err := ethtypes.SignTx(tx, ethtypes.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Error().Err(err).Msg("Unable to sign transaction")
		return nil, err
	}

	err = c.SendTransaction(ctx, stx)
	if err != nil {
		return nil, err
	}

	wait := time.Millisecond * 500
	for i := 0; i < 5; i = i + 1 {
		receipt, err := c.TransactionReceipt(ctx, stx.Hash())
		if err == nil {
			return receipt, nil
		}
		time.Sleep(wait)
		wait = time.Duration(float64(wait) * 1.5)
	}

	return nil, fmt.Errorf("unable to get tx receipt")
}
*/

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
	addr := make([]byte, 20, 20)
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
		log.Error().Err(err).Msg("could not create storage key")
		return err
	}

	var accountInfo gstypes.AccountInfo
	ok, err := c.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil {
		log.Error().Err(err).Msg("could not load storage")
		return err
	}
	if !ok {
		err = fmt.Errorf("loaded storage is not okay")
		log.Error().Err(err).Msg("loaded storage is not okay")
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
					err := rl.Wait(ctx)
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
					break
				case loadTestModeDeploy:
					startReq, endReq, err = loadtestNotImplemented(ctx, c, myNonceValue)
					break
				case loadTestModeCall:
					startReq, endReq, err = loadtestNotImplemented(ctx, c, myNonceValue)
					break
				case loadTestModeFunction:
					startReq, endReq, err = loadtestNotImplemented(ctx, c, myNonceValue)
					break
				case loadTestModeInc:
					startReq, endReq, err = loadtestNotImplemented(ctx, c, myNonceValue)
					break
				case loadTestModeStore:
					startReq, endReq, err = loadtestAvailStore(ctx, c, myNonceValue, meta, genesisHash)
					break
				case loadTestModeLong:
					startReq, endReq, err = loadtestNotImplemented(ctx, c, myNonceValue)
					break
				default:
					log.Error().Str("mode", mode).Msg("We've arrived at a load test mode that we don't recognize")
				}
				recordSample(i, j, err, startReq, endReq)
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
		log.Error().Err(err).Msg("unable to create new multi address")
		return err
	}

	if *inputLoadTestParams.PrivateKey == codeQualityPrivateKey {
		// Avail keys can use the same seed but the way the key is derived is different
		*inputLoadTestParams.PrivateKey = codeQualitySeed
	}

	kp, err := gssignature.KeyringPairFromSecret(*inputLoadTestParams.PrivateKey, uint8(*inputLoadTestParams.ChainID))
	if err != nil {
		log.Error().Err(err).Msg("could not create key pair")
		return err
	}

	amt, err := hexToBigInt(*inputLoadTestParams.HexSendAmount)
	if err != nil {
		log.Error().Err(err).Msg("couldn't parse send amount")
		return err
	}

	rv, err := c.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		log.Error().Err(err).Msg("couldn't get runtime version")
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
		Era:                gstypes.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genesisHash,
		Nonce:              gstypes.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                gstypes.NewUCompactFromUInt(100),
		TransactionVersion: rv.TransactionVersion,
		AppID:              gstypes.U32(*ltp.AvailAppID),
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

	inputData := make([]byte, *ltp.ByteCount, *ltp.ByteCount)
	hexwordRead(inputData)

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
		Era:                gstypes.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genesisHash,
		Nonce:              gstypes.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                gstypes.NewUCompactFromUInt(100),
		TransactionVersion: rv.TransactionVersion,
		AppID:              gstypes.U32(*ltp.AvailAppID),
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
