package draw

import "math"

func sqrt(v float32) float32 {
	return float32(math.Sqrt(float64(v)))
}

func min(x, y float32) float32 {
	if x <= y {
		return x
	}
	return y
}

func max(x, y float32) float32 {
	if x >= y {
		return x
	}
	return y
}
