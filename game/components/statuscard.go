package components

import (
	"fmt"
	"image/color"

	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	"github.com/allanjose001/go-battleship/internal/entity"
	"github.com/hajimehoshi/ebiten/v2"
)

// StatusCard representa um painel visual de estatísticas.
// Ele encapsula layout, tamanho e renderização.
// A regra de negócio vem exclusivamente de PlayerStats.
type StatusCard struct {
	pos        basic.Point
	currentPos basic.Point
	size       basic.Size
	body       StylableWidget // container raiz que desenha tudo
}

// SetPos define posição base lógica.
func (s *StatusCard) SetPos(point basic.Point) {
	s.pos = point
}

// GetPos retorna posição base.
func (s *StatusCard) GetPos() basic.Point {
	return s.pos
}

// GetSize retorna tamanho total do card.
func (s *StatusCard) GetSize() basic.Size {
	return s.size
}

// Update propaga offset para o corpo.
func (s *StatusCard) Update(offset basic.Point) {
	s.currentPos = s.pos.Add(offset)
	s.body.Update(s.currentPos)
}

// Draw delega renderização ao container interno.
func (s *StatusCard) Draw(screen *ebiten.Image) {
	s.body.Draw(screen)
}

// NewStatCard constrói o componente completo.
// stats contém toda regra de domínio.
// isRanking define layout compacto ou completo.
func NewStatCard(
	pos basic.Point,
	screenSize basic.Size,
	stats *entity.PlayerStats,
	isRanking bool,
	playerName string,
	rankingPosition int,
) *StatusCard {

	// Decide tamanhos estruturais do card
	contSize, cardSize := switchSizes(isRanking, screenSize)

	// Gera widgets internos a partir da entidade
	statsList := initWidgets(stats, isRanking, cardSize)

	// Widget de título (varia se for ranking)
	var titleWidget Widget

	if isRanking {

		var circleColor color.Color
		switch rankingPosition {
		case 1:
			circleColor = colors.GoldMedal
		case 2:
			circleColor = colors.SilverMedal
		case 3:
			circleColor = colors.BronzeMedal
		default:
			circleColor = colors.White
		}

		rankText := NewText(
			basic.Point{},
			fmt.Sprintf("%d°", rankingPosition),
			colors.Black,
			20,
		)

		rankCircle := NewContainer(
			basic.Point{},
			basic.Size{W: 40, H: 40},
			20,
			circleColor,
			basic.Center,
			basic.Center,
			rankText,
		)

		nameText := NewText(
			basic.Point{},
			playerName,
			colors.White,
			30,
		)

		// Row usa contSize como referência de alinhamento.
		// O size passado para Row representa o tamanho do pai imediato
		// que será usado para calcular Center / End.
		titleWidget = NewContainer( //container fantasma para manipular o alinhamento
			basic.Point{},
			basic.Size{W: contSize.W * 0.3, H: 40},
			0,
			colors.Transparent,
			basic.Center,
			basic.Center,
			NewRow(
				basic.Point{},
				15,
				basic.Size{W: contSize.W * 0.3, H: 40}, //coloca para se alinhar em 1/3 da largura
				basic.Start,
				basic.Center,
				[]Widget{rankCircle, nameText},
			),
		)

	} else {

		titleWidget = NewText(
			basic.Point{},
			playerName,
			colors.White,
			30,
		)
	}

	return &StatusCard{
		pos:        pos,
		currentPos: pos,
		size:       contSize,

		// Container raiz.
		// Ele recebe contSize e controla alinhamento global.
		body: NewContainer(
			basic.Point{},
			contSize,
			12,
			colors.Dark,
			basic.Start,
			basic.Start,

			// Column principal do card.
			// IMPORTANTE:
			// O size passado (contSize) é o tamanho do pai.
			// Column usa esse size apenas como referência
			// para calcular MainAlign e CrossAlign.
			NewColumn(
				basic.Point{},
				10,
				contSize, // size do pai para cálculo de alinhamento
				basic.Center,
				basic.Center,
				[]Widget{

					// Título
					titleWidget,

					// Container intermediário apenas para organizar a Row de stats.
					NewContainer(
						basic.Point{},
						basic.Size{
							W: contSize.W,
							H: cardSize.H,
						},
						0,
						nil,
						basic.Start,
						basic.Start,

						// Row das estatísticas.
						// O size passado aqui representa
						// o espaço disponível dentro desse container.
						NewRow(
							basic.Point{},
							20,
							basic.Size{
								W: contSize.W,
								H: cardSize.H,
							},
							basic.Center,
							basic.Center,
							statsList,
						),
					),
				},
			),
		),
	}
}

// switchSizes decide layout estrutural.
func switchSizes(isRanking bool, screenSize basic.Size) (basic.Size, basic.Size) {

	if isRanking {
		return basic.Size{
				W: 0.9 * screenSize.W,
				H: 120,
			},
			basic.Size{W: 300, H: 60}
	}

	return basic.Size{
			W: 0.9 * screenSize.W,
			H: 220,
		},
		basic.Size{W: 210, H: 100}
}

// initWidgets converte PlayerStats em widgets visuais.
// Nenhuma regra é calculada aqui — apenas formatação.
func initWidgets(
	stats *entity.PlayerStats,
	ranking bool,
	size basic.Size,
) []Widget {

	if ranking {
		return []Widget{
			createStatCard("Pontuação Total",
				fmt.Sprintf("%d", stats.TotalScore), size),
			createStatCard("Vitórias",
				fmt.Sprintf("%d", stats.Wins), size),
			createStatCard("Maior Sequência de Hits",
				fmt.Sprintf("%d", stats.HigherHitSequence), size),
		}
	}

	return []Widget{
		createStatCard("Partidas",
			fmt.Sprintf("%d", stats.Matches), size),
		createStatCard("Vitórias",
			fmt.Sprintf("%d", stats.Wins), size),
		createStatCard("Winrate",
			fmt.Sprintf("%.2f %%", stats.WinRate()), size),
		createStatCard("Hitrate",
			fmt.Sprintf("%.2f %%", stats.Accuracy()), size),
		createStatCard("Maior Score",
			fmt.Sprintf("%d", stats.HighScore), size),
	}
}

// createStatCard encapsula label + valor.
// Column interna usa o size do container pai
// para calcular alinhamento centralizado.
func createStatCard(label, value string, size basic.Size) *Container {

	labelTxt := NewText(basic.Point{}, label, colors.Black, 20)
	valueTxt := NewText(basic.Point{}, value, colors.Black, 25)

	content := NewColumn(
		basic.Point{},
		5,
		size, // size do pai usado como referência de alinhamento
		basic.Center,
		basic.Center,
		[]Widget{labelTxt, valueTxt},
	)

	return NewContainer(
		basic.Point{},
		size,
		12,
		colors.White,
		basic.Center,
		basic.Center,
		content,
	)
}
