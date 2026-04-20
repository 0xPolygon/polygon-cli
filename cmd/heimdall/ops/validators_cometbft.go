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

// cometValidatorsResp is the subset of /validators used for the table.
type cometValidatorsResp struct {
	BlockHeight string `json:"block_height"`
	Validators  []struct {
		Address string `json:"address"`
		PubKey  struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"pub_key"`
		VotingPower      string `json:"voting_power"`
		ProposerPriority string `json:"proposer_priority"`
	} `json:"validators"`
	Count string `json:"count"`
	Total string `json:"total"`
}

// newValidatorsCometBFTCmd builds `ops validators-cometbft [HEIGHT]`.
// This returns the CometBFT consensus validator set. Note that
// `heimdall validator` is the canonical way to query Heimdall's
// x/stake module — this command is deliberately distinct.
func newValidatorsCometBFTCmd() *cobra.Command {
	var fields []string
	var perPage int
	var page int
	cmd := &cobra.Command{
		Use:   "validators-cometbft [HEIGHT]",
		Short: "List CometBFT consensus validators (NOT Heimdall x/stake).",
		Long: `List validators from CometBFT's /validators endpoint at a given
height (default latest). Output is the consensus layer's view: the
20-byte consensus address, the validator's Secp256k1-eth pubkey,
voting power, and proposer priority.

This is distinct from the Heimdall x/stake validator set. Use
'polycli heimdall validator' for staking info (operator address, moniker,
jailed status, etc).`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if perPage <= 0 {
				return &client.UsageError{Msg: fmt.Sprintf("--per-page must be positive, got %d", perPage)}
			}
			if page <= 0 {
				return &client.UsageError{Msg: fmt.Sprintf("--page must be positive, got %d", page)}
			}
			rpc, cfg, err := newRPCClient(cmd)
			if err != nil {
				return err
			}
			// Hint on stderr so scripted callers piping stdout get a clean list.
			if _, werr := fmt.Fprintln(cmd.ErrOrStderr(),
				"hint: this is the CometBFT consensus set. For staking info use 'heimdall validator'."); werr != nil {
				return werr
			}

			params := map[string]any{
				"height":   nil,
				"page":     strconv.Itoa(page),
				"per_page": strconv.Itoa(perPage),
			}
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

			raw, err := rpc.Call(cmd.Context(), "validators", params)
			if err != nil {
				return fmt.Errorf("calling validators: %w", err)
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
			var vr cometValidatorsResp
			if err := json.Unmarshal(raw, &vr); err != nil {
				return fmt.Errorf("decoding validators: %w", err)
			}
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "block_height  %s\n", vr.BlockHeight); err != nil {
				return err
			}
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "count         %s / %s\n", vr.Count, vr.Total); err != nil {
				return err
			}
			records := make([]map[string]any, 0, len(vr.Validators))
			for _, v := range vr.Validators {
				records = append(records, map[string]any{
					"address":           "0x" + v.Address,
					"voting_power":      v.VotingPower,
					"proposer_priority": v.ProposerPriority,
					"pubkey_type":       v.PubKey.Type,
				})
			}
			return render.RenderTable(cmd.OutOrStdout(), records, opts)
		},
	}
	f := cmd.Flags()
	f.IntVar(&page, "page", 1, "page number (1-indexed)")
	f.IntVar(&perPage, "per-page", 100, "validators per page (CometBFT default is 30, max 100)")
	f.StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable)")
	return cmd
}
