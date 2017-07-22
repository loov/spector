package screen

import (
	"github.com/egonelbre/spector/ui"
	"github.com/egonelbre/spector/ui/g"
)

const (
	SplitterTabSize = 20
	SplitterRadius  = 1
)

var (
	HaloColor = g.Color{0x00, 0x00, 0x00, 0xFF}
)

type Area struct {
	Screen *Screen
	Parent *Area  // nil for Root
	Bounds g.Rect // last-bounds

	// Either Split or Editor
	Vertical  bool
	Splitters []*Splitter
	Editor    *Editor
}

func NewArea(screen *Screen) *Area {
	return &Area{
		Screen: screen,
		Parent: nil,
	}
}

func (area *Area) Clone() *Area {
	clone := &Area{}
	clone.Screen = area.Screen
	clone.Parent = area.Parent
	clone.Bounds = area.Bounds

	//TODO: what to do with nested splitters?
	clone.Vertical = area.Vertical
	clone.Editor = area.Editor.Clone()
	return clone
}

func (area *Area) Update(ctx *ui.Context) {
	area.Bounds = ctx.Area

	if len(area.Splitters) > 0 {
		for _, splitter := range area.Splitters {
			splitter.Update(ctx)
		}
	} else {
		area.Editor.Update(ctx)
	}
}

func (area *Area) RelativeToAbsolute(v float32) float32 {
	return g.LerpClamp(v, area.Bounds.Min.X, area.Bounds.Max.X)
}

func (area *Area) AbsoluteToRelative(v float32) float32 {
	return g.InverseLerpClamp(v, area.Bounds.Min.X, area.Bounds.Max.X)
}

type Splitter struct {
	Owner   *Area
	Content *Area
	Index   int

	RelativeCenter float32

	Active    bool
	Resizing  bool
	Splitting bool
}

func NewSplitter(owner *Area) *Splitter {
	return &Splitter{Owner: owner}
}

func (splitter *Splitter) NeighborCenters() (left, right float32) {
	rmin, rmax := float32(0), float32(1)
	if 0 <= splitter.Index-1 {
		rmin = splitter.Owner.Splitters[splitter.Index-1].RelativeCenter
	}
	if splitter.Index+1 < len(splitter.Owner.Splitters) {
		rmax = splitter.Owner.Splitters[splitter.Index+1].RelativeCenter
	}
	return splitter.Owner.RelativeToAbsolute(rmin), splitter.Owner.RelativeToAbsolute(rmax)
}

func (splitter *Splitter) MinMax() (float32, float32) {
	min, max := splitter.NeighborCenters()
	min += EditorMinSize
	max -= EditorMinSize
	center := splitter.Center()
	if min > center {
		min = center
	}
	if max < center {
		max = center
	}
	return min, max
}

func (splitter *Splitter) Center() float32 {
	return splitter.Owner.RelativeToAbsolute(splitter.RelativeCenter)
}

func (splitter *Splitter) SetCenter(absolute float32) {
	splitter.RelativeCenter = splitter.Owner.AbsoluteToRelative(absolute)
}

func (splitter *Splitter) ContentArea() g.Rect {
	left, _ := splitter.NeighborCenters()
	r := splitter.Owner.Bounds
	r.Min.X = left
	if splitter.Index > 0 {
		r.Min.X += SplitterRadius
	}
	r.Max.X = splitter.Center() - SplitterRadius
	return r
}

func (splitter *Splitter) TabRect() g.Rect {
	center := splitter.Center()
	r := splitter.Owner.Bounds
	r.Min.X = center - SplitterTabSize
	r.Max.X = center
	r.Max.Y = r.Min.Y + SplitterTabSize
	return r
}

func (splitter *Splitter) isLast() bool {
	return len(splitter.Owner.Splitters) == splitter.Index+1
}

func (splitter *Splitter) Rect() g.Rect {
	return splitter.Owner.Bounds.VerticalLine(splitter.Center(), SplitterRadius)
}

func (splitter *Splitter) Split() *Splitter {
	other := &Splitter{}
	other.Owner = splitter.Owner
	other.Content = splitter.Content.Clone()
	other.RelativeCenter = splitter.RelativeCenter

	// TODO: move splitter insertion somewhere else
	pivot := splitter.Index + 1
	prev := splitter.Owner.Splitters
	next := []*Splitter{}
	next = append(next, prev[:pivot]...)
	next = append(next, other)
	next = append(next, prev[pivot:]...)
	for i, splitter := range next {
		splitter.Index = i
	}
	splitter.Owner.Splitters = next

	splitter.RelativeCenter = splitter.Owner.AbsoluteToRelative(other.Center() - EditorMinSize)

	return other
}

func (splitter *Splitter) Update(ctx *ui.Context) {
	if splitter.Content != nil {
		content := splitter.ContentArea()
		splitter.Content.Update(ctx.Child(content))
	}

	{ // split tab drawing
		r := splitter.TabRect()
		inside := r.Contains(ctx.Input.Mouse.Pos) && ctx.Input.Mouse.Capture == nil
		if inside && ctx.Input.Mouse.Pressed {
			splitter.Splitting = true
			ctx.Input.Mouse.Capture = func() bool {
				center := splitter.Center()

				cansplit := splitter.ContentArea().Dx() > 2*EditorMinSize
				canmerge := !splitter.isLast()
				if !canmerge && !cansplit {
					splitter.Splitting = false
					return true
				}

				distance := ctx.Input.Mouse.Pos.X - center
				halo := splitter.Rect() // TODO: optimize
				if (distance < 0) && cansplit {
					halo.Min.X = ctx.Input.Mouse.Pos.X
					halo.Max.X = center

					alpha := g.Sat8(g.Abs(distance) / EditorMinSize)
					ctx.Hover.FillRect(&halo, g.Color{0xFF, 0xFF, 0xFF, alpha})

					if distance < -EditorMinSize {
						// split
						splitter.Split()
						splitter.Resizing = true
						splitter.Splitting = false

						ctx.Input.Mouse.Capture = func() bool {
							min, max := splitter.MinMax()
							splitter.SetCenter(g.Clamp(ctx.Input.Mouse.Pos.X, min, max))
							splitter.Resizing = ctx.Input.Mouse.Down
							return !ctx.Input.Mouse.Down
						}
						return !ctx.Input.Mouse.Down
					}
				} else if (distance > 0) && canmerge {
					halo.Min.X = center
					halo.Max.X = ctx.Input.Mouse.Pos.X

					alpha := g.Sat8(g.Abs(distance) / EditorMinSize)
					ctx.Hover.FillRect(&halo, g.Color{0xFF, 0xFF, 0xFF, alpha})

					if distance > EditorMinSize {
						// merge
						splitter.Splitting = false
						return true
					}
				}

				splitter.Splitting = ctx.Input.Mouse.Down
				return !ctx.Input.Mouse.Down
			}
		}

		if splitter.Splitting {
			ctx.Input.Mouse.Cursor = ui.CrosshairCursor
			distance := splitter.Center() - ctx.Input.Mouse.Pos.X
			r = r.Add(g.Vector{-distance, 0})
			ctx.Draw.FillRect(&r, g.Green)
		} else if inside {
			ctx.Input.Mouse.Cursor = ui.CrosshairCursor
			ctx.Draw.FillRect(&r, g.Red)
		} else {
			ctx.Draw.FillRect(&r, g.Color{0x80, 0x80, 0x80, 0xff})
		}
	}

	{ // separator drawing / moving
		r := splitter.Rect()
		inside := false
		if !splitter.isLast() {
			inside = r.Contains(ctx.Input.Mouse.Pos) && ctx.Input.Mouse.Capture == nil
			if inside && ctx.Input.Mouse.Pressed {
				splitter.Resizing = true
				ctx.Input.Mouse.Capture = func() bool {
					min, max := splitter.MinMax()
					splitter.SetCenter(g.Clamp(ctx.Input.Mouse.Pos.X, min, max))
					splitter.Resizing = ctx.Input.Mouse.Down
					return !ctx.Input.Mouse.Down
				}
			}
		}

		if splitter.Resizing {
			ctx.Input.Mouse.Cursor = ui.HResizeCursor
			ctx.Draw.FillRect(&r, g.Green)
		} else if inside {
			ctx.Input.Mouse.Cursor = ui.HResizeCursor
			ctx.Draw.FillRect(&r, g.Red)
		} else {
			ctx.Draw.FillRect(&r, g.Color{0x80, 0x80, 0x80, 0xff})
		}
	}
}
