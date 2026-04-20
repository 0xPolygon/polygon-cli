package msgs

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
	htx "github.com/0xPolygon/polygon-cli/internal/heimdall/tx"
)

// topupMsgShort triggers the L1-mirroring guard.
const topupMsgShort = "MsgTopupTx"

func init() {
	RegisterFactory("topup", newTopupCmd)
}

// newTopupCmd builds the `topup` subcommand (MsgTopupTx). The bridge
// produces these after observing an L1 Topup event; manual use needs
// --force.
func newTopupCmd(mode Mode, globalFlags *config.Flags) *cobra.Command {
	opts := &TxOpts{Global: globalFlags}
	var (
		proposer    string
		user        string
		fee         string
		txHashHex   string
		logIndex    uint64
		blockNumber uint64
	)
	cmd := &cobra.Command{
		Use:   "topup",
		Short: "Credit validator fee balance (MsgTopupTx, L1-mirroring).",
		Long: strings.TrimSpace(`
Build, sign, and optionally broadcast a heimdallv2.topup.MsgTopupTx.

MsgTopupTx is produced by the bridge after observing an L1 event;
manual use is a replay. Refuses without --force.
`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			prop := strings.TrimSpace(proposer)
			if prop == "" {
				signer, err := ResolveSigningKey(opts, cmd.InOrStdin())
				if err != nil {
					return err
				}
				prop = strings.ToLower(signer.Address.Hex())
			} else {
				p, err := lowerEthAddress("proposer", prop)
				if err != nil {
					return err
				}
				prop = p
			}
			if err := requireNonEmptyString("user", user); err != nil {
				return err
			}
			u, err := lowerEthAddress("user", user)
			if err != nil {
				return err
			}
			if strings.TrimSpace(fee) == "" {
				return &client.UsageError{Msg: "--fee is required"}
			}
			txHash, err := parseHexBytes("tx-hash", txHashHex, 32)
			if err != nil {
				return err
			}
			plan := &Plan{
				Msgs: []htx.Msg{&htx.TopupMsg{
					Proposer: prop, User: u, Fee: strings.TrimSpace(fee),
					TxHash: txHash, LogIndex: logIndex, BlockNumber: blockNumber,
				}},
				MsgShortType:  topupMsgShort,
				SignerAddress: prop,
			}
			return Execute(cmd, opts, mode, plan)
		},
	}
	RegisterFlags(cmd, opts, mode)
	f := cmd.Flags()
	f.StringVar(&proposer, "proposer", "", "proposer address (default: signer)")
	f.StringVar(&user, "user", "", "user address being topped up")
	f.StringVar(&fee, "fee-amount", "", "topup fee amount (decimal string)")
	f.StringVar(&txHashHex, "tx-hash", "", "L1 transaction hash (32 bytes hex)")
	f.Uint64Var(&logIndex, "log-index", 0, "L1 log index")
	f.Uint64Var(&blockNumber, "block-number", 0, "L1 block number")
	return cmd
}
