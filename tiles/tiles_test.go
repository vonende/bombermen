package tiles

import (
	"fmt"
	"testing"
	"time"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/opengl"
	"github.com/vonende/bombermen/animations"
	"github.com/vonende/bombermen/characters"
	. "github.com/vonende/bombermen/constants"
	"golang.org/x/image/colornames"
)

func run() {

	wincfg := opengl.WindowConfig{
		Title:  "GameStat Test",
		Bounds: pixel.R(-100, -100, 400, 400),
		VSync:  true,
	}
	win, err := opengl.NewWindow(wincfg)
	if err != nil {
		panic(err)
	}
	mod := 1
	win.SetMatrix(pixel.IM) //.Scaled(pixel.V(0, 0), 3))
	win.Clear(colornames.Blue)
	win.Update()

	wB := characters.NewPlayer(WhiteBomberman)

	itemBatch := pixel.NewBatch(&pixel.TrianglesData{}, animations.ItemImage)
	var ti []Tile
	var it []Item
	var bs []Bombe

	i := NewItem(SkullItem, pixel.V(0, 0))
	fmt.Println(i.GetType())
	for !win.Closed() && !win.Pressed(opengl.KeyEscape) {
		itemBatch.Clear()
		if mod%2 == 0 {
			for _, w := range ti {
				w.SetVisible(true)
			}
			for _, w := range it {
				w.SetVisible(false)
				w.SetDestroyable(false)
				w.SetTimeStamp(time.Now())
			}
		} else {
			for _, w := range ti {
				w.SetVisible(false)
			}
			for _, w := range it {
				w.SetVisible(true)
				w.SetDestroyable(true)
			}
		}
		if win.JustPressed(opengl.MouseButton1) { // Move Tile
			t := NewTile(House, win.MousePosition())
			ti = append(ti, t)
			fmt.Println("Zeichenmatrix: ", t.GetMatrix())
			fmt.Println("Position: ", t.GetPos())
			fmt.Println("Ist sichtbar: ", t.IsVisible())
			fmt.Println(t.GetType())
			mod++
		} else if win.JustPressed(opengl.MouseButton2) { // Move Item
			i := NewItem(SkullItem, win.MousePosition()) //.Scaled(1/3))
			it = append(it, i)
			mod++
		}

		if win.JustPressed(opengl.KeySpace) { // stats
			for _, w := range ti {
				fmt.Println("Tile Animation sichtbar? ", w.IsVisible())
			}
			for _, w := range it {
				fmt.Println("Item Animation sichtbar? ", w.IsVisible())
				fmt.Println("Item zerstörbar? ", w.IsDestroyable())
				fmt.Println("Item TimeStamp: ", w.GetTimeStamp())
			}
			for _, w := range bs {
				fmt.Println("Bombenpower: ", w.GetPower())
				bb, cc := w.Owner()
				fmt.Println("Bombenpower: ", bb, cc)
			}
		}

		if win.JustPressed(opengl.KeyB) {
			b := NewBomb(wB, win.MousePosition())
			b.SetPower(1)
			bs = append(bs, b)
		}

		for _, w := range ti {
			w.Draw(itemBatch)
		}
		for _, w := range it {
			w.Draw(itemBatch)
		}
		for _, w := range bs {
			w.Draw(itemBatch)
		}
		//win.Clear(colornames.Blue)
		itemBatch.Draw(win)
		win.Update()
	}

}

func TestMain(*testing.M) {
	opengl.Run(run)
}
