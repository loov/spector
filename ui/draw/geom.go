package draw

import "math"

var (
	nan32 = math.Float32frombits(0x7FBFFFFF)
	noUV  = Vector{nan32, nan32}
)

type Vector struct {
	X, Y float32
}

func (a Vector) IsInvalid() bool { return a.X != a.X }

func (a Vector) Add(b Vector) Vector    { return Vector{a.X + b.X, a.Y + b.Y} }
func (a Vector) Sub(b Vector) Vector    { return Vector{a.X - b.X, a.Y - b.Y} }
func (a Vector) Scale(s float32) Vector { return Vector{a.X * s, a.Y * s} }
func (a Vector) Div(s float32) Vector   { return Vector{a.X / s, a.Y / s} }

func (a Vector) ScaleTo(size float32) Vector {
	ilen := a.Len()
	if ilen > 0 {
		ilen = size / ilen
	}
	return a.Scale(ilen)
}

func (a Vector) Len() float32  { return sqrt(a.X*a.X + a.Y*a.Y) }
func (a Vector) Len2() float32 { return a.X*a.X + a.Y*a.Y }

func (a Vector) Min(b Vector) Vector { return Vector{min(a.X, b.X), min(a.Y, b.Y)} }
func (a Vector) Max(b Vector) Vector { return Vector{max(a.X, b.X), max(a.Y, b.Y)} }

type Rectangle struct {
	Min, Max Vector
}

func Rect(x, y, w, h float32) Rectangle {
	return Rectangle{Min: Vector{x, y}, Max: Vector{x + w, y + h}}
}

func (r Rectangle) IsInvalid() bool { return r.Min.X != r.Min.X }

func (r *Rectangle) Corners() (tl, tr, br, bl Vector) {
	tl = r.TopLeft()
	tr = r.TopRight()
	br = r.BottomRight()
	bl = r.BottomLeft()
	return
}

func (r *Rectangle) TopLeft() Vector     { return r.Min }
func (r *Rectangle) TopRight() Vector    { return Vector{r.Max.X, r.Min.Y} }
func (r *Rectangle) BottomLeft() Vector  { return Vector{r.Min.X, r.Max.Y} }
func (r *Rectangle) BottomRight() Vector { return r.Max }

func (r *Rectangle) Contains(p Vector) bool {
	return r.Min.X <= p.X && p.X < r.Max.X &&
		r.Min.Y <= p.Y && p.Y < r.Max.Y
}

func (a *Rectangle) Overlaps(b *Rectangle) bool {
	return b.Min.Y < a.Max.Y && b.Max.Y > a.Min.Y &&
		b.Min.X < a.Max.X && b.Max.X > a.Min.X
}

func (r *Rectangle) Offset(v Vector) {
	r.Min.X += v.X
	r.Min.Y += v.Y
	r.Max.X += v.X
	r.Max.Y += v.Y
}

func (r *Rectangle) ExpandFloat(amount float32) {
	r.Min.X -= amount
	r.Min.Y -= amount
	r.Max.X += amount
	r.Max.Y += amount
}

func (r *Rectangle) Expand(v Vector) {
	r.Min.X -= v.X
	r.Min.Y -= v.Y
	r.Max.X += v.X
	r.Max.Y += v.Y
}

func (r *Rectangle) Reduce(v Vector) {
	r.Min.X += v.X
	r.Min.Y += v.Y
	r.Max.X -= v.X
	r.Max.Y -= v.Y
}

func (r *Rectangle) Clip(clip *Rectangle) {
	if r.Min.X < clip.Min.X {
		r.Min.X = clip.Min.X
	}
	if r.Min.Y < clip.Min.Y {
		r.Min.Y = clip.Min.Y
	}
	if r.Max.X > clip.Max.X {
		r.Max.X = clip.Max.X
	}
	if r.Max.Y > clip.Max.Y {
		r.Max.Y = clip.Max.Y
	}
}

func (r *Rectangle) Floor() {
	r.Min.X = (float32)((int)(r.Min.X))
	r.Min.Y = (float32)((int)(r.Min.Y))
	r.Max.X = (float32)((int)(r.Max.X))
	r.Max.Y = (float32)((int)(r.Max.Y))
}

func (r *Rectangle) AsInt32() (x, y, w, h int32) {
	x = int32(r.Min.X)
	y = int32(r.Min.Y)
	w = int32(r.Max.X - r.Min.X)
	h = int32(r.Max.Y - r.Min.Y)
	return
}

func (r *Rectangle) ClosestPoint(p Vector) Vector {
	if p.X > r.Max.X {
		p.X = r.Max.X
	} else if p.X < r.Min.X {
		p.X = r.Min.X
	}

	if p.X > r.Max.X {
		p.X = r.Max.X
	} else if p.X < r.Min.X {
		p.X = r.Min.X
	}

	return p
}
