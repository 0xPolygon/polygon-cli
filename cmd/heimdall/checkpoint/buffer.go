package checkpoint

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newBufferCmd builds `checkpoint buffer` → GET /checkpoints/buffer.
// When the proposer is the zero address we print `empty` plus the
// buffer-empty hint rather than rendering the meaningless zeros.
func newBufferCmd() *cobra.Command {
	var fields []string
	cmd := &cobra.Command{
		Use:   "buffer",
		Short: "Show the in-flight (buffered) checkpoint.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			rest, cfg, err := newRESTClient(cmd)
			if err != nil {
				return err
			}
			body, status, err := rest.Get(cmd.Context(), "/checkpoints/buffer", nil)
			if err != nil {
				return err
			}
			if status == 0 && body == nil {
				return nil
			}
			opts := renderOpts(cmd, cfg, fields)
			m, err := decodeJSONMap(body, "checkpoint buffer")
			if err != nil {
				return err
			}
			// Preserve JSON passthrough; the empty-buffer hint is a
			// human-readable affordance, not a structural one.
			if opts.JSON {
				return render.RenderJSON(cmd.OutOrStdout(), m, opts)
			}
			if isBufferEmpty(m) {
				if _, werr := fmt.Fprintln(cmd.OutOrStdout(), "empty"); werr != nil {
					return werr
				}
				return render.WriteHint(cmd.OutOrStdout(), render.HintBufferEmpty, opts)
			}
			return renderCheckpointKV(cmd, m, opts)
		},
	}
	cmd.Flags().StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable)")
	return cmd
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

