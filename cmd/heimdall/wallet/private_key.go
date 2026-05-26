package wallet

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
)

// newPrivateKeyCmd builds `wallet private-key <address>`: decrypt the
// keystore entry for an address and print the plaintext private key.
// Same underlying operation as `decrypt-keystore` but addressed by
// address rather than file path.
func newPrivateKeyCmd() *cobra.Command {
	var (
		shared      keystoreSharedFlags
		acknowledge bool
	)
	cmd := &cobra.Command{
		Use:   "private-key <address-or-file>",
		Short: "Print the plaintext private key for a keystore entry.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !acknowledge {
				return &client.UsageError{Msg: "refusing to print a private key without --i-understand-the-risks"}
			}
			dir, err := resolveKeystoreDir(shared.KeystoreDir)
			if err != nil {
				return err
			}
			ks := newKeyStore(dir)
			acc, err := findAccount(ks, args[0])
			if err != nil {
				return err
			}
			password, err := readPassword(&shared, os.Stdin, false, "keystore password")
			if err != nil {
				return err
			}
			priv, err := decryptKeystoreAccount(acc, password)
			if err != nil {
				return fmt.Errorf("decrypting keystore entry: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "0x%s\n", hex.EncodeToString(crypto.FromECDSA(priv)))
			return nil
		},
	}
	bindKeystoreFlags(cmd, &shared)
	cmd.Flags().BoolVar(&acknowledge, "i-understand-the-risks", false, "required friction flag for exposing plaintext key material")
	return cmd
}
