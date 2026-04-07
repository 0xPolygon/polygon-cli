package p2p

import (
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// BlockMetrics contains Prometheus metrics for tracking head and oldest blocks.
type BlockMetrics struct {
	headBlockNumber    prometheus.Gauge
	headBlockTimestamp prometheus.Gauge
	headBlockAge       prometheus.Gauge
	oldestBlockNumber  prometheus.Gauge
	blockRange         prometheus.Gauge
}

// NewBlockMetrics creates and registers all block-related Prometheus metrics.
func NewBlockMetrics(block *types.Block) *BlockMetrics {
	m := &BlockMetrics{
		headBlockNumber: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: "sensor",
			Name:      "head_block_number",
			Help:      "Current head block number",
		}),
		headBlockTimestamp: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: "sensor",
			Name:      "head_block_timestamp",
			Help:      "Head block timestamp in Unix epoch seconds",
		}),
		headBlockAge: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: "sensor",
			Name:      "head_block_age",
			Help:      "Time since head block was received (in seconds)",
		}),
		oldestBlockNumber: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: "sensor",
			Name:      "oldest_block_number",
			Help:      "Oldest block number (floor for parent fetching)",
		}),
		blockRange: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: "sensor",
			Name:      "block_range",
			Help:      "Difference between head and oldest block numbers",
		}),
	}

	m.headBlockNumber.Set(float64(block.NumberU64()))
	m.headBlockTimestamp.Set(float64(block.Time()))
	m.headBlockAge.Set(0)
	m.oldestBlockNumber.Set(float64(block.NumberU64()))
	m.blockRange.Set(0)

	return m
}

// Update updates all block metrics.
func (m *BlockMetrics) Update(block *types.Block, oldest *types.Header) {
	if block == nil || oldest == nil {
		return
	}

	hn := block.NumberU64()
	on := oldest.Number.Uint64()
	ht := time.Unix(int64(block.Time()), 0)

	m.headBlockNumber.Set(float64(hn))
	m.headBlockTimestamp.Set(float64(block.Time()))
	m.headBlockAge.Set(time.Since(ht).Seconds())
	m.oldestBlockNumber.Set(float64(on))
	m.blockRange.Set(float64(hn - on))
}

// broadcastMetrics contains Prometheus metrics for tracking transaction broadcasts.
type broadcastMetrics struct {
	queueDepth prometheus.Gauge
	batchSize  prometheus.Histogram
	sendErrors prometheus.Counter
}

// newBroadcastMetrics creates and registers all broadcast-related Prometheus metrics.
func newBroadcastMetrics() *broadcastMetrics {
	return &broadcastMetrics{
		queueDepth: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: "sensor",
			Name:      "broadcast_queue_depth",
			Help:      "Number of transaction batches in broadcast queue",
		}),
		batchSize: promauto.NewHistogram(prometheus.HistogramOpts{
			Namespace: "sensor",
			Name:      "broadcast_batch_size",
			Help:      "Number of transactions per broadcast batch",
			Buckets:   []float64{10, 50, 100, 500, 1000, 5000, 10000},
		}),
		sendErrors: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: "sensor",
			Name:      "broadcast_send_errors",
			Help:      "Number of failed broadcast sends",
		}),
	}
}
