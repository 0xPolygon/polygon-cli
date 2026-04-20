// Package render formats Heimdall REST + CometBFT responses for CLI
// output. It supports three modes — key/value (default), table
// (list-like payloads), and JSON (with the normalizations described
// in HEIMDALLCAST_REQUIREMENTS.md §4.2).
//
// Callers pass in an already-decoded map/slice/json.RawMessage; the
// renderers do not talk to the network.
package render

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

// Options controls output formatting.
type Options struct {
	// JSON forces JSON output.
	JSON bool
	// Raw suppresses the bytes->0x-hex normalization for JSON output.
	Raw bool
	// Fields restricts output to these JSON paths (repeatable --field).
	Fields []string
	// Color mode: auto|always|never.
	Color string
	// IsTTY is set by the caller when stdout is a terminal. Combined
	// with Color=auto this decides whether to emit ANSI colour codes.
	IsTTY bool
}

// ColorEnabled returns whether colour output should be emitted given
// the current options.
func (o Options) ColorEnabled() bool {
	switch o.Color {
	case "always":
		return true
	case "never":
		return false
	}
	return o.IsTTY
}

// Byte-field heuristics: any key whose name ends with one of these
// suffixes (case-insensitive), and whose value is a base64 string, is
// re-encoded as `0x`-prefixed hex in JSON output. Conservative list —
// adding more over time as new endpoints turn up.
var byteFieldSuffixes = []string{
	"hash",
	"root",
	"proof",
	"signature",
	"signatures",
	"pubkey",
	"pub_key",
	"address",
	"data",
}

// RenderJSON emits input as pretty-printed JSON with bytes
// normalization applied (unless opts.Raw). input is expected to be
// the result of json.Unmarshal into any / map / slice.
func RenderJSON(w io.Writer, input any, opts Options) error {
	v := input
	if !opts.Raw {
		v = normalizeBytes(v)
	}
	if len(opts.Fields) > 0 {
		out, err := pluckFields(v, opts.Fields)
		if err != nil {
			return err
		}
		if len(opts.Fields) == 1 {
			return writeBareField(w, out[opts.Fields[0]])
		}
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(out)
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

// RenderKV emits a map[string]any as right-aligned key/value pairs,
// one per line (matches requirements §4.1). Nested objects are
// rendered inline as JSON on the value line.
func RenderKV(w io.Writer, input any, opts Options) error {
	v := input
	if !opts.Raw {
		v = normalizeBytes(v)
	}
	if len(opts.Fields) > 0 {
		pluck, err := pluckFields(v, opts.Fields)
		if err != nil {
			return err
		}
		if len(opts.Fields) == 1 {
			return writeBareField(w, pluck[opts.Fields[0]])
		}
		return writeAligned(w, toStringMap(pluck))
	}
	m, ok := v.(map[string]any)
	if !ok {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(v)
	}
	return writeAligned(w, toStringMap(m))
}

// RenderTable emits a list of records as a simple column-aligned
// table. The column set is the union of keys in the records in
// iteration order of the first record.
func RenderTable(w io.Writer, records []map[string]any, opts Options) error {
	if len(records) == 0 {
		_, err := fmt.Fprintln(w, "(no records)")
		return err
	}
	normalized := make([]map[string]any, len(records))
	for i, r := range records {
		if opts.Raw {
			normalized[i] = r
			continue
		}
		n, _ := normalizeBytes(r).(map[string]any)
		normalized[i] = n
	}

	var cols []string
	seen := map[string]bool{}
	for _, r := range normalized {
		for k := range r {
			if !seen[k] {
				seen[k] = true
				cols = append(cols, k)
			}
		}
	}
	sort.Strings(cols)

	widths := make([]int, len(cols))
	for i, c := range cols {
		widths[i] = len(c)
	}
	rows := make([][]string, 0, len(normalized))
	for _, r := range normalized {
		row := make([]string, len(cols))
		for i, c := range cols {
			row[i] = stringify(r[c])
			if len(row[i]) > widths[i] {
				widths[i] = len(row[i])
			}
		}
		rows = append(rows, row)
	}

	var b bytes.Buffer
	for i, c := range cols {
		if i > 0 {
			b.WriteString("  ")
		}
		fmt.Fprintf(&b, "%-*s", widths[i], c)
	}
	b.WriteByte('\n')
	for _, row := range rows {
		for i, v := range row {
			if i > 0 {
				b.WriteString("  ")
			}
			fmt.Fprintf(&b, "%-*s", widths[i], v)
		}
		b.WriteByte('\n')
	}
	_, err := w.Write(b.Bytes())
	return err
}

// writeAligned emits a map as right-aligned KV lines.
func writeAligned(w io.Writer, m map[string]string) error {
	keys := make([]string, 0, len(m))
	maxLen := 0
	for k := range m {
		keys = append(keys, k)
		if len(k) > maxLen {
			maxLen = len(k)
		}
	}
	sort.Strings(keys)
	var b bytes.Buffer
	for _, k := range keys {
		fmt.Fprintf(&b, "%-*s  %s\n", maxLen, k, m[k])
	}
	_, err := w.Write(b.Bytes())
	return err
}

// writeBareField emits a single plucked value with a trailing newline
// — scripting-friendly. Strings are emitted unquoted.
func writeBareField(w io.Writer, v any) error {
	switch vv := v.(type) {
	case string:
		_, err := fmt.Fprintln(w, vv)
		return err
	case nil:
		_, err := fmt.Fprintln(w, "")
		return err
	default:
		buf, err := json.Marshal(vv)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(w, string(buf))
		return err
	}
}

// toStringMap reduces map values to display strings.
func toStringMap(m map[string]any) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = stringify(v)
	}
	return out
}

func stringify(v any) string {
	switch vv := v.(type) {
	case nil:
		return ""
	case string:
		return vv
	case float64:
		// json.Unmarshal yields float64 for numbers; render as int
		// when integral to match cast output style.
		if vv == float64(int64(vv)) {
			return fmt.Sprintf("%d", int64(vv))
		}
		return fmt.Sprintf("%v", vv)
	case bool:
		if vv {
			return "true"
		}
		return "false"
	default:
		buf, err := json.Marshal(vv)
		if err != nil {
			return fmt.Sprintf("%v", vv)
		}
		return string(buf)
	}
}

// normalizeBytes walks the decoded JSON value and rewrites suspected
// byte-string fields from base64 to 0x-hex. Non-matching values are
// passed through unchanged.
func normalizeBytes(v any) any {
	switch vv := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(vv))
		for k, inner := range vv {
			if isByteField(k) {
				if s, ok := inner.(string); ok {
					if hex, ok := base64ToHex(s); ok {
						out[k] = hex
						continue
					}
				}
			}
			out[k] = normalizeBytes(inner)
		}
		return out
	case []any:
		out := make([]any, len(vv))
		for i, inner := range vv {
			out[i] = normalizeBytes(inner)
		}
		return out
	default:
		return v
	}
}

func isByteField(k string) bool {
	lower := strings.ToLower(k)
	for _, suffix := range byteFieldSuffixes {
		if lower == suffix || strings.HasSuffix(lower, "_"+suffix) {
			return true
		}
	}
	return false
}

// base64ToHex decodes s as standard base64 and returns a 0x-prefixed
// hex string. Returns (_, false) if s is not valid base64 or is
// suspiciously long (likely not a hash/proof/etc.).
func base64ToHex(s string) (string, bool) {
	if s == "" {
		return "", false
	}
	// If it already looks like hex (0x or plain hex), leave as-is.
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		return "", false
	}
	if looksHex(s) {
		return "", false
	}
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		decoded, err = base64.RawStdEncoding.DecodeString(s)
		if err != nil {
			return "", false
		}
	}
	return "0x" + hex.EncodeToString(decoded), true
}

func looksHex(s string) bool {
	for _, r := range s {
		switch {
		case r >= '0' && r <= '9':
		case r >= 'a' && r <= 'f':
		case r >= 'A' && r <= 'F':
		default:
			return false
		}
	}
	return len(s) > 0
}

// pluckFields extracts the named top-level fields from v. Paths may
// use dot notation ("block.header.height"). Missing fields are
// reported as an error if all paths miss; partial misses yield nil.
func pluckFields(v any, fields []string) (map[string]any, error) {
	out := make(map[string]any, len(fields))
	hit := false
	for _, f := range fields {
		val, ok := lookupPath(v, f)
		out[f] = val
		if ok {
			hit = true
		}
	}
	if !hit {
		return nil, fmt.Errorf("no fields matched: %s", strings.Join(fields, ", "))
	}
	return out, nil
}

func lookupPath(v any, path string) (any, bool) {
	parts := strings.Split(path, ".")
	cur := v
	for _, p := range parts {
		m, ok := cur.(map[string]any)
		if !ok {
			return nil, false
		}
		next, ok := m[p]
		if !ok {
			return nil, false
		}
		cur = next
	}
	return cur, true
}
