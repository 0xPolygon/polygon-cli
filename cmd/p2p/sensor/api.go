package sensor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/0xPolygon/polygon-cli/p2p"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	ethp2p "github.com/ethereum/go-ethereum/p2p"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/rs/zerolog/log"
)

// nodeInfo represents information about the sensor node.
type nodeInfo struct {
	ENR string `json:"enr"`
	URL string `json:"enode"`
}

// peerInfo represents information about a connected peer.
type peerInfo struct {
	MessagesReceived p2p.MessageCount `json:"messages_received"`
	MessagesSent     p2p.MessageCount `json:"messages_sent"`
	ConnectedAt      string           `json:"connected_at"`
	DurationSeconds  int64            `json:"duration_seconds"`
}

// handleAPI sets up the API for interacting with the sensor. The `/peers`
// endpoint returns a list of all peers connected to the sensor, including the
// types and counts of eth packets sent by and received from each peer.
func handleAPI(server *ethp2p.Server, msgsReceived, msgsSent *prometheus.CounterVec, conns *p2p.Conns) {
	http.HandleFunc("/peers", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		peers := make(map[string]peerInfo)
		for _, peer := range server.Peers() {
			url := peer.Node().URLv4()
			nodeID := peer.Node().ID().String()
			connectedAt := conns.GetPeerConnectedAt(nodeID)

			peers[url] = peerInfo{
				MessagesReceived: getPeerMessages(url, peer.Fullname(), msgsReceived),
				MessagesSent:     getPeerMessages(url, peer.Fullname(), msgsSent),
				ConnectedAt:      connectedAt.UTC().Format(time.RFC3339),
				DurationSeconds:  int64(time.Since(connectedAt).Seconds()),
			}
		}

		if err := json.NewEncoder(w).Encode(peers); err != nil {
			log.Error().Err(err).Msg("Failed to encode peers")
		}
	})

	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		info := nodeInfo{
			ENR: server.NodeInfo().ENR,
			URL: server.Self().URLv4(),
		}

		if err := json.NewEncoder(w).Encode(info); err != nil {
			log.Error().Err(err).Msg("Failed to encode node info")
		}
	})

	addr := fmt.Sprintf(":%d", inputSensorParams.APIPort)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Error().Err(err).Msg("Failed to start API handler")
	}
}

// getPeerMessages retrieves the count of various types of eth packets sent by a
// peer.
func getPeerMessages(url, name string, counter *prometheus.CounterVec) p2p.MessageCount {
	return p2p.MessageCount{
		BlockHeaders:        getCounterValue(new(eth.BlockHeadersPacket), url, name, counter),
		BlockBodies:         getCounterValue(new(eth.BlockBodiesPacket), url, name, counter),
		Blocks:              getCounterValue(new(eth.NewBlockPacket), url, name, counter),
		BlockHashes:         getCounterValue(new(eth.NewBlockHashesPacket), url, name, counter),
		BlockHeaderRequests: getCounterValue(new(eth.GetBlockHeadersPacket), url, name, counter),
		BlockBodiesRequests: getCounterValue(new(eth.GetBlockBodiesPacket), url, name, counter),
		Transactions: getCounterValue(new(eth.TransactionsPacket), url, name, counter) +
			getCounterValue(new(eth.PooledTransactionsPacket), url, name, counter),
		TransactionHashes:   getCounterValue(new(eth.NewPooledTransactionHashesPacket), url, name, counter),
		TransactionRequests: getCounterValue(new(eth.GetPooledTransactionsRequest), url, name, counter),
	}
}

// getCounterValue retrieves the count of packets for a specific type from the
// Prometheus counter.
func getCounterValue(packet eth.Packet, url, name string, counter *prometheus.CounterVec) int64 {
	metric := &dto.Metric{}

	err := counter.WithLabelValues(packet.Name(), url, name).Write(metric)
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

	// Find all matching metric families
	for _, family := range families {
		// Check for any sensor_messages metric (received, sent, etc.)
		if !strings.Contains(family.GetName(), "sensor_messages") {
			continue
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
	}

	return nil
}
