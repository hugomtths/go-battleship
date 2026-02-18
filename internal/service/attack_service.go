// AttackService: concentra a regra de combate.
// Responsável por aplicar ataques do jogador no board visual da IA
// e executar o turno da IA sincronizando o entity.Board com o board
// visual do jogador. Não orquestra turnos (isso é do BattleService).
package service

import (
	"github.com/allanjose001/go-battleship/game/shared/board"
	"github.com/allanjose001/go-battleship/internal/ai"
	"github.com/allanjose001/go-battleship/internal/entity"
)

type AttackService struct{}

func NewAttackService() *AttackService {
	return &AttackService{}
}

// PlayerAttack:
// - Ignora células já atacadas
// - Conta tentativa
// - Marca Hit/Miss conforme estado da célula
// - Retorna indicadores de acerto e fim de jogo usando totalShipCells
func (s *AttackService) PlayerAttack(aiBoard *board.Board, row, col int, attempts, hits, totalShipCells int) (int, int, bool, bool) {
	cell := &aiBoard.Cells[row][col]

	if cell.State == board.Hit || cell.State == board.Miss {
		return attempts, hits, false, false
	}

	attempts++

	if cell.State == board.Ship {
		cell.State = board.Hit
		hits++
		if hits >= totalShipCells {
			return attempts, hits, true, true
		}
		return attempts, hits, true, false
	}

	if cell.State == board.Empty {
		cell.State = board.Miss
	}

	return attempts, hits, false, false
}

// AITurn:
// - Pede para o AIPlayer atacar o entity.Board do jogador
// - Sincroniza esse ataque com o board visual (marcando Hit/Miss)
// - Checa fim de jogo com totalShipCells
func (s *AttackService) AITurn(aiPlayer *ai.AIPlayer, entityBoard *entity.Board, playerBoard *board.Board, attempts, hits, totalShipCells int) (int, int, bool) {
	if aiPlayer == nil {
		return attempts, hits, false
	}

	attempts++
	aiPlayer.Attack(entityBoard)

	for r := 0; r < board.Rows; r++ {
		for c := 0; c < board.Cols; c++ {
			entPos := entityBoard.Positions[r][c]
			cell := &playerBoard.Cells[r][c]

			if entity.IsAttacked(entPos) && cell.State != board.Hit && cell.State != board.Miss {
				if cell.State == board.Ship {
					cell.State = board.Hit
					hits++
					if hits >= totalShipCells {
						return attempts, hits, true
					}
				} else {
					cell.State = board.Miss
				}
			}
		}
	}

	return attempts, hits, false
}
