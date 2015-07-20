package main

import (
	"log"
	"math/rand"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600
)

func init() { runtime.LockOSThread() }

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.DefaultWindowHints()
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	window, err := glfw.CreateWindow(ScreenWidth, ScreenHeight, "Spector", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	setupScene()
	defer destroyScene()
	for !window.ShouldClose() {
		drawScene()
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func setupScene() {
	gl.ClearColor(1, 1, 1, 0)

	//gl.Disable(gl.DEPTH_TEST)
	gl.Disable(gl.LIGHTING)

	gl.Viewport(0, 0, ScreenWidth, ScreenHeight)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, ScreenWidth, ScreenHeight, 0, -1, 1)
}

func destroyScene() {
}

var (
	x = float32(ScreenWidth) / 2
	y = float32(ScreenHeight) / 2
)

func drawScene() {
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.Color3f(rand.Float32(), rand.Float32(), rand.Float32())
	gl.Rectf(x, y, x+32, y+32)

	x += (rand.Float32() - 0.5) * 2
	y += (rand.Float32() - 0.5) * 2
}
