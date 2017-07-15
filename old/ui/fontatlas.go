package ui

import (
	"image"
	"io/ioutil"
	"log"
	"math"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/go-gl/gl/v2.1/gl"
)

func RelBounds(r, b image.Rectangle) (n Bounds) {
	n.Min.X = float32(r.Min.X-b.Min.X) / float32(b.Dx())
	n.Min.Y = float32(r.Min.Y-b.Min.Y) / float32(b.Dy())
	n.Max.X = float32(r.Max.X-b.Min.X) / float32(b.Dx())
	n.Max.Y = float32(r.Max.Y-b.Min.Y) / float32(b.Dy())
	return n
}

const (
	glyphMargin  = 2
	glyphPadding = 1
)

type Glyph struct {
	Rune    rune
	Loc     image.Rectangle     // absolute location on image atlas
	RelLoc  Bounds              // relative location on image atlas
	Bounds  fixed.Rectangle26_6 // such that point + bounds, gives image bounds where glyph should be drawn
	Advance fixed.Int26_6       // advance from point, to the next glyph
}

type FontAtlas struct {
	Context *freetype.Context
	TTF     *truetype.Font
	Face    font.Face

	Rendered      map[rune]Glyph
	Image         *image.RGBA
	CursorX       int
	CursorY       int
	maxGlyphInRow int
	drawPadding   float32

	maxBounds  fixed.Rectangle26_6
	lineHeight float32

	Dirty   bool
	Texture uint32
}

func NewFontAtlas(filename string, dpi, fontSize float64) (*FontAtlas, error) {
	atlas := &FontAtlas{}
	atlas.Rendered = make(map[rune]Glyph, 256)

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	atlas.drawPadding = float32(fontSize * 0.5)
	atlas.lineHeight = float32(fontSize * 1.2)

	atlas.TTF, err = truetype.Parse(content)
	if err != nil {
		return nil, err
	}

	atlas.Image = image.NewRGBA(image.Rect(0, 0, 1024, 1024))

	atlas.Context = freetype.NewContext()
	atlas.Context.SetDPI(dpi)

	atlas.Context.SetFont(atlas.TTF)
	atlas.Context.SetFontSize(fontSize)

	atlas.Context.SetClip(atlas.Image.Bounds())
	atlas.Context.SetSrc(image.White)
	atlas.Context.SetDst(atlas.Image)

	atlas.maxBounds = atlas.TTF.Bounds(fixed.I(int(fontSize)))

	opts := &truetype.Options{}
	opts.Size = fontSize
	opts.Hinting = font.HintingFull

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

func (atlas *FontAtlas) loadGlyph(r rune) {
	if _, ok := atlas.Rendered[r]; ok {
		return
	}
	atlas.Dirty = true

	glyph := Glyph{}
	glyph.Rune = r

	bounds, advance, _ := atlas.Face.GlyphBounds(r)
	glyph.Bounds = bounds
	glyph.Advance = advance

	width := ceilPx(bounds.Max.X-bounds.Min.X) + glyphPadding*2
	height := ceilPx(bounds.Max.Y-bounds.Min.Y) + glyphPadding*2

	if atlas.CursorX+glyphMargin+width+glyphMargin > atlas.Image.Bounds().Dx() {
		atlas.CursorX = 0
		atlas.CursorY += glyphMargin + atlas.maxGlyphInRow
	}

	x := atlas.CursorX + glyphMargin
	y := atlas.CursorY + glyphMargin

	glyph.Loc = image.Rect(x, y, x+width, y+height)
	glyph.RelLoc = RelBounds(glyph.Loc, atlas.Image.Bounds())

	pt := fixed.P(x+glyphPadding, y+glyphPadding).Sub(bounds.Min)
	atlas.Context.DrawString(string(r), pt)

	if height > atlas.maxGlyphInRow {
		atlas.maxGlyphInRow = height
	}
	atlas.CursorX += glyphMargin + width + glyphMargin

	atlas.Rendered[r] = glyph
}

func (atlas *FontAtlas) LoadAscii() {
	for r := rune(0); r < 128; r++ {
		atlas.loadGlyph(r)
	}
	atlas.upload()
}

func (atlas *FontAtlas) LoadExtendedAscii() {
	for r := rune(0); r < 256; r++ {
		atlas.loadGlyph(r)
	}
	atlas.upload()
}

func (atlas *FontAtlas) LoadGlyphs(text string) {
	for _, r := range text {
		atlas.loadGlyph(r)
	}
	atlas.upload()
}

func (atlas *FontAtlas) upload() {
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

func (atlas *FontAtlas) Draw(text string, b Bounds) {
	atlas.LoadGlyphs(text)

	gl.Enable(gl.BLEND)
	defer gl.Disable(gl.BLEND)
	gl.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA)

	gl.Enable(gl.TEXTURE_2D)
	defer gl.Disable(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, atlas.Texture)

	x := b.Min.X + atlas.drawPadding
	y := (b.Max.Y+b.Min.Y)/2 + (ceilPxf(atlas.maxBounds.Min.Y)+ceilPxf(atlas.maxBounds.Max.Y))/2

	p := rune(0)
	for _, r := range text {
		glyph := atlas.Rendered[r]

		dx := float32(glyph.Loc.Dx())
		dy := float32(glyph.Loc.Dy())

		px := x + ceilPxf(glyph.Bounds.Min.X) - glyphPadding
		py := y + ceilPxf(glyph.Bounds.Min.Y) - glyphPadding

		// this is not the ideal way of positioning the letters
		// will create positioning artifacts
		// but it the result is more
		px = float32(math.Trunc(float64(px)))
		py = float32(math.Trunc(float64(py)))

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

		k := atlas.Face.Kern(p, r)
		p = r
		x += ceilPxf(glyph.Advance + k)
	}
}

func (atlas *FontAtlas) Measure(text string) (size Point) {
	atlas.LoadGlyphs(text)

	size.X += atlas.drawPadding * 2
	size.Y += atlas.lineHeight

	p := rune(0)
	for _, r := range text {
		glyph := atlas.Rendered[r]
		k := atlas.Face.Kern(p, r)
		p = r
		size.X += float32(ceilPx(glyph.Advance + k))
	}
	return
}
