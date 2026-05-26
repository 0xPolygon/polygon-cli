package topup

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newProofCmd builds `topup proof <ADDR>` → GET
// /topup/account-proof/{address}. Requires L1 RPC on the Heimdall
// node; on gRPC code 13 the command surfaces an L1-not-configured
// hint on stderr before propagating the error.
func newProofCmd() *cobra.Command {
	var fields []string
	cmd := &cobra.Command{
		Use:   "proof <ADDR>",
		Short: "Fetch the Merkle proof for a dividend account.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			addr, err := normalizeAddress(args[0])
			if err != nil {
				return err
			}
			rest, cfg, err := newRESTClient(cmd)
			if err != nil {
				return err
			}
			body, status, err := rest.Get(cmd.Context(), fmt.Sprintf("/topup/account-proof/%s", addr), nil)
			opts := renderOpts(cmd, cfg, fields)
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
				return fmt.Errorf("topup account-proof failed: code=%d %s", gerr.Code, gerr.Message)
			}
			m, err := decodeJSONMap(body, "topup account-proof")
			if err != nil {
				return err
			}
			if opts.JSON {
				return render.RenderJSON(cmd.OutOrStdout(), m, opts)
			}
			// Unwrap the { "proof": {...} } envelope for KV.
			if inner, ok := m["proof"].(map[string]any); ok {
				return render.RenderKV(cmd.OutOrStdout(), inner, opts)
			}
			return render.RenderKV(cmd.OutOrStdout(), m, opts)
		},
	}
	cmd.Flags().StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable)")
	return cmd
}
