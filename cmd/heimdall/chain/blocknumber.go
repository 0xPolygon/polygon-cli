package chain

import (
	"fmt"

	"github.com/spf13/cobra"
)

// newBlockNumberCmd builds `block-number` (alias `bn`). Prints
// /status.sync_info.latest_block_height as a bare integer, matching
// `cast block-number`.
func newBlockNumberCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "block-number",
		Aliases: []string{"bn"},
		Short:   "Print the latest CometBFT block height.",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			rpc, _, err := newRPCClient(cmd)
			if err != nil {
				return err
			}
			st, err := fetchStatus(cmd.Context(), rpc)
			if err != nil {
				return err
			}
			if st == nil {
				return nil // --curl
			}
			if st.SyncInfo.LatestBlockHeight == "" {
				return fmt.Errorf("status did not contain latest_block_height")
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), st.SyncInfo.LatestBlockHeight)
			return err
		},
	}
}
