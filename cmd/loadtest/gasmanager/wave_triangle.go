package gasmanager

import "math"

type TriangleWave struct {
	BaseWave
	points map[float64]float64
}

func NewTriangleWave(config WaveConfig) *TriangleWave {
	c := &TriangleWave{
		BaseWave: *NewBaseWave(config),
	}

	c.computeWave(config)

	return c
}

func (c *TriangleWave) Y() float64 {
	return c.points[c.x]
}

func (c *TriangleWave) MoveNext() {
	c.x++
	if c.x >= float64(c.config.Period) {
		c.x = 0
	}
}

func (c *TriangleWave) computeWave(config WaveConfig) {
	const start = float64(0)
	const step = float64(1.0)
	// Compute values for one full period.
	end := float64(config.Period)

	c.points = make(map[float64]float64)

	// Calculate the range of oscillation (peak-to-peak amplitude)
	peakToPeak := 2.0 * float64(config.Amplitude)

	for x := start; x <= end; x += step {
		// Calculate the time relative to the current period.
		timeInPeriod := math.Mod(x, float64(config.Period))

		// Map the time within the period to a 0.0 to 1.0 range.
		normalizedTime := timeInPeriod / float64(config.Period)

		// The core of the triangle wave generation uses math.Abs and a sawtooth-like pattern.
		// `abs(2*normalizedTime - 1)` creates a triangle wave that oscillates between 0 and 1.
		// The final result is then scaled and shifted to match the config.
		y := peakToPeak * math.Abs(2*normalizedTime-1)

		// Shift the wave vertically and adjust for the base.
		y = float64(config.Target) + float64(config.Amplitude) - y

		c.points[x] = y
	}
}
