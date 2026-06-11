package msgs

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
	htx "github.com/0xPolygon/polygon-cli/internal/heimdall/tx"
)

// validatorExitMsgShort triggers the L1-mirroring guard.
const validatorExitMsgShort = "MsgValidatorExit"

func init() {
	RegisterFactory("stake-exit", newStakeExitCmd)
}

// newStakeExitCmd builds `stake-exit` (MsgValidatorExit). L1-mirroring.
func newStakeExitCmd(mode Mode, globalFlags *config.Flags) *cobra.Command {
	opts := &TxOpts{Global: globalFlags}
	var (
		fromFlag          string
		valID             uint64
		deactivationEpoch uint64
		l1Ref             stakeL1Ref
	)
	cmd := &cobra.Command{
		Use:   "stake-exit",
		Short: "Mark validator exit (MsgValidatorExit, L1-mirroring).",
		Long: strings.TrimSpace(`
Build, sign, and optionally broadcast a heimdallv2.stake.MsgValidatorExit.

Produced by the bridge after an Unstake event; manual use requires
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
			txHash, err := l1Ref.txHash()
			if err != nil {
				return err
			}
			return executeSingleMsg(cmd, opts, mode, validatorExitMsgShort, from, &htx.ValidatorExitMsg{
				From: from, ValID: valID, DeactivationEpoch: deactivationEpoch,
				TxHash: txHash, LogIndex: l1Ref.logIndex,
				BlockNumber: l1Ref.blockNumber, Nonce: l1Ref.nonce,
			})
		},
	}
	RegisterFlags(cmd, opts, mode)
	f := cmd.Flags()
	registerFromMsgFlag(f, &fromFlag, validatorExitMsgShort)
	f.Uint64Var(&valID, "val-id", 0, "validator id")
	f.Uint64Var(&deactivationEpoch, "deactivation-epoch", 0, "deactivation epoch")
	l1Ref.registerFlags(f)
	return cmd
}
