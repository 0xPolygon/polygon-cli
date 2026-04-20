package decode

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	hproto "github.com/0xPolygon/polygon-cli/internal/heimdall/proto"
)

// newVECmd builds `decode ve <HEX>`. Parses CometBFT vote-extension
// bytes as heimdallv2.sidetxs.VoteExtension and prints a structured
// summary. Input is hex by default because vote extensions surface as
// hex strings in CometBFT logs; base64 is accepted for parity.
func newVECmd() *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "ve <vote-extension>",
		Short: "Decode CometBFT vote-extension bytes as heimdallv2.sidetxs.VoteExtension.",
		Long: strings.TrimSpace(`
Decode CometBFT vote-extension bytes as heimdallv2.sidetxs.VoteExtension.

The vote-extension protobuf is NOT wrapped in an Any on the wire; it is
passed as plain bytes through CometBFT's ExtendVote interface.
`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := decodeInput("ve", args[0])
			if err != nil {
				return err
			}
			ve, err := hproto.UnmarshalVoteExtension(raw)
			if err != nil {
				return err
			}
			env := buildVEEnvelope(ve)
			if jsonOut {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(env)
			}
			buf, err := json.MarshalIndent(env, "", "  ")
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), string(buf))
			return err
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "emit single-line JSON")
	return cmd
}

// buildVEEnvelope turns a *VoteExtension into a JSON-friendly map. We
// render bytes as 0x-hex for compactness with an extra base64 field so
// operators can round-trip the extension.
func buildVEEnvelope(ve *hproto.VoteExtension) map[string]interface{} {
	env := map[string]interface{}{
		"type_url":   hproto.VoteExtensionTypeURL,
		"block_hash": "0x" + hex.EncodeToString(ve.BlockHash),
		"height":     ve.Height,
	}
	if len(ve.SideTxResponses) > 0 {
		resps := make([]map[string]interface{}, 0, len(ve.SideTxResponses))
		for _, r := range ve.SideTxResponses {
			resps = append(resps, map[string]interface{}{
				"tx_hash": "0x" + hex.EncodeToString(r.TxHash),
				"result":  r.Result.String(),
			})
		}
		env["side_tx_responses"] = resps
	}
	if ve.MilestoneProposition != nil {
		mp := ve.MilestoneProposition
		hashes := make([]string, 0, len(mp.BlockHashes))
		for _, h := range mp.BlockHashes {
			hashes = append(hashes, "0x"+hex.EncodeToString(h))
		}
		env["milestone_proposition"] = map[string]interface{}{
			"block_hashes":       hashes,
			"start_block_number": mp.StartBlockNumber,
			"parent_hash":        "0x" + hex.EncodeToString(mp.ParentHash),
			"block_tds":          mp.BlockTDs,
		}
	}
	// Emit the raw bytes too, so callers can sanity-check they got back
	// exactly what went in.
	env["raw_b64"] = base64.StdEncoding.EncodeToString(ve.Marshal())
	return env
}
