/*
Copyright © 2022 Polygon <engineering@polygon.technology>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package crawl

import (
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/maticnetwork/polygon-cli/p2p"
)

type (
	crawlParams struct {
		Bootnodes string
		Timeout   string
		Threads   int
		NetworkID int
		NodesFile string
		Database  string
	}
)

var (
	inputCrawlParams crawlParams
)

// crawlCmd represents the crawl command
var CrawlCmd = &cobra.Command{
	Use:   "crawl [nodes file]",
	Short: "Crawl a network",
	Long:  `This is a basic function to crawl a network.`,
	Args:  cobra.MinimumNArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		inputCrawlParams.NodesFile = args[0]
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		inputSet, err := p2p.LoadNodesJSON(inputCrawlParams.NodesFile)
		if err != nil {
			return err
		}

		var cfg discover.Config
		cfg.PrivateKey, _ = crypto.GenerateKey()
		bn, err := p2p.ParseBootnodes(inputCrawlParams.Bootnodes)
		if err != nil {
			log.Error().Err(err).Msg("Unable to parse bootnodes")
			return err
		}
		cfg.Bootnodes = bn

		db, err := enode.OpenDB(inputCrawlParams.Database)
		if err != nil {
			return err
		}

		ln := enode.NewLocalNode(db, cfg.PrivateKey)
		socket, err := p2p.Listen(ln)
		if err != nil {
			return err
		}

		disc, err := discover.ListenV4(socket, ln, cfg)
		if err != nil {
			return err
		}
		defer disc.Close()

		c := newCrawler(inputSet, disc, disc.RandomNodes())
		c.revalidateInterval = 10 * time.Minute

		timeout, err := time.ParseDuration(inputCrawlParams.Timeout)
		if err != nil {
			return err
		}

		log.Info().Msg("Starting crawl")

		output := c.run(timeout, inputCrawlParams.Threads)
		return p2p.WriteNodesJSON(inputCrawlParams.NodesFile, output)
	},
}

func init() {
	CrawlCmd.PersistentFlags().StringVarP(&inputCrawlParams.Bootnodes, "bootnodes", "b", "", "Comma separated nodes used for bootstrapping. At least one bootnode is required, so other nodes in the network can discover each other.")
	if err := CrawlCmd.MarkPersistentFlagRequired("bootnodes"); err != nil {
		log.Error().Err(err).Msg("Failed to mark bootnodes as required persistent flag")
	}
	CrawlCmd.PersistentFlags().StringVarP(&inputCrawlParams.Timeout, "timeout", "t", "30m0s", "Time limit for the crawl.")
	CrawlCmd.PersistentFlags().IntVarP(&inputCrawlParams.Threads, "parallel", "p", 16, "How many parallel discoveries to attempt.")
	CrawlCmd.PersistentFlags().IntVarP(&inputCrawlParams.NetworkID, "network-id", "n", 0, "Filter discovered nodes by this network id.")
	CrawlCmd.PersistentFlags().StringVarP(&inputCrawlParams.Database, "database", "d", "", "Node database for updating and storing client information.")
}
