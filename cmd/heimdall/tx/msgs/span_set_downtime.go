package msgs

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
	htx "github.com/0xPolygon/polygon-cli/internal/heimdall/tx"
)

// setProducerDowntimeMsgShort is not L1-mirroring.
const setProducerDowntimeMsgShort = "MsgSetProducerDowntime"

func init() {
	RegisterFactory("span-set-downtime", newSpanSetDowntimeCmd)
}

// newSpanSetDowntimeCmd builds `span-set-downtime`
// (MsgSetProducerDowntime).
func newSpanSetDowntimeCmd(mode Mode, globalFlags *config.Flags) *cobra.Command {
	opts := &TxOpts{Global: globalFlags}
	var (
		producer   string
		startBlock uint64
		endBlock   uint64
	)
	cmd := &cobra.Command{
		Use:   "span-set-downtime",
		Short: "Record producer downtime window (MsgSetProducerDowntime).",
		Long: strings.TrimSpace(`
Build, sign, and optionally broadcast a heimdallv2.bor.MsgSetProducerDowntime.

Validator-only. Downtime range is inclusive [start-block, end-block].
`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(producer) == "" {
				return &client.UsageError{Msg: "--producer is required"}
			}
			p, err := lowerEthAddress("producer", producer)
			if err != nil {
				return err
			}
			if endBlock < startBlock {
				return &client.UsageError{Msg: "--end-block must be >= --start-block"}
			}
			signerAddr := p
			if opts.From != "" || opts.Account != "" || opts.KeystoreFile != "" || opts.PrivateKey != "" || opts.Mnemonic != "" {
				signer, err := ResolveSigningKey(opts, cmd.InOrStdin())
				if err != nil {
					return err
				}
				signerAddr = strings.ToLower(signer.Address.Hex())
			}
			plan := &Plan{
				Msgs: []htx.Msg{&htx.SetProducerDowntimeMsg{
					Producer:   p,
					StartBlock: startBlock,
					EndBlock:   endBlock,
				}},
				MsgShortType:  setProducerDowntimeMsgShort,
				SignerAddress: signerAddr,
			}
			return Execute(cmd, opts, mode, plan)
		},
	}
	RegisterFlags(cmd, opts, mode)
	f := cmd.Flags()
	f.StringVar(&producer, "producer", "", "producer address whose downtime is being recorded")
	f.Uint64Var(&startBlock, "start-block", 0, "bor start block (inclusive)")
	f.Uint64Var(&endBlock, "end-block", 0, "bor end block (inclusive)")
	return cmd
}
