package msgs

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
	htx "github.com/0xPolygon/polygon-cli/internal/heimdall/tx"
)

// proposeSpanMsgShort is safe (validator-only but not L1-mirroring).
const proposeSpanMsgShort = "MsgProposeSpan"

func init() {
	RegisterFactory("span-propose", newSpanProposeCmd)
}

// newSpanProposeCmd builds `span-propose` (MsgProposeSpan). Required
// flags: span-id, start/end block, chain-id, seed (32 bytes).
func newSpanProposeCmd(mode Mode, globalFlags *config.Flags) *cobra.Command {
	opts := &TxOpts{Global: globalFlags}
	var (
		spanID      uint64
		proposer    string
		startBlock  uint64
		endBlock    uint64
		chainID     string
		seedHex     string
		seedAuthor  string
	)
	cmd := &cobra.Command{
		Use:   "span-propose",
		Short: "Propose a new bor span (MsgProposeSpan).",
		Long: strings.TrimSpace(`
Build, sign, and optionally broadcast a heimdallv2.bor.MsgProposeSpan.

Validator-only; the --force flag is not required because this msg is
not an L1-mirroring type, but the on-chain handler rejects non-
validator signers.
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
			if err := requireNonEmptyString("chain-id", chainID); err != nil {
				return err
			}
			if spanID == 0 {
				return &client.UsageError{Msg: "--span-id is required"}
			}
			seed, err := parseHexBytes("seed", seedHex, 32)
			if err != nil {
				return err
			}
			if len(seed) == 0 {
				return &client.UsageError{Msg: "--seed is required"}
			}
			author := strings.TrimSpace(seedAuthor)
			if author == "" {
				author = prop
			} else {
				p, err := lowerEthAddress("seed-author", author)
				if err != nil {
					return err
				}
				author = p
			}
			plan := &Plan{
				Msgs: []htx.Msg{&htx.ProposeSpanMsg{
					SpanID: spanID, Proposer: prop,
					StartBlock: startBlock, EndBlock: endBlock,
					ChainID: strings.TrimSpace(chainID),
					Seed:    seed, SeedAuthor: author,
				}},
				MsgShortType:  proposeSpanMsgShort,
				SignerAddress: prop,
			}
			return Execute(cmd, opts, mode, plan)
		},
	}
	RegisterFlags(cmd, opts, mode)
	f := cmd.Flags()
	f.Uint64Var(&spanID, "span-id", 0, "span id to propose")
	f.StringVar(&proposer, "proposer", "", "proposer address (default: signer)")
	f.Uint64Var(&startBlock, "start-block", 0, "bor start block (inclusive)")
	f.Uint64Var(&endBlock, "end-block", 0, "bor end block (inclusive)")
	f.StringVar(&chainID, "bor-chain-id", "", "bor chain id (e.g. 137)")
	f.StringVar(&seedHex, "seed", "", "32-byte seed hash (hex)")
	f.StringVar(&seedAuthor, "seed-author", "", "seed author address (default: proposer)")
	return cmd
}
