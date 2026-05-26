package msgs

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
	htx "github.com/0xPolygon/polygon-cli/internal/heimdall/tx"
)

// checkpointMsgShort is the short Msg name matched against the
// L1-mirroring guard. MsgCheckpoint is NOT on the L1-mirroring list —
// validators legitimately propose checkpoints — so RequireForce
// returns nil. The `--i-am-a-validator` flag adds an explicit
// acknowledgement so hands-on operators cannot fire it by accident.
const checkpointMsgShort = "MsgCheckpoint"

func init() {
	RegisterFactory("checkpoint", newCheckpointCmd)
}

// newCheckpointCmd builds the `checkpoint` subcommand under the given
// umbrella (mktx / send / estimate). It constructs a MsgCheckpoint
// and hands it to Execute.
//
// The message is validator-only in practice. We guard it behind
// --i-am-a-validator instead of --force because MsgCheckpoint is not
// an L1-mirroring msg; the friction check is ours, not the guard's.
func newCheckpointCmd(mode Mode, globalFlags *config.Flags) *cobra.Command {
	opts := &TxOpts{Global: globalFlags}
	var (
		proposerFlag        string
		startBlock          uint64
		endBlock            uint64
		rootHashHex         string
		accountRootHashHex  string
		borChainID          string
		iAmAValidator       bool
	)

	cmd := &cobra.Command{
		Use:   "checkpoint",
		Short: "Propose a checkpoint (MsgCheckpoint).",
		Long: strings.TrimSpace(`
Build, sign, and optionally broadcast a heimdallv2.checkpoint.MsgCheckpoint.

This message is validator-only. --i-am-a-validator is required as an
explicit acknowledgement; pass --force to bypass if you know what you
are doing.
`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !iAmAValidator && !opts.Force {
				return &client.UsageError{Msg: "MsgCheckpoint is validator-only; re-run with --i-am-a-validator"}
			}

			proposer := strings.TrimSpace(proposerFlag)
			if proposer == "" {
				signer, err := ResolveSigningKey(opts, cmd.InOrStdin())
				if err != nil {
					return err
				}
				proposer = strings.ToLower(signer.Address.Hex())
			} else {
				p, err := lowerEthAddress("proposer", proposer)
				if err != nil {
					return err
				}
				proposer = p
			}

			if err := requireNonEmptyString("bor-chain-id", borChainID); err != nil {
				return err
			}

			rootHash, err := parseHexBytes("root-hash", rootHashHex, 32)
			if err != nil {
				return err
			}
			if len(rootHash) == 0 {
				return &client.UsageError{Msg: "--root-hash is required"}
			}
			accRootHash, err := parseHexBytes("account-root-hash", accountRootHashHex, 32)
			if err != nil {
				return err
			}

			plan := &Plan{
				Msgs: []htx.Msg{&htx.CheckpointMsg{
					Proposer:        proposer,
					StartBlock:      startBlock,
					EndBlock:        endBlock,
					RootHash:        rootHash,
					AccountRootHash: accRootHash,
					BorChainID:      strings.TrimSpace(borChainID),
				}},
				MsgShortType:  checkpointMsgShort,
				SignerAddress: proposer,
			}
			return Execute(cmd, opts, mode, plan)
		},
	}
	RegisterFlags(cmd, opts, mode)
	f := cmd.Flags()
	f.StringVar(&proposerFlag, "proposer", "", "proposer address (default: signer)")
	f.Uint64Var(&startBlock, "start-block", 0, "bor start block number (inclusive)")
	f.Uint64Var(&endBlock, "end-block", 0, "bor end block number (inclusive)")
	f.StringVar(&rootHashHex, "root-hash", "", "32-byte bor block root hash (hex)")
	f.StringVar(&accountRootHashHex, "account-root-hash", "", "32-byte account root hash (hex, optional)")
	f.StringVar(&borChainID, "bor-chain-id", "", "bor chain id the checkpoint applies to")
	f.BoolVar(&iAmAValidator, "i-am-a-validator", false, "acknowledge that MsgCheckpoint is validator-only")
	return cmd
}
