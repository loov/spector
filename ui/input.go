package ui

import (
	"time"

	"github.com/egonelbre/spector/ui/g"
)

type Input struct {
	Time  time.Time
	Mouse Mouse
}

type Cursor byte

const (
	ArrowCursor = Cursor(iota)
	IBeamCursor
	CrosshairCursor
	HandCursor
	HResizeCursor
	VResizeCursor
)

type Mouse struct {
	Pos      g.Vector
	Down     bool
	Pressed  bool
	Released bool
	Cursor   Cursor
	Last     struct {
		Pos  g.Vector
		Down bool
	}
	Capture func() (done bool)
}

func (mouse *Mouse) BeginFrame() {
	mouse.Cursor = ArrowCursor
	mouse.Pressed = !mouse.Last.Down && mouse.Down
	mouse.Released = mouse.Last.Down && !mouse.Down
}

func (mouse *Mouse) EndFrame() {
	mouse.Last.Pos = mouse.Pos
	mouse.Last.Down = mouse.Down

	if mouse.Capture != nil {
		done := mouse.Capture()
		if done {
			mouse.Capture = nil
		}
	}
}
