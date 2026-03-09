package ai

import "github.com/allanjose001/go-battleship/internal/entity"
import "fmt"

type StrategicSearchStrategy struct{}

func (s *StrategicSearchStrategy) TryAttack(ai *AIPlayer, board *entity.Board) bool {
	fmt.Println("chegou em StrategicShearch")
	if len(ai.priorityQueue) != 0 {
		fmt.Println("SSS: skipping because priorityQueue not empty")
		return false
	}

	should := ai.ShouldAttackStrategicPositions()
	fmt.Printf("SSS: ShouldAttackStrategicPositions=%v\n", should)
	if !should {
		return false
	}
	//if !ai.ShouldAttackStrategicPositions() {
	//    return false
	//}
	size := ai.SizeOfNextShip()
	fmt.Printf("SSS: SizeOfNextShip=%d\n", size)
	if size == 0 {
		fmt.Println("SSS: no ships left according to AI enemyFleet")
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
