package fixnoncegap

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/0xPolygon/polygon-cli/cmd/flag_loader"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var FixNonceGapCmd = &cobra.Command{
	Use:   "fix-nonce-gap",
	Short: "Send txs to fix the nonce gap for a specific account.",
	Long:  fixNonceGapUsage,
	Args:  cobra.NoArgs,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		rpcURL := flag_loader.GetRpcUrlFlagValue(cmd)
		inputFixNonceGapArgs.rpcURL = *rpcURL
		privateKey, err := flag_loader.GetRequiredPrivateKeyFlagValue(cmd)
		if err != nil {
			return err
		}
		inputFixNonceGapArgs.privateKey = *privateKey
		return nil
	},
	PreRunE:      prepareRpcClient,
	RunE:         fixNonceGap,
	SilenceUsage: true,
}

var (
	rpcClient *ethclient.Client
)

type fixNonceGapArgs struct {
	rpcURL     string
	privateKey string
	replace    bool
	maxNonce   uint64
}

var inputFixNonceGapArgs = fixNonceGapArgs{}

const (
	ArgPrivateKey = "private-key"
	ArgRpcURL     = "rpc-url"
	ArgReplace    = "replace"
	ArgMaxNonce   = "max-nonce"
)

//go:embed FixNonceGapUsage.md
var fixNonceGapUsage string

func prepareRpcClient(cmd *cobra.Command, args []string) error {
	var err error
	rpcURL := inputFixNonceGapArgs.rpcURL

	rpcClient, err = ethclient.Dial(rpcURL)
	if err != nil {
		log.Error().Err(err).Msgf("Unable to Dial RPC %s", rpcURL)
		return err
	}

	if _, err = rpcClient.BlockNumber(cmd.Context()); err != nil {
		log.Error().Err(err).Msg("Unable to get block number")
		return err
	}

	return nil
}

func fixNonceGap(cmd *cobra.Command, args []string) error {
	replace := inputFixNonceGapArgs.replace
	pvtKey := strings.TrimPrefix(inputFixNonceGapArgs.privateKey, "0x")
	pk, err := crypto.HexToECDSA(pvtKey)
	if err != nil {
		log.Error().Err(err).Msg("Invalid private key")
		return err
	}

	chainID, err := rpcClient.ChainID(cmd.Context())
	if err != nil {
		log.Error().Err(err).Msg("Cannot get chain ID")
		return err
	}

	opts, err := bind.NewKeyedTransactorWithChainID(pk, chainID)
	if err != nil {
		log.Error().Err(err).Msg("Cannot generate transactionOpts")
		return err
	}

	addr := opts.From

	currentNonce, err := rpcClient.NonceAt(cmd.Context(), addr, nil)
	if err != nil {
		log.Error().Err(err).Msg("Unable to get current nonce")
		return err
	}
	log.Info().Stringer("addr", addr).Msgf("Current nonce: %d", currentNonce)

	var maxNonce uint64
	if inputFixNonceGapArgs.maxNonce != 0 {
		maxNonce = inputFixNonceGapArgs.maxNonce
	} else {
		maxNonce, err = getMaxNonceFromTxPool(addr)
		if err != nil {
			if strings.Contains(err.Error(), "the method txpool_content does not exist/is not available") {
				log.Error().Err(err).Msg("The RPC doesn't provide access to txpool_content, please check --help for more information about --max-nonce")
				return nil
			}
			log.Error().Err(err).Msg("Unable to get max nonce from txpool")
			return err
		}
	}

	// check if there is a nonce gap
	if maxNonce == 0 || currentNonce >= maxNonce {
		log.Info().Stringer("addr", addr).Msg("There is no nonce gap.")
		return nil
	}
	log.Info().Stringer("addr", addr).Msgf("Nonce gap found. Max nonce: %d", maxNonce)

	gasPrice, err := rpcClient.SuggestGasPrice(cmd.Context())
	if err != nil {
		log.Error().Err(err).Msg("Unable to get suggested gas price")
		return err
	}

	to := &common.Address{}

	gas, err := rpcClient.EstimateGas(cmd.Context(), ethereum.CallMsg{
		From:     addr,
		To:       to,
		GasPrice: gasPrice,
		Value:    big.NewInt(1),
	})
	if err != nil {
		log.Error().Err(err).Msg("Unable to estimate gas")
		return err
	}

	txTemplate := &types.LegacyTx{
		To:       to,
		Gas:      gas,
		GasPrice: gasPrice,
		Value:    big.NewInt(1),
	}

	var lastTx *types.Transaction
	for i := currentNonce; i < maxNonce; i++ {
		txTemplate.Nonce = i
		tx := types.NewTx(txTemplate)
	out:
		for {
			signedTx, err := opts.Signer(opts.From, tx)
			if err != nil {
				log.Error().Err(err).Msg("Unable to sign tx")
				return err
			}
			log.Info().Stringer("hash", signedTx.Hash()).Msgf("sending tx with nonce %d", txTemplate.Nonce)

			err = rpcClient.SendTransaction(cmd.Context(), signedTx)
			if err != nil {
				if strings.Contains(err.Error(), "nonce too low") {
					log.Info().Stringer("hash", signedTx.Hash()).Msgf("another tx with nonce %d was mined while trying to increase the fee, skipping it", txTemplate.Nonce)
					break out
				} else if strings.Contains(err.Error(), "already known") {
					log.Info().Stringer("hash", signedTx.Hash()).Msgf("same tx with nonce %d already exists, skipping it", txTemplate.Nonce)
					break out
				} else if strings.Contains(err.Error(), "replacement transaction underpriced") ||
					strings.Contains(err.Error(), "INTERNAL_ERROR: could not replace existing tx") {
					if replace {
						txTemplateCopy := *txTemplate
						oldGasPrice := tx.GasPrice()
						// increase TX gas price by 10% and retry
						txTemplateCopy.GasPrice = new(big.Int).Mul(oldGasPrice, big.NewInt(11))
						txTemplateCopy.GasPrice = new(big.Int).Div(txTemplateCopy.GasPrice, big.NewInt(10))
						// if gas price didn't increase, this means the value is really small and 10% was smaller than 1.
						// This can be the case for local networks running for a long time without transactions.
						// We just add 1 wei to force gasPrice to move up to allow tx replacement
						if txTemplateCopy.GasPrice.Cmp(oldGasPrice) == 0 {
							txTemplateCopy.GasPrice = new(big.Int).Add(txTemplateCopy.GasPrice, big.NewInt(1))
						}
						tx = types.NewTx(&txTemplateCopy)
						log.Info().Stringer("hash", signedTx.Hash()).Msgf("tx with nonce %d is underpriced, increasing fee. From %d To %d", txTemplate.Nonce, oldGasPrice, txTemplateCopy.GasPrice)
						time.Sleep(time.Second)
						continue
					} else {
						log.Info().Stringer("hash", signedTx.Hash()).Msgf("another tx with nonce %d already exists, skipping it", txTemplate.Nonce)
						break out
					}
				}
				log.Error().Err(err).Msg("Unable to send tx")
				return err
			}

			// if we get here, just break the infinite loop and move to the next
			lastTx = signedTx
			break
		}
	}

	if lastTx != nil {
		log.Info().Stringer("hash", lastTx.Hash()).Msg("waiting for the last tx to get mined")
		err := WaitMineTransaction(cmd.Context(), rpcClient, lastTx, 600)
		if err != nil {
			log.Error().Err(err).Msg("Unable to wait for last tx to get mined")
			return err
		}
		log.Info().Stringer("addr", addr).Msg("Nonce gap fixed successfully")
		currentNonce, err = rpcClient.NonceAt(cmd.Context(), addr, nil)
		if err != nil {
			log.Error().Err(err).Msg("Unable to get current nonce")
			return err
		}
		log.Info().Stringer("addr", addr).Msgf("Current nonce: %d", currentNonce)
		return nil
	}

	return nil
}

func init() {
	f := FixNonceGapCmd.Flags()
	f.StringVarP(&inputFixNonceGapArgs.rpcURL, ArgRpcURL, "r", "http://localhost:8545", "the RPC endpoint URL")
	f.StringVar(&inputFixNonceGapArgs.privateKey, ArgPrivateKey, "", "private key to be used when sending txs to fix nonce gap")
	f.BoolVar(&inputFixNonceGapArgs.replace, ArgReplace, false, "replace the existing txs in the pool")
	f.Uint64Var(&inputFixNonceGapArgs.maxNonce, ArgMaxNonce, 0, "override max nonce value instead of getting it from the pool")
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

func getMaxNonceFromTxPool(addr common.Address) (uint64, error) {
	var result PoolContent
	err := rpcClient.Client().Call(&result, "txpool_content")
	if err != nil {
		return 0, err
	}

	txCollections := []PoolContentTxs{
		result.BaseFee,
		result.Pending,
		result.Queued,
	}

	maxNonceFound := uint64(0)
	for _, txCollection := range txCollections {
		// get only txs from the address we are looking for
		txs, found := txCollection[addr.String()]
		if !found {
			continue
		}

		// iterate over the transactions and get the nonce
		for nonce := range txs {
			nonceInt, ok := new(big.Int).SetString(nonce, 10)
			if !ok {
				err = fmt.Errorf("invalid nonce found: %s", nonce)
				return 0, err
			}

			if nonceInt.Uint64() > maxNonceFound {
				maxNonceFound = nonceInt.Uint64()
			}
		}
	}

	return maxNonceFound, nil
}

type PoolContent struct {
	BaseFee PoolContentTxs
	Pending PoolContentTxs
	Queued  PoolContentTxs
}

type PoolContentTxs map[string]map[string]any
