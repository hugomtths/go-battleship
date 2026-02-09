package ai

import ( "math/rand" 
		"github.com/allanjose001/go-battleship/internal/entity"
)

type RandomStrategy struct {}

const boardSize = 10;

func (s *RandomStrategy) TryAttack(ai *AIPlayer, board *entity.Board) bool {
	
	for {
		x := rand.Intn(boardSize);
		y := rand.Intn(boardSize);

		if ai.IsValid(x, y) {
			ship := board.AttackPositionB(x, y)
			ai.AdjustStrategy(board, x, y, ship)
			return true;
		}
	}
}