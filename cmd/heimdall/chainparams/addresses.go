package chainparams

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newAddressesCmd builds `chainmanager addresses` — a derived view over
// GET /chainmanager/params that surfaces just the entries of
// `params.chain_params` whose key ends in `_address`, plus the chain
// ids. This is aimed at the "paste into etherscan" workflow, where the
// confirmation depths get in the way.
//
// Default text output: one `<name>=<value>` per line, alphabetized for
// stable output. --json emits a map[string]string with the same
// entries.
func newAddressesCmd() *cobra.Command {
	var fields []string
	cmd := &cobra.Command{
		Use:   "addresses",
		Short: "Print L1 contract addresses and chain ids.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			rest, cfg, err := newRESTClient(cmd)
			if err != nil {
				return err
			}
			body, status, err := rest.Get(cmd.Context(), "/chainmanager/params", nil)
			if err != nil {
				return err
			}
			if status == 0 && body == nil {
				return nil
			}
			m, err := decodeJSONMap(body, "chainmanager params")
			if err != nil {
				return err
			}
			addrs, err := extractAddresses(m)
			if err != nil {
				return err
			}
			opts := renderOpts(cmd, cfg, fields)
			if opts.JSON {
				// Convert to map[string]any so RenderJSON honours --field
				// plucking over the derived view.
				derived := make(map[string]any, len(addrs))
				for k, v := range addrs {
					derived[k] = v
				}
				return render.RenderJSON(cmd.OutOrStdout(), derived, opts)
			}
			keys := make([]string, 0, len(addrs))
			for k := range addrs {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				if _, werr := fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, addrs[k]); werr != nil {
					return werr
				}
			}
			return nil
		},
	}
	cmd.Flags().StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable, --json only)")
	return cmd
}

// extractAddresses pulls the chain ids and `*_address` entries out of
// the `/chainmanager/params` response envelope. It does not invent
// fields: whatever the server returned under `params.chain_params`
// that ends with `_address`, plus `bor_chain_id` and
// `heimdall_chain_id`, is surfaced.
func extractAddresses(m map[string]any) (map[string]string, error) {
	params, ok := m["params"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("chainmanager params response missing `params` object (body=%v)", keysOf(m))
	}
	chainParams, ok := params["chain_params"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("chainmanager params response missing `params.chain_params` object")
	}
	out := make(map[string]string, len(chainParams))
	for k, v := range chainParams {
		if k == "bor_chain_id" || k == "heimdall_chain_id" || strings.HasSuffix(k, "_address") {
			s, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("chain_params.%s is not a string (got %T)", k, v)
			}
			out[k] = s
		}
	}
	return out, nil
}

// keysOf returns the sorted keys of m for diagnostic error messages.
// Output is JSON-encoded so a reader can paste it verbatim.
func keysOf(m map[string]any) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	b, _ := json.Marshal(keys)
	return string(b)
}
