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
	"errors"
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
		Threads   int
		NetworkID int
		NodesFile string
		Database  string
		ProjectID string
		SensorID  string
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
transactions. If no nodes.json file exists, run echo "{}" >> nodes.json to get started.`,
	Args: cobra.MinimumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		inputClientParams.NodesFile = args[0]
		if inputClientParams.NetworkID <= 0 {
			return errors.New("network ID must be greater than zero")
		}
		return nil
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
		socket, err := p2p.Listen(ln)
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

		log.Info().Msg("Starting client")

		c.run(inputClientParams.Threads)
		return nil
	},
}

func init() {
	ClientCmd.PersistentFlags().StringVarP(&inputClientParams.Bootnodes, "bootnodes", "b", "", "Comma separated nodes used for bootstrapping. At least one bootnode is required, so other nodes in the network can discover each other.")
	if err := ClientCmd.MarkPersistentFlagRequired("bootnodes"); err != nil {
		log.Error().Err(err).Msg("Failed to mark bootnodes as required persistent flag")
	}
	ClientCmd.PersistentFlags().IntVarP(&inputClientParams.Threads, "parallel", "p", 16, "How many parallel discoveries to attempt.")
	ClientCmd.PersistentFlags().IntVarP(&inputClientParams.NetworkID, "network-id", "n", 0, "Filter discovered nodes by this network ID.")
	if err := ClientCmd.MarkPersistentFlagRequired("network-id"); err != nil {
		log.Error().Err(err).Msg("Failed to mark network-id as required persistent flag")
	}
	ClientCmd.PersistentFlags().StringVarP(&inputClientParams.Database, "database", "d", "", "Node database for updating and storing client information.")
	ClientCmd.PersistentFlags().StringVarP(&inputClientParams.ProjectID, "project-id", "P", "devtools-sandbox", "GCP project ID.")
	ClientCmd.PersistentFlags().StringVarP(&inputClientParams.SensorID, "sensor-id", "s", "", "Sensor ID.")
	if err := ClientCmd.MarkPersistentFlagRequired("sensor-id"); err != nil {
		log.Error().Err(err).Msg("Failed to mark sensor-id as required persistent flag")
	}
}
