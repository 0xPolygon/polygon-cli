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

package client

import (
	"net"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/maticnetwork/polygon-cli/p2p"
)

type (
	clientParams struct {
		Bootnodes string
		Timeout   string
		Threads   int
		NetworkID int
		NodesFile string
		Database  string
		IsCrawler bool
	}
)

var (
	inputClientParams clientParams
)

// ClientCmd represents the client command. The
var ClientCmd = &cobra.Command{
	Use:   "client [nodes file]",
	Short: "devp2p client that does peer discovery and block/transaction propagation.",
	Long: `Starts a devp2p client that discovers other peers and will receive blocks and
transactions. If only peer discovery is wanted, set the --crawl and the --timeout
flags. If no nodes.json file exists, run echo "{}" >> nodes.json to get started.`,
	Args: cobra.MinimumNArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		inputClientParams.NodesFile = args[0]
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		inputSet, err := p2p.LoadNodesJSON(inputClientParams.NodesFile)
		if err != nil {
			return err
		}

		var cfg discover.Config
		cfg.PrivateKey, _ = crypto.GenerateKey()
		bn, err := p2p.ParseBootnodes(inputClientParams.Bootnodes)
		if err != nil {
			log.Error().Err(err).Msg("Unable to parse bootnodes")
			return err
		}
		cfg.Bootnodes = bn

		db, err := enode.OpenDB(inputClientParams.Database)
		if err != nil {
			return err
		}

		ln := enode.NewLocalNode(db, cfg.PrivateKey)
		socket, err := listen(ln)
		if err != nil {
			return err
		}

		disc, err := discover.ListenV4(socket, ln, cfg)
		if err != nil {
			return err
		}
		defer disc.Close()

		c := newClient(inputSet, disc, disc.RandomNodes())
		c.revalidateInterval = 10 * time.Minute

		timeout, err := time.ParseDuration(inputClientParams.Timeout)
		if err != nil {
			return err
		}

		log.Info().Msg("Starting client")

		output := c.run(timeout, inputClientParams.Threads)
		return p2p.WriteNodesJSON(inputClientParams.NodesFile, output)
	},
}

func init() {
	ClientCmd.PersistentFlags().StringVarP(&inputClientParams.Bootnodes, "bootnodes", "b", "", "Comma separated nodes used for bootstrapping. At least one bootnode is required, so other nodes in the network can discover each other.")
	if err := ClientCmd.MarkPersistentFlagRequired("bootnodes"); err != nil {
		log.Error().Err(err).Msg("Failed to mark bootnodes as required persistent flag")
	}
	ClientCmd.PersistentFlags().StringVarP(&inputClientParams.Timeout, "timeout", "t", "0", "Time limit for node discovery.")
	ClientCmd.PersistentFlags().IntVarP(&inputClientParams.Threads, "parallel", "p", 16, "How many parallel discoveries to attempt.")
	ClientCmd.PersistentFlags().IntVarP(&inputClientParams.NetworkID, "network-id", "n", 0, "Filter discovered nodes by this network id.")
	ClientCmd.PersistentFlags().StringVarP(&inputClientParams.Database, "database", "d", "", "Node database for updating and storing client information.")
	ClientCmd.PersistentFlags().BoolVarP(&inputClientParams.IsCrawler, "crawl", "c", false, "Run the client in crawl only mode.")
}

func listen(ln *enode.LocalNode) (*net.UDPConn, error) {
	addr := "0.0.0.0:0"

	socket, err := net.ListenPacket("udp4", addr)
	if err != nil {
		return nil, err
	}

	// Configure UDP endpoint in ENR from listener address.
	usocket := socket.(*net.UDPConn)
	uaddr := socket.LocalAddr().(*net.UDPAddr)

	if uaddr.IP.IsUnspecified() {
		ln.SetFallbackIP(net.IP{127, 0, 0, 1})
	} else {
		ln.SetFallbackIP(uaddr.IP)
	}

	ln.SetFallbackUDP(uaddr.Port)

	return usocket, nil
}
