package hdwallet

import (
	"fmt"
	"strings"

	"github.com/tyler-smith/go-bip39"
	"github.com/tyler-smith/go-bip39/wordlists"
)

var (
	wordsToBits    = map[int]int{12: 128, 15: 160, 18: 192, 21: 224, 24: 256}
	langToWordlist = map[string][]string{
		"chinesesimplified":  wordlists.ChineseSimplified,
		"chinesetraditional": wordlists.ChineseTraditional,
		"czech":              wordlists.Czech,
		"english":            wordlists.English,
		"french":             wordlists.French,
		"italian":            wordlists.Italian,
		"japanese":           wordlists.Japanese,
		"korean":             wordlists.Korean,
		"spanish":            wordlists.Spanish,
	}
)

func NewMnemonic(wordCount int, lang string) (string, error) {
	bits, hasKey := wordsToBits[wordCount]
	if !hasKey {
		return "", fmt.Errorf("The word count needs to be 12, 15, 18, 21, or 24. Got %d", wordCount)
	}
	wordList, hasKey := langToWordlist[strings.ToLower(lang)]
	if !hasKey {
		return "", fmt.Errorf("The language %s is not recognized.", lang)
	}

	bip39.SetWordList(wordList)

	entropy, err := bip39.NewEntropy(bits)
	if err != nil {
		return "", fmt.Errorf("There was an error getting entropy: %s", err.Error())
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", fmt.Errorf("There was an error creating the mnemonic: %s", err.Error())
	}

	return mnemonic, nil
}
