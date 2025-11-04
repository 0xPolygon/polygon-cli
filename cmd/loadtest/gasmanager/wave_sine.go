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

// computeWave calculates the wave points based on the configuration
// using the generalized sine wave formula.
// The formula used is: y = A * sin(b(x)) + k
// where:
// A = Amplitude
// b = (2π) / Period
// k = Target (vertical shift)
func (c *SineWave) computeWave(config WaveConfig) {
	const start = float64(0)
	const step = float64(1.0)
	end := float64(config.Period)

	c.points = map[float64]float64{}

	// The `b` parameter in the generalized sine formula is derived from the period.
	b := (2 * math.Pi) / float64(config.Period)

	for x := start; x <= end; x += step {
		// Apply the generalized sine formula: y = A * sin(b(x)) + k
		y := float64(config.Amplitude)*math.Sin(b*x) + float64(config.Target)
		c.points[float64(x)] = float64(y)
	}
}
