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

// newIsOldCmd builds `topup is-old <TX_HASH> <LOG_INDEX>` → GET
// /topup/is-old-tx?tx_hash=…&log_index=…. Requires L1 RPC on the
// Heimdall node; on gRPC code 13 the command surfaces an
// L1-not-configured hint on stderr before propagating the error.
//
// Default text output is a bare `true`/`false` so shell scripts can
// pipe it without parsing JSON.
func newIsOldCmd() *cobra.Command {
	var fields []string
	cmd := &cobra.Command{
		Use:   "is-old <TX_HASH> <LOG_INDEX>",
		Short: "Check whether an L1 topup tx was already processed.",
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
			body, status, err := rest.Get(cmd.Context(), "/topup/is-old-tx", q)
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
				return fmt.Errorf("topup is-old-tx failed: code=%d %s", gerr.Code, gerr.Message)
			}
			m, err := decodeJSONMap(body, "topup is-old-tx")
			if err != nil {
				return err
			}
			if opts.JSON {
				return render.RenderJSON(cmd.OutOrStdout(), m, opts)
			}
			// Default text output: bare bool when no --field filter.
			if v, ok := m["is_old"].(bool); ok && len(fields) == 0 {
				if v {
					_, err = fmt.Fprintln(cmd.OutOrStdout(), "true")
				} else {
					_, err = fmt.Fprintln(cmd.OutOrStdout(), "false")
				}
				return err
			}
			return render.RenderKV(cmd.OutOrStdout(), m, opts)
		},
	}
	cmd.Flags().StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable)")
	return cmd
}
