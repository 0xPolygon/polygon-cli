package loadtest

import (
	"context"
	"crypto/ecdsa"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"sync/atomic"

	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/holiman/uint256"

	"github.com/0xPolygon/polygon-cli/bindings/tester"
	"github.com/0xPolygon/polygon-cli/bindings/tokens"
	"github.com/0xPolygon/polygon-cli/cmd/loadtest/gasmanager"
	uniswapv3loadtest "github.com/0xPolygon/polygon-cli/cmd/loadtest/uniswapv3"

	"github.com/0xPolygon/polygon-cli/rpctypes"
	"github.com/0xPolygon/polygon-cli/util"

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
	// These constants are "stringered".
	// If you add a new constant, it fill fail to compile until you regenerate the strings.
	// There are two steps needed:
	// 1. Install stringer: `go install golang.org/x/tools/cmd/stringer`.
	// 2. Generate the string: `go generate github.com/0xPolygon/polygon-cli/cmd/loadtest`.
	// You can also use `make gen-loadtest-modes`.
	loadTestModeERC20 loadTestMode = iota
	loadTestModeERC721
	loadTestModeBlob
	loadTestModeContractCall
	loadTestModeDeploy
	loadTestModeIncrement
	loadTestModeRandom
	loadTestModeRecall
	loadTestModeRPC
	loadTestModeStore
	loadTestModeTransaction
	loadTestModeUniswapV3

	codeQualitySeed       = "code code code code code code code code code code code quality"
	codeQualityPrivateKey = "42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa"

	oneEtherInWei = 1000000000000000000 // 1 ETH in wei
)

func characterToLoadTestMode(mode string) (loadTestMode, error) {
	switch mode {
	case "2", "erc20":
		return loadTestModeERC20, nil
	case "7", "erc721":
		return loadTestModeERC721, nil
	case "b", "blob":
		return loadTestModeBlob, nil
	case "cc", "contract-call":
		return loadTestModeContractCall, nil
	case "d", "deploy":
		return loadTestModeDeploy, nil
	case "inc", "increment":
		return loadTestModeIncrement, nil
	case "r", "random":
		return loadTestModeRandom, nil
	case "R", "recall":
		return loadTestModeRecall, nil
	case "rpc":
		return loadTestModeRPC, nil
	case "s", "store":
		return loadTestModeStore, nil
	case "t", "transaction":
		return loadTestModeTransaction, nil
	case "v3", "uniswapv3":
		return loadTestModeUniswapV3, nil
	default:
		return 0, fmt.Errorf("unrecognized load test mode: %s", mode)
	}
}

func getRandomMode() loadTestMode {
	// Does not include the following modes:
	// blob, contract call, recall, rpc, uniswap v3
	modes := []loadTestMode{
		loadTestModeERC20,
		// loadTestModeERC721,
		// loadTestModeDeploy,
		// loadTestModeIncrement,
		// loadTestModeStore,
		loadTestModeTransaction,
	}
	return modes[randSrc.Intn(len(modes))]
}

func modeRequiresLoadTestContract(m loadTestMode) bool {
	if m == loadTestModeIncrement ||
		m == loadTestModeRandom ||
		m == loadTestModeStore {
		return true
	}
	return false
}
func anyModeRequiresLoadTestContract(modes []loadTestMode) bool {
	return slices.ContainsFunc(modes, modeRequiresLoadTestContract)
}
func hasMode(mode loadTestMode, modes []loadTestMode) bool {
	return slices.Contains(modes, mode)
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

	// When outputting raw transactions, we don't need to wait for anything to be mined
	if inputLoadTestParams.OutputRawTxOnly {
		inputLoadTestParams.FireAndForget = true
		log.Debug().Msg("OutputRawTxOnly mode enabled - automatically enabling FireAndForget mode")
	}

	gas, err := c.SuggestGasPrice(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Unable to retrieve gas price")
		return err
	}
	log.Trace().Interface("gasprice", gas).Msg("Retrieved current gas price")

	if !inputLoadTestParams.LegacyTransactionMode {
		gasTipCap, _err := c.SuggestGasTipCap(ctx)
		if _err != nil {
			log.Error().Err(_err).Msg("Unable to retrieve gas tip cap")
			return _err
		}
		log.Trace().Interface("gastipcap", gasTipCap).Msg("Retrieved current gas tip cap")
		inputLoadTestParams.CurrentGasTipCap = gasTipCap
	}

	trimmedHexPrivateKey := strings.TrimPrefix(inputLoadTestParams.PrivateKey, "0x")
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
	log.Trace().
		Str("addr", ethAddress.Hex()).
		Interface("balance", accountBal).
		Msg("funding account balance")

	toAddr := ethcommon.HexToAddress(inputLoadTestParams.ToAddress)

	amt := new(big.Int).SetUint64(inputLoadTestParams.EthAmountInWei)

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

	inputLoadTestParams.BigGasPriceMultiplier = big.NewFloat(inputLoadTestParams.GasPriceMultiplier)

	if inputLoadTestParams.LegacyTransactionMode && inputLoadTestParams.ForcePriorityGasPrice > 0 {
		log.Warn().Msg("Cannot set priority gas price in legacy mode")
	}
	if inputLoadTestParams.ForceGasPrice < inputLoadTestParams.ForcePriorityGasPrice {
		return errors.New("max priority fee per gas higher than max fee per gas")
	}

	if inputLoadTestParams.AdaptiveRateLimit && inputLoadTestParams.EthCallOnly {
		return errors.New("the adaptive rate limit is based on the pending transaction pool. It doesn't use this feature while also using call only")
	}

	contractAddr := ethcommon.HexToAddress(inputLoadTestParams.ContractAddress)
	inputLoadTestParams.ContractETHAddress = &contractAddr

	inputLoadTestParams.ToETHAddress = &toAddr
	inputLoadTestParams.SendAmount = amt
	inputLoadTestParams.CurrentGasPrice = gas
	inputLoadTestParams.CurrentNonce = &nonce
	inputLoadTestParams.ECDSAPrivateKey = privateKey
	inputLoadTestParams.FromETHAddress = &ethAddress
	if inputLoadTestParams.ChainID == 0 {
		inputLoadTestParams.ChainID = chainID.Uint64()
	}

	modes := inputLoadTestParams.Modes
	if len(modes) == 0 {
		return errors.New("expected at least one mode")
	}

	inputLoadTestParams.ParsedModes = make([]loadTestMode, 0)
	for _, m := range modes {
		var parsedMode loadTestMode
		parsedMode, err = characterToLoadTestMode(m)
		if err != nil {
			return err
		}
		inputLoadTestParams.ParsedModes = append(inputLoadTestParams.ParsedModes, parsedMode)
	}

	// Logic checking input parameters for specific conditions such as multiple inputs.
	if len(modes) > 1 {
		inputLoadTestParams.MultiMode = true
		if !hasUniqueModes(inputLoadTestParams.ParsedModes) {
			return errors.New("duplicate modes detected, check input modes for duplicates")
		}
	} else {
		inputLoadTestParams.MultiMode = false
		inputLoadTestParams.Mode, _ = characterToLoadTestMode((inputLoadTestParams.Modes)[0])
	}
	if hasMode(loadTestModeRandom, inputLoadTestParams.ParsedModes) && inputLoadTestParams.MultiMode {
		return errors.New("random mode can't be used in combinations with any other modes")
	}
	if hasMode(loadTestModeRPC, inputLoadTestParams.ParsedModes) && inputLoadTestParams.MultiMode && !inputLoadTestParams.EthCallOnly {
		return errors.New("rpc mode must be called with eth-call-only when multiple modes are used")
	} else if hasMode(loadTestModeRPC, inputLoadTestParams.ParsedModes) {
		log.Trace().Msg("Setting call only mode since we're doing RPC testing")
		inputLoadTestParams.EthCallOnly = true
	}
	if hasMode(loadTestModeContractCall, inputLoadTestParams.ParsedModes) && (inputLoadTestParams.ContractAddress == "" || inputLoadTestParams.ContractCallData == "") {
		return errors.New("`--contract-call` requires both a `--contract-address` and `--calldata` flags")
	}
	if inputLoadTestParams.EthCallOnly && inputLoadTestParams.AdaptiveRateLimit {
		return errors.New("using call only with adaptive rate limit doesn't make sense")
	}
	if inputLoadTestParams.EthCallOnly && inputLoadTestParams.WaitForReceipt {
		return errors.New("using call only with receipts doesn't make sense")
	}
	if inputLoadTestParams.EthCallOnly && inputLoadTestParams.Mode == loadTestModeBlob {
		return errors.New("using call only with blobs doesn't make sense")
	}
	if inputLoadTestParams.LegacyTransactionMode && inputLoadTestParams.Mode == loadTestModeBlob {
		return errors.New("blob transactions require eip-1559")
	}
	if hasMode(loadTestModeBlob, inputLoadTestParams.ParsedModes) && inputLoadTestParams.MultiMode {
		return errors.New("blob mode should only be used by itself. Blob mode will take significantly longer than other transactions to finalize, and the address will be reserved, preventing other transactions form being made")
	}
	if inputLoadTestParams.OutputRawTxOnly && inputLoadTestParams.MultiMode {
		return errors.New("Raw output is not compatible with multiple modes")
	}
	if inputLoadTestParams.OutputRawTxOnly && hasMode(loadTestModeRPC, inputLoadTestParams.ParsedModes) {
		return errors.New("Raw output is not compatible with RPC mode")
	}
	if inputLoadTestParams.OutputRawTxOnly && hasMode(loadTestModeUniswapV3, inputLoadTestParams.ParsedModes) {
		return errors.New("Raw output is not compatible with UniswapV3 mode")
	}

	randSrc = rand.New(rand.NewSource(inputLoadTestParams.Seed))

	err = initializeAccountPool(ctx, c, privateKey)
	if err != nil {
		return err
	}

	return nil
}

func initializeAccountPool(ctx context.Context, c *ethclient.Client, privateKey *ecdsa.PrivateKey) error {
	var err error
	sendingAccountsCount := inputLoadTestParams.SendingAccountsCount
	preFundSendingAccounts := inputLoadTestParams.PreFundSendingAccounts
	accountFundingAmount := inputLoadTestParams.AccountFundingAmount
	sendingAccountsFile := inputLoadTestParams.SendingAccountsFile
	callOnly := inputLoadTestParams.EthCallOnly
	rateLimit := inputLoadTestParams.RateLimit

	accountPool, err = NewAccountPool(ctx, c, privateKey, accountFundingAmount, rateLimit)
	if err != nil {
		log.Error().Err(err).Msg("Unable to create account pool")
		return fmt.Errorf("unable to create account pool. %w", err)
	}

	if len(sendingAccountsFile) > 0 {
		log.Info().
			Str("sendingAccountsFile", sendingAccountsFile).
			Msg("Adding accounts from file to the account pool")

		privateKeys, iErr := util.ReadPrivateKeysFromFile(sendingAccountsFile)
		if iErr != nil {
			log.Error().
				Err(iErr).
				Msg("Unable to read private keys from file")
			return fmt.Errorf("unable to read private keys from file. %w", iErr)
		}
		if len(privateKeys) == 0 {
			const errMsg = "no private keys found in sending accounts file"
			log.Error().Str("sendingAccountsFile", sendingAccountsFile).Msg(errMsg)
			return errors.New(errMsg)
		}

		if len(privateKeys) > 1 && inputLoadTestParams.StartNonce > 0 {
			log.Fatal().
				Str("sendingAccountsFile", sendingAccountsFile).
				Msg("nonce can't be set while using multiple sending accounts")
		}

		if len(privateKeys) == 1 {
			var nonce *uint64
			if inputLoadTestParams.StartNonce > 0 {
				nonce = &inputLoadTestParams.StartNonce
			}
			err = accountPool.Add(ctx, privateKeys[0], nonce)
		} else {
			err = accountPool.AddN(ctx, privateKeys...)
		}

		sendingAccountsCount = uint64(len(privateKeys))
	} else if sendingAccountsCount > 0 {
		log.Info().
			Uint64("sendingAccountsCount", sendingAccountsCount).
			Msg("Adding random accounts to the account pool")

		if inputLoadTestParams.StartNonce > 0 {
			log.Fatal().
				Uint64("sendingAccountsCount", sendingAccountsCount).
				Msg("nonce can't be set while using random multiple sending accounts")
		}

		err = accountPool.AddRandomN(ctx, sendingAccountsCount)
	} else {
		log.Info().
			Uint64("sendingAccountsCount", sendingAccountsCount).
			Msg("Adding single account from private key to the account pool")
		var nonce *uint64
		if inputLoadTestParams.StartNonce > 0 {
			nonce = &inputLoadTestParams.StartNonce
		}
		err = accountPool.Add(ctx, privateKey, nonce)
	}
	if err != nil {
		log.Error().Err(err).Msg("unable to set account pool")
		return fmt.Errorf("unable to set account pool. %w", err)
	}

	// wait all accounts to be ready
	for {
		rdy, rdyCount, accQty := accountPool.AllAccountsReady()
		if rdy {
			log.Info().Msg("All accounts are ready")
			break
		}
		log.Info().Int("ready", rdyCount).Int("total", accQty).Msg("waiting for all accounts to be ready")
		time.Sleep(time.Second)
	}

	// check if there are sending accounts to pre fund
	if sendingAccountsCount == 0 {
		log.Info().Msg("No sending accounts to pre-fund. Skipping pre-funding of sending accounts.")
		return nil
	}

	// checks if call only is enabled
	if callOnly {
		log.Info().Msg("call only mode is enabled. Skipping pre-funding of sending accounts.")
		return nil
	}

	// If pre-funding is disabled, we don't need to fund accounts right now
	if !preFundSendingAccounts {
		log.Info().Msg("pre-funding of sending accounts is disabled.")
		return nil
	}

	// If pre-funding is enabled, we need to fund accounts
	if accountFundingAmount.Cmp(new(big.Int)) == 0 {
		// When using multiple sending accounts and not using --eth-call-only and --account-funding-amount <= 0,
		// we need to make sure the accounts get funded. Set default funding to 1 ETH (1000000000000000000 wei)
		accountPool.fundingAmount = new(big.Int).SetUint64(oneEtherInWei)
		log.Debug().
			Msg("Multiple sending accounts detected with pre-funding enabled with zero funding amount - auto-setting funding amount to 1 ETH")
	}

	err = accountPool.FundAccounts(ctx)
	if err != nil {
		log.Error().Err(err).Msg("unable to fund sending accounts")
	}

	return nil
}

func completeLoadTest(ctx context.Context, c *ethclient.Client, rpc *ethrpc.Client) error {
	if inputLoadTestParams.FireAndForget {
		log.Info().
			Msg("FireAndForget mode enabled - skipping wait period and summarization")
		return nil
	}
	log.Debug().
		Msg("Waiting for remaining transactions to be completed and mined")

	startTime := loadTestResults[0].RequestTime
	endTime := time.Now()
	log.Debug().
		Uint64("final block number", finalBlockNumber).
		Msg("Got final block number")

	if inputLoadTestParams.EthCallOnly {
		log.Info().Msg("CallOnly mode enabled - blocks aren't mined")
		lightSummary(loadTestResults, startTime, endTime, rl)
		return nil
	}

	var err error
	finalBlockNumber, err = waitForFinalBlock(ctx, c, rpc, startBlockNumber)
	if err != nil {
		log.Error().
			Err(err).
			Msg("There was an issue waiting for all transactions to be mined")
	}
	if len(loadTestResults) == 0 {
		return errors.New("no transactions observed")
	}
	endTime = time.Now()

	err = accountPool.ReturnFunds(ctx)
	if err != nil {
		log.Error().
			Err(err).
			Msg("There was an issue returning the funds from the sending accounts back to the funding account")
	}

	if inputLoadTestParams.ShouldProduceSummary {
		err = summarizeTransactions(ctx, c, rpc, startBlockNumber, finalBlockNumber)
		if err != nil {
			log.Error().
				Err(err).
				Msg("There was an issue creating the load test summary")
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
	timeLimit := inputLoadTestParams.TimeLimit
	var overallTimer *time.Timer
	if timeLimit > 0 {
		overallTimer = time.NewTimer(time.Duration(timeLimit) * time.Second)
	} else {
		overallTimer = new(time.Timer)
	}

	// connLimit is the value we'll use to configure the connection limit within the http transport
	connLimit := 2 * int(inputLoadTestParams.Concurrency)
	// Most of these transport options are defaults. We might want to make this configurable from the CLI at some point.
	// The goal here is to avoid opening a ton of connections that go idle then get closed and eventually exhausting
	// client-side connections.
	transport := &http.Transport{
		MaxIdleConns:        connLimit,
		MaxIdleConnsPerHost: connLimit,
		MaxConnsPerHost:     connLimit,
	}
	if inputLoadTestParams.Proxy != "" {
		proxyURL, err := url.Parse(inputLoadTestParams.Proxy)
		if err != nil {
			return fmt.Errorf("invalid proxy address %s %w", inputLoadTestParams.Proxy, err)
		}
		proxyFunc := http.ProxyURL(proxyURL)
		transport.Proxy = proxyFunc
		log.Debug().Stringer("proxyURL", proxyURL).Msg("transport proxy configured")
	}
	goHttpClient := &http.Client{
		Transport: transport,
	}
	rpcOption := ethrpc.WithHTTPClient(goHttpClient)
	rpc, err := ethrpc.DialOptions(ctx, inputLoadTestParams.RPCUrl, rpcOption)
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

		log.Debug().
			Msg("Finished main load test loop")

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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func(ctx context.Context) {
		select {
		case <-ctx.Done():
			return
		default:
			errCh <- loopFunc()
		}
	}(ctx)

	// Wait for the load test to complete, either due to time limit, interrupt signal, or completion.
	select {
	case <-overallTimer.C:
		log.Info().Msg("Time's up")
	case <-sigCh:
		log.Info().Msg("Interrupted.. Stopping load test")
		if inputLoadTestParams.ShouldProduceSummary {
			finalBlockNumber, err = ec.BlockNumber(ctx)
			if err != nil {
				log.Error().Err(err).Msg("Unable to retrieve final block number")
			}
			err = summarizeTransactions(ctx, ec, rpc, startBlockNumber, finalBlockNumber)
			if err != nil {
				log.Error().Err(err).Msg("There was an issue creating the load test summary")
			}
		} else {
			if len(loadTestResults) > 0 {
				lightSummary(loadTestResults, loadTestResults[0].RequestTime, time.Now(), rl)
			}
		}
		cancel()
	case err = <-errCh:
		if err != nil {
			log.Fatal().Err(err).Msg("Received critical error while running load test")
		}
	}
	log.Info().Msg("Finished")
	return nil
}

func updateRateLimit(ctx context.Context, rl *rate.Limiter, rpc *ethrpc.Client, accountPool *AccountPool, steadyStateQueueSize uint64, rateLimitIncrement uint64, cycleDuration time.Duration, backoff float64) {
	tryTxPool := true
	ticker := time.NewTicker(cycleDuration)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			var txPoolSize uint64
			var err error
			var pendingTxs uint64
			var queuedTxs uint64
			// TODO perhaps this should be a mode rather than a fallback
			if tryTxPool {
				pendingTxs, queuedTxs, err = util.GetTxPoolStatus(rpc)
			}

			if err != nil {
				tryTxPool = false
				log.Warn().
					Err(err).
					Msg("Error getting txpool size. Falling back to latest nonce and disabling txpool check")

				pendingTxs, err = accountPool.NumberOfPendingTxs(ctx)
				if err != nil {
					log.Error().
						Err(err).
						Msg("Unable to get pending transactions to update rate limit")
					break
				}

				txPoolSize = pendingTxs
			} else {
				txPoolSize = pendingTxs + queuedTxs
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

	maxRoutines := ltp.Concurrency
	maxRequests := ltp.Requests
	chainID := new(big.Int).SetUint64(ltp.ChainID)
	privateKey := ltp.ECDSAPrivateKey
	mode := ltp.Mode
	steadyStateTxPoolSize := ltp.AdaptiveTargetSize
	adaptiveRateLimitIncrement := ltp.AdaptiveRateLimitIncrement
	rl = rate.NewLimiter(rate.Limit(ltp.RateLimit), 1)
	if ltp.RateLimit <= 0.0 {
		rl = nil
	}

	rateLimitCtx, rateLimitCancel := context.WithCancel(ctx)
	defer rateLimitCancel()
	if ltp.AdaptiveRateLimit && rl != nil {
		go updateRateLimit(rateLimitCtx, rl, rpc, accountPool, steadyStateTxPoolSize, adaptiveRateLimitIncrement, time.Duration(ltp.AdaptiveCycleDuration)*time.Second, ltp.AdaptiveBackoffFactor)
	}

	tops, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	tops = configureTransactOpts(ctx, c, tops, nil)
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
	if anyModeRequiresLoadTestContract(ltp.ParsedModes) {
		ltAddr, ltContract, err = getLoadTestContract(ctx, c, tops, cops)
		if err != nil {
			return err
		}
		log.Debug().Stringer("ltAddr", ltAddr).Msg("Obtained load test contract address")
	}

	var erc20Addr ethcommon.Address
	var erc20Contract *tokens.ERC20
	if hasMode(loadTestModeERC20, ltp.ParsedModes) || hasMode(loadTestModeRandom, ltp.ParsedModes) || hasMode(loadTestModeRPC, ltp.ParsedModes) {
		erc20Addr, erc20Contract, err = getERC20Contract(ctx, c, tops, cops)
		if err != nil {
			return err
		}
		log.Debug().Stringer("erc20Addr", erc20Addr).Msg("Obtained erc 20 contract address")
	}

	var erc721Addr ethcommon.Address
	var erc721Contract *tokens.ERC721
	if hasMode(loadTestModeERC721, ltp.ParsedModes) || hasMode(loadTestModeRandom, ltp.ParsedModes) || hasMode(loadTestModeRPC, ltp.ParsedModes) {
		erc721Addr, erc721Contract, err = getERC721Contract(ctx, c, tops, cops)
		if err != nil {
			return err
		}
		log.Debug().Stringer("erc721Addr", erc721Addr).Msg("Obtained erc 721 contract address")
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

		// Validate that the chain has enough activity for RPC mode
		if len(indexedActivity.TransactionIDs) == 0 ||
			len(indexedActivity.Addresses) == 0 ||
			len(indexedActivity.BlockIDs) == 0 ||
			indexedActivity.BlockNumber == 0 {
			return fmt.Errorf("insufficient chain activity for RPC mode: the chain must have at least some transaction history. Found %d transactions, %d addresses, %d blocks, current block number %d",
				len(indexedActivity.TransactionIDs),
				len(indexedActivity.Addresses),
				len(indexedActivity.BlockIDs),
				indexedActivity.BlockNumber)
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
			FactoryV3:                          ethcommon.HexToAddress(uniswapv3LoadTestParams.UniswapFactoryV3),
			Multicall:                          ethcommon.HexToAddress(uniswapv3LoadTestParams.UniswapMulticall),
			ProxyAdmin:                         ethcommon.HexToAddress(uniswapv3LoadTestParams.UniswapProxyAdmin),
			TickLens:                           ethcommon.HexToAddress(uniswapv3LoadTestParams.UniswapTickLens),
			NFTDescriptorLib:                   ethcommon.HexToAddress(uniswapv3LoadTestParams.UniswapNFTLibDescriptor),
			NonfungibleTokenPositionDescriptor: ethcommon.HexToAddress(uniswapv3LoadTestParams.UniswapNonfungibleTokenPositionDescriptor),
			TransparentUpgradeableProxy:        ethcommon.HexToAddress(uniswapv3LoadTestParams.UniswapUpgradeableProxy),
			NonfungiblePositionManager:         ethcommon.HexToAddress(uniswapv3LoadTestParams.UniswapNonfungiblePositionManager),
			Migrator:                           ethcommon.HexToAddress(uniswapv3LoadTestParams.UniswapMigrator),
			Staker:                             ethcommon.HexToAddress(uniswapv3LoadTestParams.UniswapStaker),
			QuoterV2:                           ethcommon.HexToAddress(uniswapv3LoadTestParams.UniswapQuoterV2),
			SwapRouter02:                       ethcommon.HexToAddress(uniswapv3LoadTestParams.UniswapSwapRouter),
			WETH9:                              ethcommon.HexToAddress(uniswapv3LoadTestParams.WETH9),
		}
		uniswapV3Config, poolConfig, err = initUniswapV3Loadtest(ctx, c, tops, cops, uniswapAddresses, *ltp.FromETHAddress)
		if err != nil {
			return err
		}
	}

	startBlockNumber, err = c.BlockNumber(ctx)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to get current block number")
		return err
	}

	if inputLoadTestParams.StartNonce <= 0 {
		err = accountPool.RefreshNonce(ctx, tops.From)
		if err != nil {
			return err
		}
	}

	mustCheckMaxBaseFee, maxBaseFeeCtxCancel, waitBaseFeeToDrop := setupBaseFeeMonitoring(ctx, c, ltp)

	log.Debug().Msg("Starting main load test loop")
	var wg sync.WaitGroup

	// setup gas budget provider
	gasVault, gasPricer, err := setupGasManager(ctx, c)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Unable to setup gas manager")
		return err
	}

	for routineID := range maxRoutines {
		log.Trace().Int64("routineID", routineID).Msg("starting concurrent routine")
		wg.Add(1)
		go func(routineID int64) {
			var startReq time.Time
			var endReq time.Time
			var tErr error
			var ltTx *types.Transaction
			for requestID := range maxRequests {
				if rl != nil {
					tErr = rl.Wait(ctx)
					if tErr != nil {
						log.Error().
							Int64("routineID", routineID).
							Int64("requestID", requestID).
							Err(tErr).
							Msg("Encountered a rate limiting error")
					}
				}

				localMode := mode
				// if there are multiple modes, iterate through them, 'r' mode is supported here
				if ltp.MultiMode {
					localMode = ltp.ParsedModes[int(routineID+requestID)%(len(ltp.ParsedModes))]
				}
				// if we're doing random, we'll just pick one based on the current index
				if localMode == loadTestModeRandom {
					localMode = getRandomMode()
				}

				account, err := accountPool.Next(ctx)
				if err != nil {
					log.Error().
						Int64("routineID", routineID).
						Int64("requestID", requestID).
						Err(err).
						Msg("Unable to get next account from account pool")
					return
				}
				chainID := new(big.Int).SetUint64(ltp.ChainID)
				sendingTops, err := bind.NewKeyedTransactorWithChainID(account.privateKey, chainID)
				if err != nil {
					log.Error().
						Int64("routineID", routineID).
						Int64("requestID", requestID).
						Err(err).
						Msg("Unable create transaction signer")
					return
				}
				sendingTops.Nonce = new(big.Int).SetUint64(account.nonce)

				if mustCheckMaxBaseFee {
					waiting := false
					for waitBaseFeeToDrop.Load() {
						if !waiting {
							waiting = true
							log.Debug().
								Int64("routineID", routineID).
								Int64("requestID", requestID).
								Msg("go routine is waiting for base fee to drop")
						}
						time.Sleep(time.Second)
					}
				}

				sendingTops = configureTransactOpts(ctx, c, sendingTops, gasPricer)

				var fixedGasLimit uint64 = sendingTops.GasLimit

				// in case the gas limit is fixed, we spend or wait for the gas budget before sending the transaction
				if fixedGasLimit > 0 {
					log.Trace().Int64("routineID", routineID).
						Int64("requestID", requestID).
						Uint64("gas", fixedGasLimit).
						Msg("spending or waiting for fixed gas limit from gas budget")
					gasVault.SpendOrWaitAvailableBudget(fixedGasLimit)
					sendingTops.GasLimit = fixedGasLimit
				}

				switch localMode {
				case loadTestModeERC20:
					startReq, endReq, ltTx, tErr = loadTestERC20(ctx, c, sendingTops, erc20Contract, ltAddr)
				case loadTestModeERC721:
					startReq, endReq, ltTx, tErr = loadTestERC721(ctx, c, sendingTops, erc721Contract, ltAddr)
				case loadTestModeBlob:
					startReq, endReq, ltTx, tErr = loadTestBlob(ctx, c, sendingTops)
				case loadTestModeContractCall:
					startReq, endReq, ltTx, tErr = loadTestContractCall(ctx, c, sendingTops)
				case loadTestModeDeploy:
					startReq, endReq, ltTx, tErr = loadTestDeploy(ctx, c, sendingTops)
				case loadTestModeIncrement:
					startReq, endReq, ltTx, tErr = loadTestIncrement(ctx, c, sendingTops, ltContract)
				case loadTestModeRecall:
					startReq, endReq, ltTx, tErr = loadTestRecall(ctx, c, sendingTops, recallTransactions[int(sendingTops.Nonce.Uint64())%len(recallTransactions)])
				case loadTestModeRPC:
					startReq, endReq, tErr = loadTestRPC(ctx, c, indexedActivity)
				case loadTestModeStore:
					startReq, endReq, ltTx, tErr = loadTestStore(ctx, c, sendingTops, ltContract)
				case loadTestModeTransaction:
					startReq, endReq, ltTx, tErr = loadTestTransaction(ctx, c, sendingTops)
				case loadTestModeUniswapV3:
					swapAmountIn := big.NewInt(int64(uniswapv3LoadTestParams.SwapAmountInput))
					startReq, endReq, ltTx, tErr = runUniswapV3Loadtest(ctx, c, sendingTops, uniswapV3Config, poolConfig, swapAmountIn)
				default:
					log.Error().Str("mode", mode.String()).Msg("We've arrived at a load test mode that we don't recognize")
				}
				if !inputLoadTestParams.FireAndForget {
					recordSample(routineID, requestID, tErr, startReq, endReq, sendingTops.Nonce.Uint64())
				}
				if tErr == nil && inputLoadTestParams.WaitForReceipt {
					receiptMaxRetries := inputLoadTestParams.ReceiptRetryMax
					receiptRetryInitialDelayMs := inputLoadTestParams.ReceiptRetryInitialDelayMs
					_, tErr = util.WaitReceiptWithRetries(ctx, c, ltTx.Hash(), receiptMaxRetries, receiptRetryInitialDelayMs)
				}

				if tErr != nil {
					log.Error().
						Int64("routineID", routineID).
						Int64("requestID", requestID).
						Err(tErr).
						Str("mode", localMode.String()).
						Str("address", sendingTops.From.String()).
						Uint64("nonce", sendingTops.Nonce.Uint64()).
						Uint64("gas", sendingTops.GasLimit).
						Any("gasPrice", sendingTops.GasPrice).
						Any("gasFeeCap", sendingTops.GasFeeCap).
						Any("gasTipCap", sendingTops.GasTipCap).
						Int64("request time", endReq.Sub(startReq).Milliseconds()).
						Msg("recorded an error while sending transactions")

					// check nonce for reuse
					// if we're not in call only mode, we want to retry
					if !ltp.EthCallOnly {
						// we start setting nonce to be reused
						reuseNonce := true

						// if it is an error that consumes the nonce, we can't retry it
						if strings.Contains(tErr.Error(), "replacement transaction underpriced") ||
							strings.Contains(tErr.Error(), "transaction underpriced") ||
							strings.Contains(tErr.Error(), "nonce too low") ||
							strings.Contains(tErr.Error(), "already known") ||
							strings.Contains(tErr.Error(), "could not replace existing") {
							reuseNonce = false
						}

						// if we can reuse the nonce, we add it back to the account pool
						// for the specific account
						if reuseNonce {
							err := accountPool.AddReusableNonce(ctx, sendingTops.From, sendingTops.Nonce.Uint64())
							if err != nil {
								log.Error().
									Str("address", sendingTops.From.String()).
									Uint64("nonce", sendingTops.Nonce.Uint64()).
									Err(err).
									Msg("Unable to add reusable nonce to account pool")
							}
						}
					}
				} else {
					if ltTx != nil {
						// if gas limit was not fixed, we ask the vault to spend the gas used after the transaction was sent
						if fixedGasLimit == 0 {
							// log.Trace().Int64("routineID", routineID).
							// 	Int64("requestID", requestID).
							// 	Uint64("gas", ltTx.Gas()).
							// 	Msg("spending gas from gas budget after transaction is sent")
							gasVault.SpendOrWaitAvailableBudget(ltTx.Gas())
						}

						log.Trace().
							Int64("routineID", routineID).
							Int64("requestID", requestID).
							Stringer("txhash", ltTx.Hash()).
							Any("nonce", sendingTops.Nonce).
							Str("mode", localMode.String()).
							Str("sendingAddress", sendingTops.From.String()).
							Uint64("gas", ltTx.Gas()).
							Any("gasPrice", sendingTops.GasPrice).
							Any("gasFeeCap", sendingTops.GasFeeCap).
							Any("gasTipCap", sendingTops.GasTipCap).
							Msg("Request")
					}
				}
			}
			wg.Done()
		}(routineID)
	}
	log.Trace().Msg("Finished starting go routines. Waiting..")
	wg.Wait()

	rateLimitCancel()
	maxBaseFeeCtxCancel()
	if ltp.EthCallOnly {
		return nil
	}

	return nil
}

func setupGasManager(ctx context.Context, c *ethclient.Client) (*gasmanager.GasVault, *gasmanager.GasPricer, error) {
	gasVault, err := setupGasVault(ctx, c)
	if err != nil {
		return nil, nil, err
	}
	gasPricer, err := setupGasPricer()
	if err != nil {
		return nil, nil, err
	}

	return gasVault, gasPricer, nil
}

func setupGasVault(ctx context.Context, c *ethclient.Client) (*gasmanager.GasVault, error) {
	log.Trace().Msg("Setting up gas limiter")
	gasVault := gasmanager.NewGasVault()

	waveLog := log.Trace().
		Uint64("Period", inputLoadTestParams.GasManagerPeriod).
		Uint64("Amplitude", inputLoadTestParams.GasManagerAmplitude).
		Uint64("Target", inputLoadTestParams.GasManagerTarget)
	var wave gasmanager.Wave
	switch inputLoadTestParams.GasManagerOscillationWave {
	case "flat":
		waveLog.Msg("Using flat wave")
		wave = gasmanager.NewFlatWave(gasmanager.WaveConfig{
			Period:    inputLoadTestParams.GasManagerPeriod,
			Amplitude: inputLoadTestParams.GasManagerAmplitude,
			Target:    inputLoadTestParams.GasManagerTarget,
		})
	case "sine":
		waveLog.Msg("Using sine wave")
		wave = gasmanager.NewSineWave(gasmanager.WaveConfig{
			Period:    inputLoadTestParams.GasManagerPeriod,
			Amplitude: inputLoadTestParams.GasManagerAmplitude,
			Target:    inputLoadTestParams.GasManagerTarget,
		})
	case "sawtooth":
		waveLog.Msg("Using sawtooth wave")
		wave = gasmanager.NewSawtoothWave(gasmanager.WaveConfig{
			Period:    inputLoadTestParams.GasManagerPeriod,
			Amplitude: inputLoadTestParams.GasManagerAmplitude,
			Target:    inputLoadTestParams.GasManagerTarget,
		})
	case "square":
		waveLog.Msg("Using square wave")
		wave = gasmanager.NewSquareWave(gasmanager.WaveConfig{
			Period:    inputLoadTestParams.GasManagerPeriod,
			Amplitude: inputLoadTestParams.GasManagerAmplitude,
			Target:    inputLoadTestParams.GasManagerTarget,
		})
	case "triangle":
		waveLog.Msg("Using triangle wave")
		wave = gasmanager.NewTriangleWave(gasmanager.WaveConfig{
			Period:    inputLoadTestParams.GasManagerPeriod,
			Amplitude: inputLoadTestParams.GasManagerAmplitude,
			Target:    inputLoadTestParams.GasManagerTarget,
		})
	default:
		err := fmt.Errorf("unknown gas oscillation wave: %s", inputLoadTestParams.GasManagerOscillationWave)
		return nil, err
	}

	gasProvider := gasmanager.NewOscillatingGasProvider(c, gasVault, wave)
	gasProvider.Start(ctx)

	return gasVault, nil
}

func setupGasPricer() (*gasmanager.GasPricer, error) {
	log.Trace().Msg("Setting up gas pricer")
	var strategy gasmanager.PriceStrategy
	switch inputLoadTestParams.GasManagerPriceStrategy {
	case "fixed":
		log.Trace().Msg("Using fixed gas price strategy")
		strategy = gasmanager.NewFixedGasPriceStrategy(gasmanager.FixedGasPriceConfig{
			GasPriceWei: inputLoadTestParams.GasManagerFixedGasPriceWei,
		})
	case "estimated":
		log.Trace().Msg("Using estimated gas price strategy")
		strategy = gasmanager.NewEstimatedGasPriceStrategy()
	case "dynamic":
		log.Trace().Msg("Using dynamic gas price strategy")

		gasPricesArr := strings.Split(inputLoadTestParams.GasManagerDynamicGasPricesWei, ",")
		var gasPrices []uint64
		if len(gasPricesArr) > 0 {
			for _, gpStr := range gasPricesArr {
				gp, err := strconv.ParseUint(strings.TrimSpace(gpStr), 10, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid gas price in dynamic gas prices list: %s", gpStr)
				}
				gasPrices = append(gasPrices, gp)
			}
			log.Trace().
				Any("GasPrices", gasPrices).
				Msg("Using custom dynamic gas prices")
		}

		strategy = gasmanager.NewDynamicGasPriceStrategy(gasmanager.DynamicGasPriceConfig{
			GasPrices: gasPrices,
			Variation: inputLoadTestParams.GasManagerDynamicGasPricesVariation,
		})
	default:
		return nil, fmt.Errorf("unknown gas price strategy: %s", inputLoadTestParams.GasManagerPriceStrategy)
	}

	gasPricer := gasmanager.NewGasPricer(strategy)
	return gasPricer, nil
}

func setupBaseFeeMonitoring(ctx context.Context, c *ethclient.Client, ltp loadTestParams) (bool, context.CancelFunc, *atomic.Bool) {
	// monitor max base fee if configured
	maxBaseFeeCtx, maxBaseFeeCtxCancel := context.WithCancel(ctx)
	mustCheckMaxBaseFee := ltp.MaxBaseFeeWei > 0
	var waitBaseFeeToDrop atomic.Bool
	waitBaseFeeToDrop.Store(false)
	if mustCheckMaxBaseFee {
		log.Info().
			Msg("max base fee monitoring enabled")

		wg := sync.WaitGroup{}
		wg.Add(1)
		// start a goroutine to monitor the base fee while load test is running
		go func(ctx context.Context, c *ethclient.Client, waitToDrop *atomic.Bool, maxBaseFeeWei uint64) {
			firstRun := true
			for {
				select {
				case <-ctx.Done():
					return
				default:
					currentBaseFeeIsGreaterThanMax, currentBaseFeeWei, err := isCurrentBaseFeeGreaterThanMaxBaseFee(ctx, c, maxBaseFeeWei)
					if err != nil {
						log.Error().
							Err(err).
							Msg("Error checking base fee during load test")
					} else {
						if currentBaseFeeIsGreaterThanMax {
							if !waitToDrop.Load() {
								log.Warn().
									Msgf("PAUSE: base fee %d Wei > limit %d Wei", currentBaseFeeWei.Uint64(), maxBaseFeeWei)
								waitToDrop.Store(true)
							}
						} else if waitToDrop.Load() {
							log.Info().
								Msgf("RESUME: base fee %d Wei ≤ limit %d Wei", currentBaseFeeWei.Uint64(), maxBaseFeeWei)
							waitToDrop.Store(false)
						}

						if firstRun {
							firstRun = false
							wg.Done()
						}
					}
					time.Sleep(time.Second)
				}
			}
		}(maxBaseFeeCtx, c, &waitBaseFeeToDrop, ltp.MaxBaseFeeWei)

		// wait for first run to complete so we know if we need to wait or not for base fee to drop
		wg.Wait()
	}
	return mustCheckMaxBaseFee, maxBaseFeeCtxCancel, &waitBaseFeeToDrop
}

func isCurrentBaseFeeGreaterThanMaxBaseFee(ctx context.Context, c *ethclient.Client, maxBaseFee uint64) (bool, *big.Int, error) {
	header, err := c.HeaderByNumber(ctx, nil)
	if errors.Is(err, context.Canceled) {
		log.Debug().Msg("max base fee monitoring context canceled")
		return false, nil, nil
	} else if err != nil {
		log.Error().Err(err).Msg("Unable to get latest block header to check base fee")
		return false, nil, err
	}

	if header.BaseFee != nil {
		currentBaseFee := header.BaseFee
		if currentBaseFee.Cmp(new(big.Int).SetUint64(maxBaseFee)) > 0 {
			return true, currentBaseFee, nil
		} else {
			return false, currentBaseFee, nil
		}
	}

	return false, nil, nil
}

func getLoadTestContract(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts, cops *bind.CallOpts) (ltAddr ethcommon.Address, ltContract *tester.LoadTester, err error) {
	ltAddr = ethcommon.HexToAddress(inputLoadTestParams.LoadtestContractAddress)

	if inputLoadTestParams.LoadtestContractAddress == "" {
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
	erc20Addr = ethcommon.HexToAddress(inputLoadTestParams.ERC20Address)
	if inputLoadTestParams.ERC20Address == "" {
		log.Info().Msg("Deploying ERC20 contract")
		erc20Addr, _, _, err = tokens.DeployERC20(tops, c)
		if err != nil {
			log.Error().Err(err).Msg("Unable to deploy ERC20 contract")
			return
		}
		// Tokens already minted and sent to the address of the deployer.
	}
	log.Info().Interface("contractaddress", erc20Addr).Msg("ERC20 contract address")

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
	erc721Addr = ethcommon.HexToAddress(inputLoadTestParams.ERC721Address)
	shouldMint := true
	if inputLoadTestParams.ERC721Address == "" {
		log.Info().Msg("Deploying ERC721 contract")
		erc721Addr, _, _, err = tokens.DeployERC721(tops, c)
		if err != nil {
			log.Error().Err(err).Msg("Unable to deploy ERC721 contract")
			return
		}
		shouldMint = false
	}
	log.Info().Interface("contractaddress", erc721Addr).Msg("ERC721 contract address")

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

	log.Info().Msg("Mint one ERC721 token")
	err = util.BlockUntilSuccessful(ctx, c, func() error {
		_, err = erc721Contract.MintBatch(tops, *inputLoadTestParams.FromETHAddress, new(big.Int).SetUint64(1))
		return err
	})
	return
}

func loadTestTransaction(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts) (t1 time.Time, t2 time.Time, tx *types.Transaction, err error) {
	ltp := inputLoadTestParams

	to := ltp.ToETHAddress
	if ltp.RandomRecipients {
		to = getRandomAddress()
	}

	const eoaTransferGasLimit = 21000
	if tops.GasLimit == 0 {
		tops.GasLimit = uint64(eoaTransferGasLimit)
	}

	amount := ltp.SendAmount
	chainID := new(big.Int).SetUint64(ltp.ChainID)

	var rtx *types.Transaction
	if ltp.LegacyTransactionMode {
		rtx = ethtypes.NewTx(&ethtypes.LegacyTx{
			Nonce:    tops.Nonce.Uint64(),
			To:       to,
			Value:    amount,
			Gas:      tops.GasLimit,
			GasPrice: tops.GasPrice,
			Data:     nil,
		})
	} else {
		dynamicFeeTx := &ethtypes.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     tops.Nonce.Uint64(),
			To:        to,
			Gas:       tops.GasLimit,
			GasFeeCap: tops.GasFeeCap,
			GasTipCap: tops.GasTipCap,
			Data:      nil,
			Value:     amount,
		}
		rtx = ethtypes.NewTx(dynamicFeeTx)
	}

	tx, err = tops.Signer(tops.From, rtx)
	if err != nil {
		log.Error().Err(err).Msg("Unable to sign transaction")
		return
	}

	t1 = time.Now()
	defer func() { t2 = time.Now() }()
	if ltp.EthCallOnly {
		_, err = c.CallContract(ctx, txToCallMsg(tx), nil)
	} else if ltp.OutputRawTxOnly {
		err = outputRawTransaction(tx)
	} else {
		err = c.SendTransaction(ctx, tx)
	}

	return
}

var (
	cachedBlockNumber           *uint64
	cachedGasPriceLock          sync.Mutex
	cachedGasPrice              *big.Int
	cachedGasTipCap             *big.Int
	cachedLatestBlockNumber     uint64
	cachedLatestBlockTime       time.Time
	cachedLatestBlockNumberLock sync.Mutex
)

func getLatestBlockNumber(ctx context.Context, c *ethclient.Client) uint64 {
	cachedLatestBlockNumberLock.Lock()
	defer cachedLatestBlockNumberLock.Unlock()
	// The case where cachedLatestBlockTime is empty should give a large Since value and cause the block to be fetched
	if time.Since(cachedLatestBlockTime) < 1*time.Second {
		return cachedLatestBlockNumber
	}
	bn, err := c.BlockNumber(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Unable to get block number while checking gas prices")
		return 0
	}
	cachedLatestBlockTime = time.Now()
	cachedLatestBlockNumber = bn
	return bn
}

func biasGasPrice(price *big.Int) *big.Int {
	gasPriceFloat := new(big.Float).SetInt(price)
	gasPriceFloat.Mul(gasPriceFloat, inputLoadTestParams.BigGasPriceMultiplier)
	result := new(big.Int)
	gasPriceFloat.Int(result)
	return result
}

func getSuggestedGasPrices(ctx context.Context, c *ethclient.Client, gasPricer *gasmanager.GasPricer) (*big.Int, *big.Int) {
	cachedGasPriceLock.Lock()
	defer cachedGasPriceLock.Unlock()

	// this should be one of the fastest RPC calls, so hopefully there isn't too much overhead calling this
	bn := getLatestBlockNumber(ctx, c)

	if gasPricer == nil { // cache is used only when gas pricer is not used
		if cachedBlockNumber != nil && bn <= *cachedBlockNumber {
			return cachedGasPrice, cachedGasTipCap
		}
	}

	// In the case of an EVM compatible system not supporting EIP-1559
	var gasPrice, gasTipCap = big.NewInt(0), big.NewInt(0)
	var pErr, tErr error
	if inputLoadTestParams.LegacyTransactionMode {
		if inputLoadTestParams.ForceGasPrice != 0 {
			gasPrice = new(big.Int).SetUint64(inputLoadTestParams.ForceGasPrice)
		} else {
			var gp *uint64
			if gasPricer != nil {
				gp = gasPricer.GetGasPrice()
			}
			if gp != nil {
				gasPrice = big.NewInt(0).SetUint64(*gp)
			} else {
				if cachedBlockNumber != nil && bn <= *cachedBlockNumber {
					return cachedGasPrice, cachedGasTipCap
				}
				gasPrice, pErr = c.SuggestGasPrice(ctx)
				if pErr != nil {
					log.Error().Err(pErr).Msg("Unable to suggest gas price")
					return cachedGasPrice, cachedGasTipCap
				}
			}
			// Bias the value up slightly
			gasPrice = biasGasPrice(gasPrice)
		}
	} else {
		var forcePriorityGasPrice *big.Int
		if inputLoadTestParams.ForcePriorityGasPrice != 0 {
			gasTipCap = new(big.Int).SetUint64(inputLoadTestParams.ForcePriorityGasPrice)
			forcePriorityGasPrice = gasTipCap
		} else if inputLoadTestParams.ChainSupportBaseFee {
			if cachedBlockNumber != nil && bn <= *cachedBlockNumber {
				gasTipCap = cachedGasTipCap
			} else {
				gasTipCap, tErr = c.SuggestGasTipCap(ctx)
				if tErr != nil {
					log.Error().Err(tErr).Msg("Unable to suggest gas tip cap")
					return cachedGasPrice, cachedGasTipCap
				}
				// Bias the value up slightly
				gasTipCap = biasGasPrice(gasTipCap)
			}
		} else {
			log.Fatal().
				Msg("Chain does not support base fee. Please set priority-gas-price flag with a value to use for gas tip cap")
		}

		if inputLoadTestParams.ForceGasPrice != 0 {
			gasPrice = new(big.Int).SetUint64(inputLoadTestParams.ForceGasPrice)
		} else if inputLoadTestParams.ChainSupportBaseFee {
			var gp *uint64
			if gasPricer != nil {
				gp = gasPricer.GetGasPrice()
			}
			if gp != nil {
				gasPrice = big.NewInt(0).SetUint64(*gp)
			} else {
				if cachedBlockNumber != nil && bn <= *cachedBlockNumber {
					return cachedGasPrice, cachedGasTipCap
				}
				gasPrice = suggestMaxFeePerGas(ctx, c, bn, forcePriorityGasPrice)
			}
		} else {
			log.Fatal().
				Msg("Chain does not support base fee. Please set gas-price flag with a value to use for max fee per gas")
		}
	}

	cachedBlockNumber = &bn
	cachedGasPrice = gasPrice
	cachedGasTipCap = gasTipCap

	l := log.Debug().
		Uint64("cachedBlockNumber", bn)

	if cachedGasPrice != nil {
		l = l.Uint64("cachedGasPrice", cachedGasPrice.Uint64())
	} else {
		l = l.Interface("cachedGasPrice", cachedGasPrice)
	}

	if cachedGasTipCap != nil {
		l = l.Uint64("cachedGasTipCap", cachedGasTipCap.Uint64())
	} else {
		l = l.Interface("cachedGasTipCap", cachedGasTipCap)
	}

	if gasPricer == nil {
		// only log when cache is used
		l.Msg("Updating gas prices")
	}

	return cachedGasPrice, cachedGasTipCap
}

func suggestMaxFeePerGas(ctx context.Context, c *ethclient.Client, blockNumber uint64, forcePriorityFee *big.Int) *big.Int {
	iHeader, iErr := c.HeaderByNumber(ctx, nil)
	if iErr != nil {
		log.Error().Err(iErr).Msg("Unable to get latest block header while checking MaxFeePerGas")
		return nil
	}

	if cachedBlockNumber != nil && blockNumber <= *cachedBlockNumber && cachedGasPrice != nil {
		return cachedGasPrice
	}

	feeHistory, iErr := c.FeeHistory(ctx, 5, nil, []float64{0.5})
	if iErr != nil {
		log.Error().Err(iErr).Msg("Unable to get fee history while checking MaxFeePerGas")
		return nil
	}

	priorityFee := forcePriorityFee
	if priorityFee == nil {
		priorityFee = feeHistory.Reward[len(feeHistory.Reward)-1][0] // 50th percentile of most recent block
	}
	baseFee := feeHistory.BaseFee[len(feeHistory.BaseFee)-1] // base fee of next block
	maxFeePerGas := new(big.Int)
	maxFeePerGas.Mul(baseFee, big.NewInt(2))
	maxFeePerGas.Add(maxFeePerGas, priorityFee)

	// in the case of a decreasing fee, the update happens only
	// after some blocks to avoid the network fee fluctuations
	const blocksToWait = 5
	isDecreasing := cachedGasPrice != nil && maxFeePerGas.Uint64() <= cachedGasPrice.Uint64()
	canDecrease := blockNumber+blocksToWait <= iHeader.Number.Uint64()
	if isDecreasing && !canDecrease && cachedGasPrice != nil {
		return cachedGasPrice
	}

	cachedGasPrice = maxFeePerGas

	log.Trace().
		Uint64("blockNumber", iHeader.Number.Uint64()).
		Str("priorityFee", priorityFee.String()).
		Str("baseFee", baseFee.String()).
		Str("maxFeePerGas", maxFeePerGas.String()).
		Uint64("cachedGasPrice", cachedGasPrice.Uint64()).
		Msg("max fee updated")

	return maxFeePerGas
}

// TODO - in the future it might be more interesting if this mode takes input or random contracts to be deployed
func loadTestDeploy(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts) (t1 time.Time, t2 time.Time, tx *types.Transaction, err error) {
	ltp := inputLoadTestParams

	t1 = time.Now()
	defer func() { t2 = time.Now() }()
	if ltp.EthCallOnly {
		msg := transactOptsToCallMsg(tops)
		msg.Data = ethcommon.FromHex(tester.LoadTesterMetaData.Bin)
		_, err = c.CallContract(ctx, msg, nil)
	} else if ltp.OutputRawTxOnly {
		// For raw output, we need to manually create and sign the deployment transaction
		tops.NoSend = true
		_, tx, _, err = tester.DeployLoadTester(tops, c)
		if err != nil {
			return
		}
		// The transaction from DeployLoadTester should already be signed
		if tx != nil {
			err = outputRawTransaction(tx)
		}
	} else {
		_, tx, _, err = tester.DeployLoadTester(tops, c)
	}
	return
}

func loadTestIncrement(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts, ltContract *tester.LoadTester) (t1 time.Time, t2 time.Time, tx *types.Transaction, err error) {
	ltp := inputLoadTestParams

	t1 = time.Now()
	defer func() { t2 = time.Now() }()
	if ltp.EthCallOnly {
		tops.NoSend = true
		tx, err = ltContract.Inc(tops)
		if err != nil {
			return
		}
		msg := txToCallMsg(tx)
		_, err = c.CallContract(ctx, msg, nil)
	} else if ltp.OutputRawTxOnly {
		tops.NoSend = true
		tx, err = ltContract.Inc(tops)
		if err != nil {
			return
		}
		// Sign the transaction manually since NoSend was true
		signedTx, signErr := tops.Signer(tops.From, tx)
		if signErr != nil {
			err = signErr
			return
		}
		err = outputRawTransaction(signedTx)
	} else {
		tx, err = ltContract.Inc(tops)
	}
	return
}

func loadTestStore(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts, ltContract *tester.LoadTester) (t1 time.Time, t2 time.Time, tx *types.Transaction, err error) {
	ltp := inputLoadTestParams

	inputData := make([]byte, ltp.StoreDataSize)
	_, _ = hexwordRead(inputData)
	t1 = time.Now()
	defer func() { t2 = time.Now() }()
	if ltp.EthCallOnly {
		tops.NoSend = true
		tx, err = ltContract.Store(tops, inputData)
		if err != nil {
			return
		}
		msg := txToCallMsg(tx)
		_, err = c.CallContract(ctx, msg, nil)
	} else if ltp.OutputRawTxOnly {
		tops.NoSend = true
		tx, err = ltContract.Store(tops, inputData)
		if err != nil {
			return
		}
		// Sign the transaction manually since NoSend was true
		signedTx, signErr := tops.Signer(tops.From, tx)
		if signErr != nil {
			err = signErr
			return
		}
		err = outputRawTransaction(signedTx)
	} else {
		tx, err = ltContract.Store(tops, inputData)
	}
	return
}

func loadTestERC20(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts, erc20Contract *tokens.ERC20, ltAddress ethcommon.Address) (t1 time.Time, t2 time.Time, tx *types.Transaction, err error) {
	ltp := inputLoadTestParams

	to := ltp.ToETHAddress
	if ltp.RandomRecipients {
		to = getRandomAddress()
	}
	amount := ltp.SendAmount

	t1 = time.Now()
	defer func() { t2 = time.Now() }()
	if ltp.EthCallOnly {
		tops.NoSend = true
		tx, err = erc20Contract.Transfer(tops, *to, amount)
		if err != nil {
			return
		}
		msg := txToCallMsg(tx)
		_, err = c.CallContract(ctx, msg, nil)
	} else if ltp.OutputRawTxOnly {
		tops.NoSend = true
		tx, err = erc20Contract.Transfer(tops, *to, amount)
		if err != nil {
			return
		}
		// Sign the transaction manually since NoSend was true
		signedTx, signErr := tops.Signer(tops.From, tx)
		if signErr != nil {
			err = signErr
			return
		}
		err = outputRawTransaction(signedTx)
	} else {
		tx, err = erc20Contract.Transfer(tops, *to, amount)
	}

	return
}

func loadTestERC721(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts, erc721Contract *tokens.ERC721, ltAddress ethcommon.Address) (t1 time.Time, t2 time.Time, tx *types.Transaction, err error) {
	ltp := inputLoadTestParams

	to := ltp.ToETHAddress
	if ltp.RandomRecipients {
		to = getRandomAddress()
	}

	t1 = time.Now()
	defer func() { t2 = time.Now() }()
	if ltp.EthCallOnly {
		tops.NoSend = true
		tx, err = erc721Contract.MintBatch(tops, *to, big.NewInt(1))
		if err != nil {
			return
		}
		msg := txToCallMsg(tx)
		_, err = c.CallContract(ctx, msg, nil)
	} else if ltp.OutputRawTxOnly {
		tops.NoSend = true
		tx, err = erc721Contract.MintBatch(tops, *to, big.NewInt(1))
		if err != nil {
			return
		}
		// Sign the transaction manually since NoSend was true
		signedTx, signErr := tops.Signer(tops.From, tx)
		if signErr != nil {
			err = signErr
			return
		}
		err = outputRawTransaction(signedTx)
	} else {
		tx, err = erc721Contract.MintBatch(tops, *to, big.NewInt(1))
	}

	return
}

func loadTestRecall(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts, originalTx rpctypes.PolyTransaction) (t1 time.Time, t2 time.Time, tx *types.Transaction, err error) {
	ltp := inputLoadTestParams

	// For EIP-1559 transactions, use GasFeeCap instead of GasPrice (which is nil for dynamic fee txs)
	gasPrice := tops.GasPrice
	if gasPrice == nil && tops.GasFeeCap != nil {
		gasPrice = tops.GasFeeCap
	}
	rtx := rawTransactionToNewTx(originalTx, tops.Nonce.Uint64(), gasPrice, tops.GasTipCap)

	tx, err = tops.Signer(tops.From, rtx)
	if err != nil {
		log.Error().Err(err).Msg("Unable to sign transaction")
		return
	}
	log.Trace().Str("txId", originalTx.Hash().String()).Bool("callOnly", ltp.EthCallOnly).Msg("Attempting to replay transaction")

	t1 = time.Now()
	defer func() { t2 = time.Now() }()
	if ltp.EthCallOnly {
		callMsg := txToCallMsg(tx)
		callMsg.From = originalTx.From()
		callMsg.Gas = originalTx.Gas()
		if ltp.EthCallOnlyLatestBlock {
			_, err = c.CallContract(ctx, callMsg, nil)
		} else {
			callMsg.GasFeeCap = new(big.Int).SetUint64(originalTx.MaxFeePerGas())
			callMsg.GasTipCap = new(big.Int).SetUint64(originalTx.MaxPriorityFeePerGas())
			if originalTx.MaxFeePerGas() == 0 && originalTx.MaxPriorityFeePerGas() == 0 {
				callMsg.GasPrice = originalTx.GasPrice()
				callMsg.GasFeeCap = nil
				callMsg.GasTipCap = nil
			} else {
				callMsg.GasPrice = nil
			}

			_, err = c.CallContract(ctx, callMsg, originalTx.BlockNumber())
		}
		if err != nil {
			log.Warn().Err(err).Msg("Recall failure")
		}
		// we're not going to return the error in the case because there is no point retrying
		err = nil
	} else if ltp.OutputRawTxOnly {
		err = outputRawTransaction(tx)
	} else {
		err = c.SendTransaction(ctx, tx)
	}
	return
}

func loadTestRPC(ctx context.Context, c *ethclient.Client, ia *IndexedActivity) (t1 time.Time, t2 time.Time, err error) {
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
		var txArgs apitypes.SendTxArgs
		if err = json.Unmarshal(rawTxData, &txArgs); err != nil {
			log.Error().Err(err).Str("txHash", pt.Hash().String()).Msg("unable to unmarshal poly transaction to json")
			return
		}
		var tx *ethtypes.Transaction
		tx, err = txArgs.ToTransaction()
		if err != nil {
			log.Error().Err(err).Str("txArgs", txArgs.String()).Msg("unable to convert the arguments to a transaction")
			return
		}
		cm := txToCallMsg(tx)
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
				Stringer("erc20addr", erc20Addr).
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
				Stringer("erc721addr", erc721Addr).
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

func loadTestContractCall(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts) (t1 time.Time, t2 time.Time, tx *types.Transaction, err error) {
	var calldata []byte

	ltp := inputLoadTestParams

	to := ltp.ContractETHAddress
	chainID := new(big.Int).SetUint64(ltp.ChainID)
	amount := big.NewInt(0)
	if ltp.ContractCallPayable {
		amount = ltp.SendAmount
	}

	if inputLoadTestParams.ContractCallData == "" {
		err = fmt.Errorf("missing calldata for function call")
		log.Error().Err(err).Msg("--calldata flag is required for contract-call mode")
		return
	}

	calldata, err = hex.DecodeString(strings.TrimPrefix(inputLoadTestParams.ContractCallData, "0x"))
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

	var rtx *types.Transaction
	if ltp.LegacyTransactionMode {
		rtx = ethtypes.NewTx(&ethtypes.LegacyTx{
			Nonce:    tops.Nonce.Uint64(),
			To:       to,
			Value:    amount,
			Gas:      tops.GasLimit,
			GasPrice: tops.GasPrice,
			Data:     calldata,
		})
	} else {
		rtx = ethtypes.NewTx(&ethtypes.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     tops.Nonce.Uint64(),
			To:        to,
			Gas:       tops.GasLimit,
			GasFeeCap: tops.GasFeeCap,
			GasTipCap: tops.GasTipCap,
			Data:      calldata,
			Value:     amount,
		})
	}
	log.Trace().Interface("rtx", rtx).Msg("Contract call data")

	tx, err = tops.Signer(tops.From, rtx)
	if err != nil {
		log.Error().Err(err).Msg("Unable to sign transaction")
		return
	}

	t1 = time.Now()
	defer func() { t2 = time.Now() }()
	if ltp.EthCallOnly {
		_, err = c.CallContract(ctx, txToCallMsg(tx), nil)
	} else if ltp.OutputRawTxOnly {
		err = outputRawTransaction(tx)
	} else {
		err = c.SendTransaction(ctx, tx)
	}
	return
}

func loadTestBlob(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts) (t1 time.Time, t2 time.Time, tx *types.Transaction, err error) {
	ltp := inputLoadTestParams

	to := ltp.ToETHAddress
	if ltp.RandomRecipients {
		to = getRandomAddress()
	}

	amount := ltp.SendAmount
	chainID := new(big.Int).SetUint64(ltp.ChainID)

	gasLimit := uint64(21000)
	// Use the gas values from tops which have been properly configured by configureTransactOpts
	// This ensures we respect ForceGasPrice, ForcePriorityGasPrice, and other overrides
	blobFeeCap := ltp.BlobFeeCap

	// Blob transactions require EIP-1559 support
	if tops.GasFeeCap == nil || tops.GasTipCap == nil {
		err = fmt.Errorf("blob transactions require EIP-1559 support (non-legacy mode)")
		log.Error().Err(err).Msg("Cannot send blob transaction in legacy mode")
		return
	}

	// Initialize blobTx with blob transaction type
	blobTx := ethtypes.BlobTx{
		ChainID:    uint256.NewInt(chainID.Uint64()),
		Nonce:      tops.Nonce.Uint64(),
		GasTipCap:  uint256.NewInt(tops.GasTipCap.Uint64()),
		GasFeeCap:  uint256.NewInt(tops.GasFeeCap.Uint64()),
		BlobFeeCap: uint256.NewInt(blobFeeCap),
		Gas:        gasLimit,
		To:         *to,
		Value:      uint256.NewInt(amount.Uint64()),
		Data:       nil,
		AccessList: nil,
		BlobHashes: make([]ethcommon.Hash, 0),
		Sidecar: &ethtypes.BlobTxSidecar{
			Blobs:       make([]kzg4844.Blob, 0),
			Commitments: make([]kzg4844.Commitment, 0),
			Proofs:      make([]kzg4844.Proof, 0),
		},
	}
	// appendBlobCommitment() will take in the blobTx struct and append values to blob transaction specific keys in the following steps:
	// The function will take in blobTx with empty BlobHashes, and Blob Sidecar variables initially.
	// Then generateRandomBlobData() is called to generate a byte slice with random values.
	// createBlob() is called to commit the randomly generated byte slice with KZG.
	// generateBlobCommitment() will do the same for the Commitment and Proof.
	// Append all the blob related computed values to the blobTx struct.
	err = appendBlobCommitment(&blobTx)
	if err != nil {
		log.Error().Err(err).Msg("Unable to parse blob")
		return
	}
	rtx := ethtypes.NewTx(&blobTx)

	tx, err = tops.Signer(tops.From, rtx)
	if err != nil {
		log.Error().Err(err).Msg("Unable to sign transaction")
		return
	}

	t1 = time.Now()
	defer func() { t2 = time.Now() }()
	if ltp.EthCallOnly {
		log.Error().Err(err).Msg("CallOnly not supported to blob transactions")
		return
	} else if ltp.OutputRawTxOnly {
		err = outputRawTransaction(tx)
	} else {
		err = c.SendTransaction(ctx, tx)
	}
	return
}

// outputRawTransaction marshals a signed transaction to hex and outputs it to stdout
func outputRawTransaction(stx *ethtypes.Transaction) error {
	rawTx, err := stx.MarshalBinary()
	if err != nil {
		log.Error().Err(err).Msg("Unable to marshal transaction to binary")
		return err
	}

	rawTxHex := "0x" + hex.EncodeToString(rawTx)
	fmt.Println(rawTxHex)

	return nil
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
	loadTestResultsMutex.Lock()
	loadTestResults = append(loadTestResults, s)
	loadTestResultsMutex.Unlock()
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

func configureTransactOpts(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts, gasPricer *gasmanager.GasPricer) *bind.TransactOpts {
	gasPrice, gasTipCap := getSuggestedGasPrices(ctx, c, gasPricer)

	tops.GasPrice = gasPrice

	ltp := inputLoadTestParams

	if ltp.ForceGasPrice != 0 {
		tops.GasPrice = big.NewInt(0).SetUint64(ltp.ForceGasPrice)
	}
	if ltp.ForceGasLimit != 0 {
		tops.GasLimit = ltp.ForceGasLimit
	}

	// if we're in legacy mode, there's no point doing anything else in this function
	if ltp.LegacyTransactionMode {
		return tops
	}
	if !ltp.ChainSupportBaseFee {
		log.Fatal().Msg("EIP-1559 not activated. Please use --legacy")
	}

	tops.GasPrice = nil
	tops.GasFeeCap = gasPrice
	tops.GasTipCap = gasTipCap

	if tops.GasTipCap.Cmp(tops.GasFeeCap) == 1 {
		tops.GasTipCap = new(big.Int).Set(tops.GasFeeCap)
	}

	return tops
}

func waitForFinalBlock(ctx context.Context, c *ethclient.Client, rpc *ethrpc.Client, startBlockNumber uint64) (uint64, error) {
	ltp := inputLoadTestParams
	var err error
	var lastBlockNumber uint64
	var checkInterval = 5 * time.Second
	var maxRetries = 30

	rateLimiter := rate.NewLimiter(rate.Limit(ltp.RateLimit), 1)
	noncesToCheck := accountPool.Nonces(ctx, true)

	retry := 0
	for {
		retry++
		if retry > maxRetries {
			log.Error().Msg("Max retries reached. Exiting...")
			return 0, fmt.Errorf("max retries reached")
		}
		lastBlockNumber, err = c.BlockNumber(ctx)
		if err != nil {
			return 0, err
		}
		if ltp.EthCallOnly {
			return lastBlockNumber, nil
		}

		wg := sync.WaitGroup{}
		remainingNoncesToCheck := atomic.Int64{}
		noncesToCheck.Range(func(key, value any) bool {
			wg.Add(1)
			remainingNoncesToCheck.Add(1)
			address := key.(ethcommon.Address)
			expectedNonce := value.(uint64)
			go func(ctx context.Context, rl *rate.Limiter) {
				defer wg.Done()
				err := rl.Wait(ctx)
				if err != nil {
					log.Error().Err(err).Msg("Rate limiter wait error")
					return
				}
				nonce, err := c.NonceAt(ctx, address, new(big.Int).SetUint64(lastBlockNumber))
				if err != nil {
					log.Error().Err(err).Str("address", address.String()).Msg("Unable to get nonce for account while checking for final block")
					return
				}
				logEvent := log.Debug().
					Str("address", address.String()).
					Uint64("nonce", nonce).
					Uint64("expectedNonce", expectedNonce).
					Uint64("lastBlockNumber", lastBlockNumber)
				if nonce < expectedNonce {
					logEvent.Msg("not all transactions for account have been mined. waiting...")
				} else {
					remainingNoncesToCheck.Add(-1)
					noncesToCheck.Delete(address)
				}
			}(ctx, rateLimiter)
			return true
		})
		wg.Wait()

		if remainingNoncesToCheck.Load() == 0 {
			log.Debug().Msg("All transactions of all accounts have been mined")
			break
		}

		log.Debug().
			Int("maxRetries", maxRetries).
			Int("retry", retry).
			Msgf("Retrying in %s...", checkInterval.String())
		time.Sleep(checkInterval)
	}

	log.Debug().
		Uint64("startblock", startBlockNumber).
		Uint64("endblock", lastBlockNumber).
		Msg("It looks like all transactions have been mined")
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
