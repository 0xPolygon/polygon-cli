package gasmanager

import "math"

// SineWave represents a sine wave pattern for gas price modulation.
type SineWave struct {
	*BaseWave
}

// NewSineWave creates a new SineWave with the given configuration.
func NewSineWave(config WaveConfig) *SineWave {
	c := &SineWave{
		BaseWave: NewBaseWave(config),
	}

	c.computeWave(config)

	return c
}

// computeWave calculates the wave points using: y = A * sin(2π/P * x) + T
// where A = Amplitude, P = Period, T = Target (vertical shift)
func (c *SineWave) computeWave(config WaveConfig) {
	period := float64(config.Period)
	amplitude := float64(config.Amplitude)
	target := float64(config.Target)
	b := (2 * math.Pi) / period

	for x := 0.0; x <= period; x++ {
		c.points[x] = amplitude*math.Sin(b*x) + target
	}
}
