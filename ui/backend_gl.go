package ui

import "github.com/go-gl/gl/v2.1/gl"

type state struct {
	back      Color
	fore      Color
	fontColor Color
	fontName  string
	fontSize  float32
}

type glbackend struct {
	atlas   *FontAtlas
	current state
}

func NewGLBackend() Backend {
	back := &glbackend{}
	var err error
	back.atlas, err = NewFontAtlas("~DejaVuSans.ttf", 72, 12)
	if err != nil {
		panic(err)
	}
	back.current.fontSize = 12
	return back
}

func (back *glbackend) Clone() Backend {
	return &glbackend{
		atlas:   back.atlas,
		current: back.current,
	}
}

func (back *glbackend) Fill(b Bounds) {
	gl.Color4ub(back.current.back.RGBA())
	gl.Begin(gl.QUADS)
	{
		gl.Vertex2f(b.Min.X, b.Min.Y)
		gl.Vertex2f(b.Max.X, b.Min.Y)
		gl.Vertex2f(b.Max.X, b.Max.Y)
		gl.Vertex2f(b.Min.X, b.Max.Y)
	}
	gl.End()
}

func (back *glbackend) Stroke(b Bounds) {
	gl.Color4ub(back.current.fore.RGBA())
	gl.Begin(gl.LINE_LOOP)
	{
		gl.Vertex2f(b.Min.X, b.Min.Y)
		gl.Vertex2f(b.Max.X, b.Min.Y)
		gl.Vertex2f(b.Max.X, b.Max.Y)
		gl.Vertex2f(b.Min.X, b.Max.Y)
	}
	gl.End()
}

func (back *glbackend) SetFore(c Color)      { back.current.fore = c }
func (back *glbackend) SetBack(c Color)      { back.current.back = c }
func (back *glbackend) SetFontColor(c Color) { back.current.fontColor = c }

func (back *glbackend) SetFont(name string, size float32) {
	back.current.fontName = name
	back.current.fontSize = size
}

func (back *glbackend) Text(text string, bounds Bounds) {
	gl.Color4ub(back.current.fontColor.RGBA())
	back.atlas.Draw(text, bounds)
}

func (back *glbackend) Measure(text string) Point {
	return back.atlas.Measure(text)
}
