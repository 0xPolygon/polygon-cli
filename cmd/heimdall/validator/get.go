package validator

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newGetCmd builds `validator get <ID>` → GET /stake/validator/{id}.
// The same code path is re-entered from ValidatorCmd's RunE when a bare
// integer is provided (`validator 4`).
func newGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <ID>",
		Short: "Fetch one validator by numeric id.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(cmd, args[0])
		},
	}
	return cmd
}

// runGet is the shared implementation used by both `validator get <ID>`
// and the bare-integer ValidatorCmd shorthand.
func runGet(cmd *cobra.Command, idArg string) error {
	id, err := strconv.ParseUint(idArg, 10, 64)
	if err != nil {
		return &client.UsageError{Msg: fmt.Sprintf("validator id must be a positive integer, got %q", idArg)}
	}
	rest, cfg, err := newRESTClient(cmd)
	if err != nil {
		return err
	}
	body, status, err := rest.Get(cmd.Context(), fmt.Sprintf("/stake/validator/%d", id), nil)
	if err != nil {
		return err
	}
	if status == 0 && body == nil {
		return nil
	}
	opts := renderOpts(cmd, cfg, nil)
	m, err := decodeJSONMap(body, "validator")
	if err != nil {
		return err
	}
	if opts.JSON {
		return render.RenderJSON(cmd.OutOrStdout(), m, opts)
	}
	return renderValidatorKV(cmd, m, opts)
}

// renderValidatorKV unwraps the { "validator": { … } } envelope for KV
// output. If the envelope is absent the map is rendered as-is.
func renderValidatorKV(cmd *cobra.Command, m map[string]any, opts render.Options) error {
	if inner, ok := m["validator"].(map[string]any); ok {
		return render.RenderKV(cmd.OutOrStdout(), inner, opts)
	}
	return render.RenderKV(cmd.OutOrStdout(), m, opts)
}
