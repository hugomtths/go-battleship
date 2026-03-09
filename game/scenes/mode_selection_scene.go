package scenes

import (
	"github.com/allanjose001/go-battleship/game/components"
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	"github.com/allanjose001/go-battleship/internal/entity"
	"github.com/hajimehoshi/ebiten/v2"
)

// ModeSelectionScene permite escolher entre Partida Clássica e Campanha.
type ModeSelectionScene struct {
	root components.LayoutWidget
	StackHandler
	profile *entity.Profile
}

func (m *ModeSelectionScene) GetMusic() string {
	return "menus"
}

func (m *ModeSelectionScene) OnEnter(prev Scene, size basic.Size) {
	// tenta obter profile do contexto (injetado pelo SceneStack)
	if m.ctx != nil && m.ctx.Profile != nil {
		m.profile = m.ctx.Profile
	}

	btnSize := basic.Size{W: 450, H: 50}
	// PARTIDA CLÁSSICA -> abre seleção de dificuldade
	classicBtn := components.NewButton(
		basic.Point{}, btnSize,
		"Partida Clássica", colors.Dark,
		nil,
		func(b *components.Button) {
			if m.ctx != nil {
				m.ctx.IsDynamicMode = false
				m.ctx.IsCampaign = false
			}
			m.ctx.SoundService.PlaySFX("click", 0.8)
			m.stack.Push(&DifficultyScene{})
		},
	)

	campaignBtn := components.NewButton(basic.Point{}, btnSize, "Campanha", colors.Dark, nil, func(b *components.Button) {
		// garante que o profile esteja no contexto
		if m.ctx != nil {
			m.ctx.IsDynamicMode = false
			m.ctx.IsCampaign = true
			if m.profile != nil {
				m.ctx.Profile = m.profile
			}
		}
		m.ctx.SoundService.PlaySFX("click", 0.8)
		m.stack.Push(&CampaignScene{})
	})

	dynamicBtn := components.NewButton(basic.Point{}, btnSize, "Dinâmico", colors.Dark, nil, func(b *components.Button) {
		// Modo Dinâmico: Dificuldade Hard, IsDynamicMode = true
		if m.ctx != nil {
			m.ctx.SetDifficulty("hard")
			m.ctx.IsCampaign = false
			m.ctx.IsDynamicMode = true
			if m.profile != nil {
				m.ctx.Profile = m.profile
			}
		}
		m.ctx.SoundService.PlaySFX("click", 0.8)
		m.stack.Push(NewPlacementSceneWithProfile(m.profile))
	})

	backBtn := components.NewButton(basic.Point{}, basic.Size{W: 220, H: 50}, "Voltar", colors.Dark, nil,
		func(b *components.Button) {
			if m.ctx.CanPopOrPush {
				m.ctx.SoundService.PlaySFX("backclick", 0.8)
				m.stack.Pop()
			}
		})
	spacer := components.NewContainer(
		basic.Point{}, basic.Size{W: 1, H: 20}, 0,
		colors.Transparent, basic.Center, basic.Center,
		nil,
	)
	spacer2 := components.NewContainer(
		basic.Point{}, basic.Size{W: 1, H: 100}, 0,
		colors.Transparent, basic.Center, basic.Center,
		nil,
	)

	screenSize := basic.Size{W: size.W, H: size.H}
	m.root = components.NewColumn(
		basic.Point{},
		20,
		screenSize,
		basic.Start,
		basic.Center,
		[]components.Widget{
			spacer,
			components.NewText(basic.Point{}, "Selecione o Modo de Jogo", colors.White, 35),
			spacer2,
			classicBtn,
			spacer,
			campaignBtn,
			spacer,
			dynamicBtn,
			spacer2,
			backBtn,
		},
	)
	m.stack.ctx.CanPopOrPush = true
	_ = m.Update()
}

func (m *ModeSelectionScene) OnExit(next Scene) {
	m.stack.ctx.CanPopOrPush = false
}

func (m *ModeSelectionScene) Update() error {
	if m.root != nil {
		m.root.Update(basic.Point{X: 0, Y: 0})
	}
	return nil
}

func (m *ModeSelectionScene) Draw(screen *ebiten.Image) {
	if m.root != nil {
		m.root.Draw(screen)
	}
}
