package milestone

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/cmdutil"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// countResponse is the shape of GET /milestones/count.
type countResponse struct {
	Count string `json:"count"`
}

// newCountCmd builds `milestone count` → GET /milestones/count.
// Default output is a bare integer (cheap liveness signal); --json
// emits the wrapper object.
func newCountCmd() *cobra.Command {
	return pkg.NewGetCmd(cmdutil.Get{
		Use:         "count",
		Short:       "Print total milestone count.",
		Path:        "/milestones/count",
		Label:       "milestone count",
		FieldsUsage: "pluck one or more fields (repeatable, --json only)",
		RenderBody: func(cmd *cobra.Command, body []byte, _ render.Options) error {
			var resp countResponse
			if jerr := json.Unmarshal(body, &resp); jerr != nil {
				return fmt.Errorf("decoding milestone count: %w", jerr)
			}
			if resp.Count == "" {
				return fmt.Errorf("milestone count response missing count (body=%q)", cmdutil.Truncate(body, 256))
			}
			_, err := fmt.Fprintln(cmd.OutOrStdout(), resp.Count)
			return err
		},
	})
}
