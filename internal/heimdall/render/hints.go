package render

import (
	"fmt"
	"io"
	"strings"
)

// Hint is a short explanatory line rendered in gray after an
// otherwise confusing response. Source catalogues the misleading
// cases from HEIMDALLCAST_REQUIREMENTS.md §4.5.
type Hint struct {
	// Key identifies the hint for logging / tests. Not rendered.
	Key string
	// Body is the hint text. Rendered in gray when color is enabled.
	Body string
}

// Hints the command packages can reference by name. Centralised so
// the wording stays consistent across subcommands.
var (
	HintIsOldRenamed = Hint{
		Key:  "is-old-renamed",
		Body: "note: upstream `is_old` renamed to `is_current` in this tool (upstream naming was misleading)",
	}
	HintL1NotConfigured = Hint{
		Key:  "l1-not-configured",
		Body: "hint: this node does not have `eth_rpc_url` configured; L1 replay checks will fail until it is set",
	}
	HintBufferEmpty = Hint{
		Key:  "buffer-empty",
		Body: "hint: the buffer is empty (no checkpoint in flight) - this is not an error",
	}
	HintMilestoneOutOfRange = Hint{
		Key:  "milestone-range",
		Body: "hint: milestone numbers are 1-indexed and bounded by `milestone count`",
	}
	HintPaginationLimit = Hint{
		Key:  "pagination-limit",
		Body: "hint: list endpoints require `pagination.limit` (try --limit)",
	}
)

// WriteHint emits a single hint line to w. When colour is enabled the
// line is wrapped in the ANSI dim sequence.
func WriteHint(w io.Writer, h Hint, opts Options) error {
	body := h.Body
	if opts.ColorEnabled() {
		body = "\x1b[2m" + body + "\x1b[0m"
	}
	_, err := fmt.Fprintln(w, body)
	return err
}

// WriteHints emits multiple hints preserving order; a nil slice is a
// no-op.
func WriteHints(w io.Writer, hints []Hint, opts Options) error {
	for _, h := range hints {
		if err := WriteHint(w, h, opts); err != nil {
			return err
		}
	}
	return nil
}

// DetectHints scans a rendered KV map for the well-known footguns and
// returns the hints that apply. Does not mutate the input.
func DetectHints(m map[string]any) []Hint {
	var out []Hint
	if _, ok := m["is_old"]; ok {
		out = append(out, HintIsOldRenamed)
	}
	if v, ok := m["proposer"]; ok {
		if s, ok := v.(string); ok && isZeroAddress(s) {
			out = append(out, HintBufferEmpty)
		}
	}
	return out
}

func isZeroAddress(s string) bool {
	s = strings.TrimPrefix(strings.TrimPrefix(strings.ToLower(s), "0x"), "0x")
	if s == "" {
		return false
	}
	for _, r := range s {
		if r != '0' {
			return false
		}
	}
	return true
}
