package loadtest

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/0xPolygon/polygon-cli/util"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

// AccountPoolConfig holds configuration for the account pool.
type AccountPoolConfig struct {
	FundingPrivateKey         *ecdsa.PrivateKey
	FundingAmount             *big.Int
	RateLimit                 float64
	EthCallOnly               bool
	RefundRemainingFunds      bool
	CheckBalanceBeforeFunding bool
	LegacyTxMode              bool
	// Gas override settings
	ForceGasPrice         uint64
	ForcePriorityGasPrice uint64
	GasPriceMultiplier    *big.Float
	ChainSupportBaseFee   bool
}

// Account represents a single account used by the load test.
type Account struct {
	ready          bool
	address        common.Address
	privateKey     *ecdsa.PrivateKey
	startNonce     uint64
	nonce          uint64
	funded         bool
	reusableNonces []uint64
}

// newAccount creates a new account with the given private key.
// The client is used to get the nonce of the account.
func newAccount(ctx context.Context, client *ethclient.Client, clientRateLimiter *rate.Limiter, privateKey *ecdsa.PrivateKey, startNonce *uint64, mu *sync.Mutex) (*Account, error) {
	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	acc := &Account{
		ready:          false,
		privateKey:     privateKey,
		address:        address,
		funded:         false,
		reusableNonces: make([]uint64, 0),
	}

	if startNonce != nil {
		acc.nonce = *startNonce
		acc.startNonce = *startNonce
		acc.ready = true
	} else {
		go func(a *Account) {
			for {
				log.Trace().Stringer("addr", acc.address).Msg("loading nonce for account in background, account not ready to be used yet")
				var err error
				err = clientRateLimiter.Wait(ctx)
				if err != nil {
					log.Error().Err(err).Stringer("addr", a.address).Msg("rate limiter wait failed getting nonce for account, retrying...")
					time.Sleep(time.Second)
					continue
				}
				acc.nonce, err = client.NonceAt(ctx, address, nil)
				if err != nil {
					log.Error().Err(err).Stringer("addr", a.address).Msg("failed to get nonce for account, retrying...")
					time.Sleep(time.Second)
					continue
				}
				acc.startNonce = acc.nonce
				mu.Lock()
				defer mu.Unlock()
				acc.ready = true
				break
			}
		}(acc)
	}

	return acc, nil
}

// Address returns the address of the account.
func (a *Account) Address() common.Address {
	return a.address
}

// PrivateKey returns the private key of the account.
func (a *Account) PrivateKey() *ecdsa.PrivateKey {
	return a.privateKey
}

// Nonce returns the current nonce of the account.
func (a *Account) Nonce() uint64 {
	return a.nonce
}

// AccountPool manages a pool of accounts used for sending transactions.
type AccountPool struct {
	accounts          []*Account
	accountsPositions map[common.Address]int

	client *ethclient.Client
	// clientRateLimiter is used to limit the rate of requests the account
	// pool needs to make to the network, like getting the nonce or balance.
	// it doesn't affect the requests made by the load test and is used only
	// internally to the account pool.
	clientRateLimiter *rate.Limiter

	mu                  sync.Mutex
	currentAccountIndex int
	fundingPrivateKey   *ecdsa.PrivateKey
	fundingAmount       *big.Int
	chainID             *big.Int

	latestBlockNumber uint64
	pendingTxsCache   *uint64

	// Configuration passed during creation
	cfg *AccountPoolConfig
}

// NewAccountPool creates a new account pool with the given configuration.
func NewAccountPool(ctx context.Context, client *ethclient.Client, cfg *AccountPoolConfig) (*AccountPool, error) {
	if cfg.FundingPrivateKey == nil {
		log.Fatal().
			Msg("fundingPrivateKey cannot be nil")
	}

	if cfg.FundingAmount == nil {
		log.Fatal().
			Msg("fundingAmount cannot be nil")
	}

	// Allow fundingAmount to be set to 0. Only check for negative fundingAmount.
	if cfg.FundingAmount.Cmp(big.NewInt(0)) < 0 {
		log.Fatal().
			Stringer("fundingAmount", cfg.FundingAmount).
			Msg("fundingAmount must be greater or equal to zero")
	}

	if client == nil {
		log.Fatal().
			Msg("client cannot be nil")
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Error().
			Err(err).
			Msg("unable to get chain id")
		return nil, fmt.Errorf("unable to get chain id: %w", err)
	}

	latestBlockNumber, err := client.BlockNumber(ctx)
	if err != nil {
		log.Error().
			Err(err).
			Msg("unable to get latestBlockNumber")
		return nil, fmt.Errorf("unable to get latestBlockNumber: %w", err)
	}

	ap := &AccountPool{
		currentAccountIndex: 0,
		client:              client,
		accounts:            make([]*Account, 0),
		fundingPrivateKey:   cfg.FundingPrivateKey,
		fundingAmount:       cfg.FundingAmount,
		chainID:             chainID,
		accountsPositions:   make(map[common.Address]int),
		latestBlockNumber:   latestBlockNumber,
		clientRateLimiter:   rate.NewLimiter(rate.Every(50*time.Millisecond), 1),
		cfg:                 cfg,
	}

	if !ap.isFundingEnabled() {
		if ap.isCallOnly() {
			log.Debug().
				Msg("sending account funding is disabled in call only mode")
		}

		if !ap.hasFundingAmount() {
			log.Debug().
				Msg("sending account funding is disabled due to funding amount being zero")
		}
	}

	return ap, nil
}

// AllAccountsReady returns whether all accounts are ready for use.
func (ap *AccountPool) AllAccountsReady() (bool, int, int) {
	ap.mu.Lock()
	defer ap.mu.Unlock()
	rdyCount := 0
	for i := range ap.accounts {
		if ap.accounts[i].ready {
			rdyCount++
		}
	}
	return rdyCount == len(ap.accounts), rdyCount, len(ap.accounts)
}

// AddRandomN adds N random accounts to the pool.
func (ap *AccountPool) AddRandomN(ctx context.Context, n uint64) error {
	for range n {
		err := ap.AddRandom(ctx)
		if err != nil {
			return fmt.Errorf("failed to add random account: %w", err)
		}
	}
	return nil
}

// AddRandom adds a random account to the pool.
func (ap *AccountPool) AddRandom(ctx context.Context) error {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return fmt.Errorf("failed to generate private key: %w", err)
	}

	forceNonce := uint64(0)
	return ap.Add(ctx, privateKey, &forceNonce)
}

// AddN adds multiple accounts to the pool with the given private keys.
func (ap *AccountPool) AddN(ctx context.Context, privateKeys ...*ecdsa.PrivateKey) error {
	for _, privateKey := range privateKeys {
		err := ap.Add(ctx, privateKey, nil)
		if err != nil {
			return fmt.Errorf("failed to add account: %w", err)
		}
	}

	return nil
}

// Add adds an account to the pool with the given private key.
func (ap *AccountPool) Add(ctx context.Context, privateKey *ecdsa.PrivateKey, startNonce *uint64) error {
	ap.mu.Lock()
	defer ap.mu.Unlock()

	account, err := newAccount(ctx, ap.client, ap.clientRateLimiter, privateKey, startNonce, &ap.mu)
	if err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}

	addressHex, privateKeyHex := util.GetAddressAndPrivateKeyHex(ctx, privateKey)
	log.Debug().
		Str("address", addressHex).
		Str("privateKey", privateKeyHex).
		Uint64("nonce", account.nonce).
		Msg("adding account to pool")

	ap.accounts = append(ap.accounts, account)
	ap.accountsPositions[account.address] = len(ap.accounts) - 1
	return nil
}

// AddReusableNonce adds a reusable nonce to the account with the given address.
func (ap *AccountPool) AddReusableNonce(ctx context.Context, address common.Address, nonce uint64) error {
	ap.mu.Lock()
	defer ap.mu.Unlock()

	accountPos, found := ap.accountsPositions[address]
	if !found {
		return fmt.Errorf("account not found in pool")
	}
	if accountPos > len(ap.accounts)-1 {
		return fmt.Errorf("account position out of bounds")
	}

	ap.accounts[accountPos].reusableNonces = append(ap.accounts[accountPos].reusableNonces, nonce)

	// sort the reusable nonces ascending because we want to use the lowest nonce first
	// and we pay the price of sorting only once when adding it
	slices.Sort(ap.accounts[accountPos].reusableNonces)

	log.Debug().
		Stringer("address", address).
		Uint64("nonce", nonce).
		Any("reusableNonces", ap.accounts[accountPos].reusableNonces).
		Msg("reusable nonce added to account")

	return nil
}

// RefreshNonce refreshes the nonce for the given address.
func (ap *AccountPool) RefreshNonce(ctx context.Context, address common.Address) error {
	ap.mu.Lock()
	defer ap.mu.Unlock()

	accountPos, found := ap.accountsPositions[address]
	if !found {
		return nil
	}
	if accountPos > len(ap.accounts)-1 {
		return fmt.Errorf("account position out of bounds")
	}

	err := ap.clientRateLimiter.Wait(ctx)
	if err != nil {
		return err
	}

	nonce, err := ap.client.NonceAt(ctx, address, nil)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %w", err)
	}

	ap.accounts[accountPos].nonce = nonce

	log.Debug().
		Stringer("address", address).
		Uint64("nonce", nonce).
		Msg("nonce refreshed")

	return nil
}

// NumberOfPendingTxs returns the difference between the internal nonce
// and the network pending nonce for all accounts in the pool.
func (ap *AccountPool) NumberOfPendingTxs(ctx context.Context) (uint64, error) {
	err := ap.clientRateLimiter.Wait(ctx)
	if err != nil {
		return 0, err
	}
	lbn, err := ap.client.BlockNumber(ctx)
	if err != nil {
		return 0, err
	}

	if lbn == ap.latestBlockNumber && ap.pendingTxsCache != nil {
		log.Debug().
			Uint64("pendingTxs", *ap.pendingTxsCache).
			Msg("returning cached pending transactions")
		return *ap.pendingTxsCache, nil
	}

	ap.mu.Lock()
	accCount := len(ap.accounts)
	accNonceMap := make(map[common.Address]uint64, accCount)
	for _, acc := range ap.accounts {
		accNonceMap[acc.address] = acc.nonce
	}
	ap.mu.Unlock()

	pendingTxCh := make(chan uint64, accCount)
	errCh := make(chan error, accCount)

	wg := sync.WaitGroup{}
	wg.Add(accCount)

	for addr, nonce := range accNonceMap {
		go func(a common.Address, n uint64) {
			defer wg.Done()
			err := ap.clientRateLimiter.Wait(ctx)
			if err != nil {
				errCh <- fmt.Errorf("failed to wait rate limit to get pending nonce for acc %s: %w", a.String(), err)
				return
			}
			pendingNonce, err := ap.client.PendingNonceAt(ctx, a)
			if err != nil {
				errCh <- fmt.Errorf("failed to get pending nonce for acc %s: %w", a.String(), err)
				return
			}
			pendingTxs := pendingNonce - n
			pendingTxCh <- pendingTxs
		}(addr, nonce)
	}
	wg.Wait()

	close(errCh)
	close(pendingTxCh)

	for err := range errCh {
		if err != nil {
			return 0, fmt.Errorf("failed to get pending transactions: %w", err)
		}
	}

	pendingTxs := uint64(0)
	for pendingTx := range pendingTxCh {
		pendingTxs += pendingTx
	}

	log.Debug().
		Uint64("pendingTxs", pendingTxs).
		Msg("number of pending transactions")

	ap.latestBlockNumber = lbn
	ap.pendingTxsCache = &pendingTxs

	return pendingTxs, nil
}

// FundAccounts funds all accounts in the pool.
func (ap *AccountPool) FundAccounts(ctx context.Context) error {
	ap.mu.Lock()
	defer ap.mu.Unlock()

	if !ap.isFundingEnabled() {
		log.Info().
			Uint64("fundingAmount", ap.fundingAmount.Uint64()).
			Msg("account funding is disabled, skipping funding of sending accounts")
		return nil
	}

	tops, err := bind.NewKeyedTransactorWithChainID(ap.fundingPrivateKey, ap.chainID)
	if err != nil {
		log.Error().Err(err).Msg("unable create transaction signer")
		return err
	}

	err = ap.clientRateLimiter.Wait(ctx)
	if err != nil {
		return err
	}
	balance, err := ap.client.BalanceAt(ctx, tops.From, nil)
	if err != nil {
		log.Error().Err(err).Msg("unable to get funding account balance")
	}

	accCount := len(ap.accounts)

	totalBalanceNeeded := new(big.Int).Mul(ap.fundingAmount, big.NewInt(int64(accCount)))
	totalFeeNeeded := new(big.Int).Mul(big.NewInt(21000), big.NewInt(int64(accCount)))
	fudgeAmountNeeded := new(big.Int).Mul(big.NewInt(1000000000), big.NewInt(int64(accCount)))

	totalNeeded := new(big.Int).Add(totalBalanceNeeded, totalFeeNeeded)
	totalNeeded.Add(totalNeeded, fudgeAmountNeeded)

	if balance.Cmp(totalBalanceNeeded) <= 0 {
		errMsg := "funding account balance can't cover the funding amount for all accounts"
		log.Error().
			Stringer("address", tops.From).
			Stringer("balance", balance).
			Stringer("totalNeeded", totalNeeded).
			Msg(errMsg)
		return errors.New(errMsg)
	}

	log.Debug().Msg("checking if multicall3 is supported")
	multicall3Addr, _ := util.IsMulticall3Supported(ctx, ap.client, true, tops, nil)
	if multicall3Addr != nil {
		log.Info().
			Stringer("address", multicall3Addr).
			Msg("multicall3 is supported and will be used to fund accounts")
	} else {
		log.Info().Msg("multicall3 is not supported, will use EOA transfers to fund accounts")
	}

	if multicall3Addr != nil {
		return ap.fundAccountsWithMulticall3(ctx, tops, multicall3Addr)
	}
	return ap.fundAccountsWithEOATransfers(ctx, tops)
}

func (ap *AccountPool) fundAccountsWithMulticall3(ctx context.Context, tops *bind.TransactOpts, multicall3Addr *common.Address) error {
	log.Debug().
		Msg("funding sending accounts with multicall3")

	const defaultAccsToFundPerTx = 400
	accsToFundPerTx, err := util.Multicall3MaxAccountsToFundPerTx(ctx, ap.client)
	if err != nil {
		log.Warn().Err(err).
			Uint64("defaultAccsToFundPerTx", defaultAccsToFundPerTx).
			Msg("failed to get multicall3 max accounts to fund per tx, falling back to default")
		accsToFundPerTx = defaultAccsToFundPerTx
	}
	log.Debug().Uint64("accsToFundPerTx", accsToFundPerTx).Msg("multicall3 max accounts to fund per tx")
	chSize := (uint64(len(ap.accounts)) / accsToFundPerTx) + 1

	txsCh := make(chan *types.Transaction, chSize)
	errCh := make(chan error, chSize)

	accs := []common.Address{}
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	for i := 0; i < len(ap.accounts); i++ {
		accountToFund := ap.accounts[i]
		// if account is the funding account, skip it
		if accountToFund.address == tops.From {
			continue
		}
		if mustBeFunded, iErr := ap.accountMustBeFunded(ctx, accountToFund); iErr != nil || !mustBeFunded {
			continue
		}

		accs = append(accs, accountToFund.address)

		if uint64(len(accs)) == accsToFundPerTx || i == len(ap.accounts)-1 {
			wg.Add(1)
			go func(tops *bind.TransactOpts, accs []common.Address) {
				defer wg.Done()
				iErr := ap.clientRateLimiter.Wait(ctx)
				if iErr != nil {
					log.Error().Err(iErr).Msg("rate limiter wait failed before funding accounts with multicall3")
					return
				}
				mu.Lock()
				defer mu.Unlock()
				tx, iErr := util.Multicall3FundAccountsWithNativeToken(ap.client, tops, accs, ap.fundingAmount, multicall3Addr)
				if iErr != nil {
					log.Error().Err(iErr).Msg("failed to fund accounts with multicall3")
					return
				}
				log.Info().
					Stringer("txHash", tx.Hash()).
					Int("done", i+1).
					Uint64("of", uint64(len(ap.accounts))).
					Msg("multicall3 transaction to fund accounts sent")
				txsCh <- tx
			}(tops, accs)
			accs = []common.Address{}
		}
	}
	wg.Wait()
	close(txsCh)
	close(errCh)

	var combinedErrors error
	for len(errCh) > 0 {
		err = <-errCh
		if combinedErrors == nil {
			combinedErrors = err
		} else {
			combinedErrors = errors.Join(combinedErrors, err)
		}
	}
	// return if there were errors sending the funding transactions
	if combinedErrors != nil {
		return combinedErrors
	}

	log.Info().Msg("all funding transactions sent, waiting for confirmation...")

	// ensure the txs to fund sending accounts using multicall3 were mined successfully
	for tx := range txsCh {
		err := ap.clientRateLimiter.Wait(ctx)
		if err != nil {
			return err
		}
		r, err := util.WaitReceipt(ctx, ap.client, tx.Hash())
		if err != nil {
			log.Error().Err(err).Msg("failed to wait for transaction to fund accounts with multicall3")
			return err
		}
		if r == nil || r.Status != types.ReceiptStatusSuccessful {
			errMsg := fmt.Sprintf("transaction to fund accounts with multicall3 failed, receipt is nil or status is not successful, txHash: %s", tx.Hash().String())
			log.Error().Msg(errMsg)
			return errors.New(errMsg)
		}
		log.Info().
			Stringer("txHash", tx.Hash()).
			Msg("transaction to fund accounts confirmed")
	}

	// mark all accounts as funded
	for i := range ap.accounts {
		ap.accounts[i].funded = true
	}

	return nil
}

func (ap *AccountPool) fundAccountsWithEOATransfers(ctx context.Context, tops *bind.TransactOpts) error {
	log.Debug().
		Msg("funding sending accounts with EOA transfers")

	accCount := len(ap.accounts)

	wg := sync.WaitGroup{}
	wg.Add(accCount)

	txCh := make(chan *types.Transaction, accCount)
	errCh := make(chan error, accCount)

	err := ap.clientRateLimiter.Wait(ctx)
	if err != nil {
		return err
	}
	nonce, err := ap.client.PendingNonceAt(ctx, tops.From)
	if err != nil {
		log.Error().Err(err).Msg("unable to get nonce")
	}

	for i := range ap.accounts {
		accountToFund := ap.accounts[i]

		// if account is the funding account, skip it
		if accountToFund.address == tops.From {
			continue
		}

		go func(forcedNonce uint64, account *Account) {
			defer wg.Done()
			if !account.funded {
				tx, err := ap.fundAccountIfNeeded(ctx, account, &forcedNonce, false)
				if err != nil {
					errCh <- fmt.Errorf("failed to fund account: %w", err)
				}
				if tx != nil {
					txCh <- tx
				}
			}
		}(nonce, accountToFund)
		nonce++
	}

	wg.Wait()

	close(errCh)
	close(txCh)

	for err := range errCh {
		if err != nil {
			return err
		}
	}

	failed := atomic.Bool{}
	for tx := range txCh {
		if tx != nil {
			log.Debug().
				Stringer("address", tx.To()).
				Stringer("txHash", tx.Hash()).
				Msg("transaction to fund account sent")
			wg.Add(1)

			go func(ctx context.Context, tx *types.Transaction, rl *rate.Limiter) {
				defer wg.Done()
				err := rl.Wait(ctx)
				if err != nil {
					failed.Store(true)
					log.Error().
						Err(err).
						Stringer("address", tx.To()).
						Stringer("txHash", tx.Hash()).
						Msg("failed to wait rate limiter before waiting for receipt of transaction to fund account")
					return
				}
				receipt, err := util.WaitReceipt(ctx, ap.client, tx.Hash())
				if err != nil {
					log.Error().
						Err(err).
						Stringer("address", tx.To()).
						Stringer("txHash", tx.Hash()).
						Msg("failed to wait for transaction to fund account")
					failed.Store(true)
					return
				}
				if receipt == nil {
					log.Error().
						Stringer("address", tx.To()).
						Stringer("txHash", tx.Hash()).
						Msg("transaction to fund account receipt is nil")
					failed.Store(true)
					return
				}
				if receipt.Status != types.ReceiptStatusSuccessful {
					log.Error().
						Stringer("address", tx.To()).
						Stringer("txHash", tx.Hash()).
						Msg("transaction to fund account has failed")
					failed.Store(true)
					return
				}
				log.Debug().
					Stringer("address", tx.To()).
					Stringer("txHash", tx.Hash()).
					Msg("transaction to fund account confirmed")
			}(ctx, tx, ap.clientRateLimiter)
		}
	}
	log.Info().Msg("all funding transactions sent, waiting for confirmation...")
	wg.Wait()
	log.Info().Msg("all funding transactions confirmed, verifying...")

	log.Debug().Msg("funding process finished")

	if failed.Load() {
		err := ap.returnFunds(ctx, false)
		if err != nil {
			log.Error().
				Err(err).
				Msg("failed to return funds from accounts after funding failure")
			return fmt.Errorf("failed to return funds from accounts after funding failure: %w", err)
		}
		return fmt.Errorf("some transactions to fund accounts failed")
	}

	log.Debug().
		Msg("all accounts funded")

	return nil
}

// ReturnFunds returns funds from all accounts back to the funding account.
func (ap *AccountPool) ReturnFunds(ctx context.Context) error {
	return ap.returnFunds(ctx, true)
}

func (ap *AccountPool) returnFunds(ctx context.Context, lock bool) error {
	if lock {
		ap.mu.Lock()
		defer ap.mu.Unlock()
	}
	if !ap.isFundingEnabled() {
		log.Debug().
			Msg("account funding is disabled, skipping returning funds from sending accounts")
		return nil
	}

	if !ap.isRefundingEnabled() {
		log.Debug().
			Bool("refundRemainingFunds", ap.cfg.RefundRemainingFunds).
			Msg("account refunding is disabled, skipping returning funds from sending accounts")
		return nil
	}

	log.Info().
		Msg("returning funds from sending accounts back to the funding account")

	ethTransferGas := big.NewInt(21000)
	err := ap.clientRateLimiter.Wait(ctx)
	if err != nil {
		return err
	}
	gasPrice, _ := ap.getSuggestedGasPrices(ctx)
	txFee := new(big.Int).Mul(ethTransferGas, gasPrice)
	// triple the txFee to account for gas price fluctuations and
	// different ways to charge transactions, like op networks
	// that charge for the l1 transaction
	biasFee := big.NewInt(0).Mul(txFee, big.NewInt(3))

	wg := sync.WaitGroup{}
	wg.Add(len(ap.accounts))

	txCh := make(chan *types.Transaction, len(ap.accounts))
	errCh := make(chan error, len(ap.accounts))

	fundingAddressHex, _ := util.GetAddressAndPrivateKeyHex(ctx, ap.fundingPrivateKey)
	fundingAddress := common.HexToAddress(fundingAddressHex)

	err = ap.clientRateLimiter.Wait(ctx)
	if err != nil {
		return err
	}
	balanceBefore, err := ap.client.BalanceAt(ctx, fundingAddress, nil)
	if err != nil {
		log.Error().Err(err).Msg("unable to get funding account balance")
		return err
	}
	log.Debug().
		Stringer("address", fundingAddress).
		Stringer("balance", balanceBefore).
		Msg("funding account balance before funds returned")

	for i := range len(ap.accounts) {
		go func(accIdx int) {
			defer wg.Done()
			// if account is the funding account, skip it
			if ap.accounts[i].address.String() == fundingAddress.String() {
				return
			}

			if ap.accounts[i].funded {
				// check if account has enough balance to pay the transfer fee
				iErr := ap.clientRateLimiter.Wait(ctx)
				if iErr != nil {
					errCh <- fmt.Errorf("failed to wait rate limit to get balance for acc %s: %w", ap.accounts[i].address.String(), iErr)
					return
				}
				balance, iErr := ap.client.BalanceAt(ctx, ap.accounts[i].address, nil)
				if iErr != nil {
					errCh <- fmt.Errorf("failed to check account balance for acc %s: %w", ap.accounts[i].address.String(), iErr)
					return
				}
				if balance.Cmp(txFee) <= 0 {
					return
				}

				// subtract the transfer fee from the balance
				amount := new(big.Int).Sub(balance, biasFee)

				// get pending nonce for account
				iErr = ap.clientRateLimiter.Wait(ctx)
				if iErr != nil {
					errCh <- fmt.Errorf("failed to wait rate limit to get nonce for acc %s: %w", ap.accounts[i].address.String(), iErr)
					return
				}
				pendingNonce, iErr := ap.client.PendingNonceAt(ctx, ap.accounts[i].address)
				if iErr != nil {
					errCh <- fmt.Errorf("failed to get nonce for acc %s: %w", ap.accounts[i].address.String(), iErr)
					return
				}
				ap.accounts[i].nonce = pendingNonce

				// loop to send tx in case we need to readjust the amount due to
				// gas price fluctuations and the tx fee
				for {
					// create the transaction to return the funds
					signedTx, iErr := ap.createEOATransferTx(ctx, ap.accounts[i].privateKey, &ap.accounts[i].nonce, fundingAddress, amount)
					if iErr != nil {
						errCh <- fmt.Errorf("failed to create tx to return balance from acc %s to %s: %w", ap.accounts[i].address.String(), fundingAddressHex, iErr)
						return
					}

					log.Debug().
						Stringer("from", ap.accounts[i].address).
						Str("to", fundingAddressHex).
						Stringer("amount", amount).
						Stringer("balance", balance).
						Stringer("txHash", signedTx.Hash()).
						Msg("returning funds")

					// send the transaction to return the funds
					iErr = ap.clientRateLimiter.Wait(ctx)
					if iErr != nil {
						errCh <- fmt.Errorf("failed to check wait rate limit to send transaction for acc %s: %w", ap.accounts[i].address.String(), iErr)
						return
					}
					iErr = ap.client.SendTransaction(ctx, signedTx)
					if iErr != nil {
						if strings.Contains(iErr.Error(), "overshot") {
							log.Info().
								Err(iErr).
								Stringer("from", ap.accounts[i].address).
								Str("to", fundingAddressHex).
								Stringer("amount", amount).
								Stringer("balance", balance).
								Msg("transaction amount overshot, adjusting amount and retrying")

							// if the amount is too high, we need to adjust it
							errArr := strings.Split(iErr.Error(), "overshot")
							if len(errArr) < 2 {
								log.Error().
									Err(iErr).
									Stringer("from", ap.accounts[i].address).
									Str("to", fundingAddressHex).
									Stringer("amount", amount).
									Stringer("balance", balance).
									Msg("unable to adjust amount due to overshot error")
								errCh <- fmt.Errorf("failed to adjust amount due to overshot error: %w", iErr)
								return
							}

							// parse the new amount from the error message
							overshotAmountStr := strings.TrimSpace(errArr[len(errArr)-1])
							overshotAmount, ok := new(big.Int).SetString(overshotAmountStr, 10)
							if !ok {
								log.Error().
									Err(iErr).
									Stringer("from", ap.accounts[i].address).
									Str("to", fundingAddressHex).
									Stringer("amount", amount).
									Stringer("balance", balance).
									Msg("unable to parse overshot amount from error message")
								errCh <- fmt.Errorf("failed to parse overshot amount from error message: %w", iErr)
								return
							}
							// reduce all overshot amount
							amount.Sub(amount, overshotAmount)
							// reduce the tx fee again to help with gas price fluctuations
							amount.Sub(amount, txFee)

							continue
						}

						log.Error().
							Err(iErr).
							Stringer("from", ap.accounts[i].address).
							Str("to", fundingAddressHex).
							Stringer("amount", amount).
							Stringer("balance", balance).
							Interface("tx", signedTx).
							Msg("unable to send return balance transaction")
						errCh <- fmt.Errorf("failed to send tx to return balance from acc %s to %s: %w", ap.accounts[i].address.String(), fundingAddressHex, iErr)
						return
					}

					txCh <- signedTx
					break
				}

			}
		}(i)
	}

	wg.Wait()

	close(errCh)
	close(txCh)

	for err := range errCh {
		if err != nil {
			return err
		}
	}

	for tx := range txCh {
		if tx != nil {
			log.Debug().
				Stringer("address", tx.To()).
				Stringer("txHash", tx.Hash()).
				Msg("transaction to return funds sent")

			_, err = util.WaitReceiptWithTimeout(ctx, ap.client, tx.Hash(), time.Minute)
			if err != nil {
				log.Error().
					Stringer("address", tx.To()).
					Stringer("txHash", tx.Hash()).
					Msg("transaction to return funds failed")
				return err
			}
		}
	}

	err = ap.clientRateLimiter.Wait(ctx)
	if err != nil {
		return err
	}
	balanceAfter, err := ap.client.BalanceAt(ctx, fundingAddress, nil)
	if err != nil {
		log.Error().Err(err).Msg("unable to get funding account balance")
		return err
	}

	log.Debug().
		Stringer("address", fundingAddress).
		Stringer("previousBalance", balanceBefore).
		Stringer("currentBalance", balanceAfter).
		Msg("all accounts funds returned")

	return nil
}

// Nonces returns the nonces of all accounts in the pool.
func (ap *AccountPool) Nonces(ctx context.Context, onlyUsed bool) *sync.Map {
	ap.mu.Lock()
	defer ap.mu.Unlock()

	nonces := &sync.Map{}
	for _, account := range ap.accounts {
		if account.nonce == account.startNonce {
			if onlyUsed {
				continue
			}
		}

		if len(account.reusableNonces) > 0 {
			nonces.Store(account.address, account.reusableNonces[len(account.reusableNonces)-1])
		} else {
			nonces.Store(account.address, account.nonce)
		}
	}
	return nonces
}

// NoncesOf returns the start nonce and current nonce for a given address.
func (ap *AccountPool) NoncesOf(address common.Address) (startNonce, nonce uint64) {
	ap.mu.Lock()
	defer ap.mu.Unlock()

	accountPos, found := ap.accountsPositions[address]
	if !found {
		return 0, 0
	}
	if accountPos > len(ap.accounts)-1 {
		return 0, 0
	}

	startNonce = ap.accounts[accountPos].startNonce
	nonce = ap.accounts[accountPos].nonce

	return startNonce, nonce
}

// Next returns the next account in the pool.
func (ap *AccountPool) Next(ctx context.Context) (Account, error) {
	ap.mu.Lock()
	defer ap.mu.Unlock()
	if len(ap.accounts) == 0 {
		return Account{}, fmt.Errorf("no accounts available")
	}
	account := ap.accounts[ap.currentAccountIndex]

	_, err := ap.fundAccountIfNeeded(ctx, account, nil, true)
	if err != nil {
		return Account{}, err
	}
	account.funded = true

	accCopy := *account

	// Check if the account has a reusable nonce
	if len(account.reusableNonces) > 0 {
		account.nonce = account.reusableNonces[0]
		account.reusableNonces = account.reusableNonces[1:]
	} else {
		account.nonce++
	}

	// move current account index to next account
	ap.currentAccountIndex++
	if ap.currentAccountIndex >= len(ap.accounts) {
		ap.currentAccountIndex = 0
	}
	return accCopy, nil
}

// SetFundingAmount updates the funding amount for the pool.
func (ap *AccountPool) SetFundingAmount(amount *big.Int) {
	ap.fundingAmount = amount
}

func (ap *AccountPool) accountMustBeFunded(ctx context.Context, account *Account) (bool, error) {
	// If funding amount is zero, skip funding entirely
	if !ap.isFundingEnabled() {
		return false, nil
	}

	// if account is funded, return it
	if account.funded {
		return false, nil
	}

	if ap.cfg.CheckBalanceBeforeFunding {
		// Check if the account has enough balance
		err := ap.clientRateLimiter.Wait(ctx)
		if err != nil {
			return false, err
		}
		balance, err := ap.client.BalanceAt(ctx, account.address, nil)
		if err != nil {
			return false, fmt.Errorf("failed to check account balance: %w", err)
		}
		// if account has enough balance
		if balance.Cmp(ap.fundingAmount) >= 0 {
			return false, nil
		}
	}
	return true, nil
}

func (ap *AccountPool) fundAccountIfNeeded(ctx context.Context, account *Account, forcedNonce *uint64, waitToFund bool) (*types.Transaction, error) {
	if mustBeFunded, err := ap.accountMustBeFunded(ctx, account); err != nil || !mustBeFunded {
		return nil, err
	}

	// Fund the account
	tx, err := ap.fund(ctx, account.address, forcedNonce, waitToFund)
	if err != nil {
		return nil, fmt.Errorf("failed to fund account: %w", err)
	}

	if waitToFund {
		err := ap.clientRateLimiter.Wait(ctx)
		if err != nil {
			return nil, err
		}
		balance, err := ap.client.BalanceAt(ctx, account.address, nil)
		if err != nil {
			return tx, fmt.Errorf("failed to check account balance: %w", err)
		}
		log.Debug().
			Stringer("address", account.address).
			Stringer("balance", balance).
			Msg("account funded")
	}
	return tx, nil
}

func (ap *AccountPool) fund(ctx context.Context, addr common.Address, forcedNonce *uint64, waitToFund bool) (*types.Transaction, error) {
	// Fund the account
	signedTx, err := ap.createEOATransferTx(ctx, ap.fundingPrivateKey, forcedNonce, addr, ap.fundingAmount)
	if err != nil {
		log.Error().Err(err).Msg("unable to create EOA Transfer tx")
		return nil, err
	}
	log.Debug().
		Stringer("address", addr).
		Uint64("amount", ap.fundingAmount.Uint64()).
		Msg("waiting account to get funded")
	err = ap.clientRateLimiter.Wait(ctx)
	if err != nil {
		return nil, err
	}
	err = ap.client.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Error().Err(err).Msg("unable to send transaction")
		return nil, err
	}

	// Wait for the transaction to be mined
	if waitToFund {
		receipt, err := util.WaitReceipt(ctx, ap.client, signedTx.Hash())
		if err != nil {
			log.Error().
				Stringer("address", addr).
				Stringer("txHash", signedTx.Hash()).
				Msg("failed to wait for transaction to be mined")
			return nil, err
		}

		if receipt.Status != types.ReceiptStatusSuccessful {
			log.Error().
				Stringer("address", addr).
				Stringer("txHash", receipt.TxHash).
				Msg("failed to wait for transaction to be mined")
			return nil, fmt.Errorf("transaction failed")
		}
	}

	return signedTx, nil
}

func (ap *AccountPool) createEOATransferTx(ctx context.Context, sender *ecdsa.PrivateKey, forcedNonce *uint64, receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	tops, err := bind.NewKeyedTransactorWithChainID(sender, ap.chainID)
	if err != nil {
		log.Error().Err(err).Msg("unable create transaction signer")
		return nil, err
	}
	tops.GasLimit = uint64(21000)
	tops = ap.configureTransactOpts(ctx, tops)

	var nonce uint64

	if forcedNonce != nil {
		nonce = *forcedNonce
	} else {
		err = ap.clientRateLimiter.Wait(ctx)
		if err != nil {
			return nil, err
		}
		nonce, err = ap.client.PendingNonceAt(ctx, tops.From)
		if err != nil {
			log.Error().
				Err(err).
				Msg("unable to get pending nonce")
			return nil, err
		}
	}

	var tx *types.Transaction
	if ap.cfg.LegacyTxMode {
		tx = types.NewTx(&types.LegacyTx{
			Nonce:    nonce,
			To:       &receiver,
			Value:    amount,
			Gas:      tops.GasLimit,
			GasPrice: tops.GasPrice,
			Data:     nil,
		})
	} else {
		dynamicFeeTx := &types.DynamicFeeTx{
			ChainID:   ap.chainID,
			Nonce:     nonce,
			To:        &receiver,
			Gas:       tops.GasLimit,
			GasFeeCap: tops.GasFeeCap,
			GasTipCap: tops.GasTipCap,
			Data:      nil,
			Value:     amount,
		}
		tx = types.NewTx(dynamicFeeTx)
	}

	signedTx, err := tops.Signer(tops.From, tx)
	if err != nil {
		log.Error().Err(err).Msg("unable to sign transaction")
		return nil, err
	}

	return signedTx, nil
}

func (ap *AccountPool) isFundingEnabled() bool {
	return !ap.isCallOnly() && ap.hasFundingAmount()
}

func (ap *AccountPool) isCallOnly() bool {
	return ap.cfg.EthCallOnly
}

func (ap *AccountPool) hasFundingAmount() bool {
	return ap.fundingAmount != nil && ap.fundingAmount.Cmp(big.NewInt(0)) > 0
}

func (ap *AccountPool) isRefundingEnabled() bool {
	if !ap.isFundingEnabled() {
		log.Debug().
			Msg("refund remaining funds is disabled because funding is disabled")
		return false
	}

	shouldRefund := ap.cfg.RefundRemainingFunds
	if !shouldRefund {
		log.Debug().
			Msg("refund remaining funds is disabled")
		return false
	}

	return true
}

func (ap *AccountPool) biasGasPrice(price *big.Int) *big.Int {
	if ap.cfg.GasPriceMultiplier == nil {
		return price
	}
	gasPriceFloat := new(big.Float).SetInt(price)
	gasPriceFloat.Mul(gasPriceFloat, ap.cfg.GasPriceMultiplier)
	result := new(big.Int)
	gasPriceFloat.Int(result)
	return result
}

func (ap *AccountPool) getSuggestedGasPrices(ctx context.Context) (*big.Int, *big.Int) {
	var gasPrice *big.Int
	gasTipCap := big.NewInt(0)
	var err error

	if ap.cfg.LegacyTxMode {
		if ap.cfg.ForceGasPrice != 0 {
			gasPrice = new(big.Int).SetUint64(ap.cfg.ForceGasPrice)
		} else {
			gasPrice, err = ap.client.SuggestGasPrice(ctx)
			if err != nil {
				log.Error().Err(err).Msg("unable to suggest gas price")
				return big.NewInt(0), big.NewInt(0)
			}
			gasPrice = ap.biasGasPrice(gasPrice)
		}
	} else {
		// Handle tip cap
		if ap.cfg.ForcePriorityGasPrice != 0 {
			gasTipCap = new(big.Int).SetUint64(ap.cfg.ForcePriorityGasPrice)
		} else if ap.cfg.ChainSupportBaseFee {
			gasTipCap, err = ap.client.SuggestGasTipCap(ctx)
			if err != nil {
				log.Error().Err(err).Msg("unable to suggest gas tip cap")
				return big.NewInt(0), big.NewInt(0)
			}
			gasTipCap = ap.biasGasPrice(gasTipCap)
		} else {
			log.Fatal().
				Msg("Chain does not support base fee. Please set priority-gas-price flag with a value to use for gas tip cap")
		}

		// Handle gas price / max fee
		if ap.cfg.ForceGasPrice != 0 {
			gasPrice = new(big.Int).SetUint64(ap.cfg.ForceGasPrice)
		} else if ap.cfg.ChainSupportBaseFee {
			gasPrice, err = ap.client.SuggestGasPrice(ctx)
			if err != nil {
				log.Error().Err(err).Msg("unable to suggest gas price")
				return big.NewInt(0), big.NewInt(0)
			}
			gasPrice = ap.biasGasPrice(gasPrice)
		} else {
			log.Fatal().
				Msg("Chain does not support base fee. Please set gas-price flag with a value to use for max fee per gas")
		}
	}

	return gasPrice, gasTipCap
}

func (ap *AccountPool) configureTransactOpts(ctx context.Context, tops *bind.TransactOpts) *bind.TransactOpts {
	gasPrice, gasTipCap := ap.getSuggestedGasPrices(ctx)
	tops.GasPrice = gasPrice

	if ap.cfg.LegacyTxMode {
		return tops
	}

	tops.GasPrice = nil
	tops.GasFeeCap = gasPrice
	tops.GasTipCap = gasTipCap

	if tops.GasTipCap.Cmp(tops.GasFeeCap) == 1 {
		tops.GasTipCap = new(big.Int).Set(tops.GasFeeCap)
	}

	return tops
}
