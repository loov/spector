package ui

import (
	"github.com/go-gl/gl/v2.1/gl"
)

type ID int32

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
	Font  *FontAtlas
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
	state.Font.Draw(b, text)
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

func (state *State) Panel(b Bounds, fn func()) {
	state.SelectColor(ButtonColor.Default)
	state.Rect(b)

	state.SelectColor(ButtonColor.Border)
	state.StrokeRect(b)

	gl.PushMatrix()
	state.Input.Mouse.Position.X -= b.Min.X
	state.Input.Mouse.Position.Y -= b.Min.Y
	defer func() {
		state.Input.Mouse.Position.X += b.Min.X
		state.Input.Mouse.Position.Y += b.Min.Y
		gl.PopMatrix()
	}()

	gl.Translatef(b.Min.X, b.Min.Y, 0)
	fn()
}
