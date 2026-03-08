package scenes

import (
	"github.com/allanjose001/go-battleship/game/components"
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	"github.com/allanjose001/go-battleship/game/state"
	"github.com/allanjose001/go-battleship/internal/medal"
	"github.com/hajimehoshi/ebiten/v2"
)

// ProfileScene representa a tela de perfil do jogador.
type ProfileScene struct {
	state *state.GameState
	root  *components.Column // O container pai que envolve toda a cena.
	StackHandler
}

func (p *ProfileScene) GetMusic() string {
	return "menus"
}

// Implementações do contrato Scene
func (p *ProfileScene) OnEnter(prev Scene, size basic.Size) {
	p.init(size)
	p.stack.ctx.CanPopOrPush = true
}

func (p *ProfileScene) OnExit(next Scene) {
	p.stack.ctx.CanPopOrPush = false
}

func (p *ProfileScene) Update() error {
	p.root.Update(basic.Point{X: 0, Y: 0})
	return nil
}

func (p *ProfileScene) Draw(screen *ebiten.Image) {
	p.root.Draw(screen)
}

// init Função que inicializa componentes
func (p *ProfileScene) init(size basic.Size) {
	playerName := p.stack.ctx.Profile.Username

	// Chamamos o method agora vinculado à struct
	medals := p.loadMedals()

	// Coluna principal que centraliza verticalmente
	p.root = components.NewColumn(
		basic.Point{},
		40,
		size,
		basic.Start,
		basic.Center,
		[]components.Widget{
			// Title
			components.NewText(basic.Point{},
				"PERFIL DE JOGADOR",
				colors.White,
				42),

			// Container com Row para estatisticas
			components.NewStatCard(
				basic.Point{},
				size,
				&p.stack.ctx.Profile.Stats,
				false,
				playerName,
				0,
			),

			// Título da seção de medalhas
			components.NewText(basic.Point{}, "MURAL DE MEDALHAS", colors.White, 28),

			// Container com Row para medalhas reais
			components.NewContainer(
				basic.Point{},
				basic.Size{W: 750, H: 100},
				0, nil,
				basic.Center, basic.Center,
				components.NewRow(
					basic.Point{},
					40,
					basic.Size{W: 750, H: 100},
					basic.Center, basic.Center,
					*medals,
				),
			),

			// Botão para acessar o histórico de partidas
			components.NewButton(
				basic.Point{},
				basic.Size{W: 300, H: 55},
				"Histórico de Partidas",
				colors.Dark,
				colors.White,
				func(b *components.Button) {
					p.stack.Push(&MatchsHistory{})
				},
			),

			// Botão Voltar
			components.NewButton(
				basic.Point{},
				basic.Size{W: 220, H: 55},
				"Voltar",
				colors.Dark,
				colors.White,
				func(b *components.Button) {
					p.stack.Pop()
				},
			),
		},
	)
	_ = p.Update()
}

// loadMedals agora é um método de ProfileScene para acessar p.stack.ctx.Profile.Stats
func (p *ProfileScene) loadMedals() *[]components.Widget {
	var widgets = []components.Widget{}

	playerMedalNames := p.stack.ctx.Profile.MedalsNames
	for i, m := range medal.GetMedals(playerMedalNames) { //isso retorna o array com posicoes preservadas
		displayIcon := medal.MedalsList[i].GrayIconPath
		displayTitle := "BLOQUEADA"
		displayDesc := "???"

		if m != nil { //posicao vazia = nao teve a medal
			displayIcon = m.IconPath
			displayTitle = m.Name
			displayDesc = m.Description
		}

		medalW := components.NewMedal(displayIcon, displayTitle, displayDesc, basic.Size{W: 230, H: 90})
		widgets = append(widgets, medalW)
	}

	return &widgets
}
