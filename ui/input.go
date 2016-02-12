package ui

type Mouse struct {
	Position Point
	Down     bool
	WasDown  bool
}

func (m *Mouse) PointsAt(b Bounds) bool {
	return b.Contains(m.Position)
}

func (m *Mouse) Clicked() bool {
	return !m.Down && m.WasDown
}

type Input struct {
	Mouse Mouse
}
