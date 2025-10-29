package gasmanager

import "math"

type SineWave struct {
	BaseWave
	points map[float64]float64
}

func NewSineWave(config WaveConfig) *SineWave {
	c := &SineWave{
		BaseWave: *NewBaseWave(config),
	}

	c.computeWave(config)

	return c
}

func (c *SineWave) Y() float64 {
	return c.points[c.x]
}

func (c *SineWave) MoveNext() {
	c.x++
	if c.x >= float64(c.config.Period) {
		c.x = 0
	}
}

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
