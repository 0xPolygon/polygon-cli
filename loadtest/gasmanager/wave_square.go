package gasmanager

import "math"

// SquareWave represents a square wave pattern for gas price modulation.
type SquareWave struct {
	*BaseWave
}

// NewSquareWave creates a new SquareWave with the given configuration.
func NewSquareWave(config WaveConfig) *SquareWave {
	c := &SquareWave{
		BaseWave: NewBaseWave(config),
	}

	c.computeWave(config)

	return c
}

// computeWave alternates between high and low values over the period.
func (c *SquareWave) computeWave(config WaveConfig) {
	period := float64(config.Period)
	highValue := float64(config.Target) + float64(config.Amplitude)
	lowValue := float64(config.Target) - float64(config.Amplitude)
	halfPeriod := period / 2.0

	for x := 0.0; x <= period; x++ {
		if math.Mod(x, period) < halfPeriod {
			c.points[x] = highValue
		} else {
			c.points[x] = lowValue
		}
	}
}
