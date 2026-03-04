package service

import (
	"errors"
	"time"

	"github.com/allanjose001/go-battleship/game/shared/board"
	"github.com/allanjose001/go-battleship/internal/ai"
	"github.com/allanjose001/go-battleship/internal/entity"
)

var (
	// ErrMatchNotFound indica que o Match em memória não foi fornecido (nil).
	ErrMatchNotFound = errors.New("match not found")

	// ErrMatchFinished indica que a partida já terminou e não aceita novas ações.
	ErrMatchFinished = errors.New("match already finished")

	// ErrMatchNotInProgress indica que a partida ainda não foi iniciada.
	ErrMatchNotInProgress = errors.New("match is not in progress")

	// ErrNotPlayersTurn indica tentativa de ataque do jogador fora do turno do jogador.
	ErrNotPlayersTurn = errors.New("not player's turn")

	// ErrNotEnemyTurn indica tentativa de execução do passo da IA fora do turno da IA.
	ErrNotEnemyTurn = errors.New("not enemy's turn")

	// ErrActionNotReady indica que o “timer” ainda não liberou o passo de ataque da IA.
	ErrActionNotReady = errors.New("next action not ready")

	// ErrNoEnemyAttackSched indica que não existe um ataque agendado para a IA.
	ErrNoEnemyAttackSched = errors.New("no enemy attack scheduled")

	// ErrMatchNotReady indica que referências runtime (boards/board lógico/fleet) não foram injetadas.
	ErrMatchNotReady = errors.New("match runtime references not set")
)

type MatchService struct {
	attack  *AttackService
	aiDelay time.Duration
}

// NewMatchService cria um MatchService.
//
// attack: pode ser nil; se nil, usa NewAttackService().
// aiDelay: delay mínimo entre ataques da IA; se <= 0, usa 1s.
func NewMatchService(attack *AttackService, aiDelay time.Duration) *MatchService {
	if aiDelay <= 0 {
		aiDelay = time.Second
	}
	if attack == nil {
		attack = NewAttackService()
	}
	return &MatchService{
		attack:  attack,
		aiDelay: aiDelay,
	}
}

// Create cria um Match em memória.
// Como o Match não é persistido, este método apenas devolve um novo Match.
func (s *MatchService) Create(id string) *entity.Match {
	return entity.NewMatch(id, nil, nil, nil, nil)
}

// Start inicializa o Match e injeta referências runtime necessárias para jogar.
//
// Importante: como a IA ataca o PLAYER, a condição de vitória da IA deve usar
// TotalPlayerShipCells (não TotalEnemyShipCells).
//
// enemyShipCells  -> total de células dos navios no tabuleiro inimigo (vitória do player).
// playerShipCells -> total de células dos navios no tabuleiro do player (vitória da IA).
func (s *MatchService) Start(
	m *entity.Match,
	now time.Time,
	playerBoard *board.Board,
	enemyBoard *board.Board,
	playerEntityBoard *entity.Board,
	playerFleet *entity.Fleet,
	enemyShipCells int,
	playerShipCells int,
) error {
	if m == nil {
		return ErrMatchNotFound
	}
	if m.IsFinished() {
		return ErrMatchFinished
	}

	// refs runtime (não persistem)
	m.PlayerBoard = playerBoard
	m.EnemyBoard = enemyBoard
	m.PlayerEntityBoard = playerEntityBoard
	m.PlayerFleet = playerFleet

	// totais para condição de vitória
	m.TotalEnemyShipCells = enemyShipCells
	m.TotalPlayerShipCells = playerShipCells

	// reseta status/turn/stats/agenda
	m.Start(now)
	return nil
}

// PlayerAttack aplica o ataque do jogador no tabuleiro da IA.
// Regras (lógica C++):
// - Se a célula já foi atacada: retorna ErrInvalidAttackCell e NÃO consome turno.
// - Se HIT: jogador continua.
// - Se MISS: turno passa para IA e agenda próximo ataque em now+aiDelay.
// - Se o ataque encerrar a partida: finaliza o match e salva MatchResult no repo.
func (s *MatchService) PlayerAttack(m *entity.Match, now time.Time, row, col int) (entity.AttackEvent, error) {
	if err := s.validatePlayerAttack(m, row, col); err != nil {
		if errors.Is(err, entity.ErrInvalidAttackCell) {
			return entity.AttackEvent{
				Attacker: entity.TurnPlayer,
				Row:      row,
				Col:      col,
				Valid:    false,
				Hit:      false,
				GameOver: false,
			}, err
		}
		return entity.AttackEvent{}, err
	}

	hit, gameOver := s.applyPlayerAttack(m, row, col)
	ev := s.makeEvent(entity.TurnPlayer, row, col, true, hit)

	if err := s.postPlayerAttack(m, now, hit, gameOver, &ev); err != nil {
		return ev, err
	}
	return ev, nil
}

// EnemyAttackStep executa UM ataque da IA quando o schedule estiver liberado.
// Regras (lógica C++):
// - Se HIT: IA continua e agenda novo ataque em now+aiDelay.
// - Se MISS: devolve turno ao jogador.
// - Se encerrar a partida: finaliza o match e salva MatchResult no repo.
func (s *MatchService) EnemyAttackStep(m *entity.Match, now time.Time, aiPlayer *ai.AIPlayer) (entity.AttackEvent, error) {
	if err := s.validateEnemyStep(m, now, aiPlayer); err != nil {
		return entity.AttackEvent{}, err
	}

	hit, gameOver := s.applyEnemyStep(m, aiPlayer)
	ev := s.makeEvent(entity.TurnEnemy, -1, -1, true, hit)

	if err := s.postEnemyStep(m, now, hit, gameOver, &ev); err != nil {
		return ev, err
	}
	return ev, nil
}

//
// -------------------------- Helpers privados --------------------------
//

func (s *MatchService) validatePlayerAttack(m *entity.Match, row, col int) error {
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
	if m.EnemyBoard == nil {
		return ErrMatchNotReady
	}

	if row < 0 || row >= board.Rows || col < 0 || col >= board.Cols {
		return entity.ErrInvalidAttackCell
	}

	cell := &m.EnemyBoard.Cells[row][col]
	if cell.State == board.Hit || cell.State == board.Miss {
		return entity.ErrInvalidAttackCell
	}

	return nil
}

func (s *MatchService) applyPlayerAttack(m *entity.Match, row, col int) (hit bool, gameOver bool) {
	m.PlayerShots, m.PlayerHits, hit, gameOver =
		s.attack.PlayerAttack(m.EnemyBoard, row, col, m.PlayerShots, m.PlayerHits, m.TotalEnemyShipCells)

	if hit {
		m.PlayerHitStreak++
		if m.PlayerHitStreak > m.PlayerMaxHitStreak {
			m.PlayerMaxHitStreak = m.PlayerHitStreak
		}
	} else {
		m.PlayerHitStreak = 0
	}
	return hit, gameOver
}

func (s *MatchService) postPlayerAttack(m *entity.Match, now time.Time, hit, gameOver bool, ev *entity.AttackEvent) error {
	if gameOver {
		s.finishAndFillWinner(m, now, entity.TurnPlayer, ev)
		return nil
	}

	if !hit {
		m.Turn = entity.TurnEnemy
		m.NextAction = entity.NextActionEnemyAttack
		m.NextActionAt = now.Add(s.aiDelay)
	}
	return nil
}

func (s *MatchService) validateEnemyStep(m *entity.Match, now time.Time, aiPlayer *ai.AIPlayer) error {
	if m == nil {
		return ErrMatchNotFound
	}
	if m.IsFinished() {
		return ErrMatchFinished
	}
	if m.Status != entity.MatchStatusInProgress {
		return ErrMatchNotInProgress
	}
	if m.Turn != entity.TurnEnemy {
		return ErrNotEnemyTurn
	}
	if m.NextAction != entity.NextActionEnemyAttack {
		return ErrNoEnemyAttackSched
	}
	if now.Before(m.NextActionAt) {
		return ErrActionNotReady
	}
	if aiPlayer == nil || m.PlayerEntityBoard == nil || m.PlayerBoard == nil {
		return ErrMatchNotReady
	}
	return nil
}

func (s *MatchService) applyEnemyStep(m *entity.Match, aiPlayer *ai.AIPlayer) (hit bool, gameOver bool) {
	// consome schedule (evita execução duplicada)
	m.ClearNextAction()

	prevHits := m.EnemyHits

	// CORREÇÃO: a IA vence quando EnemyHits >= TotalPlayerShipCells
	m.EnemyShots, m.EnemyHits, gameOver =
		s.attack.AITurn(
			aiPlayer,
			m.PlayerEntityBoard,
			m.PlayerBoard,
			m.EnemyShots,
			m.EnemyHits,
			m.TotalPlayerShipCells,
		)

	hit = m.EnemyHits > prevHits

	if hit {
		m.EnemyHitStreak++
		if m.EnemyHitStreak > m.EnemyMaxHitStreak {
			m.EnemyMaxHitStreak = m.EnemyHitStreak
		}
	} else {
		m.EnemyHitStreak = 0
	}

	return hit, gameOver
}

func (s *MatchService) postEnemyStep(m *entity.Match, now time.Time, hit, gameOver bool, ev *entity.AttackEvent) error {
	if gameOver {
		s.finishAndFillWinner(m, now, entity.TurnEnemy, ev)
		return nil
	}

	if hit {
		m.NextAction = entity.NextActionEnemyAttack
		m.NextActionAt = now.Add(s.aiDelay)
	} else {
		m.Turn = entity.TurnPlayer
		m.ClearNextAction()
	}
	return nil
}

func (s *MatchService) makeEvent(attacker entity.TurnOwner, row, col int, valid, hit bool) entity.AttackEvent {
	return entity.AttackEvent{
		Attacker: attacker,
		Row:      row,
		Col:      col,
		Valid:    valid,
		Hit:      hit,
		GameOver: false,
		Winner:   "",
	}
}

func (s *MatchService) finishAndFillWinner(m *entity.Match, now time.Time, winner entity.TurnOwner, ev *entity.AttackEvent) {
	m.Finish(now, winner)
	ev.GameOver = true
	ev.Winner = winner
}

// ResultForPlayer converte o estado final do Match em MatchResult do ponto de vista do player.
func (s *MatchService) ResultForPlayer(m *entity.Match) entity.MatchResult {
	return m.Result()
}
