package gasmanager

import "math"

// WaveType represents the type of wave pattern.
type WaveType string

const (
	WaveTypeFlat     WaveType = "flat"
	WaveTypeSine     WaveType = "sine"
	WaveTypeSquare   WaveType = "square"
	WaveTypeTriangle WaveType = "triangle"
	WaveTypeSawtooth WaveType = "sawtooth"
)

// WaveConfig holds the configuration for a wave.
type WaveConfig struct {
	Period    uint64
	Amplitude uint64
	Target    uint64
}

// Wave interface for all wave types.
type Wave interface {
	Period() uint64
	Amplitude() uint64
	Target() uint64
	X() float64
	Y() float64
	MoveNext()
}

// wave implements the Wave interface for all wave types.
type wave struct {
	config WaveConfig
	x      float64
	points map[float64]float64
}

// NewWave creates a new wave of the specified type with the given configuration.
func NewWave(waveType WaveType, config WaveConfig) Wave {
	w := &wave{
		config: config,
		x:      0,
		points: make(map[float64]float64, config.Period+1),
	}

	switch waveType {
	case WaveTypeFlat:
		computeFlat(w, config)
	case WaveTypeSine:
		computeSine(w, config)
	case WaveTypeSquare:
		computeSquare(w, config)
	case WaveTypeTriangle:
		computeTriangle(w, config)
	case WaveTypeSawtooth:
		computeSawtooth(w, config)
	default:
		computeFlat(w, config) // default to flat
	}

	return w
}

func (w *wave) MoveNext() {
	w.x++
	if w.x >= float64(w.config.Period) {
		w.x = 0
	}
}

func (w *wave) Y() float64      { return w.points[w.x] }
func (w *wave) X() float64      { return w.x }
func (w *wave) Period() uint64  { return w.config.Period }
func (w *wave) Amplitude() uint64 { return w.config.Amplitude }
func (w *wave) Target() uint64  { return w.config.Target }

// Wave computation functions

func computeFlat(w *wave, config WaveConfig) {
	target := float64(config.Target)
	for x := 0.0; x <= float64(config.Period); x++ {
		w.points[x] = target
	}
}

func computeSine(w *wave, config WaveConfig) {
	period := float64(config.Period)
	amplitude := float64(config.Amplitude)
	target := float64(config.Target)
	b := (2 * math.Pi) / period

	for x := 0.0; x <= period; x++ {
		w.points[x] = amplitude*math.Sin(b*x) + target
	}
}

func computeSquare(w *wave, config WaveConfig) {
	period := float64(config.Period)
	highValue := float64(config.Target) + float64(config.Amplitude)
	lowValue := float64(config.Target) - float64(config.Amplitude)
	halfPeriod := period / 2.0

	for x := 0.0; x <= period; x++ {
		if math.Mod(x, period) < halfPeriod {
			w.points[x] = highValue
		} else {
			w.points[x] = lowValue
		}
	}
}

func computeTriangle(w *wave, config WaveConfig) {
	period := float64(config.Period)
	amplitude := float64(config.Amplitude)
	target := float64(config.Target)
	peakToPeak := 2.0 * amplitude

	for x := 0.0; x <= period; x++ {
		normalizedTime := math.Mod(x, period) / period
		w.points[x] = target + amplitude - peakToPeak*math.Abs(2*normalizedTime-1)
	}
}

func computeSawtooth(w *wave, config WaveConfig) {
	period := float64(config.Period)
	offset := float64(config.Target - config.Amplitude)
	rangeOfWave := float64(2 * config.Amplitude)

	for x := 0.0; x <= period; x++ {
		fractionalTime := math.Mod(x, period) / period
		w.points[x] = rangeOfWave*fractionalTime + offset
	}
}

// Legacy constructors for backwards compatibility

// FlatWave is an alias for wave (kept for backwards compatibility).
type FlatWave = wave

// NewFlatWave creates a new flat wave.
func NewFlatWave(config WaveConfig) *wave {
	return NewWave(WaveTypeFlat, config).(*wave)
}

// SineWave is an alias for wave (kept for backwards compatibility).
type SineWave = wave

// NewSineWave creates a new sine wave.
func NewSineWave(config WaveConfig) *wave {
	return NewWave(WaveTypeSine, config).(*wave)
}

// SquareWave is an alias for wave (kept for backwards compatibility).
type SquareWave = wave

// NewSquareWave creates a new square wave.
func NewSquareWave(config WaveConfig) *wave {
	return NewWave(WaveTypeSquare, config).(*wave)
}

// TriangleWave is an alias for wave (kept for backwards compatibility).
type TriangleWave = wave

// NewTriangleWave creates a new triangle wave.
func NewTriangleWave(config WaveConfig) *wave {
	return NewWave(WaveTypeTriangle, config).(*wave)
}

// SawtoothWave is an alias for wave (kept for backwards compatibility).
type SawtoothWave = wave

// NewSawtoothWave creates a new sawtooth wave.
func NewSawtoothWave(config WaveConfig) *wave {
	return NewWave(WaveTypeSawtooth, config).(*wave)
}
