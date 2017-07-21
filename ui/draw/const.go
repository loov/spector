package draw

import (
	"math"

	"github.com/egonelbre/spector/ui/g"
)

var (
	NaN32  = math.Float32frombits(0x7FBFFFFF)
	NoUV   = g.Vector{NaN32, NaN32}
	NoClip = g.Rect{g.Vector{-8192, -8192}, g.Vector{+8192, +8192}}
)
