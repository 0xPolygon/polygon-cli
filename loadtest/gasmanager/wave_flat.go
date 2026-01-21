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

// computeWave sets all points to the target value.
func (c *FlatWave) computeWave(config WaveConfig) {
	target := float64(config.Target)
	for x := 0.0; x <= float64(config.Period); x++ {
		c.points[x] = target
	}
}
