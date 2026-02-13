package components

import (
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/hajimehoshi/ebiten/v2"
)

// Row organiza widgets em uma linha, com espaçamento e altura opcional fixa
// <USAR A ALTURA DO PAI PARA ALINHAMENTO NO EIXO SECUNDARIO (Y)>
type Row struct {
	Pos        basic.Point //posição inicial
	Spacing    float32     //espaçamento horizontal entre elementos
	Children   []Widget
	MainAlign  basic.Align // alinhamento dos elementos no eixo principal (center, start, end)
	CrossAlign basic.Align //eixo cruzado
	size       basic.Size  //para calculo de tamanho caso necessario
	currentPos basic.Point
}

func (r *Row) IsLayout() bool {
	return true
}

// NewRow cria uma linha e já calcula a posição de todos os widgets,
// alinhando verticalmente e no eixo secundario de acordo com o alinhamento dado
// construtor
func NewRow(
	pos basic.Point,
	spacing float32,
	parentSize basic.Size,
	mainAlign basic.Align,
	crossAlign basic.Align,
	children []Widget,
) *Row {

	r := &Row{
		Pos:        pos,
		Spacing:    spacing,
		Children:   children,
		MainAlign:  mainAlign,
		CrossAlign: crossAlign,
	}

	// 1) posiciona como Start/Start
	r.init()

	//aplica alinhamentos relativos ao retângulo do pai iniciado em r.Pos
	if mainAlign != basic.Start || crossAlign != basic.Start {
		r.align(parentSize)
	}

	return r
}

// posiciona filhos como Start/Start
func (r *Row) init() {
	cursorX := r.Pos.X

	for _, w := range r.Children {
		size := w.GetSize()
		w.SetPos(basic.Point{
			X: cursorX,
			Y: r.Pos.Y,
		})
		cursorX += size.W + r.Spacing
	}
	//calcula tamanho da row
	r.calcSize()
}

// Update chama Update de todos os filhos
func (r *Row) Update(offset basic.Point) {
	r.currentPos = r.Pos.Add(offset)
	for _, w := range r.Children {
		w.Update(r.currentPos)
	}
}

// Draw chama Draw de todos os filhos fazendo recursão com posição na arvore
func (r *Row) Draw(screen *ebiten.Image) {
	for _, w := range r.Children {
		w.Draw(screen)
	}
}

func (r *Row) align(parentSize basic.Size) {
	var offsetX float32

	// Eixo principal (horizontal): calculamos o offsetX com base na largura total
	switch r.MainAlign {
	case basic.Start:
		offsetX = 0
	case basic.Center:
		offsetX = (parentSize.W - r.GetSize().W) / 2
	case basic.End:
		offsetX = parentSize.W - r.GetSize().W
	}

	//para cada filho,calcula o x individual
	for _, w := range r.Children {

		if w == nil {
			continue
		}

		// se for layout, ele se alinha sozinho (!!!!)
		if _, ok := w.(LayoutWidget); ok {
			continue
		}

		p := w.GetPos()

		childHeight := w.GetSize().H

		//alinha eixo cruzado individualmente
		switch r.CrossAlign {
		case basic.Start:
			p.Y = r.Pos.Y // começa no início
		case basic.Center:
			p.Y = r.Pos.Y + (parentSize.H-childHeight)/2
		case basic.End:
			p.Y = r.Pos.Y + (parentSize.H - childHeight)
		}

		//adiciona o offset do main axis
		p.X = r.Pos.X + offsetX + p.X

		w.SetPos(p)
	}
}

// tamanho ocupado pela row, calculado no initWidgets
func (r *Row) calcSize() {
	var totalW float32
	var maxH float32

	for i, w := range r.Children {
		s := w.GetSize()

		if i > 0 {
			totalW += r.Spacing
		}
		totalW += s.W

		if s.H > maxH {
			maxH = s.H
		}
	}

	r.size = basic.Size{W: totalW, H: maxH}
}

func (r *Row) GetPos() basic.Point {
	return r.Pos
}

func (r *Row) GetSize() basic.Size {
	return r.size
}

func (r *Row) SetPos(p basic.Point) {
	dx := p.X - r.Pos.X
	dy := p.Y - r.Pos.Y

	r.Pos = p

	for _, w := range r.Children {
		cp := w.GetPos()
		cp.X += dx
		cp.Y += dy
		w.SetPos(cp)
	}
}
