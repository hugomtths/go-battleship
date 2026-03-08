package scenes

import (
	"fmt"
	"log"

	"github.com/allanjose001/go-battleship/game/components"
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	"github.com/hajimehoshi/ebiten/v2"
)

type HomeScreen struct {
	layout       components.LayoutWidget
	muteButton   *components.IconButton
	StackHandler //faz composição (recebe os fields e metodos
}

func (m *HomeScreen) GetMusic() string {
	return "menus"
}

func (m *HomeScreen) OnExit(_ Scene) {
	m.stack.ctx.CanPopOrPush = false
}

func (m *HomeScreen) OnEnter(_ Scene, screenSize basic.Size) {

	err := m.init(screenSize)

	if err != nil {
		log.Fatal("Erro ao carregar componentes na tela inicial: ", err)
	}

	m.stack.ctx.CanPopOrPush = true

}

func (m *HomeScreen) Update() error {
	if m.layout != nil {
		m.layout.Update(basic.Point{X: 0, Y: 0})
	}

	if m.muteButton != nil {
		// offset 0,0 porque a posição do botão já é absoluta na tela
		m.muteButton.Update(basic.Point{X: 0, Y: 0})
	}

	return nil
}

func (m *HomeScreen) Draw(screen *ebiten.Image) {
	if m.layout != nil {
		m.layout.Draw(screen)
	}
	if m.muteButton != nil {
		m.muteButton.Draw(screen) // desenha o botão fixo no canto inferior esquerdo
	}
}

func (m *HomeScreen) toggleMute() {
	fmt.Println("toggle mute")

	m.stack.ctx.SoundService.ToggleMute()

	if m.stack.ctx.SoundService.IsMuted() {
		m.muteButton.SetIcon("assets/images/mute.png")
	} else {
		m.muteButton.SetIcon("assets/images/unmute.png")
	}
}

func (m *HomeScreen) init(screenSize basic.Size) error {
	var err error
	homeImage, err := components.NewImage(
		"assets/images/home-screen.png",
		basic.Point{},
		basic.Size{W: 500, H: 500})

	if err != nil {
		return err
	}

	m.muteButton = components.NewIconButton(
		"assets/images/unmute.png",
		basic.Point{},
		basic.Size{W: 40, H: 40},
		m.toggleMute,
	)

	if m.stack.ctx.SoundService.IsMuted() {
		m.muteButton.SetIcon("assets/images/mute.png")
	}

	// ➤ Defina a posição manualmente para o canto inferior esquerdo
	m.muteButton.SetPos(basic.Point{
		X: 10,                // 10px da esquerda
		Y: screenSize.H - 50, // 10px acima da base (considerando altura 40px do botão)
	})

	m.layout = components.NewColumn(
		basic.Point{},
		20,
		screenSize,
		basic.Start,
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
					m.stack.Push(&SelectProfileScene{})

				},
			),

			components.NewButton(
				basic.Point{},
				basic.Size{W: 400, H: 50},
				"Ranking",
				colors.Dark,
				nil,
				func(bt *components.Button) {
					m.stack.Push(&RankingScene{})
				},
			),

			components.NewButton(
				basic.Point{},
				basic.Size{W: 400, H: 50},
				"Sair",
				colors.Dark,
				nil,
				func(bt *components.Button) {
					fmt.Println("sair")
					m.stack.Pop() //faz terminator em game
				},
			),
		},
	)

	// ➤ Atualize também o Update para o botão fixo
	if m.layout != nil {
		m.layout.Update(basic.Point{X: 0, Y: 0})
	}
	m.muteButton.Update(basic.Point{X: 0, Y: 0}) // offset 0,0 porque posição já está definida

	return nil
}
