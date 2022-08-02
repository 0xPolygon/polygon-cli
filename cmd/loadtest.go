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
	"math/big"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"

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

var (
	inputLoadTestParams loadTestParams
	loadTestResults     []loadTestSample
)

// loadtestCmd represents the loadtest command
var loadtestCmd = &cobra.Command{
	Use:   "loadtest [options] rpc-endpoint",
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
			return fmt.Errorf("Expected exactly one argument")
		}
		url, err := url.Parse(args[0])
		if err != nil {
			log.Error().Err(err).Msg("Unable to parse url input error")
			return err
		}
		if url.Scheme != "http" && url.Scheme != "https" {
			return fmt.Errorf("The scheme %s is not supported", url.Scheme)
		}
		inputLoadTestParams.URL = url
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
		Requests      *int64
		Concurrency   *int64
		TimeLimit     *int64
		Verbosity     *int64
		PrettyLogs    *bool
		URL           *url.URL
		ChainID       *uint64
		PrivateKey    *string
		ToAddress     *string
		HexSendAmount *string
		RateLimit     *float64

		// Computed
		CurrentGas      *big.Int
		CurrentNonce    *uint64
		ECDSAPrivateKey *ecdsa.PrivateKey
		FromETHAddress  *ethcommon.Address
		ToETHAddress    *ethcommon.Address
		SendAmount      *big.Int
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
	ltp.PrivateKey = loadtestCmd.PersistentFlags().String("private-key", "42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa", "The hex encoded private key that we'll use to sending transactions")
	ltp.ChainID = loadtestCmd.PersistentFlags().Uint64("chain-id", 1256, "The chain id for the transactions that we're going to send")
	ltp.ToAddress = loadtestCmd.PersistentFlags().String("to-address", "0xDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEF", "The address that we're going to send to")
	ltp.HexSendAmount = loadtestCmd.PersistentFlags().String("send-amount", "0x38D7EA4C68000", "The amount of wei that we'll send every transaction")
	ltp.RateLimit = loadtestCmd.PersistentFlags().Float64("rate-limit", 4, "An overall limit to the number of requests per second. Give a number less than zero to remove this limit all together")

	inputLoadTestParams = *ltp

	// TODO batch size
	// TODO Compression
	// TODO array of RPC endpoints to round robin?
}

func initalizeLoadTestParams(ctx context.Context, c *ethclient.Client) error {
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

	return nil
}

func hexToBigInt(raw any) (bi *big.Int, err error) {
	bi = big.NewInt(0)
	hexString, ok := raw.(string)
	if !ok {
		err = fmt.Errorf("Could not assert value %v as a string", raw)
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

func hexToUint64(raw any) (uint64, error) {
	hexString, ok := raw.(string)
	if !ok {
		return 0, fmt.Errorf("Could not assert %v as a string", hexString)
	}

	hexString = strings.Replace(hexString, "0x", "", -1)
	if len(hexString)%2 != 0 {
		log.Trace().Str("original", hexString).Msg("Hex of odd length")
		hexString = "0" + hexString
	}

	result, err := strconv.ParseUint(hexString, 16, 64)
	if err != nil {
		log.Error().Err(err).Msg("Unable to decode hex string")
		return 0, err
	}
	return uint64(result), nil
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

	err = initalizeLoadTestParams(ctx, ec)
	if err != nil {
		return err
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	loadTestResults = make([]loadTestSample, 0)
	errCh := make(chan error)
	go func() {
		errCh <- mainLoop(ctx, ec)
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
	currentGas := ltp.CurrentGas
	sendTo := ltp.ToETHAddress
	sendAmt := ltp.SendAmount
	chainID := new(big.Int).SetUint64(*ltp.ChainID)
	privateKey := ltp.ECDSAPrivateKey

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

	addr, _, _, err := contracts.DeployLoadTester(tops, c)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create the load testing contract")
		return err
	}
	log.Trace().Interface("contractaddress", addr).Msg("Load test contract address")
	// bump the nonce since deploying a contract should cause it to increase
	currentNonce = currentNonce + 1

	// block while the contract is pending
	bn, err := c.BlockNumber(ctx)
	if err != nil {
		return err
	}
	for {
		pendCount, err := c.PendingTransactionCount(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Unable to get pending transaction count")
			return err
		}
		bn2, err := c.BlockNumber(ctx)
		if err != nil {
			return err
		}

		if pendCount < 1 && bn != bn2 {
			break
		}
		time.Sleep(time.Second)
	}
	ltContract, err := contracts.NewLoadTester(addr, c)
	if err != nil {
		log.Error().Err(err).Msg("Unable to instantiate new contract")
		return err
	}

	ltCounter, err := ltContract.GetCallCounter(cops)
	if err != nil {
		log.Error().Err(err).Msg("Unable to call test contract")
		return err
	}
	log.Trace().Interface("counter", ltCounter).Msg("Number of contract calls")

	var currentNonceMutex sync.Mutex
	var i int64

	var wg sync.WaitGroup
	for i = 0; i < routines; i = i + 1 {
		log.Trace().Int64("routine", i).Msg("Starting Thread")
		wg.Add(1)
		go func(i int64) {
			var j int64
			for j = 0; j < requests; j = j + 1 {

				if rl != nil {
					err := rl.Wait(ctx)
					if err != nil {
						log.Error().Err(err).Msg("Encountered a rate limiting error")
					}
				}

				// TODO support different modes (transfer, contract deploy, contract call, multi, etc)
				currentNonceMutex.Lock()
				myNonceValue := currentNonce
				currentNonce = currentNonce + 1
				currentNonceMutex.Unlock()
				gasLimit := uint64(21000)

				tx := ethtypes.NewTransaction(myNonceValue, *sendTo, sendAmt, gasLimit, currentGas, nil)
				stx, err := ethtypes.SignTx(tx, ethtypes.NewEIP155Signer(chainID), privateKey)
				if err != nil {
					log.Error().Err(err).Msg("Unable to sign transaction")
					return
				}

				startReq := time.Now()
				err = c.SendTransaction(ctx, stx)
				endReq := time.Now()
				recordSample(i, j, err, startReq, endReq)
				if err != nil {
					log.Trace().Err(err).Msg("Recorded an error while sending transactions")
				}

				log.Trace().Int64("routine", i).Int64("request", j).Msg("Request")
			}
			wg.Done()
		}(i)

	}
	log.Trace().Msg("Finished starting go routines. Waiting..")
	wg.Wait()
	return nil
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

	return nil, fmt.Errorf("Unable to get tx receipt")
}
