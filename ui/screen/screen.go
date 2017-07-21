package screen

import (
	"github.com/egonelbre/spector/ui"
	"github.com/egonelbre/spector/ui/g"
)

type Screen struct {
	Root     *Area
	Registry *Registry
}

func New() *Screen {
	screen := &Screen{}
	screen.Root = NewArea(screen)
	screen.Registry = NewRegistry()
	return screen
}

func (screen *Screen) Update(ctx ui.Context) {
	screen.Root.Update(ctx)
}

type Registry struct {
	Editor map[string]EditorDefinition
}

func NewRegistry() *Registry {
	return &Registry{make(map[string]EditorDefinition)}
}

type EditorDefinition struct {
}

type Area struct {
	Screen *Screen
	Parent *Area // nil for Root

	Flex   float32 // 0.0 to 1.0
	Bounds g.Rect  // last-bounds

	// Either Split or Editor
	Vertical bool
	Split    []*Area
	Editor   *Editor
}

func NewArea(screen *Screen) *Area {
	return &Area{Screen: screen}
}

func (area *Area) Update(ctx ui.Context) {
	area.Bounds = ctx.Area
	if len(area.Split) > 0 {

	} else {
		area.Editor.Update(ctx)

		r := ctx.Area
		r.Min.X = r.Max.X - 10
		r.Max.Y = r.Min.Y + 10

		if r.Contains(ctx.Input.Mouse.Pos) {
			if ctx.Input.Mouse.Down {
				ctx.Draw.FillRect(&r, g.Green)
			} else {
				ctx.Draw.FillRect(&r, g.Blue)
			}
		} else {
			ctx.Draw.FillRect(&r, g.Color{0x80, 0x80, 0x80, 0xFF})
		}
	}
}

type Editor struct {
	Area *Area
	Name string
}

func (editor *Editor) Update(ctx ui.Context) {
	ctx.Draw.StrokeRect(&ctx.Area, 2, g.Color{0x80, 0x80, 0x80, 0xFF})
}

type Region struct {
}
