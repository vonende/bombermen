package gameStats

import (
	"../animations"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	. "github.com/vonende/bombermen/constants"
	"github.com/vonende/bombermen/gameStats"
	"github.com/vonende/bombermen/levels"
	"testing"
)

func run() {
	levelDef := levels.NewLevel("./levels/stufe_2_level_3.txt")
	pitchWidth, pitchHeight := levelDef.GetBounds()
	var zoomFactor float64
	if float64((pitchHeight+1)*TileSize+32)/float64((pitchWidth+3)*TileSize) > float64(MaxWinSizeY)/MaxWinSizeX {
		zoomFactor = MaxWinSizeY / float64((pitchHeight+1)*TileSize+32)
	} else {
		zoomFactor = MaxWinSizeX / float64((pitchWidth+3)*TileSize)
	}
	var winSizeX = zoomFactor * float64(pitchWidth+3) * TileSize
	var winSizeY = zoomFactor * (float64(pitchHeight+1)*TileSize + 32)
	//var err error

	wincfg := pixelgl.WindowConfig{
		Title:  "GameStat Test",
		Bounds: pixel.R(0, 0, winSizeX, winSizeY),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(wincfg)
	if err != nil {
		panic(err)
	}

	win.SetMatrix(pixel.IM.Scaled(pixel.V(0, 0), zoomFactor))

	win.Update()

	lv := gameStats.NewGameStat(levelDef, 1)

	fmt.Println(lv.GetBounds())

	itemBatch := pixel.NewBatch(&pixel.TrianglesData{}, animations.ItemImage)

	for !win.Closed() && !win.Pressed(pixelgl.KeyEscape) {
		itemBatch.Clear()
		if win.JustPressed(pixelgl.MouseButton1) { // Destroy
			xx, yy := lv.A().GetFieldCoord(win.MousePosition().Scaled(1 / zoomFactor))
			fmt.Println(xx, yy)
			if !lv.IsTile(xx, yy) {
				lv.RemoveItems(xx, yy, pixel.V(0, 0))
			} else if lv.IsDestroyableTile(xx, yy) {
				lv.RemoveTile(xx, yy)
			}

		}
		if win.JustPressed(pixelgl.MouseButton2) { // Info & Collect
			xx, yy := lv.A().GetFieldCoord(win.MousePosition().Scaled(1 / zoomFactor))
			fmt.Println("Ist zerstörbares Teil: ", lv.IsDestroyableTile(xx, yy))
			fmt.Println("Ist unzerstörbarTeil: ", lv.IsUndestroyableTile(xx, yy))
			fmt.Println("Ist Teil: ", lv.IsTile(xx, yy))
			if !lv.IsTile(xx, yy) {
				a, b, c := lv.GetPosOfNextTile(xx, yy, pixel.V(0, 1))
				d, e := lv.CollectItem(xx, yy)
				fmt.Println("Nächstes Teil in Richtung Up: ", a, b, c)
				fmt.Println("Item gesammelt?: ", d, e)
			}
		}
		if win.JustPressed(pixelgl.KeySpace) {
			lv.Reset()
		}

		for i := 0; i < pitchHeight; i++ {
			lv.DrawColumn(i, itemBatch)
		}
		lv.A().GetCanvas().Draw(win, *(lv.A().GetMatrix()))
		itemBatch.Draw(win)
		win.Update()
	}

}

func TestMain(*testing.M) {
	pixelgl.Run(run)
}
