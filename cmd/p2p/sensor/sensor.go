package sensor

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"net/http"
	_ "net/http/pprof"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	ethp2p "github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/nat"
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
		KeyFile                      string
		privateKey                   *ecdsa.PrivateKey
		Port                         int
		RPC                          string
		genesis                      core.Genesis
		GenesisFile                  string
		GenesisHash                  string
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
				if pprofErr := http.ListenAndServe(fmt.Sprintf("localhost:%v", inputSensorParams.PprofPort), nil); pprofErr != nil {
					log.Error().Err(pprofErr).Msg("Failed to start pprof")
				}
			}()
		}

		inputSensorParams.privateKey, err = crypto.GenerateKey()
		if err != nil {
			return err
		}

		if len(inputSensorParams.KeyFile) > 0 {
			var privateKey *ecdsa.PrivateKey
			privateKey, err = crypto.LoadECDSA(inputSensorParams.KeyFile)

			if err != nil {
				log.Warn().Err(err).Msg("Key file was not found, generating a new key file")

				err = crypto.SaveECDSA(inputSensorParams.KeyFile, inputSensorParams.privateKey)
				if err != nil {
					return err
				}
			} else {
				inputSensorParams.privateKey = privateKey
			}
		}

		inputSensorParams.genesis, err = loadGenesis(inputSensorParams.GenesisFile)
		if err != nil {
			log.Error().Err(err).Msg("Failed to load genesis file")
			return err
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		inputSet, err := p2p.LoadNodesJSON(inputSensorParams.NodesFile)
		if err != nil {
			return err
		}

		var cfg discover.Config
		cfg.PrivateKey = inputSensorParams.privateKey
		bn, err := p2p.ParseBootnodes(inputSensorParams.Bootnodes)
		if err != nil {
			return fmt.Errorf("unable to parse bootnodes: %w", err)
		}
		cfg.Bootnodes = bn

		server := ethp2p.Server{
			Config: ethp2p.Config{
				PrivateKey: inputSensorParams.privateKey,
				MaxPeers:   inputSensorParams.MaxPeers,
				ListenAddr: fmt.Sprintf(":%v", inputSensorParams.Port),
				Protocols: []ethp2p.Protocol{p2p.NewEth66Protocol(
					&inputSensorParams.genesis,
					common.HexToHash(inputSensorParams.GenesisHash),
					inputSensorParams.RPC,
					inputSensorParams.NetworkID,
				)},
				NoDial:      true,
				NoDiscovery: true,
				NAT:         nat.Any(),
			},
		}
		if err = server.Start(); err != nil {
			return err
		}
		defer server.Stop()

		ln := server.LocalNode()
		socket, err := p2p.Listen(ln, inputSensorParams.Port)
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

		log.Info().Str("enode", server.Self().URLv4()).Msg("Starting sensor")

		c.run(inputSensorParams.Threads)
		return nil
	},
}

func loadGenesis(genesisFile string) (core.Genesis, error) {
	chainConfig, err := os.ReadFile(genesisFile)

	if err != nil {
		return core.Genesis{}, err
	}
	var gen core.Genesis
	if err := json.Unmarshal(chainConfig, &gen); err != nil {
		return core.Genesis{}, err
	}
	return gen, nil
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
	SensorCmd.PersistentFlags().StringVarP(&inputSensorParams.KeyFile, "key-file", "k", "", "The file of the private key. If no key file is found then a key file will be generated.")
	SensorCmd.PersistentFlags().IntVar(&inputSensorParams.Port, "port", 30303, "The sensor's TCP and discovery port.")
	SensorCmd.PersistentFlags().StringVar(&inputSensorParams.RPC, "rpc", "https://polygon-rpc.com", "The RPC endpoint.")
	SensorCmd.PersistentFlags().StringVar(&inputSensorParams.GenesisFile, "genesis", "genesis.json", "The genesis file.")
	SensorCmd.PersistentFlags().StringVar(&inputSensorParams.GenesisHash, "genesis-hash", "0xa9c28ce2141b56c474f1dc504bee9b01eb1bd7d1a507580d5519d4437a97de1b", "The genesis block hash.")
}
