package wallet

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
)

// newPublicKeyCmd builds `wallet public-key`: print the secp256k1
// public key for a key. Default emits both the uncompressed
// (65-byte 0x04...) and compressed (33-byte 0x02/0x03...) forms on
// separate lines.
//
// Source precedence: <address> positional > --private-key > --keystore-file.
func newPublicKeyCmd() *cobra.Command {
	var (
		shared        keystoreSharedFlags
		privateKey    string
		compressedOnly bool
		uncompressedOnly bool
	)
	cmd := &cobra.Command{
		Use:   "public-key [address]",
		Short: "Print the secp256k1 public key for a key.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var priv *ecdsa.PrivateKey
			switch {
			case privateKey != "":
				p, err := parsePrivateKeyHex(privateKey)
				if err != nil {
					return err
				}
				priv = p
			case len(args) == 1 || shared.KeystoreFile != "":
				identifier := shared.KeystoreFile
				if len(args) == 1 {
					identifier = args[0]
				}
				dir, err := resolveKeystoreDir(shared.KeystoreDir)
				if err != nil {
					return err
				}
				ks := newKeyStore(dir)
				acc, err := findAccount(ks, identifier)
				if err != nil {
					return err
				}
				password, err := readPassword(&shared, os.Stdin, false, "keystore password")
				if err != nil {
					return err
				}
				p, err := decryptKeystoreAccount(acc, password)
				if err != nil {
					return err
				}
				priv = p
			default:
				return &client.UsageError{Msg: "one of address, --keystore-file, or --private-key is required"}
			}
			uncompressed := crypto.FromECDSAPub(&priv.PublicKey)
			compressed := crypto.CompressPubkey(&priv.PublicKey)
			w := cmd.OutOrStdout()
			if !compressedOnly {
				fmt.Fprintf(w, "uncompressed  0x%s\n", hex.EncodeToString(uncompressed))
			}
			if !uncompressedOnly {
				fmt.Fprintf(w, "compressed    0x%s\n", hex.EncodeToString(compressed))
			}
			return nil
		},
	}
	bindKeystoreFlags(cmd, &shared)
	f := cmd.Flags()
	f.StringVar(&privateKey, "private-key", "", "hex-encoded private key (skips the keystore)")
	f.BoolVar(&compressedOnly, "compressed", false, "print only the compressed form")
	f.BoolVar(&uncompressedOnly, "uncompressed", false, "print only the uncompressed form")
	return cmd
}
