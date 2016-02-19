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

	if text != "" {
		ctx.Backend.Text(text, ctx.Area)
	}

	return
}

func (ctx *Context) Text(text string) { Text(ctx, text) }

func Text(ctx *Context, text string) {
	if text != "" {
		ctx.Backend.Text(text, ctx.Area)
	}
}

type DynamicLayouter interface {
	LayoutContext() *Context
	NextSized(size Point) *Context
}

type Layouter interface {
	LayoutContext() *Context
	Next() *Context
}

type LayoutType byte

const (
	ToBottom = LayoutType(iota)
	ToRight
	ToLeft
	ToTop
)

type Layout struct {
	Height float32
	Width  float32
	Type   LayoutType

	*Context
}

func LayoutToRight(w float32, ctx *Context) *Layout {
	return &Layout{
		Width:   w,
		Type:    ToRight,
		Context: ctx,
	}
}

func LayoutToBottom(h float32, ctx *Context) *Layout {
	return &Layout{
		Height:  h,
		Type:    ToBottom,
		Context: ctx,
	}
}

func (layout *Layout) LayoutContext() *Context { return layout.Context }

func (layout *Layout) NextSized(size Point) *Context {
	switch layout.Type {
	case ToBottom:
		return layout.Context.Top(size.Y)
	case ToTop:
		return layout.Context.Bottom(size.Y)
	case ToRight:
		return layout.Context.Left(size.X)
	case ToLeft:
		return layout.Context.Right(size.X)
	}
	panic("shouldn't happen")
}

func (layout *Layout) Next() *Context {
	switch layout.Type {
	case ToBottom:
		return layout.Context.Top(layout.Height)
	case ToTop:
		return layout.Context.Bottom(layout.Height)
	case ToRight:
		return layout.Context.Left(layout.Width)
	case ToLeft:
		return layout.Context.Right(layout.Width)
	}
	panic("shouldn't happen")
}

type Callback struct {
	Name string
	Func func()
}

type Buttons []Callback

func (list Buttons) DoDynamic(layout DynamicLayouter) {
	context := layout.LayoutContext()
	for _, button := range list {
		size := context.Measure(button.Name)
		if layout.NextSized(size).Button(button.Name) {
			if button.Func != nil {
				button.Func()
			}
		}
	}
}

func (list Buttons) Do(layout Layouter) {
	for _, button := range list {
		if layout.Next().Button(button.Name) {
			if button.Func != nil {
				button.Func()
			}
		}
	}
}
