package basic

// Point encapsula um ponto x e y no espaço, com operações como soma
// e subtração de pontos para auxiliar em movimentação e centralização no design
type Point struct {
	X float32
	Y float32
}

func NewPoint(x, y float32) Point {
	return Point{X: x, Y: y}
}

// Add serve para mover elementos no valor dado
func (p Point) Add(other Point) Point {
	return Point{
		X: p.X + other.X,
		Y: p.Y + other.Y,
	}
}

// Sub serve para expelir elementos no valor dado
func (p Point) Sub(other Point) Point {
	return Point{
		X: p.X - other.X,
		Y: p.Y - other.Y,
	}
}
