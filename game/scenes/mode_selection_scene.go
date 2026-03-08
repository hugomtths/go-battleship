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

	btnSize := basic.Size{W: 300, H: 60}
	// PARTIDA CLÁSSICA -> abre seleção de dificuldade
	classicBtn := components.NewButton(basic.Point{}, btnSize, "Partida Clássica", colors.Dark, nil, func(b *components.Button) {
		if m.ctx != nil {
			m.ctx.IsDynamicMode = false
			m.ctx.IsCampaign = false
		}
		m.stack.Push(&DifficultyScene{})
	})

	campaignBtn := components.NewButton(basic.Point{}, btnSize, "Campanha", colors.Dark, nil, func(b *components.Button) {
		// garante que o profile esteja no contexto
		if m.ctx != nil {
			m.ctx.IsDynamicMode = false
			m.ctx.IsCampaign = true
			if m.profile != nil {
				m.ctx.Profile = m.profile
			}
		}
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
		m.stack.Push(NewPlacementSceneWithProfile(m.profile))
	})

	backBtn := components.NewButton(basic.Point{}, basic.Size{W: 220, H: 50}, "Voltar", colors.Dark, nil, func(b *components.Button) {
		m.stack.Pop()
	})

	screenSize := basic.Size{W: size.W, H: size.H}
	m.root = components.NewColumn(
		basic.Point{},
		20,
		screenSize,
		basic.Center,
		basic.Center,
		[]components.Widget{
			components.NewText(basic.Point{}, "SELECIONE O MODO DE JOGO", colors.White, 28),
			classicBtn,
			campaignBtn,
			dynamicBtn,
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

// DifficultyScene: wrapper que usa DifficultyMenu (já existente) e transforma em Scene.
// Ao selecionar dificuldade, encaminha para PlacementScene (a lógica de campaign start fica na fase de Battle/Match).
type DifficultyScene struct {
	StackHandler
	menu *DifficultyMenu
}

func (d *DifficultyScene) GetMusic() string {
	return "menus"
}

func (d *DifficultyScene) OnEnter(prev Scene, size basic.Size) {
	d.menu = NewDifficultyMenu(int(size.W), int(size.H), func(diff string) {

		var prof *entity.Profile
		if d.ctx != nil {
			d.ctx.SetDifficulty(diff)
			d.ctx.IsCampaign = false
			prof = d.ctx.Profile
		}

		d.stack.Push(NewPlacementSceneWithProfile(prof))
	}, d.stack.Pop)
	_ = d.menu.Update()
	d.stack.ctx.CanPopOrPush = true

}

func (d *DifficultyScene) OnExit(next Scene) {
	d.stack.ctx.CanPopOrPush = false
}

func (d *DifficultyScene) Update() error {

	if d.menu != nil {
		err := d.menu.Update()
		return err
	}
	return nil
}

func (d *DifficultyScene) Draw(screen *ebiten.Image) {
	if d.menu != nil {
		d.menu.Draw(screen)
	}
}
