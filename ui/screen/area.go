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

	next := resize.Target.Clone()
	resize.Screen.Areas = append(resize.Screen.Areas, next)

	resize.Y = []*float32{
		&act.Target.RelBounds.Min.Y,
		&next.RelBounds.Max.Y,
	}

	ctx.Input.Mouse.Capture = resize.Update
}

func (act *JoinSplit) SplitHorizontal(ctx *ui.Context) {
	resize := &Resize{}
	resize.Screen = act.Screen
	resize.Target = act.Target
	resize.Init(ctx)
	resize.X, resize.Y = nil, nil

	next := resize.Target.Clone()
	resize.Screen.Areas = append(resize.Screen.Areas, next)

	resize.X = []*float32{
		&act.Target.RelBounds.Max.X,
		&next.RelBounds.Min.X,
	}

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
	x, y := FindLinked(act.Screen, p, act.Screen.Bounds.ToRelative(p))

	act.X = y.Ps
	act.Y = x.Ps
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
