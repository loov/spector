package main

import (
	"flag"
	"log"
	"runtime"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"

	"github.com/egonelbre/spector/ui"
	"github.com/egonelbre/spector/ui/g"
	render "github.com/egonelbre/spector/ui/render/gl21"

	"net/http"
	_ "net/http/pprof"
)

func init() { runtime.LockOSThread() }

func main() {
	flag.Parse()

	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.Visible, glfw.False) // do not steal focus

	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	window, err := glfw.CreateWindow(800, 600, "Spector", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	window.Restore()
	window.SetPos(32, 64)

	if err := gl.Init(); err != nil {
		panic(err)
	}

	app := NewApp(window)
	app.Run()
}

type App struct {
	Window  *glfw.Window
	Context ui.Context
	View    ui.View
}

func NewApp(window *glfw.Window) *App {
	return &App{Window: window}
}

func (app *App) Run() {
	for !app.Window.ShouldClose() {
		if app.Window.GetKey(glfw.KeyEscape) == glfw.Press {
			return
		}
		if app.Window.GetKey(glfw.KeyF10) == glfw.Press {
			*app = *NewApp(app.Window)
		}

		app.UpdateFrame()

		app.Window.SwapBuffers()
		glfw.PollEvents()
	}
}

func (app *App) UpdateFrame() {
	w, h := app.Window.GetSize()
	app.Context.BeginFrame(g.Rect{
		g.Vector{0, 0},
		g.Vector{float32(w), float32(h)},
	})

	x, y := app.Window.GetCursorPos()
	app.Context.Input.Mouse.Pos = g.Vector{float32(x), float32(y)}
	app.Context.Input.Mouse.Down = app.Window.GetMouseButton(glfw.MouseButtonLeft) == glfw.Press

	app.Context.Input.Time = time.Now()

	app.RenderFrame()

	{ // reset window
		gl.MatrixMode(gl.MODELVIEW)
		gl.LoadIdentity()

		gl.Viewport(0, 0, int32(w), int32(h))
		gl.Ortho(0, float64(w), float64(h), 0, 30, -30)
		gl.ClearColor(1, 1, 1, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)
	}

	for _, list := range app.Context.Render.Frame.Lists {
		render.List(w, h, list)
	}
}

func (app *App) RenderFrame() {
	mouse := &app.Context.Input.Mouse
	app.Context.Draw.AddCircle(mouse.Pos, 5, g.Black)
}
