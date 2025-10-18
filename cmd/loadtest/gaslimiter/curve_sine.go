package gaslimiter

import "math"

type SineCurve struct {
	BaseCurve
	points map[float64]float64
}

func NewSineCurve(config CurveConfig) *SineCurve {
	c := &SineCurve{
		BaseCurve: *NewBaseCurve(config),
		points:    computeCurve(config),
	}

	return c
}

func (c *SineCurve) Y() float64 {
	return c.points[c.x]
}

func (c *SineCurve) MoveNext() {
	c.x++
	if c.x >= float64(c.config.Period) {
		c.x = 0
	}
}

func computeCurve(config CurveConfig) map[float64]float64 {
	const start = float64(0)
	const step = float64(1.0)
	end := float64(config.Period)

	points := map[float64]float64{}

	// The `b` parameter in the generalized sine formula is derived from the period.
	b := (2 * math.Pi) / float64(config.Period)

	for x := start; x <= end; x += step {
		// Apply the generalized sine formula: y = A * sin(b(x)) + k
		y := float64(config.Amplitude)*math.Sin(b*x) + float64(config.Target)
		points[float64(x)] = float64(y)
	}

	return points
}
