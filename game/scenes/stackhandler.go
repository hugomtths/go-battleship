package scenes

import (
	"github.com/allanjose001/go-battleship/game/state"
)

// StackHandler possui a stack e o context e implementa awares
type StackHandler struct {
	stack *SceneStack
	ctx   *state.GameContext
}

// SetStack serve para "passar a stack adiante" sem precisar explicitamente passar de scene em scene
func (b *StackHandler) SetStack(s *SceneStack) {
	b.stack = s
}

func (b *StackHandler) SetContext(ctx *state.GameContext) {
	b.ctx = ctx
}
