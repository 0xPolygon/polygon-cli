// Package wallet implements the `polycli heimdall wallet` umbrella
// command. It is a local-only command group: none of the subcommands
// talk to the network. Keys are stored in a go-ethereum v3 JSON
// keystore directory that is compatible with Foundry's `cast wallet`.
//
// Keystore directory precedence (highest wins):
//  1. `--keystore-dir` flag.
//  2. `ETH_KEYSTORE` environment variable.
//  3. `~/.foundry/keystores/` if it already exists (honour existing
//     cast users without migration).
//  4. `~/.polycli/keystores/` (default; created on demand).
//
// Signing uses EIP-191 personal_sign by default. `--raw` signs a
// 32-byte hash directly. Hardware wallets (`--ledger`, `--trezor`),
// `vanity`, and `sign-auth` from cast are deliberately rejected with
// a pointer at `cast wallet` — see HEIMDALLCAST_REQUIREMENTS.md §3.4.
package wallet

import (
	_ "embed"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
	sharedwallet "github.com/0xPolygon/polygon-cli/internal/heimdall/wallet"
)

//go:embed usage.md
var usage string

// flags is injected by Register. None of the wallet subcommands call
// config.Resolve — the heimdall network config is irrelevant to local
// key management — but we keep the handle for symmetry with the other
// command groups so future additions (e.g. reading the default chain
// id for tx signing hints) have it without re-plumbing.
var flags *config.Flags

// newWalletCmd builds a fresh `wallet` umbrella. Constructed per
// Register call so tests that re-wire a parent do not accumulate
// duplicate subcommands on a shared command tree.
func newWalletCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wallet",
		Short: "Manage keystores, keys, and message signatures.",
		Long:  usage,
		Args:  cobra.NoArgs,
	}
	cmd.AddCommand(
		newNewCmd(),
		newNewMnemonicCmd(),
		newAddressCmd(),
		newDeriveCmd(),
		newSignCmd(),
		newVerifyCmd(),
		newImportCmd(),
		newListCmd(),
		newRemoveCmd(),
		newPublicKeyCmd(),
		newDecryptKeystoreCmd(),
		newChangePasswordCmd(),
		newPrivateKeyCmd(),
		rejectedSubcommand("vanity", "Not supported — use `cast wallet vanity`.", "cast wallet vanity"),
		rejectedSubcommand("sign-auth", "Not supported — use `cast wallet sign-auth`.", "cast wallet sign-auth"),
	)
	return cmd
}

// Register attaches the wallet umbrella command and its subcommands to
// parent. The shared flag struct is stored for future use; wallet
// subcommands do not currently consume it.
func Register(parent *cobra.Command, f *config.Flags) {
	flags = f
	parent.AddCommand(newWalletCmd())
}

// keystoreSharedFlags holds the flags common to every subcommand that
// reads or writes the keystore. Subcommands embed an instance of this
// into their command and call resolveKeystoreDir to get the final
// on-disk directory.
type keystoreSharedFlags struct {
	KeystoreDir  string
	KeystoreFile string
	Password     string
	PasswordFile string
	Yes          bool
}

// bindKeystoreFlags attaches the shared keystore flags to cmd's
// flag set. All of these are local flags — never persistent.
func bindKeystoreFlags(cmd *cobra.Command, s *keystoreSharedFlags) {
	f := cmd.Flags()
	f.StringVar(&s.KeystoreDir, "keystore-dir", "", "keystore directory (overrides ETH_KEYSTORE, ~/.foundry/keystores, ~/.polycli/keystores)")
	f.StringVar(&s.KeystoreFile, "keystore-file", "", "explicit keystore JSON file path")
	f.StringVar(&s.Password, "password", "", "keystore password (mutually exclusive with --password-file)")
	f.StringVar(&s.PasswordFile, "password-file", "", "path to a file containing the keystore password")
	f.BoolVar(&s.Yes, "yes", false, "skip confirmation prompts")
}

// resolveKeystoreDir returns the keystore directory to use per the
// precedence rule documented on the package doc comment. It delegates
// to internal/heimdall/wallet.ResolveKeystoreDir with createDefault=true
// so the polycli fallback directory is materialised on first use — the
// wallet subcommands are keystore-management commands and an operator
// who types `polycli heimdall wallet list` on a fresh machine wants an
// empty dir, not an error.
func resolveKeystoreDir(override string) (string, error) {
	return sharedwallet.ResolveKeystoreDir(override, true)
}

// readPassword returns the password for a keystore operation. The
// precedence is: --password flag > --password-file > interactive
// prompt from stdin (if the caller's stdin is a terminal, otherwise
// the full line is read). Returning an empty password is allowed —
// go-ethereum's keystore will still accept it, which is what cast
// users expect.
func readPassword(s *keystoreSharedFlags, in io.Reader, confirm bool, label string) (string, error) {
	if s.Password != "" && s.PasswordFile != "" {
		return "", &client.UsageError{Msg: "--password and --password-file are mutually exclusive"}
	}
	if s.Password != "" {
		return s.Password, nil
	}
	if s.PasswordFile != "" {
		raw, err := os.ReadFile(s.PasswordFile)
		if err != nil {
			return "", fmt.Errorf("reading password file %s: %w", s.PasswordFile, err)
		}
		return trimTrailingNewline(string(raw)), nil
	}
	return promptPassword(in, os.Stderr, label, confirm)
}

// trimTrailingNewline strips a single trailing \n or \r\n so
// password-file contents can include a terminating newline without
// invalidating the password. A trailing whitespace character beyond a
// simple newline is preserved — operators who put intentional
// whitespace in a password file are probably not making a mistake.
func trimTrailingNewline(s string) string {
	if n := len(s); n > 0 && s[n-1] == '\n' {
		if n >= 2 && s[n-2] == '\r' {
			return s[:n-2]
		}
		return s[:n-1]
	}
	return s
}
