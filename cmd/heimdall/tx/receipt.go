package tx

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// fetchStatusLatest calls /status and returns latest_block_height as
// an int64. Used by the --confirmations wait loop.
func fetchStatusLatest(ctx context.Context, rpc *client.RPCClient) (int64, error) {
	raw, err := rpc.Call(ctx, "status", nil)
	if err != nil {
		return 0, fmt.Errorf("fetching status: %w", err)
	}
	if raw == nil {
		return 0, fmt.Errorf("receipt --confirmations is incompatible with --curl")
	}
	var st struct {
		SyncInfo struct {
			LatestBlockHeight string `json:"latest_block_height"`
		} `json:"sync_info"`
	}
	if err := json.Unmarshal(raw, &st); err != nil {
		return 0, fmt.Errorf("decoding status: %w", err)
	}
	if st.SyncInfo.LatestBlockHeight == "" {
		return 0, fmt.Errorf("status did not contain latest_block_height")
	}
	return strconv.ParseInt(st.SyncInfo.LatestBlockHeight, 10, 64)
}

// waitForConfirmations blocks until the tip height is at least
// txHeight + confirmations. Cancellable via ctx. Polls every pollInt.
func waitForConfirmations(ctx context.Context, rpc *client.RPCClient, txHeight int64, confirmations int, pollInt time.Duration) error {
	target := txHeight + int64(confirmations)
	timer := time.NewTimer(pollInt)
	defer timer.Stop()
	for {
		tip, err := fetchStatusLatest(ctx, rpc)
		if err != nil {
			return err
		}
		if tip >= target {
			return nil
		}
		timer.Reset(pollInt)
		select {
		case <-timer.C:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// newReceiptCmd builds `receipt <HASH>` (alias `re`). Same wire call
// as `tx` but renders the event/log stream prominently, mirroring
// `cast receipt`.
func newReceiptCmd() *cobra.Command {
	var confirmations int
	var fields []string
	cmd := &cobra.Command{
		Use:     "receipt <HASH>",
		Aliases: []string{"re"},
		Short:   "Show a transaction receipt (events + logs).",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if confirmations < 0 {
				return &client.UsageError{Msg: "--confirmations must be non-negative"}
			}
			hexHash, err := normalizeHash(args[0])
			if err != nil {
				return err
			}
			rpc, cfg, err := newRPCClient(cmd)
			if err != nil {
				return err
			}
			ctx := cmd.Context()

			tx, raw, err := fetchTx(ctx, rpc, hexHash)
			if err != nil {
				return err
			}
			if raw == nil {
				return nil // --curl
			}

			if confirmations > 0 {
				txHeight, perr := strconv.ParseInt(tx.Height, 10, 64)
				if perr != nil {
					return fmt.Errorf("parsing tx height %q: %w", tx.Height, perr)
				}
				if err := waitForConfirmations(ctx, rpc, txHeight, confirmations, 500*time.Millisecond); err != nil {
					return err
				}
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

			if err := render.RenderKV(cmd.OutOrStdout(), out, opts); err != nil {
				return err
			}
			return writeEvents(cmd.OutOrStdout(), tx.TxResult.Events, opts)
		},
	}
	f := cmd.Flags()
	f.IntVar(&confirmations, "confirmations", 0, "wait until tip is at least tx.height + N")
	f.StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable)")
	return cmd
}

// writeEvents emits one event block per line group, matching a cast
// receipt's Logs section. Honours --json by falling back to a struct
// dump when opts.JSON is set (JSON path already short-circuits
// earlier; this is the KV/receipt path).
func writeEvents(w interface{ Write(p []byte) (int, error) }, events []cometTxResultEvent, _ render.Options) error {
	if len(events) == 0 {
		_, err := fmt.Fprintln(w, "\nevents  (none)")
		return err
	}
	if _, err := fmt.Fprintf(w, "\nevents  (%d)\n", len(events)); err != nil {
		return err
	}
	for i, ev := range events {
		if _, err := fmt.Fprintf(w, "[%d] %s\n", i, ev.Type); err != nil {
			return err
		}
		for _, a := range ev.Attributes {
			if _, err := fmt.Fprintf(w, "    %s = %s\n", a.Key, a.Value); err != nil {
				return err
			}
		}
	}
	return nil
}
