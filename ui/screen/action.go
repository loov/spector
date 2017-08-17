package screen

import (
	"github.com/egonelbre/spector/ui"
	"github.com/egonelbre/spector/ui/g"
)

type Joiner struct {
	Screen *Screen
	Corner *Corner

	inited    bool
	splitArea g.Rect
	limit     g.Rect

	canMergeTop   bool
	canMergeRight bool
	canFullscreen bool
}

func (act *Joiner) init(ctx *ui.Context) {
	left, _ := act.Corner.BlockingHorizontal(false, true)
	_, bottom := act.Corner.BlockingVertical(true, false)

	act.splitArea = g.Rect{
		Min: left.Pos,
		Max: bottom.Pos,
	}

	// TODO: can *
}

func (act *Joiner) Capture(ctx *ui.Context) bool {
	if !ctx.Input.Mouse.Down {
		return true
	}
	if !act.inited {
		act.inited = true
		act.init(ctx)
	}

	cornerPos := act.Screen.Bounds.ToGlobal(act.Corner.Pos)
	delta := ctx.Input.Mouse.Pos.Sub(cornerPos)

	if delta.Len() < RigJoinerSize {
		// for canceling
		return false
	}

	// splitting action
	if delta.X < 0 && delta.Y > 0 {
		return act.trySplit(ctx, delta)
	}

	// merge to right
	if 0 < delta.X && 0 < delta.Y && act.canMergeRight {
		return false
	}

	// merge to top
	if delta.X < 0 && delta.Y < 0 && act.canMergeTop {
		return false
	}

	// fullscreen on release
	if 0 < delta.X && delta.Y < 0 && act.canFullscreen {
		return false
	}

	return false
}

func (act *Joiner) trySplit(ctx *ui.Context, delta g.Vector) bool {
	dx := -delta.X
	dy := delta.Y

	r := act.Screen.Bounds.Subset(act.splitArea)
	ctx.Hover.FillRect(&r, g.Color{0xFF, 0, 0, 0x20})

	if dx < dy {
		// do we have enough room for two areas?
		if r.Dy() < 2*RigTriggerSize {
			return false
		}
		ctx.Input.Mouse.Cursor = ui.CrosshairCursor

		r.Max.Y = r.Max.Y - RigTriggerSize
		if r.Min.Y+dy < r.Max.Y {
			r.Max.Y = r.Min.Y + dy
		}

		alpha := g.Sat8((dy - RigJoinerSize) / (RigTriggerSize - RigJoinerSize))
		ctx.Hover.FillRect(&r, RigBackground.WithAlpha(alpha))

		if dy > RigTriggerSize {
			// TODO: fix don't use mouse pos, it might be outside of limits
			rp := act.Screen.Bounds.ToRelative(ctx.Input.Mouse.Pos)
			split := act.Screen.Rig.SplitHorizontally(act.Corner, rp.Y)
			ctx.Input.Mouse.Capture = (&Resizer{
				Screen:     act.Screen,
				Start:      split.Center(),
				Horizontal: split,
				Vertical:   nil,
			}).Capture
		}
	} else {
		// do we have enough room for two areas?
		if r.Dx() < 2*RigTriggerSize {
			return false
		}
		ctx.Input.Mouse.Cursor = ui.CrosshairCursor

		r.Min.X = r.Min.X + RigTriggerSize
		if r.Min.X < r.Max.X-dx {
			r.Min.X = r.Max.X - dx
		}

		alpha := g.Sat8((dx - RigJoinerSize) / (RigTriggerSize - RigJoinerSize))
		ctx.Hover.FillRect(&r, RigBackground.WithAlpha(alpha))

		if dx > RigTriggerSize {
			// TODO: fix don't use mouse pos, it might be outside of limits
			rp := act.Screen.Bounds.ToRelative(ctx.Input.Mouse.Pos)
			split := act.Screen.Rig.SplitVertically(act.Corner, rp.X)
			ctx.Input.Mouse.Capture = (&Resizer{
				Screen:     act.Screen,
				Start:      split.Center(),
				Horizontal: nil,
				Vertical:   split,
			}).Capture
		}
	}

	return false
}

type Resizer struct {
	Screen     *Screen
	Start      g.Vector
	Horizontal *Border
	Vertical   *Border

	inited bool
	area   g.Rect
}

func (act *Resizer) init(ctx *ui.Context) {
	act.area = g.Rect01
	if act.Vertical != nil {
		for _, corner := range act.Vertical.Corners {
			checkTop := act.Vertical.First() != corner
			checkBottom := act.Vertical.Last() != corner

			left, right := corner.BlockingHorizontal(checkTop, checkBottom)
			if left != nil && act.area.Min.X < left.Pos.X {
				act.area.Min.X = left.Pos.X
			}
			if right != nil && right.Pos.X < act.area.Max.X {
				act.area.Max.X = right.Pos.X
			}
		}
	}
	if act.Horizontal != nil {
		for _, corner := range act.Horizontal.Corners {
			checkLeft := act.Horizontal.First() != corner
			checkRight := act.Horizontal.Last() != corner

			top, bottom := corner.BlockingVertical(checkLeft, checkRight)
			if top != nil && act.area.Min.Y < top.Pos.Y {
				act.area.Min.Y = top.Pos.Y
			}
			if bottom != nil && bottom.Pos.Y < act.area.Max.Y {
				act.area.Max.Y = bottom.Pos.Y
			}
		}
	}
}

func (act *Resizer) Capture(ctx *ui.Context) bool {
	if !ctx.Input.Mouse.Down {
		for _, border := range act.Screen.Rig.Borders {
			border.Sort()
		}
		return true
	}

	if !act.inited {
		act.inited = true
		act.init(ctx)
	}

	limiter := act.Screen.Bounds.Subset(act.area).Deflate(RigTriggerRadius)
	if limiter.Max.X < limiter.Min.X || limiter.Max.Y < limiter.Min.Y {
		return true
	}

	ctx.Hover.FillRect(&limiter, g.Color{0xFF, 0, 0, 0x20})

	clampedMouse := limiter.ClosestPoint(ctx.Input.Mouse.Pos)
	p := act.Screen.Bounds.ToRelative(clampedMouse)

	if act.Horizontal != nil {
		for _, corner := range act.Horizontal.Corners {
			corner.Pos.Y = p.Y
		}
	}

	if act.Vertical != nil {
		for _, corner := range act.Vertical.Corners {
			corner.Pos.X = p.X
		}
	}

	return false
}
