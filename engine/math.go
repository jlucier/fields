package engine

import (
	"math"
)

func Lerp(a float64, b float64, w float64) float64 {
	return (b-a)*w + a
}

func Clamp(x, a, b float64) float64 {
	return math.Min(b, math.Max(x, a))
}
