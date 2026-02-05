package sensor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/0xPolygon/polygon-cli/p2p"
	"github.com/ethereum/go-ethereum/core/types"
	ethp2p "github.com/ethereum/go-ethereum/p2p"
	"github.com/rs/zerolog/log"
)

// peerData represents the metrics and connection information for a peer.
// It includes both message counts (items sent/received) and packet counts
// (number of p2p messages), along with connection timing information.
type peerData struct {
	Name            string           `json:"name"`
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
func handleAPI(server *ethp2p.Server, conns *p2p.Conns) {
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
			connectedAt := conns.PeerConnectedAt(peerID)
			if connectedAt.IsZero() {
				continue
			}

			// Get per-peer message counts from in-memory tracking
			messages := conns.GetPeerMessages(peerID)
			if messages == nil {
				continue
			}

			peers[url] = peerData{
				Name:            conns.GetPeerName(peerID),
				Received:        messages.Received,
				Sent:            messages.Sent,
				PacketsReceived: messages.PacketsReceived,
				PacketsSent:     messages.PacketsSent,
				ConnectedAt:     connectedAt.UTC().Format(time.RFC3339),
				DurationSeconds: time.Since(connectedAt).Seconds(),
			}
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

