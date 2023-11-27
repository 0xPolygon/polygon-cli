package fund

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/maticnetwork/polygon-cli/bindings/funder"
	"github.com/maticnetwork/polygon-cli/util"
	"github.com/rs/zerolog/log"

	ethrpc "github.com/ethereum/go-ethereum/rpc"
)

// The amount in Wei to send to the Funder contract.
const funderContractFundingAmount = "1000000000000000000" // 10 ether.

var (
	// The current chain ID for transaction replay protection.
	chainID *big.Int

	// The ECDSA private key used to send transactions.
	privateKey *ecdsa.PrivateKey

	// The address of the wallet used to send transactions.
	address common.Address
)

// runFunding deploys or instantiates a `Funder` contract to bulk fund randomly generated wallets.
// Wallets' addresses and private keys are saved to a file.
func runFunding(ctx context.Context) error {
	log.Info().Msg("Starting bulk funding wallets")
	log.Trace().Interface("params", params).Msg("Input parameters")
	startTime := time.Now()

	// Set up the environment.
	c, err := dialRpc(ctx)
	if err != nil {
		return err
	}

	if err = initializeParams(ctx, c); err != nil {
		return err
	}

	var tops *bind.TransactOpts
	tops, err = bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Error().Err(err).Msg("Unable create transaction signer")
		return err
	}

	// Deploy or instantiate the Funder contract.
	var contractAddress common.Address
	var contract *funder.Funder
	contractAddress, contract, err = deployOrInstantiateFunderContract(ctx, c, tops, &bind.CallOpts{})
	if err != nil {
		return err
	}

	// Fund the Funder contract.
	amount, ok := new(big.Int).SetString(funderContractFundingAmount, 10)
	if !ok {
		return errors.New("unable to format the amount to send to the Funder contract")
	}
	if err = fundContract(ctx, c, contractAddress, amount); err != nil {
		return err
	}

	// Generate a set of wallets.
	var addresses []common.Address
	addresses, err = generateWallets(int(*params.WalletCount))
	if err != nil {
		return err
	}
	log.Debug().Interface("addresses", addresses).Msg("Address(es) of newly generated wallet(s)")

	// Fund wallets.
	if err = fundWallets(ctx, c, contract, addresses); err != nil {
		return err
	}
	log.Info().Msg("Wallet(s) funded! 💸")

	log.Info().Msgf("Total execution time: %s", time.Since(startTime))
	return nil
}

// dialRpc dials the Ethereum RPC server and return an Ethereum client.
func dialRpc(ctx context.Context) (*ethclient.Client, error) {
	rpc, err := ethrpc.DialContext(ctx, *params.RpcUrl)
	if err != nil {
		log.Error().Err(err).Msg("Unable to dial")
		return nil, err
	}
	defer rpc.Close()

	rpc.SetHeader("Accept-Encoding", "identity")
	return ethclient.NewClient(rpc), nil
}

// Initialize  parameters.
func initializeParams(ctx context.Context, c *ethclient.Client) error {
	// Get the chaind id.
	var err error
	chainID, err = c.ChainID(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Unable to fetch chain ID")
		return err
	}
	log.Trace().Uint64("chainID", chainID.Uint64()).Msg("Detected chain ID")

	// Parse the private key.
	privateKey, err = ethcrypto.HexToECDSA(strings.TrimPrefix(*params.PrivateKey, "0x"))
	if err != nil {
		log.Error().Err(err).Msg("Unable to process the private key")
		return err
	}
	address = ethcrypto.PubkeyToAddress(privateKey.PublicKey)
	return nil
}

// deployOrInstantiateFunderContract deploys or instantiates a Funder contract.
// If the pre-deployed address is specified, the contract will not be deployed.
func deployOrInstantiateFunderContract(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts, cops *bind.CallOpts) (common.Address, *funder.Funder, error) {
	// Format the funding amount.
	fundingAmount, err := util.HexToBigInt(*params.WalletFundingHexAmount)
	if err != nil {
		log.Error().Err(err).Msg("Unable to parse funding amount")
		return common.Address{}, nil, err
	}

	// Deploy the contract if no pre-deployed address flag is provided.
	var contractAddress common.Address
	if *params.FunderAddress == "" {
		// Deploy the Funder contract.
		contractAddress, _, _, err = funder.DeployFunder(tops, c, fundingAmount)
		if err != nil {
			log.Error().Err(err).Msg("Unable to deploy Funder contract")
			return common.Address{}, nil, err
		}
		log.Debug().Interface("address", contractAddress).Msg("Funder contract deployed")
	} else {
		// Use the pre-deployed address.
		contractAddress = common.HexToAddress(*params.FunderAddress)
	}

	// Instantiate the contract.
	contract, err := funder.NewFunder(contractAddress, c)
	if err != nil {
		log.Error().Err(err).Msg("Unable to instantiate Funder contract")
		return common.Address{}, nil, err
	}
	return contractAddress, contract, nil
}

func fundContract(ctx context.Context, c *ethclient.Client, contractAddress common.Address, amount *big.Int) error {
	// Get the nonce.
	nonce, err := c.PendingNonceAt(ctx, address)
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
	tx := ethtypes.NewTx(&ethtypes.LegacyTx{
		Nonce:    nonce,
		Gas:      uint64(30000),
		GasPrice: gasPrice,
		To:       &contractAddress,
		Value:    amount,
		Data:     nil,
	})
	var signedTx *ethtypes.Transaction
	signedTx, err = ethtypes.SignTx(tx, ethtypes.NewEIP155Signer(chainID), privateKey)
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

	// Get contract balance.
	var balance *big.Int
	balance, err = c.BalanceAt(ctx, contractAddress, nil)
	if err != nil {
		return err
	}
	log.Debug().Interface("balance", balance).Msg("Funder contract funded")
	return nil
}

// generateWallets generates a specified number of Ethereum wallets with random private keys.
// It returns a slice of common.Address representing the Ethereum addresses of the generated wallets.
func generateWallets(n int) ([]common.Address, error) {
	// Generate private keys.
	privateKeys := make([]*ecdsa.PrivateKey, n)
	addresses := make([]common.Address, n)
	for i := 0; i < n; i++ {
		pk, err := ethcrypto.GenerateKey()
		if err != nil {
			log.Error().Err(err).Msg("Error generating key")
			return nil, err
		}
		privateKeys[i] = pk
		addresses[i] = ethcrypto.PubkeyToAddress(pk.PublicKey)
	}

	// Save private and public keys to a file.
	go func() {
		if err := saveToFile(*params.OutputFile, privateKeys); err != nil {
			log.Error().Err(err).Msg("Unable to save keys to file")
			panic(err)
		}
		log.Info().Str("fileName", *params.OutputFile).Msg("Wallet address(es) and private key(s) saved to file")
	}()

	return addresses, nil
}

// saveToFile serializes wallet data into the specified JSON format and writes it to the designated file.
func saveToFile(fileName string, privateKeys []*ecdsa.PrivateKey) error {
	type wallet struct {
		Address    string `json:"Address"`
		PrivateKey string `json:"PrivateKey"`
	}

	// Populate the struct with addresses and private keys.
	data := make([]wallet, len(privateKeys))
	for i, privateKey := range privateKeys {
		address := ethcrypto.PubkeyToAddress(privateKey.PublicKey)
		data[i] = wallet{
			Address:    address.String(),
			PrivateKey: hex.EncodeToString(privateKey.D.Bytes()),
		}
	}

	// Save data to file.
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	if err = os.WriteFile(fileName, jsonData, 0644); err != nil {
		return err
	}
	return nil
}

// fundWallets funds multiple wallets using the provided Funder contract.
func fundWallets(ctx context.Context, c *ethclient.Client, contract *funder.Funder, wallets []common.Address) error {
	// Configure transaction options.
	tops, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Error().Err(err).Msg("Unable create transaction signer")
		return err
	}

	// Fund wallets.
	switch len(wallets) {
	case 0:
		return errors.New("no wallet to fund")
	case 1:
		// Fund a single account.
		if _, err = contract.Fund(tops, wallets[0]); err != nil {
			log.Error().Err(err).Msg("Unable to fund wallet")
			return err
		}
	default:
		// Fund multiple wallets in bulk.
		if _, err = contract.BulkFund(tops, wallets); err != nil {
			log.Error().Err(err).Msg("Unable to bulk fund wallets")
			return err
		}
	}
	return nil
}
