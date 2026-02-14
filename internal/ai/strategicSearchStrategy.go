package ai

import "github.com/allanjose001/go-battleship/internal/entity"
import "fmt"

type StrategicSearchStrategy struct{}

func (s *StrategicSearchStrategy) TryAttack(ai *AIPlayer, board *entity.Board) bool {
	fmt.Println("strategicSearchStrategy usada")

    if len(ai.priorityQueue) != 0 {
        return false
    }
    if !ai.ShouldAttackStrategicPositions() {
        return false
    }
    size := ai.SizeOfNextShip()
    if size == 0 {
        return false
    }
    // tenta vertical primeiro ou horizontal aleatoriamente
    if ai.SearchVertically(size) {
        return true
    }
    if ai.SearchHorizontally(size) {
        return true
    }
    return false
}