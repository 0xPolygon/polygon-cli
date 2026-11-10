package p2p

import (
	"testing"

	ethp2p "github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

// newTestPeerMetrics builds PeerMetrics without registering it, so each test
// gets fresh gauges and can run any number of times.
func newTestPeerMetrics() *PeerMetrics {
	return newPeerMetrics(
		prometheus.NewGauge(prometheus.GaugeOpts{Name: "peers"}),
		prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "peers_by_direction"}, []string{"direction"}),
	)
}

func assertPeerMetrics(t *testing.T, m *PeerMetrics, total, inbound, outbound float64) {
	t.Helper()
	if got := testutil.ToFloat64(m.peers); got != total {
		t.Errorf("peers = %v, want %v", got, total)
	}
	if got := testutil.ToFloat64(m.peersByDirection.WithLabelValues("inbound")); got != inbound {
		t.Errorf("inbound peers = %v, want %v", got, inbound)
	}
	if got := testutil.ToFloat64(m.peersByDirection.WithLabelValues("outbound")); got != outbound {
		t.Errorf("outbound peers = %v, want %v", got, outbound)
	}
}

func TestPeerMetricsSplitsByDirection(t *testing.T) {
	m := newTestPeerMetrics()

	m.set(10, 3)
	assertPeerMetrics(t, m, 10, 3, 7)

	m.set(4, 4)
	assertPeerMetrics(t, m, 4, 4, 0)

	m.set(0, 0)
	assertPeerMetrics(t, m, 0, 0, 0)
}

func TestPeerMetricsUpdateCountsDialedPeersAsOutbound(t *testing.T) {
	m := newTestPeerMetrics()

	// NewPeer builds a peer without the inbound flag, as a dialed connection.
	peers := []*ethp2p.Peer{
		ethp2p.NewPeer(enode.ID{1}, "bor", nil),
		ethp2p.NewPeer(enode.ID{2}, "bor", nil),
	}
	m.Update(peers)
	assertPeerMetrics(t, m, 2, 0, 2)

	m.Update(nil)
	assertPeerMetrics(t, m, 0, 0, 0)
}

func TestPeerMetricsUpdateCountsInboundPeers(t *testing.T) {
	m := newTestPeerMetrics()

	peers := []*ethp2p.Peer{
		ethp2p.NewPeer(enode.ID{1}, "bor", nil),
		ethp2p.NewPeer(enode.ID{2}, "bor", nil),
		ethp2p.NewPeer(enode.ID{3}, "bor", nil),
	}
	inbound := map[enode.ID]bool{{1}: true, {3}: true}
	m.update(peers, func(p *ethp2p.Peer) bool { return inbound[p.ID()] })
	assertPeerMetrics(t, m, 3, 2, 1)
}

func TestNewPeerMetricsExportsBothDirectionsAtZero(t *testing.T) {
	m := newTestPeerMetrics()

	if got := testutil.CollectAndCount(m.peersByDirection); got != 2 {
		t.Fatalf("peers_by_direction series = %d, want 2 before the first update", got)
	}
	assertPeerMetrics(t, m, 0, 0, 0)
}
