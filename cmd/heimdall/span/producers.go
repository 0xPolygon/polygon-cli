package span

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newProducersCmd builds `span producers <ID>` as a derived subcommand:
// fetch the span and print only the selected_producers[] array.
func newProducersCmd() *cobra.Command {
	var fields []string
	cmd := &cobra.Command{
		Use:   "producers <ID>",
		Short: "List selected producers for a span.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseSpanID("span id", args[0])
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
			opts := renderOpts(cmd, cfg, fields)
			m, err := decodeJSONMap(body, "span")
			if err != nil {
				return err
			}
			inner, ok := m["span"].(map[string]any)
			if !ok {
				return fmt.Errorf("unexpected span response: missing \"span\" envelope")
			}
			producers, _ := inner["selected_producers"].([]any)
			if opts.JSON {
				return render.RenderJSON(cmd.OutOrStdout(), producers, opts)
			}
			if len(producers) == 0 {
				_, err := fmt.Fprintln(cmd.OutOrStdout(), "(no producers)")
				return err
			}
			rows := make([]map[string]any, 0, len(producers))
			for _, p := range producers {
				if pm, ok := p.(map[string]any); ok {
					rows = append(rows, pm)
				}
			}
			return render.RenderTable(cmd.OutOrStdout(), rows, opts)
		},
	}
	cmd.Flags().StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable, --json only)")
	return cmd
}
