package ui

type OnDrag func(delta float32)

func DragX(ctx *Context, ondrag OnDrag) {
	color := ButtonColor.Default
	down := false
	if ctx.Input.Mouse.PointsAt(ctx.Area) {
		color = ButtonColor.Hot
		if ctx.Input.Mouse.Down  {
			color = ButtonColor.Active
			down = true
		}
	}

	ctx.Backend.SetBack(color)
	ctx.Backend.Fill(ctx.Area)
	ctx.Backend.Stroke(ctx.Area)

	if down && !ctx.Input.Mouse.Last.Down {
		DoDragX(ctx, ondrag)
	}
}

func DoDragX(ctx *Context, ondrag OnDrag) {
	ctx.Input.Mouse.Drag = func(ctx *Context) bool {
		delta := ctx.Input.Mouse.Last.Position.X - ctx.Input.Mouse.Position.X
		if delta != 0 {
			ondrag(delta)
		}
		return ctx.Input.Mouse.Down
	}
}

func DoDragY(ctx *Context, ondrag OnDrag) {
	ctx.Input.Mouse.Drag = func(ctx *Context) bool {
		delta := ctx.Input.Mouse.Last.Position.Y - ctx.Input.Mouse.Position.Y
		if delta != 0 {
			ondrag(delta)
		}
		return ctx.Input.Mouse.Down
	}
}