package wallet

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
)

// newRemoveCmd builds `wallet remove <address-or-file>`: delete a key
// from the keystore. Requires either --yes or an interactive y/N
// confirmation.
//
// Deletion uses keystore.Delete, which is irreversible. Operators who
// want a dry-run workflow can use `wallet list` first.
func newRemoveCmd() *cobra.Command {
	var shared keystoreSharedFlags
	cmd := &cobra.Command{
		Use:   "remove <address-or-file>",
		Short: "Remove a key from the keystore.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := resolveKeystoreDir(shared.KeystoreDir)
			if err != nil {
				return err
			}
			ks := newKeyStore(dir)
			acc, err := findAccount(ks, args[0])
			if err != nil {
				return err
			}
			if !shared.Yes {
				if !confirm(cmd, fmt.Sprintf("Delete keystore entry for %s? [y/N]: ", acc.Address.Hex())) {
					return &client.UsageError{Msg: "aborted"}
				}
			}
			password, err := readPassword(&shared, os.Stdin, false, "keystore password")
			if err != nil {
				return err
			}
			if err := ks.Delete(acc, password); err != nil {
				return fmt.Errorf("deleting keystore entry: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "removed %s\n", acc.Address.Hex())
			return nil
		},
	}
	bindKeystoreFlags(cmd, &shared)
	return cmd
}

// confirm reads a y/N answer from the command's input stream (or
// stdin if nothing is wired up). Default is No.
func confirm(cmd *cobra.Command, prompt string) bool {
	fmt.Fprint(cmd.ErrOrStderr(), prompt)
	in := cmd.InOrStdin()
	scanner := bufio.NewScanner(in)
	if !scanner.Scan() {
		return false
	}
	ans := strings.ToLower(strings.TrimSpace(scanner.Text()))
	return ans == "y" || ans == "yes"
}
