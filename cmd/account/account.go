package account

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var AccountCmd = &cobra.Command{
	Use:               "account",
	Short:             "Utilities for interacting with an account",
	Long:              "Basic utility commands for interacting with an account",
	Args:              cobra.NoArgs,
	PersistentPreRunE: prepareRpcClient,
}

var (
	rpcClient *ethclient.Client

	nonceCommand  *cobra.Command
	fixGapCommand *cobra.Command
)

type accountArgs struct {
	rpcURL     *string
	privateKey *string
}

var inputAccountArgs = accountArgs{}

const (
	ArgPrivateKey = "private-key"
	ArgRpcURL     = "rpc-url"
)

//go:embed FixGapUsage.md
var fixGapUsage string

func prepareRpcClient(cmd *cobra.Command, args []string) error {
	var err error
	rpcURL := *inputAccountArgs.rpcURL

	rpcClient, err = ethclient.Dial(rpcURL)
	if err != nil {
		log.Error().Err(err).Msgf("Unable to Dial RPC %s: %s", rpcURL, err.Error())
		return err
	}

	if _, err = rpcClient.BlockNumber(cmd.Context()); err != nil {
		log.Error().Err(err).Msgf("Unable to get block number: %s", err.Error())
		return err
	}

	return nil
}

func accountNonceFixGap(cmd *cobra.Command, args []string) error {
	pvtKey := strings.TrimPrefix(*inputAccountArgs.privateKey, "0x")
	pk, err := crypto.HexToECDSA(pvtKey)
	if err != nil {
		log.Error().Err(err).Msgf("Invalid private key: %s", err.Error())
		return err
	}

	chainID, err := rpcClient.ChainID(cmd.Context())
	if err != nil {
		log.Error().Err(err).Msgf("Cannot get chain ID: %s", err.Error())
		return err
	}

	opts, err := bind.NewKeyedTransactorWithChainID(pk, chainID)
	if err != nil {
		log.Error().Err(err).Msgf("Cannot generate transactionOpts: %s", err.Error())
		return err
	}

	addr := opts.From

	currentNonce, err := rpcClient.NonceAt(cmd.Context(), addr, nil)
	if err != nil {
		log.Error().Err(err).Msgf("Unable to get current nonce: %s", err.Error())
		return err
	}
	log.Info().Stringer("addr", addr).Msgf("Current nonce: %d", currentNonce)

	poolNonces, maxPoolNonce, err := getPoolContent(addr)
	if err != nil {
		log.Error().Err(err).Msgf("Unable to get current nonce: %s", err.Error())
		return err
	}

	// check if there is a nonce gap
	if maxPoolNonce == 0 || currentNonce >= maxPoolNonce {
		log.Info().Stringer("addr", addr).Msg("There is no nonce gap.")
		return nil
	}
	log.Info().Stringer("addr", addr).Msgf("Nonce gap found. Max pool nonce: %d", maxPoolNonce)

	gasPrice, err := rpcClient.SuggestGasPrice(cmd.Context())
	if err != nil {
		log.Error().Err(err).Msgf("Unable to get suggested gas price: %s", err.Error())
		return err
	}

	to := &common.Address{}

	gas, err := rpcClient.EstimateGas(cmd.Context(), ethereum.CallMsg{
		To:       to,
		GasPrice: gasPrice,
		Value:    big.NewInt(1),
	})
	if err != nil {
		log.Error().Err(err).Msgf("Unable to estimate gas: %s", err.Error())
		return err
	}

	txTemplate := &types.LegacyTx{
		To:       to,
		GasPrice: gasPrice,
		Value:    big.NewInt(1),
		Gas:      gas,
	}

	var lastTx *types.Transaction
	for i := currentNonce; i < maxPoolNonce; i++ {
		txTemplate.Nonce = i

		// if nonce is already in the pool, skip it
		txHash, found := poolNonces[txTemplate.Nonce]
		if found {
			log.Info().Interface("hash", txHash).Msgf("skipping tx with nonce %d because there is already a tx in the pool with this nonce", txTemplate.Nonce)
			continue
		}

		tx := types.NewTx(txTemplate)
		signedTx, err := opts.Signer(opts.From, tx)
		log.Info().Stringer("hash", signedTx.Hash()).Msgf("sending tx with nonce %d", txTemplate.Nonce)
		if err != nil {
			log.Error().Err(err).Msgf("Unable to sign tx: %s", err.Error())
			return err
		}

		lastTx = signedTx
		err = rpcClient.SendTransaction(cmd.Context(), signedTx)
		if err != nil {
			log.Error().Err(err).Msgf("Unable to send tx: %s", err.Error())
			return err
		}
	}

	if lastTx != nil {
		log.Info().Stringer("hash", lastTx.Hash()).Msg("waiting for the last tx to get mined")
		err := WaitMineTransaction(cmd.Context(), rpcClient, lastTx, 600)
		if err != nil {
			log.Error().Err(err).Msgf("Unable to wait for last tx to get mined: %s", err.Error())
			return err
		}
		log.Info().Stringer("addr", addr).Msg("Nonce gap fixed successfully")
		currentNonce, err = rpcClient.NonceAt(cmd.Context(), addr, nil)
		if err != nil {
			log.Error().Err(err).Msgf("Unable to get current nonce: %s", err.Error())
			return err
		}
		log.Info().Stringer("addr", addr).Msgf("Current nonce: %d", currentNonce)
		return nil
	}

	return nil
}

func init() {
	nonceCommand = &cobra.Command{
		Use:          "nonce",
		SilenceUsage: true,
	}
	fixGapCommand = &cobra.Command{
		Use:          "fix-gap",
		Short:        "Send txs to fix the nonce gap for a specific account",
		Long:         fixGapUsage,
		RunE:         accountNonceFixGap,
		SilenceUsage: true,
	}

	// Arguments for account
	inputAccountArgs.rpcURL = AccountCmd.PersistentFlags().StringP(ArgRpcURL, "r", "http://localhost:8545", "The RPC endpoint url")

	// Arguments for fix-gap
	inputAccountArgs.privateKey = fixGapCommand.PersistentFlags().String(ArgPrivateKey, "", "the private key to be used when sending the txs to fix the nonce gap")
	fatalIfError(fixGapCommand.MarkPersistentFlagRequired(ArgPrivateKey))

	// Top Level
	AccountCmd.AddCommand(nonceCommand)

	// Nonce
	nonceCommand.AddCommand(fixGapCommand)
}

// Wait for the transaction to be mined
func WaitMineTransaction(ctx context.Context, client *ethclient.Client, tx *types.Transaction, txTimeout uint64) error {
	timeout := time.NewTimer(time.Duration(txTimeout) * time.Second)
	defer timeout.Stop()
	for {
		select {
		case <-timeout.C:
			err := fmt.Errorf("timeout waiting for transaction to be mined")
			return err
		default:
			r, err := client.TransactionReceipt(ctx, tx.Hash())
			if err != nil {
				if !errors.Is(err, ethereum.NotFound) {
					log.Error().Err(err)
					return err
				}
				time.Sleep(1 * time.Second)
				continue
			}
			if r.Status != 0 {
				log.Info().Stringer("hash", r.TxHash).Msg("transaction successful")
				return nil
			} else if r.Status == 0 {
				log.Error().Stringer("hash", r.TxHash).Msg("transaction failed")
				return nil
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func fatalIfError(err error) {
	if err == nil {
		return
	}
	log.Fatal().Err(err).Msg("Unexpected error occurred")
}

func getPoolContent(addr common.Address) (poolTxs map[uint64]string, maxNonceFound uint64, err error) {
	var result PoolContent
	err = rpcClient.Client().Call(&result, "txpool_content")
	if err != nil {
		return
	}

	// get only txs from the address we are looking for
	poolTxs = make(map[uint64]string)
	txs, found := result.Queued[addr.String()]
	if !found {
		return
	}

	// iterate over the transactions and get the nonce
	for nonce, v := range txs {
		tx := v.(map[string]any)
		nonceInt, ok := new(big.Int).SetString(nonce, 10)
		if !ok {
			err = fmt.Errorf("invalid nonce found: %s", nonce)
			return
		}

		txHash := tx["hash"].(string)
		poolTxs[nonceInt.Uint64()] = txHash
		if nonceInt.Uint64() > maxNonceFound {
			maxNonceFound = nonceInt.Uint64()
		}
	}

	return
}

type PoolContent struct {
	Queued map[string]map[string]any
}
