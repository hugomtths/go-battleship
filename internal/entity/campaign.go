package entity

// Defina a struct aqui em cima ou no mesmo arquivo
type Campaign struct {
	ID             string                 `json:"id"`
	DifficultyStep map[string]MatchResult `json:"difficulty_step"`
	IsActive       bool                   `json:"is_active"`
}
