package gasmanager

import (
	"testing"

	"github.com/0xPolygon/polygon-cli/util"
	"github.com/stretchr/testify/assert"
)

func TestWaves(t *testing.T) {
	type TestCase struct {
		name           string
		config         WaveConfig
		expectedPoints map[float64]float64
		createWave     func(WaveConfig) Wave
	}

	testCases := []TestCase{
		{
			name:       "Flat Wave",
			config:     WaveConfig{Target: 10, Amplitude: 5, Period: 10},
			createWave: func(config WaveConfig) Wave { return NewFlatWave(config) },
			expectedPoints: map[float64]float64{
				0: 10.00000, 1: 10.00000, 2: 10.00000, 3: 10.00000, 4: 10.00000,
				5: 10.00000, 6: 10.00000, 7: 10.00000, 8: 10.00000, 9: 10.00000,
			},
		},
		{
			name:       "Sine Wave",
			config:     WaveConfig{Target: 10, Amplitude: 5, Period: 20},
			createWave: func(config WaveConfig) Wave { return NewSineWave(config) },
			expectedPoints: map[float64]float64{
				0: 10.00000, 1: 11.545085, 2: 12.938926, 3: 14.045085, 4: 14.755283,
				5: 15.000000, 6: 14.755283, 7: 14.045085, 8: 12.938926, 9: 11.545085,
				10: 10.000000, 11: 8.454915, 12: 7.061074, 13: 5.954915, 14: 5.244717,
				15: 5.000000, 16: 5.244717, 17: 5.954915, 18: 7.061074, 19: 8.454915,
			},
		},
		{
			name:       "Sawtooth Wave",
			config:     WaveConfig{Target: 10, Amplitude: 5, Period: 20},
			createWave: func(config WaveConfig) Wave { return NewSawtoothWave(config) },
			expectedPoints: map[float64]float64{
				0: 5.0, 1: 5.5, 2: 6.0, 3: 6.5, 4: 7.0,
				5: 7.5, 6: 8.0, 7: 8.5, 8: 9.0, 9: 9.5,
				10: 10.0, 11: 10.5, 12: 11.0, 13: 11.5, 14: 12.0,
				15: 12.5, 16: 13.0, 17: 13.5, 18: 14.0, 19: 14.5,
			},
		},
		{
			name:       "Triangle Wave",
			config:     WaveConfig{Target: 10, Amplitude: 5, Period: 20},
			createWave: func(config WaveConfig) Wave { return NewTriangleWave(config) },
			expectedPoints: map[float64]float64{
				0: 5.0, 1: 6.0, 2: 7.0, 3: 8.0, 4: 9.0,
				5: 10.0, 6: 11.0, 7: 12.0, 8: 13.0, 9: 14.0,
				10: 15.0, 11: 14.0, 12: 13.0, 13: 12.0, 14: 11.0,
				15: 10.0, 16: 9.0, 17: 8.0, 18: 7.0, 19: 6.0,
			},
		},
		{
			name:       "Square Wave",
			config:     WaveConfig{Target: 10, Amplitude: 5, Period: 20},
			createWave: func(config WaveConfig) Wave { return NewSquareWave(config) },
			expectedPoints: map[float64]float64{
				0: 15.00000, 1: 15.00000, 2: 15.00000, 3: 15.00000, 4: 15.00000,
				5: 15.00000, 6: 15.00000, 7: 15.00000, 8: 15.00000, 9: 15.00000,
				10: 5.00000, 11: 5.00000, 12: 5.00000, 13: 5.00000, 14: 5.00000,
				15: 5.00000, 16: 5.00000, 17: 5.00000, 18: 5.00000, 19: 5.00000,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			wave := tc.createWave(tc.config)

			assert.Equal(t, tc.config.Period, wave.Period())
			assert.Equal(t, tc.config.Amplitude, wave.Amplitude())
			assert.Equal(t, tc.config.Target, wave.Target())

			startPoint := wave.X()
			points := map[float64]float64{}
			for i := 0; i < len(tc.expectedPoints); i++ {
				points[wave.X()] = wave.Y()
				wave.MoveNext()
				if wave.X() == startPoint {
					break
				}
			}

			assert.Equal(t, len(tc.expectedPoints), len(points))

			for k, v := range tc.expectedPoints {
				if !util.CompareFloatsWithTolerance(points[k], v, 0.00001) {
					t.Errorf("At x=%f, expected y=%f, got y=%f, diff=%f", k, v, points[k], v-points[k])
				}
			}
		})
	}

}
