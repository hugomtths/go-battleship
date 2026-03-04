package scenes

import (
	"github.com/allanjose001/go-battleship/game/components"
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	"github.com/allanjose001/go-battleship/internal/service"
	"github.com/hajimehoshi/ebiten/v2"
)

type RankingScene struct {
	layout      components.Widget
	currentPage int
	StackHandler
}

func (m *RankingScene) OnExit(_ Scene) {}

func (m *RankingScene) OnEnter(_ Scene, screenSize basic.Size) {
	m.init(screenSize)
}

func (m *RankingScene) Update() error {
	if m.layout != nil {
		m.layout.Update(basic.Point{X: 0, Y: 0})
	}
	return nil
}

func (m *RankingScene) Draw(screen *ebiten.Image) {
	if m.layout != nil {
		m.layout.Draw(screen)
	}
}

func calculateRankingHeight(screenSize basic.Size) float32 {
	return screenSize.H * float32(0.5)
}

func (m *RankingScene) init(screenSize basic.Size) {

	itemsPerPage := 3
	start := m.currentPage * itemsPerPage
	end := start + itemsPerPage

	allPlayers := service.GetTopScores(9)

	if start > len(allPlayers) {
		start = len(allPlayers)
	}
	if end > len(allPlayers) {
		end = len(allPlayers)
	}

	pagePlayers := allPlayers[start:end]

	hasPrevious := m.currentPage > 0
	hasNext := end < len(allPlayers)

	var previousHandler func(*components.Button)
	var nextHandler func(*components.Button)

	prevColor := colors.Dark
	nextColor := colors.Dark

	if hasPrevious {
		previousHandler = func(bt *components.Button) {
			m.currentPage--
			m.init(screenSize)
		}
	} else {
		previousHandler = nil
		prevColor = colors.NightBlue
	}

	if hasNext {
		nextHandler = func(bt *components.Button) {
			m.currentPage++
			m.init(screenSize)
		}
	} else {
		nextHandler = nil
		nextColor = colors.NightBlue
	}

	topSpacer := components.NewContainer(
		basic.Point{},
		basic.Size{W: screenSize.W, H: screenSize.H * 0.1},
		0,
		nil,
		basic.Center,
		basic.Center,
		nil,
	)

	title := components.NewContainer(
		basic.Point{},
		basic.Size{W: screenSize.W, H: 60},
		0,
		nil,
		basic.Center,
		basic.Center,
		components.NewText(
			basic.Point{},
			"Ranking",
			colors.White,
			42,
		),
	)

	rankingHeight := calculateRankingHeight(screenSize)

	var cards []components.Widget
	for i, player := range pagePlayers {
		card := components.NewStatCard(
			basic.Point{},
			screenSize,
			&player.Stats,
			true,
			player.Username,
			start+i+1,
		)
		cards = append(cards, card)
	}

	cardsColumn := components.NewColumn(
		basic.Point{},
		10,
		basic.Size{W: screenSize.W, H: rankingHeight},
		basic.Start, 
		basic.Center,
		cards,
	)

	rankingContainer := components.NewContainer(
		basic.Point{},
		basic.Size{W: screenSize.W, H: rankingHeight},
		0,
		nil,
		basic.Start, 
		basic.Center,
		cardsColumn,
	)

	previousButton := components.NewButton(
		basic.Point{},
		basic.Size{W: 150, H: 40},
		"< Anterior",
		prevColor,
		nil,
		previousHandler,
	)

	nextButton := components.NewButton(
		basic.Point{},
		basic.Size{W: 150, H: 40},
		"Próximo >",
		nextColor,
		nil,
		nextHandler,
	)

	pagRow := components.NewRow(
		basic.Point{},
		10, 
		basic.Size{W: screenSize.W, H: 40},
		basic.Center,
		basic.Center,
		[]components.Widget{previousButton, nextButton},
	)

	paginationContainer := components.NewContainer(
		basic.Point{},
		basic.Size{W: screenSize.W, H: 50},
		0,
		nil,
		basic.Center,
		basic.Center,
		pagRow,
	)

	backButton := components.NewContainer(
		basic.Point{},
		basic.Size{W: screenSize.W, H: 60},
		0,
		nil,
		basic.Center,
		basic.Center,
		components.NewButton(
			basic.Point{},
			basic.Size{W: 400, H: 50},
			"Voltar ao menu",
			colors.Dark,
			nil,
			func(bt *components.Button) {
				m.stack.Pop()
			},
		),
	)

	m.layout = components.NewColumn(
		basic.Point{},
		15,           
		basic.Size{W: screenSize.W, H: screenSize.H},
		basic.Start,  
		basic.Center, 
		[]components.Widget{
			topSpacer,
			title,
			rankingContainer,
			paginationContainer,
			backButton,
		},
	)
}