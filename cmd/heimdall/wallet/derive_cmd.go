package wallet

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
)

// newDeriveCmd builds `wallet derive`: derive one or more addresses
// from a mnemonic.
//
// Supports two modes:
//  1. Single path / index — default, prints one address and its
//     derivation path.
//  2. `--count N` — derives N sequential addresses starting at
//     `--index`, incrementing the final path component each time.
func newDeriveCmd() *cobra.Command {
	var (
		mnemonic     string
		mnemonicFile string
		bipPass      string
		path         string
		index        uint32
		count        uint32
		showKey      bool
	)
	cmd := &cobra.Command{
		Use:   "derive",
		Short: "Derive addresses from a BIP-39 mnemonic.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if mnemonic == "" && mnemonicFile != "" {
				raw, err := os.ReadFile(mnemonicFile)
				if err != nil {
					return fmt.Errorf("reading mnemonic file %s: %w", mnemonicFile, err)
				}
				mnemonic = trimTrailingNewline(string(raw))
			}
			if mnemonic == "" {
				return &client.UsageError{Msg: "one of --mnemonic or --mnemonic-file is required"}
			}
			if count == 0 {
				count = 1
			}
			w := cmd.OutOrStdout()
			for i := uint32(0); i < count; i++ {
				priv, finalPath, addr, err := deriveFromMnemonic(mnemonic, bipPass, path, index+i)
				if err != nil {
					return err
				}
				if showKey {
					fmt.Fprintf(w, "%s\t%s\t%s\n", finalPath, addr.Hex(), "0x"+hex.EncodeToString(crypto.FromECDSA(priv)))
				} else {
					fmt.Fprintf(w, "%s\t%s\n", finalPath, addr.Hex())
				}
			}
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVar(&mnemonic, "mnemonic", "", "BIP-39 mnemonic")
	f.StringVar(&mnemonicFile, "mnemonic-file", "", "file containing a BIP-39 mnemonic")
	f.StringVar(&bipPass, "bip39-passphrase", "", "optional BIP-39 passphrase")
	f.StringVar(&path, "path", "", "derivation path (default m/44'/60'/0'/0/<index>)")
	f.Uint32Var(&index, "index", 0, "starting address index when --path is not set")
	f.Uint32Var(&count, "count", 1, "number of sequential addresses to derive")
	f.BoolVar(&showKey, "show-private-key", false, "also emit the derived private key on each line")
	return cmd
}
