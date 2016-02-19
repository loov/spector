package main

import (
	"log"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"

	"github.com/egonelbre/spector/ui"
)

type State struct {
	Backend ui.Backend
	Input   *ui.Input
}

func NewState() *State {
	state := &State{}

	state.Backend = ui.NewGLBackend()
	state.Input = &ui.Input{}

	return state
}

func (state *State) Stop() {
}

func (state *State) Update(dt float32) {
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
	x, y := window.GetCursorPos()
	down := window.GetMouseButton(glfw.MouseButtonLeft) == glfw.Press

	state.Input.Mouse.Position.X = float32(x)
	state.Input.Mouse.Position.Y = float32(y)
	state.Input.Mouse.WasDown = state.Input.Mouse.Down
	state.Input.Mouse.Down = down
}

func (state *State) Render(window *glfw.Window) {
	state.Reset(window)
	state.UpdateInput(window)

	w, h := window.GetSize()
	root := &ui.Context{
		Backend: state.Backend,
		Input:   state.Input,
		Area:    ui.Block(0, 0, float32(w), float32(h)),
	}

	state.Backend.SetBack(ui.ColorHex(0xEEEEEEFF))
	state.Backend.SetFore(ui.ColorHex(0xCCCCCCFF))
	state.Backend.SetFontColor(ui.ColorHex(0x000000FF))

	ui.Buttons{
		{"☺", nil},
		{"File", nil},
		{"Edit", nil},
		{"Help", nil},
	}.DoDynamic(ui.LayoutToRight(100, root.Top(20).Panel()))

	ui.Buttons{
		{"Alpha", nil},
		{"Beta", nil},
		{"Gamma", nil},
		{"Delta", nil},
		{"Iota", nil},
	}.Do(ui.LayoutToBottom(30, root.Right(150).Panel()))
}

type MainMenu struct{}

func (menu *MainMenu) File() {
	log.Println("File Pressed")
}

func (menu *MainMenu) Edit() {
	log.Println("Edit Pressed")
}

func (menu *MainMenu) Help() {
	log.Println("Help Pressed")
}
