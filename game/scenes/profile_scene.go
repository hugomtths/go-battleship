package scenes

import (
	"github.com/allanjose001/go-battleship/game/components"
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	"github.com/allanjose001/go-battleship/game/state"
	"github.com/allanjose001/go-battleship/internal/service"
	"github.com/hajimehoshi/ebiten/v2"
)

// ProfileScene representa a tela de perfil do jogador.
type ProfileScene struct {
	state   *state.GameState
	profile *service.Profile
	root    *components.Column // O container pai que envolve toda a cena.
}

// init Funcão que inicializa componentes
func (p *ProfileScene) init(size basic.Size) {
	// Recupera os dados do jogador do serviço
	//TODO: Isso é um state passado da cena de seleção de perfis, carregado lá
	//profile, _ := service.FindProfile("malub")

	medals := loadMedals()

	///TODO: Criar componente medal (ver se precisa back e front)

	// Coluna principal que centraliza verticalmente
	p.root = components.NewColumn(
		basic.Point{},
		40,
		size,
		basic.Start,
		basic.Center,
		[]components.Widget{

			//Title
			components.NewText(basic.Point{},
				"PERFIL DE JOGADOR",
				colors.White,
				42),

			//Container com Row para estatisticas
			components.NewStatCard(
				basic.Point{},
				//TODO: Criar o tipo datastats para facilitar isso, e facilitar carregar/salvar no json em profile
				size, //usa tamanho da tela para caso mude a resolução
				2999, 200, 90000, 62, 80,
				false,                 //para reutilizar em ranking
				"Nome do player aqui", //mock, precisa melhorar profile pra ter tudo
			),
			//medalhas
			components.NewText(basic.Point{}, "MURAL DE MEDALHAS", colors.White, 28),
			//Container com Col para medals
			components.NewContainer(
				basic.Point{},
				basic.Size{W: 750, H: 100},
				0, nil,
				basic.Center, basic.Center, //alinhamento não importa quando filho é layout
				components.NewRow(
					basic.Point{},
					40,
					basic.Size{W: 750, H: 100},
					basic.Center, basic.Center,
					*medals,
				),
			),

			//voltar

			components.NewButton(
				basic.Point{},
				basic.Size{220, 55},
				"Retornar",
				colors.Dark,
				colors.White,
				func(b *components.Button) {},
			),
		},
	)
}

// TODO: Criar isso no ProfileService para carregar de um arquivo contendo as medals -> aqui transforma em widget
func loadMedals() *[]components.Widget {
	medalData := []struct { // isso aqui pode ser as medals carregadas do json
		Icon, Title, Desc string
	}{
		{"X", "VETERANO", "10+ Partidas"},
		{"W", "SNIPER", "90% Precisão"},
		{"Q", "VELOZ", "Vitória em <5min"},
		{"S", "IMPENETRÁVEL", "0 acertos sofridos"},
	}
	var medals = []components.Widget{}
	for _, data := range medalData {
		medals = append(medals, components.NewMedal(
			data.Icon, data.Title, data.Desc, basic.Size{W: 230, H: 90}),
		)
	}
	return &medals
}

// Implementações do contrato Scene
func (p *ProfileScene) OnEnter(prev Scene, size basic.Size) {
	// Atualiza os dados do perfil ao entrar na cena
	profile, _ := service.FindProfile("malub")
	p.profile = profile
	p.init(size)

}

func (p *ProfileScene) OnExit(next Scene) {
	//aqui creio que vá passar o profile para a tela de jogo caso a proxima seja a tela de jogo
}

// Update propaga a atualização de baixo para cima na árvore de componentes.
func (p *ProfileScene) Update() error {
	p.root.Update(basic.Point{X: 0, Y: 0})
	return nil
}

// Draw renderiza recursivamente toda a cena.
func (p *ProfileScene) Draw(screen *ebiten.Image) {
	p.root.Draw(screen)
}
