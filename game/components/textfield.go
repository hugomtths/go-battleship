package components

import (
	"image/color"

	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	inputhelper "github.com/allanjose001/go-battleship/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type TextField struct {
	pos, currentPos  basic.Point
	size             basic.Size
	Text             string
	placeholder      string
	focused          bool
	body             StylableWidget
	textWidget       *Text
	textColor        color.Color
	placeholderColor color.Color
	backgroundColor  color.Color

	// 🔹 Cursor
	cursorVisible   bool
	cursorCounter   int
	cursorBlinkRate int
}

func NewTextField(pos basic.Point, size basic.Size, placeholder string) *TextField {
	tf := &TextField{
		pos:              pos,
		size:             size,
		placeholder:      placeholder,
		textColor:        colors.White,
		placeholderColor: color.RGBA{160, 170, 190, 255},
		backgroundColor:  colors.NightBlue,

		// 🔹 Configuração do cursor
		cursorVisible:   true,
		cursorCounter:   0,
		cursorBlinkRate: 30, // 30 frames ≈ 0.5s em 60 FPS
	}

	tf.textWidget = NewText(
		basic.Point{X: 16, Y: size.H/2 - 10},
		placeholder,
		tf.placeholderColor,
		20,
	)

	tf.body = NewContainer(
		basic.Point{},
		size,
		20,
		tf.backgroundColor,
		basic.Start,
		basic.Center,
		tf.textWidget,
	)

	return tf
}

func (t *TextField) GetPos() basic.Point {
	return t.pos
}

func (t *TextField) SetPos(p basic.Point) {
	t.pos = p
}

func (t *TextField) GetSize() basic.Size {
	return t.size
}

func (t *TextField) SetSize(sz basic.Size) {
	t.size = sz
	if t.body != nil {
		t.body.SetSize(sz)
	}
}

func (t *TextField) Update(offset basic.Point) {
	t.currentPos = t.pos.Add(offset)

	if t.body != nil {
		t.body.Update(t.currentPos)
	}

	mx, my := ebiten.CursorPosition()
	hovered := inputhelper.IsHovered(mx, my, t.currentPos, t.size)

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		t.focused = hovered
	}

	inputhelper.ReceiveText(&t.Text, t.focused)

	// 🔹 Controle do cursor piscando
	if t.focused {
		t.cursorCounter++
		if t.cursorCounter >= t.cursorBlinkRate {
			t.cursorVisible = !t.cursorVisible
			t.cursorCounter = 0
		}
	} else {
		t.cursorVisible = false
		t.cursorCounter = 0
	}

	// 🔹 Atualização do texto exibido
	if t.Text == "" {
		if t.focused {
			displayText := ""
			if t.cursorVisible {
				displayText = "|"
			}
			t.textWidget.Text = displayText
			t.textWidget.Color = t.textColor
		} else {
			t.textWidget.Text = t.placeholder
			t.textWidget.Color = t.placeholderColor
		}
	} else {
		displayText := t.Text
		if t.focused && t.cursorVisible {
			displayText += "|"
		}
		t.textWidget.Text = displayText
		t.textWidget.Color = t.textColor
	}
}

func (t *TextField) Draw(screen *ebiten.Image) {
	if t.body != nil {
		t.body.Draw(screen)
	}
}
