package sensor

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"net/http"
	_ "net/http/pprof"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	ethp2p "github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/nat"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/maticnetwork/polygon-cli/p2p"
	"github.com/maticnetwork/polygon-cli/p2p/database"
)

type (
	sensorParams struct {
		Bootnodes                    string
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
		ShouldRunPprof               bool
		PprofPort                    uint
		KeyFile                      string
		Port                         int
		DiscoveryPort                int
		RPC                          string
		GenesisFile                  string
		GenesisHash                  string
		DialRatio                    int
		NAT                          string

		nodes      []*enode.Node
		privateKey *ecdsa.PrivateKey
		genesis    core.Genesis
		nat        nat.Interface
	}
)

var (
	inputSensorParams sensorParams
)

// SensorCmd represents the sensor command. This is responsible for starting a
// sensor and transmitting blocks and transactions to a database.
var SensorCmd = &cobra.Command{
	Use:   "sensor [nodes file]",
	Short: "Start a devp2p sensor that discovers other peers and will receive blocks and transactions.",
	Long:  "If no nodes.json file exists, run `echo \"[]\" >> nodes.json` to get started.",
	Args:  cobra.MinimumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		inputSensorParams.NodesFile = args[0]
		inputSensorParams.nodes, err = p2p.ReadNodeSet(inputSensorParams.NodesFile)
		if err != nil {
			return err
		}

		if inputSensorParams.NetworkID == 0 {
			return errors.New("network ID must be greater than zero")
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

		inputSensorParams.nat, err = nat.Parse(inputSensorParams.NAT)
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse NAT")
			return err
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		db := database.NewDatastore(cmd.Context(), database.DatastoreOptions{
			ProjectID:                    inputSensorParams.ProjectID,
			SensorID:                     inputSensorParams.SensorID,
			MaxConcurrentWrites:          inputSensorParams.MaxConcurrentDatabaseWrites,
			ShouldWriteBlocks:            inputSensorParams.ShouldWriteBlocks,
			ShouldWriteBlockEvents:       inputSensorParams.ShouldWriteBlockEvents,
			ShouldWriteTransactions:      inputSensorParams.ShouldWriteTransactions,
			ShouldWriteTransactionEvents: inputSensorParams.ShouldWriteTransactionEvents,
		})

		bootnodes, err := p2p.ParseBootnodes(inputSensorParams.Bootnodes)
		if err != nil {
			return fmt.Errorf("unable to parse bootnodes: %w", err)
		}

		opts := p2p.Eth66ProtocolOptions{
			Context:     cmd.Context(),
			Database:    db,
			Genesis:     &inputSensorParams.genesis,
			GenesisHash: common.HexToHash(inputSensorParams.GenesisHash),
			RPC:         inputSensorParams.RPC,
			SensorID:    inputSensorParams.SensorID,
			NetworkID:   inputSensorParams.NetworkID,
			Peers:       make(chan *enode.Node),
		}

		server := ethp2p.Server{
			Config: ethp2p.Config{
				PrivateKey:     inputSensorParams.privateKey,
				BootstrapNodes: bootnodes,
				StaticNodes:    inputSensorParams.nodes,
				MaxPeers:       inputSensorParams.MaxPeers,
				ListenAddr:     fmt.Sprintf(":%d", inputSensorParams.Port),
				DiscAddr:       fmt.Sprintf(":%d", inputSensorParams.DiscoveryPort),
				Protocols:      []ethp2p.Protocol{p2p.NewEth66Protocol(opts)},
				DialRatio:      inputSensorParams.DialRatio,
				NAT:            inputSensorParams.nat,
			},
		}

		log.Info().Str("enode", server.Self().URLv4()).Msg("Starting sensor")
		if err = server.Start(); err != nil {
			return err
		}
		defer server.Stop()

		ticker := time.NewTicker(2 * time.Second)
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

		peers := make(p2p.NodeSet)
		for _, node := range inputSensorParams.nodes {
			peers[node.ID()] = node.URLv4()
		}

		for {
			select {
			case <-ticker.C:
				log.Info().Interface("peers", server.PeerCount()).Send()

				err = p2p.WriteNodeSet(inputSensorParams.NodesFile, peers)
				if err != nil {
					log.Error().Err(err).Msg("Failed to write nodes to file")
				}
			case peer := <-opts.Peers:
				if _, ok := peers[peer.ID()]; !ok {
					peers[peer.ID()] = peer.URLv4()
				}
			case <-signals:
				ticker.Stop()
				log.Info().Msg("Stopping sever...")
				return nil
			}
		}
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
	SensorCmd.PersistentFlags().IntVarP(&inputSensorParams.MaxConcurrentDatabaseWrites, "max-db-writes", "D", 10000,
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
	SensorCmd.PersistentFlags().BoolVar(&inputSensorParams.ShouldRunPprof, "pprof", false, "Whether to run pprof.")
	SensorCmd.PersistentFlags().UintVar(&inputSensorParams.PprofPort, "pprof-port", 6060, "The port to run pprof on.")
	SensorCmd.PersistentFlags().StringVarP(&inputSensorParams.KeyFile, "key-file", "k", "", "The file of the private key. If no key file is found then a key file will be generated.")
	SensorCmd.PersistentFlags().IntVar(&inputSensorParams.Port, "port", 30303, "The TCP network listening port.")
	SensorCmd.PersistentFlags().IntVar(&inputSensorParams.DiscoveryPort, "discovery-port", 30303, "The UDP P2P discovery port.")
	SensorCmd.PersistentFlags().StringVar(&inputSensorParams.RPC, "rpc", "https://polygon-rpc.com", "The RPC endpoint used to fetch the latest block.")
	SensorCmd.PersistentFlags().StringVar(&inputSensorParams.GenesisFile, "genesis", "genesis.json", "The genesis file.")
	SensorCmd.PersistentFlags().StringVar(&inputSensorParams.GenesisHash, "genesis-hash", "0xa9c28ce2141b56c474f1dc504bee9b01eb1bd7d1a507580d5519d4437a97de1b", "The genesis block hash.")
	SensorCmd.PersistentFlags().IntVar(&inputSensorParams.DialRatio, "dial-ratio", 0,
		`The ratio of inbound to dialed connections. A dial ratio of 2 allows 1/2 of
connections to be dialed. Setting this to 0 defaults it to 3.`)
	SensorCmd.PersistentFlags().StringVar(&inputSensorParams.NAT, "nat", "any", "The NAT port mapping mechanism (any|none|upnp|pmp|pmp:<IP>|extip:<IP>).")
}
