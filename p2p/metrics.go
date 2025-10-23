package p2p

import (
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// BlockMetrics contains Prometheus metrics for tracking head and oldest blocks.
type BlockMetrics struct {
	HeadBlockNumber    prometheus.Gauge
	HeadBlockTimestamp prometheus.Gauge
	HeadBlockAge       prometheus.Gauge
	OldestBlockNumber  prometheus.Gauge
	BlockRange         prometheus.Gauge
}

// NewBlockMetrics creates and registers all block-related Prometheus metrics.
func NewBlockMetrics(block *types.Block) *BlockMetrics {
	m := &BlockMetrics{
		HeadBlockNumber: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: "sensor",
			Name:      "head_block_number",
			Help:      "Current head block number",
		}),
		HeadBlockTimestamp: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: "sensor",
			Name:      "head_block_timestamp",
			Help:      "Head block timestamp in Unix epoch seconds",
		}),
		HeadBlockAge: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: "sensor",
			Name:      "head_block_age",
			Help:      "Time since head block was received (in seconds)",
		}),
		OldestBlockNumber: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: "sensor",
			Name:      "oldest_block_number",
			Help:      "Oldest block number (floor for parent fetching)",
		}),
		BlockRange: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: "sensor",
			Name:      "block_range",
			Help:      "Difference between head and oldest block numbers",
		}),
	}

	m.HeadBlockNumber.Set(float64(block.NumberU64()))
	m.HeadBlockTimestamp.Set(float64(block.Time()))
	m.HeadBlockAge.Set(0)
	m.OldestBlockNumber.Set(float64(block.NumberU64()))
	m.BlockRange.Set(0)

	return m
}

// Update updates all block metrics.
func (m *BlockMetrics) Update(block *types.Block, oldest *types.Header) {
	if m == nil {
		return
	}

	hn := block.NumberU64()
	on := oldest.Number.Uint64()
	ht := time.Unix(int64(block.Time()), 0)

	m.HeadBlockNumber.Set(float64(hn))
	m.HeadBlockTimestamp.Set(float64(block.Time()))
	m.HeadBlockAge.Set(time.Since(ht).Seconds())
	m.OldestBlockNumber.Set(float64(on))
	m.BlockRange.Set(float64(hn - on))
}
