package draw

type Color struct{ R, G, B, A uint8 }

func (c Color) Transparent() bool { return c.A == 0x0 }

var (
	Black = Color{0, 0, 0, 255}
	White = Color{255, 255, 255, 255}
	Red   = Color{255, 0, 0, 255}
	Green = Color{0, 255, 0, 255}
	Blue  = Color{0, 0, 255, 255}
)
