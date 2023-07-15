package loadtest

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "embed"
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/maticnetwork/polygon-cli/metrics"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"

	"github.com/maticnetwork/polygon-cli/contracts"
	"github.com/rs/zerolog/log"
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
	loadTestModeERC20                = "2"
	loadTestModeERC721               = "7"
	loadTestModePrecompiledContracts = "p"
	loadTestModePrecompiledContract  = "a"

	codeQualitySeed       = "code code code code code code code code code code code quality"
	codeQualityPrivateKey = "42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa"
)

func initializeLoadTestParams(ctx context.Context, c *ethclient.Client) error {
	log.Info().Msg("Connecting with RPC endpoint to initialize load test parameters")
	gas, err := c.SuggestGasPrice(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Unable to retrieve gas price")
		return err
	}
	log.Trace().Interface("gasprice", gas).Msg("Retrieved current gas price")

	if !*inputLoadTestParams.LegacyTransactionMode {
		gasTipCap, _err := c.SuggestGasTipCap(ctx)
		if _err != nil {
			log.Error().Err(_err).Msg("Unable to retrieve gas tip cap")
			return _err
		}
		log.Trace().Interface("gastipcap", gasTipCap).Msg("Retrieved current gas tip cap")
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
	log.Trace().Uint64("chainID", chainID.Uint64()).Msg("Detected Chain ID")

	if *inputLoadTestParams.LegacyTransactionMode && *inputLoadTestParams.ForcePriorityGasPrice > 0 {
		log.Warn().Msg("Cannot set priority gas price in legacy mode")
	}
	if *inputLoadTestParams.ForceGasPrice < *inputLoadTestParams.ForcePriorityGasPrice {
		return errors.New("max priority fee per gas higher than max fee per gas")
	}

	if *inputLoadTestParams.AdaptiveRateLimit && *inputLoadTestParams.CallOnly {
		return errors.New("the adaptive rate limit is based on the pending transaction pool. It doesn't use this feature while also using call only")
	}

	inputLoadTestParams.ToETHAddress = &toAddr
	inputLoadTestParams.SendAmount = amt
	inputLoadTestParams.CurrentGasPrice = gas
	inputLoadTestParams.CurrentNonce = &nonce
	inputLoadTestParams.ECDSAPrivateKey = privateKey
	inputLoadTestParams.FromETHAddress = &ethAddress
	if *inputLoadTestParams.ChainID == 0 {
		*inputLoadTestParams.ChainID = chainID.Uint64()
	}
	inputLoadTestParams.CurrentBaseFee = header.BaseFee

	randSrc = rand.New(rand.NewSource(*inputLoadTestParams.Seed))

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
		log.Info().Uint("pending", ptc).Msg("There are still outstanding transactions. There might be issues restarting with the same sending key until those transactions clear")
	}
	log.Info().Msg("Finished")
	return nil
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
	// configureTransactOpts will set some paramters meant for load testing that could interfere with the deployment of our contracts
	tops.GasLimit = 0
	tops.GasPrice = nil
	tops.GasFeeCap = nil
	tops.GasTipCap = nil

	if err != nil {
		log.Error().Err(err).Msg("Unable create transaction signer")
		return err
	}
	cops := new(bind.CallOpts)

	// deploy and instantiate the load tester contract
	var ltAddr ethcommon.Address
	var ltContract *contracts.LoadTester
	if strings.ContainsAny(mode, "rcfispas") || *inputLoadTestParams.ForceContractDeploy {
		ltAddr, ltContract, err = getLoadTestContract(ctx, c, tops, cops)
		if err != nil {
			return err
		}
		log.Debug().Str("ltAddr", ltAddr.String()).Msg("Obtained load test contract address")
	}

	var erc20Addr ethcommon.Address
	var erc20Contract *contracts.ERC20
	if mode == loadTestModeERC20 || mode == loadTestModeRandom {
		erc20Addr, erc20Contract, err = getERC20Contract(ctx, c, tops, cops)
		if err != nil {
			return err
		}
		log.Debug().Str("erc20Addr", erc20Addr.String()).Msg("Obtained erc 20 contract address")
	}

	var erc721Addr ethcommon.Address
	var erc721Contract *contracts.ERC721
	if mode == loadTestModeERC721 || mode == loadTestModeRandom {
		erc721Addr, erc721Contract, err = getERC721Contract(ctx, c, tops, cops)
		if err != nil {
			return err
		}
		log.Debug().Str("erc721Addr", erc721Addr.String()).Msg("Obtained erc 721 contract address")
	}

	var currentNonceMutex sync.Mutex
	var i int64
	startBlockNumber, err := c.BlockNumber(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get current block number")
		return err
	}

	currentNonce, err := c.NonceAt(ctx, *ltp.FromETHAddress, new(big.Int).SetUint64(startBlockNumber))
	if err != nil {
		log.Error().Err(err).Msg("Unable to get account nonce")
		return err
	}

	currentNonceMutex.Lock()
	currentNonce += 1
	currentNonceMutex.Unlock()

	startNonce := currentNonce
	log.Debug().Uint64("currentNonce", currentNonce).Msg("Starting main load test loop")
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
					startReq, endReq, err = loadTestTransaction(ctx, c, myNonceValue)
				case loadTestModeDeploy:
					startReq, endReq, err = loadTestDeploy(ctx, c, myNonceValue)
				case loadTestModeFunction, loadTestModeCall:
					startReq, endReq, err = loadTestFunction(ctx, c, myNonceValue, ltContract)
				case loadTestModeInc:
					startReq, endReq, err = loadTestInc(ctx, c, myNonceValue, ltContract)
				case loadTestModeStore:
					startReq, endReq, err = loadTestStore(ctx, c, myNonceValue, ltContract)
				case loadTestModeERC20:
					startReq, endReq, err = loadTestERC20(ctx, c, myNonceValue, erc20Contract, ltAddr)
				case loadTestModeERC721:
					startReq, endReq, err = loadTestERC721(ctx, c, myNonceValue, erc721Contract, ltAddr)
				case loadTestModePrecompiledContract:
					startReq, endReq, err = loadTestCallPrecompiledContracts(ctx, c, myNonceValue, ltContract, true)
				case loadTestModePrecompiledContracts:
					startReq, endReq, err = loadTestCallPrecompiledContracts(ctx, c, myNonceValue, ltContract, false)
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
	log.Debug().Uint64("currentNonce", currentNonce).Msg("Finished main load test loop")
	log.Debug().Msg("Waiting for transactions to actually be mined")
	if *ltp.CallOnly {
		return nil
	}
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

func getLoadTestContract(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts, cops *bind.CallOpts) (ltAddr ethcommon.Address, ltContract *contracts.LoadTester, err error) {
	ltAddr = ethcommon.HexToAddress(*inputLoadTestParams.LtAddress)

	if *inputLoadTestParams.LtAddress == "" {
		ltAddr, _, _, err = contracts.DeployLoadTester(tops, c)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create the load testing contract. Do you have the right chain id? Do you have enough funds?")
			return
		}
	}
	log.Trace().Interface("contractaddress", ltAddr).Msg("Load test contract address")

	ltContract, err = contracts.NewLoadTester(ltAddr, c)
	if err != nil {
		log.Error().Err(err).Msg("Unable to instantiate new contract")
		return
	}
	err = blockUntilSuccessful(ctx, c, func() error {
		_, err = ltContract.GetCallCounter(cops)
		return err
	})

	return
}
func getERC20Contract(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts, cops *bind.CallOpts) (erc20Addr ethcommon.Address, erc20Contract *contracts.ERC20, err error) {
	erc20Addr = ethcommon.HexToAddress(*inputLoadTestParams.ERC20Address)

	if *inputLoadTestParams.ERC20Address == "" {
		erc20Addr, _, _, err = contracts.DeployERC20(tops, c)
		if err != nil {
			log.Error().Err(err).Msg("Unable to deploy ERC20 contract")
			return
		}
	}
	log.Trace().Interface("contractaddress", erc20Addr).Msg("ERC20 contract address")

	erc20Contract, err = contracts.NewERC20(erc20Addr, c)
	if err != nil {
		log.Error().Err(err).Msg("Unable to instantiate new erc20 contract")
		return
	}
	err = blockUntilSuccessful(ctx, c, func() error {
		_, err = erc20Contract.BalanceOf(cops, *inputLoadTestParams.FromETHAddress)
		return err
	})
	if err != nil {
		return
	}

	_, err = erc20Contract.Mint(tops, metrics.UnitMegaether)
	if err != nil {
		log.Error().Err(err).Msg("There was an error minting ERC20")
		return
	}

	err = blockUntilSuccessful(ctx, c, func() error {
		var balance *big.Int
		balance, err = erc20Contract.BalanceOf(cops, *inputLoadTestParams.FromETHAddress)
		if err != nil {
			return err
		}
		if balance.Uint64() == 0 {
			err = fmt.Errorf("ERC20 Balance is Zero")
			return err
		}
		return nil
	})

	return
}
func getERC721Contract(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts, cops *bind.CallOpts) (erc721Addr ethcommon.Address, erc721Contract *contracts.ERC721, err error) {
	erc721Addr = ethcommon.HexToAddress(*inputLoadTestParams.ERC721Address)

	if *inputLoadTestParams.ERC721Address == "" {
		erc721Addr, _, _, err = contracts.DeployERC721(tops, c)
		if err != nil {
			log.Error().Err(err).Msg("Unable to deploy ERC721 contract")
			return
		}
	}
	log.Trace().Interface("contractaddress", erc721Addr).Msg("ERC721 contract address")

	erc721Contract, err = contracts.NewERC721(erc721Addr, c)
	if err != nil {
		log.Error().Err(err).Msg("Unable to instantiate new erc20 contract")
		return
	}

	err = blockUntilSuccessful(ctx, c, func() error {
		_, err = erc721Contract.BalanceOf(cops, *inputLoadTestParams.FromETHAddress)
		return err
	})
	if err != nil {
		return
	}

	err = blockUntilSuccessful(ctx, c, func() error {
		_, err = erc721Contract.MintBatch(tops, *inputLoadTestParams.FromETHAddress, new(big.Int).SetUint64(1))
		return err
	})
	return
}

func blockUntilSuccessful(ctx context.Context, c *ethclient.Client, f func() error) error {
	numberOfBlocksToWaitFor := *inputLoadTestParams.ContractCallNumberOfBlocksToWaitFor
	blockInterval := *inputLoadTestParams.ContractCallBlockInterval
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
						log.Trace().Err(err).Dur("elapsedTimeSeconds", elapsed).Msg("Function executed successfully")
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

func loadTestTransaction(ctx context.Context, c *ethclient.Client, nonce uint64) (t1 time.Time, t2 time.Time, err error) {
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
	tops.GasLimit = uint64(21000)
	tops = configureTransactOpts(tops)

	var tx *ethtypes.Transaction
	if *ltp.LegacyTransactionMode {
		tx = ethtypes.NewTx(&ethtypes.LegacyTx{
			Nonce:    nonce,
			To:       to,
			Value:    amount,
			Gas:      tops.GasLimit,
			GasPrice: tops.GasPrice,
			Data:     nil,
		})
	} else {
		dynamicFeeTx := &ethtypes.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     nonce,
			To:        to,
			Gas:       tops.GasLimit,
			GasFeeCap: tops.GasFeeCap,
			GasTipCap: tops.GasTipCap,
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
	defer func() { t2 = time.Now() }()
	if *ltp.CallOnly {
		_, err = c.CallContract(ctx, txToCallMsg(stx), nil)
	} else {
		err = c.SendTransaction(ctx, stx)
	}
	return
}

// TODO - in the future it might be more interesting if this mode takes input or random contracts to be deployed
func loadTestDeploy(ctx context.Context, c *ethclient.Client, nonce uint64) (t1 time.Time, t2 time.Time, err error) {
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
	defer func() { t2 = time.Now() }()
	if *ltp.CallOnly {
		msg := transactOptsToCallMsg(tops)
		msg.Data = ethcommon.FromHex(contracts.LoadTesterMetaData.Bin)
		_, err = c.CallContract(ctx, msg, nil)
	} else {
		_, _, _, err = contracts.DeployLoadTester(tops, c)
	}
	return
}

// getCurrentLoadTestFunction is meant to handle the business logic
// around deciding which function to execute. When we're in function
// mode where the user has provided a specific function to execute, we
// should use that function. Otherwise, we'll select random functions.
func getCurrentLoadTestFunction() uint64 {
	if loadTestModeFunction == *inputLoadTestParams.Mode {
		return *inputLoadTestParams.Function
	}
	return contracts.GetRandomOPCode()
}
func loadTestFunction(ctx context.Context, c *ethclient.Client, nonce uint64, ltContract *contracts.LoadTester) (t1 time.Time, t2 time.Time, err error) {
	ltp := inputLoadTestParams

	chainID := new(big.Int).SetUint64(*ltp.ChainID)
	privateKey := ltp.ECDSAPrivateKey
	iterations := ltp.Iterations
	f := getCurrentLoadTestFunction()

	tops, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Error().Err(err).Msg("Unable create transaction signer")
		return
	}
	tops.Nonce = new(big.Int).SetUint64(nonce)
	tops = configureTransactOpts(tops)

	t1 = time.Now()
	defer func() { t2 = time.Now() }()
	if *ltp.CallOnly {
		tops.NoSend = true
		var tx *ethtypes.Transaction
		tx, err = contracts.CallLoadTestFunctionByOpCode(f, ltContract, tops, *iterations)
		if err != nil {
			return
		}
		msg := txToCallMsg(tx)
		_, err = c.CallContract(ctx, msg, nil)
	} else {
		_, err = contracts.CallLoadTestFunctionByOpCode(f, ltContract, tops, *iterations)
	}
	return
}

func loadTestCallPrecompiledContracts(ctx context.Context, c *ethclient.Client, nonce uint64, ltContract *contracts.LoadTester, useSelectedAddress bool) (t1 time.Time, t2 time.Time, err error) {
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
	defer func() { t2 = time.Now() }()
	if *ltp.CallOnly {
		tops.NoSend = true
		var tx *ethtypes.Transaction
		tx, err = contracts.CallPrecompiledContracts(f, ltContract, tops, *iterations, privateKey)
		if err != nil {
			return
		}
		msg := txToCallMsg(tx)
		_, err = c.CallContract(ctx, msg, nil)
	} else {
		_, err = contracts.CallPrecompiledContracts(f, ltContract, tops, *iterations, privateKey)
	}
	return
}

func loadTestInc(ctx context.Context, c *ethclient.Client, nonce uint64, ltContract *contracts.LoadTester) (t1 time.Time, t2 time.Time, err error) {
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
	defer func() { t2 = time.Now() }()
	if *ltp.CallOnly {
		tops.NoSend = true
		var tx *ethtypes.Transaction
		tx, err = ltContract.Inc(tops)
		if err != nil {
			return
		}
		msg := txToCallMsg(tx)
		_, err = c.CallContract(ctx, msg, nil)
	} else {
		_, err = ltContract.Inc(tops)
	}
	return
}

func loadTestStore(ctx context.Context, c *ethclient.Client, nonce uint64, ltContract *contracts.LoadTester) (t1 time.Time, t2 time.Time, err error) {
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
	defer func() { t2 = time.Now() }()
	if *ltp.CallOnly {
		tops.NoSend = true
		var tx *ethtypes.Transaction
		tx, err = ltContract.Store(tops, inputData)
		if err != nil {
			return
		}
		msg := txToCallMsg(tx)
		_, err = c.CallContract(ctx, msg, nil)
	} else {
		_, err = ltContract.Store(tops, inputData)
	}
	return
}

func loadTestERC20(ctx context.Context, c *ethclient.Client, nonce uint64, erc20Contract *contracts.ERC20, ltAddress ethcommon.Address) (t1 time.Time, t2 time.Time, err error) {
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
	defer func() { t2 = time.Now() }()
	if *ltp.CallOnly {
		tops.NoSend = true
		var tx *ethtypes.Transaction
		tx, err = erc20Contract.Transfer(tops, *to, amount)
		if err != nil {
			return
		}
		msg := txToCallMsg(tx)
		_, err = c.CallContract(ctx, msg, nil)
	} else {
		_, err = erc20Contract.Transfer(tops, *to, amount)
	}

	return
}

func loadTestERC721(ctx context.Context, c *ethclient.Client, nonce uint64, erc721Contract *contracts.ERC721, ltAddress ethcommon.Address) (t1 time.Time, t2 time.Time, err error) {
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
	defer func() { t2 = time.Now() }()
	if *ltp.CallOnly {
		tops.NoSend = true
		var tx *ethtypes.Transaction
		tx, err = erc721Contract.MintBatch(tops, *to, new(big.Int).SetUint64(*iterations))
		if err != nil {
			return
		}
		msg := txToCallMsg(tx)
		_, err = c.CallContract(ctx, msg, nil)
	} else {
		_, err = erc721Contract.MintBatch(tops, *to, new(big.Int).SetUint64(*iterations))
	}

	return
}
func loadTestNotImplemented(ctx context.Context, c *gsrpc.SubstrateAPI, nonce uint64) (t1 time.Time, t2 time.Time, err error) {
	t1 = time.Now()
	t2 = time.Now()
	err = fmt.Errorf("this method is not implemented")
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
	n, err := randSrc.Read(addr)
	if err != nil {
		log.Error().Err(err).Msg("There was an issue getting random bytes for the address")
	}
	if n != 20 {
		log.Error().Int("n", n).Msg("Somehow we didn't read 20 random bytes")
	}
	realAddr := ethcommon.BytesToAddress(addr)
	return &realAddr
}

func configureTransactOpts(tops *bind.TransactOpts) *bind.TransactOpts {
	ltp := inputLoadTestParams

	if ltp.ForceGasPrice != nil && *ltp.ForceGasPrice != 0 {
		tops.GasPrice = big.NewInt(0).SetUint64(*ltp.ForceGasPrice)
	}
	if ltp.ForceGasLimit != nil && *ltp.ForceGasLimit != 0 {
		tops.GasLimit = *ltp.ForceGasLimit
	}

	// if we're in legacy mode, there's no point doing anything else in this function
	if *ltp.LegacyTransactionMode {
		return tops
	}

	if ltp.ForcePriorityGasPrice != nil && *ltp.ForcePriorityGasPrice != 0 {
		tops.GasTipCap = big.NewInt(0).SetUint64(*ltp.ForcePriorityGasPrice)
	}

	if ltp.CurrentBaseFee == nil {
		log.Fatal().Msg("EIP-1559 not activated. Please use --legacy")
	}

	tops.GasPrice = nil
	tops.GasFeeCap = big.NewInt(0).Add(ltp.CurrentBaseFee, ltp.CurrentGasTipCap)

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

func transactOptsToCallMsg(tops *bind.TransactOpts) ethereum.CallMsg {
	cm := new(ethereum.CallMsg)
	cm.From = *inputLoadTestParams.FromETHAddress

	cm.Gas = tops.GasLimit
	cm.GasPrice = tops.GasPrice
	cm.GasFeeCap = tops.GasFeeCap
	cm.GasTipCap = tops.GasTipCap
	cm.Value = tops.Value
	return *cm
}

func txToCallMsg(tx *ethtypes.Transaction) ethereum.CallMsg {
	cm := new(ethereum.CallMsg)
	cm.From = *inputLoadTestParams.FromETHAddress
	cm.To = tx.To()
	cm.Gas = tx.Gas()
	cm.GasPrice = tx.GasPrice()
	cm.GasFeeCap = tx.GasFeeCap()
	cm.GasTipCap = tx.GasTipCap()
	cm.Value = tx.Value()
	cm.Data = tx.Data()

	cm.AccessList = tx.AccessList()
	return *cm
}
