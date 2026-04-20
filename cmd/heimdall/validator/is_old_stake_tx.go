package validator

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newIsOldStakeTxCmd builds `validator is-old-stake-tx <TX_HASH>
// <LOG_INDEX>` → GET /stake/is-old-tx?tx_hash=…&log_index=…. On gRPC
// code 13 (or a transport-level `connection refused` / `dial tcp`)
// the command prints a human-friendly hint pointing at missing
// `eth_rpc_url` configuration before propagating the error.
func newIsOldStakeTxCmd() *cobra.Command {
	var fields []string
	cmd := &cobra.Command{
		Use:   "is-old-stake-tx <TX_HASH> <LOG_INDEX>",
		Short: "Check whether an L1 stake event was already replayed.",
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
			body, status, err := rest.Get(cmd.Context(), "/stake/is-old-tx", q)
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
				return fmt.Errorf("is-old-tx failed: code=%d %s", gerr.Code, gerr.Message)
			}
			m, err := decodeJSONMap(body, "is-old-tx")
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
