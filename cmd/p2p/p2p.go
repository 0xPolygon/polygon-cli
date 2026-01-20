package p2p

import (
	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/cmd/p2p/crawl"
	"github.com/0xPolygon/polygon-cli/cmd/p2p/nodelist"
	"github.com/0xPolygon/polygon-cli/cmd/p2p/ping"
	"github.com/0xPolygon/polygon-cli/cmd/p2p/query"
	"github.com/0xPolygon/polygon-cli/cmd/p2p/sensor"
)

var P2pCmd = &cobra.Command{
	Use:   "p2p",
	Short: "Set of commands related to devp2p.",
}

func init() {
	P2pCmd.AddCommand(crawl.CrawlCmd)
	P2pCmd.AddCommand(nodelist.NodeListCmd)
	P2pCmd.AddCommand(ping.PingCmd)
	P2pCmd.AddCommand(sensor.SensorCmd)
	P2pCmd.AddCommand(query.QueryCmd)
}
