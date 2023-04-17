package p2p

import (
	"github.com/spf13/cobra"

	"github.com/maticnetwork/polygon-cli/cmd/p2p/client"
	"github.com/maticnetwork/polygon-cli/cmd/p2p/ping"
)

var P2pCmd = &cobra.Command{
	Use:   "p2p",
	Short: "Commands related to devp2p",
}

func init() {
	P2pCmd.AddCommand(client.ClientCmd)
	P2pCmd.AddCommand(ping.PingCmd)
}
