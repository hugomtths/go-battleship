package ai

import "github.com/allanjose001/go-battleship/internal/entity"

type PartialLineStrategy struct{}

func (s *PartialLineStrategy) TryAttack(ai *AIPlayer, board *entity.Board) bool {
	if len(ai.priorityQueue) == 0 {
		return false
	}

	row, col := ai.PopPriority()
	ship := board.AttackPositionB(row, col)

	// centraliza a atualização de estado
	ai.AdjustStrategy(board, row, col, ship)

	// se não acertou nada, não tenta linha
	if ship == nil || ship.IsDestroyed() {
		return true
	}

	// tenta descobrir orientação parcial
	horizontal := false
	vertical := false

	if ai.IsValidForTesting(row, col-1) && ai.virtualBoard[row][col-1] == 2 {
		horizontal = true
	}
	if ai.IsValidForTesting(row, col+1) && ai.virtualBoard[row][col+1] == 2 {
		horizontal = true
	}
	if ai.IsValidForTesting(row-1, col) && ai.virtualBoard[row-1][col] == 2 {
		vertical = true
	}
	if ai.IsValidForTesting(row+1, col) && ai.virtualBoard[row+1][col] == 2 {
		vertical = true
	}

	ai.ClearPriorityQueue()

	if horizontal {
		ai.AddToPriorityQueue(row, col-1)
		ai.AddToPriorityQueue(row, col+1)
	} else if vertical {
		ai.AddToPriorityQueue(row-1, col)
		ai.AddToPriorityQueue(row+1, col)
	} else {
		ai.AttackNeighbors(row, col)
	}

	return true
}

