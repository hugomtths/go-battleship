//Package inputhelper funções utilitárias para lidar com entrada
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

//IsClicked verifica clique na posição (field pos) e tamanho (field size) especificado
func IsClicked(x, y int, pos basic.Point, size basic.Size) bool {
	return IsHovered(x, y, pos, size) && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
}

//IsPressed verifica se o botão está sendo pressionado na área
func IsPressed(x, y int, pos basic.Point, size basic.Size) bool {
	return IsHovered(x, y, pos, size) && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
}

// ReceiveText Atualiza o texto com os caracteres digitados e trata backspace
func ReceiveText(text *string, active bool) {
	if !active {
		return
	}

	// essa func retorna []rune, então deve ser convertida pra string
	runes := ebiten.AppendInputChars([]rune(*text))
	*text = string(runes)

	// backspace
    if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && len(*text) > 0 {
        *text = (*text)[:len(*text)-1]
    }
}
