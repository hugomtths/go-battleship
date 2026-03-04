package scenes

import (
	"github.com/allanjose001/go-battleship/game/components"
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/internal/entity"
	"github.com/hajimehoshi/ebiten/v2"
)

type HistoryCardTestScene struct {
	root components.LayoutWidget
}

func (h *HistoryCardTestScene) OnEnter(prev Scene, size basic.Size) {

	mockVitoria := entity.MatchResult{
		Win:               true,
		PlayerShots:       45,
		Hits:              17,
		HigherHitSequence: 4,
		Score:             3200,
		LostShips:         1,
		KilledShips:       5, // Afundou toda a frota inimiga
		Duration:          495000,
	}

	/*mockDerrota := entity.MatchResult{
		Win:               false,
		PlayerShots:       28,
		Hits:              6,
		HigherHitSequence: 1,
		Score:             450,
		LostShips:         5,      // Perdeu toda a frota
		KilledShips:       1,
		Duration:          320000,
	}*/

	h.root = components.NewColumn(basic.Point{},
		12,
		size,
		basic.Center, basic.Center,
		[]components.Widget{
			components.NewHistoryCard(basic.Point{}, basic.Size{
				W: size.W * 0.8, // larguinho
				H: 265,
			},
				mockVitoria,
			),
		},
	)
}

func (h *HistoryCardTestScene) OnExit(next Scene) {
}

func (h *HistoryCardTestScene) Update() error {
	h.root.Update(basic.Point{})
	return nil
}

func (h *HistoryCardTestScene) Draw(screen *ebiten.Image) {
	h.root.Draw(screen)
}
