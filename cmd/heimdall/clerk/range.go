package clerk

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// timeRangeResponse is the shape of GET /clerk/time. The endpoint
// returns a bare `event_records` array; although it is routed through
// the cosmos-pagination middleware on the server, the response does
// not carry a pagination envelope.
type timeRangeResponse struct {
	EventRecords []map[string]any `json:"event_records"`
}

// newRangeCmd builds `state-sync range --from-id ID [--to-time T]
// [--limit N]` → GET /clerk/time. --from-id is required. The server
// rejects a fully-unset query with "pagination request is empty", so
// the command refuses to call out without at least --from-id.
//
// Unlike `state-sync list` (page-based), this endpoint is wired
// through cosmos-pagination on the server side and accepts
// `pagination.limit`. We surface that as a plain `--limit`.
func newRangeCmd() *cobra.Command {
	var (
		fromID uint64
		toTime string
		limit  int
		fields []string
		base64 bool
	)
	cmd := &cobra.Command{
		Use:   "range",
		Short: "Event-records since an id, optionally bounded by a timestamp.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("from-id") {
				return &client.UsageError{Msg: "--from-id is required"}
			}
			rest, cfg, err := newRESTClient(cmd)
			if err != nil {
				return err
			}
			q := url.Values{}
			q.Set("from_id", strconv.FormatUint(fromID, 10))
			if toTime != "" {
				q.Set("to_time", toTime)
			}
			if cmd.Flags().Changed("limit") {
				q.Set("pagination.limit", strconv.Itoa(limit))
			}
			body, status, err := rest.Get(cmd.Context(), "/clerk/time", q)
			if err != nil {
				return err
			}
			if status == 0 && body == nil {
				return nil
			}
			opts := renderOpts(cmd, cfg, fields, base64)
			if opts.JSON {
				m, jerr := decodeJSONMap(body, "clerk range")
				if jerr != nil {
					return jerr
				}
				return render.RenderJSON(cmd.OutOrStdout(), m, opts)
			}
			var resp timeRangeResponse
			if jerr := json.Unmarshal(body, &resp); jerr != nil {
				return fmt.Errorf("decoding clerk range: %w", jerr)
			}
			if len(resp.EventRecords) == 0 {
				_, err = fmt.Fprintln(cmd.OutOrStdout(), "(no event records)")
				return err
			}
			return renderRecordTable(cmd, resp.EventRecords, opts)
		},
	}
	f := cmd.Flags()
	f.Uint64Var(&fromID, "from-id", 0, "lowest event-record id to return (required)")
	f.StringVar(&toTime, "to-time", "", "RFC3339 upper bound on record_time")
	f.IntVar(&limit, "limit", 0, "maximum entries to return")
	f.BoolVar(&base64, "base64", false, "preserve raw base64 for `data` (default 0x-hex)")
	f.StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable, --json only)")
	return cmd
}
