package wallet

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
)

// newAddressCmd builds `wallet address`: report the Ethereum address
// for one of several inputs.
//
// Sources, in precedence order: --private-key, --mnemonic,
// --keystore-file, or (if none of the above) list every address in
// the resolved keystore directory.
func newAddressCmd() *cobra.Command {
	var (
		shared     keystoreSharedFlags
		privateKey string
		mnemonic   string
		bipPass    string
		path       string
		index      uint32
	)
	cmd := &cobra.Command{
		Use:   "address",
		Short: "Show the address for a key or keystore file.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			w := cmd.OutOrStdout()
			switch {
			case privateKey != "":
				priv, err := parsePrivateKeyHex(privateKey)
				if err != nil {
					return err
				}
				fmt.Fprintln(w, crypto.PubkeyToAddress(priv.PublicKey).Hex())
				return nil
			case mnemonic != "":
				_, _, addr, err := deriveFromMnemonic(mnemonic, bipPass, path, index)
				if err != nil {
					return err
				}
				fmt.Fprintln(w, addr.Hex())
				return nil
			case shared.KeystoreFile != "":
				addr, err := addressFromKeystoreFile(shared.KeystoreFile)
				if err != nil {
					return err
				}
				fmt.Fprintln(w, addr.Hex())
				return nil
			}
			// No explicit source — list every address in the keystore.
			dir, err := resolveKeystoreDir(shared.KeystoreDir)
			if err != nil {
				return err
			}
			ks := newKeyStore(dir)
			accounts := ks.Accounts()
			if len(accounts) == 0 {
				return &client.UsageError{Msg: fmt.Sprintf("no keys in keystore %s; pass --private-key, --mnemonic, or --keystore-file to inspect another source", dir)}
			}
			for _, a := range accounts {
				fmt.Fprintln(w, a.Address.Hex())
			}
			return nil
		},
	}
	bindKeystoreFlags(cmd, &shared)
	f := cmd.Flags()
	f.StringVar(&privateKey, "private-key", "", "hex-encoded secp256k1 private key")
	f.StringVar(&mnemonic, "mnemonic", "", "BIP-39 mnemonic")
	f.StringVar(&bipPass, "bip39-passphrase", "", "optional BIP-39 passphrase")
	f.StringVar(&path, "path", "", "derivation path (default m/44'/60'/0'/0/<index>)")
	f.Uint32Var(&index, "index", 0, "address index used when --path is not set")
	return cmd
}

// parsePrivateKeyHex decodes a 0x-prefixed or bare 32-byte hex string
// into an ECDSA private key.
func parsePrivateKeyHex(input string) (*ecdsa.PrivateKey, error) {
	s := strings.TrimSpace(input)
	s = strings.TrimPrefix(strings.TrimPrefix(s, "0x"), "0X")
	if len(s) != 64 {
		return nil, fmt.Errorf("private key must be 32 bytes (64 hex chars), got %d", len(s))
	}
	raw, err := hex.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("decoding private key: %w", err)
	}
	return crypto.ToECDSA(raw)
}
