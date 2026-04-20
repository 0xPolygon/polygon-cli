package clerk

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newLatestIDCmd builds `state-sync latest-id` → GET
// /clerk/event-records/latest-id. Requires L1 RPC on the node; on gRPC
// code 13 (or a transport-level `connection refused`) the command
// prints a human-friendly hint pointing at missing `eth_rpc_url`
// configuration before propagating the error.
func newLatestIDCmd() *cobra.Command {
	var fields []string
	cmd := &cobra.Command{
		Use:   "latest-id",
		Short: "Latest L1 state-sync counter (requires eth_rpc_url on the node).",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			rest, cfg, err := newRESTClient(cmd)
			if err != nil {
				return err
			}
			body, status, err := rest.Get(cmd.Context(), "/clerk/event-records/latest-id", nil)
			opts := renderOpts(cmd, cfg, fields, false)
			if err != nil {
				if isL1Unreachable(body, err) {
					_ = render.WriteHint(cmd.ErrOrStderr(), render.HintL1NotConfigured, opts)
				}
				return err
			}
			if status == 0 && body == nil {
				return nil
			}
			// Some Heimdalls surface gRPC code 13 on 2xx with the
			// envelope body; check once more.
			var gerr gRPCErrorBody
			if jerr := json.Unmarshal(body, &gerr); jerr == nil && gerr.Code != 0 {
				if gerr.Code == gRPCCodeUnavailable {
					_ = render.WriteHint(cmd.ErrOrStderr(), render.HintL1NotConfigured, opts)
				}
				return fmt.Errorf("clerk latest-id failed: code=%d %s", gerr.Code, gerr.Message)
			}
			m, err := decodeJSONMap(body, "clerk latest-id")
			if err != nil {
				return err
			}
			if opts.JSON {
				return render.RenderJSON(cmd.OutOrStdout(), m, opts)
			}
			return render.RenderKV(cmd.OutOrStdout(), m, opts)
		},
	}
	cmd.Flags().StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable)")
	return cmd
}
