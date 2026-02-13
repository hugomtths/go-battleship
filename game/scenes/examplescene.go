package scenes

import (
	"github.com/allanjose001/go-battleship/game/components"
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	"github.com/hajimehoshi/ebiten/v2"
)

type ExampleScene struct {
	container components.Widget
	col       components.Column
	row       components.Row
}

func (e *ExampleScene) OnEnter(_ Scene, size basic.Size) {

	//pode-se criar tudo um dentro do outro (para quem entende a "arvore de widgets")
	/*e.row = *components.NewRow(
		basic.Point{},
		10.0,
		size,
		basic.Center,
		basic.Center,
		[]components.Widget{
			components.NewContainer(
				basic.Point{50, 50},
				basic.Size{W: 0.8 * size.W, H: 0.8 * size.H},
				10.0,
				colors.Dark,
				basic.Center, //funciona com filho "unico"
				basic.Center,
				components.NewColumn(
					basic.Point{}, //POSIÇÃO EM RELAÇÃO AO PAI !
					10.0,
					basic.Size{W: 0.8 * size.W, H: 0.8 * size.H},
					basic.Center,
					basic.Center,
					[]components.Widget{
						components.NewText(
							basic.Point{},
							"Teste dos Buttons",
							colors.White,
							40,
						),
						components.NewContainer( //aqui eu uso container apenas para ocupar espaço
							basic.Point{},
							basic.Size{H: 70, W: 1},
							0.0,
							colors.Transparent,
							basic.Center,
							basic.Center,
							nil,
						),

						components.NewButton(
							basic.Point{},
							basic.Size{W: 400, H: 70},
							"botão x",
							colors.Green,
							nil,
							func(b *components.Button) {},
						),
						components.NewButton(
							basic.Point{},
							basic.Size{W: 400, H: 70},
							"botão x",
							colors.Green,
							nil,
							func(b *components.Button) {},
						),
						components.NewButton(
							basic.Point{},
							basic.Size{W: 400, H: 70},
							"botão x",
							colors.Green,
							nil,
							func(b *components.Button) {},
						),
					},
				),
			),
		},
	)*/

	//ou pode-se criar tudo separadamente assim:

	//crio botões em uma lista
	buttons := []components.Widget{
		components.NewText(
			basic.Point{},
			"Teste dos Buttons",
			colors.White,
			40,
		),
		components.NewContainer( //aqui eu uso container apenas para ocupar espaço
			basic.Point{},
			basic.Size{H: 70, W: 1},
			0.0,
			colors.Transparent,
			basic.Center,
			basic.Center,
			nil,
		),

		components.NewButton(
			basic.Point{},
			basic.Size{W: 400, H: 70},
			"botão x",
			colors.Green,
			nil,
			func(b *components.Button) {},
		),
		components.NewButton(
			basic.Point{},
			basic.Size{W: 400, H: 70},
			"botão x",
			colors.Green,
			nil,
			func(b *components.Button) {},
		),
		components.NewButton(
			basic.Point{},
			basic.Size{W: 400, H: 70},
			"botão x",
			colors.Green,
			nil,
			func(b *components.Button) {},
		),
	}

	containerSize := basic.Size{W: 0.8 * size.W, H: 0.8 * size.H} //dado que col deve ter tamanho do pai

	e.col = *components.NewColumn(
		basic.Point{}, //POSIÇÃO EM RELAÇÃO AO PAI !
		10.0,
		containerSize, //tamanho do pai
		basic.Center,
		basic.Center,
		buttons,
	)
	e.container = components.NewContainer(
		basic.Point{50, 50},
		containerSize, //pai com seu tamanho
		10.0,
		colors.Dark,
		basic.End, //funciona quando filho não é layout (col  e row)
		basic.Start,
		&e.col,
	)

	//usado para alinhar o container na tela
	e.row = *components.NewRow(
		basic.Point{},
		10.0,
		size, //tamanho total da tela
		basic.Center,
		basic.Center,
		[]components.Widget{
			e.container,
		},
	)

}

func (e *ExampleScene) OnExit(_ Scene) {
}

func (e *ExampleScene) Update() error {
	e.row.Update(basic.Point{})
	return nil
}

func (e *ExampleScene) Draw(screen *ebiten.Image) {
	e.row.Draw(screen)
}
