package text

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"image/color"
	"image/png"
	"os"
)

/****************************************************************************************

Simples Paket zur Ausgabe von Retro-Text über ein Canvas-Element.

_________________________________
< Implementiert von Rayk von Ende >
 ---------------------------------
        \   ^__^
         \  (oo)\_______
            (__)\       )\/\
                ||----w |
                ||     ||

****************************************************************************************/

const Fire = 1

type textStruct struct {
	maxLength int
	font      *pixel.PictureData
	canvas    *pixelgl.Canvas
	sprite    *pixel.Sprite
	text      []byte
}

func NewFont(t uint8) *pixel.PictureData {
	switch t {
	case Fire:
		return loadFont("./data/graphics/firefont.png")
	default:
		return loadFont("./data/graphics/firefont.png")
	}
}

func NewTextL(font *pixel.PictureData, maxLength int) *textStruct {
	ts := new(textStruct)
	ts.maxLength = maxLength
	if font == nil {
		panic("Es wurde kein Font-Image gefunden (Nil-Pointer)!")
	}
	ts.font = font
	ts.text = make([]byte, maxLength)
	ts.canvas = pixelgl.NewCanvas(pixel.R(0, 0, float64(16*maxLength), 32))
	ts.sprite = pixel.NewSprite(font, pixel.R(0, 0, 16, 32))
	return ts
}

func NewText(font *pixel.PictureData, s string) *textStruct {
	ts := new(textStruct)
	ts.text = stringToBytes(s)
	ts.maxLength = len(ts.text)
	if font == nil {
		panic("Es wurde kein Font-Image gefunden (Nil-Pointer)!")
	}
	ts.font = font
	ts.canvas = pixelgl.NewCanvas(pixel.R(0, 0, float64(16*ts.maxLength), 32))
	ts.sprite = pixel.NewSprite(font, pixel.R(0, 0, 16, 32))
	ts.canvas.Clear(color.Transparent)
	for i, val := range ts.text {
		ts.sprite.Set(ts.font, getRectForChar(val))
		ts.sprite.Draw(ts.canvas, pixel.IM.Moved(pixel.V(8, 16)).Moved(pixel.V(float64(i*16), 0)))
	}
	return ts
}

func (ts *textStruct) Draw(target pixel.Target, matrix pixel.Matrix) {
	ts.canvas.Draw(target, pixel.IM.Moved(ts.canvas.Bounds().Center().Sub(ts.Bounds().Center())).Chained(matrix))
}

func (ts *textStruct) Bounds() pixel.Rect {
	return pixel.R(0, 0, float64(16*len(ts.text)), 32)
}

func stringToBytes(s string) []byte {
	b := make([]byte, 0)
	for _, val := range s {
		// Nicht-ASCII-Zeichen werden ausgeschlossen
		if val < 0x20 || val > 0x7E {
			val = ' '
		}
		b = append(b, byte(val))
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (ts *textStruct) Set(text string) {
	b := stringToBytes(text)
	ts.text = b[0:min(len(b), ts.maxLength)]
	ts.canvas.Clear(color.Transparent)
	for i, val := range ts.text {
		ts.sprite.Set(ts.font, getRectForChar(val))
		ts.sprite.Draw(ts.canvas, pixel.IM.Moved(pixel.V(8, 16)).Moved(pixel.V(float64(i*16), 0)))
	}
}

// Vor.: -
// Eff.: getRectForChar liefert zu einem Zeichen das zugehörige Rechteck innerhalb des font.png
func getRectForChar(c byte) pixel.Rect {
	x := float64(c) - 32
	if x < 0 || x > 96 {
		return pixel.R(0, 0, 16, 32)
	}
	x = x * 16
	return pixel.R(x, 0, x+16, 32)
}

func loadFont(s string) (fontImage *pixel.PictureData) {
	file, err := os.Open(s)
	if err != nil {
		panic(err)
	}
	img, err := png.Decode(file)
	if err != nil {
		panic(err)
	}
	err = file.Close()
	if err != nil {
		panic(err)
	}

	fontImage = pixel.PictureDataFromImage(img)
	return fontImage
}
