package ui

type Mouse struct {
	Position Point
	Down     bool
	Last     struct {
		Position Point
		Down     bool
	}
}

func (m *Mouse) Update() {
	m.Last.Position = m.Position
	m.Last.Down = m.Down
}

func (m *Mouse) PointsAt(b Bounds) bool {
	return b.Contains(m.Position)
}

func (m *Mouse) Clicked() bool {
	return !m.Down && m.Last.Down
}

type Input struct {
	Mouse Mouse
}
