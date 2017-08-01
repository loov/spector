package main

import (
	"flag"
	"log"
	"runtime"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"

	"github.com/egonelbre/spector/ui"
	"github.com/egonelbre/spector/ui/g"
	render "github.com/egonelbre/spector/ui/render/gl21"
	"github.com/egonelbre/spector/ui/screen"

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
	// glfw.WindowHint(glfw.Visible, glfw.False) // do not steal focus

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
	Context *ui.Context
	Screen  *screen.Screen

	LastCursor ui.Cursor
	Cursors    map[ui.Cursor]*glfw.Cursor
}

func NewApp(window *glfw.Window) *App {
	app := &App{}
	app.Window = window
	app.Context = ui.NewContext()
	app.Screen = screen.New()

	app.Cursors = make(map[ui.Cursor]*glfw.Cursor)
	app.Cursors[ui.ArrowCursor] = glfw.CreateStandardCursor(glfw.ArrowCursor)
	app.Cursors[ui.IBeamCursor] = glfw.CreateStandardCursor(glfw.IBeamCursor)
	app.Cursors[ui.CrosshairCursor] = glfw.CreateStandardCursor(glfw.CrosshairCursor)
	app.Cursors[ui.HandCursor] = glfw.CreateStandardCursor(glfw.HandCursor)
	app.Cursors[ui.HResizeCursor] = glfw.CreateStandardCursor(glfw.HResizeCursor)
	app.Cursors[ui.VResizeCursor] = glfw.CreateStandardCursor(glfw.VResizeCursor)

	return app
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
	fw, fh := app.Window.GetFramebufferSize()
	w, h := app.Window.GetSize()
	x, y := app.Window.GetCursorPos()

	app.Context.Input.Mouse.Pos = g.Vector{float32(x), float32(y)}
	app.Context.Input.Mouse.Down = app.Window.GetMouseButton(glfw.MouseButtonLeft) == glfw.Press

	app.Context.BeginFrame(g.Rect{
		g.Vector{0, 0},
		g.Vector{float32(w), float32(h)},
	})

	app.Context.Input.Time = time.Now()

	app.RenderFrame()
	app.Context.EndFrame()

	if app.LastCursor != app.Context.Input.Mouse.Cursor {
		app.LastCursor = app.Context.Input.Mouse.Cursor
		app.Window.SetCursor(app.Cursors[app.LastCursor])
	}

	{ // reset window
		gl.MatrixMode(gl.MODELVIEW)
		gl.LoadIdentity()

		gl.Viewport(0, 0, int32(fw), int32(fh))
		gl.Ortho(0, float64(w), float64(h), 0, 30, -30)
		gl.ClearColor(1, 1, 1, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)
	}

	for _, list := range app.Context.Render.Frame.Lists {
		render.List(w, h, list)
	}
}

func (app *App) RenderFrame() {
	app.Screen.Update(app.Context)
}
