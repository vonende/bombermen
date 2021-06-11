package text

import (
	"github.com/faiface/pixel"
	"image/png"
	"os"
)

/****************************************************************************************

Simples Paket zur Ausgabe von Retro-Text.

_________________________________
< Implementiert von Rayk von Ende >
 ---------------------------------
        \   ^__^
         \  (oo)\_______
            (__)\       )\/\
                ||----w |
                ||     ||

****************************************************************************************/

const (
	Fire = 1
	Pink = 2
	Blue = 3
)

type textStruct struct {
	font   *pixel.PictureData
	batch  *pixel.Batch
	sprite *pixel.Sprite
	text   []byte
}

func NewFont(t uint8) *pixel.PictureData {
	switch t {
	case Fire:
		return loadFont("./data/graphics/firefont.png")
	case Pink:
		return loadFont("./data/graphics/pinkfont.png")
	case Blue:
		return loadFont("./data/graphics/bluefont.png")
	default:
		return loadFont("./data/graphics/firefont.png")
	}
}

func NewText(font *pixel.PictureData, s string) *textStruct {
	ts := new(textStruct)
	ts.text = stringToBytes(s)
	if font == nil {
		panic("Es wurde kein Font-Image gefunden (Nil-Pointer)!")
	}
	ts.font = font
	ts.sprite = pixel.NewSprite(font, pixel.R(0, 0, 16, 32))
	ts.batch = pixel.NewBatch(&pixel.TrianglesData{}, font)
	return ts
}

func (ts *textStruct) Draw(target pixel.Target, matrix pixel.Matrix) {
	ts.batch.Clear()
	for i, val := range ts.text {
		ts.sprite.Set(ts.font, getRectForChar(val))
		ts.sprite.Draw(ts.batch, pixel.IM.Moved(ts.Bounds().Center().Scaled(-1)).Moved(pixel.V(8, 16)).Moved(pixel.V(float64(i*16), 0)).Chained(matrix))
	}
	ts.batch.Draw(target)
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

func (ts *textStruct) Set(text string) {
	ts.text = stringToBytes(text)
	ts.batch.Clear()
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
