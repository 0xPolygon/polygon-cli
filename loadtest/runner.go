package loadtest

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/0xPolygon/polygon-cli/bindings/tester"
	"github.com/0xPolygon/polygon-cli/bindings/tokens"
	"github.com/0xPolygon/polygon-cli/loadtest/config"
	"github.com/0xPolygon/polygon-cli/loadtest/mode"
	"github.com/0xPolygon/polygon-cli/loadtest/modes"
	"github.com/0xPolygon/polygon-cli/loadtest/uniswapv3"
	"github.com/0xPolygon/polygon-cli/util"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

// Runner handles the execution of a load test.
type Runner struct {
	cfg         *config.Config
	accountPool *AccountPool
	deps        *mode.Dependencies
	results     []Sample
	resultsMu   sync.RWMutex
	rl          *rate.Limiter
	randSrc     *rand.Rand

	startBlockNumber uint64
	finalBlockNumber uint64

	// Mode execution
	modes             []mode.Runner
	waitBaseFeeToDrop atomic.Bool

	// Clients
	client    *ethclient.Client
	rpcClient *ethrpc.Client

	// Gas price caching
	cachedBlockNumber       *uint64
	cachedGasPriceLock      sync.Mutex
	cachedGasPrice          *big.Int
	cachedGasTipCap         *big.Int
	cachedLatestBlockNumber uint64
	cachedLatestBlockTime   time.Time
	cachedLatestBlockLock   sync.Mutex
}

// NewRunner creates a new Runner with the given configuration.
func NewRunner(cfg *config.Config) (*Runner, error) {
	return &Runner{
		cfg:     cfg,
		results: make([]Sample, 0),
		randSrc: rand.New(rand.NewSource(cfg.Seed)),
	}, nil
}

// Init sets up the runner, including clients and account pool.
func (r *Runner) Init(ctx context.Context) error {
	log.Info().Msg("Initializing load test runner")

	// Configure HTTP transport
	connLimit := 2 * int(r.cfg.Concurrency)
	transport := &http.Transport{
		MaxIdleConns:        connLimit,
		MaxIdleConnsPerHost: connLimit,
		MaxConnsPerHost:     connLimit,
	}
	if r.cfg.Proxy != "" {
		proxyURL, err := url.Parse(r.cfg.Proxy)
		if err != nil {
			return errors.New("invalid proxy address: " + r.cfg.Proxy + ": " + err.Error())
		}
		transport.Proxy = http.ProxyURL(proxyURL)
		log.Debug().Stringer("proxyURL", proxyURL).Msg("transport proxy configured")
	}

	goHTTPClient := &http.Client{Transport: transport}
	rpcOption := ethrpc.WithHTTPClient(goHTTPClient)
	rpc, err := ethrpc.DialOptions(ctx, r.cfg.RPCUrl, rpcOption)
	if err != nil {
		return errors.New("unable to dial rpc: " + err.Error())
	}
	rpc.SetHeader("Accept-Encoding", "identity")
	r.rpcClient = rpc
	r.client = ethclient.NewClient(rpc)

	// Initialize dependencies early so mode parsing can use them.
	r.deps = &mode.Dependencies{
		Client:     r.client,
		RPCClient:  r.rpcClient,
		RandSource: r.randSrc,
	}

	// Initialize load test parameters
	if err := r.initParams(ctx); err != nil {
		return err
	}

	return nil
}

func (r *Runner) initParams(ctx context.Context) error {
	log.Info().Msg("Connecting with RPC endpoint to initialize load test parameters")

	// When outputting raw transactions, we don't need to wait for anything to be mined
	if r.cfg.OutputRawTxOnly {
		r.cfg.FireAndForget = true
		log.Debug().Msg("OutputRawTxOnly mode enabled - automatically enabling FireAndForget mode")
	}

	gas, err := r.client.SuggestGasPrice(ctx)
	if err != nil {
		return errors.New("unable to retrieve gas price: " + err.Error())
	}
	log.Trace().Interface("gasprice", gas).Msg("Retrieved current gas price")

	if !r.cfg.LegacyTxMode {
		var gasTipCap *big.Int
		gasTipCap, err = r.client.SuggestGasTipCap(ctx)
		if err != nil {
			return errors.New("unable to retrieve gas tip cap: " + err.Error())
		}
		log.Trace().Interface("gastipcap", gasTipCap).Msg("Retrieved current gas tip cap")
		r.cfg.CurrentGasTipCap = gasTipCap
	}

	trimmedHexPrivateKey := strings.TrimPrefix(r.cfg.PrivateKey, "0x")
	privateKey, err := crypto.HexToECDSA(trimmedHexPrivateKey)
	if err != nil {
		return errors.New("couldn't process the hex private key: " + err.Error())
	}

	blockNumber, err := r.client.BlockNumber(ctx)
	if err != nil {
		return errors.New("couldn't get the current block number: " + err.Error())
	}
	log.Trace().Uint64("blocknumber", blockNumber).Msg("Current Block Number")

	ethAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	bigBlockNumber := big.NewInt(int64(blockNumber))

	nonce, err := r.client.NonceAt(ctx, ethAddress, bigBlockNumber)
	if err != nil {
		return errors.New("unable to get account nonce: " + err.Error())
	}

	accountBal, err := r.client.BalanceAt(ctx, ethAddress, bigBlockNumber)
	if err != nil {
		return errors.New("unable to get the balance for the account: " + err.Error())
	}
	log.Trace().
		Str("addr", ethAddress.Hex()).
		Interface("balance", accountBal).
		Msg("funding account balance")

	toAddr := common.HexToAddress(r.cfg.ToAddress)
	amt := new(big.Int).SetUint64(r.cfg.EthAmountInWei)

	header, err := r.client.HeaderByNumber(ctx, nil)
	if err != nil {
		return errors.New("unable to get header: " + err.Error())
	}
	if header.BaseFee != nil {
		r.cfg.ChainSupportBaseFee = true
		log.Debug().Msg("Eip-1559 support detected")
	}

	chainID, err := r.client.ChainID(ctx)
	if err != nil {
		return errors.New("unable to fetch chain ID: " + err.Error())
	}
	log.Trace().Uint64("chainID", chainID.Uint64()).Msg("Detected Chain ID")

	r.cfg.BigGasPriceMultiplier = big.NewFloat(r.cfg.GasPriceMultiplier)

	if r.cfg.LegacyTxMode && r.cfg.ForcePriorityGasPrice > 0 {
		log.Warn().Msg("Cannot set priority gas price in legacy mode")
	}
	if r.cfg.ForceGasPrice < r.cfg.ForcePriorityGasPrice {
		return errors.New("max priority fee per gas higher than max fee per gas")
	}

	if r.cfg.AdaptiveRateLimit && r.cfg.EthCallOnly {
		return errors.New("the adaptive rate limit is based on the pending transaction pool. It doesn't use this feature while also using call only")
	}

	contractAddr := common.HexToAddress(r.cfg.ContractAddress)
	r.cfg.ContractETHAddress = &contractAddr

	r.cfg.ToETHAddress = &toAddr
	r.cfg.SendAmount = amt
	r.cfg.CurrentGasPrice = gas
	r.cfg.CurrentNonce = &nonce
	r.cfg.ECDSAPrivateKey = privateKey
	r.cfg.FromETHAddress = &ethAddress
	if r.cfg.ChainID == 0 {
		r.cfg.ChainID = chainID.Uint64()
	}

	// Parse modes before account pool init so RPC mode can force call-only behavior.
	if err := r.parseModes(ctx); err != nil {
		return err
	}

	// Initialize account pool
	if err := r.initAccountPool(ctx); err != nil {
		return err
	}

	return nil
}

func (r *Runner) initAccountPool(ctx context.Context) error {
	ecdsaPrivateKey := r.cfg.ECDSAPrivateKey

	apCfg := &AccountPoolConfig{
		FundingPrivateKey:         ecdsaPrivateKey,
		FundingAmount:             r.cfg.AccountFundingAmount,
		RateLimit:                 r.cfg.RateLimit,
		EthCallOnly:               r.cfg.EthCallOnly,
		RefundRemainingFunds:      r.cfg.RefundRemainingFunds,
		CheckBalanceBeforeFunding: r.cfg.CheckBalanceBeforeFunding,
		LegacyTxMode:              r.cfg.LegacyTxMode,
		ForceGasPrice:             r.cfg.ForceGasPrice,
		ForcePriorityGasPrice:     r.cfg.ForcePriorityGasPrice,
		GasPriceMultiplier:        r.cfg.BigGasPriceMultiplier,
		ChainSupportBaseFee:       r.cfg.ChainSupportBaseFee,
	}

	var err error
	r.accountPool, err = NewAccountPool(ctx, r.client, apCfg)
	if err != nil {
		return errors.New("unable to create account pool: " + err.Error())
	}

	// Add accounts based on configuration
	if r.cfg.SendingAccountsFile != "" {
		var privateKeys []*ecdsa.PrivateKey
		privateKeys, err = util.ReadPrivateKeysFromFile(r.cfg.SendingAccountsFile)
		if err != nil {
			return errors.New("unable to read private keys from file: " + err.Error())
		}
		if len(privateKeys) == 0 {
			return errors.New("no private keys found in sending accounts file")
		}
		if len(privateKeys) > 1 && r.cfg.StartNonce > 0 {
			log.Fatal().Msg("nonce can't be set while using multiple sending accounts")
		}
		if len(privateKeys) == 1 {
			var nonce *uint64
			if r.cfg.StartNonce > 0 {
				nonce = &r.cfg.StartNonce
			}
			err = r.accountPool.Add(ctx, privateKeys[0], nonce)
		} else {
			err = r.accountPool.AddN(ctx, privateKeys...)
		}
		r.cfg.SendingAccountsCount = uint64(len(privateKeys))
	} else if r.cfg.SendingAccountsCount > 0 {
		if r.cfg.StartNonce > 0 {
			log.Fatal().Msg("nonce can't be set while using random multiple sending accounts")
		}
		err = r.accountPool.AddRandomN(ctx, r.cfg.SendingAccountsCount)
	} else {
		var nonce *uint64
		if r.cfg.StartNonce > 0 {
			nonce = &r.cfg.StartNonce
		}
		err = r.accountPool.Add(ctx, ecdsaPrivateKey, nonce)
	}
	if err != nil {
		return errors.New("unable to set account pool: " + err.Error())
	}

	// Wait for all accounts to be ready
	for {
		rdy, rdyCount, accQty := r.accountPool.AllAccountsReady()
		if rdy {
			log.Info().Msg("All accounts are ready")
			break
		}
		log.Info().Int("ready", rdyCount).Int("total", accQty).Msg("waiting for all accounts to be ready")
		time.Sleep(time.Second)
	}

	// Pre-fund accounts if configured
	if r.cfg.SendingAccountsCount == 0 {
		log.Info().Msg("No sending accounts to pre-fund. Skipping pre-funding of sending accounts.")
		return nil
	}
	if r.cfg.EthCallOnly {
		log.Info().Msg("call only mode is enabled. Skipping pre-funding of sending accounts.")
		return nil
	}
	if !r.cfg.PreFundSendingAccounts {
		log.Info().Msg("pre-funding of sending accounts is disabled.")
		return nil
	}
	if r.cfg.AccountFundingAmount.Cmp(new(big.Int)) == 0 {
		log.Info().Msg("account funding amount is zero. Skipping pre-funding of sending accounts.")
		return nil
	}

	if err := r.accountPool.FundAccounts(ctx); err != nil {
		log.Error().Err(err).Msg("unable to fund sending accounts")
	}

	return nil
}

// Run executes the load test.
func (r *Runner) Run(ctx context.Context) error {
	log.Info().Msg("Starting Load Test")

	// Configure time limit
	var overallTimer *time.Timer
	if r.cfg.TimeLimit > 0 {
		overallTimer = time.NewTimer(time.Duration(r.cfg.TimeLimit) * time.Second)
		defer overallTimer.Stop()
	} else {
		overallTimer = new(time.Timer)
	}

	// Set up signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	defer signal.Stop(sigCh)

	errCh := make(chan error, 1)
	loadTestCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		select {
		case <-loadTestCtx.Done():
			return
		default:
			errCh <- r.mainLoop(loadTestCtx)
		}
	}()

	timedOut := false
	interrupted := false
	var mainLoopErr error

	// Wait for completion or interruption
	select {
	case <-overallTimer.C:
		log.Info().Msg("Time's up")
		timedOut = true
		cancel()
	case <-sigCh:
		log.Info().Msg("Interrupted.. Stopping load test")
		interrupted = true
		cancel()
		if r.cfg.ShouldProduceSummary {
			finalBlock, err := r.client.BlockNumber(ctx)
			if err != nil {
				log.Error().Err(err).Msg("Unable to retrieve final block number")
			} else {
				r.finalBlockNumber = finalBlock
			}
		}
	case err := <-errCh:
		mainLoopErr = err
	}

	if timedOut || interrupted {
		mainLoopErr = <-errCh
	}
	if mainLoopErr != nil {
		log.Fatal().Err(mainLoopErr).Msg("Received critical error while running load test")
	}

	if timedOut {
		log.Info().Msg("Finished")
		return nil
	}

	// Post-load-test operations use the original context (not the cancelled loadTestCtx)
	// to ensure summary/refund RPCs can complete successfully after SIGINT
	r.postLoadTest(ctx)

	log.Info().Msg("Finished")
	return nil
}

// postLoadTest handles post-load-test operations like summary and fund refunding.
func (r *Runner) postLoadTest(ctx context.Context) {
	cfg := r.cfg
	results := r.GetResults()

	// Always output a light summary if we have results
	if len(results) > 0 {
		startTime := results[0].RequestTime
		endTime := time.Now()
		LightSummary(results, startTime, endTime, r.rl)
	}

	// Skip detailed summary and refunds in fire-and-forget or call-only modes.
	// In these modes, transactions aren't tracked or no transactions are sent,
	// making detailed summaries misleading and refunds unnecessary.
	if cfg.FireAndForget || cfg.EthCallOnly {
		return
	}

	// Output detailed summary if requested
	if cfg.ShouldProduceSummary && r.startBlockNumber > 0 && r.finalBlockNumber > 0 {
		log.Info().Msg("Generating detailed summary...")
		if err := SummarizeResults(ctx, r.client, r.rpcClient, cfg, r.accountPool, results, r.startBlockNumber, r.finalBlockNumber); err != nil {
			log.Error().Err(err).Msg("Failed to generate detailed summary")
		}
	}

	// Refund remaining funds if requested
	if cfg.RefundRemainingFunds && r.accountPool != nil {
		log.Info().Msg("Refunding remaining funds...")
		if err := r.accountPool.ReturnFunds(ctx); err != nil {
			log.Error().Err(err).Msg("Failed to refund remaining funds")
		}
	}
}

func (r *Runner) mainLoop(ctx context.Context) error {
	cfg := r.cfg
	log.Trace().Interface("Input Params", cfg).Msg("Params")

	maxRoutines := cfg.Concurrency
	maxRequests := cfg.Requests
	chainID := new(big.Int).SetUint64(cfg.ChainID)
	privateKey := cfg.ECDSAPrivateKey

	r.rl = rate.NewLimiter(rate.Limit(cfg.RateLimit), 1)
	if cfg.RateLimit <= 0.0 {
		r.rl = nil
	}

	rateLimitCtx, rateLimitCancel := context.WithCancel(ctx)
	defer rateLimitCancel()
	if cfg.AdaptiveRateLimit && r.rl != nil {
		go r.updateRateLimit(rateLimitCtx)
	}

	tops, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return errors.New("unable to create transaction signer: " + err.Error())
	}
	tops = r.configureTransactOpts(ctx, tops)
	// Reset for contract deployment
	tops.GasLimit = 0
	tops.GasPrice = nil
	tops.GasFeeCap = nil
	tops.GasTipCap = nil

	// Deploy contracts if needed
	err = r.deployContracts(ctx, tops)
	if err != nil {
		return err
	}

	// Backfill ERC20/721 addresses for RPC mode if the chain has none
	if r.deps.IndexedActivity != nil {
		if len(r.deps.IndexedActivity.ERC20Addresses) == 0 && r.deps.ERC20Address != (common.Address{}) {
			r.deps.IndexedActivity.ERC20Addresses = append(r.deps.IndexedActivity.ERC20Addresses, r.deps.ERC20Address.String())
		}
		if len(r.deps.IndexedActivity.ERC721Addresses) == 0 && r.deps.ERC721Address != (common.Address{}) {
			r.deps.IndexedActivity.ERC721Addresses = append(r.deps.IndexedActivity.ERC721Addresses, r.deps.ERC721Address.String())
		}
	}

	r.startBlockNumber, err = r.client.BlockNumber(ctx)
	if err != nil {
		return errors.New("failed to get current block number: " + err.Error())
	}

	if cfg.StartNonce <= 0 {
		err = r.accountPool.RefreshNonce(ctx, tops.From)
		if err != nil {
			return err
		}
	}

	// Initialize modes
	err = r.initModes(ctx)
	if err != nil {
		return err
	}

	// Setup max base fee monitoring
	mustCheckMaxBaseFee, maxBaseFeeCtxCancel := r.setupBaseFeeMonitoring(ctx)
	defer maxBaseFeeCtxCancel()

	log.Debug().Msg("Starting main load test loop")
	var wg sync.WaitGroup
	for routineID := range maxRoutines {
		log.Trace().Int64("routineID", routineID).Msg("starting concurrent routine")
		wg.Add(1)
		go func(routineID int64) {
			defer wg.Done()
			var startReq, endReq time.Time
			var tErr error
			var ltTxHash common.Hash
			for requestID := range maxRequests {
				if ctx.Err() != nil {
					return
				}
				if r.rl != nil {
					if waitErr := r.rl.Wait(ctx); waitErr != nil {
						if errors.Is(waitErr, context.Canceled) || errors.Is(waitErr, context.DeadlineExceeded) {
							return
						}
						log.Error().Int64("routineID", routineID).Int64("requestID", requestID).Err(waitErr).Msg("Encountered a rate limiting error")
					}
				}

				if ctx.Err() != nil {
					return
				}

				// Select mode for this request
				selectedMode := r.selectMode(routineID, requestID)

				var account Account
				account, tErr = r.accountPool.Next(ctx)
				if tErr != nil {
					log.Error().Int64("routineID", routineID).Int64("requestID", requestID).Err(tErr).Msg("Unable to get next account from account pool")
					return
				}

				var sendingTops *bind.TransactOpts
				sendingTops, tErr = bind.NewKeyedTransactorWithChainID(account.PrivateKey(), chainID)
				if tErr != nil {
					log.Error().Int64("routineID", routineID).Int64("requestID", requestID).Err(tErr).Msg("Unable create transaction signer")
					return
				}
				sendingTops.Nonce = new(big.Int).SetUint64(account.Nonce())

				// Wait for base fee to drop if needed
				if mustCheckMaxBaseFee {
					waiting := false
					for r.waitBaseFeeToDrop.Load() {
						if ctx.Err() != nil {
							return
						}
						if !waiting {
							waiting = true
							log.Debug().Int64("routineID", routineID).Int64("requestID", requestID).Msg("go routine is waiting for base fee to drop")
						}
						time.Sleep(time.Second)
					}
				}

				sendingTops = r.configureTransactOpts(ctx, sendingTops)

				// Execute the selected mode
				startReq, endReq, ltTxHash, tErr = selectedMode.Execute(ctx, cfg, r.deps, sendingTops)

				// Record sample if not fire-and-forget
				if !cfg.FireAndForget {
					r.RecordSample(routineID, requestID, tErr, startReq, endReq, sendingTops.Nonce.Uint64())
				}

				// Wait for receipt if configured
				if tErr == nil && cfg.WaitForReceipt {
					_, tErr = util.WaitReceiptWithRetries(ctx, r.client, ltTxHash, cfg.ReceiptRetryMax, cfg.ReceiptRetryDelay)
				}

				// Handle errors
				if tErr != nil {
					log.Error().
						Int64("routineID", routineID).
						Int64("requestID", requestID).
						Err(tErr).
						Str("mode", selectedMode.Name()).
						Str("address", sendingTops.From.String()).
						Uint64("nonce", sendingTops.Nonce.Uint64()).
						Uint64("gas", sendingTops.GasLimit).
						Any("gasPrice", sendingTops.GasPrice).
						Any("gasFeeCap", sendingTops.GasFeeCap).
						Any("gasTipCap", sendingTops.GasTipCap).
						Int64("request time", endReq.Sub(startReq).Milliseconds()).
						Msg("recorded an error while sending transactions")

					// Check nonce for reuse
					if !cfg.EthCallOnly {
						r.handleNonceReuse(ctx, sendingTops, tErr)
					}
				}

				log.Trace().
					Int64("routineID", routineID).
					Int64("requestID", requestID).
					Stringer("txhash", ltTxHash).
					Any("nonce", sendingTops.Nonce).
					Str("mode", selectedMode.Name()).
					Str("sendingAddress", sendingTops.From.String()).
					Msg("Request")
			}
		}(routineID)
	}
	log.Trace().Msg("Finished starting go routines. Waiting..")
	wg.Wait()
	rateLimitCancel()

	if ctx.Err() != nil {
		return nil
	}

	// Wait for all transactions to be mined (unless fire-and-forget or call-only)
	if !cfg.FireAndForget && !cfg.EthCallOnly {
		log.Debug().Msg("Waiting for remaining transactions to be completed and mined")
		r.finalBlockNumber, err = r.waitForFinalBlock(ctx)
		if err != nil {
			log.Warn().Err(err).Msg("There was an issue waiting for all transactions to be mined")
		}
	} else {
		// Capture final block number for summary
		r.finalBlockNumber, err = r.client.BlockNumber(ctx)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to get final block number for summary")
		}
	}

	return nil
}

// parseModes converts mode strings to mode instances and populates cfg.ParsedModes.
func (r *Runner) parseModes(ctx context.Context) error {
	cfg := r.cfg

	if len(cfg.Modes) == 0 {
		return errors.New("expected at least one mode")
	}

	// Set multi-mode flag
	cfg.MultiMode = len(cfg.Modes) > 1

	// Parse mode strings to mode instances
	for _, modeName := range cfg.Modes {
		md, err := mode.Get(modeName)
		if err != nil {
			return err
		}
		r.modes = append(r.modes, md)

		// Parse to Mode enum for contract deployment decisions
		parsedMode, err := config.ParseMode(modeName)
		if err != nil {
			return err
		}
		cfg.ParsedModes = append(cfg.ParsedModes, parsedMode)
	}

	// Mode compatibility validation
	if cfg.MultiMode {
		if !config.HasUniqueModes(cfg.ParsedModes) {
			return errors.New("duplicate modes detected, check input modes for duplicates")
		}
		if config.HasMode(config.ModeRandom, cfg.ParsedModes) {
			return errors.New("random mode can't be used in combination with any other modes")
		}
		if config.HasMode(config.ModeBlob, cfg.ParsedModes) {
			return errors.New("blob mode should only be used by itself. Blob mode will take significantly longer than other transactions to finalize, and the address will be reserved, preventing other transactions from being made")
		}
		if config.HasMode(config.ModeRPC, cfg.ParsedModes) && !cfg.EthCallOnly {
			return errors.New("rpc mode must be called with eth-call-only when multiple modes are used")
		}
		if cfg.OutputRawTxOnly {
			return errors.New("raw output is not compatible with multiple modes")
		}
	} else if config.HasMode(config.ModeRPC, cfg.ParsedModes) {
		// RPC mode is read-only testing, automatically enable call-only mode
		// Validate incompatible options before auto-enabling
		if cfg.AdaptiveRateLimit {
			return errors.New("the adaptive rate limit is based on the pending transaction pool. It doesn't work with RPC mode which is read-only")
		}
		if cfg.WaitForReceipt {
			return errors.New("waiting for receipts doesn't make sense with RPC mode which is read-only")
		}
		if cfg.PreFundSendingAccounts || cfg.SendingAccountsFile != "" || cfg.SendingAccountsCount > 0 {
			return errors.New("pre-funding accounts doesn't make sense with RPC mode which is read-only")
		}
		log.Trace().Msg("Setting call only mode since we're doing RPC testing")
		cfg.EthCallOnly = true
	}
	if cfg.OutputRawTxOnly {
		if config.HasMode(config.ModeRPC, cfg.ParsedModes) {
			return errors.New("raw output is not compatible with RPC mode")
		}
		if config.HasMode(config.ModeUniswapV3, cfg.ParsedModes) {
			return errors.New("raw output is not compatible with UniswapV3 mode")
		}
	}
	if config.HasMode(config.ModeContractCall, cfg.ParsedModes) && (cfg.ContractAddress == "" || cfg.ContractCallData == "") {
		return errors.New("contract-call mode requires both --contract-address and --calldata flags")
	}
	if cfg.EthCallOnly && config.HasMode(config.ModeBlob, cfg.ParsedModes) {
		return errors.New("using call only with blobs doesn't make sense")
	}
	if cfg.LegacyTxMode && config.HasMode(config.ModeBlob, cfg.ParsedModes) {
		return errors.New("blob transactions require eip-1559")
	}
	if config.HasMode(config.ModeUniswapV3, cfg.ParsedModes) && cfg.UniswapV3 == nil {
		return errors.New("uniswapv3 mode requires the 'polycli loadtest uniswapv3' subcommand to configure contract addresses")
	}

	// Initialize mode-specific dependencies
	for _, parsedMode := range cfg.ParsedModes {
		switch parsedMode {
		case config.ModeRecall:
			if r.deps.RecallTransactions == nil {
				log.Info().Msg("Fetching recall transactions from recent blocks")
				txs, err := modes.GetRecallTransactions(ctx, r.client, r.rpcClient, cfg.RecallLength, cfg.BlockBatchSize)
				if err != nil {
					return errors.New("failed to fetch recall transactions: " + err.Error())
				}
				if len(txs) == 0 {
					return errors.New("we weren't able to fetch any recall transactions")
				}
				r.deps.RecallTransactions = txs
				log.Info().Int("count", len(txs)).Msg("Fetched recall transactions")
			}
		case config.ModeRPC:
			if r.deps.IndexedActivity == nil {
				log.Info().Msg("Fetching indexed activity from recent blocks")
				ia, err := modes.GetIndexedRecentActivity(ctx, r.client, r.rpcClient, cfg.RecallLength, cfg.BlockBatchSize)
				if err != nil {
					return errors.New("failed to fetch indexed activity: " + err.Error())
				}
				// Validate that the chain has enough activity for RPC mode
				if len(ia.TransactionIDs) == 0 || len(ia.Transactions) == 0 ||
					len(ia.Addresses) == 0 || len(ia.BlockIDs) == 0 ||
					ia.BlockNumber == 0 {
					return fmt.Errorf("insufficient chain activity for RPC mode: the chain must have at least some transaction history. Found %d transactions, %d addresses, %d blocks, current block number %d",
						len(ia.TransactionIDs), len(ia.Addresses), len(ia.BlockIDs), ia.BlockNumber)
				}
				r.deps.IndexedActivity = ia
				log.Info().Int("blockCount", len(ia.BlockNumbers)).Int("txCount", len(ia.TransactionIDs)).Msg("Fetched indexed activity")
			}
		}
	}

	return nil
}

func (r *Runner) initModes(ctx context.Context) error {
	for _, mode := range r.modes {
		if err := mode.Init(ctx, r.cfg, r.deps); err != nil {
			return errors.New("failed to init mode " + mode.Name() + ": " + err.Error())
		}
	}
	return nil
}

func (r *Runner) selectMode(routineID, requestID int64) mode.Runner {
	if len(r.modes) == 0 {
		return nil
	}

	// If multi-mode, cycle through modes
	if r.cfg.MultiMode && len(r.modes) > 1 {
		return r.modes[int(routineID+requestID)%len(r.modes)]
	}

	// Single mode
	return r.modes[0]
}

func (r *Runner) setupBaseFeeMonitoring(ctx context.Context) (bool, context.CancelFunc) {
	maxBaseFeeCtx, maxBaseFeeCtxCancel := context.WithCancel(ctx)
	mustCheckMaxBaseFee := r.cfg.MaxBaseFeeWei > 0

	r.waitBaseFeeToDrop.Store(false)
	if mustCheckMaxBaseFee {
		log.Info().Msg("max base fee monitoring enabled")

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			firstRun := true
			for {
				select {
				case <-maxBaseFeeCtx.Done():
					return
				default:
					currentBaseFeeIsGreaterThanMax, currentBaseFeeWei, err := r.isCurrentBaseFeeGreaterThanMax(ctx, r.cfg.MaxBaseFeeWei)
					if err != nil {
						log.Error().Err(err).Msg("Error checking base fee during load test")
					} else {
						if currentBaseFeeIsGreaterThanMax {
							if !r.waitBaseFeeToDrop.Load() {
								log.Warn().Msgf("PAUSE: base fee %d Wei > limit %d Wei", currentBaseFeeWei.Uint64(), r.cfg.MaxBaseFeeWei)
								r.waitBaseFeeToDrop.Store(true)
							}
						} else if r.waitBaseFeeToDrop.Load() {
							log.Info().Msgf("RESUME: base fee %d Wei ≤ limit %d Wei", currentBaseFeeWei.Uint64(), r.cfg.MaxBaseFeeWei)
							r.waitBaseFeeToDrop.Store(false)
						}

						if firstRun {
							firstRun = false
							wg.Done()
						}
					}
					time.Sleep(time.Second)
				}
			}
		}()

		wg.Wait()
	}

	return mustCheckMaxBaseFee, maxBaseFeeCtxCancel
}

func (r *Runner) isCurrentBaseFeeGreaterThanMax(ctx context.Context, maxBaseFee uint64) (bool, *big.Int, error) {
	header, err := r.client.HeaderByNumber(ctx, nil)
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
		}
		return false, currentBaseFee, nil
	}

	return false, nil, nil
}

func (r *Runner) handleNonceReuse(ctx context.Context, tops *bind.TransactOpts, tErr error) {
	// Start with assumption that we can reuse the nonce
	reuseNonce := true

	// If it is an error that consumes the nonce, we can't retry it
	if strings.Contains(tErr.Error(), "replacement transaction underpriced") ||
		strings.Contains(tErr.Error(), "transaction underpriced") ||
		strings.Contains(tErr.Error(), "nonce too low") ||
		strings.Contains(tErr.Error(), "already known") ||
		strings.Contains(tErr.Error(), "could not replace existing") {
		reuseNonce = false
	}

	// If we can reuse the nonce, add it back to the account pool
	if reuseNonce {
		err := r.accountPool.AddReusableNonce(ctx, tops.From, tops.Nonce.Uint64())
		if err != nil {
			log.Error().
				Str("address", tops.From.String()).
				Uint64("nonce", tops.Nonce.Uint64()).
				Err(err).
				Msg("Unable to add reusable nonce to account pool")
		}
	}
}

func (r *Runner) deployContracts(ctx context.Context, tops *bind.TransactOpts) error {
	cops := &bind.CallOpts{Context: ctx}

	// Deploy LoadTester contract if needed
	if r.cfg.LoadTestContractAddress == "" && config.AnyRequiresLoadTestContract(r.cfg.ParsedModes) {
		ltAddr, _, _, err := tester.DeployLoadTester(tops, r.client)
		if err != nil {
			return errors.New("failed to deploy load testing contract: " + err.Error())
		}
		r.deps.LoadTesterAddress = ltAddr
		ltContract, err := tester.NewLoadTester(ltAddr, r.client)
		if err != nil {
			return errors.New("unable to instantiate load tester contract: " + err.Error())
		}
		r.deps.LoadTesterContract = ltContract
		log.Debug().Stringer("ltAddr", ltAddr).Msg("Deployed load test contract")

		// Wait for contract to be mined and validated
		err = util.BlockUntilSuccessful(ctx, r.client, func() error {
			_, innerErr := ltContract.GetCallCounter(cops)
			return innerErr
		})
		if err != nil {
			return fmt.Errorf("failed waiting for load test contract to be mined: %w", err)
		}
	} else if r.cfg.LoadTestContractAddress != "" {
		ltAddr := common.HexToAddress(r.cfg.LoadTestContractAddress)
		r.deps.LoadTesterAddress = ltAddr
		ltContract, err := tester.NewLoadTester(ltAddr, r.client)
		if err != nil {
			return errors.New("unable to instantiate load tester contract: " + err.Error())
		}
		r.deps.LoadTesterContract = ltContract
	}

	// Deploy ERC20 contract if needed
	if r.cfg.ERC20Address == "" && config.AnyRequiresERC20(r.cfg.ParsedModes) {
		log.Info().Msg("Deploying ERC20 contract")
		erc20Addr, _, _, err := tokens.DeployERC20(tops, r.client)
		if err != nil {
			return errors.New("unable to deploy ERC20 contract: " + err.Error())
		}
		r.deps.ERC20Address = erc20Addr
		erc20Contract, err := tokens.NewERC20(erc20Addr, r.client)
		if err != nil {
			return errors.New("unable to instantiate ERC20 contract: " + err.Error())
		}
		r.deps.ERC20Contract = erc20Contract
		log.Info().Stringer("erc20Addr", erc20Addr).Msg("Deployed ERC20 contract")

		// Wait for contract to be mined and validate balance
		err = util.BlockUntilSuccessful(ctx, r.client, func() error {
			balance, innerErr := erc20Contract.BalanceOf(cops, *r.cfg.FromETHAddress)
			if innerErr != nil {
				return innerErr
			}
			if balance.Cmp(new(big.Int)) == 0 {
				return errors.New("ERC20 balance is zero")
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("failed waiting for ERC20 contract to be mined: %w", err)
		}
	} else if r.cfg.ERC20Address != "" {
		erc20Addr := common.HexToAddress(r.cfg.ERC20Address)
		r.deps.ERC20Address = erc20Addr
		erc20Contract, err := tokens.NewERC20(erc20Addr, r.client)
		if err != nil {
			return errors.New("unable to instantiate ERC20 contract: " + err.Error())
		}
		r.deps.ERC20Contract = erc20Contract
	}

	// Deploy ERC721 contract if needed
	if r.cfg.ERC721Address == "" && config.AnyRequiresERC721(r.cfg.ParsedModes) {
		log.Info().Msg("Deploying ERC721 contract")
		erc721Addr, _, _, err := tokens.DeployERC721(tops, r.client)
		if err != nil {
			return errors.New("unable to deploy ERC721 contract: " + err.Error())
		}
		r.deps.ERC721Address = erc721Addr
		erc721Contract, err := tokens.NewERC721(erc721Addr, r.client)
		if err != nil {
			return errors.New("unable to instantiate ERC721 contract: " + err.Error())
		}
		r.deps.ERC721Contract = erc721Contract
		log.Info().Stringer("erc721Addr", erc721Addr).Msg("Deployed ERC721 contract")

		// Wait for contract to be mined and validated
		err = util.BlockUntilSuccessful(ctx, r.client, func() error {
			_, innerErr := erc721Contract.BalanceOf(cops, *r.cfg.FromETHAddress)
			return innerErr
		})
		if err != nil {
			return fmt.Errorf("failed waiting for ERC721 contract to be mined: %w", err)
		}
	} else if r.cfg.ERC721Address != "" {
		erc721Addr := common.HexToAddress(r.cfg.ERC721Address)
		r.deps.ERC721Address = erc721Addr
		erc721Contract, err := tokens.NewERC721(erc721Addr, r.client)
		if err != nil {
			return errors.New("unable to instantiate ERC721 contract: " + err.Error())
		}
		r.deps.ERC721Contract = erc721Contract
	}

	// Initialize UniswapV3 if needed
	if config.HasMode(config.ModeUniswapV3, r.cfg.ParsedModes) && r.cfg.UniswapV3 != nil {
		log.Info().Msg("Initializing UniswapV3 contracts...")
		uniswapAddresses := uniswapv3.UniswapV3Addresses{
			FactoryV3:                          common.HexToAddress(r.cfg.UniswapV3.FactoryV3),
			Multicall:                          common.HexToAddress(r.cfg.UniswapV3.Multicall),
			ProxyAdmin:                         common.HexToAddress(r.cfg.UniswapV3.ProxyAdmin),
			TickLens:                           common.HexToAddress(r.cfg.UniswapV3.TickLens),
			NFTDescriptorLib:                   common.HexToAddress(r.cfg.UniswapV3.NFTDescriptorLib),
			NonfungibleTokenPositionDescriptor: common.HexToAddress(r.cfg.UniswapV3.NonfungibleTokenPositionDescriptor),
			TransparentUpgradeableProxy:        common.HexToAddress(r.cfg.UniswapV3.TransparentUpgradeableProxy),
			NonfungiblePositionManager:         common.HexToAddress(r.cfg.UniswapV3.NonfungiblePositionManager),
			Migrator:                           common.HexToAddress(r.cfg.UniswapV3.Migrator),
			Staker:                             common.HexToAddress(r.cfg.UniswapV3.Staker),
			QuoterV2:                           common.HexToAddress(r.cfg.UniswapV3.QuoterV2),
			SwapRouter02:                       common.HexToAddress(r.cfg.UniswapV3.SwapRouter),
			WETH9:                              common.HexToAddress(r.cfg.UniswapV3.WETH9),
		}

		cops := &bind.CallOpts{Context: ctx}
		initParams := uniswapv3.InitParams{
			PoolToken0Address: common.HexToAddress(r.cfg.UniswapV3.PoolToken0),
			PoolToken1Address: common.HexToAddress(r.cfg.UniswapV3.PoolToken1),
			PoolFees:          uniswapv3.PercentageToUniswapFeeTier(r.cfg.UniswapV3.PoolFees),
		}

		uniswapV3Config, poolConfig, err := uniswapv3.Init(ctx, r.client, tops, cops, uniswapAddresses, *r.cfg.FromETHAddress, initParams)
		if err != nil {
			return errors.New("failed to initialize UniswapV3: " + err.Error())
		}
		r.deps.UniswapV3Config = &uniswapV3Config
		r.deps.UniswapV3PoolConfig = &poolConfig
		log.Info().Msg("UniswapV3 initialized successfully")
	}

	return nil
}

func (r *Runner) updateRateLimit(ctx context.Context) {
	cfg := r.cfg
	ticker := time.NewTicker(time.Duration(cfg.AdaptiveCycleDuration) * time.Second)
	defer ticker.Stop()

	tryTxPool := true
	for {
		select {
		case <-ticker.C:
			var txPoolSize uint64
			var err error
			var pendingTxs uint64
			var queuedTxs uint64

			if tryTxPool {
				pendingTxs, queuedTxs, err = util.GetTxPoolStatus(r.rpcClient)
			}

			if err != nil {
				tryTxPool = false
				log.Warn().Err(err).Msg("Error getting txpool size. Falling back to latest nonce and disabling txpool check")

				pendingTxs, err = r.accountPool.NumberOfPendingTxs(ctx)
				if err != nil {
					log.Error().Err(err).Msg("Unable to get pending transactions to update rate limit")
					continue
				}
				txPoolSize = pendingTxs
			} else {
				txPoolSize = pendingTxs + queuedTxs
			}

			if txPoolSize < cfg.AdaptiveTargetSize {
				newRateLimit := rate.Limit(float64(r.rl.Limit()) + float64(cfg.AdaptiveRateLimitIncrement))
				r.rl.SetLimit(newRateLimit)
				log.Info().Float64("New Rate Limit (RPS)", float64(r.rl.Limit())).Uint64("Current Tx Pool Size", txPoolSize).Uint64("Steady State Tx Pool Size", cfg.AdaptiveTargetSize).Msg("Increased rate limit")
			} else if txPoolSize > cfg.AdaptiveTargetSize {
				r.rl.SetLimit(r.rl.Limit() / rate.Limit(cfg.AdaptiveBackoffFactor))
				log.Info().Float64("New Rate Limit (RPS)", float64(r.rl.Limit())).Uint64("Current Tx Pool Size", txPoolSize).Uint64("Steady State Tx Pool Size", cfg.AdaptiveTargetSize).Msg("Backed off rate limit")
			}
		case <-ctx.Done():
			return
		}
	}
}

func (r *Runner) configureTransactOpts(ctx context.Context, tops *bind.TransactOpts) *bind.TransactOpts {
	gasPrice, gasTipCap := r.getSuggestedGasPrices(ctx)
	tops.GasPrice = gasPrice

	cfg := r.cfg

	if cfg.ForceGasPrice != 0 {
		tops.GasPrice = big.NewInt(0).SetUint64(cfg.ForceGasPrice)
	}
	if cfg.ForceGasLimit != 0 {
		tops.GasLimit = cfg.ForceGasLimit
	}

	if cfg.LegacyTxMode {
		return tops
	}
	if !cfg.ChainSupportBaseFee {
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

func (r *Runner) getLatestBlockNumber(ctx context.Context) uint64 {
	r.cachedLatestBlockLock.Lock()
	defer r.cachedLatestBlockLock.Unlock()
	if time.Since(r.cachedLatestBlockTime) < 1*time.Second {
		return r.cachedLatestBlockNumber
	}
	bn, err := r.client.BlockNumber(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Unable to get block number while checking gas prices")
		return 0
	}
	r.cachedLatestBlockTime = time.Now()
	r.cachedLatestBlockNumber = bn
	return bn
}

func (r *Runner) biasGasPrice(price *big.Int) *big.Int {
	gasPriceFloat := new(big.Float).SetInt(price)
	gasPriceFloat.Mul(gasPriceFloat, r.cfg.BigGasPriceMultiplier)
	result := new(big.Int)
	gasPriceFloat.Int(result)
	return result
}

func (r *Runner) getSuggestedGasPrices(ctx context.Context) (*big.Int, *big.Int) {
	r.cachedGasPriceLock.Lock()
	defer r.cachedGasPriceLock.Unlock()

	bn := r.getLatestBlockNumber(ctx)
	if r.cachedBlockNumber != nil && bn <= *r.cachedBlockNumber {
		return r.cachedGasPrice, r.cachedGasTipCap
	}

	var gasPrice, gasTipCap = big.NewInt(0), big.NewInt(0)
	var pErr, tErr error
	cfg := r.cfg

	if cfg.LegacyTxMode {
		if cfg.ForceGasPrice != 0 {
			gasPrice = new(big.Int).SetUint64(cfg.ForceGasPrice)
		} else {
			gasPrice, pErr = r.client.SuggestGasPrice(ctx)
			if pErr == nil {
				gasPrice = r.biasGasPrice(gasPrice)
			} else {
				log.Error().Err(pErr).Msg("Unable to suggest gas price")
				return r.cachedGasPrice, r.cachedGasTipCap
			}
		}
	} else {
		var forcePriorityGasPrice *big.Int
		if cfg.ForcePriorityGasPrice != 0 {
			gasTipCap = new(big.Int).SetUint64(cfg.ForcePriorityGasPrice)
			forcePriorityGasPrice = gasTipCap
		} else if cfg.ChainSupportBaseFee {
			gasTipCap, tErr = r.client.SuggestGasTipCap(ctx)
			if tErr == nil {
				gasTipCap = r.biasGasPrice(gasTipCap)
			} else {
				log.Error().Err(tErr).Msg("Unable to suggest gas tip cap")
				return r.cachedGasPrice, r.cachedGasTipCap
			}
		} else {
			log.Fatal().Msg("Chain does not support base fee. Please set priority-gas-price flag with a value to use for gas tip cap")
		}

		if cfg.ForceGasPrice != 0 {
			gasPrice = new(big.Int).SetUint64(cfg.ForceGasPrice)
		} else if cfg.ChainSupportBaseFee {
			gasPrice = r.suggestMaxFeePerGas(ctx, bn, forcePriorityGasPrice)
		} else {
			log.Fatal().Msg("Chain does not support base fee. Please set gas-price flag with a value to use for max fee per gas")
		}
	}

	r.cachedBlockNumber = &bn
	r.cachedGasPrice = gasPrice
	r.cachedGasTipCap = gasTipCap

	log.Debug().
		Uint64("cachedBlockNumber", bn).
		Interface("cachedGasPrice", r.cachedGasPrice).
		Interface("cachedGasTipCap", r.cachedGasTipCap).
		Msg("Updating gas prices")

	return r.cachedGasPrice, r.cachedGasTipCap
}

func (r *Runner) suggestMaxFeePerGas(ctx context.Context, blockNumber uint64, forcePriorityFee *big.Int) *big.Int {
	header, err := r.client.HeaderByNumber(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg("Unable to get latest block header while checking MaxFeePerGas")
		return nil
	}

	if r.cachedBlockNumber != nil && blockNumber <= *r.cachedBlockNumber && r.cachedGasPrice != nil {
		return r.cachedGasPrice
	}

	feeHistory, err := r.client.FeeHistory(ctx, 5, nil, []float64{0.5})
	if err != nil {
		log.Error().Err(err).Msg("Unable to get fee history while checking MaxFeePerGas")
		return nil
	}

	priorityFee := forcePriorityFee
	if priorityFee == nil {
		priorityFee = feeHistory.Reward[len(feeHistory.Reward)-1][0]
	}
	baseFee := feeHistory.BaseFee[len(feeHistory.BaseFee)-1]
	maxFeePerGas := new(big.Int)
	maxFeePerGas.Mul(baseFee, big.NewInt(2))
	maxFeePerGas.Add(maxFeePerGas, priorityFee)

	const blocksToWait = 5
	isDecreasing := r.cachedGasPrice != nil && maxFeePerGas.Uint64() <= r.cachedGasPrice.Uint64()
	canDecrease := blockNumber+blocksToWait <= header.Number.Uint64()
	if isDecreasing && !canDecrease && r.cachedGasPrice != nil {
		return r.cachedGasPrice
	}

	r.cachedGasPrice = maxFeePerGas

	log.Trace().
		Uint64("blockNumber", header.Number.Uint64()).
		Str("priorityFee", priorityFee.String()).
		Str("baseFee", baseFee.String()).
		Str("maxFeePerGas", maxFeePerGas.String()).
		Msg("max fee updated")

	return maxFeePerGas
}

// RecordSample records a load test sample.
func (r *Runner) RecordSample(goRoutineID, requestID int64, err error, start, end time.Time, nonce uint64) {
	s := Sample{}
	s.GoRoutineID = goRoutineID
	s.RequestID = requestID
	s.RequestTime = start
	s.WaitTime = end.Sub(start)
	s.Nonce = nonce
	if err != nil {
		s.IsError = true
	}
	r.resultsMu.Lock()
	r.results = append(r.results, s)
	r.resultsMu.Unlock()
}

// GetResults returns all recorded samples.
func (r *Runner) GetResults() []Sample {
	r.resultsMu.RLock()
	defer r.resultsMu.RUnlock()
	result := make([]Sample, len(r.results))
	copy(result, r.results)
	return result
}

// GetAccountPool returns the account pool.
func (r *Runner) GetAccountPool() *AccountPool {
	return r.accountPool
}

// Close cleans up runner resources.
func (r *Runner) Close() {
	if r.rpcClient != nil {
		r.rpcClient.Close()
	}
}

// waitForFinalBlock waits for all transactions to be mined by checking nonces.
func (r *Runner) waitForFinalBlock(ctx context.Context) (uint64, error) {
	cfg := r.cfg
	var lastBlockNumber uint64
	var err error
	const checkInterval = 5 * time.Second
	const maxRetries = 30

	rateLimiter := rate.NewLimiter(rate.Limit(cfg.RateLimit), 1)
	noncesToCheck := r.accountPool.Nonces(ctx, true)

	for retry := 1; retry <= maxRetries; retry++ {
		lastBlockNumber, err = r.client.BlockNumber(ctx)
		if err != nil {
			return 0, err
		}
		if cfg.EthCallOnly {
			return lastBlockNumber, nil
		}

		var wg sync.WaitGroup
		var remainingNoncesToCheck atomic.Int64
		noncesToCheck.Range(func(key, value any) bool {
			wg.Add(1)
			remainingNoncesToCheck.Add(1)
			address := key.(common.Address)
			expectedNonce := value.(uint64)
			go func(ctx context.Context, rl *rate.Limiter) {
				defer wg.Done()
				if waitErr := rl.Wait(ctx); waitErr != nil {
					log.Error().Err(waitErr).Msg("Rate limiter wait error")
					return
				}
				nonce, nonceErr := r.client.NonceAt(ctx, address, new(big.Int).SetUint64(lastBlockNumber))
				if nonceErr != nil {
					log.Error().Err(nonceErr).Str("address", address.String()).Msg("Unable to get nonce for account while checking for final block")
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
			log.Debug().Uint64("lastBlockNumber", lastBlockNumber).Msg("All transactions mined")
			return lastBlockNumber, nil
		}

		log.Debug().
			Int64("remainingNonces", remainingNoncesToCheck.Load()).
			Int("retry", retry).
			Int("maxRetries", maxRetries).
			Msg("Waiting for transactions to be mined...")
		time.Sleep(checkInterval)
	}

	log.Error().Msg("Max retries reached waiting for transactions to be mined")
	return lastBlockNumber, fmt.Errorf("max retries reached waiting for transactions to be mined")
}

// SetModes sets the modes to be used during load testing.
func (r *Runner) SetModes(m []mode.Runner) {
	r.modes = m
}

// GetDependencies returns the mode dependencies.
func (r *Runner) GetDependencies() *mode.Dependencies {
	return r.deps
}

// GetClient returns the ethclient.
func (r *Runner) GetClient() *ethclient.Client {
	return r.client
}

// GetConfig returns the configuration.
func (r *Runner) GetConfig() *config.Config {
	return r.cfg
}

// Run is a convenience function that creates a runner, initializes it, and runs the load test.
// This allows both the main loadtest command and subcommands to use the same entry point.
func Run(ctx context.Context, cfg *config.Config) error {
	runner, err := NewRunner(cfg)
	if err != nil {
		return err
	}
	defer runner.Close()

	if err := runner.Init(ctx); err != nil {
		return err
	}

	return runner.Run(ctx)
}
