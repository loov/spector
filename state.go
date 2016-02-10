package main

import (
	"log"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"

	"github.com/egonelbre/spector/timeline"
	"github.com/egonelbre/spector/trace"
	"github.com/egonelbre/spector/trace/simulator"

	"github.com/egonelbre/spector/ui"
	"github.com/egonelbre/spector/ui/font"
)

type State struct {
	Timeline timeline.Timeline
	Handler  timeline.Handler

	Atlas *font.Atlas

	UI *ui.State

	Simulator *simulator.Stream
	Dirty     bool
}

func NewState() *State {
	state := &State{}
	state.Handler.Timeline = &state.Timeline
	state.Simulator = simulator.NewStream()

	var err error
	state.Atlas, err = font.NewAtlas("~DejaVuSans.ttf", 72, 12)
	if err != nil {
		panic(err)
	}

	state.Simulator.Start()
	state.Dirty = true
	state.UI = &ui.State{}

	return state
}

func (state *State) Stop() {
	state.Simulator.Stop()
}

func (state *State) Update(dt float32) {
	for _, event := range state.Simulator.Next() {
		state.Dirty = true
		state.Handler.Handle(event)
	}
}

func (state *State) Render(window *glfw.Window) {
	if !state.Dirty {
		return
	}
	state.Dirty = false

	gl.ClearColor(1, 1, 1, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	gl.Disable(gl.DEPTH)
	gl.Enable(gl.FRAMEBUFFER_SRGB)

	width, height := window.GetSize()
	gl.Viewport(0, 0, int32(width), int32(height))
	gl.Ortho(0, float64(width), float64(height), 0, 30, -30)

	view := NewView(V2{float32(width), float32(height)}, &state.Timeline)
	view.Atlas = state.Atlas
	view.Render()

	x, y := window.GetCursorPos()
	down := window.GetMouseButton(glfw.MouseButtonLeft) == glfw.Press

	root := state.UI

	root.Input.Mouse.Position = ui.Point{float32(x), float32(y)}
	root.Input.Mouse.PDown = root.Input.Mouse.Down
	root.Input.Mouse.Down = down

	if root.Button("alpha", ui.Rect(10, 10, 100, 30)) {
		log.Println("alpha pressed")
	}
	if root.Button("beta", ui.Rect(10, 50, 100, 30)) {
		log.Println("beta pressed")
	}
	if root.Button("gamma", ui.Rect(10, 90, 100, 30)) {
		log.Println("gamma pressed")
	}
}

type V2 struct{ X, Y float32 }

type View struct {
	Atlas    *font.Atlas
	Timeline *timeline.Timeline

	Size V2
	Y    float32

	Start trace.Time // time
	Stop  trace.Time // time
	Span  trace.Time // time

	TimePerPx trace.Time // time / px
}

func NewView(size V2, timeline *timeline.Timeline) *View {
	ui := &View{}
	ui.Timeline = timeline
	ui.Size = size

	ui.Start = 0
	ui.Stop = 1e5
	ui.Span = ui.Stop - ui.Start

	ui.TimePerPx = ui.Span / trace.Time(size.X)
	return ui
}

var (
	fontViewSummary = &Font{
		Height:     18,
		Foreground: 0xFFFFFFff,
		Background: 0x000000ff,
	}

	fontProcHeader = &Font{
		Height:     18,
		Foreground: 0xEEEEEEff,
		Background: 0x333333ff,
	}
)

func IDColor(id trace.ID) (r, g, b, a float32) {
	h := float32((id*31)%256) / 256.0
	return HSLA(h, 0.7, 0.8, 1)
}

func (ui *View) Render() {
	const TrackPadding = 4
	const LayerHeight = 10
	const LayerPadding = 3

	ui.H(fontViewSummary, "%d", ui.Timeline.TotalEvents)

	for _, proc := range ui.Timeline.Procs {
		ui.H(fontProcHeader, "M:%08X P:%08X", proc.MID, proc.PID)

		for _, track := range proc.Tracks {
			ui.Pad(TrackPadding)

			maxY := ui.Y
			baseY := ui.Y

			for _, thread := range track.Threads {
				ui.Y = baseY

				ui.Block(
					thread.TID,
					thread.Start, thread.Stop,
					float32(len(thread.Layers)*LayerHeight),
				)
				for i, layer := range thread.Layers {
					depth := float32(len(thread.Layers)-i) / float32(len(thread.Layers))
					ui.Spans(layer, depth, LayerHeight)
				}
				ui.Pad(LayerPadding)

				if ui.Y > maxY {
					maxY = ui.Y
				}
			}
			ui.Y = maxY
		}
	}
}

func (ui *View) Rect(p, size V2) {
	gl.Begin(gl.QUADS)
	{
		gl.Vertex2f(p.X, p.Y)
		gl.Vertex2f(p.X+size.X, p.Y)
		gl.Vertex2f(p.X+size.X, p.Y+size.Y)
		gl.Vertex2f(p.X, p.Y+size.Y)
	}
	gl.End()
}

func (ui *View) Bound(p, size V2) {
	gl.Begin(gl.LINE_LOOP)
	{
		gl.Vertex2f(p.X, p.Y)
		gl.Vertex2f(p.X+size.X-1, p.Y)
		gl.Vertex2f(p.X+size.X-1, p.Y+size.Y)
		gl.Vertex2f(p.X, p.Y+size.Y)
	}
	gl.End()
}

func RGBA(color uint32) (r, g, b, a uint8) {
	r = uint8(color >> 24)
	g = uint8(color >> 16)
	b = uint8(color >> 8)
	a = uint8(color >> 0)
	return
}

func hue(v1, v2, h float32) float32 {
	if h < 0 {
		h += 1
	}
	if h > 1 {
		h -= 1
	}
	if 6*h < 1 {
		return v1 + (v2-v1)*6*h
	} else if 2*h < 1 {
		return v2
	} else if 3*h < 2 {
		return v1 + (v2-v1)*(2.0/3.0-h)*6
	}

	return v1
}

func HSLA(h, s, l, a float32) (r, g, b, ra float32) {
	if s == 0 {
		return l, l, l, a
	}

	var v2 float32
	if l < 0.5 {
		v2 = l * (1 + s)
	} else {
		v2 = (l + s) - s*l
	}

	v1 := 2*l - v2
	r = hue(v1, v2, h+1.0/3.0)
	g = hue(v1, v2, h)
	b = hue(v1, v2, h-1.0/3.0)
	ra = a

	return
}

type Font struct {
	Height     float32
	Foreground uint32
	Background uint32
}

func (ui *View) H(font *Font, format string, args ...interface{}) {
	gl.Color4ub(RGBA(font.Background))
	ui.Rect(V2{0, ui.Y}, V2{ui.Size.X, font.Height})
	ui.Y += font.Height

	gl.Color4ub(RGBA(font.Foreground))
	ui.Atlas.Drawf(10, ui.Y-4, format, args...)
}

func (ui *View) Pad(height float32) {
	ui.Y += height
}

func (ui *View) TimeToPx(t trace.Time) float32 {
	v := float32(t-ui.Start) * ui.Size.X / float32(ui.Span)
	if v > ui.Size.X {
		v = ui.Size.X
	}
	return v
}

func (ui *View) Block(id trace.ID, start, stop trace.Time, height float32) {
	x0, x1 := ui.TimeToPx(start), ui.TimeToPx(stop)

	// ui.H at x0, x1
	ui.Pad(2)
	gl.Color4f(IDColor(id))
	ui.Rect(V2{x0, ui.Y}, V2{x1 - x0, height})
}

func (ui *View) Spans(layer *timeline.Layer, depth, height float32) {
	mingap := ui.TimePerPx
	// gl.Color3d(IDColor(depth))
	gl.Color4f(0.5, depth, depth, 1)
	for i := 0; i < len(layer.Spans); i++ {
		span := layer.Spans[i]
		if span.Stop < ui.Start {
			continue
		}
		if span.Start > ui.Stop {
			break
		}

		join := span
		for ; i < len(layer.Spans); i++ {
			last := layer.Spans[i]
			if last.Start-join.Stop > mingap {
				i--
				break
			}
			join.Stop = last.Stop
		}

		x0, x1 := ui.TimeToPx(join.Start), ui.TimeToPx(join.Stop)
		ui.Rect(V2{x0, ui.Y}, V2{x1 - x0, height})
	}
	ui.Y += height
}
