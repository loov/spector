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

func clamp(v float32, min, max float32) float32 {
	if v > max {
		return max
	} else if v < min {
		return min
	}
	return v
}

func max(x, y float32) float32 {
	if x >= y {
		return x
	}
	return y
}

func abs(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}

//TODO: optimize
func cos(x float32) float32 { return float32(math.Cos(float64(x))) }
func sin(x float32) float32 { return float32(math.Sin(float64(x))) }

const pi = math.Pi
