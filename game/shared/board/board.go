package board

import "github.com/hajimehoshi/ebiten/v2"

//estrutura do tabuleiro

const (
	Rows = 10
	Cols = 10
)

type Board struct {
	Cells           [][]Cell
	X               float64 // posição na tela
	Y               float64
	Size            float64 // tamanho total
	BackgroundImage *ebiten.Image
}

func NewBoard(x, y, size float64) *Board {
	cells := make([][]Cell, Rows)
	for i := 0; i < Rows; i++ {
		cells[i] = make([]Cell, Cols)
		for j := 0; j < Cols; j++ {
			cells[i][j] = Cell{
				Row:   i,
				Col:   j,
				State: Empty,
			}
		}
	}

	return &Board{
		Cells: cells,
		X:     x,
		Y:     y,
		Size:  size,
	}
}

type Orientation int

const (
	Horizontal Orientation = iota
	Vertical
)

func (b *Board) CanPlace(size, row, col int, orientation Orientation) bool {
	if row < 0 || row >= Rows || col < 0 || col >= Cols {
		return false
	}

	if orientation == Horizontal {
		if col+size > Cols {
			return false
		}
		for j := 0; j < size; j++ {
			if b.Cells[row][col+j].State != Empty {
				return false
			}
		}
	} else {
		if row+size > Rows {
			return false
		}
		for i := 0; i < size; i++ {
			if b.Cells[row+i][col].State != Empty {
				return false
			}
		}
	}
	return true
}

func (b *Board) PlaceShip(size, row, col int, orientation Orientation) {
	if orientation == Horizontal {
		for j := 0; j < size; j++ {
			b.Cells[row][col+j].State = Ship
		}
	} else {
		for i := 0; i < size; i++ {
			b.Cells[row+i][col].State = Ship
		}
	}
}

func (b *Board) Clear() {
	for i := 0; i < Rows; i++ {
		for j := 0; j < Cols; j++ {
			b.Cells[i][j].State = Empty
		}
	}
}
