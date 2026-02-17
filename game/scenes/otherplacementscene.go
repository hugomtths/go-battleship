package scenes

import (
	"github.com/allanjose001/go-battleship/game/components"
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	"github.com/hajimehoshi/ebiten/v2"
)

type OtherPlacementScene struct {
	root components.Widget
	StackHandler
	playButton *components.Widget // necessária essa ref separada para fazer disable

	//playButton *components.Button // para poder desabilitar, é preciso essa ref
}

func (o *OtherPlacementScene) OnEnter(_ Scene, size basic.Size) {

	majorDiv := basic.Size{W: size.W, H: size.H * 0.70}
	minorDiv := basic.Size{W: size.W, H: size.H * 0.30}

	boardAndFleetSize := basic.Size{
		W: 360,
		H: 360,
	}

	o.root = components.NewColumn(
		basic.Point{0, 15}, //posiciona column um pouco mais a baixo na tela
		10, size,
		basic.Center,
		basic.Center,
		[]components.Widget{
			//container "div" invisivel com uma row (largura da tela e altura arbitraria)
			components.NewContainer(
				basic.Point{}, majorDiv,
				0, colors.Transparent,
				basic.Center,
				basic.Center,

				//row que sera trabalhada
				components.NewRow(
					basic.Point{}, 200, majorDiv,
					basic.Center, basic.Center,
					[]components.Widget{
						//TODO: Board fica aqui -> cria ele antes, mas adiciona ele aqui para ficar alinhado
						components.NewContainer(
							basic.Point{}, boardAndFleetSize,
							0, colors.Blue,
							basic.Center,
							basic.Center,
							nil,
						),

						// lista de embarcações
						components.NewContainer(
							basic.Point{}, boardAndFleetSize,
							0, colors.Transparent,
							basic.Center,
							basic.Center,
							components.NewColumn(
								basic.Point{}, 5, boardAndFleetSize,
								basic.Center, basic.Center,
								[]components.Widget{
									components.NewText(basic.Point{}, "barcos aqui", colors.White, 20),
								},
							),
						),
					},
				),
			),
			//container "div" com outra row para os botões
			components.NewContainer(
				basic.Point{}, minorDiv,
				0, colors.Transparent,
				basic.Center,
				basic.Center,
				components.NewRow(
					basic.Point{}, 50, minorDiv,
					basic.Center, basic.Start,
					[]components.Widget{
						components.NewButton(
							basic.Point{},
							basic.Size{250, 70},
							"Aleatório",
							colors.Dark,
							colors.White,
							func(b *components.Button) {},
						),

						components.NewButton(
							basic.Point{},
							basic.Size{250, 70},
							"Rotacionar",
							colors.Dark,
							colors.White,
							func(b *components.Button) {},
						),

						components.NewContainer( // container vazio apenas para preencher espaço
							basic.Point{}, basic.Size{120, 1},
							0, colors.Transparent,
							basic.Center,
							basic.Center,
							nil),

						*o.playButton,
					},
				),
			),
		},
	)
}

// creio que aqui deva passar adiante o estado dessa cena atual
// talvez deva ser criado um struct placementstate para facilitar e ser pego pela proxima scene
func (o *OtherPlacementScene) OnExit(_ Scene) {
}

func (o *OtherPlacementScene) Update() error {

	//logica para desabilitar/habilitar botão aqui com o.playButton.SetDisabled, mas e preciso adicionar o boolean no construtor de button...

	o.root.Update(basic.Point{})
	return nil
}

func (o *OtherPlacementScene) Draw(screen *ebiten.Image) {
	o.root.Draw(screen)
}
