package msgs

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
	htx "github.com/0xPolygon/polygon-cli/internal/heimdall/tx"
)

// checkpointNoAckMsgShort triggers the L1-mirroring guard.
const checkpointNoAckMsgShort = "MsgCpNoAck"

func init() {
	RegisterFactory("checkpoint-noack", newCheckpointNoAckCmd)
}

// newCheckpointNoAckCmd builds `checkpoint-noack`. The proto only
// carries `from`; the bridge produces these when a checkpoint window
// lapses without an ack on L1.
func newCheckpointNoAckCmd(mode Mode, globalFlags *config.Flags) *cobra.Command {
	opts := &TxOpts{Global: globalFlags}
	var fromFlag string
	cmd := &cobra.Command{
		Use:   "checkpoint-noack",
		Short: "Mark missed checkpoint ack (MsgCpNoAck, L1-mirroring).",
		Long: strings.TrimSpace(`
Build, sign, and optionally broadcast a heimdallv2.checkpoint.MsgCpNoAck.

MsgCpNoAck is produced by the bridge when an L1 checkpoint window
lapses without an ack. Manual use is almost never correct; the command
refuses without --force.
`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			from, err := signerOrFlagAddress(cmd, opts, "from-msg", fromFlag)
			if err != nil {
				return err
			}
			return executeSingleMsg(cmd, opts, mode, checkpointNoAckMsgShort, from, &htx.CpNoAckMsg{From: from})
		},
	}
	RegisterFlags(cmd, opts, mode)
	registerFromMsgFlag(cmd.Flags(), &fromFlag, checkpointNoAckMsgShort)
	return cmd
}
