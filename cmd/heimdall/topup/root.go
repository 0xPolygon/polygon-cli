package topup

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newRootCmd builds `topup root` → GET /topup/dividend-account-root.
// Renders `account_root_hash` as 0x-hex by default; --raw (and --json
// with --raw) preserves the upstream base64.
func newRootCmd() *cobra.Command {
	var fields []string
	cmd := &cobra.Command{
		Use:   "root",
		Short: "Print the Merkle root of all dividend accounts.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			rest, cfg, err := newRESTClient(cmd)
			if err != nil {
				return err
			}
			body, status, err := rest.Get(cmd.Context(), "/topup/dividend-account-root", nil)
			if err != nil {
				return err
			}
			if status == 0 && body == nil {
				return nil
			}
			opts := renderOpts(cmd, cfg, fields)
			m, err := decodeJSONMap(body, "topup dividend-account-root")
			if err != nil {
				return err
			}
			if opts.JSON {
				return render.RenderJSON(cmd.OutOrStdout(), m, opts)
			}
			// Default text output: print just the root hash. Respect
			// --raw by leaving the base64 alone; otherwise re-encode to
			// 0x-hex for convenience.
			rawRoot, ok := m["account_root_hash"].(string)
			if !ok {
				// Fall back to generic KV rendering.
				return render.RenderKV(cmd.OutOrStdout(), m, opts)
			}
			if opts.Raw {
				_, err = fmt.Fprintln(cmd.OutOrStdout(), rawRoot)
				return err
			}
			// Decode base64 → 0x-hex. If the value already looks like
			// hex leave it alone.
			decoded, derr := base64.StdEncoding.DecodeString(rawRoot)
			if derr != nil {
				// Not base64; print as-is.
				_, err = fmt.Fprintln(cmd.OutOrStdout(), rawRoot)
				return err
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), "0x"+hex.EncodeToString(decoded))
			return err
		},
	}
	cmd.Flags().StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable, --json only)")
	return cmd
}
