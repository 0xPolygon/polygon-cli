package parseethwallet

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"golang.org/x/crypto/sha3"
	"io"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/spf13/cobra"
)

var (
	inputFileName *string
	inputPassword *string
)

type plainKeyJSON struct {
	Address string              `json:"address"`
	Crypto  keystore.CryptoJSON `json:"crypto"`
}
type outKey struct {
	Address    string
	PublicKey  string
	PrivateKey string
}

var ParseETHWalletCmd = &cobra.Command{
	Use:   "parseethwallet --file UTC--2023-03-03T17-26-43.371893268Z--1652e7b47af367372a7a6d7d6fe5037702860c6d",
	Short: "A simple tool to extract the private key from an eth wallet",
	Long: `
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// it would be nice to have a generic reader

		rawData, err := getInputData(cmd, args)
		if err != nil {
			return err
		}
		k := new(plainKeyJSON)
		err = json.Unmarshal(rawData, &k)
		if err != nil {
			return err
		}
		d, err := keystore.DecryptDataV3(k.Crypto, *inputPassword)
		if err != nil {
			return err
		}
		ok := toOutputKey(d)
		outData, err := json.Marshal(ok)
		if err != nil {
			return err
		}
		fmt.Println(string(outData))
		return nil
	},
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	flagSet := ParseETHWalletCmd.PersistentFlags()
	inputFileName = flagSet.String("file", "", "Provide a file with the key information ")
	inputPassword = flagSet.String("password", "", "An optional password use to unlock the key")
}

func getInputData(cmd *cobra.Command, args []string) ([]byte, error) {
	if inputFileName != nil && *inputFileName != "" {
		return os.ReadFile(*inputFileName)
	}

	if len(args) > 1 {
		concat := strings.Join(args[1:], " ")
		return []byte(concat), nil
	}

	return io.ReadAll(os.Stdin)
}

func toOutputKey(key []byte) outKey {
	ok := outKey{}
	ok.PrivateKey = hex.EncodeToString(key)
	curve := secp256k1.S256()
	x1, y1 := curve.ScalarBaseMult(key)
	concat := append(x1.Bytes(), y1.Bytes()...)
	h := sha3.NewLegacyKeccak256()
	h.Write(concat)
	b := h.Sum(nil)
	ok.Address = fmt.Sprintf("0x%s", hex.EncodeToString(b[len(b)-20:]))
	ok.PublicKey = hex.EncodeToString(concat)
	return ok
}
