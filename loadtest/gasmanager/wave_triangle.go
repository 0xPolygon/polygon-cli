package gasmanager

import "math"

// TriangleWave represents a triangle wave pattern for gas price modulation.
type TriangleWave struct {
	*BaseWave
}

// NewTriangleWave creates a new TriangleWave with the given configuration.
func NewTriangleWave(config WaveConfig) *TriangleWave {
	c := &TriangleWave{
		BaseWave: NewBaseWave(config),
	}

	c.computeWave(config)

	return c
}

// computeWave creates a triangle wave that rises and falls linearly.
func (c *TriangleWave) computeWave(config WaveConfig) {
	period := float64(config.Period)
	amplitude := float64(config.Amplitude)
	target := float64(config.Target)
	peakToPeak := 2.0 * amplitude

	for x := 0.0; x <= period; x++ {
		normalizedTime := math.Mod(x, period) / period
		// abs(2*t - 1) creates a triangle oscillating between 0 and 1
		c.points[x] = target + amplitude - peakToPeak*math.Abs(2*normalizedTime-1)
	}
}
