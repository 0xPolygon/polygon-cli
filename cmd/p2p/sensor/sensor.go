package sensor

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"
	"time"

	"net/http"
	_ "net/http/pprof"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/forkid"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	ethp2p "github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/dnsdisc"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/nat"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	dto "github.com/prometheus/client_model/go"
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
		QuickStart                   bool
		TTL                          time.Duration
		DiscoveryDNS                 string
		Persistence                  string

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
			go handlePprof()
		}

		if inputSensorParams.ShouldRunPrometheus {
			go handlePrometheus()
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
		// Initialize the database based on the persistence flag
		var db database.Database
		switch inputSensorParams.Persistence {
		case "datastore":
			db = database.NewDatastore(cmd.Context(), database.DatastoreOptions{
				ProjectID:                    inputSensorParams.ProjectID,
				DatabaseID:                   inputSensorParams.DatabaseID,
				SensorID:                     inputSensorParams.SensorID,
				MaxConcurrency:               inputSensorParams.MaxDatabaseConcurrency,
				ShouldWriteBlocks:            inputSensorParams.ShouldWriteBlocks,
				ShouldWriteBlockEvents:       inputSensorParams.ShouldWriteBlockEvents,
				ShouldWriteTransactions:      inputSensorParams.ShouldWriteTransactions,
				ShouldWriteTransactionEvents: inputSensorParams.ShouldWriteTransactionEvents,
				ShouldWritePeers:             inputSensorParams.ShouldWritePeers,
				TTL:                          inputSensorParams.TTL,
			})
		case "json":
			db = database.NewJSONDatabase(database.JSONDatabaseOptions{
				SensorID:                     inputSensorParams.SensorID,
				MaxConcurrency:               inputSensorParams.MaxDatabaseConcurrency,
				ShouldWriteBlocks:            inputSensorParams.ShouldWriteBlocks,
				ShouldWriteBlockEvents:       inputSensorParams.ShouldWriteBlockEvents,
				ShouldWriteTransactions:      inputSensorParams.ShouldWriteTransactions,
				ShouldWriteTransactionEvents: inputSensorParams.ShouldWriteTransactionEvents,
				ShouldWritePeers:             inputSensorParams.ShouldWritePeers,
			})
		case "false":
			db = database.NewNoopDatabase()
		default:
			return fmt.Errorf("invalid persistence option: %s", inputSensorParams.Persistence)
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
			Help:      "The number and type of messages the sensor has received",
		}, []string{"message", "url", "name"})

		// Create connection manager for transaction broadcasting
		connManager := p2p.NewConnectionManager()

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
			ConnManager: connManager,
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
		if err = server.Start(); err != nil {
			return err
		}
		defer server.Stop()

		events := make(chan *ethp2p.PeerEvent)
		sub := server.SubscribeEvents(events)
		defer sub.Unsubscribe()

		ticker := time.NewTicker(2 * time.Second) // Ticker for recurring tasks every 2 seconds.
		hourlyTicker := time.NewTicker(time.Hour) // Ticker for running DNS discovery every hour.
		defer ticker.Stop()
		defer hourlyTicker.Stop()

		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

		// peers represents the peer map that is used to write to the nodes.json
		// file. This is helpful when restarting the node with the --quickstart flag
		// enabled. This map does not represent the peers that are currently
		// connected to the sensor. To do that use `server.Peers()` instead.
		peers := make(map[enode.ID]string)
		var peersMutex sync.Mutex

		for _, node := range inputSensorParams.nodes {
			// Map node URLs to node IDs to avoid duplicates.
			peers[node.ID()] = node.URLv4()
		}

		go handleAPI(&server, msgCounter)
		go handleRPC(connManager, inputSensorParams.NetworkID)

		// Run DNS discovery immediately at startup.
		go handleDNSDiscovery(&server)

		for {
			select {
			case <-ticker.C:
				peersGauge.Set(float64(server.PeerCount()))
				if err := removePeerMessages(msgCounter, server.Peers()); err != nil {
					log.Error().Err(err).Msg("Failed to clean up peer messages")
				}
				db.WritePeers(context.Background(), server.Peers(), time.Now())
			case peer := <-opts.Peers:
				// Lock the peers map before modifying it.
				peersMutex.Lock()
				// Update the peer list and the nodes file.
				if _, ok := peers[peer.ID()]; !ok {
					peers[peer.ID()] = peer.URLv4()

					if err := p2p.WritePeers(inputSensorParams.NodesFile, peers); err != nil {
						log.Error().Err(err).Msg("Failed to write nodes to file")
					}
				}
				peersMutex.Unlock()
			case <-hourlyTicker.C:
				go handleDNSDiscovery(&server)
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

// handleAPI sets up the API for interacting with the sensor. The `/peers`
// endpoint returns a list of all peers connected to the sensor, including the
// types and counts of eth packets sent by each peer.
func handleAPI(server *ethp2p.Server, counter *prometheus.CounterVec) {
	http.HandleFunc("/peers", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		peers := make(map[string]p2p.MessageCount)
		for _, peer := range server.Peers() {
			url := peer.Node().URLv4()
			peers[url] = getPeerMessages(url, counter)
		}

		err := json.NewEncoder(w).Encode(peers)
		if err != nil {
			log.Error().Err(err).Msg("Failed to encode peers")
		}
	})

	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		type NodeInfo struct {
			ENR string `json:"enr"`
			URL string `json:"enode"`
		}

		info := NodeInfo{
			ENR: server.NodeInfo().ENR,
			URL: server.Self().URLv4(),
		}

		err := json.NewEncoder(w).Encode(info)
		if err != nil {
			log.Error().Err(err).Msg("Failed to encode node info")
		}
	})

	addr := fmt.Sprintf(":%d", inputSensorParams.APIPort)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Error().Err(err).Msg("Failed to start API handler")
	}
}

// handleDNSDiscovery performs DNS-based peer discovery and adds new peers to
// the p2p server. It syncs the DNS discovery tree and adds any newly discovered
// peers not already in the peers map.
func handleDNSDiscovery(server *ethp2p.Server) {
	if len(inputSensorParams.DiscoveryDNS) == 0 {
		return
	}

	log.Info().
		Str("discovery-dns", inputSensorParams.DiscoveryDNS).
		Msg("Starting DNS discovery sync")

	client := dnsdisc.NewClient(dnsdisc.Config{})
	tree, err := client.SyncTree(inputSensorParams.DiscoveryDNS)
	if err != nil {
		log.Error().Err(err).Msg("Failed to sync DNS discovery tree")
		return
	}

	// Log the number of nodes in the tree.
	log.Info().
		Int("unique_nodes", len(tree.Nodes())).
		Msg("Successfully synced DNS discovery tree")

	// Add DNS-discovered peers.
	for _, node := range tree.Nodes() {
		log.Debug().
			Str("enode", node.URLv4()).
			Msg("Discovered peer through DNS")

		// Add the peer to the static node set. The server itself handles whether to
		// connect to the peer if it's already connected. If a node is part of the
		// static peer set, the server will handle reconnecting after disconnects.
		server.AddPeer(node)
	}

	log.Info().Msg("Finished adding DNS discovery peers")
}

// handleRPC sets up the JSON-RPC server for receiving and broadcasting transactions.
// It handles eth_sendRawTransaction requests, validates transaction signatures,
// and broadcasts valid transactions to all connected peers.
func handleRPC(connManager *p2p.ConnectionManager, networkID uint64) {
	// Use network ID as chain ID for signature validation
	chainID := new(big.Int).SetUint64(networkID)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Read request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeJSONError(w, -32700, "Parse error", nil)
			return
		}
		defer r.Body.Close()

		// Parse JSON-RPC request
		var req jsonRPCRequest
		if err := json.Unmarshal(body, &req); err != nil {
			writeJSONError(w, -32700, "Parse error", nil)
			return
		}

		// Handle eth_sendRawTransaction
		if req.Method == "eth_sendRawTransaction" {
			handleSendRawTransaction(w, req, connManager, chainID)
			return
		}

		// Method not found
		writeJSONError(w, -32601, "Method not found", req.ID)
	})

	addr := fmt.Sprintf(":%d", inputSensorParams.RPCPort)
	log.Info().Str("addr", addr).Msg("Starting JSON-RPC server")
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Error().Err(err).Msg("Failed to start RPC server")
	}
}

// JSON-RPC request structure
type jsonRPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      interface{}   `json:"id"`
}

// JSON-RPC response structures
type jsonRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *rpcError   `json:"error,omitempty"`
	ID      interface{} `json:"id"`
}

type rpcError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// writeJSONError writes a JSON-RPC error response
func writeJSONError(w http.ResponseWriter, code int, message string, id interface{}) {
	w.Header().Set("Content-Type", "application/json")
	response := jsonRPCResponse{
		JSONRPC: "2.0",
		Error: &rpcError{
			Code:    code,
			Message: message,
		},
		ID: id,
	}
	json.NewEncoder(w).Encode(response)
}

// writeJSONResult writes a JSON-RPC success response
func writeJSONResult(w http.ResponseWriter, result interface{}, id interface{}) {
	w.Header().Set("Content-Type", "application/json")
	response := jsonRPCResponse{
		JSONRPC: "2.0",
		Result:  result,
		ID:      id,
	}
	json.NewEncoder(w).Encode(response)
}

// handleSendRawTransaction processes eth_sendRawTransaction requests
func handleSendRawTransaction(w http.ResponseWriter, req jsonRPCRequest, connManager *p2p.ConnectionManager, chainID *big.Int) {
	// Check params
	if len(req.Params) == 0 {
		writeJSONError(w, -32602, "Invalid params: missing raw transaction", req.ID)
		return
	}

	// Extract raw transaction hex string
	rawTxHex, ok := req.Params[0].(string)
	if !ok {
		writeJSONError(w, -32602, "Invalid params: raw transaction must be a hex string", req.ID)
		return
	}

	// Decode hex string to bytes
	txBytes, err := hexutil.Decode(rawTxHex)
	if err != nil {
		writeJSONError(w, -32602, fmt.Sprintf("Invalid transaction hex: %v", err), req.ID)
		return
	}

	// Unmarshal transaction
	tx := new(types.Transaction)
	if err := tx.UnmarshalBinary(txBytes); err != nil {
		writeJSONError(w, -32602, fmt.Sprintf("Invalid transaction encoding: %v", err), req.ID)
		return
	}

	// Validate transaction signature
	signer := types.LatestSignerForChainID(chainID)
	sender, err := types.Sender(signer, tx)
	if err != nil {
		writeJSONError(w, -32602, fmt.Sprintf("Invalid transaction signature: %v", err), req.ID)
		return
	}

	// Log the transaction
	toAddr := "nil"
	if tx.To() != nil {
		toAddr = tx.To().Hex()
	}

	log.Info().
		Str("hash", tx.Hash().Hex()).
		Str("from", sender.Hex()).
		Str("to", toAddr).
		Str("value", tx.Value().String()).
		Uint64("gas", tx.Gas()).
		Msg("Broadcasting transaction")

	// Broadcast to all peers using the connection manager
	broadcastCount := connManager.BroadcastTransaction(tx)

	log.Info().
		Str("hash", tx.Hash().Hex()).
		Int("peers", broadcastCount).
		Msg("Transaction broadcast complete")

	// Return transaction hash
	writeJSONResult(w, tx.Hash().Hex(), req.ID)
}

// getPeerMessages retrieves the count of various types of eth packets sent by a
// peer.
func getPeerMessages(url string, counter *prometheus.CounterVec) p2p.MessageCount {
	return p2p.MessageCount{
		BlockHeaders:        getCounterValue(new(eth.BlockHeadersPacket), url, counter),
		BlockBodies:         getCounterValue(new(eth.BlockBodiesPacket), url, counter),
		Blocks:              getCounterValue(new(eth.NewBlockPacket), url, counter),
		BlockHashes:         getCounterValue(new(eth.NewBlockHashesPacket), url, counter),
		BlockHeaderRequests: getCounterValue(new(eth.GetBlockHeadersPacket), url, counter),
		BlockBodiesRequests: getCounterValue(new(eth.GetBlockBodiesPacket), url, counter),
		Transactions: getCounterValue(new(eth.TransactionsPacket), url, counter) +
			getCounterValue(new(eth.PooledTransactionsPacket), url, counter),
		TransactionHashes: getCounterValue(new(eth.NewPooledTransactionHashesPacket), url, counter) +
			getCounterValue(new(eth.NewPooledTransactionHashesPacket), url, counter),
		TransactionRequests: getCounterValue(new(eth.GetPooledTransactionsRequest), url, counter),
	}
}

// getCounterValue retrieves the count of packets for a specific type from the
// Prometheus counter.
func getCounterValue(packet eth.Packet, url string, counter *prometheus.CounterVec) int64 {
	metric := &dto.Metric{}

	err := counter.WithLabelValues(packet.Name(), url).Write(metric)
	if err != nil {
		log.Error().Err(err).Send()
		return 0
	}

	return int64(metric.GetCounter().GetValue())
}

// removePeerMessages removes all the counters of peers that disconnected from
// the sensor. This prevents the metrics list from infinitely growing.
func removePeerMessages(counter *prometheus.CounterVec, peers []*ethp2p.Peer) error {
	urls := []string{}
	for _, peer := range peers {
		urls = append(urls, peer.Node().URLv4())
	}

	families, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		return err
	}

	var family *dto.MetricFamily
	for _, f := range families {
		if f.GetName() == "sensor_messages" {
			family = f
			break
		}
	}

	// During DNS-discovery or when the server is taking a while to discover
	// peers and has yet to receive a message, the sensor_messages prometheus
	// metric may not exist yet.
	if family == nil {
		log.Trace().Msg("Could not find sensor_messages metric family")
		return nil
	}

	for _, metric := range family.GetMetric() {
		for _, label := range metric.GetLabel() {
			url := label.GetValue()
			if label.GetName() != "url" || slices.Contains(urls, url) {
				continue
			}

			counter.DeletePartialMatch(prometheus.Labels{"url": url})
		}
	}

	return nil
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
	SensorCmd.Flags().IntVarP(&inputSensorParams.MaxPeers, "max-peers", "m", 2000, "Maximum number of peers to connect to")
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
	SensorCmd.Flags().BoolVar(&inputSensorParams.ShouldWritePeers, "write-peers", true, "Whether to write peers to the database")
	SensorCmd.Flags().BoolVar(&inputSensorParams.ShouldRunPprof, "pprof", false, "Whether to run pprof")
	SensorCmd.Flags().UintVar(&inputSensorParams.PprofPort, "pprof-port", 6060, "Port pprof runs on")
	SensorCmd.Flags().BoolVar(&inputSensorParams.ShouldRunPrometheus, "prom", true, "Whether to run Prometheus")
	SensorCmd.Flags().UintVar(&inputSensorParams.PrometheusPort, "prom-port", 2112, "Port Prometheus runs on")
	SensorCmd.Flags().UintVar(&inputSensorParams.APIPort, "api-port", 8080, "Port the API server will listen on")
	SensorCmd.Flags().UintVar(&inputSensorParams.RPCPort, "rpc-port", 8545, "Port for JSON-RPC server to receive transactions")
	SensorCmd.Flags().StringVarP(&inputSensorParams.KeyFile, "key-file", "k", "", "Private key file (cannot be set with --key)")
	SensorCmd.Flags().StringVar(&inputSensorParams.PrivateKey, "key", "", "Hex-encoded private key (cannot be set with --key-file)")
	SensorCmd.MarkFlagsMutuallyExclusive("key-file", "key")
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
	SensorCmd.Flags().StringVar(&inputSensorParams.DiscoveryDNS, "discovery-dns", "", "DNS discovery ENR tree url")
	SensorCmd.Flags().StringVar(&inputSensorParams.Persistence, "persistence", "datastore", "Persistence mode: datastore (Google Datastore), json (output to stdout), false (no persistence)")
}
