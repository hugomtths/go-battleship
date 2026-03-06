package scenes

import (
	"github.com/allanjose001/go-battleship/game/components"
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	"github.com/allanjose001/go-battleship/internal/ai"
	"github.com/hajimehoshi/ebiten/v2"
)

type DifficultyMenu struct {
    layout   components.Widget
    onSelect func(player *ai.AIPlayer)
    StackHandler
}

func NewDifficultyMenu(w, h int, onSelect func(player *ai.AIPlayer)) *DifficultyMenu {
	btnSize := basic.Size{W: 220, H: 60}

	btnRecruta := components.NewButton(basic.Point{}, btnSize, "Recruta", colors.Blue, colors.White, func(b *components.Button) {
		onSelect(ai.NewEasyAIPlayer())
	})

	btnImediato := components.NewButton(basic.Point{}, btnSize, "Imediato", colors.Blue, colors.White, func(b *components.Button) {
		onSelect(ai.NewMediumAIPlayer(nil))
	})

	btnAlmirante := components.NewButton(basic.Point{}, btnSize, "Almirante", colors.Blue, colors.White, func(b *components.Button) {
		onSelect(ai.NewHardAIPlayer(nil))
	})

	screenSize := basic.Size{W: float32(w), H: float32(h)}
	column := components.NewColumn(
		basic.Point{X: 0, Y: 0},
		25,
		screenSize,
		basic.Center,
		basic.Center,
		[]components.Widget{
			components.NewText(basic.Point{}, "SELEÇÃO DE DIFICULDADE", colors.White, 28),
			btnRecruta, btnImediato, btnAlmirante,
		},
	)
	return &DifficultyMenu{layout: column, onSelect: onSelect}
}

func (m *DifficultyMenu) OnEnter(prev Scene, size basic.Size) {
    // layout já criado em NewDifficultyMenu; nothing else required
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

func (m *DifficultyMenu) Layout(w, h int) (int, int) { return w, h }
