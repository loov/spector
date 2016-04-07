package ui

func DragX(ctx *Context, value *float32) {
	color := ButtonColor.Default
	pressed := false
	if ctx.Input.Mouse.PointsAt(ctx.Area) {
		color = ButtonColor.Hot
		if ctx.Input.Mouse.Down {
			color = ButtonColor.Active
			pressed = true
		}
	}

	ctx.Backend.SetBack(color)
	ctx.Backend.Fill(ctx.Area)
	ctx.Backend.Stroke(ctx.Area)

	if pressed {
		ctx.Input.NextID = ctx.ID
	}

	if ctx.Input.ActiveID == ctx.ID && ctx.Input.Mouse.Down{
		*value += ctx.Input.Mouse.Last.Position.X - ctx.Input.Mouse.Position.X
		ctx.Input.NextID = ctx.ID
	}
}
