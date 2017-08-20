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
	Rig      *Rig
	Areas    []*Area
}

func New() *Screen {
	screen := &Screen{}
	screen.Registry = NewRegistry()

	screen.Rig = NewRig()
	return screen
}

func (screen *Screen) Update(ctx *ui.Context) {
	screen.Bounds = ctx.Area

	for _, corner := range screen.Rig.Corners {
		p := ctx.Area.ToGlobal(corner.Pos)

		r := p.Inflate(RigCornerRadius)
		canCapture := !corner.IsLocked() && ctx.Input.Mouse.Capture == nil && r.Contains(ctx.Input.Mouse.Pos)
		if canCapture {
			ctx.Input.Mouse.SetCaptureCursor(ui.CrosshairCursor)
			if ctx.Input.Mouse.Pressed {
				ctx.Input.Mouse.Capture = (&Resizer{
					Screen:     screen,
					Start:      corner.Pos,
					Horizontal: corner.Horizontal,
					Vertical:   corner.Vertical,
				}).Capture
			}
		}

		if corner.SideLeft() != nil && corner.SideBottom() != nil {
			r := g.Rect{
				Min: g.Vector{
					X: p.X - RigJoinerSize,
					Y: p.Y - RigBorderRadius.Y,
				},
				Max: g.Vector{
					X: p.X - RigBorderRadius.X,
					Y: p.Y + RigJoinerSize,
				},
			}

			canCapture := ctx.Input.Mouse.Capture == nil && r.Contains(ctx.Input.Mouse.Pos)
			if canCapture {
				ctx.Input.Mouse.SetCaptureCursor(ui.CrosshairCursor)
				if !ctx.Input.Mouse.Pressed {
					ctx.Draw.FillRect(&r, RigJoinerHighlightColor)
				} else {
					ctx.Input.Mouse.Capture = (&Joiner{
						Screen: screen,
						Corner: corner,
					}).Capture
				}
			} else {
				ctx.Draw.FillRect(&r, RigJoinerColor)
			}
		}
	}

	for _, border := range screen.Rig.Borders {
		min := ctx.Area.ToGlobal(border.Min())
		max := ctx.Area.ToGlobal(border.Max())

		r := g.Rect{min, max}.Inflate(RigBorderRadius)

		canCapture := !border.Locked && ctx.Input.Mouse.Capture == nil && r.Contains(ctx.Input.Mouse.Pos)
		if canCapture {
			var horz, vert *Border
			if border.Horizontal {
				ctx.Input.Mouse.SetCaptureCursor(ui.VResizeCursor)
				horz = border
			} else {
				ctx.Input.Mouse.SetCaptureCursor(ui.HResizeCursor)
				vert = border
			}

			ctx.Draw.FillRect(&r, RigBorderHighlightColor)
			if ctx.Input.Mouse.Pressed {
				ctx.Input.Mouse.Capture = (&Resizer{
					Screen:     screen,
					Start:      border.Center(),
					Horizontal: horz,
					Vertical:   vert,
				}).Capture
			}
		} else {
			ctx.Draw.FillRect(&r, RigBorderColor)
		}
	}
}

func NewTestArea(screen *Screen) *Area {
	child := NewArea(screen)
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
