package tx

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/url"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// balanceResponse matches /cosmos/bank/v1beta1/balances/{addr}/by_denom.
type balanceResponse struct {
	Balance struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"balance"`
}

// newBalanceCmd builds `balance <ADDRESS>` (alias `b`). Default
// output is the raw 18-dec integer; --human formats with the decimals
// implied by the denom (fixed at 18 for pol/matic).
func newBalanceCmd() *cobra.Command {
	var denom string
	var human bool
	var fields []string
	cmd := &cobra.Command{
		Use:     "balance <ADDRESS>",
		Aliases: []string{"b"},
		Short:   "Show an account's balance for a denom.",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			addr, err := validateAddress(args[0])
			if err != nil {
				return err
			}
			rest, cfg, err := newRESTClient(cmd)
			if err != nil {
				return err
			}
			d := denom
			if d == "" {
				d = cfg.Denom
			}
			if d == "" {
				d = "pol"
			}
			q := url.Values{}
			q.Set("denom", d)
			body, status, err := rest.Get(cmd.Context(), "/cosmos/bank/v1beta1/balances/"+addr+"/by_denom", q)
			if err != nil {
				return err
			}
			if status == 0 && body == nil {
				return nil // --curl
			}

			opts := renderOpts(cmd, cfg, fields)
			if opts.JSON {
				var generic any
				if err := json.Unmarshal(body, &generic); err != nil {
					return fmt.Errorf("decoding balance for json: %w", err)
				}
				return render.RenderJSON(cmd.OutOrStdout(), generic, opts)
			}

			var resp balanceResponse
			if err := json.Unmarshal(body, &resp); err != nil {
				return fmt.Errorf("decoding balance: %w (body=%q)", err, truncate(body, 256))
			}
			amount := resp.Balance.Amount
			if amount == "" {
				amount = "0"
			}
			if human {
				formatted, err := formatDecimal(amount, 18)
				if err != nil {
					return err
				}
				_, err = fmt.Fprintf(cmd.OutOrStdout(), "%s %s\n", formatted, d)
				return err
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), amount)
			return err
		},
	}
	f := cmd.Flags()
	f.StringVar(&denom, "denom", "", "denom to query (default pol)")
	f.BoolVar(&human, "human", false, "format amount with decimals")
	f.StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable, --json only)")
	return cmd
}

// formatDecimal renders a raw integer string as a decimal with the
// given number of fractional digits. Trailing zeros and trailing dots
// are trimmed for readability.
func formatDecimal(amountStr string, decimals int) (string, error) {
	n, ok := new(big.Int).SetString(amountStr, 10)
	if !ok {
		return "", fmt.Errorf("balance amount %q is not an integer", amountStr)
	}
	neg := n.Sign() < 0
	if neg {
		n = new(big.Int).Neg(n)
	}
	str := n.String()
	if len(str) <= decimals {
		// Pad the integer part with leading zeros so the decimal
		// split is well-defined.
		pad := decimals - len(str) + 1
		str = leftPad(str, pad, '0')
	}
	intPart := str[:len(str)-decimals]
	fracPart := str[len(str)-decimals:]
	// Trim trailing zeros in fracPart.
	for len(fracPart) > 0 && fracPart[len(fracPart)-1] == '0' {
		fracPart = fracPart[:len(fracPart)-1]
	}
	out := intPart
	if fracPart != "" {
		out += "." + fracPart
	}
	if neg {
		out = "-" + out
	}
	return out, nil
}

func leftPad(s string, n int, c byte) string {
	if n <= 0 {
		return s
	}
	buf := make([]byte, n+len(s))
	for i := 0; i < n; i++ {
		buf[i] = c
	}
	copy(buf[n:], s)
	return string(buf)
}
