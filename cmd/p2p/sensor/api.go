package sensor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/0xPolygon/polygon-cli/p2p"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	ethp2p "github.com/ethereum/go-ethereum/p2p"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/rs/zerolog/log"
)

// peerData represents the metrics and connection information for a peer.
// It includes both message counts (items sent/received) and packet counts
// (number of p2p messages), along with connection timing information.
type peerData struct {
	Received        p2p.MessageCount `json:"received"`
	Sent            p2p.MessageCount `json:"sent"`
	PacketsReceived p2p.MessageCount `json:"packets_received"`
	PacketsSent     p2p.MessageCount `json:"packets_sent"`
	ConnectedAt     string           `json:"connected_at"`
	DurationSeconds float64          `json:"duration_seconds"`
}

// blockInfo represents basic block information.
type blockInfo struct {
	Hash   string `json:"hash"`
	Number uint64 `json:"number"`
}

// newBlockInfo creates a blockInfo from a types.Header.
// Returns nil if the header is nil.
func newBlockInfo(header *types.Header) *blockInfo {
	if header == nil {
		return nil
	}

	return &blockInfo{
		Hash:   header.Hash().Hex(),
		Number: header.Number.Uint64(),
	}
}

// apiData represents all sensor information including node info and peer data.
type apiData struct {
	ENR       string              `json:"enr"`
	URL       string              `json:"enode"`
	PeerCount int                 `json:"peer_count"`
	Peers     map[string]peerData `json:"peers"`
	Head      *blockInfo          `json:"head_block"`
	Oldest    *blockInfo          `json:"oldest_block"`
}

// handleAPI sets up the API for interacting with the sensor. All endpoints
// return information about the sensor node and all connected peers, including
// the types and counts of eth packets sent and received by each peer.
func handleAPI(server *ethp2p.Server, counter *prometheus.CounterVec, conns *p2p.Conns) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		peers := make(map[string]peerData)
		for _, peer := range server.Peers() {
			url := peer.Node().URLv4()
			peerID := peer.Node().ID().String()
			name := peer.Fullname()
			connectedAt := conns.PeerConnectedAt(peerID)
			if connectedAt.IsZero() {
				continue
			}

			msgs := peerData{
				Received:        getPeerMessages(counter, url, name, p2p.MsgReceived, false),
				Sent:            getPeerMessages(counter, url, name, p2p.MsgSent, false),
				PacketsReceived: getPeerMessages(counter, url, name, p2p.MsgReceived, true),
				PacketsSent:     getPeerMessages(counter, url, name, p2p.MsgSent, true),
				ConnectedAt:     connectedAt.UTC().Format(time.RFC3339),
				DurationSeconds: time.Since(connectedAt).Seconds(),
			}

			peers[url] = msgs
		}

		head := conns.HeadBlock()
		oldest := conns.OldestBlock()

		var headHeader *types.Header
		if head.Block != nil {
			headHeader = head.Block.Header()
		}

		data := apiData{
			ENR:       server.NodeInfo().ENR,
			URL:       server.Self().URLv4(),
			PeerCount: len(peers),
			Peers:     peers,
			Head:      newBlockInfo(headHeader),
			Oldest:    newBlockInfo(oldest),
		}

		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Error().Err(err).Msg("Failed to encode sensor data")
		}
	})

	addr := fmt.Sprintf(":%d", inputSensorParams.APIPort)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Error().Err(err).Msg("Failed to start API handler")
	}
}

// getPeerMessages retrieves the count of various types of eth packets sent by a
// peer.
func getPeerMessages(counter *prometheus.CounterVec, url, name string, direction p2p.Direction, isPacket bool) p2p.MessageCount {
	return p2p.MessageCount{
		BlockHeaders:        getCounterValue(new(eth.BlockHeadersPacket), counter, url, name, direction, isPacket),
		BlockBodies:         getCounterValue(new(eth.BlockBodiesPacket), counter, url, name, direction, isPacket),
		Blocks:              getCounterValue(new(eth.NewBlockPacket), counter, url, name, direction, isPacket),
		BlockHashes:         getCounterValue(new(eth.NewBlockHashesPacket), counter, url, name, direction, isPacket),
		BlockHeaderRequests: getCounterValue(new(eth.GetBlockHeadersPacket), counter, url, name, direction, isPacket),
		BlockBodiesRequests: getCounterValue(new(eth.GetBlockBodiesPacket), counter, url, name, direction, isPacket),
		Transactions: getCounterValue(new(eth.TransactionsPacket), counter, url, name, direction, isPacket) +
			getCounterValue(new(eth.PooledTransactionsPacket), counter, url, name, direction, isPacket),
		TransactionHashes:   getCounterValue(new(eth.NewPooledTransactionHashesPacket), counter, url, name, direction, isPacket),
		TransactionRequests: getCounterValue(new(eth.GetPooledTransactionsRequest), counter, url, name, direction, isPacket),
	}
}

// getCounterValue retrieves the count of packets for a specific type from the
// Prometheus counter.
func getCounterValue(packet eth.Packet, counter *prometheus.CounterVec, url, name string, direction p2p.Direction, isPacket bool) int64 {
	metric := &dto.Metric{}

	messageName := packet.Name()
	if isPacket {
		messageName += p2p.PacketSuffix
	}

	err := counter.WithLabelValues(messageName, url, name, string(direction)).Write(metric)
	if err != nil {
		log.Error().Err(err).Send()
		return 0
	}

	return int64(metric.GetCounter().GetValue())
}

// removePeerMessages removes all the counters of peers that disconnected from
// the sensor. This prevents the metrics list from infinitely growing.
func removePeerMessages(counter *prometheus.CounterVec, urls []string) error {
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
