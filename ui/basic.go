package ui

func (ctx *Context) Panel() *Context {
	ctx.Backend.Fill(ctx.Area)
	ctx.Backend.Stroke(ctx.Area)
	return ctx
}

func (ctx *Context) Button(text string) (pressed bool) { return Button(ctx, text) }

func Button(ctx *Context, text string) (pressed bool) {
	color := ButtonColor.Default
	if ctx.Input.Mouse.PointsAt(ctx.Area) {
		color = ButtonColor.Hot
		if ctx.Input.Mouse.Down {
			color = ButtonColor.Active
		} else if ctx.Input.Mouse.Clicked() {
			color = ButtonColor.Clicked
			pressed = true
		}
	}

	ctx.Backend.SetBack(color)
	ctx.Backend.Fill(ctx.Area)
	ctx.Backend.Stroke(ctx.Area)
	ctx.Backend.Text(text, ctx.Area)

	return
}

type Items interface {
	Len() int
	Measure(ctx *Context, i int) Point
	Render(ctx *Context, i int)
}

type Layout struct{ Context *Context }

func (lay Layout) Panel() Layout {
	lay.Context.Panel()
	return lay
}

func (lay Layout) Left(items Items) {
	for i, n := 0, items.Len(); i < n; i++ {
		size := items.Measure(lay.Context, i)
		items.Render(lay.Context.Left(size.X), i)
	}
}

func (lay Layout) Top(items Items) {
	for i, n := 0, items.Len(); i < n; i++ {
		size := items.Measure(lay.Context, i)
		items.Render(lay.Context.Top(size.Y), i)
	}
}

type Callback struct {
	Name string
	Func func()
}

type Buttons []Callback

type FixedHeight struct {
	Height float32
	Items
}

func (list FixedHeight) Measure(ctx *Context, i int) (size Point) {
	size = list.Items.Measure(ctx, i)
	size.Y = list.Height
	return
}

func (list Buttons) Len() int                                 { return len(list) }
func (list Buttons) Measure(ctx *Context, i int) (size Point) { return ctx.Measure(list[i].Name) }
func (list Buttons) Render(ctx *Context, i int) {
	if ctx.Button(list[i].Name) {
		if list[i].Func != nil {
			list[i].Func()
		}
	}
}
