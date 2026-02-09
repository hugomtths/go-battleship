package entity

import "fmt"

const BoardSize = 10

type Board struct {
	Positions [BoardSize][BoardSize]Position
}

// variação A que retorna boolean
func (b *Board) AttackPositionA(x int, y int) bool {
	fmt.Printf("atacando %v,%v\n", x, y)
	if b.CheckPosition(x, y) {
		attack(&b.Positions[x][y])

		return true
	}

	return false
}

// variação B que retorna o navio atacado (ou nil se não houver navio)
func (b *Board) AttackPositionB(x int, y int) *Ship {
	fmt.Printf("atacando %v,%v\n", x, y)
	if b.CheckPosition(x, y) {
		attack(&b.Positions[x][y])

		return GetShipReference(b.Positions[x][y])
	}

	return nil
}

func (b *Board) PlaceShip(ship *Ship, x int, y int) bool {
	if !b.CheckShipPosition(ship, x, y) {
		return false
	}

	if ship.IsHorizontal() {
		for i := y; i < y+ship.Size; i++ {
			PlaceShip(&b.Positions[x][i], ship)
		}
	} else {
		for i := x; i < x+ship.Size; i++ {
			PlaceShip(&b.Positions[i][y], ship)
		}
	}

	return true

}

func (b *Board) RemoveShipFromBoard(ship *Ship) {
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			var currentShip *Ship = GetShipReference(b.Positions[i][j])

			if currentShip == ship {
				RemoveShip(&b.Positions[i][j])

				Unblock(&b.Positions[i][j])

			}
		}
	}
}

func (b *Board) CheckShipPosition(ship *Ship, x int, y int) bool {
	if ship.IsHorizontal() { //se o navio estiver na horizontal:
		if y+ship.Size > 10 { // verifica se o navio ultrapassa os limites do tabuleiro
			return false
		}

		for i := y; i < y+ship.Size; i++ { //se a posição não está bloqueada
			if !IsValidPosition(b.Positions[x][i]) {
				return false
			}
		}
	} else { // se o navio estiver na vertical:
		if ship.Size+x > 10 {
			return false
		}

		for i := x; i < x+ship.Size; i++ {
			if !IsValidPosition(b.Positions[i][y]) {
				return false
			}
		}
	}
	// se todas as verificações passarem, a posição é válida
	return true
}

func (b *Board) CheckPosition(x int, y int) bool {
	if x < 0 || x > 9 || y < 0 || y > 9 {
		return false
	}

	return !(IsAttacked(b.Positions[x][y]))
}

func PrintBoard(b *Board) {
	for i := 0; i < 10; i++ { // itera pelas linhas
		for j := 0; j < 10; j++ { // itera pelas colunas
			if IsAttacked(b.Positions[i][j]) { // se a posição foi atacada
				if GetShipReference(b.Positions[i][j]) != nil {
					print("x ") // posição atacada com navio
					continue
				}

				print("o ") // posição atacada sem navio
				continue
			} else if GetShipReference(b.Positions[i][j]) != nil {
				print("B ") // marca como bloqueada.
				continue
			}

			//posição valida e não atacada.
			print("- ")
		}
		print("\n") // nova linha apos cada linha do tabuleiro

	}
}
