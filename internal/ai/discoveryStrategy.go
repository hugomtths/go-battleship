package ai

import "github.com/allanjose001/go-battleship/internal/entity"
import "fmt"

type DiscoveryStrategy struct{}

func (s *DiscoveryStrategy) TryAttack(ai *AIPlayer, board *entity.Board) bool {
	fmt.Println("discoveryStrategy usada")

	if ai.IsChasing() {
		return false
	}
	if len(ai.priorityQueue) == 0 {
		return false
	}

	// Pega a primeira posição da fila de prioridade
	row, col := ai.PopPriority()
	ship := board.AttackPositionB(row, col)
	ai.AdjustStrategy(board, row, col, ship)
	return true
}
