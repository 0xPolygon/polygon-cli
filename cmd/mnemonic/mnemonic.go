package mnemonic

import (
	"fmt"

	"github.com/maticnetwork/polygon-cli/hdwallet"
	"github.com/spf13/cobra"
)

var (
	inputMnemonicWords *int
	inputMnemonicLang  *string
)

// mnemonicCmd represents the mnemonic command
var MnemonicCmd = &cobra.Command{
	Use:   "mnemonic",
	Short: "Generate a BIP39 mnemonic seed.",
	Long:  "",
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
}
