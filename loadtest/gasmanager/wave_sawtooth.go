package gasmanager

import "math"

// SawtoothWave represents a wave that linearly rises from the minimum to the maximum value over its period.
type SawtoothWave struct {
	*BaseWave
}

// NewSawtoothWave creates a new SawtoothWave with the given configuration.
func NewSawtoothWave(config WaveConfig) *SawtoothWave {
	c := &SawtoothWave{
		BaseWave: NewBaseWave(config),
	}

	c.computeWave(config)

	return c
}

// computeWave creates a sawtooth that rises linearly from min to max over the period.
func (c *SawtoothWave) computeWave(config WaveConfig) {
	period := float64(config.Period)
	offset := float64(config.Target - config.Amplitude)
	rangeOfWave := float64(2 * config.Amplitude)

	for x := 0.0; x <= period; x++ {
		fractionalTime := math.Mod(x, period) / period
		c.points[x] = rangeOfWave*fractionalTime + offset
	}
}
