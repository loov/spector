// +build ignore

package ui

import (
	"image"
	"image/color"
	"io/ioutil"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/go-gl/gl/v2.1/gl"
)

type FontAtlas struct {
	TTF  *truetype.Font
	Face font.Face
}

func NewFontAtlas(filename string, dpi, fontSize float64) (*FontAtlas, error) {
	atlas := &FontAtlas{}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	atlas.TTF, err = truetype.Parse(content)
	if err != nil {
		return nil, err
	}

	opts := &truetype.Options{}
	opts.Size = fontSize

	atlas.Face = truetype.NewFace(atlas.TTF, opts)
	return atlas, nil
}

func (atlas *FontAtlas) draw(rendered *image.RGBA, b Bounds) {
	var texture uint32
	gl.Enable(gl.TEXTURE_2D)
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rendered.Bounds().Dx()),
		int32(rendered.Bounds().Dy()),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rendered.Pix))

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA)

	gl.Begin(gl.QUADS)
	{
		gl.TexCoord2f(0, 0)
		gl.Vertex2f(b.Min.X, b.Min.Y)
		gl.TexCoord2f(1, 0)
		gl.Vertex2f(b.Max.X, b.Min.Y)
		gl.TexCoord2f(1, 1)
		gl.Vertex2f(b.Max.X, b.Max.Y)
		gl.TexCoord2f(0, 1)
		gl.Vertex2f(b.Min.X, b.Max.Y)
	}
	gl.End()

	gl.Disable(gl.BLEND)

	gl.DeleteTextures(1, &texture)
	gl.Disable(gl.TEXTURE_2D)
}

func (atlas *FontAtlas) Draw(text string, b Bounds) {
	m := atlas.Face.Metrics()
	w := pow2(font.MeasureString(atlas.Face, text).Ceil())
	h := pow2(m.Ascent.Ceil() + m.Descent.Ceil())
	if w > 2048 {
		w = 2048
	}
	if h > 2048 {
		h = 2048
	}
	b.Max.X = b.Min.X + float32(w)
	b.Max.Y = b.Min.Y + float32(h)

	rendered := image.NewRGBA(image.Rect(0, 0, w, h))
	drawer := font.Drawer{
		Dst:  rendered,
		Src:  image.Black,
		Face: atlas.Face,
	}
	drawer.Dot = fixed.P(0, m.Ascent.Ceil())
	drawer.DrawString(text)

	atlas.draw(rendered, b)
}

func (atlas *FontAtlas) Measure(text string) (size Point) {
	m := atlas.Face.Metrics()
	size.X = float32(font.MeasureString(atlas.Face, text).Ceil())
	size.Y = float32(m.Ascent.Ceil() + m.Descent.Ceil())
	return size
}

func pow2(x int) int {
	x--
	x |= x >> 1
	x |= x >> 2
	x |= x >> 4
	x |= x >> 8
	x |= x >> 16
	return x + 1
}

func drawRect(rgba *image.RGBA, bounds image.Rectangle, color color.RGBA) {
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		rgba.SetRGBA(x, bounds.Min.Y, color)
		rgba.SetRGBA(x, bounds.Max.Y-1, color)
	}
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		rgba.SetRGBA(bounds.Min.X, y, color)
		rgba.SetRGBA(bounds.Max.X-1, y, color)
	}
}
