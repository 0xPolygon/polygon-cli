package wallet

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
)

// rejectHardwareFlags attaches `--ledger` and `--trezor` boolean flags
// that always error with a pointer at cast when set. Hardware-wallet
// support is explicitly out of scope (requirements §3.4); surfacing
// the flags with a helpful message is friendlier than letting cast
// muscle-memory operators hit an unrecognised-flag error.
func rejectHardwareFlags(cmd *cobra.Command) {
	var ledger, trezor bool
	cmd.Flags().BoolVar(&ledger, "ledger", false, "not supported; use `cast wallet --ledger`")
	cmd.Flags().BoolVar(&trezor, "trezor", false, "not supported; use `cast wallet --trezor`")
	wrapped := cmd.RunE
	cmd.RunE = func(c *cobra.Command, args []string) error {
		if ledger {
			return &client.UsageError{Msg: "hardware wallets are not supported by polycli; use `cast wallet --ledger`"}
		}
		if trezor {
			return &client.UsageError{Msg: "hardware wallets are not supported by polycli; use `cast wallet --trezor`"}
		}
		if wrapped != nil {
			return wrapped(c, args)
		}
		return nil
	}
}

// rejectedSubcommand returns a cobra command whose RunE always errors
// with a pointer to the equivalent cast subcommand. Use for `vanity`
// and `sign-auth` which are intentionally unimplemented.
func rejectedSubcommand(use, short, castEquivalent string) *cobra.Command {
	return &cobra.Command{
		Use:                use,
		Short:              short,
		Args:               cobra.ArbitraryArgs,
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return &client.UsageError{Msg: fmt.Sprintf("not supported by polycli; use `%s`", castEquivalent)}
		},
	}
}
