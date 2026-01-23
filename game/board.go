package game

import "fmt"

type Board struct {
	positions [10][10]Position
}

func PlaceShip(b *Board, ship *Ship, x int, y int) bool {
	if !CheckShipPosition(b, ship, x, y) {
		return false
	}

	if isHorizontal(ship) {
		for i := y; i < y+ship.Size; i++ {
			placeShip(&b.positions[x][i], ship)
		}
	} else {
		for i := x; i < x+ship.Size; i++ {
			placeShip(&b.positions[i][y], ship)
		}
	}

	return true

}

func RemoveShipFromBoard(b *Board, ship *Ship) {
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			var currentShip *Ship = getShipReference(b.positions[i][j])

			if currentShip == ship {
				removeShip(&b.positions[i][j])

				unblock(&b.positions[i][j])

			}
		}
	}
}

func AttackPosition(b *Board, x int, y int) bool {
	fmt.Printf("atacando %v,%v\n", x, y)
	if (CheckPosition(b, x, y)) {
		attack(&b.positions[x][y]);

		return true;
	}

	return false;
}

func CheckShipPosition(b *Board, ship *Ship, x int, y int) bool {
	if isHorizontal(ship) { //se o navio estiver na horizontal:
		if y+ship.Size > 10 { // verifica se o navio ultrapassa os limites do tabuleiro
			return false
		}

		for i := y; i < y+ship.Size; i++ { //se a posição não está bloqueada
			if !isValidPosition(b.positions[x][i]) {
				return false
			}
		}
	} else { // se o navio estiver na vertical:
		if ship.Size+x > 10 {
			return false
		}

		for i := x; i < x+ship.Size; i++ {
			if !isValidPosition(b.positions[i][y]) {
				return false
			}
		}
	}
	// se todas as verificações passarem, a posição é válida
	return true
}

func CheckPosition(b *Board, x int, y int) bool {
	if (x < 0 || x > 9 || y < 0 || y > 9) {
		return false;
	}

	return !(isAttacked(b.positions[x][y]));
}

func PrintBoard(b *Board) {
	for i := 0; i < 10; i++ { // itera pelas linhas
		for j := 0; j < 10; j++ { // itera pelas colunas
			if isAttacked(b.positions[i][j]) { // se a posição foi atacada
				if getShipReference(b.positions[i][j]) != nil {
					print("x ") // posição atacada com navio
					continue
				}

				print("o ") // posição atacada sem navio
				continue
			} else if getShipReference(b.positions[i][j]) != nil {
				print("B ") // marca como bloqueada.
				continue
			}

			//posição valida e não atacada.
			print("- ")
		}
		print("\n") // nova linha apos cada linha do tabuleiro

	}
}