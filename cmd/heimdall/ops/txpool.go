package ops

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// cometUnconfirmedSummary is the /num_unconfirmed_txs payload.
type cometUnconfirmedSummary struct {
	NTxs       string `json:"n_txs"`
	Total      string `json:"total"`
	TotalBytes string `json:"total_bytes"`
}

// cometUnconfirmedTxs is the /unconfirmed_txs payload. Txs is a list
// of base64 strings; each string is a CometBFT-encoded raw tx.
type cometUnconfirmedTxs struct {
	NTxs       string   `json:"n_txs"`
	Total      string   `json:"total"`
	TotalBytes string   `json:"total_bytes"`
	Txs        []string `json:"txs"`
}

// newTxPoolCmd builds `ops tx-pool`. Default output is the mempool
// summary (n_txs, total_bytes). --list additionally fetches up to
// --limit txs, decoded to `0x<sha256>` hashes one per line.
func newTxPoolCmd() *cobra.Command {
	var list bool
	var limit int
	var fields []string
	cmd := &cobra.Command{
		Use:   "tx-pool",
		Short: "Show CometBFT mempool size (--list for hashes).",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if limit <= 0 {
				return &client.UsageError{Msg: fmt.Sprintf("--limit must be positive, got %d", limit)}
			}
			rpc, cfg, err := newRPCClient(cmd)
			if err != nil {
				return err
			}
			opts := renderOpts(cmd, cfg, fields)

			if !list {
				raw, err := callEmpty(cmd.Context(), rpc, "num_unconfirmed_txs")
				if err != nil {
					return err
				}
				if raw == nil {
					return nil // --curl
				}
				if opts.JSON {
					generic, derr := decodeGeneric(raw)
					if derr != nil {
						return derr
					}
					return render.RenderJSON(cmd.OutOrStdout(), generic, opts)
				}
				var s cometUnconfirmedSummary
				if err := json.Unmarshal(raw, &s); err != nil {
					return fmt.Errorf("decoding num_unconfirmed_txs: %w", err)
				}
				out := map[string]any{
					"n_txs":       s.NTxs,
					"total":       s.Total,
					"total_bytes": s.TotalBytes,
				}
				return render.RenderKV(cmd.OutOrStdout(), out, opts)
			}

			// --list path: fetch the txs with an explicit limit.
			raw, err := rpc.Call(cmd.Context(), "unconfirmed_txs", map[string]any{"limit": strconv.Itoa(limit)})
			if err != nil {
				return fmt.Errorf("calling unconfirmed_txs: %w", err)
			}
			if raw == nil {
				return nil // --curl
			}
			if opts.JSON {
				generic, derr := decodeGeneric(raw)
				if derr != nil {
					return derr
				}
				return render.RenderJSON(cmd.OutOrStdout(), generic, opts)
			}
			var u cometUnconfirmedTxs
			if err := json.Unmarshal(raw, &u); err != nil {
				return fmt.Errorf("decoding unconfirmed_txs: %w", err)
			}
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "n_txs       %s\n", u.NTxs); err != nil {
				return err
			}
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "total       %s\n", u.Total); err != nil {
				return err
			}
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "total_bytes %s\n", u.TotalBytes); err != nil {
				return err
			}
			for i, txb64 := range u.Txs {
				hash, derr := txHashFromBase64(txb64)
				if derr != nil {
					return fmt.Errorf("hashing tx %d: %w", i, derr)
				}
				if _, err := fmt.Fprintln(cmd.OutOrStdout(), hash); err != nil {
					return err
				}
			}
			return nil
		},
	}
	f := cmd.Flags()
	f.BoolVar(&list, "list", false, "fetch pending tx payloads (up to --limit) and print their hashes")
	f.IntVar(&limit, "limit", 30, "maximum txs to request when --list is set")
	f.StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable)")
	return cmd
}

// txHashFromBase64 decodes a CometBFT base64 raw-tx blob and returns
// its `0x<sha256>` hash, mirroring how tendermint/cometbft computes
// tx hashes over the wire payload.
func txHashFromBase64(s string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", fmt.Errorf("base64 decode: %w", err)
	}
	sum := sha256.Sum256(decoded)
	return "0x" + hex.EncodeToString(sum[:]), nil
}
