package medal

import "github.com/allanjose001/go-battleship/internal/entity"

// Medal struct medalha [precisei adicionar nesse package para driblar cyclic import]
type Medal struct {
	Name         string                              `json:"name"`
	Description  string                              `json:"description"`
	IconPath     string                              `json:"icon"`
	Verification func(stats entity.PlayerStats) bool `json:"-"` //cada medalha tem seus criterios de verificação
}
