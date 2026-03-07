package scenes

import (
	"github.com/allanjose001/go-battleship/game/components"
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	"github.com/hajimehoshi/ebiten/v2"
)

type DifficultyMenu struct {
	layout   components.LayoutWidget
	onSelect func(difficulty string)
	StackHandler
}

func (m *DifficultyMenu) OnExit(_ Scene) {
	m.stack.ctx.CanPopOrPush = false
}

func (m *DifficultyMenu) GetMusic() string {
	return "menus"
}

func NewDifficultyMenu(w, h int, onSelect func(difficulty string), onBack func()) *DifficultyMenu {
	btnSize := basic.Size{W: 220, H: 60}

	btnRecruta := components.NewButton(basic.Point{}, btnSize, "Recruta", colors.Blue, colors.White, func(b *components.Button) {
		onSelect("easy")
	})

	btnImediato := components.NewButton(basic.Point{}, btnSize, "Imediato", colors.Blue, colors.White, func(b *components.Button) {
		onSelect("medium")
	})

	btnAlmirante := components.NewButton(basic.Point{}, btnSize, "Almirante", colors.Blue, colors.White, func(b *components.Button) {
		onSelect("hard")
	})

	// Botão Voltar
	btnVoltar := components.NewButton(
		basic.Point{},
		basic.Size{W: 220, H: 55},
		"Voltar",
		colors.Dark,
		colors.White,
		func(b *components.Button) {
			onBack()
		},
	)

	screenSize := basic.Size{W: float32(w), H: float32(h)}
	column := components.NewColumn(
		basic.Point{X: 0, Y: 0},
		25,
		screenSize,
		basic.Center,
		basic.Center,
		[]components.Widget{
			components.NewText(basic.Point{}, "SELEÇÃO DE DIFICULDADE", colors.White, 28),
			btnRecruta,
			btnImediato,
			btnAlmirante,
			btnVoltar,
		},
	)

	return &DifficultyMenu{layout: column, onSelect: onSelect}
}

func (m *DifficultyMenu) OnEnter(prev Scene, size basic.Size) {
	m.stack.ctx.CanPopOrPush = true
}

func (m *DifficultyMenu) Update() error {
	if m.layout != nil {
		m.layout.Update(basic.Point{X: 0, Y: 0})
	}
	return nil
}

func (m *DifficultyMenu) Draw(screen *ebiten.Image) {
	if m.layout != nil {
		m.layout.Draw(screen)
	}
}
