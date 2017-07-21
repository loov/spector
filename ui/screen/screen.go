package screen

import (
	"github.com/egonelbre/spector/ui"
	"github.com/egonelbre/spector/ui/g"
)

const (
	EditorMinSize  = 50
	SplitterTab    = 20
	SplitterRadius = 2
)

type Screen struct {
	Root     *Area
	Registry *Registry
}

func New() *Screen {
	screen := &Screen{}
	screen.Root = NewRootArea(screen)
	screen.Registry = NewRegistry()
	return screen
}

func (screen *Screen) Update(ctx *ui.Context) {
	screen.Root.Update(ctx)
}

type Area struct {
	Screen *Screen
	Parent *Area  // nil for Root
	Bounds g.Rect // last-bounds

	// Either Split or Editor
	Vertical  bool
	Splitters []*Splitter
	Areas     []*Area
	Editor    *Editor
}

type Splitter struct {
	Parent *Area
	Index  int

	RelativeCenter float32

	Active   bool
	Resizing bool
}

func (splitter *Splitter) MinMax() (float32, float32) {
	rmin, rmax := float32(0), float32(1)
	if 0 <= splitter.Index-1 {
		rmin = splitter.Parent.Splitters[splitter.Index-1].RelativeCenter
	}
	if splitter.Index+1 < len(splitter.Parent.Splitters) {
		rmax = splitter.Parent.Splitters[splitter.Index+1].RelativeCenter
	}

	min := splitter.toAbsolute(rmin) + EditorMinSize
	center := splitter.Center()
	max := splitter.toAbsolute(rmax) - EditorMinSize
	if min > center {
		min = center
	}
	if max < center {
		max = center
	}

	return min, max
}

func (splitter *Splitter) toAbsolute(v float32) float32 {
	return g.LerpClamp(v, splitter.Parent.Bounds.Min.X, splitter.Parent.Bounds.Max.X)
}
func (splitter *Splitter) toRelative(v float32) float32 {
	return g.InverseLerpClamp(v, splitter.Parent.Bounds.Min.X, splitter.Parent.Bounds.Max.X)
}

func (splitter *Splitter) Center() float32 {
	return splitter.toAbsolute(splitter.RelativeCenter)
}

func (splitter *Splitter) SetCenter(absolute float32) {
	splitter.RelativeCenter = splitter.toRelative(absolute)
}

func (splitter *Splitter) Rect() g.Rect {
	center := splitter.Center()
	return g.Rect{
		Min: g.Vector{X: center - SplitterRadius, Y: splitter.Parent.Bounds.Min.Y},
		Max: g.Vector{X: center + SplitterRadius, Y: splitter.Parent.Bounds.Max.Y},
	}
}

func NewRootArea(screen *Screen) *Area {
	root := NewArea(screen)
	root.Areas = append(root.Areas, NewArea(screen), NewArea(screen), NewArea(screen))
	for i, child := range root.Areas {
		child.Parent = root

		splitter := &Splitter{}
		splitter.Parent = root
		splitter.Index = i
		splitter.RelativeCenter = float32(i+1) / float32(len(root.Areas))

		root.Splitters = append(root.Splitters, splitter)
	}
	return root
}

func NewSplitter(parent *Area) *Splitter {
	return &Splitter{Parent: parent}
}

func NewArea(screen *Screen) *Area {
	return &Area{
		Screen: screen,
		Parent: nil,
	}
}

func (area *Area) Update(ctx *ui.Context) {
	area.Bounds = ctx.Area

	if len(area.Splitters) != len(area.Areas) {
		panic("len(splitters) != len(areas)")
	}

	if len(area.Splitters) > 0 {
		left := float32(0.0)
		for i, splitter := range area.Splitters {
			child := area.Areas[i]
			right := splitter.Center()
			child.Update(ctx.Column(left, right))
			left = right
		}

		for _, splitter := range area.Splitters {
			splitter.Update(ctx)
		}

	} else {
		area.Editor.Update(ctx)

		r := ctx.Area
		r.Min.X = r.Max.X - SplitterTab
		r.Max.Y = r.Min.Y + SplitterTab

		if r.Contains(ctx.Input.Mouse.Pos) {
			if ctx.Input.Mouse.Down {
				ctx.Draw.FillRect(&r, g.Green)
			} else {
				ctx.Draw.FillRect(&r, g.Red)
			}
		} else {
			ctx.Draw.FillRect(&r, g.Color{0x80, 0x80, 0x80, 0xFF})
		}
	}
}

func (splitter *Splitter) Update(ctx *ui.Context) {
	r := splitter.Rect()
	inside := r.Contains(ctx.Input.Mouse.Pos) && ctx.Input.Mouse.Capture == nil
	if inside && ctx.Input.Mouse.Down {
		splitter.Resizing = true
		ctx.Input.Mouse.Capture = func() bool {
			min, max := splitter.MinMax()
			splitter.SetCenter(g.Clamp(ctx.Input.Mouse.Pos.X, min, max))
			splitter.Resizing = ctx.Input.Mouse.Down
			return !ctx.Input.Mouse.Down
		}
	}

	if splitter.Resizing {
		ctx.Draw.FillRect(&r, g.Green)
	} else if inside {
		ctx.Draw.FillRect(&r, g.Red)
	} else {
		ctx.Draw.FillRect(&r, g.Color{0x80, 0x80, 0x80, 0xff})
	}
}

type Editor struct {
	Area *Area
	Name string
}

func (editor *Editor) Update(ctx *ui.Context) {
	ctx.Draw.FillRect(&ctx.Area, g.Color{0xEE, 0xEE, 0xEE, 0xFF})
}

type Region struct {
}
