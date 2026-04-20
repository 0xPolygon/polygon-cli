package wallet

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/gethkeystore"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
)

// newDecryptKeystoreCmd builds `wallet decrypt-keystore <file>`:
// decrypt an arbitrary keystore file and print the private key hex.
// Requires the `--i-understand-the-risks` friction flag to avoid
// accidental plaintext exposure.
func newDecryptKeystoreCmd() *cobra.Command {
	var (
		shared      keystoreSharedFlags
		acknowledge bool
	)
	cmd := &cobra.Command{
		Use:   "decrypt-keystore <file>",
		Short: "Decrypt a keystore file to its plaintext private key.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !acknowledge {
				return &client.UsageError{Msg: "refusing to print a private key without --i-understand-the-risks"}
			}
			data, err := os.ReadFile(args[0])
			if err != nil {
				return fmt.Errorf("reading keystore file %s: %w", args[0], err)
			}
			password, err := readPassword(&shared, os.Stdin, false, "keystore password")
			if err != nil {
				return err
			}
			priv, err := gethkeystore.DecryptKeystoreFile(data, password)
			if err != nil {
				return fmt.Errorf("decrypting keystore: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "0x%s\n", hex.EncodeToString(crypto.FromECDSA(priv)))
			return nil
		},
	}
	bindKeystoreFlags(cmd, &shared)
	cmd.Flags().BoolVar(&acknowledge, "i-understand-the-risks", false, "required friction flag for exposing plaintext key material")
	return cmd
}
