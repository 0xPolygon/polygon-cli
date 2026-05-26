package milestone

import (
	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newLatestCmd builds `milestone latest` → GET /milestones/latest.
// The single-milestone envelope is unwrapped for KV output; `hash` is
// re-encoded from base64 to `0x…`-hex by the renderer unless --raw is
// set. The envelope is addressed by the URL's `number` sequence, but
// the response body exposes only `milestone_id` (which is *not* the
// same value — see the package usage docs).
func newLatestCmd() *cobra.Command {
	var fields []string
	cmd := &cobra.Command{
		Use:   "latest",
		Short: "Show the latest milestone.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			rest, cfg, err := newRESTClient(cmd)
			if err != nil {
				return err
			}
			body, status, err := rest.Get(cmd.Context(), "/milestones/latest", nil)
			if err != nil {
				return err
			}
			if status == 0 && body == nil {
				return nil
			}
			opts := renderOpts(cmd, cfg, fields)
			m, err := decodeJSONMap(body, "milestone latest")
			if err != nil {
				return err
			}
			if opts.JSON {
				return render.RenderJSON(cmd.OutOrStdout(), m, opts)
			}
			// The server does not tell us the latest "number" — it's
			// implicitly equal to `milestone count`. Passing 0 here
			// suppresses the number-label row; renderMilestoneKV still
			// prints milestone_id from the body.
			return renderMilestoneKV(cmd, m, opts, 0)
		},
	}
	cmd.Flags().StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable)")
	return cmd
}

// renderMilestoneKV unwraps the { "milestone": {...} } envelope,
// annotates the timestamp with human-readable age, and renders with
// the shared KV formatter. When number > 0 it is prepended to the
// rendered output so the reader sees *both* `number` (the URL-path
// sequence) and `milestone_id` (the on-chain id from the body) — the
// footgun called out in HEIMDALLCAST_REQUIREMENTS.md §3.2.3.
func renderMilestoneKV(cmd *cobra.Command, m map[string]any, opts render.Options, number uint64) error {
	inner, ok := m["milestone"].(map[string]any)
	if !ok {
		return render.RenderKV(cmd.OutOrStdout(), m, opts)
	}
	if ts, ok := inner["timestamp"].(string); ok && ts != "" {
		inner["timestamp"] = render.AnnotateUnixSeconds(ts)
	}
	if number > 0 {
		// Only surface `number` when the caller knows what it is (i.e.
		// the user asked for a specific one). We splice it in rather
		// than mutate the upstream map key set in a way that would
		// conflict with a future server change.
		inner["number"] = itoa(number)
	}
	return render.RenderKV(cmd.OutOrStdout(), inner, opts)
}
