package msgs

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
	htx "github.com/0xPolygon/polygon-cli/internal/heimdall/tx"
)

// stakeUpdateMsgShort triggers the L1-mirroring guard.
const stakeUpdateMsgShort = "MsgStakeUpdate"

func init() {
	RegisterFactory("stake-update", newStakeUpdateCmd)
}

// newStakeUpdateCmd builds `stake-update` (MsgStakeUpdate).
func newStakeUpdateCmd(mode Mode, globalFlags *config.Flags) *cobra.Command {
	opts := &TxOpts{Global: globalFlags}
	var (
		fromFlag  string
		valID     uint64
		newAmount string
		l1Ref     stakeL1Ref
	)
	cmd := &cobra.Command{
		Use:   "stake-update",
		Short: "Update validator stake (MsgStakeUpdate, L1-mirroring).",
		Long: strings.TrimSpace(`
Build, sign, and optionally broadcast a heimdallv2.stake.MsgStakeUpdate.

Produced by the bridge after a StakeUpdate event; manual use requires
--force.
`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			from, err := signerOrFlagAddress(cmd, opts, "from-msg", fromFlag)
			if err != nil {
				return err
			}
			if err = requireNonZero("val-id", valID); err != nil {
				return err
			}
			if err = requireNonEmptyString("new-amount", newAmount); err != nil {
				return err
			}
			txHash, err := l1Ref.txHash()
			if err != nil {
				return err
			}
			return executeSingleMsg(cmd, opts, mode, stakeUpdateMsgShort, from, &htx.StakeUpdateMsg{
				From: from, ValID: valID, NewAmount: strings.TrimSpace(newAmount),
				TxHash: txHash, LogIndex: l1Ref.logIndex,
				BlockNumber: l1Ref.blockNumber, Nonce: l1Ref.nonce,
			})
		},
	}
	RegisterFlags(cmd, opts, mode)
	f := cmd.Flags()
	registerFromMsgFlag(f, &fromFlag, stakeUpdateMsgShort)
	f.Uint64Var(&valID, "val-id", 0, "validator id")
	f.StringVar(&newAmount, "new-amount", "", "new stake amount (decimal string)")
	l1Ref.registerFlags(f)
	return cmd
}
