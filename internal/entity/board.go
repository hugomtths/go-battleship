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

func (b *Board) MoveShipSegment(row, col int, dir Direction) (*Ship, error) {
	if !b.CheckPosition(row, col) {
		return nil, fmt.Errorf("invalid position")
	}
	ship := GetShipReference(b.Positions[row][col])
	if ship == nil {
		return nil, fmt.Errorf("no ship at position")
	}
	if IsAttacked(b.Positions[row][col]) {
		return nil, fmt.Errorf("cannot move damaged segment")
	}

	// Find connected component of intact parts
	type point struct{ r, c int }
	var segment []point
	visited := make(map[point]bool)
	queue := []point{{row, col}}
	visited[point{row, col}] = true

	// Check bounds helper
	isValidCoord := func(r, c int) bool {
		return r >= 0 && r < 10 && c >= 0 && c < 10
	}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		segment = append(segment, curr)

		neighbors := []point{
			{curr.r - 1, curr.c}, {curr.r + 1, curr.c},
			{curr.r, curr.c - 1}, {curr.r, curr.c + 1},
		}
		for _, n := range neighbors {
			if isValidCoord(n.r, n.c) {
				if !visited[n] && GetShipReference(b.Positions[n.r][n.c]) == ship && !IsAttacked(b.Positions[n.r][n.c]) {
					visited[n] = true
					queue = append(queue, n)
				}
			}
		}
	}

	// Determine move direction
	dr, dc := 0, 0
	switch dir {
	case Up:
		dr = -1
	case Down:
		dr = 1
	case Left:
		dc = -1
	case Right:
		dc = 1
	}

	// Validate target positions
	for _, p := range segment {
		tr, tc := p.r+dr, p.c+dc
		if !isValidCoord(tr, tc) {
			return nil, fmt.Errorf("out of bounds")
		}

		targetPos := b.Positions[tr][tc]
		// Manual check instead of IsValidPosition because IsValidPosition is too strict about shipReference
		if IsAttacked(targetPos) {
			return nil, fmt.Errorf("position is attacked")
		}
		if IsBlocked(targetPos) {
			return nil, fmt.Errorf("position is blocked")
		}
		targetShip := GetShipReference(targetPos)
		if targetShip != nil {
			if targetShip != ship {
				return nil, fmt.Errorf("collision with another ship")
			}
			// targetShip == ship
			// Check if target is part of the moving segment
			isMovingPart := false
			for _, mp := range segment {
				if mp.r == tr && mp.c == tc {
					isMovingPart = true
					break
				}
			}
			if !isMovingPart {
				return nil, fmt.Errorf("collision with own ship debris")
			}
		}
	}

	// Execute Move
	// 1. Remove ship from old positions
	for _, p := range segment {
		RemoveShip(&b.Positions[p.r][p.c])
		Unblock(&b.Positions[p.r][p.c])
	}

	// 2. Handle Ship Entity
	var finalShip *Ship
	// Count total parts of original ship
	// We can't count on board anymore because we removed some.
	// But ship.Size is reliable.
	if len(segment) == ship.Size {
		// Moving the whole ship
		finalShip = ship
	} else {
		// Splitting
		// Create new ship for the moving part
		finalShip = &Ship{
			Name:       ship.Name,
			Size:       len(segment),
			HitCount:   0, // Moving parts are intact
			Horizontal: ship.Horizontal,
		}
		// Update original ship (remains as debris)
		ship.Size -= len(segment)
	}

	// 3. Place finalShip at new positions
	for _, p := range segment {
		tr, tc := p.r+dr, p.c+dc
		PlaceShip(&b.Positions[tr][tc], finalShip)
	}

	return finalShip, nil
}
