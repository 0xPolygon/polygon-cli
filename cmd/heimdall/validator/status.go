package validator

import (
	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newStatusCmd builds `validator status <ADDR>` → GET
// /stake/validator-status/{addr}.
//
// The upstream response carries a field named `is_old` whose semantics
// are "still in the current validator set" — the opposite of what the
// name suggests. We rename it to `is_current` before rendering either
// the KV or JSON form, and emit a short hint once so scripters who diff
// against the upstream shape notice the rename.
func newStatusCmd() *cobra.Command {
	var fields []string
	cmd := &cobra.Command{
		Use:   "status <ADDR>",
		Short: "Check whether an address is in the current validator set.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			addr, err := normalizeSignerAddress(args[0])
			if err != nil {
				return err
			}
			rest, cfg, err := newRESTClient(cmd)
			if err != nil {
				return err
			}
			body, status, err := rest.Get(cmd.Context(), "/stake/validator-status/"+addr, nil)
			if err != nil {
				return err
			}
			if status == 0 && body == nil {
				return nil
			}
			opts := renderOpts(cmd, cfg, fields)
			m, err := decodeJSONMap(body, "validator-status")
			if err != nil {
				return err
			}
			// Rename upstream misleading field before rendering. Both
			// KV and JSON paths read from the same map, so one rewrite
			// covers both output modes.
			renamed := renameIsOldToIsCurrent(m)

			if opts.JSON {
				if err := render.RenderJSON(cmd.OutOrStdout(), m, opts); err != nil {
					return err
				}
			} else {
				if err := render.RenderKV(cmd.OutOrStdout(), m, opts); err != nil {
					return err
				}
			}
			if renamed {
				_ = render.WriteHint(cmd.ErrOrStderr(), render.HintIsOldRenamed, opts)
			}
			return nil
		},
	}
	cmd.Flags().StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable)")
	return cmd
}

// renameIsOldToIsCurrent rewrites the `is_old` key to `is_current` in m
// in place. Returns true when the rewrite happened so the caller knows
// whether to emit the rename hint.
func renameIsOldToIsCurrent(m map[string]any) bool {
	v, ok := m["is_old"]
	if !ok {
		return false
	}
	m["is_current"] = v
	delete(m, "is_old")
	return true
}
