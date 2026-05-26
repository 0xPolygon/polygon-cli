package wallet

import (
	"fmt"

	"github.com/spf13/cobra"
)

// newListCmd builds `wallet list`: print the addresses and keyfile
// paths for every key in the resolved keystore directory.
func newListCmd() *cobra.Command {
	var shared keystoreSharedFlags
	var addressesOnly bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List keys in the keystore.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := resolveKeystoreDir(shared.KeystoreDir)
			if err != nil {
				return err
			}
			ks := newKeyStore(dir)
			accounts := ks.Accounts()
			w := cmd.OutOrStdout()
			if len(accounts) == 0 {
				fmt.Fprintf(w, "(no keys in %s)\n", dir)
				return nil
			}
			for _, a := range accounts {
				if addressesOnly {
					fmt.Fprintln(w, a.Address.Hex())
					continue
				}
				fmt.Fprintf(w, "%s\t%s\n", a.Address.Hex(), a.URL.Path)
			}
			return nil
		},
	}
	bindKeystoreFlags(cmd, &shared)
	cmd.Flags().BoolVar(&addressesOnly, "addresses-only", false, "print only addresses, no keyfile paths")
	return cmd
}
