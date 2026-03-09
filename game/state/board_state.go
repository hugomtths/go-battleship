package state

import "github.com/allanjose001/go-battleship/game/shared/board"

type GameState struct {
	PlayerBoard *board.Board
	AIBoard     *board.Board
	PlayerShips interface{} // Usaremos interface{} temporariamente ou criaremos um tipo compartilhado
}

func NewGameState() *GameState {
	return &GameState{
		PlayerBoard: board.NewBoard(80, 150, 320),
		AIBoard:     board.NewBoard(500, 150, 320),
	}
}
