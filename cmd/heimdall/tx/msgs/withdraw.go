package msgs

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
	htx "github.com/0xPolygon/polygon-cli/internal/heimdall/tx"
)

// withdrawMsgShort is the short Msg name matched against the
// L1-mirroring guard. MsgWithdrawFeeTx is NOT an L1-mirroring msg
// (validators withdraw their own accumulated fees), so RequireForce
// returns nil.
const withdrawMsgShort = "MsgWithdrawFeeTx"

func init() {
	RegisterFactory("withdraw", newWithdrawCmd)
}

// newWithdrawCmd returns a cobra command that executes MsgWithdrawFeeTx
// under the given mode. Both `user` and `amount` are optional:
//
//   - When --user is omitted, the signer's address (resolved from
//     --from / --account / --private-key / --mnemonic) is used.
//   - When --amount is omitted or "0", Heimdall withdraws the full
//     balance of the account's accumulated fees. The proto field is
//     a math.Int decimal string; we pass "0" through unchanged because
//     that's the on-chain sentinel for "all".
func newWithdrawCmd(mode Mode, globalFlags *config.Flags) *cobra.Command {
	opts := &TxOpts{Global: globalFlags}
	var userFlag string
	var amountFlag string

	cmd := &cobra.Command{
		Use:   "withdraw",
		Short: "Withdraw accumulated validator fees.",
		Long: strings.TrimSpace(`
Build (or send, or estimate) a MsgWithdrawFeeTx that withdraws a
validator's accumulated Heimdall fees into the main bank balance.

The signing key and the on-chain proposer address are both derived
from --from unless --user is set explicitly. --amount defaults to
"0", which means "withdraw all" per Heimdall semantics.
`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Resolve the proposer address. --user wins; otherwise the
			// signer's address is used.
			proposer := strings.ToLower(strings.TrimSpace(userFlag))
			if proposer == "" {
				// We need a signer to derive the address. Resolve the
				// signer once here and pass its Eth-style address in as
				// the proposer; Execute will resolve the signer a second
				// time internally to actually sign, so the private key
				// never has to leave this function. We could thread the
				// signer through Plan but that would couple Plan to
				// key material for no real benefit.
				signer, err := ResolveSigningKey(opts, cmd.InOrStdin())
				if err != nil {
					return err
				}
				proposer = strings.ToLower(signer.Address.Hex())
			} else {
				if err := requireEthAddress(proposer); err != nil {
					return err
				}
			}
			amount := strings.TrimSpace(amountFlag)
			if amount == "" {
				amount = "0"
			}

			plan := &Plan{
				Msgs: []htx.Msg{&htx.WithdrawFeeMsg{
					Proposer: proposer,
					Amount:   amount,
				}},
				MsgShortType:  withdrawMsgShort,
				SignerAddress: proposer,
			}
			return Execute(cmd, opts, mode, plan)
		},
	}
	RegisterFlags(cmd, opts, mode)
	f := cmd.Flags()
	f.StringVar(&userFlag, "user", "", "address withdrawing fees (default: signer address)")
	f.StringVar(&amountFlag, "amount", "0", "amount to withdraw as decimal integer; 0 means withdraw all")
	return cmd
}

// requireEthAddress returns a usage error unless s parses as a
// 20-byte `0x`-prefixed hex address. We intentionally keep this in
// the msgs package (instead of reusing the parent `tx` helper) to
// avoid an import cycle between `tx` and `tx/msgs`.
func requireEthAddress(s string) error {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(strings.TrimPrefix(s, "0x"), "0X")
	if len(s) != 40 {
		return &client.UsageError{Msg: "--user must be a 20-byte (40 hex char) address"}
	}
	for _, c := range s {
		ok := (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
		if !ok {
			return &client.UsageError{Msg: "--user must be hex"}
		}
	}
	return nil
}
