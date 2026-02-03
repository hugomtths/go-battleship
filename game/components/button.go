package components

import (
	"image/color"
	"strings"

	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	"github.com/allanjose001/go-battleship/game/util"
	"github.com/hajimehoshi/ebiten/v2"
)

type Button struct {
	pos, scaledPos             basic.Point
	size, scaledSize           basic.Size
	label                      string
	backgroundColor, textColor color.Color
	CallBack                   func(*Button) //função que o botão chama
	hoverColor                 color.Color
	onHover                    func(*Button) //opcional
	disabled, hovered, clicked bool
	scale                      float32        // para animar click
	body                       StylableWidget //um container por ex
}

func NewButton(
	pos basic.Point, //opcional
	size basic.Size, //pode ser nil/zero
	label string,
	color color.Color,
	textColor color.Color,
	cb func(*Button), // ir para uma tela...

) *Button {
	bt := &Button{
		pos:             pos,
		scaledPos:       pos,
		size:            size,
		scaledSize:      size,
		scale:           1.0,
		label:           label,
		backgroundColor: color,
		textColor:       textColor,
		CallBack:        cb,
		hoverColor:      colors.Lighten(color, 0.25),
	}

	bt.makeBody() //cria body com container e variaveis de button

	return bt
}
func (b *Button) GetPos() basic.Point {
	return b.pos
}

func (b *Button) SetPos(point basic.Point) {
	b.pos = point
}

func (b *Button) GetSize() basic.Size {
	return b.size
}

func (b *Button) Update() {
	mouseX, mouseY := ebiten.CursorPosition() //ver como fazer com disabled

	//TODO: colocar som de hovered
	b.hovered = inputhelper.IsHovered(mouseX, mouseY, b.pos, b.size)

	b.clicked = inputhelper.IsClicked(mouseX, mouseY, b.pos, b.size)

	//b.animateClick()

	/*//dispara callback do botão
	if b.CallBack != nil &&inputhelper.IsClicked(mouseX, mouseY, b.pos, b.size) {
		//b.CallBack(b)
		println("clicked")
		b.animateClick()
	}*/

}

func (b *Button) Draw(screen *ebiten.Image) {
	b.draw(screen, basic.Point{})
}

func (b *Button) SetSize(sz basic.Size) {
	b.body.SetSize(sz)
	b.size = sz
}

func (b *Button) draw(screen *ebiten.Image, offset basic.Point) {
	pos := b.pos.Add(offset)

	if b.hovered {
		b.body.SetColor(b.hoverColor)
	} else {
		b.body.SetColor(b.backgroundColor)
	}

	b.body.draw(screen, pos)
}

func (b *Button) makeBody() {

	if b.textColor == nil {
		b.textColor = colors.White
	}
	//corpo do botão com container
	b.body = NewContainer(
		b.pos,
		b.size, //tamanho original fica guardado
		16.0,
		b.backgroundColor,
		basic.Center,
		basic.Center,
		NewText(
			basic.Point{},
			strings.ToUpper(b.label),
			b.textColor,
			18, //VER SE ESSA FONTE DA
		),
		func(c *Container) {
			//fazer aqui relação com callback
		},
	)
}

/*
func (b *Button) animateClick() {
	//ifs para alterar escala do botão
	if b.clicked {
		b.scale += 0.1
		if b.scale >= 1.2 {
			b.scale = 1.2
			b.clicked = false
		}
	} else {
		if b.scale > 1.0 {
			b.scale -= 0.05
			if b.scale < 1.0 {
				b.scale = 1.0
			}
		}
	}

	// aplica escala ao container mudando tamanho "scaled" escalando tamanho original com fator
	b.SetSize(b.GetSize().Scale(b.scale))
	b.body.SetPos(b.GetCenter(b.pos, b.size, b.scaledSize))

}

// GetCenter retorna centro do widget após escalar seu tamanho
func (b *Button) GetCenter(originPos basic.Point, originSize, newSize basic.Size) basic.Point {
	centerX := originPos.X + originSize.W/2
	centerY := originPos.Y + originSize.H/2
	return basic.Point{
		X: centerX - newSize.W/2,
		Y: centerY - newSize.H/2,
	}
}
*/
