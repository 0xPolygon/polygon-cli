package clerk

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// listResponse is the shape of GET /clerk/event-records/list. The
// endpoint returns a bare `event_records` array; it does not emit a
// Cosmos-style `pagination` envelope because it is page-based (not
// Cosmos-paginated).
type listResponse struct {
	EventRecords []map[string]any `json:"event_records"`
}

// newListCmd builds `state-sync list [--page N] [--limit N]` → GET
// /clerk/event-records/list. Upstream is PAGE-BASED, not Cosmos
// pagination — the query params are bare `page` and `limit`, and the
// server rejects `page=0` with HTTP 400. We default --page to 1 so
// the bare `state-sync list` form works; --limit is surfaced via a
// hint when omitted, because the server's default behaviour returns
// a full page and is surprising to scripts.
func newListCmd() *cobra.Command {
	var (
		page   int
		limit  int
		fields []string
		base64 bool
	)
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Paginated event-record history (page-based).",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			rest, cfg, err := newRESTClient(cmd)
			if err != nil {
				return err
			}
			opts := renderOpts(cmd, cfg, fields, base64)
			// Surface the pagination-limit hint when --limit is not
			// explicitly set. The hint catalogue is generic; we emit it
			// here so users who hit /clerk/event-records/list without a
			// limit understand why the result set is shaped as it is.
			limitSet := cmd.Flags().Changed("limit")
			if !limitSet {
				_ = render.WriteHint(cmd.ErrOrStderr(), render.HintPaginationLimit, opts)
			}
			q := url.Values{}
			q.Set("page", strconv.Itoa(page))
			if limitSet {
				q.Set("limit", strconv.Itoa(limit))
			}
			body, status, err := rest.Get(cmd.Context(), "/clerk/event-records/list", q)
			if err != nil {
				return err
			}
			if status == 0 && body == nil {
				return nil
			}
			if opts.JSON {
				m, jerr := decodeJSONMap(body, "clerk list")
				if jerr != nil {
					return jerr
				}
				return render.RenderJSON(cmd.OutOrStdout(), m, opts)
			}
			var resp listResponse
			if jerr := json.Unmarshal(body, &resp); jerr != nil {
				return fmt.Errorf("decoding clerk list: %w", jerr)
			}
			if len(resp.EventRecords) == 0 {
				_, err = fmt.Fprintln(cmd.OutOrStdout(), "(no event records)")
				return err
			}
			return renderRecordTable(cmd, resp.EventRecords, opts)
		},
	}
	f := cmd.Flags()
	// Page defaults to 1: the upstream endpoint returns HTTP 400 on
	// page=0, so a zero default would turn the bare form into an error.
	f.IntVar(&page, "page", 1, "page number (1-indexed)")
	f.IntVar(&limit, "limit", 0, "maximum entries per page")
	f.BoolVar(&base64, "base64", false, "preserve raw base64 for `data` (default 0x-hex)")
	f.StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable, --json only)")
	return cmd
}

// renderRecordTable trims each row to the scalar summary fields used
// by the table output. The full `data` blob is too wide for a table,
// so we keep id / contract / tx_hash / log_index / record_time only.
func renderRecordTable(cmd *cobra.Command, records []map[string]any, opts render.Options) error {
	summary := make([]map[string]any, 0, len(records))
	for _, r := range records {
		row := map[string]any{
			"id":           r["id"],
			"contract":     r["contract"],
			"tx_hash":      r["tx_hash"],
			"log_index":    r["log_index"],
			"bor_chain_id": r["bor_chain_id"],
			"record_time":  r["record_time"],
		}
		summary = append(summary, row)
	}
	return render.RenderTable(cmd.OutOrStdout(), summary, opts)
}
