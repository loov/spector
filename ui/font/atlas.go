package font

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"log"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/go-gl/gl/v2.1/gl"
)

type Point struct{ X, Y float32 }

type RelativeRectangle struct{ Min, Max Point }

func RelativeRect(r, b image.Rectangle) (n RelativeRectangle) {
	n.Min.X = float32(r.Min.X-b.Min.X) / float32(b.Dx())
	n.Min.Y = float32(r.Min.Y-b.Min.Y) / float32(b.Dy())
	n.Max.X = float32(r.Max.X-b.Min.X) / float32(b.Dx())
	n.Max.Y = float32(r.Max.Y-b.Min.Y) / float32(b.Dy())
	return n
}

type Glyph struct {
	Rune    rune
	Loc     image.Rectangle     // absolute location on image atlas
	RelLoc  RelativeRectangle   // relative location on image atlas
	Bounds  fixed.Rectangle26_6 // such that point + bounds, gives image bounds where glyph should be drawn
	Advance fixed.Int26_6       // advance from point, to the next glyph
}

type Atlas struct {
	Context *freetype.Context
	TTF     *truetype.Font
	Face    font.Face

	Rendered       map[rune]Glyph
	Image          *image.RGBA
	CursorX        int
	CursorY        int
	Padding        int
	MaxGlyphHeight int

	Dirty   bool
	Texture uint32
}

func NewAtlas(filename string, dpi, fontSize float64) (*Atlas, error) {
	atlas := &Atlas{}
	atlas.Rendered = make(map[rune]Glyph, 256)
	atlas.Padding = 2

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	atlas.TTF, err = truetype.Parse(content)
	if err != nil {
		return nil, err
	}

	atlas.Image = image.NewRGBA(image.Rect(0, 0, 1024, 1024))
	//draw.Draw(atlas.Image, atlas.Image.Bounds(), image.White, image.ZP, draw.Src)

	atlas.Context = freetype.NewContext()
	atlas.Context.SetDPI(dpi)

	atlas.Context.SetFont(atlas.TTF)
	atlas.Context.SetFontSize(fontSize)
	atlas.Context.SetHinting(font.HintingNone)

	atlas.Context.SetClip(atlas.Image.Bounds())
	atlas.Context.SetSrc(image.Black)
	atlas.Context.SetDst(atlas.Image)

	opts := &truetype.Options{}
	opts.Size = fontSize
	opts.Hinting = font.HintingNone

	atlas.Face = truetype.NewFace(atlas.TTF, opts)
	return atlas, nil
}

func ceilPx(i fixed.Int26_6) int {
	const ceiling = 1<<6 - 1
	return int(i+ceiling) >> 6
}

func ceilPxf(i fixed.Int26_6) float32 {
	const div = 1 << 6
	return float32(i) / div
}

func (atlas *Atlas) loadGlyph(r rune) {
	if _, ok := atlas.Rendered[r]; ok {
		return
	}
	atlas.Dirty = true

	glyph := Glyph{}
	glyph.Rune = r

	bounds, advance, _ := atlas.Face.GlyphBounds(r)
	glyph.Bounds = bounds
	glyph.Advance = advance

	width := ceilPx(bounds.Max.X - bounds.Min.X)
	height := ceilPx(bounds.Max.Y - bounds.Min.Y)

	if atlas.CursorX+atlas.Padding+width+atlas.Padding > atlas.Image.Bounds().Dx() {
		atlas.CursorX = 0
		atlas.CursorY += atlas.Padding + atlas.MaxGlyphHeight
	}

	x := atlas.CursorX + atlas.Padding
	y := atlas.CursorY + atlas.Padding

	glyph.Loc = image.Rect(x, y, x+width, y+height)
	glyph.RelLoc = RelativeRect(glyph.Loc, atlas.Image.Bounds())

	pt := fixed.P(x, y).Sub(bounds.Min)
	atlas.Context.DrawString(string(r), pt)

	if height > atlas.MaxGlyphHeight {
		atlas.MaxGlyphHeight = height
	}
	atlas.CursorX += atlas.Padding + width + atlas.Padding

	atlas.Rendered[r] = glyph
}

func (atlas *Atlas) LoadAscii() {
	for r := rune(0); r < 128; r++ {
		atlas.loadGlyph(r)
	}
	atlas.upload()
}

func (atlas *Atlas) LoadExtendedAscii() {
	for r := rune(0); r < 256; r++ {
		atlas.loadGlyph(r)
	}
	atlas.upload()
}

func (atlas *Atlas) LoadGlyphs(text string) {
	for _, r := range text {
		atlas.loadGlyph(r)
	}
	atlas.upload()
}

func (atlas *Atlas) upload() {
	if !atlas.Dirty {
		return
	}
	atlas.Dirty = false

	gl.Enable(gl.TEXTURE_2D)

	if atlas.Texture != 0 {
		gl.DeleteTextures(1, &atlas.Texture)
		atlas.Texture = 0
	}

	gl.GenTextures(1, &atlas.Texture)
	gl.BindTexture(gl.TEXTURE_2D, atlas.Texture)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(atlas.Image.Rect.Size().X),
		int32(atlas.Image.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(atlas.Image.Pix))

	if err := gl.GetError(); err != 0 {
		log.Println(err)
	}

	gl.Disable(gl.TEXTURE_2D)
}

func (atlas *Atlas) Draw(x, y float32, text string) {
	atlas.LoadGlyphs(text)

	gl.Enable(gl.BLEND)
	defer gl.Disable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.Enable(gl.TEXTURE_2D)
	defer gl.Disable(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, atlas.Texture)

	for _, r := range text {
		glyph := atlas.Rendered[r]

		dx, dy := float32(glyph.Loc.Dx()), float32(glyph.Loc.Dy())
		px, py := x+ceilPxf(glyph.Bounds.Min.X), y+ceilPxf(glyph.Bounds.Min.Y)
		gl.Begin(gl.QUADS)
		{
			gl.TexCoord2f(glyph.RelLoc.Min.X, glyph.RelLoc.Min.Y)
			gl.Vertex2f(px, py)
			gl.TexCoord2f(glyph.RelLoc.Max.X, glyph.RelLoc.Min.Y)
			gl.Vertex2f(px+dx, py)
			gl.TexCoord2f(glyph.RelLoc.Max.X, glyph.RelLoc.Max.Y)
			gl.Vertex2f(px+dx, py+dy)
			gl.TexCoord2f(glyph.RelLoc.Min.X, glyph.RelLoc.Max.Y)
			gl.Vertex2f(px, py+dy)
		}
		gl.End()

		x += ceilPxf(glyph.Advance)
	}
}

func (atlas *Atlas) Drawf(x, y float32, format string, args ...interface{}) {
	atlas.Draw(x, y, fmt.Sprintf(format, args...))
}

func drawRect(rgba *image.RGBA, bounds image.Rectangle, color color.RGBA) {
	for x := bounds.Min.X; x <= bounds.Max.X; x++ {
		rgba.SetRGBA(x, bounds.Min.Y, color)
		rgba.SetRGBA(x, bounds.Max.Y, color)
	}
	for y := bounds.Min.Y; y <= bounds.Max.Y; y++ {
		rgba.SetRGBA(bounds.Min.X, y, color)
		rgba.SetRGBA(bounds.Max.X, y, color)
	}
}

func drawVertLine(rgba *image.RGBA, x int, bounds image.Rectangle, color color.RGBA) {
	for y := bounds.Min.Y; y <= bounds.Max.Y; y++ {
		rgba.SetRGBA(x, y, color)
	}
}

func (atlas *Atlas) DrawDebug() {
	red := color.RGBA{0xFF, 0x00, 0x00, 0xFF}
	green := color.RGBA{0x00, 0xFF, 0x00, 0xFF}
	for _, glyph := range atlas.Rendered {
		drawRect(atlas.Image, glyph.Loc, red)
		drawVertLine(
			atlas.Image,
			glyph.Loc.Min.X-ceilPx(glyph.Bounds.Min.X)+ceilPx(glyph.Advance),
			glyph.Loc,
			green)
	}
}
