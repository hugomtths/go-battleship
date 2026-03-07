package scenes

// StackHandler possui a stack e o context e implementa awares
type StackHandler struct {
	stack *SceneStack
}

// SetStack serve para "passar a stack adiante" sem precisar explicitamente passar de scene em scene
func (b *StackHandler) SetStack(s *SceneStack) {
	b.stack = s
}
