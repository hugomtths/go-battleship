package service

import (
	"fmt"
	"github.com/allanjose001/go-battleship/internal/entity"
	"math/rand"
	"time"
)

func PositionShipsRandomly(b *entity.Board, f *entity.Fleet) {
	rand.Seed(time.Now().UnixNano())

restart:
	// limpa tabuleiro
	for i := 0; i < entity.BoardSize; i++ {
		for j := 0; j < entity.BoardSize; j++ {
			b.Positions[i][j] = entity.Position{}
		}
	}

	for _, ship := range f.Ships {
		if ship == nil {
			continue
		}

		placed := false
		attempts := 0
		for !placed && attempts < 1000 {
			attempts++
			ship.Horizontal = rand.Intn(2) == 0

			var row, col int
			if ship.Horizontal {
				row = rand.Intn(entity.BoardSize)
				col = rand.Intn(entity.BoardSize - ship.Size + 1)
			} else {
				row = rand.Intn(entity.BoardSize - ship.Size + 1)
				col = rand.Intn(entity.BoardSize)
			}

			if b.PlaceShip(ship, row, col) {
				placed = true
				fmt.Printf("Navio %s posicionado em %v,%v (Horizontal: %v)\n", ship.Name, row, col, ship.Horizontal)
			}
		}

		// se falhar em posicionar um navio depois de muitas tentativas, reinicia todo o processo
		if !placed {
			goto restart
		}
	}
}
