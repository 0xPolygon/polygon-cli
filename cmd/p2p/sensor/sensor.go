package sensor

import (
	"errors"
	"fmt"
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
		Bootnodes                    string
		Threads                      int
		NetworkID                    uint64
		NodesFile                    string
		Database                     string
		ProjectID                    string
		SensorID                     string
		MaxPeers                     int
		MaxConcurrentDatabaseWrites  int
		ShouldWriteBlocks            bool
		ShouldWriteBlockEvents       bool
		ShouldWriteTransactions      bool
		ShouldWriteTransactionEvents bool
		RevalidationInterval         string
		revalidationInterval         time.Duration
		ShouldRunPprof               bool
		PprofPort                    uint
	}
)

var (
	inputSensorParams sensorParams
)

// SensorCmd represents the sensor command. This is responsible for starting a
// sensor and transmitting blocks and transactions to a database.
var SensorCmd = &cobra.Command{
	Use:   "sensor [nodes file]",
	Short: "Start a devp2p sensor that discovers other peers and will receive blocks and transactions. ",
	Long:  "If no nodes.json file exists, run `echo \"{}\" >> nodes.json` to get started.",
	Args:  cobra.MinimumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		inputSensorParams.NodesFile = args[0]
		if inputSensorParams.NetworkID == 0 {
			return errors.New("network ID must be greater than zero")
		}

		inputSensorParams.revalidationInterval, err = time.ParseDuration(inputSensorParams.RevalidationInterval)
		if err != nil {
			return err
		}

		if inputSensorParams.ShouldRunPprof {
			go func() {
				if err := http.ListenAndServe(fmt.Sprintf("localhost:%v", inputSensorParams.PprofPort), nil); err != nil {
					log.Error().Err(err).Msg("Failed to start pprof")
				}
			}()
		}

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
			return fmt.Errorf("unable to parse bootnodes: %w", err)
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
		c.revalidateInterval = inputSensorParams.revalidationInterval

		log.Info().Msg("Starting sensor")

		c.run(inputSensorParams.Threads)
		return nil
	},
}

func init() {
	SensorCmd.PersistentFlags().StringVarP(&inputSensorParams.Bootnodes, "bootnodes", "b", "",
		`Comma separated nodes used for bootstrapping. At least one bootnode is
required, so other nodes in the network can discover each other.`)
	if err := SensorCmd.MarkPersistentFlagRequired("bootnodes"); err != nil {
		log.Error().Err(err).Msg("Failed to mark bootnodes as required persistent flag")
	}
	SensorCmd.PersistentFlags().IntVarP(&inputSensorParams.Threads, "parallel", "p", 16, "How many parallel discoveries to attempt.")
	SensorCmd.PersistentFlags().Uint64VarP(&inputSensorParams.NetworkID, "network-id", "n", 0, "Filter discovered nodes by this network ID.")
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
	SensorCmd.PersistentFlags().IntVarP(&inputSensorParams.MaxConcurrentDatabaseWrites, "max-db-writes", "D", 100,
		`The maximum number of concurrent database writes to perform. Increasing
this will result in less chance of missing data (i.e. broken pipes) but
can significantly increase memory usage.`)
	SensorCmd.PersistentFlags().BoolVarP(&inputSensorParams.ShouldWriteBlocks, "write-blocks", "B", true, "Whether to write blocks to the database.")
	SensorCmd.PersistentFlags().BoolVar(&inputSensorParams.ShouldWriteBlockEvents, "write-block-events", true, "Whether to write block events to the database.")
	SensorCmd.PersistentFlags().BoolVarP(&inputSensorParams.ShouldWriteTransactions, "write-txs", "t", true,
		`Whether to write transactions to the database. This option could significantly
increase CPU and memory usage.`)
	SensorCmd.PersistentFlags().BoolVar(&inputSensorParams.ShouldWriteTransactionEvents, "write-tx-events", true,
		`Whether to write transaction events to the database. This option could significantly
increase CPU and memory usage.`)
	SensorCmd.PersistentFlags().StringVarP(&inputSensorParams.RevalidationInterval, "revalidation-interval", "r", "10m", "The amount of time it takes to retry connecting to a failed peer.")
	SensorCmd.PersistentFlags().BoolVar(&inputSensorParams.ShouldRunPprof, "pprof", false, "Whether to run pprof.")
	SensorCmd.PersistentFlags().UintVar(&inputSensorParams.PprofPort, "pprof-port", 6060, "The port to run pprof on.")
}
