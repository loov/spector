package main

import (
	"flag"
	"log"
	"math"
	"runtime"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"

	"github.com/egonelbre/spector/ui/draw"

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
	var DrawList draw.List

	for !window.ShouldClose() {
		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			return
		}

		now := float64(time.Now().UnixNano()) / 1e9
		width, height := window.GetSize()

		{ // reset window
			gl.ClearColor(1, 1, 1, 1)
			gl.Clear(gl.COLOR_BUFFER_BIT)
			gl.MatrixMode(gl.MODELVIEW)
			gl.LoadIdentity()

			gl.Disable(gl.DEPTH)
			gl.Enable(gl.FRAMEBUFFER_SRGB)

			gl.Viewport(0, 0, int32(width), int32(height))
			gl.Ortho(0, float64(width), float64(height), 0, 30, -30)
		}

		DrawList.Reset()

		DrawList.BeginCommand()
		DrawList.AddRectFill(&draw.Rectangle{
			draw.Vector{10, 10},
			draw.Vector{50, 50},
		}, draw.Red)

		const LineCount = 16
		var line [LineCount]draw.Vector
		for i := range line {
			r := float64(i) / LineCount
			line[i].X = float32(r) * float32(width)
			line[i].Y = float32(height)*0.5 + float32(math.Sin(r*11.8+now*3)*100)
		}
		DrawList.AddLine(line[:], false, 10.0, draw.Blue)

		const CircleCount = 32
		var circle [CircleCount]draw.Vector
		for i := range circle {
			p := float64(i) / CircleCount
			a := now + p*math.Pi*2
			w := math.Sin(p*62)*20.0 + 100.0
			circle[i].X = float32(width)*0.5 + float32(math.Cos(a)*w)
			circle[i].Y = float32(height)*0.5 + float32(math.Sin(a)*w)
		}
		DrawList.AddLine(circle[:], true, 10.0, draw.Green)

		gl.Begin(gl.TRIANGLES)
		indices := DrawList.Indicies
		vertices := DrawList.Vertices
		for _, cmd := range DrawList.Commands {
			if cmd.Texture == 0 {
				lastColor := draw.Color{0, 0, 0, 0}
				for _, vi := range indices[:cmd.Count] {
					v := vertices[int(vi)]
					if lastColor != v.Color {
						gl.Color4ub(v.Color.R, v.Color.G, v.Color.B, v.Color.A)
						lastColor = v.Color
					}
					gl.Vertex2f(v.P.X, v.P.Y)
				}
			} else {
				lastColor := draw.Color{0, 0, 0, 0}
				for _, vi := range indices[:cmd.Count] {
					v := vertices[int(vi)]
					if lastColor != v.Color {
						gl.Color4ub(v.Color.R, v.Color.G, v.Color.B, v.Color.A)
						lastColor = v.Color
					}
					gl.TexCoord2f(v.UV.X, v.UV.Y)
					gl.Vertex2f(v.P.X, v.P.Y)
				}
			}
			indices = indices[cmd.Count:]
		}
		gl.End()

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
