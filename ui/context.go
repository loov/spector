package ui

type Widget func(ctx *Context)

type Backend interface {
	Clone() Backend

	SetFore(Color)
	SetBack(Color)

	Fill(Bounds)
	Stroke(Bounds)

	SetFont(name string, size float32)
	SetFontColor(Color)
	Text(text string, bounds Bounds)
	Measure(text string) (size Point)
}

type Context struct {
	Input *Input
	Backend
	Area Bounds

	ID    string
	Index int
	Count int
}

func (ctx *Context) Child(area Bounds) *Context {
	ctx.Count++
	return &Context{
		Input:   ctx.Input,
		Backend: ctx.Backend.Clone(),
		Area:    area,
		Index:   ctx.Count - 1,
		Count:   0,
	}
}

func (ctx *Context) Left(w float32) *Context {
	inner := ctx.Area
	inner.Max.X = inner.Min.X + w
	ctx.Area.Min.X += w
	return ctx.Child(inner)
}

func (ctx *Context) Right(w float32) *Context {
	inner := ctx.Area
	inner.Min.X = inner.Max.X - w
	ctx.Area.Max.X -= w
	return ctx.Child(inner)
}

func (ctx *Context) Top(h float32) *Context {
	inner := ctx.Area
	inner.Max.Y = inner.Min.Y + h
	ctx.Area.Min.Y += h
	return ctx.Child(inner)
}

func (ctx *Context) Bottom(h float32) *Context {
	inner := ctx.Area
	inner.Min.Y = inner.Max.Y - h
	ctx.Area.Max.Y -= h
	return ctx.Child(inner)
}

func (ctx *Context) Fill() *Context {
	return ctx.Child(ctx.Area)
}
