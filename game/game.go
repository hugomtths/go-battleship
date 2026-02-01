package game

import (
	"image/color"

	"github.com/allanjose001/go-battleship/game/components"
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	gameW = 1280
	gameH = 720
)

type Game struct {
	doubleContainerTest components.Widget
}

// teste com dois containers de exemplo
func (g *Game) init1() {
	contd := components.NewContainer(
		basic.Point{50.0, 50.0},
		basic.Size{W: 400, H: 600},
		12.0,
		colors.Blue,
		basic.Center, //main -> x
		basic.Center, //cross -> y
		components.NewText(basic.Point{}, "BATTLESHIP", color.White, 48),
		func(c *components.Container) {},
	)

	cont2 := components.NewContainer(
		basic.Point{200, 0},
		basic.Size{W: 1000, H: 700},
		0.0,
		colors.White,
		basic.Start,  //main -> x
		basic.Center, //cross -> y
		contd,
		func(c *components.Container) {},
	)
	g.doubleContainerTest = cont2

}

// text com row e container de exemplo
// (text não fica perfeito em crossalign (eixo y) start
// e end porque a renderização dá um retangulo um pouco maior ao text)
// posso tentar resolver no futuro conforme necessário
func (g *Game) init2() {
	col := components.NewRow(
		basic.Point{},
		20.0,
		basic.Size{W: gameW, H: gameH},
		basic.Center,
		basic.Center,
		[]components.Widget{
			components.NewContainer(
				basic.Point{},
				basic.Size{W: 400, H: 600},
				0.0,
				colors.Blue,
				basic.Center, //main -> x
				basic.Center, //cross -> y
				components.NewText(basic.Point{}, "BATTLESHIP", color.White, 48),
				func(c *components.Container) {},
			),
		},
	)
	g.doubleContainerTest = col
}

func NewGame() *Game {
	g := &Game{}
	g.init1() //escolher o teste aqui
	return g
}
func (g *Game) Update() error {
	g.doubleContainerTest.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 30, G: 30, B: 30, A: 255})
	g.doubleContainerTest.Draw(screen)
}

func (g *Game) Layout(_, _ int) (int, int) {
	return gameW, gameH
}
