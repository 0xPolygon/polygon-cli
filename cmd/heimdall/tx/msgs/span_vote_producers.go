package msgs

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
	htx "github.com/0xPolygon/polygon-cli/internal/heimdall/tx"
)

// voteProducersMsgShort is not L1-mirroring.
const voteProducersMsgShort = "MsgVoteProducers"

func init() {
	RegisterFactory("span-vote-producers", newSpanVoteProducersCmd)
}

// newSpanVoteProducersCmd builds `span-vote-producers`
// (MsgVoteProducers). Votes are passed as a comma-separated list of
// validator IDs.
func newSpanVoteProducersCmd(mode Mode, globalFlags *config.Flags) *cobra.Command {
	opts := &TxOpts{Global: globalFlags}
	var (
		voter    string
		voterID  uint64
		votesCSV string
	)
	cmd := &cobra.Command{
		Use:   "span-vote-producers",
		Short: "Vote for producers in the next span (MsgVoteProducers).",
		Long: strings.TrimSpace(`
Build, sign, and optionally broadcast a heimdallv2.bor.MsgVoteProducers.

--votes is a comma-separated list of validator IDs (uint64) to vote
for; order matters on-chain. Validator-only.
`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			v := strings.TrimSpace(voter)
			if v == "" {
				signer, err := ResolveSigningKey(opts, cmd.InOrStdin())
				if err != nil {
					return err
				}
				v = strings.ToLower(signer.Address.Hex())
			} else {
				p, err := lowerEthAddress("voter", v)
				if err != nil {
					return err
				}
				v = p
			}
			if strings.TrimSpace(votesCSV) == "" {
				return &client.UsageError{Msg: "--votes is required (comma-separated validator ids)"}
			}
			votes, err := parseUint64CSV(votesCSV)
			if err != nil {
				return err
			}
			if voterID == 0 {
				return &client.UsageError{Msg: "--voter-id is required"}
			}
			plan := &Plan{
				Msgs: []htx.Msg{&htx.VoteProducersMsg{
					Voter: v, VoterID: voterID, Votes: votes,
				}},
				MsgShortType:  voteProducersMsgShort,
				SignerAddress: v,
			}
			return Execute(cmd, opts, mode, plan)
		},
	}
	RegisterFlags(cmd, opts, mode)
	f := cmd.Flags()
	f.StringVar(&voter, "voter", "", "voter address (default: signer)")
	f.Uint64Var(&voterID, "voter-id", 0, "voter's validator id")
	f.StringVar(&votesCSV, "votes", "", "comma-separated validator ids to vote for")
	return cmd
}

// parseUint64CSV parses a non-empty comma-separated list of uint64s.
func parseUint64CSV(s string) ([]uint64, error) {
	parts := strings.Split(strings.TrimSpace(s), ",")
	out := make([]uint64, 0, len(parts))
	for i, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			return nil, &client.UsageError{Msg: fmt.Sprintf("--votes entry #%d is empty", i+1)}
		}
		n, err := strconv.ParseUint(p, 10, 64)
		if err != nil {
			return nil, &client.UsageError{Msg: fmt.Sprintf("--votes entry #%d %q: %v", i+1, p, err)}
		}
		out = append(out, n)
	}
	return out, nil
}
