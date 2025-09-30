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
		TransactionHashes:   getCounterValue(new(eth.NewPooledTransactionHashesPacket), url, counter),
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
