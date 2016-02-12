package main

import (
	"log"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"

	"github.com/egonelbre/spector/ui"
)

type State struct {
	Atlas *ui.FontAtlas
	UI    *ui.State
}

func NewState() *State {
	state := &State{}

	var err error
	state.Atlas, err = ui.NewFontAtlas("~DejaVuSans.ttf", 72, 12)
	if err != nil {
		panic(err)
	}

	state.UI = &ui.State{}
	state.UI.Font = state.Atlas

	return state
}

func (state *State) Stop() {
}

func (state *State) Update(dt float32) {
}

func (state *State) Render(window *glfw.Window) {
	gl.ClearColor(1, 1, 1, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	gl.Disable(gl.DEPTH)
	gl.Enable(gl.FRAMEBUFFER_SRGB)

	width, height := window.GetSize()
	gl.Viewport(0, 0, int32(width), int32(height))
	gl.Ortho(0, float64(width), float64(height), 0, 30, -30)

	x, y := window.GetCursorPos()
	down := window.GetMouseButton(glfw.MouseButtonLeft) == glfw.Press

	root := state.UI

	root.Input.Mouse.Position = ui.Point{float32(x), float32(y)}
	root.Input.Mouse.PDown = root.Input.Mouse.Down
	root.Input.Mouse.Down = down

	root.Panel(ui.Rect(float32(width-200), 0, 200, float32(height)), func() {
		r := ui.Rect(0, 0, 200, 30)
		d := ui.Point{0, r.Dy()}
		if root.Button("alpha", r) {
			log.Println("alpha pressed")
		}
		r = r.Offset(d)
		if root.Button("beta", r) {
			log.Println("beta pressed")
		}
		r = r.Offset(d)
		if root.Button("gamma", r) {
			log.Println("gamma pressed")
		}
	})
}
