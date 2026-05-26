package span

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newGetCmd builds `span get <ID>` → GET /bor/spans/{id}. The same
// code path is re-entered from SpanCmd's RunE when a bare integer is
// provided (`span 5982`).
func newGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <ID>",
		Short: "Fetch one span by id.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(cmd, args[0])
		},
	}
	return cmd
}

// runGet is the shared implementation used by both `span get <ID>`
// and the bare-integer SpanCmd shorthand.
func runGet(cmd *cobra.Command, idArg string) error {
	id, err := parseSpanID("span id", idArg)
	if err != nil {
		return err
	}
	rest, cfg, err := newRESTClient(cmd)
	if err != nil {
		return err
	}
	body, status, err := rest.Get(cmd.Context(), fmt.Sprintf("/bor/spans/%d", id), nil)
	if err != nil {
		return err
	}
	if status == 0 && body == nil {
		return nil
	}
	opts := renderOpts(cmd, cfg, nil)
	m, err := decodeJSONMap(body, "span")
	if err != nil {
		return err
	}
	if opts.JSON {
		return render.RenderJSON(cmd.OutOrStdout(), m, opts)
	}
	return renderSpanKV(cmd, m, opts)
}
