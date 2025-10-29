package gasmanager

type WaveConfig struct {
	Period    uint64
	Amplitude uint64
	Target    uint64
}

type Wave interface {
	Period() uint64
	Amplitude() uint64
	Target() uint64
	X() float64
	Y() float64
	MoveNext()
}

type BaseWave struct {
	config WaveConfig
	x      float64
}

func NewBaseWave(config WaveConfig) *BaseWave {
	return &BaseWave{
		config: config,
		x:      0,
	}
}

func (c *BaseWave) Period() uint64 {
	return c.config.Period
}

func (c *BaseWave) Amplitude() uint64 {
	return c.config.Amplitude
}

func (c *BaseWave) Target() uint64 {
	return c.config.Target
}

func (c *BaseWave) X() float64 {
	return c.x
}
