package checkpoint

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// listResponse is the shape of GET /checkpoints/list.
type listResponse struct {
	CheckpointList []map[string]any `json:"checkpoint_list"`
	Pagination     map[string]any   `json:"pagination"`
}

// newListCmd builds `checkpoint list [--limit N] [--reverse] [--page
// KEY]` → GET /checkpoints/list with Cosmos pagination parameters.
// Defaults: limit=10, reverse=true (requirements §3.2.1).
func newListCmd() *cobra.Command {
	var (
		limit   int
		reverse bool
		page    string
		fields  []string
	)
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Paginated checkpoint history.",
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
			body, status, err := rest.Get(cmd.Context(), "/checkpoints/list", q)
			if err != nil {
				return err
			}
			if status == 0 && body == nil {
				return nil
			}
			opts := renderOpts(cmd, cfg, fields)
			if opts.JSON {
				m, jerr := decodeJSONMap(body, "checkpoint list")
				if jerr != nil {
					return jerr
				}
				return render.RenderJSON(cmd.OutOrStdout(), m, opts)
			}
			var resp listResponse
			if jerr := json.Unmarshal(body, &resp); jerr != nil {
				return fmt.Errorf("decoding checkpoint list: %w", jerr)
			}
			if len(resp.CheckpointList) == 0 {
				_, err = fmt.Fprintln(cmd.OutOrStdout(), "(no checkpoints)")
				return err
			}
			// RenderTable applies byte-field normalization when opts.Raw
			// is false; pass rows through unchanged.
			if err := render.RenderTable(cmd.OutOrStdout(), resp.CheckpointList, opts); err != nil {
				return err
			}
			// Print the next_key (if any) on stderr so scripting flows
			// can capture only the table on stdout.
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
