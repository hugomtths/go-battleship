package main

import (
	"github.com/allanjose001/go-battleship/game"
	"github.com/allanjose001/go-battleship/game/components"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	components.InitFonts() //carrega a fonte apenas uma vez
	g := game.NewGame()
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}

}
