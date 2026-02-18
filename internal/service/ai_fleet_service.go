// AIFleetService: posiciona navios da IA em um board lógico.
// Responsável por:
// - Resetar o entity.Board
// - Tentar colocar cada entity.Ship aleatoriamente sem colisões
// - Definir orientação horizontal/vertical por sorteio
// É uma peça de domínio da IA (não visual).
package service

import (
	"math/rand"

	"github.com/allanjose001/go-battleship/internal/entity"
)

type AIFleetService struct{}

func NewAIFleetService() *AIFleetService {
	return &AIFleetService{}
}

// PositionShipsRandomly:
// - Limpa o board lógico
// - Itera sobre a fleet tentando colocar cada navio
// - Recomeça se algum navio não couber após várias tentativas
func (s *AIFleetService) PositionShipsRandomly(b *entity.Board, f *entity.Fleet) {
	for {
		s.resetBoard(b)

		success := true
		for _, ship := range f.Ships {
			if !s.placeRandomShip(b, ship) {
				success = false
				break
			}
		}

		if success {
			break
		}
	}
}

// resetBoard: volta todas as posições ao estado inicial
func (s *AIFleetService) resetBoard(b *entity.Board) {
	for i := 0; i < entity.BoardSize; i++ {
		for j := 0; j < entity.BoardSize; j++ {
			b.Positions[i][j] = entity.Position{}
		}
	}
}

// placeRandomShip: sorteia posição e orientação e tenta colocar o navio
func (s *AIFleetService) placeRandomShip(b *entity.Board, ship *entity.Ship) bool {
	if ship == nil {
		return true
	}

	for attempts := 0; attempts < 1000; attempts++ {
		ship.Horizontal = rand.Intn(2) == 0

		var row, col int
		if ship.Horizontal {
			row = rand.Intn(entity.BoardSize)
			col = rand.Intn(entity.BoardSize - ship.Size + 1)
		} else {
			row = rand.Intn(entity.BoardSize - ship.Size + 1)
			col = rand.Intn(entity.BoardSize)
		}

		if b.PlaceShip(ship, row, col) {
			return true
		}
	}
	return false
}
