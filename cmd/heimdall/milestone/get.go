package milestone

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newGetCmd builds `milestone get <NUMBER>` → GET /milestones/{number}.
// The same code path is re-entered from MilestoneCmd's RunE when a
// bare integer is provided (`milestone 11602043`).
func newGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <NUMBER>",
		Short: "Fetch one milestone by sequence number.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(cmd, args[0])
		},
	}
	return cmd
}

// runGet is the shared implementation used by both `milestone get
// <NUMBER>` and the bare-integer MilestoneCmd shorthand.
//
// On HTTP 404 the command fetches /milestones/count and, if the
// requested number exceeds the count (or is zero), prints a hint
// pointing at the valid range before returning the error. The hint
// travels on stderr so `--json` / `-f` output on stdout stays clean
// for scripts.
func runGet(cmd *cobra.Command, numArg string) error {
	number, err := strconv.ParseUint(numArg, 10, 64)
	if err != nil {
		return &client.UsageError{Msg: fmt.Sprintf("milestone number must be a positive integer, got %q", numArg)}
	}
	rest, cfg, err := newRESTClient(cmd)
	if err != nil {
		return err
	}
	body, status, err := rest.Get(cmd.Context(), fmt.Sprintf("/milestones/%d", number), nil)
	if err != nil {
		var hErr *client.HTTPError
		if errors.As(err, &hErr) && hErr.NotFound() {
			// 404 — try to enrich with a valid-range hint. If the
			// count lookup itself fails, return the original error.
			if count, cerr := fetchCount(cmd.Context(), rest); cerr == nil {
				opts := renderOpts(cmd, cfg, nil)
				if number == 0 || number > count {
					hint := render.Hint{
						Key:  "milestone-range",
						Body: fmt.Sprintf("hint: valid range is 1..%d", count),
					}
					_ = render.WriteHint(cmd.ErrOrStderr(), hint, opts)
				}
			}
		}
		return err
	}
	if status == 0 && body == nil {
		return nil
	}
	opts := renderOpts(cmd, cfg, nil)
	m, err := decodeJSONMap(body, "milestone")
	if err != nil {
		return err
	}
	if opts.JSON {
		// Splice `number` into the milestone envelope so --json
		// consumers can rely on it alongside `milestone_id`.
		if inner, ok := m["milestone"].(map[string]any); ok {
			inner["number"] = itoa(number)
		}
		return render.RenderJSON(cmd.OutOrStdout(), m, opts)
	}
	return renderMilestoneKV(cmd, m, opts, number)
}

// fetchCount issues a GET /milestones/count against rest and parses
// the integer out of the response body. Errors are propagated
// unchanged so the caller can decide whether to fall back.
func fetchCount(ctx context.Context, rest *client.RESTClient) (uint64, error) {
	body, _, err := rest.Get(ctx, "/milestones/count", nil)
	if err != nil {
		return 0, err
	}
	var resp countResponse
	if jerr := json.Unmarshal(body, &resp); jerr != nil {
		return 0, fmt.Errorf("decoding milestone count: %w", jerr)
	}
	if resp.Count == "" {
		return 0, fmt.Errorf("milestone count response missing count")
	}
	n, perr := strconv.ParseUint(resp.Count, 10, 64)
	if perr != nil {
		return 0, fmt.Errorf("milestone count not an unsigned integer: %w", perr)
	}
	return n, nil
}

// itoa formats an uint64 as a base-10 string. Kept as a helper so the
// only call site in latest.go / get.go doesn't have to import strconv.
func itoa(n uint64) string {
	return strconv.FormatUint(n, 10)
}
