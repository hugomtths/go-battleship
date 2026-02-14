package ai

import ( "math/rand" 
		"github.com/allanjose001/go-battleship/internal/entity"
		"fmt"
)

type RandomStrategy struct {}

const boardSize = 10;

func (s *RandomStrategy) TryAttack(ai *AIPlayer, board *entity.Board) bool {
	
	fmt.Println("randomStrategy usada")

	for {
		row := rand.Intn(boardSize);
		col := rand.Intn(boardSize);

		if ai.IsValid(row, col) {
			ship := board.AttackPositionB(row, col)
			ai.AdjustStrategy(board, row, col, ship)
			return true;
		}
	}
}