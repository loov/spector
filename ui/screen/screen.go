package screen

import (
	"github.com/egonelbre/spector/ui"
	"github.com/egonelbre/spector/ui/g"
)

const (
	EditorMinSize = 50
)

type Screen struct {
	Registry *Registry
	Bounds   g.Rect
	Areas    []*Area
}

func New() *Screen {
	screen := &Screen{}
	screen.Registry = NewRegistry()
	screen.Areas = []*Area{
		NewTestArea(screen, g.Rect{g.Vector{0.0, 0.0}, g.Vector{0.5, 0.5}}),
		NewTestArea(screen, g.Rect{g.Vector{0.5, 0.0}, g.Vector{1.0, 0.5}}),

		NewTestArea(screen, g.Rect{g.Vector{0.0, 0.5}, g.Vector{0.5, 1.0}}),
		NewTestArea(screen, g.Rect{g.Vector{0.5, 0.5}, g.Vector{1.0, 1.0}}),
	}
	return screen
}

func (screen *Screen) Update(ctx *ui.Context) {
	screen.Bounds = ctx.Area
	for _, area := range screen.Areas {
		area.Update(ctx.Child(ctx.Area.Subset(area.RelBounds)))
	}
}

func NewTestArea(screen *Screen, relbounds g.Rect) *Area {
	child := NewArea(screen)
	child.RelBounds = relbounds
	child.Editor = &Editor{
		Area:  child,
		Color: g.RandColor(0.9, 0.9),
	}
	return child
}

type Editor struct {
	Area   *Area
	Bounds g.Rect
	Color  g.Color
}

func (editor *Editor) Update(ctx *ui.Context) {
	editor.Bounds = ctx.Area
	ctx.Draw.FillRect(&editor.Bounds, editor.Color)
}

func (editor *Editor) Clone() *Editor {
	clone := &Editor{}
	clone.Area = editor.Area
	clone.Color = editor.Color
	clone.Color = g.RandColor(0.9, 0.9) // TODO: copy
	return clone
}

type Region struct {
}
