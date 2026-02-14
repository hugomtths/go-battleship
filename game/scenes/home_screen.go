package scenes

import (
	"log"

	"github.com/allanjose001/go-battleship/game/components"
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	"github.com/hajimehoshi/ebiten/v2"
)

type HomeScreen struct {
	layout       components.Widget
	StackHandler //faz composição (recebe os fields e metodos)
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
					//m.stack.Push() //TODO: colocar tela de ESCOLHA de perfis aqui
				},
			),

			components.NewButton(
				basic.Point{},
				basic.Size{W: 400, H: 50},
				"Ranking",
				colors.Dark,
				nil,
				func(bt *components.Button) {
					//m.stack.Push() //TODO: ir para tela de Ranking
					log.Println("Botão clicado!")
				},
			),

			components.NewButton(
				basic.Point{},
				basic.Size{W: 400, H: 50},
				"Sair",
				colors.Dark,
				nil,
				func(bt *components.Button) {
					m.stack.Pop() //faz terminator em game
				},
			),
		},
	)
	return nil
}
