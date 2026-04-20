package wallet

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/gethkeystore"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
)

// newImportCmd builds `wallet import`: bring an existing key into
// the polycli/cast-shared keystore directory.
//
// Accepts three mutually-exclusive sources:
//   --private-key HEX             : 32-byte secp256k1 private key
//   --keystore-file PATH          : a v3 JSON keystore (asks for its
//                                   password, then re-encrypts under
//                                   the new password)
//   --mnemonic MNEMONIC (with optional --path / --index)
func newImportCmd() *cobra.Command {
	var (
		shared       keystoreSharedFlags
		privateKey   string
		sourceFile   string
		sourcePwFile string
		mnemonic     string
		mnemonicFile string
		bipPass      string
		path         string
		index        uint32
	)
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import an existing key into the keystore.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			sources := 0
			for _, s := range []string{privateKey, sourceFile, mnemonic, mnemonicFile} {
				if s != "" {
					sources++
				}
			}
			if sources == 0 {
				return &client.UsageError{Msg: "one of --private-key, --keystore-file, --mnemonic, or --mnemonic-file is required"}
			}
			if sources > 1 {
				return &client.UsageError{Msg: "--private-key, --keystore-file, --mnemonic, and --mnemonic-file are mutually exclusive"}
			}

			dir, err := resolveKeystoreDir(shared.KeystoreDir)
			if err != nil {
				return err
			}
			ks := newKeyStore(dir)

			switch {
			case privateKey != "":
				priv, err := parsePrivateKeyHex(privateKey)
				if err != nil {
					return err
				}
				newPw, err := readPassword(&shared, os.Stdin, true, "keystore password")
				if err != nil {
					return err
				}
				acc, err := ks.ImportECDSA(priv, newPw)
				if err != nil {
					return fmt.Errorf("importing private key: %w", err)
				}
				return printImport(cmd, acc.Address.Hex(), acc.URL.Path)

			case sourceFile != "":
				data, err := os.ReadFile(sourceFile)
				if err != nil {
					return fmt.Errorf("reading source keystore %s: %w", sourceFile, err)
				}
				sourcePw, err := readSourcePassword(sourcePwFile, "source keystore password")
				if err != nil {
					return err
				}
				priv, err := gethkeystore.DecryptKeystoreFile(data, sourcePw)
				if err != nil {
					return fmt.Errorf("decrypting source keystore: %w", err)
				}
				newPw, err := readPassword(&shared, os.Stdin, true, "new keystore password")
				if err != nil {
					return err
				}
				acc, err := ks.ImportECDSA(priv, newPw)
				if err != nil {
					return fmt.Errorf("re-importing key: %w", err)
				}
				return printImport(cmd, acc.Address.Hex(), acc.URL.Path)

			default:
				// Mnemonic-based import.
				if mnemonic == "" && mnemonicFile != "" {
					raw, err := os.ReadFile(mnemonicFile)
					if err != nil {
						return fmt.Errorf("reading mnemonic file %s: %w", mnemonicFile, err)
					}
					mnemonic = trimTrailingNewline(string(raw))
				}
				priv, finalPath, _, err := deriveFromMnemonic(mnemonic, bipPass, path, index)
				if err != nil {
					return err
				}
				newPw, err := readPassword(&shared, os.Stdin, true, "keystore password")
				if err != nil {
					return err
				}
				acc, err := ks.ImportECDSA(priv, newPw)
				if err != nil {
					return fmt.Errorf("importing derived key: %w", err)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "path     %s\n", finalPath)
				return printImport(cmd, acc.Address.Hex(), acc.URL.Path)
			}
		},
	}
	bindKeystoreFlags(cmd, &shared)
	f := cmd.Flags()
	f.StringVar(&privateKey, "private-key", "", "hex-encoded secp256k1 private key")
	f.StringVar(&sourceFile, "source-keystore-file", "", "path to an existing v3 JSON keystore to import")
	f.StringVar(&sourcePwFile, "source-password-file", "", "file with the existing keystore's password")
	f.StringVar(&mnemonic, "mnemonic", "", "BIP-39 mnemonic")
	f.StringVar(&mnemonicFile, "mnemonic-file", "", "file containing a BIP-39 mnemonic")
	f.StringVar(&bipPass, "bip39-passphrase", "", "optional BIP-39 passphrase")
	f.StringVar(&path, "path", "", "derivation path (default m/44'/60'/0'/0/<index>)")
	f.Uint32Var(&index, "index", 0, "address index when --path is not set")
	return cmd
}

// readSourcePassword reads the source-keystore password from a file
// or, absent a file, from stdin. Only asked once — the operator
// already typed it to create the source keystore.
func readSourcePassword(pwFile, label string) (string, error) {
	if pwFile != "" {
		raw, err := os.ReadFile(pwFile)
		if err != nil {
			return "", fmt.Errorf("reading password file %s: %w", pwFile, err)
		}
		return trimTrailingNewline(string(raw)), nil
	}
	return promptPassword(os.Stdin, os.Stderr, label, false)
}

// printImport writes the two-line address/keyfile summary used by
// both the direct and derived import paths.
func printImport(cmd *cobra.Command, address, path string) error {
	w := cmd.OutOrStdout()
	fmt.Fprintf(w, "address  %s\n", address)
	fmt.Fprintf(w, "keyfile  %s\n", path)
	return nil
}
