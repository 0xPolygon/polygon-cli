package span

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/cmdutil"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newVotesCmd builds `span votes [VAL_ID]`:
//   - no args:  GET /bor/producer-votes          (all voters)
//   - one arg:  GET /bor/producer-votes/{val_id} (single voter)
func newVotesCmd() *cobra.Command {
	return pkg.NewGetCmd(cmdutil.Get{
		Use:   "votes [VAL_ID]",
		Short: "Show producer-set votes (all or by voter id).",
		Args:  cobra.MaximumNArgs(1),
		Label: "producer votes",
		Build: func(cmd *cobra.Command, args []string) (string, url.Values, error) {
			if len(args) == 0 {
				return "/bor/producer-votes", nil, nil
			}
			id, err := parseSpanID("validator id", args[0])
			if err != nil {
				return "", nil, err
			}
			return fmt.Sprintf("/bor/producer-votes/%d", id), nil, nil
		},
		// Both shapes are best presented as JSON: the all-votes
		// response is a nested map keyed by validator id, and the
		// single-voter response is a small object with a list.
		Render: func(cmd *cobra.Command, m map[string]any, opts render.Options) error {
			return render.RenderJSON(cmd.OutOrStdout(), m, opts)
		},
	})
}
