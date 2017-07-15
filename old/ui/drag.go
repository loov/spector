package ui

type OnDrag func(delta float32)
type OnDragXY func(dx, dy float32)

func SplitterX(ctx *Context, ondrag OnDrag) {
	color := ButtonColor.Default
	if ctx.Input.Mouse.PointsAt(ctx.Area) {
		color = ButtonColor.Hot
		if ctx.Input.Mouse.Down  {
			color = ButtonColor.Active
		}
	}

	ctx.Backend.SetBack(color)
	ctx.Backend.Fill(ctx.Area)
	ctx.Backend.Stroke(ctx.Area)

	DoDragX(ctx, ondrag)
}

func MouseJustPressed(ctx *Context) bool{
	return ctx.Input.Mouse.PointsAt(ctx.Area) && ctx.Input.Mouse.Down && !ctx.Input.Mouse.Last.Down
}

func DoDragXY(ctx *Context, ondrag OnDragXY) {
	if MouseJustPressed(ctx) {
		ctx.Input.Mouse.Drag = func(ctx *Context) bool {
			dx := ctx.Input.Mouse.Position.X - ctx.Input.Mouse.Last.Position.X
			dy := ctx.Input.Mouse.Position.Y - ctx.Input.Mouse.Last.Position.Y
			if dx != 0 || dy != 0{
				ondrag(dx, dy)
			}
			return ctx.Input.Mouse.Down
		}
	}
}

func DoDragX(ctx *Context, ondrag OnDrag) {
	if MouseJustPressed(ctx) {
		ctx.Input.Mouse.Drag = func(ctx *Context) bool {
			delta := ctx.Input.Mouse.Position.X - ctx.Input.Mouse.Last.Position.X
			if delta != 0 {
				ondrag(delta)
			}
			return ctx.Input.Mouse.Down
		}
	}
}

func DoDragY(ctx *Context, ondrag OnDrag) {
	if MouseJustPressed(ctx) {
		ctx.Input.Mouse.Drag = func(ctx *Context) bool {
			delta := ctx.Input.Mouse.Position.Y - ctx.Input.Mouse.Last.Position.Y
			if delta != 0 {
				ondrag(delta)
			}
			return ctx.Input.Mouse.Down
		}
	}
}