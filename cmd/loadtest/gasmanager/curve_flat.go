package gasmanager

type FlatCurve struct {
	BaseCurve
	points map[float64]float64
}

func NewFlatCurve(config CurveConfig) *FlatCurve {
	c := &FlatCurve{
		BaseCurve: *NewBaseCurve(config),
	}

	c.computeCurve(config)

	return c
}

func (c *FlatCurve) Y() float64 {
	return c.points[c.x]
}

func (c *FlatCurve) MoveNext() {
}

func (c *FlatCurve) computeCurve(config CurveConfig) {
	c.points = map[float64]float64{
		0: float64(config.Target),
	}
}
