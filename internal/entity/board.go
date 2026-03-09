package entity

import "fmt"


type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

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


func (b *Board) MoveShip(ship *Ship, newRow int, newCol int) error {
    if ship == nil {
        return fmt.Errorf("ship nil")
    }

    // encontra células atuais do navio
    var cells [][2]int
    for r := 0; r < BoardSize; r++ {
        for c := 0; c < BoardSize; c++ {
            if GetShipReference(b.Positions[r][c]) == ship {
                cells = append(cells, [2]int{r, c})
            }
        }
    }
    if len(cells) == 0 {
        return fmt.Errorf("barco não está no tabuleiro")
    }

    // determina coordenada top-left atual do navio (menor row e col)
    minR, minC := cells[0][0], cells[0][1]
    for _, p := range cells {
        if p[0] < minR {
            minR = p[0]
        }
        if p[1] < minC {
            minC = p[1]
        }
    }

    dRow := newRow - minR
    dCol := newCol - minC

    // permite somente movimento de 1 célula ortogonal
    if !((abs(dRow) == 1 && dCol == 0) || (abs(dCol) == 1 && dRow == 0)) {
        return fmt.Errorf("movimento inválido: deve mover exatamente 1 célula ortogonalmente")
    }

    // gera coords alvo baseado na orientação
    var targets [][2]int
    if ship.IsHorizontal() {
        for i := 0; i < ship.Size; i++ {
            r := newRow
            c := newCol + i
            targets = append(targets, [2]int{r, c})
        }
    } else {
        for i := 0; i < ship.Size; i++ {
            r := newRow + i
            c := newCol
            targets = append(targets, [2]int{r, c})
        }
    }

    // valida targets dentro do tabuleiro e não colidindo com terceiros
    for _, p := range targets {
        r, c := p[0], p[1]
        if r < 0 || r >= BoardSize || c < 0 || c >= BoardSize {
            return fmt.Errorf("alvo fora dos limites")
        }
        // pode ser válido se a posição for livre (IsValidPosition) OU se já pertencer ao mesmo navio
        ref := GetShipReference(b.Positions[r][c])
        if ref != nil && ref != ship {
            return fmt.Errorf("alvo colide com outro navio")
        }
        if !IsValidPosition(b.Positions[r][c]) && ref != ship {
            // IsValidPosition verifica attacked/blocked/shipReference==nil
            return fmt.Errorf("posição do alvo não disponível")
        }
    }

    // aplicar movimentação: remover referências antigas e colocar nas novas
    // coleta coords antigas para remoção
    var olds [][2]int
    for _, p := range cells {
        olds = append(olds, p)
    }

    // remove ship das antigas
    for _, p := range olds {
        RemoveShip(&b.Positions[p[0]][p[1]])
        // não altera attacked flag; desbloqueio aqui se desejado:
        Unblock(&b.Positions[p[0]][p[1]])
    }

    // coloca ship nas novas posições
    for _, p := range targets {
        PlaceShip(&b.Positions[p[0]][p[1]], ship)
    }

    return nil
}

func abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}
