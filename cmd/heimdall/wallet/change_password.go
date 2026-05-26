package wallet

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
)

// newChangePasswordCmd builds `wallet change-password <address-or-file>`:
// re-encrypt a keystore entry under a new password.
//
// The underlying ks.Update preserves the file path, so the entry is
// rewritten in-place rather than duplicated.
func newChangePasswordCmd() *cobra.Command {
	var (
		shared         keystoreSharedFlags
		newPassword    string
		newPasswordFile string
	)
	cmd := &cobra.Command{
		Use:   "change-password <address-or-file>",
		Short: "Change a keystore entry's password.",
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
			oldPw, err := readPassword(&shared, os.Stdin, false, "current keystore password")
			if err != nil {
				return err
			}
			newPw, err := readNewPassword(newPassword, newPasswordFile, os.Stdin, true)
			if err != nil {
				return err
			}
			if newPw == oldPw && !shared.Yes {
				return &client.UsageError{Msg: "new password matches the current one; pass --yes to keep it anyway"}
			}
			if err := ks.Update(acc, oldPw, newPw); err != nil {
				return fmt.Errorf("updating keystore entry: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "updated %s\n", acc.Address.Hex())
			return nil
		},
	}
	bindKeystoreFlags(cmd, &shared)
	f := cmd.Flags()
	f.StringVar(&newPassword, "new-password", "", "new keystore password")
	f.StringVar(&newPasswordFile, "new-password-file", "", "file containing the new keystore password")
	return cmd
}

// readNewPassword resolves the replacement password per the same
// file/flag/prompt rules as readPassword, but keyed on the separate
// --new-password / --new-password-file flags.
func readNewPassword(val, file string, in *os.File, confirm bool) (string, error) {
	if val != "" && file != "" {
		return "", &client.UsageError{Msg: "--new-password and --new-password-file are mutually exclusive"}
	}
	if val != "" {
		return val, nil
	}
	if file != "" {
		raw, err := os.ReadFile(file)
		if err != nil {
			return "", fmt.Errorf("reading new-password file %s: %w", file, err)
		}
		return trimTrailingNewline(string(raw)), nil
	}
	return promptPassword(in, os.Stderr, "new keystore password", confirm)
}
