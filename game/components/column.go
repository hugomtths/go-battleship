package components

import (
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/hajimehoshi/ebiten/v2"
)

// Column organiza widgets em uma coluna, com espaçamento e largura opcional fixa
// <USAR A LARGURA DO PAI PARA ALINHAMENTO NO EIXO SECUNDARIO (X)>
type Column struct {
	Pos, currentPos basic.Point //posição inicial

	Spacing float32 //espaçamento vertical entre elementos

	cursorY float32 //posição vertical do próximo elemento

	MainAlign basic.Align // alinhamento dos elementos no eixo principal (center, start, end)

	CrossAlign basic.Align //eixo cruzado

	Children []Widget //elementos da column

	size basic.Size //caso necessario
}

func (c *Column) IsLayout() bool {
	return true
}

// NewColumn cria uma coluna e já calcula a posição de todos os widgets
// alinhando no eixo principal (vertical) e no eixo cruzado (horizontal)
func NewColumn(
	pos basic.Point,
	spacing float32,
	parentSize basic.Size,
	mainAlign basic.Align,
	crossAlign basic.Align,
	children []Widget,
) *Column {

	c := &Column{
		Pos:        pos,
		Spacing:    spacing,
		Children:   children,
		MainAlign:  mainAlign,
		CrossAlign: crossAlign,
	}

	// posicionamento inicial (Start/Start)
	c.init()

	if mainAlign != basic.Start || crossAlign != basic.Start {
		c.align(parentSize)
	}

	for _, w := range c.Children {
		p := w.GetPos()
		p.X += c.Pos.X
		p.Y += c.Pos.Y
		w.SetPos(p)
	}

	return c
}

// Update chama Update de todos os filhos
func (c *Column) Update(offset basic.Point) {
	c.currentPos = c.Pos.Add(offset)

	for _, w := range c.Children {
		w.Update(c.currentPos)
	}
}

func (c *Column) Draw(screen *ebiten.Image) {
	for _, w := range c.Children {
		w.Draw(screen)
	}
}
func (c *Column) align(parentSize basic.Size) {
	var offsetY float32

	// Eixo principal (vertical): calculamos o offsetY com base na altura total
	switch c.MainAlign {
	case basic.Start:
		offsetY = 0
	case basic.Center:
		offsetY = (parentSize.H - c.GetSize().H) / 2
	case basic.End:
		offsetY = parentSize.H - c.GetSize().H
	}

	//para cada filho,calcula o x individual
	for _, w := range c.Children {

		if w == nil {
			continue
		}

		// se for layout, ele se alinha sozinho (!!!!)
		if _, ok := w.(LayoutWidget); ok {
			continue
		}
		p := w.GetPos()

		childWidth := w.GetSize().W

		switch c.CrossAlign {
		case basic.Start:
			p.X = c.Pos.X // começa no início
		case basic.Center:
			p.X = c.Pos.X + (parentSize.W-childWidth)/2
		case basic.End:
			p.X = c.Pos.X + (parentSize.W - childWidth)
		}

		//adiciona o offset do mainaxis
		p.Y = c.Pos.Y + offsetY + p.Y

		w.SetPos(p)
	}
}

func (c *Column) GetPos() basic.Point {
	return c.Pos
}

func (c *Column) SetPos(p basic.Point) {
	dx := p.X - c.Pos.X
	dy := p.Y - c.Pos.Y

	c.Pos = p

	for _, w := range c.Children {
		cp := w.GetPos()
		cp.X += dx
		cp.Y += dy
		w.SetPos(cp)
	}
}

// calcula tamanho (apenas no initWidgets)
func (c *Column) calcSize() {
	var totalH float32
	var maxW float32

	for i, w := range c.Children {
		s := w.GetSize()

		totalH += s.H
		if i < len(c.Children)-1 {
			totalH += c.Spacing
		}

		if s.W > maxW {
			maxW = s.W
		}
	}

	//calcula size uma unica vez aqui
	c.size = basic.Size{W: maxW, H: totalH}
}

// GetSize retorna dimensões da column
func (c *Column) GetSize() basic.Size {
	return c.size
}

// initWidgets serve para um primeiro posicionamento dos elementos (start x start)
func (c *Column) init() {
	cursorY := c.Pos.Y

	for i, w := range c.Children {
		size := w.GetSize()

		w.SetPos(basic.Point{
			X: c.Pos.X,
			Y: cursorY,
		})

		cursorY += size.H
		if i < len(c.Children)-1 {
			cursorY += c.Spacing
		}
	}
	c.calcSize()
}
