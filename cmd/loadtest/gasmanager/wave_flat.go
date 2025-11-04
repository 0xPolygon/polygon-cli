package gasmanager

// FlatWave represents a wave that maintains a constant target value over its period.
type FlatWave struct {
	*BaseWave
}

// NewFlatWave creates a new FlatWave with the given configuration.
func NewFlatWave(config WaveConfig) *FlatWave {
	c := FlatWave{
		BaseWave: NewBaseWave(config),
	}
	c.computeWave(config)
	return &c
}

// computeWave calculates the wave points based on the configuration
// by setting all points to the target value.
func (c *FlatWave) computeWave(config WaveConfig) {
	const start = float64(0)
	const step = float64(1.0)
	end := float64(config.Period)

	c.points = map[float64]float64{}

	for x := start; x <= end; x += step {
		c.points[float64(x)] = float64(config.Target)
	}
}
