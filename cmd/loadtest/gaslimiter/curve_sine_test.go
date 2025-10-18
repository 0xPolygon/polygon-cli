package gaslimiter

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSineCurve(t *testing.T) {
	config := CurveConfig{
		Target:    10,
		Amplitude: 5,
		Period:    20,
	}
	curve := NewSineCurve(config)

	expectedPoints := map[float64]float64{
		0:  10.00000,
		1:  11.545085,
		2:  12.938926,
		3:  14.045085,
		4:  14.755283,
		5:  15.000000,
		6:  14.755283,
		7:  14.045085,
		8:  12.938926,
		9:  11.545085,
		10: 10.000000,
		11: 8.454915,
		12: 7.061074,
		13: 5.954915,
		14: 5.244717,
		15: 5.000000,
		16: 5.244717,
		17: 5.954915,
		18: 7.061074,
		19: 8.454915,
	}

	assert.Equal(t, config.Period, curve.Period())
	assert.Equal(t, config.Amplitude, curve.Amplitude())
	assert.Equal(t, config.Target, curve.Target())

	for i := uint64(0); i <= config.Period; i++ {
		if !areFloatsEqualGivenTheTolerance(expectedPoints[curve.X()], curve.Y(), 0.00001) {
			t.Errorf("At x=%f, expected y=%f, got y=%f, diff=%f", curve.X(), expectedPoints[curve.X()], curve.Y(), expectedPoints[curve.X()]-curve.Y())
		}
		curve.MoveNext()
	}
}

func areFloatsEqualGivenTheTolerance(a, b, tolerance float64) bool {
	return math.Abs(a-b) < tolerance
}
