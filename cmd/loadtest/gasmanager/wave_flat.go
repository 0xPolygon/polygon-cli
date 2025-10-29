package gasmanager

type FlatWave struct {
	BaseWave
	points map[float64]float64
}

func NewFlatWave(config WaveConfig) *FlatWave {
	c := &FlatWave{
		BaseWave: *NewBaseWave(config),
	}

	c.computeWave(config)

	return c
}

func (c *FlatWave) Y() float64 {
	return c.points[c.x]
}

func (c *FlatWave) MoveNext() {
	c.x++
	if c.x >= float64(c.config.Period) {
		c.x = 0
	}
}

func (c *FlatWave) computeWave(config WaveConfig) {
	const start = float64(0)
	const step = float64(1.0)
	end := float64(config.Period)

	c.points = map[float64]float64{}

	for x := start; x <= end; x += step {
		c.points[float64(x)] = float64(config.Target)
	}
}
