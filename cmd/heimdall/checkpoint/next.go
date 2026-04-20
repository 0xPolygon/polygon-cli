package checkpoint

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// gRPCErrorBody is the standard gRPC-gateway error envelope returned
// on 4xx/5xx from Heimdall REST. Only `code` and `message` are used
// here; `details` is ignored.
type gRPCErrorBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// gRPCCodeUnavailable is the L1-unreachable code surfaced by
// /checkpoints/prepare-next when the node lacks `eth_rpc_url`.
const gRPCCodeUnavailable = 13

// newNextCmd builds `checkpoint next` → GET /checkpoints/prepare-next.
// Requires the node to have L1 RPC configured; the special gRPC-code
// 13 case is surfaced with a hint about missing `eth_rpc_url`.
func newNextCmd() *cobra.Command {
	var fields []string
	cmd := &cobra.Command{
		Use:   "next",
		Short: "Compute the next checkpoint to propose.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			rest, cfg, err := newRESTClient(cmd)
			if err != nil {
				return err
			}
			body, status, err := rest.Get(cmd.Context(), "/checkpoints/prepare-next", nil)
			if err != nil {
				// Heimdall's REST gateway surfaces the gRPC code-13 as
				// an HTTP 5xx with a gRPC envelope body. We print the
				// hint first, then propagate the original error for
				// exit-code mapping.
				opts := renderOpts(cmd, cfg, fields)
				if isL1Unreachable(body, err) {
					_ = render.WriteHint(cmd.ErrOrStderr(), render.HintL1NotConfigured, opts)
				}
				return err
			}
			if status == 0 && body == nil {
				return nil
			}
			// 2xx responses may still be the gRPC envelope in some
			// failure modes on older Heimdalls — check once more.
			opts := renderOpts(cmd, cfg, fields)
			var gerr gRPCErrorBody
			if jerr := json.Unmarshal(body, &gerr); jerr == nil && gerr.Code != 0 {
				if gerr.Code == gRPCCodeUnavailable {
					_ = render.WriteHint(cmd.ErrOrStderr(), render.HintL1NotConfigured, opts)
				}
				return fmt.Errorf("prepare-next failed: code=%d %s", gerr.Code, gerr.Message)
			}
			m, err := decodeJSONMap(body, "prepare-next")
			if err != nil {
				return err
			}
			if opts.JSON {
				return render.RenderJSON(cmd.OutOrStdout(), m, opts)
			}
			return renderCheckpointKV(cmd, m, opts)
		},
	}
	cmd.Flags().StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable)")
	return cmd
}

// isL1Unreachable inspects a non-nil error from /checkpoints/prepare-next
// and returns true if its body (HTTPError.Body) advertises gRPC code 13.
func isL1Unreachable(body []byte, err error) bool {
	var hErr *client.HTTPError
	if errors.As(err, &hErr) && len(hErr.Body) > 0 {
		body = hErr.Body
	}
	if len(body) == 0 {
		return false
	}
	var g gRPCErrorBody
	if jerr := json.Unmarshal(body, &g); jerr != nil {
		return false
	}
	return g.Code == gRPCCodeUnavailable
}
