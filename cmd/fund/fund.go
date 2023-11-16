package fund

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"time"

	_ "embed"

	"github.com/chenzhijie/go-web3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
)

var (
	//go:embed usage.md
	usage                string
	walletCount          int
	fundingWalletAddress string
	fundingWalletPK      string
	chainID              int
	chainRPC             string
	walletFundingAmt     float64
	walletFundingGas     uint64
)

func generateWalletAddresses(numWallets int) ([]*common.Address, error) {
	var addresses []*common.Address

	for i := 0; i < numWallets; i++ {
		account, err := crypto.GenerateKey()
		if err != nil {
			log.Fatal(err)
		}
		addr := crypto.PubkeyToAddress(account.PublicKey)
		addresses = append(addresses, &addr)
	}
	return addresses, nil
}

func fundWallets(web3Client *web3.Web3, wallets []*common.Address, senderAddress common.Address, senderPrivateKey *ecdsa.PrivateKey, amountWei *big.Int, walletFundingGas uint64) error {
	nonce, err := web3Client.Eth.GetNonce(senderAddress, nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, wallet := range wallets {
		_, err := web3Client.Eth.SyncSendRawTransaction(
			common.HexToAddress((*wallet).Hex()),
			amountWei,
			nonce,
			walletFundingGas,
			web3Client.Utils.ToGWei(1),
			nil,
		)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Funded", (*wallet).Hex(), "with", amountWei, "wei")
		nonce++
	}
	return nil
}

// fundCmd represents the fund command
var FundCmd = &cobra.Command{
	Use:   fmt.Sprintf("fund"),
	Short: "Bulk fund many crypto wallets automatically.",
	Long:  usage,
	Run: func(cmd *cobra.Command, args []string) {
		runFunding(cmd)
	},
}

func runFunding(cmd *cobra.Command) error {
	// setup new web3 session with remote rpc node
	web3Client, clientErr := web3.NewWeb3(chainRPC)
	if clientErr != nil {
		cmd.PrintErrf("There was an error creating web3 client: %s", clientErr.Error())
		return clientErr
	}

	// add pk to session for sending signed transactions
	web3Client.Eth.SetAccount(fundingWalletPK)
	if setAcctErr := web3Client.Eth.SetAccount(fundingWalletPK); setAcctErr != nil {
		cmd.PrintErrf("There was an error setting account with pk: %s", setAcctErr.Error())
		return setAcctErr
	}

	// set proper chainId for corresponding chainRPC
	cdkChainId := int64(chainID)
	web3Client.Eth.SetChainId(cdkChainId)

	// convert wallet address and pk format for downstream processing
	fundingWalletAddressParsed := common.HexToAddress(fundingWalletAddress)
	fundingWalletECDSA, ecdsaErr := crypto.HexToECDSA(fundingWalletPK)
	if ecdsaErr != nil {
		cmd.PrintErrf("There was an error getting ECDSA: %s", ecdsaErr.Error())
		return ecdsaErr
	}

	// generate set of new wallet addresses
	addresses, genWalletErr := generateWalletAddresses(walletCount)
	if genWalletErr != nil {
		cmd.PrintErrf("There was an error generating wallet addresses: %s", genWalletErr.Error())
		return genWalletErr
	}

	// fund all crypto wallets
	fmt.Println("Funding all loadtest wallets...")
	fundWalletErr := fundWallets(web3Client, addresses, fundingWalletAddressParsed, fundingWalletECDSA, big.NewInt(int64(walletFundingAmt*1e18)), uint64(walletFundingGas))
	if fundWalletErr != nil {
		cmd.PrintErrf("There was an error funding wallets: %s", fundWalletErr.Error())
		return fundWalletErr
	}
	// small pause to let all funds land and state propogate across network
	time.Sleep(10 * time.Second)
	return nil
}

func init() {
	FundCmd.Flags().IntVar(&walletCount, "wallet-count", 2, "Number of wallets to fund")
	FundCmd.Flags().StringVar(&fundingWalletAddress, "funding-wallet-address", "", "Origin wallet that will be doing the funding")
	FundCmd.Flags().StringVar(&fundingWalletPK, "funding-wallet-pk", "", "Corresponding private key for funding wallet address, ensure you remove leading 0x")
	FundCmd.Flags().IntVar(&chainID, "chain-id", 100, "Blockchain network chain id")
	FundCmd.Flags().StringVar(&chainRPC, "chain-rpc", "", "Blockchain RPC node endpoint for sending funding transactions")
	FundCmd.Flags().Float64Var(&walletFundingAmt, "wallet-funding-amt", 0.05, "Amount to fund each wallet with")
	FundCmd.Flags().Uint64Var(&walletFundingGas, "wallet-funding-gas", 50000, "Gas for each wallet funding transaction")
}
