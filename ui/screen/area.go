package screen

import (
	"github.com/egonelbre/spector/ui"
	"github.com/egonelbre/spector/ui/g"
)

const (
	JoinSplitSize    = 20
	AreaBorderRadius = 3
)

var (
	AreaBackground       = g.Color{0x80, 0x80, 0x80, 0xFF}
	BorderColor          = g.Color{0x80, 0x80, 0x80, 0xFF}
	BorderHighlightColor = g.Color{0xFF, 0x80, 0x80, 0xFF}
)

type Area struct {
	Screen *Screen
	Bounds g.Rect // last-bounds
	Editor *Editor
}

func NewArea(screen *Screen) *Area {
	return &Area{
		Screen: screen,
	}
}

func (area *Area) Clone() *Area {
	clone := &Area{}
	clone.Screen = area.Screen
	clone.Bounds = area.Bounds
	clone.Editor = area.Editor.Clone()

	return clone
}

func (area *Area) Update(ctx *ui.Context) {
	area.Bounds = ctx.Area
	area.Editor.Update(ctx)
}
