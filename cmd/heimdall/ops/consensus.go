package ops

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// cometDumpConsensus is the minimal shape we peel off /dump_consensus_state
// for the default summary. Everything we don't summarise stays in the
// raw payload and is reachable via --json / --field.
type cometDumpConsensus struct {
	RoundState struct {
		Height     string `json:"height"`
		Round      int    `json:"round"`
		Step       int    `json:"step"`
		StartTime  string `json:"start_time"`
		CommitTime string `json:"commit_time"`
		Votes      []struct {
			Round              int      `json:"round"`
			Prevotes           []string `json:"prevotes"`
			PrevotesBitArray   string   `json:"prevotes_bit_array"`
			Precommits         []string `json:"precommits"`
			PrecommitsBitArray string   `json:"precommits_bit_array"`
		} `json:"votes"`
		Validators struct {
			Proposer struct {
				Address string `json:"address"`
			} `json:"proposer"`
		} `json:"validators"`
	} `json:"round_state"`
}

// newConsensusCmd builds `ops consensus`. Default output is a concise
// KV summary (height/round/step + per-round vote bit-arrays). --json
// dumps the full /dump_consensus_state payload.
//
// The full payload is dense and expensive to generate on a busy node;
// we surface a stderr warning whenever we issue the default summary or
// --json call so operators don't blow up their peer's load.
func newConsensusCmd() *cobra.Command {
	var fields []string
	cmd := &cobra.Command{
		Use:   "consensus",
		Short: "Summarise CometBFT /dump_consensus_state.",
		Long: `Summarise the CometBFT consensus round state (height, round, step,
proposer, per-round vote bit-arrays).

WARNING: /dump_consensus_state is an expensive RPC on a busy node and
is frequently disabled via RPC.EnableConsensusEndpoints=false in
config.toml. If the node rejects the call with a method-not-enabled
error, that's a node-side configuration, not a bug in polycli.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			rpc, cfg, err := newRPCClient(cmd)
			if err != nil {
				return err
			}
			// Warn before the call so the operator sees the warning
			// even when the call hangs or fails slowly.
			if _, werr := fmt.Fprintln(cmd.ErrOrStderr(),
				"warning: /dump_consensus_state is expensive; avoid on a node under load"); werr != nil {
				return werr
			}
			raw, err := callEmpty(cmd.Context(), rpc, "dump_consensus_state")
			if err != nil {
				return err
			}
			if raw == nil {
				return nil // --curl
			}
			opts := renderOpts(cmd, cfg, fields)
			if opts.JSON {
				generic, derr := decodeGeneric(raw)
				if derr != nil {
					return derr
				}
				return render.RenderJSON(cmd.OutOrStdout(), generic, opts)
			}
			var dc cometDumpConsensus
			if err := json.Unmarshal(raw, &dc); err != nil {
				return fmt.Errorf("decoding dump_consensus_state: %w", err)
			}
			rs := dc.RoundState
			out := map[string]any{
				"height":           rs.Height,
				"round":            rs.Round,
				"step":             rs.Step,
				"start_time":       rs.StartTime,
				"commit_time":      rs.CommitTime,
				"proposer_address": "0x" + rs.Validators.Proposer.Address,
				"num_vote_rounds":  len(rs.Votes),
			}
			if len(rs.Votes) > 0 {
				latest := rs.Votes[len(rs.Votes)-1]
				out["prevotes_bit_array"] = latest.PrevotesBitArray
				out["precommits_bit_array"] = latest.PrecommitsBitArray
			}
			return render.RenderKV(cmd.OutOrStdout(), out, opts)
		},
	}
	f := cmd.Flags()
	f.StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable)")
	return cmd
}
