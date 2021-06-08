package levels

import (
	"fmt"
	"testing"
)

func TestMain(*testing.M) {
	var level1 Level
	level1 = NewLevel("./level1.txt")
	fmt.Println(level1.GetTilePos())
	fmt.Println(len(level1.GetTilePos()))
	fmt.Println(level1.GetLevelItems())
	fmt.Println(level1.GetLevelEnemies())
	fmt.Println(level1.GetBounds())
	fmt.Println(level1.GetArenaType())
}
