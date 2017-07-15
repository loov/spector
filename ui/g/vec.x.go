package g

func (a Vector) Rotate() Vector { return Vector{-a.Y, a.X} }

func (a Vector) ScaleTo(size float32) Vector {
	ilen := a.Len()
	if ilen > 0 {
		ilen = size / ilen
	}
	return a.Scale(ilen)
}

func (r Rect) AsInt32() (x, y, w, h int32) {
	x = int32(r.Min.X)
	y = int32(r.Min.Y)
	w = int32(r.Max.X - r.Min.X)
	h = int32(r.Max.Y - r.Min.Y)
	return
}

// Corners returns top-left, top-right, bottom-right, bottom-left vectors
func (r Rect) Corners() (tl, tr, br, bl Vector) {
	tl = r.TopLeft()
	tr = r.TopRight()
	br = r.BottomRight()
	bl = r.BottomLeft()
	return
}

func (r Rect) Clip(clip Rect) {
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

func (r Rect) Floor() {
	r.Min.X = (float32)((int)(r.Min.X))
	r.Min.Y = (float32)((int)(r.Min.Y))
	r.Max.X = (float32)((int)(r.Max.X))
	r.Max.Y = (float32)((int)(r.Max.Y))
}

func (r Rect) ClosestPoint(p Vector) Vector {
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

func (r Rect) TopLeft() Vector     { return r.Min }
func (r Rect) TopRight() Vector    { return Vector{r.Max.X, r.Min.Y} }
func (r Rect) BottomLeft() Vector  { return Vector{r.Min.X, r.Max.Y} }
func (r Rect) BottomRight() Vector { return r.Max }

func SegmentNormal(a, b Vector) Vector {
	return b.Sub(a).Rotate()
}
