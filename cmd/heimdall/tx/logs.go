package tx

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// txSearchResult is the decoded shape of a CometBFT /tx_search reply.
type txSearchResult struct {
	Txs        []txSearchEntry `json:"txs"`
	TotalCount string          `json:"total_count"`
}

type txSearchEntry struct {
	Hash     string            `json:"hash"`
	Height   string            `json:"height"`
	Index    int               `json:"index"`
	TxResult cometTxResultBody `json:"tx_result"`
}

// newLogsCmd builds `logs <QUERY>`. Wraps CometBFT's /tx_search RPC
// for pagination + full-text querying over the tx index. Default
// output lists `<height>  <hash>`; --json emits the full envelope.
func newLogsCmd() *cobra.Command {
	var limit, page int
	var fields []string
	cmd := &cobra.Command{
		Use:   "logs <QUERY>",
		Short: "Query the CometBFT tx index.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if limit <= 0 {
				limit = 30
			}
			if page <= 0 {
				page = 1
			}
			rpc, cfg, err := newRPCClient(cmd)
			if err != nil {
				return err
			}
			params := map[string]any{
				"query":    args[0],
				"prove":    false,
				"page":     fmt.Sprintf("%d", page),
				"per_page": fmt.Sprintf("%d", limit),
				"order_by": "desc",
			}
			raw, err := rpc.Call(cmd.Context(), "tx_search", params)
			if err != nil {
				return fmt.Errorf("tx_search: %w", err)
			}
			if raw == nil {
				return nil // --curl
			}
			var res txSearchResult
			if err := json.Unmarshal(raw, &res); err != nil {
				return fmt.Errorf("decoding tx_search: %w", err)
			}

			opts := renderOpts(cmd, cfg, fields)
			if opts.JSON {
				var generic any
				if err := json.Unmarshal(raw, &generic); err != nil {
					return fmt.Errorf("decoding tx_search for json: %w", err)
				}
				return render.RenderJSON(cmd.OutOrStdout(), generic, opts)
			}

			if len(res.Txs) == 0 {
				_, err := fmt.Fprintln(cmd.OutOrStdout(), "(no matches)")
				return err
			}
			for _, t := range res.Txs {
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s  0x%s\n", t.Height, t.Hash); err != nil {
					return err
				}
			}
			if _, err := fmt.Fprintf(cmd.ErrOrStderr(), "total_count=%s  page=%d  per_page=%d\n", res.TotalCount, page, limit); err != nil {
				return err
			}
			return nil
		},
	}
	f := cmd.Flags()
	f.IntVar(&limit, "limit", 30, "max results per page")
	f.IntVar(&page, "page", 1, "page number (1-indexed)")
	f.StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable, --json only)")
	return cmd
}
