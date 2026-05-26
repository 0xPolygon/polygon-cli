package tx

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newRPCCmd builds `rpc <METHOD> [ARGS...]`, a raw JSON-RPC
// passthrough that mirrors `cast rpc`. ARGS are interpreted as
// alternating `key=value` pairs. Values are passed as JSON when they
// parse as JSON, otherwise as strings.
func newRPCCmd() *cobra.Command {
	var fields []string
	cmd := &cobra.Command{
		Use:   "rpc <METHOD> [KEY=VALUE...]",
		Short: "Invoke an arbitrary CometBFT JSON-RPC method.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			method := args[0]
			params, err := parseRPCArgs(args[1:])
			if err != nil {
				return err
			}
			rpc, cfg, err := newRPCClient(cmd)
			if err != nil {
				return err
			}
			raw, err := rpc.Call(cmd.Context(), method, params)
			if err != nil {
				return fmt.Errorf("rpc %s: %w", method, err)
			}
			if raw == nil {
				return nil // --curl
			}
			opts := renderOpts(cmd, cfg, fields)
			var generic any
			if err := json.Unmarshal(raw, &generic); err != nil {
				return fmt.Errorf("decoding rpc response: %w", err)
			}
			// Default to JSON output for rpc passthrough — KV doesn't
			// make sense when the caller doesn't know the schema.
			opts.JSON = true
			return render.RenderJSON(cmd.OutOrStdout(), generic, opts)
		},
	}
	cmd.Flags().StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable)")
	return cmd
}

// parseRPCArgs converts a flat list of "key=value" strings into a
// CometBFT params map. Values that parse as JSON are inserted as
// JSON; everything else is inserted as a string.
func parseRPCArgs(args []string) (map[string]any, error) {
	if len(args) == 0 {
		return nil, nil
	}
	out := make(map[string]any, len(args))
	for _, a := range args {
		eq := strings.IndexByte(a, '=')
		if eq <= 0 {
			return nil, &client.UsageError{Msg: fmt.Sprintf("rpc arg %q must be KEY=VALUE", a)}
		}
		k := a[:eq]
		v := a[eq+1:]
		// Try JSON first so numbers, bools, objects and arrays travel
		// correctly over the wire. Unquoted strings stay as strings.
		var parsed any
		if err := json.Unmarshal([]byte(v), &parsed); err == nil {
			out[k] = parsed
			continue
		}
		out[k] = v
	}
	return out, nil
}
