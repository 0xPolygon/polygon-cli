package span

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sort"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/cmdutil"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// scoresResponse is the shape of GET /bor/validator-performance-score.
// Upstream returns a map keyed by validator id with string-encoded
// integer scores.
type scoresResponse struct {
	ValidatorPerformanceScore map[string]string `json:"validator_performance_score"`
}

// newScoresCmd builds `span scores` → GET
// /bor/validator-performance-score. Prints the map sorted by score
// descending (tie-broken by ascending validator id for determinism),
// with one validator per line.
func newScoresCmd() *cobra.Command {
	return pkg.NewGetCmd(cmdutil.Get{
		Use:   "scores",
		Short: "Show validator performance scores (desc).",
		Path:  "/bor/validator-performance-score",
		Label: "validator performance scores",
		RenderBody: func(cmd *cobra.Command, body []byte, opts render.Options) error {
			var resp scoresResponse
			if jerr := json.Unmarshal(body, &resp); jerr != nil {
				return fmt.Errorf("decoding validator performance scores: %w", jerr)
			}
			rows := sortScoresDesc(resp.ValidatorPerformanceScore)
			if len(rows) == 0 {
				_, err := fmt.Fprintln(cmd.OutOrStdout(), "(no scores)")
				return err
			}
			table := make([]map[string]any, 0, len(rows))
			for _, r := range rows {
				table = append(table, map[string]any{
					"val_id": r.id,
					"score":  r.score,
				})
			}
			return render.RenderTable(cmd.OutOrStdout(), table, opts)
		},
	})
}

type scoreRow struct {
	id    string
	score string
}

// sortScoresDesc returns rows sorted by score (big.Int, descending)
// with tie-break on ascending numeric validator id for determinism.
// Non-numeric ids or scores compare as zero — we don't fail loudly
// because the upstream shape is stable.
func sortScoresDesc(in map[string]string) []scoreRow {
	rows := make([]scoreRow, 0, len(in))
	for k, v := range in {
		rows = append(rows, scoreRow{id: k, score: v})
	}
	sort.Slice(rows, func(i, j int) bool {
		si := parseBigInt(rows[i].score)
		sj := parseBigInt(rows[j].score)
		if c := sj.Cmp(si); c != 0 {
			return c < 0
		}
		// Tie-break on ascending validator id.
		ii := parseBigInt(rows[i].id)
		ij := parseBigInt(rows[j].id)
		return ii.Cmp(ij) < 0
	})
	return rows
}

func parseBigInt(s string) *big.Int {
	n := new(big.Int)
	if _, ok := n.SetString(s, 10); !ok {
		return new(big.Int)
	}
	return n
}
