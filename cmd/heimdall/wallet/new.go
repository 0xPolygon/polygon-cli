package wallet

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
)

// newNewCmd builds `wallet new`: generate a random secp256k1 key,
// encrypt it with a password, and store it in the resolved keystore
// directory. Prints the address and keyfile path on success.
func newNewCmd() *cobra.Command {
	var shared keystoreSharedFlags
	cmd := &cobra.Command{
		Use:   "new",
		Short: "Generate a new key in the keystore.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := resolveKeystoreDir(shared.KeystoreDir)
			if err != nil {
				return err
			}
			password, err := readPassword(&shared, os.Stdin, true, "keystore password")
			if err != nil {
				return err
			}
			priv, err := crypto.GenerateKey()
			if err != nil {
				return fmt.Errorf("generating key: %w", err)
			}
			ks := newKeyStore(dir)
			acc, err := ks.ImportECDSA(priv, password)
			if err != nil {
				return fmt.Errorf("writing keystore entry: %w", err)
			}
			w := cmd.OutOrStdout()
			fmt.Fprintf(w, "address  %s\n", acc.Address.Hex())
			fmt.Fprintf(w, "keyfile  %s\n", acc.URL.Path)
			return nil
		},
	}
	bindKeystoreFlags(cmd, &shared)
	rejectHardwareFlags(cmd)
	return cmd
}
