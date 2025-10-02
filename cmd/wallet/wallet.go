package wallet

import (
	"encoding/json"
	"fmt"
	"os"

	_ "embed"

	"github.com/0xPolygon/polygon-cli/hdwallet"
	"github.com/spf13/cobra"
)

var (
	//go:embed usage.md
	usage                    string
	inputWords               *int
	inputLang                *string
	inputKDFIterations       *uint
	inputPassword            *string
	inputPasswordFile        *string
	inputMnemonic            *string
	inputMnemonicFile        *string
	inputPath                *string
	inputAddressesToGenerate *uint
	inputUseRawEntropy       *bool
	inputRootOnly            *bool
)

// WalletCmd represents the wallet command
var WalletCmd = &cobra.Command{
	Use:   "wallet [create|inspect]",
	Short: "Create or inspect BIP39(ish) wallets.",
	Long:  usage,
	RunE: func(cmd *cobra.Command, args []string) error {
		mode := args[0]
		var err error
		var mnemonic string
		if mode == "inspect" {
			// in the case of inspect, we'll partse a mnemonic and then continue
			mnemonic, err = getFileOrFlag(inputMnemonicFile, inputMnemonic)
			if err != nil {
				return err
			}
		} else {
			mnemonic, err = hdwallet.NewMnemonic(*inputWords, *inputLang)
			if err != nil {
				return err
			}
		}
		// mnemonic = "maid palace spring laptop shed when text taxi pupil movie athlete tag"
		// mnemonic = "crop cash unable insane eight faith inflict route frame loud box vibrant"
		// mnemonic = "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
		password, err := getFileOrFlag(inputPasswordFile, inputPassword)
		if err != nil {
			return err
		}
		pw, err := hdwallet.NewPolyWallet(mnemonic, password)
		if err != nil {
			return err
		}
		err = pw.SetPath(*inputPath)
		if err != nil {
			return err
		}
		err = pw.SetIterations(*inputKDFIterations)
		if err != nil {
			return err
		}
		err = pw.SetUseRawEntropy(*inputUseRawEntropy)
		if err != nil {
			return err
		}

		if *inputRootOnly {
			var key *hdwallet.PolyWalletExport
			key, err = pw.ExportRootAddress()
			if err != nil {
				return err
			}
			out, _ := json.MarshalIndent(key, " ", " ")
			fmt.Println(string(out))
			return nil
		}
		key, err := pw.ExportHDAddresses(int(*inputAddressesToGenerate))
		if err != nil {
			return err
		}
		// TODO support json vs txt out
		out, _ := json.MarshalIndent(key, " ", " ")
		fmt.Println(string(out))
		return nil
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expected exactly one argument: create or inspect")
		}
		if args[0] != "create" && args[0] != "inspect" {
			return fmt.Errorf("expected argument to be create or inspect. Got: %s", args[0])
		}
		return nil
	},
}

func getFileOrFlag(filename *string, flag *string) (string, error) {
	if filename == nil && flag == nil {
		return "", fmt.Errorf("both the filename and the flag pointers are nil")
	}
	if filename != nil && *filename != "" {
		filedata, err := os.ReadFile(*filename)
		if err != nil {
			return "", fmt.Errorf("could not open the specified file %s. Got error %s", *filename, err.Error())
		}
		return string(filedata), nil
	}
	if flag != nil {
		return *flag, nil
	}
	return "", fmt.Errorf("unable to determine flat or filename")
}

func init() {
	inputKDFIterations = WalletCmd.Flags().Uint("iterations", 2048, "number of pbkdf2 iterations to perform")
	inputWords = WalletCmd.Flags().Int("words", 24, "number of words to use in mnemonic")
	inputAddressesToGenerate = WalletCmd.Flags().Uint("addresses", 10, "number of addresses to generate")
	inputLang = WalletCmd.Flags().String("language", "english", "which language to use [ChineseSimplified, ChineseTraditional, Czech, English, French, Italian, Japanese, Korean, Spanish]")
	// https://github.com/satoshilabs/slips/blob/master/slip-0044.md
	// 0 - bitcoin
	// 60 - ether
	// 966 - matic
	inputPath = WalletCmd.Flags().String("path", "m/44'/60'/0'", "what would you like the derivation path to be")
	inputPassword = WalletCmd.Flags().String("password", "", "password used along with the mnemonic")
	inputPasswordFile = WalletCmd.Flags().String("password-file", "", "password stored in a file used along with the mnemonic")
	inputMnemonic = WalletCmd.Flags().String("mnemonic", "", "mnemonic phrase used to generate entropy")
	inputMnemonicFile = WalletCmd.Flags().String("mnemonic-file", "", "mnemonic phrase written in file used to generate entropy")
	inputUseRawEntropy = WalletCmd.Flags().Bool("raw-entropy", false, "substrate and polka dot don't follow strict bip39 and use raw entropy")
	inputRootOnly = WalletCmd.Flags().Bool("root-only", false, "don't produce HD accounts. Just produce a single wallet")
}
