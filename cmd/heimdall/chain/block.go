package chain

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newBlockCmd builds the `block [HEIGHT]` command (alias `bl`). The
// HEIGHT arg accepts an integer, `latest`, or `earliest`. Default
// output is the summary keys chain_id/height/time/proposer/num_txs;
// --full includes the tx list.
func newBlockCmd() *cobra.Command {
	var full bool
	var fields []string
	cmd := &cobra.Command{
		Use:     "block [HEIGHT]",
		Aliases: []string{"bl"},
		Short:   "Show a CometBFT block by height (or latest).",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rpc, cfg, err := newRPCClient(cmd)
			if err != nil {
				return err
			}
			ctx := cmd.Context()

			heightArg := ""
			if len(args) == 1 {
				heightArg = args[0]
			}
			height, err := resolveHeight(ctx, rpc, heightArg)
			if err != nil {
				return err
			}
			blk, raw, err := fetchBlock(ctx, rpc, height)
			if err != nil {
				return err
			}
			if raw == nil {
				return nil // --curl
			}

			opts := renderOpts(cmd, cfg, fields)

			// --json passes through the full result with render's
			// byte-field normalization.
			if opts.JSON {
				var generic any
				if err := json.Unmarshal(raw, &generic); err != nil {
					return fmt.Errorf("decoding block for json: %w", err)
				}
				return render.RenderJSON(cmd.OutOrStdout(), generic, opts)
			}

			out := map[string]any{
				"chain_id": blk.Block.Header.ChainID,
				"height":   blk.Block.Header.Height,
				"time":     blk.Block.Header.Time,
				"proposer": "0x" + blk.Block.Header.ProposerAddress,
				"num_txs":  len(blk.Block.Data.Txs),
				"hash":     "0x" + blk.BlockID.Hash,
			}
			if full {
				out["txs"] = blk.Block.Data.Txs
			}
			if err := render.RenderKV(cmd.OutOrStdout(), out, opts); err != nil {
				return err
			}
			// Hint path: zero-proposer is the only generic trigger
			// DetectHints can spot on a block summary.
			hints := render.DetectHints(out)
			return render.WriteHints(cmd.ErrOrStderr(), hints, opts)
		},
	}
	f := cmd.Flags()
	f.BoolVar(&full, "full", false, "include the full tx list in output")
	f.StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable)")
	return cmd
}
