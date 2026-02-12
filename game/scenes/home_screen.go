package scenes

import (
	"image/color"
	"log"

	"github.com/allanjose001/go-battleship/game/components"
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type HomeScreen struct {
	layout       components.Widget
	homeImage    *ebiten.Image
	alturaDoTopo float32
	sairDoJogo   bool
}

func (m *HomeScreen) OnEnter(_ Scene, screenSize basic.Size) {
	var err error
	m.homeImage, _, err = ebitenutil.NewImageFromFile("assets/images/home-screen.png")

	if err != nil {
		log.Fatal("Erro ao carregar a imagem:", err)
	}

	m.alturaDoTopo = 450.0

	m.layout = components.NewColumn(
		basic.Point{X: 0, Y: m.alturaDoTopo},
		20,
		basic.Size{W: screenSize.W, H: screenSize.H - m.alturaDoTopo},
		basic.Start,
		basic.Center,
		[]components.Widget{
			components.NewButton(
				basic.Point{},
				basic.Size{W: 300, H: 50},
				"Jogar",
				color.RGBA{R: 48, G: 67, B: 103, A: 255},
				nil,
				func(bt *components.Button) {
					log.Println("Botão clicado!") // Aqui ficará a função que inicia o jogo (mudar para a tela de jogo)
				},
			),

			components.NewButton(
				basic.Point{},
				basic.Size{W: 300, H: 50},
				"Ranking",
				color.RGBA{R: 48, G: 67, B: 103, A: 255},
				nil,
				func(bt *components.Button) {
					log.Println("Botão clicado!") // Aqui ficará a função que mostra o ranking (mudar para a tela de ranking)
				},
			),

			components.NewButton(
				basic.Point{},
				basic.Size{W: 300, H: 50},
				"Sair",
				color.RGBA{R: 48, G: 67, B: 103, A: 255},
				nil,
				func(bt *components.Button) {
					m.sairDoJogo = true
					log.Println("Saindo do jogo...")
				},
			),
		},
	)
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
	screen.Fill(color.RGBA{R: 13, G: 27, B: 42, A: 255})

	if m.homeImage != nil {
		escalaImagem := 0.8

		larguraTela := screen.Bounds().Dx()
		larguraImagem := m.homeImage.Bounds().Dx()
		larguraImagemEscalada := float64(larguraImagem) * escalaImagem

		posicaoX := float64((larguraTela - int(larguraImagemEscalada)) / 2)

		posicaoY := 60.0

		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Scale(escalaImagem, escalaImagem)
		opts.GeoM.Translate(posicaoX, posicaoY)
		screen.DrawImage(m.homeImage, opts)
	}

	if m.layout != nil {
		m.layout.Draw(screen)
	}
}

func (m *HomeScreen) OnExit(_ Scene) {}
