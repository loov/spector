package g

import "math"

const (
	Pi  = math.Pi
	Tau = 2 * math.Pi
)

func Pow(base, e float32) float32 { return float32(math.Pow(float64(base), float64(e))) }
func Mod(x, y float32) float32    { return float32(math.Mod(float64(x), float64(y))) }
func Sqr(v float32) float32       { return v * v }
func Sqrt(v float32) float32      { return float32(math.Sqrt(float64(v))) }

func Ceil(v float32) float32  { return float32(math.Ceil(float64(v))) }
func Floor(v float32) float32 { return float32(math.Floor(float64(v))) }

func Sin(v float32) float32 { return float32(math.Sin(float64(v))) }
func Cos(v float32) float32 { return float32(math.Cos(float64(v))) }

func Sincos(v float32) (float32, float32) {
	sn, cs := math.Sincos(float64(v))
	return float32(sn), float32(cs)
}

func Abs(v float32) float32 {
	if v < 0 {
		return -v
	}
	return v
}

func Min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func Max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func MinMax(a, b float32) (float32, float32) {
	if a < b {
		return a, b
	}
	return b, a
}

func Lerp(a, b, p float32) float32 {
	return a + (b-a)*p
}

func LerpClamp(a, b, p float32) float32 {
	if p < 0 {
		return a
	} else if p > 1 {
		return b
	}
	return a + (b-a)*p
}

func Clamp(v, min, max float32) float32 {
	if v < min {
		return min
	} else if v > max {
		return max
	}
	return v
}

func Clamp01(v float32) float32 {
	if v < 0 {
		return 0
	} else if v > 1 {
		return 1
	}
	return v
}

func Clamp1(v float32) float32 {
	if v < -1 {
		return -1
	} else if v > 1 {
		return 1
	}
	return v
}
