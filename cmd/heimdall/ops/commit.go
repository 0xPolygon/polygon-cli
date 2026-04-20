package ops

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// cometCommit is a minimal /commit summary: the header fields and the
// canonical flag for whether the block is signed.
type cometCommit struct {
	SignedHeader struct {
		Header struct {
			ChainID         string `json:"chain_id"`
			Height          string `json:"height"`
			Time            string `json:"time"`
			ProposerAddress string `json:"proposer_address"`
			AppHash         string `json:"app_hash"`
			DataHash        string `json:"data_hash"`
		} `json:"header"`
		Commit struct {
			Height  string `json:"height"`
			Round   int    `json:"round"`
			BlockID struct {
				Hash string `json:"hash"`
			} `json:"block_id"`
		} `json:"commit"`
	} `json:"signed_header"`
	Canonical bool `json:"canonical"`
}

// newCommitCmd builds `ops commit [HEIGHT]`. An empty or omitted
// height means latest. --json passes the full /commit response
// through; default is the KV summary.
func newCommitCmd() *cobra.Command {
	var fields []string
	cmd := &cobra.Command{
		Use:   "commit [HEIGHT]",
		Short: "Fetch a signed CometBFT commit header at height (default latest).",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rpc, cfg, err := newRPCClient(cmd)
			if err != nil {
				return err
			}
			// Build params: CometBFT rejects nil params, so default to
			// an explicit `height: nil` (meaning latest) and only set
			// a concrete height when the arg is supplied.
			params := map[string]any{"height": nil}
			if len(args) == 1 {
				hArg := strings.TrimSpace(args[0])
				if hArg != "" && !strings.EqualFold(hArg, "latest") {
					h, perr := strconv.ParseInt(hArg, 10, 64)
					if perr != nil || h <= 0 {
						return &client.UsageError{Msg: fmt.Sprintf("invalid height %q (want positive integer or `latest`)", hArg)}
					}
					params["height"] = strconv.FormatInt(h, 10)
				}
			}
			raw, err := rpc.Call(cmd.Context(), "commit", params)
			if err != nil {
				return fmt.Errorf("calling commit: %w", err)
			}
			if raw == nil {
				return nil // --curl
			}
			opts := renderOpts(cmd, cfg, fields)
			if opts.JSON {
				generic, derr := decodeGeneric(raw)
				if derr != nil {
					return derr
				}
				return render.RenderJSON(cmd.OutOrStdout(), generic, opts)
			}
			var c cometCommit
			if err := json.Unmarshal(raw, &c); err != nil {
				return fmt.Errorf("decoding commit: %w", err)
			}
			h := c.SignedHeader.Header
			out := map[string]any{
				"chain_id":         h.ChainID,
				"height":           h.Height,
				"time":             h.Time,
				"proposer_address": "0x" + h.ProposerAddress,
				"app_hash":         "0x" + h.AppHash,
				"data_hash":        "0x" + h.DataHash,
				"block_hash":       "0x" + c.SignedHeader.Commit.BlockID.Hash,
				"commit_round":     c.SignedHeader.Commit.Round,
				"canonical":        c.Canonical,
			}
			return render.RenderKV(cmd.OutOrStdout(), out, opts)
		},
	}
	f := cmd.Flags()
	f.StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable)")
	return cmd
}
