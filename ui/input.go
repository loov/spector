package ui

import (
	"time"

	"github.com/egonelbre/spector/ui/g"
)

type Input struct {
	Time  time.Time
	Mouse Mouse
}

type Mouse struct {
	Pos  g.Vector
	Down bool
	Last struct {
		Pos  g.Vector
		Down bool
	}
	Capture func() (done bool)
}

func (mouse *Mouse) BeginFrame() {
	mouse.Last.Pos = mouse.Pos
	mouse.Last.Down = mouse.Down
}

func (mouse *Mouse) EndFrame() {
	if mouse.Capture != nil {
		done := mouse.Capture()
		if done {
			mouse.Capture = nil
		}
	}
}
