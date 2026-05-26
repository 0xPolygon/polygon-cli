package checkpoint

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// lastNoAckResponse is the shape of GET /checkpoints/last-no-ack.
// `last_no_ack_id` is unix seconds despite the name.
type lastNoAckResponse struct {
	LastNoAckID string `json:"last_no_ack_id"`
}

// newLastNoAckCmd builds `checkpoint last-no-ack` → GET
// /checkpoints/last-no-ack. Prints the unix-seconds value plus a
// human-readable age (via the shared timestamp annotator).
func newLastNoAckCmd() *cobra.Command {
	var fields []string
	cmd := &cobra.Command{
		Use:   "last-no-ack",
		Short: "Print the timestamp of the last no-ack.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			rest, cfg, err := newRESTClient(cmd)
			if err != nil {
				return err
			}
			body, status, err := rest.Get(cmd.Context(), "/checkpoints/last-no-ack", nil)
			if err != nil {
				return err
			}
			if status == 0 && body == nil {
				return nil
			}
			opts := renderOpts(cmd, cfg, fields)
			if opts.JSON {
				m, jerr := decodeJSONMap(body, "last-no-ack")
				if jerr != nil {
					return jerr
				}
				return render.RenderJSON(cmd.OutOrStdout(), m, opts)
			}
			var resp lastNoAckResponse
			if jerr := json.Unmarshal(body, &resp); jerr != nil {
				return fmt.Errorf("decoding last-no-ack: %w", jerr)
			}
			if resp.LastNoAckID == "" {
				return fmt.Errorf("last-no-ack response missing last_no_ack_id (body=%q)", truncate(body, 256))
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), render.AnnotateUnixSeconds(resp.LastNoAckID))
			return err
		},
	}
	cmd.Flags().StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable, --json only)")
	return cmd
}
