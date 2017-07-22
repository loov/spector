package screen

import (
	"github.com/egonelbre/spector/ui"
	"github.com/egonelbre/spector/ui/g"
)

const (
	EditorMinSize = 50
)

type Screen struct {
	Root     *Area
	Registry *Registry
}

func New() *Screen {
	screen := &Screen{}
	screen.Root = NewRootTestArea(screen)
	screen.Registry = NewRegistry()
	return screen
}

func (screen *Screen) Update(ctx *ui.Context) {
	screen.Root.Update(ctx)
}

func NewRootTestArea(screen *Screen) *Area {
	root := NewArea(screen)
	areas := []*Area{
		NewTestArea(screen, root),
		NewTestArea(screen, root),
		NewTestArea(screen, root),
	}
	for i, child := range areas {
		child.Parent = root

		splitter := &Splitter{}
		splitter.Owner = root
		splitter.Content = areas[i]
		splitter.Index = i
		splitter.RelativeCenter = float32(i+1) / float32(len(areas))

		root.Splitters = append(root.Splitters, splitter)
	}
	return root
}

func NewTestArea(screen *Screen, parent *Area) *Area {
	child := NewArea(screen)
	child.Parent = parent
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
	return clone
}

type Region struct {
}
