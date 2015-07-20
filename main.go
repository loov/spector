package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"

	"github.com/egonelbre/trace-spector/render"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600
)

var (
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

func init() { runtime.LockOSThread() }

func startprof() {
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
	}
}

func stopprof() {
	pprof.StopCPUProfile()
}

func main() {
	flag.Parse()

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

	start, stop := time.Now(), time.Now()

	px := time.Duration(ScreenHeight * ScreenHeight)

	buffer := render.NewBuffer(32 << 20)
	startprof()
	for !window.ShouldClose() {
		start = stop
		stop = time.Now()
		elapsed := stop.Sub(start)
		fmt.Printf("%.3fFPS  %v PR %v\n", 1/elapsed.Seconds(), elapsed/px, elapsed)

		buffer.Reset()
		DrawToBuffer(buffer)
		RenderBuffer(buffer)

		//RenderDirect()
		window.SwapBuffers()
		glfw.PollEvents()
	}
	stopprof()
}

func setupScene() {
	gl.ClearColor(1, 1, 1, 0)

	gl.Enable(gl.COLOR_MATERIAL)
	gl.Disable(gl.LIGHTING)
	gl.Disable(gl.DEPTH_TEST)

	gl.Viewport(0, 0, ScreenWidth, ScreenHeight)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, ScreenWidth, ScreenHeight, 0, -1, 1)
}

func destroyScene() {
}

func DrawToBuffer(buf *render.Buffer) {
	t := int(time.Now().UnixNano() / 100000)
	for x := 0; x < ScreenWidth; x++ {
		for y := 0; y < ScreenHeight; y++ {
			*buf.Color() = render.Color{
				R: byte(x & 0xFF),
				G: byte(y & 0xFF),
				B: byte((x*y + t) & 0xFF),
			}
			xf := float32(x)
			yf := float32(y)

			*buf.Rect() = render.Rect{
				render.Point{xf, yf},
				render.Point{xf + 1, yf + 1},
			}
		}
	}
}

func RenderDirect() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	t := int(time.Now().UnixNano() / 100000)
	for x := 0; x < ScreenWidth; x++ {
		for y := 0; y < ScreenHeight; y++ {
			gl.Color3ub(
				byte(x&0xFF),
				byte(y&0xFF),
				byte((x*y+t)&0xFF),
			)

			xf := float32(x)
			yf := float32(y)
			gl.Rectf(xf, yf, xf+1, yf+1)
		}
	}
}

func RenderBuffer(buf *render.Buffer) {
	gl.Clear(gl.COLOR_BUFFER_BIT)

	rd := buf.Reader()
	for rd.Next() {
		switch rd.Op() {
		case render.OpColor:
			c := rd.Color()
			gl.Color3ub(c.R, c.G, c.B)
		case render.OpRect:
			r := rd.Rect()
			gl.Rectf(r.A.X, r.A.Y, r.B.X, r.B.Y)
		case render.OpTri:
			gl.Begin(gl.TRIANGLES)
			t := rd.Tri()
			gl.Vertex2f(t.A.X, t.A.Y)
			gl.Vertex2f(t.B.X, t.B.Y)
			gl.Vertex2f(t.C.X, t.C.Y)
			gl.End()
		case render.OpQuad:
			gl.Begin(gl.QUADS)
			t := rd.Quad()
			gl.Vertex2f(t.A.X, t.A.Y)
			gl.Vertex2f(t.B.X, t.B.Y)
			gl.Vertex2f(t.C.X, t.C.Y)
			gl.Vertex2f(t.D.X, t.D.Y)
			gl.End()
		case render.OpLine:
			gl.Begin(gl.LINES)
			t := rd.Line()
			gl.Vertex2f(t.A.X, t.A.Y)
			gl.Vertex2f(t.B.X, t.B.Y)
			gl.End()
		default:
			panic("unimplemented")
		}
	}
}
