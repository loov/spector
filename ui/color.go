package ui

type Color struct{ R, G, B, A uint8 }

func (c Color) RGBA() (r, g, b, a uint8) { return c.R, c.G, c.B, c.A }

type StateColors struct {
	Text    Color
	Default Color
	Hot     Color
	Active  Color
	Clicked Color
	Border  Color
}

func ColorHex(hex uint32) Color {
	return Color{
		R: uint8(hex >> 24),
		G: uint8(hex >> 16),
		B: uint8(hex >> 8),
		A: uint8(hex >> 0),
	}
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
