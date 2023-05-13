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

package sensor

import (
	"errors"
	"time"

	"net/http"
	_ "net/http/pprof"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/maticnetwork/polygon-cli/p2p"
)

type (
	sensorParams struct {
		Bootnodes string
		Threads   int
		NetworkID int
		NodesFile string
		Database  string
		ProjectID string
		SensorID  string
		MaxPeers  int
	}
)

var (
	inputSensorParams sensorParams
)

// SensorCmd represents the sensor command. This is responsible for starting a
// sensor and transmitting blocks and transactions to a database.
var SensorCmd = &cobra.Command{
	Use:   "sensor [nodes file]",
	Short: "devp2p sensor that does peer discovery and block/transaction propagation.",
	Long: `Starts a devp2p sensor that discovers other peers and will receive blocks and
transactions. If no nodes.json file exists, run echo "{}" >> nodes.json to get started.`,
	Args: cobra.MinimumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		inputSensorParams.NodesFile = args[0]
		if inputSensorParams.NetworkID <= 0 {
			return errors.New("network ID must be greater than zero")
		}
		go func() {
			http.ListenAndServe("localhost:6060", nil)
		}()

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		inputSet, err := p2p.LoadNodesJSON(inputSensorParams.NodesFile)
		if err != nil {
			return err
		}

		var cfg discover.Config
		cfg.PrivateKey, _ = crypto.GenerateKey()
		bn, err := p2p.ParseBootnodes(inputSensorParams.Bootnodes)
		if err != nil {
			log.Error().Err(err).Msg("Unable to parse bootnodes")
			return err
		}
		cfg.Bootnodes = bn

		db, err := enode.OpenDB(inputSensorParams.Database)
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

		c := newSensor(inputSet, disc, disc.RandomNodes())
		c.revalidateInterval = 10 * time.Minute

		log.Info().Msg("Starting client")

		c.run(inputSensorParams.Threads)
		return nil
	},
}

func init() {
	SensorCmd.PersistentFlags().StringVarP(&inputSensorParams.Bootnodes, "bootnodes", "b", "", "Comma separated nodes used for bootstrapping. At least one bootnode is required, so other nodes in the network can discover each other.")
	if err := SensorCmd.MarkPersistentFlagRequired("bootnodes"); err != nil {
		log.Error().Err(err).Msg("Failed to mark bootnodes as required persistent flag")
	}
	SensorCmd.PersistentFlags().IntVarP(&inputSensorParams.Threads, "parallel", "p", 16, "How many parallel discoveries to attempt.")
	SensorCmd.PersistentFlags().IntVarP(&inputSensorParams.NetworkID, "network-id", "n", 0, "Filter discovered nodes by this network ID.")
	if err := SensorCmd.MarkPersistentFlagRequired("network-id"); err != nil {
		log.Error().Err(err).Msg("Failed to mark network-id as required persistent flag")
	}
	SensorCmd.PersistentFlags().StringVarP(&inputSensorParams.Database, "database", "d", "", "Node database for updating and storing client information.")
	SensorCmd.PersistentFlags().StringVarP(&inputSensorParams.ProjectID, "project-id", "P", "", "GCP project ID.")
	SensorCmd.PersistentFlags().StringVarP(&inputSensorParams.SensorID, "sensor-id", "s", "", "Sensor ID.")
	if err := SensorCmd.MarkPersistentFlagRequired("sensor-id"); err != nil {
		log.Error().Err(err).Msg("Failed to mark sensor-id as required persistent flag")
	}
	SensorCmd.PersistentFlags().IntVarP(&inputSensorParams.MaxPeers, "max-peers", "m", 200, "Maximum number of peers to connect to.")
}
