package validator

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/cmdutil"
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
	return pkg.NewGetCmd(cmdutil.Get{
		Use:         "total-power",
		Short:       "Print aggregate validator voting power.",
		Path:        "/stake/total-power",
		Label:       "total-power",
		FieldsUsage: "pluck one or more fields (repeatable, --json only)",
		RenderBody: func(cmd *cobra.Command, body []byte, _ render.Options) error {
			var resp totalPowerResponse
			if jerr := json.Unmarshal(body, &resp); jerr != nil {
				return fmt.Errorf("decoding total-power: %w", jerr)
			}
			if resp.TotalPower == "" {
				return fmt.Errorf("total-power response missing total_power (body=%q)", cmdutil.Truncate(body, 256))
			}
			_, err := fmt.Fprintln(cmd.OutOrStdout(), resp.TotalPower)
			return err
		},
	})
}
