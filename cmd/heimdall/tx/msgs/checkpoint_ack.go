package msgs

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
	htx "github.com/0xPolygon/polygon-cli/internal/heimdall/tx"
)

// checkpointAckMsgShort triggers the L1-mirroring guard: unless --force
// is set, Execute refuses before building.
const checkpointAckMsgShort = "MsgCpAck"

func init() {
	RegisterFactory("checkpoint-ack", newCheckpointAckCmd)
}

// newCheckpointAckCmd builds the `checkpoint-ack` subcommand. The
// message is produced by the bridge after observing an L1 event; the
// CLI requires --l1-tx so operators think twice before sending one.
func newCheckpointAckCmd(mode Mode, globalFlags *config.Flags) *cobra.Command {
	opts := &TxOpts{Global: globalFlags}
	var (
		fromFlag    string
		number      uint64
		proposer    string
		startBlock  uint64
		endBlock    uint64
		rootHashHex string
		l1TxHex     string
	)

	cmd := &cobra.Command{
		Use:   "checkpoint-ack",
		Short: "Acknowledge a checkpoint on L2 (MsgCpAck, L1-mirroring).",
		Long: strings.TrimSpace(`
Build, sign, and optionally broadcast a heimdallv2.checkpoint.MsgCpAck.

MsgCpAck is produced by the bridge after observing an L1 event. Manual
use is a replay that competes with the real bridge path; the command
refuses to run without --force. --l1-tx identifies the L1 tx hash the
operator intends to mirror (advisory — not part of the proto).
`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(l1TxHex) == "" {
				return &client.UsageError{Msg: "--l1-tx is required (even with --force) to cite the L1 tx being mirrored"}
			}
			if _, err := parseHexBytes("l1-tx", l1TxHex, 32); err != nil {
				return err
			}

			from := strings.TrimSpace(fromFlag)
			if from == "" {
				signer, err := ResolveSigningKey(opts, cmd.InOrStdin())
				if err != nil {
					return err
				}
				from = strings.ToLower(signer.Address.Hex())
			} else {
				p, err := lowerEthAddress("from", from)
				if err != nil {
					return err
				}
				from = p
			}
			rootHash, err := parseHexBytes("root-hash", rootHashHex, 32)
			if err != nil {
				return err
			}

			plan := &Plan{
				Msgs: []htx.Msg{&htx.CpAckMsg{
					From:       from,
					Number:     number,
					Proposer:   strings.ToLower(strings.TrimSpace(proposer)),
					StartBlock: startBlock,
					EndBlock:   endBlock,
					RootHash:   rootHash,
				}},
				MsgShortType:  checkpointAckMsgShort,
				SignerAddress: from,
			}
			return Execute(cmd, opts, mode, plan)
		},
	}
	RegisterFlags(cmd, opts, mode)
	f := cmd.Flags()
	f.StringVar(&fromFlag, "from-msg", "", "MsgCpAck.from address (default: signer)")
	f.Uint64Var(&number, "number", 0, "checkpoint number on Heimdall")
	f.StringVar(&proposer, "proposer", "", "original proposer address of the checkpoint")
	f.Uint64Var(&startBlock, "start-block", 0, "bor start block number")
	f.Uint64Var(&endBlock, "end-block", 0, "bor end block number")
	f.StringVar(&rootHashHex, "root-hash", "", "32-byte root hash (hex)")
	f.StringVar(&l1TxHex, "l1-tx", "", "L1 transaction hash being mirrored (32 bytes hex)")
	return cmd
}
