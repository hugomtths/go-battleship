package inputhelper

import (
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// IsHovered verifica se o mouse está sobre a área
// recebe posição do mouse, posição do alvo e tamanho do alvo
func IsHovered(x, y int, pos basic.Point, size basic.Size) bool {
	mx, my := float32(x), float32(y)
	return mx >= pos.X &&
		mx <= pos.X+size.W &&
		my >= pos.Y &&
		my <= pos.Y+size.H
}

func IsClicked(x, y int, pos basic.Point, size basic.Size) bool {
	return IsHovered(x, y, pos, size) && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
}

func IsPressed(x, y int, pos basic.Point, size basic.Size) bool {
	return IsHovered(x, y, pos, size) && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
}

// UpdateText Atualiza o texto com os caracteres digitados e backspace
// Atualiza o texto com os caracteres digitados e backspace
// active = teclado ativo
func UpdateText(text *string, active bool) {
	if !active {
		return
	}

	// essa func retorna []rune, então deve ser convertida pra string
	runes := ebiten.AppendInputChars([]rune(*text))
	*text = string(runes)

	// backspace
	if ebiten.IsKeyPressed(ebiten.KeyBackspace) && len(*text) > 0 {
		*text = (*text)[:len(*text)-1]
	}
}
