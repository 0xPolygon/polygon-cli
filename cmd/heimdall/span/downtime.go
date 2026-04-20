package span

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newDowntimeCmd builds `span downtime <PRODUCER_ID>` → GET
// /bor/producers/planned-downtime/{id}. On HTTP 404 ("no planned
// downtime found") the command prints `none` and exits 0, because for
// operators the absence of a planned downtime record is a normal
// answer rather than a failure.
func newDowntimeCmd() *cobra.Command {
	var fields []string
	cmd := &cobra.Command{
		Use:   "downtime <PRODUCER_ID>",
		Short: "Show planned downtime for a producer (or `none`).",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseSpanID("producer id", args[0])
			if err != nil {
				return err
			}
			rest, cfg, err := newRESTClient(cmd)
			if err != nil {
				return err
			}
			body, status, err := rest.Get(cmd.Context(), fmt.Sprintf("/bor/producers/planned-downtime/%d", id), nil)
			if err != nil {
				var hErr *client.HTTPError
				if errors.As(err, &hErr) && hErr.NotFound() {
					_, werr := fmt.Fprintln(cmd.OutOrStdout(), "none")
					return werr
				}
				return err
			}
			if status == 0 && body == nil {
				return nil
			}
			opts := renderOpts(cmd, cfg, fields)
			m, err := decodeJSONMap(body, "planned downtime")
			if err != nil {
				return err
			}
			if opts.JSON {
				return render.RenderJSON(cmd.OutOrStdout(), m, opts)
			}
			if inner, ok := m["downtime_range"].(map[string]any); ok {
				return render.RenderKV(cmd.OutOrStdout(), inner, opts)
			}
			return render.RenderKV(cmd.OutOrStdout(), m, opts)
		},
	}
	cmd.Flags().StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable)")
	return cmd
}
