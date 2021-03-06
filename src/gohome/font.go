package gohome

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/math/fixed"
	"image"
)

const (
	DPI = 72
)

// A struct representing a truetype font
type Font struct {
	ttf         *truetype.Font
	drawContext *freetype.Context
	// The size in which the font should be drawn
	FontSize    int
}

// Draw a string and returns it as a Texture
func (this *Font) DrawString(str string) Texture {
	var rv Texture

	imgWidth, imgHeight := this.getTextureSize(str)

	rect := image.Rect(0, 0, int(imgWidth), int(imgHeight))
	img := image.NewRGBA(rect)

	this.drawContext.SetFontSize(float64(this.FontSize))
	this.drawContext.SetClip(img.Bounds())
	this.drawContext.SetDst(img)
	this.drawContext.DrawString(str, freetype.Pt(0, int(imgHeight*7/9)))

	rv = Render.CreateTexture("Text "+str, false)
	rv.LoadFromImage(img)

	return rv
}

// Initialises the font with a truetype font object
func (this *Font) Init(ttf *truetype.Font) {
	this.ttf = ttf
	this.drawContext = freetype.NewContext()
	this.drawContext.SetDPI(DPI)
	this.drawContext.SetFont(this.ttf)
	this.drawContext.SetSrc(image.White)
	this.drawContext.SetFontSize(float64(this.FontSize))
}

func (this *Font) getTextureSize(str string) (width, height int) {
	runeString := []rune(str)

	scale := fixed.Int26_6(this.FontSize)
	bounds := this.ttf.Bounds(scale)
	height = int(bounds.Max.Y - bounds.Min.Y)

	for i := 0; i < len(runeString); i++ {
		hmetric := this.ttf.HMetric(scale, this.ttf.Index(runeString[i]))
		if i == 0 {
			width += int(fixed.Int26_6(hmetric.LeftSideBearing) + fixed.Int26_6(hmetric.AdvanceWidth))
		}
		width += int(fixed.Int26_6(hmetric.AdvanceWidth))
	}

	return width, height
}

// Returns the maximum height of a character
func (this *Font) GetGlyphMaxHeight() int {
	rect := this.ttf.Bounds(fixed.Int26_6(this.FontSize))
	return int(rect.Max.Y - rect.Min.Y)
}

// Returns the maximum width of a character
func (this *Font) GetGlyphMaxWidth() int {
	rect := this.ttf.Bounds(fixed.Int26_6(this.FontSize))
	return int(rect.Max.X - rect.Min.X)
}
