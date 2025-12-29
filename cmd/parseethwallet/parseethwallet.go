package parseethwallet

import (
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/0xPolygon/polygon-cli/gethkeystore"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/sha3"
)

var (
	//go:embed usage.md
	usage                  string
	inputFileName          string
	inputPassword          string
	inputRawHexPrivateKey  string
	inputKeyStoreDirectory string
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
		if inputRawHexPrivateKey != "" {
			trimmedHexPrivateKey := strings.TrimPrefix(inputRawHexPrivateKey, "0x")
			ks := keystore.NewKeyStore(inputKeyStoreDirectory, keystore.StandardScryptN, keystore.StandardScryptP)
			pk, err := crypto.HexToECDSA(trimmedHexPrivateKey)
			if err != nil {
				return err
			}
			_, err = ks.ImportECDSA(pk, inputPassword)
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
		d, err := keystore.DecryptDataV3(k.Crypto, inputPassword)
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
	f := ParseETHWalletCmd.Flags()
	f.StringVar(&inputFileName, "file", "", "file with key information")
	f.StringVar(&inputPassword, "password", "", "optional password to unlock key")
	f.StringVar(&inputRawHexPrivateKey, "hexkey", "", "optional hexkey to use for generating geth style key")
	f.StringVar(&inputKeyStoreDirectory, "keystore", "/tmp/keystore", "directory where keys will be stored when importing raw hex")
}

func getInputData(cmd *cobra.Command, args []string) ([]byte, error) {
	if inputFileName != "" {
		return os.ReadFile(inputFileName)
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
