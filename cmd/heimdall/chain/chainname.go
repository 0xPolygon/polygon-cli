package chain

import (
	"fmt"

	"github.com/spf13/cobra"
)

// newChainCmd builds `chain`. Looks up the /status chain id against
// the built-in table to print a human-readable chain name. Unknown
// ids fall through to `unknown chain <id>`.
func newChainCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "chain",
		Short: "Print the human-readable chain name.",
		Args:  cobra.NoArgs,
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
			id := st.NodeInfo.Network
			if name, ok := chainNames[id]; ok {
				_, err = fmt.Fprintln(cmd.OutOrStdout(), name)
				return err
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "unknown chain %s\n", id)
			return err
		},
	}
}
