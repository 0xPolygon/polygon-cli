package loadtest

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
)

type Account struct {
	address    common.Address
	privateKey *ecdsa.PrivateKey
	nonce      uint64
	funded     bool
	used       bool
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
		privateKey: privateKey,
		address:    address,
		nonce:      nonce,
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
	mu                  sync.Mutex
	client              *ethclient.Client
	accounts            []Account
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
	return nil
}

func (ap *AccountPool) FundAccounts() error {
	ap.mu.Lock()
	defer ap.mu.Unlock()

	for i := range ap.accounts {
		account := ap.accounts[i]
		if !account.funded {
			err := ap.fundAccountIfNeeded(context.Background(), account)
			if err != nil {
				return fmt.Errorf("failed to fund account: %w", err)
			}
		}
	}

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

func (ap *AccountPool) Run(ctx context.Context, f func(account Account) error) error {
	account, err := ap.next(ctx)
	log.Debug().
		Str("address", account.address.Hex()).
		Str("nonce", fmt.Sprintf("%d", account.nonce)).
		Msg("account returned from pool")
	if err != nil {
		return err
	}
	err = f(account)
	if err != nil {
		return err
	}
	return nil
}

func (ap *AccountPool) next(ctx context.Context) (Account, error) {
	ap.mu.Lock()
	defer ap.mu.Unlock()
	if len(ap.accounts) == 0 {
		return Account{}, fmt.Errorf("no accounts available")
	}
	account := ap.accounts[ap.currentAccountIndex]

	err := ap.fundAccountIfNeeded(ctx, account)
	if err != nil {
		return Account{}, err
	}
	ap.accounts[ap.currentAccountIndex].funded = true
	ap.accounts[ap.currentAccountIndex].nonce++

	// move current account index to next account
	ap.currentAccountIndex++
	if ap.currentAccountIndex >= len(ap.accounts) {
		ap.currentAccountIndex = 0
	}

	return account, nil
}

func (ap *AccountPool) fundAccountIfNeeded(ctx context.Context, account Account) error {
	// if account is funded, return it
	if account.funded {
		return nil
	}

	// Check if the account must be funded
	balance, err := ap.client.BalanceAt(context.Background(), account.address, nil)
	if err != nil {
		return fmt.Errorf("failed to check account balance: %w", err)
	}
	// if account has enough balance
	if balance.Cmp(ap.fundingAmount) >= 0 {
		return nil
	}

	// Fund the account
	log.Debug().
		Str("address", account.address.Hex()).
		Str("balance", balance.String()).
		Msg("account needs to be funded")
	_, err = ap.fund(ctx, account, true)
	if err != nil {
		return fmt.Errorf("failed to fund account: %w", err)
	}

	balance, err = ap.client.BalanceAt(context.Background(), account.address, nil)
	if err != nil {
		return fmt.Errorf("failed to check account balance: %w", err)
	}
	log.Debug().
		Str("address", account.address.Hex()).
		Str("balance", balance.String()).
		Msg("account funded")

	return nil
}

func (ap *AccountPool) fund(ctx context.Context, acc Account, waitToFund bool) (*types.Transaction, error) {
	// Fund the account
	ltp := inputLoadTestParams

	tops, err := bind.NewKeyedTransactorWithChainID(ap.fundingPrivateKey, ap.chainID)
	if err != nil {
		log.Error().Err(err).Msg("Unable create transaction signer")
		return nil, err
	}
	tops.GasLimit = uint64(21000)
	tops = configureTransactOpts(ctx, ap.client, tops)

	nonce, err := ap.client.PendingNonceAt(ctx, tops.From)
	if err != nil {
		log.Error().Err(err).Msg("Unable to get nonce")
	}

	var tx *ethtypes.Transaction
	if *ltp.LegacyTransactionMode {
		tx = ethtypes.NewTx(&ethtypes.LegacyTx{
			Nonce:    nonce,
			To:       &acc.address,
			Value:    ap.fundingAmount,
			Gas:      tops.GasLimit,
			GasPrice: tops.GasPrice,
			Data:     nil,
		})
	} else {
		dynamicFeeTx := &ethtypes.DynamicFeeTx{
			ChainID:   ap.chainID,
			Nonce:     nonce,
			To:        &acc.address,
			Gas:       tops.GasLimit,
			GasFeeCap: tops.GasFeeCap,
			GasTipCap: tops.GasTipCap,
			Data:      nil,
			Value:     ap.fundingAmount,
		}
		tx = ethtypes.NewTx(dynamicFeeTx)
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
				Str("txHash", receipt.TxHash.Hex()).
				Msgf("transaction to fund account failed")
			return nil, err
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
	if receipt.Status != ethtypes.ReceiptStatusSuccessful {
		return nil, fmt.Errorf("transaction failed")
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
