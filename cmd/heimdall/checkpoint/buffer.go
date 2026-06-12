package checkpoint

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/cmdutil"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newBufferCmd builds `checkpoint buffer` → GET /checkpoints/buffer.
// When the proposer is the zero address we print `empty` plus the
// buffer-empty hint rather than rendering the meaningless zeros.
// JSON passthrough is preserved; the empty-buffer hint is a
// human-readable affordance, not a structural one.
func newBufferCmd() *cobra.Command {
	return pkg.NewGetCmd(cmdutil.Get{
		Use:   "buffer",
		Short: "Show the in-flight (buffered) checkpoint.",
		Path:  "/checkpoints/buffer",
		Label: "checkpoint buffer",
		Render: func(cmd *cobra.Command, m map[string]any, opts render.Options) error {
			if isBufferEmpty(m) {
				if _, werr := fmt.Fprintln(cmd.OutOrStdout(), "empty"); werr != nil {
					return werr
				}
				return render.WriteHint(cmd.OutOrStdout(), render.HintBufferEmpty, opts)
			}
			return renderCheckpointKV(cmd, m, opts)
		},
	})
}

// isBufferEmpty returns true if the inner checkpoint has a zero-address
// or empty-string proposer. Heimdall uses both spellings depending on
// the release: the zero-20-byte form `0x00…00` and, on v2, a literal
// empty string.
func isBufferEmpty(m map[string]any) bool {
	inner, ok := m["checkpoint"].(map[string]any)
	if !ok {
		return false
	}
	p, ok := inner["proposer"].(string)
	if !ok {
		return false
	}
	return isZeroOrEmptyAddress(p)
}

func isZeroOrEmptyAddress(s string) bool {
	lower := strings.ToLower(strings.TrimSpace(s))
	if lower == "" {
		return true
	}
	lower = strings.TrimPrefix(lower, "0x")
	if lower == "" {
		return true
	}
	for _, r := range lower {
		if r != '0' {
			return false
		}
	}
	return true
}
