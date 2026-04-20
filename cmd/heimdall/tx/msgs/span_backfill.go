package msgs

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
	htx "github.com/0xPolygon/polygon-cli/internal/heimdall/tx"
)

// backfillSpansMsgShort is not L1-mirroring.
const backfillSpansMsgShort = "MsgBackfillSpans"

func init() {
	RegisterFactory("span-backfill", newSpanBackfillCmd)
}

// newSpanBackfillCmd builds `span-backfill` (MsgBackfillSpans).
func newSpanBackfillCmd(mode Mode, globalFlags *config.Flags) *cobra.Command {
	opts := &TxOpts{Global: globalFlags}
	var (
		proposer        string
		chainID         string
		latestSpanID    uint64
		latestBorSpanID uint64
	)
	cmd := &cobra.Command{
		Use:   "span-backfill",
		Short: "Trigger span backfill (MsgBackfillSpans).",
		Long: strings.TrimSpace(`
Build, sign, and optionally broadcast a heimdallv2.bor.MsgBackfillSpans.

Requests Heimdall to resync spans when the chain's view of the latest
span drifts from bor's. Validator-only.
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
			if err := requireNonEmptyString("bor-chain-id", chainID); err != nil {
				return err
			}
			plan := &Plan{
				Msgs: []htx.Msg{&htx.BackfillSpansMsg{
					Proposer: prop, ChainID: strings.TrimSpace(chainID),
					LatestSpanID: latestSpanID, LatestBorSpanID: latestBorSpanID,
				}},
				MsgShortType:  backfillSpansMsgShort,
				SignerAddress: prop,
			}
			return Execute(cmd, opts, mode, plan)
		},
	}
	RegisterFlags(cmd, opts, mode)
	f := cmd.Flags()
	f.StringVar(&proposer, "proposer", "", "proposer address (default: signer)")
	f.StringVar(&chainID, "bor-chain-id", "", "bor chain id")
	f.Uint64Var(&latestSpanID, "latest-span-id", 0, "latest heimdall span id")
	f.Uint64Var(&latestBorSpanID, "latest-bor-span-id", 0, "latest bor span id")
	return cmd
}
