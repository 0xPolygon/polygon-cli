package sensor

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"net/http"
	_ "net/http/pprof"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/forkid"
	"github.com/ethereum/go-ethereum/crypto"
	ethp2p "github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/dnsdisc"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/nat"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/p2p"
	"github.com/0xPolygon/polygon-cli/p2p/database"
	"github.com/0xPolygon/polygon-cli/rpctypes"
)

type (
	sensorParams struct {
		Bootnodes                    string
		NetworkID                    uint64
		NodesFile                    string
		StaticNodesFile              string
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
		ShouldWritePeers             bool
		ShouldRunPprof               bool
		PprofPort                    uint
		ShouldRunPrometheus          bool
		PrometheusPort               uint
		APIPort                      uint
		RPCPort                      uint
		KeyFile                      string
		PrivateKey                   string
		Port                         int
		DiscoveryPort                int
		RPC                          string
		GenesisHash                  string
		ForkID                       []byte
		DialRatio                    int
		NAT                          string
		TTL                          time.Duration
		DiscoveryDNS                 string
		Database                     string
		NoDiscovery                  bool
		RequestsCache                p2p.CacheOptions
		ParentsCache                 p2p.CacheOptions
		BlocksCache                  p2p.CacheOptions

		bootnodes    []*enode.Node
		staticNodes  []*enode.Node
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
		_, err = p2p.ReadNodeSet(inputSensorParams.NodesFile)
		if err != nil {
			log.Warn().Err(err).Msgf("Creating nodes file %v because it does not exist", inputSensorParams.NodesFile)
		}

		if len(inputSensorParams.StaticNodesFile) > 0 {
			inputSensorParams.staticNodes, err = p2p.ReadNodeSet(inputSensorParams.StaticNodesFile)
			if err != nil {
				log.Warn().Err(err).Msgf("Static nodes file %v not found", inputSensorParams.StaticNodesFile)
			}
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

		if len(inputSensorParams.PrivateKey) > 0 {
			inputSensorParams.privateKey, err = crypto.HexToECDSA(inputSensorParams.PrivateKey)
			if err != nil {
				log.Error().Err(err).Msg("Failed to parse PrivateKey")
				return err
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
		db, err := newDatabase(cmd.Context())
		if err != nil {
			return err
		}

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
			Help:      "The number and type of messages the sensor has sent and received",
		}, []string{"message", "url", "name", "direction"})

		// Create peer connection manager for broadcasting transactions
		// and managing the global blocks cache
		conns := p2p.NewConns(p2p.ConnsOptions{
			BlocksCache: inputSensorParams.BlocksCache,
			Head:        head,
		})

		opts := p2p.EthProtocolOptions{
			Context:       cmd.Context(),
			Database:      db,
			GenesisHash:   common.HexToHash(inputSensorParams.GenesisHash),
			RPC:           inputSensorParams.RPC,
			SensorID:      inputSensorParams.SensorID,
			NetworkID:     inputSensorParams.NetworkID,
			Conns:         conns,
			ForkID:        forkid.ID{Hash: [4]byte(inputSensorParams.ForkID)},
			MsgCounter:    msgCounter,
			RequestsCache: inputSensorParams.RequestsCache,
			ParentsCache:  inputSensorParams.ParentsCache,
		}

		config := ethp2p.Config{
			PrivateKey:     inputSensorParams.privateKey,
			BootstrapNodes: inputSensorParams.bootnodes,
			StaticNodes:    inputSensorParams.staticNodes,
			TrustedNodes:   inputSensorParams.trustedNodes,
			MaxPeers:       inputSensorParams.MaxPeers,
			ListenAddr:     fmt.Sprintf(":%d", inputSensorParams.Port),
			DiscAddr:       fmt.Sprintf(":%d", inputSensorParams.DiscoveryPort),
			DialRatio:      inputSensorParams.DialRatio,
			NAT:            inputSensorParams.nat,
			DiscoveryV4:    !inputSensorParams.NoDiscovery,
			DiscoveryV5:    !inputSensorParams.NoDiscovery,
			Protocols: []ethp2p.Protocol{
				p2p.NewEthProtocol(66, opts),
				p2p.NewEthProtocol(67, opts),
				p2p.NewEthProtocol(68, opts),
			},
		}

		server := ethp2p.Server{Config: config}

		log.Info().Str("enode", server.Self().URLv4()).Msg("Starting sensor")

		// Starting the server isn't actually a blocking call so the sensor needs to
		// have something that waits for it. This is implemented by the for {} loop
		// seen below.
		if err = server.Start(); err != nil {
			return err
		}
		defer server.Stop()

		events := make(chan *ethp2p.PeerEvent)
		sub := server.SubscribeEvents(events)
		defer sub.Unsubscribe()

		ticker := time.NewTicker(2 * time.Second) // Ticker for recurring tasks every 2 seconds.
		ticker1h := time.NewTicker(time.Hour)     // Ticker for running DNS discovery every hour.
		defer ticker.Stop()
		defer ticker1h.Stop()

		dnsLock := make(chan struct{}, 1)
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

		if inputSensorParams.ShouldRunPprof {
			go handlePprof()
		}

		if inputSensorParams.ShouldRunPrometheus {
			go handlePrometheus()
		}

		go handleAPI(&server, msgCounter, conns)

		// Start the RPC server for receiving transactions
		go handleRPC(conns, inputSensorParams.NetworkID)

		// Run DNS discovery immediately at startup.
		go handleDNSDiscovery(&server, dnsLock)

		for {
			select {
			case <-ticker.C:
				peersGauge.Set(float64(server.PeerCount()))
				db.WritePeers(cmd.Context(), server.Peers(), time.Now())

				urls := []string{}
				for _, peer := range server.Peers() {
					urls = append(urls, peer.Node().URLv4())
				}

				if err := removePeerMessages(msgCounter, urls); err != nil {
					log.Error().Err(err).Msg("Failed to clean up peer messages")
				}

				if err := p2p.WritePeers(inputSensorParams.NodesFile, urls); err != nil {
					log.Error().Err(err).Msg("Failed to write nodes to file")
				}
			case <-ticker1h.C:
				go handleDNSDiscovery(&server, dnsLock)
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

// handlePprof starts a server for performance profiling using pprof on the
// specified port. This allows for real-time monitoring and analysis of the
// sensor's performance. The port number is configured through
// inputSensorParams.PprofPort. An error is logged if the server fails to start.
func handlePprof() {
	addr := fmt.Sprintf(":%d", inputSensorParams.PprofPort)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Error().Err(err).Msg("Failed to start pprof")
	}
}

// handlePrometheus starts a server to expose Prometheus metrics at the /metrics
// endpoint. This enables Prometheus to scrape and collect metrics data for
// monitoring purposes. The port number is configured through
// inputSensorParams.PrometheusPort. An error is logged if the server fails to
// start.
func handlePrometheus() {
	http.Handle("/metrics", promhttp.Handler())
	addr := fmt.Sprintf(":%d", inputSensorParams.PrometheusPort)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Error().Err(err).Msg("Failed to start Prometheus handler")
	}
}

// handleDNSDiscovery performs DNS-based peer discovery and adds new peers to
// the p2p server. It uses an iterator to discover peers incrementally rather
// than loading all nodes at once. The lock channel prevents concurrent runs.
func handleDNSDiscovery(server *ethp2p.Server, lock chan struct{}) {
	if len(inputSensorParams.DiscoveryDNS) == 0 {
		return
	}

	select {
	case lock <- struct{}{}:
		defer func() { <-lock }()
	default:
		log.Warn().Msg("DNS discovery already running, skipping")
		return
	}

	log.Info().
		Str("discovery-dns", inputSensorParams.DiscoveryDNS).
		Msg("Starting DNS discovery")

	client := dnsdisc.NewClient(dnsdisc.Config{})
	iter, err := client.NewIterator(inputSensorParams.DiscoveryDNS)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create DNS discovery iterator")
		return
	}
	defer iter.Close()

	// Add DNS-discovered peers using the iterator.
	count := 0
	for iter.Next() {
		node := iter.Node()
		log.Debug().
			Str("enode", node.URLv4()).
			Msg("Discovered peer through DNS")

		// Add the peer to the static node set. The server itself handles whether to
		// connect to the peer if it's already connected. If a node is part of the
		// static peer set, the server will handle reconnecting after disconnects.
		server.AddPeer(node)
		count++
	}

	log.Info().
		Int("discovered_peers", count).
		Msg("Finished DNS discovery")
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

// newDatabase creates and configures the appropriate database backend based
// on the sensor parameters.
func newDatabase(ctx context.Context) (database.Database, error) {
	switch inputSensorParams.Database {
	case "datastore":
		return database.NewDatastore(ctx, database.DatastoreOptions{
			ProjectID:                    inputSensorParams.ProjectID,
			DatabaseID:                   inputSensorParams.DatabaseID,
			SensorID:                     inputSensorParams.SensorID,
			ChainID:                      inputSensorParams.NetworkID,
			MaxConcurrency:               inputSensorParams.MaxDatabaseConcurrency,
			ShouldWriteBlocks:            inputSensorParams.ShouldWriteBlocks,
			ShouldWriteBlockEvents:       inputSensorParams.ShouldWriteBlockEvents,
			ShouldWriteTransactions:      inputSensorParams.ShouldWriteTransactions,
			ShouldWriteTransactionEvents: inputSensorParams.ShouldWriteTransactionEvents,
			ShouldWritePeers:             inputSensorParams.ShouldWritePeers,
			TTL:                          inputSensorParams.TTL,
		}), nil
	case "json":
		return database.NewJSONDatabase(database.JSONDatabaseOptions{
			SensorID:                     inputSensorParams.SensorID,
			ChainID:                      inputSensorParams.NetworkID,
			MaxConcurrency:               inputSensorParams.MaxDatabaseConcurrency,
			ShouldWriteBlocks:            inputSensorParams.ShouldWriteBlocks,
			ShouldWriteBlockEvents:       inputSensorParams.ShouldWriteBlockEvents,
			ShouldWriteTransactions:      inputSensorParams.ShouldWriteTransactions,
			ShouldWriteTransactionEvents: inputSensorParams.ShouldWriteTransactionEvents,
			ShouldWritePeers:             inputSensorParams.ShouldWritePeers,
		}), nil
	case "none":
		return database.NoDatabase(), nil
	default:
		return nil, fmt.Errorf("invalid database option: %s", inputSensorParams.Database)
	}
}

func init() {
	f := SensorCmd.Flags()
	f.StringVarP(&inputSensorParams.Bootnodes, "bootnodes", "b", "", "comma separated nodes used for bootstrapping")
	f.Uint64VarP(&inputSensorParams.NetworkID, "network-id", "n", 0, "filter discovered nodes by this network ID")
	if err := SensorCmd.MarkFlagRequired("network-id"); err != nil {
		log.Error().Err(err).Msg("Failed to mark network-id as required persistent flag")
	}
	f.StringVarP(&inputSensorParams.ProjectID, "project-id", "p", "", "GCP project ID")
	f.StringVarP(&inputSensorParams.DatabaseID, "database-id", "d", "", "datastore database ID")
	f.StringVarP(&inputSensorParams.SensorID, "sensor-id", "s", "", "sensor ID when writing block/tx events")
	if err := SensorCmd.MarkFlagRequired("sensor-id"); err != nil {
		log.Error().Err(err).Msg("Failed to mark sensor-id as required persistent flag")
	}
	f.IntVarP(&inputSensorParams.MaxPeers, "max-peers", "m", 2000, "maximum number of peers to connect to")
	f.IntVarP(&inputSensorParams.MaxDatabaseConcurrency, "max-db-concurrency", "D", 10000,
		`maximum number of concurrent database operations to perform (increasing this
will result in less chance of missing data but can significantly increase memory usage)`)
	f.BoolVarP(&inputSensorParams.ShouldWriteBlocks, "write-blocks", "B", true, "write blocks to database")
	f.BoolVar(&inputSensorParams.ShouldWriteBlockEvents, "write-block-events", true, "write block events to database")
	f.BoolVarP(&inputSensorParams.ShouldWriteTransactions, "write-txs", "t", true,
		`write transactions to database (this option can significantly increase CPU and memory usage)`)
	f.BoolVar(&inputSensorParams.ShouldWriteTransactionEvents, "write-tx-events", true,
		`write transaction events to database (this option can significantly increase CPU and memory usage)`)
	f.BoolVar(&inputSensorParams.ShouldWritePeers, "write-peers", true, "write peers to database")
	f.BoolVar(&inputSensorParams.ShouldRunPprof, "pprof", false, "run pprof server")
	f.UintVar(&inputSensorParams.PprofPort, "pprof-port", 6060, "port pprof runs on")
	f.BoolVar(&inputSensorParams.ShouldRunPrometheus, "prom", true, "run Prometheus server")
	f.UintVar(&inputSensorParams.PrometheusPort, "prom-port", 2112, "port Prometheus runs on")
	f.UintVar(&inputSensorParams.APIPort, "api-port", 8080, "port API server will listen on")
	f.UintVar(&inputSensorParams.RPCPort, "rpc-port", 8545, "port for JSON-RPC server to receive transactions")
	f.StringVarP(&inputSensorParams.KeyFile, "key-file", "k", "", "private key file (cannot be set with --key)")
	f.StringVar(&inputSensorParams.PrivateKey, "key", "", "hex-encoded private key (cannot be set with --key-file)")
	SensorCmd.MarkFlagsMutuallyExclusive("key-file", "key")
	f.IntVar(&inputSensorParams.Port, "port", 30303, "TCP network listening port")
	f.IntVar(&inputSensorParams.DiscoveryPort, "discovery-port", 30303, "UDP P2P discovery port")
	f.StringVar(&inputSensorParams.RPC, "rpc", "https://polygon-rpc.com", "RPC endpoint used to fetch latest block")
	f.StringVar(&inputSensorParams.GenesisHash, "genesis-hash", "0xa9c28ce2141b56c474f1dc504bee9b01eb1bd7d1a507580d5519d4437a97de1b", "genesis block hash")
	f.BytesHexVar(&inputSensorParams.ForkID, "fork-id", []byte{240, 151, 188, 19}, "hex encoded fork ID (omit 0x)")
	f.IntVar(&inputSensorParams.DialRatio, "dial-ratio", 0,
		`ratio of inbound to dialed connections (dial ratio of 2 allows 1/2 of connections to be dialed, setting to 0 defaults to 3)`)
	f.StringVar(&inputSensorParams.NAT, "nat", "any", "NAT port mapping mechanism (any|none|upnp|pmp|pmp:<IP>|extip:<IP>)")
	f.StringVar(&inputSensorParams.StaticNodesFile, "static-nodes", "", "static nodes file")
	f.StringVar(&inputSensorParams.TrustedNodesFile, "trusted-nodes", "", "trusted nodes file")
	f.DurationVar(&inputSensorParams.TTL, "ttl", 14*24*time.Hour, "time to live")
	f.StringVar(&inputSensorParams.DiscoveryDNS, "discovery-dns", "", "DNS discovery ENR tree URL")
	f.StringVar(&inputSensorParams.Database, "database", "none",
		`which database to persist data to, options are:
  - datastore (GCP Datastore)
  - json (output to stdout)
  - none (no persistence)`)
	f.BoolVar(&inputSensorParams.NoDiscovery, "no-discovery", false, "disable P2P peer discovery")
	f.IntVar(&inputSensorParams.RequestsCache.MaxSize, "max-requests", 2048, "maximum request IDs to track per peer (0 for no limit)")
	f.DurationVar(&inputSensorParams.RequestsCache.TTL, "requests-cache-ttl", 5*time.Minute, "time to live for requests cache entries (0 for no expiration)")
	f.IntVar(&inputSensorParams.ParentsCache.MaxSize, "max-parents", 1024, "maximum parent block hashes to track per peer (0 for no limit)")
	f.DurationVar(&inputSensorParams.ParentsCache.TTL, "parents-cache-ttl", 5*time.Minute, "time to live for parent hash cache entries (0 for no expiration)")
	f.IntVar(&inputSensorParams.BlocksCache.MaxSize, "max-blocks", 1024, "maximum blocks to track across all peers (0 for no limit)")
	f.DurationVar(&inputSensorParams.BlocksCache.TTL, "blocks-cache-ttl", 10*time.Minute, "time to live for block cache entries (0 for no expiration)")
}
