package checkpoint

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// countResponse is the shape of GET /checkpoints/count.
type countResponse struct {
	AckCount string `json:"ack_count"`
}

// newCountCmd builds `checkpoint count` → GET /checkpoints/count.
// Default output is a bare integer (cheap liveness signal); --json
// emits the wrapper object.
func newCountCmd() *cobra.Command {
	var fields []string
	cmd := &cobra.Command{
		Use:   "count",
		Short: "Print total acked checkpoint count.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			rest, cfg, err := newRESTClient(cmd)
			if err != nil {
				return err
			}
			body, status, err := rest.Get(cmd.Context(), "/checkpoints/count", nil)
			if err != nil {
				return err
			}
			if status == 0 && body == nil {
				return nil
			}
			opts := renderOpts(cmd, cfg, fields)
			if opts.JSON {
				m, jerr := decodeJSONMap(body, "checkpoint count")
				if jerr != nil {
					return jerr
				}
				return render.RenderJSON(cmd.OutOrStdout(), m, opts)
			}
			var resp countResponse
			if jerr := json.Unmarshal(body, &resp); jerr != nil {
				return fmt.Errorf("decoding checkpoint count: %w", jerr)
			}
			if resp.AckCount == "" {
				return fmt.Errorf("checkpoint count response missing ack_count (body=%q)", truncate(body, 256))
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), resp.AckCount)
			return err
		},
	}
	cmd.Flags().StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable, --json only)")
	return cmd
}
