package ui

type Point struct{ X, Y float32 }

func (p Point) Offset(by Point) Point {
	return Point{p.X + by.X, p.Y + by.Y}
}

type Bounds struct{ Min, Max Point }

func Rect(x, y, w, h float32) Bounds {
	if w < 0 {
		x, w = x+w, -w
	}
	if h < 0 {
		y, h = y+h, -h
	}
	return Bounds{
		Min: Point{x, y},
		Max: Point{x + w, y + h},
	}
}

func (b Bounds) Dx() float32 { return b.Max.X - b.Min.X }
func (b Bounds) Dy() float32 { return b.Max.Y - b.Min.Y }

func (b Bounds) Offset(by Point) Bounds {
	return Bounds{
		Min: b.Min.Offset(by),
		Max: b.Max.Offset(by),
	}
}

func (b Bounds) Contains(p Point) bool {
	return b.Min.X <= p.X && p.X < b.Max.X &&
		b.Min.Y <= p.Y && p.Y < b.Max.Y
}
