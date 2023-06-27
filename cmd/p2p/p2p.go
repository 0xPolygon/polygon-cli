package p2p

import (
	"github.com/spf13/cobra"

	"github.com/maticnetwork/polygon-cli/cmd/p2p/crawl"
	"github.com/maticnetwork/polygon-cli/cmd/p2p/ping"
	"github.com/maticnetwork/polygon-cli/cmd/p2p/sensor"
)

var P2pCmd = &cobra.Command{
	Use:   "p2p",
	Short: "Commands related to devp2p",
}

func init() {
	P2pCmd.AddCommand(sensor.SensorCmd)
	P2pCmd.AddCommand(crawl.CrawlCmd)
	P2pCmd.AddCommand(ping.PingCmd)
}
