package main

import (
	"flag"
	"log"
	"runtime"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"

	"github.com/egonelbre/spector/ui"

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

	state := NewState()
	for !window.ShouldClose() {
		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			return
		}
		if window.GetKey(glfw.KeyR) == glfw.Press {
			state.Stop()
			state = NewState()
		}

		start := time.Now()
		state.Update(1.0 / 60.0)
		updateTime := time.Since(start)

		start = time.Now()
		state.Render(window)
		renderTime := time.Since(start)

		text := "update:" + updateTime.String() + " render: " + renderTime.String()
		w, h := window.GetSize()
		state.Backend.SetFontColor(ui.ColorHex(0xFF0000FF))
		size := state.Backend.Measure(text)
		state.Backend.Text(text, ui.Block(float32(w)-size.X, float32(h)-size.Y, size.X, size.Y))

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
