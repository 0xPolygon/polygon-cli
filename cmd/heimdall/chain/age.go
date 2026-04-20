package chain

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/render"
)

// newAgeCmd builds `age [HEIGHT]`. Prints the block timestamp via
// the shared timestamp helper, matching `cast age`.
func newAgeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "age [HEIGHT]",
		Short: "Show the timestamp of a CometBFT block.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rpc, _, err := newRPCClient(cmd)
			if err != nil {
				return err
			}
			heightArg := ""
			if len(args) == 1 {
				heightArg = args[0]
			}
			height, err := resolveHeight(cmd.Context(), rpc, heightArg)
			if err != nil {
				return err
			}
			blk, raw, err := fetchBlock(cmd.Context(), rpc, height)
			if err != nil {
				return err
			}
			if raw == nil {
				return nil
			}
			unix, err := unixFromRFC3339Nano(blk.Block.Header.Time)
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), render.AnnotateUnixSeconds(unix))
			return err
		},
	}
}
