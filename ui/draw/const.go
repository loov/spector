package draw

import (
	"math"

	"github.com/egonelbre/spector/ui/g"
)

var (
	nan32  = math.Float32frombits(0x7FBFFFFF)
	noUV   = g.Vector{nan32, nan32}
	noClip = g.Rect{g.Vector{-8192, -8192}, g.Vector{+8192, +8192}}
)
