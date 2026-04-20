package clerk

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newGetCmd builds `state-sync get <ID>` → GET
// /clerk/event-records/{id}. The same code path is re-entered from
// ClerkCmd's RunE when a bare integer is provided (`state-sync 36610`).
//
// The record's `data` field is rendered as `0x…`-hex by default; pass
// --base64 (or the global --raw) to preserve the upstream base64.
func newGetCmd() *cobra.Command {
	var (
		fields []string
		base64 bool
	)
	cmd := &cobra.Command{
		Use:   "get <ID>",
		Short: "Fetch one event-record by id.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(cmd, args[0], base64, fields...)
		},
	}
	f := cmd.Flags()
	f.BoolVar(&base64, "base64", false, "preserve raw base64 for `data` (default 0x-hex)")
	f.StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable)")
	return cmd
}

// runGet is the shared implementation used by both `state-sync get
// <ID>` and the bare-integer ClerkCmd shorthand. The bare-integer
// caller does not expose flags, so base64/fields default to zero values
// and the command falls back to the global `--raw` / default KV output.
func runGet(cmd *cobra.Command, idArg string, base64 bool, fields ...string) error {
	id, err := strconv.ParseUint(idArg, 10, 64)
	if err != nil {
		return &client.UsageError{Msg: fmt.Sprintf("event-record id must be a positive integer, got %q", idArg)}
	}
	rest, cfg, err := newRESTClient(cmd)
	if err != nil {
		return err
	}
	body, status, err := rest.Get(cmd.Context(), fmt.Sprintf("/clerk/event-records/%d", id), nil)
	if err != nil {
		return err
	}
	if status == 0 && body == nil {
		return nil
	}
	opts := renderOpts(cmd, cfg, fields, base64)
	m, err := decodeJSONMap(body, "clerk event-record")
	if err != nil {
		return err
	}
	if opts.JSON {
		return render.RenderJSON(cmd.OutOrStdout(), m, opts)
	}
	return renderRecordKV(cmd, m, opts)
}

// renderRecordKV unwraps the { "record": {...} } envelope and renders
// with the shared KV formatter. The `record_time` timestamp from this
// endpoint is RFC3339 (not unix seconds) so we leave it untouched.
func renderRecordKV(cmd *cobra.Command, m map[string]any, opts render.Options) error {
	inner, ok := m["record"].(map[string]any)
	if !ok {
		return render.RenderKV(cmd.OutOrStdout(), m, opts)
	}
	return render.RenderKV(cmd.OutOrStdout(), inner, opts)
}
