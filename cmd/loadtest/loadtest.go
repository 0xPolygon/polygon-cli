package loadtest

import (
	"context"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"math/rand"

	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/holiman/uint256"

	"github.com/maticnetwork/polygon-cli/bindings/tester"
	"github.com/maticnetwork/polygon-cli/bindings/tokens"
	uniswapv3loadtest "github.com/maticnetwork/polygon-cli/cmd/loadtest/uniswapv3"

	"github.com/maticnetwork/polygon-cli/abi"
	"github.com/maticnetwork/polygon-cli/rpctypes"
	"github.com/maticnetwork/polygon-cli/util"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

//go:generate stringer -type=loadTestMode
type (
	loadTestMode int
)

const (
	// these constants are stringered. If you add a new constant it fill fail to compile until you regenerate the strings. There are two steps needed.
	// 1. Install stringer with something like `go install golang.org/x/tools/cmd/stringer`
	// 2. now that its installed (make sure your GOBIN is on the PATH) you can run `go generate github.com/maticnetwork/polygon-cli/cmd/loadtest`
	loadTestModeTransaction loadTestMode = iota
	loadTestModeDeploy
	loadTestModeCall
	loadTestModeFunction
	loadTestModeInc
	loadTestModeStore
	loadTestModeERC20
	loadTestModeERC721
	loadTestModePrecompiledContracts
	loadTestModePrecompiledContract

	// All the modes AFTER random mode will not be used when mode random is selected
	loadTestModeRandom
	loadTestModeRecall
	loadTestModeRPC
	loadTestModeContractCall
	loadTestModeInscription
	loadTestModeUniswapV3
	loadTestModeBlob

	codeQualitySeed       = "code code code code code code code code code code code quality"
	codeQualityPrivateKey = "42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa"
)

func characterToLoadTestMode(mode string) (loadTestMode, error) {
	switch mode {
	case "t", "transaction":
		return loadTestModeTransaction, nil
	case "d", "deploy":
		return loadTestModeDeploy, nil
	case "c", "call":
		return loadTestModeCall, nil
	case "f", "function":
		return loadTestModeFunction, nil
	case "i", "inc", "increment":
		return loadTestModeInc, nil
	case "r", "random":
		return loadTestModeRandom, nil
	case "s", "store":
		return loadTestModeStore, nil
	case "2", "erc20":
		return loadTestModeERC20, nil
	case "7", "erc721":
		return loadTestModeERC721, nil
	case "p", "precompile":
		return loadTestModePrecompiledContract, nil
	case "P", "precompiles":
		return loadTestModePrecompiledContracts, nil
	case "R", "recall":
		return loadTestModeRecall, nil
	case "v3", "uniswapv3":
		return loadTestModeUniswapV3, nil
	case "rpc":
		return loadTestModeRPC, nil
	case "cc", "contract-call":
		return loadTestModeContractCall, nil
	case "inscription":
		return loadTestModeInscription, nil
	case "blob":
		return loadTestModeBlob, nil
	default:
		return 0, fmt.Errorf("unrecognized load test mode: %s", mode)
	}
}

func getRandomMode() loadTestMode {
	maxMode := int(loadTestModeRandom)
	return loadTestMode(randSrc.Intn(maxMode))
}

func modeRequiresLoadTestContract(m loadTestMode) bool {
	if m == loadTestModeCall ||
		m == loadTestModeFunction ||
		m == loadTestModeInc ||
		m == loadTestModeRandom ||
		m == loadTestModeStore ||
		m == loadTestModePrecompiledContract ||
		m == loadTestModePrecompiledContracts {
		return true
	}
	return false
}
func anyModeRequiresLoadTestContract(modes []loadTestMode) bool {
	for _, m := range modes {
		if modeRequiresLoadTestContract(m) {
			return true
		}
	}
	return false
}
func hasMode(mode loadTestMode, modes []loadTestMode) bool {
	for _, m := range modes {
		if m == mode {
			return true
		}
	}
	return false
}

func hasUniqueModes(modes []loadTestMode) bool {
	seen := make(map[loadTestMode]bool, len(modes))
	for _, m := range modes {
		if !seen[m] {
			seen[m] = true
		} else {
			return false
		}
	}
	return true
}

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

	trimmedHexPrivateKey := strings.TrimPrefix(*inputLoadTestParams.PrivateKey, "0x")
	privateKey, err := ethcrypto.HexToECDSA(trimmedHexPrivateKey)
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

	amt := util.EthToWei(*inputLoadTestParams.EthAmountInWei)

	header, err := c.HeaderByNumber(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg("Unable to get header")
		return err
	}
	if header.BaseFee != nil {
		inputLoadTestParams.ChainSupportBaseFee = true
		log.Debug().Msg("Eip-1559 support detected")
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

	contractAddr := ethcommon.HexToAddress(*inputLoadTestParams.ContractAddress)
	inputLoadTestParams.ContractETHAddress = &contractAddr

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

	modes := *inputLoadTestParams.Modes
	if len(modes) == 0 {
		return errors.New("expected at least one mode")
	}

	inputLoadTestParams.ParsedModes = make([]loadTestMode, 0)
	for _, m := range modes {
		parsedMode, err := characterToLoadTestMode(m)
		if err != nil {
			return err
		}
		inputLoadTestParams.ParsedModes = append(inputLoadTestParams.ParsedModes, parsedMode)
	}

	// Logic checking input parameters for specific conditions such as multiple inputs.
	if len(modes) > 1 {
		inputLoadTestParams.MultiMode = true
		if !hasUniqueModes(inputLoadTestParams.ParsedModes) {
			return errors.New("Duplicate modes detected, check input modes for duplicates")
		}
	} else {
		inputLoadTestParams.MultiMode = false
		inputLoadTestParams.Mode, _ = characterToLoadTestMode((*inputLoadTestParams.Modes)[0])
	}
	if hasMode(loadTestModeRandom, inputLoadTestParams.ParsedModes) && inputLoadTestParams.MultiMode {
		return errors.New("random mode can't be used in combinations with any other modes")
	}
	if hasMode(loadTestModeRPC, inputLoadTestParams.ParsedModes) && inputLoadTestParams.MultiMode && !*inputLoadTestParams.CallOnly {
		return errors.New("rpc mode must be called with call-only when multiple modes are used")
	} else if hasMode(loadTestModeRPC, inputLoadTestParams.ParsedModes) {
		log.Trace().Msg("Setting call only mode since we're doing RPC testing")
		*inputLoadTestParams.CallOnly = true
	}
	if hasMode(loadTestModeContractCall, inputLoadTestParams.ParsedModes) && (*inputLoadTestParams.ContractAddress == "" || (*inputLoadTestParams.ContractCallData == "" && *inputLoadTestParams.ContractCallFunctionSignature == "")) {
		return errors.New("`--contract-call` requires both a `--contract-address` and calldata, either with `--calldata` or `--function-signature --function-arg` flags.")
	}
	if *inputLoadTestParams.CallOnly && *inputLoadTestParams.AdaptiveRateLimit {
		return errors.New("using call only with adaptive rate limit doesn't make sense")
	}
	if hasMode(loadTestModeBlob, inputLoadTestParams.ParsedModes) && inputLoadTestParams.MultiMode {
		return errors.New("Blob mode should only be used by itself. Blob mode will take significantly longer than other transactions to finalize, and the address will be reserved, preventing other transactions form being made.")
	}

	randSrc = rand.New(rand.NewSource(*inputLoadTestParams.Seed))

	return nil
}

func initNonce(ctx context.Context, c *ethclient.Client, rpc *ethrpc.Client) error {
	var err error
	startBlockNumber, err = c.BlockNumber(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get current block number")
		return err
	}

	// Get pending nonce to be prevent nonce collision (if tx from same sender is already present)
	currentNonce, err = c.PendingNonceAt(ctx, *inputLoadTestParams.FromETHAddress)
	startNonce = currentNonce

	if err != nil {
		log.Error().Err(err).Msg("Unable to get account nonce")
		return err
	}

	return nil
}

func completeLoadTest(ctx context.Context, c *ethclient.Client, rpc *ethrpc.Client) error {
	log.Debug().Uint64("startNonce", startNonce).Uint64("lastNonce", currentNonce).Msg("Finished main load test loop")
	if *inputLoadTestParams.SendOnly {
		log.Info().Uint64("transactionsSent", currentNonce-startNonce).Msg("SendOnly mode enabled - skipping wait period and summarization")
		return nil
	}
	log.Debug().Msg("Waiting for remaining transactions to be completed and mined")

	var err error
	finalBlockNumber, err = waitForFinalBlock(ctx, c, rpc, startBlockNumber, startNonce, currentNonce)
	if err != nil {
		log.Error().Err(err).Msg("There was an issue waiting for all transactions to be mined")
	}
	if len(loadTestResults) == 0 {
		return errors.New("no transactions observed")
	}

	startTime := loadTestResults[0].RequestTime
	endTime := time.Now()
	log.Debug().Uint64("currentNonce", currentNonce).Uint64("final block number", finalBlockNumber).Msg("Got final block number")

	if *inputLoadTestParams.CallOnly {
		log.Info().Msg("CallOnly mode enabled - blocks aren't mined")
		lightSummary(loadTestResults, startTime, endTime, rl)
		return nil
	}

	if *inputLoadTestParams.ShouldProduceSummary {
		err = summarizeTransactions(ctx, c, rpc, startBlockNumber, startNonce, finalBlockNumber, currentNonce)
		if err != nil {
			log.Error().Err(err).Msg("There was an issue creating the load test summary")
		}
	}
	lightSummary(loadTestResults, startTime, endTime, rl)

	return nil
}

// runLoadTest initiates and runs the entire load test process, including initialization,
// the main load test loop, and the completion steps. It takes a context for cancellation signals.
// The function returns an error if there are issues during the load test process.
func runLoadTest(ctx context.Context) error {
	log.Info().Msg("Starting Load Test")

	// Configure the overall time limit for the load test.
	timeLimit := *inputLoadTestParams.TimeLimit
	var overallTimer *time.Timer
	if timeLimit > 0 {
		overallTimer = time.NewTimer(time.Duration(timeLimit) * time.Second)
	} else {
		overallTimer = new(time.Timer)
	}

	// Dial the Ethereum RPC server.
	rpc, err := ethrpc.DialContext(ctx, *inputLoadTestParams.RPCUrl)
	if err != nil {
		log.Error().Err(err).Msg("Unable to dial rpc")
		return err
	}
	defer rpc.Close()
	rpc.SetHeader("Accept-Encoding", "identity")
	ec := ethclient.NewClient(rpc)

	// Define the main loop function.
	// Make sure to define any logic associated to the load test (initialization, main load test loop
	// or completion steps) in this function in order to handle cancellation signals properly.
	loopFunc := func() error {
		if err = initializeLoadTestParams(ctx, ec); err != nil {
			log.Error().Err(err).Msg("Error initializing load test parameters")
			return err
		}

		if err = mainLoop(ctx, ec, rpc); err != nil {
			log.Error().Err(err).Msg("Error during the main load test loop")
			return err
		}

		if err = completeLoadTest(ctx, ec, rpc); err != nil {
			log.Error().Err(err).Msg("Encountered error while wrapping up loadtest")
		}
		return nil
	}

	// Set up signal handling for interrupts.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	// Initialize channels for handling errors and running the main loop.
	loadTestResults = make([]loadTestSample, 0)
	errCh := make(chan error)
	go func() {
		errCh <- loopFunc()
	}()

	// Wait for the load test to complete, either due to time limit, interrupt signal, or completion.
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
	log.Info().Msg("Finished")
	return nil
}

func updateRateLimit(ctx context.Context, rl *rate.Limiter, rpc *ethrpc.Client, steadyStateQueueSize uint64, rateLimitIncrement uint64, cycleDuration time.Duration, backoff float64) {
	ticker := time.NewTicker(cycleDuration)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			pendingTx, queuedTx, err := util.GetTxPoolStatus(rpc)
			if err != nil {
				log.Error().Err(err).Msg("Error getting txpool size")
				return
			}

			txPoolSize := pendingTx + queuedTx
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
	mode := ltp.Mode
	steadyStateTxPoolSize := *ltp.SteadyStateTxPoolSize
	adaptiveRateLimitIncrement := *ltp.AdaptiveRateLimitIncrement
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
	// configureTransactOpts will set some parameters meant for load testing that could interfere with the deployment of our contracts
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
	var ltContract *tester.LoadTester
	if anyModeRequiresLoadTestContract(ltp.ParsedModes) || *inputLoadTestParams.ForceContractDeploy {
		ltAddr, ltContract, err = getLoadTestContract(ctx, c, tops, cops)
		if err != nil {
			return err
		}
		log.Debug().Str("ltAddr", ltAddr.String()).Msg("Obtained load test contract address")
	}

	var erc20Addr ethcommon.Address
	var erc20Contract *tokens.ERC20
	if hasMode(loadTestModeERC20, ltp.ParsedModes) || hasMode(loadTestModeRandom, ltp.ParsedModes) || hasMode(loadTestModeRPC, ltp.ParsedModes) {
		erc20Addr, erc20Contract, err = getERC20Contract(ctx, c, tops, cops)
		if err != nil {
			return err
		}
		log.Debug().Str("erc20Addr", erc20Addr.String()).Msg("Obtained erc 20 contract address")
	}

	var erc721Addr ethcommon.Address
	var erc721Contract *tokens.ERC721
	if hasMode(loadTestModeERC721, ltp.ParsedModes) || hasMode(loadTestModeRandom, ltp.ParsedModes) || hasMode(loadTestModeRPC, ltp.ParsedModes) {
		erc721Addr, erc721Contract, err = getERC721Contract(ctx, c, tops, cops)
		if err != nil {
			return err
		}
		log.Debug().Str("erc721Addr", erc721Addr.String()).Msg("Obtained erc 721 contract address")
	}

	var recallTransactions []rpctypes.PolyTransaction
	if hasMode(loadTestModeRecall, ltp.ParsedModes) {
		recallTransactions, err = getRecallTransactions(ctx, c, rpc)
		if err != nil {
			return err
		}
		if len(recallTransactions) == 0 {
			return errors.New("we weren't able to fetch any recall transactions")
		}
		log.Debug().Int("txs", len(recallTransactions)).Msg("Retrieved transactions for total recall")
	}

	var indexedActivity *IndexedActivity
	if hasMode(loadTestModeRPC, ltp.ParsedModes) {
		indexedActivity, err = getIndexedRecentActivity(ctx, c, rpc)
		if err != nil {
			return err
		}
		if len(indexedActivity.ERC20Addresses) == 0 {
			indexedActivity.ERC20Addresses = append(indexedActivity.ERC20Addresses, erc20Addr.String())
		}

		if len(indexedActivity.ERC721Addresses) == 0 {
			indexedActivity.ERC721Addresses = append(indexedActivity.ERC721Addresses, erc721Addr.String())
		}

		log.Debug().
			Int("transactions", len(indexedActivity.TransactionIDs)).
			Int("blocks", len(indexedActivity.BlockNumbers)).
			Int("addresses", len(indexedActivity.Addresses)).
			Int("erc20s", len(indexedActivity.ERC20Addresses)).
			Int("erc721", len(indexedActivity.ERC721Addresses)).
			Int("contracts", len(indexedActivity.Contracts)).
			Msg("Retrieved recent indexed activity")
	}

	var uniswapV3Config uniswapv3loadtest.UniswapV3Config
	var poolConfig uniswapv3loadtest.PoolConfig
	if hasMode(loadTestModeUniswapV3, ltp.ParsedModes) {
		uniswapAddresses := uniswapv3loadtest.UniswapV3Addresses{
			FactoryV3:                          ethcommon.HexToAddress(*uniswapv3LoadTestParams.UniswapFactoryV3),
			Multicall:                          ethcommon.HexToAddress(*uniswapv3LoadTestParams.UniswapMulticall),
			ProxyAdmin:                         ethcommon.HexToAddress(*uniswapv3LoadTestParams.UniswapProxyAdmin),
			TickLens:                           ethcommon.HexToAddress(*uniswapv3LoadTestParams.UniswapTickLens),
			NFTDescriptorLib:                   ethcommon.HexToAddress(*uniswapv3LoadTestParams.UniswapNFTLibDescriptor),
			NonfungibleTokenPositionDescriptor: ethcommon.HexToAddress(*uniswapv3LoadTestParams.UniswapNonfungibleTokenPositionDescriptor),
			TransparentUpgradeableProxy:        ethcommon.HexToAddress(*uniswapv3LoadTestParams.UniswapUpgradeableProxy),
			NonfungiblePositionManager:         ethcommon.HexToAddress(*uniswapv3LoadTestParams.UniswapNonfungiblePositionManager),
			Migrator:                           ethcommon.HexToAddress(*uniswapv3LoadTestParams.UniswapMigrator),
			Staker:                             ethcommon.HexToAddress(*uniswapv3LoadTestParams.UniswapStaker),
			QuoterV2:                           ethcommon.HexToAddress(*uniswapv3LoadTestParams.UniswapQuoterV2),
			SwapRouter02:                       ethcommon.HexToAddress(*uniswapv3LoadTestParams.UniswapSwapRouter),
			WETH9:                              ethcommon.HexToAddress(*uniswapv3LoadTestParams.WETH9),
		}
		uniswapV3Config, poolConfig, err = initUniswapV3Loadtest(ctx, c, tops, cops, uniswapAddresses, *ltp.FromETHAddress)
		if err != nil {
			return err
		}
	}

	var i int64
	err = initNonce(ctx, c, rpc)
	if err != nil {
		return err
	}
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
			var tErr error

			for j = 0; j < requests; j = j + 1 {
				if rl != nil {
					tErr = rl.Wait(ctx)
					if tErr != nil {
						log.Error().Err(tErr).Msg("Encountered a rate limiting error")
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
				if ltp.MultiMode {
					localMode = ltp.ParsedModes[int(i+j)%(len(ltp.ParsedModes))]
				}
				// if we're doing random, we'll just pick one based on the current index
				if localMode == loadTestModeRandom {
					localMode = getRandomMode()
				}
				switch localMode {
				case loadTestModeTransaction:
					startReq, endReq, tErr = loadTestTransaction(ctx, c, myNonceValue)
				case loadTestModeDeploy:
					startReq, endReq, tErr = loadTestDeploy(ctx, c, myNonceValue)
				case loadTestModeFunction, loadTestModeCall:
					startReq, endReq, tErr = loadTestFunction(ctx, c, myNonceValue, ltContract)
				case loadTestModeInc:
					startReq, endReq, tErr = loadTestInc(ctx, c, myNonceValue, ltContract)
				case loadTestModeStore:
					startReq, endReq, tErr = loadTestStore(ctx, c, myNonceValue, ltContract)
				case loadTestModeERC20:
					startReq, endReq, tErr = loadTestERC20(ctx, c, myNonceValue, erc20Contract, ltAddr)
				case loadTestModeERC721:
					startReq, endReq, tErr = loadTestERC721(ctx, c, myNonceValue, erc721Contract, ltAddr)
				case loadTestModePrecompiledContract:
					startReq, endReq, tErr = loadTestCallPrecompiledContracts(ctx, c, myNonceValue, ltContract, true)
				case loadTestModePrecompiledContracts:
					startReq, endReq, tErr = loadTestCallPrecompiledContracts(ctx, c, myNonceValue, ltContract, false)
				case loadTestModeRecall:
					startReq, endReq, tErr = loadTestRecall(ctx, c, myNonceValue, recallTransactions[int(currentNonce)%len(recallTransactions)])
				case loadTestModeUniswapV3:
					swapAmountIn := big.NewInt(int64(*uniswapv3LoadTestParams.SwapAmountInput))
					startReq, endReq, tErr = runUniswapV3Loadtest(ctx, c, myNonceValue, uniswapV3Config, poolConfig, swapAmountIn)
				case loadTestModeRPC:
					startReq, endReq, tErr = loadTestRPC(ctx, c, myNonceValue, indexedActivity)
				case loadTestModeContractCall:
					startReq, endReq, tErr = loadTestContractCall(ctx, c, myNonceValue)
				case loadTestModeInscription:
					startReq, endReq, tErr = loadTestInscription(ctx, c, myNonceValue)
				case loadTestModeBlob:
					startReq, endReq, tErr = loadTestBlob(ctx, c, myNonceValue)
				default:
					log.Error().Str("mode", mode.String()).Msg("We've arrived at a load test mode that we don't recognize")
				}
				recordSample(i, j, tErr, startReq, endReq, myNonceValue)
				if tErr != nil {
					log.Error().Err(tErr).Uint64("nonce", myNonceValue).Int64("request time", endReq.Sub(startReq).Milliseconds()).Msg("Recorded an error while sending transactions")
					// The nonce is used to index the recalled transactions in call-only mode. We don't want to retry a transaction if it legit failed on the chain
					if !*ltp.CallOnly {
						retryForNonce = true
					}
					if strings.Contains(tErr.Error(), "replacement transaction underpriced") && retryForNonce {
						retryForNonce = false
					}
					if strings.Contains(tErr.Error(), "transaction underpriced") && retryForNonce {
						retryForNonce = false
					}
					if strings.Contains(tErr.Error(), "nonce too low") && retryForNonce {
						retryForNonce = false
					}
				}

				log.Trace().Uint64("nonce", myNonceValue).Int64("routine", i).Str("mode", localMode.String()).Int64("request", j).Msg("Request")
			}
			wg.Done()
		}(i)
	}
	log.Trace().Msg("Finished starting go routines. Waiting..")
	wg.Wait()
	cancel()
	if *ltp.CallOnly {
		return nil
	}

	return nil
}

func getLoadTestContract(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts, cops *bind.CallOpts) (ltAddr ethcommon.Address, ltContract *tester.LoadTester, err error) {
	ltAddr = ethcommon.HexToAddress(*inputLoadTestParams.LtAddress)

	if *inputLoadTestParams.LtAddress == "" {
		ltAddr, _, _, err = tester.DeployLoadTester(tops, c)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create the load testing contract. Do you have the right chain id? Do you have enough funds?")
			return
		}
	}
	log.Trace().Interface("contractaddress", ltAddr).Msg("Load test contract address")

	ltContract, err = tester.NewLoadTester(ltAddr, c)
	if err != nil {
		log.Error().Err(err).Msg("Unable to instantiate new contract")
		return
	}
	err = util.BlockUntilSuccessful(ctx, c, func() error {
		_, err = ltContract.GetCallCounter(cops)
		return err
	})

	return
}
func getERC20Contract(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts, cops *bind.CallOpts) (erc20Addr ethcommon.Address, erc20Contract *tokens.ERC20, err error) {
	erc20Addr = ethcommon.HexToAddress(*inputLoadTestParams.ERC20Address)
	if *inputLoadTestParams.ERC20Address == "" {
		erc20Addr, _, _, err = tokens.DeployERC20(tops, c)
		if err != nil {
			log.Error().Err(err).Msg("Unable to deploy ERC20 contract")
			return
		}
		// Tokens already minted and sent to the address of the deployer.
	}
	log.Trace().Interface("contractaddress", erc20Addr).Msg("ERC20 contract address")

	erc20Contract, err = tokens.NewERC20(erc20Addr, c)
	if err != nil {
		log.Error().Err(err).Msg("Unable to instantiate new erc20 contract")
		return
	}

	err = util.BlockUntilSuccessful(ctx, c, func() error {
		var balance *big.Int
		balance, err = erc20Contract.BalanceOf(cops, *inputLoadTestParams.FromETHAddress)
		if err != nil {
			return err
		}
		if balance.Uint64() == 0 {
			err = errors.New("ERC20 Balance is Zero")
			return err
		}
		return nil
	})

	return
}
func getERC721Contract(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts, cops *bind.CallOpts) (erc721Addr ethcommon.Address, erc721Contract *tokens.ERC721, err error) {
	erc721Addr = ethcommon.HexToAddress(*inputLoadTestParams.ERC721Address)
	shouldMint := true
	if *inputLoadTestParams.ERC721Address == "" {
		erc721Addr, _, _, err = tokens.DeployERC721(tops, c)
		if err != nil {
			log.Error().Err(err).Msg("Unable to deploy ERC721 contract")
			return
		}
		shouldMint = false
	}
	log.Trace().Interface("contractaddress", erc721Addr).Msg("ERC721 contract address")

	erc721Contract, err = tokens.NewERC721(erc721Addr, c)
	if err != nil {
		log.Error().Err(err).Msg("Unable to instantiate new erc721 contract")
		return
	}

	err = util.BlockUntilSuccessful(ctx, c, func() error {
		_, err = erc721Contract.BalanceOf(cops, *inputLoadTestParams.FromETHAddress)
		return err
	})
	if err != nil {
		return
	}
	if !shouldMint {
		return
	}

	err = util.BlockUntilSuccessful(ctx, c, func() error {
		_, err = erc721Contract.MintBatch(tops, *inputLoadTestParams.FromETHAddress, new(big.Int).SetUint64(1))
		return err
	})
	return
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
	gasPrice, gasTipCap := getSuggestedGasPrices(ctx, c)

	var tx *ethtypes.Transaction
	if *ltp.LegacyTransactionMode {
		tx = ethtypes.NewTx(&ethtypes.LegacyTx{
			Nonce:    nonce,
			To:       to,
			Value:    amount,
			Gas:      tops.GasLimit,
			GasPrice: gasPrice,
			Data:     nil,
		})
	} else {
		dynamicFeeTx := &ethtypes.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     nonce,
			To:        to,
			Gas:       tops.GasLimit,
			GasFeeCap: gasPrice,
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
	defer func() { t2 = time.Now() }()
	if *ltp.CallOnly {
		_, err = c.CallContract(ctx, txToCallMsg(stx), nil)
	} else {
		err = c.SendTransaction(ctx, stx)
	}
	return
}

var (
	cachedBlockNumber  uint64
	cachedGasPriceLock sync.Mutex
	cachedGasPrice     *big.Int
	cachedGasTipCap    *big.Int
)

func getSuggestedGasPrices(ctx context.Context, c *ethclient.Client) (*big.Int, *big.Int) {
	// this should be one of the fastest RPC calls, so hopefully there isn't too much overhead calling this
	bn, err := c.BlockNumber(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Unable to get block number while checking gas prices")
		return nil, nil
	}
	isDynamic := inputLoadTestParams.ChainSupportBaseFee

	cachedGasPriceLock.Lock()
	defer cachedGasPriceLock.Unlock()
	if bn <= cachedBlockNumber {
		return cachedGasPrice, cachedGasTipCap
	}
	gp, pErr := c.SuggestGasPrice(ctx)
	gt, tErr := c.SuggestGasTipCap(ctx)

	// In the case of an EVM compatible system not supporting EIP-1559
	if *inputLoadTestParams.LegacyTransactionMode {
		gt = big.NewInt(0)
	}

	if pErr == nil && (tErr == nil || !isDynamic) {
		cachedBlockNumber = bn
		cachedGasPrice = gp
		cachedGasTipCap = gt
		if inputLoadTestParams.ForceGasPrice != nil && *inputLoadTestParams.ForceGasPrice != 0 {
			cachedGasPrice = new(big.Int).SetUint64(*inputLoadTestParams.ForceGasPrice)
		}
		if inputLoadTestParams.ForcePriorityGasPrice != nil && *inputLoadTestParams.ForcePriorityGasPrice != 0 {
			cachedGasTipCap = new(big.Int).SetUint64(*inputLoadTestParams.ForcePriorityGasPrice)
		}
		if cachedGasTipCap.Cmp(cachedGasPrice) == 1 {
			cachedGasTipCap = cachedGasPrice
		}
		l := log.Debug().Uint64("cachedBlockNumber", bn).Uint64("cachedgasPrice", cachedGasPrice.Uint64())
		if cachedGasTipCap != nil {
			l = l.Uint64("cachedGasTipCap", cachedGasTipCap.Uint64())
		}
		l.Msg("Updating gas prices")

		return cachedGasPrice, cachedGasTipCap
	}

	// Something went wrong
	if pErr != nil {
		log.Error().Err(pErr).Msg("Unable to suggest gas price")
		return cachedGasPrice, cachedGasTipCap
	}
	if tErr != nil && isDynamic {
		log.Error().Err(tErr).Msg("Unable to suggest gas tip cap")
		return cachedGasPrice, cachedGasTipCap
	}
	log.Error().Err(tErr).Msg("This error should not have happened. We got a gas tip price error in an environment that is not dynamic")
	return cachedGasPrice, cachedGasTipCap

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
		msg.Data = ethcommon.FromHex(tester.LoadTesterMetaData.Bin)
		_, err = c.CallContract(ctx, msg, nil)
	} else {
		_, _, _, err = tester.DeployLoadTester(tops, c)
	}
	return
}

// getCurrentLoadTestFunction is meant to handle the business logic
// around deciding which function to execute. When we're in function
// mode where the user has provided a specific function to execute, we
// should use that function. Otherwise, we'll select random functions.
func getCurrentLoadTestFunction() uint64 {
	if loadTestModeFunction == inputLoadTestParams.Mode {
		return *inputLoadTestParams.Function
	}
	return tester.GetRandomOPCode()
}
func loadTestFunction(ctx context.Context, c *ethclient.Client, nonce uint64, ltContract *tester.LoadTester) (t1 time.Time, t2 time.Time, err error) {
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
		tx, err = tester.CallLoadTestFunctionByOpCode(f, ltContract, tops, *iterations)
		if err != nil {
			return
		}
		msg := txToCallMsg(tx)
		_, err = c.CallContract(ctx, msg, nil)
	} else {
		_, err = tester.CallLoadTestFunctionByOpCode(f, ltContract, tops, *iterations)
	}
	return
}

func loadTestCallPrecompiledContracts(ctx context.Context, c *ethclient.Client, nonce uint64, ltContract *tester.LoadTester, useSelectedAddress bool) (t1 time.Time, t2 time.Time, err error) {
	var f int
	ltp := inputLoadTestParams

	chainID := new(big.Int).SetUint64(*ltp.ChainID)
	privateKey := ltp.ECDSAPrivateKey
	iterations := ltp.Iterations
	if useSelectedAddress {
		f = int(*ltp.Function)
	} else {
		f = tester.GetRandomPrecompiledContractAddress()
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
		tx, err = tester.CallPrecompiledContracts(f, ltContract, tops, *iterations, privateKey)
		if err != nil {
			return
		}
		msg := txToCallMsg(tx)
		_, err = c.CallContract(ctx, msg, nil)
	} else {
		_, err = tester.CallPrecompiledContracts(f, ltContract, tops, *iterations, privateKey)
	}
	return
}

func loadTestInc(ctx context.Context, c *ethclient.Client, nonce uint64, ltContract *tester.LoadTester) (t1 time.Time, t2 time.Time, err error) {
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

func loadTestStore(ctx context.Context, c *ethclient.Client, nonce uint64, ltContract *tester.LoadTester) (t1 time.Time, t2 time.Time, err error) {
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

func loadTestERC20(ctx context.Context, c *ethclient.Client, nonce uint64, erc20Contract *tokens.ERC20, ltAddress ethcommon.Address) (t1 time.Time, t2 time.Time, err error) {
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

func loadTestERC721(ctx context.Context, c *ethclient.Client, nonce uint64, erc721Contract *tokens.ERC721, ltAddress ethcommon.Address) (t1 time.Time, t2 time.Time, err error) {
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

func loadTestRecall(ctx context.Context, c *ethclient.Client, nonce uint64, originalTx rpctypes.PolyTransaction) (t1 time.Time, t2 time.Time, err error) {
	ltp := inputLoadTestParams

	chainID := new(big.Int).SetUint64(*ltp.ChainID)
	privateKey := ltp.ECDSAPrivateKey

	tops, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Error().Err(err).Msg("Unable create transaction signer")
		return
	}
	gasPrice, gasTipCap := getSuggestedGasPrices(ctx, c)
	tx := rawTransactionToNewTx(originalTx, nonce, gasPrice, gasTipCap)
	tops = configureTransactOpts(tops)

	stx, err := tops.Signer(*ltp.FromETHAddress, tx)
	if err != nil {
		log.Error().Err(err).Msg("Unable to sign transaction")
		return
	}
	log.Trace().Str("txId", originalTx.Hash().String()).Bool("callOnly", *ltp.CallOnly).Msg("Attempting to replay transaction")

	t1 = time.Now()
	defer func() { t2 = time.Now() }()
	if *ltp.CallOnly {
		callMsg := txToCallMsg(stx)
		callMsg.From = originalTx.From()
		callMsg.Gas = originalTx.Gas()
		if *ltp.CallOnlyLatestBlock {
			_, err = c.CallContract(ctx, callMsg, nil)
		} else {
			callMsg.GasPrice = originalTx.GasPrice()
			callMsg.GasFeeCap = new(big.Int).SetUint64(originalTx.MaxFeePerGas())
			callMsg.GasTipCap = new(big.Int).SetUint64(originalTx.MaxPriorityFeePerGas())
			_, err = c.CallContract(ctx, callMsg, originalTx.BlockNumber())
		}
		if err != nil {
			log.Warn().Err(err).Msg("Recall failure")
		}
		// we're not going to return the error in the case because there is no point retrying
		err = nil
	} else {
		err = c.SendTransaction(ctx, stx)
	}
	return
}

func loadTestRPC(ctx context.Context, c *ethclient.Client, nonce uint64, ia *IndexedActivity) (t1 time.Time, t2 time.Time, err error) {
	funcNum := randSrc.Intn(300)
	t1 = time.Now()
	defer func() { t2 = time.Now() }()
	if funcNum < 10 {
		log.Trace().Msg("eth_gasPrice")
		_, err = c.SuggestGasPrice(ctx)
	} else if funcNum < 21 {
		log.Trace().Msg("eth_estimateGas")
		var rawTxData []byte
		pt := ia.Transactions[randSrc.Intn(len(ia.TransactionIDs))]
		rawTxData, err = pt.MarshalJSON()
		if err != nil {
			log.Error().Err(err).Str("txHash", pt.Hash().String()).Msg("issue converting poly transaction to json")
			return
		}
		var tx apitypes.SendTxArgs
		err = json.Unmarshal(rawTxData, &tx)
		if err != nil {
			log.Error().Err(err).Str("txHash", pt.Hash().String()).Msg("unable to unmarshal poly transaction to json.")
			return
		}
		cm := txToCallMsg(tx.ToTransaction())
		cm.From = pt.From()
		_, err = c.EstimateGas(ctx, cm)
	} else if funcNum < 33 {
		log.Trace().Msg("eth_getTransactionCount")
		_, err = c.NonceAt(ctx, ethcommon.HexToAddress(ia.Addresses[randSrc.Intn(len(ia.Addresses))]), nil)
	} else if funcNum < 47 {
		log.Trace().Msg("eth_getCode")
		_, err = c.CodeAt(ctx, ethcommon.HexToAddress(ia.Contracts[randSrc.Intn(len(ia.Contracts))]), nil)
	} else if funcNum < 64 {
		log.Trace().Msg("eth_getBlockByNumber")
		_, err = c.BlockByNumber(ctx, big.NewInt(int64(randSrc.Intn(int(ia.BlockNumber)))))
	} else if funcNum < 84 {
		log.Trace().Msg("eth_getTransactionByHash")
		_, _, err = c.TransactionByHash(ctx, ethcommon.HexToHash(ia.TransactionIDs[randSrc.Intn(len(ia.TransactionIDs))]))
	} else if funcNum < 109 {
		log.Trace().Msg("eth_getBalance")
		_, err = c.BalanceAt(ctx, ethcommon.HexToAddress(ia.Addresses[randSrc.Intn(len(ia.Addresses))]), nil)
	} else if funcNum < 142 {
		log.Trace().Msg("eth_getTransactionReceipt")
		_, err = c.TransactionReceipt(ctx, ethcommon.HexToHash(ia.TransactionIDs[randSrc.Intn(len(ia.TransactionIDs))]))
	} else if funcNum < 192 {
		log.Trace().Msg("eth_getLogs")
		h := ethcommon.HexToHash(ia.BlockIDs[randSrc.Intn(len(ia.BlockIDs))])
		_, err = c.FilterLogs(ctx, ethereum.FilterQuery{BlockHash: &h})
	} else {

		log.Trace().Msg("eth_call")

		if len(ia.ERC20Addresses) != 0 {
			erc20Str := string(ia.ERC20Addresses[randSrc.Intn(len(ia.ERC20Addresses))])
			erc20Addr := ethcommon.HexToAddress(erc20Str)

			log.Trace().
				Str("erc20str", erc20Str).
				Str("erc20addr", erc20Addr.String()).
				Msg("Retrieve contract addresses")
			cops := new(bind.CallOpts)
			cops.Context = ctx
			var erc20Contract *tokens.ERC20

			erc20Contract, err = tokens.NewERC20(erc20Addr, c)
			if err != nil {
				log.Error().Err(err).Msg("Unable to instantiate new erc20 contract")
				return
			}
			t1 = time.Now()

			_, err = erc20Contract.BalanceOf(cops, *inputLoadTestParams.FromETHAddress)
			if err != nil && err == bind.ErrNoCode {
				err = nil
			}
			// tokenURI would be the next most popular call, but it's not very complex
		} else {
			log.Warn().Msg("Unable to find deployed erc20 contract, skipping making calls...")
		}

		if len(ia.ERC721Addresses) != 0 {
			erc721Str := string(ia.ERC721Addresses[randSrc.Intn(len(ia.ERC721Addresses))])
			erc721Addr := ethcommon.HexToAddress(erc721Str)

			log.Trace().
				Str("erc721str", erc721Str).
				Str("erc721addr", erc721Addr.String()).
				Msg("Retrieve contract addresses")
			cops := new(bind.CallOpts)
			cops.Context = ctx
			var erc721Contract *tokens.ERC721

			erc721Contract, err = tokens.NewERC721(erc721Addr, c)
			if err != nil {
				log.Error().Err(err).Msg("Unable to instantiate new erc721 contract")
				return
			}
			t1 = time.Now()

			_, err = erc721Contract.BalanceOf(cops, *inputLoadTestParams.FromETHAddress)
			if err != nil && err == bind.ErrNoCode {
				err = nil
			}
			// tokenURI would be the next most popular call, but it's not very complex
		} else {
			log.Warn().Msg("Unable to find deployed erc721 contract, skipping making calls...")
		}
	}

	return
}

func loadTestContractCall(ctx context.Context, c *ethclient.Client, nonce uint64) (t1 time.Time, t2 time.Time, err error) {
	ltp := inputLoadTestParams

	to := ltp.ContractETHAddress

	chainID := new(big.Int).SetUint64(*ltp.ChainID)
	privateKey := ltp.ECDSAPrivateKey

	tops, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Error().Err(err).Msg("Unable create transaction signer")
		return
	}

	amount := big.NewInt(0)
	if *ltp.ContractCallPayable {
		amount = ltp.SendAmount
	}

	tops = configureTransactOpts(tops)
	gasPrice, gasTipCap := getSuggestedGasPrices(ctx, c)

	var stringCallData string

	if *inputLoadTestParams.ContractCallData == "" && *inputLoadTestParams.ContractCallFunctionSignature == "" {
		log.Error().Err(fmt.Errorf("Missing calldata for function call"))
		return
	}

	if *inputLoadTestParams.ContractCallData != "" {
		stringCallData = *inputLoadTestParams.ContractCallData
	} else {
		stringCallData, err = abi.AbiEncode(*inputLoadTestParams.ContractCallFunctionSignature, *inputLoadTestParams.ContractCallFunctionArgs)
		if err != nil {
			log.Error().Err(err).Msg("Failed to encode calldata")
			return
		}
	}

	calldata, err := hex.DecodeString(strings.TrimPrefix(stringCallData, "0x"))
	if err != nil {
		log.Error().Err(err).Msg("Unable to decode calldata string")
		return
	}

	if tops.GasLimit == 0 {
		estimateInput := ethereum.CallMsg{
			From:      tops.From,
			To:        to,
			Value:     amount,
			GasPrice:  tops.GasPrice,
			GasTipCap: tops.GasTipCap,
			GasFeeCap: tops.GasFeeCap,
			Data:      calldata,
		}
		tops.GasLimit, err = c.EstimateGas(ctx, estimateInput)
		if err != nil {
			log.Error().Err(err).Msg("Unable to estimate gas for transaction. Manually setting gas-limit might be required")
			return
		}
	}

	var tx *ethtypes.Transaction
	if *ltp.LegacyTransactionMode {
		tx = ethtypes.NewTx(&ethtypes.LegacyTx{
			Nonce:    nonce,
			To:       to,
			Value:    amount,
			Gas:      tops.GasLimit,
			GasPrice: gasPrice,
			Data:     calldata,
		})
	} else {
		tx = ethtypes.NewTx(&ethtypes.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     nonce,
			To:        to,
			Gas:       tops.GasLimit,
			GasFeeCap: gasPrice,
			GasTipCap: gasTipCap,
			Data:      calldata,
			Value:     amount,
		})
	}
	log.Trace().Interface("tx", tx).Msg("Contract call data")

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

func loadTestInscription(ctx context.Context, c *ethclient.Client, nonce uint64) (t1 time.Time, t2 time.Time, err error) {
	ltp := inputLoadTestParams

	to := ltp.FromETHAddress

	chainID := new(big.Int).SetUint64(*ltp.ChainID)
	privateKey := ltp.ECDSAPrivateKey

	tops, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Error().Err(err).Msg("Unable create transaction signer")
		return
	}

	amount := big.NewInt(0)
	tops = configureTransactOpts(tops)
	gasPrice, gasTipCap := getSuggestedGasPrices(ctx, c)

	calldata := []byte(*ltp.InscriptionContent)
	if tops.GasLimit == 0 {
		estimateInput := ethereum.CallMsg{
			From:      tops.From,
			To:        to,
			Value:     amount,
			GasPrice:  tops.GasPrice,
			GasTipCap: tops.GasTipCap,
			GasFeeCap: tops.GasFeeCap,
			Data:      calldata,
		}
		tops.GasLimit, err = c.EstimateGas(ctx, estimateInput)
		if err != nil {
			log.Error().Err(err).Msg("Unable to estimate gas for transaction. Manually setting gas-limit might be required")
			return
		}
	}

	var tx *ethtypes.Transaction
	if *ltp.LegacyTransactionMode {
		tx = ethtypes.NewTx(&ethtypes.LegacyTx{
			Nonce:    nonce,
			To:       to,
			Value:    amount,
			Gas:      tops.GasLimit,
			GasPrice: gasPrice,
			Data:     calldata,
		})
	} else {
		tx = ethtypes.NewTx(&ethtypes.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     nonce,
			To:        to,
			Gas:       tops.GasLimit,
			GasFeeCap: gasPrice,
			GasTipCap: gasTipCap,
			Data:      calldata,
			Value:     amount,
		})
	}
	log.Trace().Interface("tx", tx).Msg("Contract call data")

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

func loadTestBlob(ctx context.Context, c *ethclient.Client, nonce uint64) (t1 time.Time, t2 time.Time, err error) {
	ltp := inputLoadTestParams

	to := ltp.ToETHAddress
	if *ltp.ToRandom {
		to = getRandomAddress()
	}

	amount := ltp.SendAmount
	chainID := new(big.Int).SetUint64(*ltp.ChainID)
	privateKey := ltp.ECDSAPrivateKey

	gasLimit := uint64(21000)
	gasPrice, gasTipCap := getSuggestedGasPrices(ctx, c)
	// blobFeeCap := uint64(1000000000) // 1eth
	blobFeeCap := ltp.BlobFeeCap

	// Initialize blobTx with blob transaction type
	blobTx := ethtypes.BlobTx{
		ChainID:    uint256.NewInt(chainID.Uint64()),
		Nonce:      nonce,
		GasTipCap:  uint256.NewInt(gasTipCap.Uint64()),
		GasFeeCap:  uint256.NewInt(gasPrice.Uint64()),
		BlobFeeCap: uint256.NewInt(*blobFeeCap),
		Gas:        gasLimit,
		To:         *to,
		Value:      uint256.NewInt(amount.Uint64()),
		Data:       nil,
		AccessList: nil,
		BlobHashes: make([]common.Hash, 0),
		Sidecar: &types.BlobTxSidecar{
			Blobs:       make([]kzg4844.Blob, 0),
			Commitments: make([]kzg4844.Commitment, 0),
			Proofs:      make([]kzg4844.Proof, 0),
		},
	}
	// appendBlobCommitment() will take in the blobTx struct and append values to blob transaction specific keys in the following steps:
	// The function will take in blobTx with empty BlobHashses, and Blob Sidecar variables initially.
	// Then generateRandomBlobData() is called to generate a byte slice with random values.
	// createBlob() is called to commit the randomly generated byte slice with KZG.
	// generateBlobCommitment() will do the same for the Commitment and Proof.
	// Append all the blob related computed values to the blobTx struct.
	err = appendBlobCommitment(&blobTx)
	if err != nil {
		log.Error().Err(err).Msg("Unable to parse blob")
		return
	}
	tx := types.NewTx(&blobTx)

	stx, err := types.SignTx(tx, types.LatestSignerForChainID(chainID), privateKey)
	if err != nil {
		log.Error().Err(err).Msg("Unable to sign transaction")
		return
	}

	t1 = time.Now()
	defer func() { t2 = time.Now() }()
	if *ltp.CallOnly {
		log.Error().Err(err).Msg("CallOnly not supported to blob transactions")
		return
	} else {
		err = c.SendTransaction(ctx, stx)
	}
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
	var prevNonceForFinalBlock uint64
	var currentNonceForFinalBlock uint64
	var initialWaitCount = 20
	var maxWaitCount = initialWaitCount
	for {
		lastBlockNumber, err = c.BlockNumber(ctx)
		if err != nil {
			return 0, err
		}
		if *ltp.CallOnly {
			return lastBlockNumber, nil
		}
		currentNonceForFinalBlock, err = c.NonceAt(ctx, *ltp.FromETHAddress, new(big.Int).SetUint64(lastBlockNumber))
		if err != nil {
			return 0, err
		}
		if currentNonceForFinalBlock < endNonce && maxWaitCount > 0 {
			log.Trace().Uint64("endNonce", endNonce).Uint64("currentNonceForFinalBlock", currentNonceForFinalBlock).Uint64("prevNonceForFinalBlock", prevNonceForFinalBlock).Msg("Not all transactions have been mined. Waiting")
			time.Sleep(5 * time.Second)
			if currentNonceForFinalBlock == prevNonceForFinalBlock {
				maxWaitCount = maxWaitCount - 1 // only decrement if currentNonceForFinalBlock doesn't progress
			}
			prevNonceForFinalBlock = currentNonceForFinalBlock
			log.Trace().Int("Remaining Attempts", maxWaitCount).Msg("Retrying...")
			continue
		}
		if maxWaitCount <= 0 {
			return 0, fmt.Errorf("waited for %d attempts for the transactions to be mined", initialWaitCount)
		}
		break
	}

	log.Trace().Uint64("currentNonceForFinalBlock", currentNonceForFinalBlock).Uint64("startblock", startBlockNumber).Uint64("endblock", lastBlockNumber).Msg("It looks like all transactions have been mined")
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
