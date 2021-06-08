package constants

import "testing"

func TestHello(t *testing.T) {
	if MaxWinSizeX < 320 {
		t.Errorf("MaxWinSizeX (%q) ist zu klein. Bitte größer/gleich als 320 wählen.", MaxWinSizeX)
	}
	if MaxWinSizeY < 200 {
		t.Errorf("MaxWinSizeY (%q) ist zu klein. Bitte größer/gleich als 200 wählen.", MaxWinSizeY)
	}
}
