package scenes

import (
	"github.com/allanjose001/go-battleship/game/components"
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	"github.com/allanjose001/go-battleship/internal/entity"
	"github.com/hajimehoshi/ebiten/v2"
)

type CampaignHistoryScene struct {
	root           components.Widget
	StackHandler
	difficultyKey  string
	difficultyName string
	currentPage    int
}

func NewCampaignHistoryScene(diffKey, diffName string) *CampaignHistoryScene {
	return &CampaignHistoryScene{
		difficultyKey:  diffKey,
		difficultyName: diffName,
	}
}

func (s *CampaignHistoryScene) OnEnter(prev Scene, size basic.Size) {
	s.init(size)
}

func (s *CampaignHistoryScene) init(size basic.Size) {
	profile := s.ctx.Profile
	var results []entity.MatchResult

	// Coleta resultados de campanhas passadas
	for _, c := range profile.Campaigns {
		if res, ok := c.DifficultyStep[s.difficultyKey]; ok {
			results = append(results, res)
		}
	}
	// Coleta resultado da campanha atual
	if profile.CurrentCampaign != nil {
		if res, ok := profile.CurrentCampaign.DifficultyStep[s.difficultyKey]; ok {
			results = append(results, res)
		}
	}

	// Paginação
	itemsPerPage := 2 
	start := s.currentPage * itemsPerPage
	end := start + itemsPerPage
	if start > len(results) {
		start = len(results)
	}
	if end > len(results) {
		end = len(results)
	}
	pageResults := results[start:end]

	// Construção da UI
	title := components.NewText(basic.Point{}, "Histórico - "+s.difficultyName, colors.White, 32)

	var cards []components.Widget
	for _, res := range pageResults {
		card := components.NewHistoryCard(
			basic.Point{},
			basic.Size{W: size.W * 0.6, H: 200}, // Tamanho reduzido
			res,
		)
		cards = append(cards, card)
	}

	// Se não houver histórico
	if len(results) == 0 {
		cards = append(cards, components.NewText(basic.Point{}, "Nenhuma partida registrada.", colors.White, 24))
	}

	listColumn := components.NewColumn(
		basic.Point{},
		20,
		basic.Size{W: size.W, H: size.H * 0.7},
		basic.Center,
		basic.Center,
		cards,
	)

	// Botões de navegação
	hasPrev := s.currentPage > 0
	hasNext := end < len(results)

	prevBtn := components.NewButton(basic.Point{}, basic.Size{W: 120, H: 40}, "< Ant", colors.Dark, nil, func(b *components.Button) {
		if hasPrev {
			s.currentPage--
			s.init(size)
		}
	})
	if !hasPrev { prevBtn.SetDisabled(true) }

	nextBtn := components.NewButton(basic.Point{}, basic.Size{W: 120, H: 40}, "Prox >", colors.Dark, nil, func(b *components.Button) {
		if hasNext {
			s.currentPage++
			s.init(size)
		}
	})
	if !hasNext { nextBtn.SetDisabled(true) }

	navRow := components.NewRow(basic.Point{}, 20, basic.Size{W: size.W, H: 50}, basic.Center, basic.Center, []components.Widget{prevBtn, nextBtn})

	backBtn := components.NewButton(basic.Point{}, basic.Size{W: 200, H: 50}, "Voltar", colors.Dark, nil, func(b *components.Button) {
		s.stack.Pop()
	})

	s.root = components.NewColumn(
		basic.Point{},
		20,
		size,
		basic.Center,
		basic.Center,
		[]components.Widget{
			title,
			listColumn,
			navRow,
			backBtn,
		},
	)
}

func (s *CampaignHistoryScene) OnExit(next Scene) {}

func (s *CampaignHistoryScene) Update() error {
	if s.root != nil {
		s.root.Update(basic.Point{})
	}
	return nil
}

func (s *CampaignHistoryScene) Draw(screen *ebiten.Image) {
	if s.root != nil {
		s.root.Draw(screen)
	}
}