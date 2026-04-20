package clerk

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// countResponse is the shape of GET /clerk/event-records/count.
type countResponse struct {
	Count string `json:"count"`
}

// newCountCmd builds `state-sync count` → GET /clerk/event-records/count.
// Default output is a bare integer (cheap liveness signal); --json
// emits the wrapper object.
func newCountCmd() *cobra.Command {
	var fields []string
	cmd := &cobra.Command{
		Use:   "count",
		Short: "Print total state-sync (clerk) event-record count.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			rest, cfg, err := newRESTClient(cmd)
			if err != nil {
				return err
			}
			body, status, err := rest.Get(cmd.Context(), "/clerk/event-records/count", nil)
			if err != nil {
				return err
			}
			if status == 0 && body == nil {
				return nil
			}
			opts := renderOpts(cmd, cfg, fields, false)
			if opts.JSON {
				m, jerr := decodeJSONMap(body, "clerk count")
				if jerr != nil {
					return jerr
				}
				return render.RenderJSON(cmd.OutOrStdout(), m, opts)
			}
			var resp countResponse
			if jerr := json.Unmarshal(body, &resp); jerr != nil {
				return fmt.Errorf("decoding clerk count: %w", jerr)
			}
			if resp.Count == "" {
				return fmt.Errorf("clerk count response missing count (body=%q)", truncate(body, 256))
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), resp.Count)
			return err
		},
	}
	cmd.Flags().StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable, --json only)")
	return cmd
}
