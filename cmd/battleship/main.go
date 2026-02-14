package main

import (
	"github.com/allanjose001/go-battleship/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	g := game.NewGame()
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}

}
