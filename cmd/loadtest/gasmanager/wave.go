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
	points map[float64]float64
}

func NewBaseWave(config WaveConfig) *BaseWave {
	return &BaseWave{
		config: config,
		x:      0,
		points: make(map[float64]float64),
	}
}

// MoveNext advances the wave to the next position.
func (w *BaseWave) MoveNext() {
	w.x++
	if w.x >= float64(w.config.Period) {
		w.x = 0
	}
}

// Y returns the current value of the wave at position x.
func (w *BaseWave) Y() float64 {
	return w.points[w.x]
}

func (w *BaseWave) Period() uint64 {
	return w.config.Period
}

func (w *BaseWave) Amplitude() uint64 {
	return w.config.Amplitude
}

func (w *BaseWave) Target() uint64 {
	return w.config.Target
}

func (w *BaseWave) X() float64 {
	return w.x
}
