package gasmanager

type CurveConfig struct {
	Period    uint64
	Amplitude uint64
	Target    uint64
}

type Curve interface {
	Period() uint64
	Amplitude() uint64
	Target() uint64
	X() float64
	Y() float64
	MoveNext()
}

type BaseCurve struct {
	config CurveConfig
	x      float64
}

func NewBaseCurve(config CurveConfig) *BaseCurve {
	return &BaseCurve{
		config: config,
		x:      0,
	}
}

func (c *BaseCurve) Period() uint64 {
	return c.config.Period
}

func (c *BaseCurve) Amplitude() uint64 {
	return c.config.Amplitude
}

func (c *BaseCurve) Target() uint64 {
	return c.config.Target
}

func (c *BaseCurve) X() float64 {
	return c.x
}
