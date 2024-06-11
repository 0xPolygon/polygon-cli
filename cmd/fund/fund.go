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
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/maticnetwork/polygon-cli/bindings/funder"
	"github.com/maticnetwork/polygon-cli/hdwallet"
	"github.com/maticnetwork/polygon-cli/util"
	"github.com/rs/zerolog/log"
)

// The initial balance of Ether to send to the Funder contract.
const funderContractBalanceInEth = 1_000.0

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

	var privateKey *ecdsa.PrivateKey
	var chainID *big.Int
	privateKey, chainID, err = initializeParams(ctx, c)
	if err != nil {
		return err
	}

	var tops *bind.TransactOpts
	tops, err = bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Error().Err(err).Msg("Unable create transaction signer")
		return err
	}

	// Deploy or instantiate the Funder contract.
	var contract *funder.Funder
	contract, err = deployOrInstantiateFunderContract(ctx, c, tops, privateKey)
	if err != nil {
		return err
	}

	// Derive or generate a set of wallets.
	var addresses []common.Address
	if params.WalletAddresses != nil && *params.WalletAddresses != nil {
		log.Info().Msg("Using addresses provided by the user")
		addresses = make([]common.Address, len(*params.WalletAddresses))
		for i, address := range *params.WalletAddresses {
			addresses[i] = common.HexToAddress(address)
		}
	} else if *params.UseHDDerivation {
		log.Info().Msg("Deriving wallets from the default mnemonic")
		addresses, err = deriveHDWallets(int(*params.WalletsNumber))
	} else {
		log.Info().Msg("Generating random wallets")
		addresses, err = generateWallets(int(*params.WalletsNumber))
	}
	if err != nil {
		return err
	}

	// Fund wallets.
	if err = fundWallets(ctx, c, tops, contract, addresses); err != nil {
		return err
	}
	log.Info().Msg("Wallet(s) funded! 💸")

	log.Info().Msgf("Total execution time: %s", time.Since(startTime))
	return nil
}

// dialRpc dials the Ethereum RPC server and return an Ethereum client.
func dialRpc(ctx context.Context) (*ethclient.Client, error) {
	rpc, err := rpc.DialContext(ctx, *params.RpcUrl)
	if err != nil {
		log.Error().Err(err).Msg("Unable to dial")
		return nil, err
	}
	defer rpc.Close()

	rpc.SetHeader("Accept-Encoding", "identity")
	return ethclient.NewClient(rpc), nil
}

// Initialize  parameters.
func initializeParams(ctx context.Context, c *ethclient.Client) (*ecdsa.PrivateKey, *big.Int, error) {
	// Parse the private key.
	trimmedHexPrivateKey := strings.TrimPrefix(*params.PrivateKey, "0x")
	privateKey, err := crypto.HexToECDSA(trimmedHexPrivateKey)
	if err != nil {
		log.Error().Err(err).Msg("Unable to process the private key")
		return nil, nil, err
	}

	// Get the chaind id.
	var chainID *big.Int
	chainID, err = c.ChainID(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Unable to fetch chain ID")
		return nil, nil, err
	}
	log.Trace().Uint64("chainID", chainID.Uint64()).Msg("Detected chain ID")
	return privateKey, chainID, nil
}

// deployOrInstantiateFunderContract deploys or instantiates a Funder contract.
// If the pre-deployed address is specified, the contract will not be deployed.
func deployOrInstantiateFunderContract(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts, privateKey *ecdsa.PrivateKey) (*funder.Funder, error) {
	// Deploy the contract if no pre-deployed address flag is provided.
	var contractAddress common.Address
	var err error
	if *params.FunderAddress == "" {
		// Deploy the Funder contract.
		// Note: `fundingAmountInWei` reprensents the amount the Funder contract will send to each newly generated wallets.
		fundingAmountInWei := util.EthToWei(*params.FundingAmountInEth)
		contractAddress, _, _, err = funder.DeployFunder(tops, c, fundingAmountInWei)
		if err != nil {
			log.Error().Err(err).Msg("Unable to deploy Funder contract")
			return nil, err
		}
		log.Debug().Interface("address", contractAddress).Msg("Funder contract deployed")

		// Fund the Funder contract.
		// Note: `funderContractBalanceInWei` reprensents the initial balance of the Funder contract.
		// The contract needs initial funds to be able to fund wallets.
		funderContractBalanceInWei := util.EthToWei(funderContractBalanceInEth)
		if err = util.SendTx(ctx, c, privateKey, &contractAddress, funderContractBalanceInWei, nil, uint64(30000)); err != nil {
			return nil, err
		}
	} else {
		// Use the pre-deployed address.
		contractAddress = common.HexToAddress(*params.FunderAddress)
	}

	// Instantiate the contract.
	contract, err := funder.NewFunder(contractAddress, c)
	if err != nil {
		log.Error().Err(err).Msg("Unable to instantiate Funder contract")
		return nil, err
	}
	return contract, nil
}

// deriveWallets generates and exports a specified number of HD wallet addresses.
func deriveHDWallets(n int) ([]common.Address, error) {
	wallet, err := hdwallet.NewPolyWallet(defaultMnemonic, defaultPassword)
	if err != nil {
		return nil, err
	}

	var derivedWallets *hdwallet.PolyWalletExport
	derivedWallets, err = wallet.ExportHDAddresses(n)
	if err != nil {
		return nil, err
	}

	addresses := make([]common.Address, n)
	for i, wallet := range derivedWallets.Addresses {
		addresses[i] = common.HexToAddress(wallet.ETHAddress)
		log.Trace().Interface("address", addresses[i]).Str("privateKey", wallet.HexPrivateKey).Str("path", wallet.Path).Msg("New wallet derived")
	}
	log.Info().Int("count", n).Msg("Wallet(s) derived")
	return addresses, nil
}

// generateWallets generates a specified number of Ethereum wallets with random private keys.
// It returns a slice of common.Address representing the Ethereum addresses of the generated wallets.
func generateWallets(n int) ([]common.Address, error) {
	// Generate private keys.
	privateKeys := make([]*ecdsa.PrivateKey, n)
	addresses := make([]common.Address, n)
	for i := 0; i < n; i++ {
		pk, err := crypto.GenerateKey()
		if err != nil {
			log.Error().Err(err).Msg("Error generating key")
			return nil, err
		}
		privateKeys[i] = pk
		addresses[i] = crypto.PubkeyToAddress(pk.PublicKey)
		log.Trace().Interface("address", addresses[i]).Str("privateKey", hex.EncodeToString(pk.D.Bytes())).Msg("New wallet generated")
	}
	log.Info().Int("count", n).Msg("Wallet(s) generated")

	// Save private and public keys to a file.
	go func() {
		if err := saveToFile(*params.OutputFile, privateKeys); err != nil {
			log.Error().Err(err).Msg("Unable to save keys to file")
			panic(err)
		}
		log.Info().Str("fileName", *params.OutputFile).Msg("Wallets' address(es) and private key(s) saved to file")
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
		address := crypto.PubkeyToAddress(privateKey.PublicKey)
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
func fundWallets(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts, contract *funder.Funder, wallets []common.Address) error {
	// Fund wallets.
	switch len(wallets) {
	case 0:
		return errors.New("no wallet to fund")
	case 1:
		// Fund a single account.
		if _, err := contract.Fund(tops, wallets[0]); err != nil {
			log.Error().Err(err).Msg("Unable to fund wallet")
			return err
		}
	default:
		// Fund multiple wallets in bulk.
		if _, err := contract.BulkFund(tops, wallets); err != nil {
			log.Error().Err(err).Msg("Unable to bulk fund wallets")
			return err
		}
	}
	return nil
}
