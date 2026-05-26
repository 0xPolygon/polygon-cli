package topup

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newSequenceCmd builds `topup sequence <TX_HASH> <LOG_INDEX>` → GET
// /topup/sequence?tx_hash=…&log_index=…. Requires L1 RPC on the
// Heimdall node; on gRPC code 13 the command surfaces an
// L1-not-configured hint on stderr before propagating the error.
func newSequenceCmd() *cobra.Command {
	var fields []string
	cmd := &cobra.Command{
		Use:   "sequence <TX_HASH> <LOG_INDEX>",
		Short: "Dedup sequence key for an L1 topup tx.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			hash, err := normalizeTxHash(args[0])
			if err != nil {
				return err
			}
			logIndex, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return &client.UsageError{Msg: fmt.Sprintf("log_index must be a non-negative integer, got %q", args[1])}
			}
			rest, cfg, err := newRESTClient(cmd)
			if err != nil {
				return err
			}
			q := url.Values{}
			q.Set("tx_hash", hash)
			q.Set("log_index", strconv.FormatUint(logIndex, 10))
			body, status, err := rest.Get(cmd.Context(), "/topup/sequence", q)
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
			var gerr gRPCErrorBody
			if jerr := json.Unmarshal(body, &gerr); jerr == nil && gerr.Code != 0 {
				if gerr.Code == gRPCCodeUnavailable {
					_ = render.WriteHint(cmd.ErrOrStderr(), render.HintL1NotConfigured, opts)
				}
				return fmt.Errorf("topup sequence failed: code=%d %s", gerr.Code, gerr.Message)
			}
			m, err := decodeJSONMap(body, "topup sequence")
			if err != nil {
				return err
			}
			if opts.JSON {
				return render.RenderJSON(cmd.OutOrStdout(), m, opts)
			}
			// Default text output: print the bare sequence if present.
			if seq, ok := m["sequence"].(string); ok && len(fields) == 0 {
				_, err = fmt.Fprintln(cmd.OutOrStdout(), seq)
				return err
			}
			return render.RenderKV(cmd.OutOrStdout(), m, opts)
		},
	}
	cmd.Flags().StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable)")
	return cmd
}
