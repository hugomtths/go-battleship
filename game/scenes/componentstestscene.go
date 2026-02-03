package scenes

import (
	"image/color"

	"github.com/allanjose001/go-battleship/game/components"
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	"github.com/hajimehoshi/ebiten/v2"
)

type ComponentsTestScene struct {
	containerTest components.Widget
}

func (c ComponentsTestScene) OnEnter(prev Scene) {
	c.init1() //ou init dois, só escolher o teste
}

func (c ComponentsTestScene) OnExit(next Scene) {
}

func (c ComponentsTestScene) Update() error {
	return nil
}

func (c ComponentsTestScene) Draw(screen *ebiten.Image) {
	c.containerTest.Draw(screen)
}

// teste com dois containers de exemplo
func (c ComponentsTestScene) init1() {
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
	c.containerTest = cont2

}

// text com row e container de exemplo
// (text não fica perfeito em crossalign (eixo y) start
// e end porque a renderização dá um retangulo um pouco maior ao text)
// posso tentar resolver no futuro conforme necessário
func (c ComponentsTestScene) init2() {
	row := components.NewRow(
		basic.Point{},
		20.0,
		basic.Size{},
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
	c.containerTest = row
}
