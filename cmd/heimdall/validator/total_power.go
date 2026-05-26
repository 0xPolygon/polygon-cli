package validator

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// totalPowerResponse is the shape of GET /stake/total-power.
type totalPowerResponse struct {
	TotalPower string `json:"total_power"`
}

// newTotalPowerCmd builds `validator total-power` → GET
// /stake/total-power. Default output is a bare integer (cheap liveness
// signal); --json emits the wrapper object.
func newTotalPowerCmd() *cobra.Command {
	var fields []string
	cmd := &cobra.Command{
		Use:   "total-power",
		Short: "Print aggregate validator voting power.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			rest, cfg, err := newRESTClient(cmd)
			if err != nil {
				return err
			}
			body, status, err := rest.Get(cmd.Context(), "/stake/total-power", nil)
			if err != nil {
				return err
			}
			if status == 0 && body == nil {
				return nil
			}
			opts := renderOpts(cmd, cfg, fields)
			if opts.JSON {
				m, jerr := decodeJSONMap(body, "total-power")
				if jerr != nil {
					return jerr
				}
				return render.RenderJSON(cmd.OutOrStdout(), m, opts)
			}
			var resp totalPowerResponse
			if jerr := json.Unmarshal(body, &resp); jerr != nil {
				return fmt.Errorf("decoding total-power: %w", jerr)
			}
			if resp.TotalPower == "" {
				return fmt.Errorf("total-power response missing total_power (body=%q)", truncate(body, 256))
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), resp.TotalPower)
			return err
		},
	}
	cmd.Flags().StringArrayVarP(&fields, "field", "f", nil, "pluck one or more fields (repeatable, --json only)")
	return cmd
}
