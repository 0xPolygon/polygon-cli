package util

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// SendTx is a simple wrapper to send a transaction from one Ethereum address to another.
func SendTx(ctx context.Context, c *ethclient.Client, privateKey *ecdsa.PrivateKey, to *common.Address, amount *big.Int, data []byte, gasLimit uint64) error {
	// Get the chaind id.
	chainID, err := c.ChainID(ctx)
	if err != nil {
		return err
	}

	// Get the nonce.
	from := crypto.PubkeyToAddress(privateKey.PublicKey)
	var nonce uint64
	nonce, err = c.PendingNonceAt(ctx, from)
	if err != nil {
		return err
	}

	// Get suggested gas price.
	var gasPrice *big.Int
	gasPrice, err = c.SuggestGasPrice(ctx)
	if err != nil {
		return err
	}

	// Create and sign the transaction.
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		To:       to,
		Value:    amount,
		Data:     data,
	})
	var signedTx *types.Transaction
	signedTx, err = types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return err
	}

	// Send the transaction.
	if err = c.SendTransaction(ctx, signedTx); err != nil {
		return err
	}
	if _, err = bind.WaitMined(ctx, c, signedTx); err != nil {
		return err
	}
	return nil
}
