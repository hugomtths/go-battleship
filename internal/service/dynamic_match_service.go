package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/allanjose001/go-battleship/game/scenes/audio"
	"github.com/allanjose001/go-battleship/game/shared/board"
	"github.com/allanjose001/go-battleship/internal/entity"
)

type DynamicMatchService struct {
	*MatchService
}

// Corrigido: recebe attack e aiDelay (mesma semântica de NewMatchService)
func NewDynamicMatchService(attack *AttackService, aiDelay time.Duration, ss *audio.SoundService) *DynamicMatchService {
	return &DynamicMatchService{
		MatchService: NewMatchService(attack, aiDelay, ss),
	}
}

// MovePlayerShip tenta mover `ship` do jogador para (newRow,newCol).
// now: tempo atual usado para agendamento do próximo passo da IA.
// Observação: MoveShip foi implementado em internal/entity/Board (PlayerEntityBoard),
// portanto chamamos ali. Se quiser refletir no PlayerBoard (visual), sincronize depois.
func (s *DynamicMatchService) MovePlayerShip(m *entity.Match, ship *entity.Ship, newRow int, newCol int, now time.Time) error {
	if m == nil {
		return ErrMatchNotFound
	}
	if m.IsFinished() {
		return ErrMatchFinished
	}
	if m.Status != entity.MatchStatusInProgress {
		return ErrMatchNotInProgress
	}
	if m.Turn != entity.TurnPlayer {
		return ErrNotPlayersTurn
	}
	// precisa das referências runtime presentes
	if m.PlayerEntityBoard == nil || m.PlayerBoard == nil {
		return ErrMatchNotReady
	}
	if ship == nil {
		return errors.New("ship nil")
	}

	// checa se ship pertence à frota do jogador
	found := false
	if m.PlayerFleet != nil {
		for _, sh := range m.PlayerFleet.GetFleetShips() {
			if sh == ship {
				found = true
				break
			}
		}
	}
	if !found {
		return fmt.Errorf("ship not in player fleet")
	}

	// delega a movimentação para PlayerEntityBoard (onde MoveShip existe)
	if err := m.PlayerEntityBoard.MoveShip(ship, newRow, newCol); err != nil {
		return err
	}

	// sincroniza o board visual para refletir a nova posição dos navios
	s.syncVisualShipPositions(m)

	// consumir o turno do jogador: passa para IA e agenda próximo ataque
	m.Turn = entity.TurnEnemy
	m.NextAction = entity.NextActionEnemyAttack
	m.NextActionAt = now.Add(s.aiDelay)

	return nil
}

// syncVisualShipPositions garante que o board visual (usado pelo AttackService e Renderer)
// esteja em sincronia com o board lógico (onde os navios realmente se movem).
func (s *DynamicMatchService) syncVisualShipPositions(m *entity.Match) {
	if m.PlayerBoard == nil || m.PlayerEntityBoard == nil {
		return
	}

	for r := 0; r < 10; r++ {
		for c := 0; c < 10; c++ {
			entPos := m.PlayerEntityBoard.Positions[r][c]
			cell := &m.PlayerBoard.Cells[r][c]

			// Se a célula já foi atacada (Hit/Miss no visual), não alteramos seu estado.
			// O AttackService já cuida de manter Hit/Miss sincronizado.
			// Aqui só nos importamos com células não atacadas que podem ter ganhado ou perdido um navio.
			if cell.State != board.Hit && cell.State != board.Miss {
				if entity.GetShipReference(entPos) != nil {
					cell.State = board.Ship
				} else {
					cell.State = board.Empty
				}
			}
		}
	}
}
