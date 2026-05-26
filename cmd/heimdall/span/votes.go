package span

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newVotesCmd builds `span votes [VAL_ID]`:
//   - no args:  GET /bor/producer-votes          (all voters)
//   - one arg:  GET /bor/producer-votes/{val_id} (single voter)
func newVotesCmd() *cobra.Command {
	var fields []string
	cmd := &cobra.Command{
		Use:   "votes [VAL_ID]",
		Short: "Show producer-set votes (all or by voter id).",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rest, cfg, err := newRESTClient(cmd)
			if err != nil {
				return err
			}
			var path string
			if len(args) == 0 {
				path = "/bor/producer-votes"
			} else {
				id, perr := parseSpanID("validator id", args[0])
				if perr != nil {
					return perr
				}
				path = fmt.Sprintf("/bor/producer-votes/%d", id)
			}
			body, status, err := rest.Get(cmd.Context(), path, nil)
			if err != nil {
				return err
			}
			if status == 0 && body == nil {
				return nil
			}
			opts := renderOpts(cmd, cfg, fields)
			m, err := decodeJSONMap(body, "producer votes")
			if err != nil {
				return err
			}
			// Both shapes are best presented as JSON: the all-votes
			// response is a nested map keyed by validator id, and the
			// single-voter response is a small object with a list.
			return render.RenderJSON(cmd.OutOrStdout(), m, opts)
		},
	}
	cmd.Flags().StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable)")
	return cmd
}
