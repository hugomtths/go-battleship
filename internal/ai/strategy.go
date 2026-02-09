package ai

import "github.com/allanjose001/go-battleship/internal/entity"

type Strategy interface {
	TryAttack(ai *AIPlayer, board *entity.Board) bool 
}