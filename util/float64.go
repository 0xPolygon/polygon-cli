package util

import "math"

func CompareFloatsWithTolerance(a, b, tolerance float64) bool {
	return math.Abs(a-b) < tolerance
}
