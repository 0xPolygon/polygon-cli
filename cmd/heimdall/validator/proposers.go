package validator

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// proposersResponse is the shape of GET /stake/proposers/{N}.
type proposersResponse struct {
	Proposers []map[string]any `json:"proposers"`
}

// newProposersCmd builds `validator proposers [N]` → GET
// /stake/proposers/{N}. N defaults to 1 when omitted — matching the
// cast-style "no arg == single-shot" ergonomic without colliding with
// `validator proposer` (singular).
func newProposersCmd() *cobra.Command {
	var fields []string
	cmd := &cobra.Command{
		Use:   "proposers [N]",
		Short: "Show the next N proposers (default 1).",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			n := uint64(1)
			if len(args) == 1 {
				parsed, err := strconv.ParseUint(args[0], 10, 64)
				if err != nil || parsed == 0 {
					return &client.UsageError{Msg: fmt.Sprintf("proposers N must be a positive integer, got %q", args[0])}
				}
				n = parsed
			}
			rest, cfg, err := newRESTClient(cmd)
			if err != nil {
				return err
			}
			body, status, err := rest.Get(cmd.Context(), fmt.Sprintf("/stake/proposers/%d", n), nil)
			if err != nil {
				return err
			}
			if status == 0 && body == nil {
				return nil
			}
			opts := renderOpts(cmd, cfg, fields)
			if opts.JSON {
				m, jerr := decodeJSONMap(body, "proposers")
				if jerr != nil {
					return jerr
				}
				return render.RenderJSON(cmd.OutOrStdout(), m, opts)
			}
			var resp proposersResponse
			if jerr := json.Unmarshal(body, &resp); jerr != nil {
				return fmt.Errorf("decoding proposers: %w", jerr)
			}
			if len(resp.Proposers) == 0 {
				_, werr := fmt.Fprintln(cmd.OutOrStdout(), "(no proposers)")
				return werr
			}
			return render.RenderTable(cmd.OutOrStdout(), resp.Proposers, opts)
		},
	}
	cmd.Flags().StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable)")
	return cmd
}
