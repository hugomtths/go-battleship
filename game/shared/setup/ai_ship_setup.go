package setup

import (
	"math/rand"

	"github.com/allanjose001/go-battleship/game/shared/board"
)

// RandomlyPlaceAIShips posiciona navios aleatoriamente em um tabuleiro.
// Útil para configurar o tabuleiro da IA.
func RandomlyPlaceAIShips(b *board.Board) {
	b.Clear()

	shipSizes := []int{6, 4, 3, 3, 1} // Mesmos tamanhos do jogador

	for _, sz := range shipSizes {
		for {
			row := rand.Intn(board.Rows)
			col := rand.Intn(board.Cols)
			or := board.Orientation(rand.Intn(2))

			if b.CanPlace(sz, row, col, or) {
				b.PlaceShip(sz, row, col, or)
				break
			}
		}
	}
}
