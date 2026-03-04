package state

import "github.com/allanjose001/go-battleship/internal/entity"

// GameContext possui dados de interesse do jogo (tela de jogo, perfis, etc)
type GameContext struct {
	Profile       *entity.Profile
	Match         *entity.Match
	BattleService BattleService
}

type ContextAware interface {
	SetContext(*GameContext)
}

// BattleService define a interface para interação com a lógica de batalha.
// Essa interface é duplicada aqui para evitar ciclos de importação com internal/service.
type BattleService interface {
	HandlePlayerClick(row, col int) (*entity.MatchResult, error)
	HandleEnemyTurn() (*entity.MatchResult, error)
	Stats() (playerShots, playerHits, enemyShots, enemyHits int, isPlayerTurn bool)
	WinnerName() string
}

func NewGameContext() *GameContext {
	return &GameContext{}
}

func (c *GameContext) SetProfile(p *entity.Profile) {
	c.Profile = p
}

func (c *GameContext) SetMatch(m *entity.Match) {
	c.Match = m
}

func (c *GameContext) SetBattleService(s BattleService) {
	c.BattleService = s
}
