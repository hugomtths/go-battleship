// BattleSetupService: inicializa as estruturas que a IA usa na batalha.
// Responsável por:
// - Criar a Fleet (entidade lógica de navios) usada pela IA
// - Mapear placements visuais do jogador (board.Board + ShipPlacement)
//   para a representação lógica (entity.Board + entity.Ship)
// - Instanciar o AIPlayer com a fleet resultante
package service

import (
	"github.com/allanjose001/go-battleship/game/shared/board"
	"github.com/allanjose001/go-battleship/game/shared/placement"
	"github.com/allanjose001/go-battleship/internal/ai"
	"github.com/allanjose001/go-battleship/internal/entity"
)

type BattleSetupService struct{}

func NewBattleSetupService() *BattleSetupService {
	return &BattleSetupService{}
}

// InitBattleAI:
// - Varre os navios posicionados pelo jogador e procura correspondentes na Fleet
// - Respeita orientação e posição (X/Y) ao colocar no entity.Board
// - Cria um AIPlayer “hard” com estratégias combinadas
// - Retorna o AIPlayer, o board lógico e a fleet que ele rastreia
func (s *BattleSetupService) InitBattleAI(playerShips []*placement.ShipPlacement) (*ai.AIPlayer, *entity.Board, *entity.Fleet) {
	fleet := entity.NewFleet()
	entityBoard := &entity.Board{}

	usedShips := make(map[int]bool)

	for _, ps := range playerShips {
		if !ps.Placed {
			continue
		}

		var entShip *entity.Ship
		for i, ship := range fleet.Ships {
			if !usedShips[i] && ship.Size == ps.Size {
				entShip = ship
				usedShips[i] = true
				break
			}
		}

		if entShip != nil {
			entShip.Horizontal = ps.Orientation == board.Horizontal
			entityBoard.PlaceShip(entShip, ps.Y, ps.X)
		}
	}

	aiPlayer := ai.NewHardAIPlayer(fleet)

	return aiPlayer, entityBoard, fleet
}
