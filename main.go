package main

import (
	"./animations"
	"./arena"
	"./characters"
	. "./constants"
	"./level1"
	"./sounds"
	"./tiles"
	"./titlebar"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	//"golang.org/x/image/colornames"
	"math"
	"math/rand"
	"time"
)

var bombs []tiles.Bombe
var turfNtreesArena arena.Arena
var lev1 level1.Level
var tempAniSlice [][]interface{} // [Animation][Matrix]
var monster []characters.Enemy
var whiteBomberman characters.Player

// Vor: ...
// Eff: Ist der Counddown der Bombe abgelaufen passiert folgendes:
//     		- Eine neue Explosionsanimation ist erstellt und an die Position der ehemaligen bombe gesetzt.
//      	- Es ertönt der Explosionssound.
//      Ist der Countdown nicht abgelaufen, passiert nichts.

func checkForExplosions() {

	for _, item := range bombs {
		if ((item).GetTimeStamp()).Before(time.Now()) {
			//bomben = append (bomben,item)

			item.Ani().Die()

			x, y := turfNtreesArena.GetFieldCoord(item.GetPos())
			power := int(item.GetPower())
			l, r, u, d := power, power, power, power

			// Explosion darf nicht über Spielfeldrand hinausragen
			if l > x {
				l = x
			}
			if turfNtreesArena.GetWidth()-1-r < x {
				r = turfNtreesArena.GetWidth() - 1 - x
			}
			if d > y {
				d = y
			}
			if turfNtreesArena.GetHeight()-1-u < y {
				u = turfNtreesArena.GetHeight() - 1 - y
			}

			// Falls es Hindernisse gibt, die zerstörbar oder unzerstörbar sind
			bl, xl, yl := lev1.GetPosOfNextTile(x, y, pixel.V(float64(-l), 0))
			if bl {
				if lev1.IsDestroyableTile(xl, yl) {
					l = x - xl
				} else {
					l = x - xl - 1
				}
			}
			br, xr, yr := lev1.GetPosOfNextTile(x, y, pixel.V(float64(r), 0))
			if br {
				if lev1.IsDestroyableTile(xr, yr) {
					r = xr - x
				} else {
					r = xr - x - 1
				}
			}
			bd, xd, yd := lev1.GetPosOfNextTile(x, y, pixel.V(0, float64(-d)))
			if bd {
				if lev1.IsDestroyableTile(xd, yd) {
					d = y - yd
				} else {
					d = y - yd - 1
				}
			}
			bu, xu, yu := lev1.GetPosOfNextTile(x, y, pixel.V(0, float64(u)))
			if bu {
				if lev1.IsDestroyableTile(xu, yu) {
					u = yu - y
				} else {
					u = yu - y - 1
				}
			}

			// falls sich ein Monster oder Player im Explosionsradius befindet

		A:
			for i := 1; i <= l; i++ {
				for _, m := range monster {
					xx, yy := turfNtreesArena.GetFieldCoord(m.GetPos())
					if x-i == xx && y == yy {
						l = i
						m.Ani().Die()
						break A
					}
				}
			}

		B:
			for i := 1; i <= r; i++ {
				for _, m := range monster {
					xx, yy := turfNtreesArena.GetFieldCoord(m.GetPos())
					if x+i == xx && y == yy {
						r = i
						m.Ani().Die()
						break B
					}
				}
			}

		C:
			for i := 1; i <= u; i++ {
				for _, m := range monster {
					xx, yy := turfNtreesArena.GetFieldCoord(m.GetPos())
					if y+i == yy && x == xx {
						u = i
						m.Ani().Die()
						break C
					}
				}
			}

		D:
			for i := 1; i <= d; i++ {
				for _, m := range monster {
					xx, yy := turfNtreesArena.GetFieldCoord(m.GetPos())
					if y-i == yy && x == xx {
						d = i
						m.Ani().Die()
						break D
					}
				}
			}

			if xl+l == x {
				lev1.RemoveTile(xl, yl)
			}
			if xr-r == x {
				lev1.RemoveTile(xr, yr)
			}
			if yd+d == y {
				lev1.RemoveTile(xd, yd)
			}
			if yu-u == y {
				lev1.RemoveTile(xu, yu)
			}

			// Items, die im Expolsionsradius liegen werden zerstört, die Expolion wird aber nicht kleiner!

			lev1.RemoveItems(x, y, pixel.V(float64(-l), 0))
			lev1.RemoveItems(x, y, pixel.V(float64(r), 0))
			lev1.RemoveItems(x, y, pixel.V(0, float64(-d)))
			lev1.RemoveItems(x, y, pixel.V(0, float64(u)))

			// falls weitere Bomben im Explosionsradius liegen, werden auch gleich explodieren

			for i := 1; i <= l; i++ {
				b, bom := isThereABomb(item.GetPos().Add(pixel.V(float64(-i)*TileSize, 0)))
				if b {
					bom.SetTimeStamp(time.Now())
				}
			}
			for i := 1; i <= r; i++ {
				b, bom := isThereABomb(item.GetPos().Add(pixel.V(float64(i)*TileSize, 0)))
				if b {
					bom.SetTimeStamp(time.Now())
				}
			}
			for i := 1; i <= u; i++ {
				b, bom := isThereABomb(item.GetPos().Add(pixel.V(0, float64(i)*TileSize)))
				if b {
					bom.SetTimeStamp(time.Now())
				}
			}
			for i := 1; i <= d; i++ {
				b, bom := isThereABomb(item.GetPos().Add(pixel.V(0, float64(-i)*TileSize)))
				if b {
					bom.SetTimeStamp(time.Now())
				}
			}

			ani := animations.NewExplosion(uint8(l), uint8(r), uint8(u), uint8(d))
			ani.Show()
			tempAni := make([]interface{}, 2)
			tempAni[0] = ani
			tempAni[1] = (item.GetMatrix()).Moved(ani.ToCenter())
			tempAniSlice = append(tempAniSlice, tempAni)
			s2 := sounds.NewSound(Deathflash)
			go s2.PlaySound()
		}
	}

}

// Vor.:...
// Eff.: Nicht explodierte Bomben aus dem Slice existingBombs werden in den
//       Slice remainingBombs kopiert
func removeExplodedBombs(existingBombs []tiles.Bombe) (remainingBombs []tiles.Bombe) {
	j := 0
	for i, bomb := range existingBombs {
		if !bomb.IsVisible() {
			remainingBombs = append(remainingBombs, existingBombs[j:i]...)
			j = i + 1
		}
	}
	remainingBombs = append(remainingBombs, existingBombs[j:]...)
	return remainingBombs
}

func showExplosions(win *pixelgl.Window) {
	for _, a := range tempAniSlice {
		ani := (a[0]).(animations.Animation)
		ani.Update()
		ani.GetSprite().Draw(win, (a[1]).(pixel.Matrix))
	}
}

func clearExplosions(existingExplosions [][]interface{}) (remainingExplosions [][]interface{}) {
	for _, exp := range existingExplosions {
		if exp[0].(animations.Animation).IsVisible() {
			remainingExplosions = append(remainingExplosions, exp)
		}
	}
	return remainingExplosions
}

func isThereABomb(v pixel.Vec) (bool, tiles.Bombe) {
	for _, item := range bombs {
		if item.GetPos() == v {
			return true, item
		}
	}
	return false, nil
}

func getGrantedDirections(c characters.Character) [4]bool {
	var b [4]bool
	b[0] = true
	b[1] = true
	b[2] = true
	b[3] = true
	pb := c.GetPosBox()
	ll := pb.Min.Sub(turfNtreesArena.GetLowerLeft())
	ur := pb.Max.Sub(turfNtreesArena.GetLowerLeft())
	if lev1.IsTile(int((ll.X-1)/TileSize), int(ll.Y/TileSize)) || lev1.IsTile(int((ll.X-1)/TileSize), int(ur.Y/TileSize)) || ll.X-1 < 0 {
		b[0] = false
	}
	if int((ur.X+1)/TileSize) > turfNtreesArena.GetWidth()-1 {
		b[1] = false
	} else if lev1.IsTile(int((ur.X+1)/TileSize), int(ll.Y/TileSize)) || lev1.IsTile(int((ur.X+1)/TileSize), int(ur.Y/TileSize)) {
		b[1] = false
	}
	if int((ur.Y+1)/TileSize) > turfNtreesArena.GetHeight()-1 {
		b[2] = false
	} else if lev1.IsTile(int(ll.X/TileSize), int((ur.Y+1)/TileSize)) || lev1.IsTile(int(ur.X/TileSize), int((ur.Y+1)/TileSize)) {
		b[2] = false
	}
	if lev1.IsTile(int(ll.X/TileSize), int((ll.Y-1)/TileSize)) || lev1.IsTile(int(ur.X/TileSize), int((ll.Y-1)/TileSize)) || ll.Y-1 < 0 {
		b[3] = false
	}

	return b
}

func moveCharacter(c characters.Character, dt float64, dir uint8) {
	switch dir {
	case Left:
		dist := -c.GetSpeed() * dt
		if dist <= -TileSize {
			dist = -TileSize + 0.1
		}
		pb := c.GetPosBox()
		ll := pb.Min.Sub(turfNtreesArena.GetLowerLeft())
		ur := pb.Max.Sub(turfNtreesArena.GetLowerLeft())
		bl, xl, _ := lev1.GetPosOfNextTile(int(ll.X/TileSize), int(ll.Y/TileSize), pixel.V(-TileSize, 0))
		bu, xu, _ := lev1.GetPosOfNextTile(int(ll.X/TileSize), int(ur.Y/TileSize), pixel.V(-TileSize, 0))
		if bl || bu {
			if bl && (xl >= xu || xu == -1) {
				if ll.X+dist < float64((xl+1)*TileSize) {
					dist = float64((xl+1)*TileSize) - ll.X + 0.1
				}
			} else if bu && (xu >= xl || xl == -1) {
				if ll.X+dist < float64((xu+1)*TileSize) {
					dist = float64((xu+1)*TileSize) - ll.X + 0.1
				}
			}
		}
		c.Move(pixel.V(dist, 0))
	case Right:
		dist := c.GetSpeed() * dt
		if dist >= TileSize {
			dist = TileSize - 0.1
		}
		pb := c.GetPosBox()
		ll := pb.Min.Sub(turfNtreesArena.GetLowerLeft())
		ur := pb.Max.Sub(turfNtreesArena.GetLowerLeft())
		bl, xl, _ := lev1.GetPosOfNextTile(int((ur.X)/TileSize), int(ll.Y/TileSize), pixel.V(TileSize, 0))
		bu, xu, _ := lev1.GetPosOfNextTile(int((ur.X)/TileSize), int(ur.Y/TileSize), pixel.V(TileSize, 0))
		if bl || bu {
			if bl && (xl <= xu || xu == -1) {
				if ur.X+dist > float64((xl)*TileSize) {
					dist = float64((xl)*TileSize) - ur.X - 0.1
				}
			} else if bu && (xu <= xl || xl == -1) {
				if ur.X+dist > float64((xu)*TileSize) {
					dist = float64((xu)*TileSize) - ur.X - 0.1
				}
			}
		}
		c.Move(pixel.V(dist, 0))
	case Up:
		dist := c.GetSpeed() * dt
		if dist >= TileSize {
			dist = TileSize - 0.1
		}
		pb := c.GetPosBox()
		ll := pb.Min.Sub(turfNtreesArena.GetLowerLeft())
		ur := pb.Max.Sub(turfNtreesArena.GetLowerLeft())
		bl, _, yl := lev1.GetPosOfNextTile(int((ll.X)/TileSize), int((ur.Y)/TileSize), pixel.V(0, TileSize))
		br, _, yr := lev1.GetPosOfNextTile(int((ur.X)/TileSize), int((ur.Y)/TileSize), pixel.V(0, TileSize))
		if bl || br {
			if bl && (yl <= yr || yr == -1) {
				if ur.Y+dist > float64((yl)*TileSize) {
					dist = float64((yl)*TileSize) - ur.Y - 0.1
				}
			} else if br && (yr <= yl || yl == -1) {
				if ur.Y+dist > float64((yr)*TileSize) {
					dist = float64((yr)*TileSize) - ur.Y - 0.1
				}
			}
		}
		c.Move(pixel.V(0, dist))
	case Down:
		dist := -c.GetSpeed() * dt
		if dist <= -TileSize {
			dist = -TileSize + 0.1
		}
		pb := c.GetPosBox()
		ll := pb.Min.Sub(turfNtreesArena.GetLowerLeft())
		ur := pb.Max.Sub(turfNtreesArena.GetLowerLeft())
		bl, _, yl := lev1.GetPosOfNextTile(int((ll.X)/TileSize), int((ll.Y)/TileSize), pixel.V(0, -TileSize))
		br, xr, yr := lev1.GetPosOfNextTile(int((ur.X)/TileSize), int((ll.Y)/TileSize), pixel.V(0, -TileSize))
		if bl || br {
			fmt.Println(br, xr, yr)
			if bl && (yl >= yr || yr == -1) {
				if ll.Y+dist < float64((yl+1)*TileSize) {
					dist = float64((yl+1)*TileSize) - ll.Y + 0.1
				}
			} else if br && (yr >= yl || yl == -1) {
				if ll.Y+dist < float64((yr+1)*TileSize) {
					dist = float64((yr+1)*TileSize) - ll.Y + 0.1
				}
			}
			fmt.Println(dist, ll.Y, float64((yr)*TileSize))
		}
		c.Move(pixel.V(0, dist))
	}
	c.Ani().SetView(dir)
}

func sun() {
	const zoomFactor = 3
	const typ = 2
	const pitchWidth = 13
	const pitchHeight = 11
	var winSizeX float64 = zoomFactor * ((3 + pitchWidth) * TileSize) // TileSize = 16
	var winSizeY float64 = zoomFactor * ((1+pitchHeight)*TileSize + 32)

	wincfg := pixelgl.WindowConfig{
		Title:  "Bomberman 2021",
		Bounds: pixel.R(0, 0, winSizeX, winSizeY),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(wincfg)
	if err != nil {
		panic(err)
	}

	turfNtreesArena = arena.NewArena(typ, pitchWidth, pitchHeight)

	lev1 = level1.NewBlankLevel(turfNtreesArena, 1)
	lev1.SetRandomTilesAndItems(20, 10)

	whiteBomberman = characters.NewPlayer(WhiteBomberman)
	whiteBomberman.Ani().Show()
	whiteBomberman.IncPower()
	whiteBomberman.IncPower()

	tb := titlebar.New((3 + pitchWidth) * TileSize)
	tb.SetMatrix(pixel.IM.Moved(pixel.V((3+pitchWidth)*TileSize/2, (1+pitchHeight)*TileSize+16)))
	tb.SetLifePointers(whiteBomberman.GetLifePointer())
	tb.SetPointsPointer(whiteBomberman.GetPointsPointer())
	tb.SetPlayers(1)
	go tb.Manager()
	tb.SetSeconds(5 * 60)
	tb.StartCountdown()
	tb.Update()

	// 2 Enemys

	monster = append(monster, characters.NewEnemy(YellowPopEye))
	monster = append(monster, characters.NewEnemy(Drop))

	// not a real random number

	rand.Seed(42)

	// Put character at free space with at least two free neighbours in a row
	/*
		A:
			for i := 2 * turfNtreesArena.GetWidth(); i < len(turfNtreesArena.GetPassability())-2*turfNtreesArena.GetWidth(); i++ { // Einschränkung des Wertebereichs von i um index out of range Probleme zu vermeiden
				if turfNtreesArena.GetPassability()[i] && turfNtreesArena.GetPassability()[i-1] && turfNtreesArena.GetPassability()[i-2] || // checke links, rechts, oben, unten
					turfNtreesArena.GetPassability()[i] && turfNtreesArena.GetPassability()[i+1] && turfNtreesArena.GetPassability()[i+2] ||
					turfNtreesArena.GetPassability()[i] && turfNtreesArena.GetPassability()[i+turfNtreesArena.GetWidth()] &&
						turfNtreesArena.GetPassability()[i+2*turfNtreesArena.GetWidth()] ||
					turfNtreesArena.GetPassability()[i] && turfNtreesArena.GetPassability()[i-turfNtreesArena.GetWidth()] &&
						turfNtreesArena.GetPassability()[i-2*turfNtreesArena.GetWidth()] {
					whiteBomberman.MoveTo(turfNtreesArena.GetLowerLeft().Add(pixel.V(float64(i%turfNtreesArena.GetWidth())*
						TileSize, float64(i/turfNtreesArena.GetWidth())*TileSize)))
					break A
				}
			}
	*/

	// Bomberman is in lowleft Corner
	whiteBomberman.MoveTo(turfNtreesArena.GetLowerLeft().Add(pixel.V(0, 0)))

	///////////////////////// ToDo Enyemys should be a Part of Level //////////////////////////////////////////////
	xx, yy := turfNtreesArena.GetFieldCoord(whiteBomberman.GetPos())

	for _, m := range monster {
		for {
			i := rand.Intn(turfNtreesArena.GetWidth())
			j := rand.Intn(turfNtreesArena.GetHeight())
			if !turfNtreesArena.IsTile(i, j) && xx != i && yy != j {
				m.MoveTo(turfNtreesArena.GetLowerLeft().Add(pixel.V(float64(i)*
					TileSize, float64(j)*TileSize)))
				m.Ani().SetVisible(true)
				break
			}
		}
	}

	/////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	itemBatch := pixel.NewBatch(&pixel.TrianglesData{}, animations.ItemImage)
	win.SetMatrix(pixel.IM.Scaled(pixel.V(0, 0), zoomFactor))
	win.Update()
	last := time.Now()
	dt := time.Since(last).Seconds()

	for !win.Closed() && !win.Pressed(pixelgl.KeyEscape) {
		keypressed := false
		dt = time.Since(last).Seconds()
		last = time.Now()
		if win.Pressed(pixelgl.KeyLeft) {
			moveCharacter(whiteBomberman, dt, Left)
			keypressed = true
		}
		if win.Pressed(pixelgl.KeyRight) {
			moveCharacter(whiteBomberman, dt, Right)
			keypressed = true
		}
		if win.Pressed(pixelgl.KeyUp) {
			moveCharacter(whiteBomberman, dt, Up)
			keypressed = true
		}
		if win.Pressed(pixelgl.KeyDown) {
			moveCharacter(whiteBomberman, dt, Down)
			keypressed = true
		}
		if !keypressed {
			whiteBomberman.Ani().SetView(Stay)
		}
		if win.JustPressed(pixelgl.KeyB) {
			pb := whiteBomberman.GetPosBox()
			loleft := turfNtreesArena.GetLowerLeft()
			b, _ := isThereABomb(pixel.Vec{math.Round(pb.Center().X/TileSize) * TileSize, math.Round(pb.Center().Y/TileSize) * TileSize})
			c := lev1.IsTile(int((pb.Min.X-loleft.X)/TileSize), int((pb.Min.Y-loleft.Y)/TileSize))
			if !b && !c {
				bombs = append(bombs, tiles.NewBomb(whiteBomberman))
			}
		}

		/////////////////////////////////////ToDO Moving Enemys ///////////////////////////////////////////////////////////

		/*
			for _,m := range(monster) {
				xx,yy := turfNtreesArena.GetFieldCoord(m.GetPos())
				x,y := turfNtreesArena.GetFieldCoord(whiteBomberman.GetPos())
				if x == xx && y == yy {
					whiteBomberman.DecLife()
				}
				if !m.IsFollowing() {
					dir := rand.Intn(4)
					switch dir {
						case 0:									// l
							if !turfNtreesArena.IsTile(xx-1,yy) {
								m.Move(pixel.V(-stepSize,0))
								m.Ani().SetView(Left)
							}
						case 1:									// r
							if !turfNtreesArena.IsTile(xx+1,yy) {
								m.Move(pixel.V(stepSize,0))
								m.Ani().SetView(Right)
							}
						case 2:									// up
							if !turfNtreesArena.IsTile(xx,yy+1) {
								m.Move(pixel.V(0,stepSize))
								m.Ani().SetView(Up)
							}
						case 3:
							if	!turfNtreesArena.IsTile(xx,yy-1) {
								m.Move(pixel.V(0,-stepSize))
								m.Ani().SetView(Down)
							}
					}
				}
			}
		*/

		/////////////////////////////////////////////////////////////////////////////////////////////////////////////7

		turfNtreesArena.GetCanvas().Draw(win, *(turfNtreesArena.GetMatrix()))

		checkForExplosions()
		bombs = removeExplodedBombs(bombs)

		itemBatch.Clear()

		//lev.DrawItems(win)
		//lev.DrawTiles(win)
		for i := 0; i < pitchHeight; i++ {
			lev1.DrawColumn(i, itemBatch)
		}

		for _, item := range bombs {
			item.Draw(itemBatch)
		}

		itemBatch.Draw(win)

		showExplosions(win)
		tempAniSlice = clearExplosions(tempAniSlice)

		whiteBomberman.Draw(win)

		for _, m := range monster {
			m.Draw(win)
		}

		tb.Draw(win)

		win.Update()
	}
}

func main() {
	pixelgl.Run(sun)
}
