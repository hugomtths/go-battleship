package setup

import (
	"math/rand"

	"github.com/allanjose001/go-battleship/game/shared/board"
	"github.com/allanjose001/go-battleship/game/shared/placement"
)

// RandomlyPlaceAIShips posiciona navios aleatoriamente em um tabuleiro.
// Útil para configurar o tabuleiro da IA.
func RandomlyPlaceAIShips(b *board.Board) []*placement.ShipPlacement {
	b.Clear()

	shipSizes := []int{6, 6, 4, 4, 3, 1} // Mesmos tamanhos do jogador
	var placements []*placement.ShipPlacement

	for _, sz := range shipSizes {
		for {
			row := rand.Intn(board.Rows)
			col := rand.Intn(board.Cols)
			or := board.Orientation(rand.Intn(2))

			if b.CanPlace(sz, row, col, or) {
				b.PlaceShip(sz, row, col, or)
				placements = append(placements, &placement.ShipPlacement{
					Size:        sz,
					X:           col,
					Y:           row,
					Orientation: or,
					Placed:      true,
				})
				break
			}
		}
	}
	return placements
}
