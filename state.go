package main

import (
	"fmt"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"

	"github.com/egonelbre/spector/timeline"
	"github.com/egonelbre/spector/trace"
	"github.com/egonelbre/spector/trace/simulator"

	"github.com/egonelbre/spector/ui"
)

var Highlight = ui.Color{0xff, 0, 0xff, 0xff}

type State struct {
	Timeline timeline.Timeline
	Handler  timeline.Handler

	Backend ui.Backend
	Time    time.Time
	Input   *ui.Input

	Simulator *simulator.Stream
	Dirty     bool
}

func NewState() *State {
	state := &State{}
	state.Handler.Timeline = &state.Timeline
	state.Simulator = simulator.NewStream()

	state.Simulator.Start()
	state.Dirty = true

	state.Backend = ui.NewGLBackend()
	state.Time = time.Now()
	state.Input = &ui.Input{}

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

func (state *State) Reset(window *glfw.Window) {
	gl.ClearColor(1, 1, 1, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	gl.Disable(gl.DEPTH)
	gl.Enable(gl.FRAMEBUFFER_SRGB)

	width, height := window.GetSize()
	gl.Viewport(0, 0, int32(width), int32(height))
	gl.Ortho(0, float64(width), float64(height), 0, 30, -30)
}

func (state *State) UpdateInput(window *glfw.Window) {
	state.Input.Update()

	x, y := window.GetCursorPos()
	state.Input.Mouse.Position.X = float32(x)
	state.Input.Mouse.Position.Y = float32(y)
	state.Input.Mouse.Down = window.GetMouseButton(glfw.MouseButtonLeft) == glfw.Press
}

func (state *State) Render(window *glfw.Window) {
	if !state.Dirty {
		return
	}
	state.Dirty = false

	state.Reset(window)
	state.Time = time.Now()
	state.UpdateInput(window)

	hue := float32(state.Time.UnixNano()/1e6%360) / 360.0
	Highlight = ui.ColorHSLA(hue, 0.7, 0.7, 1.0)

	w, h := window.GetSize()
	root := &ui.Context{
		Backend: state.Backend,
		Input:   state.Input,
		Area:    ui.Block(0, 0, float32(w), float32(h)),
	}

	if root.Input.Mouse.Drag != nil {
		if !root.Input.Mouse.Drag(root) {
			root.Input.Mouse.Drag = nil
		}
	}

	state.Backend.SetBack(ui.ColorHex(0xEEEEEEFF))
	state.Backend.SetFore(ui.ColorHex(0xCCCCCCFF))
	state.Backend.SetFontColor(ui.ColorHex(0x000000FF))

	view := NewView(state, root, &state.Timeline)
	view.Render()
}

type Camera struct {
	Start     trace.Time
	Stop      trace.Time
	Span      trace.Time
	TimePerPx trace.Time
}

type View struct {
	State    *State
	Context  *ui.Context
	Timeline *timeline.Timeline

	Size ui.Point
	Y    float32

	Camera
	Target Camera
}

func NewView(state *State, context *ui.Context, timeline *timeline.Timeline) *View {
	ui := &View{}
	ui.State = state
	ui.Timeline = timeline
	ui.Context = context
	ui.Size = context.Area.Size()

	ui.Start = 0
	ui.Stop = 1e5
	if len(timeline.Procs) > 0 {
		ui.Stop = timeline.Procs[0].Time
	}
	ui.Span = ui.Stop - ui.Start

	ui.TimePerPx = ui.Span / trace.Time(context.Area.Dx())
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

func IDColor(id trace.ID) ui.Color {
	h := float32((id*31)%256) / 256.0
	return ui.ColorHSLA(h, 0.7, 0.8, 1)
}

func min(a, b trace.Time) trace.Time {
	if a < b {
		return a
	}
	return b
}

func (view *View) Render() {
	const TrackPadding = 4
	const LayerHeight = 10
	const LayerPadding = 3

	view.H(fontViewSummary, "%d", view.Timeline.TotalEvents)

	for _, proc := range view.Timeline.Procs {
		view.H(fontProcHeader, "M:%08X P:%08X", proc.MID, proc.PID)

		for _, track := range proc.Tracks {
			view.Pad(TrackPadding)

			maxY := view.Y
			baseY := view.Y

			for _, thread := range track.Threads {
				view.Y = baseY

				view.Block(
					thread.TID,
					thread.Start, min(thread.Stop, proc.Time),
					float32(len(thread.Layers)*LayerHeight),
				)
				for i, layer := range thread.Layers {
					depth := float32(len(thread.Layers)-i) / float32(len(thread.Layers))
					view.Spans(proc, layer, depth, LayerHeight)
				}
				view.Pad(LayerPadding)

				if view.Y > maxY {
					maxY = view.Y
				}
			}
			view.Y = maxY
		}
	}
}

type Font struct {
	Height     float32
	Foreground uint32
	Background uint32
}

func (view *View) H(font *Font, format string, args ...interface{}) {
	view.Context.Backend.SetBack(ui.ColorHex(font.Background))
	view.Context.Backend.SetFontColor(ui.ColorHex(font.Foreground))
	bounds := ui.Block(0, view.Y, view.Size.X, font.Height)
	view.Context.Backend.Fill(bounds)
	text := fmt.Sprintf(format, args...)
	view.Context.Backend.Text(text, bounds)

	view.Y += font.Height
}

func (view *View) Pad(height float32) {
	view.Y += height
}

func (view *View) TimeToPx(t trace.Time) float32 {
	v := float32(t-view.Start) * view.Size.X / float32(view.Span)
	if v > view.Size.X {
		v = view.Size.X
	}
	return v
}

func (view *View) Block(id trace.ID, start, stop trace.Time, height float32) {
	x0, x1 := view.TimeToPx(start), view.TimeToPx(stop)

	view.Pad(2)
	block := ui.Block(x0, view.Y, x1-x0, height)
	view.Context.SetBack(IDColor(id))
	view.Context.Backend.Fill(block)
	if view.State.Input.Mouse.PointsAt(block) {
		view.Context.Backend.SetFore(Highlight)
		view.Context.Backend.Stroke(block)
	}
}

func (view *View) Spans(proc *timeline.Proc, layer *timeline.Layer, depth, height float32) {
	mingap := view.TimePerPx * 2

	view.Context.SetBack(ui.ColorFloat(0.5, depth, depth, 1))
	for i := 0; i < len(layer.Spans); i++ {
		span := layer.Spans[i]
		if span.Stop < view.Start {
			continue
		}
		if span.Start > view.Stop {
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

		x0 := view.TimeToPx(join.Start)
		x1 := view.TimeToPx(min(join.Stop, proc.Time))
		block := ui.Block(x0, view.Y, x1-x0, height)
		view.Context.Backend.Fill(block)
		if view.State.Input.Mouse.PointsAt(block) {
			view.Context.Backend.SetFore(Highlight)
			view.Context.Backend.Stroke(block)
		}
	}
	view.Y += height
}
