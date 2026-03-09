package ai

import (
	"fmt"
	"math/rand"

	"github.com/allanjose001/go-battleship/internal/entity"
)

type EvasionStrategy struct{}

func (s *EvasionStrategy) TryAttack(ai *AIPlayer, board *entity.Board) bool {
	fmt.Println("chegou em evasionStrategy")

	if ai.ownBoard == nil || len(ai.evasionQueue) == 0 {
		fmt.Println("evasionStrategy: sem navios para evadir")
		return false
	}

	ship := ai.dequeueEvasion()

	// Navio destruído entre o hit e o turno da IA: descarta silenciosamente
	if ship.IsDestroyed() {
		fmt.Printf("evasionStrategy: '%s' destruído antes da evasão, descartando\n", ship.Name)
		return false
	}

	topR, topC := findShipTopLeft(ai.ownBoard, ship)
	if topR == -1 {
		fmt.Printf("evasionStrategy: '%s' não encontrado no board\n", ship.Name)
		return false
	}

	dirs := []entity.Direction{entity.Up, entity.Down, entity.Left, entity.Right}
	rand.Shuffle(len(dirs), func(i, j int) { dirs[i], dirs[j] = dirs[j], dirs[i] })

	for _, dir := range dirs {
		dr, dc := dirToDeltas(dir)
		newRow := topR + dr
		newCol := topC + dc

		if err := ai.ownBoard.MoveShip(ship, newRow, newCol); err == nil {
			fmt.Printf("evasionStrategy: '%s' movido de (%d,%d) para (%d,%d)\n",
				ship.Name, topR, topC, newRow, newCol)
			return true // <- consome o turno: IA moveu, não ataca
		}
	}

	fmt.Printf("evasionStrategy: '%s' não pôde ser movido (sem espaço), atacando normalmente\n", ship.Name)
	return false // <- só cai aqui se nenhuma direção foi possível
}

func findShipTopLeft(b *entity.Board, ship *entity.Ship) (int, int) {
	minR, minC := -1, -1
	for r := 0; r < entity.BoardSize; r++ {
		for c := 0; c < entity.BoardSize; c++ {
			if entity.GetShipReference(b.Positions[r][c]) == ship {
				if minR == -1 || r < minR || (r == minR && c < minC) {
					minR, minC = r, c
				}
			}
		}
	}
	return minR, minC
}

func dirToDeltas(dir entity.Direction) (int, int) {
	switch dir {
	case entity.Up:
		return -1, 0
	case entity.Down:
		return 1, 0
	case entity.Left:
		return 0, -1
	case entity.Right:
		return 0, 1
	}
	return 0, 0
}
