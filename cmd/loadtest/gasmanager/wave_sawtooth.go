package gasmanager

import "math"

type SawtoothWave struct {
	BaseWave
	points map[float64]float64
}

func NewSawtoothWave(config WaveConfig) *SawtoothWave {
	c := &SawtoothWave{
		BaseWave: *NewBaseWave(config),
	}

	c.computeWave(config)

	return c
}

func (c *SawtoothWave) Y() float64 {
	return c.points[c.x]
}

func (c *SawtoothWave) MoveNext() {
	c.x++
	if c.x >= float64(c.config.Period) {
		c.x = 0
	}
}

func (c *SawtoothWave) computeWave(config WaveConfig) {
	const start = float64(0)
	const step = float64(1.0)
	// The end is defined by the Period, but the loop needs to go up to but not
	// including the end to avoid a double point at the end of the period.
	end := float64(config.Period)

	c.points = map[float64]float64{}

	// Calculate the minimum value of the wave, which is the offset.
	offset := float64(config.Target - config.Amplitude)

	// Calculate the total peak-to-peak range.
	rangeOfWave := float64(2 * config.Amplitude)

	for x := start; x <= end; x += step {
		// math.Mod finds the remainder of the division of x by the period.
		// This causes the value to repeat every 'config.Period' units.
		// Dividing by config.Period scales this to a 0.0 to 1.0 range.
		fractionalTime := math.Mod(x, float64(config.Period)) / float64(config.Period)

		// Scale the fractional time to the desired amplitude range.
		// Add the offset to shift the wave to the correct target value.
		y := (rangeOfWave * fractionalTime) + offset

		c.points[x] = y
	}
}
