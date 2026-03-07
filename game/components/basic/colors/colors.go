package colors

import "image/color"

// cores devem ser adicionadas aqui
var (
	White       = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	Black       = color.RGBA{A: 255}
	Red         = color.RGBA{R: 255, A: 255}
	Green       = color.RGBA{R: 0, G: 255, B: 0, A: 255}
	DarkBlue    = color.RGBA{R: 13, G: 89, B: 99, A: 255}
	Blue        = color.RGBA{R: 43, G: 112, B: 121, A: 255}
	NightBlue   = color.RGBA{R: 40, G: 40, B: 50, A: 255}
	PlayerInput = color.RGBA{R: 65, G: 81, B: 100, A: 255} // #415164
	Dark        = color.RGBA{R: 48, G: 67, B: 103, A: 255}
	Transparent = color.RGBA{}
	Background  = color.RGBA{R: 13, G: 27, B: 42, A: 255}

	// --- NOVAS CORES PARA O MENU DE DIFICULDADES ---

	// NavyBlue O azul escuro profundo para o fundo do menu (#0D5963 aproximado)
	NavyBlue = color.RGBA{R: 13, G: 89, B: 99, A: 255}

	// SeaCyan O ciano vibrante para os botões e medalhas (#2B7079 aproximado)
	SeaCyan = color.RGBA{R: 43, G: 112, B: 121, A: 255}

	// DeepWater Um tom intermediário para contrastes navais
	DeepWater = color.RGBA{R: 10, G: 45, B: 50, A: 255}

	GoldMedal   = color.RGBA{R: 255, G: 215, B: 0, A: 255}
	SilverMedal = color.RGBA{R: 192, G: 192, B: 192, A: 255}
	BronzeMedal = color.RGBA{R: 205, G: 127, B: 50, A: 255}
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

// GrayOut deixa a cor acinzentada (para botões disabled)
func GrayOut(c color.Color, factor float64) color.Color {
	r16, g16, b16, a16 := c.RGBA()

	// converte de 16-bit para 8-bit
	r := float64(uint8(r16 >> 8))
	g := float64(uint8(g16 >> 8))
	b := float64(uint8(b16 >> 8))
	a := uint8(a16 >> 8)

	gray := (r + g + b) / 3

	r = r*(1-factor) + gray*factor
	g = g*(1-factor) + gray*factor
	b = b*(1-factor) + gray*factor

	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: a}
}
