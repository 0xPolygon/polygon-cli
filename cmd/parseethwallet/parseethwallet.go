package parseethwallet

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/maticnetwork/polygon-cli/gethkeystore"
	"io"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"golang.org/x/crypto/sha3"

	_ "embed"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/spf13/cobra"
)

var (
	//go:embed usage.md
	usage                  string
	inputFileName          *string
	inputPassword          *string
	inputRawHexPrivateKey  *string
	inputKeyStoreDirectory *string
)

type outKey struct {
	Address    string
	PublicKey  string
	PrivateKey string
}

var ParseETHWalletCmd = &cobra.Command{
	Use:   "parseethwallet [flags]",
	Short: "Extract the private key from an eth wallet.",
	Long:  usage,
	RunE: func(cmd *cobra.Command, args []string) error {
		// it would be nice to have a generic reader
		if *inputRawHexPrivateKey != "" {
			trimmedHexPrivateKey := strings.TrimPrefix(*inputRawHexPrivateKey, "0x")
			ks := keystore.NewKeyStore(*inputKeyStoreDirectory, keystore.StandardScryptN, keystore.StandardScryptP)
			pk, err := crypto.HexToECDSA(trimmedHexPrivateKey)
			if err != nil {
				return err
			}
			_, err = ks.ImportECDSA(pk, *inputPassword)
			if err != nil {
				return err
			}
			return nil
		}

		rawData, err := getInputData(cmd, args)
		if err != nil {
			return err
		}
		k := new(gethkeystore.RawKeystoreData)
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
	inputRawHexPrivateKey = flagSet.String("hexkey", "", "An optional hexkey that would be use to generate a geth style key")
	inputKeyStoreDirectory = flagSet.String("keystore", "/tmp/keystore", "The directory where keys would be stored when importing a raw hex")
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
