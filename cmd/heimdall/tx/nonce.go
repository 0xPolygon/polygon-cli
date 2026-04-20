package tx

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// accountResponse matches the shape of the Cosmos SDK auth REST
// gateway response for /cosmos/auth/v1beta1/accounts/{addr}. Only
// the BaseAccount-like fields are decoded; extra types land as raw
// via --json.
type accountResponse struct {
	Account struct {
		Type          string `json:"@type"`
		Address       string `json:"address"`
		AccountNumber string `json:"account_number"`
		Sequence      string `json:"sequence"`
	} `json:"account"`
}

// fetchAccount queries /cosmos/auth/v1beta1/accounts/{addr} and
// returns the decoded response plus the raw body for --json passthrough.
// Returns (nil, nil, nil) under --curl.
func fetchAccount(ctx context.Context, rest *client.RESTClient, addr string) (*accountResponse, []byte, error) {
	body, status, err := rest.Get(ctx, "/cosmos/auth/v1beta1/accounts/"+addr, nil)
	if err != nil {
		return nil, body, err
	}
	if status == 0 && body == nil {
		return nil, nil, nil // --curl
	}
	var out accountResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, body, fmt.Errorf("decoding account: %w (body=%q)", err, truncate(body, 256))
	}
	return &out, body, nil
}

// newNonceCmd builds `nonce <ADDRESS>`. Prints the bare sequence
// number by default; --json returns the full account object with
// bytes normalization applied.
func newNonceCmd() *cobra.Command {
	return newNonceLikeCmd("nonce", "Print an account's sequence number.")
}

// newSequenceAliasCmd registers `sequence <ADDRESS>` as a top-level
// synonym for `nonce` — both spellings are common (Cosmos → sequence,
// cast → nonce).
func newSequenceAliasCmd() *cobra.Command {
	return newNonceLikeCmd("sequence", "Alias of nonce; print an account's sequence.")
}

// newNonceLikeCmd is the shared constructor so `nonce` and `sequence`
// are byte-identical apart from Use.
func newNonceLikeCmd(use, short string) *cobra.Command {
	var fields []string
	cmd := &cobra.Command{
		Use:   use + " <ADDRESS>",
		Short: short,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			addr, err := validateAddress(args[0])
			if err != nil {
				return err
			}
			rest, cfg, err := newRESTClient(cmd)
			if err != nil {
				return err
			}
			acc, raw, err := fetchAccount(cmd.Context(), rest, addr)
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
					return fmt.Errorf("decoding account for json: %w", err)
				}
				return render.RenderJSON(cmd.OutOrStdout(), generic, opts)
			}
			// Bare output (scripting-friendly): just the sequence.
			if acc.Account.Sequence == "" {
				return fmt.Errorf("account %s has no sequence field (type=%s)", addr, acc.Account.Type)
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), acc.Account.Sequence)
			return err
		},
	}
	cmd.Flags().StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable, --json only)")
	return cmd
}
