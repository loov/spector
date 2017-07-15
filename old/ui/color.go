package ui

type Color struct{ R, G, B, A uint8 }

func ColorFloat(r, g, b, a float32) Color {
	return Color{
		R: sat8(r * 255.0),
		G: sat8(g * 255.0),
		B: sat8(b * 255.0),
		A: sat8(a * 255.0),
	}
}

func sat8(v float32) uint8 {
	if v >= 255 {
		return 255
	}
	if v <= 0 {
		return 0
	}
	return uint8(v)
}

func ColorHex(hex uint32) Color {
	return Color{
		R: uint8(hex >> 24),
		G: uint8(hex >> 16),
		B: uint8(hex >> 8),
		A: uint8(hex >> 0),
	}
}

func (c Color) RGBA() (r, g, b, a uint8) {
	return c.R, c.G, c.B, c.A
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

func ColorHSLA(h, s, l, a float32) Color { return ColorFloat(hsla(h, s, l, a)) }

type StateColors struct {
	Text    Color
	Default Color
	Hot     Color
	Active  Color
	Clicked Color
	Border  Color
}

var (
	ButtonColor = StateColors{
		Text:    ColorHex(0x000000ff),
		Default: ColorHex(0xEEEEECff),
		Hot:     ColorHex(0xD3D7CFff),
		Active:  ColorHex(0xFCE94Fff),
		Clicked: ColorHex(0xFF0000ff),
		Border:  ColorHex(0xAAAAAAff),
	}
)
