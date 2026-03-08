package scenes

import (

	"github.com/allanjose001/go-battleship/game/components"
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	"github.com/allanjose001/go-battleship/internal/entity"
	"github.com/hajimehoshi/ebiten/v2"
)

type MatchsHistory struct {
	layout      components.Widget
	currentPage int
	StackHandler
}

func (m *MatchsHistory) GetMusic() string {
	return "menus"
}

func (m *MatchsHistory) OnExit(_ Scene) {
	m.stack.ctx.CanPopOrPush = false
}

func (m *MatchsHistory) OnEnter(_ Scene, screenSize basic.Size) {
	m.init(screenSize)

	_ = m.Update()
	m.stack.ctx.CanPopOrPush = true
}

func (m *MatchsHistory) Update() error {
	if m.layout != nil {
		m.layout.Update(basic.Point{X: 0, Y: 0})
	}
	return nil
}

func (m *MatchsHistory) Draw(screen *ebiten.Image) {
	if m.layout != nil {
		m.layout.Draw(screen)
	}
}

func (m *MatchsHistory) init(screenSize basic.Size) {
	var allMatches []entity.MatchResult
	if m.stack != nil && m.stack.ctx.Profile != nil && m.stack.ctx.Profile.History != nil {
		allMatches = m.stack.ctx.Profile.History
	}

	itemsPerPage := 2
	start := m.currentPage * itemsPerPage
	end := start + itemsPerPage

	if start > len(allMatches) {
		start = len(allMatches)
	}
	if end > len(allMatches) {
		end = len(allMatches)
	}

	pageMatches := allMatches[start:end]

	hasPrevious := m.currentPage > 0
	hasNext := end < len(allMatches)

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

	topSpacer := components.NewContainer(
		basic.Point{},
		basic.Size{W: screenSize.W, H: 20},
		0, nil, basic.Center, basic.Center, nil,
	)

	title := components.NewContainer(
		basic.Point{},
		basic.Size{W: screenSize.W, H: 60},
		0, nil, basic.Center, basic.Center,
		components.NewText(basic.Point{}, "Histórico de Partidas", colors.White, 42),
	)

	var cards []components.Widget
	for _, match := range pageMatches {
		card := components.NewHistoryCard(
			basic.Point{},
			basic.Size{W: screenSize.W * 0.8, H: 265},
			match,
		)
		cards = append(cards, card)
	}

	// Altura fixa para a área de cards evita estouro do layout
	cardsAreaHeight := float32(265*itemsPerPage + 20*(itemsPerPage-1))

	cardsColumn := components.NewColumn(
		basic.Point{},
		20,
		basic.Size{W: screenSize.W, H: cardsAreaHeight},
		basic.Start,
		basic.Center,
		cards,
	)

	cardsContainer := components.NewContainer(
		basic.Point{},
		basic.Size{W: screenSize.W, H: cardsAreaHeight},
		0, nil,
		basic.Start, basic.Center,
		cardsColumn,
	)


	backButton := components.NewContainer(
		basic.Point{},
		basic.Size{W: screenSize.W, H: 60},
		0, nil, basic.Center, basic.Center,
		components.NewButton(
			basic.Point{},
			basic.Size{W: 400, H: 50},
			"Voltar ao perfil",
			colors.Dark,
			nil,
			func(bt *components.Button) {
				m.stack.Pop()
			},
		),
	)

	var mainWidgets []components.Widget
	mainWidgets = append(mainWidgets, topSpacer, title)
	mainWidgets = append(mainWidgets, cardsContainer)

	mainWidgets = append(mainWidgets, paginationContainer, backButton)

	m.layout = components.NewColumn(
		basic.Point{},
		15,
		basic.Size{W: screenSize.W, H: screenSize.H},
		basic.Start,
		basic.Center,
		mainWidgets,
	)
}
