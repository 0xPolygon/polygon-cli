package fund

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
	"sync"
	"time"

	_ "embed"

	"github.com/chenzhijie/go-web3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	//go:embed usage.md
	usage                  string
	walletCount            int
	fundingWalletPK        string
	fundingWalletPublicKey *ecdsa.PublicKey
	chainID                int
	chainRPC               string
	concurrencyLevel       int
	walletFundingAmt       float64
	walletFundingGas       uint64
	verbosityEnabled       bool
	nonceMutex             sync.Mutex
	globalNonce            uint64
	nonceInitialized       bool
	outputFileFlag         string
)

// Wallet struct to hold public key, private key, and address
type Wallet struct {
	PublicKey  *ecdsa.PublicKey
	PrivateKey *ecdsa.PrivateKey
	Address    common.Address
}

func generateNonce(web3Client *web3.Web3) (uint64, error) {
	nonceMutex.Lock()
	defer nonceMutex.Unlock()

	if nonceInitialized == true {
		globalNonce++
	} else {
		// Derive the public key from the funding wallet's private key
		fundingWalletECDSA, ecdsaErr := crypto.HexToECDSA(fundingWalletPK)
		if ecdsaErr != nil {
			log.Error().Err(ecdsaErr).Msg("Error getting ECDSA from funding wallet private key")
			return 0, ecdsaErr
		}

		fundingWalletPublicKey = &fundingWalletECDSA.PublicKey
		// Convert ecdsa.PublicKey to common.Address
		fundingAddress := crypto.PubkeyToAddress(*fundingWalletPublicKey)

		nonce, err := web3Client.Eth.GetNonce(fundingAddress, nil)
		if err != nil {
			log.Error().Err(err).Msg("Error getting nonce")
			return 0, err
		}
		globalNonce = nonce
		nonceInitialized = true
	}

	return globalNonce, nil
}

func generateWallets(numWallets int) ([]Wallet, error) {
	wallets := make([]Wallet, 0, numWallets)

	for i := 0; i < numWallets; i++ {
		account, err := crypto.GenerateKey()
		if err != nil {
			log.Error().Err(err).Msg("Error generating key")
			return nil, err
		}

		addr := crypto.PubkeyToAddress(account.PublicKey)
		wallet := Wallet{
			PublicKey:  &account.PublicKey,
			PrivateKey: account,
			Address:    addr,
		}

		wallets = append(wallets, wallet)
	}
	return wallets, nil
}

func fundWallets(web3Client *web3.Web3, wallets []Wallet, amountWei *big.Int, walletFundingGas uint64, concurrency int) error {
	// Create a channel to control concurrency
	walletChan := make(chan Wallet, len(wallets))
	for _, wallet := range wallets {
		walletChan <- wallet
	}
	close(walletChan)

	// Wait group to ensure all goroutines finish before returning
	var wg sync.WaitGroup
	wg.Add(concurrency)

	// Function to fund wallets
	fundWallet := func() {
		defer wg.Done()
		for wallet := range walletChan {
			nonce, err := generateNonce(web3Client)
			if err != nil {
				log.Error().Err(err).Msg("Error getting nonce")
				return
			}

			// Fund the wallet using the obtained nonce
			_, err = web3Client.Eth.SyncSendRawTransaction(
				wallet.Address,
				amountWei,
				nonce,
				walletFundingGas,
				web3Client.Utils.ToGWei(1),
				nil,
			)
			if err != nil {
				log.Error().Err(err).Str("wallet", wallet.Address.Hex()).Msg("Error funding wallet")
				return
			}

			log.Info().Str("wallet", wallet.Address.Hex()).Msgf("Funded with %s wei", amountWei.String())
		}
	}

	// Start funding the wallets concurrently
	for i := 0; i < concurrency; i++ {
		go fundWallet()
	}

	// Wait for all goroutines to finish
	wg.Wait()
	return nil
}

// fundCmd represents the fund command
var FundCmd = &cobra.Command{
	Use:   "fund",
	Short: "Bulk fund many crypto wallets automatically.",
	Long:  usage,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runFunding(cmd); err != nil {
			log.Error().Err(err).Msg("Error funding wallets")
		}
	},
}

func runFunding(cmd *cobra.Command) error {
	// Capture the start time
	startTime := time.Now()

	// Remove '0x' prefix from fundingWalletPK if present
	if strings.HasPrefix(fundingWalletPK, "0x") {
		fundingWalletPK = fundingWalletPK[2:]
	}

	// setup new web3 session with remote rpc node
	web3Client, clientErr := web3.NewWeb3(chainRPC)
	if clientErr != nil {
		cmd.PrintErrf("There was an error creating web3 client: %s", clientErr.Error())
		return clientErr
	}

	// add pk to session for sending signed transactions
	if setAcctErr := web3Client.Eth.SetAccount(fundingWalletPK); setAcctErr != nil {
		cmd.PrintErrf("There was an error setting account with pk: %s", setAcctErr.Error())
		return setAcctErr
	}

	// set proper chainId for corresponding chainRPC
	cdkChainId := int64(chainID)
	web3Client.Eth.SetChainId(cdkChainId)

	// generate set of new wallet objects
	wallets, genWalletErr := generateWallets(walletCount)
	if genWalletErr != nil {
		cmd.PrintErrf("There was an error generating wallet objects: %s", genWalletErr.Error())
		return genWalletErr
	}

	// fund all crypto wallets
	log.Info().Msg("Starting to fund loadtest wallets...")
	fundWalletErr := fundWallets(web3Client, wallets, big.NewInt(int64(walletFundingAmt*1e18)), uint64(walletFundingGas), concurrencyLevel)
	if fundWalletErr != nil {
		log.Error().Err(fundWalletErr).Msg("Error funding wallets")
		return fundWalletErr
	}

	// Save wallet details to a file
	outputFile := outputFileFlag // You can modify the file format or name as needed

	type WalletDetails struct {
		Address    string `json:"Address"`
		PrivateKey string `json:"PrivateKey"`
	}

	walletDetails := make([]WalletDetails, len(wallets))
	for i, w := range wallets {
		privateKey := hex.EncodeToString(w.PrivateKey.D.Bytes()) // Convert private key to hex
		walletDetails[i] = WalletDetails{
			Address:    w.Address.Hex(),
			PrivateKey: privateKey,
		}
	}

	// Convert walletDetails to JSON
	walletsJSON, jsonErr := json.MarshalIndent(walletDetails, "", "  ")
	if jsonErr != nil {
		log.Error().Err(jsonErr).Msg("Error converting wallet details to JSON")
		return jsonErr
	}

	// Write JSON data to a file
	writeErr := ioutil.WriteFile(outputFile, walletsJSON, 0644)
	if writeErr != nil {
		log.Error().Err(writeErr).Msg("Error writing wallet details to file")
		return writeErr
	}

	log.Info().Msgf("Wallet details have been saved to %s", outputFile)

	// Calculate the duration
	duration := time.Since(startTime)
	log.Info().Msgf("Total execution time: %s", duration)

	return nil
}

func init() {
	// Configure zerolog to output to os.Stdout
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	FundCmd.Flags().IntVar(&walletCount, "wallet-count", 2, "Number of wallets to fund")
	FundCmd.Flags().StringVar(&fundingWalletPK, "funding-wallet-pk", "", "Corresponding private key for funding wallet address, ensure you remove leading 0x")
	FundCmd.Flags().IntVar(&chainID, "chain-id", 0, "The chain id for the transactions.")
	FundCmd.Flags().StringVar(&chainRPC, "rpc-url", "http://localhost:8545", "The RPC endpoint url")
	FundCmd.Flags().IntVar(&concurrencyLevel, "concurrency", 5, "Concurrency level for speeding up funding wallets")
	FundCmd.Flags().Float64Var(&walletFundingAmt, "wallet-funding-amt", 0.05, "Amount to fund each wallet with")
	FundCmd.Flags().Uint64Var(&walletFundingGas, "wallet-funding-gas", 50000, "Gas for each wallet funding transaction")
	FundCmd.Flags().BoolVar(&verbosityEnabled, "verbosity", true, "Global verbosity flag (true/false)")
	FundCmd.Flags().StringVar(&outputFileFlag, "output-file", "wallets.csv", "Specify the output CSV file name")

	// Set the global log level based on a verbosity flag
	if verbosityEnabled {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
