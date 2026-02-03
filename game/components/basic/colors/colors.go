package colors

import "image/color"

// cores devem ser adicionadas aqui
var (
	White       = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	Black       = color.RGBA{A: 255}
	Red         = color.RGBA{R: 255, A: 255}
	Green       = color.RGBA{G: 255, A: 255}
	Blue        = color.RGBA{G: 200, B: 255, A: 1}
	Dark        = color.RGBA{R: 2, G: 94, B: 120, A: 1}
	Transparent = color.RGBA{}
)

// Lighten função que clareia cor (usado em hover e click em botão)
func Lighten(c color.Color, t float64) color.Color {
	r, g, b, a := c.RGBA()

	lerp := func(v uint32) uint8 {
		f := float64(v >> 8) // converte 16-bit → 8-bit
		return uint8(f + (255-f)*t)
	}

	return color.RGBA{
		R: lerp(r),
		G: lerp(g),
		B: lerp(b),
		A: uint8(a >> 8),
	}
}
