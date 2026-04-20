package tx

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// cometTxResult is the decoded shape of a CometBFT /tx response.
type cometTxResult struct {
	Hash     string            `json:"hash"`
	Height   string            `json:"height"`
	Index    int               `json:"index"`
	Tx       string            `json:"tx"`
	TxResult cometTxResultBody `json:"tx_result"`
}

type cometTxResultBody struct {
	Code      int                  `json:"code"`
	Data      string               `json:"data"`
	Log       string               `json:"log"`
	Info      string               `json:"info"`
	GasWanted string               `json:"gas_wanted"`
	GasUsed   string               `json:"gas_used"`
	Events    []cometTxResultEvent `json:"events"`
	Codespace string               `json:"codespace"`
}

type cometTxResultEvent struct {
	Type       string                        `json:"type"`
	Attributes []cometTxResultEventAttribute `json:"attributes"`
}

type cometTxResultEventAttribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Index bool   `json:"index"`
}

// fetchTx calls CometBFT /tx at the given hash. The JSON-RPC
// "hash" argument is base64-encoded bytes, not hex — CometBFT's
// reflect-based RPC decodes string fields into `[]byte` via base64.
// Callers pass hex (with or without 0x prefix) and this function
// translates. Returns (nil, nil, nil) under --curl.
func fetchTx(ctx context.Context, rpc *client.RPCClient, hexHash string) (*cometTxResult, json.RawMessage, error) {
	h := strings.TrimPrefix(strings.TrimPrefix(hexHash, "0x"), "0X")
	raw, err := hex.DecodeString(h)
	if err != nil {
		return nil, nil, fmt.Errorf("decoding hash %q: %w", hexHash, err)
	}
	params := map[string]any{"hash": base64.StdEncoding.EncodeToString(raw), "prove": false}
	resRaw, err := rpc.Call(ctx, "tx", params)
	if err != nil {
		return nil, nil, fmt.Errorf("fetching tx: %w", err)
	}
	if resRaw == nil {
		return nil, nil, nil
	}
	var out cometTxResult
	if err := json.Unmarshal(resRaw, &out); err != nil {
		return nil, nil, fmt.Errorf("decoding tx: %w", err)
	}
	return &out, resRaw, nil
}

// newTxCmd builds `tx <HASH>` (alias `t`). Prints a summary keyed by
// hash/height/code/gas plus any log. Event / log details go to
// `receipt`. --raw preserves the base64 TxRaw body.
func newTxCmd() *cobra.Command {
	var fields []string
	cmd := &cobra.Command{
		Use:     "tx <HASH>",
		Aliases: []string{"t"},
		Short:   "Show a transaction by hash.",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			hexHash, err := normalizeHash(args[0])
			if err != nil {
				return err
			}
			rpc, cfg, err := newRPCClient(cmd)
			if err != nil {
				return err
			}
			tx, raw, err := fetchTx(cmd.Context(), rpc, hexHash)
			if err != nil {
				return err
			}
			if raw == nil {
				return nil // --curl
			}

			opts := renderOpts(cmd, cfg, fields)
			if opts.JSON {
				var generic any
				if err := json.Unmarshal(raw, &generic); err != nil {
					return fmt.Errorf("decoding tx for json: %w", err)
				}
				return render.RenderJSON(cmd.OutOrStdout(), generic, opts)
			}

			out := map[string]any{
				"hash":       "0x" + tx.Hash,
				"height":     tx.Height,
				"index":      tx.Index,
				"code":       tx.TxResult.Code,
				"gas_used":   tx.TxResult.GasUsed,
				"gas_wanted": tx.TxResult.GasWanted,
				"num_events": len(tx.TxResult.Events),
			}
			if tx.TxResult.Log != "" {
				out["raw_log"] = tx.TxResult.Log
			}
			if tx.TxResult.Codespace != "" {
				out["codespace"] = tx.TxResult.Codespace
			}
			if cfg.Raw {
				out["tx"] = tx.Tx
			}
			return render.RenderKV(cmd.OutOrStdout(), out, opts)
		},
	}
	cmd.Flags().StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable)")
	return cmd
}
