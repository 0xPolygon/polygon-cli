package validator

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// setResponse is the shape of GET /stake/validators-set.
type setResponse struct {
	ValidatorSet struct {
		Validators       []map[string]any `json:"validators"`
		Proposer         map[string]any   `json:"proposer"`
		TotalVotingPower string           `json:"total_voting_power"`
	} `json:"validator_set"`
}

// validSortOrders lists the supported --sort values. Kept exported-like
// as a package-level for easy validation in both the umbrella and the
// alias commands.
var validSortOrders = map[string]bool{
	"power":  true,
	"id":     true,
	"signer": true,
}

// newSetCmd builds `validator set [--sort …] [--limit N]` →
// GET /stake/validators-set. Defaults to power-desc ordering.
func newSetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Print the current validator set.",
		Args:  cobra.NoArgs,
		RunE:  runSet,
	}
	attachSetFlags(cmd.Flags())
	return cmd
}

// attachSetFlags binds --sort, --limit, and --field onto the given flag
// set. Called by both `validator set` and the top-level `validators`
// alias so both commands expose the same surface.
func attachSetFlags(f *pflag.FlagSet) {
	f.StringVar(&setFlags.sort, "sort", "power", "sort order: power|id|signer (power is descending)")
	f.IntVar(&setFlags.limit, "limit", 0, "truncate output to the first N validators (0 = unlimited)")
	f.StringArrayVarP(&setFlags.fields, "field", "f", nil, "pluck one or more fields (repeatable, --json only)")
}

// runSet is the shared RunE for `validator set` and `validators`.
func runSet(cmd *cobra.Command, _ []string) error {
	if !validSortOrders[setFlags.sort] {
		return &client.UsageError{Msg: fmt.Sprintf("--sort must be one of power|id|signer, got %q", setFlags.sort)}
	}
	if setFlags.limit < 0 {
		return &client.UsageError{Msg: fmt.Sprintf("--limit must be non-negative, got %d", setFlags.limit)}
	}
	rest, cfg, err := newRESTClient(cmd)
	if err != nil {
		return err
	}
	body, status, err := rest.Get(cmd.Context(), "/stake/validators-set", nil)
	if err != nil {
		return err
	}
	if status == 0 && body == nil {
		return nil
	}
	opts := renderOpts(cmd, cfg, setFlags.fields)

	var resp setResponse
	if jerr := json.Unmarshal(body, &resp); jerr != nil {
		return fmt.Errorf("decoding validator set: %w", jerr)
	}
	validators := resp.ValidatorSet.Validators
	sortValidators(validators, setFlags.sort)
	if setFlags.limit > 0 && setFlags.limit < len(validators) {
		validators = validators[:setFlags.limit]
	}

	if opts.JSON {
		// Preserve the envelope shape but apply sort/limit to the
		// validators array so scripts see the same view as the KV table.
		full, jerr := decodeJSONMap(body, "validator set")
		if jerr != nil {
			return jerr
		}
		if inner, ok := full["validator_set"].(map[string]any); ok {
			inner["validators"] = toAnySlice(validators)
		}
		return render.RenderJSON(cmd.OutOrStdout(), full, opts)
	}
	if err := render.RenderTable(cmd.OutOrStdout(), validators, opts); err != nil {
		return err
	}
	if resp.ValidatorSet.TotalVotingPower != "" {
		if _, werr := fmt.Fprintf(cmd.ErrOrStderr(), "total_voting_power=%s\n", resp.ValidatorSet.TotalVotingPower); werr != nil {
			return werr
		}
	}
	return nil
}

// sortValidators orders rows in-place by the requested key. `power` is
// descending (biggest first), `id` and `signer` ascending. The row
// fields `voting_power` and `val_id` arrive as JSON strings, so we
// parse them before comparing.
func sortValidators(rows []map[string]any, order string) {
	switch order {
	case "id":
		sort.SliceStable(rows, func(i, j int) bool {
			return intField(rows[i], "val_id") < intField(rows[j], "val_id")
		})
	case "signer":
		sort.SliceStable(rows, func(i, j int) bool {
			return strings.ToLower(stringField(rows[i], "signer")) <
				strings.ToLower(stringField(rows[j], "signer"))
		})
	case "power":
		sort.SliceStable(rows, func(i, j int) bool {
			return intField(rows[i], "voting_power") > intField(rows[j], "voting_power")
		})
	}
}

func stringField(m map[string]any, k string) string {
	if v, ok := m[k].(string); ok {
		return v
	}
	return ""
}

// intField parses a string or float64 field into an int64 for sort
// comparisons. Returns 0 on any parse failure so sort order remains
// total.
func intField(m map[string]any, k string) int64 {
	switch v := m[k].(type) {
	case string:
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0
		}
		return n
	case float64:
		return int64(v)
	}
	return 0
}

// toAnySlice widens a []map[string]any to []any so it fits back into a
// decoded JSON envelope as a slice value.
func toAnySlice(rows []map[string]any) []any {
	out := make([]any, len(rows))
	for i, r := range rows {
		out[i] = r
	}
	return out
}
