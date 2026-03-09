package state

import (
	"github.com/allanjose001/go-battleship/game/scenes/audio"
	"github.com/allanjose001/go-battleship/internal/entity"
)

// GameContext possui dados de interesse do jogo (tela de jogo, perfis, etc)
type GameContext struct {
	Profile              *entity.Profile
	Match                *entity.Match
	BattleService        BattleService
	DynamicBattleService DynamicBattleService
	SoundService         *audio.SoundService
	Difficulty           string
	IsCampaign           bool
	IsDynamicMode        bool
	CanPopOrPush         bool
}

type ContextAware interface {
	SetContext(*GameContext)
}

// deve ser inicializado agora
func NewGameContext() *GameContext {
	ss := audio.NewSoundService()

	ss.LoadMusic("menus", "assets/audio/music/menus.ogg")
	ss.LoadMusic("loss", "assets/audio/music/loss.ogg")
	ss.LoadMusic("battle", "assets/audio/music/battle-scene.ogg")
	ss.LoadMusic("victory", "assets/audio/music/victory.ogg")
	ss.LoadSFX("attack", "assets/audio/sfx/attack2.ogg")
	ss.LoadSFX("watersplash", "assets/audio/sfx/waterSplash1.ogg")
	ss.LoadSFX("fah", "assets/audio/sfx/fah.ogg")
	ss.LoadSFX("click", "assets/audio/sfx/click.ogg")
	ss.LoadSFX("backclick", "assets/audio/sfx/backclick.ogg")

	return &GameContext{
		SoundService: ss,
		CanPopOrPush: true,
	}
}

// BattleService define a interface para interação com a lógica de batalha.
// Essa interface é duplicada aqui para evitar ciclos de importação com internal/service.
type BattleService interface {
	HandlePlayerClick(row, col int) (*entity.MatchResult, error)
	HandleEnemyTurn() (*entity.MatchResult, error)
	Stats() (playerShots, playerHits, enemyShots, enemyHits int, isPlayerTurn bool)
	WinnerName() string
}

type DynamicBattleService interface {
	HandlePlayerMove(row, col int, dir entity.Direction) error
	HandleEnemyMove() error
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

func (c *GameContext) SetDynamicBattleService(s DynamicBattleService) {
	c.DynamicBattleService = s
}

func (c *GameContext) SetDifficulty(d string) {
	c.Difficulty = d
}
