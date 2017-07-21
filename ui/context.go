package ui

import (
	"time"

	"github.com/egonelbre/spector/ui/draw"
	"github.com/egonelbre/spector/ui/g"
)

type Context struct {
	*Render
	Input *Input

	Area g.Rect

	ID    string
	Index int
	Count int
}

func NewContext() *Context {
	return &Context{
		Render: &Render{},
		Input:  &Input{},
	}
}

// TODO: rename to Layers
type Render struct {
	Frame  draw.Frame
	Draw   *draw.List
	Hover  *draw.List
	Cursor *draw.List
}

type Input struct {
	Time  time.Time
	Mouse Mouse
}

type Mouse struct {
	Pos  g.Vector
	Down bool
	Last struct {
		Pos  g.Vector
		Down bool
	}
	Drag func(*Context) (done bool)
}

func (context *Context) BeginFrame(area g.Rect) {
	context.Area = area
	context.Render.BeginFrame()
	context.Input.Mouse.BeginFrame()
}

func (render *Render) BeginFrame() {
	render.Frame.Reset()
	render.Draw = render.Frame.Layer()
	render.Hover = render.Frame.Layer()
	render.Cursor = render.Frame.Layer()
}

func (mouse *Mouse) BeginFrame() {
	mouse.Last.Pos = mouse.Pos
	mouse.Last.Down = mouse.Down
}
