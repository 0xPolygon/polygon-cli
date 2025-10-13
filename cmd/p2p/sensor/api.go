package sensor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

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

// peerMessages represents the messages sent and received by a peer.
type peerMessages struct {
	Received        p2p.MessageCount `json:"received"`
	Sent            p2p.MessageCount `json:"sent"`
	PacketsReceived p2p.MessageCount `json:"packets_received"`
	PacketsSent     p2p.MessageCount `json:"packets_sent"`
}

// handleAPI sets up the API for interacting with the sensor. The `/peers`
// endpoint returns a list of all peers connected to the sensor, including the
// types and counts of eth packets sent and received by each peer.
func handleAPI(server *ethp2p.Server, counter *prometheus.CounterVec) {
	http.HandleFunc("/peers", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		peers := make(map[string]peerMessages)
		for _, peer := range server.Peers() {
			url := peer.Node().URLv4()
			peers[url] = peerMessages{
				Received:        getPeerMessages(counter, url, peer.Fullname(), p2p.MsgReceived, false),
				Sent:            getPeerMessages(counter, url, peer.Fullname(), p2p.MsgSent, false),
				PacketsReceived: getPeerMessages(counter, url, peer.Fullname(), p2p.MsgReceived, true),
				PacketsSent:     getPeerMessages(counter, url, peer.Fullname(), p2p.MsgSent, true),
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
