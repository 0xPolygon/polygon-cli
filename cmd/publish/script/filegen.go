// This script generates a file containing a specified number of Ethereum transactions
// that transfer a small amount of ETH from a funding account to a newly created account.
//
// Transactions are encoded in RLP format and written to a file.
//
// The funding account is specified by its private key, and the transactions are created
// with a fixed gas limit and fee structure.
//
// The script logs the progress of transaction creation and file writing, including
// the address and private key of the newly created account.

package main

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	// number of transactions to add to the file
	numberOfTransactions = uint64(1000000)

	// chain ID for the Ethereum mainnet
	chainID = uint64(1337)

	// funding account private key in hex format
	fundingPrivateKeyHex = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

	// gas limit for the ETH transfer transactions
	transferGasLimit = uint64(21000) // Standard gas limit for ETH transfer

	// fee per gas for the transactions
	feePerGas = uint64(1000000000)

	// fee tip cap for the transactions
	feeTipCap = uint64(1000000000)

	// log every % of the transactions added to the file
	logEveryPercent = uint64(20)
	logEveryTxs     = numberOfTransactions * logEveryPercent / 100
)

type account struct {
	addr common.Address
	key  ecdsa.PrivateKey
}

func (a account) keyHex() string {
	privateKeyBytes := crypto.FromECDSA(&a.key)
	privateKeyHex := fmt.Sprintf("0x%x", privateKeyBytes)

	return privateKeyHex
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	fundingAccPrivateKey, err := loadPrivateKey(fundingPrivateKeyHex)
	checkErr(err, "failed to load funding account private key")

	fundingAcc, err := loadAccount(*fundingAccPrivateKey)
	checkErr(err, "failed to load funding account")

	acc, err := newAccount()
	checkErr(err, "failed to create new account")
	log.Info().
		Stringer("address", acc.addr).
		Str("privateKey", acc.keyHex()).
		Msg("new account created")

	transferAmount := big.NewInt(1)

	outputFilename := "txs"
	outputFilePath := filepath.Join(".", outputFilename)
	outputFilePath, err = filepath.Abs(outputFilePath)
	checkErr(err, "failed to get absolute path for file")

	f, err := deleteAndCreateFile(outputFilePath)
	checkErr(err, "failed to create temp file")
	defer f.Close()
	log.Info().
		Str("filePath", outputFilePath).
		Msg("file created for transactions")

	log.Info().
		Uint64("numberOfTransactions", numberOfTransactions).
		Msg("adding transactions to file")

	logPageSize := logEveryTxs
	if logPageSize == 0 {
		logPageSize = 1
	}

	nonce := uint64(0)
	for nonce < numberOfTransactions {
		nonce++

		tx, err := createTxToTransferETH(acc, fundingAcc, transferAmount, nonce)
		checkErr(err, "failed to write to file")

		txRLP, err := rlp(tx)
		checkErr(err, "failed to encode transaction to RLP")

		err = writeToFile(f, txRLP)
		checkErr(err, "failed to write transaction to file")

		txJSONB, err := tx.MarshalJSON()
		checkErr(err, "failed to marshal transaction to JSON")
		txJSON := string(txJSONB)

		log.Trace().
			Str("tx", txJSON).
			Str("rlp", txRLP).
			Msg("transaction written to file")

		if nonce%logPageSize == 0 {
			log.Info().
				Uint64("numberOfTransactions", nonce).
				Msg("transactions added to file")
		}
	}
	if nonce%logPageSize != 0 {
		log.Info().
			Uint64("numberOfTransactions", nonce).
			Msg("transactions added to file")
	}

	fundingAmountInWei := computeFundingAmountInWei(transferAmount)
	castSendCommand := fmt.Sprintf("cast send --private-key %s --value %d %s --rpc-url http://127.0.0.1:8545", fundingPrivateKeyHex, fundingAmountInWei, acc.addr.String())
	log.Info().
		Str("cmd", castSendCommand).
		Msg("cast send command to fund account created")
}

func computeFundingAmountInWei(transferAmount *big.Int) *big.Int {
	numberOfTransactionsBig := big.NewInt(0).SetUint64(numberOfTransactions)
	transferGasLimitBig := big.NewInt(0).SetUint64(transferGasLimit)

	feePerGasBig := big.NewInt(0).SetUint64(feePerGas)
	feeTipCapBig := big.NewInt(0).SetUint64(feeTipCap)

	feePerTx := big.NewInt(0).Mul(feePerGasBig, transferGasLimitBig)

	costPerTx := big.NewInt(0)
	costPerTx.Add(transferAmount, feePerTx)
	costPerTx.Add(costPerTx, feeTipCapBig)

	fundingAmount := big.NewInt(0).Mul(costPerTx, numberOfTransactionsBig)
	return fundingAmount
}

func deleteAndCreateFile(filePath string) (*os.File, error) {
	if _, err := os.Stat(filePath); err == nil {
		if iErr := os.Remove(filePath); iErr != nil {
			return nil, fmt.Errorf("failed to delete existing file: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to check if file exists: %w", err)
	}

	f, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	return f, nil
}

func newAccount() (*account, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate new private key: %w", err)
	}

	return loadAccount(*privateKey)
}

func loadAccount(privateKey ecdsa.PrivateKey) (*account, error) {
	pubKey := privateKey.Public().(*ecdsa.PublicKey)
	address := crypto.PubkeyToAddress(*pubKey)

	acc := &account{
		addr: address,
		key:  privateKey,
	}

	return acc, nil
}

func loadPrivateKey(privateKeyHex string) (*ecdsa.PrivateKey, error) {
	privateKeyBytes, err := crypto.HexToECDSA(privateKeyHex[2:])
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %w", err)
	}
	return privateKeyBytes, nil
}

func rlp(tx *types.Transaction) (string, error) {
	txBytes, err := tx.MarshalBinary()
	if err != nil {
		return "", fmt.Errorf("failed to marshal transaction to binary: %w", err)
	}
	return fmt.Sprintf("0x%x", txBytes), nil
}

func createTxToTransferETH(fromAcc, toAcc *account, value *big.Int, forceNonce uint64) (*types.Transaction, error) {
	tx := types.NewTx(&types.DynamicFeeTx{
		// destination
		To:    &toAcc.addr,
		Nonce: forceNonce,
		Value: value,

		// fees
		Gas:       transferGasLimit, // Standard gas limit for ETH transfer
		GasFeeCap: big.NewInt(0).SetUint64(feePerGas),
		GasTipCap: big.NewInt(0).SetUint64(feeTipCap),

		// identification
		ChainID: big.NewInt(0).SetUint64(chainID),
	})

	// signature
	signer := types.LatestSignerForChainID(big.NewInt(0).SetUint64(chainID))
	signedTx, err := types.SignTx(tx, signer, &fromAcc.key)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	return signedTx, nil
}

func writeToFile(f *os.File, data string) error {
	if _, err := f.WriteString(data + "\n"); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	return nil
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatal().
			Err(err).
			Msg(msg)
	}
}
