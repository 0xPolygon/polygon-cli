package metrics

import (
	"sync"
	"time"

	"github.com/0xPolygon/polygon-cli/rpctypes"
	"github.com/rs/zerolog/log"
)

// MetricPlugin defines the interface for all metric calculators
type MetricPlugin interface {
	// Name returns the unique identifier for this metric
	Name() string

	// ProcessBlock processes a new block for metric calculation
	ProcessBlock(block rpctypes.PolyBlock)

	// GetMetric returns the current calculated metric value
	GetMetric() interface{}

	// GetUpdateInterval returns how often this metric should be recalculated
	GetUpdateInterval() time.Duration
}

// MetricUpdate represents an update from a metric plugin
type MetricUpdate struct {
	Name  string
	Value interface{}
	Time  time.Time
}

// MetricsSystem manages all metric plugins and coordinates updates
type MetricsSystem struct {
	mu       sync.RWMutex
	plugins  map[string]MetricPlugin
	updateCh chan MetricUpdate
	stopCh   chan struct{}
	wg       sync.WaitGroup
}

// NewMetricsSystem creates a new metrics system
func NewMetricsSystem() *MetricsSystem {
	return &MetricsSystem{
		plugins:  make(map[string]MetricPlugin),
		updateCh: make(chan MetricUpdate, 10000), // Increased buffer to handle bursts
		stopCh:   make(chan struct{}),
	}
}

// RegisterPlugin registers a new metric plugin
func (m *MetricsSystem) RegisterPlugin(plugin MetricPlugin) {
	m.mu.Lock()
	defer m.mu.Unlock()

	name := plugin.Name()
	if _, exists := m.plugins[name]; exists {
		log.Warn().Str("plugin", name).Msg("Metric plugin already registered, replacing")
	}

	m.plugins[name] = plugin
	log.Info().Str("plugin", name).Msg("Registered metric plugin")
}

// ProcessBlock sends a new block to all registered plugins
func (m *MetricsSystem) ProcessBlock(block rpctypes.PolyBlock) {
	m.mu.RLock()
	plugins := make([]MetricPlugin, 0, len(m.plugins))
	for _, plugin := range m.plugins {
		plugins = append(plugins, plugin)
	}
	m.mu.RUnlock()

	// Process block in each plugin
	for _, plugin := range plugins {
		m.wg.Add(1)
		go func(p MetricPlugin) {
			defer m.wg.Done()

			// Process the block
			p.ProcessBlock(block)

			// Send update
			select {
			case m.updateCh <- MetricUpdate{
				Name:  p.Name(),
				Value: p.GetMetric(),
				Time:  time.Now(),
			}:
			case <-m.stopCh:
				return
			}
		}(plugin)
	}
}

// GetUpdateChannel returns the channel for metric updates
func (m *MetricsSystem) GetUpdateChannel() <-chan MetricUpdate {
	return m.updateCh
}

// GetMetric returns the current value of a specific metric
func (m *MetricsSystem) GetMetric(name string) (interface{}, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if plugin, exists := m.plugins[name]; exists {
		return plugin.GetMetric(), true
	}
	return nil, false
}

// Stop gracefully shuts down the metrics system
func (m *MetricsSystem) Stop() {
	close(m.stopCh)
	m.wg.Wait()
	close(m.updateCh)
}
