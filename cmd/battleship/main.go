package main

import (
	"github.com/allanjose001/go-battleship/game"
	"github.com/allanjose001/go-battleship/internal/bootstrap"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	bootstrap.InitRandom()

	g := game.NewGame()
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := game.SetGameWindowIcon("assets/icons/skull.png"); err != nil {
		panic(err)
	}

	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
