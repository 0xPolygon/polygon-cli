package loadtest

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

// Structure used by the account pool to control the
// current state of an account
type Account struct {
	address        common.Address
	privateKey     *ecdsa.PrivateKey
	startNonce     uint64
	nonce          uint64
	funded         bool
	reusableNonces []uint64
}

// Creates a new account with the given private key.
// The client is used to get the nonce of the account.
func newAccount(ctx context.Context, client *ethclient.Client, privateKey *ecdsa.PrivateKey) (*Account, error) {
	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(ctx, address)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	return &Account{
		privateKey:     privateKey,
		address:        address,
		startNonce:     nonce,
		nonce:          nonce,
		funded:         false,
		reusableNonces: make([]uint64, 0),
	}, nil
}

// Returns the address of the account
func (a *Account) Address(ctx context.Context) common.Address {
	return a.address
}

// Returns the private key of the account
func (a *Account) PrivateKey(ctx context.Context) *ecdsa.PrivateKey {
	return a.privateKey
}

// Returns the nonce of the account
func (a *Account) Nonce(ctx context.Context) uint64 {
	return a.nonce
}

// Structure to control accounts used by the tests
type AccountPool struct {
	accounts          []Account
	accountsPositions map[common.Address]int

	client            *ethclient.Client
	clientRateLimiter *rate.Limiter

	mu                  sync.Mutex
	currentAccountIndex int
	fundingPrivateKey   *ecdsa.PrivateKey
	fundingAmount       *big.Int
	chainID             *big.Int

	latestBlockNumber uint64
	pendingTxsCache   *uint64
}

// Creates a new account pool with the given funding private key.
// The funding private key is used to fund the accounts in the pool.
// The funding amount is the amount of ether to send to each account.
// The client is used to interact with the network to get account information
// and also to send transactions to fund accounts.
func NewAccountPool(ctx context.Context, client *ethclient.Client, fundingPrivateKey *ecdsa.PrivateKey, fundingAmount *big.Int) (*AccountPool, error) {
	if fundingPrivateKey == nil {
		panic("fundingPrivateKey cannot be nil")
	}

	if fundingAmount == nil {
		panic("fundingAmount cannot be nil")
	}

	if fundingAmount.Cmp(big.NewInt(0)) <= 0 {
		panic("fundingAmount must be greater than 0")
	}

	if client == nil {
		panic("client cannot be nil")
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

	return &AccountPool{
		currentAccountIndex: 0,
		client:              client,
		accounts:            make([]Account, 0),
		fundingPrivateKey:   fundingPrivateKey,
		fundingAmount:       fundingAmount,
		chainID:             chainID,
		accountsPositions:   make(map[common.Address]int),
		latestBlockNumber:   latestBlockNumber,
		clientRateLimiter:   rate.NewLimiter(rate.Every(50*time.Millisecond), 1),
	}, nil
}

// Adds N random accounts to the pool
func (ap *AccountPool) AddRandomN(ctx context.Context, n uint64) error {
	for i := uint64(0); i < n; i++ {
		err := ap.AddRandom(ctx)
		if err != nil {
			return fmt.Errorf("failed to add random account: %w", err)
		}
	}
	return nil
}

// Adds a random account to the pool
func (ap *AccountPool) AddRandom(ctx context.Context) error {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return fmt.Errorf("failed to generate private key: %w", err)
	}

	return ap.Add(ctx, privateKey)
}

// Adds multiple accounts to the pool with the given private keys
func (ap *AccountPool) AddN(ctx context.Context, privateKeys ...*ecdsa.PrivateKey) error {
	for _, privateKey := range privateKeys {
		err := ap.Add(ctx, privateKey)
		if err != nil {
			return fmt.Errorf("failed to add account: %w", err)
		}
	}

	return nil
}

// Adds an account to the pool with the given private key
func (ap *AccountPool) Add(ctx context.Context, privateKey *ecdsa.PrivateKey) error {
	ap.mu.Lock()
	defer ap.mu.Unlock()

	account, err := newAccount(ctx, ap.client, privateKey)
	if err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}

	addressHex, privateKeyHex := getAddressAndPrivateKeyHex(ctx, privateKey)
	log.Debug().
		Str("address", addressHex).
		Str("privateKey", privateKeyHex).
		Uint64("nonce", account.nonce).
		Msg("adding account to pool")

	ap.accounts = append(ap.accounts, *account)
	ap.accountsPositions[account.address] = len(ap.accounts) - 1
	return nil
}

// Adds a reusable nonce to the account with the given address
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
		Str("address", address.String()).
		Uint64("nonce", nonce).
		Any("reusableNonces", ap.accounts[accountPos].reusableNonces).
		Msg("reusable nonce added to account")

	return nil
}

// Refreshes the nonce with the PendingNonceAt for the given address
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
	nonce, err := ap.client.PendingNonceAt(context.Background(), address)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %w", err)
	}

	ap.accounts[accountPos].nonce = nonce

	log.Debug().
		Str("address", address.String()).
		Uint64("nonce", nonce).
		Msg("nonce refreshed")

	return nil
}

// for each account, using the internally controlled nonce, compares it to the
// network pending nonce to knows how many transactions the network behind for the
// specific account, them sum all the pending transactions differences
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
	accountsClone := make([]Account, len(ap.accounts))
	copy(accountsClone, ap.accounts)
	ap.mu.Unlock()

	pendingTxCh := make(chan uint64, len(accountsClone))
	errCh := make(chan error, len(accountsClone))

	wg := sync.WaitGroup{}
	wg.Add(len(accountsClone))

	for i := range accountsClone {
		go func(account Account) {
			defer wg.Done()
			err := ap.clientRateLimiter.Wait(ctx)
			if err != nil {
				errCh <- fmt.Errorf("failed to wait rate limit to get pending nonce for acc %s: %w", account.address.String(), err)
				return
			}
			pendingNonce, err := ap.client.NonceAt(ctx, account.address, nil)
			if err != nil {
				errCh <- fmt.Errorf("failed to get pending nonce for acc %s: %w", account.address.String(), err)
				return
			}
			pendingTxs := pendingNonce - account.nonce
			pendingTxCh <- pendingTxs
		}(accountsClone[i])
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

// Funds all accounts in the pool
func (ap *AccountPool) FundAccounts(ctx context.Context) error {
	ap.mu.Lock()
	defer ap.mu.Unlock()

	log.Debug().
		Msg("Funding all sending accounts")

	wg := sync.WaitGroup{}
	wg.Add(len(ap.accounts))

	txCh := make(chan *types.Transaction, len(ap.accounts))
	errCh := make(chan error, len(ap.accounts))

	tops, err := bind.NewKeyedTransactorWithChainID(ap.fundingPrivateKey, ap.chainID)
	if err != nil {
		log.Error().Err(err).Msg("Unable create transaction signer")
		return err
	}

	err = ap.clientRateLimiter.Wait(ctx)
	if err != nil {
		return err
	}
	nonce, err := ap.client.PendingNonceAt(ctx, tops.From)
	if err != nil {
		log.Error().Err(err).Msg("Unable to get nonce")
	}

	err = ap.clientRateLimiter.Wait(ctx)
	if err != nil {
		return err
	}
	balance, err := ap.client.BalanceAt(ctx, tops.From, nil)
	if err != nil {
		log.Error().Err(err).Msg("Unable to get funding address balance")
	}

	totalBalanceNeeded := new(big.Int).Mul(ap.fundingAmount, big.NewInt(int64(len(ap.accounts))))
	totalFeeNeeded := new(big.Int).Mul(big.NewInt(21000), big.NewInt(int64(len(ap.accounts))))
	fudgeAmountNeeded := new(big.Int).Mul(big.NewInt(1000000000), big.NewInt(int64(len(ap.accounts))))

	totalNeeded := new(big.Int).Add(totalBalanceNeeded, totalFeeNeeded)
	totalNeeded.Add(totalNeeded, fudgeAmountNeeded)

	if balance.Cmp(totalBalanceNeeded) <= 0 {
		errMsg := "Funding account balance can't cover the funding amount for all accounts"
		log.Error().
			Str("address", tops.From.Hex()).
			Str("balance", balance.String()).
			Str("totalNeeded", totalNeeded.String()).
			Msg(errMsg)
		return errors.New(errMsg)
	}

	for i := range ap.accounts {
		accountToFund := ap.accounts[i]

		// if account is the funding account, skip it
		if accountToFund.address == tops.From {
			continue
		}

		go func(forcedNonce uint64, account Account) {
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

	for tx := range txCh {
		if tx != nil {
			log.Debug().
				Str("address", tx.To().Hex()).
				Str("txHash", tx.Hash().Hex()).
				Msg("transaction to fund account sent")

			_, err := ap.waitMined(ctx, tx)
			if err != nil {
				log.Error().
					Str("address", tx.To().Hex()).
					Str("txHash", tx.Hash().Hex()).
					Msgf("transaction to fund account failed")
				return err
			}
		}
	}

	log.Debug().
		Msg("All accounts funded")

	return nil
}

// Return the funds from all accounts in the pool to the funding account
func (ap *AccountPool) ReturnFunds(ctx context.Context) error {
	ap.mu.Lock()
	defer ap.mu.Unlock()

	ltp := inputLoadTestParams

	log.Debug().
		Msg("Returning funds from sending addresses to funding address")

	ethTransferGas := big.NewInt(21000)
	err := ap.clientRateLimiter.Wait(ctx)
	if err != nil {
		return err
	}
	var pricePerGas *big.Int
	if *ltp.LegacyTransactionMode {
		gasPrice, iErr := ap.client.SuggestGasPrice(ctx)
		if iErr != nil {
			log.Error().Err(iErr).Msg("Unable to get gas price")
			return iErr
		}
		pricePerGas = gasPrice
	} else {
		feeMutex.RLock()
		pricePerGas = ltp.MaxFeePerGas
		feeMutex.RUnlock()
	}
	txFee := new(big.Int).Mul(ethTransferGas, pricePerGas)
	// double the txFee to account for gas price fluctuations and
	// different ways to charge transactions, like op networks
	// that charge for the l1 transaction
	txFee.Add(txFee, txFee)

	wg := sync.WaitGroup{}
	wg.Add(len(ap.accounts))

	txCh := make(chan *types.Transaction, len(ap.accounts))
	errCh := make(chan error, len(ap.accounts))

	fundingAddressHex, _ := getAddressAndPrivateKeyHex(ctx, ap.fundingPrivateKey)
	fundingAddress := common.HexToAddress(fundingAddressHex)

	err = ap.clientRateLimiter.Wait(ctx)
	if err != nil {
		return err
	}
	balanceBefore, err := ap.client.BalanceAt(ctx, fundingAddress, nil)
	if err != nil {
		log.Error().Err(err).Msg("Unable to get funding address balance")
		return err
	}
	log.Debug().
		Str("address", fundingAddress.Hex()).
		Str("balance", balanceBefore.String()).
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
				amount := new(big.Int).Sub(balance, txFee)

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

				// create the transaction to return the funds
				signedTx, iErr := ap.createEOATransferTx(ctx, ap.accounts[i].privateKey, &ap.accounts[i].nonce, fundingAddress, amount)
				if iErr != nil {
					errCh <- fmt.Errorf("failed to create tx to return balance from acc %s to %s: %w", ap.accounts[i].address.String(), fundingAddressHex, iErr)
					return
				}

				log.Debug().
					Str("from", ap.accounts[i].address.Hex()).
					Str("to", fundingAddressHex).
					Str("amount", amount.String()).
					Str("balance", balance.String()).
					Str("txHash", signedTx.Hash().String()).
					Msg("returning funds")

				// send the transaction to return the funds
				iErr = ap.clientRateLimiter.Wait(ctx)
				if iErr != nil {
					errCh <- fmt.Errorf("failed to check wait rate limit to send transaction for acc %s: %w", ap.accounts[i].address.String(), iErr)
					return
				}
				iErr = ap.client.SendTransaction(ctx, signedTx)
				if iErr != nil {
					log.Debug().
						Str("from", ap.accounts[i].address.Hex()).
						Str("to", fundingAddressHex).
						Str("amount", amount.String()).
						Str("balance", balance.String()).
						Interface("tx", signedTx).
						Msg("Unable to send return balance transaction")
					errCh <- fmt.Errorf("failed to send tx to return balance from acc %s to %s: %w", ap.accounts[i].address.String(), fundingAddressHex, iErr)
					return
				}

				txCh <- signedTx
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
				Str("address", tx.To().Hex()).
				Str("txHash", tx.Hash().Hex()).
				Msg("transaction to return funds sent")

			_, err = ap.waitMined(ctx, tx)
			if err != nil {
				log.Error().
					Str("address", tx.To().Hex()).
					Str("txHash", tx.Hash().Hex()).
					Msgf("transaction to return funds failed")
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
		log.Error().Err(err).Msg("Unable to get funding address balance")
		return err
	}

	log.Debug().
		Str("address", fundingAddress.Hex()).
		Str("previousBalance", balanceBefore.String()).
		Str("currentBalance", balanceAfter.String()).
		Msg("All accounts funds returned")

	return nil
}

// Returns the nonces of all accounts in the pool
func (ap *AccountPool) Nonces(ctx context.Context) map[common.Address]uint64 {
	ap.mu.Lock()
	defer ap.mu.Unlock()

	nonces := make(map[common.Address]uint64)
	for _, account := range ap.accounts {
		if len(account.reusableNonces) > 0 {
			nonces[account.address] = account.reusableNonces[len(account.reusableNonces)-1]
		} else {
			nonces[account.address] = account.nonce
		}
	}
	return nonces
}

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

// Returns the next account in the pool
func (ap *AccountPool) Next(ctx context.Context) (Account, error) {
	ap.mu.Lock()
	defer ap.mu.Unlock()
	if len(ap.accounts) == 0 {
		return Account{}, fmt.Errorf("no accounts available")
	}
	account := ap.accounts[ap.currentAccountIndex]

	// if test is call only, there is no need to fund accounts, return it
	if !*inputLoadTestParams.CallOnly {
		_, err := ap.fundAccountIfNeeded(ctx, account, nil, true)
		if err != nil {
			return Account{}, err
		}
	}
	ap.accounts[ap.currentAccountIndex].funded = true

	// Check if the account has a reusable nonce
	if len(account.reusableNonces) > 0 {
		account.nonce = ap.accounts[ap.currentAccountIndex].reusableNonces[0]
		ap.accounts[ap.currentAccountIndex].reusableNonces = ap.accounts[ap.currentAccountIndex].reusableNonces[1:]
	} else {
		ap.accounts[ap.currentAccountIndex].nonce++
	}

	// move current account index to next account
	ap.currentAccountIndex++
	if ap.currentAccountIndex >= len(ap.accounts) {
		ap.currentAccountIndex = 0
	}
	log.Debug().
		Str("address", account.address.Hex()).
		Str("nonce", fmt.Sprintf("%d", account.nonce)).
		Msg("account returned from pool")
	return account, nil
}

// Checks multiple conditions of the account and funds it if needed
func (ap *AccountPool) fundAccountIfNeeded(ctx context.Context, account Account, forcedNonce *uint64, waitToFund bool) (*types.Transaction, error) {
	// if account is funded, return it
	if account.funded {
		return nil, nil
	}

	// Check if the account must be funded
	err := ap.clientRateLimiter.Wait(ctx)
	if err != nil {
		return nil, err
	}
	balance, err := ap.client.BalanceAt(ctx, account.address, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check account balance: %w", err)
	}
	// if account has enough balance
	if balance.Cmp(ap.fundingAmount) >= 0 {
		return nil, nil
	}

	// Fund the account
	log.Debug().
		Str("address", account.address.Hex()).
		Str("balance", balance.String()).
		Msg("account needs to be funded")
	tx, err := ap.fund(ctx, account, forcedNonce, waitToFund)
	if err != nil {
		return nil, fmt.Errorf("failed to fund account: %w", err)
	}

	if waitToFund {
		err := ap.clientRateLimiter.Wait(ctx)
		if err != nil {
			return nil, err
		}
		balance, err = ap.client.BalanceAt(ctx, account.address, nil)
		if err != nil {
			return tx, fmt.Errorf("failed to check account balance: %w", err)
		}
		log.Debug().
			Str("address", account.address.Hex()).
			Str("balance", balance.String()).
			Msg("account funded")
	}
	return tx, nil
}

// Funds the account
func (ap *AccountPool) fund(ctx context.Context, acc Account, forcedNonce *uint64, waitToFund bool) (*types.Transaction, error) {
	// Fund the account
	signedTx, err := ap.createEOATransferTx(ctx, ap.fundingPrivateKey, forcedNonce, acc.address, ap.fundingAmount)
	if err != nil {
		log.Error().Err(err).Msg("Unable to create EOA Transfer tx")
		return nil, err
	}
	log.Debug().
		Str("address", acc.address.Hex()).
		Uint64("amount", ap.fundingAmount.Uint64()).
		Msgf("waiting account to get funded")
	err = ap.clientRateLimiter.Wait(ctx)
	if err != nil {
		return nil, err
	}
	err = ap.client.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Error().Err(err).Msg("Unable to send transaction")
		return nil, err
	}

	// Wait for the transaction to be mined
	if waitToFund {
		receipt, err := ap.waitMined(ctx, signedTx)
		if err != nil {
			log.Error().
				Str("address", acc.address.Hex()).
				Str("txHash", signedTx.Hash().Hex()).
				Msgf("failed to wait for transaction to be mined")
			return nil, err
		}

		if receipt.Status != types.ReceiptStatusSuccessful {
			log.Error().
				Str("address", acc.address.Hex()).
				Str("txHash", receipt.TxHash.Hex()).
				Msgf("failed to wait for transaction to be mined")
			return nil, fmt.Errorf("transaction failed")
		}
	}

	return signedTx, nil
}

func (ap *AccountPool) createEOATransferTx(ctx context.Context, sender *ecdsa.PrivateKey, forcedNonce *uint64, receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	ltp := inputLoadTestParams

	tops, err := bind.NewKeyedTransactorWithChainID(sender, ap.chainID)
	if err != nil {
		log.Error().Err(err).Msg("Unable create transaction signer")
		return nil, err
	}
	tops.GasLimit = uint64(21000)
	tops = configureTransactOpts(ctx, ap.client, tops)

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
				Msg("Unable to get pending nonce")
			return nil, err
		}
	}

	var tx *types.Transaction
	if *ltp.LegacyTransactionMode {
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
		log.Error().Err(err).Msg("Unable to sign transaction")
		return nil, err
	}

	return signedTx, nil
}

// Waits for the transaction to be mined
func (ap *AccountPool) waitMined(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	receipt, err := bind.WaitMined(ctxTimeout, ap.client, tx)
	if err != nil {
		log.Error().
			Str("txHash", tx.Hash().Hex()).
			Err(err).
			Msg("Unable to wait for transaction to be mined")
		return nil, err
	}
	return receipt, nil
}

// Returns the address and private key of the given private key
func getAddressAndPrivateKeyHex(ctx context.Context, privateKey *ecdsa.PrivateKey) (string, string) {
	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyHex := fmt.Sprintf("0x%x", privateKeyBytes)

	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	address := crypto.PubkeyToAddress(*publicKey)

	return address.String(), privateKeyHex
}
