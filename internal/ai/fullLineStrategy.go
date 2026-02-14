package ai

import (
	"fmt"

	"github.com/allanjose001/go-battleship/internal/entity"
)

type FullLineStrategy struct{}

func (s *FullLineStrategy) TryAttack(ai *AIPlayer, board *entity.Board) bool {
	fmt.Println("fullLineStrategy usada")

	if len(ai.priorityQueue) == 0 {
		return false
	}
	row, col := ai.PopPriority()
	ship := board.AttackPositionB(row, col)
	ai.AdjustStrategy(board, row, col, ship)
	if ship == nil || ship.IsDestroyed() {
		return true
	}

	// Detecta orientação baseada em acertos anteriores
	horizontal := false
	vertical := false
	if ai.IsValidForTesting(row, col-1) && ai.virtualBoard[row][col-1] == 2 { horizontal = true }
	if ai.IsValidForTesting(row, col+1) && ai.virtualBoard[row][col+1] == 2 { horizontal = true }
	if ai.IsValidForTesting(row-1, col) && ai.virtualBoard[row-1][col] == 2 { vertical = true }
	if ai.IsValidForTesting(row+1, col) && ai.virtualBoard[row+1][col] == 2 { vertical = true }

	ai.ClearPriorityQueue()

	if horizontal {
		// encontra o bloco contíguo de acertos à esquerda
		c := col
		for ai.IsValidForTesting(row, c-1) && ai.virtualBoard[row][c-1] == 2 {
			c--
		}
		// adiciona a célula imediatamente antes do bloco (se válida)
		if ai.IsValid(row, c-1) {
			ai.AddToPriorityQueue(row, c-1)
		}

		// encontra o bloco contíguo de acertos à direita
		c = col
		for ai.IsValidForTesting(row, c+1) && ai.virtualBoard[row][c+1] == 2 {
			c++
		}
		// adiciona a célula imediatamente depois do bloco (se válida)
		if ai.IsValid(row, c+1) {
			ai.AddToPriorityQueue(row, c+1)
		}

		ai.StartChase()
	} else if vertical {
		r := row
		for ai.IsValidForTesting(r-1, col) && ai.virtualBoard[r-1][col] == 2 {
			r--
		}
		if ai.IsValid(r-1, col) {
			ai.AddToPriorityQueue(r-1, col)
		}

		r = row
		for ai.IsValidForTesting(r+1, col) && ai.virtualBoard[r+1][col] == 2 {
			r++
		}
		if ai.IsValid(r+1, col) {
			ai.AddToPriorityQueue(r+1, col)
		}

		ai.StartChase()
	} else {
		ai.AttackNeighbors(row, col)
	}
	return true
}
