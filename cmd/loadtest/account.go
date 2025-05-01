package loadtest

import (
	"context"
	"crypto/ecdsa"
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
)

type Account struct {
	address        common.Address
	privateKey     *ecdsa.PrivateKey
	nonce          uint64
	funded         bool
	reusableNonces []uint64
}

func newAccount(client *ethclient.Client, privateKey *ecdsa.PrivateKey) (*Account, error) {
	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), address)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	return &Account{
		privateKey:     privateKey,
		address:        address,
		nonce:          nonce,
		funded:         false,
		reusableNonces: make([]uint64, 0),
	}, nil
}

func (a *Account) Address() common.Address {
	return a.address
}
func (a *Account) PrivateKey() *ecdsa.PrivateKey {
	return a.privateKey
}
func (a *Account) Nonce() uint64 {
	return a.nonce
}

type AccountPool struct {
	accounts          []Account
	accountsPositions map[common.Address]int

	mu                  sync.Mutex
	client              *ethclient.Client
	currentAccountIndex int
	fundingPrivateKey   *ecdsa.PrivateKey
	fundingAmount       *big.Int
	chainID             *big.Int
}

func NewAccountPool(ctx context.Context, client *ethclient.Client, fundingPrivateKey *ecdsa.PrivateKey, fundingAmount *big.Int) *AccountPool {
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
		log.Error().Err(err).Msg("Unable to get chain ID")
		return nil
	}

	return &AccountPool{
		currentAccountIndex: 0,
		client:              client,
		accounts:            make([]Account, 0),
		fundingPrivateKey:   fundingPrivateKey,
		fundingAmount:       fundingAmount,
		chainID:             chainID,
		accountsPositions:   make(map[common.Address]int),
	}
}

func (ap *AccountPool) AddRandomN(n uint64) error {
	for i := uint64(0); i < n; i++ {
		err := ap.AddRandom()
		if err != nil {
			return fmt.Errorf("failed to add random account: %w", err)
		}
	}
	return nil
}

func (ap *AccountPool) AddRandom() error {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return fmt.Errorf("failed to generate private key: %w", err)
	}

	return ap.Add(privateKey)
}

func (ap *AccountPool) Add(privateKey *ecdsa.PrivateKey) error {
	ap.mu.Lock()
	defer ap.mu.Unlock()

	addressHex, privateKeyHex := getAddressAndPrivateKeyHex(privateKey)
	log.Debug().
		Str("address", addressHex).
		Str("privateKey", privateKeyHex).
		Msg("adding account to pool")

	account, err := newAccount(ap.client, privateKey)
	if err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}

	ap.accounts = append(ap.accounts, *account)
	ap.accountsPositions[account.address] = len(ap.accounts) - 1
	return nil
}

func (ap *AccountPool) AddReusableNonce(address common.Address, nonce uint64) error {
	ap.mu.Lock()
	defer ap.mu.Unlock()

	accountPos, found := ap.accountsPositions[address]
	if !found {
		return fmt.Errorf("account not found in pool")
	}
	if accountPos >= len(ap.accounts)-1 {
		return fmt.Errorf("account position out of bounds")
	}

	ap.accounts[accountPos].reusableNonces = append(ap.accounts[accountPos].reusableNonces, nonce)

	// sort the reusable nonces ascending because we want to use the lowest nonce first
	// and we pay the price of sorting only once when adding it
	slices.Sort(ap.accounts[accountPos].reusableNonces)

	return nil
}

func (ap *AccountPool) FundAccounts(ctx context.Context) error {
	ap.mu.Lock()
	defer ap.mu.Unlock()

	wg := sync.WaitGroup{}
	wg.Add(len(ap.accounts))

	txCh := make(chan *types.Transaction, len(ap.accounts))
	errCh := make(chan error, len(ap.accounts))

	tops, err := bind.NewKeyedTransactorWithChainID(ap.fundingPrivateKey, ap.chainID)
	if err != nil {
		log.Error().Err(err).Msg("Unable create transaction signer")
		return err
	}

	nonce, err := ap.client.PendingNonceAt(ctx, tops.From)
	if err != nil {
		log.Error().Err(err).Msg("Unable to get nonce")
	}

	for i := range ap.accounts {
		accountToFund := ap.accounts[i]
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

func (ap *AccountPool) Nonces() map[common.Address]uint64 {
	ap.mu.Lock()
	defer ap.mu.Unlock()

	nonces := make(map[common.Address]uint64)
	for _, account := range ap.accounts {
		nonces[account.address] = account.nonce
	}
	return nonces
}

func (ap *AccountPool) Next(ctx context.Context) (Account, error) {
	account, err := ap.next(ctx)
	log.Debug().
		Str("address", account.address.Hex()).
		Str("nonce", fmt.Sprintf("%d", account.nonce)).
		Msg("account returned from pool")
	if err != nil {
		return Account{}, err
	}
	return account, nil
}

func (ap *AccountPool) next(ctx context.Context) (Account, error) {
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

	return account, nil
}

func (ap *AccountPool) fundAccountIfNeeded(ctx context.Context, account Account, forcedNonce *uint64, waitToFund bool) (*types.Transaction, error) {
	// if account is funded, return it
	if account.funded {
		return nil, nil
	}

	// Check if the account must be funded
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

func (ap *AccountPool) fund(ctx context.Context, acc Account, forcedNonce *uint64, waitToFund bool) (*types.Transaction, error) {
	// Fund the account
	ltp := inputLoadTestParams

	tops, err := bind.NewKeyedTransactorWithChainID(ap.fundingPrivateKey, ap.chainID)
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
		nonce, err = ap.client.PendingNonceAt(ctx, tops.From)
		if err != nil {
			log.Error().Err(err).Msg("Unable to get nonce")
		}
	}

	var tx *types.Transaction
	if *ltp.LegacyTransactionMode {
		tx = types.NewTx(&types.LegacyTx{
			Nonce:    nonce,
			To:       &acc.address,
			Value:    ap.fundingAmount,
			Gas:      tops.GasLimit,
			GasPrice: tops.GasPrice,
			Data:     nil,
		})
	} else {
		dynamicFeeTx := &types.DynamicFeeTx{
			ChainID:   ap.chainID,
			Nonce:     nonce,
			To:        &acc.address,
			Gas:       tops.GasLimit,
			GasFeeCap: tops.GasFeeCap,
			GasTipCap: tops.GasTipCap,
			Data:      nil,
			Value:     ap.fundingAmount,
		}
		tx = types.NewTx(dynamicFeeTx)
	}

	signedTx, err := tops.Signer(*ltp.FromETHAddress, tx)
	if err != nil {
		log.Error().Err(err).Msg("Unable to sign transaction")
		return nil, err
	}

	log.Debug().
		Str("address", acc.address.Hex()).
		Uint64("amount", ap.fundingAmount.Uint64()).
		Msgf("waiting account to get funded")
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

func getAddressAndPrivateKeyHex(privateKey *ecdsa.PrivateKey) (string, string) {
	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyHex := fmt.Sprintf("0x%x", privateKeyBytes)

	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	address := crypto.PubkeyToAddress(*publicKey)

	return address.String(), privateKeyHex
}
