package parseethwallet

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"golang.org/x/crypto/sha3"
	"io"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/spf13/cobra"
)

var (
	inputFileName          *string
	inputPassword          *string
	inputRawHexPrivateKey  *string
	inputKeyStoreDirectory *string
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
This function can take a geth style wallet file and extract the private key as hex. It can also do the opposite

This command would take the private key and import it into a local keystore with no password

$ polycli parseethwallet --hexkey 42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa

$ cat UTC--2023-05-09T22-48-57.582848385Z--85da99c8a7c2c95964c8efd687e95e632fc533d6  | jq '.'
{
  "address": "85da99c8a7c2c95964c8efd687e95e632fc533d6",
  "crypto": {
    "cipher": "aes-128-ctr",
    "ciphertext": "d0b4377a4ae5ebc9a5bef06ce4be99565d10cb0dedc2f7ff5aaa07ea68e7b597",
    "cipherparams": {
      "iv": "8ecd172ff7ace15ed5bc44ea89473d8e"
    },
    "kdf": "scrypt",
    "kdfparams": {
      "dklen": 32,
      "n": 262144,
      "p": 1,
      "r": 8,
      "salt": "cd6ec772dc43225297412809feaae441d578642c6a67cabf4e29bcaf594f575b"
    },
    "mac": "c992128ed466ad15a9648f4112af22929b95f511f065b12a80abcfb7e4d39a79"
  },
  "id": "82af329d-2af5-41a6-ae6b-624f3e1c224b",
  "version": 3
}

If we wanted to go the opposite direction, we could run a command like this:

polycli parseethwallet --file /tmp/keystore/UTC--2023-05-09T22-48-57.582848385Z--85da99c8a7c2c95964c8efd687e95e632fc533d6  | jq '.'
{
  "Address": "0x85da99c8a7c2c95964c8efd687e95e632fc533d6",
  "PublicKey": "507cf9a75e053cda6922467721ddb10412da9bec30620347d9529cc77fca24334a4cf59685be4a2fdeabf4e7753350e42d2d3a20250fd9dc554d226463c8a3d5",
  "PrivateKey": "42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa"
}

`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// it would be nice to have a generic reader
		if *inputRawHexPrivateKey != "" {
			ks := keystore.NewKeyStore(*inputKeyStoreDirectory, keystore.StandardScryptN, keystore.StandardScryptP)
			pk, err := crypto.HexToECDSA(*inputRawHexPrivateKey)
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
