package bootstrap

import (
	"math/rand"
	"time"
)

// InitRandom inicializa a seed global de rand
// para toda a aplicação, evitando chamadas repetidas
// em cada service individual.
func InitRandom() {
	rand.Seed(time.Now().UnixNano())
}
