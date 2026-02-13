package game

import (
	"image/color"

	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/scenes"
	"github.com/hajimehoshi/ebiten/v2"
)

var windowSize = basic.Size{W: 1280, H: 800}

var currentGame *Game

type Game struct {
	scene scenes.Scene
}

func ChangeScene(s scenes.Scene) {
	if currentGame != nil {
		if currentGame.scene != nil {
			currentGame.scene.OnExit(s)
		}
		prev := currentGame.scene
		currentGame.scene = s
		currentGame.scene.OnEnter(prev, windowSize)
	}
}

func NewGame() *Game {
	scenes.SwitchTo = ChangeScene
	g := &Game{
	}
	currentGame = g
	g.scene.OnEnter(nil, windowSize)
	return g
}
func (g *Game) Update() error {
	err := g.scene.Update()
	return err
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 13, G: 27, B: 42, A: 255})
	g.scene.Draw(screen)
}

func (g *Game) Layout(_, _ int) (int, int) {
	return int(windowSize.W), int(windowSize.H)
}
