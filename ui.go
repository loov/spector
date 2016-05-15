// +build ignore

package main

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"

	"github.com/egonelbre/spector/ui"
)

type State struct {
	Backend ui.Backend
	Input   *ui.Input

	MemStats      runtime.MemStats
	SidePanelSize float32

	Box ui.Point
}

func NewState() *State {
	state := &State{}

	state.Backend = ui.NewGLBackend()
	state.Input = &ui.Input{}
	state.SidePanelSize = 350

	state.Box = ui.Pt(10, 10)

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
	state.Input.Update()

	x, y := window.GetCursorPos()
	state.Input.Mouse.Position.X = float32(x)
	state.Input.Mouse.Position.Y = float32(y)
	state.Input.Mouse.Down = window.GetMouseButton(glfw.MouseButtonLeft) == glfw.Press
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

	if root.Input.Mouse.Drag != nil {
		if !root.Input.Mouse.Drag(root) {
			root.Input.Mouse.Drag = nil
		}
	}

	state.Backend.SetBack(ui.ColorHex(0xEEEEEEFF))
	state.Backend.SetFore(ui.ColorHex(0xCCCCCCFF))
	state.Backend.SetFontColor(ui.ColorHex(0x000000FF))

	mm := &MainMenu{}
	ui.Buttons{
		{"☺", mm.File},
		{"Load", nil},
		{"Save", nil},
		{"Quit", nil},
	}.DoDynamic(ui.LayoutToRight(50, root.Top(20).Panel()))

	boxbounds := ui.Bounds{
		state.Box,
		state.Box.Offset(ui.Point{30, 30})}

	runtime.ReadMemStats(&state.MemStats)
	root.Right(state.SidePanelSize).Reflect("Mem", &state.MemStats)

	ui.SplitterX(root.Right(5), func(delta float32) {
		state.SidePanelSize -= delta
	})

	box := root.Child(boxbounds)
	box.Backend.SetBack(ui.ColorHex(0xFF88EEFF))
	box.Backend.SetFore(ui.ColorHex(0xFF88EEFF))
	box.Backend.Fill(box.Area)
	box.Backend.Stroke(box.Area)

	ui.DoDragXY(box, func(dx, dy float32) {
		state.Box.X += dx
		state.Box.Y += dy
	})
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
