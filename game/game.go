package game

import (
	"image/color"

	"github.com/allanjose001/go-battleship/game/scenes"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	gameW = 1280
	gameH = 720
)

type Game struct {
	scene scenes.Scene
}

func NewGame() *Game {
	g := &Game{
		scene: &scenes.ButtonScene{},
	}
	g.scene.OnEnter(nil) //escolher o teste no OnEnter dessa struct
	return g
}
func (g *Game) Update() error {
	err := g.scene.Update()
	return err
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 30, G: 30, B: 30, A: 255})
	g.scene.Draw(screen)
}

func (g *Game) Layout(_, _ int) (int, int) {
	return gameW, gameH
}
