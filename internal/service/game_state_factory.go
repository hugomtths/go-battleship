// GameStateFactory: cria o estado de batalha (GameState) a partir
// do tabuleiro já configurado do jogador e da lista de navios
// posicionados. Não conhece regras de ataque nem IA; foca apenas
// em construir a estrutura de dados usada pela fase de batalha.
package service

import (
	"github.com/allanjose001/go-battleship/game/shared/board"
	"github.com/allanjose001/go-battleship/game/shared/placement"
	"github.com/allanjose001/go-battleship/game/shared/setup"
	"github.com/allanjose001/go-battleship/game/state"
)

type GameService struct{}

func NewGameService() *GameService {
	return &GameService{}
}

// NewBattleGameState:
// - Reaproveita o board do jogador e clona as dimensões para o board da IA
// - Posiciona os navios da IA no tabuleiro dela (visual) via setup
// - Devolve um GameState pronto para a BattleScene consumir
func (g *GameService) NewBattleGameState(playerBoard *board.Board, ships []*placement.ShipPlacement) *state.GameState {
	gs := state.NewGameState()
	gs.PlayerBoard = playerBoard
	gs.PlayerShips = ships

	gs.AIBoard.X = 1280 - playerBoard.X - playerBoard.Size
	gs.AIBoard.Y = playerBoard.Y
	gs.AIBoard.Size = playerBoard.Size

	setup.RandomlyPlaceAIShips(gs.AIBoard)

	return gs
}
