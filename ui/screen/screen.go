package screen

import (
	"github.com/egonelbre/spector/ui"
	"github.com/egonelbre/spector/ui/g"
)

const (
	EditorMinSize  = 20
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
	Parent *Area // nil for Root

	Flex   float32 // 0.0 to 1.0
	Bounds g.Rect  // last-bounds

	State struct {
		Resizing  bool
		Active    int
		Splitters []Splitter
		Min, Max  float32
	}

	// Either Split or Editor
	Vertical bool
	Split    []*Area
	Editor   *Editor
}

type Splitter struct {
	Center float32
	Active bool
}

func NewRootArea(screen *Screen) *Area {
	root := NewArea(screen)
	root.Split = append(root.Split,
		NewArea(screen),
		NewArea(screen),
		NewArea(screen))

	for _, area := range root.Split {
		area.Parent = root
	}

	return root
}

func NewArea(screen *Screen) *Area {
	return &Area{
		Screen: screen,
		Parent: nil,
		Flex:   1,
	}
}

func (area *Area) Update(ctx *ui.Context) {
	area.Bounds = ctx.Area
	if !ctx.Input.Mouse.Down {
		area.State.Resizing = false
		area.State.Active = -1
	}

	if len(area.Split) > 0 {
		if len(area.State.Splitters) != len(area.Split) {
			area.State.Splitters = make([]Splitter, len(area.Split))
		}

		total := float32(0.0)
		for _, area := range area.Split {
			total += area.Flex
		}

		width := area.Bounds.Dx()
		left := float32(0.0)
		for i, child := range area.Split {
			right := left + width*child.Flex/total
			area.State.Splitters[i].Center = right
			area.State.Splitters[i].Active = false
			child.Update(ctx.Column(left, right))
			left = right
		}

		if area.State.Resizing && area.State.Active >= 0 && area.State.Active < len(area.State.Splitters) {
			area.State.Splitters[area.State.Active].Center = ctx.Input.Mouse.Pos.X
			area.State.Splitters[area.State.Active].Active = true
		} else {
			if area.Bounds.Contains(ctx.Input.Mouse.Pos) {
				for i := range area.State.Splitters {
					splitter := &area.State.Splitters[i]
					if g.Abs(splitter.Center-ctx.Input.Mouse.Pos.X) < SplitterRadius {
						splitter.Active = true
						if ctx.Input.Mouse.Down {
							area.State.Resizing = true
							area.State.Active = i
							area.State.Min, area.State.Max = area.Bounds.Min.X, area.Bounds.Max.X
							if i-1 >= 0 {
								area.State.Min = area.State.Splitters[i-1].Center + EditorMinSize
								if area.State.Min > splitter.Center {
									area.State.Min = splitter.Center
								}
							}
							if i+1 < len(area.State.Splitters) {
								area.State.Max = area.State.Splitters[i+1].Center - EditorMinSize
								if area.State.Max < splitter.Center {
									area.State.Max = splitter.Center
								}
							}
						}
					}
				}
			}
		}

		for i := range area.State.Splitters {
			splitter := &area.State.Splitters[i]

			r := g.Rect{}
			r.Min.X = splitter.Center - SplitterRadius
			r.Min.Y = area.Bounds.Min.Y
			r.Max.X = splitter.Center + SplitterRadius
			r.Max.Y = area.Bounds.Max.Y

			if splitter.Active {
				if ctx.Input.Mouse.Down {
					ctx.Draw.FillRect(&r, g.Green)
				} else {
					ctx.Draw.FillRect(&r, g.Red)
				}
			} else {
				ctx.Draw.FillRect(&r, g.Color{0x80, 0x80, 0x80, 0xff})
			}
		}
	} else {
		area.Editor.Update(ctx)

		r := ctx.Area
		r.Min.X = r.Max.X - EditorMinSize
		r.Max.Y = r.Min.Y + EditorMinSize

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

type Editor struct {
	Area *Area
	Name string
}

func (editor *Editor) Update(ctx *ui.Context) {
	ctx.Draw.FillRect(&ctx.Area, g.Color{0xEE, 0xEE, 0xEE, 0xFF})
}

type Region struct {
}
