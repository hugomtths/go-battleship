package entity

import (
	"errors"
	"time"

	"github.com/allanjose001/go-battleship/game/shared/board"
	"github.com/allanjose001/go-battleship/game/shared/placement"
)

// Regras / estados

type MatchStatus string

const (
	MatchStatusWaiting    MatchStatus = "waiting"
	MatchStatusInProgress MatchStatus = "in_progress"
	MatchStatusFinished   MatchStatus = "finished"
)

type TurnOwner string

const (
	TurnPlayer TurnOwner = "player"
	TurnEnemy  TurnOwner = "enemy"
)

type NextAction string

const (
	NextActionNone        NextAction = ""
	NextActionEnemyAttack NextAction = "enemy_attack"
)

var (
	// Jogada inválida (ex.: célula já atacada).
	ErrInvalidAttackCell = errors.New("invalid attack cell")
	// Partida não iniciada ou configuração inválida.
	ErrMatchNotReady = errors.New("match not ready")
)

// AttackEvent é um evento “do match” para o front consumir (opcional, mas útil).
// Ele não carrega referências complexas; é puro dado.
type AttackEvent struct {
	Attacker TurnOwner `json:"attacker"`
	Row      int       `json:"row"`
	Col      int       `json:"col"`
	Valid    bool      `json:"valid"`
	Hit      bool      `json:"hit"`
	GameOver bool      `json:"game_over"`
	Winner   TurnOwner `json:"winner"` // preenchido só se GameOver=true
}

// Match é a partida.
// Observação: Boards e IA NÃO são serializáveis e ficam com json:"-".
type Match struct {
	ID            string      `json:"id"`
	Status        MatchStatus `json:"status"`
	Difficulty    string      `json:"difficulty"`
	IsDynamicMode bool        `json:"is_dynamic_mode"`

	Turn   TurnOwner `json:"turn"`
	Winner TurnOwner `json:"winner"` // "" enquanto não terminou

	// Agendamento do turno da IA (equivalente ao QTimer::singleShot)
	NextAction   NextAction `json:"next_action"`
	NextActionAt time.Time  `json:"next_action_at"`

	// Timing da partida (para MatchResult.Duration)
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`

	// Estatísticas básicas (para MatchResult)
	PlayerShots int `json:"player_shots"`
	PlayerHits  int `json:"player_hits"`

	EnemyShots int `json:"enemy_shots"`
	EnemyHits  int `json:"enemy_hits"`

	// Sequência de acertos (para HigherHitSequence)
	PlayerHitStreak    int `json:"player_hit_streak"`
	PlayerMaxHitStreak int `json:"player_max_hit_streak"`

	EnemyHitStreak    int `json:"enemy_hit_streak"`
	EnemyMaxHitStreak int `json:"enemy_max_hit_streak"`

	// Total de células ocupadas por navios do inimigo (usado no seu AttackService.PlayerAttack)
	TotalEnemyShipCells  int `json:"total_enemy_ship_cells"`
	TotalPlayerShipCells int `json:"total_player_ship_cells"`

	LastAttackAt time.Time `json:"-"`     // momento do último ataque do player
	Score        int       `json:"score"` // score atual atualizado a cada tiro

	// Estado runtime (não persistir)
	PlayerBoard *board.Board               `json:"-"`
	EnemyBoard  *board.Board               `json:"-"`
	PlayerShips []*placement.ShipPlacement `json:"-"`
	Profile     *Profile                   `json:"-"`

	// Visão lógica do jogador para a IA (entity.Board é o que seu AIPlayer ataca)
	PlayerEntityBoard *Board `json:"-"`
	PlayerFleet       *Fleet `json:"-"`
}

func NewMatch(id string, difficulty string, playerBoard, aiBoard *board.Board, ships []*placement.ShipPlacement, profile *Profile) *Match {
	return &Match{
		ID:          id,
		Difficulty:  difficulty,
		Status:      MatchStatusWaiting,
		Turn:        TurnPlayer,
		Winner:      "",
		PlayerBoard: playerBoard,
		EnemyBoard:  aiBoard,
		PlayerShips: ships,
		Profile:     profile,
	}
}

func (m *Match) IsFinished() bool {
	return m.Status == MatchStatusFinished
}

func (m *Match) ClearNextAction() {
	m.NextAction = NextActionNone
	m.NextActionAt = time.Time{}
}

func (m *Match) Start(now time.Time) {
	m.Status = MatchStatusInProgress
	m.Turn = TurnPlayer
	m.Winner = ""
	m.StartedAt = now
	m.FinishedAt = time.Time{}
	m.ClearNextAction()

	// reseta stats
	m.PlayerShots = 0
	m.PlayerHits = 0
	m.EnemyShots = 0
	m.EnemyHits = 0
	m.PlayerHitStreak = 0
	m.PlayerMaxHitStreak = 0
	m.EnemyHitStreak = 0
	m.EnemyMaxHitStreak = 0
}

func (m *Match) Finish(now time.Time, winner TurnOwner) {
	m.Status = MatchStatusFinished
	m.Winner = winner
	m.FinishedAt = now
	m.ClearNextAction()
}

// Result gera MatchResult a partir do estado do Match.
// Score/killedShips podem ser ajustados depois (depende do teu design de pontuação).
func (m *Match) Result() MatchResult {
	var dur int64
	if !m.StartedAt.IsZero() {
		end := m.FinishedAt
		if end.IsZero() {
			end = time.Now()
		}
		dur = end.Sub(m.StartedAt).Milliseconds()
	}

	win := m.Winner == TurnPlayer

	// Nota: killedShips com precisão exige rastrear frota lógica da IA.
	// Como hoje seu inimigo está só no board visual (setup.RandomlyPlaceAIShips),
	// eu deixo killedShips = 0 por enquanto.
	killedShips := 0

	// LostShips dá para obter da Fleet lógica do player (a IA mantém hitcount).
	lostShips := 0
	if m.PlayerFleet != nil {
		for _, sh := range m.PlayerFleet.GetFleetShips() {
			if sh != nil && sh.IsDestroyed() {
				lostShips++
			}
		}
	}

	score := m.Score

	return MatchResult{
		Win:               win,
		PlayerShots:       m.PlayerShots,
		Hits:              m.PlayerHits,
		HigherHitSequence: m.PlayerMaxHitStreak,
		Score:             score,
		LostShips:         lostShips,
		KilledShips:       killedShips,
		Duration:          dur,
	}
}

func (m *Match) UpdateScore(hit bool, now time.Time) {
	if !hit {
		return
	}

	const (
		basePoints    = 10
		hitBonus      = 5
		streakFactor  = 10
		maxDeltaSec   = 5   // limite de 5 segundos entre tiros
		maxDuration   = 300 // 5 minutos em segundos
		minMultiplier = 0.1 // mínimo de pontos
	)

	// --- tempo desde o último ataque ---
	deltaSec := 1.0
	if !m.LastAttackAt.IsZero() {
		deltaSec = now.Sub(m.LastAttackAt).Seconds()
		if deltaSec <= 0 {
			deltaSec = 0.1
		}
	}
	if deltaSec > maxDeltaSec {
		deltaSec = maxDeltaSec
	}

	// multiplicador entre tiros (0s = máximo, 5s = mínimo)
	deltaMultiplier := 1 + (maxDeltaSec-deltaSec)/maxDeltaSec

	// --- tempo total da partida ---
	elapsedSec := 1.0
	if !m.StartedAt.IsZero() {
		elapsedSec = now.Sub(m.StartedAt).Seconds()
		if elapsedSec <= 0 {
			elapsedSec = 1
		}
	}
	if elapsedSec > maxDuration {
		elapsedSec = maxDuration
	}

	// multiplicador de duração, linear até 5 minutos, mínimo fixo
	durationMultiplier := 1 - (elapsedSec/maxDuration)*(1-minMultiplier)
	if durationMultiplier < minMultiplier {
		durationMultiplier = minMultiplier
	}

	// multiplicador de sequência
	streakMultiplier := 1 + float64(m.PlayerHitStreak)*float64(streakFactor)/100.0

	// cálculo final
	points := float64(basePoints+hitBonus) * streakMultiplier * deltaMultiplier * durationMultiplier
	if points < 1 {
		points = 1
	}

	m.Score += int(points)
	m.LastAttackAt = now
}
