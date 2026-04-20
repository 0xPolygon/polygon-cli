package chain

import (
	"fmt"

	"github.com/spf13/cobra"
)

// newChainIDCmd builds `chain-id` (alias `ci`). Short-circuits
// CometBFT /status and prints the network id.
func newChainIDCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "chain-id",
		Aliases: []string{"ci"},
		Short:   "Print the CometBFT chain id.",
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
				return nil
			}
			if st.NodeInfo.Network == "" {
				return fmt.Errorf("status did not contain node_info.network")
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), st.NodeInfo.Network)
			return err
		},
	}
}
