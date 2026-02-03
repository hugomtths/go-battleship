package components

import (
	"image/color"

	"github.com/allanjose001/go-battleship/game/components/basic"
)

type StylableWidget interface {
	//herda methods de widget
	Widget
	//SetColor serve para mudar a cor quando necessário
	SetColor(color.Color)

	//muda tamanho
	SetSize(size basic.Size)
}
