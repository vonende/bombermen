package titlebars

import (
	"testing"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/opengl"
)

func run() {
	wincfg := opengl.WindowConfig{
		Title:  "Titlebar Test",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}

	win, err := opengl.NewWindow(wincfg)
	if err != nil {
		panic(err)
	}
	defer win.Destroy()

	win.SetMatrix(pixel.IM.Moved(win.Bounds().Center()).Scaled(win.Bounds().Center(), 3))

	bar := New(16 * 16)
	bar.SetPlayers(4)
	var playerOneLifes uint8 = 3
	var playerTwoLifes uint8 = 1
	var playerThreeLifes uint8 = 0
	var playerFourLifes uint8 = 5
	bar.SetLifePointers(&playerOneLifes, &playerTwoLifes, &playerThreeLifes, &playerFourLifes)
	var points uint32 = 1243
	bar.SetPointsPointer(&points)
	go bar.Manager()
	bar.SetSeconds(200)
	bar.StartCountdown()
	win.Update()
	for {
		bar.Draw(win)
		win.Update()
		if win.Pressed(opengl.KeyEscape) {
			break
		}
	}
}

func TestMain(*testing.M) {
	opengl.Run(run)
}
