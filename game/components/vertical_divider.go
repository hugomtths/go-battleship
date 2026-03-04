package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// VerticalDivider é um componente visual simples usado para desenhar uma linha vertical na tela.
// Geralmente utilizado para separar áreas da interface, como os dois tabuleiros de batalha.
type VerticalDivider struct {
	// X define a posição horizontal da linha.
	X float64
	// Y1 define o ponto inicial vertical da linha (topo).
	Y1 float64
	// Y2 define o ponto final vertical da linha (base).
	Y2 float64
	// Color define a cor da linha a ser desenhada.
	Color color.Color
}

// NewVerticalDivider cria uma nova instância de um divisor vertical.
// Recebe a posição X, os limites verticais (Y1, Y2) e a cor desejada.
func NewVerticalDivider(x, y1, y2 float64, c color.Color) *VerticalDivider {
	return &VerticalDivider{
		X:     x,
		Y1:    y1,
		Y2:    y2,
		Color: c,
	}
}

// Draw renderiza a linha vertical na tela usando as coordenadas e cor configuradas.
// Utiliza a função utilitária do Ebiten para desenhar linhas.
func (v *VerticalDivider) Draw(screen *ebiten.Image) {
	ebitenutil.DrawLine(screen, v.X, v.Y1, v.X, v.Y2, v.Color)
}
