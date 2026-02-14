package entity

import "fmt"

const BoardSize = 10

type Board struct {
	Positions [BoardSize][BoardSize]Position
}

// variação A que retorna boolean
func (b *Board) AttackPositionA(row int, col int) bool {
	fmt.Printf("atacando %v,%v\n", row, col)
	if b.CheckPosition(row, col) {
		attack(&b.Positions[row][col])

		return true
	}

	return false
}

// variação B que retorna o navio atacado (ou nil se não houver navio)
func (b *Board) AttackPositionB(row int, col int) *Ship {
	fmt.Printf("atacando %v,%v\n", row, col)
	if b.CheckPosition(row, col) {
		attack(&b.Positions[row][col])

		return GetShipReference(b.Positions[row][col])
	}

	return nil
}

func (b *Board) PlaceShip(ship *Ship, row int, col int) bool {
	if !b.CheckShipPosition(ship, row, col) {
		return false
	}

	if ship.IsHorizontal() {
		for i := col; i < col+ship.Size; i++ {
			PlaceShip(&b.Positions[row][i], ship)
		}
	} else {
		for i := row; i < row+ship.Size; i++ {
			PlaceShip(&b.Positions[i][col], ship)
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

func (b *Board) CheckShipPosition(ship *Ship, row int, col int) bool {
	if ship.IsHorizontal() { //se o navio estiver na horizontal:
		if col+ship.Size > 10 { // verifica se o navio ultrapassa os limites do tabuleiro
			return false
		}

		for i := col; i < col+ship.Size; i++ { //se a posição não está bloqueada
			if !IsValidPosition(b.Positions[row][i]) {
				return false
			}
		}
	} else { // se o navio estiver na vertical:
		if ship.Size+row > 10 {
			return false
		}

		for i := row; i < row+ship.Size; i++ {
			if !IsValidPosition(b.Positions[i][col]) {
				return false
			}
		}
	}
	// se todas as verificações passarem, a posição é válida
	return true
}

func (b *Board) CheckPosition(row int, col int) bool {
	if row < 0 || row > 9 || col < 0 || col > 9 {
		return false
	}

	return !(IsAttacked(b.Positions[row][col]))
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
