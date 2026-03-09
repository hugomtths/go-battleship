package scenes

import (
	"strings"

	"github.com/allanjose001/go-battleship/game/components"
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	"github.com/allanjose001/go-battleship/internal/entity"
	"github.com/allanjose001/go-battleship/internal/service"
	"github.com/hajimehoshi/ebiten/v2"
)

type CampaignHistoryScene struct {
	root components.Widget
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

func (s *CampaignHistoryScene) GetMusic() string {
	return "menus"
}

func (s *CampaignHistoryScene) OnEnter(prev Scene, size basic.Size) {
	s.init(size)
	_ = s.Update()
	s.stack.ctx.CanPopOrPush = true
}

func (s *CampaignHistoryScene) init(size basic.Size) {
	// Atualiza o perfil com os dados mais recentes do serviço
	if s.ctx != nil && s.ctx.Profile != nil {
		if p, err := service.FindProfile(s.ctx.Profile.Username); err == nil {
			s.ctx.Profile = p
		}
	}

	profile := s.ctx.Profile
	var results []entity.MatchResult

	// Coleta histórico geral filtrando por modo Campanha e dificuldade selecionada
	for _, match := range profile.History {
		if strings.Contains(match.Mode, "Campanha") && match.Difficulty == s.difficultyKey {
			results = append(results, match)
		}
	}

	// Inverte a ordem para mostrar os mais recentes primeiro
	for i, j := 0, len(results)-1; i < j; i, j = i+1, j-1 {
		results[i], results[j] = results[j], results[i]
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
	titleContainer := components.NewContainer(
		basic.Point{},
		basic.Size{W: size.W, H: 80},
		0, nil,
		basic.Center, basic.Center,
		components.NewText(basic.Point{}, "Histórico - "+s.difficultyName, colors.White, 32),
	)

	var cards []components.Widget
	for _, res := range pageResults {
		card := components.NewHistoryCard(
			basic.Point{},
			basic.Size{W: size.W * 0.9, H: 240},
			res,
		)
		cards = append(cards, card)
	}

	// Se não houver histórico
	if len(results) == 0 {
		cards = append(cards, components.NewText(basic.Point{}, "Nenhuma partida registrada.", colors.White, 24))
	}

	// Calcula altura disponível para a lista para não empurrar os botões para fora
	// Altura total - Título(80) - Paginação(60) - Voltar(80) - Espaçamentos(~30)
	listHeight := size.H - 250

	listColumn := components.NewColumn(
		basic.Point{},
		20,
		basic.Size{W: size.W, H: listHeight},
		basic.Start, // Alinha itens no topo
		basic.Center,
		cards,
	)

	listContainer := components.NewContainer(
		basic.Point{},
		basic.Size{W: size.W, H: listHeight},
		0, nil,
		basic.Start, basic.Center,
		listColumn,
	)

	// Botões de navegação
	hasPrev := s.currentPage > 0
	hasNext := end < len(results)

	var previousHandler func(*components.Button)
	var nextHandler func(*components.Button)

	prevColor := colors.Dark
	nextColor := colors.Dark

	if hasPrev {
		previousHandler = func(b *components.Button) {
			s.ctx.SoundService.PlaySFX("click", 0.8)
			s.currentPage--
			s.init(size)
		}
	} else {
		previousHandler = nil
		prevColor = colors.NightBlue
	}

	if hasNext {
		nextHandler = func(b *components.Button) {
			s.ctx.SoundService.PlaySFX("click", 0.8)
			s.currentPage++
			s.init(size)
		}
	} else {
		nextHandler = nil
		nextColor = colors.NightBlue
	}

	prevBtn := components.NewButton(basic.Point{}, basic.Size{W: 150, H: 40}, "< Anterior", prevColor, nil, previousHandler)
	nextBtn := components.NewButton(basic.Point{}, basic.Size{W: 150, H: 40}, "Próximo >", nextColor, nil, nextHandler)

	pagRow := components.NewRow(
		basic.Point{},
		10,
		basic.Size{W: size.W, H: 40},
		basic.Center,
		basic.Center,
		[]components.Widget{prevBtn, nextBtn},
	)

	paginationContainer := components.NewContainer(
		basic.Point{},
		basic.Size{W: size.W, H: 50},
		0,
		nil,
		basic.Center,
		basic.Center,
		pagRow,
	)

	backBtnContainer := components.NewContainer(
		basic.Point{},
		basic.Size{W: size.W, H: 80},
		0,
		nil,
		basic.Center,
		basic.Center,
		components.NewButton(basic.Point{}, basic.Size{W: 400, H: 50}, "Voltar", colors.Dark, nil, func(b *components.Button) {
			s.ctx.SoundService.PlaySFX("backclick", 0.8)
			s.stack.Pop()
		}),
	)

	s.root = components.NewColumn(
		basic.Point{},
		10,
		size,
		basic.Start, // Começa do topo da tela
		basic.Center,
		[]components.Widget{
			titleContainer,
			listContainer,
			paginationContainer,
			backBtnContainer,
		},
	)
}

func (s *CampaignHistoryScene) OnExit(next Scene) {
	s.stack.ctx.CanPopOrPush = false
}

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
