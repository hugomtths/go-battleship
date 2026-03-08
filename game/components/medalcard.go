package components

import (
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/allanjose001/go-battleship/game/components/basic/colors"
	"github.com/hajimehoshi/ebiten/v2"
)

// MedalCard encapsula um card visual de medalha reutilizável.
// Ele funciona como um widget composto: toda a renderização real é delegada
// para o campo `body`, que contém a árvore de containers/layouts internos.
//
// (ex: cor diferente) via StylableWidget.
type MedalCard struct {
	pos, currentPos basic.Point    // posição base e posição final com offset aplicado
	size            basic.Size     // tamanho total do card (usado pelo layout pai)
	body            StylableWidget // widget raiz que contém tod o layout interno
}

// NewMedal constrói a estrutura completa do card.
// O layout interno é proporcional ao tamanho recebido:
//
//	ATUALMENTE:  10% da largura -> área do ícone, 50% da largura -> área de texto
//
// O restante funciona como padding visual.
func NewMedal(icon, title, desc string, size basic.Size) *MedalCard {

	var iconSize basic.Size

	if title == "BLOQUEADA" {
		iconSize = basic.Size{W: 40, H: 55}
	} else {
		iconSize = basic.Size{W: 40, H: 70}
	}

	iconImage, err := NewImage(icon, basic.Point{}, iconSize)
	titleTxt := NewText(basic.Point{}, title, colors.GoldMedal, 16)

	// Largura máxima da descrição = mesma do container de texto original (size.W * 0.6)
	// O texto vai quebrar linha e centralizar automaticamente dentro dessa área
	descTxt := NewTextWrap(basic.Point{}, desc, colors.GoldMedal, 12, size.W*0.6)

	var iconHandler Widget

	if err != nil {
		iconHandler = NewText(basic.Point{}, "ERROR", colors.Red, 16)
	} else {
		iconHandler = iconImage
	}

	return &MedalCard{
		size: size,
		body: NewContainer(
			basic.Point{},
			size,
			15,
			colors.White,
			basic.Center,
			basic.Center,
			NewRow(
				basic.Point{}, 10,
				size,
				basic.Center,
				basic.Center,
				[]Widget{
					// Container do ícone — igual ao original
					NewContainer(
						basic.Point{},
						basic.Size{
							W: size.W * 0.1,
							H: size.H,
						},
						0, colors.Transparent, basic.Center, basic.Center,
						iconHandler,
					),

					// Container de texto — igual ao original
					NewContainer(
						basic.Point{},
						basic.Size{
							W: size.W * 0.5,
							H: size.H,
						},
						0, colors.Transparent, basic.Start, basic.Start,
						NewColumn(
							basic.Point{}, 12,
							basic.Size{
								W: size.W * 0.6,
								H: size.H,
							},
							basic.Center,
							basic.Center,
							[]Widget{
								titleTxt,
								descTxt, // <- única mudança: Text -> TextWrap
							},
						),
					),
				},
			),
		),
	}
}

// GetPos retorna a posição base do card.
func (m *MedalCard) GetPos() basic.Point {
	return m.pos
}

// SetPos define a posição base do card.
func (m *MedalCard) SetPos(point basic.Point) {
	m.pos = point
}

// GetSize retorna o tamanho total do card.
func (m *MedalCard) GetSize() basic.Size {
	return m.size
}

// Update aplica offset acumulado e propaga para o layout interno.
func (m *MedalCard) Update(offset basic.Point) {
	m.currentPos = m.pos.Add(offset)
	m.body.Update(m.currentPos)
}

// Draw delega o desenho para o widget raiz.
func (m *MedalCard) Draw(screen *ebiten.Image) {
	m.body.Draw(screen)
}
