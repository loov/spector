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
	Screen    *Screen
	RelBounds g.Rect
	Bounds    g.Rect // last-bounds
	Editor    *Editor
}

func NewArea(screen *Screen) *Area {
	return &Area{
		Screen:    screen,
		RelBounds: g.Rect{g.V0, g.V1},
	}
}

func (area *Area) Clone() *Area {
	clone := &Area{}
	clone.Screen = area.Screen
	clone.RelBounds = area.RelBounds
	clone.Bounds = area.Bounds
	clone.Editor = area.Editor.Clone()

	return clone
}

func (area *Area) Update(ctx *ui.Context) {
	area.Bounds = ctx.Area
	area.Editor.Update(ctx)

	{
		r := area.JoinSplitRect()
		canCapture := ctx.Input.Mouse.Capture == nil && r.Contains(ctx.Input.Mouse.Pos)
		if canCapture && ctx.Input.Mouse.Pressed {
			split := &JoinSplit{Target: area}
			split.Init(ctx)
			ctx.Input.Mouse.Capture = split.Update
		}

		if canCapture {
			ctx.Input.Mouse.Cursor = ui.CrosshairCursor
			ctx.Draw.FillRect(&r, BorderHighlightColor)
		} else {
			ctx.Draw.FillRect(&r, BorderColor)
		}
	}

	{
		r := area.Bounds
		p := ctx.Input.Mouse.Pos

		canCapture := ctx.Input.Mouse.Capture == nil

		nearLeft := g.Abs(r.Min.X-p.X) <= AreaBorderRadius
		nearRight := g.Abs(r.Max.X-p.X) <= AreaBorderRadius

		nearTop := g.Abs(r.Min.Y-p.Y) <= AreaBorderRadius
		nearBottom := g.Abs(r.Max.Y-p.Y) <= AreaBorderRadius

		nearEdge := nearLeft || nearRight || nearTop || nearBottom

		if canCapture && nearEdge && ctx.Input.Mouse.Pressed {
			resize := &Resize{Target: area}
			resize.Init(ctx)
			ctx.Input.Mouse.Capture = resize.Update
		}

		if (nearLeft || nearRight) && (nearTop || nearBottom) {
			ctx.Input.Mouse.Cursor = ui.CrosshairCursor
			ctx.Draw.StrokeRect(&area.Bounds, 1, BorderHighlightColor)
		} else if nearLeft || nearRight {
			ctx.Input.Mouse.Cursor = ui.HResizeCursor
			ctx.Draw.StrokeRect(&area.Bounds, 1, BorderHighlightColor)
		} else if nearTop || nearBottom {
			ctx.Input.Mouse.Cursor = ui.VResizeCursor
			ctx.Draw.StrokeRect(&area.Bounds, 1, BorderHighlightColor)
		} else {
			ctx.Draw.StrokeRect(&area.Bounds, 1, BorderColor)
		}

	}
}
func (area *Area) JoinSplitRect() g.Rect {
	r := area.Bounds
	r.Min.X = r.Max.X - JoinSplitSize
	r.Max.Y = r.Min.Y + JoinSplitSize
	return r
}

type JoinSplit struct {
	Screen *Screen
	Target *Area
}

func (act *JoinSplit) Init(ctx *ui.Context) {
	act.Screen = act.Target.Screen
}

func (act *JoinSplit) Update(ctx *ui.Context) bool {
	if !ctx.Input.Mouse.Down {
		return true
	}

	r := act.Target.Bounds
	p := ctx.Input.Mouse.Pos
	d := p.Sub(r.TopRight())
	if d.X < 0 && d.Y > 0 {
		r = r.Deflate(g.Vector{AreaBorderRadius / 2, AreaBorderRadius / 2})
		if -d.X > d.Y {
			r.Min.X = p.X
			dist := -d.X
			alpha := g.Sat8(g.Abs(dist) / EditorMinSize)
			ctx.Hover.FillRect(&r, AreaBackground.WithAlpha(alpha))
			if dist > EditorMinSize {
				act.SplitHorizontal(ctx)
			}
		} else {
			r.Max.Y = p.Y
			dist := d.Y
			alpha := g.Sat8(g.Abs(dist) / EditorMinSize)
			ctx.Hover.FillRect(&r, AreaBackground.WithAlpha(alpha))
			if dist > EditorMinSize {
				act.SplitVertical(ctx)
			}
		}

		// TODO: draw plus sign in center
	}

	ctx.Input.Mouse.Cursor = ui.HandCursor
	return false
}

func (act *JoinSplit) SplitVertical(ctx *ui.Context) {
	resize := &Resize{}
	resize.Screen = act.Screen
	resize.Target = act.Target
	resize.Init(ctx)
	resize.X, resize.Y = nil, nil

	resize.Y = []*float32{
		&act.Target.RelBounds.Min.Y,
	}
	//TODO

	ctx.Input.Mouse.Capture = resize.Update
}

func (act *JoinSplit) SplitHorizontal(ctx *ui.Context) {
	resize := &Resize{}
	resize.Screen = act.Screen
	resize.Target = act.Target
	resize.Init(ctx)
	resize.X, resize.Y = nil, nil

	resize.X = []*float32{
		&act.Target.RelBounds.Max.X,
	}
	//TODO

	ctx.Input.Mouse.Capture = resize.Update
}

type Resize struct {
	Screen *Screen
	Target *Area

	X []*float32
	Y []*float32
}

func (act *Resize) Init(ctx *ui.Context) {
	act.Screen = act.Target.Screen

	p := ctx.Input.Mouse.Pos
	for _, area := range act.Screen.Areas {
		if g.Abs(area.Bounds.Min.X-p.X) <= AreaBorderRadius {
			act.X = append(act.X, &area.RelBounds.Min.X)
		}
		if g.Abs(area.Bounds.Max.X-p.X) <= AreaBorderRadius {
			act.X = append(act.X, &area.RelBounds.Max.X)
		}

		if g.Abs(area.Bounds.Min.Y-p.Y) <= AreaBorderRadius {
			act.Y = append(act.Y, &area.RelBounds.Min.Y)
		}
		if g.Abs(area.Bounds.Max.Y-p.Y) <= AreaBorderRadius {
			act.Y = append(act.Y, &area.RelBounds.Max.Y)
		}
	}
}

func (act *Resize) Update(ctx *ui.Context) bool {
	if !ctx.Input.Mouse.Down {
		return true
	}

	if len(act.X) > 0 && len(act.Y) > 0 {
		ctx.Input.Mouse.Cursor = ui.CrosshairCursor
	} else if len(act.X) > 0 {
		ctx.Input.Mouse.Cursor = ui.HResizeCursor
	} else if len(act.Y) > 0 {
		ctx.Input.Mouse.Cursor = ui.VResizeCursor
	} else {
		return true
	}

	rp := act.Screen.Bounds.ToRelative(ctx.Input.Mouse.Pos)
	for _, px := range act.X {
		*px = rp.X
	}
	for _, py := range act.Y {
		*py = rp.Y
	}

	return false
}

/*
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
	r.Min.X = center - JoinSplitSize
	r.Max.X = center
	r.Max.Y = r.Min.Y + JoinSplitSize
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
*/
