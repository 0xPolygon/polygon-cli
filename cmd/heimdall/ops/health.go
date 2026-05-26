package ops

import (
	"fmt"

	"github.com/spf13/cobra"
)

// newHealthCmd builds `ops health`. Prints "OK" on success, returns a
// non-nil error otherwise (cobra bubbles it up and callers map to a
// cast-style exit code via client.ExitCode).
func newHealthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "health",
		Short: "Probe CometBFT /health; exit 0 on success.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			rpc, _, err := newRPCClient(cmd)
			if err != nil {
				return err
			}
			raw, err := callEmpty(cmd.Context(), rpc, "health")
			if err != nil {
				return err
			}
			if raw == nil {
				return nil // --curl
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), "OK")
			return err
		},
	}
	return cmd
}
