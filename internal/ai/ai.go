package ai

import ( "github.com/allanjose001/go-battleship/internal/entity" )


type AIPlayer struct { 
	virtualBoard [10][10]int
	priorityQueue []Pair
	Strategies []Strategy
	enemyFleet *entity.Fleet
	chaseMode bool
}

func (ai *AIPlayer) Attack(enemyBoard *entity.Board) {
	for _, strat := range ai.Strategies { // verifica estrategias disponiveis
		if strat.TryAttack(ai, enemyBoard) {
			return;
		}
	}
}

func (ai *AIPlayer) AdjustStrategy(board *entity.Board, row, y int, ship *entity.Ship) {
	if ship == nil {
		ai.virtualBoard[row][y] = 1
		return
	}

	if ship.IsDestroyed() {
		ai.virtualBoard[row][y] = 3
		ai.WreckedShipAdjustment(board, row, y)
		ai.ClearPriorityQueue()
		ai.FleetShipDestroyed(ship.Size)
		ai.StopChase()
	} else {
		ai.virtualBoard[row][y] = 2
		ai.AttackNeighbors(row, y)
	}
}

// retorna o tamanho do proximo navio da enemyFleet
func (ai *AIPlayer) SizeOfNextShip() int {
	for _, ship := range ai.enemyFleet.Ships {
		if ship != nil && !ship.IsDestroyed() {
			return ship.Size
		}
	}
	return 0
}

// marca o navio como destruído na fleet interna do AI
func (ai *AIPlayer) FleetShipDestroyed(size int) {
	for _, ship := range ai.enemyFleet.Ships {
		if ship != nil && ship.Size == size && !ship.IsDestroyed() {
			ship.HitCount = ship.Size
			return
		}
	}
}

// Retorna a posição inicial do navio (mais à esquerda se horizontal, mais acima se vertical)
func (ai *AIPlayer) LocateShipStart(board *entity.Board, row, col int) (startRow, startCol int) {
	ship := entity.GetShipReference(board.Positions[row][col])
	if ship == nil {
		return row, col
	}
	startRow, startCol = row, col

	if ship.IsHorizontal() {
		for startCol > 0 && entity.GetShipReference(board.Positions[row][startCol-1]) == ship {
			startCol--
		}
	} else {
		for startRow > 0 && entity.GetShipReference(board.Positions[startRow-1][col]) == ship {
			startRow--
		}
	}

	return
}

// Procura verticalmente por uma sequência de posições vazias com tamanho suficiente para o próximo navio
func (ai *AIPlayer) SearchVertically(size int) bool {
	for j := 0; j < 10; j++ {
		contiguous := 0
		for i := 0; i < 10; i++ {
			if ai.virtualBoard[i][j] == 0 {
				contiguous++
			} else {
				contiguous = 0
			}
			if contiguous >= size {
				targetRow := i - size/2
				if targetRow < 0 {
					targetRow = 0
				} else if targetRow > 9 {
					targetRow = 9
				}
				ai.AddToPriorityQueue(targetRow, j)
				return true
			}
		}
	}
	return false
}

// Procura horizontalmente por uma sequência de posições vazias com tamanho suficiente para o próximo navio
func (ai *AIPlayer) SearchHorizontally(size int) bool {
	for i := 0; i < 10; i++ {
		contiguous := 0
		for j := 0; j < 10; j++ {
			if ai.virtualBoard[i][j] == 0 {
				contiguous++
			} else {
				contiguous = 0
			}
			if contiguous >= size {
				targetCol := j - size/2
				if targetCol < 0 {
					targetCol = 0
				} else if targetCol > 9 {
					targetCol = 9
				}
				ai.AddToPriorityQueue(i, targetCol)
				return true
			}
		}
	}
	return false
}

func (ai *AIPlayer) WreckedShipAdjustment(board *entity.Board, row, col int) {
	ship := entity.GetShipReference(board.Positions[row][col])
	if ship == nil {
		return
	}

	startRow, endRow := row, row
	startCol, endCol := col, col

	// vertical
	for startRow > 0 && entity.GetShipReference(board.Positions[startRow-1][col]) == ship {
		startRow--
	}
	for endRow < 9 && entity.GetShipReference(board.Positions[endRow+1][col]) == ship {
		endRow++
	}

	// horizontal
	for startCol > 0 && entity.GetShipReference(board.Positions[row][startCol-1]) == ship {
		startCol--
	}
	for endCol < 9 && entity.GetShipReference(board.Positions[row][endCol+1]) == ship {
		endCol++
	}

	adjStartRow := max(0, startRow-1)
	adjEndRow := min(9, endRow+1)
	adjStartCol := max(0, startCol-1)
	adjEndCol := min(9, endCol+1)

	for i := adjStartRow; i <= adjEndRow; i++ {
		for j := adjStartCol; j <= adjEndCol; j++ {
			ai.virtualBoard[i][j] = 3
		}
	}
}


// Adiciona uma posição válida à fila de prioridade
func (ai *AIPlayer) AddToPriorityQueue(row, col int) {
	if ai.IsValid(row, col) {
		ai.priorityQueue = append(ai.priorityQueue, Pair{row, col})
	}
}

// Limpa a fila de prioridade
func (ai *AIPlayer) ClearPriorityQueue() {
	ai.priorityQueue = []Pair{}
}

// Adiciona posições vizinhas à fila de prioridade
func (ai *AIPlayer) AttackNeighbors(row, col int) {
	ai.AddToPriorityQueue(row-1, col)
	ai.AddToPriorityQueue(row+1, col)
	ai.AddToPriorityQueue(row, col-1)
	ai.AddToPriorityQueue(row, col+1)
}

// Verifica se a posição é válida para atacar (dentro do tabuleiro e ainda não marcada)
func (ai *AIPlayer) IsValid(row, col int) bool {
	if row < 0 || row >= 10 || col < 0 || col >= 10 {
		return false
	}
	return ai.virtualBoard[row][col] == 0
}

// Apenas checa se a posição está dentro do tabuleiro
func (ai *AIPlayer) IsValidForTesting(row, col int) bool {
	return row >= 0 && row < 10 && col >= 0 && col < 10
}

func (ai *AIPlayer) PopPriority() (row, y int) {
	if len(ai.priorityQueue) == 0 {
		return -1, -1
	}

	p := ai.priorityQueue[0]
	ai.priorityQueue = ai.priorityQueue[1:]
	return p.row, p.col
}

func (ai *AIPlayer) StartChase() {
	ai.chaseMode = true
}

func (ai *AIPlayer) StopChase() {
	ai.chaseMode = false
}

func (ai *AIPlayer) IsChasing() bool {
	return ai.chaseMode
}

// Ataca posições estratégicas quando 30% do tabuleiro estiver preenchido
func (ai *AIPlayer) ShouldAttackStrategicPositions() bool {
	filled := 0

	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if ai.virtualBoard[i][j] != 0 {
				filled++
			}
		}
	}
	return float64(filled)/100.0 >= 0.3 
}