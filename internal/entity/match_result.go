package entity

import "fmt"

// MatchResult struct que encapsula resultado da partida para histórico e estatisticas do jogo
type MatchResult struct {
	Win               bool  `json:"win"`
	PlayerShots       int   `json:"player_shots"`
	Hits              int   `json:"hits"`
	HigherHitSequence int   `json:"higher_hit_sequence"`
	Score             int   `json:"score"`
	LostShips         int   `json:"lost_ships"`
	KilledShips       int   `json:"killed_ships"`
	Duration          int64 `json:"duration"` //-> em milissegundos
	Mode			  string `json:"mode"`
}

// FormattedDuration retorna string para ser usada no front
func (m MatchResult) FormattedDuration() string {
	totalSec := m.Duration / 1000
	minute := totalSec / 60
	sec := totalSec % 60

	return fmt.Sprintf("%02dm:%02ds", minute, sec)
}
