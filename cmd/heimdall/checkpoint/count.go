package checkpoint

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/cmdutil"
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
	return pkg.NewGetCmd(cmdutil.Get{
		Use:         "count",
		Short:       "Print total acked checkpoint count.",
		Path:        "/checkpoints/count",
		Label:       "checkpoint count",
		FieldsUsage: "pluck one or more fields (repeatable, --json only)",
		RenderBody: func(cmd *cobra.Command, body []byte, _ render.Options) error {
			var resp countResponse
			if jerr := json.Unmarshal(body, &resp); jerr != nil {
				return fmt.Errorf("decoding checkpoint count: %w", jerr)
			}
			if resp.AckCount == "" {
				return fmt.Errorf("checkpoint count response missing ack_count (body=%q)", cmdutil.Truncate(body, 256))
			}
			_, err := fmt.Fprintln(cmd.OutOrStdout(), resp.AckCount)
			return err
		},
	})
}
