/*
Copyright © 2022 Polygon <engineering@polygon.technology>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/maticnetwork/polygon-cli/hdwallet"
	"github.com/spf13/cobra"
)

var (
	inputCurve               *string
	inputKDFIterations       *uint
	inputPassword            *string
	inputPasswordFile        *string
	inputMnemonic            *string
	inputMnemonicFile        *string
	inputPath                *string
	inputAddressesToGenerate *uint
	inputUseRawEntropy       *bool
)

// walletCmd represents the wallet command
var walletCmd = &cobra.Command{
	Use:   "wallet [create|inspect]",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		mode := args[0]
		if mode == "inspect" {
			// in the case of inspect, we'll partse a mnemonic and then continue
			return fmt.Errorf("Not implemented yet")
		}

		mnemonic, err := hdwallet.NewMnemonic(*inputWords, *inputLang)
		if err != nil {
			return err
		}
		// TODO remove this once we implement inspect
		// mnemonic = "maid palace spring laptop shed when text taxi pupil movie athlete tag"
		mnemonic = "crop cash unable insane eight faith inflict route frame loud box vibrant"
		mnemonic = "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
		password, err := getFileOrFlag(inputPasswordFile, inputPassword)
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

		key, err := pw.ExportAddresses(int(*inputAddressesToGenerate))
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
			return fmt.Errorf("Expected exactly one argument: create or inspect.")
		}
		if args[0] != "create" && args[0] != "inspect" {
			return fmt.Errorf("Expected argument to be create or inspect. Got: %s", args[0])
		}
		return nil
	},
}

func getFileOrFlag(filename *string, flag *string) (string, error) {
	if filename == nil && flag == nil {
		return "", fmt.Errorf("Both the filename and the flag pointers are nil")
	}
	if filename != nil && *filename != "" {
		filedata, err := ioutil.ReadFile(*filename)
		if err != nil {
			return "", fmt.Errorf("Could not open the specified file %s. Got error %s", *filename, err.Error())
		}
		return string(filedata), nil
	}
	if flag != nil {
		return *flag, nil
	}
	return "", fmt.Errorf("Unable to determine flat or filename")
}

func init() {
	rootCmd.AddCommand(walletCmd)
	inputCurve = walletCmd.PersistentFlags().String("curve", "secp256k1", "ed25519, sr25519, or secp256k1")
	inputKDFIterations = walletCmd.PersistentFlags().Uint("iterations", 2048, "Number of pbkdf2 iterations to perform")
	inputWords = walletCmd.PersistentFlags().Int("words", 24, "The number of words to use in the mnemonic")
	inputAddressesToGenerate = walletCmd.PersistentFlags().Uint("addresses", 10, "The number of addresses to generate")
	inputLang = walletCmd.PersistentFlags().String("language", "english", "Which language to use [ChineseSimplified, ChineseTraditional, Czech, English, French, Italian, Japanese, Korean, Spanish]")
	// https://github.com/satoshilabs/slips/blob/master/slip-0044.md
	// 0 - bitcoin
	// 60 - ether
	// 966 - matic
	inputPath = walletCmd.PersistentFlags().String("path", "m/44'/60'/0'", "What would you like the derivation path to be")
	inputPassword = walletCmd.PersistentFlags().String("password", "", "Password used along with the mnemonic")
	inputPasswordFile = walletCmd.PersistentFlags().String("password-file", "", "Password stored in a file used along with the mnemonic")
	inputMnemonic = walletCmd.PersistentFlags().String("mnemonic", "", "A mnemonic phrase used to generate entropy")
	inputMnemonicFile = walletCmd.PersistentFlags().String("mnemonic-file", "", "A mneomonic phrase written in a file used to generate entropy")
	inputUseRawEntropy = walletCmd.PersistentFlags().Bool("raw-entropy", false, "substrate and polkda dot don't follow strict bip39 and use raw entropy")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// walletCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// walletCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
