package checkpoint

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/cmdutil"
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
	return pkg.NewGetCmd(cmdutil.Get{
		Use:         "last-no-ack",
		Short:       "Print the timestamp of the last no-ack.",
		Path:        "/checkpoints/last-no-ack",
		Label:       "last-no-ack",
		FieldsUsage: "pluck one or more fields (repeatable, --json only)",
		RenderBody: func(cmd *cobra.Command, body []byte, _ render.Options) error {
			var resp lastNoAckResponse
			if jerr := json.Unmarshal(body, &resp); jerr != nil {
				return fmt.Errorf("decoding last-no-ack: %w", jerr)
			}
			if resp.LastNoAckID == "" {
				return fmt.Errorf("last-no-ack response missing last_no_ack_id (body=%q)", cmdutil.Truncate(body, 256))
			}
			_, err := fmt.Fprintln(cmd.OutOrStdout(), render.AnnotateUnixSeconds(resp.LastNoAckID))
			return err
		},
	})
}
