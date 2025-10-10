package fund

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/0xPolygon/polygon-cli/bindings/funder"
	"github.com/0xPolygon/polygon-cli/hdwallet"
	"github.com/0xPolygon/polygon-cli/util"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
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

	var addresses []common.Address
	var privateKeys []*ecdsa.PrivateKey

	if len(params.KeyFile) > 0 { // get addresses from key-file
		addresses, privateKeys, err = getAddressesAndKeysFromKeyFile(params.KeyFile)
	} else if len(params.Seed) > 0 { // get addresses from seed
		addresses, privateKeys, err = getAddressesAndKeysFromSeed(params.Seed, int(params.WalletsNumber))
	} else { // get addresses from private key
		addresses, privateKeys, err = getAddressesAndKeysFromPrivateKey(ctx, c)
	}
	// check errors after getting addresses
	if err != nil {
		return err
	}

	// Save private and public keys to a file if we have private keys.
	if len(privateKeys) > 0 {
		go func() {
			if err = saveToFile(params.OutputFile, privateKeys); err != nil {
				log.Error().Err(err).Msg("Unable to save keys to file")
				panic(err)
			}
			log.Info().Str("fileName", params.OutputFile).Msg("Wallets' address(es) and private key(s) saved to file")
		}()
	}

	// Fund wallets.
	log.Debug().Msg("checking if multicall3 is supported")
	var multicall3Addr *common.Address
	if len(params.Multicall3Address) > 0 {
		addr := common.HexToAddress(params.Multicall3Address)
		if addr == (common.Address{}) {
			log.Warn().
				Str("address", params.Multicall3Address).
				Msg("invalid multicall3 address provided, will try to detect or deploy multicall3")
			multicall3Addr = &addr
		}
	}

	multicall3Addr, _ = util.IsMulticall3Supported(ctx, c, true, tops, multicall3Addr)
	if multicall3Addr != nil {
		log.Info().
			Stringer("address", multicall3Addr).
			Msg("multicall3 is supported and will be used to fund wallets")
		return fundWalletsWithMulticall3(ctx, c, tops, addresses, multicall3Addr)
	} else {
		log.Info().Msg("multicall3 is not supported, will use funder contract to fund wallets")
		// Deploy or instantiate the Funder contract.
		var contract *funder.Funder
		contract, err = deployOrInstantiateFunderContract(ctx, c, tops, privateKey, len(addresses))
		if err != nil {
			return err
		}
		if err = fundWallets(ctx, c, tops, contract, addresses); err != nil {
			return err
		}
	}

	log.Info().Msg("Wallet(s) funded! 💸")

	log.Info().Msgf("Total execution time: %s", time.Since(startTime))
	return nil
}

func getAddressesAndKeysFromKeyFile(keyFilePath string) ([]common.Address, []*ecdsa.PrivateKey, error) {
	if len(keyFilePath) == 0 {
		return nil, nil, errors.New("the key file path is empty")
	}

	log.Trace().
		Str("keyFilePath", keyFilePath).
		Msg("getting addresses from key file")

	privateKeys, iErr := util.ReadPrivateKeysFromFile(keyFilePath)
	if iErr != nil {
		log.Error().
			Err(iErr).
			Msg("Unable to read private keys from key file")
		return nil, nil, fmt.Errorf("unable to read private keys from key file. %w", iErr)
	}
	addresses := make([]common.Address, len(privateKeys))
	for i, privateKey := range privateKeys {
		addresses[i] = util.GetAddress(context.Background(), privateKey)
		log.Trace().
			Interface("address", addresses[i]).
			Str("privateKey", hex.EncodeToString(privateKey.D.Bytes())).
			Msg("New wallet derived from key file")
	}
	log.Info().Int("count", len(addresses)).Msg("Wallet(s) derived from key file")
	return addresses, privateKeys, nil
}

func getAddressesAndKeysFromPrivateKey(ctx context.Context, c *ethclient.Client) ([]common.Address, []*ecdsa.PrivateKey, error) {
	// Derive or generate a set of wallets.
	var addresses []common.Address
	var privateKeys []*ecdsa.PrivateKey
	var err error
	if len(params.WalletAddresses) > 0 {
		log.Info().Msg("Using addresses provided by the user")
		addresses = make([]common.Address, len(params.WalletAddresses))
		for i, address := range params.WalletAddresses {
			addresses[i] = common.HexToAddress(address)
		}
		// No private keys available when using provided addresses
		privateKeys = nil
	} else if params.UseHDDerivation {
		log.Info().Msg("Deriving wallets from the default mnemonic")
		addresses, privateKeys, err = deriveHDWalletsWithKeys(int(params.WalletsNumber))
	} else {
		log.Info().Msg("Generating random wallets")
		addresses, privateKeys, err = generateWalletsWithKeys(int(params.WalletsNumber))
	}
	if err != nil {
		return nil, nil, err
	}
	return addresses, privateKeys, nil
}

// dialRpc dials the Ethereum RPC server and return an Ethereum client.
func dialRpc(ctx context.Context) (*ethclient.Client, error) {
	rpc, err := rpc.DialContext(ctx, params.RpcUrl)
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
	trimmedHexPrivateKey := strings.TrimPrefix(params.PrivateKey, "0x")
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
func deployOrInstantiateFunderContract(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts, privateKey *ecdsa.PrivateKey, numAddresses int) (*funder.Funder, error) {
	// Deploy the contract if no pre-deployed address flag is provided.
	var contractAddress common.Address
	var err error
	if params.FunderAddress == "" {
		// Deploy the Funder contract.
		// Note: `fundingAmountInWei` represents the amount the Funder contract will send to each newly generated wallets.
		fundingAmountInWei := params.FundingAmountInWei
		contractAddress, _, _, err = funder.DeployFunder(tops, c, fundingAmountInWei)
		if err != nil {
			log.Error().Err(err).Msg("Unable to deploy Funder contract")
			return nil, err
		}
		log.Debug().Interface("address", contractAddress).Msg("Funder contract deployed")

		// Fund the Funder contract.
		// Calculate the total amount needed to fund the contract based on the number of addresses.
		// Note: `funderContractBalanceInWei` represents the initial balance of the Funder contract.
		// The contract needs initial funds to be able to fund wallets.
		funderContractBalanceInWei := new(big.Int).Mul(fundingAmountInWei, big.NewInt(int64(numAddresses)))
		if err = util.SendTx(ctx, c, privateKey, &contractAddress, funderContractBalanceInWei, nil, uint64(30000)); err != nil {
			return nil, err
		}
	} else {
		// Use the pre-deployed address.
		contractAddress = common.HexToAddress(params.FunderAddress)
	}

	// Instantiate the contract.
	contract, err := funder.NewFunder(contractAddress, c)
	if err != nil {
		log.Error().Err(err).Msg("Unable to instantiate Funder contract")
		return nil, err
	}
	return contract, nil
}

// deriveHDWalletsWithKeys generates and exports a specified number of HD wallet addresses and their private keys.
func deriveHDWalletsWithKeys(n int) ([]common.Address, []*ecdsa.PrivateKey, error) {
	wallet, err := hdwallet.NewPolyWallet(defaultMnemonic, defaultPassword)
	if err != nil {
		return nil, nil, err
	}

	var derivedWallets *hdwallet.PolyWalletExport
	derivedWallets, err = wallet.ExportHDAddresses(n)
	if err != nil {
		return nil, nil, err
	}

	addresses := make([]common.Address, n)
	privateKeys := make([]*ecdsa.PrivateKey, n)
	for i, wallet := range derivedWallets.Addresses {
		addresses[i] = common.HexToAddress(wallet.ETHAddress)
		// Parse the private key
		trimmedHexPrivateKey := strings.TrimPrefix(wallet.HexPrivateKey, "0x")
		privateKey, err := crypto.HexToECDSA(trimmedHexPrivateKey)
		if err != nil {
			return nil, nil, fmt.Errorf("unable to parse private key for wallet %d: %w", i, err)
		}
		privateKeys[i] = privateKey
		log.Trace().Interface("address", addresses[i]).Str("privateKey", wallet.HexPrivateKey).Str("path", wallet.Path).Msg("New wallet derived")
	}
	log.Info().Int("count", n).Msg("Wallet(s) derived")
	return addresses, privateKeys, nil
}

// generateWalletsWithKeys generates a specified number of Ethereum wallets with random private keys.
// It returns a slice of common.Address representing the Ethereum addresses and their corresponding private keys.
func generateWalletsWithKeys(n int) ([]common.Address, []*ecdsa.PrivateKey, error) {
	// Generate private keys.
	privateKeys := make([]*ecdsa.PrivateKey, n)
	addresses := make([]common.Address, n)
	for i := 0; i < n; i++ {
		pk, err := crypto.GenerateKey()
		if err != nil {
			log.Error().Err(err).Msg("Error generating key")
			return nil, nil, err
		}
		privateKeys[i] = pk
		addresses[i] = crypto.PubkeyToAddress(pk.PublicKey)
		log.Trace().Interface("address", addresses[i]).Str("privateKey", hex.EncodeToString(pk.D.Bytes())).Msg("New wallet generated")
	}
	log.Info().Int("count", n).Msg("Wallet(s) generated")

	return addresses, privateKeys, nil
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

func fundWalletsWithMulticall3(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts, wallets []common.Address, multicall3Addr *common.Address) error {
	log.Debug().
		Msg("funding wallets with multicall3")

	const defaultAccsToFundPerTx = 400
	accsToFundPerTx, err := util.Multicall3MaxAccountsToFundPerTx(ctx, c)
	if err != nil {
		log.Warn().Err(err).
			Uint64("defaultAccsToFundPerTx", defaultAccsToFundPerTx).
			Msg("failed to get multicall3 max accounts to fund per tx, falling back to default")
		accsToFundPerTx = defaultAccsToFundPerTx
	}
	log.Debug().Uint64("accsToFundPerTx", accsToFundPerTx).Msg("multicall3 max accounts to fund per tx")
	chSize := (uint64(len(wallets)) / accsToFundPerTx) + 1

	txsCh := make(chan *types.Transaction, chSize)

	accs := []common.Address{}
	wg := sync.WaitGroup{}
	rl := rate.NewLimiter(rate.Limit(params.RateLimit), 1)
	if params.RateLimit <= 0.0 {
		rl = nil
	}
	for i := 0; i < len(wallets); i++ {
		wallet := wallets[i]
		// if account is the funding account, skip it
		if wallet == tops.From {
			continue
		}
		accs = append(accs, wallet)

		if uint64(len(accs)) == accsToFundPerTx || i == len(wallets)-1 {
			wg.Add(1)
			go func(tops *bind.TransactOpts, accs []common.Address) {
				defer wg.Done()
				if rl != nil {
					err := rl.Wait(ctx)
					if err != nil {
						log.Error().Err(err).Msg("rate limiter wait failed before funding accounts with multicall3")
						return
					}
				}
				tx, err := util.Multicall3FundAccountsWithNativeToken(c, tops, accs, params.FundingAmountInWei, multicall3Addr)
				if err != nil {
					log.Error().Err(err).Msg("failed to fund accounts with multicall3")
					return
				}
				log.Debug().
					Stringer("txHash", tx.Hash()).
					Int("done", i+1).
					Uint64("of", uint64(len(wallets))).
					Msg("multicall3 transaction to fund accounts sent")
				txsCh <- tx
			}(tops, accs)
			accs = []common.Address{}
		}
	}
	wg.Wait()
	close(txsCh)

	log.Info().Msg("all funding transactions sent, waiting for confirmation...")

	// ensure the txs to fund sending accounts using multicall3 were mined successfully
	for tx := range txsCh {
		if rl != nil {
			err := rl.Wait(ctx)
			if err != nil {
				return err
			}
		}

		r, err := util.WaitReceipt(ctx, c, tx.Hash())
		if err != nil {
			log.Error().Err(err).Msg("failed to wait for transaction to fund accounts with multicall3")
			return err
		}
		if r == nil || r.Status != types.ReceiptStatusSuccessful {
			errMsg := fmt.Sprintf("transaction to fund accounts with multicall3 failed, receipt is nil or status is not successful, txHash: %s", tx.Hash().String())
			log.Error().Msg(errMsg)
			return errors.New(errMsg)
		}
		log.Debug().
			Stringer("txHash", tx.Hash()).
			Msg("transaction to fund accounts confirmed")
	}

	return nil
}

func getAddressesAndKeysFromSeed(seed string, numWallets int) ([]common.Address, []*ecdsa.PrivateKey, error) {
	if len(seed) == 0 {
		return nil, nil, errors.New("the seed string is empty")
	}
	if numWallets <= 0 {
		return nil, nil, errors.New("number of wallets must be greater than 0")
	}

	log.Info().
		Str("seed", seed).
		Int("numWallets", numWallets).
		Msg("Generating wallets from seed")

	addresses := make([]common.Address, numWallets)
	privateKeys := make([]*ecdsa.PrivateKey, numWallets)

	for i := 0; i < numWallets; i++ {
		// Create a deterministic string by combining seed with index and current date
		// Format: seed_index_YYYYMMDD (e.g., "ephemeral_test_0_20241010")
		currentDate := time.Now().Format("20060102") // YYYYMMDD format
		seedWithIndex := fmt.Sprintf("%s_%d_%s", seed, i, currentDate)

		// Generate SHA256 hash of the seed+index+date
		hash := sha256.Sum256([]byte(seedWithIndex))
		hashHex := hex.EncodeToString(hash[:])

		// Create private key from hash
		privateKey, err := crypto.HexToECDSA(hashHex)
		if err != nil {
			return nil, nil, fmt.Errorf("unable to create private key from seed for wallet %d: %w", i, err)
		}

		privateKeys[i] = privateKey
		addresses[i] = crypto.PubkeyToAddress(privateKey.PublicKey)

		log.Trace().
			Interface("address", addresses[i]).
			Str("privateKey", hashHex).
			Str("seedWithIndex", seedWithIndex).
			Msg("New wallet generated from seed")
	}

	log.Info().Int("count", numWallets).Msg("Wallet(s) generated from seed")
	return addresses, privateKeys, nil
}
