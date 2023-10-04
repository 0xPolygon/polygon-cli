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

	"github.com/maticnetwork/polygon-cli/contracts/tokens"
	"github.com/maticnetwork/polygon-cli/contracts/uniswapv3"
	"github.com/maticnetwork/polygon-cli/rpctypes"

	_ "embed"

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
	loadTestModeUniswapV3

	// All the modes AFTER random mode will not be used when mode random is selected
	loadTestModeRandom
	loadTestModeRecall
	loadTestModeRPC

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
		m == loadTestModeStore {
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
	if header.BaseFee != nil {
		inputLoadTestParams.ChainSupportBaseFee = true
		log.Debug().Msg("eip-1559 support detected")
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

	modes := *inputLoadTestParams.Modes
	if len(modes) == 0 {
		return fmt.Errorf("expected at least one mode")
	}

	inputLoadTestParams.ParsedModes = make([]loadTestMode, 0)
	for _, m := range modes {
		parsedMode, err := characterToLoadTestMode(m)
		if err != nil {
			return err
		}
		inputLoadTestParams.ParsedModes = append(inputLoadTestParams.ParsedModes, parsedMode)
	}

	if len(modes) > 1 {
		inputLoadTestParams.MultiMode = true
	} else {
		inputLoadTestParams.MultiMode = false
		inputLoadTestParams.Mode, _ = characterToLoadTestMode((*inputLoadTestParams.Modes)[0])
	}

	if hasMode(loadTestModeRandom, inputLoadTestParams.ParsedModes) && inputLoadTestParams.MultiMode {
		return fmt.Errorf("random mode can't be used in combinations with any other modes")
	}
	if hasMode(loadTestModeRPC, inputLoadTestParams.ParsedModes) && inputLoadTestParams.MultiMode && !*inputLoadTestParams.CallOnly {
		return fmt.Errorf("rpc mode must be called with call-only when multiple modes are used")
	} else if hasMode(loadTestModeRPC, inputLoadTestParams.ParsedModes) {
		log.Trace().Msg("setting call only mode since we're doing RPC testing")
		*inputLoadTestParams.CallOnly = true
	}
	// TODO check for duplicate modes?

	if *inputLoadTestParams.CallOnly && *inputLoadTestParams.AdaptiveRateLimit {
		return fmt.Errorf("using call only with adaptive rate limit doesn't make sense")
	}

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

	loopFunc := func() error {
		err = initializeLoadTestParams(ctx, ec)
		if err != nil {
			return err
		}

		return mainLoop(ctx, ec, rpc)
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
	mode := ltp.Mode
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
	if anyModeRequiresLoadTestContract(ltp.ParsedModes) || *inputLoadTestParams.ForceContractDeploy {
		ltAddr, ltContract, err = getLoadTestContract(ctx, c, tops, cops)
		if err != nil {
			return err
		}
		log.Debug().Str("ltAddr", ltAddr.String()).Msg("Obtained load test contract address")
	}

	var erc20Addr ethcommon.Address
	var erc20Contract *tokens.ERC20
	if mode == loadTestModeERC20 || mode == loadTestModeRandom {
		erc20Addr, erc20Contract, err = getERC20Contract(ctx, c, tops, cops)
		if err != nil {
			return err
		}
		log.Debug().Str("erc20Addr", erc20Addr.String()).Msg("Obtained erc 20 contract address")
	}

	var erc721Addr ethcommon.Address
	var erc721Contract *tokens.ERC721
	if mode == loadTestModeERC721 || mode == loadTestModeRandom {
		erc721Addr, erc721Contract, err = getERC721Contract(ctx, c, tops, cops)
		if err != nil {
			return err
		}
		log.Debug().Str("erc721Addr", erc721Addr.String()).Msg("Obtained erc 721 contract address")
	}

	uniswapAddresses := UniswapV3Addresses{
		FactoryV3:                   ethcommon.HexToAddress(*uniswapv3LoadTestParams.UniswapFactoryV3),
		Multicall:                   ethcommon.HexToAddress(*uniswapv3LoadTestParams.UniswapMulticall),
		ProxyAdmin:                  ethcommon.HexToAddress(*uniswapv3LoadTestParams.UniswapProxyAdmin),
		TickLens:                    ethcommon.HexToAddress(*uniswapv3LoadTestParams.UniswapTickLens),
		NFTDescriptorLib:            ethcommon.HexToAddress(*uniswapv3LoadTestParams.UniswapNFTLibDescriptor),
		NFTPositionDescriptor:       ethcommon.HexToAddress(*uniswapv3LoadTestParams.UniswapNFTPositionDescriptor),
		TransparentUpgradeableProxy: ethcommon.HexToAddress(*uniswapv3LoadTestParams.UniswapUpgradeableProxy),
		NFPositionManager:           ethcommon.HexToAddress(*uniswapv3LoadTestParams.UniswapNFPositionManager),
		Migrator:                    ethcommon.HexToAddress(*uniswapv3LoadTestParams.UniswapMigrator),
		Staker:                      ethcommon.HexToAddress(*uniswapv3LoadTestParams.UniswapStaker),
		QuoterV2:                    ethcommon.HexToAddress(*uniswapv3LoadTestParams.UniswapQuoterV2),
		SwapRouter02:                ethcommon.HexToAddress(*uniswapv3LoadTestParams.UniswapSwapRouter),
		WETH9:                       ethcommon.HexToAddress(*uniswapv3LoadTestParams.WETH9),
	}
	var uniswapV3Config UniswapV3Config
	var poolConfig PoolConfig
	if mode == loadTestModeUniswapV3 || mode == loadTestModeRandom {
		uniswapV3Config, err = deployUniswapV3(ctx, c, tops, cops, uniswapAddresses, *ltp.FromETHAddress)
		if err != nil {
			return nil
		}
		log.Debug().Interface("config", uniswapV3Config.ToAddresses()).Msg("UniswapV3 deployment config")

		tokensAToMint := big.NewInt(1_000_000_000_000_000_000)
		var token0Config contractConfig[uniswapv3.Swapper]
		token0Config, err = deploySwapperContract(ctx, c, tops, cops, uniswapV3Config, "Token0", "A", tokensAToMint, *ltp.FromETHAddress, ethcommon.HexToAddress(*uniswapv3LoadTestParams.UniswapPoolToken0))
		if err != nil {
			return nil
		}

		tokensBToMint := big.NewInt(1_000_000_000_000_000_000)
		var token1Config contractConfig[uniswapv3.Swapper]
		token1Config, err = deploySwapperContract(ctx, c, tops, cops, uniswapV3Config, "Token1", "B", tokensBToMint, *ltp.FromETHAddress, ethcommon.HexToAddress(*uniswapv3LoadTestParams.UniswapPoolToken1))
		if err != nil {
			return nil
		}

		fees := big.NewInt(3_000)
		poolConfig := PoolConfig{Fees: fees}
		if token0Config.Address.Hex() < token1Config.Address.Hex() {
			poolConfig.Token0 = token0Config
			poolConfig.ReserveA = tokensAToMint
			poolConfig.Token1 = token1Config
			poolConfig.ReserveB = tokensBToMint
		} else {
			poolConfig.Token0 = token1Config
			poolConfig.ReserveA = tokensBToMint
			poolConfig.Token1 = token0Config
			poolConfig.ReserveB = tokensAToMint
		}
		if err = createPool(ctx, c, tops, cops, uniswapV3Config, poolConfig, *ltp.FromETHAddress); err != nil {
			return nil
		}
	}

	var recallTransactions []rpctypes.PolyTransaction
	if mode == loadTestModeRecall {
		recallTransactions, err = getRecallTransactions(ctx, c, rpc)
		if err != nil {
			return err
		}
		if len(recallTransactions) == 0 {
			return fmt.Errorf("we weren't able to fetch any recall transactions")
		}
		log.Debug().Int("txs", len(recallTransactions)).Msg("retrieved transactions for total recall")
	}

	var indexedActivity *IndexedActivity
	if mode == loadTestModeRPC || mode == loadTestModeRandom {
		indexedActivity, err = getIndexedRecentActivity(ctx, c, rpc)
		if err != nil {
			return err
		}
		log.Debug().
			Int("transactions", len(indexedActivity.TransactionIDs)).
			Int("blocks", len(indexedActivity.BlockNumbers)).
			Int("addresses", len(indexedActivity.Addresses)).
			Int("erc20s", len(indexedActivity.ERC20Addresses)).
			Int("erc721", len(indexedActivity.ERC721Addresses)).
			Int("contracts", len(indexedActivity.Contracts)).
			Msg("retrieved recent indexed activity")
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
					startReq, endReq, tErr = loadTestUniswapV3(ctx, c, myNonceValue, uniswapV3Config, poolConfig)
				case loadTestModeRPC:
					startReq, endReq, tErr = loadTestRPC(ctx, c, myNonceValue, indexedActivity)
				default:
					log.Error().Str("mode", mode.String()).Msg("We've arrived at a load test mode that we don't recognize")
				}
				recordSample(i, j, tErr, startReq, endReq, myNonceValue)
				if tErr != nil {
					log.Error().Err(tErr).Uint64("nonce", myNonceValue).Msg("Recorded an error while sending transactions")
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
func getERC20Contract(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts, cops *bind.CallOpts) (erc20Addr ethcommon.Address, erc20Contract *tokens.ERC20, err error) {
	erc20Addr = ethcommon.HexToAddress(*inputLoadTestParams.ERC20Address)
	shouldMint := false
	if *inputLoadTestParams.ERC20Address == "" {
		erc20Addr, _, _, err = tokens.DeployERC20(tops, c, "ERC20TestToken", "T20")
		if err != nil {
			log.Error().Err(err).Msg("Unable to deploy ERC20 contract")
			return
		}
		// if we're deploying a new ERC 20 we should mint tokens
		shouldMint = true
	}
	log.Trace().Interface("contractaddress", erc20Addr).Msg("ERC20 contract address")

	erc20Contract, err = tokens.NewERC20(erc20Addr, c)
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

	if !shouldMint {
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

	err = blockUntilSuccessful(ctx, c, func() error {
		_, err = erc721Contract.BalanceOf(cops, *inputLoadTestParams.FromETHAddress)
		return err
	})
	if err != nil {
		return
	}
	if !shouldMint {
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
	if pErr == nil && (tErr == nil || !isDynamic) {
		cachedBlockNumber = bn
		cachedGasPrice = gp
		cachedGasTipCap = gt
		if inputLoadTestParams.ForceGasPrice != nil && *inputLoadTestParams.ForcePriorityGasPrice != 0 {
			cachedGasPrice = new(big.Int).SetUint64(*inputLoadTestParams.ForcePriorityGasPrice)
		}
		if inputLoadTestParams.ForcePriorityGasPrice != nil && *inputLoadTestParams.ForcePriorityGasPrice != 0 {
			cachedGasTipCap = new(big.Int).SetUint64(*inputLoadTestParams.ForcePriorityGasPrice)
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
	if loadTestModeFunction == inputLoadTestParams.Mode {
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
		erc20Str := string(ia.ERC20Addresses[randSrc.Intn(len(ia.ERC20Addresses))])
		erc721Str := string(ia.ERC721Addresses[randSrc.Intn(len(ia.ERC721Addresses))])
		erc20Addr := ethcommon.HexToAddress(erc20Str)
		erc721Addr := ethcommon.HexToAddress(erc721Str)
		log.Trace().
			Str("erc20str", erc20Str).
			Str("erc721str", erc721Str).
			Str("erc20addr", erc20Addr.String()).
			Str("erc721addr", erc721Addr.String()).
			Msg("retrieve contract addresses")
		cops := new(bind.CallOpts)
		cops.Context = ctx
		var erc721Contract *tokens.ERC721
		var erc20Contract *tokens.ERC20

		erc721Contract, err = tokens.NewERC721(erc721Addr, c)
		if err != nil {
			log.Error().Err(err).Msg("Unable to instantiate new erc721 contract")
			return
		}
		erc20Contract, err = tokens.NewERC20(erc20Addr, c)
		if err != nil {
			log.Error().Err(err).Msg("Unable to instantiate new erc20 contract")
			return
		}
		t1 = time.Now()

		_, err = erc721Contract.BalanceOf(cops, *inputLoadTestParams.FromETHAddress)
		if err != nil && err == bind.ErrNoCode {
			err = nil
		}
		_, err = erc20Contract.BalanceOf(cops, *inputLoadTestParams.FromETHAddress)
		if err != nil && err == bind.ErrNoCode {
			err = nil
		}
		// tokenURI would be the next most popular call, but it's not very complex

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
