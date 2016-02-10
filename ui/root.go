package ui

import "github.com/go-gl/gl/v2.1/gl"

type ID int32

type Point struct{ X, Y float32 }
type Bounds struct{ Min, Max Point }

func Rect(x, y, w, h float32) Bounds {
	if w < 0 {
		x, w = x+w, -w
	}
	if h < 0 {
		y, h = y+h, -h
	}
	return Bounds{
		Min: Point{x, y},
		Max: Point{x + w, y + h},
	}
}

func (b Bounds) Contains(p Point) bool {
	return b.Min.X <= p.X && p.X <= b.Max.X &&
		b.Min.Y <= p.Y && p.Y <= b.Max.Y
}

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
		Text:    ColorHex(0x333333ff),
		Default: ColorHex(0xEEEEECff),
		Hot:     ColorHex(0xD3D7CFff),
		Active:  ColorHex(0xFCE94Fff),
		Clicked: ColorHex(0xFF0000ff),
		Border:  ColorHex(0xD3D7CFff),
	}
)

type Mouse struct {
	Position Point
	Down     bool
	PDown    bool
}

func (m *Mouse) Released() bool {
	return !m.Down && m.PDown
}

type Input struct {
	Mouse  Mouse
	Hot    ID
	Active ID
}

type State struct {
	Input Input
}

func (state *State) SelectColor(c Color) {
	gl.Color4ub(c.RGBA())
}

func (b Bounds) render() {
	gl.Vertex2f(b.Min.X, b.Min.Y)
	gl.Vertex2f(b.Max.X, b.Min.Y)
	gl.Vertex2f(b.Max.X, b.Max.Y)
	gl.Vertex2f(b.Min.X, b.Max.Y)
}

func (state *State) Rect(b Bounds) {
	gl.Begin(gl.QUADS)
	b.render()
	gl.End()
}

func (state *State) StrokeRect(b Bounds) {
	gl.Begin(gl.LINE_LOOP)
	b.render()
	gl.End()
}

func (state *State) Text(text string, b Bounds) {

}

func (state *State) Button(text string, b Bounds) (pressed bool) {
	color := ButtonColor.Default
	if b.Contains(state.Input.Mouse.Position) {
		color = ButtonColor.Hot
		if state.Input.Mouse.Down {
			color = ButtonColor.Active
		} else if state.Input.Mouse.Released() {
			color = ButtonColor.Clicked
			pressed = true
		}
	}

	state.SelectColor(color)
	state.Rect(b)

	state.SelectColor(ButtonColor.Border)
	state.StrokeRect(b)

	state.SelectColor(ButtonColor.Text)
	state.Text(text, b)

	return
}
