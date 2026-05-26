package checkpoint

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newGetCmd builds `checkpoint get <ID>` → GET /checkpoints/{id}. The
// same code path is re-entered from CheckpointCmd's RunE when a bare
// integer is provided (`checkpoint 38871`).
func newGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <ID>",
		Short: "Fetch one checkpoint by id.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(cmd, args[0])
		},
	}
	return cmd
}

// runGet is the shared implementation used by both `checkpoint get
// <ID>` and the bare-integer CheckpointCmd shorthand.
func runGet(cmd *cobra.Command, idArg string) error {
	id, err := strconv.ParseUint(idArg, 10, 64)
	if err != nil {
		return &client.UsageError{Msg: fmt.Sprintf("checkpoint id must be a positive integer, got %q", idArg)}
	}
	rest, cfg, err := newRESTClient(cmd)
	if err != nil {
		return err
	}
	body, status, err := rest.Get(cmd.Context(), fmt.Sprintf("/checkpoints/%d", id), nil)
	if err != nil {
		return err
	}
	if status == 0 && body == nil {
		return nil
	}
	// --field is valid here but not a flag on the umbrella; honour the
	// global --json and --raw only.
	opts := renderOpts(cmd, cfg, nil)
	m, err := decodeJSONMap(body, "checkpoint")
	if err != nil {
		return err
	}
	if opts.JSON {
		return render.RenderJSON(cmd.OutOrStdout(), m, opts)
	}
	return renderCheckpointKV(cmd, m, opts)
}
