package ai

import (
	"fmt"
	"math/rand"

	"github.com/allanjose001/go-battleship/internal/entity"
)

// RandomMoveStrategy move aleatoriamente um navio da IA com uma certa probabilidade,
// sem depender de ter sido atacado. Usada no modo dinâmico.
type RandomMoveStrategy struct {
	// Chance de 0 a 100 de tentar mover um navio por turno (ex: 40 = 40%)
	Chance int
}

func (s *RandomMoveStrategy) TryAttack(ai *AIPlayer, board *entity.Board) bool {
	fmt.Println("chegou em randomMoveStrategy")

	if ai.ownBoard == nil {
		fmt.Println("randomMoveStrategy: ownBoard nil, pulando")
		return false
	}

	// Rolagem de dado: só age com a probabilidade definida
	chance := s.Chance
	if chance <= 0 {
		chance = 15 // padrão: 40%
	}
	if rand.Intn(100) >= chance {
		fmt.Println("randomMoveStrategy: não ativou neste turno (sorte)")
		return false
	}

	// Coleta todos os navios ainda vivos no próprio board da IA
	aliveShips := collectAliveShips(ai.ownBoard)
	if len(aliveShips) == 0 {
		fmt.Println("randomMoveStrategy: nenhum navio vivo encontrado")
		return false
	}

	// Embaralha para escolher um navio aleatório
	rand.Shuffle(len(aliveShips), func(i, j int) {
		aliveShips[i], aliveShips[j] = aliveShips[j], aliveShips[i]
	})

	// Tenta mover o navio escolhido em uma direção aleatória
	dirs := []entity.Direction{entity.Up, entity.Down, entity.Left, entity.Right}

	for _, ship := range aliveShips {
		topR, topC := findShipTopLeft(ai.ownBoard, ship)
		if topR == -1 {
			continue
		}

		// Embaralha direções para evitar viés
		rand.Shuffle(len(dirs), func(i, j int) { dirs[i], dirs[j] = dirs[j], dirs[i] })

		for _, dir := range dirs {
			dr, dc := dirToDeltas(dir)
			newRow := topR + dr
			newCol := topC + dc

			if err := ai.ownBoard.MoveShip(ship, newRow, newCol); err == nil {
				fmt.Printf("randomMoveStrategy: '%s' movido de (%d,%d) para (%d,%d)\n",
					ship.Name, topR, topC, newRow, newCol)
				return true // consumiu o turno: IA moveu, não ataca
			}
		}
	}

	fmt.Println("randomMoveStrategy: nenhum navio pôde ser movido")
	return false
}

// collectAliveShips percorre o board e retorna ponteiros únicos de navios não destruídos.
func collectAliveShips(b *entity.Board) []*entity.Ship {
	seen := make(map[*entity.Ship]bool)
	var result []*entity.Ship

	for r := 0; r < entity.BoardSize; r++ {
		for c := 0; c < entity.BoardSize; c++ {
			ship := entity.GetShipReference(b.Positions[r][c])
			if ship != nil && !ship.IsDestroyed() && !seen[ship] {
				seen[ship] = true
				result = append(result, ship)
			}
		}
	}
	return result
}
