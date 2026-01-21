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

// computeWave calculates the wave points based on the configuration
// The square wave alternates between high and low values over its period.
func (c *SquareWave) computeWave(config WaveConfig) {
	const start = float64(0)
	const step = float64(1.0)
	end := float64(config.Period)

	c.points = make(map[float64]float64)

	// Calculate the high and low values of the square wave.
	highValue := float64(config.Target) + float64(config.Amplitude)
	lowValue := float64(config.Target) - float64(config.Amplitude)

	// The duration of each state (high and low) is half the period.
	halfPeriod := float64(config.Period) / 2.0

	for x := start; x <= end; x += step {
		// math.Mod finds the remainder of the division of x by the period.
		timeInPeriod := math.Mod(x, float64(config.Period))

		// Determine if the wave is in the first half or second half of the period.
		if timeInPeriod < halfPeriod {
			c.points[x] = highValue
		} else {
			c.points[x] = lowValue
		}
	}
}
