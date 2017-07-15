package draw

type Color struct{ R, G, B, A uint8 }

func (c Color) Transparent() bool { return c.A == 0x0 }

func (c Color) RGBA() (r, g, b, a uint8) { return c.R, c.G, c.B, c.A }

var (
	Black = Color{0, 0, 0, 255}
	White = Color{255, 255, 255, 255}
	Red   = Color{255, 0, 0, 255}
	Green = Color{0, 255, 0, 255}
	Blue  = Color{0, 0, 255, 255}
)

func ColorFloat(r, g, b, a float32) Color {
	return Color{sat8(r), sat8(g), sat8(b), sat8(a)}
}

func ColorHex(hex uint32) Color {
	return Color{
		R: uint8(hex >> 24),
		G: uint8(hex >> 16),
		B: uint8(hex >> 8),
		A: uint8(hex >> 0),
	}
}

func ColorHSLA(h, s, l, a float32) Color { return ColorFloat(hsla(h, s, l, a)) }
func ColorHSL(h, s, l float32) Color     { return ColorHSLA(h, s, l, 1) }

func sat8(v float32) uint8 {
	v *= 255.0
	if v >= 255 {
		return 255
	} else if v <= 0 {
		return 0
	}
	return uint8(v)
}

func hue(v1, v2, h float32) float32 {
	if h < 0 {
		h += 1
	}
	if h > 1 {
		h -= 1
	}
	if 6*h < 1 {
		return v1 + (v2-v1)*6*h
	} else if 2*h < 1 {
		return v2
	} else if 3*h < 2 {
		return v1 + (v2-v1)*(2.0/3.0-h)*6
	}

	return v1
}

func hsla(h, s, l, a float32) (r, g, b, ra float32) {
	if s == 0 {
		return l, l, l, a
	}

	var v2 float32
	if l < 0.5 {
		v2 = l * (1 + s)
	} else {
		v2 = (l + s) - s*l
	}

	v1 := 2*l - v2
	r = hue(v1, v2, h+1.0/3.0)
	g = hue(v1, v2, h)
	b = hue(v1, v2, h-1.0/3.0)
	ra = a

	return
}
