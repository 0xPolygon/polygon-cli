package span

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// listResponse is the shape of GET /bor/spans/list. The table
// renderer only needs a handful of scalar fields, so we keep the rest
// nested as JSON in the original map.
type listResponse struct {
	SpanList   []map[string]any `json:"span_list"`
	Pagination map[string]any   `json:"pagination"`
}

// newListCmd builds `span list [--limit N] [--reverse] [--page KEY]`
// → GET /bor/spans/list with Cosmos pagination parameters. Defaults:
// limit=10, reverse=true (newest-first, mirroring checkpoint list).
func newListCmd() *cobra.Command {
	var (
		limit   int
		reverse bool
		page    string
		fields  []string
	)
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Paginated span history.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if limit <= 0 {
				limit = 10
			}
			rest, cfg, err := newRESTClient(cmd)
			if err != nil {
				return err
			}
			q := url.Values{}
			q.Set("pagination.limit", strconv.Itoa(limit))
			q.Set("pagination.reverse", strconv.FormatBool(reverse))
			if page != "" {
				q.Set("pagination.key", page)
			}
			body, status, err := rest.Get(cmd.Context(), "/bor/spans/list", q)
			if err != nil {
				return err
			}
			if status == 0 && body == nil {
				return nil
			}
			opts := renderOpts(cmd, cfg, fields)
			if opts.JSON {
				m, jerr := decodeJSONMap(body, "span list")
				if jerr != nil {
					return jerr
				}
				return render.RenderJSON(cmd.OutOrStdout(), m, opts)
			}
			var resp listResponse
			if jerr := json.Unmarshal(body, &resp); jerr != nil {
				return fmt.Errorf("decoding span list: %w", jerr)
			}
			if len(resp.SpanList) == 0 {
				_, err = fmt.Fprintln(cmd.OutOrStdout(), "(no spans)")
				return err
			}
			// Trim each row to the scalar summary fields — the full
			// validator set and producer list are too wide for a table.
			summary := make([]map[string]any, 0, len(resp.SpanList))
			for _, s := range resp.SpanList {
				row := map[string]any{
					"id":           s["id"],
					"start_block":  s["start_block"],
					"end_block":    s["end_block"],
					"bor_chain_id": s["bor_chain_id"],
					"producers":    producerCount(s),
				}
				summary = append(summary, row)
			}
			if err := render.RenderTable(cmd.OutOrStdout(), summary, opts); err != nil {
				return err
			}
			if nk, ok := resp.Pagination["next_key"].(string); ok && nk != "" {
				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "next_key=%s\n", nk)
			}
			return nil
		},
	}
	f := cmd.Flags()
	f.IntVar(&limit, "limit", 10, "maximum entries to return")
	f.BoolVar(&reverse, "reverse", true, "newest-first ordering")
	f.StringVar(&page, "page", "", "pagination key from a previous response")
	f.StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable, --json only)")
	return cmd
}

// producerCount returns the size of span["selected_producers"] or 0 if
// missing/malformed. Used only for summary-table rendering.
func producerCount(span map[string]any) int {
	ps, ok := span["selected_producers"].([]any)
	if !ok {
		return 0
	}
	return len(ps)
}
