package scenes

import (
	"log"

	"github.com/allanjose001/go-battleship/game/components"
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	"github.com/hajimehoshi/ebiten/v2"
)

type HomeScreen struct {
	layout     components.Widget //deixei tudo no layout
	sairDoJogo bool
}

func (m *HomeScreen) OnExit(_ Scene) {}

func (m *HomeScreen) OnEnter(_ Scene, screenSize basic.Size) {
	err := m.init(screenSize)

	if err != nil {
		log.Fatal("Erro ao carregar componentes na tela inicial: ", err)
	}

}

func (m *HomeScreen) Update() error {
	if m.layout != nil {
		m.layout.Update(basic.Point{X: 0, Y: 0})
	}
	if m.sairDoJogo {
		return ebiten.Termination
	}

	return nil
}

func (m *HomeScreen) Draw(screen *ebiten.Image) {
	if m.layout != nil {
		m.layout.Draw(screen)
	}
}

// init Inicializa componentes
func (m *HomeScreen) init(screenSize basic.Size) error {
	var err error
	homeImage, err := components.NewImage(
		"assets/images/home-screen.png",
		basic.Point{},
		basic.Size{W: 580, H: 400})

	if err != nil {
		return err
	}
	m.layout = components.NewColumn(
		basic.Point{},
		20,
		basic.Size{W: screenSize.W, H: screenSize.H},
		basic.Center,
		basic.Center,
		[]components.Widget{
			homeImage,
			components.NewButton(
				basic.Point{},
				basic.Size{W: 400, H: 50},
				"Jogar",
				colors.Dark,
				nil,
				func(bt *components.Button) {
					log.Println("Botão clicado!") // Aqui ficará a função que inicia o jogo (mudar para a tela de jogo)
				},
			),

			components.NewButton(
				basic.Point{},
				basic.Size{W: 400, H: 50},
				"Ranking",
				colors.Dark,
				nil,
				func(bt *components.Button) {
					log.Println("Botão clicado!") // Aqui ficará a função que mostra o ranking (mudar para a tela de ranking)
				},
			),

			components.NewButton(
				basic.Point{},
				basic.Size{W: 400, H: 50},
				"Sair",
				colors.Dark,
				nil,
				func(bt *components.Button) {
					m.sairDoJogo = true
					log.Println("Saindo do jogo...")
				},
			),
		},
	)
	return nil
}
