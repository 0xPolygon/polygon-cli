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
package mnemonic

import (
	_ "embed"
	"fmt"

	"github.com/maticnetwork/polygon-cli/hdwallet"
	"github.com/spf13/cobra"
)

var (
	//go:embed usage.md
	usage              string
	inputMnemonicWords *int
	inputMnemonicLang  *string
)

// mnemonicCmd represents the mnemonic command
var MnemonicCmd = &cobra.Command{
	Use:   "mnemonic",
	Short: "Generate a BIP39 mnemonic seed.",
	Long:  usage,
	RunE: func(cmd *cobra.Command, args []string) error {
		mnemonic, err := hdwallet.NewMnemonic(*inputMnemonicWords, *inputMnemonicLang)
		if err != nil {
			return err
		}
		cmd.Println(mnemonic)
		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if *inputMnemonicWords < 12 {
			return fmt.Errorf("the number of words in the mnemonic must be 12 or more. Given: %d", *inputMnemonicWords)
		}
		if *inputMnemonicWords > 24 {
			return fmt.Errorf("the number of words in the mnemonic must be 24 or less. Given: %d", *inputMnemonicWords)
		}
		if *inputMnemonicWords%3 != 0 {
			return fmt.Errorf("the number of words in the mnemonic must be a multiple of 3")
		}
		return nil

	},
}

func init() {
	inputMnemonicWords = MnemonicCmd.PersistentFlags().Int("words", 24, "The number of words to use in the mnemonic")
	inputMnemonicLang = MnemonicCmd.PersistentFlags().String("language", "english", "Which language to use [ChineseSimplified, ChineseTraditional, Czech, English, French, Italian, Japanese, Korean, Spanish]")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mnemonicCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mnemonicCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
