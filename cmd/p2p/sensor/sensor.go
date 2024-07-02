package sensor

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"net/http"
	_ "net/http/pprof"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/forkid"
	"github.com/ethereum/go-ethereum/crypto"
	ethp2p "github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/nat"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/maticnetwork/polygon-cli/p2p"
	"github.com/maticnetwork/polygon-cli/p2p/database"
	"github.com/maticnetwork/polygon-cli/rpctypes"
)

type (
	sensorParams struct {
		Bootnodes                    string
		NetworkID                    uint64
		NodesFile                    string
		TrustedNodesFile             string
		ProjectID                    string
		DatabaseID                   string
		SensorID                     string
		MaxPeers                     int
		MaxDatabaseConcurrency       int
		ShouldWriteBlocks            bool
		ShouldWriteBlockEvents       bool
		ShouldWriteTransactions      bool
		ShouldWriteTransactionEvents bool
		ShouldRunPprof               bool
		PprofPort                    uint
		ShouldRunPrometheus          bool
		PrometheusPort               uint
		KeyFile                      string
		Port                         int
		DiscoveryPort                int
		RPC                          string
		GenesisHash                  string
		ForkID                       []byte
		DialRatio                    int
		NAT                          string
		QuickStart                   bool
		TTL                          time.Duration

		bootnodes    []*enode.Node
		nodes        []*enode.Node
		trustedNodes []*enode.Node
		privateKey   *ecdsa.PrivateKey
		nat          nat.Interface
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
	Long:  "If no nodes.json file exists, it will be created.",
	Args:  cobra.MinimumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		inputSensorParams.NodesFile = args[0]
		inputSensorParams.nodes, err = p2p.ReadNodeSet(inputSensorParams.NodesFile)
		if err != nil {
			log.Warn().Err(err).Msgf("Creating nodes file %v because it does not exist", inputSensorParams.NodesFile)
		}

		if len(inputSensorParams.TrustedNodesFile) > 0 {
			inputSensorParams.trustedNodes, err = p2p.ReadNodeSet(inputSensorParams.TrustedNodesFile)
			if err != nil {
				log.Warn().Err(err).Msgf("Trusted nodes file %v not found", inputSensorParams.TrustedNodesFile)
			}
		}

		if len(inputSensorParams.Bootnodes) > 0 {
			inputSensorParams.bootnodes, err = p2p.ParseBootnodes(inputSensorParams.Bootnodes)
			if err != nil {
				return fmt.Errorf("unable to parse bootnodes: %w", err)
			}
		}

		if inputSensorParams.NetworkID == 0 {
			return errors.New("network ID must be greater than zero")
		}

		if inputSensorParams.ShouldRunPprof {
			go func() {
				addr := fmt.Sprintf(":%v", inputSensorParams.PprofPort)
				if pprofErr := http.ListenAndServe(addr, nil); pprofErr != nil {
					log.Error().Err(pprofErr).Msg("Failed to start pprof")
				}
			}()
		}

		if inputSensorParams.ShouldRunPrometheus {
			go func() {
				http.Handle("/metrics", promhttp.Handler())
				addr := fmt.Sprintf(":%v", inputSensorParams.PrometheusPort)
				if promErr := http.ListenAndServe(addr, nil); promErr != nil {
					log.Error().Err(promErr).Msg("Failed to start Prometheus handler")
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
			DatabaseID:                   inputSensorParams.DatabaseID,
			SensorID:                     inputSensorParams.SensorID,
			MaxConcurrency:               inputSensorParams.MaxDatabaseConcurrency,
			ShouldWriteBlocks:            inputSensorParams.ShouldWriteBlocks,
			ShouldWriteBlockEvents:       inputSensorParams.ShouldWriteBlockEvents,
			ShouldWriteTransactions:      inputSensorParams.ShouldWriteTransactions,
			ShouldWriteTransactionEvents: inputSensorParams.ShouldWriteTransactionEvents,
			TTL:                          inputSensorParams.TTL,
		})

		// Fetch the latest block which will be used later when crafting the status
		// message. This call will only be made once and stored in the head field
		// until the sensor receives a new block it can overwrite it with.
		block, err := getLatestBlock(inputSensorParams.RPC)
		if err != nil {
			return err
		}
		head := p2p.HeadBlock{
			Hash:            block.Hash.ToHash(),
			TotalDifficulty: block.TotalDifficulty.ToBigInt(),
			Number:          block.Number.ToUint64(),
		}

		peersGauge := promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: "sensor",
			Name:      "peers",
			Help:      "The number of peers the sensor is connected to",
		})

		msgCounter := promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: "sensor",
			Name:      "messages",
			Help:      "The number and type of messages the sensor has received",
		}, []string{"code", "message"})

		opts := p2p.EthProtocolOptions{
			Context:     cmd.Context(),
			Database:    db,
			GenesisHash: common.HexToHash(inputSensorParams.GenesisHash),
			RPC:         inputSensorParams.RPC,
			SensorID:    inputSensorParams.SensorID,
			NetworkID:   inputSensorParams.NetworkID,
			Peers:       make(chan *enode.Node),
			Head:        &head,
			HeadMutex:   &sync.RWMutex{},
			ForkID:      forkid.ID{Hash: [4]byte(inputSensorParams.ForkID)},
			MsgCounter:  msgCounter,
		}

		config := ethp2p.Config{
			PrivateKey:     inputSensorParams.privateKey,
			BootstrapNodes: inputSensorParams.bootnodes,
			TrustedNodes:   inputSensorParams.trustedNodes,
			MaxPeers:       inputSensorParams.MaxPeers,
			ListenAddr:     fmt.Sprintf(":%d", inputSensorParams.Port),
			DiscAddr:       fmt.Sprintf(":%d", inputSensorParams.DiscoveryPort),
			DialRatio:      inputSensorParams.DialRatio,
			NAT:            inputSensorParams.nat,
			DiscoveryV4:    true,
			DiscoveryV5:    true,
			Protocols: []ethp2p.Protocol{
				p2p.NewEthProtocol(66, opts),
				p2p.NewEthProtocol(67, opts),
				p2p.NewEthProtocol(68, opts),
			},
		}

		if inputSensorParams.QuickStart {
			config.StaticNodes = inputSensorParams.nodes
		}

		server := ethp2p.Server{Config: config}

		log.Info().Str("enode", server.Self().URLv4()).Msg("Starting sensor")

		// Starting the server isn't actually a blocking call so the sensor needs to
		// have something that waits for it. This is implemented by the for {} loop
		// seen below.
		if err := server.Start(); err != nil {
			return err
		}
		defer server.Stop()

		events := make(chan *ethp2p.PeerEvent)
		sub := server.SubscribeEvents(events)
		defer sub.Unsubscribe()

		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

		peers := make(map[enode.ID]string)
		for _, node := range inputSensorParams.nodes {
			// Because the node URLs can change, map them to the node ID to prevent
			// duplicates.
			peers[node.ID()] = node.URLv4()
		}

		for {
			select {
			case <-ticker.C:
				peersGauge.Set(float64(server.PeerCount()))
			case peer := <-opts.Peers:
				// Update the peer list and the nodes file.
				if _, ok := peers[peer.ID()]; !ok {
					peers[peer.ID()] = peer.URLv4()

					if err := p2p.WritePeers(inputSensorParams.NodesFile, peers); err != nil {
						log.Error().Err(err).Msg("Failed to write nodes to file")
					}
				}
			case <-signals:
				// This gracefully stops the sensor so that the peers can be written to
				// the nodes file.
				log.Info().Msg("Stopping sensor...")
				return nil
			case event := <-events:
				log.Debug().Any("event", event).Send()
			case err := <-sub.Err():
				log.Error().Err(err).Send()
			}
		}
	},
}

// getLatestBlock will get the latest block from an RPC provider.
func getLatestBlock(url string) (*rpctypes.RawBlockResponse, error) {
	client, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	var block rpctypes.RawBlockResponse
	err = client.Call(&block, "eth_getBlockByNumber", "latest", true)
	if err != nil {
		return nil, err
	}

	return &block, nil
}

func init() {
	SensorCmd.Flags().StringVarP(&inputSensorParams.Bootnodes, "bootnodes", "b", "", "Comma separated nodes used for bootstrapping")
	SensorCmd.Flags().Uint64VarP(&inputSensorParams.NetworkID, "network-id", "n", 0, "Filter discovered nodes by this network ID")
	if err := SensorCmd.MarkFlagRequired("network-id"); err != nil {
		log.Error().Err(err).Msg("Failed to mark network-id as required persistent flag")
	}
	SensorCmd.PersistentFlags().StringVarP(&inputSensorParams.ProjectID, "project-id", "p", "", "GCP project ID")
	SensorCmd.PersistentFlags().StringVarP(&inputSensorParams.DatabaseID, "database-id", "d", "", "Datastore database ID")
	SensorCmd.Flags().StringVarP(&inputSensorParams.SensorID, "sensor-id", "s", "", "Sensor ID when writing block/tx events")
	if err := SensorCmd.MarkFlagRequired("sensor-id"); err != nil {
		log.Error().Err(err).Msg("Failed to mark sensor-id as required persistent flag")
	}
	SensorCmd.Flags().IntVarP(&inputSensorParams.MaxPeers, "max-peers", "m", 200, "Maximum number of peers to connect to")
	SensorCmd.Flags().IntVarP(&inputSensorParams.MaxDatabaseConcurrency, "max-db-concurrency", "D", 10000,
		`Maximum number of concurrent database operations to perform. Increasing this
will result in less chance of missing data (i.e. broken pipes) but can
significantly increase memory usage.`)
	SensorCmd.Flags().BoolVarP(&inputSensorParams.ShouldWriteBlocks, "write-blocks", "B", true, "Whether to write blocks to the database")
	SensorCmd.Flags().BoolVar(&inputSensorParams.ShouldWriteBlockEvents, "write-block-events", true, "Whether to write block events to the database")
	SensorCmd.Flags().BoolVarP(&inputSensorParams.ShouldWriteTransactions, "write-txs", "t", true,
		`Whether to write transactions to the database. This option could significantly
increase CPU and memory usage.`)
	SensorCmd.Flags().BoolVar(&inputSensorParams.ShouldWriteTransactionEvents, "write-tx-events", true,
		`Whether to write transaction events to the database. This option could
significantly increase CPU and memory usage.`)
	SensorCmd.Flags().BoolVar(&inputSensorParams.ShouldRunPprof, "pprof", false, "Whether to run pprof")
	SensorCmd.Flags().UintVar(&inputSensorParams.PprofPort, "pprof-port", 6060, "Port pprof runs on")
	SensorCmd.Flags().BoolVar(&inputSensorParams.ShouldRunPrometheus, "prom", true, "Whether to run Prometheus")
	SensorCmd.Flags().UintVar(&inputSensorParams.PrometheusPort, "prom-port", 2112, "Port Prometheus runs on")
	SensorCmd.Flags().StringVarP(&inputSensorParams.KeyFile, "key-file", "k", "", "Private key file")
	SensorCmd.Flags().IntVar(&inputSensorParams.Port, "port", 30303, "TCP network listening port")
	SensorCmd.Flags().IntVar(&inputSensorParams.DiscoveryPort, "discovery-port", 30303, "UDP P2P discovery port")
	SensorCmd.Flags().StringVar(&inputSensorParams.RPC, "rpc", "https://polygon-rpc.com", "RPC endpoint used to fetch the latest block")
	SensorCmd.Flags().StringVar(&inputSensorParams.GenesisHash, "genesis-hash", "0xa9c28ce2141b56c474f1dc504bee9b01eb1bd7d1a507580d5519d4437a97de1b", "The genesis block hash")
	SensorCmd.Flags().BytesHexVar(&inputSensorParams.ForkID, "fork-id", []byte{240, 151, 188, 19}, "The hex encoded fork id (omit the 0x)")
	SensorCmd.Flags().IntVar(&inputSensorParams.DialRatio, "dial-ratio", 0,
		`Ratio of inbound to dialed connections. A dial ratio of 2 allows 1/2 of
connections to be dialed. Setting this to 0 defaults it to 3.`)
	SensorCmd.Flags().StringVar(&inputSensorParams.NAT, "nat", "any", "NAT port mapping mechanism (any|none|upnp|pmp|pmp:<IP>|extip:<IP>)")
	SensorCmd.Flags().BoolVar(&inputSensorParams.QuickStart, "quick-start", false,
		`Whether to load the nodes.json as static nodes to quickly start the network.
This produces faster development cycles but can prevent the sensor from being to
connect to new peers if the nodes.json file is large.`)
	SensorCmd.Flags().StringVar(&inputSensorParams.TrustedNodesFile, "trusted-nodes", "", "Trusted nodes file")
	SensorCmd.Flags().DurationVar(&inputSensorParams.TTL, "ttl", 14*24*time.Hour, "Time to live")
}
