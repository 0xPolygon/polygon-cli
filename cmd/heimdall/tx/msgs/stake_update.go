package msgs

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
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
		fromFlag    string
		valID       uint64
		newAmount   string
		txHashHex   string
		logIndex    uint64
		blockNumber uint64
		nonce       uint64
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
			from := strings.TrimSpace(fromFlag)
			if from == "" {
				signer, err := ResolveSigningKey(opts, cmd.InOrStdin())
				if err != nil {
					return err
				}
				from = strings.ToLower(signer.Address.Hex())
			} else {
				p, err := lowerEthAddress("from-msg", from)
				if err != nil {
					return err
				}
				from = p
			}
			if valID == 0 {
				return &client.UsageError{Msg: "--val-id is required"}
			}
			if strings.TrimSpace(newAmount) == "" {
				return &client.UsageError{Msg: "--new-amount is required"}
			}
			txHash, err := parseHexBytes("tx-hash", txHashHex, 32)
			if err != nil {
				return err
			}
			plan := &Plan{
				Msgs: []htx.Msg{&htx.StakeUpdateMsg{
					From: from, ValID: valID, NewAmount: strings.TrimSpace(newAmount),
					TxHash: txHash, LogIndex: logIndex,
					BlockNumber: blockNumber, Nonce: nonce,
				}},
				MsgShortType:  stakeUpdateMsgShort,
				SignerAddress: from,
			}
			return Execute(cmd, opts, mode, plan)
		},
	}
	RegisterFlags(cmd, opts, mode)
	f := cmd.Flags()
	f.StringVar(&fromFlag, "from-msg", "", "MsgStakeUpdate.from address (default: signer)")
	f.Uint64Var(&valID, "val-id", 0, "validator id")
	f.StringVar(&newAmount, "new-amount", "", "new stake amount (decimal string)")
	f.StringVar(&txHashHex, "tx-hash", "", "L1 tx hash (32 bytes hex)")
	f.Uint64Var(&logIndex, "log-index", 0, "L1 log index")
	f.Uint64Var(&blockNumber, "block-number", 0, "L1 block number")
	f.Uint64Var(&nonce, "nonce-l1", 0, "L1 stake nonce")
	return cmd
}
