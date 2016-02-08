package font

import (
	"image"
	"image/color"
	"io/ioutil"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type Glyph struct {
	Rune    rune
	Bounds  image.Rectangle // location on atlas
	Advance int
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

func (atlas *Atlas) LoadGlyph(r rune) {
	if _, ok := atlas.Rendered[r]; ok {
		return
	}

	glyph := Glyph{}
	glyph.Rune = r

	bounds, advance, _ := atlas.Face.GlyphBounds(r)
	glyph.Advance = int((advance + 63) >> 6)

	width := int(bounds.Max.X-bounds.Min.X) >> 6
	height := int(bounds.Max.Y-bounds.Min.Y) >> 6

	if atlas.CursorX+atlas.Padding+width+atlas.Padding > atlas.Image.Bounds().Dx() {
		atlas.CursorX = 0
		atlas.CursorY += atlas.Padding + atlas.MaxGlyphHeight
	}

	x := atlas.CursorX + atlas.Padding
	y := atlas.CursorY + atlas.Padding

	glyph.Bounds = image.Rect(x, y, x+width, y+height)

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
		atlas.LoadGlyph(r)
	}
}

func (atlas *Atlas) LoadExtendedAscii() {
	for r := rune(0); r < 256; r++ {
		atlas.LoadGlyph(r)
	}
}

func (atlas *Atlas) LoadGlyphs(text string) {
	for _, r := range text {
		atlas.LoadGlyph(r)
	}
}

func (atlas *Atlas) Draw(text string) {
	atlas.LoadGlyphs(text)
}

func DrawRect(rgba *image.RGBA, bounds image.Rectangle, color color.RGBA) {
	for x := bounds.Min.X; x <= bounds.Max.X; x++ {
		rgba.SetRGBA(x, bounds.Min.Y, color)
		rgba.SetRGBA(x, bounds.Max.Y, color)
	}
	for y := bounds.Min.Y; y <= bounds.Max.Y; y++ {
		rgba.SetRGBA(bounds.Min.X, y, color)
		rgba.SetRGBA(bounds.Max.X, y, color)
	}
}

func DrawVertLine(rgba *image.RGBA, x int, bounds image.Rectangle, color color.RGBA) {
	for y := bounds.Min.Y; y <= bounds.Max.Y; y++ {
		rgba.SetRGBA(x, y, color)
	}
}

func (atlas *Atlas) DrawDebug() {
	red := color.RGBA{0xFF, 0x00, 0x00, 0xFF}
	green := color.RGBA{0x00, 0xFF, 0x00, 0xFF}
	for _, glyph := range atlas.Rendered {
		DrawRect(atlas.Image, glyph.Bounds, red)
		DrawVertLine(atlas.Image, glyph.Bounds.Min.X+glyph.Advance, glyph.Bounds, green)
	}
}
