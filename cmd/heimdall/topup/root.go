package topup

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/cmdutil"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newRootCmd builds `topup root` → GET /topup/dividend-account-root.
// Renders `account_root_hash` as 0x-hex by default; --raw (and --json
// with --raw) preserves the upstream base64.
func newRootCmd() *cobra.Command {
	return pkg.NewGetCmd(cmdutil.Get{
		Use:         "root",
		Short:       "Print the Merkle root of all dividend accounts.",
		Path:        "/topup/dividend-account-root",
		Label:       "topup dividend-account-root",
		FieldsUsage: "pluck one or more fields (repeatable, --json only)",
		Render: func(cmd *cobra.Command, m map[string]any, opts render.Options) error {
			// Default text output: print just the root hash. Respect
			// --raw by leaving the base64 alone; otherwise re-encode to
			// 0x-hex for convenience.
			rawRoot, ok := m["account_root_hash"].(string)
			if !ok {
				// Fall back to generic KV rendering.
				return render.RenderKV(cmd.OutOrStdout(), m, opts)
			}
			if opts.Raw {
				_, err := fmt.Fprintln(cmd.OutOrStdout(), rawRoot)
				return err
			}
			// Decode base64 → 0x-hex. If the value already looks like
			// hex leave it alone.
			decoded, derr := base64.StdEncoding.DecodeString(rawRoot)
			if derr != nil {
				// Not base64; print as-is.
				_, err := fmt.Fprintln(cmd.OutOrStdout(), rawRoot)
				return err
			}
			_, err := fmt.Fprintln(cmd.OutOrStdout(), "0x"+hex.EncodeToString(decoded))
			return err
		},
	})
}
