package wallet

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tyler-smith/go-bip39"
)

// newNewMnemonicCmd builds `wallet new-mnemonic`: generate a fresh
// BIP-39 mnemonic, derive the key at m/44'/60'/0'/0/<index>, import
// it into the keystore, and print the mnemonic once to stderr with a
// prominent warning.
func newNewMnemonicCmd() *cobra.Command {
	var (
		shared     keystoreSharedFlags
		words      int
		bipPass    string
		path       string
		index      uint32
		printOnly  bool
	)
	cmd := &cobra.Command{
		Use:   "new-mnemonic",
		Short: "Generate a new BIP-39 mnemonic and derive a key.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			bits, ok := mnemonicWordCountBits(words)
			if !ok {
				return fmt.Errorf("--words must be 12, 15, 18, 21, or 24; got %d", words)
			}
			entropy, err := bip39.NewEntropy(bits)
			if err != nil {
				return fmt.Errorf("generating entropy: %w", err)
			}
			mnemonic, err := bip39.NewMnemonic(entropy)
			if err != nil {
				return fmt.Errorf("generating mnemonic: %w", err)
			}
			priv, finalPath, addr, err := deriveFromMnemonic(mnemonic, bipPass, path, index)
			if err != nil {
				return err
			}
			w := cmd.OutOrStdout()
			errW := cmd.ErrOrStderr()

			if printOnly {
				// Operator just wants the mnemonic + derived key, no
				// keystore side-effect. Useful for scripting.
				fmt.Fprintf(w, "address     %s\n", addr.Hex())
				fmt.Fprintf(w, "path        %s\n", finalPath)
				fmt.Fprintf(w, "mnemonic    %s\n", mnemonic)
				fmt.Fprintln(errW, "WARNING: the mnemonic above is the only copy. Record it now; it will not be shown again.")
				return nil
			}

			dir, err := resolveKeystoreDir(shared.KeystoreDir)
			if err != nil {
				return err
			}
			password, err := readPassword(&shared, os.Stdin, true, "keystore password")
			if err != nil {
				return err
			}
			ks := newKeyStore(dir)
			acc, err := ks.ImportECDSA(priv, password)
			if err != nil {
				return fmt.Errorf("writing keystore entry: %w", err)
			}
			if acc.Address != addr {
				// Should never happen; defensive check against an
				// accidental shift in the derivation logic.
				return fmt.Errorf("keystore recorded %s but derivation produced %s", acc.Address.Hex(), addr.Hex())
			}
			fmt.Fprintf(w, "address     %s\n", acc.Address.Hex())
			fmt.Fprintf(w, "path        %s\n", finalPath)
			fmt.Fprintf(w, "keyfile     %s\n", acc.URL.Path)
			fmt.Fprintf(w, "mnemonic    %s\n", mnemonic)
			fmt.Fprintln(errW, strings.Repeat("!", 70))
			fmt.Fprintln(errW, "WARNING: the mnemonic above is the ONLY copy polycli will print.")
			fmt.Fprintln(errW, "Record it somewhere safe; losing it means losing the key.")
			fmt.Fprintln(errW, strings.Repeat("!", 70))
			return nil
		},
	}
	bindKeystoreFlags(cmd, &shared)
	f := cmd.Flags()
	f.IntVar(&words, "words", 12, "mnemonic word count (12, 15, 18, 21, 24)")
	f.StringVar(&bipPass, "bip39-passphrase", "", "optional BIP-39 passphrase (not the keystore password)")
	f.StringVar(&path, "path", "", "derivation path (default m/44'/60'/0'/0/<index>)")
	f.Uint32Var(&index, "index", 0, "address index used when --path is not set")
	f.BoolVar(&printOnly, "print-only", false, "print mnemonic and derived address without writing to keystore")
	rejectHardwareFlags(cmd)
	return cmd
}

// mnemonicWordCountBits maps a BIP-39 mnemonic word count to the bit
// length of entropy required. 12/15/18/21/24 -> 128/160/192/224/256.
func mnemonicWordCountBits(words int) (int, bool) {
	switch words {
	case 12:
		return 128, true
	case 15:
		return 160, true
	case 18:
		return 192, true
	case 21:
		return 224, true
	case 24:
		return 256, true
	}
	return 0, false
}
